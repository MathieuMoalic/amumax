{
  description = "A flake for the amumax project";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true; # if your dependencies are unfree
          };
        };

        buildAmumax = pkgs:
          pkgs.buildGoModule rec {
            pname = "amumax";
            version = "2024.04.09";

            src = pkgs.fetchFromGitHub {
              owner = "MathieuMoalic";
              repo = "amumax";
              rev = version;
              hash = "sha256-vsOBj8CcGVonHbQJvPupVUbHCqdqZB3Ro+BCWEBtiiA=";
            };

            vendorHash = "sha256-YqB7EofpTqDnqOQ+ARDJNvZVFltAy0j210lbSwEvifw=";

            nativeBuildInputs = [
              pkgs.cudaPackages.cuda_nvcc
              pkgs.addDriverRunpath
            ];

            buildInputs = [
              pkgs.cudaPackages.cuda_cudart
              pkgs.cudaPackages.cuda_nvcc.dev
              pkgs.cudaPackages.libcufft
              pkgs.cudaPackages.libcurand
            ];

            CGO_CFLAGS = [
              "-lcufft"
              "-lcurand"
            ];

            CGO_LDFLAGS = ["-L${pkgs.cudaPackages.cuda_cudart.lib}/lib/stubs/"];

            ldflags = [
              "-s"
              "-w"
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
            pkgs.cudaPackages.cuda_cudart
            pkgs.cudaPackages.cuda_nvcc.dev
            pkgs.cudaPackages.cuda_nvcc
            pkgs.cudaPackages.libcufft
            pkgs.cudaPackages.libcurand
            pkgs.addDriverRunpath
            # Add any other dependencies you need for development
          ];
          CGO_CFLAGS = [
              "-lcufft"
              "-lcurand"
            ];

          CGO_LDFLAGS = ["-L${pkgs.cudaPackages.cuda_cudart.lib}/lib/stubs/"];

          ldflags = [
              "-s"
              "-w"
            ];
          # Set up any environment variables required for development
          # For example, you might need to specify paths for CUDA libraries
          # Environment variables like CGO_CFLAGS and CGO_LDFLAGS may go here if needed for development
        };
      in {
        packages.amumax = buildAmumax pkgs;
        defaultPackage.amumax = buildAmumax pkgs;
        devShell = devEnv; # Provide the development environment for use with `nix develop`
      }
    );
}
