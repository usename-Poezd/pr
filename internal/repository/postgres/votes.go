package postgres

import (
	"context"
	"poll-app/internal/domain"
	"poll-app/internal/repository"
)

func (r *Repository) Vote(ctx context.Context, pollID string, answers []domain.VoteAnswer) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var n int
	if err = tx.QueryRow(ctx, `SELECT count(*) FROM questions WHERE poll_id=$1`, pollID).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	if n != len(answers) {
		return repository.ErrInvalidVote
	}
	for _, a := range answers {
		var ok bool
		if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM options o JOIN questions q ON q.id=o.question_id WHERE q.poll_id=$1 AND q.id=$2 AND o.id=$3)`, pollID, a.QuestionID, a.OptionID).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return repository.ErrInvalidVote
		}
		if _, err = tx.Exec(ctx, `INSERT INTO votes(poll_id,question_id,option_id) VALUES($1,$2,$3)`, pollID, a.QuestionID, a.OptionID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
func (r *Repository) Results(ctx context.Context, id string) (domain.Results, error) {
	rows, err := r.db.Query(ctx, `SELECT p.id::text,p.title,q.id::text,q.text,o.id::text,o.text,count(v.id),coalesce(sum(count(v.id)) OVER (PARTITION BY q.id),0) FROM polls p JOIN questions q ON q.poll_id=p.id JOIN options o ON o.question_id=q.id LEFT JOIN votes v ON v.option_id=o.id AND v.question_id=q.id WHERE p.id=$1 GROUP BY p.id,p.title,q.id,q.text,q.position,o.id,o.text,o.position ORDER BY q.position,o.position`, id)
	if err != nil {
		return domain.Results{}, err
	}
	defer rows.Close()
	var out domain.Results
	var last string
	for rows.Next() {
		var q domain.QuestionResults
		var o domain.OptionResults
		var total int64
		if err = rows.Scan(&out.ID, &out.Title, &q.ID, &q.Text, &o.ID, &o.Text, &o.Votes, &total); err != nil {
			return out, err
		}
		if out.ID == "" {
			return out, repository.ErrNotFound
		}
		if q.ID != last {
			out.Questions = append(out.Questions, q)
			last = q.ID
		} else {
			out.Questions[len(out.Questions)-1].Options = append(out.Questions[len(out.Questions)-1].Options, o)
			continue
		}
		if total > 0 {
			o.Percentage = float64(o.Votes) * 100 / float64(total)
		}
		out.Questions[len(out.Questions)-1].Options = []domain.OptionResults{o}
	}
	if err = rows.Err(); err != nil {
		return out, err
	}
	if out.ID == "" {
		return out, repository.ErrNotFound
	}
	return out, nil
}
