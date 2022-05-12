FROM matmoa/amumax:base
WORKDIR /src
ENV PATH /usr/local/go/bin:$PATH
ENV NVCC_CCBIN=/usr/bin/gcc
ENV CGO_CFLAGS="-I${LD_LIBRARY_PATH}"
ENV CGO_LDFLAGS="-lcufft -lcurand -lcuda -L${LD_LIBRARY_PATH}"
ENV CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'
ENV NVCCFLAGS="-std=c++03 -ccbin=$NVCC_CCBIN --compiler-options -Werror --compiler-options -Wall -Xptxas -O3 -ptx"
RUN apt-get install -y git
COPY . .
RUN bash server/build_node.sh
ENTRYPOINT ["/src/amumax"]
# CMD ["/src/amumax"]
# CMD ["/bin/bash","/src/server/build_node.sh"]