FROM alpine:3.6
MAINTAINER Alex Kern <alex@pavlov.ai>

RUN mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    apk add --no-cache curl ca-certificates && \
    curl -o /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.6.4/bin/linux/amd64/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    kubectl version --client

COPY build /build
COPY docker-entrypoint.sh /

EXPOSE 80
CMD /docker-entrypoint.sh
