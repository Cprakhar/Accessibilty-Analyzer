package jobs

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"

	"backend/models"
	"backend/services"
	"backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalyzeJob struct {
	ReportID primitive.ObjectID
	URL      string
	HTML     string
}

var analyzeJobQueue = make(chan AnalyzeJob, 100)

func EnqueueAnalyzeJob(job AnalyzeJob) {
	analyzeJobQueue <- job
}

func StartAnalyzeWorker() {
	go func() {
		for job := range analyzeJobQueue {
			processAnalyzeJob(job)
		}
	}()
}

func processAnalyzeJob(job AnalyzeJob) {
	userID := "unknown"
	// Try to fetch userId from report for logging
	report, err := services.GetReportByID(context.Background(), job.ReportID)
	if err == nil {
		userID = report.UserID.Hex()
	}
	// Debug: log docker path (optional, keep for troubleshooting only)
	// path, pathErr := exec.LookPath("docker")
	// log.Printf("[DEBUG] Docker path: %s, err: %v", path, pathErr)
	// if pathErr != nil {
	// 	utils.LogAction(userID, "analyze", "failure", "docker not found: "+pathErr.Error())
	// 	return
	// }

	input := make(map[string]string)
	if job.URL != "" {
		input["url"] = job.URL
	}
	if job.HTML != "" {
		input["html"] = job.HTML
	}
	jsonInput, err := json.Marshal(input)
	if err != nil {
		utils.LogAction(userID, "analyze", "failure", "Failed to marshal input: "+err.Error())
		_ = services.UpdateReportResults(context.Background(), job.ReportID, map[string]interface{}{"error": "Failed to marshal input"}, models.ReportStatusFailed)
		return
	}
	cmd := exec.Command("docker", "run", "-i", "--rm", "axe-runner")
	cmd.Stdin = strings.NewReader(string(jsonInput))
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			utils.LogAction(userID, "analyze", "failure", "axe-runner stderr: "+string(exitErr.Stderr))
		}
		utils.LogAction(userID, "analyze", "failure", "axe-runner failed: "+err.Error())
		_ = services.UpdateReportResults(context.Background(), job.ReportID, map[string]interface{}{"error": err.Error()}, models.ReportStatusFailed)
		return
	}
	var results map[string]interface{}
	err = json.Unmarshal(output, &results)
	if err != nil {
		utils.LogAction(userID, "analyze", "failure", "Invalid axe-runner output: "+err.Error())
		_ = services.UpdateReportResults(context.Background(), job.ReportID, map[string]interface{}{"error": "Invalid axe-runner output"}, models.ReportStatusFailed)
		return
	}
	err = services.UpdateReportResults(context.Background(), job.ReportID, results, models.ReportStatusComplete)
	if err != nil {
		utils.LogAction(userID, "analyze", "failure", "Failed to update report: "+err.Error())
		return
	}
	utils.LogAction(userID, "analyze", "success", "Analysis complete for report "+job.ReportID.Hex())
	suggestions, err := services.GenerateSuggestionsFromLLM(results)
	if err != nil {
		utils.LogAction(userID, "llm_suggestion", "failure", "LLM error: "+err.Error())
	} else if len(suggestions) > 0 {
		err2 := services.CreateSuggestion(context.Background(), job.ReportID, suggestions)
		if err2 != nil {
			utils.LogAction(userID, "llm_suggestion", "failure", "Failed to save suggestions: "+err2.Error())
		} else {
			utils.LogAction(userID, "llm_suggestion", "success", "Suggestions saved for report "+job.ReportID.Hex())
		}
	} else {
		utils.LogAction(userID, "llm_suggestion", "failure", "No suggestions returned from LLM")
	}
}
