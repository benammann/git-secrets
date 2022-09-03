# Build the manager binary
FROM golang:1.19 as builder

WORKDIR /git-secrets
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

ARG BUILD_VERSION="n/a"
ARG BUILD_COMMIT="n/a"
ARG BUILD_DATE="n/a"

ENV BUILD_VERSION="$BUILD_VERSION"
ENV BUILD_COMMIT="$BUILD_COMMIT"
ENV BUILD_DATE="$BUILD_DATE"

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o git-secrets -ldflags "-X main.version=$BUILD_VERSION -X main.commit=$BUILD_COMMIT -X main.date=$BUILD_DATE"

FROM alpine:latest
WORKDIR /git-secrets

RUN apk --no-cache add ca-certificates git

RUN mkdir -p bin

COPY --from=builder /git-secrets/git-secrets bin/git-secrets

RUN adduser -D gitsecrets
RUN chown -R 1000:1000 .
USER gitsecrets

ENV PATH="/git-secrets/bin:${PATH}"

ENTRYPOINT ["git-secrets"]
