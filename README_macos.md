# MacOS required packages
The machine provisioner requires some basic packages.

* python >= 3.10
* pip3 >= 20.3

For MacOS these packages can be installed using xcode-select,  eg. run `xcode-select --install` to launch the installer.

Afterwards you can continue with [Usage/Inventory creation](README.md) 

Important: In order to have nix-installed packages in the PATH (e.g. devspace), you have to put
`eval "$(devbox global shellenv)`
in your .zshrc