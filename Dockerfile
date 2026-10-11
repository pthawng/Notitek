# syntax=docker/dockerfile:1
# One Dockerfile for every runtime: docker build --build-arg RUNTIME=worker .

FROM golang:1.26-alpine AS build
ARG RUNTIME=api
WORKDIR /src
ENV CGO_ENABLED=0 GOFLAGS=-trimpath

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o /out/notitek ./cmd/${RUNTIME}

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/notitek /notitek
USER nonroot:nonroot
EXPOSE 8080 9090
ENTRYPOINT ["/notitek"]
