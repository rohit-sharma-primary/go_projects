package model

type Event struct {
	EventId   string    `json:"event_id"`
	UserId    int       `json:"user_id"`
	Type      EventType `json:"type"`
	Timestamp int64     `json:"timestamp"`
}

type EventType string

const (
	Click    EventType = "click"
	View     EventType = "view"
	Purchase EventType = "purchase"
)

type Metrics struct {
	TotalEvents    int               `json:"total_events"`
	EventTypeCount map[EventType]int `json:"event_type_count"`
	PerUserCount   map[int]int       `json:"per_user_count"`
}
