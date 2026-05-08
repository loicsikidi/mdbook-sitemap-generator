let
  # golang pinned to 1.25.0
  nixpkgs =
    fetchTarball
    # go to https://www.nixhub.io/packages/go to the list of available versions
    "https://github.com/NixOS/nixpkgs/archive/f4b140d5b253f5e2a1ff4e5506edbf8267724bde.tar.gz";
  pkgs = import nixpkgs {
    config = {};
    overlays = [];
  };
  helpers = import (builtins.fetchTarball
    "https://github.com/loicsikidi/nix-shell-toolbox/tarball/main") {
    inherit pkgs;
    hooksConfig = {
      treefmt.enable = true;
      gofmt.enable = false;
      gotest.settings.flags = "-race";
    };
  };
in
  pkgs.mkShell {
    buildInputs = helpers.packages;

    shellHook = ''
      ${helpers.shellHook}
      echo "Development environment ready!"
      echo "  - Go version: $(go version)"
    '';

    # to enable debugging with delve
    hardeningDisable = ["fortify"];

    env = {
      # Required to run tests with -race flag
      CGO_ENABLED = "1";
    };
  }
