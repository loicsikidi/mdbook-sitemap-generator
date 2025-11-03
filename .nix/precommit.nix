{pkgs, ...}: let
  nix-pre-commit-hooks = import (
    builtins.fetchGit {
      url = "https://github.com/cachix/git-hooks.nix";
      ref = "refs/heads/master";
      rev = "ca5b894d3e3e151ffc1db040b6ce4dcc75d31c37"; # 2025-17-10
    }
  );
in
  nix-pre-commit-hooks.run {
    src = ./.;
    hooks = {
      # common
      end-of-file-fixer = {
        enable = true;
        package = pkgs.python3Packages.pre-commit-hooks;
      };
      # nix
      alejandra = {
        enable = true;
        package = pkgs.alejandra;
      };
      # golang
      gofmt = {
        enable = true;
        package = pkgs.go;
      };
      govet = {
        enable = true;
        package = pkgs.go;
      };
      golangci-lint = {
        enable = true;
        package = pkgs.golangci-lint;
        extraPackages = [pkgs.go];
        stages = ["pre-push"]; # because it takes a while
      };
      gotest = {
        enable = true;
        package = pkgs.go;
        settings.flags = "-race -failfast -v";
        stages = ["pre-push"]; # because it takes a while
      };
    };
  }
