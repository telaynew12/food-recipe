package image

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type FileUploadRequest struct {
	Input struct {
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Base64Str string `json:"base64str" binding:"required"`
		RecipeID  string `json:"RecipeID" binding:"required"` // Changed RecipeID to string
	} `json:"input"`
}

type FileUploadResponse struct {
	FilePath string `json:"file_path"`
	RecipeID string `json:"RecipeID"` // Changed RecipeID to string
}

// Update the featured_image column for the specified recipeID
func updateFeaturedImageInDatabase(db *sql.DB, filePath string, RecipeID string) error {
	query := `UPDATE recipes SET featured_image = $1 WHERE id = $2`
	_, err := db.Exec(query, filePath, RecipeID)
	if err != nil {
		return fmt.Errorf("failed to update featured image: %w", err)
	}
	return nil
}

func FileUploadHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body
	var req FileUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Decode Base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(req.Input.Base64Str)
	if err != nil {
		http.Error(w, "Failed to decode base64 string", http.StatusBadRequest)
		return
	}

	// Save the file
	filePath := fmt.Sprintf("./uploads/%s", req.Input.Name)
	err = ioutil.WriteFile(filePath, decodedBytes, 0644)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Update the featured image URL in the recipes table for the specified recipeID
	// Store the URL path instead of the local file path
	imageURL := fmt.Sprintf("/uploads/%s", req.Input.Name)
	err = updateFeaturedImageInDatabase(db, imageURL, req.Input.RecipeID)
	if err != nil {
		http.Error(w, "Failed to update file path in database", http.StatusInternalServerError)
		return
	}

	// Respond with the file path and recipeID
	response := FileUploadResponse{
		FilePath: imageURL,
		RecipeID: req.Input.RecipeID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func init() {
	// Create uploads directory if not exists
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("./uploads", 0755); err != nil {
			log.Fatalf("Failed to create uploads directory: %v", err)
		}
	}
}
