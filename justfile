default: build copy_local copy_pcss

build:
	docker build -t matmoa/amumax:build .
	docker run --rm -v $PWD:/src matmoa/amumax:build

copy_pcss:
	scp -r ./build/* pcss:grant_398/scratch/bin

copy_local:
	cp -r ./build/* ~/.local/bin