# SAUI — Server-Authoritative UI
## Claude Code Configuration & Project Reference

---

## What This Project Is

A website built to document, demonstrate, and build consensus around **Server-Authoritative UI (SAUI)** — an architectural approach to web development where the server is the exclusive owner of application state and the frontend is a stateless display layer.

The site is itself built using the SAUI stack. It demonstrates by existing.

---

## The Core Proposition

State belongs where it can be secured, validated, and made authoritative: the server.

The browser renders truth. It does not hold it.

A user interface that displays server state is always correct by definition. A user interface that maintains its own copy of state is always eventually wrong.

---

## Why This Exists

The SPA era solved real problems in 2013. Gmail and Google Maps genuinely required client-side state. The mistake was applying that architecture universally — to food ordering apps, banking dashboards, booking systems, CMS tools — products that share none of the requirements that justified the SPA pattern in the first place.

The cost is now visible:

- Build toolchains of compounding complexity
- Entire engineering disciplines (state management, cache invalidation, optimistic updates, hydration) exist solely to synchronise two copies of the same state
- Security incidents caused by client-side state manipulation
- Test suites that test synchronisation rather than business logic
- Frontend teams blocked on backend contracts because they simulate server logic locally

SAUI is not a new idea. It is the web's original model, re-expressed with modern Go tooling that removes the historical limitations that made SPAs appealing.

---

## Philosophy

- The server owns state. Always. Without exception.
- The frontend is a projection renderer. It receives truth and displays it.
- The frontend submits intents (actions). It never asserts state.
- Security is structural, not bolted on. If the client never holds state, it cannot tamper with it.
- Complexity belongs where it can be tested: the backend.
- Frontend richness (animations, transitions, micro-interactions) is a CSS and interaction design concern. It has no dependency on who owns state.
- Progressive enhancement is the baseline. The UI works without JavaScript for reading.
- Zero external frontend dependencies unless explicitly justified.
- Boring, proven infrastructure over clever, novel infrastructure.

---

## Tech Stack

### Backend
- **Language:** Go — standard library first, no frameworks
- **Database:** SQLite via `modernc.org/sqlite` (pure Go, no cgo, single binary)
- **Templates:** `html/template` (standard library) or `templ` if type-safe templates are warranted
- **HTTP:** `net/http` standard library only
- **Real-time:** SSE (Server-Sent Events) via standard `net/http` — no WebSocket unless justified
- **State persistence:** Append-only event log in SQLite
- **In-memory state:** `sync.RWMutex`-guarded map for ephemeral session state
- **Deployment:** Single binary, Fly.io

### Frontend
- **Markup:** Semantic HTML5
- **Interaction:** htmx for partial DOM updates — no custom JS framework
- **Styling:** Vanilla CSS, CSS custom properties for theming
- **JavaScript:** Vanilla only, no build step, no bundler, no transpilation
- **Progressive enhancement:** All content readable without JavaScript

### Infrastructure
- **Hosting:** Fly.io (single binary deployment, SQLite volume mount)
- **TLS:** Fly.io managed
- **No CDN dependency for functionality** — CDN is acceptable for static assets only

---

## Project Structure

```
saui-web/
├── main.go
├── statestore/
│   ├── store.go          — StateStore struct, sole state authority
│   ├── ephemeral.go      — in-memory session state, RWMutex guarded
│   ├── durable.go        — append-only event log, SQLite
│   ├── projection.go     — derives read models from event log
│   └── gateway.go        — all external state access enforced here
├── actions/
│   ├── action.go         — Action interface: Validate() + Apply()
│   └── [domain].go       — one file per state transition type
├── projections/
│   └── [name].go         — named read models per frontend slice
├── handlers/
│   ├── pages.go          — full page renders
│   └── partials.go       — htmx partial renders
├── middleware/
│   ├── authn.go
│   ├── authz.go
│   └── ratelimit.go
├── static/
│   ├── css/
│   └── js/
└── templates/
    ├── layout/
    ├── pages/
    └── partials/
```

---

## State Architecture

### Two-Layer Model

**Durable state** — what actually happened
- Append-only event log in SQLite
- Never mutated, only appended
- Source of truth that survives restart
- Auditable by definition

**Ephemeral state** — current session context
- In-memory only, keyed by session ID
- UI state: current page context, in-progress flows, unsaved drafts
- Derived from durable state on session creation
- Never used as input to durable state transitions
- Discarded on session end

### The Gateway Contract

All state access — read or write — goes through `statestore/gateway.go`.

No handler, no middleware, no other package reads or writes state directly.

