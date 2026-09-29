/* Progressive enhancement: one form, one save, freely navigable sections. */
(function () {
    'use strict';
    function init() {
        var form = document.getElementById('product-form');
        if (!form) return;
        var sections = Array.from(form.querySelectorAll(':scope > [data-section]'));
        if (sections.length !== 5) return;
        var labels = ['Basics', 'Content & media', 'Display & SEO'];
        var groups = [[sections[0]], [sections[1], sections[2]], [sections[3], sections[4]]];
        var current = 0, showAll = false;
        var controls = document.createElement('div');
        controls.className = 'editor-progress';
        controls.innerHTML = '<p class="editor-progress-note">Complete the essentials, then add details. Changes are saved only when you save the product.</p><nav aria-label="Product editing stages" class="editor-stages"></nav><div class="editor-progress-meta"><p role="status" aria-live="polite" data-stage-status></p><button type="button" data-show-all aria-pressed="false">Show all sections</button></div>';
        form.insertBefore(controls, sections[0]);
        var nav = controls.querySelector('nav');
        var status = controls.querySelector('[data-stage-status]');
        var panels = groups.map(function (items, index) {
            var panel = document.createElement('section');
            panel.id = 'product-stage-' + index;
            panel.className = 'editor-stage';
            panel.setAttribute('aria-label', labels[index]);
            form.insertBefore(panel, items[0]);
            items.forEach(function (item) {
                panel.appendChild(item);
                item.querySelector('.section-body').classList.remove('hidden');
                item.querySelector(':scope > button').setAttribute('aria-expanded', 'true');
                item.querySelector('.section-chevron').style.transform = 'rotate(0deg)';
            });
            var button = document.createElement('button');
            button.type = 'button';
            button.textContent = (index + 1) + '. ' + labels[index];
            button.setAttribute('aria-controls', panel.id);
            button.addEventListener('click', function () { showAll = false; show(index, true); });
            nav.appendChild(button);
            return panel;
        });
        var actions = form.querySelector(':scope > .pt-2');
        actions.classList.add('editor-actions');
        var pager = document.createElement('div');
        pager.className = 'editor-pager';
        pager.innerHTML = '<button type="button" data-stage-back>← Back</button><button type="button" data-stage-next>Continue →</button>';
        actions.before(pager);
        var back = pager.querySelector('[data-stage-back]');
        var next = pager.querySelector('[data-stage-next]');
        function show(index, focus) {
            current = index;
            panels.forEach(function (panel, i) { panel.hidden = !showAll && i !== current; });
            Array.from(nav.children).forEach(function (button, i) {
                if (i === current && !showAll) button.setAttribute('aria-current', 'step');
                else button.removeAttribute('aria-current');
            });
            controls.querySelector('[data-show-all]').setAttribute('aria-pressed', String(showAll));
            controls.querySelector('[data-show-all]').textContent = showAll ? 'Use stages' : 'Show all sections';
            status.textContent = showAll ? 'All sections shown' : 'Stage ' + (current + 1) + ' of 3 · ' + labels[current];
            pager.hidden = showAll;
            back.hidden = current === 0;
            next.hidden = current === 2;
            next.textContent = current === 0 ? 'Continue to content →' : 'Continue to display & SEO →';
            if (focus) {
                nav.children[current].focus({preventScroll: true});
                controls.scrollIntoView({block: 'start', behavior: 'instant'});
            }
        }
        controls.querySelector('[data-show-all]').addEventListener('click', function () { showAll = !showAll; show(current, false); });
        back.addEventListener('click', function () { show(current - 1, true); });
        next.addEventListener('click', function () {
            var invalid = Array.from(panels[current].querySelectorAll('input,textarea,select')).find(function (field) { return field.willValidate && !field.validity.valid; });
            if (invalid) { invalid.reportValidity(); return; }
            show(current + 1, true);
        });
        form.addEventListener('admin:reveal-field', function (event) {
            var index = panels.findIndex(function (panel) { return panel.contains(event.detail.field); });
            if (index >= 0 && !showAll) show(index, false);
        });
        // Error-summary links and deep links must expose their destination too.
        function revealHash() {
            var id; try { id = decodeURIComponent(location.hash.slice(1)); } catch (_) { return; }
            var field = document.getElementById(id);
            if (!field || !form.contains(field)) return;
            form.dispatchEvent(new CustomEvent('admin:reveal-field', {detail: {field: field}}));
            var body = field.closest('.section-body');
            if (body && body.classList.contains('hidden')) body.previousElementSibling.click();
            field.focus();
        }
        window.addEventListener('hashchange', revealHash);
        form.addEventListener('click', function (event) {
            if (event.target.closest('a[href^="#"]')) setTimeout(revealHash, 0);
        });
        show(0, false);
        revealHash();
    }
    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', init);
    else init();
})();
