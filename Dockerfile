FROM harbor.digitalservice.id/proxyjds/library/golang:1.21.6-alpine AS builder
WORKDIR /app 
COPY . .
RUN go mod tidy && go build -o main main.go

# Run stage
FROM alpine
WORKDIR /app

COPY --from=builder /app/ ./

EXPOSE 3333
CMD ["/app/main"]