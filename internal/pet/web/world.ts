import * as T from 'three';
import type {Item, World} from './ui';

const palettes = [
  ['#f9d392','#f2a754','#809f62','#356b52'],
  ['#9aaac4','#fbeca9','#626d96','#333f6d'],
  ['#f4cecd','#ffedb5','#96b870','#4c845c'],
  ['#b2d4d3','#fbdf9d','#6aab9d','#356d74'],
];
const stoneColors = ['#e1b07e','#91b0b5','#bf9fbd','#c8cd83'];

export function paint(ctx: CanvasRenderingContext2D, item: Item, width: number, height: number) {
  const [sky, sun, hill, grass] = palettes[item.style];
  ctx.fillStyle = sky;
  ctx.fillRect(0,0,width,height);
  const oval = (fill:string,x:number,y:number,rx:number,ry:number) => {
    ctx.fillStyle=fill;
    ctx.beginPath();
    ctx.ellipse(x*width,y*height,rx*width,ry*height,0,0,2*Math.PI);
    ctx.fill();
  };
  oval(sun,item.style===1?.26:.70,.25,item.style===1?.075:.10,.14);
  oval(hill,.24,.83,.68,.36);
  oval(grass,.80,1.02,.68,.48);
  if (item.kind==='build') {
    const levels = (item.level||1)+2;
    if(item.style===1) {
      for(let i=0;i<8;i++) {
        const a=Math.PI+i*Math.PI/7;
        oval(stoneColors[1],.5+Math.cos(a)*.25,.75+Math.sin(a)*.47,.095,.12);
      }
    } else if(item.style===2) {
      ctx.fillStyle=stoneColors[0];
      ctx.fillRect(width*.36,height*(.75-levels*.095),width*.28,height*levels*.095);
      ctx.fillStyle=stoneColors[2];
      ctx.fillRect(width*.31,height*(.72-levels*.095),width*.38,height*.12);
    } else {
      for(let i=0;i<levels;i++) {
        oval(stoneColors[item.style],.52+Math.sin(i*1.9)*.03,.76-i*.10,.25-i*.032,.072);
      }
    }
    return;
  }
  if(item.style===1 || item.style===3) {
    oval(item.style===1?'#839ab0':'#b4d9c7',.46,.76,.34,.13);
    for(let i=0;i<4;i++) oval(sun,.34+i*.037,.71+i*.028,.058,.006);
    if(item.style===1) {
      for(let i=0;i<12;i++) oval('#faf0cf',.13+(i*.173)%0.78,.08+(i*.267)%0.35,.007,.009);
    } else {
      for(let i=0;i<6;i++) oval('#e0eee2',.08+i*.11,.33+Math.sin(i*.7)*.05,.13,.041);
    }
  } else {
    for(let i=0;i<17;i++) {
      const x=.06+(i*.369)%0.87,y=.63+(i*.173)%0.31;
      if(item.style===2) {
        oval(grass,x,y+.024,.006,.036);
        oval(i%2===0?'#d76f8d':sun,x,y,.02,.022);
      } else {
        oval(i%3===0?sun:hill,x,y,.012,.015);
      }
    }
  }
}

