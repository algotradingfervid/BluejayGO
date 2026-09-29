const {test} = require('node:test');
const assert = require('node:assert/strict');
const vm = require('node:vm');
const fs = require('node:fs');
function fixture() {
    let doc;
    class Element {
        constructor(hidden=false) { this.hidden=hidden; this.attrs={}; this.events={};this.isConnected=true;this.style={}; const classes=new Set(hidden?['hidden']:[]); this.classList={contains:s=>classes.has(s),add:s=>classes.add(s),remove:s=>classes.delete(s),toggle(s,on){if(on)classes.add(s);else classes.delete(s);}}; }
        setAttribute(k,v){this.attrs[k]=v;} removeAttribute(k){delete this.attrs[k];}
        addEventListener(k,fn){this.events[k]=fn;} fire(k,e={}){this.events[k]?.(e);}
        contains(el){return el===this;} focus(){doc.activeElement=this;}
    }
    doc=new Element();doc.body=new Element();doc.body.style.overflow='auto';
    const toggle=new Element(), menu=new Element(true), modal=new Element(true),input=new Element(),results=new Element(),close=new Element(),search=new Element(),desktop=new Element();
    const nodes={'mobile-navigation':menu,'search-modal':modal,'search-input':input,'search-results':results};
    doc.getElementById=id=>nodes[id];doc.querySelector=()=>toggle;doc.querySelectorAll=()=>[search];doc.activeElement=toggle;
    modal.querySelector=()=>close;modal.querySelectorAll=()=>[input,close];
    vm.runInNewContext(fs.readFileSync('public/js/public-navigation.js','utf8'),{document:doc,window:{matchMedia:()=>desktop}});
    return {doc,toggle,menu,modal,input,results,close,search,desktop,key(key,extra={}){let prevented=false;doc.fire('keydown',{key,preventDefault(){prevented=true;},...extra});return prevented;}};
}
test('mobile navigation toggles, Escape closes and returns focus, desktop resize resets',()=>{
    const f=fixture();f.toggle.fire('click');assert.equal(f.menu.hidden,false);assert.equal(f.toggle.attrs['aria-expanded'],'true');
    f.doc.activeElement=f.menu;f.key('Escape');assert.equal(f.menu.hidden,true);assert.equal(f.doc.activeElement,f.toggle);
    f.toggle.fire('click');f.desktop.fire('change',{matches:true});assert.equal(f.menu.hidden,true);assert.equal(f.toggle.attrs['aria-expanded'],'false');
});
test('search closes mobile navigation, traps focus and restores opener plus scrolling',()=>{
    const f=fixture();f.toggle.fire('click');f.search.focus();f.search.fire('click');
    assert.equal(f.menu.hidden,true);assert.equal(f.modal.classList.contains('hidden'),false);assert.equal(f.doc.activeElement,f.input);assert.equal(f.doc.body.style.overflow,'hidden');
    assert.equal(f.key('Tab',{shiftKey:true}),true);assert.equal(f.doc.activeElement,f.close);
    assert.equal(f.key('Tab'),true);assert.equal(f.doc.activeElement,f.input);
    f.key('Escape');assert.equal(f.doc.activeElement,f.search);assert.equal(f.doc.body.style.overflow,'auto');
});
test('pending suggestions signal busy without taking search focus',()=>{
    const f=fixture();f.search.fire('click');f.input.fire('htmx:beforeRequest');assert.equal(f.results.attrs['aria-busy'],'true');
    f.input.fire('htmx:afterRequest');assert.equal(f.results.attrs['aria-busy'],undefined);assert.equal(f.doc.activeElement,f.input);
});
