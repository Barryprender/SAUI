// Include the CSRF token from the meta tag in all mutating htmx requests.
document.body.addEventListener('htmx:configRequest', (e) => {
  if (e.detail.verb !== 'get') {
    const meta = document.querySelector('meta[name="csrf-token"]');
    if (meta) e.detail.headers['X-CSRF-Token'] = meta.content;
  }
});

// Zoom arch diagram into the node column on mobile by swapping viewBox.
const archSvg = document.querySelector('.arch-diagram svg');
if (archSvg) {
  const mq = window.matchMedia('(width < 600px)');
  const apply = (e) => archSvg.setAttribute('viewBox', e.matches ? '250 0 360 760' : '0 0 860 760');
  mq.addEventListener('change', apply);
  apply(mq);
}

// Mark the active step button after the arch step detail panel updates.
document.body.addEventListener('htmx:afterSwap', (e) => {
  if (e.target.id !== 'arch-step-detail') return;
  document.querySelectorAll('.arch-step-btn').forEach((btn) => {
    btn.setAttribute('aria-pressed', 'false');
  });
  const trigger = e.detail.requestConfig?.elt;
  if (trigger?.classList.contains('arch-step-btn')) {
    trigger.setAttribute('aria-pressed', 'true');
  }
});
