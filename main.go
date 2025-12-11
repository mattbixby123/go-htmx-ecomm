package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       int64 // in cents
	ImageURL    string
}

var products = []Product{
	{
		ID:          "1",
		Name:        "Premium Headphones",
		Description: "High-quality wireless headphones with noise cancellation",
		Price:       29900,
		ImageURL:    "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400",
	},
	{
		ID:          "2",
		Name:        "Smart Watch",
		Description: "Fitness tracking smartwatch with heart rate monitor",
		Price:       19900,
		ImageURL:    "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400",
	},
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Create template with custom function
	tmpl := template.New("home.html").Funcs(template.FuncMap{
		"divf": func(a, b int64) float64 {
			return float64(a) / float64(b)
		},
	})

	tmpl, err := tmpl.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Products": products,
	}
	tmpl.Execute(w, data)
}
