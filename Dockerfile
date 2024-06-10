# ======
# Build
# ======
FROM golang:1.22-bookworm AS build
ENV TZ=Asia/Tokyo

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY src/ ./internal/

RUN go build -o /bin/main -ldflags="-s -w" -trimpath ./internal

# ======
# Deploy
# ======
FROM gcr.io/distroless/static-debian12 AS deploy
ENV TZ=Asia/Tokyo

COPY --from=build /bin/main /main

EXPOSE 8080
USER nonroot:nonroot

CMD ["/main"]
