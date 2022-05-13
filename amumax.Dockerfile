FROM matmoa/amumax:build as build
FROM debian:buster-slim 
COPY --from=build /src/amumax /src/amumax
COPY --from=build /src/libcufft.so.10 /src/libcufft.so.10
COPY --from=build /src/libcurand.so.10 /src/libcurand.so.10
ENTRYPOINT ["/src/amumax"]