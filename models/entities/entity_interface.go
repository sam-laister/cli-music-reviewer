package entities

import "time"

type EntityInterface interface {
	SetID(uint64)
	GetID() uint64
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	MarkCreated()
	MarkUpdated()
	TableName() string
}
