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
            version = "2024.07.23";

            src = pkgs.fetchFromGitHub {
              owner = "MathieuMoalic";
              repo = "amumax";
              rev = version;
              hash = "sha256-KfagOhaVmps5WLANatQPNaDELPyjXzBwyZ3EBuGtExw=";
            };
            vendorHash = "sha256-SHUBKLKV8lwjyXlhM5OyHpwvm1s/yo9I3+Bow+MwRc0=";
            # src = builtins.fetchGit {
            #   path = ./.;
            #   rev = "84848b4b467e4948e753af393ac83ef90a076373";
            # };
            # src = builtins.path {
            #   path = ./.;
            # };

            nativeBuildInputs = [
              cuda.cuda_nvcc
              pkgs.addDriverRunpath
              pkgs.bun
            ];

            buildInputs = [
              cuda.cuda_cudart
              cuda.libcufft
              cuda.libcurand
            ];

            CGO_CFLAGS = ["-lcufft" "-lcurand"];
            CGO_LDFLAGS = ["-L${cuda.cuda_cudart.lib}/lib/stubs/"];
            ldflags = [
              "-s"
              "-w"
              "-X github.com/MathieuMoalic/amumax/engine.VERSION=${version}"
              "-X github.com/MathieuMoalic/amumax/util.VERSION=${version}"
            ];

            doCheck = false;

            # preBuild = ''
            #   cd ${src}/frontend
            #   ls -la
            #   bun install
            #   bun run build
            # '';

            postFixup = ''
              addDriverRunpath $out/bin/*
            '';

            meta = with pkgs.lib; {
              description = "Fork of mumax3";
              homepage = "https://github.com/MathieuMoalic/amumax";
              license = licenses.gpl3;
              maintainers = ["MathieuMoalic"];
              mainProgram = "amumax";
            };
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
        packages.amumax = buildAmumax pkgs;
        packages.default = buildAmumax pkgs;
        defaultPackage.amumax = buildAmumax pkgs;
        devShell = devEnv;
      }
    );
}
