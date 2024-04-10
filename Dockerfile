FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22-alpine as base

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  small-user

WORKDIR $GOPATH/src/broadlinkac/app/

COPY . .

RUN go mod download
RUN go mod verify

ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/main .

FROM scratch

COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=base /out/main .
COPY ./config/config.yml ./config/config.yml

USER small-user:small-user

CMD ["/main"]
