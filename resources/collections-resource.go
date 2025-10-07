package resources

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Epistemic-Technology/zotero/zotero"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var CollectionsResource = mcp.Resource{
	Name:        "zotero.collections",
	Description: "List of collections in my Zotero library",
	MIMEType:    "application/json",
	Title:       "My Zotero Collections",
	URI:         "file://zotero/collections.json",
}

func ZoteroListCollectionsResourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	libraryID := os.Getenv("ZOTERO_LIBRARY_ID")
	apiKey := os.Getenv("ZOTERO_API_KEY")
	client := zotero.NewClient(libraryID, zotero.LibraryTypeUser, zotero.WithAPIKey(apiKey))
	collections, err := client.Collections(ctx, nil)
	if err != nil {
		return nil, err
	}
	collectionsJSON, err := json.Marshal(collections)
	if err != nil {
		return nil, err
	}
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(collectionsJSON),
			},
		},
	}, nil
}
