// Run on a saved blog editor in a disposable local app using playwright-cli.
async (page) => {
    const editor = page.locator('trix-editor[aria-label="Body"]');
    await editor.waitFor();
    await page.waitForFunction(() => document.querySelector('textarea[data-rich-editor]').hidden);
    const initial = await editor.innerText();
    if (await page.evaluate(() => AdminForms.isDirty())) throw new Error('Saved editor starts dirty');
    await editor.fill(initial + ' Synthetic regression edit.');
    await page.waitForFunction(() => AdminForms.isDirty(), null, {timeout: 5000});
    await editor.fill(initial);
    await page.waitForFunction(() => !AdminForms.isDirty(), null, {timeout: 5000});
    return {passed: 3, checks: ['saved editor starts clean', 'real edits are protected', 'reverting text returns to clean']};
}
