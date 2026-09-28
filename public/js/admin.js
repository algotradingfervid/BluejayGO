/* ============================================
   Bluejay CMS — Admin Sidebar JS
   ============================================ */

(function() {
    'use strict';

    var STORAGE_KEY = 'bluejay_sidebar_groups';

    // Get saved group states from localStorage
    function getSavedStates() {
        try {
            var raw = localStorage.getItem(STORAGE_KEY);
            return raw ? JSON.parse(raw) : {};
        } catch(e) {
            return {};
        }
    }

    // Save group states to localStorage
    function saveStates(states) {
        try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(states));
        } catch(e) {}
    }

    function syncGroup(group) {
        var open = group.classList.contains('open');
        var button = group.querySelector('.sidebar-group-header');
        var items = group.querySelector('.sidebar-group-items');
        if (button) button.setAttribute('aria-expanded', String(open));
        if (items) {
            items.inert = !open;
            items.id = items.id || 'sidebar-group-' + group.getAttribute('data-group');
            if (button) button.setAttribute('aria-controls', items.id);
        }
    }
    // Toggle a collapsible group
    window.toggleGroup = function(groupName) {
        var group = document.querySelector('[data-group="' + groupName + '"]');
        if (!group) return;

        var isOpen = group.classList.contains('open');
        if (isOpen) {
            group.classList.remove('open');
        } else {
            group.classList.add('open');
        }

        syncGroup(group);
        // Persist state
        var states = getSavedStates();
        states[groupName] = !isOpen;
        saveStates(states);
    };

    var sidebarQuery = window.matchMedia('(max-width: 1023px)');
    var sidebarOpener = null;
    function setSidebar(open, restoreFocus) {
        var sidebar = document.querySelector('.admin-sidebar');
        if (!sidebar) return;
        var compact = sidebarQuery.matches;
        sidebar.classList.toggle('open', compact && open);
        sidebar.inert = compact && !open;
        sidebar.setAttribute('aria-hidden', String(compact && !open));
        var overlay = document.querySelector('.sidebar-overlay');
        if (overlay) overlay.classList.toggle('active', compact && open);
        document.querySelectorAll('[data-sidebar-toggle]').forEach(function(button) {
            button.setAttribute('aria-expanded', String(compact && open));
            button.setAttribute('aria-label', open ? 'Close navigation' : 'Open navigation');
        });
        if (compact && open) {
            var close = sidebar.querySelector('[data-sidebar-close]');
            if (close) close.focus();
        } else if (restoreFocus && sidebarOpener && sidebarOpener.isConnected) {
            sidebarOpener.focus();
        }
    }
    window.toggleSidebar = function() {
        var sidebar = document.querySelector('.admin-sidebar');
        if (!sidebar || !sidebarQuery.matches) return;
        var open = !sidebar.classList.contains('open');
        if (open) sidebarOpener = document.activeElement;
        setSidebar(open, !open);
    };
    sidebarQuery.addEventListener('change', function() { setSidebar(false, false); });
    document.addEventListener('keydown', function(event) {
        var sidebar = document.querySelector('.admin-sidebar');
        if (!sidebarQuery.matches || !sidebar || !sidebar.classList.contains('open')) return;
        if (event.key === 'Escape') { event.preventDefault(); setSidebar(false, true); }
        if (event.key === 'Tab') {
            var items = Array.from(sidebar.querySelectorAll('a[href],button:not([disabled])')).filter(function(el) { return !el.closest('[inert]') && el.getClientRects().length && getComputedStyle(el).visibility !== 'hidden'; });
            var first = items[0], last = items[items.length - 1];
            if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus(); }
            if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus(); }
        }
    });

    // Initialize sidebar on page load
    function initSidebar() {
        var sidebar = document.querySelector('.admin-sidebar');
        if (sidebar) sidebar.parentElement.classList.add('admin-shell');
        setSidebar(false, false);
        var currentPath = window.location.pathname;
        var savedStates = getSavedStates();

        // Mark active link
        var allLinks = document.querySelectorAll('#sidebar-nav [data-path]');
        for (var i = 0; i < allLinks.length; i++) {
            var link = allLinks[i];
            var linkPath = link.getAttribute('data-path');
            if (currentPath === linkPath || currentPath.indexOf(linkPath + '/') === 0) {
                link.classList.add('active');
            }
        }

        // Find which group the active link belongs to and auto-expand it
        var activeLink = document.querySelector('#sidebar-nav [data-path].active');
        var activeGroupName = null;
        if (activeLink) {
            var parentGroup = activeLink.closest('.sidebar-group');
            if (parentGroup) {
                activeGroupName = parentGroup.getAttribute('data-group');
            }
        }

        // Apply saved states + auto-expand active group
        var groups = document.querySelectorAll('.sidebar-group');
        for (var j = 0; j < groups.length; j++) {
            var group = groups[j];
            var name = group.getAttribute('data-group');

            // Auto-expand if it contains the active page
            if (name === activeGroupName) {
                group.classList.add('open');
                // Also mark the group header as active
                var header = group.querySelector('.sidebar-group-header');
                if (header) header.classList.add('active');
            }
            // Or restore saved state
            else if (savedStates[name]) {
                group.classList.add('open');
            }
            syncGroup(group);
        }
    }

    // Run on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initSidebar);
    } else {
        initSidebar();
    }
})();

