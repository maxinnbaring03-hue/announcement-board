# ==========================================
# STAGE 1: Build the Svelte Frontend
# ==========================================
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend

# Copy package files and install dependencies
COPY frontend/package*.json ./
RUN npm install

# Copy the rest of the frontend code and build it
# (Assuming Vite builds your Svelte app into a 'dist' folder)
COPY frontend/ .
RUN npm run build


# ==========================================
# STAGE 2: Build the Go Backend
# ==========================================
FROM golang:alpine AS backend-builder
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go code from the root directory
COPY . .

# Compile the Go application
RUN go build -o server .


# ==========================================
# STAGE 3: The Final Single Container
# ==========================================
FROM alpine:latest
WORKDIR /root/

# 1. Grab the compiled Go server from Stage 2
COPY --from=backend-builder /app/server .

# 2. Grab the compiled Svelte UI from Stage 1 and put it in a folder called 'public'
COPY --from=frontend-builder /app/frontend/dist ./public

# Expose the single port your Go server uses
EXPOSE 8080

# Turn on the Go server (which now powers the backend AND serves the Svelte frontend)
CMD ["./server"]