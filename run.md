# Tarcin Project Execution Guide

This document explains how to launch and execute the backend, agent, and monitoring components.

---

## 1. Setup Phase

First, ensure that all Go and Python dependencies are fully installed:
```bash
./setup.sh
```

Also, ensure that your PostgreSQL database credentials are correct in `backend/.env` and that Ollama is running with the correct model pulled:
```bash
ollama pull phi4-mini:latest
```

---

## 2. Launch the gRPC Backend Server

Run the Go gRPC backend server. It listens on port `:50051` and writes logs to `backend/server.log`:
```bash
cd backend
go build -o server ./cmd/server/main.go
./server
```

---

## 3. Run the AI Performance Evaluator Agent

Activate the virtual environment and run the pipeline.

To run for **all** interns:
```bash
source .venv/bin/activate
python agents/main.py --week 2026-W22
```

To run for a **subset** of interns (e.g., limit to 3):
```bash
source .venv/bin/activate
python agents/main.py --week 2026-W22 --limit 3
```

- Local Markdown and HTML reports are written to `agents/reports/`.
- PostgreSQL tables `ReportDraft` (legacy) and `Summary` (new detailed GFM table) are updated automatically.
- Logs are appended to `agents/agent.log`.

---

## 4. Run the System Metrics Monitor

To monitor CPU, RAM, and GPU utilization dynamically and save them to a dedicated PostgreSQL schema (`metrics.system_metrics`):
```bash
source .venv/bin/activate
python agents/metrics.py
```
- Live resource stats are displayed dynamically in the terminal.
- Press `Ctrl + C` to stop monitoring gracefully at any time.
- Logs are appended to `agents/metrics.log`.
