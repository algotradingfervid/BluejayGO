(function () {
    'use strict';
    var columnID = 0;
    var requestNotices = new WeakMap();
    function unavailable(image) {
        if (!(image instanceof HTMLImageElement) || !image.getAttribute('src') || image.dataset.fallbackApplied) return;
        image.dataset.fallbackApplied = 'true';
        image.alt = (image.alt ? image.alt + ' — ' : '') + 'Image unavailable';
        image.src = '/public/images/image-unavailable.svg';
    }
    function requestFailed(event) {
        var detail = event.detail || {};
        var target = detail.target || detail.elt;
        if (!(target instanceof Element) || target === document.body) return;
        // Contact validation responses already provide specific inline messages.
        if (target.id === 'contact-form-feedback' && detail.xhr && detail.xhr.status === 400) return;
        var notice = requestNotices.get(target);
        if (!notice || !notice.isConnected) {
            notice = document.createElement('p');
            notice.className = 'ui-request-feedback';
            notice.setAttribute('role', 'alert');
            target.before(notice);
            requestNotices.set(target, notice);
        }
        notice.textContent = 'The request could not be completed. Your entries are still here. Please try again.';
    }
    document.addEventListener('htmx:responseError', requestFailed);
    document.addEventListener('htmx:sendError', requestFailed);
    document.addEventListener('htmx:timeout', requestFailed);
    document.addEventListener('htmx:afterRequest', function (event) {
        var detail = event.detail || {};
        if (!detail.successful) return;
        var notice = requestNotices.get(detail.target);
        if (notice) { notice.remove(); requestNotices.delete(detail.target); }
    });
    document.addEventListener('error', function (event) { unavailable(event.target); }, true);
    function enhance() {
        document.querySelectorAll('img').forEach(function (img) { if (img.complete && !img.naturalWidth) unavailable(img); });
        var sidebar = document.querySelector('.admin-sidebar');
        var shell = sidebar && sidebar.parentElement;
        document.querySelectorAll('trix-toolbar input[data-trix-input]').forEach(function (input) { input.setAttribute('aria-label', 'Link URL'); });
        if (!shell) return;
        var main = shell.querySelector(':scope > .flex-1');
        if (main) { main.id = 'admin-main'; main.setAttribute('role', 'main'); main.tabIndex = -1; }
        shell.querySelectorAll('table').forEach(function (table) {
            var headings = Array.from(table.querySelectorAll('thead th'));
            if (!headings.length) return;
            table.classList.add('responsive-records');
            table.setAttribute('role', 'table');
            table.querySelectorAll('thead,tbody').forEach(function (group) { group.setAttribute('role', 'rowgroup'); });
            table.querySelectorAll('tr').forEach(function (row) { row.setAttribute('role', 'row'); });
            headings.forEach(function (heading) {
                heading.id = heading.id || 'record-column-' + (++columnID);
                heading.scope = 'col'; heading.setAttribute('role', 'columnheader');
            });
            table.querySelectorAll('tbody tr').forEach(function (row) {
                Array.from(row.cells).forEach(function (cell, index) {
                    cell.setAttribute('role', 'cell');
                    if (cell.colSpan !== 1 || !headings[index] || cell.querySelector('.record-label')) return;
                    cell.setAttribute('headers', headings[index].id);
                    var label = document.createElement('span');
                    label.className = 'record-label'; label.setAttribute('aria-hidden', 'true');
                    label.textContent = headings[index].textContent.trim();
                    var value = document.createElement('div');
                    value.className = 'record-value';
                    while (cell.firstChild) value.appendChild(cell.firstChild);
                    cell.append(label, value);
                });
            });
        });
        shell.querySelectorAll('[data-section] > button').forEach(function (button, index) {
            var body = button.nextElementSibling;
            if (!body) return;
            body.id = body.id || 'editor-section-' + index;
            button.setAttribute('aria-controls', body.id);
            button.setAttribute('aria-expanded', String(!body.classList.contains('hidden')));
            if (!button.dataset.disclosureReady) {
                button.dataset.disclosureReady = 'true';
                button.addEventListener('click', function () { button.setAttribute('aria-expanded', String(!body.classList.contains('hidden'))); });
            }
        });
    }
    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', enhance);
    else enhance();
    document.addEventListener('htmx:afterSwap', enhance);
    document.addEventListener('trix-initialize', enhance);
})();
