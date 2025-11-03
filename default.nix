{pkgs ? import <nixpkgs> {}}: {
  binary = pkgs.callPackage ./.nix/binary.nix {};
}
