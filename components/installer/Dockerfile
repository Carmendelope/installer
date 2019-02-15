#
# Copyright (C) 2019 Nalej Group - All Rights Reserved
#

FROM debian:stretch

# Now just add the binary\
RUN mkdir -p /nalej/bin

COPY dockerenv/rke /nalej/bin/
COPY dockerenv/zerotier-idtool /nalej/bin/
COPY installer /nalej/


RUN chmod +x /nalej/bin/rke
RUN chmod +x /nalej/bin/zerotier-idtool

ENTRYPOINT ["./nalej/installer"]