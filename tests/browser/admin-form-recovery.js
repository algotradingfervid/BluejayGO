// Run with playwright-cli run-code --filename after opening the local app.
// Uses real browser constraint validation; all data lives in a disposable DOM.
async (page) => {
    const origin = await page.evaluate(() => location.origin);
    await page.goto(origin + '/__qa-form-fixture__');
    await page.setContent(`
        <style>.hidden { display:none } button,input { margin:8px; padding:8px }</style>
        <form id="editor" method="post"><input name="title" value="Saved"></form>
        <div id="detail"><form id="inline" method="post"><input name="spec" value="2 USB-C"></form></div>
        <div id="suggestions"></div><button id="tab">Features</button>
        <form id="validation" method="post">
          <button type="button" id="collapse"><span class="section-chevron">Content</span></button>
          <div id="section" class="hidden"><input id="required" name="description" required></div>
          <button>Save content</button>
        </form>
        <form id="settings" method="post">
          <div id="tab-general" class="tab-content"><input name="site" value="Site"></div>
          <div id="tab-contact" class="tab-content hidden"><input id="email" name="email" type="email" value="bad" required></div>
          <div id="tab-social" class="tab-content hidden"><input name="url" type="url" value="also-bad" required></div>
          <button>Save settings</button>
        </form>`);
    await page.evaluate(() => {
        document.querySelector('#collapse').onclick = () => document.querySelector('#section').classList.toggle('hidden');
        window.switchTab = name => document.querySelectorAll('.tab-content').forEach(el => el.classList.toggle('hidden', el.id !== 'tab-' + name));
        document.querySelectorAll('form').forEach(form => form.addEventListener('submit', event => event.preventDefault()));
    });
    await page.addScriptTag({url: origin + '/public/js/admin.js'});
    const checks = [];
    const check = (name, value) => { if (!value) throw new Error(name); checks.push(name); };
    await page.getByRole('button', {name: 'Save content'}).click();
    check('invalid collapsed section opens and focuses its field', await page.locator('#required').evaluate(el => !!el.getClientRects().length && document.activeElement === el));
    await page.locator('#required').fill('Restored description');
    await page.getByRole('button', {name: 'Save content'}).click();
    check('corrected field is valid', await page.locator('#validation').evaluate(el => el.checkValidity()));
    await page.getByRole('button', {name: 'Save settings'}).click();
    check('first hidden invalid tab opens and receives focus', await page.locator('#email').evaluate(el => !!el.getClientRects().length && document.activeElement === el));
    check('later invalid tab does not hide first error', await page.locator('#tab-social').evaluate(el => el.classList.contains('hidden')));
    await page.locator('#email').fill('qa@example.com');
    await page.getByRole('button', {name: 'Save settings'}).click();
    check('next invalid field is reachable on the next submission', await page.locator('#tab-social').evaluate(el => !el.classList.contains('hidden')));

    await page.locator('#inline input').fill('4 USB-C');
    const cancel = await page.evaluate(() => {
        window.confirm = () => false;
        let issued = false;
        const event = new CustomEvent('htmx:confirm', {bubbles:true,cancelable:true,detail:{verb:'get',target:document.querySelector('#detail'),issueRequest:()=>{issued=true;}}});
        document.querySelector('#tab').dispatchEvent(event);
        return {blocked:event.defaultPrevented,issued,value:document.querySelector('#inline input').value};
    });
    check('Cancel prevents replacing dirty detail form', cancel.blocked && !cancel.issued && cancel.value === '4 USB-C');
    const accept = await page.evaluate(() => {
        window.confirm = () => true;
        let issued = false;
        const event = new CustomEvent('htmx:confirm', {bubbles:true,cancelable:true,detail:{verb:'get',target:document.querySelector('#detail'),issueRequest:skip=>{issued=skip;}}});
        document.querySelector('#tab').dispatchEvent(event);
        return event.defaultPrevented && issued;
    });
    check('Discard resumes the requested detail replacement', accept);
    const adjacent = await page.evaluate(() => {
        window.confirm = () => { throw new Error('unrelated request prompted'); };
        const request = (verb,target) => { const event=new CustomEvent('htmx:confirm',{bubbles:true,cancelable:true,detail:{verb,target,elt:document.querySelector('#inline')}});document.querySelector('#tab').dispatchEvent(event);return event.defaultPrevented; };
        return {suggestion:request('get',document.querySelector('#suggestions')),save:request('post',document.querySelector('#detail'))};
    });
    check('suggestion GET does not discard surrounding edits', !adjacent.suggestion);
    check('save POST is not intercepted as navigation', !adjacent.save);
    const sibling = await page.evaluate(() => {
        window.confirm = () => false;
        let issued = false;
        const event = new CustomEvent('htmx:confirm', {bubbles:true,cancelable:true,detail:{verb:'post',elt:document.querySelector('#tab'),target:document.querySelector('#detail'),issueRequest:()=>{issued=true;}}});
        document.querySelector('#tab').dispatchEvent(event);
        return event.defaultPrevented && !issued;
    });
    check('sibling mutation cannot silently replace another dirty form', sibling);
    await page.locator('#inline input').fill('2 USB-C');
    const clean = await page.evaluate(() => {
        const event=new CustomEvent('htmx:confirm',{bubbles:true,cancelable:true,detail:{verb:'get',target:document.querySelector('#detail')}});
        document.querySelector('#tab').dispatchEvent(event); return event.defaultPrevented;
    });
    check('undoing inline edits allows clean navigation without prompt', !clean);
    return {passed:checks.length,checks};
}
