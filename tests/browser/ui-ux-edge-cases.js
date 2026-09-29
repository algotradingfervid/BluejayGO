// Run in playwright-cli against a disposable local DB after signing in.
async page=>{
 if (!['127.0.0.1','localhost'].includes(new URL(page.url()).hostname)) throw new Error('Use a disposable local test server only.');
 const base=new URL(page.url()).origin,checks=[];
 const check=(name,value)=>{if(!value)throw new Error(name);checks.push(name)};
 const screenshot=async name=>page.screenshot({path:'output/playwright/ui-ux-20260929/after/'+name+'.png'});
 await page.setViewportSize({width:390,height:844});await page.goto(base+'/');
 await page.keyboard.press('Tab');check('public skip link is first keyboard stop',await page.getByRole('link',{name:'Skip to content'}).evaluate(el=>el===document.activeElement));await page.keyboard.press('Enter');check('public skip link moves focus into main',await page.locator('#main-content').evaluate(el=>el===document.activeElement));
 await page.getByRole('button',{name:'Open navigation menu',exact:true}).click();await page.keyboard.press('Escape');check('mobile menu Escape restores toggle focus',await page.locator('[data-mobile-menu-toggle]').evaluate(el=>el===document.activeElement&&!document.querySelector('#mobile-navigation').getClientRects().length));
 await page.getByRole('button',{name:'Search',exact:true}).filter({visible:true}).click();
 await page.keyboard.press('Escape');check('search Escape restores opener',await page.locator('[data-open-search]').filter({visible:true}).evaluate(el=>el===document.activeElement));
 await page.getByRole('button',{name:'Show testimonial 2',exact:true}).click();check('testimonial selection announces current state',await page.getByRole('button',{name:'Show testimonial 2',exact:true}).getAttribute('aria-pressed')==='true'&&await page.locator('[data-testimonial]:visible').count()===1);
 check('testimonial touch target is 44px',await page.locator('[data-testimonial-btn]').first().evaluate(el=>el.getBoundingClientRect().width>=44&&el.getBoundingClientRect().height>=44));
 await page.goto(base+'/admin/products/new');await page.waitForFunction(()=>document.querySelector('textarea[data-rich-editor]').hidden);
 await page.keyboard.press('Tab');check('admin skip link is first keyboard stop',await page.getByRole('link',{name:'Skip to content'}).evaluate(el=>el===document.activeElement));await page.keyboard.press('Enter');check('admin skip reaches main',await page.locator('#admin-main').evaluate(el=>el===document.activeElement));
 await page.getByRole('button',{name:'Open navigation',exact:true}).click();
 check('drawer isolates background from interaction',await page.locator('#admin-main').evaluate(el=>el.inert));
 check('drawer has modal semantics',await page.locator('#admin-navigation').getAttribute('aria-modal')==='true');
 await page.keyboard.press('Shift+Tab');check('drawer reverse-tab remains inside',await page.locator('#admin-navigation').evaluate(el=>el.contains(document.activeElement)));
 await page.keyboard.press('Escape');check('drawer Escape restores menu focus',await page.locator('[data-sidebar-toggle]').evaluate(el=>el===document.activeElement&&!document.querySelector('#admin-main').inert));
 await page.getByRole('button',{name:'Open navigation',exact:true}).click();await page.setViewportSize({width:1440,height:1000});await page.waitForFunction(()=>!document.querySelector('#admin-main').inert);check('desktop resize clears modal and inert states',await page.locator('#admin-navigation').evaluate(el=>!el.inert&&!el.hasAttribute('aria-modal')&&!document.querySelector('#admin-main').inert));
 await page.emulateMedia({reducedMotion:'reduce'});
 await page.getByRole('button',{name:'2. Content & media',exact:true}).click();check('stage transition disabled for reduced motion',await page.locator('#product-stage-1').evaluate(el=>el.getAnimations().length===0&&getComputedStyle(el).animationName==='none'));
 await page.emulateMedia({reducedMotion:'no-preference'});await page.getByRole('button',{name:'3. Display & SEO',exact:true}).click();check('final stage timing is 180ms',await page.locator('#product-stage-2').evaluate(el=>getComputedStyle(el).animationDuration==='0.18s'));
 await page.emulateMedia({reducedMotion:'reduce'});check('live reduced-motion change cancels motion',await page.locator('#product-stage-2').evaluate(el=>el.getAnimations().length===0));await page.emulateMedia({reducedMotion:'no-preference'});
 for(let i=0;i<12;i++)await page.locator('.editor-stages button').nth(i%3).click();check('rapid stage changes leave one panel and one current marker',await page.locator('.editor-stage:visible').count()===1&&await page.locator('.editor-stages [aria-current=step]').count()===1);
 await page.goto(base+'/admin/solutions');const edit=await page.locator('a[href$="/edit"]').first().getAttribute('href');await page.goto(base+edit);
 const editor=page.getByRole('textbox',{name:'Full Description',exact:true});await editor.waitFor();await page.waitForFunction(()=>document.querySelector('textarea[data-rich-editor]').hidden);
 check('solution editor starts clean with saved HTML intact',!await page.evaluate(()=>AdminForms.isDirty()));const original=await page.locator('#overview-input').inputValue();
 await editor.fill('Synthetic solution editor save and reload verification.');
 await page.getByRole('button',{name:'Update Solution',exact:true}).click();await page.waitForLoadState('load');await page.goto(base+edit);await editor.waitFor();check('solution rich-text edit survives server round trip',(await editor.innerText()).includes('Synthetic solution editor save'));
 // Restore the original HTML through the real form without Trix normalizing it.
 await page.locator('#overview-input').evaluate((el,value)=>{el.value=value;el.dispatchEvent(new Event('input',{bubbles:true}))},original);await page.getByRole('button',{name:'Update Solution',exact:true}).click();await page.waitForLoadState('load');
 await page.goto(base+'/contact');await page.locator('#contact-name').fill('Failure recovery');await page.locator('#contact-email').fill('review@example.com');await page.locator('#contact-phone').fill('+91 98765 43210');await page.locator('#contact-company').fill('Local test');await page.locator('#contact-message').fill('Preserve my enquiry.');
 await page.route('**/contact/submit',route=>route.fulfill({status:500,body:'Synthetic server failure'}));await page.getByRole('button',{name:/Send Message/}).click();await page.locator('.ui-request-feedback').waitFor();
 check('failed enquiry shows retry guidance and preserves input',await page.locator('#contact-message').inputValue()==='Preserve my enquiry.'&&(await page.locator('.ui-request-feedback').innerText()).includes('try again'));
 await screenshot('contact-failure-recovery');await page.unroute('**/contact/submit');await page.getByRole('button',{name:/Send Message/}).click();await page.waitForFunction(()=>document.querySelector('#contact-message').value==='');check('successful retry clears failure notice',await page.locator('.ui-request-feedback').count()===0);
 // A fresh no-script context confirms that stage enhancement is optional.
 const context=await page.context().browser().newContext({javaScriptEnabled:false,storageState:await page.context().storageState(),viewport:{width:1440,height:1000}});
 try {const fallback=await context.newPage();await fallback.goto(base+'/admin/products/new');check('no-script product form retains core inputs',await fallback.locator('[name=name]').isVisible()&&await fallback.locator('[name=description]').isVisible());check('no-script product optional sections remain usable',await fallback.locator('[name=video_url]').isVisible()&&await fallback.locator('[name=meta_title]').isVisible());check('no-script rich HTML textarea stays editable',await fallback.locator('textarea[data-rich-editor]').isVisible());await fallback.screenshot({path:'output/playwright/ui-ux-20260929/after/no-script-product.png'});}finally{await context.close()}
 await page.goto(base+'/admin/products/settings');
 for (const badState of ['null','[true]','{broken']) {
  await page.evaluate(value=>localStorage.setItem('bluejay_sidebar_groups',value),badState);await page.reload();
  check('sidebar recovers from stored state '+badState,await page.locator('#sidebar-nav a[aria-current=page]').getAttribute('href')==='/admin/products/settings');
 }
 await page.evaluate(()=>localStorage.removeItem('bluejay_sidebar_groups'));
 return {passed:checks.length,checks};
}
