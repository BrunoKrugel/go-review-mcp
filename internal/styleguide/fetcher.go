package styleguide

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// StyleGuideURL represents a style guide source
type StyleGuideURL struct {
	Name        string
	URL         string
	Description string
}

// GetGoogleStyleGuideURLs returns all Google style guide URLs
func GetGoogleStyleGuideURLs() []StyleGuideURL {
	return []StyleGuideURL{
		{
			Name:        "google-guide",
			URL:         "https://google.github.io/styleguide/go/guide",
			Description: "Google Go Style Guide",
		},
		{
			Name:        "google-decisions",
			URL:         "https://google.github.io/styleguide/go/decisions",
			Description: "Google Go Style Decisions",
		},
		{
			Name:        "google-practices",
			URL:         "https://google.github.io/styleguide/go/best-practices",
			Description: "Google Go Best Practices",
		},
	}
}

// GetUberStyleGuideURL returns the Uber style guide URL
func GetUberStyleGuideURL() StyleGuideURL {
	return StyleGuideURL{
		Name:        "uber",
		URL:         "https://raw.githubusercontent.com/uber-go/guide/master/style.md",
		Description: "Uber Go Style Guide",
	}
}

const (
	// DefaultHTTPTimeout is the default timeout for HTTP requests
	DefaultHTTPTimeout = 30 * time.Second
)

// Fetcher handles fetching and caching style guide content
type Fetcher struct {
	cache     map[string]*CachedContent
	client    *http.Client
	userAgent string
	cacheTTL  time.Duration
	cacheMu   sync.RWMutex
}

// CachedContent represents cached style guide content
type CachedContent struct {
	Content   string
	FetchedAt time.Time
	URL       string
}

// NewFetcher creates a new style guide fetcher
func NewFetcher(cacheTTL time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: DefaultHTTPTimeout,
		},
		cache:     make(map[string]*CachedContent),
		cacheTTL:  cacheTTL,
		userAgent: "go-review-mcp/1.0",
	}
}

// FetchStyleGuide fetches a style guide from URL with caching
func (f *Fetcher) FetchStyleGuide(ctx context.Context, url string) (string, error) {
	f.cacheMu.RLock()
	cached, exists := f.cache[url]
	f.cacheMu.RUnlock()

	if exists && time.Since(cached.FetchedAt) < f.cacheTTL {
		return cached.Content, nil
	}

	content, err := f.fetchFromURL(ctx, url)
	if err != nil {
		if exists {
			return cached.Content, nil
		}
		return "", fmt.Errorf("failed to fetch %s: %w", url, err)
	}

	f.cacheMu.Lock()
	f.cache[url] = &CachedContent{
		Content:   content,
		FetchedAt: time.Now(),
		URL:       url,
	}
	f.cacheMu.Unlock()

	return content, nil
}

func (f *Fetcher) fetchFromURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}

// FetchAll fetches all style guides concurrently
func (f *Fetcher) FetchAll(ctx context.Context) (map[string]string, error) {
	type result struct {
		err     error
		name    string
		content string
	}

	googleGuides := GetGoogleStyleGuideURLs()
	uberGuide := GetUberStyleGuideURL()
	totalGuides := len(googleGuides) + 1
	results := make(chan result, totalGuides)
	var wg sync.WaitGroup

	for _, guide := range googleGuides {
		wg.Add(1)
		go func(g StyleGuideURL) {
			defer wg.Done()
			content, err := f.FetchStyleGuide(ctx, g.URL)
			results <- result{name: g.Name, content: content, err: err}
		}(guide)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		content, err := f.FetchStyleGuide(ctx, uberGuide.URL)
		results <- result{name: uberGuide.Name, content: content, err: err}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	guides := make(map[string]string)
	var errors []error

	for r := range results {
		if r.err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", r.name, r.err))
			continue
		}
		guides[r.name] = r.content
	}

	if len(errors) > 0 && len(guides) == 0 {
		return nil, fmt.Errorf("failed to fetch any style guides: %v", errors)
	}

	return guides, nil
}

// InvalidateCache clears the cache for a specific URL or all
func (f *Fetcher) InvalidateCache(url string) {
	f.cacheMu.Lock()
	defer f.cacheMu.Unlock()

	if url == "" {
		f.cache = make(map[string]*CachedContent)
	} else {
		delete(f.cache, url)
	}
}

// GetCacheInfo returns information about cached content
func (f *Fetcher) GetCacheInfo() map[string]time.Time {
	f.cacheMu.RLock()
	defer f.cacheMu.RUnlock()

	info := make(map[string]time.Time)
	for url, cached := range f.cache {
		info[url] = cached.FetchedAt
	}
	return info
}
