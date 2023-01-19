ARG GO_VERSION
FROM golang:${GO_VERSION} as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build
COPY . .
RUN make build


FROM alpine
RUN apk --no-cache add util-linux coreutils && apk update && apk upgrade
COPY --from=builder /build/bin/hostpathplugin /hostpathplugin
ENTRYPOINT [ "/hostpathplugin" ]
