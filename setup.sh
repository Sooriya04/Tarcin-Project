#!/bin/bash
set -e

echo "========================================="
echo " Starting Tarcin Project Dependency Setup"
echo "========================================="

# 1. Setup Go backend dependencies
echo -e "\n1. Installing Go backend modules..."
if [ -d "backend" ]; then
    cd backend
    go mod tidy
    cd ..
    echo "Go modules installed successfully!"
else
    echo "Warning: 'backend' directory not found. Skipping Go modules setup."
fi

# 2. Setup Python virtual environment & dependencies
echo -e "\n2. Initializing Python virtual environment..."
if [ ! -d ".venv" ]; then
    python3 -m venv .venv
    echo "Created virtual environment in .venv/"
else
    echo "Virtual environment (.venv/) already exists."
fi

# Activate virtual environment
source .venv/bin/activate

echo -e "\n3. Installing/Upgrading python packages..."
pip install --upgrade pip
pip install langchain langchain-ollama psycopg2-binary python-dotenv grpcio grpcio-tools markdown-it-py psutil

echo -e "\n========================================="
echo " Tarcin Dependency Setup Completed Successfully!"
echo "========================================="
