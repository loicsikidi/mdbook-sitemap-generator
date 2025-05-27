{ pkgs ? import <nixpkgs> { } }: {
  binary = pkgs.callPackage ./.nix/binary.nix { };
  pre-commit-check = pkgs.callPackage ./.nix/pre-commit.nix { };
}
