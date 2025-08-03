# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Developpement

- If you don't have all the information, ASK.
- If you don't know, ASK.
- No TODOs in the code, no unused functions
- Comments must be in ENGLISH.
- "Everything should be made as simple as possible, but not simpler. Nothing is more simple than greatness; indeed, to be simple is to be great" Albert Einstein

## Testing

- TEST-DRIVEN DEVELOPMENT IS NON-NEGOTIABLE.
- Run all tests: `go test -v ./...`
- Run specific package tests: `go test -v ./internal/app/web-page`
- Tests use Ginkgo BDD framework with Testify assertions

## Building

- Build all packages: `go build ./...`
- The project is structured as a library with no main executable

## Module Management

- Module name: `mmm-osint`
- Go version: 1.23.0

## Architecture

WHEN THERE IS A CHANGE IN THE ARCHITECTURE, UPDATE THE CLAUDE.MD FILE.

This is a Go-based OSINT (Open Source Intelligence) toolkit with a modular architecture focused on web scraping, caching, and queuing capabilities.

### Core Components

**Web Scraping (`internal/app/web-page/`)**

- Main interface: `WebScraper` with implementations using Colly framework
- Supports comprehensive data extraction: HTML, text, links, images, forms, scripts, meta tags
- Configurable scraping options: timeouts, user agents, rate limiting, selective extraction
- Factory pattern for scraper instantiation
- BDD-style tests using Ginkgo and Testify

**Caching System (`internal/pkg/cache/`)**

- Generic `Cache` interface supporting Set/Get/Delete/Exists operations
- Implementations: LRU cache and Redis cache
- Context-aware operations with expiration support

**Queue System (`internal/pkg/queue/`)**

- Generic queue interface `Queue[T]` for type-safe message handling
- Request-Response pattern with `RequestResponseQueue[T, R]`
- Redis-based implementation for distributed queuing
- Support for both fire-and-forget and request-reply messaging patterns

### Key Design Patterns

- **Interface-driven design**: All major components define interfaces first
- **Factory pattern**: Used for component instantiation (web scrapers, queues)
- **Generic types**: Queue system uses Go generics for type safety
- **Context propagation**: All operations support context for cancellation/timeouts

### Dependencies

Key external dependencies:

- `github.com/gocolly/colly/v2`: Web scraping framework
- `github.com/redis/go-redis/v9`: Redis client for caching and queuing
- `github.com/hashicorp/golang-lru/v2`: LRU cache implementation
- `github.com/onsi/ginkgo/v2` + `github.com/stretchr/testify`: BDD testing framework

### Testing Strategy

- BDD-style tests using Ginkgo's Describe/Context/It structure
- Testify assertions for readable test expectations
- Mock HTTP servers for web scraping tests
- Test coverage includes timeout handling, error scenarios, and configuration options
