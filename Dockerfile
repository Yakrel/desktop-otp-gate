FROM golang:alpine as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"

FROM scratch
ENV EXIT_RESET=1
COPY --from=builder /build/simple-nginx-otp /simple-nginx-otp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/simple-nginx-otp"]
