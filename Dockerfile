FROM gcr.io/distroless/static:nonroot

COPY mailchimp-mcp-server /mailchimp-mcp-server

ENTRYPOINT ["/mailchimp-mcp-server"]
