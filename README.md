# SAUI — Server-Authoritative UI

A website that documents an architectural argument and demonstrates it by being
built with it: **the server owns application state; the browser renders it and
submits intents, but never holds it.**

Live at [saui.fly.dev](https://saui.fly.dev). The reasoning is on the site; this file is
for running and changing the code.

## Stack

Go standard library, `net/http` routing, SQLite through `modernc.org/sqlite`
(pure Go, no cgo, one file), `templ` for compiled type-checked markup, and htmx
for partial updates. No frontend framework, no bundler, no npm.

Two direct dependencies. Everything else in `sbom.json` is transitive.

## Run it

Requires Go 1.25.14 or later — the floor is a security floor, see below.

```sh
go run .
```

Serves on `:8080` against `./saui.db`, created on first run. Override with
`SAUI_ADDR`, `SAUI_DB`, `SAUI_SECURE=true` (sets the `Secure` cookie flag; use
it behind TLS), and `SAUI_TRUSTED_PROXY`.

After editing any `.templ` file:

```sh
templ generate
```

## Verify it

`./verify.sh` is the only definition of a passing build. CI runs that file, not
a copy of it, so a green run on your machine and a green run on GitHub mean the
same thing.

```sh
./verify.sh                 # every gate
./verify.sh --update-sbom   # regenerate sbom.json, then every gate
```

It fails — never warns — on unformatted source, `go vet` findings, `_templ.go`
files behind their `.templ` sources, any failing test, any `govulncheck`
finding, or a `sbom.json` that no longer matches a regeneration under the
shipped build constraints. It installs its own pinned tooling.

Statement coverage is 18.3%. That is thin, and stated rather than hidden: the
state store, the feedback action and the projections are tested; the demo
handlers and the middleware are largely not.

## The Go version floor

The `go` line in `go.mod`, the CI toolchain and the `Dockerfile` builder are all
pinned to 1.25.14. Earlier 1.25 patches ship standard-library vulnerabilities
that `govulncheck` reports as reachable from this code. Raise all three
together, or the container ships a toolchain the gate never inspected.

## Layout

| Path | Contents |
| --- | --- |
| `actions/` | Intents the client may submit; each validates before it applies |
| `statestore/` | SQLite, the append-only event log, and the gateway over it |
| `projections/` | Read models rebuilt from events |
| `handlers/` | HTTP handlers, one per page or partial |
| `middleware/` | Session, CSRF, rate limiting, security headers |
| `templates/` | `.templ` sources and their committed generated `_templ.go` |
| `demo/<domain>/` | One package per demo domain, templates co-located |
| `locale/` | English and Spanish strings |
| `docs/` | Technical documentation and architecture decision records |

## Security

Report vulnerabilities privately through GitHub's **Report a vulnerability**
button on the Security tab. Response windows, the data the product holds, and
the disclosure policy are in [SECURITY.md](SECURITY.md).

Compliance artifacts for Regulation (EU) 2024/2847 are in
[docs/technical-documentation.md](docs/technical-documentation.md) and
`sbom.json` (CycloneDX 1.6).

## Contributing

Read [CLAUDE.md](CLAUDE.md) first — it carries the project's constraints and the
traps that make a green run lie. Run `./verify.sh` before opening a pull
request. Decisions that close off an alternative get a record in
[docs/adr/](docs/adr/).
