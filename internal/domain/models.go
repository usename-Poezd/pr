package domain

type Poll struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	ResultsVisible bool       `json:"results_visible"`
	Questions      []Question `json:"questions"`
	AdminToken     string     `json:"-"`
}
type Question struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	Multiple bool     `json:"multiple"`
	Options  []Option `json:"options"`
}
type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type VoteAnswer struct {
	QuestionID string   `json:"question_id"`
	OptionID   string   `json:"option_id,omitempty"`
	OptionIDs  []string `json:"option_ids,omitempty"`
}
type Results struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Questions []QuestionResults `json:"questions"`
}
type QuestionResults struct {
	ID       string          `json:"id"`
	Text     string          `json:"text"`
	Multiple bool            `json:"multiple"`
	Options  []OptionResults `json:"options"`
}
type OptionResults struct {
	ID         string  `json:"id"`
	Text       string  `json:"text"`
	Votes      int64   `json:"votes"`
	Percentage float64 `json:"percentage"`
}
