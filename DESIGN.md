---
name: SAUI
description: Server-Authoritative UI — a reference implementation of server-owned state, built to be read as much as used.
colors:
  warm-paper:      "oklch(97% 0.008 80)"
  warm-surface:    "oklch(93% 0.010 80)"
  warm-border:     "oklch(87% 0.015 80)"
  near-black:      "oklch(14% 0.000 0)"
  muted-text:      "oklch(50% 0.010 80)"
  ink-blue:        "oklch(46% 0.120 240)"
  ink-blue-deep:   "oklch(40% 0.140 240)"
  ink-blue-subtle: "oklch(94% 0.040 240)"
  code-dark:       "oklch(17% 0.020 240)"
  code-light:      "oklch(88% 0.030 200)"
  error-bg:        "oklch(95% 0.04 25)"
  error-border:    "oklch(55% 0.18 25)"
  error-text:      "oklch(35% 0.14 25)"
  statement-void:  "oklch(10% 0.030 245)"
typography:
  display:
    fontFamily: "'Lora', 'Lora-Fallback', Georgia, serif"
    fontSize: "clamp(2.20rem, 7.0vw, 3.80rem)"
    fontWeight: 600
    lineHeight: 1.15
    letterSpacing: "-0.02em"
  headline:
    fontFamily: "'Lora', 'Lora-Fallback', Georgia, serif"
    fontSize: "clamp(1.80rem, 5.0vw, 2.80rem)"
    fontWeight: 600
    lineHeight: 1.15
    letterSpacing: "-0.01em"
  title:
    fontFamily: "'Lora', 'Lora-Fallback', Georgia, serif"
    fontSize: "clamp(1.50rem, 4.0vw, 2.10rem)"
    fontWeight: 600
    lineHeight: 1.35
  body:
    fontFamily: "system-ui, -apple-system, sans-serif"
    fontSize: "clamp(1.00rem, 1.8vw, 1.125rem)"
    fontWeight: 400
    lineHeight: 1.72
  label:
    fontFamily: "system-ui, -apple-system, sans-serif"
    fontSize: "clamp(0.70rem, 1.2vw, 0.75rem)"
    fontWeight: 700
    letterSpacing: "0.18em"
rounded:
  sm: "3px"
  md: "6px"
  lg: "12px"
spacing:
  1:  "0.25rem"
  2:  "0.50rem"
  3:  "0.75rem"
  4:  "1.00rem"
  5:  "1.50rem"
  6:  "2.00rem"
  7:  "3.00rem"
  8:  "4.00rem"
  9:  "6.00rem"
  10: "8.00rem"
components:
  button-primary:
    backgroundColor: "{colors.ink-blue}"
    textColor: "oklch(98% 0.005 80)"
    rounded: "{rounded.md}"
    padding: "0.75rem 1.5rem"
  button-primary-hover:
    backgroundColor: "{colors.ink-blue-deep}"
    textColor: "oklch(98% 0.005 80)"
  button-secondary:
    backgroundColor: "transparent"
    textColor: "{colors.ink-blue}"
    rounded: "{rounded.md}"
    padding: "0.75rem 1.5rem"
  button-secondary-hover:
    backgroundColor: "{colors.ink-blue-subtle}"
    textColor: "{colors.ink-blue}"
  nav-card:
    backgroundColor: "{colors.warm-surface}"
    textColor: "{colors.near-black}"
    rounded: "{rounded.md}"
    padding: "1.50rem"
  nav-card-hover:
    backgroundColor: "{colors.warm-surface}"
    textColor: "{colors.near-black}"
---

# Design System: SAUI

## 1. Overview

**Creative North Star: "The Technical Manual"**

Not the corporate kind. The spiral-bound kind that lived on an engineer's desk for a decade: marked up in pencil, pages softened with use, cover slightly faded. Typeset with care because the content demanded respect, not because someone was trying to impress. Serious and functional, with the warmth of analogue production beneath the precision.

