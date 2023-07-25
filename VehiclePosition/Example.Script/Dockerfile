FROM ubuntu:focal as builder

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Berlin
RUN apt-get update && apt-get install -y build-essential git autoconf libtool libssl-dev pkg-config

WORKDIR /app
RUN git clone https://github.com/obgm/libcoap --recursive .

RUN ./autogen.sh && ./configure --with-openssl --disable-doxygen --disable-manpages --disable-shared && make && make install

FROM ubuntu:focal

RUN apt-get update && apt-get install -y mosquitto-clients jq protobuf-compiler

COPY --from=builder /usr/local/bin/coap-client-openssl /usr/local/bin/coap-client-openssl

WORKDIR /app

COPY ./ ./

CMD ["./post_vehicle_position"]
