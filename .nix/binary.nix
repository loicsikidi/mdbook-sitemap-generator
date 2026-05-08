{pkgs}:
pkgs.buildGoModule {
  name = "mdbook-sitemap-generator";
  src = pkgs.lib.cleanSource ../.;
  vendorHash = "sha256-5uDi/9YlUNRDQKFROlN8sLvLRFFOvZvZHcmR6fARG5Q=";
}
