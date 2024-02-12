package handlers

import (
	"database/sql"
	"net/http"

	"github.com/clodevo/raven-proxy/pkg/database"
	"github.com/clodevo/raven-proxy/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// Adjust the import path based on your project structure
)

// tenantsHandler to handle both collection and individual tenant actions
func TenantsHandler(c *gin.Context) {
	tenantIDStr := c.Param("tenantID")
	if tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
			return
		}

		// Add the tenantID to Gin context or handle it directly
		c.Set("tenantID", tenantID)

		// Based on the request method, call the specific function
		switch c.Request.Method {
		case "GET":
			getTenant(c)
		case "PUT":
			updateTenant(c)
		case "DELETE":
			deleteTenant(c)
		default:
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		}
	} else {
		// Handle collection actions
		switch c.Request.Method {
		case "GET":
			getAllTenants(c)
		case "POST":
			createTenant(c)
		default:
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		}
	}
}

// @Summary Create a new tenant
// @Description Create a new tenant with the provided name
// @Tags tenants
// @Accept json
// @Produce json
// @Param body body models.CreateTenantRequest true "Tenant creation request"
// @Success 200 {object} models.Tenant
// @Failure 400 {object} models.ErrorResponse
// @Router /tenants [post]
// @Security ApiKeyAuth
func createTenant(c *gin.Context) {
	var req models.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request: " + err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: name"})
		return
	}

	tenant_id := uuid.New()

	var exists bool = true

	for exists {
		tenant_id = uuid.New()
		// Check if the generated api_key_id exists in the database
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tenants WHERE tenant_id = ? )", tenant_id).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking tenant_id  existence: " + err.Error()})
			return
		}
	}

	stmt, err := database.DB.Prepare("INSERT INTO tenants (tenant_id,tenant_name) VALUES (?,?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(tenant_id, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	tenant := models.Tenant{
		ID:   tenant_id,
		Name: req.Name,
	}

	c.JSON(http.StatusOK, tenant)
}

// @Summary Get all tenants
// @Description Retrieve a list of all tenants
// @Tags tenants
// @Accept json
// @Produce json
// @Success 200 {array} models.Tenant
// @Failure 500 {object} models.ErrorResponse
// @Router /tenants [get]
// @Security ApiKeyAuth
func getAllTenants(c *gin.Context) {
	rows, err := database.DB.Query("SELECT tenant_id, tenant_name FROM tenants")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tenant: " + err.Error()})
		return
	}
	defer rows.Close()

	// Initialize tenants as an empty slice, not nil.
	tenants := make([]models.Tenant, 0)
	for rows.Next() {
		var t models.Tenant
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tenant: " + err.Error()})
			return
		}
		tenants = append(tenants, t)
	}

	// This will return an empty array [] if no tenants are found.
	c.JSON(http.StatusOK, tenants)
}

// @Summary Get a tenant by ID
// @Description Retrieve a tenant by its ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenantID path string true "Tenant ID"
// @Success 200 {object} models.Tenant
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /tenants/{tenantID} [get]
// @Security ApiKeyAuth
func getTenant(c *gin.Context) {
	tenantIDStr := c.Param("tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var tenant models.Tenant
	err = database.DB.QueryRow("SELECT tenant_id, tenant_name FROM tenants WHERE tenant_id = ?", tenantID).Scan(&tenant.ID, &tenant.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// @Summary Update a tenant by ID
// @Description Update a tenant's name by its ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenantID path string true "Tenant ID"
// @Param body body models.CreateTenantRequest true "Tenant object"
// @Success 200 {string} string "Tenant updated successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tenants/{tenantID} [put]
// @Security ApiKeyAuth
func updateTenant(c *gin.Context) {
	tenantIDStr := c.Param("tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req models.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request: " + err.Error()})
		return
	}

	_, err = database.DB.Exec("UPDATE tenants SET tenant_name = ? WHERE tenant_id = ?", req.Name, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant updated successfully"})
}

// @Summary Delete a tenant by ID
// @Description Delete a tenant by its ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenantID path string true "Tenant ID"
// @Success 200 {string} string "Tenant deleted successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tenants/{tenantID} [delete]
// @Security ApiKeyAuth
func deleteTenant(c *gin.Context) {
	tenantIDStr := c.Param("tenantID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM tenants WHERE tenant_id = ?", tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}
