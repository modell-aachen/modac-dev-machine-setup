{
  description = "Modac development machine2 provisioner CLI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages = {
          machine2 = pkgs.buildGoModule {
            pname = "machine2";
            version = "1.0.0";

            src = ./.;

            # Go module vendoring hash
            # Run: nix build .#machine2 2>&1 | grep "got:" to get the actual hash
            vendorHash = "sha256-hocnLCzWN8srQcO3BMNkd2lt0m54Qe7sqAhUxVZlz1k=";

            # Install provision scripts and templates alongside binary
            postInstall = ''
              mkdir -p $out/share/machine2
              cp -r scripts/provision $out/share/machine2/provision-scripts
              cp -r scripts/templates $out/share/machine2/templates

              mkdir -p $out/share/bash-completion/completions
              $out/bin/machine2 completion bash > $out/share/bash-completion/completions/machine2.bash
            '';

            ldflags = [
              "-s"
              "-w"
              "-X main.version=1.0.0"
              "-X main.commit=${self.rev or "dev"}"
            ];

            meta = with pkgs.lib; {
              description = "Modac development machine2 provisioner CLI";
              homepage = "https://github.com/modell-aachen/modac-dev-machine-setup";
              license = licenses.mit;
              platforms = platforms.unix;
              mainProgram = "machine2";
            };
          };

          default = self.packages.${system}.machine2;
        };

        # Development shell with Go tooling
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_22
            gopls
            gotools
            golangci-lint
            jq
          ];

          shellHook = ''
            echo "Machine CLI development environment"
            echo "Go version: $(go version)"
            echo ""
            echo "Available commands:"
            echo "  go build -o bin/machine2 ./cmd/machine2  - Build binary"
            echo "  go test ./...                           - Run tests"
            echo "  go run ./cmd/machine2                    - Run CLI"
            echo "  golangci-lint run                       - Lint code"
            echo ""
            echo "Nix commands:"
            echo "  nix build .#machine2                     - Build with Nix"
            echo "  nix run .#machine2 -- --help             - Run with Nix"
          '';
        };

        # CLI app for running with nix run
        apps.default = {
          type = "app";
          program = "${self.packages.${system}.machine2}/bin/machine2";
        };
      }
    );
}
