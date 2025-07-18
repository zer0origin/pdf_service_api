package integration

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/testutil"
	"strings"
	"testing"
)

func TestSelectionsIntegration(t *testing.T) {
	t.Run("getSelections", getSelections)
	t.Run("deleteSelections", deleteSelections)
	t.Run("deleteSelectionsUuidDoesNotExist", deleteSelectionsUuidDoesNotExist)
	t.Run("createNewSelection", createNewSelection)
}

// /api/v1/documents/:id/selections/
func getSelections(t *testing.T) {
	expectedJsonResponse := `[{"selectionUUID":"a5fdea38-0a86-4c19-ae4f-c87a01bc860d","documentID":"b66fd223-515f-4503-80cc-2bdaa50ef474"},{"selectionUUID":"335a6b95-6707-4e2b-9c37-c76d017f6f97","documentID":"b66fd223-515f-4503-80cc-2bdaa50ef474"}]`
	t.Parallel()
	router := testutil.CreateV1RouterAndPostgresContainer(t, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/documents/%s/selections/", "b66fd223-515f-4503-80cc-2bdaa50ef474"),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotNil(t, w.Body.String())
	assert.NotContains(t, w.Body.String(), "Error")
	assert.Equal(t, expectedJsonResponse, w.Body.String())
}

// /api/v1/documents/:id/selections/
func deleteSelections(t *testing.T) {
	t.Parallel()
	router := testutil.CreateV1RouterAndPostgresContainer(t, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/selections/%s", "a5fdea38-0a86-4c19-ae4f-c87a01bc860d"),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

// /api/v1/documents/:id/selections/
func deleteSelectionsUuidDoesNotExist(t *testing.T) {
	t.Parallel()
	router := testutil.CreateV1RouterAndPostgresContainer(t, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/selections/%s", uuid.New().String()),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func createNewSelection(t *testing.T) {
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	t.Parallel()
	router := testutil.CreateV1RouterAndPostgresContainer(t, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)

	request := &v1.AddNewSelectionRequest{
		DocumentID:      func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		IsComplete:      false,
		Settings:        nil,
		SelectionBounds: nil,
	}

	requestJSON, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/selections/",
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotContains(t, w.Body.String(), "Error")
}
