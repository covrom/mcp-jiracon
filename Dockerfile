FROM golang:alpine AS build

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath \
    -ldflags "-s -w" \
    -o /out/mcp-atlassian ./cmd/mcp-jiracon

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=build /out/mcp-atlassian /usr/local/bin/mcp-atlassian

ENTRYPOINT ["/usr/local/bin/mcp-atlassian"]
CMD ["-transport=stdio"]