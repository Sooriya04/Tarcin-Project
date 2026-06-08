package grpc

import (
	"backend/internal/grpc/pb"
	"context"
)

func (s *PerformanceService) GetEngagementAnalytics(ctx context.Context, req *pb.Empty) (*pb.EngagementAnalyticsResponse, error) {
	resp := &pb.EngagementAnalyticsResponse{}

	// 1. Interns with no recent attendance (past 7 days)
	q1 := `
		SELECT u.id, u.name, CAST(COUNT(a.id) AS FLOAT) as att_count
		FROM "User" u
		LEFT JOIN "Attendance" a ON u.id = a."userId" AND a.date::timestamp > NOW() - INTERVAL '7 days'
		WHERE u.role = 'INTERN'
		GROUP BY u.id, u.name
		HAVING COUNT(a.id) = 0
		LIMIT 5
	`
	rows1, err := s.DB.QueryContext(ctx, q1)
	if err == nil {
		defer rows1.Close()
		for rows1.Next() {
			var id, name string
			var val float32
			if err := rows1.Scan(&id, &name, &val); err == nil {
				resp.NoRecentAttendance = append(resp.NoRecentAttendance, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 2. Declining participation (consecutiveNegativeWeeks)
	q2 := `
		SELECT u.id, u.name, CAST(w."consecutiveNegativeWeeks" AS FLOAT)
		FROM "User" u
		JOIN "WeeklySnapshot" w ON u.id = w."userId"
		WHERE w."consecutiveNegativeWeeks" > 0 AND u.role = 'INTERN'
		ORDER BY w."consecutiveNegativeWeeks" DESC
		LIMIT 5
	`
	rows2, err := s.DB.QueryContext(ctx, q2)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var id, name string
			var val float32
			if err := rows2.Scan(&id, &name, &val); err == nil {
				resp.DecliningParticipation = append(resp.DecliningParticipation, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 3. Stopped daily reports
	q3 := `
		SELECT u.id, u.name, CAST(COUNT(dr.id) AS FLOAT) as rep_count
		FROM "User" u
		LEFT JOIN "DailyReport" dr ON u.id = dr."internId" AND dr.date::timestamp > NOW() - INTERVAL '7 days'
		WHERE u.role = 'INTERN'
		GROUP BY u.id, u.name
		HAVING COUNT(dr.id) = 0
		LIMIT 5
	`
	rows3, err := s.DB.QueryContext(ctx, q3)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var id, name string
			var val float32
			if err := rows3.Scan(&id, &name, &val); err == nil {
				resp.StoppedDailyReports = append(resp.StoppedDailyReports, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 4. Lost active streaks recently
	q4 := `
		SELECT u.id, u.name, CAST(s."longestStreak" AS FLOAT)
		FROM "User" u
		JOIN "Streak" s ON u.id = s."userId"
		WHERE u.role = 'INTERN' AND s."currentStreak" = 0 AND s."longestStreak" > 3
		ORDER BY s."longestStreak" DESC
		LIMIT 5
	`
	rows4, err := s.DB.QueryContext(ctx, q4)
	if err == nil {
		defer rows4.Close()
		for rows4.Next() {
			var id, name string
			var val float32
			if err := rows4.Scan(&id, &name, &val); err == nil {
				resp.LostActiveStreaks = append(resp.LostActiveStreaks, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 5. Most engaged interns
	q5 := `
		SELECT id, name, CAST(xp AS FLOAT)
		FROM "User"
		WHERE role = 'INTERN'
		ORDER BY xp DESC
		LIMIT 5
	`
	rows5, err := s.DB.QueryContext(ctx, q5)
	if err == nil {
		defer rows5.Close()
		for rows5.Next() {
			var id, name string
			var val float32
			if err := rows5.Scan(&id, &name, &val); err == nil {
				resp.MostEngagedInterns = append(resp.MostEngagedInterns, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 6. At risk interns (escalationLevel > 0)
	q6 := `
		SELECT u.id, u.name, CAST(w."escalationLevel" AS FLOAT)
		FROM "User" u
		JOIN "WeeklySnapshot" w ON u.id = w."userId"
		WHERE w."escalationLevel" > 0 AND u.role = 'INTERN'
		ORDER BY w."escalationLevel" DESC
		LIMIT 5
	`
	rows6, err := s.DB.QueryContext(ctx, q6)
	if err == nil {
		defer rows6.Close()
		for rows6.Next() {
			var id, name string
			var val float32
			if err := rows6.Scan(&id, &name, &val); err == nil {
				resp.AtRiskInterns = append(resp.AtRiskInterns, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 7. Sudden drops in activity levels (high XP but 0 streak)
	q7 := `
		SELECT u.id, u.name, CAST(u.xp AS FLOAT)
		FROM "User" u
		JOIN "Streak" s ON u.id = s."userId"
		WHERE u.role = 'INTERN' AND s."currentStreak" = 0 AND u.xp > 100
		ORDER BY u.xp DESC
		LIMIT 5
	`
	rows7, err := s.DB.QueryContext(ctx, q7)
	if err == nil {
		defer rows7.Close()
		for rows7.Next() {
			var id, name string
			var val float32
			if err := rows7.Scan(&id, &name, &val); err == nil {
				resp.SuddenActivityDrops = append(resp.SuddenActivityDrops, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 8. Colleges producing most engaged interns
	q8 := `
		SELECT c.name, AVG(u.xp) as avg_xp
		FROM "User" u
		JOIN "InternProfile" ip ON u.id = ip."userId"
		JOIN "College" c ON ip."collegeId" = c.id
		WHERE u.role = 'INTERN'
		GROUP BY c.name
		ORDER BY avg_xp DESC
		LIMIT 5
	`
	rows8, err := s.DB.QueryContext(ctx, q8)
	if err == nil {
		defer rows8.Close()
		for rows8.Next() {
			var name string
			var val float32
			if err := rows8.Scan(&name, &val); err == nil {
				resp.TopEngagedColleges = append(resp.TopEngagedColleges, &pb.CollegeEngagementStat{CollegeName: name, AvgEngagementScore: val})
			}
		}
	}

	// 9. Stuck or blocked
	q9 := `
		SELECT u.id, u.name, CAST(COUNT(w.id) AS FLOAT) as stuck_count
		FROM "User" u
		JOIN "WeeklySnapshot" w ON u.id = w."userId"
		WHERE w.momentum IN ('Stuck', 'Blocked')
		GROUP BY u.id, u.name
		ORDER BY stuck_count DESC
		LIMIT 5
	`
	rows9, err := s.DB.QueryContext(ctx, q9)
	if err == nil {
		defer rows9.Close()
		for rows9.Next() {
			var id, name string
			var val float32
			if err := rows9.Scan(&id, &name, &val); err == nil {
				resp.StuckOrBlocked = append(resp.StuckOrBlocked, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 10. Missing or partial signal
	q10 := `
		SELECT u.id, u.name, 1.0 as val
		FROM "User" u
		JOIN "WeeklySnapshot" w ON u.id = w."userId"
		WHERE w.signal IN ('Missing', 'Partial') AND w."weekStartDate"::timestamp > NOW() - INTERVAL '7 days'
		LIMIT 5
	`
	rows10, err := s.DB.QueryContext(ctx, q10)
	if err == nil {
		defer rows10.Close()
		for rows10.Next() {
			var id, name string
			var val float32
			if err := rows10.Scan(&id, &name, &val); err == nil {
				resp.MissingOrPartialSignal = append(resp.MissingOrPartialSignal, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	// 11. High messages low tasks
	q11 := `
		SELECT u.id, u.name, CAST(COUNT(t.id) AS FLOAT) as task_count
		FROM "User" u
		LEFT JOIN "Task" t ON u.id = t."assignedToId" AND t."submissionDate"::timestamp > NOW() - INTERVAL '7 days'
		WHERE u.role = 'INTERN'
		GROUP BY u.id, u.name
		ORDER BY task_count ASC
		LIMIT 5
	`
	rows11, err := s.DB.QueryContext(ctx, q11)
	if err == nil {
		defer rows11.Close()
		for rows11.Next() {
			var id, name string
			var val float32
			if err := rows11.Scan(&id, &name, &val); err == nil {
				resp.HighMessagesLowTasks = append(resp.HighMessagesLowTasks, &pb.EngagementStat{InternId: id, InternName: name, StatValue: val})
			}
		}
	}

	return resp, nil
}
