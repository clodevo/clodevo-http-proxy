package utils

import (
	"encoding/base64"
	"strings"

	"github.com/clodevo/raven-proxy/pkg/database"
	"github.com/valyala/fasthttp"
)

var TenantNameRequest string

// Create a global instance of AuthCache
var authCache = NewAuthCache()

func Authenticate(ctx *fasthttp.RequestCtx) bool {
	auth := string(ctx.Request.Header.Peek("Proxy-Authorization"))
	if auth == "" {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Response.SetBodyString("Unauthorized: Authorization header required")
		return false
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Response.SetBodyString("Unauthorized: Invalid Authorization format")
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Response.SetBodyString("Unauthorized: Invalid Base64 encoding")
		return false
	}

	// Expecting the format to be "tenant_name:api_key"
	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Response.SetBodyString("Unauthorized: tenant_name and api_key required")
		return false
	}

	tenantName, apiKey := creds[0], creds[1]

	if valid, cachedTenantName := authCache.Check(tenantName, apiKey); valid {
		// Cache hit and not expired, consider authenticated
		TenantNameRequest = cachedTenantName
		return true
	}

	// Verify tenant_name and api_key against the database
	var dbTenantName string
	err = database.DB.QueryRow(`
        SELECT tenants.tenant_name 
        FROM api_keys 
        JOIN tenants ON api_keys.tenant_id = tenants.tenant_id 
        WHERE api_keys.api_key = ? AND tenants.tenant_name = ?`, apiKey, tenantName).Scan(&dbTenantName)

	if err != nil {
		// If the query fails, the API key or tenant_name is invalid
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Response.SetBodyString("Unauthorized: Invalid tenant_name or api_key")
		return false
	}

	if err == nil {
		// Update cache on successful authentication
		authCache.Update(tenantName, apiKey)
	}

	// If the API key and tenant_name are valid, optionally add the tenant_name to the response header
	TenantNameRequest = tenantName

	return true
}
