package model

import (
	"time"

	"github.com/google/uuid"
)

type Meta struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
	File      File
	Type      MetaType  `db:"type"`
	UserID    int       `db:"user_id"`
	ID        int       `db:"id"`
	Uuid      uuid.UUID `db:"uuid"`
}

type File struct {
	CreatedAt time.Time `db:"created_at"`
	Name      string    `db:"name"`
	Data      []byte    `db:"data"`
	MetaID    int       `db:"meta_id"`
	ID        int       `db:"id"`
}

type MetaType int

const (
	MetaTypeFile MetaType = iota + 1
	MetaTypeCreds
	MetaTypeCard
)

func (m MetaType) String() string {
	switch m {
	case MetaTypeCreds:
		return "creds"
	case MetaTypeCard:
		return "card"
	default:
		return "file"
	}
}

func MetaTypeFromString(t string) MetaType {
	switch t {
	case "creds":
		return MetaTypeCreds
	case "card":
		return MetaTypeCard
	default:
		return MetaTypeFile
	}
}
