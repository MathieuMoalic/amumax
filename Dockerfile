FROM nvidia/cuda:10.2-devel-ubuntu18.04
WORKDIR /src
RUN apt-get update
RUN apt-get install -y wget
RUN wget https://go.dev/dl/go1.17.7.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.7.linux-amd64.tar.gz
ENV PATH /usr/local/go/bin:$PATH
ENV NAME=amumax
ENV CUDAVERSION=10.2
ENV OS=linux
ENV BUILDDIR=./build/${NAME}_${OS}_${CUDAVERSION}
ENV RPATH=lib 
ENV NVCC_CCBIN=/usr/bin/gcc
ENV PATH=/usr/local/go/bin/:$PATH
ENV CGO_CFLAGS="-I${LD_LIBRARY_PATH}"
ENV CGO_LDFLAGS="-lcufft -lcurand -lcuda -L${LD_LIBRARY_PATH} -Wl,-rpath -Wl,\$ORIGIN/$RPATH"
ENV CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'
ENV NVCCFLAGS="-std=c++03 -ccbin=$NVCC_CCBIN --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"

# SHELL ["/bin/bash", "-c"]
CMD ["/bin/bash","/src/build_release.sh"]