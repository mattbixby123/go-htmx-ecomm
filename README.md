# TechStore - Full Stack Go E-commerce with Square Payments

A complete e-commerce application built with Go, featuring Square payment integration, HTMX for dynamic interactions, and Tailwind CSS for styling.

## Features

- 🛍️ Product catalog with detailed views
- 🛒 Shopping cart functionality
- 💳 Secure payment processing with Square
- 📱 Responsive design with Tailwind CSS
- ⚡ Dynamic updates with HTMX (no page reloads)
- 🎨 Modern UI with product images
- 📧 Order confirmation system
- 🔒 PCI-compliant payment handling

## Tech Stack

- **Backend**: Go (Golang)
- **Payment Processing**: Square SDK
- **Frontend**: HTML Templates, Tailwind CSS, HTMX
- **No Database**: In-memory storage (easily extensible to PostgreSQL/MySQL)

## Prerequisites

- Go 1.21 or higher
- Square Developer Account (free)
- Git

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/techstore.git
cd techstore
```

### 2. Set Up Square Account

1. Go to [Square Developer Dashboard](https://developer.squareup.com/apps)
2. Create a new application
3. Get your credentials:
   - **Application ID**: Found in your app's credentials
   - **Access Token**: Generate a sandbox access token
   - **Location ID**: Found in Locations section

### 3. Configure Environment Variables

```bash
cp .env.example .env
```

Edit `.env` and add your Square credentials:

```env
SQUARE_APPLICATION_ID=your_application_id
SQUARE_ACCESS_TOKEN=your_access_token
SQUARE_LOCATION_ID=your_location_id
SQUARE_ENVIRONMENT=sandbox
PORT=8080
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run main.go
```

Visit `http://localhost:8080` in your browser.

## Project Structure

```
techstore/
├── main.go                 # Main application logic
├── go.mod                  # Go dependencies
├── .env                    # Environment variables (create from .env.example)
├── .env.example           # Example environment configuration
├── templates/
│   ├── home.html          # Homepage with product listing
│   ├── products.html      # Products page
│   ├── cart.html          # Shopping cart
│   ├── checkout.html      # Checkout with Square payment form
│   └── order-confirmation.html  # Order success page
└── README.md
```

## Using Square Payments

### Sandbox Testing

Square provides test card numbers for sandbox testing:

- **Success**: `4111 1111 1111 1111`
- **CVV**: Any 3 digits
- **Expiration**: Any future date
- **ZIP**: Any 5 digits

### Going to Production

1. Switch to production credentials in `.env`:
   ```env
   SQUARE_ENVIRONMENT=production
   SQUARE_ACCESS_TOKEN=your_production_token
   ```

2. Update the Square.js script in `templates/checkout.html`:
   ```html
   <!-- Change from sandbox to production -->
   <script type="text/javascript" src="https://web.squarecdn.com/v1/square.js"></script>
   ```

3. Complete Square's verification process for production access

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Homepage with products |
| GET | `/products` | Products listing page |
| GET | `/cart` | View shopping cart |
| POST | `/add-to-cart` | Add item to cart |
| POST | `/remove-from-cart` | Remove item from cart |
| GET | `/checkout` | Checkout page with payment form |
| POST | `/process-payment` | Process Square payment |
| GET | `/order-confirmation` | Order success page |

## Key Features Explained

### 1. Square Payment Integration

The checkout flow uses Square's Web Payments SDK:

```javascript
// Initialize Square payments
const payments = window.Square.payments(appId, locationId);
const card = await payments.card();
await card.attach('#card-container');

// Tokenize and process payment
const token = await card.tokenize();
await createPayment(token);
```

### 2. In-Memory Data Storage

Current implementation uses Go maps for storage:

```go
var (
    cart   = make(map[string]CartItem)
    orders = make(map[string]Order)
)
```

**Note**: Data resets on server restart. For production, implement database persistence (PostgreSQL, MySQL, or MongoDB).

### 3. HTMX Dynamic Updates

Add to cart without page reload:

```html
<button onclick="addToCart('product-id')">Add to Cart</button>
```

## Extending the Application

### Add Database Persistence

Replace in-memory maps with a database:

```go
// Example with PostgreSQL
import "database/sql"
import _ "github.com/lib/pq"

db, err := sql.Open("postgres", "connection-string")
```

### Add User Authentication

Implement user sessions:

```go
import "github.com/gorilla/sessions"

var store = sessions.NewCookieStore([]byte("secret-key"))
```

### Add Product Categories

Extend the `Product` struct:

```go
type Product struct {
    ID          string
    Name        string
    Description string
    Price       int64
    ImageURL    string
    Category    string  // New field
    Stock       int     // New field
}
```

### Add Email Notifications

Use a mail service:

```go
import "net/smtp"

// Send order confirmation email
func sendConfirmationEmail(email, orderID string) error {
    // Implementation
}
```

## Security Best Practices

1. **Never commit `.env` file** - It contains sensitive credentials
2. **Use HTTPS in production** - Required by Square
3. **Validate user inputs** - Prevent injection attacks
4. **Rate limit API endpoints** - Prevent abuse
5. **Implement CSRF protection** - For production apps

## Common Issues

### "Failed to initialize payment form"

- Check your `SQUARE_APPLICATION_ID` is correct
- Verify `SQUARE_LOCATION_ID` matches your account
- Ensure you're using sandbox credentials for testing

### Payment always fails

- Use Square's test card numbers in sandbox mode
- Check `SQUARE_ACCESS_TOKEN` has payment permissions
- Verify `SQUARE_ENVIRONMENT` is set to "sandbox"

### Port already in use

```bash
# Change port in .env
PORT=3000

# Or kill existing process
lsof -ti:8080 | xargs kill
```

## Testing

Test the full checkout flow:

1. Add products to cart
2. Proceed to checkout
3. Use test card: `4111 1111 1111 1111`
4. Submit payment
5. Verify order confirmation page

## Deployment

### Deploy to a VPS (DigitalOcean, AWS, etc.)

```bash
# Build the binary
go build -o techstore

# Run with environment variables
export $(cat .env | xargs) && ./techstore
```

### Deploy with Docker

Create `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o techstore
CMD ["./techstore"]
```

Build and run:

```bash
docker build -t techstore .
docker run -p 8080:8080 --env-file .env techstore
```

## Resources

- [Square Developer Docs](https://developer.squareup.com/docs)
- [Square Go SDK](https://github.com/square/square-go-sdk)
- [HTMX Documentation](https://htmx.org/)
- [Tailwind CSS](https://tailwindcss.com/)

## License

MIT License - feel free to use this project for learning or commercial purposes.

## Support

For issues or questions:
- Open an issue on GitHub
- Check Square's developer forums
- Review Go documentation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a pull request

## Future Enhancements

- [ ] User authentication and accounts
- [ ] Product search and filtering
- [ ] Order history
- [ ] Admin dashboard
- [ ] Inventory management
- [ ] Email notifications
- [ ] Database persistence
- [ ] Product reviews
- [ ] Wishlist functionality
- [ ] Multiple payment methods

---

Built with ❤️ using Go and Square