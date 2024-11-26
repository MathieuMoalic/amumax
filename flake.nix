{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = {nixpkgs, ...}: let
    system = "x86_64-linux";
    pkgs = import nixpkgs {
      inherit system;
      config = {
        allowUnfree = true; # cuda is unfree
      };
    };

    CGO_CFLAGS = ["-lcufft" "-lcurand"]; # needed to build ptx
    CGO_LDFLAGS = ["-lcuda -lcurand -lcufft -Wl,-rpath -Wl,\$ORIGIN"];
    CGO_CFLAGS_ALLOW = "(-fno-schedule-insns|-malign-double|-ffast-math)";

    cuda = pkgs.cudaPackages_11;
    basepkgs = [
      cuda.cuda_nvcc
      cuda.cuda_cudart
      cuda.libcufft
      cuda.libcurand
      pkgs.bun
    ];

    # Common function to build NPM packages
    buildFrontend = {
      src,
      npmDepsHash,
      version,
    }:
      pkgs.buildNpmPackage {
        inherit version src npmDepsHash;
        pname = "frontend";

        npmBuild = ''
          npm run build
        '';

        installPhase = ''
          mv dist $out
        '';
      };

    # Common function to build Amumax
    buildAmumax = {
      src,
      frontend,
      vendorHash,
      version,
    }:
      pkgs.buildGoModule {
        inherit version CGO_CFLAGS CGO_LDFLAGS CGO_CFLAGS_ALLOW vendorHash src;
        pname = "amumax";

        buildInputs = basepkgs ++ [pkgs.addDriverRunpath];

        buildPhase = ''
          cp -r ${frontend} src/api/static
          go build -v -o $out/bin/amumax -ldflags '-s -w -X github.com/MathieuMoalic/amumax/src/engine.VERSION=${version}' .
        '';

        doCheck = false;

        postFixup = ''
          addDriverRunpath $out/bin/*
        '';
      };

    #################### GIT ############################
    gitVersion = "git"; # Set the version for the Git build

    GitFrontend = buildFrontend {
      src = ./frontend;
      npmDepsHash = "sha256-wD5xGMgpdgVfTMmoIM5ed1Gg4eOkQB7Uo8kA9+ptSKY=";
      version = gitVersion;
    };
    GitBuildAmumax = with pkgs.lib.fileset;
      buildAmumax {
        src = toSource {
          root = ./.;
          fileset = unions [./src ./go.mod ./go.sum ./main.go];
        };
        frontend = GitFrontend;
        vendorHash = "sha256-TayaztdGoJ+WiuMumJ4CQf+Bbr86PjkfXLEBEJ2JV28=";
        version = gitVersion;
      };

    #################### RELEASE ########################
    releaseVersion = "2024.11.26"; # Set the version for the Release build

    ReleaseSrc = pkgs.fetchFromGitHub {
      owner = "MathieuMoalic";
      repo = "amumax";
      rev = releaseVersion;
      hash = "sha256-6IcEr4sZRK7/B3JdbM8vSQHWXe21et9WtvqZNMHGwwY=";
    };

    ReleaseFrontend = buildFrontend {
      src = "${ReleaseSrc}/frontend";
      npmDepsHash = "sha256-wD5xGMgpdgVfTMmoIM5ed1Gg4eOkQB7Uo8kA9+ptSKY=";
      version = releaseVersion;
    };

    ReleaseBuildAmumax = buildAmumax {
      src = ReleaseSrc;
      frontend = ReleaseFrontend;
      vendorHash = "sha256-TayaztdGoJ+WiuMumJ4CQf+Bbr86PjkfXLEBEJ2JV28=";
      version = releaseVersion;
    };

    #################### DEVELOPMENT ENVIRONMENT ########################
    devEnv = pkgs.mkShell {
      inherit CGO_CFLAGS CGO_LDFLAGS CGO_CFLAGS_ALLOW;
      buildInputs =
        basepkgs
        ++ [
          pkgs.go
          pkgs.gopls
          pkgs.golangci-lint
          pkgs.gcc11
          pkgs.nodejs_22
          pkgs.nix-prefetch-github
          pkgs.prefetch-npm-deps
          pkgs.nix-prefetch
          pkgs.jq
          pkgs.podman
          pkgs.delve
          pkgs.gomodifytags
          pkgs.websocat
        ];

      LD_LIBRARY_PATH = "${cuda.libcufft}/lib:${cuda.libcurand}/lib:/run/opengl-driver/lib/";

      shellHook = ''
        export PATH="${pkgs.gcc11}/bin:$PATH"
        export GOPATH=$(pwd)/.go/path
        export GOCACHE=$(pwd)/.go/cache
        export GOENV=$(pwd)/.go/env
        export VITE_WS_URL=http://localhost:35367/ws
        mkdir -p $GOPATH $GOCACHE
      '';
    };
  in {
    packages.${system} = {
      default = ReleaseBuildAmumax;
      git = GitBuildAmumax;
    };
    devShell.${system} = devEnv;
  };
}
