package webpage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCollyScraper_Scrape(t *testing.T) {
	// Read the mock HTML file
	mockHTML, err := os.ReadFile("mock/testing-web-page.html")
	if err != nil {
		t.Fatalf("Failed to read mock HTML file: %v", err)
	}

	// Create a test server with the mock HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(mockHTML)
	}))
	defer server.Close()

	scraper := NewCollyScraper()
	ctx := context.Background()
	options := DefaultScrapingOptions()

	result, err := scraper.Scrape(ctx, server.URL, options)
	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	// Test basic fields
	if result.URL != server.URL {
		t.Errorf("Expected URL %s, got %s", server.URL, result.URL)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}

	// Test title extraction
	expectedTitle := "Meteo Lille (59000) - Nord : Prévisions Meteo GRATUITE à 15 jours - La Chaîne Météo"
	if result.Title != expectedTitle {
		t.Errorf("Expected title '%s', got '%s'", expectedTitle, result.Title)
	}

	// Test HTML body extraction
	if result.HTMLBody == "" {
		t.Error("HTMLBody should not be empty")
	}

	if !strings.Contains(result.HTMLBody, "<!DOCTYPE html>") {
		t.Error("HTMLBody should contain DOCTYPE declaration")
	}

	// Test meta tags extraction
	if len(result.MetaTags) == 0 {
		t.Error("MetaTags should not be empty")
	}

	// Check specific meta tags
	if result.MetaTags["description"] == "" {
		t.Error("Description meta tag should be extracted")
	}

	if result.MetaTags["og:title"] == "" {
		t.Error("OpenGraph title should be extracted")
	}

	// Test links extraction
	if len(result.Links) == 0 {
		t.Error("Links should be extracted")
	}

	// Test that we have some links (not necessarily canonical since it's in <link> tag, not <a>)
	// Note: canonical links are in <link> tags which are not extracted as Links but as stylesheets/scripts
	if len(result.Links) < 5 { // Should have plenty of links in the weather page
		t.Logf("Only found %d links, expected more", len(result.Links))
	}

	// Test text extraction
	if result.Text == "" {
		t.Error("Text should not be empty")
	}

	// Test timestamps
	if result.ScrapedAt.IsZero() {
		t.Error("ScrapedAt timestamp should be set")
	}
}

