FROM ubuntu:xenial
RUN apt-get update -y && apt-get install -y ca-certificates curl git && rm -rf /var/lib/apt/lists/*
RUN \
  curl -Lo /usr/local/bin/meta4 https://github.com/dpb587/metalink/releases/download/v0.1.0/meta4-0.1.0-linux-amd64 \
  && echo '5f09dc2fbf0b766b3645fb25781977303863893b9385b1074b8f989ca77e68ec  /usr/local/bin/meta4' | sha256sum -c \
  && chmod +x /usr/local/bin/meta4
RUN \
  curl -Lo /usr/local/bin/bosh https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.26-linux-amd64 \
  && echo '81673e467a90a22a1365a5981c137b48eee10db8ed2af6adb67191ce35d12c73  /usr/local/bin/bosh' | sha256sum -c \
  && chmod +x /usr/local/bin/bosh
RUN \
  curl -Lo /tmp/download https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz \
  && echo '1862f4c3d3907e59b04a757cfda0ea7aa9ef39274af99a784f5be843c80c6772  /tmp/download' | sha256sum -c \
  && mkdir /usr/local/go \
  && tar -xzf /tmp/download -C /usr/local/go --strip-components=1 \
  && rm /tmp/download
ENV GOROOT /usr/local/go
ENV PATH $GOROOT/bin:$PATH
ADD bcr /usr/local/bin/bcr
