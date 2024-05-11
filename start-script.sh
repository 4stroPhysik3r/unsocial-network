#!/bin/bash
clear

# Start backend
cd backend
go run main.go &

cd ..
# Start frontend
cd frontend
npm install
npm run serve
