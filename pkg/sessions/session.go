package sessions

import "time"

const ExpireTime = 120 * time.Second

type Session struct {
	ID     string
	UserId int
	Email  string
	Expiry time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
