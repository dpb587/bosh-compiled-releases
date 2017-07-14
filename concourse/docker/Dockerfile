FROM alpine:3.6
RUN apk --no-cache add bash ca-certificates curl git
RUN \
  curl -Lo /usr/local/bin/meta4 https://github.com/dpb587/metalink/releases/download/v0.1.0/meta4-0.1.0-linux-amd64 \
  && echo '5f09dc2fbf0b766b3645fb25781977303863893b9385b1074b8f989ca77e68ec  /usr/local/bin/meta4' | sha256sum -c \
  && chmod +x /usr/local/bin/meta4
RUN \
  curl -Lo /usr/local/bin/bosh https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.26-linux-amd64 \
  && echo '81673e467a90a22a1365a5981c137b48eee10db8ed2af6adb67191ce35d12c73  /usr/local/bin/bosh' | sha256sum -c \
  && chmod +x /usr/local/bin/bosh
ADD bcr /usr/local/bin/bcr
