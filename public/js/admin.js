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

        // Persist state
        var states = getSavedStates();
        states[groupName] = !isOpen;
        saveStates(states);
    };

    // Toggle mobile sidebar
    window.toggleSidebar = function() {
        var sidebar = document.querySelector('.admin-sidebar');
        var overlay = document.querySelector('.sidebar-overlay');
        if (sidebar) sidebar.classList.toggle('open');
        if (overlay) overlay.classList.toggle('active');
    };

    // Initialize sidebar on page load
    function initSidebar() {
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
        }
    }

    // Run on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initSidebar);
    } else {
        initSidebar();
    }
})();
