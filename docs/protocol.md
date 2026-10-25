# Protocol: x360-dashboard

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

`X360_DASHBOARD_SIGNING_KEY` is a mounted, uncommitted RSA-2048/65537 private key.
Home payloads use the native YTGR header, RSA PKCS#1 v1.5 SHA-1 header signature, and
independent payload SHA-1. Wharf must pin the matching public modulus via `--xsigner_public_key`.
Key generation and distribution are deployment steps; no private key is in the image or repository.
The matching trust update on real consoles is still pending.