```go
// Reading state — returns a projection
gateway.Project(ctx, sessionID, projectionName) (Projection, error)

// Writing state — submits a validated intent
gateway.Dispatch(ctx, sessionID, action) error
```

The gateway enforces:
- Authentication context is valid
- Actor has permission for this action or projection
- Action validates against current state before application
- Every state transition is logged

### Action Pattern

State changes are validated intents, not direct mutations.

```go
type Action interface {
    Validate(current State) error
    Apply(current State) (State, Event)
}
```

The frontend submits an action. The gateway validates it against current state. If valid: applies it, appends the event to the durable log, updates in-memory projection, returns updated projection. If invalid: rejects with current truth. No partial application. No silent failure.

---

## Security Model

Security is structural in SAUI, not additive.

- **No state tampering** — the client submits intents, never state. There is no state object on the client to tamper with.
- **No privilege escalation** — the server computes what an actor is allowed to see. The classic `isAdmin: true` client-side attack is structurally impossible.
- **No stale reads** — projections are computed from current state on request. There is no client cache to be stale.
- **Audit trail is free** — every state transition is a server-side event. Logged once at the gateway.
- **Session compromise is bounded** — a stolen session sees only what the server projects for it.
- **Input validation is singular** — validated once at the gateway. No client-side validation path to maintain in parallel.

Middleware stack (applied in order):
1. TLS (Fly.io managed)
2. Rate limiting (per route, per session)
3. Authentication (session cookie, HttpOnly, Secure, SameSite=Strict)
4. Authorisation (per action, per projection)
5. Request validation (HX-Request header verification for htmx routes)
6. CSRF (double-submit cookie pattern for mutations)

---

## Frontend Constraints

These are hard rules for this codebase.

- No JavaScript framework. No React, Vue, Svelte, Angular.
- No npm. No node_modules. No build step. No bundler.
- No client-side routing.
- No client-side state management.
- No localStorage or sessionStorage for meaningful state.
- htmx loaded from a pinned version. Hash-verified.
- All interactivity via htmx attributes or minimal vanilla JS.
- CSS only for animations and transitions — no JS animation libraries.
- The site must be fully readable with JavaScript disabled.

**JavaScript language target:** Baseline 2023. ES modules via `type="module"`. No polyfills. No transpilation. No ES5 target. Compatibility matrix is modern evergreen browsers — the audience is developers.

**TypeScript:** Not used directly. JSDoc type annotations used in JS files for type safety and IDE inference without a compilation step. TypeScript is not introduced unless a separate pre-compiled tools directory is explicitly warranted — output committed as plain JS, never in the live build path.

**View Transitions API:** Used as progressive enhancement only where natively supported. No polyfill. Feature-detected at runtime:
```js
if (document.startViewTransition) {
    htmx.config.globalViewTransitions = true;
}
```

Rationale: the site must demonstrate the stack it advocates. A Next.js front end would be a self-refutation.

---

## Progressive Enhancement — Degradation Hierarchy

Progressive enhancement is not an afterthought. It is tested explicitly at each level. All four levels must produce correct state — the server owns truth regardless of client capability.

```
Level 0 — No JS, no CSS
  Server renders complete HTML. All content readable and navigable.
  Forms submit via native POST. Page reloads on state-changing action.
  PRG (Post-Redirect-Get) pattern enforced on all mutation handlers.
  This level must always work. It is the non-negotiable foundation.

Level 1 — CSS only, no JS
  Full layout, typography, colour, spacing applied.
  Visually complete for all static content.
  No visual states that depend on JS class toggling to be meaningful.
  CSS handles :hover, :focus, :disabled — no JS required.

Level 2 — JS available, htmx not yet loaded
  Native form submission active (Level 0 fallback holds).
  Full page responses served — server detects absence of HX-Request header.
  No partial updates. Correct state delivered via full page render.

Level 3 — htmx loaded
  Partial DOM updates active.
  Form submissions intercepted, partials swapped in place.
  hx-boost active on navigation — no full reload for page transitions.
  Server returns partial on HX-Request: true, full page otherwise.

Level 4 — Full baseline JS active
  SSE connection established for real-time projection updates.
  View Transitions API active if natively supported.
  JSDoc-typed enhancement modules active.
  Full experience delivered.
```

**Handler pattern — every mutation handler:**
```go
func (h *Handler) HandleAction(w http.ResponseWriter, r *http.Request) {
    // dispatch regardless of request type — state change always happens
    err := h.gateway.Dispatch(r.Context(), sessionID, action)
    if err != nil { /* handle */ }

    if r.Header.Get("HX-Request") == "true" {
        // return partial projection
        h.templates.ExecuteTemplate(w, "partial.html", projection)
        return
    }
    // PRG — prevents double submission on reload
    http.Redirect(w, r, "/confirmed-destination", http.StatusSeeOther)
}
```

