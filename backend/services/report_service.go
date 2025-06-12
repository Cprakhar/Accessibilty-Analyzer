package services

import (
	"backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reportCollection *mongo.Collection

func InitReportService(db *mongo.Database) {
	reportCollection = db.Collection("reports")
}

func CreateReport(ctx context.Context, userId primitive.ObjectID, url, html string) (*models.Report, error) {
	report := &models.Report{
		UserID:          userId,
		URL:             url,
		HTMLSnapshot:    html,
		AnalysisResults: nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.ReportStatusPending,
	}
	res, err := reportCollection.InsertOne(ctx, report)
	if err != nil {
		return nil, err
	}
	report.ID = res.InsertedID.(primitive.ObjectID)
	return report, nil
}

func UpdateReportResults(ctx context.Context, reportId primitive.ObjectID, results interface{}, status models.ReportStatus) error {
	update := bson.M{
		"$set": bson.M{
			"analysisResults": results,
			"status":          status,
			"updatedAt":       time.Now(),
		},
	}
	_, err := reportCollection.UpdateByID(ctx, reportId, update)
	return err
}

func GetReportByID(ctx context.Context, reportId primitive.ObjectID) (*models.Report, error) {
	var report models.Report
	err := reportCollection.FindOne(ctx, bson.M{"_id": reportId}).Decode(&report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}
