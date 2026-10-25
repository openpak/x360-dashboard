# Xbox 360 Dashboard Server

OpenPak's Xbox 360 Dashboard Server.

| | |
| --- | --- |
| Network | Xbox 360 |
| Role | Public Home content for OpenPak's Xbox 360 network: serves the signed `dashhome.xml` Home manifest |
| Image | `ghcr.io/openpak/x360-dashboard` |
| Licence | AGPL-3.0-only |

## What it serves

| Endpoint | What |
| --- | --- |
| `GET /epix/<locale>/dashhome.xml` | Signed Home manifest (locale varies, e.g. `en-US`) |
| `GET /healthz` | Health check, returns 200 |

## Run

```sh
X360_DASHBOARD_SIGNING_KEY=/path/to/content-key.pem GOWORK=off go run ./cmd/x360-dashboard
```

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
