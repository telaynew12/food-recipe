package chapa_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yourproject/chapa" // Replace with the correct import path
)

func TestHandleChapaPayment(t *testing.T) {
	// Mock request payload
	mockRequest := chapa.Request{
		Amount:      100.0,
		Currency:    "ETB",
		Email:       "testuser@example.com",
		FirstName:   "Test",
		LastName:    "User",
		TxRef:       "test-tx-ref",
		CallbackURL: "https://localhost:8085/api/payment/callback",
	}

	// Encode the mock request into JSON
	requestBody, err := json.Marshal(mockRequest)
	if err != nil {
		t.Fatalf("Failed to encode request: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest("POST", "/process-chapa-payment", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the output
	recorder := httptest.NewRecorder()

	// Call the handler
	chapa.HandlePayment(recorder, req)

	// Check the status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", recorder.Code)
	}

	// Parse the response body
	var response chapa.Response
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify the response content
	if response.Status != "success" {
		t.Errorf("Expected success status, got %s", response.Status)
	}

	if response.CheckoutURL == "" {
		t.Error("Expected a valid checkout URL, got an empty string")
	}
}
