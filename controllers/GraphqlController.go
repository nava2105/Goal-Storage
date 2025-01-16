package controllers

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"net/http"
)

// HandleGraphQL handles the GraphQL requests.
func HandleGraphQL(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse GraphQL query from request.
		var requestBody struct {
			Query string `json:"query"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Execute query on schema.
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: requestBody.Query,
			Context:       r.Context(),
		})

		// Handle errors.
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": result.Errors,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result.Data)
	}
}
