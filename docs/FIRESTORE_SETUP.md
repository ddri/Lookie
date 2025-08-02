# Firestore Setup Guide

## Overview

Lookie has been migrated from SQLite to Google Cloud Firestore for better scalability, real-time capabilities, and cost-effectiveness. This guide walks you through setting up Firestore for development and production.

## Prerequisites

- Google Cloud Platform account
- `gcloud` CLI installed and configured
- Go 1.23+ 
- Firebase CLI (optional, for advanced features)

## Step-by-Step Setup

### 1. Create Firebase Project

1. **Go to Firebase Console**: https://console.firebase.google.com/
2. **Click "Create a project"**
3. **Choose project name**: e.g., `lookie-quantum-intelligence`
4. **Enable Google Analytics**: Optional, recommended for production
5. **Wait for project creation**: Usually takes 1-2 minutes

### 2. Enable Firestore Database

1. **Navigate to Firestore**: Build > Firestore Database
2. **Click "Create database"**
3. **Choose security rules**:
   - **Production mode**: Recommended (we'll configure rules later)
   - **Test mode**: Only for development (open to all)
4. **Select location**: Choose closest to your users
   - `us-central1` (Iowa) - Good for US/global
   - `europe-west1` (Belgium) - Good for Europe
   - `asia-northeast1` (Tokyo) - Good for Asia
5. **Wait for database creation**: Usually takes 2-3 minutes

### 3. Configure Authentication

#### Option A: Application Default Credentials (Recommended)

```bash
# Install Google Cloud CLI if not already installed
# https://cloud.google.com/sdk/docs/install

# Authenticate with your Google account
gcloud auth login

# Set your project
gcloud config set project your-firebase-project-id

# Set up Application Default Credentials
gcloud auth application-default login
```

This will open a browser window for authentication and store credentials locally.

#### Option B: Service Account Key (Production)

1. **Create Service Account**:
   - Go to Google Cloud Console > IAM & Admin > Service Accounts
   - Click "Create Service Account"
   - Name: `lookie-service-account`
   - Description: `Service account for Lookie app`

2. **Assign Roles**:
   - `Cloud Datastore User` (for Firestore access)
   - `Firebase Admin` (for Firebase features)

3. **Generate Key**:
   - Click on the service account
   - Go to "Keys" tab
   - Click "Add Key" > "Create new key"
   - Choose JSON format
   - Download and save securely

4. **Set Environment Variable**:
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
   ```

### 4. Configure Firestore Security Rules

Create proper security rules for your Firestore database:

1. **Go to Firestore Console** > Rules
2. **Replace default rules** with:

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    // Allow read access to companies for all authenticated users
    match /companies/{companyId} {
      allow read: if true;  // Public read access for company data
      allow write: if request.auth != null;  // Only authenticated users can write
    }
    
    // Articles collection - authenticated access only
    match /articles/{articleId} {
      allow read, write: if request.auth != null;
    }
    
    // System configuration - admin only
    match /system/{document} {
      allow read: if true;
      allow write: if request.auth != null && 
                     request.auth.token.email in ['admin@yourcompany.com'];
    }
    
    // Scraping runs - service account only
    match /scraping_runs/{runId} {
      allow read, write: if request.auth != null;
    }
  }
}
```

3. **Click "Publish"** to apply the rules

### 5. Set Up Firestore Indexes

For optimal query performance, create these indexes:

1. **Go to Firestore Console** > Indexes
2. **Create Composite Indexes**:

   **Articles by Company and Date**:
   - Collection: `articles`
   - Fields: `company_id` (Ascending), `published_at` (Descending)
   
   **Unprocessed Articles**:
   - Collection: `articles`
   - Fields: `is_processed` (Ascending), `scraped_at` (Ascending)

3. **Or use the provided indexes file**:
   ```bash
   # If you have Firebase CLI installed
   firebase deploy --only firestore:indexes
   ```

### 6. Initialize Sample Data

Run the migration utility to populate initial company data:

```bash
# Navigate to your project directory
cd /path/to/lookie

# Run the Firebase migration
go run cmd/migrate-firebase/simple.go
```

This will create the 4 quantum computing companies:
- Q-CTRL (software, has RSS)
- Diraq (hardware, no RSS)
- PsiQuantum (hardware, no RSS)  
- IonQ (hardware, has RSS)

### 7. Configure Application

Update your `config.yaml`:

```yaml
firestore:
  project_id: "your-firebase-project-id"
  credentials_file: ""  # Leave empty if using ADC

server:
  port: 8080
  host: "localhost"

scraping:
  default_delay: "10s"
  max_retries: 3
  timeout: "30s"

ai:
  provider: "gemini"
  model: "gemini-pro"

logging:
  level: "info"
  format: "json"
```

### 8. Test the Setup

Build and run the application:

```bash
# Build the application
go build -o bin/lookie cmd/lookie/main.go

# Run the application
./bin/lookie
```

You should see:
```
time="2025-08-02T08:00:00+10:00" level=info msg="Starting Lookie with Firestore"
2025/08/02 08:00:00 Connected to Firestore database: your-firebase-project-id
time="2025-08-02T08:00:00+10:00" level=info msg="Starting HTTP server on port 8080"
```

Test the endpoints:
```bash
# Health check
curl http://localhost:8080/health

# Get companies
curl http://localhost:8080/companies

# Trigger scraping
curl -X POST http://localhost:8080/scrape
```

## Firestore Data Model

### Collections Structure

```
/companies/{companyId}              # Company documents
├── id: string                      # Document ID (e.g., "q-ctrl")
├── name: string                    # Company name
├── domain: string                  # Company domain
├── rss_url: string                 # RSS feed URL (nullable)
├── quantum_focus: string           # "hardware", "software", "cloud"
├── market_segment: string          # "startup", "enterprise", "public_company"
├── is_active: boolean              # Whether to scrape this company
└── stats: object                   # Embedded statistics

/articles/{articleId}               # Article documents  
├── id: string                      # Document ID
├── company_id: string              # Reference to company
├── company_name: string            # Denormalized company name
├── url: string                     # Article URL
├── title: string                   # Article title
├── content: string                 # Article content
├── content_hash: string            # SHA256 hash for deduplication
├── published_at: timestamp         # When article was published
├── scraped_at: timestamp           # When we scraped it
├── is_processed: boolean           # Whether AI has processed it
├── classification: object          # AI classification results
└── top_entities: array             # Extracted entities

/system/{document}                  # App configuration
└── main                            # Main system config document

/scraping_runs/{runId}              # Operational tracking
├── company_id: string              # Which company was scraped
├── status: string                  # "running", "completed", "failed"
├── articles_found: number          # How many articles found
├── started_at: timestamp           # When scraping started
└── completed_at: timestamp         # When scraping finished
```

### Key Design Decisions

1. **Document IDs**: Use descriptive strings instead of auto-generated IDs
   - Companies: `q-ctrl`, `diraq`, `psiquantum`, `ionq`
   - Articles: `article_{companyId}_{hash}`

2. **Denormalization**: Embed frequently-accessed data
   - Article documents include `company_name` for fast display
   - Company documents include embedded `stats` object

3. **Content Hashing**: SHA256 hash of URL + title for deduplication
   - Prevents storing duplicate articles
   - Fast duplicate checking with indexed hash field

4. **Real-time Ready**: Structure supports real-time subscriptions
   - Front-end can subscribe to new articles
   - Live updates when scraping completes

## Troubleshooting

### Common Issues

1. **Authentication Failed**
   ```
   Error: failed to create Firestore client: unexpected end of JSON input
   ```
   **Solution**: Run `gcloud auth application-default login`

2. **Permission Denied**
   ```
   Error: rpc error: code = PermissionDenied
   ```
   **Solution**: Check Firestore security rules or user permissions

3. **Project Not Found**
   ```
   Error: project not found
   ```
   **Solution**: Verify project ID in config and `gcloud config`

4. **Firestore Not Enabled**
   ```
   Error: Firestore API has not been used
   ```
   **Solution**: Enable Firestore in Firebase Console

### Useful Commands

```bash
# Check current gcloud configuration
gcloud config list

# Check ADC status
gcloud auth application-default print-access-token

# Test Firestore access
gcloud firestore databases list

# Check project permissions
gcloud projects get-iam-policy your-project-id

# Reset authentication
gcloud auth revoke --all
gcloud auth login
gcloud auth application-default login
```

### Firebase Console URLs

- **Main Console**: https://console.firebase.google.com/
- **Firestore Database**: https://console.firebase.google.com/project/YOUR_PROJECT/firestore
- **Security Rules**: https://console.firebase.google.com/project/YOUR_PROJECT/firestore/rules
- **Indexes**: https://console.firebase.google.com/project/YOUR_PROJECT/firestore/indexes
- **Usage & Billing**: https://console.firebase.google.com/project/YOUR_PROJECT/usage

## Production Considerations

### Security

1. **Use Service Accounts**: For production deployments
2. **Restrict Rules**: Limit write access appropriately
3. **Enable Audit Logs**: Track all database access
4. **Monitor Usage**: Set up billing alerts

### Performance

1. **Index Strategy**: Create indexes for all query patterns
2. **Batch Operations**: Use batch writes for bulk operations
3. **Connection Pooling**: Reuse Firestore client instances
4. **Cache Reads**: Cache frequently-accessed data

### Monitoring

1. **Cloud Logging**: Enable structured logging
2. **Error Reporting**: Set up error alerts
3. **Performance Monitoring**: Track query performance
4. **Usage Monitoring**: Monitor read/write operations

## Cost Optimization

Firestore pricing is based on:
- **Document reads**: $0.06 per 100K reads
- **Document writes**: $0.18 per 100K writes
- **Document deletes**: $0.02 per 100K deletes
- **Storage**: $0.18 per GB/month

**Tips to minimize costs**:
1. Use efficient queries (avoid scanning large collections)
2. Implement proper indexing
3. Cache frequently-accessed data
4. Use batch operations for bulk writes
5. Monitor usage in Firebase Console

## Next Steps

Once Firestore is set up and working:

1. **Deploy to Cloud Run**: Scale automatically with traffic
2. **Set up CI/CD**: Automate deployments
3. **Add monitoring**: Track performance and errors
4. **Implement caching**: Improve response times
5. **Add real-time features**: Live dashboards and notifications

For deployment instructions, see `docs/DEPLOYMENT.md`.
For API documentation, see `docs/API.md`.