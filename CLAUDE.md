# modac-dev-machine-setup

Development machine provisioner for the Modell-Aachen team. Automates setup of complete dev environments using a Go CLI tool (`machine`), devbox (Nix-based), and 1Password for secrets.

## Architecture

```
modac-dev-machine-setup/
├── devbox/plugins/modac/plugin.json   # Devbox plugin (packages, env vars, shell hooks)
├── nixpkgs/modac-dev-machine/         # Go CLI source code
│   ├── cmd/machine/                   # CLI commands (Cobra)
│   ├── internal/
│   │   ├── provision/                 # 20+ provisioning modules
│   │   ├── config/                    # Config handling (devbox.json)
│   │   ├── output/                    # Logging & terminal formatting
│   │   ├── platform/                  # OS detection (Darwin/Ubuntu)
│   │   ├── backup/                    # 1Password backup integration
│   │   └── util/                      # Filesystem utilities
│   ├── scripts/templates/             # Template files (devbox.json, etc.)
│   ├── scripts/bash/                  # Shell helper scripts
│   ├── flake.nix                      # Nix build definition
│   └── go.mod
├── install                            # Bootstrap script
└── renovate.json
```

## Key Concepts

### Provisioning Flow
`machine provision` runs modules sequentially via `internal/provision/executor.go`. Each module is a Go package with a `Run(out, platform)` function. Modules can be filtered with `-f MODULE`.

Module order: devbox-update → onepassword → restore-backup → packages → **setup-envs** → asdf-packages → asdf → kubectl-krew → setup-k8s-cluster → node → certificates → setup-dev → completions → claude → github-auth-login → install-modac-shell-helper → orbstack → docker-packages → docker

### Secrets Management (setup-envs)
- `op_secrets_tpl` in devbox.json maps env var names to `op://` 1Password references
- `setup-envs` module generates `~/.secrets/env.tpl` from these mappings
- Runs `op inject` to resolve references into `~/.secrets/.env`
- devbox loads this via `env_from` field in global devbox.json
- Currently only handles the **global** devbox config, not per-project configs

### Devbox Plugin (`devbox/plugins/modac/plugin.json`)
- Defines ~25 packages (kubectl, helm, nodejs, gh, claude-code, etc.)
- Sets env vars (REPOS_DIRECTORY, QWIKI_API_HOST, PATH additions, etc.)
- Shell init hooks: devbox completions, direnv, machine aliases, custom completions
- Referenced by global devbox.json via GitHub flake URL with pinned rev hash

### Config Paths
- Global devbox config: `~/.local/share/devbox/global/default/devbox.json`
- Secrets template: `~/.secrets/env.tpl`
- Resolved secrets: `~/.secrets/.env`
- Provision logs: `~/.machine/logs/provision-YYYYMMDD-HHMMSS.log`
- Repos: `~/qwiki-repos/`

### Template System
Templates are loaded from `share/machine/templates/` (relative to binary) with fallback to `scripts/templates/` (repo). Key template: `devbox.json` defines default `op_secrets_tpl` entries.

## Build & Test

```bash
# Build via nix flake
nix build .#machine -L

# Run Go tests
cd nixpkgs/modac-dev-machine && go test ./...

# Lint
cd nixpkgs/modac-dev-machine && go vet ./...
```

## Commands

```bash
machine provision              # Full provisioning
machine provision -f MODULE    # Run specific module
machine provision list-modules # List modules
machine edit-config            # Edit global devbox.json
machine backup create          # Backup config to 1Password
machine backup restore         # Restore from 1Password
machine aliases                # Print shell aliases
```

## Plugin Version Update
When updating `plugin.json`, the rev hash in the global devbox.json package reference must be updated to match the latest commit. Use `machine update-machine-hash` or the `/update-machine-hash` skill.

## Platform Support
- **macOS (Darwin)**: Homebrew for system packages, Orbstack for containers
- **Ubuntu (Linux)**: apt for system packages, supports Distrobox environments
- Detection via `internal/platform/detect.go`
