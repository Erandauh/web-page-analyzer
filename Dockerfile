# -------- Builder stage --------
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/google/wire/cmd/wire@latest
RUN wire ./   # generate wire_gen.go

# Build with docker
#  Go binary was built against a newer version than runtime container  so need additional flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o web-page-analyzer

# Runtime
FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /app/web-page-analyzer .

ENV TZ=Asia/Colombo
EXPOSE 8080

ENTRYPOINT ["./web-page-analyzer"]