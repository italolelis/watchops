FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY ./ ./

RUN apk add --update --no-cache git
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/publisher" github.com/italolelis/watchops/cmd/publisher

# ---

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/dist/publisher /

EXPOSE 8080 9090
ENTRYPOINT ["/publisher"]
