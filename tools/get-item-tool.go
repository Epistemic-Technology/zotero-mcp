package tools

import (
	"context"
	"os"

	"github.com/Epistemic-Technology/zotero/zotero"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetItemQuery struct {
	ID string `json:"id"`
}

type GetItemResult struct {
	Item     *zotero.Item `json:"item"`
	Children *[]ItemChild `json:"children,omitempty"`
}

type ItemChild struct {
	ID    string      `json:"id,omitempty"`
	Type  string      `json:"type,omitempty"`
	Title string      `json:"title,omitempty"`
	Link  zotero.Link `json:"link,omitempty"`
}

func GetItemTool() *mcp.Tool {
	inputschema, err := jsonschema.For[GetItemQuery](nil)
	if err != nil {
		panic(err)
	}
	getTool := mcp.Tool{
		Name:        "zotero.get-item",
		Description: "Get a Zotero item by ID",
		InputSchema: inputschema,
	}
	return &getTool
}

func GetItemToolHandler(ctx context.Context, req *mcp.CallToolRequest, query GetItemQuery) (*mcp.CallToolResult, *GetItemResult, error) {
	libraryID := os.Getenv("ZOTERO_LIBRARY_ID")
	apiKey := os.Getenv("ZOTERO_API_KEY")
	client := zotero.NewClient(libraryID, zotero.LibraryTypeUser, zotero.WithAPIKey(apiKey))
	item, err := client.Item(ctx, query.ID, nil)
	if err != nil {
		return nil, nil, err
	}
	children, err := client.Children(ctx, query.ID, nil)
	if err != nil {
		return nil, nil, err
	}
	childrenArray := make([]ItemChild, len(children))
	for i, child := range children {
		childrenArray[i] = ItemChild{
			ID:    child.Data.Key,
			Type:  child.Data.ItemType,
			Title: child.Data.Title,
		}
		if child.Links.Self.Href != "" {
			childrenArray[i].Link = child.Links.Self
		}
	}

	return &mcp.CallToolResult{}, &GetItemResult{Item: item, Children: &childrenArray}, nil
}