**Form pattern — every form:**
```html
<!-- native action always set — htmx enhances, never replaces -->
<form method="POST" action="/actions/[domain]/[verb]"
      hx-post="/actions/[domain]/[verb]"
      hx-target="#[target-id]"
      hx-swap="outerHTML">
```

**CSS state pattern — no JS dependency for visibility:**
```css
/* wrong — content invisible until JS adds class */
.panel { display: none; }
.panel.is-open { display: block; }

/* right — content visible by default, JS enhances */
.panel { display: block; }
.panel[aria-expanded="false"] { display: none; }
```

**Native HTML elements preferred for interactive patterns:**
- `<details>` / `<summary>` for disclosure — no JS required
- `<dialog>` for modals — native open/close via `form[method="dialog"]`
- `<progress>` for progress indication
- Form validation via `required`, `pattern`, `type` attributes first

**SSE degradation:**
```js
if (typeof EventSource !== 'undefined') {
    // connect SSE for real-time projection updates
} else {
    // poll every 30s as fallback — correct state on each poll
}
```

---

## Content Structure

```
/                   — the proposition, stated concisely
/why                — diagnosis of SPA complexity cost
/architecture       — SAUI pattern explained with diagrams
/stack              — the Go implementation in detail
/cases              — real application walkthroughs
  /cases/food-ordering
  /cases/banking
  /cases/healthcare
  /cases/saas-dashboard
  /cases/distributed-systems
  /cases/micro-frontends
/testing            — how SAUI collapses frontend test surface
/limits             — where SAUI is the wrong choice (honest)
/blog               — ongoing posts, counterarguments, responses
/code               — reference implementation, GitHub links
```

---

## Case Study Scope

Each case study covers:
1. Current client-side state model and the problems it causes
2. What moves to the server under SAUI
3. What stays on the client (legitimate UI-only state)
4. Honest degradation cost if any
5. Security properties gained
6. Testing surface change

### Confirmed case studies
- Food ordering (cart integrity, price manipulation, availability staleness)
- Banking (balance accuracy, transfer state durability, audit requirements)
- Healthcare booking (slot locking, double booking elimination)
- SaaS dashboards (filter state, bulk action correctness, unsaved edit durability)
- Distributed systems / microservices (BFF pattern, per-service authority, no shared state bus)
- Micro frontends (team isolation via stateless MFEs, elimination of cross-MFE state contracts)

---

## Distributed Systems Position

SAUI scales to distributed systems. "Centralised" means authoritative per service domain — not a single global database.

```
Frontend (stateless display)
    ↓
BFF — Backend for Frontend (Go, owns session + ephemeral state, aggregates projections)
    ↓
Order Service | Inventory Service | Payment Service | User Service
(each owns its domain state authoritatively)
```

The frontend never talks to multiple services. The BFF is its single source of projected truth. Distribution complexity lives behind the BFF — the frontend is unaware of it.

Micro frontends under SAUI are genuinely independent because there is no shared client state to couple them. Each MFE consumes its own projection. No shared state bus. No cross-MFE event contracts. No integration test surface between MFEs.

---

## Testing Position

Under SPA architecture, the majority of the test surface exists to verify that client state correctly reflects server state. This is testing the synchronisation problem.

Under SAUI, that synchronisation problem does not exist.

**Frontend tests collapse to:** does this HTML render correctly given this projection? Stable, fast, not coupled to business logic.

**Backend tests cover everything meaningful:** action validation, state transitions, authorisation, projection correctness. Pure Go. Fast. Deterministic. No browser required.

**Team isolation consequence:** backend teams change business logic without breaking frontend tests. Frontend teams redesign the UI without breaking backend tests. The seam is the projection contract. Version it. The teams are independent.

---

## Honest Limits

SAUI is the wrong choice for:

- **Real-time collaborative editing** — concurrent multi-user mutation of the same document requires conflict resolution (CRDTs, operational transforms) that SAUI does not provide.
- **Canvas / creative tools** — local state is inherent to the interaction model.
- **Offline-first PWAs** — client-side state with sync-on-reconnect is a structural requirement. SAUI is online-first.
- **Games and simulations** — high-frequency local state is a performance requirement.
- **High-frequency financial data** (live ticker, order book) — SSE/WebSocket projection streams work, but at scale this becomes an infrastructure problem requiring careful design.

These are genuine exceptions, not contrived ones. The site states them plainly.

---

## Writing Tone

