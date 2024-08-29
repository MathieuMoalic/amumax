{
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
        buildAmumax = pkgs.buildGoModule rec {
          pname = "amumax";
          version = "2024.08.29";

          # src = pkgs.fetchFromGitHub {
          #   owner = "MathieuMoalic";
          #   repo = "amumax";
          #   rev = version;
          #   hash = "sha256-KfagOhaVmps5WLANatQPNaDELPyjXzBwyZ3EBuGtExw=";
          # };
          vendorHash = "sha256-SHUBKLKV8lwjyXlhM5OyHpwvm1s/yo9I3+Bow+MwRc0=";
          src = ./.;

          buildInputs = [
            cuda.cuda_nvcc
            cuda.cuda_cudart
            cuda.libcufft
            cuda.libcurand
            pkgs.addDriverRunpath
            pkgs.bun
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
        };

        devEnv = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.golangci-lint
            cuda.cuda_cudart
            cuda.cuda_nvcc
            cuda.libcufft
            cuda.libcurand
            pkgs.gcc11
            pkgs.bun
          ];
          CGO_LDFLAGS = "-lcufft -lcuda -lcurand -Wl,-rpath -Wl,\$ORIGIN";
          CGO_CFLAGS_ALLOW = "(-fno-schedule-insns|-malign-double|-ffast-math)";
          LD_LIBRARY_PATH = "${cuda.libcufft}/lib:${cuda.libcurand}/lib:/run/opengl-driver/lib/";
          ldflags = ["-s" "-w"];
          shellHook = ''
            export GOPATH=$(pwd)/.go/path
            export GOCACHE=$(pwd)/.go/cache
            mkdir -p $GOPATH $GOCACHE
          '';
        };
      in {
        packages.default = buildAmumax;
        devShell = devEnv;
      }
    );
}
