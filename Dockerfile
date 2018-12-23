FROM scratch
MAINTAINER payneio "paul@payne.io"

ADD ca-certificates.crt /etc/ssl/certs/
ADD build/linux-amd64/enom-ddns /

CMD [ "/enom-ddns" ]
