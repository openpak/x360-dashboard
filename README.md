# x360-dashboard

Public Home content for OpenPak's Xbox 360 network. One service, repository, and image.
The authored `dashhome.xml` references UI layouts, disc handling, and an icon supplied by
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

Status: bootstrap Home with disc and OpenPak tiles. XML and HTTP checks pass locally;
console parsing and rendering must be verified on Wharf before claiming Home works.
Persistent profile settings, store listings, game catalogs, and personalized Home content
are separate work. Licensed AGPL-3.0; see LICENSE.
