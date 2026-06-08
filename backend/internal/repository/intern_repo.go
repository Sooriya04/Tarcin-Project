package repository

import (
	"backend/internal/models"
	"database/sql"
	"log"
)

type InternRepository struct {
	DB *sql.DB
}

func NewInternRepository(db *sql.DB) *InternRepository {
	return &InternRepository{DB: db}
}

func (r *InternRepository) GetHealthMetrics() models.InternHealthResponse {
	var resp models.InternHealthResponse
	resp.TopColleges = make(map[string]int)
	resp.FastestDomains = make(map[string]int)

	// 1. Active interns
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM "User" WHERE role = 'INTERN' AND status = 'APPROVED'`).Scan(&resp.ActiveInterns)
	if err != nil {
		log.Printf("Query error (ActiveInterns): %v\n", err)
	}

	// 2. Completed internships
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM "InternProfile" WHERE "conversionStatus" IN ('COMPLETED', 'CONVERTED')`).Scan(&resp.CompletedInterns)
	if err != nil {
		log.Printf("Query error (CompletedInterns): %v\n", err)
	}

	// 3. Inactive interns
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM "User" WHERE role = 'INTERN' AND status IN ('PENDING', 'REJECTED')`).Scan(&resp.InactiveInterns)
	if err != nil {
		log.Printf("Query error (InactiveInterns): %v\n", err)
	}

	// 4. Top colleges
	collegeRows, err := r.DB.Query(`
		SELECT c.name, COUNT(ip.id) 
		FROM "InternProfile" ip
		JOIN "College" c ON ip."collegeId" = c.id
		GROUP BY c.name
		ORDER BY COUNT(ip.id) DESC
		LIMIT 5
	`)
	if err == nil {
		defer collegeRows.Close()
		for collegeRows.Next() {
			var name string
			var count int
			if err := collegeRows.Scan(&name, &count); err == nil {
				resp.TopColleges[name] = count
			}
		}
	} else {
		log.Printf("Query error (TopColleges): %v\n", err)
	}

	// 5. Fastest growing domains
	domainRows, err := r.DB.Query(`
		SELECT "preferredDomain", COUNT(id)
		FROM "InternProfile"
		WHERE "preferredDomain" IS NOT NULL
		GROUP BY "preferredDomain"
		ORDER BY COUNT(id) DESC
		LIMIT 5
	`)
	if err == nil {
		defer domainRows.Close()
		for domainRows.Next() {
			var domain string
			var count int
			if err := domainRows.Scan(&domain, &count); err == nil {
				resp.FastestDomains[domain] = count
			}
		}
	} else {
		log.Printf("Query error (FastestDomains): %v\n", err)
	}

	return resp
}
