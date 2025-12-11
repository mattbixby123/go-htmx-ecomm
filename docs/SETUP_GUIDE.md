# Square Payment Setup Guide

This guide will walk you through setting up Square payments for your TechStore application.

## Step 1: Create a Square Developer Account

1. Go to [https://developer.squareup.com/](https://developer.squareup.com/)
2. Click "Sign Up" in the top right
3. Fill in your information and create an account
4. Verify your email address

## Step 2: Create a New Application

1. Once logged in, go to the [Developer Dashboard](https://developer.squareup.com/apps)
2. Click "Create App" or the "+" button
3. Enter your application name (e.g., "TechStore")
4. Click "Create Application"

## Step 3: Get Your Credentials

### Application ID

1. In your app's dashboard, click on "Credentials" in the left sidebar
2. You'll see two sections: **Sandbox** and **Production**
3. Under **Sandbox**, find the **Application ID**
4. Copy this ID - it looks like: `sandbox-sq0idb-xxxxxxxxxxxxxxxxxxxxxx`

### Access Token

1. Still in the **Credentials** section under **Sandbox**
2. Find the **Access Token** section
3. Click "Show" to reveal the token
4. Copy the entire token - it's a long string starting with something like `EAAAl...`

**Important**: Never share this token or commit it to version control!

### Location ID

1. In the left sidebar, click on "Locations"
2. You should see at least one location (likely "Default Test Location" for sandbox)
3. Copy the **Location ID** - it looks like: `LXXXXXXXXXXXXXX`

If you don't see any locations:
- Go to your [Square Dashboard](https://squareup.com/dashboard/)
- Make sure you have at least one location set up
- The location will sync to the developer dashboard

## Step 4: Configure Your Application

1. Open your `.env` file (or create it from `.env.example`)
2. Add your credentials:

```env
SQUARE_APPLICATION_ID=sandbox-sq0idb-your-app-id-here
SQUARE_ACCESS_TOKEN=EAAAl-your-access-token-here
SQUARE_LOCATION_ID=LYOUR-LOCATION-ID-HERE
SQUARE_ENVIRONMENT=sandbox
PORT=8080
```

## Step 5: Test Your Integration

### Using Sandbox Test Cards

Square provides test card numbers that work in sandbox mode:

| Card Type | Card Number | Result |
|-----------|-------------|--------|
| Visa | 4111 1111 1111 1111 | Success |
| Visa | 4532 7597 5970 6077 | Success |
| Mastercard | 5105 1051 0510 5100 | Success |
| American Express | 3782 822463 10005 | Success |
| Discover | 6011 1111 1111 1117 | Success |

**For all test cards:**
- CVV: Any 3 digits (4 for Amex)
- Expiration: Any future date
- ZIP: Any 5 digits

### Test Cards for Specific Scenarios

| Card Number | Scenario |
|-------------|----------|
| 4000 0000 0000 0002 | Card declined |
| 4000 0000 0000 0101 | Declined - insufficient funds |
| 4000 0000 0000 0069 | Expired card |
| 4000 0000 0000 0119 | Processing error |

## Step 6: Test the Complete Flow

1. Start your application:
   ```bash
   go run main.go
   ```

2. Open `http://localhost:8080` in your browser

3. Add products to your cart

4. Go to checkout

5. Fill in the billing information:
   - Name: Test User
   - Email: test@example.com
   - Address: 123 Test St
   - City: San Francisco
   - State: CA
   - ZIP: 94103

6. Enter payment information:
   - Card: 4111 1111 1111 1111
   - CVV: 123
   - Expiration: 12/25
   - ZIP: 94103

7. Click "Pay Now"

8. You should be redirected to the order confirmation page

## Step 7: Verify the Payment in Square Dashboard

1. Go to your [Square Dashboard](https://squareup.com/dashboard/)
2. Switch to "Sandbox" mode (toggle in top right)
3. Click on "Transactions" in the left sidebar
4. You should see your test payment listed

## Moving to Production

When you're ready to accept real payments:

### 1. Complete Square Account Verification

1. Go to your [Square Dashboard](https://squareup.com/dashboard/)
2. Complete all required verification steps:
   - Business information
   - Bank account for deposits
   - Identity verification

### 2. Get Production Credentials

1. Go back to [Developer Dashboard](https://developer.squareup.com/apps)
2. Click on your application
3. Under **Credentials**, switch to the **Production** tab
4. Copy your production **Application ID** and **Access Token**
5. Get your production **Location ID** from the Locations page

### 3. Update Your Configuration

```env
SQUARE_APPLICATION_ID=sq0idp-your-production-app-id
SQUARE_ACCESS_TOKEN=EAAAl-your-production-token
SQUARE_LOCATION_ID=YOUR-PRODUCTION-LOCATION-ID
SQUARE_ENVIRONMENT=production
```

### 4. Update Square.js in checkout.html

Change the script source from sandbox to production:

```html
<!-- FROM: -->
<script type="text/javascript" src="https://sandbox.web.squarecdn.com/v1/square.js"></script>

<!-- TO: -->
<script type="text/javascript" src="https://web.squarecdn.com/v1/square.js"></script>
```

### 5. Enable HTTPS

Square **requires HTTPS** for production. Options:

**Option A: Use a reverse proxy (Nginx)**
```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
    }
}
```

**Option B: Use Caddy (automatic HTTPS)**
```
yourdomain.com {
    reverse_proxy localhost:8080
}
```

**Option C: Deploy to a platform with HTTPS**
- Heroku
- Railway
- Fly.io
- DigitalOcean App Platform

## Troubleshooting

### "Failed to initialize payment form"

**Problem**: Square Web Payments SDK not loading

**Solutions**:
- Check that your Application ID is correct
- Verify you're using the right environment (sandbox/production)
- Check browser console for JavaScript errors
- Ensure you have a stable internet connection

### "Payment failed: Invalid credentials"

**Problem**: Access token or location ID incorrect

**Solutions**:
- Double-check your `.env` file
- Make sure there are no extra spaces in the credentials
- Verify the environment (sandbox vs production) matches your tokens
- Regenerate your access token if needed

### "Location ID not found"

**Problem**: The location doesn't exist or doesn't match the token

**Solutions**:
- Go to Square Dashboard → Locations
- Verify the location ID exists
- Make sure the location is active
- Check that the location belongs to the account that generated the token

### Payments work in sandbox but fail in production

**Solutions**:
- Verify all production credentials are correct
- Ensure HTTPS is enabled
- Check that Square account verification is complete
- Review Square's production access requirements

### Card is declined

**In Sandbox**:
- Use official Square test card numbers
- Check that you're entering the card number correctly

**In Production**:
- This is likely a real card decline
- Ask the customer to check with their bank
- Try a different payment method

## Security Best Practices

1. **Never commit credentials** - Always use `.env` and add it to `.gitignore`

2. **Use environment variables** - Don't hardcode sensitive data

3. **Implement HTTPS** - Required for production, recommended for sandbox

4. **Validate server-side** - Always verify payments on the server

5. **Log errors** - But never log sensitive card data

6. **Rate limit** - Prevent abuse of your payment endpoint

7. **Use idempotency keys** - Already implemented with UUID in the code

## Additional Resources

- [Square Developer Documentation](https://developer.squareup.com/docs)
- [Web Payments SDK Guide](https://developer.squareup.com/docs/web-payments/overview)
- [Square API Reference](https://developer.squareup.com/reference/square)
- [Square Community Forums](https://developer.squareup.com/forums)
- [Square Status Page](https://status.squareup.com/)

## Support

If you encounter issues:

1. Check the [Square Developer Forums](https://developer.squareup.com/forums)
2. Review [Square's documentation](https://developer.squareup.com/docs)
3. Contact Square Developer Support (if you have a Square Developer account)
4. Open an issue in this repository

## Next Steps

Once you have payments working:

1. Test all error scenarios
2. Implement proper error logging
3. Add email notifications for orders
4. Set up webhook handling for payment events
5. Implement refund functionality
6. Add recurring payment support (if needed)
7. Consider PCI compliance requirements for your specific use case

---

Happy coding! 🚀