That tension — serious but with a memory of simpler, more deliberate things — runs through every design decision. Serif headings carry authority and a quiet nostalgia for physical documents. Near-white backgrounds have a faint warmth, like paper rather than backlit glass. The Ink Blue accent is the color of a technical pen annotation: purposeful, deliberate, never decorative. Animation exists only where it clarifies; silence is the default.

SAUI makes an argument. The design's job is to make that argument legible, not to perform trustworthiness. This system explicitly rejects the vocabulary of modern SaaS marketing and the ready-made feel of dashboard UI kits. It should look authored, not assembled. Nothing should look like it rolled out of a design system template.

**Key Characteristics:**
- Serif headings (Lora), system-ui body: reference text paired with operational text
- Warm amber-tinted neutrals throughout: paper, not screen
- Ink Blue held to approximately 10% of any surface — its rarity is its authority
- Flat surfaces at rest; depth through tonal layering only, not shadow
- Motion vocabulary limited to: scroll-entry reveals, hover lifts, and one architectural animation. Nothing ambient.

## 2. Colors: The Ledger Palette

A restrained palette: warm amber-tinted neutrals carrying the surface, a single Ink Blue marking structure and action, and deep slate reserved for code. The warmth is deliberate — paper, not a rendered white.

### Primary
- **Ink Blue** (`oklch(46% 0.120 240)`): The sole action color. Links, primary button fills, active nav states, focus rings, card hover borders, nav card heading text. Never used as a decorative background. Its scarcity is its meaning.
- **Ink Blue Deep** (`oklch(40% 0.140 240)`): Ink Blue on hover. Slightly darker and more saturated; signals state change.
- **Ink Blue Subtle** (`oklch(94% 0.040 240)`): The ghost fill — used for secondary button hover backgrounds and textarea focus rings. A blue so faint it reads as a warm surface tint.

### Neutral
- **Warm Paper** (`oklch(97% 0.008 80)`): Page background. The amber tint reads as paper rather than white light. Never replaced with pure white.
- **Warm Surface** (`oklch(93% 0.010 80)`): Elevated surfaces — nav cards, form fields. Slightly darker than the page.
- **Warm Border** (`oklch(87% 0.015 80)`): Dividers, card strokes, input borders at rest. The structural skeleton.
- **Near Black** (`oklch(14% 0.000 0)`): All display and heading text. Lightness is 14%, not 0 — deliberate, not a browser default.
- **Muted Text** (`oklch(50% 0.010 80)`): Secondary text, nav links at rest, captions, form labels. The voice beneath the voice.
- **Deep Slate** (`oklch(17% 0.020 240)`): Code block backgrounds only. A deep blue-tinted dark that reads as a terminal interior.
- **Code Light** (`oklch(88% 0.030 200)`): Code text on Deep Slate. Cool off-white with a faint cyan pull.

### Named Rules
**The Ink Reserve Rule.** Ink Blue appears on approximately 10% of any screen. Spreading it across decorative elements erases its authority. If Ink Blue is everywhere, it is nowhere.

**The Warm Surface Rule.** No surface is pure gray or pure white. Every background, border, and text neutral carries a faint warm tint toward hue 80 (amber). `oklch(X% 0 0)` and `#808080` are prohibited.

## 3. Typography

**Display Font:** Lora (Google Fonts; weights 400 and 600; italic variant), with a metric-matched `Lora-Fallback` face (Georgia-based, `size-adjust: 97%`) to eliminate layout shift on load.
**Body Font:** `system-ui, -apple-system, sans-serif` — the reader's own OS type, for legibility in running prose without a second network request.
**Mono Font:** `ui-monospace, 'Cascadia Code', 'Fira Code', monospace` — code blocks and inline code only.

**Character:** The serif/sans pairing is editorial. Lora carries authority and an Old Style warmth that recalls printed technical documentation; system-ui carries the immediacy and legibility of the machine. The contrast between the two registers is the system's typographic personality: headlines feel set, body text feels typed.

