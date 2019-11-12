FROM alpine:3.9.4

COPY dockerenv/rke /nalej/bin/
COPY dockerenv/zerotier-idtool /nalej/bin/
COPY installer /nalej/

RUN chmod +x /nalej/bin/*

ENTRYPOINT ["./nalej/installer"]
