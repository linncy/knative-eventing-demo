# ---- build stage ----
FROM golang:1.22 AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server/main.go

# ---- runtime stage ----
FROM gcr.io/distroless/static:nonroot
USER 65532:65532
COPY --from=builder /out/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