### Hierarchy
- **Display** (weight 600, `clamp(2.20rem, 7.0vw, 3.80rem)`, line-height 1.15, letter-spacing -0.02em): Page titles and hero h1s only. Never for section headings.
- **Headline** (weight 600, `clamp(1.80rem, 5.0vw, 2.80rem)`, line-height 1.15, letter-spacing -0.01em): Section h2s and major content divisions.
- **Title** (weight 600, `clamp(1.50rem, 4.0vw, 2.10rem)`, line-height 1.35): h3s, content card headings.
- **Body** (weight 400, `clamp(1.00rem, 1.8vw, 1.125rem)`, line-height 1.72): All running prose. Maximum line length 68ch. Dark mode line-height opens to 1.78.
- **Label** (weight 700, `clamp(0.70rem, 1.2vw, 0.75rem)`, letter-spacing 0.18em, uppercase, system-ui): Navigation labels, eyebrows, form field labels, section markers. The annotation layer.

### Named Rules
**The Two-Register Rule.** Lora for reading; system-ui for operating. Headings, body copy, blockquotes, and captions use the serif. Nav items, labels, buttons, and form elements use the sans. Never invert this.

**The Weight Contract.** Lora has exactly two weights in use: 400 and 600. Weight 500 is excluded — it produces visually indistinct hierarchy in Lora's design. The Google Fonts request reflects this.

## 4. Elevation

SAUI is a flat-by-default system. Surfaces are coplanar at rest. Depth is communicated entirely through tonal layering: three lightness steps in the neutral scale (page at 97%, surface at 93%, border at 87%). Shadow is not a structural tool.

The sole exception is the hover state on nav cards: `0 4px 16px -4px color-mix(in oklch, var(--c-accent) 20%, transparent)`. A whisper of Ink Blue shadow that confirms the card is interactive. It appears only on hover and disappears on active (mousedown). It is affordance feedback, not hierarchy.

### Named Rules
**The Flat-By-Default Rule.** No component carries a resting shadow. A shadow at rest means something is trying to look important. Flat surfaces let the content do that work. If a resting shadow is present, remove it and shift the background lightness instead.

## 5. Components

### Buttons

Measured and deliberate. Padding proportional to content, modest radius, structural border (2px solid), no gradient fills. Primary is Ink Blue with near-white text. Secondary is a ghost variant — transparent background, accent border, accent text.

- **Shape:** 6px radius — softened but not rounded
- **Primary:** Ink Blue fill, `oklch(98% 0.005 80)` text, 2px solid Ink Blue border, padding `0.75rem 1.5rem`. System-ui, weight 600, tracking 0.01em.
- **Hover:** Ink Blue Deep background, `translateY(-1px)`. Active resets to Y(0). No spring, no bounce, no elastic.
- **Focus:** 2px solid Ink Blue outline at 4px offset. Never absent; never replaced with a glow.
- **Secondary (ghost):** Transparent background, 2px solid Ink Blue border, Ink Blue text. Hover fills with Ink Blue Subtle.
- **Disabled:** 50% opacity, `pointer-events: none`. No additional styling.

### Inputs / Fields

Minimal at rest: 1px stroke in Warm Border on a Warm Paper background. On focus the stroke shifts to Ink Blue and a low-opacity ring appears. Precise state signal without visual noise.

- **Style:** 1px solid Warm Border, Warm Paper background, 3px radius
- **Focus:** Border to Ink Blue; `0 0 0 3px` outer ring in Ink Blue Subtle. Not a glow — a boundary.
- **Placeholder:** Muted Text at 60% opacity
- **Textarea:** Vertical resize only

### Nav Cards

Grid of wayfinding cards: flat and bordered at rest, lifting and accenting on hover. The internal heading is Ink Blue — the one instance of Ink Blue as text on a card surface.

