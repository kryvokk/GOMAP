FROM golang:1.27.3-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/gomal-scan ./cmd/gomal-scan

FROM alpine:3.22
RUN addgroup -S scanner && adduser -S -G scanner scanner
COPY --from=build /out/gomal-scan /usr/local/bin/gomal-scan
USER scanner
ENTRYPOINT ["/usr/local/bin/gomal-scan"]
