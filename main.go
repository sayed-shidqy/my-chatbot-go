package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Text string `json:"text"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := strings.ToLower(msg.Text)
	var response string

	// Logika Bot yang lebih interaktif
	if strings.Contains(input, "halo") || strings.Contains(input, "hai") {
		response = "Halo! Saya adalah asisten virtual Sayed Shidqy. Senang bertemu dengan Anda!"
	} else if strings.Contains(input, "siapa") && strings.Contains(input, "anda") {
		response = "Saya adalah bot cerdas yang dideploy menggunakan Docker. Saya dibuat untuk membantu demonstrasi tugas ini."
	} else if strings.Contains(input, "dosen") || strings.Contains(input, "pak") || strings.Contains(input, "bu") {
		response = "Selamat datang, Bapak/Ibu Dosen! Terima kasih sudah meluangkan waktu untuk memeriksa proyek Docker Sayed."
	} else if strings.Contains(input, "jam") || strings.Contains(input, "waktu") {
		response = "Saat ini menunjukkan pukul " + time.Now().Format("15:04:05") + " WIB."
	} else if strings.Contains(input, "docker") {
		response = "Aplikasi ini berjalan di dalam container Docker menggunakan image Alpine Linux yang sangat ringan!"
	} else {
		response = "Pertanyaan yang menarik! Sayangnya saya masih dalam tahap pengembangan. Ada hal lain yang ingin Bapak/Ibu ketahui tentang sistem ini?"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Message{Text: response})
}

func main() {
	// Sajikan file statis dari folder templates
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	
	// API endpoint untuk chat
	http.HandleFunc("/chat", chatHandler)

	fmt.Println("-------------------------------------------")
	fmt.Println("🚀 Server Bot Sayed Aktif!")
	fmt.Println("🌍 Akses di: http://localhost:8080")
	fmt.Println("-------------------------------------------")
	
	http.ListenAndServe(":8080", nil)
}
