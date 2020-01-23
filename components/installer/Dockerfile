FROM alpine:3.9.4

COPY dockerenv/rke /nalej/bin/
COPY dockerenv/zerotier-idtool /nalej/bin/
COPY installer /nalej/

RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

RUN curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.4.1 sh -  \
    && mv istio-1.4.1 istio

RUN chmod +x /nalej/bin/*

ENTRYPOINT ["./nalej/installer"]
