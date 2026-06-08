package grpc

import (
	"backend/internal/grpc/pb"
	"context"
	"log"
)

func (s *PerformanceService) GetTaskAnalytics(ctx context.Context, req *pb.Empty) (*pb.TaskAnalyticsResponse, error) {
	resp := &pb.TaskAnalyticsResponse{}

	// 1. Highest task completion rate
	q1 := `
		SELECT u.id, u.name, 
			(CAST(SUM(CASE WHEN t.status = 'APPROVED' THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)) * 100 as completion_rate,
			COUNT(*) as total
		FROM "Task" t
		JOIN "User" u ON t."assignedToId" = u.id
		GROUP BY u.id, u.name
		HAVING COUNT(*) > 0
		ORDER BY completion_rate DESC
		LIMIT 5
	`
	rows1, err := s.DB.QueryContext(ctx, q1)
	if err == nil {
		defer rows1.Close()
		for rows1.Next() {
			var id, name string
			var rate float32
			var count int32
			if err := rows1.Scan(&id, &name, &rate, &count); err == nil {
				resp.HighestCompletionRate = append(resp.HighestCompletionRate, &pb.InternTaskStat{InternId: id, InternName: name, CompletionRate: rate, Count: count})
			}
		}
	} else {
		log.Printf("Q1 Error: %v\n", err)
	}

	// 2. Most late tasks
	q2 := `
		SELECT u.id, u.name, COUNT(*) as late_count
		FROM "Task" t
		JOIN "User" u ON t."assignedToId" = u.id
		WHERE t."submissionDate" IS NOT NULL 
		  AND t."submissionDate" != '' 
		  AND t.deadline != ''
		  AND t."submissionDate"::timestamp > t.deadline::timestamp
		GROUP BY u.id, u.name
		ORDER BY late_count DESC
		LIMIT 5
	`
	rows2, err := s.DB.QueryContext(ctx, q2)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var id, name string
			var count int32
			if err := rows2.Scan(&id, &name, &count); err == nil {
				resp.MostLateTasks = append(resp.MostLateTasks, &pb.InternTaskStat{InternId: id, InternName: name, Count: count})
			}
		}
	} else {
		log.Printf("Q2 Error: %v\n", err)
	}

	// 3. Most frequently rejected tasks
	q3 := `
		SELECT title, COUNT(*) as rejection_count
		FROM "Task"
		WHERE status = 'REJECTED'
		GROUP BY title
		ORDER BY rejection_count DESC
		LIMIT 5
	`
	rows3, err := s.DB.QueryContext(ctx, q3)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var title string
			var count int32
			if err := rows3.Scan(&title, &count); err == nil {
				resp.FrequentlyRejectedTasks = append(resp.FrequentlyRejectedTasks, &pb.TaskTitleStat{Title: title, Count: count})
			}
		}
	} else {
		log.Printf("Q3 Error: %v\n", err)
	}

	// 4. Percentage of incomplete tasks
	q4 := `
		SELECT COALESCE((CAST(SUM(CASE WHEN status NOT IN ('APPROVED', 'REJECTED', 'SUBMITTED') THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(*), 0)) * 100, 0)
		FROM "Task"
	`
	err = s.DB.QueryRowContext(ctx, q4).Scan(&resp.IncompleteTasksPercentage)
	if err != nil {
		log.Printf("Q4 Error: %v\n", err)
	}

	// 5. Domains with highest task approval rates
	q5 := `
		SELECT ip."preferredDomain", 
			(CAST(SUM(CASE WHEN t.status = 'APPROVED' THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(*), 0)) * 100 as approval_rate
		FROM "Task" t
		JOIN "InternProfile" ip ON t."assignedToId" = ip."userId"
		WHERE ip."preferredDomain" IS NOT NULL
		GROUP BY ip."preferredDomain"
		ORDER BY approval_rate DESC
		LIMIT 5
	`
	rows5, err := s.DB.QueryContext(ctx, q5)
	if err == nil {
		defer rows5.Close()
		for rows5.Next() {
			var domain string
			var rate float32
			if err := rows5.Scan(&domain, &rate); err == nil {
				resp.DomainApprovalRates = append(resp.DomainApprovalRates, &pb.DomainApprovalStat{Domain: domain, ApprovalRate: rate})
			}
		}
	} else {
		log.Printf("Q5 Error: %v\n", err)
	}

	// 6. Mentors with interns having fastest task completion times
	q6 := `
		SELECT u.name as mentor_name, 
			AVG(EXTRACT(EPOCH FROM (t."submissionDate"::timestamp - t."createdAt"::timestamp))/3600) as avg_hours
		FROM "Task" t
		JOIN "User" u ON t."assignedById" = u.id
		WHERE t."submissionDate" IS NOT NULL AND t."submissionDate" != '' AND t.status = 'APPROVED'
		GROUP BY u.id, u.name
		ORDER BY avg_hours ASC
		LIMIT 5
	`
	rows6, err := s.DB.QueryContext(ctx, q6)
	if err == nil {
		defer rows6.Close()
		for rows6.Next() {
			var mentor string
			var hours float32
			if err := rows6.Scan(&mentor, &hours); err == nil {
				resp.FastestMentorInterns = append(resp.FastestMentorInterns, &pb.MentorSpeedStat{MentorName: mentor, AvgCompletionHours: hours})
			}
		}
	} else {
		log.Printf("Q6 Error: %v\n", err)
	}

	// 8. Top reasons for task rejection
	q8 := `
		SELECT "mentorFeedback"
		FROM "Task"
		WHERE status = 'REJECTED' AND "mentorFeedback" IS NOT NULL AND "mentorFeedback" != ''
		GROUP BY "mentorFeedback"
		ORDER BY COUNT(*) DESC
		LIMIT 5
	`
	rows8, err := s.DB.QueryContext(ctx, q8)
	if err == nil {
		defer rows8.Close()
		for rows8.Next() {
			var feedback string
			if err := rows8.Scan(&feedback); err == nil {
				resp.TopRejectionReasons = append(resp.TopRejectionReasons, feedback)
			}
		}
	} else {
		log.Printf("Q8 Error: %v\n", err)
	}

	// 9. Consistently submit tasks before deadlines
	q9 := `
		SELECT u.id, u.name, COUNT(*) as early_count
		FROM "Task" t
		JOIN "User" u ON t."assignedToId" = u.id
		WHERE t."submissionDate" IS NOT NULL AND t."submissionDate" != '' AND t.deadline != '' AND t."submissionDate"::timestamp <= t.deadline::timestamp
		GROUP BY u.id, u.name
		ORDER BY early_count DESC
		LIMIT 5
	`
	rows9, err := s.DB.QueryContext(ctx, q9)
	if err == nil {
		defer rows9.Close()
		for rows9.Next() {
			var id, name string
			var count int32
			if err := rows9.Scan(&id, &name, &count); err == nil {
				resp.ConsistentEarlySubmitters = append(resp.ConsistentEarlySubmitters, &pb.InternTaskStat{InternId: id, InternName: name, Count: count})
			}
		}
	} else {
		log.Printf("Q9 Error: %v\n", err)
	}

	// 10. Pending for > 14 days
	q10 := `
		SELECT id, title, EXTRACT(DAY FROM (NOW() - t."createdAt"::timestamp)) as days_pending
		FROM "Task" t
		WHERE status IN ('TODO', 'IN_PROGRESS', 'REVISION_NEEDED') 
		  AND NOW() - t."createdAt"::timestamp > INTERVAL '14 days'
		ORDER BY days_pending DESC
		LIMIT 10
	`
	rows10, err := s.DB.QueryContext(ctx, q10)
	if err == nil {
		defer rows10.Close()
		for rows10.Next() {
			var id, title string
			var days int32
			if err := rows10.Scan(&id, &title, &days); err == nil {
				resp.LongPendingTasks = append(resp.LongPendingTasks, &pb.PendingTask{TaskId: id, Title: title, DaysPending: days})
			}
		}
	} else {
		log.Printf("Q10 Error: %v\n", err)
	}

	// 11. Tasks without evidence
	q11 := `
		SELECT COALESCE(CAST(SUM(CASE WHEN "submissionUrl" IS NULL OR "submissionUrl" = '' THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(*), 0) * 100, 0)
		FROM "Task"
		WHERE status IN ('APPROVED', 'SUBMITTED')
	`
	_ = s.DB.QueryRowContext(ctx, q11).Scan(&resp.TasksWithoutEvidencePercentage)

	// 12. Tasks with highest blockers
	q12 := `
		SELECT title, CAST(COUNT(*) AS FLOAT) as blocker_count
		FROM "Task"
		WHERE status = 'BLOCKERS'
		GROUP BY title
		ORDER BY blocker_count DESC
		LIMIT 5
	`
	rows12, err := s.DB.QueryContext(ctx, q12)
	if err == nil {
		defer rows12.Close()
		for rows12.Next() {
			var title string
			var count float32
			if err := rows12.Scan(&title, &count); err == nil {
				resp.TasksWithMostBlockers = append(resp.TasksWithMostBlockers, &pb.TaskTitleStat{Title: title, Count: int32(count)})
			}
		}
	}

	// 13. Average time spent vs approval
	q13 := `
		SELECT status, CAST(AVG(EXTRACT(EPOCH FROM ("submissionDate"::timestamp - "createdAt"::timestamp)) / 60) AS FLOAT) as avg_minutes
		FROM "Task"
		WHERE status IN ('APPROVED', 'REVISION_NEEDED', 'REJECTED') AND "submissionDate" IS NOT NULL AND "submissionDate" != ''
		GROUP BY status
	`
	rows13, err := s.DB.QueryContext(ctx, q13)
	if err == nil {
		defer rows13.Close()
		for rows13.Next() {
			var status string
			var avg float32
			if err := rows13.Scan(&status, &avg); err == nil {
				resp.TimeVsApprovalCorrelation = append(resp.TimeVsApprovalCorrelation, &pb.TaskTimeCorrelation{Status: status, AvgMinutes: avg})
			}
		}
	}

	return resp, nil
}
