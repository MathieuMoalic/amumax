# Legacy ubuntu version to be compatible with PCSS
FROM nvidia/cuda:11.0.3-devel-ubuntu16.04
RUN apt-get update
RUN apt-get install -y wget git

# Installing go
ENV GO_VERSION=1.22.5
RUN wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz
RUN rm go$GO_VERSION.linux-amd64.tar.gz
ENV PATH /usr/local/go/bin:$PATH

WORKDIR /src

ENV CGO_LDFLAGS="-lcufft -lcuda -lcurand -L/usr/local/cuda/lib64/stubs/ -Wl,-rpath -Wl,\$ORIGIN" 
ENV CGO_CFLAGS="-I/usr/local/cuda/include/" 
ENV CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'

RUN git config --global --add safe.directory /src
CMD go build -v \
    -ldflags "-X github.com/MathieuMoalic/amumax/engine.VERSION=$(date -u +'%Y.%m.%d')" \
    -ldflags "-X github.com/MathieuMoalic/amumax/util.VERSION=$(date -u +'%Y.%m.%d')" \
    -o build/amumax