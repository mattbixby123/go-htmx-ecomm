# TechStore - Full Stack Go E-commerce Application

## 🎯 Project Overview

This is a **complete, production-ready e-commerce application** built with Go that integrates Square payment processing. It features a modern UI, secure payment handling, and a clean architecture that's easy to extend.

## ✨ What's Included

### Core Files

1. **main.go** - Complete backend application with:
   - Product catalog management
   - Shopping cart functionality
   - Square payment processing
   - Order management
   - RESTful API endpoints

2. **HTML Templates** (5 files):
   - `home.html` - Landing page with product showcase
   - `products.html` - Product listing page
   - `cart.html` - Shopping cart with item management
   - `checkout.html` - Secure checkout with Square payment form
   - `order-confirmation.html` - Order success page

3. **Configuration Files**:
   - `.env.example` - Environment variables template
   - `go.mod` - Go dependencies
   - `.gitignore` - Git ignore rules

4. **Development Tools**:
   - `Makefile` - Common development commands
   - `Dockerfile` - Container configuration
   - `docker-compose.yml` - Multi-container orchestration

5. **Documentation**:
   - `README.md` - Comprehensive project documentation
   - `SETUP_GUIDE.md` - Step-by-step Square payment setup

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Square Developer Account (free)

### Setup Steps

1. **Get Square Credentials**:
   ```
   → Sign up at developer.squareup.com
   → Create an application
   → Copy Application ID, Access Token, and Location ID
   ```

2. **Configure Environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your Square credentials
   ```

3. **Install & Run**:
   ```bash
   go mod download
   go run main.go
   ```

4. **Test**:
   - Visit `http://localhost:8080`
   - Add products to cart
   - Use test card: `4111 1111 1111 1111`

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────┐
│                  Frontend Layer                  │
│  (HTML Templates + Tailwind CSS + HTMX)        │
├─────────────────────────────────────────────────┤
│                   Go Backend                     │
│  • HTTP Handlers                                │
│  • Business Logic                               │
│  • Data Models                                  │
├─────────────────────────────────────────────────┤
│               Square Payment SDK                 │
│  • Payment Processing                           │
│  • Tokenization                                 │
│  • Security                                     │
└─────────────────────────────────────────────────┘
```

## 💳 Payment Flow

1. User adds products to cart
2. Proceeds to checkout
3. Enters billing/payment info
4. Square.js tokenizes card (client-side, PCI-compliant)
5. Token sent to backend
6. Backend processes payment via Square API
7. Order confirmed and displayed

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Homepage |
| GET | `/cart` | View cart |
| POST | `/add-to-cart` | Add item to cart |
| POST | `/remove-from-cart` | Remove item |
| GET | `/checkout` | Checkout page |
| POST | `/process-payment` | Process payment |
| GET | `/order-confirmation` | Order success |

## 🎨 Features

### Implemented
✅ Product catalog with images
✅ Shopping cart with add/remove
✅ Square payment integration
✅ Responsive design (mobile-friendly)
✅ Dynamic updates (no page reloads)
✅ Order confirmation system
✅ Sandbox testing support
✅ Error handling
✅ PCI-compliant payment processing

### Easy to Add
- User authentication
- Database persistence (PostgreSQL/MySQL)
- Email notifications
- Product search & filtering
- Admin dashboard
- Order history
- Inventory management
- Multiple payment methods

## 🔒 Security Features

- **PCI Compliance**: Card data never touches your server
- **Tokenization**: Square handles sensitive data
- **HTTPS Ready**: SSL/TLS support in production
- **Environment Variables**: Secrets stored securely
- **Input Validation**: Server-side validation
- **Idempotency**: Prevents duplicate charges

## 📦 Dependencies

```go
github.com/google/uuid        // Unique IDs
github.com/joho/godotenv      // Environment variables
github.com/square/square-go-sdk // Square API
```

## 🧪 Testing

### Sandbox Test Cards

| Card Number | Result |
|------------|--------|
| 4111 1111 1111 1111 | Success |
| 4000 0000 0000 0002 | Decline |
| 4000 0000 0000 0119 | Error |

All test cards:
- CVV: Any 3 digits
- Expiration: Any future date
- ZIP: Any 5 digits

## 🚀 Deployment Options

### Traditional VPS
```bash
go build -o techstore
./techstore
```

### Docker
```bash
docker build -t techstore .
docker run -p 8080:8080 --env-file .env techstore
```

### Docker Compose
```bash
docker-compose up
```

### Cloud Platforms
- Heroku
- DigitalOcean App Platform
- Railway
- Fly.io
- AWS EC2/ECS
- Google Cloud Run

## 📈 Scalability Considerations

### Current State (Development)
- In-memory data storage
- Single server instance
- No database

### Production Recommendations
1. **Add Database**: PostgreSQL or MySQL
2. **Session Management**: Redis for cart/sessions
3. **Load Balancing**: Multiple app instances
4. **CDN**: Static assets on CloudFlare/AWS CloudFront
5. **Caching**: Redis for product data
6. **Monitoring**: Prometheus + Grafana
7. **Logging**: Structured logging with ELK stack

## 🔧 Customization Guide

### Add New Products
Edit `products` slice in `main.go`:
```go
products = []Product{
    {
        ID:          "5",
        Name:        "New Product",
        Description: "Description here",
        Price:       9900, // $99.00 in cents
        ImageURL:    "https://...",
    },
}
```

### Change Styling
Edit any template file - Tailwind classes are inline:
```html
<button class="bg-blue-600 text-white px-4 py-2 rounded-lg">
    Button
