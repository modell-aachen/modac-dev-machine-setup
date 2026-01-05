{ stdenv, bash }:

stdenv.mkDerivation {
  pname = "modac-dev-provisioner";
  version = "0.3.7";

  # The devbox directory (containing bin/modac-provision, provision-scripts, etc.)
  src = ./.;

  buildInputs = [ bash ];

  installPhase = ''
    mkdir -p "$out/bin"
    cp -r bin/* "$out/bin/"

    mkdir -p "$out/share"
    cp -r share/* "$out/share"
  '';

  meta = {
    description = "Modac Devbox modac-provision entrypoint and scripts for Darwin and Ubuntu";
    platforms = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
  };
}
