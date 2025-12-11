# Quick Start Guide - TechStore

Get your e-commerce store running in **5 minutes**!

## Prerequisites

- Go 1.21+ installed ([download here](https://go.dev/dl/))
- Web browser
- Text editor

## Step 1: Get Square Credentials (2 minutes)

1. Visit [https://developer.squareup.com/](https://developer.squareup.com/)
2. Sign up for free
3. Create a new application
4. Go to **Credentials** → **Sandbox**
5. Copy these three values:
   - Application ID
   - Access Token
   - Location ID (from Locations tab)

## Step 2: Setup Project (1 minute)

```bash
# Navigate to the project directory
cd techstore

# Create .env file from example
cp .env.example .env

# Edit .env and paste your Square credentials
nano .env  # or use any text editor
```

Your `.env` should look like:
```env
SQUARE_APPLICATION_ID=sandbox-sq0idb-xxxxxx
SQUARE_ACCESS_TOKEN=EAAAl-xxxxxx
SQUARE_LOCATION_ID=Lxxxxxx
SQUARE_ENVIRONMENT=sandbox
PORT=8080
```

## Step 3: Run the Application (30 seconds)

```bash
# Install dependencies
go mod download

# Start the server
go run main.go
```

You should see:
```
Server starting on port 8080
```

## Step 4: Test It! (1 minute)

1. Open your browser to `http://localhost:8080`
2. Click **Add to Cart** on any product
3. Click **Cart** in the navigation
4. Click **Proceed to Checkout**
5. Fill in the form:
   - Name: Test User
   - Email: test@example.com
   - Address: 123 Test St
   - City: San Francisco
   - State: CA
   - ZIP: 94103
6. Enter payment info:
   - Card: `4111 1111 1111 1111`
   - CVV: `123`
   - Expiration: `12/25`
   - ZIP: `94103`
7. Click **Pay Now**
8. See your order confirmation! 🎉

## Common Issues

### "Failed to initialize payment form"
- Double-check your `SQUARE_APPLICATION_ID` in `.env`
- Make sure there are no extra spaces

### "Port 8080 already in use"
```bash
# Change PORT in .env to something else
PORT=3000
```

### Payment fails
- Make sure you're using the test card: `4111 1111 1111 1111`
- Check that all Square credentials are correct

## What's Next?

✅ **Customize**: Edit product data in `main.go`
✅ **Style**: Modify Tailwind classes in HTML templates
✅ **Deploy**: See `README.md` for deployment options
✅ **Production**: Follow `SETUP_GUIDE.md` for production setup

## Quick Commands

```bash
# Run the app
make run

# Build binary
make build

# View all available commands
make help
```

## File Structure

```
techstore/
├── main.go              # Backend logic
├── templates/           # HTML pages
│   ├── home.html
│   ├── cart.html
│   ├── checkout.html
│   └── ...
├── .env                 # Your Square credentials
├── README.md           # Full documentation
├── SETUP_GUIDE.md      # Detailed Square setup
└── PROJECT_OVERVIEW.md # Architecture overview
```

## Need Help?

- 📖 Read `SETUP_GUIDE.md` for detailed Square setup
- 📚 Check `README.md` for full documentation
- 🐛 Found a bug? Open an issue on GitHub
- 💬 Questions? Check Square Developer Forums

## Test Cards

All these work in sandbox mode:

| Card Number | Result |
|------------|--------|
| 4111 1111 1111 1111 | ✅ Success |
| 5105 1051 0510 5100 | ✅ Success |
| 4000 0000 0000 0002 | ❌ Declined |

**For all cards**: CVV = any 3 digits, Expiration = any future date

---

**Congratulations! You now have a working e-commerce store! 🎊**

Next steps:
1. Customize the products
2. Change the styling
3. Add your logo
4. Deploy to production

Happy coding! 🚀