# Firebase Authentication Setup

## Current Status
- ✅ GCP project: `lookie-quantum-intelligence` exists  
- ✅ Firestore database is created and accessible via gcloud
- ✅ User authentication: `david@openqase.com` is active
- ❌ Application Default Credentials: Not set up for Firebase SDK

## The Issue
The Firebase Admin SDK for Go requires **Application Default Credentials (ADC)** to be configured. This is different from the user credentials we have with `gcloud auth login`.

## Solution: Set Up ADC

### Option 1: Interactive Setup (Recommended)
```bash
# This will open a browser window for authentication
gcloud auth application-default login
```

After running this command:
1. Browser opens for authentication
2. Login with your Google account (`david@openqase.com`)
3. Grant permissions for Application Default Credentials
4. Credentials are stored in `~/.config/gcloud/application_default_credentials.json`

### Option 2: Non-Interactive (If browser isn't available)
```bash
# Generate a command for browser-based auth
gcloud auth application-default login --no-browser
# Follow the provided instructions
```

## Test Authentication
After setting up ADC, test the connection:
```bash
go run cmd/test-firestore/main.go
```

Expected output:
```
🔥 Testing Firestore connection to project: lookie-quantum-intelligence
✅ Connected to Firestore successfully!
📝 Creating test company document...
✅ Test company 'Test Quantum Corp' created successfully!
...
🎉 All Firestore tests passed!
```

## Why This is Needed
- **gcloud commands** use user credentials (which we have)
- **Firebase Admin SDK** uses Application Default Credentials (which we need to set up)
- **ADC** allows your Go application to authenticate with Google Cloud services

## Next Steps After Authentication
1. ✅ Test Firestore connection
2. ✅ Migrate SQLite data to Firestore  
3. ✅ Update scraper to use Firestore
4. ✅ Continue with Epic 2 implementation

## Security Note
ADC credentials are stored locally and are scoped to your user account. They're safe for development use and don't require managing service account keys.