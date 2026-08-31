# Builder
FROM golang:1.26-alpine3.23 AS builder

ARG LDFLAGS

WORKDIR /demo

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o demo -ldflags "${LDFLAGS}"

# Runner
FROM alpine:3.24 AS runner

WORKDIR /demo

COPY --from=builder /demo/demo ./demo

RUN adduser -u 1000 -H -D demo

EXPOSE 3000

USER demo

ENV PATH=$PATH:/demo

ENTRYPOINT ["demo"]
