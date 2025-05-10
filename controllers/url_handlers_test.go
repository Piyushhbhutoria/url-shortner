package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Piyushhbhutoria/url-shortner/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() (*gin.Engine, *URLController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	urlService := services.NewURLService()
	URLController := NewURLController(urlService)
	return router, URLController
}

func TestHandleShortenURL(t *testing.T) {
	router, URLController := setupRouter()
	router.POST("/api/shorten", URLController.HandleShortenURL)

	// Test valid URL
	t.Run("Valid URL", func(t *testing.T) {
		// Create request body
		reqBody := ShortenURLRequest{
			URL: "https://www.example.com",
		}
		jsonBody, _ := json.Marshal(reqBody)

		// Create request
		req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var response ShortenURLResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.ShortURL)
	})

	// Test invalid request (missing URL)
	t.Run("Missing URL", func(t *testing.T) {
		// Create empty request body
		reqBody := ShortenURLRequest{
			URL: "",
		}
		jsonBody, _ := json.Marshal(reqBody)

		// Create request
		req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandleRedirect(t *testing.T) {
	router, URLController := setupRouter()
	router.GET("/r/:shortURL", URLController.HandleRedirect)

	// First, create a short URL
	originalURL := "https://www.example.com"
	shortURL, _ := URLController.urlService.ShortenURL(originalURL)

	// Test valid redirect
	t.Run("Valid Redirect", func(t *testing.T) {
		// Create request
		req, _ := http.NewRequest("GET", "/r/"+shortURL, nil)
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check redirect
		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, originalURL, w.Header().Get("Location"))
	})

	// Test invalid short URL
	t.Run("Invalid Short URL", func(t *testing.T) {
		// Create request with non-existent short URL
		req, _ := http.NewRequest("GET", "/r/nonexistent", nil)
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandleGetTopDomains(t *testing.T) {
	router, URLController := setupRouter()
	router.GET("/api/metrics/top-domains", URLController.HandleGetTopDomains)

	// First, create some short URLs to populate domain counts
	urls := []string{
		"https://www.youtube.com/watch?v=video1",
		"https://www.youtube.com/watch?v=video2",
		"https://www.udemy.com/course/course1",
		"https://www.udemy.com/course/course2",
		"https://www.udemy.com/course/course3",
		"https://en.wikipedia.org/wiki/URL_shortening",
	}

	for _, url := range urls {
		_, _ = URLController.urlService.ShortenURL(url)
	}

	// Test getting top domains
	t.Run("Get Top Domains", func(t *testing.T) {
		// Create request
		req, _ := http.NewRequest("GET", "/api/metrics/top-domains", nil)
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var response TopDomainsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// We should have 3 domains at most
		assert.LessOrEqual(t, len(response.Domains), 3)

		// Top domain should be udemy with count 3
		if len(response.Domains) > 0 {
			assert.Equal(t, "udemy", response.Domains[0].Domain)
			assert.Equal(t, 3, response.Domains[0].Count)
		}
	})
}
