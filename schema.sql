--
-- PostgreSQL database dump
--

\restrict 2TODWod3bUECuZ9tI9YQTchgQvnb2ADZXchkcr6QbzBcoudWPyUgedutDmn5ysa

-- Dumped from database version 16.14 (Ubuntu 16.14-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.14 (Ubuntu 16.14-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: AILog; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."AILog" (
    id text NOT NULL,
    "adminId" text NOT NULL,
    "targetInternId" text,
    action text NOT NULL,
    "promptTokens" integer,
    "responseTokens" integer,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."AILog" OWNER TO postgres;

--
-- Name: Achievement; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Achievement" (
    id text NOT NULL,
    "userId" text NOT NULL,
    type text NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    icon text NOT NULL,
    "xpAwarded" integer NOT NULL,
    "unlockedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."Achievement" OWNER TO postgres;

--
-- Name: Attendance; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Attendance" (
    id text NOT NULL,
    "userId" text NOT NULL,
    checkin text NOT NULL,
    checkout text,
    date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."Attendance" OWNER TO postgres;

--
-- Name: Certificate; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Certificate" (
    id text NOT NULL,
    "userId" text NOT NULL,
    "issueDate" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "fileUrl" text,
    "pdfData" bytea
);


ALTER TABLE public."Certificate" OWNER TO postgres;

--
-- Name: College; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."College" (
    id text NOT NULL,
    name text NOT NULL,
    department text,
    "hodName" text,
    "hodEmail" text,
    city text,
    state text,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."College" OWNER TO postgres;

--
-- Name: DailyReport; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."DailyReport" (
    id text NOT NULL,
    "internId" text NOT NULL,
    date text NOT NULL,
    activities text NOT NULL,
    "timeSpent" integer NOT NULL,
    blockers text,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."DailyReport" OWNER TO postgres;

--
-- Name: DecisionWorkflow; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."DecisionWorkflow" (
    id text NOT NULL,
    type text NOT NULL,
    "subjectType" text NOT NULL,
    "subjectId" text,
    status text DEFAULT 'OPEN'::text NOT NULL,
    severity text,
    "proposedDecision" text,
    outcome text,
    "openedBy" text NOT NULL,
    "openedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "closedAt" text,
    "dueDate" text,
    "weekId" text
);


ALTER TABLE public."DecisionWorkflow" OWNER TO postgres;

--
-- Name: Domain; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Domain" (
    id text NOT NULL,
    "domainName" text NOT NULL
);


ALTER TABLE public."Domain" OWNER TO postgres;

--
-- Name: Evaluation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Evaluation" (
    id text NOT NULL,
    "internId" text NOT NULL,
    "evaluatorId" text NOT NULL,
    "weekNumber" integer NOT NULL,
    "technicalScore" real NOT NULL,
    "consistencyScore" real NOT NULL,
    "leadershipScore" real NOT NULL,
    "initiativeScore" real NOT NULL,
    "totalScore" real NOT NULL,
    feedback text,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."Evaluation" OWNER TO postgres;

--
-- Name: FlagResponse; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."FlagResponse" (
    id text NOT NULL,
    "mentorEmail" text NOT NULL,
    "internName" text NOT NULL,
    action text NOT NULL,
    "emailBody" text,
    status text DEFAULT 'PENDING_REVIEW'::text NOT NULL,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."FlagResponse" OWNER TO postgres;

--
-- Name: InternProfile; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."InternProfile" (
    id text NOT NULL,
    "userId" text NOT NULL,
    phone text,
    college text,
    "collegeDistrict" text,
    "batchYear" text,
    department text,
    skills text DEFAULT '[]'::text NOT NULL,
    linkedin text,
    "portfolioUrl" text,
    "resumeData" bytea,
    "resumeType" text,
    "profileImageData" bytea,
    "profileImageType" text,
    "preferredDomain" text,
    "internshipDuration" text,
    availability text,
    "assignedMentorId" text,
    "hodEmail" text,
    "aiPerformanceSummary" text,
    "internshipEndDate" text,
    "conversionStatus" text DEFAULT 'NOT_STARTED'::text,
    "collegeId" text
);


ALTER TABLE public."InternProfile" OWNER TO postgres;

--
-- Name: Meeting; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Meeting" (
    id text NOT NULL,
    title text NOT NULL,
    description text,
    date text NOT NULL,
    link text NOT NULL,
    "mentorId" text NOT NULL,
    "internId" text,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."Meeting" OWNER TO postgres;

--
-- Name: Mentor; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Mentor" (
    id text NOT NULL,
    "userId" text NOT NULL,
    "domainId" text
);


ALTER TABLE public."Mentor" OWNER TO postgres;

--
-- Name: MentorAssessment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."MentorAssessment" (
    id text NOT NULL,
    "mentorId" text NOT NULL,
    "menteeId" text NOT NULL,
    "weekId" text NOT NULL,
    rating text NOT NULL,
    comment text,
    attestations text DEFAULT '[]'::text NOT NULL,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."MentorAssessment" OWNER TO postgres;

--
-- Name: Message; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Message" (
    id text NOT NULL,
    content text NOT NULL,
    "senderId" text NOT NULL,
    "receiverId" text NOT NULL,
    read boolean DEFAULT false NOT NULL,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."Message" OWNER TO postgres;

--
-- Name: Notification; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Notification" (
    id text NOT NULL,
    "userId" text NOT NULL,
    message text NOT NULL,
    type text NOT NULL,
    read boolean DEFAULT false NOT NULL,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."Notification" OWNER TO postgres;

--
-- Name: ReportDraft; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."ReportDraft" (
    id text NOT NULL,
    "weekId" text NOT NULL,
    "reportType" text NOT NULL,
    "recipientUserId" text,
    "recipientName" text NOT NULL,
    "recipientEmail" text NOT NULL,
    "renderedHtml" text NOT NULL,
    status text DEFAULT 'DRAFT'::text NOT NULL,
    "heldReason" text,
    "approvedBy" text,
    "approvedAt" text,
    "sentAt" text,
    "failureReason" text,
    "isDryRun" boolean DEFAULT true NOT NULL,
    "feedbackStatus" text,
    "feedbackNotes" text,
    "feedbackBy" text,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."ReportDraft" OWNER TO postgres;

--
-- Name: Streak; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Streak" (
    id text NOT NULL,
    "userId" text NOT NULL,
    "currentStreak" integer DEFAULT 0 NOT NULL,
    "longestStreak" integer DEFAULT 0 NOT NULL,
    "lastCheckIn" text
);


ALTER TABLE public."Streak" OWNER TO postgres;

--
-- Name: Task; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Task" (
    id text NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    deadline text NOT NULL,
    status text DEFAULT 'TODO'::text NOT NULL,
    priority text DEFAULT 'medium'::text NOT NULL,
    score integer,
    "assignedById" text NOT NULL,
    "assignedToId" text NOT NULL,
    "submissionUrl" text,
    "submissionDate" text,
    "internComment" text,
    "mentorFeedback" text,
    "reviewedAt" text,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "Task_status_check" CHECK ((status = ANY (ARRAY['TODO'::text, 'IN_PROGRESS'::text, 'BLOCKERS'::text, 'SUBMITTED'::text, 'REVISION_NEEDED'::text, 'REVISION'::text, 'APPROVED'::text, 'REJECTED'::text])))
);


ALTER TABLE public."Task" OWNER TO postgres;

--
-- Name: User; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."User" (
    id text NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    role text DEFAULT 'INTERN'::text NOT NULL,
    status text DEFAULT 'PENDING'::text NOT NULL,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "referralCode" text,
    "referredById" text,
    xp integer DEFAULT 0 NOT NULL,
    level integer DEFAULT 1 NOT NULL,
    "profileImageData" bytea,
    "profileImageType" text,
    CONSTRAINT "User_role_check" CHECK ((role = ANY (ARRAY['ADMIN'::text, 'MENTOR'::text, 'INTERN'::text]))),
    CONSTRAINT "User_status_check" CHECK ((status = ANY (ARRAY['PENDING'::text, 'APPROVED'::text, 'REJECTED'::text])))
);


ALTER TABLE public."User" OWNER TO postgres;

--
-- Name: WeeklySnapshot; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."WeeklySnapshot" (
    id text NOT NULL,
    "userId" text NOT NULL,
    "weekId" text NOT NULL,
    participation text NOT NULL,
    "outputEvidence" text NOT NULL,
    momentum text NOT NULL,
    quality text NOT NULL,
    signal text NOT NULL,
    "escalationLevel" integer DEFAULT 0 NOT NULL,
    "consecutiveNegativeWeeks" integer DEFAULT 0 NOT NULL,
    "overrideSignal" text,
    "overrideReason" text,
    "overriddenBy" text,
    "overriddenAt" text,
    "verifiedBy" text,
    "verifiedAt" text,
    "internFacingLabel" text,
    "computedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."WeeklySnapshot" OWNER TO postgres;

--
-- Name: WorkLog; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."WorkLog" (
    id text NOT NULL,
    "userId" text NOT NULL,
    "logDate" text NOT NULL,
    "taskItems" text DEFAULT '[]'::text NOT NULL,
    "daySummary" text,
    "hasEvidence" boolean DEFAULT false NOT NULL,
    "hasBlocker" boolean DEFAULT false NOT NULL,
    "isLate" boolean DEFAULT false NOT NULL,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."WorkLog" OWNER TO postgres;

--
-- Name: WorkflowAction; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."WorkflowAction" (
    id text NOT NULL,
    "workflowId" text NOT NULL,
    description text NOT NULL,
    "assignedTo" text NOT NULL,
    "dueDate" text NOT NULL,
    status text DEFAULT 'OPEN'::text NOT NULL,
    notes text,
    "completedAt" text,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."WorkflowAction" OWNER TO postgres;

--
-- Name: WorkflowApprovalStep; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."WorkflowApprovalStep" (
    id text NOT NULL,
    "workflowId" text NOT NULL,
    "stepNumber" integer NOT NULL,
    "approverId" text NOT NULL,
    role text NOT NULL,
    status text DEFAULT 'PENDING'::text NOT NULL,
    decision text,
    reason text,
    "decidedAt" text
);


ALTER TABLE public."WorkflowApprovalStep" OWNER TO postgres;

--
-- Name: XPLog; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."XPLog" (
    id text NOT NULL,
    "userId" text NOT NULL,
    amount integer NOT NULL,
    reason text NOT NULL,
    metadata text,
    "createdAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public."XPLog" OWNER TO postgres;

--
-- Name: AILog AILog_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."AILog"
    ADD CONSTRAINT "AILog_pkey" PRIMARY KEY (id);


--
-- Name: Achievement Achievement_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Achievement"
    ADD CONSTRAINT "Achievement_pkey" PRIMARY KEY (id);


--
-- Name: Achievement Achievement_userId_type_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Achievement"
    ADD CONSTRAINT "Achievement_userId_type_key" UNIQUE ("userId", type);


--
-- Name: Attendance Attendance_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Attendance"
    ADD CONSTRAINT "Attendance_pkey" PRIMARY KEY (id);


--
-- Name: Certificate Certificate_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Certificate"
    ADD CONSTRAINT "Certificate_pkey" PRIMARY KEY (id);


--
-- Name: College College_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."College"
    ADD CONSTRAINT "College_name_key" UNIQUE (name);


--
-- Name: College College_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."College"
    ADD CONSTRAINT "College_pkey" PRIMARY KEY (id);


--
-- Name: DailyReport DailyReport_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."DailyReport"
    ADD CONSTRAINT "DailyReport_pkey" PRIMARY KEY (id);


--
-- Name: DecisionWorkflow DecisionWorkflow_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."DecisionWorkflow"
    ADD CONSTRAINT "DecisionWorkflow_pkey" PRIMARY KEY (id);


--
-- Name: Domain Domain_domainName_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Domain"
    ADD CONSTRAINT "Domain_domainName_key" UNIQUE ("domainName");


--
-- Name: Domain Domain_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Domain"
    ADD CONSTRAINT "Domain_pkey" PRIMARY KEY (id);


--
-- Name: Evaluation Evaluation_internId_weekNumber_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Evaluation"
    ADD CONSTRAINT "Evaluation_internId_weekNumber_key" UNIQUE ("internId", "weekNumber");


--
-- Name: Evaluation Evaluation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Evaluation"
    ADD CONSTRAINT "Evaluation_pkey" PRIMARY KEY (id);


--
-- Name: FlagResponse FlagResponse_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."FlagResponse"
    ADD CONSTRAINT "FlagResponse_pkey" PRIMARY KEY (id);


--
-- Name: InternProfile InternProfile_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InternProfile"
    ADD CONSTRAINT "InternProfile_pkey" PRIMARY KEY (id);


--
-- Name: InternProfile InternProfile_userId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InternProfile"
    ADD CONSTRAINT "InternProfile_userId_key" UNIQUE ("userId");


--
-- Name: Meeting Meeting_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Meeting"
    ADD CONSTRAINT "Meeting_pkey" PRIMARY KEY (id);


--
-- Name: MentorAssessment MentorAssessment_mentorId_menteeId_weekId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."MentorAssessment"
    ADD CONSTRAINT "MentorAssessment_mentorId_menteeId_weekId_key" UNIQUE ("mentorId", "menteeId", "weekId");


--
-- Name: MentorAssessment MentorAssessment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."MentorAssessment"
    ADD CONSTRAINT "MentorAssessment_pkey" PRIMARY KEY (id);


--
-- Name: Mentor Mentor_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Mentor"
    ADD CONSTRAINT "Mentor_pkey" PRIMARY KEY (id);


--
-- Name: Mentor Mentor_userId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Mentor"
    ADD CONSTRAINT "Mentor_userId_key" UNIQUE ("userId");


--
-- Name: Message Message_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_pkey" PRIMARY KEY (id);


--
-- Name: Notification Notification_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Notification"
    ADD CONSTRAINT "Notification_pkey" PRIMARY KEY (id);


--
-- Name: ReportDraft ReportDraft_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ReportDraft"
    ADD CONSTRAINT "ReportDraft_pkey" PRIMARY KEY (id);


--
-- Name: Streak Streak_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Streak"
    ADD CONSTRAINT "Streak_pkey" PRIMARY KEY (id);


--
-- Name: Streak Streak_userId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Streak"
    ADD CONSTRAINT "Streak_userId_key" UNIQUE ("userId");


--
-- Name: Task Task_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Task"
    ADD CONSTRAINT "Task_pkey" PRIMARY KEY (id);


--
-- Name: User User_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_email_key" UNIQUE (email);


--
-- Name: User User_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_pkey" PRIMARY KEY (id);


--
-- Name: User User_referralCode_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_referralCode_key" UNIQUE ("referralCode");


--
-- Name: WeeklySnapshot WeeklySnapshot_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WeeklySnapshot"
    ADD CONSTRAINT "WeeklySnapshot_pkey" PRIMARY KEY (id);


--
-- Name: WeeklySnapshot WeeklySnapshot_userId_weekId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WeeklySnapshot"
    ADD CONSTRAINT "WeeklySnapshot_userId_weekId_key" UNIQUE ("userId", "weekId");


--
-- Name: WorkLog WorkLog_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WorkLog"
    ADD CONSTRAINT "WorkLog_pkey" PRIMARY KEY (id);


--
-- Name: WorkLog WorkLog_userId_logDate_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WorkLog"
    ADD CONSTRAINT "WorkLog_userId_logDate_key" UNIQUE ("userId", "logDate");


--
-- Name: WorkflowAction WorkflowAction_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WorkflowAction"
    ADD CONSTRAINT "WorkflowAction_pkey" PRIMARY KEY (id);


--
-- Name: WorkflowApprovalStep WorkflowApprovalStep_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WorkflowApprovalStep"
    ADD CONSTRAINT "WorkflowApprovalStep_pkey" PRIMARY KEY (id);


--
-- Name: WorkflowApprovalStep WorkflowApprovalStep_workflowId_stepNumber_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WorkflowApprovalStep"
    ADD CONSTRAINT "WorkflowApprovalStep_workflowId_stepNumber_key" UNIQUE ("workflowId", "stepNumber");


--
-- Name: XPLog XPLog_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."XPLog"
    ADD CONSTRAINT "XPLog_pkey" PRIMARY KEY (id);


--
-- Name: AILog_adminId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "AILog_adminId_idx" ON public."AILog" USING btree ("adminId");


--
-- Name: AILog_createdAt_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "AILog_createdAt_idx" ON public."AILog" USING btree ("createdAt");


--
-- Name: Achievement_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Achievement_userId_idx" ON public."Achievement" USING btree ("userId");


--
-- Name: Attendance_date_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Attendance_date_idx" ON public."Attendance" USING btree (date);


--
-- Name: Attendance_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Attendance_userId_idx" ON public."Attendance" USING btree ("userId");


--
-- Name: Certificate_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Certificate_userId_idx" ON public."Certificate" USING btree ("userId");


--
-- Name: DailyReport_date_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "DailyReport_date_idx" ON public."DailyReport" USING btree (date);


--
-- Name: DailyReport_internId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "DailyReport_internId_idx" ON public."DailyReport" USING btree ("internId");


--
-- Name: DecisionWorkflow_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "DecisionWorkflow_status_idx" ON public."DecisionWorkflow" USING btree (status);


--
-- Name: DecisionWorkflow_subjectId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "DecisionWorkflow_subjectId_idx" ON public."DecisionWorkflow" USING btree ("subjectId");


--
-- Name: DecisionWorkflow_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "DecisionWorkflow_type_idx" ON public."DecisionWorkflow" USING btree (type);


--
-- Name: DecisionWorkflow_weekId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "DecisionWorkflow_weekId_idx" ON public."DecisionWorkflow" USING btree ("weekId");


--
-- Name: Evaluation_evaluatorId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Evaluation_evaluatorId_idx" ON public."Evaluation" USING btree ("evaluatorId");


--
-- Name: Evaluation_internId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Evaluation_internId_idx" ON public."Evaluation" USING btree ("internId");


--
-- Name: FlagResponse_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "FlagResponse_status_idx" ON public."FlagResponse" USING btree (status);


--
-- Name: InternProfile_preferredDomain_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "InternProfile_preferredDomain_idx" ON public."InternProfile" USING btree ("preferredDomain");


--
-- Name: Meeting_internId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Meeting_internId_idx" ON public."Meeting" USING btree ("internId");


--
-- Name: Meeting_mentorId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Meeting_mentorId_idx" ON public."Meeting" USING btree ("mentorId");


--
-- Name: MentorAssessment_menteeId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "MentorAssessment_menteeId_idx" ON public."MentorAssessment" USING btree ("menteeId");


--
-- Name: MentorAssessment_mentorId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "MentorAssessment_mentorId_idx" ON public."MentorAssessment" USING btree ("mentorId");


--
-- Name: MentorAssessment_weekId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "MentorAssessment_weekId_idx" ON public."MentorAssessment" USING btree ("weekId");


--
-- Name: Message_createdAt_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Message_createdAt_idx" ON public."Message" USING btree ("createdAt");


--
-- Name: Message_receiverId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Message_receiverId_idx" ON public."Message" USING btree ("receiverId");


--
-- Name: Message_senderId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Message_senderId_idx" ON public."Message" USING btree ("senderId");


--
-- Name: Notification_read_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Notification_read_idx" ON public."Notification" USING btree (read);


--
-- Name: Notification_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Notification_userId_idx" ON public."Notification" USING btree ("userId");


--
-- Name: ReportDraft_recipientUserId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "ReportDraft_recipientUserId_idx" ON public."ReportDraft" USING btree ("recipientUserId");


--
-- Name: ReportDraft_reportType_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "ReportDraft_reportType_idx" ON public."ReportDraft" USING btree ("reportType");


--
-- Name: ReportDraft_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "ReportDraft_status_idx" ON public."ReportDraft" USING btree (status);


--
-- Name: ReportDraft_weekId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "ReportDraft_weekId_idx" ON public."ReportDraft" USING btree ("weekId");


--
-- Name: Streak_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Streak_userId_idx" ON public."Streak" USING btree ("userId");


--
-- Name: Task_assignedById_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Task_assignedById_idx" ON public."Task" USING btree ("assignedById");


--
-- Name: Task_assignedToId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Task_assignedToId_idx" ON public."Task" USING btree ("assignedToId");


--
-- Name: Task_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "Task_status_idx" ON public."Task" USING btree (status);


--
-- Name: User_createdAt_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "User_createdAt_idx" ON public."User" USING btree ("createdAt");


--
-- Name: User_role_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "User_role_idx" ON public."User" USING btree (role);


--
-- Name: User_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "User_status_idx" ON public."User" USING btree (status);


--
-- Name: User_xp_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "User_xp_idx" ON public."User" USING btree (xp);


--
-- Name: WeeklySnapshot_escalationLevel_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WeeklySnapshot_escalationLevel_idx" ON public."WeeklySnapshot" USING btree ("escalationLevel");


--
-- Name: WeeklySnapshot_signal_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WeeklySnapshot_signal_idx" ON public."WeeklySnapshot" USING btree (signal);


--
-- Name: WeeklySnapshot_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WeeklySnapshot_userId_idx" ON public."WeeklySnapshot" USING btree ("userId");


--
-- Name: WeeklySnapshot_weekId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WeeklySnapshot_weekId_idx" ON public."WeeklySnapshot" USING btree ("weekId");


--
-- Name: WorkLog_logDate_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WorkLog_logDate_idx" ON public."WorkLog" USING btree ("logDate");


--
-- Name: WorkLog_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WorkLog_userId_idx" ON public."WorkLog" USING btree ("userId");


--
-- Name: WorkLog_userId_logDate_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WorkLog_userId_logDate_idx" ON public."WorkLog" USING btree ("userId", "logDate");


--
-- Name: WorkflowAction_assignedTo_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WorkflowAction_assignedTo_idx" ON public."WorkflowAction" USING btree ("assignedTo");


--
-- Name: WorkflowAction_workflowId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WorkflowAction_workflowId_idx" ON public."WorkflowAction" USING btree ("workflowId");


--
-- Name: WorkflowApprovalStep_workflowId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "WorkflowApprovalStep_workflowId_idx" ON public."WorkflowApprovalStep" USING btree ("workflowId");


--
-- Name: XPLog_createdAt_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "XPLog_createdAt_idx" ON public."XPLog" USING btree ("createdAt");


--
-- Name: XPLog_userId_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "XPLog_userId_idx" ON public."XPLog" USING btree ("userId");


--
-- Name: AILog AILog_adminId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."AILog"
    ADD CONSTRAINT "AILog_adminId_fkey" FOREIGN KEY ("adminId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Achievement Achievement_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Achievement"
    ADD CONSTRAINT "Achievement_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Attendance Attendance_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Attendance"
    ADD CONSTRAINT "Attendance_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Certificate Certificate_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Certificate"
    ADD CONSTRAINT "Certificate_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: DailyReport DailyReport_internId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."DailyReport"
    ADD CONSTRAINT "DailyReport_internId_fkey" FOREIGN KEY ("internId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Evaluation Evaluation_evaluatorId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Evaluation"
    ADD CONSTRAINT "Evaluation_evaluatorId_fkey" FOREIGN KEY ("evaluatorId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Evaluation Evaluation_internId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Evaluation"
    ADD CONSTRAINT "Evaluation_internId_fkey" FOREIGN KEY ("internId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: InternProfile InternProfile_collegeId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InternProfile"
    ADD CONSTRAINT "InternProfile_collegeId_fkey" FOREIGN KEY ("collegeId") REFERENCES public."College"(id) ON DELETE SET NULL;


--
-- Name: InternProfile InternProfile_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InternProfile"
    ADD CONSTRAINT "InternProfile_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Meeting Meeting_internId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Meeting"
    ADD CONSTRAINT "Meeting_internId_fkey" FOREIGN KEY ("internId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Meeting Meeting_mentorId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Meeting"
    ADD CONSTRAINT "Meeting_mentorId_fkey" FOREIGN KEY ("mentorId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: MentorAssessment MentorAssessment_menteeId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."MentorAssessment"
    ADD CONSTRAINT "MentorAssessment_menteeId_fkey" FOREIGN KEY ("menteeId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: MentorAssessment MentorAssessment_mentorId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."MentorAssessment"
    ADD CONSTRAINT "MentorAssessment_mentorId_fkey" FOREIGN KEY ("mentorId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Mentor Mentor_domainId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Mentor"
    ADD CONSTRAINT "Mentor_domainId_fkey" FOREIGN KEY ("domainId") REFERENCES public."Domain"(id) ON DELETE SET NULL;


--
-- Name: Mentor Mentor_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Mentor"
    ADD CONSTRAINT "Mentor_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Message Message_receiverId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_receiverId_fkey" FOREIGN KEY ("receiverId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Message Message_senderId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_senderId_fkey" FOREIGN KEY ("senderId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Notification Notification_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Notification"
    ADD CONSTRAINT "Notification_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Streak Streak_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Streak"
    ADD CONSTRAINT "Streak_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Task Task_assignedById_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Task"
    ADD CONSTRAINT "Task_assignedById_fkey" FOREIGN KEY ("assignedById") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Task Task_assignedToId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Task"
    ADD CONSTRAINT "Task_assignedToId_fkey" FOREIGN KEY ("assignedToId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: User User_referredById_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_referredById_fkey" FOREIGN KEY ("referredById") REFERENCES public."User"(id) ON DELETE SET NULL;


--
-- Name: WeeklySnapshot WeeklySnapshot_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WeeklySnapshot"
    ADD CONSTRAINT "WeeklySnapshot_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: WorkLog WorkLog_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WorkLog"
    ADD CONSTRAINT "WorkLog_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: WorkflowAction WorkflowAction_workflowId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WorkflowAction"
    ADD CONSTRAINT "WorkflowAction_workflowId_fkey" FOREIGN KEY ("workflowId") REFERENCES public."DecisionWorkflow"(id) ON DELETE CASCADE;


--
-- Name: WorkflowApprovalStep WorkflowApprovalStep_workflowId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."WorkflowApprovalStep"
    ADD CONSTRAINT "WorkflowApprovalStep_workflowId_fkey" FOREIGN KEY ("workflowId") REFERENCES public."DecisionWorkflow"(id) ON DELETE CASCADE;


--
-- Name: XPLog XPLog_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."XPLog"
    ADD CONSTRAINT "XPLog_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict 2TODWod3bUECuZ9tI9YQTchgQvnb2ADZXchkcr6QbzBcoudWPyUgedutDmn5ysa

