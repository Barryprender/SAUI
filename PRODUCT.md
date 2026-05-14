# Design Context — SAUI

## Target Audience
Experienced web developers and software architects evaluating the SAUI pattern (Server-Authoritative UI). They are opinionated, sceptical, and technically literate. They read code before prose.

## Use Cases
- Learning the architectural argument for server-owned state
- Evaluating whether SAUI fits their stack
- Reading implementation detail (Go, htmx, SQLite, templ)
- Submitting counterarguments via the feedback form

## Brand Personality & Tone
**Authoritative, precise, minimal.** This is an architectural argument, not a product page. The design should feel like a well-typeset technical paper — confident, opinionated, not trying to sell anything. No decorative flourishes that aren't earning their keep. Restraint is the aesthetic.

Secondary quality: **legible at density**. Code blocks, SQL, and Go interfaces appear throughout. The design must handle monospace content without breaking rhythm.

## Design Direction
**Editorial / technical reference.** Light-dominant (dark mode adapts via `prefers-color-scheme`). Serif headings (Lora) against sans body for typographic contrast. The accent colour (oklch blue ~240°) is used sparingly — links, focus rings, active states, callout borders.

## Aesthetic Constraints
- No decorative elements that don't carry meaning
- No gradients on text or as decorative backgrounds
- No glassmorphism, glow effects, or neon
- No rounded-rectangle cards with drop shadows as decoration
- The hero diagram (SVG) and code blocks are the visual centrepieces

## Technical Constraints
- Go `templ` templates — no framework JS, no build step
- htmx for partial updates; progressive enhancement L0–L4 mandatory
- CSS custom properties as design tokens (see below)
- `@layer base, component, utility` — new styles go in `component`
- No inline styles except dynamic values via templ expressions
- WCAG 2.2 AA minimum contrast

## Design Tokens (use these, never hardcode)

### Spacing (8pt grid)
`--sp-1` 0.25rem · `--sp-2` 0.5rem · `--sp-3` 0.75rem · `--sp-4` 1rem
`--sp-5` 1.5rem · `--sp-6` 2rem · `--sp-7` 3rem · `--sp-8` 4rem
`--sp-9` 6rem · `--sp-10` 8rem

### Type scale (fluid)
`--text-xs` · `--text-sm` · `--text-base` · `--text-lg` · `--text-xl`
`--text-2xl` · `--text-3xl` · `--text-4xl`

### Leading
`--leading-tight` 1.15 · `--leading-snug` 1.35 · `--leading-prose` 1.72

### Colour tokens
`--c-bg` · `--c-surface` · `--c-border` · `--c-text` · `--c-text-muted`
`--c-accent` · `--c-accent-hover` · `--c-accent-subtle`
`--c-code-bg` · `--c-code-text`

### Border & radius
`--border` (1px solid var(--c-border)) · `--radius-sm` 3px · `--radius-md` 6px · `--radius-lg` 12px

### Transitions
`--t-fast` 150ms cubic-bezier(0.25,1,0.5,1) · `--t-base` 250ms cubic-bezier(0.16,1,0.3,1)

### Layout
`--max-width` 1400px · `--gutter` clamp(1.25rem, 5vw, 3.5rem) · `--prose-width` 68ch

## Quality Bar
**Flagship.** This site IS the reference implementation. Every UI detail is a claim about the pattern's quality.

## What "Extraordinary" Means Here
Not sensory spectacle — functional precision. A form that validates clearly, a panel that loads without jank, spacing that creates hierarchy without shouting. The extraordinary quality here is the absence of anything unnecessary.
