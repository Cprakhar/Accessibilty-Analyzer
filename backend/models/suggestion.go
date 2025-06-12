package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SuggestionSummary struct {
	Problem       string `bson:"problem" json:"problem"`
	Impact        string `bson:"impact" json:"impact"`
	AffectedUsers string `bson:"affectedUsers" json:"affectedUsers"`
}

type SuggestionWhyMatters struct {
	UserImpact            string `bson:"userImpact" json:"userImpact"`
	AssistiveTechAffected string `bson:"assistiveTechAffected" json:"assistiveTechAffected"`
}

type SuggestionHowToFix struct {
	Step1       string `bson:"step1" json:"step1"`
	CodeExample string `bson:"codeExample" json:"codeExample"`
}

type SuggestionTestingInstructions struct {
	Verify string `bson:"verify" json:"verify"`
	Tools  string `bson:"tools" json:"tools"`
}

type SuggestionPriorityLevel struct {
	Urgency   string `bson:"urgency" json:"urgency"`
	WCAGLevel string `bson:"wcagLevel" json:"wcagLevel"`
}

type SuggestionItem struct {
	Issue               string                        `bson:"issue" json:"issue"`
	Summary             SuggestionSummary             `bson:"summary" json:"summary"`
	WhyMatters          SuggestionWhyMatters          `bson:"whyMatters" json:"whyMatters"`
	HowToFix            SuggestionHowToFix            `bson:"howToFix" json:"howToFix"`
	TestingInstructions SuggestionTestingInstructions `bson:"testingInstructions" json:"testingInstructions"`
	PriorityLevel       SuggestionPriorityLevel       `bson:"priorityLevel" json:"priorityLevel"`
}

type Suggestion struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ReportID    primitive.ObjectID `bson:"reportId" json:"reportId"`
	Suggestions []SuggestionItem   `bson:"suggestions" json:"suggestions"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}
