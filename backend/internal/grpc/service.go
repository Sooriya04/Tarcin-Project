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

	// 1. Fetch main User record (joined with referrer's details)
	user, err := s.queryRowToMap(ctx, `
		SELECT u.id, u.name, u.email, u.role, u.status, u."createdAt", u."updatedAt", u."referralCode", u."referredById", u.xp, u.level,
		       ref.name AS "referredByName", ref.email AS "referredByEmail"
		FROM "User" u
		LEFT JOIN "User" ref ON u."referredById" = ref.id
		WHERE u.id = $1 OR u.email = $1`, emailOrId)
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

	// Fetch InternProfile along with College and Mentor details
	runQuerySingle("InternProfile", `
		SELECT ip.*, 
		       c.name AS "collegeName", c.department AS "collegeDept", c."hodName" AS "collegeHodName", c."hodEmail" AS "collegeHodEmail", c.city AS "collegeCity", c.state AS "collegeState",
		       m_user.name AS "assignedMentorName", m_user.email AS "assignedMentorEmail"
		FROM "InternProfile" ip
		LEFT JOIN "College" c ON ip."collegeId" = c.id
		LEFT JOIN "User" m_user ON ip."assignedMentorId" = m_user.id
		WHERE ip."userId" = $1`, userId)

	// Fetch Mentor along with Domain details
	runQuerySingle("Mentor", `
		SELECT m.*, d."domainName"
		FROM "Mentor" m
		LEFT JOIN "Domain" d ON m."domainId" = d.id
		WHERE m."userId" = $1`, userId)

	// Fetch referred users
	runQuery("ReferredUsers", `
		SELECT u.id, u.name, u.email, u.role, u.status, u.xp, u.level,
		       ref.name AS "referrerName"
		FROM "User" u
		LEFT JOIN "User" ref ON u."referredById" = ref.id
		WHERE u."referredById" = $1`, userId)

	// Fetch Tasks Assigned To (with Assigner details)
	runQuery("TasksAssignedTo", `
		SELECT t.*, u.name AS "assignedByName", u.email AS "assignedByEmail"
		FROM "Task" t
		LEFT JOIN "User" u ON t."assignedById" = u.id
		WHERE t."assignedToId" = $1`, userId)

	// Fetch Tasks Assigned By (with Assignee details)
	runQuery("TasksAssignedBy", `
		SELECT t.*, u.name AS "assignedToName", u.email AS "assignedToEmail"
		FROM "Task" t
		LEFT JOIN "User" u ON t."assignedToId" = u.id
		WHERE t."assignedById" = $1`, userId)

	// Fetch Attendance
	runQuery("Attendance", `
		SELECT a.*, u.name AS "userName", u.email AS "userEmail"
		FROM "Attendance" a
		LEFT JOIN "User" u ON a."userId" = u.id
		WHERE a."userId" = $1`, userId)

	// Fetch Certificates
	runQuery("Certificates", `
		SELECT cert.*, u.name AS "userName", u.email AS "userEmail"
		FROM "Certificate" cert
		LEFT JOIN "User" u ON cert."userId" = u.id
		WHERE cert."userId" = $1`, userId)

	// Fetch Notifications
	runQuery("Notifications", `
		SELECT n.*, u.name AS "userName", u.email AS "userEmail"
		FROM "Notification" n
		LEFT JOIN "User" u ON n."userId" = u.id
		WHERE n."userId" = $1`, userId)

	// Fetch Evaluations as Intern (with Evaluator details)
	runQuery("EvaluationsAsIntern", `
		SELECT e.*, u.name AS "evaluatorName", u.email AS "evaluatorEmail"
		FROM "Evaluation" e
		LEFT JOIN "User" u ON e."evaluatorId" = u.id
		WHERE e."internId" = $1`, userId)

	// Fetch Evaluations as Evaluator (with Intern details)
	runQuery("EvaluationsAsEvaluator", `
		SELECT e.*, u.name AS "internName", u.email AS "internEmail"
		FROM "Evaluation" e
		LEFT JOIN "User" u ON e."internId" = u.id
		WHERE e."evaluatorId" = $1`, userId)

	// Fetch Meetings as Intern (with Mentor details)
	runQuery("MeetingsAsIntern", `
		SELECT m.*, u.name AS "mentorName", u.email AS "mentorEmail"
		FROM "Meeting" m
		LEFT JOIN "User" u ON m."mentorId" = u.id
		WHERE m."internId" = $1`, userId)

	// Fetch Meetings as Mentor (with Intern details)
	runQuery("MeetingsAsMentor", `
		SELECT m.*, u.name AS "internName", u.email AS "internEmail"
		FROM "Meeting" m
		LEFT JOIN "User" u ON m."internId" = u.id
		WHERE m."mentorId" = $1`, userId)

	// Fetch Achievements
	runQuery("Achievements", `
		SELECT ach.*, u.name AS "userName", u.email AS "userEmail"
		FROM "Achievement" ach
		LEFT JOIN "User" u ON ach."userId" = u.id
		WHERE ach."userId" = $1`, userId)

	// Fetch Streak
	runQuerySingle("Streak", `
		SELECT str.*, u.name AS "userName", u.email AS "userEmail"
		FROM "Streak" str
		LEFT JOIN "User" u ON str."userId" = u.id
		WHERE str."userId" = $1`, userId)

	// Fetch XPLogs
	runQuery("XPLogs", `
		SELECT xl.*, u.name AS "userName", u.email AS "userEmail"
		FROM "XPLog" xl
		LEFT JOIN "User" u ON xl."userId" = u.id
		WHERE xl."userId" = $1`, userId)

	// Fetch DailyReports (with Intern details)
	runQuery("DailyReports", `
		SELECT dr.*, u.name AS "internName", u.email AS "internEmail"
		FROM "DailyReport" dr
		LEFT JOIN "User" u ON dr."internId" = u.id
		WHERE dr."internId" = $1`, userId)

	// Fetch AILogs as Admin (with Target Intern details)
	runQuery("AILogsAsAdmin", `
		SELECT al.*, u.name AS "targetInternName", u.email AS "targetInternEmail"
		FROM "AILog" al
		LEFT JOIN "User" u ON al."targetInternId" = u.id
		WHERE al."adminId" = $1`, userId)

	// Fetch AILogs as Target (with Admin details)
	runQuery("AILogsAsTarget", `
		SELECT al.*, u.name AS "adminName", u.email AS "adminEmail"
		FROM "AILog" al
		LEFT JOIN "User" u ON al."adminId" = u.id
		WHERE al."targetInternId" = $1`, userId)

	// Fetch Messages Sent (with Receiver details)
	runQuery("MessagesSent", `
		SELECT msg.*, u.name AS "receiverName", u.email AS "receiverEmail"
		FROM "Message" msg
		LEFT JOIN "User" u ON msg."receiverId" = u.id
		WHERE msg."senderId" = $1`, userId)

	// Fetch Messages Received (with Sender details)
	runQuery("MessagesReceived", `
		SELECT msg.*, u.name AS "senderName", u.email AS "senderEmail"
		FROM "Message" msg
		LEFT JOIN "User" u ON msg."senderId" = u.id
		WHERE msg."receiverId" = $1`, userId)

	// Fetch WorkLogs
	runQuery("WorkLogs", `
		SELECT wl.*, u.name AS "userName", u.email AS "userEmail"
		FROM "WorkLog" wl
		LEFT JOIN "User" u ON wl."userId" = u.id
		WHERE wl."userId" = $1`, userId)

	// Fetch MentorAssessments as Mentor (with Mentee details)
	runQuery("MentorAssessmentsAsMentor", `
		SELECT ma.*, u.name AS "menteeName", u.email AS "menteeEmail"
		FROM "MentorAssessment" ma
		LEFT JOIN "User" u ON ma."menteeId" = u.id
		WHERE ma."mentorId" = $1`, userId)

	// Fetch MentorAssessments as Mentee (with Mentor details)
	runQuery("MentorAssessmentsAsMentee", `
		SELECT ma.*, u.name AS "mentorName", u.email AS "mentorEmail"
		FROM "MentorAssessment" ma
		LEFT JOIN "User" u ON ma."mentorId" = u.id
		WHERE ma."menteeId" = $1`, userId)

	// Fetch DecisionWorkflows as Subject (with Creator details)
	runQuery("DecisionWorkflowsAsSubject", `
		SELECT dw.*, u.name AS "openedByName", u.email AS "openedByEmail"
		FROM "DecisionWorkflow" dw
		LEFT JOIN "User" u ON dw."openedBy" = u.id
		WHERE dw."subjectId" = $1 AND dw."subjectType" = 'USER'`, userId)

	// Fetch DecisionWorkflows Opened By (with Subject details)
	runQuery("DecisionWorkflowsOpenedBy", `
		SELECT dw.*, u.name AS "subjectName", u.email AS "subjectEmail"
		FROM "DecisionWorkflow" dw
		LEFT JOIN "User" u ON dw."subjectId" = u.id AND dw."subjectType" = 'USER'
		WHERE dw."openedBy" = $1`, userId)
	
	// Fetch WorkflowActions with associated DecisionWorkflow and Assignee details
	runQuery("WorkflowActions", `
		SELECT wa.*, 
		       dw.type AS "workflowType", dw.status AS "workflowStatus", dw.severity AS "workflowSeverity",
		       u.name AS "assigneeName", u.email AS "assigneeEmail"
		FROM "WorkflowAction" wa
		LEFT JOIN "DecisionWorkflow" dw ON wa."workflowId" = dw.id
		LEFT JOIN "User" u ON wa."assignedTo" = u.id
		WHERE wa."assignedTo" = $1`, userId)

	// Fetch WorkflowApprovalSteps with associated DecisionWorkflow and Approver details
	runQuery("WorkflowApprovalSteps", `
		SELECT was.*, 
		       dw.type AS "workflowType", dw.status AS "workflowStatus", dw.severity AS "workflowSeverity",
		       u.name AS "approverName", u.email AS "approverEmail"
		FROM "WorkflowApprovalStep" was
		LEFT JOIN "DecisionWorkflow" dw ON was."workflowId" = dw.id
		LEFT JOIN "User" u ON was."approverId" = u.id
		WHERE was."approverId" = $1`, userId)

	// Fetch WeeklySnapshots
	runQuery("WeeklySnapshots", `
		SELECT ws.*, u.name AS "userName", u.email AS "userEmail"
		FROM "WeeklySnapshot" ws
		LEFT JOIN "User" u ON ws."userId" = u.id
		WHERE ws."userId" = $1`, userId)

	// Fetch ReportDrafts (with Recipient details)
	runQuery("ReportDrafts", `
		SELECT rd.*, u.name AS "recipientName", u.email AS "recipientEmail"
		FROM "ReportDraft" rd
		LEFT JOIN "User" u ON rd."recipientUserId" = u.id
		WHERE rd."recipientUserId" = $1`, userId)

	// Fetch FlagResponses as Mentor
	runQuery("FlagResponsesAsMentor", `
		SELECT fr.*, u.name AS "mentorName"
		FROM "FlagResponse" fr
		LEFT JOIN "User" u ON fr."mentorEmail" = u.email
		WHERE fr."mentorEmail" = $1`, userEmail)

	// Fetch FlagResponses as Intern
	runQuery("FlagResponsesAsIntern", `
		SELECT fr.*, u.email AS "internEmail"
		FROM "FlagResponse" fr
		LEFT JOIN "User" u ON fr."internName" = u.name
		WHERE fr."internName" = $1`, userName)

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

