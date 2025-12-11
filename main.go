package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       int64 // in cents
	ImageURL    string
}

type CartItem struct {
	Product  Product
	Quantity int
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

var cart = make(map[string]CartItem)

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
	http.HandleFunc("/cart", cartHandler)
	http.HandleFunc("/add-to-cart", addToCartHandler)

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

func cartHandler(w http.ResponseWriter, r *http.Request) {
	// Create template with custom function (same as homeHandler)
	tmpl := template.New("cart.html").Funcs(template.FuncMap{
		"divf": func(a, b int64) float64 {
			return float64(a) / float64(b)
		},
	})

	tmpl, err := tmpl.ParseFiles("templates/cart.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	cartItems := make([]CartItem, 0, len(cart))
	var total int64
	for _, item := range cart {
		cartItems = append(cartItems, item)
		total += item.Product.Price * int64(item.Quantity)
	}

	data := map[string]interface{}{
		"CartItems": cartItems,
		"Total":     total,
		"CartCount": len(cart),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func addToCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productID := r.FormValue("product_id")
	quantityStr := r.FormValue("quantity")
	quantity, _ := strconv.Atoi(quantityStr)
	if quantity < 1 {
		quantity = 1
	}

	// Find product
	var product *Product
	for _, p := range products {
		if p.ID == productID {
			product = &p
			break
		}
	}

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Add to cart
	if existingItem, exists := cart[productID]; exists {
		existingItem.Quantity += quantity
		cart[productID] = existingItem
	} else {
		cart[productID] = CartItem{
			Product:  *product,
			Quantity: quantity,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"cartCount": len(cart),
	})
}
