package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/api-gateway/config"
)

func main() {
	cfg := config.LoadConfig()

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "api-gateway"})
	})

	// Proxy handler
	router.Any("/api/v1/*path", func(c *gin.Context) {
		path := c.Param("path")
		
		var targetURL string

		// Procurement Service Routes (check these first due to overlap)
		if strings.HasPrefix(path, "/suppliers") || 
		   strings.HasPrefix(path, "/purchase-orders") || 
		   strings.HasPrefix(path, "/grns") {
			targetURL = cfg.ProcurementServiceURL
		} else if strings.HasPrefix(path, "/organizations") {
			// Check for procurement sub-resources under organizations
			// /organizations/:id/suppliers
			// /organizations/:id/purchase-orders
			// /organizations/:id/grns
			// We can check if the path contains these segments
			if strings.Contains(path, "/suppliers") || 
			   strings.Contains(path, "/purchase-orders") || 
			   strings.Contains(path, "/grns") {
				targetURL = cfg.ProcurementServiceURL
			} else {
				// Default to Org Service for other /organizations routes
				targetURL = cfg.OrgServiceURL
			}
		} else if strings.HasPrefix(path, "/auth") {
			targetURL = cfg.AuthServiceURL
		} else if strings.HasPrefix(path, "/companies") || strings.HasPrefix(path, "/locations") || strings.HasPrefix(path, "/users") {
			// Users/me/access is in org service
			targetURL = cfg.OrgServiceURL
		} else if strings.HasPrefix(path, "/products") || strings.HasPrefix(path, "/categories") || strings.HasPrefix(path, "/brands") || strings.HasPrefix(path, "/attributes")  {
			// Assuming categories, brands, attributes are in product service
			targetURL = cfg.ProductServiceURL
		} else if strings.HasPrefix(path, "/inventory") || strings.HasPrefix(path, "/stock") || strings.HasPrefix(path, "/warehouses") {
			// Assuming stock, warehouses are in inventory service
			targetURL = cfg.InventoryServiceURL
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		remote, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		
		// Update the director to keeping the original path
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = remote.Host
			// We can strip prefix if needed, but here we forward the full /api/v1/... path 
			// if the services expect /api/v1/...
			// The services in this project seem to define groups like router.Group("/api/v1") in their main.go
			// So forwarding the full path is correct.
		}
        
        // Error handler
        proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
            log.Printf("Proxy error: %v", err)
            w.WriteHeader(http.StatusBadGateway)
            w.Write([]byte("Bad Gateway"))
        }

		proxy.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("API Gateway starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
