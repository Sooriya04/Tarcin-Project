package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"backend/internal/grpc/pb"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to local gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to gRPC server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewPerformanceServiceClient(conn)

	// Create MCP server
	s := server.NewMCPServer(
		"TarcinBridgeMCP",
		"1.0.0",
	)

	// Register tools
	registerTools(s, client)

	// Start stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "MCP Server error: %v\n", err)
	}
}

func registerTools(s *server.MCPServer, client pb.PerformanceServiceClient) {
	// 1. Execute SQL Query
	toolSQL := mcp.NewTool("execute_sql_query",
		mcp.WithDescription("Run a read-only SQL query against the Postgres database to retrieve custom metrics, counts, averages, and profiles."),
		mcp.WithString("query", mcp.Description("The SELECT SQL query to execute. Must be read-only."), mcp.Required()),
	)
	s.AddTool(toolSQL, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args struct {
			Query string `json:"query"`
		}
		argBytes, _ := json.Marshal(request.Params.Arguments)
		json.Unmarshal(argBytes, &args)

		resp, err := client.ExecuteSQLQuery(ctx, &pb.SQLQueryRequest{Query: args.Query})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if resp.Error != "" {
			return mcp.NewToolResultError(resp.Error), nil
		}
		return mcp.NewToolResultText(resp.JsonResult), nil
	})
}

