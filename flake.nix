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
            allowUnfree = true; # if your dependencies are unfree
          };
        };
        buildAmumax = pkgs:
          pkgs.buildGoModule rec {
            pname = "amumax";
            version = "2023.12.14";

            src = pkgs.fetchFromGitHub {
              owner = "MathieuMoalic";
              repo = "amumax";
              rev = version;
              hash = "sha256-U9e8DvgAb5/e2JTDI0yXPF9ollixax3JjeyEFiJbesM=";
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
      in {
        packages.amumax = buildAmumax pkgs;
        defaultPackage = buildAmumax pkgs;
      }
    );
}
