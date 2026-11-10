'use strict';
const spec = JSON.parse(atob(document.body.dataset.spec));
const board = document.querySelector('canvas');
const ctx = board.getContext('2d');
const status = document.querySelector('#status');
const controls = document.querySelector('#controls');
const start = document.querySelector('#start');
const heroes = {rocket:'🚀',cat:'🐱',robot:'🤖',frog:'🐸',owl:'🦉',dragon:'🐲'};
const items = {star:'⭐',gem:'💎',apple:'🍎',moon:'🌙',cookie:'🍪',flower:'🌼'};
const themes = {space:['#141743','#283e81'],forest:['#103a32','#285d42'],candy:['#722873','#a54388'],ocean:['#092f63','#11668c'],sunset:['#562054','#ad4e54']};
const tokens = [items[spec.item], '🍇', '🦋', '🍄', '🐚', '🍉'];
const directions = {ArrowLeft:[-1,0], ArrowUp:[0,-1], ArrowDown:[0,1], ArrowRight:[1,0]};
const help = {
  catch:'Move left and right to catch the falling treasures.',
  maze:'Follow the path to the treasure. Use the arrow keys or buttons.',
  runner:'Jump over the blocks. Press Space or tap Jump.',
  memory:'Flip cards to find six matching pairs.'
};
document.querySelector('#help').textContent = help[spec.template];
let active, score, heroX, heroY, falling, obstacle, jump, tile, deck, revealed, matched, wait, last, keys = new Set();
let previousAnnouncement = '';
const maze = [
  '..........',
  '.##.#.###.',
  '....#...#.',
  '.####.#.#.',
  '......#...',
  '.###....#.'
];
function announce(message) {
  if (message !== previousAnnouncement) {
    status.textContent = message;
    previousAnnouncement = message;
  }
}
function reset() {
  active = true;
  start.textContent = 'Restart';
  score = 0;
  heroX = 300;
  heroY = 306;
  falling = {x:300,y:-20};
  obstacle = 650;
  jump = 0;
  tile = [0,0];
  deck = [...tokens, ...tokens];
  for(let i = deck.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [deck[i],deck[j]] = [deck[j],deck[i]];
  }
  revealed = [];
  matched = new Set();
  wait = 0;
  keys.clear();
  announce(spec.template === 'memory' ? 'Find the first pair!' : 'Let’s go!');
}
function finish(message) {
  active = false;
  start.textContent = 'Play again';
  announce(message);
}
function move(key) {
  if (!active) return;
  if (spec.template === 'runner' && key === ' ' && jump === 0) jump = .86;
  if (spec.template !== 'maze' || !directions[key]) return;
  const [dx,dy] = directions[key];
  const x = tile[0] + dx, y = tile[1] + dy;
  if (maze[y]?.[x] !== '.') return;
  tile = [x,y];
  if (x === 9 && y === 5) finish('You found the treasure!');
}
function flip(i) {
  if (!active || wait || matched.has(i) || revealed.includes(i) || !deck[i]) return;
  revealed.push(i);
  if (revealed.length < 2) {
    announce('Card ' + (i + 1) + ' shows ' + deck[i]);
    return;
  }
  if (deck[revealed[0]] === deck[i]) {
    revealed.forEach(card => matched.add(card));
    revealed = [];
    if (matched.size === deck.length) finish('You matched every pair!');
    else announce(matched.size / 2 + ' of 6 pairs found');
  } else {
    wait = .75;
    announce('Different cards. Try another pair.');
  }
}
function text(value,x,y,size=35) {
  ctx.fillStyle = '#fff';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.font = 'bold ' + size + 'px system-ui';
  ctx.fillText(value,x,y);
}
function step(dt) {
  if (!active) return;
  if (spec.template === 'catch') {
    heroX = Math.max(28, Math.min(572, heroX + ((keys.has('ArrowRight') ? 1 : 0) - (keys.has('ArrowLeft') ? 1 : 0)) * dt * 340));
    falling.y += dt * (93 + spec.speed * 27);
    if (falling.y > 285 && falling.y < 329 && Math.abs(falling.x - heroX) < 47) {
      score++;
      announce(score + ' of ' + spec.goal + ' treasures');
      falling.y = 370;
      if (score >= spec.goal) finish('You caught all the treasures!');
    }
    if (falling.y > 355) falling = {x:35+Math.random()*530,y:-20};
  }
  if (spec.template === 'runner') {
    jump = Math.max(0, jump - dt);
    heroY = 306 - Math.sin(jump / .86 * Math.PI) * 125;
    obstacle -= dt * (170 + spec.speed * 25);
    if (obstacle > 70 && obstacle < 123 && heroY > 275) finish('You cleared ' + score + ' blocks. Try again!');
    if (obstacle < -30) {
      obstacle = 630;
      score++;
      announce(score + ' of ' + spec.goal + ' blocks');
      if (score >= spec.goal) finish('You cleared every block!');
    }
  }
  if (wait) {
    wait = Math.max(0, wait - dt);
    if (!wait) revealed = [];
  }
}
function draw(now) {
  step(Math.min(.04,(now - (last || now))/1000));
  last = now;
  ctx.fillStyle = themes[spec.theme][0];
  ctx.fillRect(0,0,600,360);
  ctx.fillStyle = themes[spec.theme][1];
  for (let x=28;x<600;x+=56) for(let y=24;y<340;y+=55) ctx.fillRect(x,y,3,3);
  if (!deck) {
    text(heroes[spec.hero],300,145,80);
    text('Press Start to play',300,237,28);
  } else if (spec.template === 'catch') {
    text(heroes[spec.hero],heroX,306,48);
    text(items[spec.item],falling.x,falling.y);
  } else if (spec.template === 'runner') {
    ctx.fillRect(0,327,600,33);
    ctx.fillStyle = '#ffba71';
    ctx.fillRect(obstacle,289,29,39);
    text(heroes[spec.hero],96,heroY,47);
  } else if (spec.template === 'maze') {
    maze.forEach((row,y) => [...row].forEach((cell,x) => {
      if (cell === '#') ctx.fillRect(x*60+3,y*60+3,54,54);
    }));
    text(items[spec.item],570,330);
    text(heroes[spec.hero],tile[0]*60+30,tile[1]*60+30);
  } else {
    deck.forEach((symbol,i) => {
      const x = (i%4)*150, y = Math.floor(i/4)*120;
      ctx.fillStyle = matched.has(i) ? '#337666' : themes[spec.theme][1];
      ctx.fillRect(x+8,y+8,134,104);
      text(matched.has(i) || revealed.includes(i) ? symbol : '?',x+75,y+60,47);
    });
  }
  requestAnimationFrame(draw);
}
function button(label, action) {
  const node = document.createElement('button');
  node.textContent = label;
  node.addEventListener('click',action);
  controls.append(node);
  return node;
}
function position(e) {
  const rect = board.getBoundingClientRect();
  return [(e.clientX-rect.left)*600/rect.width,(e.clientY-rect.top)*360/rect.height];
}
function boardAction(e) {
  const [x,y] = position(e);
  if (!active) return;
  if (spec.template === 'catch') heroX = x;
  if (spec.template === 'runner') move(' ');
  if (spec.template === 'memory' && e.type === 'pointerdown') flip(Math.floor(x/150) + Math.floor(y/120)*4);
}
start.addEventListener('click',reset);
board.addEventListener('pointerdown',boardAction);
if (spec.template === 'catch') {
  board.addEventListener('pointermove',boardAction);
  for (const [label,key] of [['Left','ArrowLeft'],['Right','ArrowRight']]) {
    const node = button(label, () => {if(active) heroX = Math.max(28, Math.min(572,heroX + directions[key][0]*42));});
    node.addEventListener('pointerdown', e => {node.setPointerCapture(e.pointerId);keys.add(key);});
    node.addEventListener('lostpointercapture', () => keys.delete(key));
  }
}
if (spec.template === 'runner') button('Jump',() => move(' '));
if (spec.template === 'maze') for (const key of Object.keys(directions)) button(key.slice(5),() => move(key));
if (spec.template === 'memory') for (let i=0;i<12;i++) button('Card '+(i+1),() => flip(i));
window.addEventListener('keydown', e => {
  if (!active || (!directions[e.key] && !(e.key === ' ' && spec.template === 'runner' && e.target.tagName !== 'BUTTON'))) return;
  e.preventDefault();
  keys.add(e.key);
  if (!e.repeat) move(e.key);
});
window.addEventListener('keyup',e => keys.delete(e.key));
window.addEventListener('blur',() => keys.clear());
requestAnimationFrame(draw);
