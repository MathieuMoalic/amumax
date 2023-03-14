default: build copy_local copy_pcss

build:
	DOCKER_HOST= docker build -t matmoa/amumax:build .
	DOCKER_HOST= docker run --rm -v $PWD:/src matmoa/amumax:build
	mv ./build/amumax ./build/amumax$(date -I)


copy_pcss:
	scp -r ./build/* pcss:grant_398/scratch/bin/amumax_versions
	ssh pcss "cd ~/grant_398/scratch/bin && ln -sf amumax_versions/amumax$(date -I) amumax"

copy_local:
	cp -r ./build/* ~/.local/bin
	ln -sf ~/.local/bin/amumax$(date -I) ~/.local/bin/amumax

test:
	go build -v . && ./amumax -f mytest/chunky.mx3
	ls -l mytest/chunky.zarr/m/
	cat mytest/chunky.zarr/m/.zarray
	. ~/ws/.env/bin/activate
	python mytest/chunky.py