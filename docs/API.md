# Lookie API Documentation

## Overview

Lookie provides a RESTful API for quantum computing intelligence operations. The API is built on Firebase/Firestore and provides endpoints for managing companies, triggering scraping operations, and monitoring system health.

**Base URL**: `http://localhost:8080` (development)  
**Production URL**: TBD (when deployed to Cloud Run)

## Authentication

Currently, the API uses Firebase Application Default Credentials for backend operations. No user authentication is required for public endpoints.

**Future**: User authentication will be added using Firebase Auth for private endpoints.

## Endpoints

### Health Check

#### `GET /health`

Simple health check endpoint to verify the service is running.

**Response**:
```json
{
  "status": "healthy"
}
```

**Status Codes**:
- `200`: Service is healthy

**Example**:
```bash
curl http://localhost:8080/health
```

---

### Companies

#### `GET /companies`

Retrieves all active quantum computing companies from Firestore.

**Response**:
```json
[
  {
    "id": "q-ctrl",
    "name": "Q-CTRL",
    "domain": "q-ctrl.com", 
    "rss_url": "https://q-ctrl.com/blog/feed",
    "news_page_url": "https://q-ctrl.com/blog",
    "quantum_focus": "software",
    "market_segment": "enterprise",
    "description": "AI-powered quantum infrastructure software",
    "is_active": true,
    "robots_txt_compliant": true,
    "created_at": "2025-08-01T12:00:00Z",
    "updated_at": "2025-08-01T12:00:00Z",
    "stats": {
      "total_articles": 0,
      "articles_this_month": 0,
      "avg_confidence_score": 0.0
    }
  },
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
    "created_at": "2025-08-01T12:00:00Z",
    "updated_at": "2025-08-01T12:00:00Z",
    "stats": {
      "total_articles": 0,
      "articles_this_month": 0,
      "avg_confidence_score": 0.0
    }
  }
]
```

**Status Codes**:
- `200`: Success
- `500`: Internal server error

**Example**:
```bash
curl http://localhost:8080/companies
```

---

### Scraping

#### `POST /scrape`

Triggers manual scraping of all active companies with RSS feeds.

**Request**: No body required

**Response**:
```json
{
  "message": "Scraping completed successfully",
  "status": "completed"
}
```

**Status Codes**:
- `200`: Scraping completed successfully
- `500`: Scraping failed (check logs for details)

**Example**:
```bash
curl -X POST http://localhost:8080/scrape
```

**Notes**:
- Only companies with valid RSS feeds will be scraped
- Articles are automatically deduplicated using content hashes
- New articles are stored in Firestore for later processing
- The operation may take 10-30 seconds depending on RSS feed sizes

---

## Data Models

### Company Model

```json
{
  "id": "string",                    // Unique company identifier
  "name": "string",                  // Company display name
  "domain": "string",                // Company domain (e.g., "q-ctrl.com")
  "rss_url": "string",              // RSS feed URL (nullable)
  "news_page_url": "string",        // Company news page URL
  "quantum_focus": "string",        // "hardware", "software", "cloud"
  "market_segment": "string",       // "startup", "enterprise", "public_company"
  "description": "string",          // Company description
  "is_active": "boolean",           // Whether to scrape this company
  "robots_txt_compliant": "boolean", // Whether scraping complies with robots.txt
  "created_at": "timestamp",        // When company was added
  "updated_at": "timestamp",        // Last modification time
  "stats": {
    "total_articles": "number",         // Total articles scraped
    "articles_this_month": "number",    // Articles scraped this month
    "avg_confidence_score": "number"    // Average AI confidence score
  }
}
```

### Article Model

```json
{
  "id": "string",                    // Unique article identifier
  "company_id": "string",            // Reference to company
  "company_name": "string",          // Denormalized company name
  "url": "string",                   // Article URL
  "title": "string",                 // Article title
  "content": "string",               // Article content/description
  "summary": "string",               // AI-generated summary (future)
  "published_at": "timestamp",       // When article was published
  "scraped_at": "timestamp",         // When we scraped it
  "source_type": "string",           // "rss", "web", "manual"
  "content_hash": "string",          // SHA256 hash for deduplication
  "word_count": "number",            // Number of words in content
  "language": "string",              // Content language (default: "en")
  "is_processed": "boolean",         // Whether AI has processed it
  "classification": {                // AI classification results (future)
    "category": "string",            // "case_study", "product_news", etc.
    "confidence_score": "number",    // AI confidence (0.0-1.0)
    "reasoning": "string"            // Why AI chose this classification
  },
  "top_entities": [                  // Extracted entities (future)
    {
      "type": "string",              // "technology", "person", "company"
      "value": "string",             // Entity value
      "confidence": "number"         // Extraction confidence
    }
  ]
}
```

