# Epic 3: Core Services Refactor - SQLite to Firestore Migration

## Overview

This epic focuses on migrating the core scraper service from SQLite to Firestore. Our research revealed that the current SQLite implementation is **incomplete** - articles are generated but not persisted to the database.

## Current State Analysis

### ✅ What Works
- RSS feed parsing with `gofeed` library
- Company data retrieval from database
- Article model generation with content hashing
- HTTP service with Gin framework
- Configuration management with Viper

### ❌ Critical Issues Discovered
- **Missing Article CRUD Operations**: No functions to save articles to SQLite
- **No Deduplication Logic**: Articles are hashed but not checked against existing data
- **Incomplete Data Flow**: Articles are parsed but never persisted
- **No Processing Queue**: No mechanism to track article processing status

### 🏗️ Architecture Gap

**Current (Broken) Flow:**
```
RSS Feed → Parse → Generate Article Model → [STOPS HERE - NO PERSISTENCE]
```

**Required Complete Flow:**
```
RSS Feed → Parse → Check Duplicates → Store Article → Queue for Processing → AI Classification → Update Status
```

## Migration Strategy

We'll use a **3-Phase Approach** to ensure stability and proper testing:

### Phase 1: Complete SQLite Implementation (Foundation)
**Goal**: Fix the broken SQLite implementation before migrating

**Tasks:**
1. **Complete Article CRUD Operations** in `internal/storage/database.go`:
   ```go
   func (db *Database) CreateArticle(article *models.Article) error
   func (db *Database) GetArticlesByCompany(companyID int, limit int) ([]models.Article, error)
   func (db *Database) CheckArticleExists(contentHash string) (bool, error)
   func (db *Database) UpdateArticleStatus(articleID int, isProcessed bool) error
   func (db *Database) GetUnprocessedArticles(limit int) ([]models.Article, error)
   ```

2. **Fix Scraper Service** in `internal/scrapers/scraper.go`:
   - Add article persistence to `ScrapeCompany` method
   - Implement deduplication checking before insert
   - Add proper error handling and logging

3. **Test Complete SQLite Flow**:
   - Verify articles are saved to database
   - Test deduplication prevents duplicates
   - Validate processing status updates

### Phase 2: Storage Interface Abstraction
**Goal**: Create clean abstraction layer supporting both storage backends

**Tasks:**
1. **Create Storage Interface** in `internal/storage/interface.go`:
   ```go
   type Storage interface {
       // Company operations
       GetAllCompanies() ([]Company, error)
       GetCompanyByID(id string) (*Company, error)
       
       // Article operations  
       CreateArticle(article *Article) error
       GetArticlesByCompany(companyID string, limit int) ([]Article, error)
       CheckArticleExists(contentHash string) (bool, error)
       UpdateArticleStatus(articleID string, isProcessed bool) error
       GetUnprocessedArticles(limit int) ([]Article, error)
       
       // Classification operations
       UpdateArticleClassification(articleID string, classification *Classification) error
   }
   ```

2. **Implement Interface for SQLite** (`sqlite_storage.go`)
3. **Implement Interface for Firestore** (`firestore_storage.go`)
4. **Update Scraper to Use Interface** (dependency injection pattern)

### Phase 3: Firestore Migration & Enhanced Features  
**Goal**: Complete migration with enhanced capabilities

**Tasks:**
1. **Data Model Mapping**:
   - Convert integer IDs to string document IDs
   - Implement denormalization for performance
   - Add embedded company metadata to articles

2. **Enhanced Scraper Features**:
   - Real-time article subscriptions
   - Better error handling and retries
   - Improved deduplication with Firestore queries
   - Batch operations for performance

3. **Migration Tooling**:
   - Dual-write capability during transition
   - Data validation and integrity checks
   - Rollback procedures if needed

## Technical Specifications

### Data Model Changes

