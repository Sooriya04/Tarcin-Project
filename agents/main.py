import sys
import os
import argparse
import json

# Adjust sys.path to find sibling imports
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from config import get_db_connection
from logger import logger
from client import GrpcClient
from report_generator import ReportGenerator

# Ensure Summary table exists inside the tarcin database
def init_summary_database(conn):
    logger.info("Checking and initializing 'Summary' table in the database...")
    with conn.cursor() as cur:
        cur.execute("""
            CREATE TABLE IF NOT EXISTS "Summary" (
                "id" TEXT NOT NULL PRIMARY KEY,
                "internId" TEXT NOT NULL,
                "weekId" TEXT NOT NULL,
                "markdownReport" TEXT NOT NULL,
                "htmlReport" TEXT NOT NULL,
                "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY ("internId") REFERENCES "User"("id") ON DELETE CASCADE
            );
            CREATE INDEX IF NOT EXISTS "Summary_internId_idx" ON "Summary"("internId");
            CREATE INDEX IF NOT EXISTS "Summary_weekId_idx" ON "Summary"("weekId");
        """)
    conn.commit()
    logger.info("'Summary' table check completed.")

# Save to the traditional ReportDraft table (for backward compatibility)
def save_to_report_draft(conn, intern_id, intern_name, intern_email, week_id, html_report):
    report_id = f"rd-{intern_id}-{week_id}"
    with conn.cursor() as cur:
        cur.execute(
            """INSERT INTO "ReportDraft" (
                "id", "weekId", "reportType", "recipientUserId", "recipientName", "recipientEmail", "renderedHtml", "status", "isDryRun", "createdAt"
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, NOW())
            ON CONFLICT (id) DO UPDATE SET
                "renderedHtml" = EXCLUDED."renderedHtml",
                "status" = 'DRAFT',
                "createdAt" = NOW()""",
            (report_id, week_id, 'WEEKLY_PERFORMANCE', intern_id, intern_name, intern_email, html_report, 'DRAFT', True)
        )
    conn.commit()
    logger.info(f"Database Updated: Legacy ReportDraft written for {intern_name} (ID: {report_id})")

# Save to the new Summary table
def save_to_summary_table(conn, intern_id, intern_name, week_id, markdown_report, html_report):
    summary_id = f"sum-{intern_id}-{week_id}"
    with conn.cursor() as cur:
        cur.execute(
            """INSERT INTO "Summary" (
                "id", "internId", "weekId", "markdownReport", "htmlReport", "createdAt"
            ) VALUES (%s, %s, %s, %s, %s, NOW())
            ON CONFLICT (id) DO UPDATE SET
                "markdownReport" = EXCLUDED."markdownReport",
                "htmlReport" = EXCLUDED."htmlReport",
                "createdAt" = NOW()""",
            (summary_id, intern_id, week_id, markdown_report, html_report)
        )
    conn.commit()
    logger.info(f"Database Updated: New Summary table entry written for {intern_name} (ID: {summary_id})")

def main():
    parser = argparse.ArgumentParser(description="AI Intern Performance Report Agent (Production-Ready)")
    parser.add_argument("--week", default="2026-W22", help="The target week (e.g., 2026-W22)")
    parser.add_argument("--limit", type=int, default=None, help="Limit the number of interns to evaluate")
    args = parser.parse_args()

    week_id = args.week
    logger.info(f"==================================================================")
    logger.info(f"Starting AI Performance Agentic Pipeline for week: {week_id}")
    logger.info(f"==================================================================")

    try:
        conn = get_db_connection()
    except Exception as e:
        logger.error(f"Failed to connect to PostgreSQL: {e}")
        sys.exit(1)

    # Initialize the Summary database table
    try:
        init_summary_database(conn)
    except Exception as e:
        logger.error(f"Failed to initialize 'Summary' database table: {e}")
        conn.close()
        sys.exit(1)

    # Fetch all active interns
    try:
        with conn.cursor() as cur:
            cur.execute('SELECT id, name, email FROM "User" WHERE role = \'INTERN\' AND status = \'APPROVED\'')
            interns = [dict(zip(['id', 'name', 'email'], row)) for row in cur.fetchall()]
    except Exception as e:
        logger.error(f"Failed to fetch interns list: {e}")
        conn.close()
        sys.exit(1)

    if not interns:
        logger.info("No active interns found in the database. Exiting.")
        conn.close()
        return

    if args.limit is not None:
        logger.info(f"Applying limit: only evaluating the first {args.limit} interns.")
        interns = interns[:args.limit]

    logger.info(f"Discovered {len(interns)} interns to evaluate.")

    grpc_client = GrpcClient()
    report_gen = ReportGenerator()

    # Create reports output directory
    reports_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)), "reports")
    os.makedirs(reports_dir, exist_ok=True)

    for intern in interns:
        intern_id = intern['id']
        intern_name = intern['name']
        intern_email = intern['email']
        
        logger.info(f"Evaluating intern: {intern_name} (ID: {intern_id})")
        try:
            # Fetch data from gRPC client
            json_result_str = grpc_client.fetch_user_information(intern_id)
            intern_data = json.loads(json_result_str)
            
            # Generate detailed Markdown report
            logger.info(f"Invoking phi4-mini local LLM via LangChain...")
            markdown_report = report_gen.generate_detailed_markdown(intern_data, week_id)
            
            # Render Markdown to styled HTML
            html_report = report_gen.render_markdown_to_html(markdown_report)
            
            # Save to standard ReportDraft table
            save_to_report_draft(conn, intern_id, intern_name, intern_email, week_id, html_report)
            
            # Save to new Summary table
            save_to_summary_table(conn, intern_id, intern_name, week_id, markdown_report, html_report)

            # Save Markdown locally
            md_file_path = os.path.join(reports_dir, f"{intern_id}_{week_id}.md")
            with open(md_file_path, "w", encoding="utf-8") as f:
                f.write(markdown_report)
                
            # Save HTML locally
            html_file_path = os.path.join(reports_dir, f"{intern_id}_{week_id}.html")
            with open(html_file_path, "w", encoding="utf-8") as f:
                f.write(html_report)
                
            logger.info(f"Saved local files: {md_file_path} and {html_file_path}")
            
        except Exception as e:
            logger.error(f"Failed pipeline processing for {intern_name}: {e}")

    conn.close()
    logger.info(f"AI Performance Agentic Pipeline completed successfully!")

if __name__ == "__main__":
    main()
