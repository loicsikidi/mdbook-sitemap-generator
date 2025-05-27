{ pkgs }:
pkgs.buildGo124Module rec {
  name = "mdbook-sitemap-generator";
  src = pkgs.lib.cleanSource ../.;
  vendorHash = "sha256-h5UXs7ujP0YIBKusstTDXvWBhqKpYs6HtzzdaB1+6Wg=";
}
