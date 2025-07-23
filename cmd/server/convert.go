// this file is used in web/api/convert.go again
// in server sub-dir is a norman web handler
// in web/api is a serverless function
// but with different package main

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"

	"github.com/Sn0rt/secret2es/pkg/converter"
)

type ConvertRequest struct {
	Content         string            `json:"content"`
	StoreType       string            `json:"storeType"`
	StoreName       string            `json:"storeName"`
	CreationPolicy  string            `json:"creationPolicy"`
	RefreshPolicy   string            `json:"refreshPolicy"`
	RefreshInterval string            `json:"refreshInterval"`
	Resolve         bool              `json:"resolve"`
	EnvVars         map[string]string `json:"envVars,omitempty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var request ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var missingFields []string
	if request.Content == "" {
		missingFields = append(missingFields, "content")
	}
	if request.StoreType == "" {
		missingFields = append(missingFields, "storeType")
	}
	if request.StoreName == "" {
		missingFields = append(missingFields, "storeName")
	}
	if request.CreationPolicy == "" {
		missingFields = append(missingFields, "creationPolicy")
	}

	if len(missingFields) > 0 {
		errorResponse := map[string]string{
			"error": fmt.Sprintf("Missing required fields: %s", strings.Join(missingFields, ", ")),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if request.Resolve && len(request.EnvVars) == 0 {
		errorResponse := map[string]string{
			"error": "Resolve is set to true but no environment variables provided",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	refreshPolicy := esv1.ExternalSecretRefreshPolicy(request.RefreshPolicy)
	result, warn, err := converter.ConvertSecretContent(
		[]byte(request.Content),
		request.StoreType,
		request.StoreName,
		esv1.ExternalSecretCreationPolicy(request.CreationPolicy),
		request.Resolve,
		refreshPolicy,
		request.RefreshInterval,
		request.EnvVars,
	)

	if err != nil {
		errorResponse := map[string]string{
			"error": "Conversion error: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := map[string]interface{}{
		"result": result,
	}

	if warn != "" {
		response["warnings"] = warn
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
