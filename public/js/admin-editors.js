// The named textarea remains the submitted value and a usable no-script fallback.
(function () {
    'use strict';
    function enhance() {
        if (!window.customElements || !customElements.get('trix-editor')) return;
        document.querySelectorAll('textarea[data-rich-editor]').forEach(function (source) {
            if (source.dataset.editorReady) return;
            source.dataset.editorReady = 'true';
            var editor = document.createElement('trix-editor');
            editor.id = source.id + '-editor';
            editor.setAttribute('input', source.id);
            editor.setAttribute('aria-label', source.getAttribute('aria-label'));
            editor.className = 'trix-content border-2 border-black min-h-[200px] text-sm';
            editor.addEventListener('trix-file-accept', function (event) {
                event.preventDefault();
                var notice = source.parentElement.querySelector('[data-editor-upload-notice]');
                if (!notice) {
                    notice = document.createElement('p');
                    notice.dataset.editorUploadNotice = '';
                    notice.setAttribute('role', 'alert');
                    notice.className = 'text-sm text-red-700 mt-2';
                    notice.textContent = 'Inline file uploads are not supported. Use this page’s image or media fields instead.';
                    editor.insertAdjacentElement('afterend', notice);
                }
            });
            var original = source.value;
            var initialDocument = null;
            editor.addEventListener('trix-initialize', function () {
                source.hidden = true;
                // There is no attachment upload endpoint for these editors.
                var attachmentButton = editor.toolbarElement.querySelector('[data-trix-action="attachFiles"]');
                if (attachmentButton) {
                    attachmentButton.disabled = true;
                    attachmentButton.hidden = true;
                    attachmentButton.style.display = 'none';
                }
                var help = source.parentElement.querySelector('[data-editor-fallback]');
                if (help) help.hidden = true;
                document.querySelectorAll('label').forEach(function (label) {
                    if (label.htmlFor === source.id) label.htmlFor = editor.id;
                });
                // Loading the editor must not change the saved HTML until the user edits.
                initialDocument = JSON.stringify(editor.editor.getDocument().toJSON());
                source.value = original;
                if (window.AdminForms) window.AdminForms.refresh(source.form);
            }, {once: true});
            editor.addEventListener('trix-change', function () {
                // Undo can restore the original document after Trix has normalized
                // its HTML. Preserve the saved HTML (including unsupported IDs)
                // when both text and formatting return to their initial state.
                if (initialDocument !== null && JSON.stringify(editor.editor.getDocument().toJSON()) === initialDocument) {
                    source.value = original;
                }
                source.dispatchEvent(new Event('input', {bubbles: true}));
            });
            source.insertAdjacentElement('afterend', editor);
        });
    }
    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', enhance);
    else enhance();
    // Trix registers its custom elements on a timer; no timeout hides the fallback.
    if (window.customElements) customElements.whenDefined('trix-editor').then(enhance);
})();
