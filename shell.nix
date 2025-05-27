let
  # golang pinned to 1.24.0
  nixpkgs = fetchTarball
    "https://github.com/NixOS/nixpkgs/archive/2d068ae5c6516b2d04562de50a58c682540de9bf.tar.gz";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };
  pre-commit = import ./default.nix { };
in pkgs.mkShellNoCC {
  packages = with pkgs; [
    go # v1.24.0
    delve
  ];
  shellHook = ''
    ${pre-commit.pre-commit-check.shellHook}
  '';
  buildInputs = pre-commit.pre-commit-check.enabledPackages;
}
