(function () {
    'use strict';
    var toggle = document.querySelector('[data-mobile-menu-toggle]');
    var menu = document.getElementById('mobile-navigation');
    var desktop = window.matchMedia('(min-width: 1024px)');
    function setMenuOpen(open, restoreFocus) {
        menu.hidden = !open;
        menu.classList.toggle('hidden', !open);
        toggle.setAttribute('aria-expanded', String(open));
        toggle.setAttribute('aria-label', open ? 'Close navigation menu' : 'Open navigation menu');
        if (restoreFocus) toggle.focus({preventScroll: true});
    }
    toggle.addEventListener('click', function () { setMenuOpen(menu.hidden); });
    desktop.addEventListener('change', function (event) { if (event.matches) setMenuOpen(false); });
    document.addEventListener('click', function (event) {
        if (!menu.hidden && !menu.contains(event.target) && !toggle.contains(event.target)) setMenuOpen(false);
    });

    var modal = document.getElementById('search-modal');
    var input = document.getElementById('search-input');
    var results = document.getElementById('search-results');
    var previousFocus, previousOverflow;
    function setSearchOpen(open) {
        if (open) {
            if (!modal.classList.contains('hidden')) return;
            previousFocus = document.activeElement;
            previousOverflow = document.body.style.overflow;
            setMenuOpen(false);
            modal.classList.remove('hidden');
            document.body.style.overflow = 'hidden';
            input.focus();
        } else {
            if (modal.classList.contains('hidden')) return;
            modal.classList.add('hidden');
            document.body.style.overflow = previousOverflow;
            if (previousFocus && previousFocus.isConnected) previousFocus.focus({preventScroll: true});
        }
    }
    document.querySelectorAll('[data-open-search]').forEach(function (button) {
        button.addEventListener('click', function () { setSearchOpen(true); });
    });
    modal.querySelector('[data-close-search]').addEventListener('click', function () { setSearchOpen(false); });
    modal.addEventListener('click', function (event) { if (event.target === modal) setSearchOpen(false); });
    input.addEventListener('htmx:beforeRequest', function () {
        results.setAttribute('aria-busy', 'true');
    });
    input.addEventListener('htmx:afterRequest', function () { results.removeAttribute('aria-busy'); });
    input.addEventListener('htmx:responseError', function () {
        results.textContent = 'Search is unavailable right now. Please try again.';
        var link = document.createElement('a');
        link.href = '/products';
        link.className = 'block underline font-bold py-3';
        link.textContent = 'Browse all products';
        results.appendChild(link);
    });
    document.addEventListener('keydown', function (event) {
        if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
            event.preventDefault();
            setSearchOpen(modal.classList.contains('hidden'));
        }
        if (event.key === 'Escape') {
            if (!modal.classList.contains('hidden')) setSearchOpen(false);
            else if (!menu.hidden) setMenuOpen(false, true);
        }
        if (modal.classList.contains('hidden') || event.key !== 'Tab') return;
        var items = modal.querySelectorAll('input, button, a[href]');
        var first = items[0], last = items[items.length - 1];
        if (event.shiftKey && document.activeElement === first) {
            event.preventDefault(); last.focus();
        } else if (!event.shiftKey && document.activeElement === last) {
            event.preventDefault(); first.focus();
        }
    });
})();
