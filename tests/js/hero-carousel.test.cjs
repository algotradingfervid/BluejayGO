const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
class Element {
    constructor() { this.attrs = {}; this.style = {}; this.events = {}; this.classList = {toggle() {}}; }
    setAttribute(k, v) { this.attrs[k] = v; }
    getAttribute(k) { return this.attrs[k]; }
    addEventListener(k, fn) { this.events[k] = fn; }
    fire(k, event = {}) { this.events[k]?.(event); }
}
function fixture(reduce = false, autoplay = '1') {
    const root = new Element(), track = new Element(), rotation = new Element(), status = new Element();
    const prev = new Element(), next = new Element(), slides = [new Element(),new Element(),new Element()], dots = slides.map(() => new Element());
    const doc = new Element(), motion = new Element(), timers = new Map(); let serial = 0;
    root.attrs = {'data-autoplay': autoplay, 'data-interval': '5'};
    root.querySelector = s => ({'[data-hero-track]':track,'[data-hero-rotation]':rotation,'[data-hero-status]':status,'[data-hero-prev]':prev,'[data-hero-next]':next}[s]);
    root.querySelectorAll = s => s === '[data-hero-slide]' ? slides : dots;
    root.contains = el => slides.includes(el) || el === rotation;
    doc.querySelector = () => root; doc.hidden = false; motion.matches = reduce;
    vm.runInNewContext(fs.readFileSync('public/js/hero-carousel.js', 'utf8'), {document:doc, window:{matchMedia:()=>motion},setInterval:fn=>{timers.set(++serial,fn);return serial;}, clearInterval:id=>timers.delete(id)});
    return {root,track,rotation,status,prev,next,slides,dots,doc,motion,timers,tick(){ [...timers.values()].forEach(fn=>fn()); }};
}
test('explicit pause survives hover, focus and manual navigation; Play resumes', () => {
    const f = fixture(); assert.equal(f.timers.size,1); f.tick();
    assert.equal(f.track.style.transform,'translateX(-100%)');
    f.rotation.fire('click'); assert.equal(f.timers.size,0);
    f.root.fire('mouseenter'); f.root.fire('mouseleave'); f.root.fire('focusin'); f.root.fire('focusout',{relatedTarget:null}); f.next.fire('click');
    assert.equal(f.timers.size,0); assert.equal(f.rotation.attrs['aria-label'],'Play slideshow');
    f.rotation.fire('click'); assert.equal(f.timers.size,1);
});
test('reduced motion starts manual, disables transitions and responds to preference changes', () => {
    const f = fixture(true); assert.equal(f.timers.size,0); assert.equal(f.track.style.transitionDuration,'0s');
    f.next.fire('click'); assert.equal(f.track.style.transform,'translateX(-100%)'); assert.equal(f.timers.size,0);
    f.rotation.fire('click'); assert.equal(f.timers.size,1);
    f.motion.fire('change'); assert.equal(f.timers.size,0);
});
test('inactive slides are inert and hidden; arrows and dots wrap correctly', () => {
    const f = fixture(false,'0'); assert.equal(f.timers.size,0); f.prev.fire('click');
    assert.deepEqual(f.slides.map(s=>s.inert),[true,true,false]);
    assert.deepEqual(f.slides.map(s=>s.attrs['aria-hidden']),['true','true','false']);
    f.dots[0].fire('click'); assert.equal(f.track.style.transform,'translateX(-0%)');
    assert.equal(f.dots[0].attrs['aria-current'],'true'); assert.equal(f.status.textContent,'Slide 1 of 3');
});
test('hover, focus and background tabs suspend automatic changes', () => {
    const f = fixture(); f.root.fire('mouseenter'); assert.equal(f.timers.size,0);
    f.root.fire('focusin'); f.root.fire('mouseleave'); assert.equal(f.timers.size,0);
    f.root.fire('focusout',{relatedTarget:null}); assert.equal(f.timers.size,1);
    f.doc.hidden=true; f.doc.fire('visibilitychange'); assert.equal(f.timers.size,0);
});
