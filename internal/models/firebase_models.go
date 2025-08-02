package models

import (
	"time"
)

// FirestoreCompany represents a quantum computing company being monitored in Firestore
type FirestoreCompany struct {
	ID                  string    `firestore:"id" json:"id"`
	Name                string    `firestore:"name" json:"name"`
	Domain              string    `firestore:"domain" json:"domain"`
	RSSURL              string    `firestore:"rss_url" json:"rss_url"`
	NewsPageURL         string    `firestore:"news_page_url" json:"news_page_url"`
	QuantumFocus        string    `firestore:"quantum_focus" json:"quantum_focus"`
	MarketSegment       string    `firestore:"market_segment" json:"market_segment"`
	Description         string    `firestore:"description" json:"description"`
	IsActive            bool      `firestore:"is_active" json:"is_active"`
	RobotsTxtCompliant  bool      `firestore:"robots_txt_compliant" json:"robots_txt_compliant"`
	LastScrapedAt       *time.Time `firestore:"last_scraped_at" json:"last_scraped_at"`
	CreatedAt           time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt           time.Time `firestore:"updated_at" json:"updated_at"`
	
	// Embedded stats for fast queries (denormalization)
	Stats CompanyStats `firestore:"stats" json:"stats"`
}

// CompanyStats holds aggregated statistics for a company
type CompanyStats struct {
	TotalArticles      int       `firestore:"total_articles" json:"total_articles"`
	ArticlesThisMonth  int       `firestore:"articles_this_month" json:"articles_this_month"`
	LastCaseStudy      *time.Time `firestore:"last_case_study" json:"last_case_study"`
	LastFundingNews    *time.Time `firestore:"last_funding_news" json:"last_funding_news"`
	AvgConfidenceScore float64   `firestore:"avg_confidence_score" json:"avg_confidence_score"`
}

// FirestoreArticle represents a scraped article with embedded classification in Firestore
type FirestoreArticle struct {
	ID           string    `firestore:"id" json:"id"`
	CompanyID    string    `firestore:"company_id" json:"company_id"`
	CompanyName  string    `firestore:"company_name" json:"company_name"` // Denormalized for fast queries
	URL          string    `firestore:"url" json:"url"`
	Title        string    `firestore:"title" json:"title"`
	Content      string    `firestore:"content" json:"content"`
	Summary      string    `firestore:"summary" json:"summary"`
	PublishedAt  *time.Time `firestore:"published_at" json:"published_at"`
	ScrapedAt    time.Time `firestore:"scraped_at" json:"scraped_at"`
	SourceType   string    `firestore:"source_type" json:"source_type"` // "rss", "web", "manual"
	ContentHash  string    `firestore:"content_hash" json:"content_hash"`
	WordCount    int       `firestore:"word_count" json:"word_count"`
	Language     string    `firestore:"language" json:"language"`
	IsProcessed  bool      `firestore:"is_processed" json:"is_processed"`
	
	// Embedded classification (most recent/primary)
	Classification *Classification `firestore:"classification,omitempty" json:"classification,omitempty"`
	
	// Embedded entities (top 5 for quick access)
	TopEntities []Entity `firestore:"top_entities" json:"top_entities"`
}

// Classification represents AI analysis of an article
type Classification struct {
	ID              string    `firestore:"id" json:"id"`
	ArticleID       string    `firestore:"article_id" json:"article_id"`
	Category        string    `firestore:"category" json:"category"` // "case_study", "product_news", "research", "partnership", "funding", "other"
	Subcategory     string    `firestore:"subcategory" json:"subcategory"` // "quantum_advantage", "hardware_breakthrough", "software_release"
	ConfidenceScore float64   `firestore:"confidence_score" json:"confidence_score"`
	AIModel         string    `firestore:"ai_model" json:"ai_model"`
	Reasoning       string    `firestore:"reasoning" json:"reasoning"`
	IsValidated     bool      `firestore:"is_validated" json:"is_validated"`
	ValidatedBy     string    `firestore:"validated_by" json:"validated_by"`
	CreatedAt       time.Time `firestore:"created_at" json:"created_at"`
}

// Entity represents extracted entities from article content
type Entity struct {
	ID              string  `firestore:"id" json:"id"`
	ArticleID       string  `firestore:"article_id" json:"article_id"`
	EntityType      string  `firestore:"entity_type" json:"entity_type"` // "technology", "person", "company", "location", "metric"
	EntityValue     string  `firestore:"entity_value" json:"entity_value"`
	ConfidenceScore float64 `firestore:"confidence_score" json:"confidence_score"`
	StartPosition   int     `firestore:"start_position" json:"start_position"`
	EndPosition     int     `firestore:"end_position" json:"end_position"`
	CreatedAt       time.Time `firestore:"created_at" json:"created_at"`
}

// ScrapingRun tracks operational metadata for scraping operations
type ScrapingRun struct {
	ID                string    `firestore:"id" json:"id"`
	CompanyID         string    `firestore:"company_id" json:"company_id"`
	CompanyName       string    `firestore:"company_name" json:"company_name"` // Denormalized
	RunType           string    `firestore:"run_type" json:"run_type"` // "scheduled", "manual", "retry"
	Status            string    `firestore:"status" json:"status"` // "running", "completed", "failed"
	ArticlesFound     int       `firestore:"articles_found" json:"articles_found"`
	ArticlesNew       int       `firestore:"articles_new" json:"articles_new"`
	ArticlesProcessed int       `firestore:"articles_processed" json:"articles_processed"`
	ErrorMessage      string    `firestore:"error_message" json:"error_message"`
	StartedAt         time.Time `firestore:"started_at" json:"started_at"`
	CompletedAt       *time.Time `firestore:"completed_at" json:"completed_at"`
	DurationSeconds   int       `firestore:"duration_seconds" json:"duration_seconds"`
}

// Notification represents a notification to be sent
type Notification struct {
	ID               string    `firestore:"id" json:"id"`
	ArticleID        string    `firestore:"article_id" json:"article_id"`
	ArticleTitle     string    `firestore:"article_title" json:"article_title"` // Denormalized
	CompanyName      string    `firestore:"company_name" json:"company_name"` // Denormalized
	NotificationType string    `firestore:"notification_type" json:"notification_type"` // "email", "webhook"
	Recipient        string    `firestore:"recipient" json:"recipient"`
	Subject          string    `firestore:"subject" json:"subject"`
	Body             string    `firestore:"body" json:"body"`
	Status           string    `firestore:"status" json:"status"` // "pending", "sent", "failed"
	SentAt           *time.Time `firestore:"sent_at" json:"sent_at"`
	ErrorMessage     string    `firestore:"error_message" json:"error_message"`
	CreatedAt        time.Time `firestore:"created_at" json:"created_at"`
}

// SystemConfig holds application-wide configuration
type SystemConfig struct {
	ID                string    `firestore:"id" json:"id"`
	ScrapingInterval  string    `firestore:"scraping_interval" json:"scraping_interval"`
	LastSystemHealth  time.Time `firestore:"last_system_health" json:"last_system_health"`
	ActiveCompanies   int       `firestore:"active_companies" json:"active_companies"`
	TotalArticles     int       `firestore:"total_articles" json:"total_articles"`
	ProcessingBacklog int       `firestore:"processing_backlog" json:"processing_backlog"`
	UpdatedAt         time.Time `firestore:"updated_at" json:"updated_at"`
}