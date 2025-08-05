# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

YOUR MOTTO: "Everything should be made as simple as possible, but not simpler. Nothing is more simple than greatness; indeed, to be simple is to be great" **Albert Einstein**

FOLLOW K.I.S.S principle

## Development

- If you don't have all the information, ASK.
- If you don't know, ASK.
- No TODOs in the code, no unused functions
- Comments must be in ENGLISH.
- For each package that is intended to be used by others, always create interfaces. Make sure they are as SIMPLE as possible.

## Testing

- TEST-DRIVEN DEVELOPMENT IS NON-NEGOTIABLE.
- Run all tests: `go test -v ./...`
- Run specific package tests: `go test -v ./internal/pkg/cache`
- Run tests with coverage: `go test -v -cover ./...`
- Run tests for specific file: `go test -v ./internal/pkg/cache -run TestLRUCache`
- Tests use Ginkgo BDD framework with Gomega assertions

### Test Example Structure

All tests must follow the Ginkgo BDD pattern with Gomega assertions:

```go
package mypackage_test

import (
    "context"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "mmm-osint/internal/pkg/mypackage"
)

var _ = Describe("MyComponent", func() {
    var (
        component mypackage.Interface
        ctx       context.Context
    )

    BeforeEach(func() {
        component = mypackage.New()
        ctx = context.Background()
    })

    AfterEach(func() {
        if component != nil {
            component.Close()
        }
    })

    Describe("MethodName", func() {
        Context("with valid input", func() {
            It("should return expected result", func() {
                result, err := component.MethodName(ctx, "input")

                Expect(err).NotTo(HaveOccurred())
                Expect(result).NotTo(BeEmpty())
                Expect(result).To(ContainSubstring("expected"))
            })
        })

        Context("with invalid input", func() {
            It("should handle errors gracefully", func() {
                result, err := component.MethodName(ctx, "")

                Expect(err).To(HaveOccurred())
                Expect(result).To(BeEmpty())
            })
        })
    })
})
```

## Building

- Build all packages: `go build ./...`
- Check for Go formatting issues: `go fmt ./...`
- Check for common Go issues: `go vet ./...`
- The project is structured as a library with no main executable

## Module Management

- Module name: `mmm-osint`
- Go version: 1.23.0

## Architecture

WHEN THERE IS A CHANGE IN THE ARCHITECTURE, UPDATE THE CLAUDE.MD FILE.

This is a Go-based OSINT (Open Source Intelligence) toolkit with a modular architecture focused on web scraping, caching, and queuing capabilities.

### Core Components

**Web Scraping (`internal/pkg/web-page/`)**

- Main interface: `WebScraper` with implementations using Colly framework
- Supports comprehensive data extraction: HTML, text, links, images, forms, scripts, meta tags
- Configurable scraping options: timeouts, user agents, rate limiting, selective extraction
- Factory pattern for scraper instantiation
- BDD-style tests using Ginkgo and Gomega

**Caching System (`internal/pkg/cache/`)**

- Generic `Cache` interface supporting Set/Get/Delete/Exists operations
- Implementations: LRU cache and Redis cache
- Context-aware operations with expiration support

**Queue System (`internal/pkg/queue/`)**

- Generic queue interface `Queue[T]` for type-safe message handling
- Request-Response pattern with `RequestResponseQueue[T, R]`
- Redis-based implementation for distributed queuing
- Support for both fire-and-forget and request-reply messaging patterns

**Keyword Extraction (`internal/pkg/keyword/`)**

- Simple interface for extracting keywords from text
- Prose library integration for natural language processing
- Configurable options: minimum word length, maximum keywords, stop word filtering
- Two extraction methods: simple string list or keywords with frequency scores

**PII Extraction (`internal/pkg/pii/`)**

- Simple interface for extracting Personally Identifiable Information from text
- Regex-based detection of emails, phones, credit cards, SSNs, IP addresses, IBANs
- Built on intMeric/pii-extractor library
- Returns structured results with entity types, values, counts, and contexts

**Graph Database (`internal/pkg/graph/`)**

- Interface for graph database operations with Neo4j implementation
- Typed nodes (URL, User) with validation for displayName and ID fields
- Separate concerns: Neo4j stores relationships, MongoDB stores data
- Methods: CreateNode, CreateRelation, GetNode, NodeExists, Close
- Connection pooling and automatic reconnection handling
- Factory pattern for graph instantiation with environment variable configuration

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
- `github.com/jdkato/prose/v2`: Natural language processing for keyword extraction
- `github.com/intMeric/pii-extractor`: PII detection and extraction
- `github.com/onsi/ginkgo/v2` + `github.com/onsi/gomega`: BDD testing framework
- `github.com/neo4j/neo4j-go-driver/v5`: Neo4j database driver for graph operations

### Testing Strategy

- BDD-style tests using Ginkgo's Describe/Context/It structure
- Gomega assertions for readable test expectations
- Mock HTTP servers for web scraping tests
- Test coverage includes timeout handling, error scenarios, and configuration options

### Directory Structure

- `cmd/`: Command-line executables (currently empty - library project)
- `internal/app/`: Application-specific code
  - `services/`: Business logic services (future: OSINT-specific orchestration)
  - `usecases/`: Use cases and business workflows (future: investigation workflows)
- `internal/pkg/`: Reusable internal packages
  - `web-page/`: Web scraping functionality with Colly integration
  - `cache/`: Generic caching interfaces and implementations
  - `queue/`: Message queue interfaces and Redis implementations
  - `keyword/`: Keyword extraction from text using prose library
  - `pii/`: PII extraction using intMeric/pii-extractor library
  - `env/`: Environment configuration utilities (hostname, env vars with defaults)
  - `graph/`: Graph database interfaces and Neo4j implementation with typed nodes (URL, User) and strict validation
