package grpc

import (
	"backend/internal/grpc/pb"
	"context"
)

func (s *PerformanceService) GetConversionAnalytics(ctx context.Context, req *pb.Empty) (*pb.ConversionAnalyticsResponse, error) {
	resp := &pb.ConversionAnalyticsResponse{}

	// 1. Not started approaching end
	q1 := `
		SELECT u.name, 1.0, ip."conversionStatus"
		FROM "User" u
		JOIN "InternProfile" ip ON u.id = ip."userId"
		WHERE ip."conversionStatus" = 'NOT_STARTED' AND ip."endDate"::timestamp < NOW() + INTERVAL '14 days'
		LIMIT 5
	`
	rows1, err := s.DB.QueryContext(ctx, q1)
	if err == nil {
		defer rows1.Close()
		for rows1.Next() {
			var name, info string
			var val float32
			if err := rows1.Scan(&name, &val, &info); err == nil {
				resp.NotStartedApproachingEnd = append(resp.NotStartedApproachingEnd, &pb.ConversionStat{InternName: name, StatValue: val, ExtraInfo: info})
			}
		}
	}

	// 2. High eval ready hire
	q2 := `
		SELECT u.name, CAST(AVG(e."totalScore") AS FLOAT) as avg_score, ip."conversionStatus"
		FROM "User" u
		JOIN "InternProfile" ip ON u.id = ip."userId"
		JOIN "Evaluation" e ON u.id = e."internId"
		GROUP BY u.id, u.name, ip."conversionStatus"
		HAVING AVG(e."totalScore") > 80 AND ip."conversionStatus" = 'CONVERTED'
		LIMIT 5
	`
	rows2, err := s.DB.QueryContext(ctx, q2)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var name, info string
			var val float32
			if err := rows2.Scan(&name, &val, &info); err == nil {
				resp.HighEvalReadyHire = append(resp.HighEvalReadyHire, &pb.ConversionStat{InternName: name, StatValue: val, ExtraInfo: info})
			}
		}
	}

	// 3. Top conversion colleges
	q3 := `
		SELECT c.name, CAST(SUM(CASE WHEN ip."conversionStatus" = 'CONVERTED' THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(ip.id), 0) * 100 as rate, 'Conversion Rate'
		FROM "College" c
		JOIN "InternProfile" ip ON c.id = ip."collegeId"
		GROUP BY c.name
		ORDER BY rate DESC
		LIMIT 5
	`
	rows3, err := s.DB.QueryContext(ctx, q3)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var name, info string
			var val float32
			if err := rows3.Scan(&name, &val, &info); err == nil {
				resp.TopConversionColleges = append(resp.TopConversionColleges, &pb.ConversionStat{InternName: name, StatValue: val, ExtraInfo: info})
			}
		}
	}

	return resp, nil
}
