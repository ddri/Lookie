package scrapers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"

	"github.com/dryan/lookie/internal/models"
	"github.com/dryan/lookie/internal/storage"
)

type FirestoreScraperService struct {
	db     *storage.FirestoreDB
	logger *logrus.Logger
	client *http.Client
}

func NewFirestoreScraperService(db *storage.FirestoreDB, logger *logrus.Logger) *FirestoreScraperService {
	return &FirestoreScraperService{
		db:     db,
		logger: logger,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *FirestoreScraperService) ScrapeCompany(company models.FirestoreCompany) (int, error) {
	if company.RSSURL == "" {
		return 0, fmt.Errorf("no RSS feed URL for company %s", company.Name)
	}

	s.logger.WithFields(logrus.Fields{
		"company": company.Name,
		"url":     company.RSSURL,
	}).Info("Scraping RSS feed")

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(company.RSSURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	newArticleCount := 0
	duplicateCount := 0
	
	for _, item := range feed.Items {
		content := item.Description
		if content == "" {
			content = item.Content
		}

		// Create content hash for deduplication
		hash := sha256.Sum256([]byte(item.Link + item.Title))
		contentHash := fmt.Sprintf("%x", hash)

		// Check if article already exists
		exists, err := s.db.CheckArticleExists(contentHash)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"company": company.Name,
				"title":   item.Title,
				"error":   err,
			}).Warn("Failed to check article existence, skipping")
			continue
		}

		if exists {
			duplicateCount++
			s.logger.WithFields(logrus.Fields{
				"company": company.Name,
				"title":   item.Title,
				"hash":    contentHash[:8],
			}).Debug("Article already exists, skipping")
			continue
		}

		var publishedAt *time.Time
		if item.PublishedParsed != nil {
			publishedAt = item.PublishedParsed
		}

		// No need for SQLite model conversion anymore

		// Create Firestore article directly
		firestoreArticle := models.FirestoreArticle{
			ID:             fmt.Sprintf("article_%s_%s", company.ID, contentHash[:8]),
			CompanyID:      company.ID,
			CompanyName:    company.Name,
			URL:            item.Link,
			Title:          item.Title,
			Content:        content,
			Summary:        "",
			PublishedAt:    publishedAt,
			ScrapedAt:      time.Now(),
			SourceType:     "rss",
			ContentHash:    contentHash,
			WordCount:      len(content),
			Language:       "en",
			IsProcessed:    false,
			Classification: nil,
			TopEntities:    []models.Entity{},
		}

		// Save article to Firestore
		err = s.saveFirestoreArticle(firestoreArticle)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"company": company.Name,
				"title":   item.Title,
				"error":   err,
			}).Error("Failed to save article to Firestore")
			continue
		}

		newArticleCount++
		
		s.logger.WithFields(logrus.Fields{
			"company":    company.Name,
			"article_id": firestoreArticle.ID,
			"title":      item.Title,
			"hash":       contentHash[:8],
		}).Debug("New article saved successfully")
	}

	s.logger.WithFields(logrus.Fields{
		"company":      company.Name,
		"total_items":  len(feed.Items),
		"new_articles": newArticleCount,
		"duplicates":   duplicateCount,
	}).Info("RSS feed scraping completed")

	return newArticleCount, nil
}

func (s *FirestoreScraperService) saveFirestoreArticle(article models.FirestoreArticle) error {
	ctx := context.Background()
	
	// Check for duplicate by content hash first
	if article.ContentHash != "" {
		iter := s.db.Client().Collection("articles").Where("content_hash", "==", article.ContentHash).Limit(1).Documents(ctx)
		defer iter.Stop()
		
		if doc, err := iter.Next(); err == nil {
			return fmt.Errorf("article already exists with hash %s: %s", article.ContentHash, doc.Ref.ID)
		}
	}

	_, err := s.db.Client().Collection("articles").Doc(article.ID).Set(ctx, article)
	if err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}

	return nil
}

func (s *FirestoreScraperService) ScrapeAllCompanies() error {
	ctx := context.Background()
	
	// Get all companies from Firestore
	iter := s.db.Client().Collection("companies").Where("is_active", "==", true).Documents(ctx)
	defer iter.Stop()

	var companies []models.FirestoreCompany
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		var company models.FirestoreCompany
		if err := doc.DataTo(&company); err != nil {
			s.logger.WithFields(logrus.Fields{
				"company_id": doc.Ref.ID,
				"error":      err,
			}).Warn("Failed to parse company document")
			continue
		}

		companies = append(companies, company)
	}

	s.logger.WithField("count", len(companies)).Info("Starting scraping run for all companies")

	totalNewArticles := 0
	totalErrors := 0

	for _, company := range companies {
		newArticles, err := s.ScrapeCompany(company)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"company": company.Name,
				"error":   err,
			}).Error("Failed to scrape company")
			totalErrors++
			continue
		}

		totalNewArticles += newArticles
		
		// Add a small delay between companies to be respectful
		time.Sleep(2 * time.Second)
	}

	s.logger.WithFields(logrus.Fields{
		"companies_processed": len(companies),
		"total_new_articles":  totalNewArticles,
		"errors":              totalErrors,
	}).Info("Scraping run completed")

	return nil
}