package tools

import (
	"context"
	"os"

	"github.com/Epistemic-Technology/zotero/zotero"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CreateCollectionQuery struct {
	Name string `json:"name"`
}

type CreateCollectionResult struct {
	ID string `json:"id"`
}

func CreateCollectionTool() *mcp.Tool {
	inputschema, err := jsonschema.For[CreateCollectionQuery](nil)
	if err != nil {
		panic(err)
	}
	return &mcp.Tool{
		Name:        "CreateCollection",
		Description: "Create a new collection in Zotero",
		InputSchema: inputschema,
	}
}

func CreateCollectionToolHandler(ctx context.Context, req *mcp.CallToolRequest, query CreateCollectionQuery) (*mcp.CallToolResult, *CreateCollectionResult, error) {
	libraryID := os.Getenv("ZOTERO_LIBRARY_ID")
	apiKey := os.Getenv("ZOTERO_API_KEY")
	client := zotero.NewClient(libraryID, zotero.LibraryTypeUser, zotero.WithAPIKey(apiKey))
	collection := zotero.Collection{
		Data: zotero.CollectionData{
			Name: query.Name,
		},
	}
	collectionID := ""
	response, err := client.CreateCollections(ctx, []zotero.Collection{collection})
	if err != nil {
		return nil, nil, err
	}
	if len(response.Success) > 0 {
		for _, key := range response.Success {
			if keyStr, ok := key.(string); ok {
				collectionID = keyStr
			}
		}
	}

	return &mcp.CallToolResult{}, &CreateCollectionResult{ID: collectionID}, nil
}
