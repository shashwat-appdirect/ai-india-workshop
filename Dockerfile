# Multi-stage Dockerfile for building both frontend and backend

# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Install dependencies
RUN npm ci

# Copy frontend source code
COPY frontend/ .

# Set API base URL to relative path for production (frontend and backend on same domain)
ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# Build the frontend
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app/backend

# Install git (required for some Go dependencies)
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ .

# Build the backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Stage 3: Production Image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy backend binary from builder
COPY --from=backend-builder /app/backend/server .

# Copy frontend static files from builder
COPY --from=frontend-builder /app/frontend/dist ./static

# Expose port (Cloud Run uses PORT environment variable)
EXPOSE 8080

# Set environment variable for static files directory
ENV STATIC_DIR=/app/static
ENV PORT=8080

# Run the server
CMD ["./server"]

