// Include the CSRF token from the meta tag in all mutating htmx requests.
document.body.addEventListener('htmx:configRequest', (e) => {
  if (e.detail.verb !== 'get') {
    const meta = document.querySelector('meta[name="csrf-token"]');
    if (meta) e.detail.headers['X-CSRF-Token'] = meta.content;
  }
});

// Mobile nav overlay
const mobileNav = document.getElementById('mobile-nav');
const navToggle = document.querySelector('.nav-toggle');

if (mobileNav && navToggle) {
  const navClose = mobileNav.querySelector('.mobile-nav__close');

  const openNav = () => {
    mobileNav.showModal();
    navToggle.setAttribute('aria-expanded', 'true');
  };

  const closeNav = () => {
    mobileNav.classList.add('is-closing');
    mobileNav.addEventListener('transitionend', () => {
      mobileNav.close();
      mobileNav.classList.remove('is-closing');
      navToggle.setAttribute('aria-expanded', 'false');
      navToggle.focus();
    }, { once: true });
  };

  navToggle.addEventListener('click', openNav);
  navClose?.addEventListener('click', closeNav);

  mobileNav.addEventListener('cancel', (e) => {
    e.preventDefault();
    closeNav();
  });

  window.addEventListener('resize', () => {
    if (mobileNav.open && window.matchMedia('(width >= 800px)').matches) {
      mobileNav.close();
      navToggle.setAttribute('aria-expanded', 'false');
    }
  });
}

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