- **Corner Style:** 6px radius
- **Background:** Warm Surface at rest and on hover
- **Border:** Warm Border at rest; shifts to Ink Blue on hover/focus
- **Shadow:** Ink Blue whisper-shadow on hover only (`0 4px 16px -4px` at 20% Ink Blue via `color-mix`)
- **Lift:** `translateY(-2px)` on hover; resets to 0 on active
- **Internal Padding:** 1.5rem
- **Heading Color:** Ink Blue — the card title is an action pointer, not decorative text

### Navigation

Sticky header, 64px block-size, Warm Paper background, Warm Border separator below. Brand wordmark: system-ui, weight 700, letter-spacing 0.14em, uppercase. Nav links: system-ui, Muted Text at rest, Ink Blue on hover and `aria-current="page"`. Touch target minimum 44px block-size.

Mobile: horizontal overflow scroll, scrollbar hidden. No hamburger menu — navigation always reachable.

### Code Blocks

A hard contrast register against the warm-paper page: Deep Slate background, Code Light text, monospace type at `--text-sm` with 1.65 line-height. Padding generous (1.5rem vertical, 2rem horizontal). Scrolls horizontally, never wraps.

Inline code: Warm Surface background, Warm Border stroke, 3px radius. Keeps the warm register rather than pulling in the dark code context.

### The Architecture Diagram (Signature Component)

The SVG lifecycle diagram is one of two visual centrepieces on the site (alongside code blocks). Four labelled nodes (Browser, Gateway, Event Log, Browser) connected by animated data packets. The Gateway node uses a scan-line animation — a blue fill sweeping left to right — one of only two ambient animations in the system. Node rectangles use 8px radius. The SVG is viewBox-scaled and responds to container width. On viewports below 600px, the side annotations are hidden to preserve readability of the core flow.

## 6. Do's and Don'ts

### Do:
- **Do** use Lora for all headings, display text, blockquotes, and captions. Use system-ui for all operational text: nav items, labels, buttons, and form elements.
- **Do** apply Ink Blue exclusively to interactive affordances: links, primary buttons, focus rings, active nav states, hover borders. Its scarcity is semantic.
- **Do** use the three-step tonal layering system (page 97% → surface 93% → border 87%) to communicate depth. These steps replace shadow.
- **Do** gate scroll-driven animations inside `@supports (animation-timeline: scroll())`. The full content must be readable if the animation is absent.
- **Do** implement all interactive states on every component: default, hover, focus-visible, active, disabled. Focus rings are always 2px solid Ink Blue.
- **Do** keep the motion vocabulary constrained: `translateY()` reveal on scroll entry, `translateY(-1px/-2px)` hover lift, data-packet animations in the architecture diagram. Nothing else moves.
- **Do** tint every neutral. All backgrounds and borders carry a faint warm tint (hue ~80). Code-context surfaces carry a cool blue tint (hue ~240). Pure gray (`oklch(X% 0 0)`) is prohibited.

### Don't:
- **Don't** use gradient text (`background-clip: text` on a gradient). Prohibited — decorative, never meaningful.
- **Don't** use side-stripe borders (border-left or border-right wider than 1px as a colored accent). Use full borders, background tints, or nothing.
- **Don't** add a resting shadow to any component. The Flat-By-Default Rule is non-negotiable.
- **Don't** use glassmorphism (blurred, translucent surfaces) anywhere, for any reason.
- **Don't** use Ink Blue as a background on non-button surfaces. It is an action marker, not a brand field.
- **Don't** produce the vocabulary of SaaS marketing: no hero metric templates (big number, small label, gradient accent), no purple/teal/neon AI color schemes, no glow borders.
- **Don't** let this look like it came out of a UI kit. No identical icon-heading-text card grids. No ready-made component feel. Every element should look considered, not assembled.
- **Don't** use modal dialogs as a first response. Exhaust inline and progressive alternatives.
