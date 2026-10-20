# x360-dashboard

Public Home content for OpenPak's Xbox 360 network. One service, repository, and image.
The authored `dashhome.xml`, wrapped in a signed native XSigner envelope, references UI layouts, disc handling, and an icon supplied by
stock dashboard 17559. No firmware, extracted manifests, images, or Aurora files are bundled.

LiveHive configuration:

```
EpixXManUriRoot=http://content.xboxlive.com
EpixXManManifestUriPath=/epix/$locale/
```

The dashboard requests `/epix/en-US/dashhome.xml` (locale varies). Both keys are required:
missing values cause 8B010002 before any HTTP request. This public manifest carries no user
identity or bearer token. DNS must send content.xboxlive.com to the primary; merge the
Traefik fragment in deploy/. Wharf uses its local DNS override until real-console DNS exists.

`X360_DASHBOARD_ADDR` defaults to `:10043`. Run `GOWORK=off go test ./...`.
Tag releases to publish ghcr.io/openpak/x360-dashboard for amd64 and arm64.

Status: Home, Games and Settings channels use the stock dashboard renderer. Authored tiles
invoke native My Games, Recent, My Apps, profile, system, preferences and family screens.
No online Social package is supplied yet. XML, HTTP and signature checks pass locally;
this expanded navigation still needs client verification. Installed-title discovery and
launching depend on the console content store and remain separate from opening My Games.
Persistent profile settings, store listings, game catalogs, and personalized Home content
are separate work. Licensed AGPL-3.0; see LICENSE.

`X360_DASHBOARD_SIGNING_KEY` is a mounted, uncommitted RSA-2048/65537 private key.
Home payloads use the native YTGR header, RSA PKCS#1 v1.5 SHA-1 header signature, and
independent payload SHA-1. Wharf must pin the matching public modulus via `--xsigner_public_key`.
Key generation and distribution are deployment steps; no private key is in the image or repository.
The matching trust update on real consoles is still pending.
