# Lookie - Quantum Computing Intelligence Service

Lookie is a cloud-native quantum computing news monitoring and AI classification service that automatically scrapes RSS feeds from quantum companies, uses Google Gemini AI to classify content, and provides alerts for important developments.

> **🚀 Production Ready**: Lookie is deployed and running on Google Cloud Run with Firestore database and real Gemini AI integration. See [Live Demo](#live-demo) below.

## Features

- **Automated RSS Scraping**: Monitors RSS feeds from quantum computing companies
- **AI-Powered Classification**: Uses Google Gemini AI to classify content by type, importance, and relevance  
- **Cloud-Native Architecture**: Runs on Google Cloud Run with Firestore database
- **Real-time API**: REST API for company management, article retrieval, and scraping control
- **Content Deduplication**: SHA256 hashing prevents duplicate articles
- **Secure Configuration**: API keys stored in Google Secret Manager

## Live Demo

**Production Service**: `https://lookie-727276629029.us-central1.run.app`

Try these endpoints:
- **Health Check**: `GET /health`
- **List Companies**: `GET /companies` 
- **Manual Scrape**: `POST /scrape`
- **Update RSS URL**: `PUT /companies/{id}/rss`

## Current Quantum Companies

- **IonQ** (ionq.com) - Trapped ion quantum computing, public company
- **Q-CTRL** (q-ctrl.com) - AI-powered quantum infrastructure software  
- **PsiQuantum** (psiquantum.com) - Photonic quantum computing, enterprise
- **Diraq** (diraq.com) - Silicon-based quantum computing, startup

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                   GOOGLE CLOUD PLATFORM                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Cloud Run     │  │   Firestore     │  │  Secret Manager │ │
│  │  (Lookie API)   │◄─┤   (Database)    │  │  (API Keys)     │ │
│  │                 │  │                 │  │                 │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│           │                     │                     │         │
│           ▼                     ▼                     ▼         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │ RSS Scraper     │  │   Gemini AI     │  │ Cloud Scheduler │ │
│  │ (gofeed lib)    │  │ (Classification)│  │   (Future)      │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Technology Stack

**Backend:**
- **Language**: Go 1.23
- **Framework**: Gin HTTP framework
- **Database**: Google Firestore (NoSQL document database)
- **AI**: Google Gemini API for content classification
- **RSS Parsing**: gofeed library
- **Configuration**: Viper with YAML + environment variables

**Infrastructure:**
- **Hosting**: Google Cloud Run (serverless containers)
- **Container Registry**: Google Container Registry (GCR)
- **Secrets**: Google Secret Manager
- **Build**: Google Cloud Build
- **Authentication**: Application Default Credentials

## Data Model

### Companies Collection
```json
{
  "id": "ionq",
  "name": "IonQ",
  "domain": "ionq.com", 
  "rss_url": "https://ionq.com/blog/feed",
  "news_page_url": "https://ionq.com/news",
  "quantum_focus": "hardware",
  "market_segment": "public_company",
  "description": "Trapped ion quantum computing",
  "is_active": true,
  "robots_txt_compliant": true,
  "last_scraped_at": "2025-08-03T07:30:00Z",
  "created_at": "2025-07-28T09:03:38Z",
  "updated_at": "2025-08-03T07:30:00Z",
  "stats": {
    "total_articles": 0,
    "articles_this_month": 0,
    "last_case_study": null,
    "last_funding_news": null,
    "avg_confidence_score": 0
  }
}
```

### Articles Collection
```json
{
  "id": "article_ionq_abc12345",
  "company_id": "ionq",
  "company_name": "IonQ",
  "url": "https://ionq.com/blog/quantum-breakthrough",
  "title": "Major Quantum Breakthrough Achieved",
  "content": "Full article content...",
  "summary": "",
  "published_at": "2025-08-03T07:00:00Z",
  "scraped_at": "2025-08-03T07:30:00Z",
  "source_type": "rss",
  "content_hash": "sha256_hash_of_content",
  "word_count": 1500,
  "language": "en",
  "is_processed": false,
  "classification": null,
  "top_entities": []
}
```

## Quick Start

### Option 1: Use Live Service (Recommended)

The service is already running in production! Try these commands:

```bash
# Check service health
curl https://lookie-727276629029.us-central1.run.app/health

# List quantum companies being monitored
curl https://lookie-727276629029.us-central1.run.app/companies

# Trigger manual scraping
curl -X POST https://lookie-727276629029.us-central1.run.app/scrape
```

### Option 2: Local Development

**Prerequisites:**
- Go 1.23+
- Google Cloud Platform account with Firestore enabled
- Google Gemini API key

**Setup:**
```bash
# 1. Clone repository
git clone https://github.com/dryan/lookie.git
cd lookie

# 2. Install dependencies  
go mod download

# 3. Set up Google Cloud authentication
gcloud auth application-default login
gcloud config set project your-gcp-project-id

# 4. Configure application
cp config.example.yaml config.yaml
# Edit config.yaml with your project settings

# 5. Set environment variables
export GEMINI_API_KEY="your-gemini-api-key"
export GOOGLE_CLOUD_PROJECT="your-gcp-project-id"

# 6. Build and run
go build -o bin/lookie cmd/lookie/main.go
./bin/lookie
```

## API Endpoints

### Core Operations

**Health Check**
```bash
GET /health
# Returns: {"status":"healthy"}
```

**List Companies**
```bash  
GET /companies
# Returns: Array of quantum companies being monitored
```

**Manual Scrape**
```bash
POST /scrape  
# Triggers immediate scraping of all active companies
# Returns: {"message":"Scraping completed successfully","status":"completed"}
```

**Update Company RSS URL**
```bash
PUT /companies/:id/rss
Content-Type: application/json

{
  "rss_url": "https://example.com/feed"
}
```

### Example API Usage

```bash
# Check what companies are monitored
curl https://lookie-727276629029.us-central1.run.app/companies | jq '.[].name'

# Update a company's RSS feed 
curl -X PUT https://lookie-727276629029.us-central1.run.app/companies/ionq/rss \
  -H "Content-Type: application/json" \
  -d '{"rss_url": "https://techcrunch.com/feed/"}'

# Trigger scraping
curl -X POST https://lookie-727276629029.us-central1.run.app/scrape
```

## Configuration

### Production Config (config.production.yaml)
```yaml
firestore:
  project_id: "lookie-quantum-intelligence"
  credentials_file: ""  # Uses Application Default Credentials

server:
  port: 8080
  host: "0.0.0.0"

scraping:
  default_delay: "10s"
  max_retries: 3
  timeout: "30s"
  user_agent: "Lookie/1.0 (Quantum News Monitoring)"

ai:
  gemini_api_key: "${GEMINI_API_KEY}"
  max_tokens: 1000
  temperature: 0.1

notifications:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  smtp_username: "${SMTP_USERNAME}"
  smtp_password_env: "SMTP_PASSWORD"
  from_email: "${SMTP_USERNAME}"
  to_email: "alerts@yourcompany.com"

logging:
  level: "info"
  format: "json"
```

### Environment Variables (Cloud Run)
```bash
GOOGLE_CLOUD_PROJECT=lookie-quantum-intelligence
GEMINI_API_KEY=AIza... # (stored in Secret Manager)
SMTP_USERNAME=alerts@yourcompany.com
SMTP_PASSWORD=app_password
```

## How It Works

### RSS Scraping Process
1. **Company Management**: Firestore stores quantum companies with RSS feed URLs
2. **Content Fetching**: Uses `gofeed` library to parse RSS XML from company feeds  
3. **Deduplication**: SHA256 hashing prevents storing duplicate articles
4. **Storage**: Articles saved to Firestore with metadata (company, timestamp, etc.)
5. **AI Classification**: (Ready) Gemini AI will analyze and categorize content
6. **Notifications**: (Future) Alerts for high-value content like funding news

### Current Status

**✅ Working Components:**
- Cloud Run deployment with auto-scaling
- Firestore database with company and article collections
- RSS feed parsing and content extraction
- Content deduplication using SHA256 hashing
- RESTful API for management and monitoring
- Real Gemini API key integration (ready for classification)
- Secure secret management in Google Secret Manager

**⚠️ Current Status:**
- Infrastructure is production-ready and fully functional
- RSS parsing works perfectly with valid feeds (demonstrated with TechCrunch)
- Quantum company RSS feeds are broken/non-existent (data source issue, not code issue)
- Waiting for working RSS feed URLs to test full AI classification pipeline

## Development Phases

### 📋 Phase 1: RSS-Based Monitoring (Current)
**Goal**: Establish baseline functionality with RSS feed monitoring and AI classification

**Scope**:
- RSS feed parsing and content extraction
- AI-powered content classification using Gemini
- Automated scraping and storage in Firestore
- Basic alerting for high-value content (case studies, funding news)
- Web dashboard for viewing classified content

**Status**: 🟡 Infrastructure complete, waiting for valid RSS feed URLs

**Remaining Tasks**:
1. Obtain working RSS feed URLs for quantum companies
2. Test full pipeline: RSS → Firestore → Gemini classification  
3. **EPIC: Build comprehensive web interface** (see below)
4. Add Cloud Scheduler for automated scraping
5. Configure alerting system

### 🚀 Phase 2: Comprehensive Monitoring (Future)
**Goal**: Full quantum industry intelligence with multi-source data collection

**Scope**:
- **Website Monitoring**: Direct HTML scraping of company news pages
- **Social Media Monitoring**: Twitter, LinkedIn, company blogs
- **Research Paper Tracking**: arXiv, IEEE, academic publications  
- **Patent Monitoring**: USPTO, international patent databases
- **Financial Data**: Funding rounds, IPOs, earnings reports
- **Job Posting Analysis**: Hiring trends and talent movement
- **Conference/Event Tracking**: QKD, Q2B, quantum conferences
- **Regulatory News**: Government quantum initiatives, policy changes

**Technical Additions**:
- HTML parsing and content extraction engines
- Social media APIs (Twitter, LinkedIn)
- Academic database integrations
- Real-time web monitoring with change detection
- Advanced entity recognition (people, companies, technologies)
- Trend analysis and predictive insights
- Multi-language content support

**Data Sources**:
- Company websites and news pages
- Social media platforms
- Academic and research databases  
- Patent offices and legal databases
- Financial news and SEC filings
- Job boards and recruiting sites
- Industry publications and trade media

## 🎯 EPIC: Comprehensive Web Interface

**Problem**: Currently Lookie has no user interface - all interactions are via API calls. Users need both a dashboard to consume intelligence and an admin panel to manage the system.

### 📊 User Dashboard (Intelligence Consumption)
**Goal**: Professional quantum intelligence dashboard for daily use

**Core Features**:
- **Article Feed**: Chronological view of all scraped and classified articles
- **Company Insights**: Company-specific intelligence pages with article history
- **Search & Filtering**: Full-text search with filters (date, company, classification, confidence)
- **AI Classifications**: Visual indicators for article types (funding, case studies, research, etc.)
- **Trending Topics**: Word clouds and trending themes in quantum industry
- **Alert Center**: High-priority notifications for important developments

**Advanced Features**:
- **Comparative Analysis**: Side-by-side company intelligence comparison
- **Timeline View**: Company milestone tracking and funding history
- **Export Functions**: PDF reports, CSV data exports, email digests
- **Saved Searches**: Bookmark complex queries and get notifications
- **Analytics Dashboard**: Industry trends, classification statistics, source health

**Technical Stack**:
- **Frontend**: React/Next.js with Tailwind CSS
- **Data Fetching**: REST API integration with our existing endpoints
- **Real-time Updates**: WebSocket connections for live article feeds
- **Charts/Visualizations**: Chart.js or D3.js for trend analysis
- **Authentication**: Firebase Auth for user management

### ⚙️ Admin Panel (System Management)
**Goal**: Complete system administration and maintenance interface

**Company Management**:
- **Company Directory**: Add, edit, delete quantum companies
- **RSS Feed Management**: Update feed URLs, test feed validity, view scraping status
- **Company Profiles**: Edit descriptions, focus areas, market segments
- **Bulk Operations**: Import companies from CSV, bulk RSS updates

**Content Management**:
- **Article Review**: Review and edit AI classifications manually
- **Content Moderation**: Mark articles as relevant/irrelevant, merge duplicates
- **Classification Training**: Provide feedback to improve AI accuracy
- **Bulk Actions**: Delete spam, recategorize articles, export datasets

**System Operations**:
- **Scraping Control**: Manual scraping triggers, scheduling configuration
- **Health Monitoring**: RSS feed status, API health, database performance
- **User Management**: Admin user accounts, permissions, audit logs
- **Configuration**: AI classification settings, notification rules, system parameters

**Analytics & Reporting**:
- **System Metrics**: Scraping success rates, classification accuracy, storage usage
- **Content Analytics**: Articles per company, source reliability, trending topics
- **User Activity**: Dashboard usage, popular searches, export frequency
- **Performance Monitoring**: API response times, error rates, uptime statistics

**Technical Implementation**:
- **Admin Framework**: React Admin or custom React interface
- **Database Access**: Direct Firestore admin operations
- **File Uploads**: Company logos, bulk import functionality
- **Role-Based Access**: Different permission levels for different admin users
- **Audit Logging**: Track all admin actions for security and compliance

### 🚀 Implementation Phases

**Phase 1A: Core User Dashboard (2-3 weeks)**
1. Basic article listing with search and filters
2. Company-specific views
3. AI classification display
4. Responsive design and mobile support

**Phase 1B: Advanced User Features (1-2 weeks)**
1. Trending analysis and charts
2. Export functionality
3. Saved searches and alerts
4. Real-time updates

**Phase 1C: Admin Panel (2-3 weeks)**
1. Company and RSS feed management
2. Content moderation tools
3. System monitoring dashboard
4. User management and permissions

**Phase 1D: Polish & Production (1 week)**
1. Performance optimization
2. Security hardening
3. Documentation and help system
4. Production deployment

### 📱 User Experience Goals
- **Fast**: Dashboard loads in <2 seconds
- **Intuitive**: No training required for basic usage
- **Mobile-Friendly**: Works perfectly on tablets and phones
- **Professional**: Clean, modern design suitable for business use
- **Accessible**: WCAG compliance for screen readers
- **Reliable**: 99.9% uptime with graceful error handling

### 🎨 Design Principles
- **Clean & Modern**: Minimalist design focused on content
- **Data-Dense**: Maximum information in minimum space
- **Actionable**: Every view leads to clear next actions
- **Customizable**: User preferences for layout and notifications
- **Consistent**: Unified design language across all interfaces

This Epic transforms Lookie from an API-only service into a complete, usable quantum intelligence platform.

## Development

### Project Structure

```
lookie/
├── cmd/
│   └── lookie/                    # Main application entry point
├── internal/
│   ├── models/                    # Firestore data models  
│   ├── scrapers/                  # RSS scraping with gofeed
│   └── storage/                   # Firestore database layer
├── pkg/
│   └── config/                    # Configuration management
├── config.production.yaml         # Production configuration
├── config.example.yaml           # Template configuration  
├── Dockerfile                     # Container build definition
├── go.mod & go.sum               # Go dependencies
└── README.md                     # This file
```

### Core Components

**cmd/lookie/main.go**
- HTTP server with Gin framework
- API endpoints for health, companies, scraping
- Firestore database initialization
- Environment variable handling

**internal/scrapers/firestore_scraper.go**  
- RSS feed parsing using gofeed library
- Content deduplication with SHA256 hashing
- Firestore document creation and storage
- Error handling and logging

**internal/storage/firestore.go**
- Firestore client initialization and management
- Company and article data access methods
- Authentication via Application Default Credentials

**pkg/config/config.go**
- YAML configuration loading with Viper
- Environment variable override support
- Structured configuration for all services

### Local Development

```bash
# 1. Build application
go build -o bin/lookie cmd/lookie/main.go

# 2. Run locally (requires GCP authentication)
gcloud auth application-default login
export GEMINI_API_KEY="your-key"
./bin/lookie

# 3. Test endpoints
curl localhost:8080/health
curl localhost:8080/companies
```

### Deployment to Cloud Run

```bash
# 1. Build and push container image
gcloud builds submit --tag gcr.io/lookie-quantum-intelligence/lookie

# 2. Deploy to Cloud Run
gcloud run deploy lookie \
  --image gcr.io/lookie-quantum-intelligence/lookie \
  --platform managed \
  --region us-central1 \
  --set-env-vars "GOOGLE_CLOUD_PROJECT=lookie-quantum-intelligence,GEMINI_API_KEY=your-key"
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