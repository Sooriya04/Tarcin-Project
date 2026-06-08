package grpc

import (
	"backend/internal/grpc/pb"
	"context"
)

func (s *PerformanceService) GetMentorAnalytics(ctx context.Context, req *pb.Empty) (*pb.MentorAnalyticsResponse, error) {
	resp := &pb.MentorAnalyticsResponse{}

	// 1. Mentor with largest number of interns
	q1 := `
		SELECT u.name, CAST(COUNT(ip.id) AS FLOAT) as intern_count
		FROM "InternProfile" ip
		JOIN "User" u ON ip."assignedMentorId" = u.id
		GROUP BY u.name
		ORDER BY intern_count DESC
		LIMIT 5
	`
	rows1, err := s.DB.QueryContext(ctx, q1)
	if err == nil {
		defer rows1.Close()
		for rows1.Next() {
			var name string
			var count float32
			if err := rows1.Scan(&name, &count); err == nil {
				resp.LargestInternCount = append(resp.LargestInternCount, &pb.MentorStat{MentorName: name, StatValue: count})
			}
		}
	}

	// 2. Mentor's interns with highest average evaluation score
	q2 := `
		SELECT u.name, AVG(e."totalScore") as avg_score
		FROM "InternProfile" ip
		JOIN "Evaluation" e ON ip."userId" = e."internId"
		JOIN "User" u ON ip."assignedMentorId" = u.id
		GROUP BY u.name
		ORDER BY avg_score DESC
		LIMIT 5
	`
	rows2, err := s.DB.QueryContext(ctx, q2)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var name string
			var score float32
			if err := rows2.Scan(&name, &score); err == nil {
				resp.HighestAvgEvaluation = append(resp.HighestAvgEvaluation, &pb.MentorStat{MentorName: name, StatValue: score})
			}
		}
	}

	// 3. Lowest attendance rate (Using average attendances per intern)
	q3 := `
		SELECT u.name, CAST(COUNT(a.id) AS FLOAT) / NULLIF(COUNT(DISTINCT ip."userId"), 0) as avg_attendances
		FROM "InternProfile" ip
		JOIN "User" u ON ip."assignedMentorId" = u.id
		LEFT JOIN "Attendance" a ON ip."userId" = a."userId"
		GROUP BY u.name
		ORDER BY avg_attendances ASC
		LIMIT 5
	`
	rows3, err := s.DB.QueryContext(ctx, q3)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var name string
			var rate float32
			if err := rows3.Scan(&name, &rate); err == nil {
				resp.LowestAttendanceRate = append(resp.LowestAttendanceRate, &pb.MentorStat{MentorName: name, StatValue: rate})
			}
		}
	}

	// 4. Highest task approval rates
	q4 := `
		SELECT u.name, 
			(CAST(SUM(CASE WHEN t.status = 'APPROVED' THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(t.id), 0)) * 100 as approval_rate
		FROM "InternProfile" ip
		JOIN "Task" t ON ip."userId" = t."assignedToId"
		JOIN "User" u ON ip."assignedMentorId" = u.id
		GROUP BY u.name
		ORDER BY approval_rate DESC
		LIMIT 5
	`
	rows4, err := s.DB.QueryContext(ctx, q4)
	if err == nil {
		defer rows4.Close()
		for rows4.Next() {
			var name string
			var rate float32
			if err := rows4.Scan(&name, &rate); err == nil {
				resp.HighestTaskApproval = append(resp.HighestTaskApproval, &pb.MentorStat{MentorName: name, StatValue: rate})
			}
		}
	}

	// 5. Highest intern retention rates
	q5 := `
		SELECT u.name, 
			(CAST(SUM(CASE WHEN ip."conversionStatus" IN ('COMPLETED', 'CONVERTED') THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(ip.id), 0)) * 100 as retention_rate
		FROM "InternProfile" ip
		JOIN "User" u ON ip."assignedMentorId" = u.id
		GROUP BY u.name
		ORDER BY retention_rate DESC
		LIMIT 5
	`
	rows5, err := s.DB.QueryContext(ctx, q5)
	if err == nil {
		defer rows5.Close()
		for rows5.Next() {
			var name string
			var rate float32
			if err := rows5.Scan(&name, &rate); err == nil {
				resp.HighestRetentionRate = append(resp.HighestRetentionRate, &pb.MentorStat{MentorName: name, StatValue: rate})
			}
		}
	}

	// 6. Consistently miss deadlines
	q6 := `
		SELECT u.name, 
			(CAST(SUM(CASE WHEN t."submissionDate" IS NOT NULL AND t."submissionDate" != '' AND t.deadline != '' AND t."submissionDate"::timestamp > t.deadline::timestamp THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(t.id), 0)) * 100 as late_rate
		FROM "InternProfile" ip
		JOIN "Task" t ON ip."userId" = t."assignedToId"
		JOIN "User" u ON ip."assignedMentorId" = u.id
		GROUP BY u.name
		ORDER BY late_rate DESC
		LIMIT 5
	`
	rows6, err := s.DB.QueryContext(ctx, q6)
	if err == nil {
		defer rows6.Close()
		for rows6.Next() {
			var name string
			var rate float32
			if err := rows6.Scan(&name, &rate); err == nil {
				resp.ConsistentlyMissDeadlines = append(resp.ConsistentlyMissDeadlines, &pb.MentorStat{MentorName: name, StatValue: rate})
			}
		}
	}

	// 7. Most balanced workload (closet to average intern count)
	q7 := `
		WITH mentor_counts AS (
			SELECT "assignedMentorId", CAST(COUNT(*) AS FLOAT) as cnt
			FROM "InternProfile"
			WHERE "assignedMentorId" IS NOT NULL
			GROUP BY "assignedMentorId"
		),
		avg_cnt AS (
			SELECT AVG(cnt) as average FROM mentor_counts
		)
		SELECT u.name, m.cnt, ABS(m.cnt - a.average) as diff
		FROM mentor_counts m
		CROSS JOIN avg_cnt a
		JOIN "User" u ON m."assignedMentorId" = u.id
		ORDER BY diff ASC
		LIMIT 5
	`
	rows7, err := s.DB.QueryContext(ctx, q7)
	if err == nil {
		defer rows7.Close()
		for rows7.Next() {
			var name string
			var cnt float32
			var diff float32
			if err := rows7.Scan(&name, &cnt, &diff); err == nil {
				resp.BalancedWorkload = append(resp.BalancedWorkload, &pb.MentorStat{MentorName: name, StatValue: cnt})
			}
		}
	}

	// 8. Mentor requires intervention (lowest evaluation score)
	q8 := `
		SELECT u.name, AVG(e."totalScore") as avg_score
		FROM "InternProfile" ip
		JOIN "Evaluation" e ON ip."userId" = e."internId"
		JOIN "User" u ON ip."assignedMentorId" = u.id
		GROUP BY u.name
		ORDER BY avg_score ASC
		LIMIT 5
	`
	rows8, err := s.DB.QueryContext(ctx, q8)
	if err == nil {
		defer rows8.Close()
		for rows8.Next() {
			var name string
			var score float32
			if err := rows8.Scan(&name, &score); err == nil {
				resp.NeedsIntervention = append(resp.NeedsIntervention, &pb.MentorStat{MentorName: name, StatValue: score})
			}
		}
	}

	// 9. Mentor effectiveness across domains (avg score per domain)
	q9 := `
		SELECT d."domainName", u.name, AVG(e."totalScore") as avg_score
		FROM "InternProfile" ip
		JOIN "Evaluation" e ON ip."userId" = e."internId"
		JOIN "User" u ON ip."assignedMentorId" = u.id
		JOIN "Mentor" m ON u.id = m."userId"
		JOIN "Domain" d ON m."domainId" = d.id
		GROUP BY d."domainName", u.name
		ORDER BY d."domainName", avg_score DESC
	`
	rows9, err := s.DB.QueryContext(ctx, q9)
	if err == nil {
		defer rows9.Close()
		for rows9.Next() {
			var domain, name string
			var score float32
			if err := rows9.Scan(&domain, &name, &score); err == nil {
				resp.EffectivenessAcrossDomains = append(resp.EffectivenessAcrossDomains, &pb.DomainMentorStat{DomainName: domain, MentorName: name, EffectivenessScore: score})
			}
		}
	}

	// 10. Lagging assessments
	q10 := `
		SELECT u.name, CAST(COUNT(e.id) AS FLOAT) / NULLIF(COUNT(DISTINCT ip.id), 0) as eval_ratio
		FROM "User" u
		JOIN "InternProfile" ip ON u.id = ip."assignedMentorId"
		LEFT JOIN "Evaluation" e ON ip."userId" = e."internId" AND e."createdAt"::timestamp > NOW() - INTERVAL '7 days'
		GROUP BY u.name
		ORDER BY eval_ratio ASC
		LIMIT 5
	`
	rows10, err := s.DB.QueryContext(ctx, q10)
	if err == nil {
		defer rows10.Close()
		for rows10.Next() {
			var name string
			var val float32
			if err := rows10.Scan(&name, &val); err == nil {
				resp.LaggingAssessments = append(resp.LaggingAssessments, &pb.MentorStat{MentorName: name, StatValue: val})
			}
		}
	}

	// 11. Negative empty feedback
	q11 := `
		SELECT u.name, CAST(COUNT(e.id) AS FLOAT) as negative_empty_feedback
		FROM "User" u
		JOIN "Evaluation" e ON u.id = e."mentorId"
		WHERE e."totalScore" < 50 AND (e.comments IS NULL OR e.comments = '')
		GROUP BY u.name
		ORDER BY negative_empty_feedback DESC
		LIMIT 5
	`
	rows11, err := s.DB.QueryContext(ctx, q11)
	if err == nil {
		defer rows11.Close()
		for rows11.Next() {
			var name string
			var val float32
			if err := rows11.Scan(&name, &val); err == nil {
				resp.NegativeEmptyFeedback = append(resp.NegativeEmptyFeedback, &pb.MentorStat{MentorName: name, StatValue: val})
			}
		}
	}

	// 12. Disputed evidence flags
	q12 := `
		SELECT u.name, CAST(COUNT(t.id) AS FLOAT) as disputed_flags
		FROM "User" u
		JOIN "Task" t ON u.id = t."assignedById"
		WHERE t.status IN ('REJECTED', 'REVISION_NEEDED')
		GROUP BY u.name
		ORDER BY disputed_flags DESC
		LIMIT 5
	`
	rows12, err := s.DB.QueryContext(ctx, q12)
	if err == nil {
		defer rows12.Close()
		for rows12.Next() {
			var name string
			var val float32
			if err := rows12.Scan(&name, &val); err == nil {
				resp.DisputedEvidenceFlags = append(resp.DisputedEvidenceFlags, &pb.MentorStat{MentorName: name, StatValue: val})
			}
		}
	}

	return resp, nil
}
