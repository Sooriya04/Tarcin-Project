package grpc

import (
	"backend/internal/grpc/pb"
	"context"
)

func (s *PerformanceService) GetCollegeAnalytics(ctx context.Context, req *pb.Empty) (*pb.CollegeAnalyticsResponse, error) {
	resp := &pb.CollegeAnalyticsResponse{}

	// 1. Highest number of interns
	q1 := `
		SELECT c.name, CAST(COUNT(ip.id) AS FLOAT) as intern_count
		FROM "College" c
		JOIN "InternProfile" ip ON c.id = ip."collegeId"
		GROUP BY c.name
		ORDER BY intern_count DESC
		LIMIT 5
	`
	rows1, err := s.DB.QueryContext(ctx, q1)
	if err == nil {
		defer rows1.Close()
		for rows1.Next() {
			var name string
			var val float32
			if err := rows1.Scan(&name, &val); err == nil {
				resp.HighestInternCount = append(resp.HighestInternCount, &pb.CollegeStat{CollegeName: name, StatValue: val})
			}
		}
	}

	// 2. Best performing (Average Evaluation score)
	q2 := `
		SELECT c.name, CAST(AVG(e."totalScore") AS FLOAT) as avg_score
		FROM "College" c
		JOIN "InternProfile" ip ON c.id = ip."collegeId"
		JOIN "Evaluation" e ON ip."userId" = e."internId"
		GROUP BY c.name
		ORDER BY avg_score DESC
		LIMIT 5
	`
	rows2, err := s.DB.QueryContext(ctx, q2)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var name string
			var val float32
			if err := rows2.Scan(&name, &val); err == nil {
				resp.BestPerforming = append(resp.BestPerforming, &pb.CollegeStat{CollegeName: name, StatValue: val})
			}
		}
	}

	// 3. Highest completion rates
	q3 := `
		SELECT c.name, 
			COALESCE(CAST(SUM(CASE WHEN ip."conversionStatus" IN ('COMPLETED', 'CONVERTED') THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(ip.id), 0) * 100, 0) as completion_rate
		FROM "College" c
		JOIN "InternProfile" ip ON c.id = ip."collegeId"
		GROUP BY c.name
		ORDER BY completion_rate DESC
		LIMIT 5
	`
	rows3, err := s.DB.QueryContext(ctx, q3)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var name string
			var val float32
			if err := rows3.Scan(&name, &val); err == nil {
				resp.HighestCompletionRate = append(resp.HighestCompletionRate, &pb.CollegeStat{CollegeName: name, StatValue: val})
			}
		}
	}

	// 4. Highest dropout risk
	q4 := `
		SELECT c.name, CAST(AVG(w."escalationLevel") AS FLOAT) as risk_level
		FROM "College" c
		JOIN "InternProfile" ip ON c.id = ip."collegeId"
		JOIN "WeeklySnapshot" w ON ip."userId" = w."userId"
		GROUP BY c.name
		ORDER BY risk_level DESC
		LIMIT 5
	`
	rows4, err := s.DB.QueryContext(ctx, q4)
	if err == nil {
		defer rows4.Close()
		for rows4.Next() {
			var name string
			var val float32
			if err := rows4.Scan(&name, &val); err == nil {
				resp.HighestDropoutRisk = append(resp.HighestDropoutRisk, &pb.CollegeStat{CollegeName: name, StatValue: val})
			}
		}
	}

	// 5. Compare attendance trends (Avg attendances per intern)
	q5 := `
		SELECT c.name, CAST(COUNT(a.id) AS FLOAT) / NULLIF(COUNT(DISTINCT ip."userId"), 0) as avg_attendances
		FROM "College" c
		JOIN "InternProfile" ip ON c.id = ip."collegeId"
		LEFT JOIN "Attendance" a ON ip."userId" = a."userId"
		GROUP BY c.name
		ORDER BY avg_attendances DESC
		LIMIT 5
	`
	rows5, err := s.DB.QueryContext(ctx, q5)
	if err == nil {
		defer rows5.Close()
		for rows5.Next() {
			var name string
			var val float32
			if err := rows5.Scan(&name, &val); err == nil {
				resp.AttendanceTrends = append(resp.AttendanceTrends, &pb.CollegeStat{CollegeName: name, StatValue: val})
			}
		}
	}

	// 6. Should receive more slots (Highest Avg XP)
	q6 := `
		SELECT c.name, CAST(AVG(u.xp) AS FLOAT) as performance_score
		FROM "College" c
		JOIN "InternProfile" ip ON c.id = ip."collegeId"
		JOIN "User" u ON ip."userId" = u.id
		GROUP BY c.name
		ORDER BY performance_score DESC
		LIMIT 5
	`
	rows6, err := s.DB.QueryContext(ctx, q6)
	if err == nil {
		defer rows6.Close()
		for rows6.Next() {
			var name string
			var val float32
			if err := rows6.Scan(&name, &val); err == nil {
				resp.RecommendedForMoreSlots = append(resp.RecommendedForMoreSlots, &pb.CollegeStat{CollegeName: name, StatValue: val})
			}
		}
	}

	// 7. Declining performance (Avg consecutiveNegativeWeeks)
	q7 := `
		SELECT c.name, CAST(AVG(w."consecutiveNegativeWeeks") AS FLOAT) as decline_severity
		FROM "College" c
		JOIN "InternProfile" ip ON c.id = ip."collegeId"
		JOIN "WeeklySnapshot" w ON ip."userId" = w."userId"
		GROUP BY c.name
		ORDER BY decline_severity DESC
		LIMIT 5
	`
	rows7, err := s.DB.QueryContext(ctx, q7)
	if err == nil {
		defer rows7.Close()
		for rows7.Next() {
			var name string
			var val float32
			if err := rows7.Scan(&name, &val); err == nil {
				resp.DecliningPerformance = append(resp.DecliningPerformance, &pb.CollegeStat{CollegeName: name, StatValue: val})
			}
		}
	}

	return resp, nil
}
