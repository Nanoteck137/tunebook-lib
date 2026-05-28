package timer

import "time"

type Simple struct {
	start time.Time

	duration time.Duration
}

func (s *Simple) Start() {
	s.start = time.Now()
}

func (s *Simple) Stop() time.Duration {
	t := time.Since(s.start)
	s.duration = t

	return t
}

func (s *Simple) Duration() time.Duration {
	return s.duration
}
