package domain

type Poll struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Questions  []Question `json:"questions"`
	AdminToken string     `json:"-"`
}
type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}
type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type VoteAnswer struct {
	QuestionID string `json:"question_id"`
	OptionID   string `json:"option_id"`
}
type Results struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Questions []QuestionResults `json:"questions"`
}
type QuestionResults struct {
	ID      string          `json:"id"`
	Text    string          `json:"text"`
	Options []OptionResults `json:"options"`
}
type OptionResults struct {
	ID         string  `json:"id"`
	Text       string  `json:"text"`
	Votes      int64   `json:"votes"`
	Percentage float64 `json:"percentage"`
}
