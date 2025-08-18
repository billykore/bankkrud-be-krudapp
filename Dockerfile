FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN go install github.com/google/wire/cmd/wire@latest
RUN go mod download
RUN wire ./cmd
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ./binary ./cmd

FROM alpine:3.21
COPY --from=builder /app/binary /
EXPOSE 8000
ENTRYPOINT ["./binary"]
