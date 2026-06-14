import json
from langchain_ollama import ChatOllama
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from markdown_it import MarkdownIt

class ReportGenerator:
    def __init__(self, model_name="phi4-mini:latest"):
        self.llm = ChatOllama(
            model=model_name,
            temperature=0.2,
            base_url="http://localhost:11434"
        )
        self.output_parser = StrOutputParser()
        self.md_parser = MarkdownIt()

    def generate_detailed_markdown(self, intern_data, week_id):
        user_info = intern_data.get('User', {})
        intern_name = user_info.get('name', 'Unknown')
        intern_email = user_info.get('email', 'unknown')
        xp = user_info.get('xp', 0)
        level = user_info.get('level', 1)

        tasks = intern_data.get('TasksAssignedTo', [])
        daily_reports = intern_data.get('DailyReports', [])
        evaluations = intern_data.get('EvaluationsAsIntern', [])
        mentor_assessments = intern_data.get('MentorAssessmentsAsMentee', [])
        weekly_snapshots = intern_data.get('WeeklySnapshots', [])

        prompt = ChatPromptTemplate.from_messages([
            ("system", (
                "You are an experienced Internship Program Manager responsible for generating professional weekly performance reports. "
                "Your task is to analyze the provided intern data and generate a factual, evidence-based, manager-ready performance report in standard Markdown (GFM). "
                "Do NOT include conversational intros/outros. Return only the final Markdown code, beginning with '# Intern Performance Report'."
            )),
            ("user", (
                "Analyze the provided intern data and generate a factual, evidence-based, manager-ready performance report for intern {intern_name} ({intern_email}) for Week {week_id}.\n\n"
                "Here is the database context gathered from the gRPC backend for this intern:\n"
                "- Current XP: {xp} (Level {level})\n"
                "- Weekly Snapshots: {weekly_snapshots}\n"
                "- Tasks: {tasks}\n"
                "- Daily Activity Logs & Blockers: {daily_reports}\n"
                "- Rubric Evaluations: {evaluations}\n"
                "- Mentor Assessments: {mentor_assessments}\n\n"
                "IMPORTANT RULES:\n"
                "* Never invent information that does not exist in the provided data.\n"
                "* Every observation must be supported by evidence from tasks, work logs, evaluations, mentor assessments, attendance, meetings, achievements, XP logs, weekly snapshots, or daily reports.\n"
                "* Avoid generic AI phrases such as: 'showed great potential', 'demonstrated excellence', 'continued to grow', 'performed admirably', unless supported by actual evidence.\n"
                "* Focus on measurable performance.\n"
                "* Maintain a professional and constructive tone.\n"
                "* Highlight both strengths and weaknesses.\n"
                "* If data is missing, explicitly state that insufficient evidence exists.\n"
                "* Do not repeat the same information in multiple sections.\n"
                "* Be specific and actionable.\n\n"
                "Generate the report using the exact structure below:\n\n"
                "# Intern Performance Report\n\n"
                "## 1. Executive Summary\n"
                "Provide a concise overview of the intern's week.\n"
                "Include:\n"
                "* Overall performance assessment\n"
                "* Major accomplishments\n"
                "* Major concerns\n"
                "* Overall performance category (Must be one of: Exceptional, Strong, Satisfactory, Developing, At Risk)\n"
                "* Confidence level (Must be one of: High Confidence, Medium Confidence, Low Confidence) based on the amount of available evidence.\n\n"
                "## 2. Performance Scorecard\n"
                "Display a table:\n"
                "| Category | Score | Assessment |\n"
                "| --- | --- | --- |\n"
                "| Technical Skills | X/5 | ... |\n"
                "| Initiative | X/5 | ... |\n"
                "| Leadership | X/5 | ... |\n"
                "| Consistency | X/5 | ... |\n"
                "| Overall | X/20 | ... |\n\n"
                "After the table explain:\n"
                "* Why the score was earned\n"
                "* Which evidence contributed most to the score\n\n"
                "## 3. Key Strengths\n"
                "Identify 3-5 strengths. For each strength include:\n"
                "* Observation\n"
                "* Supporting evidence\n"
                "* Impact\n\n"
                "## 4. Areas for Improvement\n"
                "Identify the most important improvement opportunities. For each area include:\n"
                "* Issue\n"
                "* Evidence\n"
                "* Potential impact\n"
                "* Recommended corrective action\n"
                "Do not create improvement areas unless supported by evidence.\n\n"
                "## 5. Deliverables & Evidence Review\n"
                "Create a detailed table:\n"
                "| Deliverable | Status | Quality Assessment | Evidence |\n"
                "| --- | --- | --- | --- |\n"
                "Include tasks, daily reports, work logs, submissions, and achievements. Provide commentary on overall output quality.\n\n"
                "## 6. Mentor & Evaluation Insights\n"
                "Summarize mentor feedback, evaluation scores, assessment comments, and attestations. Highlight patterns. If mentor feedback conflicts with evaluation scores, explain the discrepancy.\n\n"
                "## 7. Blockers & Risk Assessment\n"
                "### Blockers\n"
                "List reported blockers, severity, resolution status, and time impact.\n"
                "### Risk Level\n"
                "Assign Low Risk, Medium Risk, or High Risk and justify using evidence (rejected tasks, missing submissions, negative evaluations, attendance concerns, repeated blockers).\n\n"
                "## 8. Momentum & Trend Analysis\n"
                "Analyze participation, work quality, output consistency, weekly snapshot signals, XP progression, and streak information. Determine momentum (Improving, Stable, or Declining) and explain why. If previous-week data is unavailable, explicitly state that trend analysis is limited.\n\n"
                "## 9. Next Week Action Plan\n"
                "Provide:\n"
                "### Priority Actions\n"
                "### Learning Goals\n"
                "### Expected Outcomes\n"
                "Recommendations must be specific and measurable.\n\n"
                "## 10. Final Verdict\n"
                "Summarize current standing, readiness for increased responsibility, main strength, and main concern.\n"
                "End with:\n"
                "Overall Status: [Exceptional / Strong / Satisfactory / Developing / At Risk]\n"
                "Confidence: [High / Medium / Low]\n"
                "Evidence Count:\n"
                "* Tasks Reviewed: X\n"
                "* Work Logs Reviewed: X\n"
                "* Daily Reports Reviewed: X\n"
                "* Evaluations Reviewed: X\n"
                "* Mentor Assessments Reviewed: X\n"
            ))
        ])

        chain = prompt | self.llm | self.output_parser
        response = chain.invoke({
            "intern_name": intern_name,
            "intern_email": intern_email,
            "week_id": week_id,
            "xp": xp,
            "level": level,
            "weekly_snapshots": json.dumps(weekly_snapshots, indent=2),
            "tasks": json.dumps(tasks, indent=2),
            "daily_reports": json.dumps(daily_reports, indent=2),
            "evaluations": json.dumps(evaluations, indent=2),
            "mentor_assessments": json.dumps(mentor_assessments, indent=2)
        })

        clean_md = response.strip()
        if clean_md.startswith("```markdown"):
            clean_md = clean_md[11:]
        if clean_md.endswith("```"):
            clean_md = clean_md[:-3]
        return clean_md.strip()

    def render_markdown_to_html(self, markdown_content):
        # Convert Markdown to HTML
        html_body = self.md_parser.render(markdown_content)
        
        # Wrap in a premium HTML/CSS dashboard layout
        styled_html = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
    body {{
        font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
        line-height: 1.6;
        color: #e2e8f0;
        background-color: #0f172a;
        margin: 0;
        padding: 40px;
    }}
    .container {{
        max-width: 900px;
        margin: 0 auto;
        background: #1e293b;
        padding: 40px;
        border-radius: 12px;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.15), 0 2px 4px -1px rgba(0, 0, 0, 0.1);
        border: 1px solid #334155;
    }}
    h1, h2, h3 {{
        color: #38bdf8;
        border-bottom: 1px solid #334155;
        padding-bottom: 8px;
        margin-top: 30px;
    }}
    h1 {{
        color: #0ea5e9;
        font-size: 2.2em;
        margin-top: 0;
    }}
    table {{
        width: 100%;
        border-collapse: collapse;
        margin: 20px 0;
    }}
    th, td {{
        padding: 12px;
        border-bottom: 1px solid #334155;
        text-align: left;
        font-size: 0.95em;
    }}
    th {{
        background-color: #334155;
        color: #f8fafc;
        font-weight: 600;
    }}
    tr:hover {{
        background-color: #1e293b;
    }}
    blockquote {{
        border-left: 4px solid #38bdf8;
        background: #0f172a;
        margin: 20px 0;
        padding: 15px 20px;
        border-radius: 0 8px 8px 0;
    }}
    blockquote p {{
        margin: 0;
        font-style: italic;
    }}
    code {{
        background-color: #0f172a;
        color: #f43f5e;
        padding: 2px 6px;
        border-radius: 4px;
        font-family: Consolas, Monaco, monospace;
        font-size: 0.9em;
    }}
    pre {{
        background-color: #0f172a;
        padding: 15px;
        border-radius: 8px;
        overflow-x: auto;
        border: 1px solid #334155;
    }}
    ul, ol {{
        padding-left: 20px;
    }}
    li {{
        margin-bottom: 8px;
    }}
</style>
</head>
<body>
    <div class="container">
        {html_body}
    </div>
</body>
</html>
"""
        return styled_html
