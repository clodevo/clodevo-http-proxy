package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"path/filepath"

	"github.com/google/uuid"
)

// GenerateRandomString creates a random string of a specified length using a given charset.
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatal(err)
		}
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

// // AuthenticateAdminInit checks for an existing admin API key in the database or environment variable,
// // and generates a new one if none is found, storing it in the database.
// func AuthenticateAdminInit(db *sql.DB, adminAPIKeyEnvVar string, adminAPIKey *string) {
// 	var err error

// 	// Check if an API key exists in the environment variable
// 	if adminAPIKeyEnv := os.Getenv(adminAPIKeyEnvVar); adminAPIKeyEnv != "" {
// 		*adminAPIKey = adminAPIKeyEnv
// 		log.Printf("Using admin API key from environment variable: %s", *adminAPIKey)
// 	} else {
// 		// If not in environment, check the database
// 		err = db.QueryRow("SELECT apikey FROM admin LIMIT 1").Scan(adminAPIKey)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				// No API key found in the database, generate a new one
// 				*adminAPIKey = GenerateRandomString(32)
// 				log.Printf("Generated default admin API key: %s", *adminAPIKey)

// 				// Store the new API key in the database
// 				_, err = db.Exec("INSERT INTO admin (apikey) VALUES (?)", *adminAPIKey)
// 				if err != nil {
// 					log.Fatalf("Failed to store default admin API key in database: %v", err)
// 				}
// 			} else {
// 				// An error occurred while querying the database
// 				log.Fatalf("Error querying admin API key from database: %v", err)
// 			}
// 		} else {
// 			log.Printf("Using admin API key from database: %s", *adminAPIKey)
// 		}
// 	}
// }

func InitializeDefaultTenant(db *sql.DB, aclDataPath string) {
	err := os.MkdirAll(aclDataPath, 0755)
	if err != nil {
		log.Fatalf("Failed to create aclDataPath directory: %v", err)
	}

	defaultTenantPath := filepath.Join(aclDataPath, "default.json")
	if _, err := os.Stat(defaultTenantPath); os.IsNotExist(err) {
		defaultList := struct {
			Whitelist []string `json:"Whitelist"`
			Blacklist []string `json:"Blacklist"`
		}{
			Whitelist: []string{"*.example.com"},
			Blacklist: []string{"restricted.example.com"},
		}
		fileContent, err := json.MarshalIndent(defaultList, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal default tenant list: %v", err)
		}

		err = ioutil.WriteFile(defaultTenantPath, fileContent, 0644)
		if err != nil {
			log.Fatalf("Failed to write default tenant list file: %v", err)
		}
		log.Println("Default tenant configuration created successfully.")
	} else if err != nil {
		log.Fatalf("Error checking for default tenant configuration: %v", err)
	} else {
		log.Println("Default tenant configuration already exists.")
	}

	// New logic to ensure a default tenant exists in the database
	const defaultTenantName = "default"
	var tenantID uuid.UUID
	err = db.QueryRow("SELECT tenant_id FROM tenants WHERE tenant_name = ?", defaultTenantName).Scan(&tenantID)
	if err == sql.ErrNoRows {
		// Default tenant does not exist, create it
		tenantID = uuid.New()
		_, err := db.Exec("INSERT INTO tenants (tenant_id, tenant_name) VALUES (?, ?)", tenantID, defaultTenantName)
		if err != nil {
			log.Fatalf("Failed to insert default tenant into database: %v", err)
		}
		log.Println("Default tenant created successfully in the database.")
	} else if err != nil {
		log.Fatalf("Error checking for default tenant in database: %v", err)
	} else {
		log.Println("Default tenant already exists in the database.")
	}
}
