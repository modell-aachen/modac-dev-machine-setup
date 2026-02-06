{
  description = "Modac development machine provisioner CLI";

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
          machine = pkgs.buildGoModule {
            pname = "machine";
            version = "1.0.0";

            src = ./.;

            # Go module vendoring hash
            vendorHash = "sha256-7K17JaXFsjf163g5PXCb5ng2gYdotnZ2IDKk8KFjNj0=";

            # Install templates alongside binary (scripts are now in Go)
            postInstall = ''
              mkdir -p $out/share/machine
              cp -r scripts/templates $out/share/machine/templates

              mkdir -p $out/share/bash-completion/completions
              $out/bin/machine completion bash > $out/share/bash-completion/completions/machine.bash

              cp -r scripts/bash/* $out/bin
              cp -r scripts/completions/* $out/share/bash-completion/completions
            '';

            ldflags = [
              "-s"
              "-w"
              "-X main.version=1.0.0"
              "-X main.commit=${self.rev or "dev"}"
            ];

            meta = with pkgs.lib; {
              description = "Modac development machine provisioner CLI";
              homepage = "https://github.com/modell-aachen/modac-dev-machine-setup";
              license = licenses.mit;
              platforms = platforms.unix;
              mainProgram = "machine";
            };
          };

          default = self.packages.${system}.machine;
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
            echo "  go build -o bin/machine ./cmd/machine  - Build binary"
            echo "  go test ./...                           - Run tests"
            echo "  go run ./cmd/machine                    - Run CLI"
            echo "  golangci-lint run                       - Lint code"
            echo ""
            echo "Nix commands:"
            echo "  nix build .#machine                     - Build with Nix"
            echo "  nix run .#machine -- --help             - Run with Nix"
          '';
        };

        # CLI app for running with nix run
        apps.default = {
          type = "app";
          program = "${self.packages.${system}.machine}/bin/machine";
        };
      }
    );
}
