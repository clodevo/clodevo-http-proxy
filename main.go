package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/clodevo/raven-proxy/docs" // Assuming this is for documentation purposes.

	"github.com/clodevo/raven-proxy/pkg/handlers"

	"github.com/clodevo/raven-proxy/pkg/acl"
	"github.com/clodevo/raven-proxy/pkg/config"
	"github.com/clodevo/raven-proxy/pkg/database"
	"github.com/clodevo/raven-proxy/pkg/proxy"
	"github.com/clodevo/raven-proxy/pkg/utils"

	"github.com/appleboy/graceful"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // SQLite driver for database interactions.

	// Configuration management.
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/valyala/fasthttp"
)

var (
	db          *sql.DB // Global database connection.
	adminAPIKey string  // Admin API key for authentication.
	aclDataPath string  // Path to ACL data files.
)

// @title Clodevo Forward Proxy API
// @version 1.0
// @description This is a forward proxy server made by Clodevo.
// @contact.name API Support
// @contact.url https://www.clodevo.com/
// @contact.email support@clodevo.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Admin-API-Key
// @BasePath /
// @schemes http
func main() {
	appConfig := config.LoadAppConfig()

	// Simplify database initialization
	database.InitDB(appConfig.DatabaseConfig)
	defer db.Close() // Ensure db is a properly initialized *sql.DB in the database package

	// Convert string log level to LogLevel type
	logLevel := StringToLogLevel(appConfig.LogLevel) // You need to implement this function
	utils.GetLogger().SetLogLevel(logLevel)

	// Initialize ACLManager with the logger
	aclManager := acl.NewACLManager(appConfig.ACLDataPath, utils.GetLogger())

	// Admin API key and ACL data path are now directly accessible
	adminAPIKey = appConfig.AdminAPIKey
	aclDataPath = appConfig.ACLDataPath

	// Git synchronization setup (if needed)
	if appConfig.GitSyncConfig.RepoURL != "" {
		utils.Sync(&appConfig.GitSyncConfig)
	}

	// Initialize services (e.g., HTTP servers)

	router := setupRouter()
	startAdminServer(router, appConfig)
	startProxyServer(&appConfig.ProxyConfig, aclManager)

	// Wait for graceful shutdown
	waitForShutdown()
}

// setupRoutes configures the API endpoints.
func setupRouter() *gin.Engine {
	router := gin.Default()

	// Swagger documentation endpoint.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1), ginSwagger.PersistAuthorization(true),
	))

	// Public welcome route.
	router.GET("/", handlers.AdminPageHandler)

	// Authenticated routes setup.
	authenticatedRoutes := router.Group("/", authMiddleware)
	setupAuthenticatedRoutes(authenticatedRoutes)

	return router
}

// authMiddleware handles API key authentication.
func authMiddleware(c *gin.Context) {
	apiKey := c.GetHeader("X-Admin-API-Key")
	if apiKey != adminAPIKey {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"message":       "Invalid API key",
		})
		return
	}
	c.Next()
}

// setupAuthenticatedRoutes defines routes that require authentication.
func setupAuthenticatedRoutes(group *gin.RouterGroup) {
	group.GET("/tenants", handlers.TenantsHandler)
	group.GET("/tenants/:tenantID", handlers.TenantsHandler)
	group.POST("/tenants", handlers.TenantsHandler)
	group.PUT("/tenants/:tenantID", handlers.TenantsHandler)
	group.DELETE("/tenants/:tenantID", handlers.TenantsHandler)

	group.GET("/:tenantID/api-keys", handlers.GetTenantAPIKey)
	group.POST("/:tenantID/api-keys", handlers.CreateAPIKey)
	group.PUT("/:tenantID/api-keys/:apiKeyID/rotate", handlers.RotateAPIKey)
	group.DELETE("/:tenantID/api-keys/:apiKeyID", handlers.DeleteAPIKey)
}

// startAdminServer initializes and starts the Gin HTTP server.
func startAdminServer(router *gin.Engine, adminConfig *config.AppConfig) {
	adminServer := &http.Server{
		Addr:         adminConfig.AdminAddr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Printf("Admin server started at %s\n", adminConfig.AdminAddr)
		if err := adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error in admin ListenAndServe: %s\n", err)
		}
	}()
}

// startProxyServer initializes and starts the FastHTTP server.
func startProxyServer(proxyConfig *config.ProxyConfig, aclManager *acl.ACLManager) {
	server := &fasthttp.Server{
		Handler:            fasthttp.CompressHandler(proxy.FastHTTPHandler(proxyConfig, aclManager)),
		ReadTimeout:        proxyConfig.Timeout,
		WriteTimeout:       proxyConfig.Timeout,
		MaxConnsPerIP:      1000, // Consider making this configurable as well
		MaxRequestsPerConn: 1000, // Consider making this configurable as well
		IdleTimeout:        3 * proxyConfig.Timeout,
		ReduceMemoryUsage:  true,
		CloseOnShutdown:    true,
		Concurrency:        proxyConfig.MaxConcurrent,
	}

	go func() {
		fmt.Printf("Proxy server started at %s\n", proxyConfig.Addr)
		if err := server.ListenAndServe(proxyConfig.Addr); err != nil { // Use the address from the configuration
			fmt.Printf("Error in Proxy ListenAndServe: %s\n", err)
		}
	}()
}

// waitForShutdown handles graceful shutdown on interrupt signals.
func waitForShutdown() {
	graceful.NewManager().AddRunningJob(func(ctx context.Context) error {
		<-ctx.Done()
		fmt.Println("Servers are shutting down")
		// Implement additional shutdown logic if necessary.
		return nil
	})

	<-graceful.NewManager().Done()
	fmt.Println("Gracefully stopped")
}

// This function should map string representations of log levels to the LogLevel constants
func StringToLogLevel(level string) utils.LogLevel {
	switch level {
	case "debug":
		return utils.LogLevelDebug
	case "trace":
		return utils.LogLevelTrace
	default:
		return utils.LogLevelInfo
	}
}
