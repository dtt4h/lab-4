FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY backend ./backend
RUN CGO_ENABLED=0 go build -o /api ./backend/cmd
FROM alpine:3.21
COPY --from=build /api /api
COPY backend/migrations /migrations
EXPOSE 8080
CMD ["/api"]
