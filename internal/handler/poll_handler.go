package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"poll-app/internal/domain"
	"poll-app/internal/repository"
	"poll-app/internal/service"
)

type Handler struct{ s *service.PollService }

func New(s *service.PollService) *Handler { return &Handler{s: s} }

type createQuestion struct {
	Text     string   `json:"text"`
	Multiple bool     `json:"multiple"`
	Options  []string `json:"options"`
}
type createRequest struct {
	Title     string           `json:"title"`
	Questions []createQuestion `json:"questions"`
}
type voteRequest struct {
	Answers []domain.VoteAnswer `json:"answers"`
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/polls", h.create)
	r.Get("/polls/{id}", h.get)
	r.Post("/polls/{id}/vote", h.vote)
	r.Get("/polls/{id}/results", h.results)
	return r
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, e error) {
	status := 500
	msg := "internal server error"
	switch {
	case errors.Is(e, service.ErrInvalid):
		status = 400
		msg = "invalid request"
	case errors.Is(e, service.ErrForbidden):
		status = 403
		msg = "forbidden"
	case errors.Is(e, repository.ErrNotFound):
		status = 404
		msg = "poll not found"
	}
	write(w, status, map[string]string{"error": msg})
}
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in createRequest
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		fail(w, service.ErrInvalid)
		return
	}
	qs := make([]domain.Question, len(in.Questions))
	for i, q := range in.Questions {
		qs[i].Text = q.Text
		qs[i].Multiple = q.Multiple
		for _, text := range q.Options {
			qs[i].Options = append(qs[i].Options, domain.Option{Text: text})
		}
	}
	p, e := h.s.Create(r.Context(), in.Title, qs)
	if e != nil {
		fail(w, e)
		return
	}
	id := p.ID
	write(w, 201, map[string]string{"id": id, "admin_token": p.AdminToken, "public_link": "/polls/" + id, "admin_link": "/polls/" + id + "/results?admin_token=" + p.AdminToken})
}
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	p, e := h.s.Get(r.Context(), chi.URLParam(r, "id"))
	if e != nil {
		fail(w, e)
		return
	}
	write(w, 200, p)
}
func (h *Handler) vote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, e := r.Cookie("voted_" + id); e == nil {
		write(w, 409, map[string]string{"error": "already voted"})
		return
	}
	var in voteRequest
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		fail(w, service.ErrInvalid)
		return
	}
	if e := h.s.Vote(r.Context(), id, in.Answers); e != nil {
		fail(w, e)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "voted_" + id, Value: "true", Path: "/", MaxAge: 31536000, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	write(w, 201, map[string]string{"message": "vote recorded"})
}
func (h *Handler) results(w http.ResponseWriter, r *http.Request) {
	v, e := h.s.Results(r.Context(), chi.URLParam(r, "id"), r.URL.Query().Get("admin_token"))
	if e != nil {
		fail(w, e)
		return
	}
	write(w, 200, v)
}
