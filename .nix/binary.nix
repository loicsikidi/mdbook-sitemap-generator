{ pkgs }:
pkgs.buildGo124Module rec {
  name = "mdbook-sitemap-generator";
  src = pkgs.lib.cleanSource ../.;
  vendorHash = "sha256-WUQW8EDJ7kT2CUZsNtlVUVwwqFRHkpkU6pFmx7/MDGg=";
}
