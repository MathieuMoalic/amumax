{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = {
    self,
    nixpkgs,
    ...
  }: let
    system = "x86_64-linux";
    pkgs = import nixpkgs {
      inherit system;
      config = {
        allowUnfree = true; # cuda is unfree
      };
    };

    frontend = pkgs.buildNpmPackage {
      pname = "frontend";
      version = "2024.08.29";
      src = ./frontend;
      npmDepsHash = "sha256-DJOiaPDiWJEkcon/Lc3TD/5cS5v5ArORnpp7HDEpa4E=";

      npmBuild = ''
        npm run build
      '';

      installPhase = ''
        mv dist $out
      '';
    };

    cuda = pkgs.cudaPackages_11;
    buildAmumax = pkgs.buildGoModule rec {
      pname = "amumax";
      version = "2024.08.29";
      vendorHash = "sha256-vJKjIjcw+yUzwY43BKkXn2exVQxH6ZHnar7MaJHr9x4=";
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
      CGO_LDFLAGS = ["-L${cuda.cuda_cudart.lib}/lib/stubs/ -lcuda -lcurand -lcufft"];
      CGO_CFLAGS_ALLOW = "(-fno-schedule-insns|-malign-double|-ffast-math)";
      ldflags = [
        "-s"
        "-w"
        "-X github.com/MathieuMoalic/amumax/engine.VERSION=${version}"
      ];
      buildPhase = ''
        cp -r ${frontend} api/static
        go build -v -o $out/bin/amumax .
      '';

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
    packages.${system}.default = buildAmumax;
    devShell.${system} = devEnv;
  };
}
