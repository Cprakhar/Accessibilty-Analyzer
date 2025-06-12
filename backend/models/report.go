package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"
	ReportStatusComplete ReportStatus = "complete"
	ReportStatusFailed   ReportStatus = "failed"
)

type Report struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID          primitive.ObjectID `bson:"userId" json:"userId"`
	URL             string             `bson:"url" json:"url"`
	HTMLSnapshot    string             `bson:"htmlSnapshot" json:"htmlSnapshot"`
	AnalysisResults interface{}        `bson:"analysisResults" json:"analysisResults"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	Status          ReportStatus       `bson:"status" json:"status"`
}
