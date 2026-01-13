package model

// Event структура записи аудита
type Event struct {
	TimeStamp   int64  `json:"ts"`
	Action      string `json:"action"`
	UserID      string `json:"user_id"`
	OriginalURL string `json:"url"`
}

// EventArray список Event
type EventArray []Event
