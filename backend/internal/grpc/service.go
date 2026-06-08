package grpc

import (
	"backend/internal/grpc/pb"
	"backend/internal/repository"
	"context"
	"database/sql"
	"log"
)

type PerformanceService struct {
	pb.UnimplementedPerformanceServiceServer
	DB   *sql.DB
	Repo *repository.InternRepository
}

func NewPerformanceService(db *sql.DB, repo *repository.InternRepository) *PerformanceService {
	return &PerformanceService{DB: db, Repo: repo}
}

func (s *PerformanceService) GetPerformanceMetrics(ctx context.Context, req *pb.Empty) (*pb.PerformanceMetricsResponse, error) {
	resp := &pb.PerformanceMetricsResponse{
		DomainWisePerformance:  make(map[string]float32),
		CollegeWisePerformance: make(map[string]float32),
	}

	// 1. Average Evaluation Score
	err := s.DB.QueryRowContext(ctx, `SELECT COALESCE(AVG("totalScore"), 0) FROM "Evaluation"`).Scan(&resp.AverageEvaluationScore)
	if err != nil {
		log.Printf("Query error (AverageEvaluationScore): %v", err)
	}

	// 2. Top Performers
	topRows, err := s.DB.QueryContext(ctx, `
		SELECT e."internId", u.name, AVG(e."totalScore") as avg_score
		FROM "Evaluation" e
		JOIN "User" u ON e."internId" = u.id
		GROUP BY e."internId", u.name
		ORDER BY avg_score DESC
		LIMIT 5
	`)
	if err == nil {
		defer topRows.Close()
		for topRows.Next() {
			var id, name string
			var score float32
			if err := topRows.Scan(&id, &name, &score); err == nil {
				resp.TopPerformers = append(resp.TopPerformers, &pb.InternScore{Id: id, Name: name, Score: score})
			}
		}
	} else {
		log.Printf("Query error (TopPerformers): %v", err)
	}

	// 3. Lowest Performers
	lowRows, err := s.DB.QueryContext(ctx, `
		SELECT e."internId", u.name, AVG(e."totalScore") as avg_score
		FROM "Evaluation" e
		JOIN "User" u ON e."internId" = u.id
		GROUP BY e."internId", u.name
		ORDER BY avg_score ASC
		LIMIT 5
	`)
	if err == nil {
		defer lowRows.Close()
		for lowRows.Next() {
			var id, name string
			var score float32
			if err := lowRows.Scan(&id, &name, &score); err == nil {
				resp.LowestPerformers = append(resp.LowestPerformers, &pb.InternScore{Id: id, Name: name, Score: score})
			}
		}
	} else {
		log.Printf("Query error (LowestPerformers): %v", err)
	}

	// 4. Domain-wise Performance
	domainRows, err := s.DB.QueryContext(ctx, `
		SELECT ip."preferredDomain", AVG(e."totalScore")
		FROM "Evaluation" e
		JOIN "InternProfile" ip ON e."internId" = ip."userId"
		WHERE ip."preferredDomain" IS NOT NULL
		GROUP BY ip."preferredDomain"
	`)
	if err == nil {
		defer domainRows.Close()
		for domainRows.Next() {
			var domain string
			var score float32
			if err := domainRows.Scan(&domain, &score); err == nil {
				resp.DomainWisePerformance[domain] = score
			}
		}
	} else {
		log.Printf("Query error (DomainWisePerformance): %v", err)
	}

	// 5. College-wise Performance
	collegeRows, err := s.DB.QueryContext(ctx, `
		SELECT c.name, AVG(e."totalScore")
		FROM "Evaluation" e
		JOIN "InternProfile" ip ON e."internId" = ip."userId"
		JOIN "College" c ON ip."collegeId" = c.id
		GROUP BY c.name
	`)
	if err == nil {
		defer collegeRows.Close()
		for collegeRows.Next() {
			var college string
			var score float32
			if err := collegeRows.Scan(&college, &score); err == nil {
				resp.CollegeWisePerformance[college] = score
			}
		}
	} else {
		log.Printf("Query error (CollegeWisePerformance): %v", err)
	}

	return resp, nil
}

func (s *PerformanceService) GetInternHealth(ctx context.Context, req *pb.Empty) (*pb.InternHealthResponse, error) {
	metrics := s.Repo.GetHealthMetrics()

	resp := &pb.InternHealthResponse{
		ActiveInterns:         int32(metrics.ActiveInterns),
		CompletedInternships:  int32(metrics.CompletedInterns),
		InactiveInterns:       int32(metrics.InactiveInterns),
		TopColleges:           make(map[string]int32),
		FastestGrowingDomains: make(map[string]int32),
	}

	for k, v := range metrics.TopColleges {
		resp.TopColleges[k] = int32(v)
	}

	for k, v := range metrics.FastestDomains {
		resp.FastestGrowingDomains[k] = int32(v)
	}

	return resp, nil
}
