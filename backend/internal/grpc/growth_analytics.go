package grpc

import (
	"backend/internal/grpc/pb"
	"context"
)

func (s *PerformanceService) GetGrowthAnalytics(ctx context.Context, req *pb.Empty) (*pb.GrowthAnalyticsResponse, error) {
	resp := &pb.GrowthAnalyticsResponse{}

	// 1. Most XP gained this week
	q1 := `
		SELECT u.name, CAST(w."xpGained" AS FLOAT)
		FROM "User" u
		JOIN "WeeklySnapshot" w ON u.id = w."userId"
		WHERE w."weekStartDate"::timestamp > NOW() - INTERVAL '7 days' AND u.role = 'INTERN'
		ORDER BY w."xpGained" DESC
		LIMIT 5
	`
	rows1, err := s.DB.QueryContext(ctx, q1)
	if err == nil {
		defer rows1.Close()
		for rows1.Next() {
			var name string
			var val float32
			if err := rows1.Scan(&name, &val); err == nil {
				resp.MostXpThisWeek = append(resp.MostXpThisWeek, &pb.GrowthStat{ItemName: name, StatValue: val})
			}
		}
	}

	// 2. Most frequent achievements
	q2 := `
		SELECT a.title, CAST(COUNT(ua.id) AS FLOAT) as frequency
		FROM "Achievement" a
		JOIN "UserAchievement" ua ON a.id = ua."achievementId"
		GROUP BY a.title
		ORDER BY frequency DESC
		LIMIT 5
	`
	rows2, err := s.DB.QueryContext(ctx, q2)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var name string
			var val float32
			if err := rows2.Scan(&name, &val); err == nil {
				resp.FrequentAchievements = append(resp.FrequentAchievements, &pb.GrowthStat{ItemName: name, StatValue: val})
			}
		}
	}

	// 3. Do high-XP interns perform better in evaluations?
	q3 := `
		SELECT u.name, CAST(u.xp AS FLOAT), CAST(AVG(e."totalScore") AS FLOAT)
		FROM "User" u
		JOIN "Evaluation" e ON u.id = e."internId"
		WHERE u.role = 'INTERN'
		GROUP BY u.id, u.name, u.xp
		ORDER BY u.xp DESC
		LIMIT 5
	`
	rows3, err := s.DB.QueryContext(ctx, q3)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var name string
			var xp, eval float32
			if err := rows3.Scan(&name, &xp, &eval); err == nil {
				resp.XpVsEvaluations = append(resp.XpVsEvaluations, &pb.GrowthStat{ItemName: name, StatValue: xp, SecondaryValue: eval})
			}
		}
	}

	// 4. Streak vs task completion
	q4 := `
		SELECT u.name, CAST(s."longestStreak" AS FLOAT), 
		       COALESCE(CAST((CAST(SUM(CASE WHEN t.status = 'APPROVED' THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(t.id), 0)) * 100 AS FLOAT), 0)
		FROM "User" u
		JOIN "Streak" s ON u.id = s."userId"
		LEFT JOIN "Task" t ON u.id = t."assignedToId"
		WHERE u.role = 'INTERN'
		GROUP BY u.id, u.name, s."longestStreak"
		ORDER BY s."longestStreak" DESC
		LIMIT 5
	`
	rows4, err := s.DB.QueryContext(ctx, q4)
	if err == nil {
		defer rows4.Close()
		for rows4.Next() {
			var name string
			var streak, comp float32
			if err := rows4.Scan(&name, &streak, &comp); err == nil {
				resp.StreakVsTasks = append(resp.StreakVsTasks, &pb.GrowthStat{ItemName: name, StatValue: streak, SecondaryValue: comp})
			}
		}
	}

	// 5. Domains earn XP fastest
	q5 := `
		SELECT ip."preferredDomain", CAST(AVG(u.xp) AS FLOAT) as avg_xp
		FROM "InternProfile" ip
		JOIN "User" u ON ip."userId" = u.id
		WHERE u.role = 'INTERN' AND ip."preferredDomain" IS NOT NULL
		GROUP BY ip."preferredDomain"
		ORDER BY avg_xp DESC
		LIMIT 5
	`
	rows5, err := s.DB.QueryContext(ctx, q5)
	if err == nil {
		defer rows5.Close()
		for rows5.Next() {
			var name string
			var xp float32
			if err := rows5.Scan(&name, &xp); err == nil {
				resp.FastestXpDomains = append(resp.FastestXpDomains, &pb.GrowthStat{ItemName: name, StatValue: xp})
			}
		}
	}

	// 6. Closest to level up
	q6 := `
		SELECT name, CAST((1000 - (xp % 1000)) AS FLOAT) as xp_to_next
		FROM "User"
		WHERE role = 'INTERN'
		ORDER BY xp_to_next ASC
		LIMIT 5
	`
	rows6, err := s.DB.QueryContext(ctx, q6)
	if err == nil {
		defer rows6.Close()
		for rows6.Next() {
			var name string
			var xp float32
			if err := rows6.Scan(&name, &xp); err == nil {
				resp.ClosestToLevelUp = append(resp.ClosestToLevelUp, &pb.GrowthStat{ItemName: name, StatValue: xp})
			}
		}
	}

	// 7. Active gamification percentage
	q7 := `
		SELECT COALESCE(CAST(SUM(CASE WHEN xp > 0 THEN 1 ELSE 0 END) AS FLOAT) / NULLIF(COUNT(*), 0) * 100, 0)
		FROM "User"
		WHERE role = 'INTERN'
	`
	_ = s.DB.QueryRowContext(ctx, q7).Scan(&resp.ActiveGamificationPercentage)

	// 8. Stagnated XP growth
	q8 := `
		SELECT u.name, CAST(u.xp AS FLOAT)
		FROM "User" u
		JOIN "WeeklySnapshot" w ON u.id = w."userId"
		WHERE u.role = 'INTERN' AND w."weekStartDate"::timestamp > NOW() - INTERVAL '7 days' AND w."xpGained" = 0 AND u.xp > 500
		ORDER BY u.xp DESC
		LIMIT 5
	`
	rows8, err := s.DB.QueryContext(ctx, q8)
	if err == nil {
		defer rows8.Close()
		for rows8.Next() {
			var name string
			var xp float32
			if err := rows8.Scan(&name, &xp); err == nil {
				resp.StagnatedXpGrowth = append(resp.StagnatedXpGrowth, &pb.GrowthStat{ItemName: name, StatValue: xp})
			}
		}
	}

	return resp, nil
}
