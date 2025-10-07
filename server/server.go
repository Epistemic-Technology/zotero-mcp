package server

import (
	"github.com/Epistemic-Technology/zotero-mcp/resources"
	"github.com/Epistemic-Technology/zotero-mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func CreateServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "zotero-mcp", Version: "v0.0.1"}, nil)
	mcp.AddTool(server, tools.SearchTool(), tools.SearchToolHandler)
	server.AddResource(&resources.CollectionsResource, resources.ZoteroListCollectionsResourceHandler)
	return server
}
