FROM golang:alpine AS builder

WORKDIR /work

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" -o /bin/tf-changelog-validator \
  ./cmd/tf-changelog-validator

FROM scratch
LABEL maintiner="Charles Kenney <me@ch4z.io>"

COPY --chown=0:0 --from=builder /bin/tf-changelog-validator /bin/tf-changelog-validator

ENTRYPOINT ["/bin/tf-changelog-validator", "--repoPath=/var/repo"]
