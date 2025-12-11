package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       int64 // in cents
	ImageURL    string
}

type Order struct {
	ID        string
	Items     []CartItem
	Total     int64
	Status    string
	CreatedAt time.Time
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
var orders = make(map[string]Order)

func init() {
	// Load .env file
	godotenv.Load()
}

func main() {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/cart", cartHandler)
	http.HandleFunc("/add-to-cart", addToCartHandler)
	http.HandleFunc("/checkout", checkoutHandler)
	http.HandleFunc("/process-payment", processPaymentHandler)
	http.HandleFunc("/order-confirmation", orderConfirmationHandler)

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
		"Products":  products,
		"CartCount": len(cart),
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

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to home if cart is empty
	if len(cart) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Create template with custom function
	tmpl := template.New("checkout.html").Funcs(template.FuncMap{
		"divf": func(a, b int64) float64 {
			return float64(a) / float64(b)
		},
	})

	tmpl, err := tmpl.ParseFiles("templates/checkout.html")
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
		"CartItems":        cartItems,
		"Total":            total,
		"SquareAppID":      os.Getenv("SQUARE_APPLICATION_ID"),
		"SquareLocationID": os.Getenv("SQUARE_LOCATION_ID"),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func processPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		SourceID string `json:"sourceId"`
		Email    string `json:"email"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate total
	var total int64
	for _, item := range cart {
		total += item.Product.Price * int64(item.Quantity)
	}

	// Create payment with Square using direct HTTP API call
	idempotencyKey := uuid.New().String()
	locationID := os.Getenv("SQUARE_LOCATION_ID")
	accessToken := os.Getenv("SQUARE_ACCESS_TOKEN")
	environment := os.Getenv("SQUARE_ENVIRONMENT")

	// Determine API endpoint based on environment
	apiURL := "https://connect.squareup.com/v2/payments"
	if environment == "sandbox" {
		apiURL = "https://connect.squareupsandbox.com/v2/payments"
	}

	// Create payment request body
	paymentBody := map[string]interface{}{
		"source_id":       requestBody.SourceID,
		"idempotency_key": idempotencyKey,
		"amount_money": map[string]interface{}{
			"amount":   total,
			"currency": "USD",
		},
		"location_id": locationID,
	}

	if requestBody.Email != "" {
		paymentBody["buyer_email_address"] = requestBody.Email
	}
	if requestBody.Name != "" {
		paymentBody["note"] = fmt.Sprintf("Order for %s", requestBody.Name)
	}

	// Marshal to JSON
	jsonBody, err := json.Marshal(paymentBody)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal error",
		})
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Request creation error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal error",
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Square-Version", "2024-12-18")

	// Make API call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Payment API error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Payment failed: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Parse response
	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("Response decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal error",
		})
		return
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		log.Printf("Payment failed: %+v", apiResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Payment failed",
		})
		return
	}

	// Extract payment ID
	payment := apiResponse["payment"].(map[string]interface{})
	paymentID := payment["id"].(string)

	// Create order record
	orderID := uuid.New().String()
	cartItems := make([]CartItem, 0, len(cart))
	for _, item := range cart {
		cartItems = append(cartItems, item)
	}

	order := Order{
		ID:        orderID,
		Items:     cartItems,
		Total:     total,
		Status:    "completed",
		CreatedAt: time.Now(),
	}
	orders[orderID] = order

	// Clear cart
	cart = make(map[string]CartItem)

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"orderID":   orderID,
		"paymentID": paymentID,
	})
}

func orderConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	order, exists := orders[orderID]

	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	tmpl := template.New("order-confirmation.html").Funcs(template.FuncMap{
		"divf": func(a, b int64) float64 {
			return float64(a) / float64(b)
		},
	})

	tmpl, err := tmpl.ParseFiles("templates/order-confirmation.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, order)
}
