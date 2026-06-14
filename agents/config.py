import os
import psycopg2
from dotenv import load_dotenv

# Load database configuration from backend/.env
load_dotenv('backend/.env')

DB_URL = os.getenv("DB_URL", "postgres://postgres:localhost@localhost:5432/tarcin?sslmode=disable")
GRPC_SERVER_ADDR = os.getenv("GRPC_SERVER_ADDR", "localhost:50051")

def get_db_connection():
    return psycopg2.connect(DB_URL)
