# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM            golang:1.15.2-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/zapconfig
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM alpine:3.12
LABEL           org.label-schema.build-date=$BUILD_DATE \
                org.label-schema.name="zapconfig" \
                org.label-schema.description="" \
                org.label-schema.url="https://moul.io/zapconfig/" \
                org.label-schema.vcs-ref=$VCS_REF \
                org.label-schema.vcs-url="https://github.com/moul/zapconfig" \
                org.label-schema.vendor="Manfred Touron" \
                org.label-schema.version=$VERSION \
                org.label-schema.schema-version="1.0" \
                org.label-schema.cmd="docker run -i -t --rm moul/zapconfig" \
                org.label-schema.help="docker exec -it $CONTAINER zapconfig --help"
COPY            --from=builder /go/bin/zapconfig /bin/
ENTRYPOINT      ["/bin/zapconfig"]
#CMD             []
