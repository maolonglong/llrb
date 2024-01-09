{
  description = "llrb - provides a simple, ordered, in-memory data structure for Go programs";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    nur.url = "github:nix-community/NUR";
    maolonglong-nur.url = "github:maolonglong/nur-packages";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    nur,
    maolonglong-nur,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        overlays = [
          (final: prev: {
            nur = import nur {
              nurpkgs = prev;
              pkgs = prev;
              repoOverrides = {
                maolonglong = import maolonglong-nur {pkgs = prev;};
              };
            };
          })
        ];
        pkgs = import nixpkgs {inherit system overlays;};
      in {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs =
            (with pkgs; [
              just
              go
              golines
              gosimports
            ])
            ++ (with pkgs.nur.repos.maolonglong; [
              gofumpt
              skywalking-eyes
            ]);
        };
      }
    );
}
