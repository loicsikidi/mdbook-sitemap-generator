{pkgs}:
pkgs.buildGo124Module rec {
  name = "mdbook-sitemap-generator";
  src = pkgs.lib.cleanSource ../.;
  vendorHash = "sha256-cEgvwog50izBOyMlCdLI2KvwSHPKZsp1wSw6a59V1yw=";
}
