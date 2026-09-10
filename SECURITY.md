# Security policy

## What this product is, and what it holds

SAUI is a public website that documents an architectural pattern and
demonstrates it with interactive demos. It has no user accounts and no
authentication.

What the server stores, in a single SQLite file:

- An anonymous session identifier, issued in an `HttpOnly`, `SameSite=Strict`
  cookie and expiring after 24 hours.
- An append-only event log keyed by that session identifier: the demo
  interactions a visitor performs, and free-text messages submitted through the
  feedback form (2000 characters maximum).

No names, email addresses, passwords or payment data are collected. Visitors may
type personal data into the feedback box; that text is stored as written.

## The worst plausible defect

The demos are the attack surface. The pattern's central claim is that a client
cannot assert state, so a defect that let one session read or mutate another
session's event log would disprove the product's own thesis as well as leak
whatever a visitor typed into the feedback form. Cross-session leakage is
therefore treated as critical regardless of the volume of data involved.

Stored cross-site scripting through the feedback form is the second concern: the
site is a reference implementation, and a payload rendered back to other readers
would be served from an origin they were invited to trust.

## Reporting a vulnerability

Report privately through GitHub's **Report a vulnerability** button on the
Security tab of this repository. That opens a private advisory visible only to
the maintainer.

Do not open a public issue for a suspected vulnerability, and do not test
against the live site in a way that degrades it for other visitors.

## Response windows

These are commitments, not aspirations. The project has one maintainer; the
windows are set to what one person can actually meet.

| Stage | Window |
| --- | --- |
| Acknowledgement of the report | 3 working days |
| First assessment, with a severity and a plan | 10 working days |
| Fix or documented mitigation, high or critical | 30 days from assessment |
| Fix or documented mitigation, everything else | 90 days from assessment |
| Public disclosure | After the fix ships, or 90 days, whichever is first |

Credit is given to the reporter in the advisory unless they ask otherwise.

## Actively exploited vulnerabilities

If a vulnerability in a released version is found to be under active
exploitation, the EU Cyber Resilience Act early-warning clock starts at
**discovery**, not at triage: ENISA must receive an early warning within 24
hours, a fuller notification within 72 hours, and a final report within 14 days
of a corrective measure being available. That obligation applies to this project
from 11 September 2026.

## Supported versions

Only the currently deployed version is supported. There are no release branches;
fixes ship forward. The intended support period is recorded in
[docs/technical-documentation.md](docs/technical-documentation.md).

## What the automated gates do and do not cover

`./verify.sh` runs on every push and pull request. It fails the build on
unformatted code, `go vet` findings, stale generated templates, failing tests, a
`govulncheck` finding, or a stale SBOM.

It does not substitute for a threat model, and it cannot make untested code
correct. Statement coverage across the application packages is currently 18.3%,
which is thin: most of the demo handlers and all of the middleware are exercised
only indirectly, if at all. Treat a green pipeline as evidence that nothing
known-bad is present, not as evidence that the code is right.
