# Firebase Local Emulator Setup Guide

## Prerequisites ✅
- Node.js v23.11.0 (installed)
- Firebase CLI 13.29.1 (installed)
- Go 1.23+ (updated)
- Java 21+ (required for Firestore emulator)

## Setup Steps

### 1. Firebase Authentication
```bash
# Login to Firebase (will open browser)
firebase login
```

### 2. Initialize Firebase Project  
```bash
# Initialize in project root
firebase init

# Select options:
# ✅ Firestore: Configure security rules and indexes files
# ✅ Emulators: Set up local emulators for Firebase products
# 
# Choose existing project: lookie-quantum-intelligence
# 
# Firestore setup:
# - Rules file: firestore.rules (we already have this)
# - Indexes file: firestore.indexes.json (we already have this)
#
# Emulator setup:
# - Select: Firestore Emulator
# - Accept default ports (8080 for Firestore)
# - Enable Emulator UI (port 4000)
```

### 3. Verify Java Version
```bash
java -version
# Should be Java 21 or higher for latest emulator
```

### 4. Start Emulator
```bash
# Start Firestore emulator
firebase emulators:start --only firestore

# Or start with UI
firebase emulators:start --only firestore,ui
```

### 5. Test Connection
```bash
# Set environment variable for Go app
export FIRESTORE_EMULATOR_HOST=127.0.0.1:8080

# Test our migration
go run cmd/migrate-firebase/simple.go
```

## Expected File Structure After Init

```
lookie/
├── firebase.json          # Firebase project config  
├── firestore.rules        # Security rules (already exists)
├── firestore.indexes.json # Database indexes (already exists)
├── .firebaserc            # Project aliases
└── firebase-debug.log     # Debug logs
```

## Development Workflow

1. **Start emulator**: `firebase emulators:start --only firestore`
2. **Set env var**: `export FIRESTORE_EMULATOR_HOST=127.0.0.1:8080`
3. **Run Go app**: Development app connects to local emulator
4. **View UI**: http://localhost:4000 for Emulator Suite UI
5. **Test & iterate**: Modify code, restart, test

## Next Steps After Setup

1. ✅ Test Firestore models with emulator
2. ✅ Run migration utility against emulator
3. ✅ Validate data structure and queries
4. ✅ Test security rules
5. ✅ Export emulator data for testing
6. 🚀 Deploy to production when ready

## Troubleshooting

**Java Version Issues:**
```bash
# Install Java 21 via Homebrew (macOS)
brew install openjdk@21
```

**Port Conflicts:**
```bash
# Use different ports if needed
firebase emulators:start --only firestore --port 9000
```

**Authentication Issues:**
```bash
# Logout and login again
firebase logout
firebase login
```