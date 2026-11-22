import {createContext, createElement as h, useContext} from 'react';
import {safeLink, type DisplayArticle, type Source} from './model';

export const Controls = createContext<any>(null);
const tag = (name: string) => (p: any) => h(name, {className: p.variant || undefined}, p.children);
export const Page = tag('main'), Row = tag('div'), Card = tag('section'), Section = tag('section');
export const Text = tag('p');
export const Divider = () => h('hr');
export const Alert = (p: any) => h('p', {role: 'status', className: 'notice'}, p.children);
export const Heading = (p: any) => h('h' + (p.level || 2), {}, p.children);
export const Link = (p: any) => h('a', {href: p.href, className: p.variant,
  'data-workspace-apps-trigger': p.variant === 'apps' ? 'true' : undefined}, p.children);
export function Button(p: any) {
  const controls = useContext(Controls);
  return h('button', {type: 'button', className: p.variant, disabled: p.disabled,
    onClick: () => { controls.item.current = p.item; p.onClick?.(); }}, p.children);
}
export function Input(p: any) {
  const controls = useContext(Controls);
  return h('label', {className: 'field'}, p.label, h('input', {type: p.type || 'text',
    placeholder: p.placeholder, maxLength: p.name === 'query' ? 128 : 512,
    autoComplete: p.type === 'password' ? 'new-password' : 'off',
    value: controls.fields[p.name] || '', onChange: (event: any) => controls.field(p.name, event.target.value),
    onKeyDown: (event: any) => { if (event.key === 'Enter' && p.name === 'query') controls.action('search'); }}));
}
export function Select(p: any) {
  const controls = useContext(Controls);
  return h('label', {className: 'field'}, p.label, h('select',
    {value: controls.fields[p.name], onChange: (event: any) => controls.field(p.name, event.target.value)},
    p.options.map((value: any) => h('option', {key: value.id, value: value.id, disabled: value.disabled}, value.label))));
}
export function Cover({tone}: {tone: string}) {
  return h('div', {className: 'cover ' + tone, 'aria-hidden': true},
    h('i'), h('i'), h('i'), h('i'), h('b'));
}
export function Story({item}: {item: DisplayArticle}) {
  const controls = useContext(Controls);
  return h('article', {className: 'story'},
    h(Cover, {tone: item.tone}),
    h('div', {}, h('p', {className: 'kicker'}, item.label),
      h('h3', {}, h('button', {className: 'story-title', onClick: () => controls.action('read', item.id)}, item.title)),
      h('p', {className: 'story-summary'}, item.summary), h('p', {className: 'byline'}, item.byline)),
    h('button', {className: 'save-icon', disabled: controls.bookmarkDisabled,
      'aria-label': (item.saved ? 'Remove ' : 'Save ') + item.title,
      'aria-pressed': item.saved, onClick: () => controls.action('bookmark', item.id)}, item.saved ? '★' : '☆'));
}
export function Paragraphs({body, refs}: {body: string; refs: number[]}) {
  return h('div', {},
    body.split(/\n\s*\n/).filter(Boolean).map((part, i) => h('p', {key: i}, part)),
    refs.length > 0 && h('p', {className: 'ref-tags'}, [...new Set(refs)].map(i =>
      h('a', {key: i, href: '#source-' + i, 'aria-label': 'Go to source ' + i,
        onClick: event => { event.preventDefault(); const target = document.getElementById('source-' + i); target?.focus(); target?.scrollIntoView({block: 'center'}); }}, '[' + i + ']'))));
}
export function SourceReference({item}: {item: Source & {index: number}}) {
  const controls = useContext(Controls);
  return h('div', {className: 'reference', id: 'source-' + item.index, tabIndex: -1},
    h('span', {}, '[' + item.index + '] '),
    !controls.offline && safeLink(item.url)
      ? h('a', {href: item.url, target: '_blank', rel: 'noopener noreferrer'}, item.title)
      : h('strong', {}, item.title),
    h('p', {}, item.date ? item.date.slice(0, 10) + ' · ' + item.url : item.url));
}
export function SourceChoice({item}: {item: Source & {index: number; selected: boolean}}) {
  const controls = useContext(Controls);
  return h('div', {className: 'source-choice'},
    h('label', {}, h('input', {type: 'checkbox', checked: item.selected,
      disabled: controls.busy || (!item.selected && controls.selectedCount >= 5),
      onChange: () => controls.pickSource(item.index)}), item.title),
    h('a', {href: item.url, target: '_blank', rel: 'noopener noreferrer'}, new URL(item.url).hostname),
    item.date && h('span', {}, ' · ' + item.date.slice(0, 10)),
    h('p', {}, item.text));
}
