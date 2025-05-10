package services

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
)

// URLService handles URL shortening operations
type URLService struct {
	urlMap       map[string]string // shortURL -> originalURL
	reverseMap   map[string]string // originalURL -> shortURL
	domainCounts map[string]int    // domain -> count
	mutex        sync.RWMutex
}

// NewURLService creates a new URL service
func NewURLService() *URLService {
	return &URLService{
		urlMap:       make(map[string]string),
		reverseMap:   make(map[string]string),
		domainCounts: make(map[string]int),
		mutex:        sync.RWMutex{},
	}
}

// ShortenURL shortens a URL and returns the shortened version
func (s *URLService) ShortenURL(originalURL string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if we've already shortened this URL
	if shortURL, exists := s.reverseMap[originalURL]; exists {
		return shortURL, nil
	}

	// Generate a shorter URL
	shortURL := s.generateShortURL(originalURL)

	// Store the mapping
	s.urlMap[shortURL] = originalURL
	s.reverseMap[originalURL] = shortURL

	// Update domain metrics
	domain, err := extractDomain(originalURL)
	if err == nil {
		s.domainCounts[domain]++
	}

	return shortURL, nil
}

// GetOriginalURL retrieves the original URL from a shortened URL
func (s *URLService) GetOriginalURL(shortURL string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	originalURL, exists := s.urlMap[shortURL]
	return originalURL, exists
}

// GetTopDomains returns the top N domains by frequency
func (s *URLService) GetTopDomains(n int) []DomainCount {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Convert map to slice for sorting
	var domainStats []DomainCount
	for domain, count := range s.domainCounts {
		domainStats = append(domainStats, DomainCount{Domain: domain, Count: count})
	}

	// Sort by count (descending)
	sort.Slice(domainStats, func(i, j int) bool {
		return domainStats[i].Count > domainStats[j].Count
	})

	// Return top N or all if less than N
	if len(domainStats) < n {
		return domainStats
	}
	return domainStats[:n]
}

// DomainCount represents a domain and its count
type DomainCount struct {
	Domain string
	Count  int
}

// generateShortURL creates a short URL string from the original URL
func (s *URLService) generateShortURL(originalURL string) string {
	// Generate SHA-256 hash of the URL
	hasher := sha256.New()
	hasher.Write([]byte(originalURL))
	hash := hasher.Sum(nil)

	// Encode first 6 bytes of hash to base64 to get a short URL
	encoded := base64.URLEncoding.EncodeToString(hash[:6])

	// Use only the first 8 characters for shorter URL
	return encoded[:8]
}

// extractDomain extracts the domain from a URL
func extractDomain(rawURL string) (string, error) {
	// Try to parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Check if scheme is data: or host is empty
	if parsedURL.Scheme == "data" || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL format or cannot extract domain")
	}

	// Extract domain from host (remove port number if present)
	domain := parsedURL.Hostname()

	// Further simplify by keeping only the main domain
	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		// For domains like example.co.uk or subdomain.example.com
		// This is a simplistic approach, more robust would use a public suffix list
		if len(parts) > 2 && parts[len(parts)-2] == "co" {
			// Special case for .co.uk, .co.jp, etc.
			return parts[len(parts)-3], nil
		}
		return parts[len(parts)-2], nil
	}

	return domain, nil
}
