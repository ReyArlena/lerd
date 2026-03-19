# Node

## Commands

| Command | Description |
|---|---|
| `lerd isolate:node <version>` | Pin Node version for cwd — writes `.node-version`, runs `fnm install` |
| `lerd node [args...]` | Run node using the project's version via fnm |
| `lerd npm [args...]` | Run npm using the project's version via fnm |
| `lerd npx [args...]` | Run npx using the project's version via fnm |

---

## Version resolution

1. `.nvmrc` in the project root
2. `.node-version` in the project root
3. `package.json` — `engines.node` field
4. Global default in `~/.config/lerd/config.yaml`

To pin a project:

```bash
cd ~/Lerd/my-app
lerd isolate:node 20
# writes .node-version and runs: fnm install 20
```

---

## fnm

Node version management is handled by [fnm](https://github.com/Schniz/fnm), which is bundled and installed automatically. The `node`, `npm`, and `npx` shims in `~/.local/share/lerd/bin/` invoke the correct version via fnm for each project.
