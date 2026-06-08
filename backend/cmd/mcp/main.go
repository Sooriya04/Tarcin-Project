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
	// 1. Task Analytics
	toolTask := mcp.NewTool("get_task_analytics",
		mcp.WithDescription("Task completions, blockers, timelines."),
	)
	s.AddTool(toolTask, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetTaskAnalytics(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// 2. Mentor Analytics
	toolMentor := mcp.NewTool("get_mentor_analytics",
		mcp.WithDescription("Mentor workload, effectiveness, feedback quality."),
	)
	s.AddTool(toolMentor, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetMentorAnalytics(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// 3. Engagement Analytics
	toolEng := mcp.NewTool("get_engagement_analytics",
		mcp.WithDescription("Attendance, daily reports, drop-out risk."),
	)
	s.AddTool(toolEng, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetEngagementAnalytics(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// 4. Growth Analytics
	toolGrowth := mcp.NewTool("get_growth_analytics",
		mcp.WithDescription("Gamification, XP, leveling progress."),
	)
	s.AddTool(toolGrowth, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetGrowthAnalytics(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// 5. College Analytics
	toolCol := mcp.NewTool("get_college_analytics",
		mcp.WithDescription("College completion rates, institutional stats."),
	)
	s.AddTool(toolCol, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetCollegeAnalytics(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// 6. Workflow Analytics
	toolWf := mcp.NewTool("get_workflow_analytics",
		mcp.WithDescription("Escalations, negative trends, root causes."),
	)
	s.AddTool(toolWf, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetWorkflowAnalytics(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// 7. Conversion Analytics
	toolConv := mcp.NewTool("get_conversion_analytics",
		mcp.WithDescription("Hire-ready interns, conversion rates."),
	)
	s.AddTool(toolConv, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetConversionAnalytics(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// 8. Intern Health
	toolHealth := mcp.NewTool("get_intern_health",
		mcp.WithDescription("Active, completed, inactive headcount."),
	)
	s.AddTool(toolHealth, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetInternHealth(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// 9. Performance Metrics
	toolPerf := mcp.NewTool("get_performance_metrics",
		mcp.WithDescription("Top/lowest performers, evaluation scores."),
	)
	s.AddTool(toolPerf, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := client.GetPerformanceMetrics(ctx, &pb.Empty{})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(resp, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}
