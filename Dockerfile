FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY ./ ./

RUN apk add --update --no-cache git
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/fourkeys" github.com/italolelis/fourkeys/cmd/fourkeys

# ---

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/dist/fourkeys /

EXPOSE 8080 9090
ENTRYPOINT ["/fourkeys"]
