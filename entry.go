package frontend

import "time"

type Entry struct {
	ID        int
	Date      time.Time
	ProjectID int
	User      string
	Comment   string
	Duration  int64
}