export function renderHabitat(
  root: HTMLDivElement, state: () => World, trackingChanged: (tracking:boolean) => void, failed: () => void,
) {
  const renderer = new T.WebGLRenderer({antialias:true, alpha:false, powerPreference:'low-power'});
  renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));
  renderer.outputColorSpace = T.SRGBColorSpace;
  renderer.toneMapping = T.ACESFilmicToneMapping;
  root.append(renderer.domElement);
  const canvas = renderer.domElement;
  canvas.tabIndex=0;
  canvas.setAttribute('aria-label','Clover Hollow 3D view. Focus to walk with WASD or arrow keys.');
  const scene = new T.Scene();
  const sky = new T.Color('#c7e9ce');
  const fog = new T.Fog(sky,25,110);
  scene.background=sky;
  scene.fog=fog;
  const camera = new T.PerspectiveCamera(55,1,.07,150);
  camera.rotation.order='YXZ';

  const hemisphere = new T.HemisphereLight('#ecfce8','#837b57',2.0);
  const sunlight = new T.DirectionalLight('#ffebba',2.7);
  sunlight.position.set(-11,24,12);
  scene.add(hemisphere,sunlight);
  const skyTop=new T.Color('#80c2e2');
  const domeGeo=new T.SphereGeometry(127,24,12);
  const domeMat=new T.ShaderMaterial({
    side:T.BackSide,depthWrite:false,
    uniforms:{upper:{value:skyTop},horizon:{value:sky}},
    vertexShader:`varying float height;
      void main(){height=position.y/127.0;gl_Position=projectionMatrix*modelViewMatrix*vec4(position,1.0);}`,
    fragmentShader:`uniform vec3 upper;uniform vec3 horizon;varying float height;
      void main(){gl_FragColor=vec4(mix(horizon,upper,smoothstep(-0.06,0.58,height)),1.0);
      #include <colorspace_fragment>
      }`,
  });
  const dome=new T.Mesh(domeGeo,domeMat);
  dome.renderOrder=-20;
  scene.add(dome);

  const geometries = {
    globe:new T.IcosahedronGeometry(1,2), rock:new T.IcosahedronGeometry(1,0),
    branch:new T.CylinderGeometry(.55,1,1,8), block:new T.BoxGeometry(1,1,1),
    pine:new T.ConeGeometry(1,1,8), plane:new T.CircleGeometry(1,48),
  };
  const materials = new Map<string,T.Material>();
  const textures = new Set<T.Texture>();
  const mesh = (
    type:keyof typeof geometries, color:string, x:number,y:number,z:number,
    sx:number,sy:number,sz:number,parent:T.Object3D=scene,
  ) => {
    let mat = materials.get(color);
    if(!mat) {
      mat = new T.MeshStandardMaterial({color,roughness:.92,flatShading:type==='rock'||type==='pine'});
      materials.set(color,mat);
    }
    const result = new T.Mesh(geometries[type],mat);
    result.position.set(x,y,z);
    result.scale.set(sx,sy,sz);
    parent.add(result);
    return result;
  };
  const disc = (color:string,x:number,y:number,z:number,rx:number,rz:number,parent:T.Object3D=scene) => {
    const ground = mesh('plane',color,x,y,z,rx,rz,1,parent);
    ground.rotation.x=-Math.PI/2;
    return ground;
  };
  const shade = (x:number,z:number,rx:number,rz:number,parent:T.Object3D=scene) => disc('#526f35',x,.016,z,rx,rz,parent);
  let seed = Number(BigInt(state().seed) % 2147483647n)||7;
  const random = () => ((seed = (seed*48271)%2147483647)/2147483647);

  mesh('block','#7fa64c',0,-.49,0,200,.95,200);
  for(let i=0;i<23;i++) {
    const t=i*2.399;
    const r=20+random()*39;
    disc(i%2?'#84ab55':'#79a04c',Math.cos(t)*r,.001+i*.00006,Math.sin(t)*r,3+random()*7,2+random()*5);
  }
  // Stepping stones lead from the art clearing to the burrow.
  for(let i=0;i<21;i++) {
    const t=i/20;
    disc(i%3?'#b2b275':'#c0bb85',.8+8*t+Math.sin(t*7)*2,.009+i*.0001,9-17*t,.54+random()*.29,.72);
  }

  const tree = (x:number,z:number,s:number,pine:boolean) => {
    shade(x+.35,z+.25,1.55*s,.95*s);
    mesh('branch','#66553d',x,1.4*s,z,.23*s,2.8*s,.22*s);
    if(pine) {
      const colors=['#356749','#427754','#5b8b55'];
      for(let i=0;i<3;i++) {
        const scale=(2.1-i*.46)*s;
        const crown=mesh('pine',colors[i],x,(2.8+i*.93)*s,z,scale,2.65*s,scale);
        crown.rotation.y=i*1.12;
      }
      return;
    }
    mesh('globe','#4b8247',x,3.3*s,z,1.8*s,1.7*s,1.72*s);
    mesh('globe','#659851',x-.87*s,3.7*s,z-.25*s,1.18*s,1.12*s,1.17*s);
    mesh('globe','#87ad5b',x+.43*s,4.28*s,z+.12*s,1.15*s,1.2*s,1.02*s);
  };
  for(let i=0;i<52;i++) {
    const angle=(i+.45)*Math.PI*2/52;
    const r=23+random()*20;
    tree(Math.cos(angle)*r,Math.sin(angle)*r,.76+random()*.63,i%3!==0);
  }
  [[-10,-1,1],[-10,-13,1.2],[13,-12,1.1],[15,6,.83],[-4,-12,.74],[5,-12,.88],[-9,11,1],[-19,-1,1.15]].forEach(([x,z,s],i)=>tree(x,z,s,i%3===0));
  for(let i=0;i<20;i++) {
    const a=i*.38;
    const r=59+random()*22;
    const mountain=mesh('pine',i%2?'#769c8e':'#8eae9a',Math.cos(a)*r,2,Math.sin(a)*r,11+random()*8,15+random()*18,12);
    mountain.rotation.y=a;
  }
  const clouds = new T.Group();
  scene.add(clouds);
  for(let i=0;i<15;i++) {
    const a=i*2.399;
    const r=36+random()*40;
    const x=Math.cos(a)*r, z=Math.sin(a)*r, y=14+random()*10;
    for(let j=0;j<4;j++) mesh('globe','#f9f9de',x+j*2.15,y+Math.sin(j*1.5),z,3.4,1.2,1.9,clouds);
  }

  disc('#5a916c',-15,.012,-9,5.07,4.67);
  disc('#70c5bb',-15,.025,-9,4.48,4.03);
  const rings = new T.Group();
  scene.add(rings);
  for(let i=0;i<10;i++) {
    const a=i*2.399;
    const r=Math.sqrt(random())*3.4;
    disc('#a7dbb9',-15+Math.cos(a)*r,.031,-9+Math.sin(a)*r,.54,.045,rings);
  }
  for(let i=0;i<12;i++) {
    const a=i*.51;
    mesh('rock','#a0a38a',-15+Math.cos(a)*4.6,.12,-9+Math.sin(a)*4.25,.34,.25,.27);
  }

  const grassGeom=new T.ConeGeometry(.06,.28,3);
  const grassMat=new T.MeshStandardMaterial({color:'#7da34e'});
  const grass=new T.InstancedMesh(grassGeom,grassMat,1150);
  const blade=new T.Object3D();
  const tint=new T.Color();
  let grassIndex=0;
  for(let i=0;i<430;i++) {
    const x=(random()-.5)*51,z=(random()-.5)*44;
    if(Math.hypot(x+15,z+9)<5.2 || Math.hypot(x-9,z+8)<3.4) continue;
    for(let j=0;j<3&&grassIndex<1150;j++) {
      const scale=.45+random()*.65;
      blade.position.set(x+j*.075,.12*scale,z+(random()-.5)*.12);
      blade.scale.setScalar(scale);
      blade.rotation.set(0,random()*6.28,(j-1)*.34);
      blade.updateMatrix();
      grass.setMatrixAt(grassIndex,blade.matrix);
      grass.setColorAt(grassIndex++,tint.set(j%2?'#d9ef9e':'#b6d97b'));
    }
  }
  grass.count=grassIndex;
  scene.add(grass);

  const flowerColors = ['#e5cd62','#f6efd2','#d78992','#a696c2'];
  for(let i=0;i<130;i++) {
    const x=(random()-.5)*48,z=(random()-.5)*41;
    if(Math.hypot(x+15,z+9)<5.3 || Math.hypot(x-9,z+8)<3.4) continue;
    const size=.08+random()*.09;
    const color=flowerColors[i%flowerColors.length];
    mesh('branch','#4a7b3e',x,.10,z,.026,.24,.026);
    mesh('globe',color,x,.21,z,size,.068,size);
  }

  // Moss-roofed burrow and bed.
  shade(9,-8.9,2.9,2.2);
  mesh('globe','#af845b',9,1.15,-9,2.42,1.83,1.87);
  mesh('globe','#5d8649',9,2.26,-9.22,2.56,.99,1.92);
  mesh('globe','#364837',9,.89,-7.43,.89,.91,.33);
  mesh('globe','#645644',9,.51,-7.13,.88,.32,.31);
  mesh('globe','#dbb27f',9,.17,-7.08,.91,.18,.63);
  mesh('block','#eccb87',9.08,.26,-7.01,1.19,.14,.72);
  mesh('globe','#e8cb76',7.49,1.23,-7.73,.31,.32,.22);
  mesh('globe','#e8cb76',10.56,1.17,-7.73,.31,.32,.22);
  for(let i=0;i<7;i++) mesh('rock','#83a158',7.2+i*.57,2.72+Math.sin(i*.7)*.33,-8.82,.27,.16,.31);

  const ball = mesh('globe','#e9b05b',-6,.34,5,.34,.34,.34);
  const ballBand = mesh('globe','#d96653',-6,.34,5,.348,.08,.348);
  shade(-6,5,.46,.38);
  const board = new T.Group();
  board.position.set(7,0,3.55);
  scene.add(board);
  shade(0,0,1.1,.58,board);
  mesh('block','#92764c',-.52,.85,0,.12,1.73,.13,board).rotation.z=-.22;
  mesh('block','#92764c',.52,.85,0,.12,1.73,.13,board).rotation.z=.22;
  mesh('block','#92764c',0,1.1,.02,1.3,.11,.23,board);
  mesh('block','#e1c391',0,1.67,.04,1.32,1.05,.15,board);
  const picture = document.createElement('canvas');
  picture.width=192;
  picture.height=144;
  const pictureCtx=picture.getContext('2d')!;
  const artwork=new T.CanvasTexture(picture);
  textures.add(artwork);
  artwork.colorSpace=T.SRGBColorSpace;
  const front=mesh('block','#f8eed5',0,1.67,.125,1.19,.89,.021,board);
  const artMaterial=new T.MeshStandardMaterial({map:artwork,roughness:1});
  front.material=artMaterial;
  materials.set('artMaterial',artMaterial);
  pictureCtx.fillStyle='#f8eed5';
  pictureCtx.fillRect(0,0,192,144);
  artwork.needsUpdate=true;
  let artKey='';

  const pet = new T.Group();
  const shadow = new T.Group();
  scene.add(pet,shadow);
  shade(0,0,1.03,.69,shadow);
  const creature = new T.Group();
  pet.add(creature);
  const fur='#d97d47',cream='#fff0cc',dark='#392b27';
  const torso=mesh('globe',fur,0,.62,-.14,.66,.67,.84,creature);
  mesh('globe',cream,0,.58,.50,.46,.48,.28,creature);
  const face = new T.Group();
  face.position.set(0,1.21,.29);
  creature.add(face);
  mesh('globe',fur,0,0,0,.71,.60,.61,face);
  mesh('globe',cream,0,-.24,.40,.55,.29,.28,face);
  const ears:T.Object3D[]=[];
  const eyes:T.Object3D[]=[];
  for(const s of [-1,1]) {
    const ear=new T.Group();
    face.add(ear);
    ear.position.set(s*.48,.37,-.035);
    ear.rotation.z=s*-.25;
    mesh('globe','#b85c35',0,.23,0,.20,.40,.19,ear);
    mesh('globe','#efb58c',0,.23,.14,.12,.27,.07,ear);
    ears.push(ear);
    eyes.push(mesh('globe',dark,s*.25,.035,.55,.073,.10,.047,face));
    mesh('globe','#fff8e1',s*.25-.016,.075,.585,.023,.029,.015,face);
    mesh('globe','#eea277',s*.43,-.16,.48,.12,.065,.036,face);
  }
  mesh('globe',dark,0,-.12,.683,.105,.066,.045,face);
  const tail=new T.Group();
  tail.position.set(0,.60,-.67);
  creature.add(tail);
  mesh('globe','#b8653c',0,.24,-.37,.27,.32,.57,tail).rotation.x=.4;
  mesh('globe',cream,0,.43,-.73,.23,.24,.25,tail);
  const paws=[[-.44,.17,.36],[.44,.17,.36],[-.40,.16,-.54],[.40,.16,-.54]].map(([x,y,z])=>
    mesh('globe','#b3623b',x,y,z,.21,.24,.32,creature));
  const stars=new T.Group();
  scene.add(stars);
  const starMat=new T.MeshBasicMaterial({color:'#fff2c8'});
  materials.set('stars',starMat);
  for(let i=0;i<50;i++) {
    const a=i*2.399, y=14+random()*58;
    const star=new T.Mesh(geometries.rock,starMat);
    star.position.set(Math.cos(a)*58,y,Math.sin(a)*58);
    star.scale.setScalar(.07+random()*.13);
    stars.add(star);
  }
  const rainGeometry = new T.BufferGeometry();
  const rainPositions=new Float32Array(112*6);
  rainGeometry.setAttribute('position',new T.BufferAttribute(rainPositions,3));
  const rainMaterial=new T.LineBasicMaterial({color:'#bde4da',transparent:true,opacity:.49});
  materials.set('rain',rainMaterial);
  const rain=new T.LineSegments(rainGeometry,rainMaterial);
  rain.frustumCulled=false;
  scene.add(rain);
  const drops=Array.from({length:112},()=>[random()*25-12.5,random()*13,random()*25-12.5]);

  const objects=new Map<number,{group:T.Group; key:string}>();
  const addItem=(item:Item) => {
    const group=new T.Group();
    group.position.set(item.x,0,item.z);
    scene.add(group);
    if(item.kind==='build') {
      shade(0,0,1.28,.89,group);
      const levels=(item.level||1)+2;
      if(item.style===1) {
        for(let i=0;i<8;i++) {
          const a=Math.PI-i*Math.PI/7;
          mesh('rock',stoneColors[i%2?1:3],Math.cos(a)*.96,.24+Math.sin(a)*(.98+levels*.12),0,.34,.39,.36,group).rotation.z=a;
        }
      } else if(item.style===2) {
        for(let i=0;i<levels;i++) {
          mesh('block','#ac794e',0,.19+i*.37,0,.91-i*.08,.32,.69-i*.05,group).rotation.y=(i%2)*.19;
        }
        mesh('pine',stoneColors[2],0,.43+levels*.36,0,.61,.54,.61,group);
      } else {
        for(let i=0;i<levels;i++) {
          const s=.91-i*.13;
          const color=stoneColors[(item.style+(i===0?1:0))%4];
          mesh('rock',color,Math.sin(i*1.8)*.13,.27+i*.40,0,s,.30,s*.78,group).rotation.y=i*.65;
        }
      }
    } else {
      shade(0,0,.43,.32,group);
      const color=item.kind==='berry'?'#a74e87':item.kind==='mushroom'?'#db785d':'#edab50';
      mesh('branch','#e6d3a0',0,.21,0,.12,.38,.12,group);
      if(item.kind==='berry') {
        for(let i=0;i<4;i++) mesh('globe',color,Math.cos(i*1.7)*.15,.24+Math.sin(i)*.11,Math.sin(i*1.7)*.15,.16,.17,.15,group);
      } else {
        mesh('globe',color,0,.42,0,.28,.17,.28,group);
      }
      if(item.kind!=='mushroom') mesh('globe','#b3c66e',.13,.49,-.07,.18,.07,.09,group).rotation.z=.45;
    }
    return group;
  };
  let seenTick = -1;
  const syncItems=(w:World) => {
    const items=[...w.food,...w.creations.filter(x=>x.kind==='build')];
    const ids=new Set(items.map(item=>item.id));
    for(const [id,entry] of objects) {
      if(!ids.has(id)) { scene.remove(entry.group); objects.delete(id); }
    }
    for(const item of items) {
      const key=item.kind+item.style+item.level;
      const existing=objects.get(item.id);
      if(existing?.key===key) continue;
      if(existing) scene.remove(existing.group);
      objects.set(item.id,{group:addItem(item),key});
    }
    const painting=[...w.creations].reverse().find(x=>x.kind==='art');
    const working=w.action==='draw' && w.until>0;
    const progress=working?Math.max(.12,1-(w.until-w.tick)/16):1;
    const key=working?'draft-'+w.tick:String(painting?.id ?? '');
    if(key!==artKey) {
      pictureCtx.fillStyle='#f8eed5';
      pictureCtx.fillRect(0,0,192,144);
      if(working||painting) {
        pictureCtx.save();
        pictureCtx.beginPath();
        pictureCtx.rect(0,0,192*progress,144);
        pictureCtx.clip();
        paint(pictureCtx,working?{kind:'art',style:w.focusStyle,id:0,x:0,z:0}:painting!,192,144);
        pictureCtx.restore();
      }
      artwork.needsUpdate=true;
      artKey=key;
    }
  };

  let tracking=true, yaw=0, pitch=0, pointer:number|null=null, lastX=0, lastY=0;
  const keys=new Set<string>();
  const touchKeys=new Set<string>();
  const directions:Record<string,string>={w:'up',a:'left',s:'down',d:'right',ArrowUp:'up',ArrowLeft:'left',ArrowDown:'down',ArrowRight:'right'};
  const track=(on:boolean) => {
    if(tracking!==on) { tracking=on; trackingChanged(on); }
  };
  const keydown=(e:KeyboardEvent) => {
    const dir=directions[e.key.length===1?e.key.toLowerCase():e.key];
    if(!dir || e.altKey || e.ctrlKey || e.metaKey) return;
    e.preventDefault();
    keys.add(dir);
    track(false);
  };
  const keyup=(e:KeyboardEvent) => keys.delete(directions[e.key.length===1?e.key.toLowerCase():e.key]);
  const clear=() => { keys.clear(); touchKeys.clear(); };
  const clearKeyboard=() => keys.clear();
  const down=(e:PointerEvent) => {
    if(pointer!==null || e.button!==0) return;
    pointer=e.pointerId;
    lastX=e.clientX;
    lastY=e.clientY;
    canvas.focus({preventScroll:true});
    canvas.setPointerCapture(e.pointerId);
  };
  const look=(e:PointerEvent) => {
    if(e.pointerId!==pointer) return;
    const dx=e.clientX-lastX,dy=e.clientY-lastY;
    if(!dx&&!dy) return;
    track(false);
    yaw-=dx*.006;
    pitch=T.MathUtils.clamp(pitch-dy*.005,-1.18,1.05);
    lastX=e.clientX;
    lastY=e.clientY;
  };
  const release=(e:PointerEvent) => { if(e.pointerId===pointer) pointer=null; };
  const lose=(e:Event) => { e.preventDefault(); failed(); renderer.setAnimationLoop(null); };
  canvas.addEventListener('pointerdown',down);
  canvas.addEventListener('pointermove',look);
  canvas.addEventListener('lostpointercapture',release);
  canvas.addEventListener('keydown',keydown);
  canvas.addEventListener('keyup',keyup);
  canvas.addEventListener('blur',clearKeyboard);
  canvas.addEventListener('webglcontextlost',lose);
  window.addEventListener('blur',clear);

  const resize=new ResizeObserver(() => {
    const width=root.clientWidth,height=root.clientHeight;
    if(!width||!height) return;
    camera.aspect=width/height;
    camera.fov=camera.aspect<1?63:55;
    camera.updateProjectionMatrix();
    renderer.setSize(width,height,false);
  });
  resize.observe(root);
  const target=new T.Vector3(), cameraTarget=new T.Vector3(), lookAt=new T.Vector3(), movement=new T.Vector3();
  const dayColor=new T.Color('#d4ecd8'), nightColor=new T.Color('#202d4b'), rainColor=new T.Color('#819d9b');
  const dayTop=new T.Color('#8ec9df'), nightTop=new T.Color('#101a33'), rainTop=new T.Color('#526c80');
  const walk=(x:number,z:number,distance:number) => {
    movement.set(x,0,z);
    if(!movement.lengthSq()) return;
    movement.normalize().applyAxisAngle(T.Object3D.DEFAULT_UP,yaw).multiplyScalar(distance);
    camera.position.add(movement);
    camera.position.x=T.MathUtils.clamp(camera.position.x,-23,23);
    camera.position.z=T.MathUtils.clamp(camera.position.z,-21,22);
  };
  const first=state().position;
  pet.position.set(first.x,0,first.z);
  camera.position.set(first.x+2.6,1.74,first.z+5.6);
  let previous=0;
  let lastFrame=0;
  renderer.setAnimationLoop((ms:number) => {
    if(document.hidden || ms-lastFrame<30) return;
    lastFrame=ms;
    const dt=previous?Math.min(.08,(ms-previous)/1000):.033;
    previous=ms;
    const time=ms/1000;
    const w=state();
    if(w.tick!==seenTick) { syncItems(w); seenTick=w.tick; }
    const daylight=T.MathUtils.smoothstep(Math.sin((w.minute-345)/1440*Math.PI*2),-.18,.36);
    const wet=w.weather==='Rain';
    sky.copy(nightColor).lerp(wet?rainColor:dayColor,daylight);
    skyTop.copy(nightTop).lerp(wet?rainTop:dayTop,daylight);
    fog.color.copy(sky);
    hemisphere.intensity=.78+daylight*1.39;
    sunlight.intensity=.16+daylight*(wet?.8:2.35);
    stars.visible=daylight<.27 && !wet;
    clouds.visible=daylight>.16;
    clouds.position.x=Math.sin(time*.022)*6;
    rings.position.x=Math.sin(time*.43)*.15;

    target.set(w.position.x,0,w.position.z);
    const walking=pet.position.distanceTo(target)>.035;
    const mix=1-Math.exp(-dt*(walking?2.65:7));
    const heading=walking ? Math.atan2(target.x-pet.position.x,target.z-pet.position.z) :
      tracking ? Math.atan2(camera.position.x-pet.position.x,camera.position.z-pet.position.z) : pet.rotation.y;
    pet.rotation.y += Math.atan2(Math.sin(heading-pet.rotation.y),Math.cos(heading-pet.rotation.y))*(1-Math.exp(-dt*5));
    pet.position.lerp(target,mix);
    shadow.position.copy(pet.position);
    const sleeping=w.mood==='Dreaming';
    const active=w.until>0;
    creature.scale.y += ((sleeping?.54:1)-creature.scale.y)*(1-Math.exp(-dt*5));
    creature.position.y=sleeping?.04:Math.abs(Math.sin(time*(walking?9:3.1)))*(walking?.15:.035);
    torso.rotation.x=sleeping?.15:0;
    tail.rotation.y=Math.sin(time*(w.happiness>76?8:4))*(sleeping?.08:.40);
    face.rotation.z=Math.sin(time*2)*.08;
    face.rotation.x=active&&(w.action==='draw'||w.action==='eat')?Math.sin(time*5)*.17:0;
    const blink=sleeping || Math.sin(time*1.42)>.996;
    for(let i=0;i<2;i++) {
      eyes[i].scale.y=blink?.021:.10;
      ears[i].rotation.x=Math.sin(time*2.3+i*2)*.12;
    }
    for(let i=0;i<4;i++) paws[i].position.y=.17+(walking?Math.max(0,Math.sin(time*9+i*Math.PI/2))*.18:0);
    const playing=w.action==='play'&&active;
    ball.position.y=.34+(playing?Math.abs(Math.sin(time*5))*.75:0);
    ballBand.position.y=ball.position.y;

    if(tracking) {
      cameraTarget.set(pet.position.x+2.45,1.71,pet.position.z+5.15);
      camera.position.lerp(cameraTarget,1-Math.exp(-dt*3.3));
      lookAt.set(pet.position.x,.99,pet.position.z);
      camera.lookAt(lookAt);
      yaw=camera.rotation.y;
      pitch=camera.rotation.x;
    } else {
      camera.rotation.set(pitch,yaw,0,'YXZ');
      const pressed=(dir:string) => +(keys.has(dir)||touchKeys.has(dir));
      walk(pressed('right')-pressed('left'),pressed('down')-pressed('up'),dt*4.8);
    }
    rain.visible=wet;
    if(wet) {
      drops.forEach(([x,y,z],i) => {
        const a=x+camera.position.x,b=(y-time*8%13+13)%13,c=z+camera.position.z;
        rainPositions.set([a,b,c,a+.11,b-.57,c],i*6);
      });
      rainGeometry.attributes.position.needsUpdate=true;
    }
    dome.position.copy(camera.position);
    renderer.render(scene,camera);
  });

  return {
    follow:() => track(true),
    move:(dir:string,pressed:boolean) => {
      if(pressed) {
        if(!touchKeys.has(dir)) walk(+(dir==='right')-+(dir==='left'),+(dir==='down')-+(dir==='up'),.28);
        touchKeys.add(dir);
        track(false);
      } else touchKeys.delete(dir);
    },
    dispose:() => {
      resize.disconnect();
      renderer.setAnimationLoop(null);
      window.removeEventListener('blur',clear);
      for(const geo of Object.values(geometries)) geo.dispose();
      rainGeometry.dispose();
      grassGeom.dispose();
      grassMat.dispose();
      domeGeo.dispose();
      domeMat.dispose();
      for(const mat of materials.values()) mat.dispose();
      for(const tex of textures) tex.dispose();
      renderer.dispose();
      canvas.removeEventListener('webglcontextlost',lose);
      renderer.forceContextLoss();
      canvas.remove();
    },
  };
}
