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
				Expect(err).NotTo(HaveOccurred(), "Failed to read mock HTML file")

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
				Expect(err).NotTo(HaveOccurred())

				Expect(result.URL).To(Equal(server.URL))
				Expect(result.StatusCode).To(Equal(200))
			})

			It("should extract the correct title", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				Expect(err).NotTo(HaveOccurred())

				expectedTitle := "Meteo Lille (59000) - Nord : Prévisions Meteo GRATUITE à 15 jours - La Chaîne Météo"
				Expect(result.Title).To(Equal(expectedTitle))
			})

			It("should extract HTML body content", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.HTMLBody).NotTo(BeEmpty())
				Expect(result.HTMLBody).To(ContainSubstring("<!DOCTYPE html>"))
			})

			It("should extract meta tags", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.MetaTags).NotTo(BeEmpty())
				Expect(result.MetaTags["description"]).NotTo(BeEmpty())
				Expect(result.MetaTags["og:title"]).NotTo(BeEmpty())
			})

			It("should extract links", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Links).NotTo(BeEmpty())
			})

			It("should extract text content", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Text).NotTo(BeEmpty())
			})

			It("should set scraped timestamp", func() {
				result, err := scraper.Scrape(ctx, server.URL, options)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.ScrapedAt.IsZero()).To(BeFalse())
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
				Expect(err.Error()).To(ContainSubstring("timeout"))
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
				Expect(results).To(HaveLen(2))
				Expect(results[0].Title).To(Equal("Page 1"))
				Expect(results[1].Title).To(Equal("Page 2"))
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
				Expect(err).NotTo(HaveOccurred())

				Expect(result.HTMLBody).NotTo(BeEmpty(), "HTMLBody should be extracted")
				Expect(result.MetaTags).NotTo(BeEmpty(), "MetaTags should be extracted")

				Expect(result.Text).To(BeEmpty(), "Text should not be extracted when ExtractText is false")
				Expect(result.Links).To(BeEmpty(), "Links should not be extracted when ExtractLinks is false")
				Expect(result.Images).To(BeEmpty(), "Images should not be extracted when ExtractImages is false")
				Expect(result.Scripts).To(BeEmpty(), "Scripts should not be extracted when ExtractScripts is false")
			})
		})
	})

	Describe("Default Options", func() {
		It("should have correct default values", func() {
			defaultOptions := DefaultScrapingOptions()

			Expect(defaultOptions.Timeout).To(Equal(30 * time.Second))
			Expect(defaultOptions.ExtractText).To(BeTrue())
			Expect(defaultOptions.ExtractHTML).To(BeTrue())
			Expect(defaultOptions.ExtractLinks).To(BeTrue())
			Expect(defaultOptions.ExtractImages).To(BeTrue())
			Expect(defaultOptions.ExtractForms).To(BeTrue())
			Expect(defaultOptions.ExtractScripts).To(BeTrue())
			Expect(defaultOptions.ExtractMeta).To(BeTrue())
		})
	})

	Describe("Configuration", func() {
		Describe("SetUserAgent", func() {
			It("should set custom user agent", func() {
				customUA := "Custom-Bot/2.0"
				scraper.SetUserAgent(customUA)

				Expect(scraper.userAgent).To(Equal(customUA))
			})
		})

		Describe("SetTimeout", func() {
			It("should set custom timeout", func() {
				customTimeout := 45 * time.Second
				scraper.SetTimeout(customTimeout)

				Expect(scraper.timeout).To(Equal(customTimeout))
			})
		})
	})
})
