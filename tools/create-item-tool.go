package tools

import (
	"context"
	"log"
	"os"

	"github.com/Epistemic-Technology/zotero/zotero"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CreateItemQuery struct {
	ItemType    string           `json:"item_type,omitempty" pattern:"^(book|bookSection|journalArticle|magazineArticle|conferencePaper|thesis|report|webpage|blogPost|preprint|manuscript|presentation)$"`
	Title       string           `json:"title,omitempty"`
	Creators    []zotero.Creator `json:"creators,omitempty"`
	Abstract    string           `json:"abstract,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Collections []string         `json:"collections,omitempty"`
	FilePath    string           `json:"file_path,omitempty"`
}

type CreateItemResult struct {
	ItemID       string `json:"item_id"`
	AttachmentID string `json:"attachment_id"`
}

func CreateItemTool() *mcp.Tool {
	inputschema, err := jsonschema.For[CreateItemQuery](nil)
	if err != nil {
		panic(err)
	}
	createItemTool := mcp.Tool{
		Name:        "zotero.create-item",
		Description: "Create a new item in Zotero",
		InputSchema: inputschema,
	}
	return &createItemTool
}

func CreateItemToolHandler(ctx context.Context, req *mcp.CallToolRequest, query CreateItemQuery) (*mcp.CallToolResult, *CreateItemResult, error) {
	libraryID := os.Getenv("ZOTERO_LIBRARY_ID")
	apiKey := os.Getenv("ZOTERO_API_KEY")
	tags := make([]zotero.Tag, len(query.Tags))
	for i, tag := range query.Tags {
		tags[i] = zotero.Tag{Tag: tag}
	}
	item := zotero.Item{
		Data: zotero.ItemData{
			Title:        query.Title,
			ItemType:     query.ItemType,
			Creators:     query.Creators,
			AbstractNote: query.Abstract,
			Tags:         tags,
			Collections:  query.Collections,
		},
	}
	createItemResult := CreateItemResult{}
	client := zotero.NewClient(libraryID, zotero.LibraryTypeUser, zotero.WithAPIKey(apiKey))
	resp, err := client.CreateItems(ctx, []zotero.Item{item})
	if err != nil {
		return nil, nil, err
	}
	var itemKey string
	if len(resp.Success) > 0 {
		for _, key := range resp.Success {
			if keyStr, ok := key.(string); ok {
				itemKey = keyStr
				log.Printf("Item created successfully with key: %s", itemKey)
				createItemResult.ItemID = itemKey
			}
		}
	}
	if len(resp.Failed) > 0 {
		log.Println("Failed items:")
		for idx, failure := range resp.Failed {
			log.Printf("Failed item at index %s: %d - %s", idx, failure.Code, failure.Message)
		}
	}

	// Upload file attachment if specified
	if query.FilePath != "" && itemKey != "" {
		attachment, err := client.UploadAttachment(ctx, itemKey, query.FilePath, "", "")
		if err != nil {
			log.Printf("Error uploading attachment: %v", err)
			log.Println("Note: Item was created successfully, but attachment upload failed")
		}
		log.Printf("Successfully attached file!")
		log.Printf("Attachment Key: %s", attachment.Key)
		log.Printf("Filename: %s", attachment.Data.Filename)
		createItemResult.AttachmentID = attachment.Key
	}
	return &mcp.CallToolResult{}, &createItemResult, nil
}
