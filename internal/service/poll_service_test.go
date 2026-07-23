package service

import (
	"context"
	"errors"
	"poll-app/internal/domain"
	"poll-app/internal/repository"
	"testing"
)

type testRepo struct{}

func (testRepo) Create(context.Context, string, string, bool, []domain.Question) (domain.Poll, error) {
	return domain.Poll{ID: "id"}, nil
}
func (testRepo) Get(context.Context, string) (domain.Poll, error)        { return domain.Poll{}, nil }
func (testRepo) AdminToken(context.Context, string) (string, error)      { return "", nil }
func (testRepo) ResultsAccess(context.Context, string) (bool, error)     { return true, nil }
func (testRepo) Vote(context.Context, string, []domain.VoteAnswer) error { return nil }
func (testRepo) Results(context.Context, string) (domain.Results, error) {
	return domain.Results{}, nil
}

type voteRepo struct {
	answers []domain.VoteAnswer
	err     error
	single  bool
}

func (r *voteRepo) Create(context.Context, string, string, bool, []domain.Question) (domain.Poll, error) {
	return domain.Poll{ID: "id"}, nil
}
func (r *voteRepo) Get(context.Context, string) (domain.Poll, error)   { return domain.Poll{}, nil }
func (r *voteRepo) AdminToken(context.Context, string) (string, error) { return "", nil }
func (r *voteRepo) ResultsAccess(context.Context, string) (bool, error) {
	return true, nil
}
func (r *voteRepo) Vote(_ context.Context, _ string, a []domain.VoteAnswer) error {
	r.answers = a
	if r.single && len(a[0].OptionIDs) != 1 {
		return repository.ErrInvalidVote
	}
	return r.err
}
func (r *voteRepo) Results(context.Context, string) (domain.Results, error) {
	return domain.Results{}, nil
}

func TestCreateValidatesQuestions(t *testing.T) {
	s := NewPollService(testRepo{})
	if _, err := s.Create(context.Background(), "title", true, []domain.Question{{Text: "question", Options: []domain.Option{{Text: "only one"}}}}); err != ErrInvalid {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

type resultsRepo struct {
	testRepo
	adminToken string
	visible    bool
}

func (r resultsRepo) AdminToken(context.Context, string) (string, error) { return r.adminToken, nil }
func (r resultsRepo) ResultsAccess(context.Context, string) (bool, error) {
	return r.visible, nil
}

func TestResultsHideAnonymousPollFromPublicAccess(t *testing.T) {
	r := resultsRepo{visible: false}
	if _, err := NewPollService(r).Results(context.Background(), "poll", ""); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestResultsAllowAnonymousVoterVisibility(t *testing.T) {
	r := resultsRepo{visible: true}
	if _, err := NewPollService(r).Results(context.Background(), "poll", ""); err != nil {
		t.Fatal(err)
	}
}

func TestResultsAdminAlwaysAllowed(t *testing.T) {
	r := resultsRepo{visible: false, adminToken: "secret"}
	if _, err := NewPollService(r).Results(context.Background(), "poll", "secret"); err != nil {
		t.Fatal(err)
	}
}

func TestVoteAcceptsMultipleSelections(t *testing.T) {
	r := &voteRepo{}
	if err := NewPollService(r).Vote(context.Background(), "poll", []domain.VoteAnswer{{QuestionID: "q", OptionIDs: []string{"a", "b"}}}); err != nil {
		t.Fatal(err)
	}
	if len(r.answers[0].OptionIDs) != 2 {
		t.Fatalf("expected two options, got %#v", r.answers)
	}
}

func TestVoteRejectsDuplicateOption(t *testing.T) {
	if err := NewPollService(&voteRepo{}).Vote(context.Background(), "poll", []domain.VoteAnswer{{QuestionID: "q", OptionIDs: []string{"a", "a"}}}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestVoteRejectsMultiselectRepositoryError(t *testing.T) {
	r := &voteRepo{single: true}
	if err := NewPollService(r).Vote(context.Background(), "poll", []domain.VoteAnswer{{QuestionID: "q", OptionIDs: []string{"a", "b"}}}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}
