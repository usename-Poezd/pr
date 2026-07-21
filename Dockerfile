FROM node:20-alpine AS frontend
WORKDIR /src/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /src/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /server ./cmd/server
FROM alpine:3.19
WORKDIR /app
COPY --from=build /server ./server
COPY --from=build /src/migrations ./migrations
COPY --from=build /src/frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./server"]
