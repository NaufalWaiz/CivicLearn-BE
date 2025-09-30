package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// AskGemini mengirim pertanyaan ke Gemini dan mengembalikan jawaban.
// Aman terhadap error, response kosong, atau format JSON baru.
func AskGemini(question string) (string, error) {
	payload := map[string]interface{}{
		"model": "models/gemini-pro",
		"input": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]string{
					{
						"type": "text",
						"text": "Kamu adalah asisten hukum Indonesia. Jawablah pertanyaan berikut hanya terkait hukum Indonesia: " + question,
					},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key="+os.Getenv("GEMINI_API_KEY"),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", fmt.Errorf("failed to request Gemini: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("DEBUG Gemini response:", string(respBody)) // debug

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to decode Gemini response: %v", err)
	}

	outputs, ok := result["output"].([]interface{})
	if !ok || len(outputs) == 0 {
		return "", fmt.Errorf("no output returned from Gemini")
	}

	firstOutput, ok := outputs[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid output format")
	}

	contentArray, ok := firstOutput["content"].([]interface{})
	if !ok || len(contentArray) == 0 {
		return "", fmt.Errorf("no content found in Gemini output")
	}

	firstContent, ok := contentArray[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content format")
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("text field not found in Gemini content")
	}

	return text, nil
}
