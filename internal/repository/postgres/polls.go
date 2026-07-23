package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"poll-app/internal/domain"
	"poll-app/internal/repository"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, title, token string, resultsVisible bool, qs []domain.Question) (domain.Poll, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Poll{}, err
	}
	defer tx.Rollback(ctx)
	var p domain.Poll
	if err = tx.QueryRow(ctx, `INSERT INTO polls(title,admin_token,results_visible) VALUES($1,$2,$3) RETURNING id::text`, title, token, resultsVisible).Scan(&p.ID); err != nil {
		return p, err
	}
	p.Title = title
	p.ResultsVisible = resultsVisible
	for i := range qs {
		var q domain.Question
		if err = tx.QueryRow(ctx, `INSERT INTO questions(poll_id,text,multiple,position) VALUES($1,$2,$3,$4) RETURNING id::text`, p.ID, qs[i].Text, qs[i].Multiple, i).Scan(&q.ID); err != nil {
			return p, err
		}
		q.Text = qs[i].Text
		q.Multiple = qs[i].Multiple
		for j := range qs[i].Options {
			var o domain.Option
			if err = tx.QueryRow(ctx, `INSERT INTO options(question_id,text,position) VALUES($1,$2,$3) RETURNING id::text`, q.ID, qs[i].Options[j].Text, j).Scan(&o.ID); err != nil {
				return p, err
			}
			o.Text = qs[i].Options[j].Text
			q.Options = append(q.Options, o)
		}
		p.Questions = append(p.Questions, q)
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Poll{}, err
	}
	return p, nil
}

func (r *Repository) Get(ctx context.Context, id string) (domain.Poll, error) {
	rows, err := r.db.Query(ctx, `SELECT p.id::text,p.title,p.results_visible,q.id::text,q.text,q.multiple,o.id::text,o.text FROM polls p JOIN questions q ON q.poll_id=p.id JOIN options o ON o.question_id=q.id WHERE p.id=$1 ORDER BY q.position,o.position`, id)
	if err != nil {
		return domain.Poll{}, err
	}
	defer rows.Close()
	var p domain.Poll
	var lastQ string
	for rows.Next() {
		var q domain.Question
		var o domain.Option
		if err = rows.Scan(&p.ID, &p.Title, &p.ResultsVisible, &q.ID, &q.Text, &q.Multiple, &o.ID, &o.Text); err != nil {
			return p, err
		}
		if p.ID == "" {
			return p, repository.ErrNotFound
		}
		if q.ID != lastQ {
			p.Questions = append(p.Questions, q)
			lastQ = q.ID
		} else {
			p.Questions[len(p.Questions)-1].Options = append(p.Questions[len(p.Questions)-1].Options, o)
			continue
		}
		p.Questions[len(p.Questions)-1].Options = []domain.Option{o}
	}
	if err = rows.Err(); err != nil {
		return p, err
	}
	if p.ID == "" {
		return p, repository.ErrNotFound
	}
	return p, nil
}

func (r *Repository) ResultsAccess(ctx context.Context, id string) (bool, error) {
	var visible bool
	err := r.db.QueryRow(ctx, `SELECT results_visible FROM polls WHERE id=$1`, id).Scan(&visible)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, repository.ErrNotFound
	}
	return visible, err
}
func (r *Repository) AdminToken(ctx context.Context, id string) (string, error) {
	var token string
	err := r.db.QueryRow(ctx, `SELECT admin_token FROM polls WHERE id=$1`, id).Scan(&token)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", repository.ErrNotFound
	}
	return token, err
}
