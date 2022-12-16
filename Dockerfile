FROM golang:1.19.4-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o teslanotify ./cmd/teslanotify

FROM alpine:latest
WORKDIR /app/
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=build /app/teslanotify .
USER appuser
ENTRYPOINT ["/app/teslanotify"]
