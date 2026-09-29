// Run on a saved blog editor in a disposable local app using playwright-cli.
async (page) => {
    const editor = page.locator('trix-editor[aria-label="Body"]');
    await editor.waitFor();
    await page.waitForFunction(() => document.querySelector('textarea[data-rich-editor]').hidden);
    const originalHTML = await page.locator('textarea[data-rich-editor]').inputValue();
    if (await page.evaluate(() => AdminForms.isDirty())) throw new Error('Saved editor starts dirty');
    await editor.click();
    await page.keyboard.press('ControlOrMeta+End');
    await page.keyboard.type(' Synthetic regression edit.');
    await page.waitForFunction(() => AdminForms.isDirty(), null, {timeout: 5000});
    await page.keyboard.press('ControlOrMeta+z');
    await page.waitForFunction(() => !AdminForms.isDirty(), null, {timeout: 5000});
    if (await page.locator('textarea[data-rich-editor]').inputValue() !== originalHTML) throw new Error('Undo changed original HTML');
    await page.keyboard.press('ControlOrMeta+Shift+z');
    await page.waitForFunction(() => AdminForms.isDirty(), null, {timeout: 5000});
    await page.keyboard.press('ControlOrMeta+z');
    await page.waitForFunction(() => !AdminForms.isDirty(), null, {timeout: 5000});
    // The old test filled innerText back into the editor, which discards rich
    // formatting. That is a real change, not a valid test of Undo.
    await editor.fill(await editor.innerText());
    const changedHTML = await page.locator('textarea[data-rich-editor]').inputValue();
    if (changedHTML !== originalHTML && !await page.evaluate(() => AdminForms.isDirty())) throw new Error('Formatting changes must stay dirty');
    return {passed: 5, checks: ['saved editor starts clean', 'real edits are protected', 'Undo restores exact original HTML and clean state', 'Redo restores dirty state', 'formatting changes are not mistaken for Undo']};
}
