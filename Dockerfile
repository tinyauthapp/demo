# ---- build stage -----------------------------------------------------------
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go page.html ./
# Static, stripped binary so it can run in a scratch image with no libc.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /tinyauth-demo .

# ---- runtime stage ---------------------------------------------------------
FROM scratch
COPY --from=build /tinyauth-demo /tinyauth-demo
# Run as an unprivileged UID (no /etc/passwd needed in scratch).
USER 65534:65534
ENV LOGOUT_URL=/logout
EXPOSE 3000
ENTRYPOINT ["/tinyauth-demo"]
