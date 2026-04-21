package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Text string `json:"text"`
}

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var userMsg Message
	if err := json.NewDecoder(r.Body).Decode(&userMsg); err != nil {
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	apiUrl := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey

	reqBody := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{{Text: userMsg.Text}},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var geminiResp GeminiResponse
	json.Unmarshal(body, &geminiResp)

	// Perbaikan logika pengambilan balasan:
	botReply := "Maaf, AI tidak memberikan respon."
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates.Content.Parts) > 0 {
		botReply = geminiResp.Candidates.Content.Parts.Text
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Message{Text: botReply})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	http.HandleFunc("/chat", chatHandler)

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	fmt.Printf("Server running on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
