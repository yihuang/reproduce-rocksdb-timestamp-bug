{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = self.overlays.default;
          config = { };
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go_1_22
            pkgs.rocksdb
          ];
        };
      }
    ))
    // {
      overlays.default = [
        (final: prev: {
          rocksdb = final.callPackage ./rocksdb.nix { };
        })
      ];
    };
}
