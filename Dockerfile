ARG GO_VERSION=1.24
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH
ARG TARGETOS

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        build-base \
        gcc-aarch64-linux-gnu \
        binutils-aarch64-linux-gnu \
        && \
        update-ca-certificates

COPY . .

RUN if [ "$TARGETARCH" = "arm64" ]; then \
        export CC=aarch64-linux-gnu-gcc \
        CXX=aarch64-linux-gnu-g++; \
    fi; \
    CGO_ENABLED=1 GOARCH=$TARGETARCH GOOS=$TARGETOS go build -o /bin/server ./cmd/main.go

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        curl \
        && \
        update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /src/docs /docs
COPY --from=build /bin/server /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]