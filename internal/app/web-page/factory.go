package webpage

func NewWebScraper() WebScraper {
	return NewCollyScraper()
}

func NewWebScraperWithDebug() WebScraper {
	return NewCollyScraperWithDebug()
}