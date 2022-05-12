FROM nvidia/cuda:11.6.0-devel-ubuntu20.04
RUN apt-get update
RUN apt-get install -y wget git
RUN wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz
RUN rm go1.18.1.linux-amd64.tar.gz
ENV PATH /usr/local/go/bin:$PATH