import os
import sys
import time
import psycopg2
import subprocess
import psutil
from dotenv import load_dotenv

# Load database configuration from backend/.env
load_dotenv(os.path.join(os.path.dirname(os.path.abspath(__file__)), '../backend/.env'))

DB_URL = os.getenv("DB_URL", "postgres://postgres:localhost@localhost:5432/tarcin?sslmode=disable")
LOG_PATH = os.path.join(os.path.dirname(os.path.abspath(__file__)), "metrics.log")

def log_to_file(message):
    timestamp = time.strftime("[%Y-%m-%d %H:%M:%S]")
    with open(LOG_PATH, "a") as f:
        f.write(f"{timestamp} {message}\n")

def init_db():
    conn = psycopg2.connect(DB_URL)
    cursor = conn.cursor()
    # Create a completely separate schema and table in postgres to avoid mingling with public schema
    cursor.execute("""
        CREATE SCHEMA IF NOT EXISTS metrics;
        CREATE TABLE IF NOT EXISTS metrics.system_metrics (
            id SERIAL PRIMARY KEY,
            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            cpu_usage REAL,
            ram_usage REAL,
            gpu_usage REAL,
            gpu_mem_usage REAL
        );
    """)
    conn.commit()
    cursor.close()
    conn.close()

def get_gpu_metrics():
    # Attempt to query NVIDIA GPU metrics using nvidia-smi
    try:
        cmd = ["nvidia-smi", "--query-gpu=utilization.gpu,utilization.memory", "--format=csv,noheader,nounits"]
        output = subprocess.check_output(cmd, stderr=subprocess.DEVNULL).decode("utf-8").strip()
        gpu, gpu_mem = map(float, output.split(","))
        return gpu, gpu_mem
    except Exception:
        # Fallback if no Nvidia GPU or nvidia-smi is missing
        return 0.0, 0.0

def main():
    try:
        init_db()
    except Exception as e:
        print(f"Error initializing PostgreSQL metrics database: {e}")
        sys.exit(1)

    log_to_file("Metrics service started using PostgreSQL schema.")
    print("==================================================================")
    print(" Tarcin System Metrics Monitoring Service Started (PostgreSQL) ")
    print(" Schema: metrics | Table: system_metrics")
    print(f" Log File: {LOG_PATH}")
    print(" Press Ctrl+C to stop monitoring.")
    print("==================================================================")
    print(f"{'Timestamp':<21} | {'CPU (%)':<8} | {'RAM (%)':<8} | {'GPU (%)':<8} | {'GPU Mem (%)':<11}")
    print("-" * 65)

    try:
        conn = psycopg2.connect(DB_URL)
        cursor = conn.cursor()
        while True:
            # Gather metrics (cpu_percent with interval=1 consumes 1 second)
            cpu = psutil.cpu_percent(interval=1)
            ram = psutil.virtual_memory().percent
            gpu, gpu_mem = get_gpu_metrics()
            
            timestamp = time.strftime("%Y-%m-%d %H:%M:%S")

            # Store in PostgreSQL under metrics schema
            cursor.execute("""
                INSERT INTO metrics.system_metrics (cpu_usage, ram_usage, gpu_usage, gpu_mem_usage)
                VALUES (%s, %s, %s, %s)
            """, (cpu, ram, gpu, gpu_mem))
            conn.commit()

            # Format and output metrics
            msg = f"{timestamp} | {cpu:<8.1f} | {ram:<8.1f} | {gpu:<8.1f} | {gpu_mem:<11.1f}"
            print(msg)
            log_to_file(f"CPU: {cpu}%, RAM: {ram}%, GPU: {gpu}%, GPU Mem: {gpu_mem}%")

            # Sleep remaining time of the 5-second interval
            time.sleep(4)

    except KeyboardInterrupt:
        print("\nMonitoring stopped by user request.")
        log_to_file("Metrics service stopped by user request.")
    except Exception as e:
        print(f"\nError occurred in metrics loop: {e}")
        log_to_file(f"Error occurred in metrics loop: {e}")
    finally:
        if 'cursor' in locals():
            cursor.close()
        if 'conn' in locals():
            conn.close()

if __name__ == "__main__":
    main()
