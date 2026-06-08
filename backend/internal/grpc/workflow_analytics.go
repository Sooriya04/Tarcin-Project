package grpc

import (
	"backend/internal/grpc/pb"
	"context"
)

func (s *PerformanceService) GetWorkflowAnalytics(ctx context.Context, req *pb.Empty) (*pb.WorkflowAnalyticsResponse, error) {
	resp := &pb.WorkflowAnalyticsResponse{}

	// 1. Escalated interns and triggers
	q1 := `
		SELECT u.name, CAST(w."escalationLevel" AS FLOAT), COALESCE(w.trigger, 'Unknown')
		FROM "User" u
		JOIN "WeeklySnapshot" w ON u.id = w."userId"
		WHERE w."escalationLevel" > 0 AND w."weekStartDate"::timestamp > NOW() - INTERVAL '7 days'
		LIMIT 5
	`
	rows1, err := s.DB.QueryContext(ctx, q1)
	if err == nil {
		defer rows1.Close()
		for rows1.Next() {
			var name, trigger string
			var level float32
			if err := rows1.Scan(&name, &level, &trigger); err == nil {
				resp.EscalatedInterns = append(resp.EscalatedInterns, &pb.WorkflowStat{InternName: name, EscalationLevel: level, RootCause: trigger})
			}
		}
	}

	// 2. Trapped in consecutive negative weeks
	q2 := `
		SELECT u.name, CAST(w."consecutiveNegativeWeeks" AS FLOAT), 'Negative Trend'
		FROM "User" u
		JOIN "WeeklySnapshot" w ON u.id = w."userId"
		WHERE w."consecutiveNegativeWeeks" > 2
		LIMIT 5
	`
	rows2, err := s.DB.QueryContext(ctx, q2)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var name, trigger string
			var level float32
			if err := rows2.Scan(&name, &level, &trigger); err == nil {
				resp.NegativeWeeksTrend = append(resp.NegativeWeeksTrend, &pb.WorkflowStat{InternName: name, EscalationLevel: level, RootCause: trigger})
			}
		}
	}

	// 3. Most common root causes
	q3 := `
		SELECT 'All Cohorts', 1.0, COALESCE(trigger, 'Performance') as trigger_reason
		FROM "WeeklySnapshot"
		WHERE trigger IS NOT NULL AND trigger != ''
		GROUP BY trigger_reason
		ORDER BY COUNT(*) DESC
		LIMIT 5
	`
	rows3, err := s.DB.QueryContext(ctx, q3)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var name, trigger string
			var level float32
			if err := rows3.Scan(&name, &level, &trigger); err == nil {
				resp.CommonRootCauses = append(resp.CommonRootCauses, &pb.WorkflowStat{InternName: name, EscalationLevel: level, RootCause: trigger})
			}
		}
	}

	// 4. Past due decisions (proxy: pending tasks > 14 days)
	q4 := `
		SELECT u.name, CAST(EXTRACT(DAY FROM NOW() - t."createdAt"::timestamp) AS FLOAT), t.title
		FROM "Task" t
		JOIN "User" u ON t."assignedToId" = u.id
		WHERE t.status = 'PENDING' AND t."createdAt"::timestamp < NOW() - INTERVAL '14 days'
		LIMIT 5
	`
	rows4, err := s.DB.QueryContext(ctx, q4)
	if err == nil {
		defer rows4.Close()
		for rows4.Next() {
			var name, trigger string
			var level float32
			if err := rows4.Scan(&name, &level, &trigger); err == nil {
				resp.PastDueDecisions = append(resp.PastDueDecisions, &pb.WorkflowStat{InternName: name, EscalationLevel: level, RootCause: trigger})
			}
		}
	}

	return resp, nil
}
