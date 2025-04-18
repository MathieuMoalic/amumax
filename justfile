run-dev:
	sudo podman run -it --rm -p 35367:35367 -v $PWD:/src \
	--device=nvidia.com/gpu=all \
	matmoa/amumax:build bash

image:
	sudo podman build -t matmoa/amumax:build .

build-cuda: 
	sudo podman run --rm -v $PWD:/src matmoa/amumax:build sh src/cuda/build_cuda.sh

copy-pcss:
	scp -r ./build/amumax pcss:grant_398/scratch/bin/amumax_versions/amumax$(date -I)
	ssh pcss "cd ~/grant_398/scratch/bin && ln -sf amumax_versions/amumax$(date -I) amumax"

build-frontend: 
	cd frontend && npm run build && rm -rf ../src/api/static && mv dist ../src/api/static
	# rm -rf api/static
	# sudo podman run --rm \
	# 	-v .:/src \
	# 	-w /src/frontend \
	# 	--entrypoint /bin/sh \
	# 	docker.io/node:18.20.4-alpine3.20 -c 'npm install && npm run build && rm -rf ../src/api/static && mv dist ../src/api/static'

build:
	sudo podman run --rm -v $PWD:/src matmoa/amumax:build

update-flake-gh-hash VERSION:
	#!/usr/bin/env sh
	set -euxo pipefail
	sed -i 's/releaseVersion = "[^"]*"/releaseVersion = "'"{{VERSION}}"'"/' flake.nix

	GH_HASH=$(nix-prefetch-github MathieuMoalic amumax --rev {{VERSION}} | jq -r '.hash')
	escaped_hash=$(printf '%s' "$GH_HASH" | sed 's/[&/\]/\\&/g')
	sed -i "s/hash = pkgs.lib.fakeHash;/hash = \"$escaped_hash\";/" flake.nix

test:
	go test ./src/...
	
release: 
	#!/usr/bin/env sh
	set -euxo pipefail
	git checkout main

		if [ -n "$(git status --porcelain)" ]; then
		echo "Working directory is not clean. Please commit or stash your changes."
		exit 1
	fi
	
	git pull
	VERSION=$(date -u +'%Y.%m.%d')
	gh release view $VERSION &>/dev/null && gh release delete $VERSION -y
	git show-ref --tags $VERSION &>/dev/null && git tag -d $VERSION && git push --tags

	just test
	
	just image build-cuda build-frontend build

	# We need to commit before the release
	git add .
	if git diff-index --quiet HEAD --; then
		echo "No changes to commit. Skipping commit step."
	else
		git commit -m "Release of $VERSION"
	fi
	git push
	gh release create $VERSION ./build/* --title $VERSION --notes "Release of ${VERSION}"
	just copy-pcss
	just flake-release

flake-release:
	#!/usr/bin/env sh
	set -euxo pipefail
	VERSION=$(date -u +'%Y.%m.%d')
	just update-flake-hashes-git
	just update-flake-gh-hash ${VERSION}
	nix run . -- -v
	git add .
	git commit -m "Update github hash for the release of ${VERSION}"
	git push

update-flake-hashes-git:
	#!/usr/bin/env sh
	set -euxo pipefail

	echo "Resetting npmDepsHash and vendorHash to placeholder values..."
	sed -i 's/npmDepsHash = "sha256-[^\"]*";/npmDepsHash = pkgs.lib.fakeHash;/' flake.nix
	sed -i 's/vendorHash = "sha256-[^\"]*";/vendorHash = pkgs.lib.fakeHash;/' flake.nix
	sed -i 's/hash = "sha256-[^\"]*";/hash = pkgs.lib.fakeHash;/' flake.nix

	echo "Starting the hash update process..."

	update_hashes() {
		echo "Running nix command to capture output and find new hashes..."
		output=$(nix run .#git -- -v 2>&1 || true)

		new_hash=$(echo "$output" | grep 'got:' | awk '{print $2}')
		escaped_hash=$(printf '%s' "$new_hash" | sed 's/[&/\]/\\&/g')

		if [[ -n "$new_hash" ]]; then
			echo "New hash found: $new_hash"
			if [[ "$output" == *"frontend-git-npm-deps.drv':"* ]]; then
				echo "Updating npmDepsHash in flake.nix..."
				sed -i "s/npmDepsHash = pkgs.lib.fakeHash;/npmDepsHash = \"$escaped_hash\";/" flake.nix
			elif [[ "$output" == *"git-go-modules.drv':"* ]]; then
				echo "Updating vendorHash in flake.nix..."
				sed -i "s/vendorHash = pkgs.lib.fakeHash;/vendorHash = \"$escaped_hash\";/" flake.nix
			else
				echo "Error: None of the expected patterns found in the output." >&2
				return 1
			fi
		else
			echo "Error: No new hash found in the output." >&2
			return 1
		fi
	}

	echo "Updating hashes..."
	update_hashes
	echo "First update completed. Running the update again..."
	update_hashes

	echo "Running final test to verify updated hashes..."
	nix run .#git -- -v

	echo "Hash update process completed."
