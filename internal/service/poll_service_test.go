package service

import (
	"context"
	"poll-app/internal/domain"
	"testing"
)

type testRepo struct{}

func (testRepo) Create(context.Context, string, string, []domain.Question) (domain.Poll, error) {
	return domain.Poll{ID: "id"}, nil
}
func (testRepo) Get(context.Context, string) (domain.Poll, error)        { return domain.Poll{}, nil }
func (testRepo) AdminToken(context.Context, string) (string, error)      { return "", nil }
func (testRepo) Vote(context.Context, string, []domain.VoteAnswer) error { return nil }
func (testRepo) Results(context.Context, string) (domain.Results, error) {
	return domain.Results{}, nil
}

func TestCreateValidatesQuestions(t *testing.T) {
	s := NewPollService(testRepo{})
	if _, err := s.Create(context.Background(), "title", []domain.Question{{Text: "question", Options: []domain.Option{{Text: "only one"}}}}); err != ErrInvalid {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}
