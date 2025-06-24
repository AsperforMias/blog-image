FROM golang:1.24.4 AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-image .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/blog-image /blog-image

ENTRYPOINT ["/blog-image"]
