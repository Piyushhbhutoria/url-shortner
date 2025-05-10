package controllers

import (
	"net/http"

	"github.com/Piyushhbhutoria/url-shortner/services"
	"github.com/gin-gonic/gin"
)

// URLController handles the URL shortening API endpoints
type URLController struct {
	urlService *services.URLService
}

// NewURLController creates a new URL handler
func NewURLController(urlService *services.URLService) *URLController {
	return &URLController{
		urlService: urlService,
	}
}

// ShortenURLRequest is the request body for shortening a URL
type ShortenURLRequest struct {
	URL string `json:"url" binding:"required"`
}

// ShortenURLResponse is the response body for shortening a URL
type ShortenURLResponse struct {
	ShortURL string `json:"short_url"`
}

// TopDomainsResponse is the response for top domains metrics
type TopDomainsResponse struct {
	Domains []DomainMetric `json:"domains"`
}

// DomainMetric represents a domain with its count
type DomainMetric struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

// HandleShortenURL handles the URL shortening request
func (h *URLController) HandleShortenURL(c *gin.Context) {
	var request ShortenURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Basic validation
	if request.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	// Shorten the URL
	shortURL, err := h.urlService.ShortenURL(request.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL: " + err.Error()})
		return
	}

	// Create full short URL with host
	baseURL := getBaseURL(c.Request)
	fullShortURL := baseURL + "/r/" + shortURL

	c.JSON(http.StatusOK, ShortenURLResponse{
		ShortURL: fullShortURL,
	})
}

// HandleRedirect handles the redirection from short URL to original URL
func (h *URLController) HandleRedirect(c *gin.Context) {
	shortURL := c.Param("shortURL")

	originalURL, exists := h.urlService.GetOriginalURL(shortURL)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// HandleGetTopDomains returns the top 3 domains that have been shortened
func (h *URLController) HandleGetTopDomains(c *gin.Context) {
	topDomains := h.urlService.GetTopDomains(3)

	// Convert to response format
	var response TopDomainsResponse
	for _, domain := range topDomains {
		response.Domains = append(response.Domains, DomainMetric{
			Domain: domain.Domain,
			Count:  domain.Count,
		})
	}

	c.JSON(http.StatusOK, response)
}

// getBaseURL extracts the base URL from the request
func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
