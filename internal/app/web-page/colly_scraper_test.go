package webpage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollyScraper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CollyScraper Suite")
}

var _ = Describe("CollyScraper", func() {
	var (
		scraper *CollyScraper
		ctx     context.Context
		options *ScrapingOptions
	)

	BeforeEach(func() {
		scraper = NewCollyScraper()
		ctx = context.Background()
		options = DefaultScrapingOptions()
	})

	Describe("Scrape", func() {
		Context("when scraping a mock HTML page", func() {
			var (
				server   *httptest.Server
				mockHTML []byte
			)

			BeforeEach(func() {
				var err error
				mockHTML, err = os.ReadFile("mock/testing-web-page.html")
				require.NoError(GinkgoT(), err, "Failed to read mock HTML file")

				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusOK)
					w.Write(mockHTML)
				}))
			})

			AfterEach(func() {
				server.Close()
			})

			It("should successfully scrape the page", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).ToNot(BeNil())
			})

			It("should extract basic page information", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				require.NoError(GinkgoT(), err)

				assert.Equal(GinkgoT(), server.URL, result.URL)
				assert.Equal(GinkgoT(), 200, result.StatusCode)
			})

			It("should extract the correct title", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				require.NoError(GinkgoT(), err)

				expectedTitle := "Meteo Lille (59000) - Nord : Prévisions Meteo GRATUITE à 15 jours - La Chaîne Météo"
				assert.Equal(GinkgoT(), expectedTitle, result.Title)
			})

			It("should extract HTML body content", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				require.NoError(GinkgoT(), err)

				assert.NotEmpty(GinkgoT(), result.HTMLBody)
				assert.Contains(GinkgoT(), result.HTMLBody, "<!DOCTYPE html>")
			})

			It("should extract meta tags", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				require.NoError(GinkgoT(), err)

				assert.NotEmpty(GinkgoT(), result.MetaTags)
				assert.NotEmpty(GinkgoT(), result.MetaTags["description"])
				assert.NotEmpty(GinkgoT(), result.MetaTags["og:title"])
			})

			It("should extract links", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				require.NoError(GinkgoT(), err)

				assert.NotEmpty(GinkgoT(), result.Links)
			})

			It("should extract text content", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				require.NoError(GinkgoT(), err)

				assert.NotEmpty(GinkgoT(), result.Text)
			})

			It("should set scraped timestamp", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				require.NoError(GinkgoT(), err)

				assert.False(GinkgoT(), result.ScrapedAt.IsZero())
			})
		})

		Context("when handling errors", func() {
			It("should handle invalid URLs gracefully", func() {
				result, err := scraper.Scrape(ctx, "http://invalid-url-that-does-not-exist.local", options)

				Expect(err).To(HaveOccurred())
				Expect(result).ToNot(BeNil())
			})
		})

		Context("when testing timeouts", func() {
			var slowServer *httptest.Server

			BeforeEach(func() {
				slowServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2 * time.Second)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("<html><head><title>Slow Page</title></head><body>Content</body></html>"))
				}))

				options = &ScrapingOptions{
					Timeout:         500 * time.Millisecond,
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
			})

			AfterEach(func() {
				slowServer.Close()
			})

			It("should timeout on slow responses", func() {
				_, err := scraper.Scrape(ctx, slowServer.URL, options)

				Expect(err).To(HaveOccurred())
				assert.Contains(GinkgoT(), err.Error(), "timeout")
			})
		})
	})

	Describe("ScrapeMultiple", func() {
		Context("when scraping multiple pages", func() {
			var (
				server1 *httptest.Server
				server2 *httptest.Server
			)

			BeforeEach(func() {
				server1 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("<html><head><title>Page 1</title></head><body><h1>Hello Page 1</h1></body></html>"))
				}))

				server2 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("<html><head><title>Page 2</title></head><body><h1>Hello Page 2</h1></body></html>"))
				}))

				options = &ScrapingOptions{
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
			})

			AfterEach(func() {
				server1.Close()
				server2.Close()
			})

			It("should scrape multiple URLs successfully", func() {
				urls := []string{server1.URL, server2.URL}
				results, err := scraper.ScrapeMultiple(ctx, urls, options)

				Expect(err).ToNot(HaveOccurred())
				assert.Len(GinkgoT(), results, 2)
				assert.Equal(GinkgoT(), "Page 1", results[0].Title)
				assert.Equal(GinkgoT(), "Page 2", results[1].Title)
			})
		})
	})

	Describe("Scraping Options", func() {
		Context("when using selective extraction options", func() {
			var server *httptest.Server

			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

				options = &ScrapingOptions{
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
			})

			AfterEach(func() {
				server.Close()
			})

			It("should extract only selected content types", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				require.NoError(GinkgoT(), err)

				assert.NotEmpty(GinkgoT(), result.HTMLBody, "HTMLBody should be extracted")
				assert.NotEmpty(GinkgoT(), result.MetaTags, "MetaTags should be extracted")

				assert.Empty(GinkgoT(), result.Text, "Text should not be extracted when ExtractText is false")
				assert.Empty(GinkgoT(), result.Links, "Links should not be extracted when ExtractLinks is false")
				assert.Empty(GinkgoT(), result.Images, "Images should not be extracted when ExtractImages is false")
				assert.Empty(GinkgoT(), result.Scripts, "Scripts should not be extracted when ExtractScripts is false")
			})
		})
	})

	Describe("Default Options", func() {
		It("should have correct default values", func() {
			defaultOptions := DefaultScrapingOptions()

			assert.Equal(GinkgoT(), 30*time.Second, defaultOptions.Timeout)
			assert.True(GinkgoT(), defaultOptions.ExtractText)
			assert.True(GinkgoT(), defaultOptions.ExtractHTML)
			assert.True(GinkgoT(), defaultOptions.ExtractLinks)
			assert.True(GinkgoT(), defaultOptions.ExtractImages)
			assert.True(GinkgoT(), defaultOptions.ExtractForms)
			assert.True(GinkgoT(), defaultOptions.ExtractScripts)
			assert.True(GinkgoT(), defaultOptions.ExtractMeta)
		})
	})

	Describe("Configuration", func() {
		Describe("SetUserAgent", func() {
			It("should set custom user agent", func() {
				customUA := "Custom-Bot/2.0"
				scraper.SetUserAgent(customUA)

				assert.Equal(GinkgoT(), customUA, scraper.userAgent)
			})
		})

		Describe("SetTimeout", func() {
			It("should set custom timeout", func() {
				customTimeout := 45 * time.Second
				scraper.SetTimeout(customTimeout)

				assert.Equal(GinkgoT(), customTimeout, scraper.timeout)
			})
		})
	})
})
