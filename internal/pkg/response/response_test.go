package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"warehouse/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"name": "test"}
	Success(c, data)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != errors.CodeSuccess {
		t.Errorf("Code = %d, want %d", resp.Code, errors.CodeSuccess)
	}
	if resp.Message != "success" {
		t.Errorf("Message = %s, want %s", resp.Message, "success")
	}
	if resp.Data == nil {
		t.Error("Data should not be nil")
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		message    string
		statusCode int
	}{
		{"BadRequest", errors.CodeBadRequest, "invalid input", http.StatusBadRequest},
		{"Unauthorized", errors.CodeUnauthorized, "not authenticated", http.StatusUnauthorized},
		{"Forbidden", errors.CodeForbidden, "access denied", http.StatusForbidden},
		{"NotFound", errors.CodeNotFound, "resource not found", http.StatusNotFound},
		{"InternalError", errors.CodeInternalError, "server error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, tt.code, tt.message)

			if w.Code != tt.statusCode {
				t.Errorf("Status = %d, want %d", w.Code, tt.statusCode)
			}

			var resp Response
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if resp.Code != tt.code {
				t.Errorf("Code = %d, want %d", resp.Code, tt.code)
			}
			if resp.Message != tt.message {
				t.Errorf("Message = %s, want %s", resp.Message, tt.message)
			}
		})
	}
}

func TestSuccessWithPage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []map[string]string{
		{"id": "1", "name": "item1"},
		{"id": "2", "name": "item2"},
	}

	SuccessWithPage(c, items, 100, 1, 10)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != errors.CodeSuccess {
		t.Errorf("Code = %d, want %d", resp.Code, errors.CodeSuccess)
	}

	pageData, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Data is not a map")
	}

	if pageData["total"].(float64) != 100 {
		t.Errorf("Total = %v, want 100", pageData["total"])
	}
	if pageData["page"].(float64) != 1 {
		t.Errorf("Page = %v, want 1", pageData["page"])
	}
	if pageData["page_size"].(float64) != 10 {
		t.Errorf("PageSize = %v, want 10", pageData["page_size"])
	}
}

func TestSuccessWithNilData(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != errors.CodeSuccess {
		t.Errorf("Code = %d, want %d", resp.Code, errors.CodeSuccess)
	}
}

func TestSuccessWithEmptySlice(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []string{}
	Success(c, items)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}
