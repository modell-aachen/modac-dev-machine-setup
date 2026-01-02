{
  description = "Modac devbox provision scripts";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f:
        nixpkgs.lib.genAttrs systems (system:
          f {
            inherit system;
            pkgs = import nixpkgs { inherit system; };
          }
        );
    in {
      packages = forAllSystems ({ pkgs, system }: {
        modac-provision = pkgs.callPackage ./default.nix {};
        default = self.packages.${system}.modac-provision;
      });
    };
}
