package strategy

import "github.com/google/uuid"

type Writer interface {
	Insert(s *Strategy) error
}

type Reader interface {
	Get(q *ReaderQuery) ([]*Strategy, error)
}

type ReaderQuery struct {
	UserID     *uuid.UUID
	Market     *string
	Runner     *string
	Price      *float32
	CompetitionID *uint64
	Side       *string
	Status     *string
	Visibility *string
	OrderBy    *string
}

