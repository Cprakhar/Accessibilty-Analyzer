package services

import (
	"backend/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kaptinlin/jsonrepair"
)

type LLMResponse struct {
	Suggestions []string `json:"suggestions"`
}

func GenerateSuggestionsFromLLM(analysisResults interface{}) ([]models.SuggestionItem, error) {
	llmApiUrl := os.Getenv("LLM_API_URL")
	llmApiKey := os.Getenv("LLM_API_KEY")
	if llmApiUrl == "" || llmApiKey == "" {
		return nil, errors.New("LLM API URL or KEY not set")
	}

	// Format concise prompt for LLM (OpenAI-compatible)
	prompt := `You are an expert web accessibility consultant. Analyze the following axe-core accessibility violations and provide actionable, developer-friendly suggestions for each issue as a JSON array.

Now consider the following TypeScript interface for the JSON schema:

interface ViolationSummary {
    problem: string;
    impact: "critical" | "serious" | "moderate" | "minor";
    affectedUsers: string;
}

interface WhyMatters {
    userImpact: string;
    assistiveTechAffected: string;
}

interface HowToFix {
    step1: string;
    codeExample: string;
}

interface TestingInstructions {
    verify: string;
    tools: string;
}

interface PriorityLevel {
    urgency: "high" | "medium" | "low";
    wcagLevel: "A" | "AA" | "AAA";
}

interface AccessibilityViolation {
    issue: string;
    summary: ViolationSummary;
    whyMatters: WhyMatters;
    howToFix: HowToFix;
    testingInstructions: TestingInstructions;
    priorityLevel: PriorityLevel;
}

interface AccessibilityAnalysis {
    violations: AccessibilityViolation[];
}

IMPORTANT FORMATTING RULES:
- Use single quotes for all HTML attribute values in codeExample fields e.g.: <html lang='en'>
- Ensure all JSON is valid and parseable
- Match the exact field names from the TypeScript interface
- Use only the specified enum values for impact, urgency, and wcagLevel

Violations JSON:
<analysisResults>

Write the accessibility analysis according to the AccessibilityAnalysis schema.
On the response, include only the JSON. No additional text, explanations, or formatting.`
	// Only send up to 5 violations to the LLM to avoid exceeding request length
	var violations interface{} = analysisResults
	if m, ok := analysisResults.(map[string]interface{}); ok {
		if v, exists := m["violations"]; exists {
			if arr, ok := v.([]interface{}); ok && len(arr) > 5 {
				violations = arr[:5]
			} else {
				violations = v
			}
		}
	}

	// Wrap violations in the expected schema for the LLM prompt
	analysis := map[string]interface{}{
		"violations": violations,
	}
	jsonResults, _ := json.Marshal(analysis)
	prompt = strings.Replace(prompt, "<analysisResults>", string(jsonResults), 1)

	payload := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"model": os.Getenv("GROQ_MODEL"),
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", llmApiUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llmApiKey)
	var resp *http.Response
	maxRetries := 5
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 429 {
			resp.Body.Close()
			// Exponential backoff: 1s, 2s, 4s, 8s, 16s
			delay := 1 << attempt
			println("[LLM] 429 Too Many Requests, retrying in", delay, "seconds...")
			time.Sleep(time.Duration(delay) * time.Second)
			continue
		}
		break
	}
	if resp.StatusCode == 429 {
		return nil, errors.New("LLM API error: 429 Too Many Requests (rate limited after retries)")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("LLM API error: " + resp.Status)
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil || len(chatResp.Choices) == 0 {
		return nil, errors.New("invalid LLM response")
	}

	raw := chatResp.Choices[0].Message.Content
	// Log the raw LLM response for debugging
	println("[LLM RAW RESPONSE]", raw)

	// Attempt to clean up common LLM formatting issues
	clean := raw
	clean = strings.TrimSpace(clean)
	// Remove code block markers if present
	if strings.HasPrefix(clean, "```") {
		clean = strings.TrimPrefix(clean, "```")
		clean = strings.TrimSpace(clean)
		if idx := strings.Index(clean, "\n"); idx != -1 {
			clean = clean[idx+1:]
		}
		clean = strings.TrimSuffix(clean, "```")
	}
	// Remove leading/trailing brackets if duplicated
	clean = strings.TrimPrefix(clean, "[")
	clean = strings.TrimSuffix(clean, "]")
	clean = "[" + clean + "]"

	var suggestions []models.SuggestionItem
	err = json.Unmarshal([]byte(clean), &suggestions)
	if err != nil {
		// Try to repair the JSON using github.com/kaptinlin/jsonrepair
		repaired, repairErr := jsonrepair.JSONRepair(clean)
		if repairErr == nil {
			err = json.Unmarshal([]byte(repaired), &suggestions)
			if err == nil {
				return suggestions, nil
			}
		}
		// If still fails, return error with raw response for debugging
		return nil, errors.New("LLM did not return valid JSON array of suggestions. Raw response: " + raw)
	}
	return suggestions, nil
}
