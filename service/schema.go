package service

import (
	"Goal-Storage/utils"
	"errors"
	"github.com/graphql-go/graphql"
	"log"
)

// CreateGraphQLSchema defines and returns the GraphQL schema
func CreateGraphQLSchema() graphql.Schema {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"userId": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Extract the token from the context
					ctx := p.Context
					token, ok := ctx.Value("Authorization").(string)
					if !ok || token == "" {
						return nil, errors.New("no authorization token found")
					}
					// Fetch the userId using the token
					userId, err := utils.FetchUserIdFromAuthAPI(token)
					if err != nil {
						return nil, err
					}
					return userId, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}
	return schema
}
