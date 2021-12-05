FROM nvidia/cuda:11.4.2-devel-ubuntu20.04

RUN apt update && apt-get install -y wget
RUN wget https://go.dev/dl/go1.17.4.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.4.linux-amd64.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"

ENTRYPOINT tail -f /dev/null
