package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	// Initialize the service
	service := NewURLService()

	// Test cases
	testCases := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "Valid URL",
			url:         "https://www.example.com",
			expectError: false,
		},
		{
			name:        "Another Valid URL",
			url:         "https://github.com/piyushhbhutoria/infracloud",
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortURL, err := service.ShortenURL(tc.url)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, shortURL)

				// Verify we can get the original URL back
				originalURL, exists := service.GetOriginalURL(shortURL)
				assert.True(t, exists)
				assert.Equal(t, tc.url, originalURL)
			}
		})
	}
}

func TestShortenURLIdempotence(t *testing.T) {
	// Initialize the service
	service := NewURLService()

	// Shorten the same URL twice
	url := "https://www.example.com"

	shortURL1, err1 := service.ShortenURL(url)
	assert.NoError(t, err1)
	assert.NotEmpty(t, shortURL1)

	shortURL2, err2 := service.ShortenURL(url)
	assert.NoError(t, err2)
	assert.NotEmpty(t, shortURL2)

	// Verify both operations return the same short URL
	assert.Equal(t, shortURL1, shortURL2)
}

func TestGetTopDomains(t *testing.T) {
	// Initialize the service
	service := NewURLService()

	// Shorten URLs from different domains
	domains := map[string]int{
		"https://www.youtube.com/watch?v=video1":       1,
		"https://www.youtube.com/watch?v=video2":       1,
		"https://www.youtube.com/watch?v=video3":       1,
		"https://www.youtube.com/watch?v=video4":       1, // youtube = 4
		"https://stackoverflow.com/questions/123":      1, // stackoverflow = 1
		"https://en.wikipedia.org/wiki/Go":             1,
		"https://en.wikipedia.org/wiki/URL_shortening": 1, // wikipedia = 2
		"https://www.udemy.com/course/course1":         1,
		"https://www.udemy.com/course/course2":         1,
		"https://www.udemy.com/course/course3":         1,
		"https://www.udemy.com/course/course4":         1,
		"https://www.udemy.com/course/course5":         1,
		"https://www.udemy.com/course/course6":         1, // udemy = 6
	}

	// Shorten all the URLs
	for url := range domains {
		_, err := service.ShortenURL(url)
		assert.NoError(t, err)
	}

	// Get top domains
	topDomains := service.GetTopDomains(3)

	// Verify the order: Udemy, YouTube, Wikipedia
	assert.Equal(t, 3, len(topDomains))
	assert.Equal(t, "udemy", topDomains[0].Domain)
	assert.Equal(t, 6, topDomains[0].Count)
	assert.Equal(t, "youtube", topDomains[1].Domain)
	assert.Equal(t, 4, topDomains[1].Count)
	assert.Equal(t, "wikipedia", topDomains[2].Domain)
	assert.Equal(t, 2, topDomains[2].Count)
}

func TestExtractDomain(t *testing.T) {
	testCases := []struct {
		name           string
		url            string
		expectedDomain string
		expectError    bool
	}{
		{
			name:           "Standard domain",
			url:            "https://www.example.com",
			expectedDomain: "example",
			expectError:    false,
		},
		{
			name:           "Domain without subdomain",
			url:            "https://example.com",
			expectedDomain: "example",
			expectError:    false,
		},
		{
			name:           "Domain with subdomain",
			url:            "https://sub.example.com",
			expectedDomain: "example",
			expectError:    false,
		},
		{
			name:           "UK domain",
			url:            "https://www.example.co.uk",
			expectedDomain: "example",
			expectError:    false,
		},
		{
			name:           "Invalid URL format",
			url:            "data:text/plain,invalid", // Use a data URL that's valid but can't extract a domain
			expectedDomain: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			domain, err := extractDomain(tc.url)

			if tc.expectError {
				assert.Error(t, err, "Expected error for URL: %s", tc.url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDomain, domain)
			}
		})
	}
}
