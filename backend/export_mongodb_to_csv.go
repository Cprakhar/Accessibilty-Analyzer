package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Violation struct {
	ID          string   `bson:"id"`
	Impact      string   `bson:"impact"`
	Tags        []string `bson:"tags"`
	Help        string   `bson:"help"`
	HelpURL     string   `bson:"helpUrl"`
	Description string   `bson:"description"`
	Nodes       []struct {
		HTML string `bson:"html"`
	} `bson:"nodes"`
}

type Report struct {
	ID              string    `bson:"_id"`
	URL             string    `bson:"url"`
	Domain          string    `bson:"domain"`
	CreatedAt       time.Time `bson:"createdAt"`
	AnalysisResults struct {
		Violations []Violation `bson:"violations"`
	} `bson:"analysisResults"`
}

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/accessibility_analyser"
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := client.Database("accessibility_analyser")
	reports := db.Collection("reports")

	cur, err := reports.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("Failed to query reports: %v", err)
	}
	defer cur.Close(context.Background())

	file, err := os.Create("violations_export.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"scan_id", "scan_date", "url", "domain", "violation_id", "violation_type", "impact", "wcag_tags", "element_count", "help_url", "description", "node_html"}
	w.Write(headers)

	for cur.Next(context.Background()) {
		var r Report
		if err := cur.Decode(&r); err != nil {
			continue
		}
		for _, v := range r.AnalysisResults.Violations {
			tags := ""
			if len(v.Tags) > 0 {
				tags = fmt.Sprintf("%v", v.Tags)
			}
			elementCount := fmt.Sprintf("%d", len(v.Nodes))
			for _, node := range v.Nodes {
				row := []string{
					r.ID,
					r.CreatedAt.Format("2006-01-02 15:04:05"),
					r.URL,
					r.Domain,
					v.ID,
					v.Help,
					v.Impact,
					tags,
					elementCount,
					v.HelpURL,
					v.Description,
				}
				if node.HTML != "" {
					row = append(row, node.HTML)
				} else {
					row = append(row, "")
				}
				w.Write(row)
			}
		}
	}
	log.Println("Export complete: violations_export.csv")
}
