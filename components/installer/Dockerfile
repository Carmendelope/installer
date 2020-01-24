FROM alpine:3.9.4

COPY dockerenv/rke /nalej/bin/
COPY dockerenv/zerotier-idtool /nalej/bin/
COPY installer /nalej/

# Download Istioctl and the kubectl. This is required by istioctl to run.

ARG KUBE_VERSION="1.15.1"

RUN apk add --update ca-certificates && \
    apk add --update -t deps curl && \
    curl -L https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.4.1 sh -  && \
    mv istio-1.4.1 istio && \
    apk del --purge deps && \
    rm /var/cache/apk/*




RUN chmod +x -R /istio/bin \
    && chmod +x /nalej/bin/*

ENTRYPOINT ["./nalej/installer"]
