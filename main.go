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

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
}

func main() {
	// Initialize database
	InitDatabase()

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Public routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	// Protected routes (require authentication)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/profile", authMiddleware(profileHandler))
	http.HandleFunc("/update-password", authMiddleware(updatePasswordHandler))
	http.HandleFunc("/cart", authMiddleware(cartHandler))
	http.HandleFunc("/add-to-cart", authMiddleware(addToCartHandler))
	http.HandleFunc("/remove-from-cart", authMiddleware(removeFromCartHandler))
	http.HandleFunc("/checkout", authMiddleware(checkoutHandler))
	http.HandleFunc("/process-payment", authMiddleware(processPaymentHandler))
	http.HandleFunc("/order-confirmation", authMiddleware(orderConfirmationHandler))

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

	// Get products from database
	var products []Product
	DB.Find(&products)

	// Check if user is logged in
	user, _ := getCurrentUser(r)

	// Get cart count for logged-in users
	var cartCount int64
	if user != nil {
		DB.Model(&CartItem{}).Where("user_id = ?", user.ID).Count(&cartCount)
	}

	data := map[string]interface{}{
		"Products":  products,
		"CartCount": cartCount,
		"User":      user,
	}
	tmpl.Execute(w, data)
}

func cartHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Create template with custom function
	tmpl := template.New("cart.html").Funcs(template.FuncMap{
		"divf": func(a, b int64) float64 {
			return float64(a) / float64(b)
		},
	})

	tmpl, err = tmpl.ParseFiles("templates/cart.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	// Get cart items from database with product info
	var cartItems []CartItem
	DB.Preload("Product").Where("user_id = ?", user.ID).Find(&cartItems)

	// Calculate total
	var total int64
	for _, item := range cartItems {
		total += item.Product.Price * int64(item.Quantity)
	}

	data := map[string]interface{}{
		"CartItems": cartItems,
		"Total":     total,
		"CartCount": len(cartItems),
		"User":      user,
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

	user, err := getCurrentUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not authenticated"})
		return
	}

	productID := r.FormValue("product_id")
	quantityStr := r.FormValue("quantity")
	quantity, _ := strconv.Atoi(quantityStr)
	if quantity < 1 {
		quantity = 1
	}

	// Check if product exists
	var product Product
	if err := DB.Where("id = ?", productID).First(&product).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Check if item already in cart
	var existingItem CartItem
	result := DB.Where("user_id = ? AND product_id = ?", user.ID, productID).First(&existingItem)

	if result.Error == nil {
		// Update quantity
		existingItem.Quantity += quantity
		DB.Save(&existingItem)
	} else {
		// Create new cart item
		newItem := CartItem{
			UserID:    user.ID,
			ProductID: productID,
			Quantity:  quantity,
		}
		DB.Create(&newItem)
	}

	// Get updated cart count
	var cartCount int64
	DB.Model(&CartItem{}).Where("user_id = ?", user.ID).Count(&cartCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"cartCount": cartCount,
	})
}

func removeFromCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	productID := r.FormValue("product_id")
	
	// Delete cart item
	DB.Where("user_id = ? AND product_id = ?", user.ID, productID).Delete(&CartItem{})

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get cart items
	var cartItems []CartItem
	DB.Preload("Product").Where("user_id = ?", user.ID).Find(&cartItems)

	if len(cartItems) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Create template with custom function
	tmpl := template.New("checkout.html").Funcs(template.FuncMap{
		"divf": func(a, b int64) float64 {
			return float64(a) / float64(b)
		},
	})

	tmpl, err = tmpl.ParseFiles("templates/checkout.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	// Calculate total
	var total int64
	for _, item := range cartItems {
		total += item.Product.Price * int64(item.Quantity)
	}

	data := map[string]interface{}{
		"CartItems":        cartItems,
		"Total":            total,
		"SquareAppID":      os.Getenv("SQUARE_APPLICATION_ID"),
		"SquareLocationID": os.Getenv("SQUARE_LOCATION_ID"),
		"User":             user,
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

	user, err := getCurrentUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not authenticated"})
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

	// Get cart items
	var cartItems []CartItem
	DB.Preload("Product").Where("user_id = ?", user.ID).Find(&cartItems)

	if len(cartItems) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cart is empty"})
		return
	}

	// Calculate total
	var total int64
	for _, item := range cartItems {
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

	// Create order in database
	order := Order{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Total:     total,
		Status:    "completed",
		PaymentID: paymentID,
		CreatedAt: time.Now(),
	}
	DB.Create(&order)

	// Create order items
	for _, item := range cartItems {
		orderItem := OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		}
		DB.Create(&orderItem)
	}

	// Clear user's cart
	DB.Where("user_id = ?", user.ID).Delete(&CartItem{})

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"orderID":   order.ID,
		"paymentID": paymentID,
	})
}

func orderConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	orderID := r.URL.Query().Get("id")

	// Get order from database
	var order Order
	result := DB.Preload("Items.Product").Where("id = ? AND user_id = ?", orderID, user.ID).First(&order)
	if result.Error != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	tmpl := template.New("order-confirmation.html").Funcs(template.FuncMap{
		"divf": func(a, b int64) float64 {
			return float64(a) / float64(b)
		},
	})

	tmpl, err = tmpl.ParseFiles("templates/order-confirmation.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"ID":        order.ID,
		"Items":     order.Items,
		"Total":     order.Total,
		"Status":    order.Status,
		"CreatedAt": order.CreatedAt,
		"User":      user,
	}

	tmpl.Execute(w, data)
}