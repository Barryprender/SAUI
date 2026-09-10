# ADR-0001: Render markup with templ rather than html/template

## Status
**Accepted** — 10 September 2026.

The project charter names the choice as open ("`html/template` (standard
library) or `templ` if type-safe templates are warranted"), and the code has
long since settled it. The reasoning lived only in the shape of the source, and
the decision runs against the charter's own preference for the standard library,
so it needs a record rather than an inference.

## Context

Every page on this site is an argument that the server owns state and the
browser renders it. That argument fails visibly if a page renders wrongly, and
the site is dense with the kind of markup that renders wrongly: six demo domains
in `demo/`, each with its own projections, all rendered through shared partials
in `templates/partials/`, and every string on the site duplicated in English and
Spanish through `locale`.

With `html/template` the failure mode is a runtime one. A partial that expects a
`FoodOrder` and receives a `BankingAccount` produces a blank region or a 500 on
the request that hits it, not a compile error. A renamed field silently renders
empty. On a site whose entire claim is that the server is authoritative, a
silently empty region is the worst available outcome: it looks like working
software.

The counter-pressure is real. The charter commits to zero external frontend
dependencies without justification and to no build step, and `templ` is both a
dependency and a code generation step.

## Decision

Use `github.com/a-h/templ` for all markup. Generated `_templ.go` files are
committed alongside their `.templ` sources.

Because generated code that lags its source is exactly the silent failure this
decision was taken to eliminate, `./verify.sh` runs `templ generate` and fails
the build if any `_templ.go` file changes. The generator version is pinned in
`verify.sh` so the gate's own output cannot move without a commit.

## Alternatives considered

**`html/template` from the standard library.** Rejected because its type errors
surface as blank output at request time. The site's density of near-identical
domain shapes across six demos makes passing the wrong struct to the right
partial an ordinary mistake, not an exotic one, and the standard library gives
no signal when it happens.

**`html/template` with hand-written wrapper functions per partial.** This
recovers compile-time checking by writing, by hand, what `templ` generates. It
trades one dependency for an open-ended volume of boilerplate that must be kept
in step with the templates by discipline alone. Discipline is what the gate
exists to replace.

**Rendering markup from Go directly, with no template layer.** Rejected because
the markup is the deliverable here. The site is read by people who will view
source, and WCAG 2.2 AA compliance is easier to hold in a template than in
string concatenation.

## Consequences

**Negative.** The charter's "no build step" claim is now qualified: contributors
must install `templ` at the pinned version and regenerate before committing.
A stale `_templ.go` file is a new class of mistake — one the gate catches, but
only after the fact, and only for contributors who run the gate before pushing.
The dependency is single-maintainer, and a `templ` release that changes
generated output requires a coordinated version bump in `verify.sh` and a full
regeneration commit. Diff noise on `_templ.go` files makes review of markup
changes harder than it would be with plain templates.

**Positive.** Passing the wrong type to a partial fails at compile time, before
CI and long before a reader sees a blank region. Escaping is applied by the
generator at every interpolation, which removes stored cross-site scripting
through the feedback form as a class rather than as a case. The bilingual
`locale` calls type-check, so a missing translation is a build failure instead
of an English string on a Spanish page. Rendering is compiled Go, so no template
parsing happens at request time.

## Follow-up

The decision is reversible if all of the following become true:

1. `html/template` gains compile-time type checking of template data, or a
   standard-library equivalent ships.
2. The generated-code freshness gate is proven to have failed to catch a stale
   `_templ.go` file that reached production.
3. The `templ` project stops publishing security fixes within the windows stated
   in `SECURITY.md`.
