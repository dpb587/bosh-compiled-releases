FROM alpine:3.6
RUN apk --no-cache add bash ca-certificates curl git
RUN \
  curl -Lo /usr/local/bin/meta4 https://github.com/dpb587/metalink/releases/download/v0.1.0/meta4-0.1.0-linux-amd64 \
  && echo '5f09dc2fbf0b766b3645fb25781977303863893b9385b1074b8f989ca77e68ec  /usr/local/bin/meta4' | sha256sum -c \
  && chmod +x /usr/local/bin/meta4
ADD bcr /usr/local/bin/bcr
