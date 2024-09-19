image:
	podman build -t matmoa/amumax:build .

build_cuda: 
	podman run --rm -v $PWD:/src matmoa/amumax:build sh cuda/build_cuda.sh

copy_pcss:
	scp -r ./build/amumax pcss:grant_398/scratch/bin/amumax_versions/amumax$(date -I)
	ssh pcss "cd ~/grant_398/scratch/bin && ln -sf amumax_versions/amumax$(date -I) amumax"

build-frontend: 
	rm -rf api/static
	podman run --rm \
		-v .:/src \
		-w /src/frontend \
		--entrypoint /bin/sh \
		docker.io/node:18.20.4-alpine3.20 -c 'npm install && npm run build && mv dist ../api/static'

build:
	podman run --rm -v $PWD:/src matmoa/amumax:build

release: image build_cuda build-frontend build copy_pcss
	VERSION=$(date -u +'%Y.%m.%d') && \
	echo $VERSION && \
	sed -i 's/releaseVersion = "[^"]*"/releaseVersion = "'"$VERSION"'"/' flake.nix && \
	git add flake.nix flake.lock && \
	git commit -m "Release of ${VERSION}" && \
	git push && \
	gh release create $VERSION ./build/* --title $VERSION --notes "Release of ${VERSION}"
	just copy_pcss