/* Shared form state. Compare values, not event counts, so undoing an edit is clean. */
(function() {
    'use strict';
    var records = new Map();
    var leaving = false;
    function fields(form) {
        return Array.from(form.elements).filter(function(el) {
            return el.name && !['submit', 'button', 'reset'].includes(el.type) && !el.closest('trix-toolbar, [data-dirty-ignore]');
        });
    }
    function value(el, initial) {
        if (el.type === 'checkbox' || el.type === 'radio') return [el.value, initial ? el.defaultChecked : el.checked];
        if (el.type === 'file') return initial ? [] : Array.from(el.files || []).map(function(f) { return [f.name, f.size, f.lastModified]; });
        if (el.tagName === 'SELECT') {
            var options = Array.from(el.options);
            var selected = options.filter(function(o) { return initial ? o.defaultSelected : o.selected; });
            if (initial && !selected.length && !el.multiple && options.length) selected = [options[0]];
            return selected.map(function(o) { return o.value; });
        }
        return initial ? el.defaultValue : el.value;
    }
    function snapshot(form, initial) {
        return JSON.stringify(fields(form).map(function(el) { return [el.name, value(el, initial)]; }));
    }
    function eligible(form) {
        var action = form.getAttribute('action') || '';
        return !form.hasAttribute('data-dirty-ignore') && !/\/admin\/(login|logout)$/.test(action) &&
            (form.hasAttribute('data-dirty-track') || form.method.toLowerCase() === 'post' || form.hasAttribute('hx-post') || form.hasAttribute('hx-put'));
    }
    function add(form) {
        if (!eligible(form) || records.has(form)) return;
        records.set(form, {baseline: snapshot(form, true), force: form.hasAttribute('data-dirty-on-load'), dirty: false, banner: null});
    }
    function refresh(form) {
        if (!form || !form.isConnected) return false;
        add(form);
        var record = records.get(form);
        if (!record) return false;
        record.dirty = record.force || snapshot(form, false) !== record.baseline;
        if (!record.banner) {
            record.banner = form.querySelector('[data-dirty-banner]') || (form.id === 'settings-form' ? document.querySelector('[data-dirty-banner]') : null);
            if (!record.banner && record.dirty) {
                record.banner = document.createElement('p');
                record.banner.className = 'admin-dirty-notice';
                record.banner.setAttribute('data-dirty-banner', '');
                record.banner.setAttribute('role', 'status');
                record.banner.textContent = 'You have unsaved changes.';
                form.prepend(record.banner);
            }
        }
        if (record.banner) {
            record.banner.hidden = !record.dirty;
            record.banner.classList.toggle('hidden', !record.dirty);
            if (record.banner.id === 'unsaved-banner') record.banner.classList.toggle('flex', record.dirty);
        }
        form.setAttribute('data-dirty', String(record.dirty));
        form.dispatchEvent(new CustomEvent('admin:dirtychange', {detail: {dirty: record.dirty}}));
        return record.dirty;
    }
    function refreshAll() {
        document.querySelectorAll('form').forEach(add);
        var dirty = false;
        records.forEach(function(record, form) {
            if (!form.isConnected) records.delete(form);
            else if (refresh(form)) dirty = true;
        });
        return dirty;
    }
    window.AdminForms = {
        refresh: refresh,
        isDirty: function(form) { return form ? refresh(form) : refreshAll(); },
        markSaved: function(form) {
            add(form);
            var record = records.get(form);
            if (record) { record.force = false; record.baseline = snapshot(form, false); refresh(form); }
        },
        discard: function(form) {
            // Explicit discard may reload immediately; permit that navigation only.
            if (form) { form.reset(); var record = records.get(form); if (record) record.force = false; }
            refreshAll();
            leaving = true;
        }
    };
    ['input', 'change', 'trix-change'].forEach(function(type) {
        document.addEventListener(type, function(event) {
            var form = event.target.form || event.target.closest('form');
            if (form) refresh(form);
        });
    });
    document.addEventListener('reset', function(event) { setTimeout(function() { refresh(event.target); }, 0); });
    document.addEventListener('click', function(event) {
        var link = event.target.closest('a[href]');
        if (!link || event.defaultPrevented || event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey || link.hasAttribute('download') || (link.target && link.target !== '_self')) return;
        var url = new URL(link.href, window.location.href);
        if (!['http:', 'https:'].includes(url.protocol) || (url.pathname === location.pathname && url.search === location.search && url.hash)) return;
        if (refreshAll()) {
            if (!window.confirm('You have unsaved changes. Discard them and leave this page? Choose Cancel to stay.')) {
                event.preventDefault();
                event.stopImmediatePropagation();
                return;
            }
            leaving = true;
        }
    }, true);
    document.addEventListener('submit', function(event) {
        // Native validation has passed before submit fires. HTMX prevents the default
        // and has its own completion event; client-cancelled submits stay dirty.
        queueMicrotask(function() { if (!event.defaultPrevented) leaving = true; });
    });
    document.addEventListener('htmx:afterRequest', function(event) {
        var detail = event.detail || {};
        var form = detail.elt && (detail.elt.tagName === 'FORM' ? detail.elt : detail.elt.closest('form'));
        var verb = (detail.requestConfig && detail.requestConfig.verb || '').toLowerCase();
        // Suggestions and child actions do not save the surrounding editor.
        if (form && detail.elt === form && ['post', 'put', 'patch'].includes(verb) && detail.successful) window.AdminForms.markSaved(form);
    });
    document.addEventListener('htmx:afterSwap', refreshAll);
    window.addEventListener('beforeunload', function(event) {
        if (!leaving && refreshAll()) { event.preventDefault(); event.returnValue = ''; }
    });
    window.addEventListener('pageshow', function() {
        leaving = false;
        refreshAll();
        // Some browser history restoration occurs after pageshow.
        setTimeout(refreshAll, 0);
    });
    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', refreshAll);
    else refreshAll();
})();
