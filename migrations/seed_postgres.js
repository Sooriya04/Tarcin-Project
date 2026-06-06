const { Client } = require('pg');

const postgresDbUrl = 'postgres://postgres:localhost@localhost:5432/postgres';
const tarcinDbUrl = 'postgres://postgres:localhost@localhost:5432/tarcin';

async function ensureDatabaseExists() {
  console.log('Connecting to default postgres database to verify if "tarcin" database exists...');
  const client = new Client({ connectionString: postgresDbUrl });
  try {
    await client.connect();
    
    // Check if the database already exists
    const checkRes = await client.query("SELECT 1 FROM pg_database WHERE datname = 'tarcin'");
    if (checkRes.rowCount === 0) {
      console.log('Database "tarcin" does not exist. Creating it now...');
      // CREATE DATABASE cannot run in a transaction, so we run it directly
      await client.query('CREATE DATABASE tarcin');
      console.log('Database "tarcin" created successfully.');
    } else {
      console.log('Database "tarcin" already exists.');
    }
  } catch (err) {
    console.error('Error ensuring database exists:', err);
    throw err;
  } finally {
    await client.end();
  }
}

async function seed() {
  await ensureDatabaseExists();

  const client = new Client({ connectionString: tarcinDbUrl });
  try {
    console.log('Connecting to "tarcin" database...');
    await client.connect();
    console.log('Connected successfully!');

    // Begin Transaction
    await client.query('BEGIN');

    // 2. Drop all tables in reverse order of dependencies
    console.log('Cleaning/dropping existing tables...');
    const tablesToDrop = [
      'FlagResponse', 'ReportDraft', 'WeeklySnapshot', 'WorkflowApprovalStep',
      'WorkflowAction', 'DecisionWorkflow', 'MentorAssessment', 'WorkLog',
      'Message', 'AILog', 'DailyReport', 'XPLog', 'Streak', 'Achievement',
      'Meeting', 'Evaluation', 'Notification', 'Certificate', 'Attendance',
      'Task', 'Mentor', 'InternProfile', 'User', 'College', 'Domain'
    ];
    for (const table of tablesToDrop) {
      await client.query(`DROP TABLE IF EXISTS "${table}" CASCADE;`);
    }
    console.log('Existing tables dropped.');

    // 3. Create all tables
    console.log('Creating database tables...');

    await client.query(`
      CREATE TABLE "Domain" (
        "id"         TEXT NOT NULL PRIMARY KEY,
        "domainName" TEXT NOT NULL UNIQUE
      );
    `);

    await client.query(`
      CREATE TABLE "College" (
        "id"         TEXT NOT NULL PRIMARY KEY,
        "name"       TEXT NOT NULL UNIQUE,
        "department" TEXT,
        "hodName"    TEXT,
        "hodEmail"   TEXT,
        "city"       TEXT,
        "state"      TEXT,
        "createdAt"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updatedAt"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
      );
    `);

    await client.query(`
      CREATE TABLE "User" (
        "id"               TEXT NOT NULL PRIMARY KEY,
        "name"             TEXT NOT NULL,
        "email"            TEXT NOT NULL UNIQUE,
        "role"             TEXT NOT NULL DEFAULT 'INTERN' CHECK ("role" IN ('ADMIN','MENTOR','INTERN')),
        "status"           TEXT NOT NULL DEFAULT 'PENDING' CHECK ("status" IN ('PENDING','APPROVED','REJECTED')),
        "createdAt"        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updatedAt"        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "referralCode"     TEXT UNIQUE,
        "referredById"     TEXT,
        "xp"               INTEGER NOT NULL DEFAULT 0,
        "level"            INTEGER NOT NULL DEFAULT 1,
        "profileImageData" BYTEA,
        "profileImageType" TEXT,
        FOREIGN KEY ("referredById") REFERENCES "User"("id") ON DELETE SET NULL
      );
      CREATE INDEX "User_role_idx"      ON "User"("role");
      CREATE INDEX "User_status_idx"    ON "User"("status");
      CREATE INDEX "User_createdAt_idx" ON "User"("createdAt");
      CREATE INDEX "User_xp_idx"        ON "User"("xp");
    `);

    await client.query(`
      CREATE TABLE "InternProfile" (
        "id"                   TEXT NOT NULL PRIMARY KEY,
        "userId"               TEXT NOT NULL UNIQUE,
        "phone"                TEXT,
        "college"              TEXT,
        "collegeDistrict"      TEXT,
        "batchYear"            TEXT,
        "department"           TEXT,
        "skills"               TEXT NOT NULL DEFAULT '[]',
        "linkedin"             TEXT,
        "portfolioUrl"         TEXT,
        "resumeData"           BYTEA,
        "resumeType"           TEXT,
        "profileImageData"     BYTEA,
        "profileImageType"     TEXT,
        "preferredDomain"      TEXT,
        "internshipDuration"   TEXT,
        "availability"         TEXT,
        "assignedMentorId"     TEXT,
        "hodEmail"             TEXT,
        "aiPerformanceSummary" TEXT,
        "internshipEndDate"    TEXT,
        "conversionStatus"     TEXT DEFAULT 'NOT_STARTED',
        "collegeId"            TEXT,
        FOREIGN KEY ("userId")    REFERENCES "User"("id") ON DELETE CASCADE,
        FOREIGN KEY ("collegeId") REFERENCES "College"("id") ON DELETE SET NULL
      );
      CREATE INDEX "InternProfile_preferredDomain_idx" ON "InternProfile"("preferredDomain");
    `);

    await client.query(`
      CREATE TABLE "Mentor" (
        "id"       TEXT NOT NULL PRIMARY KEY,
        "userId"   TEXT NOT NULL UNIQUE,
        "domainId" TEXT,
        FOREIGN KEY ("userId")   REFERENCES "User"("id") ON DELETE CASCADE,
        FOREIGN KEY ("domainId") REFERENCES "Domain"("id") ON DELETE SET NULL
      );
    `);

    await client.query(`
      CREATE TABLE "Task" (
        "id"             TEXT NOT NULL PRIMARY KEY,
        "title"          TEXT NOT NULL,
        "description"    TEXT NOT NULL,
        "deadline"       TEXT NOT NULL,
        "status"         TEXT NOT NULL DEFAULT 'TODO' CHECK ("status" IN ('TODO','IN_PROGRESS','BLOCKERS','SUBMITTED','REVISION_NEEDED','REVISION','APPROVED','REJECTED')),
        "priority"       TEXT NOT NULL DEFAULT 'medium',
        "score"          INTEGER,
        "assignedById"   TEXT NOT NULL,
        "assignedToId"   TEXT NOT NULL,
        "submissionUrl"  TEXT,
        "submissionDate" TEXT,
        "internComment"  TEXT,
        "mentorFeedback" TEXT,
        "reviewedAt"     TEXT,
        "createdAt"      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updatedAt"      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("assignedById") REFERENCES "User"("id") ON DELETE CASCADE,
        FOREIGN KEY ("assignedToId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Task_status_idx"       ON "Task"("status");
      CREATE INDEX "Task_assignedToId_idx"  ON "Task"("assignedToId");
      CREATE INDEX "Task_assignedById_idx"  ON "Task"("assignedById");
    `);

    await client.query(`
      CREATE TABLE "Attendance" (
        "id"       TEXT NOT NULL PRIMARY KEY,
        "userId"   TEXT NOT NULL,
        "checkin"  TEXT NOT NULL,
        "checkout" TEXT,
        "date"     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Attendance_userId_idx" ON "Attendance"("userId");
      CREATE INDEX "Attendance_date_idx"   ON "Attendance"("date");
    `);

    await client.query(`
      CREATE TABLE "Certificate" (
        "id"        TEXT NOT NULL PRIMARY KEY,
        "userId"    TEXT NOT NULL,
        "issueDate" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "fileUrl"   TEXT,
        "pdfData"   BYTEA,
        FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Certificate_userId_idx" ON "Certificate"("userId");
    `);

    await client.query(`
      CREATE TABLE "Notification" (
        "id"        TEXT NOT NULL PRIMARY KEY,
        "userId"    TEXT NOT NULL,
        "message"   TEXT NOT NULL,
        "type"      TEXT NOT NULL,
        "read"      BOOLEAN NOT NULL DEFAULT false,
        "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Notification_userId_idx" ON "Notification"("userId");
      CREATE INDEX "Notification_read_idx"   ON "Notification"("read");
    `);

    await client.query(`
      CREATE TABLE "Evaluation" (
        "id"               TEXT NOT NULL PRIMARY KEY,
        "internId"         TEXT NOT NULL,
        "evaluatorId"      TEXT NOT NULL,
        "weekNumber"       INTEGER NOT NULL,
        "technicalScore"   REAL NOT NULL,
        "consistencyScore" REAL NOT NULL,
        "leadershipScore"  REAL NOT NULL,
        "initiativeScore"  REAL NOT NULL,
        "totalScore"       REAL NOT NULL,
        "feedback"         TEXT,
        "createdAt"        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updatedAt"        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        UNIQUE ("internId", "weekNumber"),
        FOREIGN KEY ("internId")    REFERENCES "User"("id") ON DELETE CASCADE,
        FOREIGN KEY ("evaluatorId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Evaluation_internId_idx"    ON "Evaluation"("internId");
      CREATE INDEX "Evaluation_evaluatorId_idx" ON "Evaluation"("evaluatorId");
    `);

    await client.query(`
      CREATE TABLE "Meeting" (
        "id"          TEXT NOT NULL PRIMARY KEY,
        "title"       TEXT NOT NULL,
        "description" TEXT,
        "date"        TEXT NOT NULL,
        "link"        TEXT NOT NULL,
        "mentorId"    TEXT NOT NULL,
        "internId"    TEXT,
        "createdAt"   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("mentorId") REFERENCES "User"("id") ON DELETE CASCADE,
        FOREIGN KEY ("internId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Meeting_mentorId_idx" ON "Meeting"("mentorId");
      CREATE INDEX "Meeting_internId_idx" ON "Meeting"("internId");
    `);

    await client.query(`
      CREATE TABLE "Achievement" (
        "id"          TEXT NOT NULL PRIMARY KEY,
        "userId"      TEXT NOT NULL,
        "type"        TEXT NOT NULL,
        "title"       TEXT NOT NULL,
        "description" TEXT NOT NULL,
        "icon"        TEXT NOT NULL,
        "xpAwarded"   INTEGER NOT NULL,
        "unlockedAt"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        UNIQUE ("userId", "type"),
        FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Achievement_userId_idx" ON "Achievement"("userId");
    `);

    await client.query(`
      CREATE TABLE "Streak" (
        "id"            TEXT NOT NULL PRIMARY KEY,
        "userId"        TEXT NOT NULL UNIQUE,
        "currentStreak" INTEGER NOT NULL DEFAULT 0,
        "longestStreak" INTEGER NOT NULL DEFAULT 0,
        "lastCheckIn"   TEXT,
        FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Streak_userId_idx" ON "Streak"("userId");
    `);

    await client.query(`
      CREATE TABLE "XPLog" (
        "id"        TEXT NOT NULL PRIMARY KEY,
        "userId"    TEXT NOT NULL,
        "amount"    INTEGER NOT NULL,
        "reason"    TEXT NOT NULL,
        "metadata"  TEXT,
        "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "XPLog_userId_idx"    ON "XPLog"("userId");
      CREATE INDEX "XPLog_createdAt_idx" ON "XPLog"("createdAt");
    `);

    await client.query(`
      CREATE TABLE "DailyReport" (
        "id"         TEXT NOT NULL PRIMARY KEY,
        "internId"   TEXT NOT NULL,
        "date"       TEXT NOT NULL,
        "activities" TEXT NOT NULL,
        "timeSpent"  INTEGER NOT NULL,
        "blockers"   TEXT,
        "createdAt"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updatedAt"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("internId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "DailyReport_internId_idx" ON "DailyReport"("internId");
      CREATE INDEX "DailyReport_date_idx"     ON "DailyReport"("date");
    `);

    await client.query(`
      CREATE TABLE "AILog" (
        "id"             TEXT NOT NULL PRIMARY KEY,
        "adminId"        TEXT NOT NULL,
        "targetInternId" TEXT,
        "action"         TEXT NOT NULL,
        "promptTokens"   INTEGER,
        "responseTokens" INTEGER,
        "createdAt"      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("adminId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "AILog_adminId_idx"   ON "AILog"("adminId");
      CREATE INDEX "AILog_createdAt_idx" ON "AILog"("createdAt");
    `);

    await client.query(`
      CREATE TABLE "Message" (
        "id"         TEXT NOT NULL PRIMARY KEY,
        "content"    TEXT NOT NULL,
        "senderId"   TEXT NOT NULL,
        "receiverId" TEXT NOT NULL,
        "read"       BOOLEAN NOT NULL DEFAULT false,
        "createdAt"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("senderId")   REFERENCES "User"("id") ON DELETE CASCADE,
        FOREIGN KEY ("receiverId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "Message_senderId_idx"   ON "Message"("senderId");
      CREATE INDEX "Message_receiverId_idx" ON "Message"("receiverId");
      CREATE INDEX "Message_createdAt_idx"  ON "Message"("createdAt");
    `);

    await client.query(`
      CREATE TABLE "WorkLog" (
        "id"          TEXT NOT NULL PRIMARY KEY,
        "userId"      TEXT NOT NULL,
        "logDate"     TEXT NOT NULL,
        "taskItems"   TEXT NOT NULL DEFAULT '[]',
        "daySummary"  TEXT,
        "hasEvidence" BOOLEAN NOT NULL DEFAULT false,
        "hasBlocker"  BOOLEAN NOT NULL DEFAULT false,
        "isLate"      BOOLEAN NOT NULL DEFAULT false,
        "createdAt"   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updatedAt"   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        UNIQUE ("userId", "logDate"),
        FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "WorkLog_userId_idx"         ON "WorkLog"("userId");
      CREATE INDEX "WorkLog_logDate_idx"        ON "WorkLog"("logDate");
      CREATE INDEX "WorkLog_userId_logDate_idx" ON "WorkLog"("userId","logDate");
    `);

    await client.query(`
      CREATE TABLE "MentorAssessment" (
        "id"           TEXT NOT NULL PRIMARY KEY,
        "mentorId"     TEXT NOT NULL,
        "menteeId"     TEXT NOT NULL,
        "weekId"       TEXT NOT NULL,
        "rating"       TEXT NOT NULL,
        "comment"      TEXT,
        "attestations" TEXT NOT NULL DEFAULT '[]',
        "createdAt"    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updatedAt"    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        UNIQUE ("mentorId", "menteeId", "weekId"),
        FOREIGN KEY ("mentorId") REFERENCES "User"("id") ON DELETE CASCADE,
        FOREIGN KEY ("menteeId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "MentorAssessment_menteeId_idx" ON "MentorAssessment"("menteeId");
      CREATE INDEX "MentorAssessment_weekId_idx"   ON "MentorAssessment"("weekId");
      CREATE INDEX "MentorAssessment_mentorId_idx" ON "MentorAssessment"("mentorId");
    `);

    await client.query(`
      CREATE TABLE "DecisionWorkflow" (
        "id"               TEXT NOT NULL PRIMARY KEY,
        "type"             TEXT NOT NULL,
        "subjectType"      TEXT NOT NULL,
        "subjectId"        TEXT,
        "status"           TEXT NOT NULL DEFAULT 'OPEN',
        "severity"         TEXT,
        "proposedDecision" TEXT,
        "outcome"          TEXT,
        "openedBy"         TEXT NOT NULL,
        "openedAt"         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "closedAt"         TEXT,
        "dueDate"          TEXT,
        "weekId"           TEXT
      );
      CREATE INDEX "DecisionWorkflow_type_idx"      ON "DecisionWorkflow"("type");
      CREATE INDEX "DecisionWorkflow_status_idx"    ON "DecisionWorkflow"("status");
      CREATE INDEX "DecisionWorkflow_subjectId_idx" ON "DecisionWorkflow"("subjectId");
      CREATE INDEX "DecisionWorkflow_weekId_idx"    ON "DecisionWorkflow"("weekId");
    `);

    await client.query(`
      CREATE TABLE "WorkflowAction" (
        "id"          TEXT NOT NULL PRIMARY KEY,
        "workflowId"  TEXT NOT NULL,
        "description" TEXT NOT NULL,
        "assignedTo"  TEXT NOT NULL,
        "dueDate"     TEXT NOT NULL,
        "status"      TEXT NOT NULL DEFAULT 'OPEN',
        "notes"       TEXT,
        "completedAt" TEXT,
        "createdAt"   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY ("workflowId") REFERENCES "DecisionWorkflow"("id") ON DELETE CASCADE
      );
      CREATE INDEX "WorkflowAction_workflowId_idx" ON "WorkflowAction"("workflowId");
      CREATE INDEX "WorkflowAction_assignedTo_idx" ON "WorkflowAction"("assignedTo");
    `);

    await client.query(`
      CREATE TABLE "WorkflowApprovalStep" (
        "id"         TEXT NOT NULL PRIMARY KEY,
        "workflowId" TEXT NOT NULL,
        "stepNumber" INTEGER NOT NULL,
        "approverId" TEXT NOT NULL,
        "role"       TEXT NOT NULL,
        "status"     TEXT NOT NULL DEFAULT 'PENDING',
        "decision"   TEXT,
        "reason"     TEXT,
        "decidedAt"  TEXT,
        UNIQUE ("workflowId", "stepNumber"),
        FOREIGN KEY ("workflowId") REFERENCES "DecisionWorkflow"("id") ON DELETE CASCADE
      );
      CREATE INDEX "WorkflowApprovalStep_workflowId_idx" ON "WorkflowApprovalStep"("workflowId");
    `);

    await client.query(`
      CREATE TABLE "WeeklySnapshot" (
        "id"                       TEXT NOT NULL PRIMARY KEY,
        "userId"                   TEXT NOT NULL,
        "weekId"                   TEXT NOT NULL,
        "participation"            TEXT NOT NULL,
        "outputEvidence"           TEXT NOT NULL,
        "momentum"                 TEXT NOT NULL,
        "quality"                  TEXT NOT NULL,
        "signal"                   TEXT NOT NULL,
        "escalationLevel"          INTEGER NOT NULL DEFAULT 0,
        "consecutiveNegativeWeeks" INTEGER NOT NULL DEFAULT 0,
        "overrideSignal"           TEXT,
        "overrideReason"           TEXT,
        "overriddenBy"             TEXT,
        "overriddenAt"             TEXT,
        "verifiedBy"               TEXT,
        "verifiedAt"               TEXT,
        "internFacingLabel"        TEXT,
        "computedAt"               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        UNIQUE ("userId", "weekId"),
        FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
      );
      CREATE INDEX "WeeklySnapshot_userId_idx"          ON "WeeklySnapshot"("userId");
      CREATE INDEX "WeeklySnapshot_weekId_idx"          ON "WeeklySnapshot"("weekId");
      CREATE INDEX "WeeklySnapshot_signal_idx"          ON "WeeklySnapshot"("signal");
      CREATE INDEX "WeeklySnapshot_escalationLevel_idx" ON "WeeklySnapshot"("escalationLevel");
    `);

    await client.query(`
      CREATE TABLE "ReportDraft" (
        "id"              TEXT NOT NULL PRIMARY KEY,
        "weekId"          TEXT NOT NULL,
        "reportType"      TEXT NOT NULL,
        "recipientUserId" TEXT,
        "recipientName"   TEXT NOT NULL,
        "recipientEmail"  TEXT NOT NULL,
        "renderedHtml"    TEXT NOT NULL,
        "status"          TEXT NOT NULL DEFAULT 'DRAFT',
        "heldReason"      TEXT,
        "approvedBy"      TEXT,
        "approvedAt"      TEXT,
        "sentAt"          TEXT,
        "failureReason"   TEXT,
        "isDryRun"        BOOLEAN NOT NULL DEFAULT true,
        "feedbackStatus"  TEXT,
        "feedbackNotes"   TEXT,
        "feedbackBy"      TEXT,
        "createdAt"       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
      );
      CREATE INDEX "ReportDraft_weekId_idx"          ON "ReportDraft"("weekId");
      CREATE INDEX "ReportDraft_status_idx"          ON "ReportDraft"("status");
      CREATE INDEX "ReportDraft_reportType_idx"      ON "ReportDraft"("reportType");
      CREATE INDEX "ReportDraft_recipientUserId_idx" ON "ReportDraft"("recipientUserId");
    `);

    await client.query(`
      CREATE TABLE "FlagResponse" (
        "id"          TEXT NOT NULL PRIMARY KEY,
        "mentorEmail" TEXT NOT NULL,
        "internName"  TEXT NOT NULL,
        "action"      TEXT NOT NULL,
        "emailBody"   TEXT,
        "status"      TEXT NOT NULL DEFAULT 'PENDING_REVIEW',
        "createdAt"   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
      );
      CREATE INDEX "FlagResponse_status_idx" ON "FlagResponse"("status");
    `);

    console.log('Tables and indexes created successfully.');

    // 4. Seed Data
    console.log('Inserting dummy data (10+ rows per table)...');

    // --- Domain (10 rows) ---
    console.log('Seeding Domain...');
    const domains = [
      ['domain-001', 'Full Stack'],
      ['domain-002', 'AI / ML'],
      ['domain-003', 'Data Science'],
      ['domain-004', 'Cybersecurity'],
      ['domain-005', 'Cloud / DevOps'],
      ['domain-006', 'Web Development'],
      ['domain-007', 'Mobile Dev'],
      ['domain-008', 'UI/UX Design'],
      ['domain-009', 'QA Testing'],
      ['domain-010', 'Product Management']
    ];
    for (const d of domains) {
      await client.query('INSERT INTO "Domain" ("id", "domainName") VALUES ($1, $2)', d);
    }

    // --- College (10 rows) ---
    console.log('Seeding College...');
    const colleges = [
      ['college-001', 'IIT Bombay', 'Computer Science', 'Dr. A. Sharma', 'hod1@college.edu', 'Mumbai', 'Maharashtra'],
      ['college-002', 'IIT Delhi', 'Information Technology', 'Dr. B. Patel', 'hod2@college.edu', 'Delhi', 'Delhi'],
      ['college-003', 'IIT Madras', 'Software Engineering', 'Dr. C. Iyer', 'hod3@college.edu', 'Chennai', 'Tamil Nadu'],
      ['college-004', 'BITS Pilani', 'Electronics', 'Dr. D. Sen', 'hod4@college.edu', 'Pilani', 'Rajasthan'],
      ['college-005', 'NIT Trichy', 'Information Science', 'Dr. E. Reddy', 'hod5@college.edu', 'Trichy', 'Tamil Nadu'],
      ['college-006', 'VIT Vellore', 'Computer Applications', 'Dr. F. Das', 'hod6@college.edu', 'Vellore', 'Tamil Nadu'],
      ['college-007', 'SRM University', 'Software Systems', 'Dr. G. Nair', 'hod7@college.edu', 'Chennai', 'Tamil Nadu'],
      ['college-008', 'Delhi Technological University', 'Computing', 'Dr. H. Joshi', 'hod8@college.edu', 'Delhi', 'Delhi'],
      ['college-009', 'RV College of Engineering', 'Data Systems', 'Dr. I. Rao', 'hod9@college.edu', 'Bangalore', 'Karnataka'],
      ['college-010', 'PES University', 'Artificial Intelligence', 'Dr. J. Roy', 'hod10@college.edu', 'Bangalore', 'Karnataka']
    ];
    for (const c of colleges) {
      await client.query('INSERT INTO "College" ("id", "name", "department", "hodName", "hodEmail", "city", "state") VALUES ($1, $2, $3, $4, $5, $6, $7)', c);
    }

    // --- User (25 rows) ---
    console.log('Seeding User...');
    const users = [
      // Admins (5 rows)
      ['user-admin-1', 'Admin Alice', 'admin1@tarcin.com', 'ADMIN', 'APPROVED', 'REF-A1', null, 0, 1],
      ['user-admin-2', 'Admin Bob', 'admin2@tarcin.com', 'ADMIN', 'APPROVED', 'REF-A2', null, 0, 1],
      ['user-admin-3', 'Admin Charlie', 'admin3@tarcin.com', 'ADMIN', 'APPROVED', 'REF-A3', null, 0, 1],
      ['user-admin-4', 'Admin David', 'admin4@tarcin.com', 'ADMIN', 'APPROVED', 'REF-A4', null, 0, 1],
      ['user-admin-5', 'Admin Eva', 'admin5@tarcin.com', 'ADMIN', 'APPROVED', 'REF-A5', null, 0, 1],
      // Mentors (10 rows)
      ['user-mentor-1', 'Mentor Frank', 'mentor1@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M1', null, 0, 1],
      ['user-mentor-2', 'Mentor Grace', 'mentor2@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M2', null, 0, 1],
      ['user-mentor-3', 'Mentor Heidi', 'mentor3@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M3', null, 0, 1],
      ['user-mentor-4', 'Mentor Ivan', 'mentor4@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M4', null, 0, 1],
      ['user-mentor-5', 'Mentor Judy', 'mentor5@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M5', null, 0, 1],
      ['user-mentor-6', 'Mentor Ken', 'mentor6@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M6', null, 0, 1],
      ['user-mentor-7', 'Mentor Leo', 'mentor7@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M7', null, 0, 1],
      ['user-mentor-8', 'Mentor Mallory', 'mentor8@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M8', null, 0, 1],
      ['user-mentor-9', 'Mentor Niaj', 'mentor9@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M9', null, 0, 1],
      ['user-mentor-10', 'Mentor Olivia', 'mentor10@tarcin.com', 'MENTOR', 'APPROVED', 'REF-M10', null, 0, 1],
      // Interns (10 rows)
      ['user-intern-1', 'Intern Paul', 'intern1@tarcin.com', 'INTERN', 'APPROVED', 'REF-I1', null, 150, 2],
      ['user-intern-2', 'Intern Peggy', 'intern2@tarcin.com', 'INTERN', 'APPROVED', 'REF-I2', 'user-intern-1', 250, 3],
      ['user-intern-3', 'Intern Sybil', 'intern3@tarcin.com', 'INTERN', 'APPROVED', 'REF-I3', 'user-intern-2', 50, 1],
      ['user-intern-4', 'Intern Ted', 'intern4@tarcin.com', 'INTERN', 'APPROVED', 'REF-I4', null, 300, 4],
      ['user-intern-5', 'Intern Victor', 'intern5@tarcin.com', 'INTERN', 'APPROVED', 'REF-I5', null, 400, 5],
      ['user-intern-6', 'Intern Wendy', 'intern6@tarcin.com', 'INTERN', 'APPROVED', 'REF-I6', null, 120, 2],
      ['user-intern-7', 'Intern Walter', 'intern7@tarcin.com', 'INTERN', 'APPROVED', 'REF-I7', null, 80, 1],
      ['user-intern-8', 'Intern Xenia', 'intern8@tarcin.com', 'INTERN', 'APPROVED', 'REF-I8', null, 220, 3],
      ['user-intern-9', 'Intern Yuri', 'intern9@tarcin.com', 'INTERN', 'PENDING', 'REF-I9', null, 0, 1],
      ['user-intern-10', 'Intern Zoe', 'intern10@tarcin.com', 'INTERN', 'REJECTED', 'REF-I10', null, 0, 1]
    ];
    for (const u of users) {
      await client.query(
        `INSERT INTO "User" (
          "id", "name", "email", "role", "status", "referralCode", "referredById", "xp", "level"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
        u
      );
    }

    // --- InternProfile (10 rows) ---
    console.log('Seeding InternProfile...');
    const internProfiles = [
      ['profile-001', 'user-intern-1', '9876543210', 'IIT Bombay', 'Mumbai', '2026', 'Computer Science', '["React", "Node.js"]', 'https://linkedin.com/in/paul', 'https://paul.dev', 'Full Stack', '6 months', 'Full Time', 'user-mentor-1', 'hod1@college.edu', 'Excellent developer', '2026-11-30', 'NOT_STARTED', 'college-001'],
      ['profile-002', 'user-intern-2', '9876543211', 'IIT Delhi', 'Delhi', '2026', 'Information Technology', '["Python", "Django"]', 'https://linkedin.com/in/peggy', 'https://peggy.dev', 'AI / ML', '6 months', 'Full Time', 'user-mentor-2', 'hod2@college.edu', 'Great research orientation', '2026-11-30', 'NOT_STARTED', 'college-002'],
      ['profile-003', 'user-intern-3', '9876543212', 'IIT Madras', 'Chennai', '2026', 'Software Engineering', '["Go", "Docker"]', 'https://linkedin.com/in/sybil', 'https://sybil.dev', 'Cloud / DevOps', '6 months', 'Part Time', 'user-mentor-3', 'hod3@college.edu', 'Good infrastructure skills', '2026-11-30', 'NOT_STARTED', 'college-003'],
      ['profile-004', 'user-intern-4', '9876543213', 'BITS Pilani', 'Pilani', '2026', 'Electronics', '["C++", "Embedded Systems"]', 'https://linkedin.com/in/ted', 'https://ted.dev', 'AI / ML', '3 months', 'Full Time', 'user-mentor-4', 'hod4@college.edu', 'Strong math background', '2026-08-31', 'NOT_STARTED', 'college-004'],
      ['profile-005', 'user-intern-5', '9876543214', 'NIT Trichy', 'Trichy', '2026', 'Information Science', '["Kotlin", "Android"]', 'https://linkedin.com/in/victor', 'https://victor.dev', 'Mobile Dev', '6 months', 'Full Time', 'user-mentor-5', 'hod5@college.edu', 'Quality Android code', '2026-11-30', 'NOT_STARTED', 'college-005'],
      ['profile-006', 'user-intern-6', '9876543215', 'VIT Vellore', 'Vellore', '2026', 'Computer Applications', '["HTML", "CSS", "Figma"]', 'https://linkedin.com/in/wendy', 'https://wendy.dev', 'UI/UX Design', '6 months', 'Full Time', 'user-mentor-6', 'hod6@college.edu', 'Great visual aesthetics', '2026-11-30', 'NOT_STARTED', 'college-006'],
      ['profile-007', 'user-intern-7', '9876543216', 'SRM University', 'Chennai', '2026', 'Software Systems', '["Java", "Spring Boot"]', 'https://linkedin.com/in/walter', 'https://walter.dev', 'Full Stack', '6 months', 'Part Time', 'user-mentor-7', 'hod7@college.edu', 'Solid backend understanding', '2026-11-30', 'NOT_STARTED', 'college-007'],
      ['profile-008', 'user-intern-8', '9876543217', 'Delhi Technological University', 'Delhi', '2026', 'Computing', '["SQL", "Tableau"]', 'https://linkedin.com/in/xenia', 'https://xenia.dev', 'Data Science', '6 months', 'Full Time', 'user-mentor-8', 'hod8@college.edu', 'Accurate data modeling', '2026-11-30', 'NOT_STARTED', 'college-008'],
      ['profile-009', 'user-intern-9', '9876543218', 'RV College of Engineering', 'Bangalore', '2026', 'Data Systems', '["Typescript", "React Native"]', 'https://linkedin.com/in/yuri', 'https://yuri.dev', 'Web Development', '6 months', 'Full Time', 'user-mentor-9', 'hod9@college.edu', 'Fast learning ability', '2026-11-30', 'NOT_STARTED', 'college-009'],
      ['profile-010', 'user-intern-10', '9876543219', 'PES University', 'Bangalore', '2026', 'Artificial Intelligence', '["Pytorch", "NLP"]', 'https://linkedin.com/in/zoe', 'https://zoe.dev', 'AI / ML', '6 months', 'Full Time', 'user-mentor-10', 'hod10@college.edu', 'Requires some revision on fundamentals', '2026-11-30', 'NOT_STARTED', 'college-010']
    ];
    for (const ip of internProfiles) {
      await client.query(
        `INSERT INTO "InternProfile" (
          "id", "userId", "phone", "college", "collegeDistrict", "batchYear", "department", "skills", "linkedin", "portfolioUrl",
          "preferredDomain", "internshipDuration", "availability", "assignedMentorId", "hodEmail", "aiPerformanceSummary",
          "internshipEndDate", "conversionStatus", "collegeId"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
        ip
      );
    }

    // --- Mentor (10 rows) ---
    console.log('Seeding Mentor...');
    const mentors = [
      ['mentor-profile-001', 'user-mentor-1', 'domain-001'],
      ['mentor-profile-002', 'user-mentor-2', 'domain-002'],
      ['mentor-profile-003', 'user-mentor-3', 'domain-003'],
      ['mentor-profile-004', 'user-mentor-4', 'domain-004'],
      ['mentor-profile-005', 'user-mentor-5', 'domain-005'],
      ['mentor-profile-006', 'user-mentor-6', 'domain-006'],
      ['mentor-profile-007', 'user-mentor-7', 'domain-007'],
      ['mentor-profile-008', 'user-mentor-8', 'domain-008'],
      ['mentor-profile-009', 'user-mentor-9', 'domain-009'],
      ['mentor-profile-010', 'user-mentor-10', 'domain-010']
    ];
    for (const m of mentors) {
      await client.query('INSERT INTO "Mentor" ("id", "userId", "domainId") VALUES ($1, $2, $3)', m);
    }

    // --- Task (10 rows) ---
    console.log('Seeding Task...');
    const tasks = [
      ['task-001', 'Workspace Setup', 'Configure IDE, database connections, and Git repositories.', '2026-06-15', 'APPROVED', 'medium', 100, 'user-mentor-1', 'user-intern-1', 'https://github.com/intern1/tarcin', '2026-06-12', 'Done with setup', 'Great execution.', '2026-06-13'],
      ['task-002', 'API Refactoring', 'Migrate endpoints from Express callback to async/await.', '2026-06-15', 'TODO', 'high', null, 'user-mentor-2', 'user-intern-2', null, null, null, null, null],
      ['task-003', 'Database Optimization', 'Add indexing for foreign keys and common query paths.', '2026-06-16', 'IN_PROGRESS', 'medium', null, 'user-mentor-3', 'user-intern-3', null, null, null, null, null],
      ['task-004', 'UI Component Design', 'Develop re-usable input forms with proper Tailwind state styles.', '2026-06-18', 'SUBMITTED', 'low', null, 'user-mentor-4', 'user-intern-4', 'https://github.com/intern4/tarcin/pull/1', '2026-06-14', 'Please review.', null, null],
      ['task-005', 'Security Audit', 'Scan dependencies for known vulnerabilities and upgrade outdated packages.', '2026-06-20', 'BLOCKERS', 'high', null, 'user-mentor-5', 'user-intern-5', null, null, 'Facing issues with package locks', null, null],
      ['task-006', 'Integration Tests', 'Write unit and integration tests using Jest and Supertest.', '2026-06-22', 'REVISION_NEEDED', 'medium', null, 'user-mentor-6', 'user-intern-6', 'https://github.com/intern6/tarcin', '2026-06-13', 'Tests written', 'Need to add edge cases.', '2026-06-14'],
      ['task-007', 'Deploy Script', 'Set up GitHub Actions to auto-deploy to dev instance.', '2026-06-25', 'REVISION', 'medium', null, 'user-mentor-7', 'user-intern-7', 'https://github.com/intern7/tarcin', '2026-06-14', 'Revised scripts', null, null],
      ['task-008', 'SEO Analytics', 'Configure Google Analytics tag and verify routing triggers.', '2026-06-30', 'REJECTED', 'low', 30, 'user-mentor-8', 'user-intern-8', 'https://github.com/intern8/tarcin', '2026-06-12', 'SEO setup', 'Incorrect tracking ID, resubmit.', '2026-06-13'],
      ['task-009', 'Documentation', 'Write technical specs for APIs and database design.', '2026-07-02', 'TODO', 'low', null, 'user-mentor-9', 'user-intern-9', null, null, null, null, null],
      ['task-010', 'Final Showcase', 'Present live application demo to mentors.', '2026-07-10', 'TODO', 'high', null, 'user-mentor-10', 'user-intern-10', null, null, null, null, null]
    ];
    for (const t of tasks) {
      await client.query(
        `INSERT INTO "Task" (
          "id", "title", "description", "deadline", "status", "priority", "score", "assignedById", "assignedToId",
          "submissionUrl", "submissionDate", "internComment", "mentorFeedback", "reviewedAt"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
        t
      );
    }

    // --- Attendance (10 rows) ---
    console.log('Seeding Attendance...');
    const attendance = [
      ['att-001', 'user-intern-1', '09:00:00', '18:00:00'],
      ['att-002', 'user-intern-2', '09:15:00', '18:15:00'],
      ['att-003', 'user-intern-3', '08:50:00', '17:50:00'],
      ['att-004', 'user-intern-4', '09:00:00', '18:00:00'],
      ['att-005', 'user-intern-5', '09:05:00', '18:05:00'],
      ['att-006', 'user-intern-6', '09:00:00', '18:00:00'],
      ['att-007', 'user-intern-7', '09:10:00', '18:10:00'],
      ['att-008', 'user-intern-8', '08:55:00', '17:55:00'],
      ['att-009', 'user-intern-9', '09:00:00', '18:00:00'],
      ['att-010', 'user-intern-10', '09:20:00', '18:20:00']
    ];
    for (const a of attendance) {
      await client.query('INSERT INTO "Attendance" ("id", "userId", "checkin", "checkout") VALUES ($1, $2, $3, $4)', a);
    }

    // --- Certificate (10 rows) ---
    console.log('Seeding Certificate...');
    const certificates = [
      ['cert-001', 'user-intern-1', 'http://tarcin.com/certs/cert-001.pdf'],
      ['cert-002', 'user-intern-2', 'http://tarcin.com/certs/cert-002.pdf'],
      ['cert-003', 'user-intern-3', 'http://tarcin.com/certs/cert-003.pdf'],
      ['cert-004', 'user-intern-4', 'http://tarcin.com/certs/cert-004.pdf'],
      ['cert-005', 'user-intern-5', 'http://tarcin.com/certs/cert-005.pdf'],
      ['cert-006', 'user-intern-6', 'http://tarcin.com/certs/cert-006.pdf'],
      ['cert-007', 'user-intern-7', 'http://tarcin.com/certs/cert-007.pdf'],
      ['cert-008', 'user-intern-8', 'http://tarcin.com/certs/cert-008.pdf'],
      ['cert-009', 'user-intern-9', 'http://tarcin.com/certs/cert-009.pdf'],
      ['cert-010', 'user-intern-10', 'http://tarcin.com/certs/cert-010.pdf']
    ];
    for (const cert of certificates) {
      await client.query('INSERT INTO "Certificate" ("id", "userId", "fileUrl") VALUES ($1, $2, $3)', cert);
    }

    // --- Notification (10 rows) ---
    console.log('Seeding Notification...');
    const notifications = [
      ['notif-001', 'user-intern-1', 'Your workspace task has been approved.', 'TASK_APPROVED', false],
      ['notif-002', 'user-intern-2', 'You have a new task assigned: API Refactoring.', 'TASK_ASSIGNED', false],
      ['notif-003', 'user-intern-3', 'Your weekly evaluation has been posted.', 'EVALUATION_POSTED', true],
      ['notif-004', 'user-intern-4', 'Meeting scheduled: Weekly Sync with Mentor.', 'MEETING_SCHEDULED', false],
      ['notif-005', 'user-intern-5', 'Please submit your daily activity report.', 'REMINDER', false],
      ['notif-006', 'user-intern-6', 'Your test cases task needs revision.', 'TASK_REVISION', true],
      ['notif-007', 'user-intern-7', 'Welcome to the Tarcin Internship portal!', 'WELCOME', false],
      ['notif-008', 'user-intern-8', 'SEO task submitted. Awaiting review.', 'TASK_SUBMITTED', false],
      ['notif-009', 'user-intern-9', 'A meeting link has been updated.', 'MEETING_UPDATED', false],
      ['notif-010', 'user-intern-10', 'Evaluation complete for week 1.', 'EVALUATION_POSTED', true]
    ];
    for (const n of notifications) {
      await client.query('INSERT INTO "Notification" ("id", "userId", "message", "type", "read") VALUES ($1, $2, $3, $4, $5)', n);
    }

    // --- Evaluation (10 rows) ---
    console.log('Seeding Evaluation...');
    const evaluations = [
      ['eval-001', 'user-intern-1', 'user-mentor-1', 1, 4.5, 4.8, 4.0, 4.5, 17.8, 'Very structured workspace and good start.'],
      ['eval-002', 'user-intern-2', 'user-mentor-2', 1, 4.0, 4.2, 4.2, 4.0, 16.4, 'Approaches the task details very systematically.'],
      ['eval-003', 'user-intern-3', 'user-mentor-3', 1, 3.8, 4.0, 3.5, 4.2, 15.5, 'Quick learner but needs to align with guidelines.'],
      ['eval-004', 'user-intern-4', 'user-mentor-4', 1, 4.8, 4.5, 4.5, 4.8, 18.6, 'Outstanding logic and code format.'],
      ['eval-005', 'user-intern-5', 'user-mentor-5', 1, 3.5, 3.8, 3.0, 3.5, 13.8, 'Good progress but communicates less on blockers.'],
      ['eval-006', 'user-intern-6', 'user-mentor-6', 1, 4.2, 4.0, 4.0, 4.5, 16.7, 'Excellent eye for design details.'],
      ['eval-007', 'user-intern-7', 'user-mentor-7', 1, 4.0, 4.2, 3.8, 4.0, 16.0, 'Solid performance in backend configurations.'],
      ['eval-008', 'user-intern-8', 'user-mentor-8', 1, 3.0, 3.5, 3.0, 3.2, 12.7, 'Needs to check configuration requirements closely.'],
      ['eval-009', 'user-intern-9', 'user-mentor-9', 1, 4.2, 4.5, 4.0, 4.5, 17.2, 'Very receptive to feedback and acts on it quickly.'],
      ['eval-010', 'user-intern-10', 'user-mentor-10', 1, 2.5, 3.0, 2.0, 2.5, 10.0, 'NEEDS IMPROVEMENT. Needs revision on basic javascript.']
    ];
    for (const ev of evaluations) {
      await client.query(
        `INSERT INTO "Evaluation" (
          "id", "internId", "evaluatorId", "weekNumber", "technicalScore", "consistencyScore", "leadershipScore", "initiativeScore", "totalScore", "feedback"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
        ev
      );
    }

    // --- Meeting (10 rows) ---
    console.log('Seeding Meeting...');
    const meetings = [
      ['meet-001', 'Weekly Setup review', 'Brief check-in on workspace installation.', '2026-06-08', 'https://meet.google.com/abc-def-1', 'user-mentor-1', 'user-intern-1'],
      ['meet-002', 'API Design Sync', 'Design discussions for target REST schema.', '2026-06-09', 'https://meet.google.com/abc-def-2', 'user-mentor-2', 'user-intern-2'],
      ['meet-003', 'Database review', 'Indices setup discussion.', '2026-06-10', 'https://meet.google.com/abc-def-3', 'user-mentor-3', 'user-intern-3'],
      ['meet-004', 'Component feedback', 'Going over UI sketches and inputs.', '2026-06-11', 'https://meet.google.com/abc-def-4', 'user-mentor-4', 'user-intern-4'],
      ['meet-005', 'Security blocker sync', 'Debugging dependencies conflicts.', '2026-06-12', 'https://meet.google.com/abc-def-5', 'user-mentor-5', 'user-intern-5'],
      ['meet-006', 'Integration tests sync', 'Reviewing Supertest implementation.', '2026-06-13', 'https://meet.google.com/abc-def-6', 'user-mentor-6', 'user-intern-6'],
      ['meet-007', 'CI/CD pipeline help', 'Helping configure GitHub Actions secrets.', '2026-06-14', 'https://meet.google.com/abc-def-7', 'user-mentor-7', 'user-intern-7'],
      ['meet-008', 'SEO revision talk', 'Correcting Google tag scripts.', '2026-06-15', 'https://meet.google.com/abc-def-8', 'user-mentor-8', 'user-intern-8'],
      ['meet-009', 'API spec finalization', 'Making sure Swagger templates align.', '2026-06-16', 'https://meet.google.com/abc-def-9', 'user-mentor-9', 'user-intern-9'],
      ['meet-010', 'Week 1 summary review', 'Overall summary and target alignment.', '2026-06-17', 'https://meet.google.com/abc-def-10', 'user-mentor-10', 'user-intern-10']
    ];
    for (const m of meetings) {
      await client.query(
        `INSERT INTO "Meeting" (
          "id", "title", "description", "date", "link", "mentorId", "internId"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
        m
      );
    }

    // --- Achievement (10 rows) ---
    console.log('Seeding Achievement...');
    const achievements = [
      ['ach-001', 'user-intern-1', 'FIRST_SUBMISSION', 'Fast Starter', 'Submitted first task within 24 hours', 'rocket', 100],
      ['ach-002', 'user-intern-2', 'STREAK_5', 'Persistent Learner', 'Active daily check-in streak of 5 days', 'fire', 150],
      ['ach-003', 'user-intern-3', 'XP_100', 'Centurion', 'Earned 100 XP overall', 'medal', 100],
      ['ach-004', 'user-intern-4', 'PERFECT_EVAL', 'Golden Standard', 'Scored maximum points in technical review', 'star', 200],
      ['ach-005', 'user-intern-5', 'FIRST_SUBMISSION', 'Quick Start', 'Completed initial setup task', 'rocket', 100],
      ['ach-006', 'user-intern-6', 'STREAK_5', 'Dedicated Developer', 'Maintained 5 day workspace commits', 'fire', 150],
      ['ach-007', 'user-intern-7', 'XP_100', 'Accumulator', 'Reached 100 XP milestone', 'medal', 100],
      ['ach-008', 'user-intern-8', 'BUG_FINDER', 'Hawk Eye', 'Reported schema bug in baseline file', 'bug', 200],
      ['ach-009', 'user-intern-9', 'FIRST_SUBMISSION', 'Punctual submission', 'Submitted task within deadline', 'rocket', 100],
      ['ach-010', 'user-intern-10', 'STREAK_5', 'Regular Check-in', 'Kept active check-in streak', 'fire', 150]
    ];
    for (const ac of achievements) {
      await client.query(
        `INSERT INTO "Achievement" (
          "id", "userId", "type", "title", "description", "icon", "xpAwarded"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
        ac
      );
    }

    // --- Streak (10 rows) ---
    console.log('Seeding Streak...');
    const streaks = [
      ['streak-001', 'user-intern-1', 4, 10, '2026-06-05'],
      ['streak-002', 'user-intern-2', 5, 5, '2026-06-05'],
      ['streak-003', 'user-intern-3', 2, 8, '2026-06-04'],
      ['streak-004', 'user-intern-4', 7, 12, '2026-06-05'],
      ['streak-005', 'user-intern-5', 1, 4, '2026-06-05'],
      ['streak-006', 'user-intern-6', 5, 9, '2026-06-05'],
      ['streak-007', 'user-intern-7', 3, 6, '2026-06-05'],
      ['streak-008', 'user-intern-8', 0, 4, '2026-06-03'],
      ['streak-009', 'user-intern-9', 4, 4, '2026-06-05'],
      ['streak-010', 'user-intern-10', 1, 3, '2026-06-05']
    ];
    for (const s of streaks) {
      await client.query(
        `INSERT INTO "Streak" (
          "id", "userId", "currentStreak", "longestStreak", "lastCheckIn"
        ) VALUES ($1, $2, $3, $4, $5)`,
        s
      );
    }

    // --- XPLog (10 rows) ---
    console.log('Seeding XPLog...');
    const xplogs = [
      ['xp-001', 'user-intern-1', 50, 'Workspace Setup approval', '{"taskId": "task-001"}'],
      ['xp-002', 'user-intern-2', 100, 'Streak milestone', '{"days": 5}'],
      ['xp-003', 'user-intern-3', 50, 'Daily check-in', '{}'],
      ['xp-004', 'user-intern-4', 150, 'High quality code submission', '{"taskId": "task-004"}'],
      ['xp-005', 'user-intern-5', 50, 'Workspace setup', '{"taskId": "task-005"}'],
      ['xp-006', 'user-intern-6', 50, 'First test scripts', '{"taskId": "task-006"}'],
      ['xp-007', 'user-intern-7', 50, 'CI pipeline configure', '{"taskId": "task-007"}'],
      ['xp-008', 'user-intern-8', 50, 'Bug submission', '{"issue": "Schema syntax"}'],
      ['xp-009', 'user-intern-9', 50, 'Weekly review meeting completion', '{}'],
      ['xp-010', 'user-intern-10', 20, 'Daily attendance', '{}']
    ];
    for (const x of xplogs) {
      await client.query(
        `INSERT INTO "XPLog" (
          "id", "userId", "amount", "reason", "metadata"
        ) VALUES ($1, $2, $3, $4, $5)`,
        x
      );
    }

    // --- DailyReport (10 rows) ---
    console.log('Seeding DailyReport...');
    const dailyReports = [
      ['report-001', 'user-intern-1', '2026-06-05', 'Configured pg connection and seeded tables.', 480, 'None'],
      ['report-002', 'user-intern-2', '2026-06-05', 'Refactored user profiles mock script.', 360, 'Faced NPM sync issues.'],
      ['report-003', 'user-intern-3', '2026-06-05', 'Wrote database creation functions.', 420, 'None'],
      ['report-004', 'user-intern-4', '2026-06-05', 'Created react components and layouts.', 500, 'None'],
      ['report-005', 'user-intern-5', '2026-06-05', 'Upgraded dependencies to fix audit warnings.', 300, 'Tailwind v4 upgrade issues.'],
      ['report-006', 'user-intern-6', '2026-06-05', 'Added test scenarios for task model.', 450, 'None'],
      ['report-007', 'user-intern-7', '2026-06-05', 'Tested deployment scripts on local VM.', 400, 'None'],
      ['report-008', 'user-intern-8', '2026-06-05', 'Configured metadata tagging for SEO.', 360, 'Tag script had layout shifts.'],
      ['report-009', 'user-intern-9', '2026-06-05', 'Parsed and reviewed endpoints documentation.', 410, 'None'],
      ['report-010', 'user-intern-10', '2026-06-05', 'Studied Prisma models and mapped constraints.', 240, 'Required mentorship on async actions.']
    ];
    for (const r of dailyReports) {
      await client.query(
        `INSERT INTO "DailyReport" (
          "id", "internId", "date", "activities", "timeSpent", "blockers"
        ) VALUES ($1, $2, $3, $4, $5, $6)`,
        r
      );
    }

    // --- AILog (10 rows) ---
    console.log('Seeding AILog...');
    const ailogs = [
      ['ailog-001', 'user-admin-1', 'user-intern-1', 'Generate Performance Summary', 1200, 400],
      ['ailog-002', 'user-admin-2', 'user-intern-2', 'Check Code quality pattern', 800, 300],
      ['ailog-003', 'user-admin-3', 'user-intern-3', 'Resolve query structure', 1500, 500],
      ['ailog-004', 'user-admin-4', 'user-intern-4', 'Evaluate design proposal', 900, 250],
      ['ailog-005', 'user-admin-5', 'user-intern-5', 'Dependency lock resolution', 2000, 800],
      ['ailog-006', 'user-admin-1', 'user-intern-6', 'Check unit coverage statistics', 1100, 350],
      ['ailog-007', 'user-admin-2', 'user-intern-7', 'Build pipeline debug instructions', 1400, 450],
      ['ailog-008', 'user-admin-3', 'user-intern-8', 'SEO tagging improvements suggestions', 1000, 300],
      ['ailog-009', 'user-admin-4', 'user-intern-9', 'Format Swagger documentation', 1600, 600],
      ['ailog-010', 'user-admin-5', 'user-intern-10', 'Review fundamental concepts exercises', 1300, 500]
    ];
    for (const ai of ailogs) {
      await client.query(
        `INSERT INTO "AILog" (
          "id", "adminId", "targetInternId", "action", "promptTokens", "responseTokens"
        ) VALUES ($1, $2, $3, $4, $5, $6)`,
        ai
      );
    }

    // --- Message (10 rows) ---
    console.log('Seeding Message...');
    const messages = [
      ['msg-001', 'Hi Frank, I completed the setup.', 'user-intern-1', 'user-mentor-1', true],
      ['msg-002', 'Approved it, good job.', 'user-mentor-1', 'user-intern-1', true],
      ['msg-003', 'Grace, could you guide on API paths?', 'user-intern-2', 'user-mentor-2', false],
      ['msg-004', 'Heidi, I have added indices to the logs.', 'user-intern-3', 'user-mentor-3', true],
      ['msg-005', 'Excellent, let me review it.', 'user-mentor-3', 'user-intern-3', true],
      ['msg-006', 'Ivan, component design PR is ready.', 'user-intern-4', 'user-mentor-4', false],
      ['msg-007', 'Judy, facing npm audit issues.', 'user-intern-5', 'user-mentor-5', false],
      ['msg-008', 'Ken, integration test cases are updated.', 'user-intern-6', 'user-mentor-6', true],
      ['msg-009', 'Leo, build pipeline configured.', 'user-intern-7', 'user-mentor-7', false],
      ['msg-010', 'Mallory, Google tag setup is verified.', 'user-intern-8', 'user-mentor-8', false]
    ];
    for (const msg of messages) {
      await client.query(
        `INSERT INTO "Message" (
          "id", "content", "senderId", "receiverId", "read"
        ) VALUES ($1, $2, $3, $4, $5)`,
        msg
      );
    }

    // --- WorkLog (10 rows) ---
    console.log('Seeding WorkLog...');
    const workLogs = [
      ['wl-001', 'user-intern-1', '2026-06-05', '["Setup pg", "Wrote seed scripts"]', 'Successfully configured the PostgreSQL client and prepared table schemas.', true, false, false],
      ['wl-002', 'user-intern-2', '2026-06-05', '["Wrote user endpoint mocks"]', 'Added routing support for mock files.', true, true, false],
      ['wl-003', 'user-intern-3', '2026-06-05', '["Setup database index triggers"]', 'Verified database query metrics.', true, false, false],
      ['wl-004', 'user-intern-4', '2026-06-05', '["Created buttons and form controls"]', 'Wrote responsive React elements.', true, false, false],
      ['wl-005', 'user-intern-5', '2026-06-05', '["Upgraded vulnerability lock files"]', 'Fixed dependencies warnings.', false, true, false],
      ['wl-006', 'user-intern-6', '2026-06-05', '["Added API integrations tests"]', 'Ensured coverage satisfies guidelines.', true, false, false],
      ['wl-007', 'user-intern-7', '2026-06-05', '["Configured deployment actions"]', 'CI/CD pipeline test successful.', true, false, false],
      ['wl-008', 'user-intern-8', '2026-06-05', '["Parsed tags structure"]', 'Aligned metadata values.', true, false, true],
      ['wl-009', 'user-intern-9', '2026-06-05', '["Wrote API schema specs"]', 'Prepared Swagger documentation draft.', true, false, false],
      ['wl-010', 'user-intern-10', '2026-06-05', '["Studied PostgreSQL features"]', 'Read pg documentation.', false, false, false]
    ];
    for (const wl of workLogs) {
      await client.query(
        `INSERT INTO "WorkLog" (
          "id", "userId", "logDate", "taskItems", "daySummary", "hasEvidence", "hasBlocker", "isLate"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
        wl
      );
    }

    // --- MentorAssessment (10 rows) ---
    console.log('Seeding MentorAssessment...');
    const assessments = [
      ['ma-001', 'user-mentor-1', 'user-intern-1', '2026-W22', 'Excellent', 'Very fast delivery and good code discipline.', '["Attested code structure"]'],
      ['ma-002', 'user-mentor-2', 'user-intern-2', '2026-W22', 'Good', 'Solid understanding of routing constraints.', '["Verified user handlers"]'],
      ['ma-003', 'user-mentor-3', 'user-intern-3', '2026-W22', 'Satisfactory', 'Satisfactorily configured db variables.', '["Attested db indices"]'],
      ['ma-004', 'user-mentor-4', 'user-intern-4', '2026-W22', 'Excellent', 'Outstanding interface style implementation.', '["Verified React layouts"]'],
      ['ma-005', 'user-mentor-5', 'user-intern-5', '2026-W22', 'Needs Improvement', 'Facing issues with locks, needs support.', '[]'],
      ['ma-006', 'user-mentor-6', 'user-intern-6', '2026-W22', 'Good', 'Good test cases coverage.', '["Attested tests coverage"]'],
      ['ma-007', 'user-mentor-7', 'user-intern-7', '2026-W22', 'Good', 'Successfully deployed local workflows.', '["Verified Actions configurations"]'],
      ['ma-008', 'user-mentor-8', 'user-intern-8', '2026-W22', 'Needs Improvement', 'Failed Google tag verification due to wrong ID.', '[]'],
      ['ma-009', 'user-mentor-9', 'user-intern-9', '2026-W22', 'Good', 'Clean API specs draft.', '["Attested swagger specs"]'],
      ['ma-010', 'user-mentor-10', 'user-intern-10', '2026-W22', 'Needs Improvement', 'Struggling with database concepts.', '[]']
    ];
    for (const ma of assessments) {
      await client.query(
        `INSERT INTO "MentorAssessment" (
          "id", "mentorId", "menteeId", "weekId", "rating", "comment", "attestations"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
        ma
      );
    }

    // --- DecisionWorkflow (10 rows) ---
    console.log('Seeding DecisionWorkflow...');
    const workflows = [
      ['dw-001', 'PROBATION_REVIEW', 'USER', 'user-intern-1', 'CLOSED', 'MEDIUM', 'End probation early', 'Approved', 'user-admin-1', null, '2026-06-08', '2026-W22'],
      ['dw-002', 'PERFORMANCE_ALERT', 'USER', 'user-intern-10', 'OPEN', 'HIGH', 'Initiate extra training modules', null, 'user-admin-2', null, null, '2026-W22'],
      ['dw-003', 'ROLE_CHANGE', 'USER', 'user-intern-5', 'OPEN', 'LOW', 'Assigned mentor support session', null, 'user-admin-3', null, null, '2026-W22'],
      ['dw-004', 'PERFORMANCE_ALERT', 'USER', 'user-intern-8', 'CLOSED', 'MEDIUM', 'Issue warning letter', 'Approved', 'user-admin-4', null, '2026-06-09', '2026-W22'],
      ['dw-005', 'PROBATION_REVIEW', 'USER', 'user-intern-2', 'OPEN', 'MEDIUM', 'Review final tasks submission', null, 'user-admin-1', null, null, '2026-W22'],
      ['dw-006', 'PROBATION_REVIEW', 'USER', 'user-intern-3', 'OPEN', 'LOW', 'Observe consistency', null, 'user-admin-2', null, null, '2026-W22'],
      ['dw-007', 'PROBATION_REVIEW', 'USER', 'user-intern-4', 'CLOSED', 'LOW', 'Confirm transition', 'Confirmed', 'user-admin-3', null, '2026-06-10', '2026-W22'],
      ['dw-008', 'ROLE_CHANGE', 'USER', 'user-intern-6', 'OPEN', 'LOW', 'Recommend UI Lead role', null, 'user-admin-4', null, null, '2026-W22'],
      ['dw-009', 'ROLE_CHANGE', 'USER', 'user-intern-7', 'OPEN', 'LOW', 'Observe backend tasks', null, 'user-admin-5', null, null, '2026-W22'],
      ['dw-010', 'PERFORMANCE_ALERT', 'USER', 'user-intern-9', 'OPEN', 'HIGH', 'Request hod email clarification', null, 'user-admin-1', null, null, '2026-W22']
    ];
    for (const w of workflows) {
      await client.query(
        `INSERT INTO "DecisionWorkflow" (
          "id", "type", "subjectType", "subjectId", "status", "severity", "proposedDecision", "outcome", "openedBy", "closedAt", "dueDate", "weekId"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
        w
      );
    }

    // --- WorkflowAction (10 rows) ---
    console.log('Seeding WorkflowAction...');
    const workflowActions = [
      ['wa-001', 'dw-001', 'Confirm attendance records', 'user-mentor-1', '2026-06-08', 'COMPLETED', 'All records check out', '2026-06-08'],
      ['wa-002', 'dw-002', 'Schedule JS fundamentals check', 'user-mentor-10', '2026-06-15', 'OPEN', null, null],
      ['wa-003', 'dw-003', 'Conduct one-on-one help session', 'user-mentor-5', '2026-06-15', 'OPEN', null, null],
      ['wa-004', 'dw-004', 'Draft and send formal notice', 'user-admin-4', '2026-06-10', 'COMPLETED', 'Notice sent via system email', '2026-06-09'],
      ['wa-005', 'dw-005', 'Verify final pull request commits', 'user-mentor-2', '2026-06-15', 'OPEN', null, null],
      ['wa-006', 'dw-006', 'Verify daily log consistency', 'user-mentor-3', '2026-06-15', 'OPEN', null, null],
      ['wa-007', 'dw-007', 'Issue completion voucher', 'user-admin-3', '2026-06-12', 'COMPLETED', 'Voucher issued successfully', '2026-06-10'],
      ['wa-008', 'dw-008', 'Evaluate UI portfolio layouts', 'user-mentor-6', '2026-06-15', 'OPEN', null, null],
      ['wa-009', 'dw-009', 'Verify backend logic tests', 'user-mentor-7', '2026-06-15', 'OPEN', null, null],
      ['wa-010', 'dw-010', 'Send email query to HOD', 'user-admin-1', '2026-06-15', 'OPEN', null, null]
    ];
    for (const wa of workflowActions) {
      await client.query(
        `INSERT INTO "WorkflowAction" (
          "id", "workflowId", "description", "assignedTo", "dueDate", "status", "notes", "completedAt"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
        wa
      );
    }

    // --- WorkflowApprovalStep (10 rows) ---
    console.log('Seeding WorkflowApprovalStep...');
    const approvalSteps = [
      ['was-001', 'dw-001', 1, 'user-admin-1', 'ADMIN', 'APPROVED', 'APPROVED', 'Performance is clear.', '2026-06-08'],
      ['was-002', 'dw-002', 1, 'user-admin-2', 'ADMIN', 'PENDING', null, null, null],
      ['was-003', 'dw-003', 1, 'user-admin-3', 'ADMIN', 'PENDING', null, null, null],
      ['was-004', 'dw-004', 1, 'user-admin-4', 'ADMIN', 'APPROVED', 'APPROVED', 'Warning verified.', '2026-06-09'],
      ['was-005', 'dw-005', 1, 'user-admin-1', 'ADMIN', 'PENDING', null, null, null],
      ['was-006', 'dw-006', 1, 'user-admin-2', 'ADMIN', 'PENDING', null, null, null],
      ['was-007', 'dw-007', 1, 'user-admin-3', 'ADMIN', 'APPROVED', 'APPROVED', 'Completion approved.', '2026-06-10'],
      ['was-008', 'dw-008', 1, 'user-admin-4', 'ADMIN', 'PENDING', null, null, null],
      ['was-009', 'dw-009', 1, 'user-admin-5', 'ADMIN', 'PENDING', null, null, null],
      ['was-010', 'dw-010', 1, 'user-admin-1', 'ADMIN', 'PENDING', null, null, null]
    ];
    for (const was of approvalSteps) {
      await client.query(
        `INSERT INTO "WorkflowApprovalStep" (
          "id", "workflowId", "stepNumber", "approverId", "role", "status", "decision", "reason", "decidedAt"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
        was
      );
    }

    // --- WeeklySnapshot (10 rows) ---
    console.log('Seeding WeeklySnapshot...');
    const weeklySnapshots = [
      ['ws-001', 'user-intern-1', '2026-W22', 'Consistent', 'Provided', 'Progressing', 'Substantive', 'Green', 0, 0, null, null, null, null, 'user-mentor-1', '2026-06-05', 'On Track'],
      ['ws-002', 'user-intern-2', '2026-W22', 'Consistent', 'Provided', 'Progressing', 'Substantive', 'Green', 0, 0, null, null, null, null, 'user-mentor-2', '2026-06-05', 'On Track'],
      ['ws-003', 'user-intern-3', '2026-W22', 'Consistent', 'Provided', 'Progressing', 'Substantive', 'Green', 0, 0, null, null, null, null, 'user-mentor-3', '2026-06-05', 'On Track'],
      ['ws-004', 'user-intern-4', '2026-W22', 'Consistent', 'Provided', 'Progressing', 'Substantive', 'Green', 0, 0, null, null, null, null, 'user-mentor-4', '2026-06-05', 'On Track'],
      ['ws-005', 'user-intern-5', '2026-W22', 'Partial', 'Provided', 'Stuck', 'Substantive', 'Yellow', 1, 0, null, null, null, null, 'user-mentor-5', '2026-06-05', 'Needs Support'],
      ['ws-006', 'user-intern-6', '2026-W22', 'Consistent', 'Provided', 'Progressing', 'Substantive', 'Green', 0, 0, null, null, null, null, 'user-mentor-6', '2026-06-05', 'On Track'],
      ['ws-007', 'user-intern-7', '2026-W22', 'Consistent', 'Provided', 'Progressing', 'Substantive', 'Green', 0, 0, null, null, null, null, 'user-mentor-7', '2026-06-05', 'On Track'],
      ['ws-008', 'user-intern-8', '2026-W22', 'Consistent', 'Provided', 'Progressing', 'Substantive', 'Yellow', 1, 0, null, null, null, null, 'user-mentor-8', '2026-06-05', 'Needs Revision'],
      ['ws-009', 'user-intern-9', '2026-W22', 'Consistent', 'Provided', 'Progressing', 'Substantive', 'Green', 0, 0, null, null, null, null, 'user-mentor-9', '2026-06-05', 'On Track'],
      ['ws-010', 'user-intern-10', '2026-W22', 'Missing', 'None', 'Blocked', 'Empty', 'Red', 2, 1, null, null, null, null, 'user-mentor-10', '2026-06-05', 'Critically Stuck']
    ];
    for (const ws of weeklySnapshots) {
      await client.query(
        `INSERT INTO "WeeklySnapshot" (
          "id", "userId", "weekId", "participation", "outputEvidence", "momentum", "quality", "signal",
          "escalationLevel", "consecutiveNegativeWeeks", "overrideSignal", "overrideReason", "overriddenBy",
          "overriddenAt", "verifiedBy", "verifiedAt", "internFacingLabel"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
        ws
      );
    }

    // --- ReportDraft (10 rows) ---
    console.log('Seeding ReportDraft...');
    const reportDrafts = [
      ['rd-001', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-1', 'Paul Intern', 'intern1@tarcin.com', '<h1>Performance Report: Paul</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-002', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-2', 'Peggy Intern', 'intern2@tarcin.com', '<h1>Performance Report: Peggy</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-003', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-3', 'Sybil Intern', 'intern3@tarcin.com', '<h1>Performance Report: Sybil</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-004', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-4', 'Ted Intern', 'intern4@tarcin.com', '<h1>Performance Report: Ted</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-005', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-5', 'Victor Intern', 'intern5@tarcin.com', '<h1>Performance Report: Victor</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-006', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-6', 'Wendy Intern', 'intern6@tarcin.com', '<h1>Performance Report: Wendy</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-007', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-7', 'Walter Intern', 'intern7@tarcin.com', '<h1>Performance Report: Walter</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-008', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-8', 'Xenia Intern', 'intern8@tarcin.com', '<h1>Performance Report: Xenia</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-009', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-9', 'Yuri Intern', 'intern9@tarcin.com', '<h1>Performance Report: Yuri</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null],
      ['rd-010', '2026-W22', 'WEEKLY_PERFORMANCE', 'user-intern-10', 'Zoe Intern', 'intern10@tarcin.com', '<h1>Performance Report: Zoe</h1>', 'DRAFT', null, null, null, null, null, true, null, null, null]
    ];
    for (const rd of reportDrafts) {
      await client.query(
        `INSERT INTO "ReportDraft" (
          "id", "weekId", "reportType", "recipientUserId", "recipientName", "recipientEmail", "renderedHtml",
          "status", "heldReason", "approvedBy", "approvedAt", "sentAt", "failureReason", "isDryRun",
          "feedbackStatus", "feedbackNotes", "feedbackBy"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
        rd
      );
    }

    // --- FlagResponse (10 rows) ---
    console.log('Seeding FlagResponse...');
    const flagResponses = [
      ['fr-001', 'mentor1@tarcin.com', 'Paul Intern', 'SEND_WARNING', 'Dear Paul, please improve attendance.', 'PENDING_REVIEW'],
      ['fr-002', 'mentor2@tarcin.com', 'Peggy Intern', 'NO_ACTION', null, 'PENDING_REVIEW'],
      ['fr-003', 'mentor3@tarcin.com', 'Sybil Intern', 'SEND_WARNING', 'Dear Sybil, check your task submission logs.', 'PENDING_REVIEW'],
      ['fr-004', 'mentor4@tarcin.com', 'Ted Intern', 'SEND_WARNING', 'Dear Ted, submit evidence on time.', 'PENDING_REVIEW'],
      ['fr-005', 'mentor5@tarcin.com', 'Victor Intern', 'NO_ACTION', null, 'PENDING_REVIEW'],
      ['fr-006', 'mentor6@tarcin.com', 'Wendy Intern', 'SEND_WARNING', 'Dear Wendy, clean the Figma folder structure.', 'PENDING_REVIEW'],
      ['fr-007', 'mentor7@tarcin.com', 'Walter Intern', 'NO_ACTION', null, 'PENDING_REVIEW'],
      ['fr-008', 'mentor8@tarcin.com', 'Xenia Intern', 'SEND_WARNING', 'Dear Xenia, configure correct GA tag ids.', 'PENDING_REVIEW'],
      ['fr-009', 'mentor9@tarcin.com', 'Yuri Intern', 'NO_ACTION', null, 'PENDING_REVIEW'],
      ['fr-010', 'mentor10@tarcin.com', 'Zoe Intern', 'SEND_WARNING', 'Dear Zoe, please schedule a basic check session with HOD.', 'PENDING_REVIEW']
    ];
    for (const fr of flagResponses) {
      await client.query(
        `INSERT INTO "FlagResponse" (
          "id", "mentorEmail", "internName", "action", "emailBody", "status"
        ) VALUES ($1, $2, $3, $4, $5, $6)`,
        fr
      );
    }

    // 5. Commit Transaction
    await client.query('COMMIT');
    console.log('Database cleanup and dummy seeding completed successfully!');

  } catch (error) {
    console.error('Error seeding database:', error);
    await client.query('ROLLBACK');
  } finally {
    await client.end();
    console.log('Database client closed.');
  }
}

seed();
