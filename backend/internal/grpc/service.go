package grpc

import (
	"backend/internal/grpc/pb"
	"backend/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type PerformanceService struct {
	pb.UnimplementedPerformanceServiceServer
	DB   *sql.DB
	Repo *repository.InternRepository
}

func NewPerformanceService(db *sql.DB, repo *repository.InternRepository) *PerformanceService {
	return &PerformanceService{DB: db, Repo: repo}
}

func (s *PerformanceService) ExecuteSQLQuery(ctx context.Context, req *pb.SQLQueryRequest) (*pb.SQLQueryResponse, error) {
	query := req.GetQuery()
	log.Printf("[ExecuteSQLQuery] Incoming query: %s", query)

	// Basic check to prevent modifications
	trimmed := strings.TrimSpace(strings.ToUpper(query))
	if !strings.HasPrefix(trimmed, "SELECT") && !strings.HasPrefix(trimmed, "WITH") {
		log.Printf("[ExecuteSQLQuery] Security error: Query does not start with SELECT or WITH")
		return &pb.SQLQueryResponse{
			Error: "Security Error: Only read-only SELECT or WITH queries are allowed.",
		}, nil
	}

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[ExecuteSQLQuery] SQL error: %v for query: %s", err, query)
		return &pb.SQLQueryResponse{
			Error: fmt.Sprintf("SQL Error: %v", err),
		}, nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Printf("[ExecuteSQLQuery] Columns error: %v", err)
		return &pb.SQLQueryResponse{
			Error: fmt.Sprintf("SQL Columns Error: %v", err),
		}, nil
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("[ExecuteSQLQuery] Scan error: %v", err)
			return &pb.SQLQueryResponse{
				Error: fmt.Sprintf("SQL Scan Error: %v", err),
			}, nil
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}
		result = append(result, rowMap)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ExecuteSQLQuery] Rows error: %v", err)
		return &pb.SQLQueryResponse{
			Error: fmt.Sprintf("SQL Rows Iteration Error: %v", err),
		}, nil
	}

	if result == nil {
		result = []map[string]interface{}{}
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("[ExecuteSQLQuery] JSON Marshal error: %v", err)
		return &pb.SQLQueryResponse{
			Error: fmt.Sprintf("JSON Marshal Error: %v", err),
		}, nil
	}

	log.Printf("[ExecuteSQLQuery] Success: returned %d rows", len(result))
	return &pb.SQLQueryResponse{
		JsonResult: string(jsonData),
	}, nil
}

