# Tarcin AI Performance Agentic Pipeline & Metrics Monitor

This package contains the production-ready AI Reporting Agent and System Metrics Monitor.

---

## 1. AI Reporting Agent

The agent retrieves intern performance metrics directly from the Go gRPC backend (`localhost:50051`) and uses LangChain with the local Ollama model `phi4-mini:latest` to generate detailed, weekly summary dashboards.

### File Structure
- `agents/main.py`: Main entry point orchestrator.
- `agents/config.py`: Connection configuration.
- `agents/client.py`: gRPC client for fetching user details.
- `agents/report_generator.py`: Runs LangChain chain, prompts, and Markdown-to-HTML rendering using `markdown-it-py`.
- `agents/logger.py`: Handles file/console logging separate from the backend.

### Output Tables
Generated reports are stored in:
1. **`Summary` Table** (New): Stores the detailed raw `markdownReport` and compiled `htmlReport` for each intern.
2. **`ReportDraft` Table** (Legacy): Stores the compiled `htmlReport` for compatibility.
3. **Local Files**: HTML/Markdown outputs saved to `agents/reports/`.

### How to Run
```bash
# Activate virtual environment
source .venv/bin/activate

# Execute pipeline (defaults to 2026-W22)
python agents/main.py --week 2026-W22
```

---

## 2. System Metrics Monitor (`agents/metrics.py`)

A standalone daemon script that continuously polls the host's system metrics and records them.

### Data Storage & Isolation
- **No public schema pollution**: All metrics are stored in a dedicated PostgreSQL schema (`metrics`) and table (`system_metrics`), isolated from public application data.
- **Logs**: Periodically recorded in `agents/metrics.log`.

### Schema (PostgreSQL `metrics.system_metrics` table)
- `id` (SERIAL PRIMARY KEY)
- `timestamp` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
- `cpu_usage` (REAL, %)
- `ram_usage` (REAL, %)
- `gpu_usage` (REAL, % - queries nvidia-smi if Nvidia GPU is present)
- `gpu_mem_usage` (REAL, % - queries nvidia-smi if Nvidia GPU is present)

### How to Run
```bash
python agents/metrics.py
```
- The terminal displays live metric updates dynamically.
- Press `Ctrl + C` to stop monitoring gracefully at any time.

---

## 3. Logs Location
- **Go Backend Server Logs**: Written to `backend/server.log`.
- **Python Agent Pipeline Logs**: Written to `agents/agent.log`.
- **System Metrics Logs**: Written to `agents/metrics.log`.
