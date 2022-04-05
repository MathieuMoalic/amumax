build:
	go build -v
release:
	sudo docker build -t amumax . 
	sudo docker run --rm -v "$$PWD":/src amumax
	scp -r ./build/amumax_linux_10.2/* pcss:bin

win-release:
	go build -v
	# cp 