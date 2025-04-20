package model

import (
	"time"

	"github.com/google/uuid"
)

type Meta struct {
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Uuid      uuid.UUID `json:"uuid" db:"uuid"`
	Type      MetaType  `json:"type" db:"type"`
	Name      string    `json:"name" db:"name"`
	UserID    int       `json:"-" db:"user_id"`
	ID        int       `json:"id" db:"id"`
	File      File
}

type File struct {
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Name      string    `json:"name" db:"name"`
	MetaID    int       `json:"-" db:"meta_id"`
	ID        int       `json:"id" db:"id"`
	Data      []byte    `json:"data" db:"data"`
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
