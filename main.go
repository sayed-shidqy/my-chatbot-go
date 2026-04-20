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

// Struktur untuk request ke Gemini
type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

// Struktur untuk respon dari Gemini
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
	json.NewDecoder(r.Body).Decode(&userMsg)

	// AMBIL API KEY DARI ENV ATAU MASUKKAN LANGSUNG (Saran: Pakai ENV di Railway)
	apiKey := os.Getenv("GEMINI_API_KEY") 
	if apiKey == "" {
		apiKey = "AIzaSyDxxBRXNq1vbrW2TuuuDMqMTHjCSpMSAjA" // Ganti dengan key Anda
	}

	apiUrl := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey

	// Siapkan payload untuk Gemini
	geminiReq := GeminiRequest{}
	geminiReq.Contents = append(geminiReq.Contents, struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{
		Parts: []struct {
			Text string `json:"text"`
		}{{Text: userMsg.Text}},
	})

	jsonData, _ := json.Marshal(geminiReq)
	
	// Kirim request ke Google
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Gagal kontak Gemini", 500)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var geminiResp GeminiResponse
	json.Unmarshal(body, &geminiResp)

	// Ambil teks jawaban dari Gemini
	botReply := "Maaf, saya sedang lelah."
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
	
	fmt.Printf("Server berjalan di port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
