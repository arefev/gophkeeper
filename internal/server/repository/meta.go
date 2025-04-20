package repository

import (
	"go.uber.org/zap"
)

type Meta struct {
	log *zap.Logger
	*Base
}

func NewMeta(tr TxGetter, log *zap.Logger) *Meta {
	return &Meta{
		log:  log,
		Base: NewBase(tr, log),
	}
}
