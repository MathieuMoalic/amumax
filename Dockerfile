FROM nvidia/cuda:12.1.0-devel-ubuntu22.04
ENV GO_VERSION=1.20.1
RUN apt-get update
RUN apt-get install -y wget git
RUN wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz
RUN rm go$GO_VERSION.linux-amd64.tar.gz
ENV PATH /usr/local/go/bin:$PATH

WORKDIR /src
ENV NVCC_CCBIN=/usr/bin/gcc
ENV CGO_CFLAGS="-I${LD_LIBRARY_PATH}"
ENV CGO_LDFLAGS="-lcufft -lcurand -lcuda -L${LD_LIBRARY_PATH} -Wl,-rpath -Wl,\$ORIGIN/$RPATH"
ENV CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'
ENV NVCCFLAGS="-std=c++03 -ccbin=$NVCC_CCBIN --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"
RUN git config --global --add safe.directory /src
CMD cd cuda && ./build_cuda.sh && cd .. && \
  go build -v -ldflags "-X github.com/MathieuMoalic/amumax/engine.VERSION=$(date -u +'%Y.%m.%d')" && \
  rm -rfd /src/build && \
  mkdir /src/build && \
  cp /src/amumax /src/build && \
  cp $( ldd /src/amumax | grep libcufft | awk '{print $3}' ) /src/build && \
  cp $( ldd /src/amumax | grep libcurand | awk '{print $3}' ) /src/build 