## Error Handling

### Error Response Format

All errors return a consistent JSON format:

```json
{
  "error": "Error message describing what went wrong"
}
```

### Common Error Codes

- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error (check logs)
- `503 Service Unavailable`: Firestore connectivity issues

### Example Error Response

```json
{
  "error": "failed to initialize Firestore: unexpected end of JSON input"
}
```

## Rate Limiting

Currently, no rate limiting is implemented. The application is designed for:
- **Manual scraping**: Occasional manual triggers
- **Automated scraping**: Scheduled periodic runs
- **Company queries**: Frequent reads for dashboards

**Future**: Rate limiting will be added based on usage patterns.

## Monitoring and Health

### Application Logs

The application uses structured JSON logging with these fields:

```json
{
  "level": "info|warn|error",
  "time": "2025-08-02T08:00:00Z",
  "msg": "Human-readable message",
  "company": "Q-CTRL",
  "url": "https://example.com",
  "error": "error details if applicable"
}
```

### Health Monitoring

Monitor these aspects for production health:

1. **Health Endpoint**: `GET /health` should return 200
2. **Company Count**: `/companies` should return expected companies
3. **Scraping Success**: `/scrape` should complete without errors
4. **Firestore Connectivity**: Check logs for connection errors

## Future API Endpoints

These endpoints are planned for future releases:

### Articles
- `GET /articles` - List all articles with pagination
- `GET /articles/{id}` - Get specific article
- `GET /articles/search?q={query}` - Search articles
- `GET /companies/{id}/articles` - Get articles for specific company

### AI Classification
- `POST /classify` - Trigger AI classification for unprocessed articles
- `GET /classifications/stats` - Get classification statistics

### System Management
- `GET /system/health` - Detailed system health
- `GET /system/metrics` - Performance metrics
- `GET /system/stats` - Usage statistics

### Real-time Subscriptions
- WebSocket endpoints for real-time article updates
- Server-Sent Events for live scraping status

## Development

### Running Locally

```bash
# Start the server
go run cmd/lookie/main.go

# The API will be available at http://localhost:8080
```

### Testing with curl

```bash
# Test health
curl http://localhost:8080/health

# Get companies
curl http://localhost:8080/companies

# Trigger scraping
curl -X POST http://localhost:8080/scrape

# Test with verbose output
curl -v http://localhost:8080/health
```

### Testing with HTTPie

```bash
# Install HTTPie
pip install httpie

# Test endpoints
http GET localhost:8080/health
http GET localhost:8080/companies
http POST localhost:8080/scrape
```

## Production Deployment

When deployed to Cloud Run, the API will be available at:
- Production URL: `https://lookie-[hash]-uc.a.run.app`
- Custom Domain: TBD

### Cloud Run Configuration

```yaml
# Cloud Run service configuration
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: lookie
spec:
  template:
    spec:
      containers:
      - image: gcr.io/project-id/lookie
        ports:
        - containerPort: 8080
        env:
        - name: FIRESTORE_PROJECT_ID
          value: "lookie-quantum-intelligence"
```

## Security Considerations

### Current Security
- **No user authentication**: All endpoints are public
- **Internal authentication**: Uses Firebase Admin SDK with service account
- **CORS**: Not configured (same-origin only)

### Production Security Recommendations
1. **Add authentication**: Implement Firebase Auth for user sessions
2. **Rate limiting**: Prevent abuse of scraping endpoints  
3. **CORS configuration**: Allow specific origins only
4. **Input validation**: Sanitize all user inputs
5. **API keys**: Require API keys for programmatic access

## Contact

For API questions or issues:
- Create an issue on GitHub
- Check application logs for detailed error information
- Use the `/health` endpoint for basic diagnostics

---

*Last updated: August 2, 2025*  
*API Version: 1.0 (Firestore)*