- Practitioner voice, not evangelist voice
- Architectural argument, not framework marketing
- Honest about limits — credibility comes from acknowledging where the approach fails
- No tribal language ("React is dead", "SPAs are a mistake")
- The diagnosis is structural, not personal
- Invite disagreement — the goal is consensus, not conversion

---

## Relation to Secure-UI

Secure-UI solves the client/server trust boundary — who is making the request and can it be trusted.

SAUI solves the state ownership question — where does truth live and who can change it.

They are complementary. Secure-UI is the appropriate security layer for a SAUI application. They share the same server-first, zero-dependency philosophy and were designed with the same threat model in mind.

---

## Claude Code Rules

- Go standard library first. No frameworks. No ORM.
- `modernc.org/sqlite` only for SQLite. No cgo.
- All state access through `statestore/gateway.go` exclusively.
- No handler imports `statestore` directly — only via injected `*Gateway`.
- Errors handled explicitly. No panic outside main.
- `context.Context` threaded correctly through all calls.
- No globals except logger and config at startup.
- HTML templates use `html/template` auto-escaping. No raw HTML injection.
- No `hx-*` attributes in user-generated content — strip at the gateway response layer.
- CSS custom properties for all theme values. No inline styles.
- No JavaScript build step. No npm. No bundler. JS files served directly.
- JSDoc type annotations for type safety — no TypeScript compilation in the live path.
- JS targets Baseline 2023. ES modules via type="module". No polyfills. No transpilation.
- All four degradation levels tested explicitly before any feature is considered complete.
- All htmx routes validate `HX-Request: true` header in middleware.
- Projection endpoints are GET only. Action endpoints are POST only.
- SQLite WAL mode enabled. Single writer enforced.
- Fly.io deployment: single binary, volume-mounted SQLite, no external dependencies.

---

## Case Study Demo Package Pattern

Each case study has a live interactive demo. All six demos follow the same structure.

### Directory layout — one package per demo

```
demo/<domain>/
  events.go        — event type constants + event-log replay helpers (unexported)
  actions.go       — action types implementing statestore.Action
  projections.go   — projection types + init() registration with statestore
  helpers.go       — shared pure functions (e.g. formatPrice)
  layout.templ     — standalone app layout (no SAUI docs nav)
  *.templ          — page and partial templates
  *_templ.go       — generated; do not edit by hand

statestore/<domain>.go   — DB migration + Store methods for domain reference data
handlers/<domain>_demo.go — HTTP handlers (one import: saui/demo/<domain>)
static/demo/<domain>/    — CSS and assets
```

### Why one package per demo (not actions/<domain>/ + projections/<domain>/)

Two sub-packages both named after the domain would both declare `package <domain>`,
forcing import aliases everywhere. Co-locating everything in `demo/<domain>/` gives
a single clean import, zero aliases, and keeps domain code self-contained.

### Template co-location

Templates live inside `demo/<domain>/`, not in `templates/demo/<domain>/`.
Because they are in the same package as the projection types, they reference
types directly with no import — `MenuProjection`, `CartProjection`, etc.

### Naming conventions inside a demo package

Drop the domain prefix on types — they are already scoped by package:
- `food.AddToCart` not `food.FoodAddToCart`
- `food.MenuProjection` not `food.FoodMenuProjection`
- `food.CartItemAdded` not `food.FoodCartItemAdded`

### Handler wiring

Demo routes use the same `page` middleware (session + CSRF) as doc pages.
Routes follow the pattern `/demo/<domain>/...`

### Adding a new demo

1. `mkdir demo/<domain>/`
2. Write `events.go`, `actions.go`, `projections.go`, `helpers.go`
3. Write `layout.templ` + page/partial templates
4. Run `templ generate ./demo/<domain>/...`
5. Add `statestore/<domain>.go` with migration + Store methods; call migration from `statestore/store.go` `New()`
6. Add `handlers/<domain>_demo.go`
7. Add `static/demo/<domain>/<domain>.css`
8. Wire routes in `main.go`
9. `go build ./...` to verify

### Completed demos

- `demo/food/` — food ordering (`/demo/food-ordering`)

### Planned demos (case studies)

- `demo/banking/` — `/demo/banking`
- `demo/healthcare/` — `/demo/healthcare`
- `demo/saas/` — `/demo/saas-dashboard`
- `demo/distributed/` — `/demo/distributed-systems`
- `demo/mfe/` — `/demo/micro-frontends`

---

## Non-Negotiable Constraints

1. The site itself must run on the SAUI stack. No exceptions.
2. No client-side state for anything with business meaning.
3. The `/limits` page must exist and be honest.
4. The reference implementation must be open source and linkable.
5. No venture capital positioning. No SaaS. No monetisation of the idea itself.
6. The goal is consensus and adoption, not a product.
