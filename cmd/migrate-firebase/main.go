package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/dryan/lookie/internal/models"
	"github.com/dryan/lookie/internal/storage"
)

// SQLiteCompany represents the old SQLite company structure
type SQLiteCompany struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Domain             string `json:"domain"`
	RSSURL             string `json:"rss_url"`
	NewsPageURL        string `json:"news_page_url"`
	QuantumFocus       string `json:"quantum_focus"`
	MarketSegment      string `json:"market_segment"`
	Description        string `json:"description"`
	IsActive           int    `json:"is_active"`
	RobotsTxtCompliant int    `json:"robots_txt_compliant"`
	LastScrapedAt      string `json:"last_scraped_at"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(".env.firebase"); err != nil {
		log.Printf("Warning: Could not load .env.firebase: %v", err)
	}

	ctx := context.Background()

	// Get configuration from environment
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID environment variable is required")
	}

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	
	fmt.Printf("🔥 Starting Firebase migration for project: %s\n", projectID)
	if credentialsFile != "" {
		fmt.Printf("📄 Using credentials file: %s\n", credentialsFile)
	} else {
		fmt.Println("📄 Using Application Default Credentials")
	}

	// Initialize Firestore
	db, err := storage.NewFirestoreDB(ctx, projectID, credentialsFile)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer db.Close()

	fmt.Println("✅ Connected to Firestore")

	// Migrate companies
	if err := migrateCompanies(db); err != nil {
		log.Fatalf("Failed to migrate companies: %v", err)
	}

	// Initialize system configuration
	if err := initializeSystemConfig(db); err != nil {
		log.Fatalf("Failed to initialize system config: %v", err)
	}

	fmt.Println("🎉 Migration completed successfully!")
	fmt.Println("")
	fmt.Println("📊 Migration Summary:")
	fmt.Println("  - Companies migrated from SQLite to Firestore")
	fmt.Println("  - Document IDs converted from integers to strings")
	fmt.Println("  - Company stats initialized")
	fmt.Println("  - System configuration created")
	fmt.Println("")
	fmt.Printf("🔗 View your data: https://console.firebase.google.com/project/%s/firestore\n", projectID)
}

func migrateCompanies(db *storage.FirestoreDB) error {
	fmt.Println("📋 Migrating companies...")

	// Read companies from exported JSON
	data, err := os.ReadFile("migration/data/companies.json")
	if err != nil {
		return fmt.Errorf("failed to read companies.json: %w", err)
	}

	var sqliteCompanies []SQLiteCompany
	if err := json.Unmarshal(data, &sqliteCompanies); err != nil {
		return fmt.Errorf("failed to parse companies JSON: %w", err)
	}

	fmt.Printf("  Found %d companies to migrate\n", len(sqliteCompanies))

	for _, sqliteCompany := range sqliteCompanies {
		// Convert SQLite company to Firestore company
		company := models.FirestoreCompany{
			ID:                 generateCompanyID(sqliteCompany.Name),
			Name:               sqliteCompany.Name,
			Domain:             sqliteCompany.Domain,
			RSSURL:             sqliteCompany.RSSURL,
			NewsPageURL:        sqliteCompany.NewsPageURL,
			QuantumFocus:       sqliteCompany.QuantumFocus,
			MarketSegment:      sqliteCompany.MarketSegment,
			Description:        sqliteCompany.Description,
			IsActive:           sqliteCompany.IsActive == 1,
			RobotsTxtCompliant: sqliteCompany.RobotsTxtCompliant == 1,
			CreatedAt:          parseTime(sqliteCompany.CreatedAt),
			UpdatedAt:          parseTime(sqliteCompany.UpdatedAt),
			Stats: models.CompanyStats{
				TotalArticles:      0,
				ArticlesThisMonth:  0,
				LastCaseStudy:      nil,
				LastFundingNews:    nil,
				AvgConfidenceScore: 0.0,
			},
		}

		// Parse last scraped at if present
		if sqliteCompany.LastScrapedAt != "" && sqliteCompany.LastScrapedAt != "null" {
			lastScraped := parseTime(sqliteCompany.LastScrapedAt)
			company.LastScrapedAt = &lastScraped
		}

		// Create company in Firestore
		if err := db.CreateCompany(&company); err != nil {
			return fmt.Errorf("failed to create company %s: %w", company.Name, err)
		}

		fmt.Printf("  ✅ Migrated: %s (ID: %s)\n", company.Name, company.ID)
	}

	fmt.Printf("✅ Successfully migrated %d companies\n", len(sqliteCompanies))
	return nil
}

func initializeSystemConfig(db *storage.FirestoreDB) error {
	fmt.Println("⚙️ Initializing system configuration...")

	config := models.SystemConfig{
		ID:                "main",
		ScrapingInterval:  "2h",
		LastSystemHealth:  time.Now(),
		ActiveCompanies:   4, // Based on our seed data
		TotalArticles:     0,
		ProcessingBacklog: 0,
		UpdatedAt:         time.Now(),
	}

	// Create system config document
	collections := db.Collections()
	_, err := collections.System.Doc(config.ID).Set(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create system config: %w", err)
	}

	fmt.Println("✅ System configuration initialized")
	return nil
}

// generateCompanyID creates a URL-friendly ID from company name
func generateCompanyID(name string) string {
	// Convert to lowercase and replace spaces with hyphens
	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, ".", "")
	return id
}

// parseTime converts SQLite datetime strings to Go time.Time
func parseTime(timeStr string) time.Time {
	if timeStr == "" || timeStr == "null" {
		return time.Now()
	}

	// Try different time formats that SQLite might use
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t
		}
	}

	// If all formats fail, return current time
	log.Printf("Warning: Could not parse time '%s', using current time", timeStr)
	return time.Now()
}