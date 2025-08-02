# Lookie - Quantum Computing Intelligence Service

Lookie is an automated quantum computing news monitoring and intelligence service that scrapes, classifies, and alerts on important developments in the quantum computing industry.

> **🔥 Now Running on Firebase/Firestore**: This project has been migrated from SQLite to Firebase/Firestore for better scalability, cost-effectiveness, and modern cloud-native architecture. See [Firestore Setup](#firestore-setup) below.

## Features

- **Automated News Scraping**: Monitors RSS feeds and websites from major quantum computing companies
- **AI-Powered Classification**: Uses Google Gemini to classify content by type, importance, and relevance
- **Smart Notifications**: Sends immediate alerts for high-value content like case studies and funding announcements
- **Comprehensive API**: REST API for managing articles, classifications, and system monitoring
- **Health Monitoring**: Built-in system health checks and performance metrics
- **Flexible Deployment**: Can run as server-only, scheduler-only, or full mode

## Monitored Companies

- **Q-CTRL**: RSS feed monitoring
- **Diraq**: HTML scraping with custom selectors
- **PsiQuantum**: HTML scraping with custom selectors  
- **IonQ**: RSS feed monitoring

## Current Firestore Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     FIREBASE PROJECT                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Cloud Run     │  │ Cloud Functions │  │   Firestore     │ │
│  │   (API Server)  │◄─┤   (Background)  │◄─┤   (Database)    │ │
│  │                 │  │                 │  │                 │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│           │                     │                     │         │
│           ▼                     ▼                     ▼         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │ Firebase Auth   │  │ Cloud Scheduler │  │ Vertex AI       │ │
│  │ (Admin SDK)     │  │ (Triggers)      │  │ (Gemini)        │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.23 or higher
- Google Cloud Platform account
- Firebase project with Firestore enabled
- Google Gemini API key
- Gmail account for notifications (optional)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/dryan/lookie.git
   cd lookie
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up Firebase authentication**:
   ```bash
   # Option 1: Use Application Default Credentials (recommended)
   gcloud auth application-default login
   
   # Option 2: Set environment variable for service account key
   export GOOGLE_APPLICATION_CREDENTIALS=path/to/service-account-key.json
   ```

4. **Set up configuration**:
   ```bash
   cp config.example.yaml config.yaml
   # Edit config.yaml with your Firebase project settings
   ```

5. **Create environment file**:
   ```bash
   cp .env.example .env
   # Add your API keys and credentials
   ```

6. **Build the application**:
   ```bash
   go build -o bin/lookie cmd/lookie/main.go
   ```

7. **Run the service**:
   ```bash
   ./bin/lookie
   ```

## Configuration

### config.yaml

```yaml
firestore:
  project_id: "your-firebase-project-id"
  credentials_file: "./config/service-account-key.json"  # Optional if using ADC

server:
  port: 8080
  host: "0.0.0.0"

ai:
  gemini_api_key: "${GEMINI_API_KEY}"
  max_tokens: 1000
  temperature: 0.1

notifications:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  smtp_username: "${SMTP_USERNAME}"
  smtp_password_env: "SMTP_PASSWORD"
  from_email: "your-email@gmail.com"
  to_email: "alerts@yourcompany.com"

scraping:
  user_agent: "Lookie/1.0 (+https://yoursite.com/lookie)"
  rate_limit_delay: "2s"
  request_timeout: "30s"

logging:
  level: "info"
  format: "json"
```

### Environment Variables

Create a `.env` file:

```bash
# Required
GEMINI_API_KEY=your_gemini_api_key_here

# Optional - Firebase service account (if not using ADC)
GOOGLE_APPLICATION_CREDENTIALS=./config/service-account-key.json

# Optional - for email notifications
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your_app_password
```

## Usage

### Running Different Modes

**Full mode** (default - runs both HTTP server and scheduler):
```bash
./bin/lookie
```

**Server only** (API endpoints only, no automated scraping):
```bash
./bin/lookie -server-only
```

**Scheduler only** (background scraping, no HTTP server):
```bash
./bin/lookie -scheduler-only
```

### API Endpoints

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Manual Scraping
```bash
# Scrape all companies
curl -X POST http://localhost:8080/scrape
```

#### Get Companies
```bash
# Get all companies
curl http://localhost:8080/companies
```

#### Articles
```bash
# Get all articles
curl http://localhost:8080/api/v1/articles

# Get specific article
curl http://localhost:8080/api/v1/articles/123

# Search articles
curl "http://localhost:8080/api/v1/articles/search?q=quantum&category=case_study"
```

#### Classifications
```bash
# Classify unprocessed articles
curl -X POST http://localhost:8080/api/v1/classify/batch

# Get classification stats
curl http://localhost:8080/api/v1/classifications/stats
```

#### System Monitoring
```bash
# System health
curl http://localhost:8080/api/v1/system/health

# Performance metrics
curl http://localhost:8080/api/v1/system/metrics

# Service dependencies
curl http://localhost:8080/api/v1/system/dependencies
```

## Development

### Project Structure

```
lookie/
├── cmd/
│   ├── lookie/           # Main application
│   └── migrate/          # Database migrations
├── internal/
│   ├── ai/              # AI classification service
│   ├── models/          # Data models
│   ├── monitoring/      # Health monitoring
│   ├── notifications/   # Email notifications
│   ├── scrapers/        # Web scraping
│   ├── scheduler/       # Task scheduling
│   ├── server/          # HTTP server
│   └── storage/         # Database layer
├── pkg/
│   └── config/          # Configuration management
├── migrations/          # SQL migration files
├── config.yaml          # Configuration file
└── README.md
```

### Building

```bash
# Development build
go build -o bin/lookie cmd/lookie/main.go

# Production build with version info
go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.GitHash=$(git rev-parse HEAD)" -o bin/lookie cmd/lookie/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/scrapers/
```

### Database Migrations

Create new migration:
```bash
# This creates timestamped up/down migration files
touch migrations/$(date +%Y%m%d_%H%M%S)_your_migration_name.up.sql
touch migrations/$(date +%Y%m%d_%H%M%S)_your_migration_name.down.sql
```

Run migrations:
```bash
go run cmd/migrate/main.go
```

## Deployment

### Local Development

1. Follow the Quick Start guide above
2. Use `config.yaml` for local configuration
3. Run with `go run cmd/lookie/main.go`

### Production (Linux Server)

1. **Build for Linux**:
   ```bash
   GOOS=linux GOARCH=amd64 go build -o lookie-linux cmd/lookie/main.go
   ```

2. **Deploy files**:
   ```bash
   scp lookie-linux user@server:/opt/lookie/
   scp config.yaml user@server:/opt/lookie/
   scp -r migrations/ user@server:/opt/lookie/
   ```

3. **Create systemd service** (`/etc/systemd/system/lookie.service`):
   ```ini
   [Unit]
   Description=Lookie Quantum Intelligence Service
   After=network.target

   [Service]
   Type=simple
   User=lookie
   WorkingDirectory=/opt/lookie
   ExecStart=/opt/lookie/lookie-linux
   Restart=always
   RestartSec=10
   Environment=GEMINI_API_KEY=your_key_here
   Environment=SMTP_PASSWORD=your_password_here

   [Install]
   WantedBy=multi-user.target
   ```

4. **Start service**:
   ```bash
   sudo systemctl enable lookie
   sudo systemctl start lookie
   sudo systemctl status lookie
   ```

### GCP Cloud Run

1. **Create Dockerfile**:
   ```dockerfile
   FROM golang:1.21-alpine AS builder
   WORKDIR /app
   COPY . .
   RUN go build -o lookie cmd/lookie/main.go

   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=builder /app/lookie .
   COPY --from=builder /app/config.yaml .
   COPY --from=builder /app/migrations ./migrations/
   CMD ["./lookie", "-server-only"]
   ```

2. **Deploy to Cloud Run**:
   ```bash
   gcloud run deploy lookie --source . --platform managed --region us-central1
   ```

## Monitoring

### Logs

The application uses structured JSON logging. Key log fields:
- `level`: Log level (info, warn, error)
- `timestamp`: ISO 8601 timestamp
- `component`: Which service generated the log
- `message`: Human-readable message

### Health Checks

Monitor the `/health` endpoint for basic health status:
```bash
curl http://localhost:8080/health
```

For detailed system health:
```bash
curl http://localhost:8080/api/v1/system/health
```

### Performance Metrics

Access performance metrics at:
```bash
curl http://localhost:8080/api/v1/system/metrics
```

## Troubleshooting

### Common Issues

1. **Database initialization failed**:
   ```bash
   # Run migrations manually
   go run cmd/migrate/main.go
   ```

2. **Gemini API quota exceeded**:
   - Check your Google Cloud console for API usage
   - Verify your API key has Gemini API access enabled

3. **Email notifications not working**:
   - Verify Gmail App Password is correctly set
   - Check SMTP configuration in config.yaml
   - Ensure "Less secure app access" is enabled (if not using App Password)

4. **Scraping failures**:
   - Check robots.txt compliance
   - Verify company URLs are accessible
   - Review rate limiting settings

### Debug Mode

Run with debug logging:
```bash
LOOKIE_LOG_LEVEL=debug ./bin/lookie
```

### Checking Dependencies

Verify external service connectivity:
```bash
curl http://localhost:8080/api/v1/system/dependencies
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- Create an issue on GitHub for bugs or feature requests
- Check the logs for detailed error information
- Use the `/api/v1/system/health` endpoint for system diagnostics

---

## Firestore Setup

### 🔥 **Firebase/Firestore Integration**

Lookie now runs entirely on Firebase/Firestore, providing:
- ✅ **Serverless scaling** - Automatically handles traffic spikes
- ✅ **Real-time data** - Live updates across all clients
- ✅ **Cost-effective** - Pay only for what you use
- ✅ **Modern AI integration** - Native Google Cloud AI support
- ✅ **No database maintenance** - Fully managed by Google

### 🚀 **Quick Firebase Setup**

1. **Create Firebase Project**:
   ```bash
   # Go to https://console.firebase.google.com/
   # Click "Create a project"
   # Choose project name (e.g., "lookie-quantum-intelligence")
   # Enable Google Analytics (optional)
   ```

2. **Enable Firestore**:
   ```bash
   # In Firebase Console > Build > Firestore Database
   # Click "Create database"
   # Choose "Start in production mode"
   # Select location (e.g., us-central1)
   ```

3. **Set up Authentication**:
   ```bash
   # Install Google Cloud CLI
   gcloud auth login
   gcloud config set project your-firebase-project-id
   
   # Set up Application Default Credentials
   gcloud auth application-default login
   ```

4. **Initialize with Sample Data**:
   ```bash
   # The app will create initial company data automatically
   # Or import existing data using the migration tools
   go run cmd/migrate-firebase/simple.go
   ```

### 🏗️ **New Firestore Data Model**

**Collections Structure:**
```
/companies/{companyId}              # Company documents
/articles/{articleId}               # Article documents  
/scraping_runs/{runId}              # Operational tracking
/system/{document}                  # App configuration

# Subcollections (nested under articles)
/articles/{articleId}/classifications/{classificationId}
/articles/{articleId}/entities/{entityId}
/articles/{articleId}/notifications/{notificationId}
```

**Key Design Changes:**
- **Document-based** instead of relational tables
- **Denormalized data** for fast queries (embedded company name in articles)
- **Real-time subscriptions** for live updates
- **Embedded classifications** for single-read access
- **Auto-scaling** collections based on usage

### 📋 **Implementation Epics**

#### **🎯 Epic 1: Firebase Foundation (Week 1)**
- Set up Firebase project & GCP integration
- Configure Firestore database with security rules
- Create development/production environments
- Set up CI/CD pipeline for Cloud Run deployment

**Deliverables:**
- Firebase project configured
- Firestore database initialized
- Basic security rules in place
- Development environment ready

#### **🎯 Epic 2: Data Migration & Models (Week 2)**
- Export existing SQLite data to JSON format
- Design new Firestore document models in Go
- Implement Firebase Admin SDK integration (v4)
- Create data migration utility
- Import seed data (companies + existing articles) to Firestore

**Deliverables:**
- New Go models for Firestore documents
- Firebase Admin SDK integrated
- Migration utility working
- Seed data imported

#### **🎯 Epic 3: Core Services Refactor (Week 3)**
- Replace SQLite storage layer with Firestore SDK
- Refactor scraper service for document-based storage
- Update API endpoints for Firestore queries
- Implement real-time data subscriptions
- Add deduplication and content hashing

**Deliverables:**
- Scraper writes to Firestore
- Basic API endpoints working
- Real-time query capabilities
- Content deduplication working

#### **🎯 Epic 4: AI Classification System (Week 4)**
- Integrate Vertex AI/Gemini for article classification
- Create Cloud Function for automated processing
- Implement entity extraction and sentiment analysis
- Add classification confidence scoring
- Set up automated processing pipeline

**Deliverables:**
- AI classification working
- Cloud Function processing new articles
- Entity extraction implemented
- Confidence scoring system

#### **🎯 Epic 5: Serverless Architecture (Week 5)**
- Deploy API service to Cloud Run
- Create Cloud Functions for background processing
- Set up Cloud Scheduler for periodic scraping
- Configure auto-scaling and monitoring
- Implement health checks and error handling

**Deliverables:**
- Production deployment on Cloud Run
- Scheduled scraping working
- Auto-scaling configured
- Monitoring and alerting set up

#### **🎯 Epic 6: Enhanced Intelligence Features (Week 6)**
- Build real-time dashboard capabilities
- Add advanced search and filtering
- Implement notification system for high-value content
- Create analytics and reporting endpoints
- Performance optimization and caching

**Deliverables:**
- Advanced search working
- Notification system active
- Analytics dashboard
- Performance optimized

### 🔧 **Technical Stack Updates**

**Current (Legacy):**
- Go 1.24 + Gin framework
- SQLite database
- Manual deployment
- Basic Gemini API integration

**Target (Firebase):**
- Go 1.24 + Firebase Admin SDK v4
- Cloud Firestore (NoSQL document database)
- Cloud Run (serverless containers)
- Cloud Functions (event-driven processing)
- Cloud Scheduler (cron jobs)
- Vertex AI (enhanced Gemini integration)

### 💰 **Cost Comparison**

**Current Approach (VPS/Compute Engine):**
- Compute: $20-30/month (always running)
- Storage: $2-5/month
- Gemini API: $5-15/month
- **Total: $27-50/month**

**Firebase Approach:**
- Cloud Run: $0-5/month (pay per request)
- Firestore: $0-3/month (free tier covers initial usage)
- Cloud Functions: $0-2/month 
- Gemini API: $5-15/month
- **Total: $5-25/month**

### 🚀 **Getting Started with Migration**

1. **Review the plan** and provide feedback on priorities
2. **Start Epic 1** - Firebase project setup
3. **Validate approach** after each epic completion
4. **Iterate and improve** based on learnings

**Current Status**: ✅ Research and planning complete, ready to begin Epic 1

---