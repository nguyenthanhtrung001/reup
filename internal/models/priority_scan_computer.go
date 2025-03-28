package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PriorityScanComputer struct {
	ID           primitive.ObjectID `bson:"_id"`
	ComputerName string             `bson:"computer_name"`
	CycleSeconds int                `bson:"cycle_seconds"`
	ScanNumbers  int                `bson:"scan_numbers"`
	DomainAPI    string             `bson:"domain_api,omitempty"` // domain_api có thể null, nên dùng con trỏ
	CreatedAt    *time.Time         `bson:"created_at,omitempty"`
	UpdatedAt    *time.Time         `bson:"updated_at,omitempty"`
}
