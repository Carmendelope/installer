#
# Copyright (C) 2018 Nalej Group - All Rights Reserved
#

# iron/go is the alpine image with only ca-certificates added
FROM iron/go

# Now just add the binary\
RUN mkdir -p /nalej/bin

COPY dockerenv/rke /nalej/bin/
COPY installer /nalej/


RUN chmod +x /nalej/bin/rke

ENTRYPOINT ["./nalej/installer"]