**SQLite Article Model (Current):**
```go
type Article struct {
    ID          int       `db:"id" json:"id"`
    CompanyID   int       `db:"company_id" json:"company_id"`
    URL         string    `db:"url" json:"url"`
    Title       string    `db:"title" json:"title"`
    Content     string    `db:"content" json:"content"`
    PublishedAt time.Time `db:"published_at" json:"published_at"`
    ContentHash string    `db:"content_hash" json:"content_hash"`
    IsProcessed bool      `db:"is_processed" json:"is_processed"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
```

**Firestore Article Model (Target):**
```go
type FirestoreArticle struct {
    ID           string                 `firestore:"id" json:"id"`
    CompanyID    string                 `firestore:"company_id" json:"company_id"`
    CompanyName  string                 `firestore:"company_name" json:"company_name"`
    URL          string                 `firestore:"url" json:"url"`
    Title        string                 `firestore:"title" json:"title"`
    Content      string                 `firestore:"content" json:"content"`
    Summary      string                 `firestore:"summary,omitempty" json:"summary,omitempty"`
    PublishedAt  time.Time              `firestore:"published_at" json:"published_at"`
    ContentHash  string                 `firestore:"content_hash" json:"content_hash"`
    Status       string                 `firestore:"status" json:"status"` // pending, processing, classified
    Classifications []Classification    `firestore:"classifications,omitempty" json:"classifications,omitempty"`
    ConfidenceScore float64            `firestore:"confidence_score,omitempty" json:"confidence_score,omitempty"`
    CreatedAt    time.Time              `firestore:"created_at" json:"created_at"`
    UpdatedAt    time.Time              `firestore:"updated_at" json:"updated_at"`
}
```

### Key Improvements in Firestore Model

1. **Denormalization**: Embedded company name for efficient queries
2. **Enhanced Status**: More granular processing states
3. **Embedded Classifications**: AI results stored directly in article
4. **Content Summary**: AI-generated summaries for better UX
5. **Confidence Scoring**: Quantitative quality metrics

## Risk Mitigation

### Rollback Strategy
- Maintain SQLite implementation as fallback
- Feature flags to switch between storage backends
- Data export capabilities for emergency migration

### Testing Strategy
- Unit tests for all storage operations
- Integration tests with real RSS feeds
- Performance benchmarks comparing SQLite vs Firestore
- End-to-end scraping workflow validation

### Monitoring
- Scraping success/failure rates
- Article processing latency
- Database operation performance
- Error rates and patterns

## Success Criteria

### Phase 1 Complete:
- [ ] All Article CRUD operations implemented in SQLite
- [ ] Scraper successfully saves articles to database
- [ ] Deduplication prevents duplicate articles
- [ ] Processing status updates work correctly

### Phase 2 Complete:
- [ ] Storage interface abstraction implemented  
- [ ] Both SQLite and Firestore backends work with same interface
- [ ] Scraper service uses dependency injection
- [ ] All tests pass with both storage backends

### Phase 3 Complete:
- [ ] Full migration to Firestore completed
- [ ] Enhanced features (real-time, better performance) working
- [ ] All existing functionality preserved
- [ ] Performance improvements measurable
- [ ] Documentation updated

## Timeline Estimate

- **Phase 1**: 2-3 days (Critical foundation work)
- **Phase 2**: 2-3 days (Interface abstraction)
- **Phase 3**: 3-4 days (Migration & enhancements)
- **Total**: ~1.5 weeks

## Dependencies

- Existing Firestore connection (✅ Complete from Epic 2)
- Company data migration (✅ Complete from Epic 2)
- Firebase Admin SDK setup (✅ Complete from Epic 2)

## Next Steps

1. Start with Phase 1: Complete the broken SQLite implementation
2. Create comprehensive tests for article operations
3. Validate complete scraping workflow end-to-end
4. Proceed with interface abstraction once foundation is solid

---

*This document will be updated as we progress through each phase.*