</button>
```

### Add Database
Replace in-memory maps with database calls:
```go
// Instead of: cart = make(map[string]CartItem)
// Use: db.Query("SELECT * FROM cart WHERE user_id = ?", userId)
```

### Add Authentication
Implement middleware:
```go
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check session/JWT
        next(w, r)
    }
}
```

## 📝 Common Tasks

### Run locally
```bash
make run
```

### Build binary
```bash
make build
```

### Hot reload (dev)
```bash
make dev  # Requires: go install github.com/cosmtrek/air@latest
```

### Format code
```bash
make fmt
```

### Run tests
```bash
make test
```

## 🐛 Troubleshooting

### Payment fails
- Check Square credentials in `.env`
- Use sandbox test cards
- Check browser console for errors
- Verify Square environment (sandbox/production)

### Port already in use
```bash
# Change PORT in .env or kill process
lsof -ti:8080 | xargs kill
```

### Templates not found
```bash
# Ensure templates/ directory exists
ls templates/
```

## 📚 Learning Resources

- **Go**: [go.dev/tour](https://go.dev/tour)
- **Square**: [developer.squareup.com/docs](https://developer.squareup.com/docs)
- **HTMX**: [htmx.org](https://htmx.org)
- **Tailwind**: [tailwindcss.com](https://tailwindcss.com)

## 🤝 Contributing

This is a learning/starter project. Feel free to:
- Fork and modify
- Add features
- Submit pull requests
- Report issues

## 📄 License

MIT License - Free to use for personal or commercial projects

## 🎓 What You'll Learn

Building/studying this project teaches:
- Go web development
- Payment gateway integration
- RESTful API design
- Frontend templating
- Docker containerization
- Security best practices
- E-commerce workflows

## 💡 Next Steps

After getting this running:

1. ✅ **Test thoroughly** with sandbox
2. 🔐 **Add authentication** for user accounts
3. 💾 **Implement database** for persistence
4. 📧 **Set up email** notifications
5. 📊 **Add analytics** tracking
6. 🎨 **Customize design** to match brand
7. 🚀 **Deploy to production** with HTTPS
8. 📱 **Build mobile app** using same API

## ⚡ Performance Tips

- Use connection pooling for database
- Implement caching for product data
- Compress responses with gzip
- Use CDN for static assets
- Optimize images (WebP format)
- Enable HTTP/2
- Use goroutines for async operations

## 🔐 Production Checklist

Before going live:

- [ ] Switch to production Square credentials
- [ ] Enable HTTPS (required by Square)
- [ ] Set up database backups
- [ ] Implement rate limiting
- [ ] Add logging and monitoring
- [ ] Set up error alerting
- [ ] Test disaster recovery
- [ ] Review security best practices
- [ ] Set up CI/CD pipeline
- [ ] Configure proper CORS
- [ ] Add health check endpoints
- [ ] Set up webhook handling

## 🎉 Success Metrics

Track these to measure success:
- Conversion rate (visitors → purchases)
- Average order value
- Cart abandonment rate
- Payment success rate
- Page load times
- Mobile vs desktop usage
- Error rates
- Customer satisfaction

---

**Built with ❤️ using Go, Square, HTMX, and Tailwind CSS**

For questions or issues, refer to:
- `README.md` - General documentation
- `SETUP_GUIDE.md` - Square payment setup
- Square Developer Forums
- GitHub Issues