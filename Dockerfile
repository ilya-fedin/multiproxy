FROM alpine

ENV VERSION 0.8.12
ENV DOCKER true

RUN apk add --no-cache --virtual .build-deps gcc linux-headers make musl-dev tar openssl && \
  apk add --no-cache s6 su-exec && \
  wget -O 3proxy.tar.gz "https://github.com/z3APA3A/3proxy/archive/${VERSION}.tar.gz" && \
  mkdir -p /usr/src/3proxy && \
  tar -xzf 3proxy.tar.gz -C /usr/src/3proxy --strip-components=1 && \
  rm 3proxy.tar.gz && \
  make -C /usr/src/3proxy -f Makefile.Linux && \
  make -C /usr/src/3proxy -f Makefile.Linux install && \
  rm -r /usr/src/3proxy && \
  apk del .build-deps

COPY run.sh /usr/local/bin
COPY 3proxy.cfg /etc
COPY instances /etc

ENTRYPOINT ["run.sh"]
