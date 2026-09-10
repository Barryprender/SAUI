# Technical documentation

Prepared to satisfy Annex VII of Regulation (EU) 2024/2847 (the Cyber
Resilience Act). It is updated on every material release; a change to the
dependency set, the data held, or the risk assessment is material.

- **Product:** SAUI — Server-Authoritative UI
- **Version:** tracked by the deployed git commit; see `sbom.json`
- **Manufacturer:** Barry Prendergast
- **Document revised:** 10 September 2026

## 1. Product description

A public website that documents the Server-Authoritative UI pattern and
demonstrates it with interactive demos across six domains (food ordering,
banking, healthcare, SaaS, micro-frontends, distributed systems). The site is
built with the pattern it describes, so the deployed artifact is also the
reference implementation.

It ships as a single statically linked Linux binary plus a `static/` directory,
in an Alpine container, deployed to Fly.io. There is no installer, no client
agent and no update mechanism on the user's machine: the browser is the only
client, and it holds no application state.

## 2. Intended purpose and reasonably foreseeable misuse

**Intended purpose.** To be read by web developers and architects evaluating the
pattern, and to let them exercise the demos and submit written feedback.

**Not intended for.** Storing anything a visitor would mind losing or having
read. The demos reset, the database is not backed up for visitor benefit, and
there is no account to recover.

**Foreseeable misuse.** Visitors may paste personal data or credentials into the
feedback box. The field is capped at 2000 characters and its contents are never
executed or echoed to other visitors, but the text is stored as written. This is
the residual risk in section 4.

## 3. Architecture and data

| Layer | Choice | Note |
| --- | --- | --- |
| Language | Go, standard library first | No web framework |
| HTTP | `net/http` | Go 1.22+ routing patterns |
| Database | `modernc.org/sqlite` | Pure Go, no cgo, single file |
| Templates | `templ` | Compiled, type-checked — see [adr/0001-use-templ-for-templates.md](adr/0001-use-templ-for-templates.md) |
| Client | htmx, progressive enhancement | Reading works with JavaScript disabled |

Persistent data is one SQLite file with an append-only `events` table keyed by an
anonymous session identifier, plus per-demo projection tables rebuilt from it.
The full inventory of personal data held is in
[SECURITY.md](../SECURITY.md#what-this-product-is-and-what-it-holds).

## 4. Cybersecurity risk assessment

Assessed against the OWASP Top 10 (2021). Risks are listed with the control that
addresses them and the residual exposure that remains.

| Risk | Control | Residual |
| --- | --- | --- |
| Broken access control — one session reading another's events | Every query is scoped by the server-held session identifier; the client never supplies it | Thin test coverage of the demo handlers (section 6) |
| Injection — SQL | All statements are parameterised; no string-built SQL | None known |
| Injection — cross-site scripting | `templ` escapes interpolated values at compile time; feedback text is never rendered as HTML | None known |
| Cross-site request forgery | Double-submit token, `SameSite=Strict` on both cookies | None known |
| Session hijacking | `HttpOnly`, `Secure` in production, `SameSite=Strict`, 24-hour expiry | Session fixation is not separately mitigated; the session carries no privileges |
| Vulnerable components | `govulncheck` gate on every push; Go toolchain floor pinned at 1.25.14 across `go.mod`, CI and the Dockerfile | Database lag between disclosure and publication |
| Security misconfiguration | Security-header middleware; container runs as a non-root user; no debug endpoints | None known |
| Denial of service | Rate-limiting middleware; 2000-character cap on feedback | Single instance, no autoscaling — availability is best-effort and stated as such |
| Logging failures | Structured JSON logs of requests and errors; message bodies are not logged | Logs are not shipped off-host or retained centrally |
| Server-side request forgery | The server makes no outbound requests on behalf of a visitor | Not applicable |

The product processes no special-category personal data and performs no
automated decision-making.

## 5. Standards and controls applied

- OWASP Top 10 (2021) — applied as the risk framework in section 4
- WCAG 2.2 AA — the accessibility target for all markup
- CycloneDX 1.6 — SBOM format, `sbom.json` at the repository root
- Regulation (EU) 2024/2847 Annex I Part I — secure-by-default configuration, no
  default credentials, no known exploitable vulnerabilities at release
- Regulation (EU) 2024/2847 Annex I Part II — coordinated disclosure policy and
  SBOM, see [SECURITY.md](../SECURITY.md)

## 6. Verification

`./verify.sh` is the single definition of a passing build; CI runs that file
rather than a copy of it. It fails, never warns, on:

1. Unformatted Go source
2. `go vet` findings
3. Generated `_templ.go` files behind their `.templ` sources — this repository's
   silent lie, because stale generated code lets the whole suite pass against
   markup that is no longer in the source
4. Any failing test, run with `-count=1`
5. Any `govulncheck` finding
6. A `sbom.json` that no longer matches a regeneration under the shipped build
   constraints (`CGO_ENABLED=0 GOOS=linux GOARCH=amd64`)

**Known weakness.** Statement coverage across the application packages is 18.3%.
The state store, the feedback action and the projection layer are tested; the
demo handlers and the middleware are largely not. The pipeline is trustworthy
about what it checks, and what it checks is narrower than the product.

## 7. Software bill of materials

`sbom.json` at the repository root, CycloneDX 1.6, generated by
`cyclonedx-gomod` v1.9.0 under the constraints of the shipped container.
Regenerating it is a gate, so it cannot silently drift: see section 6.

Direct dependencies are `github.com/a-h/templ` and `modernc.org/sqlite`.
Everything else in the SBOM is transitive.

## 8. Support lifecycle

- **Support start:** 10 September 2026
- **Support end:** 10 September 2031

Five years, the Cyber Resilience Act minimum. Security fixes are provided for
the deployed version throughout that period, within the windows stated in
SECURITY.md. If the site is retired earlier, the end date is revised in this
document and announced on the site before the change takes effect — sunsetting
is a decision to record, not something to let lapse.
