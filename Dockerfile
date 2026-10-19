FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/x360-dashboard ./cmd/x360-dashboard

FROM alpine:3.22
RUN addgroup -S app && adduser -S -G app app
COPY --from=build /out/x360-dashboard /usr/local/bin/x360-dashboard
USER app
EXPOSE 10043/tcp
ENTRYPOINT ["/usr/local/bin/x360-dashboard"]
