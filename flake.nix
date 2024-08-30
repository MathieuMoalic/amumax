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

    version = "2024.08.30";

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

    frontend = pkgs.buildNpmPackage {
      inherit version;
      pname = "frontend";
      src = ./frontend;
      npmDepsHash = "sha256-DJOiaPDiWJEkcon/Lc3TD/5cS5v5ArORnpp7HDEpa4E=";

      npmBuild = ''
        npm run build
      '';

      installPhase = ''
        mv dist $out
      '';
    };

    buildAmumax = pkgs.buildGoModule {
      inherit version CGO_CFLAGS CGO_LDFLAGS CGO_CFLAGS_ALLOW;
      pname = "amumax";
      vendorHash = "sha256-ly7mLulUon9XIztddOtP6VEGJZk6A6xa5rK/pYwAP2A=";
      src = ./.;

      buildInputs =
        basepkgs
        ++ [
          pkgs.addDriverRunpath
        ];

      # strip symbols and add version
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
      inherit CGO_CFLAGS CGO_LDFLAGS CGO_CFLAGS_ALLOW;
      buildInputs =
        basepkgs
        ++ [
          pkgs.go
          pkgs.gopls
          pkgs.golangci-lint
          pkgs.gcc11
        ];

      LD_LIBRARY_PATH = "${cuda.libcufft}/lib:${cuda.libcurand}/lib:/run/opengl-driver/lib/";

      shellHook = ''
        export PATH="${pkgs.gcc11}/bin:$PATH"
        export GOPATH=$(pwd)/.go/path
        export GOCACHE=$(pwd)/.go/cache
        mkdir -p $GOPATH $GOCACHE
      '';
    };
  in {
    packages.${system} = {
      default = buildAmumax;
      git = buildAmumax;
    };
    devShell.${system} = devEnv;
  };
}
