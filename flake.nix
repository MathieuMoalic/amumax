{
  description = "A flake for the amumax project";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true; # cuda is unfree
          };
        };
        cuda = pkgs.cudaPackages_11;
        buildAmumax = pkgs:
          pkgs.buildGoModule rec {
            pname = "amumax";
            version = "2024.07.16";

            src = pkgs.fetchFromGitHub {
              owner = "MathieuMoalic";
              repo = "amumax";
              rev = version;
              hash = "sha256-tJSu77e7755nNQ/bvBZgv5hqiYNHW0hNGNhaNrvZWGM=";
            };

            vendorHash = "sha256-GAtFL46BvfI/s3coVGBsBtelZAC8xJpRfSjhwhODQNk=";

            nativeBuildInputs = [
              cuda.cuda_nvcc
              pkgs.addDriverRunpath
            ];

            buildInputs = [
              cuda.cuda_cudart
              cuda.cuda_nvcc.dev
              cuda.libcufft
              cuda.libcurand
            ];
            CGO_CFLAGS = ["-lcufft" "-lcurand"];
            CGO_LDFLAGS = ["-L${cuda.cuda_cudart.lib}/lib/stubs/"];
            ldflags = [
              "-s"
              "-w"
              "-X github.com/MathieuMoalic/amumax/engine.VERSION=${version}"
            ];

            doCheck = false;

            postFixup = ''
              addDriverRunpath $out/bin/*
            '';

            meta = with pkgs.lib; {
              description = "Fork of mumax3";
              homepage = "https://github.com/MathieuMoalic/amumax/tree/main";
              license = licenses.gpl3;
              maintainers = [];
              mainProgram = "amumax";
            };
          };

        devEnv = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.golangci-lint
            cuda.cuda_cudart
            cuda.cuda_nvcc.dev
            cuda.cuda_nvcc
            cuda.libcufft
            cuda.libcurand
            pkgs.addDriverRunpath
          ];
          CGO_CFLAGS = ["-lcufft" "-lcurand"];
          CGO_LDFLAGS = ["-L${cuda.cuda_cudart.lib}/lib/stubs/"];
          ldflags = ["-s" "-w"];

          shellHook = ''
            export LD_LIBRARY_PATH=${cuda.libcufft}/lib:${cuda.libcurand}/lib:/run/opengl-driver/lib/:$LD_LIBRARY_PATH
          '';
        };
      in {
        packages.amumax = buildAmumax pkgs;
        packages.default = buildAmumax pkgs;
        defaultPackage.amumax = buildAmumax pkgs;
        devShell = devEnv;
      }
    );
}
