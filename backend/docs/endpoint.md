# Tarcin Performance Backend Endpoints

The backend is a high-performance gRPC microservice running on `localhost:50051` with Server Reflection enabled.

## 1. Performance Metrics
**Method:** `tarcin.performance.PerformanceService/GetPerformanceMetrics`
**Description:** Returns top and lowest performers, average evaluation scores, and performance sliced by domain and college.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetPerformanceMetrics`

## 2. Intern Health
**Method:** `tarcin.performance.PerformanceService/GetInternHealth`
**Description:** High-level metrics on active, completed, and inactive interns, plus top colleges and growing domains.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetInternHealth`

## 3. Task Analytics
**Method:** `tarcin.performance.PerformanceService/GetTaskAnalytics`
**Description:** Deep dive into task completion rates, late tasks, rejection reasons, and mentor turnaround speed.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetTaskAnalytics`

## 4. Mentor Analytics
**Method:** `tarcin.performance.PerformanceService/GetMentorAnalytics`
**Description:** Analyzes mentor workload, effectiveness, task approval rates, and dropout risk among their mentees.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetMentorAnalytics`

## 5. Engagement Analytics
**Method:** `tarcin.performance.PerformanceService/GetEngagementAnalytics`
**Description:** Tracks intern attendance streaks, declining participation, daily reports, and drop-out risks.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetEngagementAnalytics`

## 6. Growth Analytics
**Method:** `tarcin.performance.PerformanceService/GetGrowthAnalytics`
**Description:** Gamification metrics including XP gained, frequent achievements, and leveling up progress.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetGrowthAnalytics`

## 7. College Analytics
**Method:** `tarcin.performance.PerformanceService/GetCollegeAnalytics`
**Description:** Institutional performance tracking, completion rates by college, and cohort evaluation scores.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetCollegeAnalytics`

## 8. Workflow Analytics
**Method:** `tarcin.performance.PerformanceService/GetWorkflowAnalytics`
**Description:** Tracks escalated interns, negative trends, root causes for blockers, and past due operational decisions.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetWorkflowAnalytics`

## 9. Conversion Analytics
**Method:** `tarcin.performance.PerformanceService/GetConversionAnalytics`
**Description:** Highlights interns ready for hire, lagging conversions approaching deadlines, and college conversion rates.
**Test:**
`grpcurl -plaintext localhost:50051 tarcin.performance.PerformanceService/GetConversionAnalytics`
