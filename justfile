image:
	podman build -t matmoa/amumax:build .

build-cuda: 
	podman run --rm -v $PWD:/src matmoa/amumax:build sh src/cuda/build-cuda.sh

copy-pcss:
	scp -r ./build/amumax pcss:grant_398/scratch/bin/amumax_versions/amumax$(date -I)
	ssh pcss "cd ~/grant_398/scratch/bin && ln -sf amumax_versions/amumax$(date -I) amumax"

build-frontend: 
	rm -rf api/static
	podman run --rm \
		-v .:/src \
		-w /src/frontend \
		--entrypoint /bin/sh \
		docker.io/node:18.20.4-alpine3.20 -c 'npm install && npm run build && mv dist ../src/api/static'

build:
	podman run --rm -v $PWD:/src matmoa/amumax:build

update-flake-hashes VERSION:
	#!/usr/bin/env sh
	set -euxo pipefail
	sed -i 's/releaseVersion = "[^"]*"/releaseVersion = "'"{{VERSION}}"'"/' flake2.nix

	GH_HASH=$(nix-prefetch-github MathieuMoalic amumax --rev {{VERSION}} | jq -r '.hash')
	sed -i "/# gh hash/ s|hash = \".*\";|hash = \"$GH_HASH\";|" flake2.nix

	NPM_HASH=$(prefetch-npm-deps frontend/package-lock.json)
	sed -i "/# npm hash/ s|npmDepsHash = \".*\";|npmDepsHash = \"$NPM_HASH\";|" flake2.nix

pre-commit:
	#!/usr/bin/env sh
	if git diff --cached --name-only | grep -q -e "frontend/package-lock.json" -e "go.sum"; then
		# Check if flake.nix is also staged for commit
		if ! git diff --cached --name-only | grep -q "flake.nix"; then
			echo "Error: lock files have changed, but flake.nix is not updated."
			echo "Please update flake.nix accordingly."
			exit 1  # Block the commit
		fi
	fi
	# If the check passes, allow the commit
	exit 0

test:
	go test ./src/...
	
release: 
	#!/usr/bin/env sh
	just test
	VERSION=$(date -u +'%Y.%m.%d')
	just update-flake $VERSION
	git add flake.nix flake.lock
	just image build-cuda build-frontend build
	just pre-commit
	git commit -m "Release of ${VERSION}"
	git push
	gh release create $VERSION ./build/* --title $VERSION --notes "Release of ${VERSION}"
	just copy-pcss
