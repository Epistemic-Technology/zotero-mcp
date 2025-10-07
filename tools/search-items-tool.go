package tools

import (
	"context"
	"os"

	"github.com/Epistemic-Technology/zotero/zotero"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchQuery struct {
	Query        string   `json:"q" jsonschema:"Search text"`
	QueryMode    string   `json:"qmode" jsonschema:"What fields to search. Use 'everything' for full-text search." pattern:"^(titleCreatorYear|everything)$"`
	ItemType     []string `json:"itemType" jsonschema:"Filter to this type of item. Use '-' prefix to exclude this type. Examples: 'book', '-attachment', 'journalArticle'"`
	Tags         []string `json:"tags" jsonschema:"Filter to these tags. Use '-' prefix to exclude this tag."`
	CollectionID string   `json:"collectionId" jsonschema:"Filter to this collection ID"`
	Limit        int      `json:"limit" jsonschema:"Maximum number of results to return (default: 100)"`
}

type SearchResult struct {
	Items *[]ItemSummaryData `json:"items"`
}

type ItemSummaryData struct {
	Key      string           `json:"key,omitempty" jsonschema:"Zotero item ID"`
	ItemType string           `json:"itemType,omitempty" jsonschema:"Zotero item type"`
	Title    string           `json:"title,omitempty" jsonschema:"Title of the item"`
	Creators []zotero.Creator `json:"creators,omitempty" jsonschema:"Creators of the item"`
	Abstract string           `json:"abstractNote,omitempty" jsonschema:"Abstract of the item"`
}

func SearchTool() *mcp.Tool {
	inputschema, err := jsonschema.For[SearchQuery](nil)
	if err != nil {
		panic(err)
	}
	searchTool := mcp.Tool{
		Name:        "zotero.search",
		Description: "Search Zotero library",
		InputSchema: inputschema,
	}
	return &searchTool
}

func SearchToolHandler(ctx context.Context, req *mcp.CallToolRequest, query SearchQuery) (*mcp.CallToolResult, *SearchResult, error) {
	libraryID := os.Getenv("ZOTERO_LIBRARY_ID")
	apiKey := os.Getenv("ZOTERO_API_KEY")
	client := zotero.NewClient(libraryID, zotero.LibraryTypeUser, zotero.WithAPIKey(apiKey))
	zoteroParams := zotero.QueryParams{
		ItemType: query.ItemType,
		Tag:      query.Tags,
		Limit:    query.Limit,
		Q:        query.Query,
		QMode:    query.QueryMode,
	}
	var items []zotero.Item
	var err error
	if query.CollectionID != "" {
		items, err = client.CollectionItems(ctx, query.CollectionID, &zoteroParams)
	} else {
		items, err = client.Items(ctx, &zoteroParams)
	}
	if err != nil {
		return nil, nil, err
	}
	var itemSummaries []ItemSummaryData
	for _, item := range items {
		itemSummary := ItemSummaryData{
			Key:      item.Data.Key,
			ItemType: item.Data.ItemType,
			Title:    item.Data.Title,
			Creators: item.Data.Creators,
			Abstract: item.Data.AbstractNote,
		}
		itemSummaries = append(itemSummaries, itemSummary)
	}
	result := SearchResult{
		Items: &itemSummaries,
	}
	return &mcp.CallToolResult{}, &result, nil
}
