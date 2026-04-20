FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /mailchimp-mcp-server ./cmd/mailchimp-mcp-server

FROM gcr.io/distroless/static:nonroot

COPY --from=build /mailchimp-mcp-server /mailchimp-mcp-server

ENTRYPOINT ["/mailchimp-mcp-server"]
