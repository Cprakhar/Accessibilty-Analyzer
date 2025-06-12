package services

import (
	"backend/models"
	"context"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reportCollection *mongo.Collection
var suggestionCollection *mongo.Collection

func InitReportService(db *mongo.Database) {
	reportCollection = db.Collection("reports")
}

func InitSuggestionService(db *mongo.Database) {
	suggestionCollection = db.Collection("suggestions")
}

func CreateReport(ctx context.Context, userId primitive.ObjectID, urlStr, html string) (*models.Report, error) {
	domain := ""
	if parsed, err := url.Parse(urlStr); err == nil {
		domain = parsed.Hostname()
	}
	report := &models.Report{
		UserID:          userId,
		URL:             urlStr,
		Domain:          domain,
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

func ListReportsByUser(ctx context.Context, userId primitive.ObjectID) ([]map[string]interface{}, error) {
	cur, err := reportCollection.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var reports []map[string]interface{}
	for cur.Next(ctx) {
		var r models.Report
		if err := cur.Decode(&r); err != nil {
			continue
		}
		reports = append(reports, map[string]interface{}{
			"_id":       r.ID,
			"url":       r.URL,
			"createdAt": r.CreatedAt,
			"status":    r.Status,
		})
	}
	return reports, nil
}

func DeleteReportByID(ctx context.Context, reportId primitive.ObjectID) error {
	_, err := reportCollection.DeleteOne(ctx, bson.M{"_id": reportId})
	return err
}

func CreateSuggestion(ctx context.Context, reportId primitive.ObjectID, suggestions []models.SuggestionItem) error {
	s := &models.Suggestion{
		ReportID:    reportId,
		Suggestions: suggestions,
		CreatedAt:   time.Now(),
	}
	_, err := suggestionCollection.InsertOne(ctx, s)
	return err
}

func GetSuggestionsByReportID(ctx context.Context, reportId primitive.ObjectID) (map[string]interface{}, error) {
	var s models.Suggestion
	err := suggestionCollection.FindOne(ctx, bson.M{"reportId": reportId}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"reportId":    s.ReportID,
		"suggestions": s.Suggestions,
		"createdAt":   s.CreatedAt,
	}, nil
}