func TestCollyScraper_ScrapeMultiple(t *testing.T) {
	// Create multiple test servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Page 1</title></head><body><h1>Hello Page 1</h1></body></html>"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Page 2</title></head><body><h1>Hello Page 2</h1></body></html>"))
	}))
	defer server2.Close()

	scraper := NewCollyScraper()
	ctx := context.Background()
	options := &ScrapingOptions{
		Timeout:         5 * time.Second,
		UserAgent:       "Test-Bot/1.0",
		RateLimitDelay:  100 * time.Millisecond,
		ExtractText:     true,
		ExtractHTML:     true,
		ExtractLinks:    true,
		ExtractImages:   true,
		ExtractForms:    true,
		ExtractScripts:  true,
		ExtractMeta:     true,
		FollowRedirects: true,
	}

	urls := []string{server1.URL, server2.URL}
	results, err := scraper.ScrapeMultiple(ctx, urls, options)

	if err != nil {
		t.Fatalf("ScrapeMultiple failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Test first result
	if results[0].Title != "Page 1" {
		t.Errorf("Expected title 'Page 1', got '%s'", results[0].Title)
	}

	// Test second result
	if results[1].Title != "Page 2" {
		t.Errorf("Expected title 'Page 2', got '%s'", results[1].Title)
	}
}

func TestCollyScraper_ScrapingOptions(t *testing.T) {
	// Test with different options
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<head>
					<title>Test Page</title>
					<meta name="description" content="Test description">
				</head>
				<body>
					<h1>Test Content</h1>
					<a href="/link1">Link 1</a>
					<img src="/image1.jpg" alt="Image 1">
					<script src="/script1.js"></script>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	scraper := NewCollyScraper()
	ctx := context.Background()

	// Test with selective extraction
	options := &ScrapingOptions{
		Timeout:         5 * time.Second,
		UserAgent:       "Test-Bot/1.0",
		ExtractText:     false,
		ExtractHTML:     true,
		ExtractLinks:    false,
		ExtractImages:   false,
		ExtractForms:    false,
		ExtractScripts:  false,
		ExtractMeta:     true,
		FollowRedirects: true,
	}

	result, err := scraper.Scrape(ctx, server.URL, options)
	if err != nil {
		t.Fatalf("Scrape failed: %v", err)
	}

	// Should have HTML and meta tags
	if result.HTMLBody == "" {
		t.Error("HTMLBody should be extracted")
	}

	if len(result.MetaTags) == 0 {
		t.Error("MetaTags should be extracted")
	}

	// Should NOT have text, links, images, scripts
	if result.Text != "" {
		t.Error("Text should not be extracted when ExtractText is false")
	}

	if len(result.Links) > 0 {
		t.Error("Links should not be extracted when ExtractLinks is false")
	}

	if len(result.Images) > 0 {
		t.Error("Images should not be extracted when ExtractImages is false")
	}

	if len(result.Scripts) > 0 {
		t.Error("Scripts should not be extracted when ExtractScripts is false")
	}
}

func TestCollyScraper_ErrorHandling(t *testing.T) {
	scraper := NewCollyScraper()
	ctx := context.Background()
	options := DefaultScrapingOptions()

	// Test invalid URL
	result, err := scraper.Scrape(ctx, "http://invalid-url-that-does-not-exist.local", options)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}

	// Result should still be returned with error info
	if result == nil {
		t.Error("Result should not be nil even on error")
	}
}

func TestCollyScraper_Timeout(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Slow Page</title></head><body>Content</body></html>"))
	}))
	defer server.Close()

	scraper := NewCollyScraper()
	ctx := context.Background()
	options := &ScrapingOptions{
		Timeout:         500 * time.Millisecond, // Short timeout
		UserAgent:       "Test-Bot/1.0",
		ExtractText:     true,
		ExtractHTML:     true,
		ExtractLinks:    true,
		ExtractImages:   true,
		ExtractForms:    true,
		ExtractScripts:  true,
		ExtractMeta:     true,
		FollowRedirects: true,
	}

	_, err := scraper.Scrape(ctx, server.URL, options)
	if err == nil {
		t.Error("Expected timeout error")
	}

	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func TestDefaultScrapingOptions(t *testing.T) {
	options := DefaultScrapingOptions()

	if options.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", options.Timeout)
	}

	if !options.ExtractText {
		t.Error("ExtractText should be true by default")
	}

	if !options.ExtractHTML {
		t.Error("ExtractHTML should be true by default")
	}

	if !options.ExtractLinks {
		t.Error("ExtractLinks should be true by default")
	}

	if !options.ExtractImages {
		t.Error("ExtractImages should be true by default")
	}

	if !options.ExtractForms {
		t.Error("ExtractForms should be true by default")
	}

	if !options.ExtractScripts {
		t.Error("ExtractScripts should be true by default")
	}

	if !options.ExtractMeta {
		t.Error("ExtractMeta should be true by default")
	}
}

func TestCollyScraper_SetUserAgent(t *testing.T) {
	scraper := NewCollyScraper()
	customUA := "Custom-Bot/2.0"

	scraper.SetUserAgent(customUA)

	if scraper.userAgent != customUA {
		t.Errorf("Expected user agent '%s', got '%s'", customUA, scraper.userAgent)
	}
}

func TestCollyScraper_SetTimeout(t *testing.T) {
	scraper := NewCollyScraper()
	customTimeout := 45 * time.Second

	scraper.SetTimeout(customTimeout)

	if scraper.timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, scraper.timeout)
	}
}
