FROM nvidia/cuda:11.6.0-devel-ubuntu20.04
RUN apt-get update
RUN apt-get install -y wget git
RUN wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz
RUN rm go1.18.1.linux-amd64.tar.gz
ENV PATH /usr/local/go/bin:$PATH

WORKDIR /src
ENV PATH /usr/local/go/bin:$PATH
ENV NVCC_CCBIN=/usr/bin/gcc
ENV CGO_CFLAGS="-I${LD_LIBRARY_PATH}"
ENV CGO_LDFLAGS="-lcufft -lcurand -lcuda -L${LD_LIBRARY_PATH} -Wl,-rpath -Wl,\$ORIGIN/$RPATH"
ENV CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'
ENV NVCCFLAGS="-std=c++03 -ccbin=$NVCC_CCBIN --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"
COPY . .
RUN cd cuda && ./build_cuda.sh && cd .. 
RUN go build -v
RUN cp $( ldd /src/amumax | grep libcufft | awk '{print $3}' ) /src
RUN cp $( ldd /src/amumax | grep libcurand | awk '{print $3}' ) /src