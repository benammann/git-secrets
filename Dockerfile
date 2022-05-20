# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /git-secrets
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o git-secrets main.go

FROM bash
WORKDIR /git-secrets

RUN apk add --no-cache bash
RUN apk add curl
RUN apk add git

RUN mkdir -p bin

COPY --from=builder /git-secrets/git-secrets bin/git-secrets

RUN adduser -D gitsecrets
RUN chown -R 1000:1000 .
USER gitsecrets

ENV PATH="/git-secrets/bin:${PATH}"

ENTRYPOINT ["git-secrets"]
