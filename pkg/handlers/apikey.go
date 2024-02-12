package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/clodevo/raven-proxy/pkg/database"
	"github.com/clodevo/raven-proxy/pkg/models"
	"github.com/clodevo/raven-proxy/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// Adjust the import path based on your project structure
	// Assuming utility functions are stored here
)

// type APIKeyGetter struct {
// 	ID        uuid.UUID  `json:"api_key_id"`
// 	TenantID  uuid.UUID  `json:"tenant_id"`
// 	CreatedAt *time.Time `json:"created_at,omitempty"`
// 	UpdatedAt *time.Time `json:"updated_at,omitempty"`
// }

// @Summary Create an API key
// @Description Create a new API key for a tenant
// @Tags api-keys
// @Accept json
// @Produce json
// @Param tenantID path string true "Tenant ID"
// @Success 200 {object} models.APIKey
// @Router /{tenantID}/api-keys [post]
// @Security ApiKeyAuth
func CreateAPIKey(c *gin.Context) {

	tenantIDStr := c.Param("tenantID")
	if tenantIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Tenant ID"})
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tenant ID: " + err.Error()})
		return
	}

	api_key_id := uuid.New()
	// Generate a random API key
	apiKey := utils.GenerateRandomString(32)

	var exists bool = true

	for exists {
		api_key_id = uuid.New()
		// Check if the generated api_key_id exists in the database
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM api_keys WHERE api_key_id = ? AND tenant_id = ?)", api_key_id, tenantID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking api_key_id existence: " + err.Error()})
			return
		}
	}

	stmt, err := database.DB.Prepare("INSERT INTO api_keys (api_key_id,api_key, tenant_id) VALUES ( ? ,?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(api_key_id, apiKey, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	now := time.Now() // Store the current time in a variable
	apiKeyObj := models.APIKey{
		ID:        api_key_id,
		Key:       apiKey,
		TenantID:  tenantID,
		CreatedAt: &now,
	}

	message := "API key created successfully. Please ensure you copy and securely store this API key immediately; it will not be displayed again for security reasons."
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    apiKeyObj,
	})
}

// @Summary Rotate API key
// @Description Rotate the API key associated with the given ID and tenant, generating a new key while keeping the same ID
// @Tags api-keys
// @Accept json
// @Produce json
// @Param tenantID path string true "Tenant ID"
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} models.APIKey
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router  /{tenantID}/api-keys/{apiKeyID}/rotate [put]
// @Security ApiKeyAuth
func RotateAPIKey(c *gin.Context) {

	tenantIDStr := c.Param("tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tenant ID: " + err.Error()})
		return
	}

	// Parse the apiKeyID from string to UUID
	apiKeyIDStr := c.Param("apiKeyID")
	apiKeyID, err := uuid.Parse(apiKeyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API Key ID format"})
		return
	}

	// Generate a new random API key
	newAPIKey := utils.GenerateRandomString(32)

	updated_at := time.Now() // Store the current time in a variable

	// Update the API key in the database
	stmt, err := database.DB.Prepare("UPDATE api_keys SET api_key = ?, updated_at = ? WHERE api_key_id = ? AND tenant_id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}
	defer stmt.Close()

	// Note: Ensure that the database driver supports UUIDs directly or convert apiKeyID to string if necessary
	result, err := stmt.Exec(newAPIKey, updated_at, apiKeyID.String(), tenantID.String()) // Convert UUID to string if needed
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found for the given ID and tenant"})
		return
	}

	// Respond with the updated API key
	updatedAPIKey := models.APIKey{
		ID:        apiKeyID, // This assumes the ID field of APIKey is of type uuid.UUID
		Key:       newAPIKey,
		TenantID:  tenantID,
		UpdatedAt: &updated_at,
		CreatedAt: nil,
	}

	c.JSON(http.StatusOK, updatedAPIKey)
}

// @Summary Delete API key
// @Description Delete the API key associated with the given ID and tenant
// @Tags api-keys
// @Accept json
// @Produce json
// @Param tenantID path string true "Tenant ID"
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {string} string "API key deleted successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router  /{tenantID}/api-keys/{apiKeyID} [delete]
// @Security ApiKeyAuth
func DeleteAPIKey(c *gin.Context) {
	tenantID := c.GetInt("tenantID")
	apiKeyIDStr := c.Param("apiKeyID")

	// Parse the apiKeyID from string to UUID
	apiKeyID, err := uuid.Parse(apiKeyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API Key ID format"})
		return
	}

	// Delete the API key from the database
	stmt, err := database.DB.Prepare("DELETE FROM api_keys WHERE api_key_id = ? AND tenant_id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}
	defer stmt.Close()

	// Execute the deletion, converting apiKeyID to string if necessary
	result, err := stmt.Exec(apiKeyID.String(), tenantID) // Convert UUID to string if needed
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found for the given ID and tenant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}

// @Summary Get Tenant API keys IDs
// @Description Get the API keys associated with the given tenantID, only API keys id are shown.
// @Tags api-keys
// @Accept json
// @Produce json
// @Param tenantID path string true "Tenant ID"
// @Success 200 {array} models.APIKey
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router  /{tenantID}/api-keys [get]
// @Security ApiKeyAuth
func GetTenantAPIKey(c *gin.Context) {

	tenantIDStr := c.Param("tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tenant ID: " + err.Error()})
		return
	}

	rows, err := database.DB.Query("SELECT api_key_id, api_key, created_at, updated_at FROM api_keys  WHERE  tenant_id = ? ", tenantID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tenant API keys : " + err.Error()})
		return
	}
	defer rows.Close()

	// Initialize tenant apikeys as an empty slice, not nil.
	apikeys := make([]models.APIKey, 0)
	for rows.Next() {
		var tap models.APIKey
		if err := rows.Scan(&tap.ID, &tap.Key, &tap.CreatedAt, &tap.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tenant API keys: " + err.Error()})
			return
		}
		tap.TenantID = tenantID
		if os.Getenv("TENANT_GET_API_KEY_REVEAL") != "true" {
			tap.Key = "***********************"
		}
		apikeys = append(apikeys, tap)
	}

	// This will return an empty array [] if no tenants are found.
	c.JSON(http.StatusOK, apikeys)
}
