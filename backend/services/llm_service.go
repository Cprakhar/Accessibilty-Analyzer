package services

import (
	"backend/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
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
	prompt := `You are an expert web accessibility consultant. Analyze the following axe-core accessibility violations and provide actionable, developer-friendly suggestions for each issue as a JSON array. For each violation, use these exact attribute names:
- issue (string)
- summary (object: problem, impact, affectedUsers)
- whyMatters (object: userImpact, assistiveTechAffected)
- howToFix (object: step1, codeExample)
- testingInstructions (object: verify, tools)
- priorityLevel (object: urgency, wcagLevel)

Example:
[
  {
    "issue": "html-has-lang",
    "summary": {"problem": "...", "impact": "...", "affectedUsers": "..."},
    "whyMatters": {"userImpact": "...", "assistiveTechAffected": "..."},
    "howToFix": {"step1": "...", "codeExample": "..."},
    "testingInstructions": {"verify": "...", "tools": "..."},
    "priorityLevel": {"urgency": "...", "wcagLevel": "..."}
  }
]

Violations JSON:
<analysisResults>

Respond ONLY with a JSON array of suggestions, no extra explanation or formatting.`
	// Only send 'violations' to the LLM
	var violations interface{} = analysisResults
	if m, ok := analysisResults.(map[string]interface{}); ok {
		if v, exists := m["violations"]; exists {
			violations = v
		}
	}

	jsonResults, _ := json.Marshal(violations)
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
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
		// If still fails, return error with raw response for debugging
		return nil, errors.New("LLM did not return valid JSON array of suggestions. Raw response: " + raw)
	}
	return suggestions, nil
}
