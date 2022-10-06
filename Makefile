release:
	docker build -t matmoa/amumax:build -f build.Dockerfile .
	docker run --rm -v "$$PWD":/src matmoa/amumax:build 

pcss:
	docker build -t matmoa/amumax:pcss -f pcss.Dockerfile .
	docker run --rm -v "$$PWD":/src matmoa/amumax:pcss 
	scp -r ./build/11.2/* pcss:grant_398/scratch/bin

dbuild:
	docker build -t matmoa/amumax:build -f build.Dockerfile .
damumax:
	docker build -t matmoa/amumax:amumax -f amumax.Dockerfile .

.PHONY: cuda
cuda:
	cd cuda && ./build_cuda.sh
