release:
	# sudo docker build -t amumax . 
	sudo docker run --rm -v "$$PWD":/src amumax
	scp -r ./build/amumax_linux_10.2/* pcss:bin

dbuild:
	docker build -t matmoa/amumax:build -f build.Dockerfile .
damumax:
	docker build -t matmoa/amumax:amumax -f amumax.Dockerfile .

.PHONY: cuda
cuda:
	cd cuda && ./build_cuda.sh
