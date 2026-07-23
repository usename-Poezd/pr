package repository

import (
	"context"
	"errors"
	"poll-app/internal/domain"
)

var ErrNotFound = errors.New("poll not found")
var ErrInvalidAdminToken = errors.New("invalid admin token")
var ErrInvalidVote = errors.New("invalid vote")

type PollRepository interface {
	Create(context.Context, string, string, bool, []domain.Question) (domain.Poll, error)
	Get(context.Context, string) (domain.Poll, error)
	AdminToken(context.Context, string) (string, error)
	ResultsAccess(context.Context, string) (bool, error)
	Vote(context.Context, string, []domain.VoteAnswer) error
	Results(context.Context, string) (domain.Results, error)
}
