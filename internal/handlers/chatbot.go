package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/naufalwaiz/hukum-chatbot/config"
)

type ChatRequest struct {
	Question string `json:"question"`
	UserID   string `json:"user_id"`
}

func Chatbot(c *gin.Context) {
	var req ChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Payload terbaru menggunakan "messages"
	payload := map[string]interface{}{
		"model": "models/gemini-pro",
		"input": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]string{
					{
						"type": "text",
						"text": "Kamu adalah asisten hukum Indonesia. Jawablah pertanyaan berikut hanya terkait hukum Indonesia: " + req.Question,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gemini request failed: %v", err)})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("DEBUG Gemini response:", string(respBody)) // debug

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Gemini API returned non-200 status",
			"raw_response": string(respBody),
		})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Failed to decode Gemini response",
			"raw_response": string(respBody),
		})
		return
	}

	// Ambil jawaban dari output → content → text
	outputs, ok := result["output"].([]interface{})
	if !ok || len(outputs) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "No output returned from Gemini",
			"raw_response": string(respBody),
		})
		return
	}

	firstOutput, ok := outputs[0].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Invalid output format",
			"raw_response": string(respBody),
		})
		return
	}

	contentArray, ok := firstOutput["content"].([]interface{})
	if !ok || len(contentArray) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "No content found in Gemini output",
			"raw_response": string(respBody),
		})
		return
	}

	firstContent, ok := contentArray[0].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Invalid content format",
			"raw_response": string(respBody),
		})
		return
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Text field not found in Gemini content",
			"raw_response": string(respBody),
		})
		return
	}

	answer := text

	// Simpan ke database
	_, _ = config.DB.Exec(
		"INSERT INTO chats (user_id, question, answer) VALUES ($1, $2, $3)",
		req.UserID, req.Question, answer,
	)

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
