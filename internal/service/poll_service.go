package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"poll-app/internal/domain"
	"poll-app/internal/repository"
	"strings"
)

var ErrInvalid = errors.New("invalid poll input")
var ErrForbidden = errors.New("invalid admin token")

type PollService struct{ repo repository.PollRepository }

func NewPollService(r repository.PollRepository) *PollService { return &PollService{repo: r} }
func (s *PollService) Create(ctx context.Context, title string, resultsVisible bool, qs []domain.Question) (domain.Poll, error) {
	title = strings.TrimSpace(title)
	if title == "" || len(title) > 500 || len(qs) == 0 || len(qs) > 50 {
		return domain.Poll{}, ErrInvalid
	}
	for i := range qs {
		qs[i].Text = strings.TrimSpace(qs[i].Text)
		if qs[i].Text == "" || len(qs[i].Text) > 1000 || len(qs[i].Options) < 2 || len(qs[i].Options) > 50 {
			return domain.Poll{}, ErrInvalid
		}
		seen := map[string]bool{}
		for j := range qs[i].Options {
			qs[i].Options[j].Text = strings.TrimSpace(qs[i].Options[j].Text)
			key := strings.ToLower(qs[i].Options[j].Text)
			if key == "" || len(key) > 500 || seen[key] {
				return domain.Poll{}, ErrInvalid
			}
			seen[key] = true
		}
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return domain.Poll{}, err
	}
	token := fmt.Sprintf("%x", b)
	p, err := s.repo.Create(ctx, title, token, resultsVisible, qs)
	p.AdminToken = token
	return p, err
}
func (s *PollService) Get(ctx context.Context, id string) (domain.Poll, error) {
	if id == "" {
		return domain.Poll{}, ErrInvalid
	}
	return s.repo.Get(ctx, id)
}
func (s *PollService) Vote(ctx context.Context, id string, a []domain.VoteAnswer) error {
	if id == "" || len(a) == 0 {
		return ErrInvalid
	}
	seen := map[string]bool{}
	options := map[string]bool{}
	for i := range a {
		v := &a[i]
		if v.QuestionID == "" || seen[v.QuestionID] {
			return ErrInvalid
		}
		if len(v.OptionIDs) > 0 && v.OptionID != "" {
			return ErrInvalid
		}
		if len(v.OptionIDs) == 0 && v.OptionID != "" {
			v.OptionIDs = []string{v.OptionID}
			v.OptionID = ""
		}
		if len(v.OptionIDs) == 0 {
			return ErrInvalid
		}
		for _, optionID := range v.OptionIDs {
			if optionID == "" || options[optionID] {
				return ErrInvalid
			}
			options[optionID] = true
		}
		seen[v.QuestionID] = true
	}
	err := s.repo.Vote(ctx, id, a)
	if errors.Is(err, repository.ErrInvalidVote) {
		return ErrInvalid
	}
	return err
}
func (s *PollService) Results(ctx context.Context, id, token string) (domain.Results, error) {
	if id == "" {
		return domain.Results{}, ErrInvalid
	}
	if token != "" {
		expected, err := s.repo.AdminToken(ctx, id)
		if err != nil {
			return domain.Results{}, err
		}
		if token != expected {
			return domain.Results{}, ErrForbidden
		}
	} else {
		visible, err := s.repo.ResultsAccess(ctx, id)
		if err != nil {
			return domain.Results{}, err
		}
		if !visible {
			return domain.Results{}, ErrForbidden
		}
	}
	return s.repo.Results(ctx, id)
}
