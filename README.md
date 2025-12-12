# Go HTMX E-Commerce Store

A full-stack e-commerce application built with Go, HTMX, PostgreSQL, and Square Payments.

## 🔐 Security Features

- **Client-side password hashing**: PBKDF2 with 1000 iterations (email-based salt)
- **Server-side password hashing**: bcrypt with cost factor 12
- **JWT authentication**: Secure token-based authentication
- **User-specific carts**: Each user has their own isolated cart
- **Session management**: Database-backed session tracking with expiration

## 🚀 Features

- User registration and authentication
- Product browsing
- Shopping cart (user-specific)
- Square payment integration (sandbox & production)
- Order history
- Password management
- Responsive design with Tailwind CSS

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 17 (or 15+)
- Square Developer Account (for payment processing)

## 🛠️ Setup Instructions

### 1. Update PostgreSQL PATH

Edit your `~/.zshrc`:

```bash
code ~/.zshrc
```

Change this line:
```bash
export PATH="$HOMEBREW_PREFIX/opt/postgresql@15/bin:$PATH"
```

To:
```bash
export PATH="$HOMEBREW_PREFIX/opt/postgresql@17/bin:$PATH"
```

Then reload:
```bash
source ~/.zshrc
```

### 2. Start PostgreSQL and Create Database

```bash
# Start PostgreSQL
brew services start postgresql@17

# Create database
createdb techstore

# Verify connection
psql techstore
# Inside psql:
SELECT version();
\q
```

### 3. Install Go Dependencies

```bash
go mod download
```

Or install individually:
```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/google/uuid
go get github.com/joho/godotenv
```

### 4. Configure Environment Variables

Copy the example env file:
```bash
cp .env.example .env
```

Edit `.env` and add your Square credentials:
```env
DATABASE_URL=postgres://localhost:5432/techstore?sslmode=disable

SQUARE_APPLICATION_ID=your_app_id
SQUARE_ACCESS_TOKEN=your_access_token
SQUARE_LOCATION_ID=your_location_id
SQUARE_ENVIRONMENT=sandbox

JWT_SECRET=generate-a-random-secret-here

PORT=8080
```

**Generate JWT Secret:**
```bash
openssl rand -base64 32
```

**Get Square Credentials:**
1. Go to https://developer.squareup.com/
2. Create an application
3. Get Application ID and Access Token
4. Use API Explorer to get Location ID

### 5. Run the Application

```bash
go run main.go auth.go database.go models.go
```

Or build and run:
```bash
go build -o ecommerce
./ecommerce
```

The server will start on http://localhost:8080

## 📁 Project Structure

```
.
├── main.go              # Main application logic, routes, cart, payment
├── auth.go              # JWT authentication and user management
├── database.go          # PostgreSQL connection and migrations
├── models.go            # Database models (User, Product, Order, etc.)
├── templates/           # HTML templates
│   ├── home.html
│   ├── login.html
│   ├── register.html
│   ├── profile.html
│   ├── cart.html
│   ├── checkout.html
│   └── order-confirmation.html
├── .env                 # Environment variables (not in git)
├── .env.example         # Example environment file
└── go.mod               # Go dependencies
```

## 🗄️ Database Schema

The application automatically creates these tables:

- **users**: User accounts with bcrypt-hashed passwords
- **sessions**: JWT session tracking
- **products**: Product catalog (auto-seeded with 4 products)
- **cart_items**: User-specific shopping carts
- **orders**: Completed orders
- **order_items**: Products in each order

## 🔄 How It Works

### Authentication Flow

1. **Registration**:
   - Client hashes password with PBKDF2 (1000 iterations, email-based salt)
   - Server hashes again with bcrypt (cost 12, automatic salt)
   - JWT token generated and returned
   - Token stored in HttpOnly cookie

2. **Login**:
   - Client hashes password with PBKDF2
   - Server verifies against bcrypt hash
   - JWT token generated
   - User data returned

3. **Protected Routes**:
   - JWT token extracted from cookie or Authorization header
   - Token validated and user verified
   - User data loaded from database

### Shopping Flow

1. User browses products
2. Adds items to cart (stored in PostgreSQL with user_id)
3. Views cart (user-specific items only)
4. Proceeds to checkout
5. Enters payment info (Square Web Payments SDK)
6. Payment processed via Square API
7. Order created in database
8. Cart cleared
9. Order confirmation displayed

## 🧪 Testing

### Test Square Payment

Use these test card numbers in sandbox mode:

- **Card Number**: `4111 1111 1111 1111`
- **CVV**: `111`
- **Expiration**: Any future date
- **Zip**: Any 5 digits

## 📊 Database Management

### View data in psql:

```bash
psql techstore

# List tables
\dt

# View users
SELECT * FROM users;

# View products
SELECT * FROM products;

# View cart items
SELECT * FROM cart_items;

# View orders
SELECT * FROM orders;
```

### Reset database:

```bash
# Drop database
dropdb techstore

# Recreate
createdb techstore

# Restart app (migrations will run automatically)
go run main.go auth.go database.go models.go
```

## 🔒 Security Notes

**Development:**
- JWT tokens stored in HttpOnly cookies
- HTTPS disabled (set `Secure: true` in production)
- Client-side PBKDF2 prevents plain-text password transmission

**Production:**
- Use HTTPS (set `Secure: true` on cookies)
- Change JWT_SECRET to a strong random value
- Use Square production environment
- Consider rate limiting
- Add CSRF protection
- Enable GORM query logging only in development

## 🚧 Production Deployment

1. Set up HTTPS (Let's Encrypt, Cloudflare, etc.)
2. Update `.env`:
   ```env
   SQUARE_ENVIRONMENT=production
   JWT_SECRET=<strong-random-secret>
   ```
3. Update cookie settings in `auth.go`:
   ```go
   Secure: true,  // Enable in production
   ```
4. Use a production PostgreSQL database
5. Enable database backups
6. Set up monitoring and logging

## 🐛 Troubleshooting

**"Database connection failed":**
- Ensure PostgreSQL is running: `brew services list`
- Check DATABASE_URL in `.env`
- Verify database exists: `psql -l`

**"Payment failed":**
- Check Square credentials in `.env`
- Verify SQUARE_ENVIRONMENT is set to "sandbox"
- Check Square dashboard for errors

**"JWT token invalid":**
- Check JWT_SECRET is set
- Clear browser cookies and re-login
- Verify token hasn't expired (7 days)

## 📝 License

MIT

## 👨‍💻 Author

Built with Go, PostgreSQL, HTMX, Tailwind CSS, and Square Payments API