package grpc

import (
	"backend/internal/grpc/pb"
	"backend/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

func (s *PerformanceService) GetUserInformation(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	emailOrId := req.GetEmailOrId()
	log.Printf("[GetUserInformation] Fetching information for: %s", emailOrId)

	// 1. Fetch main User record
	user, err := s.queryRowToMap(ctx, `SELECT id, name, email, role, status, "createdAt", "updatedAt", "referralCode", "referredById", xp, level FROM "User" WHERE id = $1 OR email = $1`, emailOrId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[GetUserInformation] User not found: %s", emailOrId)
			return &pb.UserResponse{
				Error: fmt.Sprintf("User not found: %s", emailOrId),
			}, nil
		}
		log.Printf("[GetUserInformation] DB Error: %v", err)
		return &pb.UserResponse{
			Error: fmt.Sprintf("DB Error: %v", err),
		}, nil
	}

	userId := fmt.Sprintf("%v", user["id"])
	userName := fmt.Sprintf("%v", user["name"])
	userEmail := fmt.Sprintf("%v", user["email"])

	// 2. Fetch all other records related to this user
	result := make(map[string]interface{})
	result["User"] = user

	// Helper function to query and set list or object
	runQuery := func(key, query string, args ...interface{}) {
		rowsMap, err := s.queryRowsToMaps(ctx, query, args...)
		if err != nil {
			log.Printf("[GetUserInformation] Error querying %s: %v", key, err)
			result[key] = []map[string]interface{}{}
		} else {
			result[key] = rowsMap
		}
	}

	runQuerySingle := func(key, query string, args ...interface{}) {
		rowsMap, err := s.queryRowsToMaps(ctx, query, args...)
		if err != nil || len(rowsMap) == 0 {
			result[key] = nil
		} else {
			result[key] = rowsMap[0]
		}
	}

	// Fetch InternProfile along with College details (joining College table)
	runQuerySingle("InternProfile", `
		SELECT ip.*, c.name AS "collegeName", c.department AS "collegeDept", c."hodName" AS "collegeHodName", c."hodEmail" AS "collegeHodEmail", c.city AS "collegeCity", c.state AS "collegeState"
		FROM "InternProfile" ip
		LEFT JOIN "College" c ON ip."collegeId" = c.id
		WHERE ip."userId" = $1`, userId)

	// Fetch Mentor along with Domain details (joining Domain table)
	runQuerySingle("Mentor", `
		SELECT m.*, d."domainName"
		FROM "Mentor" m
		LEFT JOIN "Domain" d ON m."domainId" = d.id
		WHERE m."userId" = $1`, userId)

	// Fetch referred users (joining User table on itself)
	runQuery("ReferredUsers", `SELECT id, name, email, role, status, xp, level FROM "User" WHERE "referredById" = $1`, userId)

	runQuery("TasksAssignedTo", `SELECT * FROM "Task" WHERE "assignedToId" = $1`, userId)
	runQuery("TasksAssignedBy", `SELECT * FROM "Task" WHERE "assignedById" = $1`, userId)
	runQuery("Attendance", `SELECT * FROM "Attendance" WHERE "userId" = $1`, userId)
	runQuery("Certificates", `SELECT * FROM "Certificate" WHERE "userId" = $1`, userId)
	runQuery("Notifications", `SELECT * FROM "Notification" WHERE "userId" = $1`, userId)
	runQuery("EvaluationsAsIntern", `SELECT * FROM "Evaluation" WHERE "internId" = $1`, userId)
	runQuery("EvaluationsAsEvaluator", `SELECT * FROM "Evaluation" WHERE "evaluatorId" = $1`, userId)
	runQuery("MeetingsAsIntern", `SELECT * FROM "Meeting" WHERE "internId" = $1`, userId)
	runQuery("MeetingsAsMentor", `SELECT * FROM "Meeting" WHERE "mentorId" = $1`, userId)
	runQuery("Achievements", `SELECT * FROM "Achievement" WHERE "userId" = $1`, userId)
	runQuerySingle("Streak", `SELECT * FROM "Streak" WHERE "userId" = $1`, userId)
	runQuery("XPLogs", `SELECT * FROM "XPLog" WHERE "userId" = $1`, userId)
	runQuery("DailyReports", `SELECT * FROM "DailyReport" WHERE "internId" = $1`, userId)
	runQuery("AILogsAsAdmin", `SELECT * FROM "AILog" WHERE "adminId" = $1`, userId)
	runQuery("AILogsAsTarget", `SELECT * FROM "AILog" WHERE "targetInternId" = $1`, userId)
	runQuery("MessagesSent", `SELECT * FROM "Message" WHERE "senderId" = $1`, userId)
	runQuery("MessagesReceived", `SELECT * FROM "Message" WHERE "receiverId" = $1`, userId)
	runQuery("WorkLogs", `SELECT * FROM "WorkLog" WHERE "userId" = $1`, userId)
	runQuery("MentorAssessmentsAsMentor", `SELECT * FROM "MentorAssessment" WHERE "mentorId" = $1`, userId)
	runQuery("MentorAssessmentsAsMentee", `SELECT * FROM "MentorAssessment" WHERE "menteeId" = $1`, userId)
	runQuery("DecisionWorkflowsAsSubject", `SELECT * FROM "DecisionWorkflow" WHERE "subjectId" = $1`, userId)
	runQuery("DecisionWorkflowsOpenedBy", `SELECT * FROM "DecisionWorkflow" WHERE "openedBy" = $1`, userId)
	
	// Fetch WorkflowActions with associated DecisionWorkflow details
	runQuery("WorkflowActions", `
		SELECT wa.*, dw.type AS "workflowType", dw.status AS "workflowStatus", dw.severity AS "workflowSeverity"
		FROM "WorkflowAction" wa
		LEFT JOIN "DecisionWorkflow" dw ON wa."workflowId" = dw.id
		WHERE wa."assignedTo" = $1`, userId)

	// Fetch WorkflowApprovalSteps with associated DecisionWorkflow details
	runQuery("WorkflowApprovalSteps", `
		SELECT was.*, dw.type AS "workflowType", dw.status AS "workflowStatus", dw.severity AS "workflowSeverity"
		FROM "WorkflowApprovalStep" was
		LEFT JOIN "DecisionWorkflow" dw ON was."workflowId" = dw.id
		WHERE was."approverId" = $1`, userId)

	runQuery("WeeklySnapshots", `SELECT * FROM "WeeklySnapshot" WHERE "userId" = $1`, userId)
	runQuery("ReportDrafts", `SELECT * FROM "ReportDraft" WHERE "recipientUserId" = $1`, userId)
	runQuery("FlagResponsesAsMentor", `SELECT * FROM "FlagResponse" WHERE "mentorEmail" = $1`, userEmail)
	runQuery("FlagResponsesAsIntern", `SELECT * FROM "FlagResponse" WHERE "internName" = $1`, userName)

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("[GetUserInformation] JSON Marshal error: %v", err)
		return &pb.UserResponse{
			Error: fmt.Sprintf("JSON Marshal Error: %v", err),
		}, nil
	}

	log.Printf("[GetUserInformation] Success: gathered all information for user %s", userId)
	return &pb.UserResponse{
		JsonResult: string(jsonData),
	}, nil
}

func (s *PerformanceService) queryRowsToMaps(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
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
		return nil, err
	}
	if result == nil {
		result = []map[string]interface{}{}
	}
	return result, nil
}

func (s *PerformanceService) queryRowToMap(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	results, err := s.queryRowsToMaps(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	return results[0], nil
}

