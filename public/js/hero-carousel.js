(function () {
    'use strict';
    var root = document.querySelector('[data-hero-carousel]');
    if (!root) return;
    var track = root.querySelector('[data-hero-track]');
    var slides = root.querySelectorAll('[data-hero-slide]');
    if (slides.length < 2) return;
    var dots = root.querySelectorAll('[data-hero-dot]');
    var rotation = root.querySelector('[data-hero-rotation]');
    var status = root.querySelector('[data-hero-status]');
    var reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    var playing = root.getAttribute('data-autoplay') === '1' && !reducedMotion.matches;
    var hovering = false, focused = false, idx = 0, timer = null;
    var interval = Math.max(2, parseInt(root.getAttribute('data-interval'), 10) || 5) * 1000;

    function show(n) {
        idx = (n + slides.length) % slides.length;
        track.style.transform = 'translateX(-' + idx * 100 + '%)';
        slides.forEach(function (slide, i) {
            slide.setAttribute('aria-hidden', i === idx ? 'false' : 'true');
            slide.inert = i !== idx;
        });
        dots.forEach(function (dot, i) {
            dot.classList.toggle('bg-primary', i === idx);
            dot.classList.toggle('bg-gray-200', i !== idx);
            dot.setAttribute('aria-current', i === idx ? 'true' : 'false');
        });
        status.textContent = 'Slide ' + (idx + 1) + ' of ' + slides.length;
    }
    function sync() {
        if (timer !== null) clearInterval(timer);
        timer = null;
        rotation.textContent = playing ? 'Pause' : 'Play';
        rotation.setAttribute('aria-label', playing ? 'Pause slideshow' : 'Play slideshow');
        status.setAttribute('aria-live', playing ? 'off' : 'polite');
        if (playing && !hovering && !focused && !document.hidden) {
            timer = setInterval(function () { show(idx + 1); }, interval);
        }
    }
    function manual(n) { show(n); sync(); }
    root.querySelector('[data-hero-prev]').addEventListener('click', function () { manual(idx - 1); });
    root.querySelector('[data-hero-next]').addEventListener('click', function () { manual(idx + 1); });
    dots.forEach(function (dot, i) { dot.addEventListener('click', function () { manual(i); }); });
    rotation.addEventListener('click', function () {
        playing = !playing;
        // Explicit Play opts into rotation; subsequent entry into content pauses it.
        if (playing) { hovering = false; focused = false; }
        sync();
    });
    root.addEventListener('mouseenter', function () { hovering = true; sync(); });
    root.addEventListener('mouseleave', function () { hovering = false; sync(); });
    root.addEventListener('focusin', function () { focused = true; sync(); });
    root.addEventListener('focusout', function (event) {
        if (!root.contains(event.relatedTarget)) { focused = false; sync(); }
    });
    document.addEventListener('visibilitychange', sync);
    function updateMotion() {
        track.style.transitionDuration = reducedMotion.matches ? '0s' : '';
        if (reducedMotion.matches) playing = false;
        sync();
    }
    reducedMotion.addEventListener('change', updateMotion);
    show(0);
    updateMotion();
})();
