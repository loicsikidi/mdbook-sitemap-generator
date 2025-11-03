let
  # golang pinned to 1.24.0
  nixpkgs =
    fetchTarball
    "https://github.com/NixOS/nixpkgs/archive/2d068ae5c6516b2d04562de50a58c682540de9bf.tar.gz";
  pkgs = import nixpkgs {
    config = {};
    overlays = [];
  };
  pre-commit = pkgs.callPackage ./.nix/precommit.nix {};
in
  pkgs.mkShellNoCC {
    shellHook = ''
      ${pre-commit.shellHook}
    '';
    buildInputs = pre-commit.enabledPackages;

    packages = with pkgs; [
      go # v1.24.0
      delve

      # Required to run tests with -race flag
      gcc # 14.3.0
    ];

    env = {
      # Required to run tests with -race flag
      CGO_ENABLED = "1";
    };
  }
