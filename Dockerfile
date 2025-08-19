ARG GO_VERSION=1.24
FROM golang:${GO_VERSION} AS build
WORKDIR /src

# install build tools (Debian apt)
RUN --mount=type=cache,target=/go/pkg/mod \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential gcc ca-certificates pkg-config \
        libssl-dev && \
    rm -rf /var/lib/apt/lists/*

# download modules (cacheable)
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

# build with CGO enabled so wasmtime-go can link
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=1 go build -o /bin/server ./cmd/main.go

FROM debian:bookworm-slim AS final

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata curl && \
    rm -rf /var/lib/apt/lists/*

ARG UID=10001
RUN adduser --disabled-password --gecos "" --home "/nonexistent" \
    --shell "/sbin/nologin" --no-create-home --uid "${UID}" appuser
USER appuser

COPY --from=build /bin/server /bin/

EXPOSE 8080
ENTRYPOINT [ "/bin/server" ]
