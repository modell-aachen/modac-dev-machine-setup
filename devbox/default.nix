{ stdenv, bash }:

stdenv.mkDerivation {
  pname = "modac-provision";
  version = "0.3.6";

  # The devbox directory (containing bin/modac-provision, provision-scripts, etc.)
  src = ./.;

  buildInputs = [ bash ];

  installPhase = ''
    mkdir -p "$out/bin"
    cp -r bin/* "$out/bin/"

    mkdir -p "$out/templates"
    cp -r templates/* "$out/templates/"

    mkdir -p "$out/share/bash-completion/completions"
    cp -r completions/* "$out/share/bash-completion/completions"

    mkdir -p "$out/provision-scripts"
    cp -r provision-scripts/* "$out/provision-scripts/"
  '';

  meta = {
    description = "Modac Devbox modac-provision entrypoint and scripts for Darwin and Ubuntu";
    platforms = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
  };
}
