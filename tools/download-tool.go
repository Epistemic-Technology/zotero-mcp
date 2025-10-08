package tools

import (
	"context"
	"os"

	"github.com/Epistemic-Technology/zotero/zotero"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type DownloadQuery struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type DownloadResult struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

func DownloadTool() *mcp.Tool {
	inputschema, err := jsonschema.For[DownloadQuery](nil)
	if err != nil {
		panic(err)
	}
	downloadTool := mcp.Tool{
		Name:        "zotero.download",
		Description: "Download a Zotero file attachment",
		InputSchema: inputschema,
	}
	return &downloadTool
}

func DownloadToolHandler(ctx context.Context, req *mcp.CallToolRequest, query DownloadQuery) (*mcp.CallToolResult, *DownloadResult, error) {
	libraryID := os.Getenv("ZOTERO_LIBRARY_ID")
	apiKey := os.Getenv("ZOTERO_API_KEY")
	client := zotero.NewClient(libraryID, zotero.LibraryTypeUser, zotero.WithAPIKey(apiKey))
	downloadedFilePath, err := client.Dump(ctx, query.ID, query.Path, "")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{}, &DownloadResult{
		ID:   query.ID,
		Path: downloadedFilePath,
	}, nil
}
