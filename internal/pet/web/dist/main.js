var pT=Object.create;var qv=Object.defineProperty;var mT=Object.getOwnPropertyDescriptor;var gT=Object.getOwnPropertyNames;var vT=Object.getPrototypeOf,_T=Object.prototype.hasOwnProperty;var Zi=(e,t)=>()=>{try{return t||e((t={exports:{}}).exports,t),t.exports}catch(n){throw t=0,n}};var yT=(e,t,n,i)=>{if(t&&typeof t=="object"||typeof t=="function")for(let a of gT(t))!_T.call(e,a)&&a!==n&&qv(e,a,{get:()=>t[a],enumerable:!(i=mT(t,a))||i.enumerable});return e};var Al=(e,t,n)=>(n=e!=null?pT(vT(e)):{},yT(t||!e||!e.__esModule?qv(n,"default",{value:e,enumerable:!0}):n,e));var s_=Zi(kt=>{"use strict";var Sp=Symbol.for("react.transitional.element"),xT=Symbol.for("react.portal"),ST=Symbol.for("react.fragment"),bT=Symbol.for("react.strict_mode"),MT=Symbol.for("react.profiler"),TT=Symbol.for("react.consumer"),ET=Symbol.for("react.context"),AT=Symbol.for("react.forward_ref"),wT=Symbol.for("react.suspense"),CT=Symbol.for("react.memo"),Qv=Symbol.for("react.lazy"),RT=Symbol.for("react.activity"),DT=Symbol.for("react.view_transition"),Yv=Symbol.iterator;function NT(e){return e===null||typeof e!="object"?null:(e=Yv&&e[Yv]||e["@@iterator"],typeof e=="function"?e:null)}var jv={isMounted:function(){return!1},enqueueForceUpdate:function(){},enqueueReplaceState:function(){},enqueueSetState:function(){}},$v=Object.assign,t_={};function kr(e,t,n){this.props=e,this.context=t,this.refs=t_,this.updater=n||jv}kr.prototype.isReactComponent={};kr.prototype.setState=function(e,t){if(typeof e!="object"&&typeof e!="function"&&e!=null)throw Error("takes an object of state variables to update or a function which returns an object of state variables.");this.updater.enqueueSetState(this,e,t,"setState")};kr.prototype.forceUpdate=function(e){this.updater.enqueueForceUpdate(this,e,"forceUpdate")};function e_(){}e_.prototype=kr.prototype;function bp(e,t,n){this.props=e,this.context=t,this.refs=t_,this.updater=n||jv}var Mp=bp.prototype=new e_;Mp.constructor=bp;$v(Mp,kr.prototype);Mp.isPureReactComponent=!0;var Zv=Array.isArray;function xp(){}var De={H:null,A:null,T:null,S:null},n_=Object.prototype.hasOwnProperty;function Tp(e,t,n){var i=n.ref;return{$$typeof:Sp,type:e,key:t,ref:i!==void 0?i:null,props:n}}function UT(e,t){return Tp(e.type,t,e.props)}function Ep(e){return typeof e=="object"&&e!==null&&e.$$typeof===Sp}function LT(e){var t={"=":"=0",":":"=2"};return"$"+e.replace(/[=:]/g,function(n){return t[n]})}var Jv=/\/+/g;function yp(e,t){return typeof e=="object"&&e!==null&&e.key!=null?LT(""+e.key):t.toString(36)}function IT(e){switch(e.status){case"fulfilled":return e.value;case"rejected":throw e.reason;default:switch(typeof e.status=="string"?e.then(xp,xp):(e.status="pending",e.then(function(t){e.status==="pending"&&(e.status="fulfilled",e.value=t)},function(t){e.status==="pending"&&(e.status="rejected",e.reason=t)})),e.status){case"fulfilled":return e.value;case"rejected":throw e.reason}}throw e}function Gr(e,t,n,i,a){var s=typeof e;(s==="undefined"||s==="boolean")&&(e=null);var r=!1;if(e===null)r=!0;else switch(s){case"bigint":case"string":case"number":r=!0;break;case"object":switch(e.$$typeof){case Sp:case xT:r=!0;break;case Qv:return r=e._init,Gr(r(e._payload),t,n,i,a)}}if(r)return a=a(e),r=i===""?"."+yp(e,0):i,Zv(a)?(n="",r!=null&&(n=r.replace(Jv,"$&/")+"/"),Gr(a,t,n,"",function(c){return c})):a!=null&&(Ep(a)&&(a=UT(a,n+(a.key==null||e&&e.key===a.key?"":(""+a.key).replace(Jv,"$&/")+"/")+r)),t.push(a)),1;r=0;var o=i===""?".":i+":";if(Zv(e))for(var l=0;l<e.length;l++)i=e[l],s=o+yp(i,l),r+=Gr(i,t,n,s,a);else if(l=NT(e),typeof l=="function")for(e=l.call(e),l=0;!(i=e.next()).done;)i=i.value,s=o+yp(i,l++),r+=Gr(i,t,n,s,a);else if(s==="object"){if(typeof e.then=="function")return Gr(IT(e),t,n,i,a);throw t=String(e),Error("Objects are not valid as a React child (found: "+(t==="[object Object]"?"object with keys {"+Object.keys(e).join(", ")+"}":t)+"). If you meant to render a collection of children, use an array instead.")}return r}function Au(e,t,n){if(e==null)return e;var i=[],a=0;return Gr(e,i,"","",function(s){return t.call(n,s,a++)}),i}function OT(e){if(e._status===-1){var t=e._result,n=t();n.then(function(i){(e._status===0||e._status===-1)&&(e._status=1,e._result=i,n.status===void 0&&(n.status="fulfilled",n.value=i))},function(i){(e._status===0||e._status===-1)&&(e._status=2,e._result=i,n.status===void 0&&(n.status="rejected",n.reason=i))}),e._status===-1&&(e._status=0,e._result=n)}if(e._status===1)return e._result.default;throw e._result}var Kv=typeof reportError=="function"?reportError:function(e){if(typeof window=="object"&&typeof window.ErrorEvent=="function"){var t=new window.ErrorEvent("error",{bubbles:!0,cancelable:!0,message:typeof e=="object"&&e!==null&&typeof e.message=="string"?String(e.message):String(e),error:e});if(!window.dispatchEvent(t))return}else if(typeof process=="object"&&typeof process.emit=="function"){process.emit("uncaughtException",e);return}console.error(e)};function i_(e){var t=De.T,n={};n.types=t!==null?t.types:null,De.T=n;try{var i=e(),a=De.S;a!==null&&a(n,i),typeof i=="object"&&i!==null&&typeof i.then=="function"&&i.then(xp,Kv)}catch(s){Kv(s)}finally{t!==null&&n.types!==null&&(t.types=n.types),De.T=t}}function a_(e){var t=De.T;if(t!==null){var n=t.types;n===null?t.types=[e]:n.indexOf(e)===-1&&n.push(e)}else i_(a_.bind(null,e))}var PT={map:Au,forEach:function(e,t,n){Au(e,function(){t.apply(this,arguments)},n)},count:function(e){var t=0;return Au(e,function(){t++}),t},toArray:function(e){return Au(e,function(t){return t})||[]},only:function(e){if(!Ep(e))throw Error("React.Children.only expected to receive a single React element child.");return e}};kt.Activity=RT;kt.Children=PT;kt.Component=kr;kt.Fragment=ST;kt.Profiler=MT;kt.PureComponent=bp;kt.StrictMode=bT;kt.Suspense=wT;kt.ViewTransition=DT;kt.__CLIENT_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE=De;kt.__COMPILER_RUNTIME={__proto__:null,c:function(e){return De.H.useMemoCache(e)}};kt.addTransitionType=a_;kt.cache=function(e){return function(){return e.apply(null,arguments)}};kt.cacheSignal=function(){return null};kt.cloneElement=function(e,t,n){if(e==null)throw Error("The argument must be a React element, but you passed "+e+".");var i=$v({},e.props),a=e.key;if(t!=null)for(s in t.key!==void 0&&(a=""+t.key),t)!n_.call(t,s)||s==="key"||s==="__self"||s==="__source"||s==="ref"&&t.ref===void 0||(i[s]=t[s]);var s=arguments.length-2;if(s===1)i.children=n;else if(1<s){for(var r=Array(s),o=0;o<s;o++)r[o]=arguments[o+2];i.children=r}return Tp(e.type,a,i)};kt.createContext=function(e){return e={$$typeof:ET,_currentValue:e,_currentValue2:e,_threadCount:0,Provider:null,Consumer:null},e.Provider=e,e.Consumer={$$typeof:TT,_context:e},e};kt.createElement=function(e,t,n){var i,a={},s=null;if(t!=null)for(i in t.key!==void 0&&(s=""+t.key),t)n_.call(t,i)&&i!=="key"&&i!=="__self"&&i!=="__source"&&(a[i]=t[i]);var r=arguments.length-2;if(r===1)a.children=n;else if(1<r){for(var o=Array(r),l=0;l<r;l++)o[l]=arguments[l+2];a.children=o}if(e&&e.defaultProps)for(i in r=e.defaultProps,r)a[i]===void 0&&(a[i]=r[i]);return Tp(e,s,a)};kt.createRef=function(){return{current:null}};kt.forwardRef=function(e){return{$$typeof:AT,render:e}};kt.isValidElement=Ep;kt.lazy=function(e){return{$$typeof:Qv,_payload:{_status:-1,_result:e},_init:OT}};kt.memo=function(e,t){return{$$typeof:CT,type:e,compare:t===void 0?null:t}};kt.startTransition=i_;kt.unstable_useCacheRefresh=function(){return De.H.useCacheRefresh()};kt.use=function(e){return De.H.use(e)};kt.useActionState=function(e,t,n){return De.H.useActionState(e,t,n)};kt.useCallback=function(e,t){return De.H.useCallback(e,t)};kt.useContext=function(e){return De.H.useContext(e)};kt.useDebugValue=function(){};kt.useDeferredValue=function(e,t){return De.H.useDeferredValue(e,t)};kt.useEffect=function(e,t){return De.H.useEffect(e,t)};kt.useEffectEvent=function(e){return De.H.useEffectEvent(e)};kt.useId=function(){return De.H.useId()};kt.useImperativeHandle=function(e,t,n){return De.H.useImperativeHandle(e,t,n)};kt.useInsertionEffect=function(e,t){return De.H.useInsertionEffect(e,t)};kt.useLayoutEffect=function(e,t){return De.H.useLayoutEffect(e,t)};kt.useMemo=function(e,t){return De.H.useMemo(e,t)};kt.useOptimistic=function(e,t){return De.H.useOptimistic(e,t)};kt.useReducer=function(e,t,n){return De.H.useReducer(e,t,n)};kt.useRef=function(e){return De.H.useRef(e)};kt.useState=function(e){return De.H.useState(e)};kt.useSyncExternalStore=function(e,t,n){return De.H.useSyncExternalStore(e,t,n)};kt.useTransition=function(){return De.H.useTransition()};kt.version="19.3.1"});var wl=Zi((PN,r_)=>{"use strict";r_.exports=s_()});var g_=Zi(Fe=>{"use strict";function Rp(e,t){var n=e.length;e.push(t);t:for(;0<n;){var i=n-1>>>1,a=e[i];if(0<wu(a,t))e[i]=t,e[n]=a,n=i;else break t}}function Ji(e){return e.length===0?null:e[0]}function Ru(e){if(e.length===0)return null;var t=e[0],n=e.pop();if(n!==t){e[0]=n;t:for(var i=0,a=e.length,s=a>>>1;i<s;){var r=2*(i+1)-1,o=e[r],l=r+1,c=e[l];if(0>wu(o,n))l<a&&0>wu(c,o)?(e[i]=c,e[l]=n,i=l):(e[i]=o,e[r]=n,i=r);else if(l<a&&0>wu(c,n))e[i]=c,e[l]=n,i=l;else break t}}return t}function wu(e,t){var n=e.sortIndex-t.sortIndex;return n!==0?n:e.id-t.id}Fe.unstable_now=void 0;typeof performance=="object"&&typeof performance.now=="function"?(o_=performance,Fe.unstable_now=function(){return o_.now()}):(Ap=Date,l_=Ap.now(),Fe.unstable_now=function(){return Ap.now()-l_});var o_,Ap,l_,Ma=[],ts=[],BT=1,xi=null,Dn=3,Dp=!1,Cl=!1,Rl=!1,Np=!1,f_=typeof setTimeout=="function"?setTimeout:null,h_=typeof clearTimeout=="function"?clearTimeout:null,c_=typeof setImmediate<"u"?setImmediate:null;function Cu(e){for(var t=Ji(ts);t!==null;){if(t.callback===null)Ru(ts);else if(t.startTime<=e)Ru(ts),t.sortIndex=t.expirationTime,Rp(Ma,t);else break;t=Ji(ts)}}function Up(e){if(Rl=!1,Cu(e),!Cl)if(Ji(Ma)!==null)Cl=!0,Wr||(Wr=!0,Xr());else{var t=Ji(ts);t!==null&&Lp(Up,t.startTime-e)}}var Wr=!1,Dl=-1,d_=5,p_=-1;function m_(){return Np?!0:!(Fe.unstable_now()-p_<d_)}function wp(){if(Np=!1,Wr){var e=Fe.unstable_now();p_=e;var t=!0;try{t:{Cl=!1,Rl&&(Rl=!1,h_(Dl),Dl=-1),Dp=!0;var n=Dn;try{e:{for(Cu(e),xi=Ji(Ma);xi!==null&&!(xi.expirationTime>e&&m_());){var i=xi.callback;if(typeof i=="function"){xi.callback=null,Dn=xi.priorityLevel;var a=i(xi.expirationTime<=e);if(e=Fe.unstable_now(),typeof a=="function"){xi.callback=a,Cu(e),t=!0;break e}xi===Ji(Ma)&&Ru(Ma),Cu(e)}else Ru(Ma);xi=Ji(Ma)}if(xi!==null)t=!0;else{var s=Ji(ts);s!==null&&Lp(Up,s.startTime-e),t=!1}}break t}finally{xi=null,Dn=n,Dp=!1}t=void 0}}finally{t?Xr():Wr=!1}}}var Xr;typeof c_=="function"?Xr=function(){c_(wp)}:typeof MessageChannel<"u"?(Cp=new MessageChannel,u_=Cp.port2,Cp.port1.onmessage=wp,Xr=function(){u_.postMessage(null)}):Xr=function(){f_(wp,0)};var Cp,u_;function Lp(e,t){Dl=f_(function(){e(Fe.unstable_now())},t)}Fe.unstable_IdlePriority=5;Fe.unstable_ImmediatePriority=1;Fe.unstable_LowPriority=4;Fe.unstable_NormalPriority=3;Fe.unstable_Profiling=null;Fe.unstable_UserBlockingPriority=2;Fe.unstable_cancelCallback=function(e){e.callback=null};Fe.unstable_forceFrameRate=function(e){0>e||125<e?console.error("forceFrameRate takes a positive int between 0 and 125, forcing frame rates higher than 125 fps is not supported"):d_=0<e?Math.floor(1e3/e):5};Fe.unstable_getCurrentPriorityLevel=function(){return Dn};Fe.unstable_next=function(e){switch(Dn){case 1:case 2:case 3:var t=3;break;default:t=Dn}var n=Dn;Dn=t;try{return e()}finally{Dn=n}};Fe.unstable_requestPaint=function(){Np=!0};Fe.unstable_runWithPriority=function(e,t){switch(e){case 1:case 2:case 3:case 4:case 5:break;default:e=3}var n=Dn;Dn=e;try{return t()}finally{Dn=n}};Fe.unstable_scheduleCallback=function(e,t,n){var i=Fe.unstable_now();switch(typeof n=="object"&&n!==null?(n=n.delay,n=typeof n=="number"&&0<n?i+n:i):n=i,e){case 1:var a=-1;break;case 2:a=250;break;case 5:a=1073741823;break;case 4:a=1e4;break;default:a=5e3}return a=n+a,e={id:BT++,callback:t,priorityLevel:e,startTime:n,expirationTime:a,sortIndex:-1},n>i?(e.sortIndex=n,Rp(ts,e),Ji(Ma)===null&&e===Ji(ts)&&(Rl?(h_(Dl),Dl=-1):Rl=!0,Lp(Up,n-i))):(e.sortIndex=a,Rp(Ma,e),Cl||Dp||(Cl=!0,Wr||(Wr=!0,Xr()))),e};Fe.unstable_shouldYield=m_;Fe.unstable_wrapCallback=function(e){var t=Dn;return function(){var n=Dn;Dn=t;try{return e.apply(this,arguments)}finally{Dn=n}}}});var __=Zi((zN,v_)=>{"use strict";v_.exports=g_()});var S_=Zi(Nn=>{"use strict";var zT=wl();function x_(e){var t="https://react.dev/errors/"+e;if(1<arguments.length){t+="?args[]="+encodeURIComponent(arguments[1]);for(var n=2;n<arguments.length;n++)t+="&args[]="+encodeURIComponent(arguments[n])}return"Minified React error #"+e+"; visit "+t+" for the full message or use the non-minified dev environment for full errors and additional helpful warnings."}function es(){}var Vn={d:{f:es,r:function(){throw Error(x_(522))},D:es,C:es,L:es,m:es,X:es,S:es,M:es},p:0,findDOMNode:null},FT=Symbol.for("react.portal"),VT=Symbol.for("react.recoverable"),y_=Symbol.for("react.optimistic_key");function HT(e,t,n){var i=3<arguments.length&&arguments[3]!==void 0?arguments[3]:null;return{$$typeof:FT,key:i==null?null:i===y_?y_:""+i,children:e,containerInfo:t,implementation:n}}var Nl=zT.__CLIENT_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE;function Du(e,t){if(e==="font")return"";if(typeof t=="string")return t==="use-credentials"?t:""}Nn.__DOM_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE=Vn;Nn.browser=function(e){return{$$typeof:VT,_reason:e}};Nn.createPortal=function(e,t){var n=2<arguments.length&&arguments[2]!==void 0?arguments[2]:null;if(!t||t.nodeType!==1&&t.nodeType!==9&&t.nodeType!==11)throw Error(x_(299));return HT(e,t,null,n)};Nn.flushSync=function(e){var t=Nl.T,n=Vn.p;try{if(Nl.T=null,Vn.p=2,e)return e()}finally{Nl.T=t,Vn.p=n,Vn.d.f()}};Nn.preconnect=function(e,t){typeof e=="string"&&(t?(t=t.crossOrigin,t=typeof t=="string"?t==="use-credentials"?t:"":void 0):t=null,Vn.d.C(e,t))};Nn.prefetchDNS=function(e){typeof e=="string"&&Vn.d.D(e)};Nn.preinit=function(e,t){if(typeof e=="string"&&t&&typeof t.as=="string"){var n=t.as,i=Du(n,t.crossOrigin),a=typeof t.integrity=="string"?t.integrity:void 0,s=typeof t.fetchPriority=="string"?t.fetchPriority:void 0;n==="style"?Vn.d.S(e,typeof t.precedence=="string"?t.precedence:void 0,{crossOrigin:i,integrity:a,fetchPriority:s}):n==="script"&&Vn.d.X(e,{crossOrigin:i,integrity:a,fetchPriority:s,nonce:typeof t.nonce=="string"?t.nonce:void 0})}};Nn.preinitModule=function(e,t){if(typeof e=="string")if(typeof t=="object"&&t!==null){if(t.as==null||t.as==="script"){var n=Du(t.as,t.crossOrigin);Vn.d.M(e,{crossOrigin:n,integrity:typeof t.integrity=="string"?t.integrity:void 0,nonce:typeof t.nonce=="string"?t.nonce:void 0,fetchPriority:typeof t.fetchPriority=="string"?t.fetchPriority:void 0})}}else t==null&&Vn.d.M(e)};Nn.preload=function(e,t){if(typeof e=="string"&&typeof t=="object"&&t!==null&&typeof t.as=="string"){var n=t.as,i=Du(n,t.crossOrigin);Vn.d.L(e,n,{crossOrigin:i,integrity:typeof t.integrity=="string"?t.integrity:void 0,nonce:typeof t.nonce=="string"?t.nonce:void 0,type:typeof t.type=="string"?t.type:void 0,fetchPriority:typeof t.fetchPriority=="string"?t.fetchPriority:void 0,referrerPolicy:typeof t.referrerPolicy=="string"?t.referrerPolicy:void 0,imageSrcSet:typeof t.imageSrcSet=="string"?t.imageSrcSet:void 0,imageSizes:typeof t.imageSizes=="string"?t.imageSizes:void 0,media:typeof t.media=="string"?t.media:void 0})}};Nn.preloadModule=function(e,t){if(typeof e=="string")if(t){var n=Du(t.as,t.crossOrigin);Vn.d.m(e,{as:typeof t.as=="string"&&t.as!=="script"?t.as:void 0,crossOrigin:n,integrity:typeof t.integrity=="string"?t.integrity:void 0,nonce:typeof t.nonce=="string"?t.nonce:void 0,fetchPriority:typeof t.fetchPriority=="string"?t.fetchPriority:void 0})}else Vn.d.m(e)};Nn.requestFormReset=function(e){Vn.d.r(e)};Nn.unstable_batchedUpdates=function(e,t){return e(t)};Nn.useFormState=function(e,t,n){return Nl.H.useFormState(e,t,n)};Nn.useFormStatus=function(){return Nl.H.useHostTransitionStatus()};Nn.version="19.3.1"});var T_=Zi((VN,M_)=>{"use strict";function b_(){if(!(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__>"u"||typeof __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE!="function"))try{__REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE(b_)}catch(e){console.error(e)}}b_(),M_.exports=S_()});var uM=Zi(hh=>{"use strict";var an=__(),cx=wl(),GT=T_();function J(e){var t="https://react.dev/errors/"+e;if(1<arguments.length){t+="?args[]="+encodeURIComponent(arguments[1]);for(var n=2;n<arguments.length;n++)t+="&args[]="+encodeURIComponent(arguments[n])}return"Minified React error #"+e+"; visit "+t+" for the full message or use the non-minified dev environment for full errors and additional helpful warnings."}function ux(e){return!(!e||e.nodeType!==1&&e.nodeType!==9&&e.nodeType!==11)}function _c(e){for(var t=e,n=t;n&&!n.alternate;)t=n,(t.flags&4098)!==0&&(e=t.return),n=t.return;for(;t.return;)t=t.return;return t.tag===3?e:null}function fx(e){if(e.tag===13){var t=e.memoizedState;if(t===null&&(e=e.alternate,e!==null&&(t=e.memoizedState)),t!==null)return t.dehydrated}return null}function hx(e){if(e.tag===31){var t=e.memoizedState;if(t===null&&(e=e.alternate,e!==null&&(t=e.memoizedState)),t!==null)return t.dehydrated}return null}function E_(e){if(_c(e)!==e)throw Error(J(188))}function kT(e){var t=e.alternate;if(!t){if(t=_c(e),t===null)throw Error(J(188));return t!==e?null:e}for(var n=e,i=t;;){var a=n.return;if(a===null)break;var s=a.alternate;if(s===null){if(i=a.return,i!==null){n=i;continue}break}if(a.child===s.child){for(s=a.child;s;){if(s===n)return E_(a),e;if(s===i)return E_(a),t;s=s.sibling}throw Error(J(188))}if(n.return!==i.return)n=a,i=s;else{for(var r=!1,o=a.child;o;){if(o===n){r=!0,n=a,i=s;break}if(o===i){r=!0,i=a,n=s;break}o=o.sibling}if(!r){for(o=s.child;o;){if(o===n){r=!0,n=s,i=a;break}if(o===i){r=!0,i=s,n=a;break}o=o.sibling}if(!r)throw Error(J(189))}}if(n.alternate!==i)throw Error(J(190))}if(n.tag!==3)throw Error(J(188));return n.stateNode.current===n?e:t}function dx(e){var t=e.tag;if(t===5||t===26||t===27||t===6)return e;for(e=e.child;e!==null;){if(t=dx(e),t!==null)return t;e=e.sibling}return null}function $n(e,t,n,i,a,s){for(;e!==null;){if((e.tag===5||e.tag===27||e.tag===6)&&n(e,i,a,s)||(e.tag!==22||e.memoizedState===null)&&(t||e.tag!==5&&e.tag!==27)&&$n(e.child,t,n,i,a,s))return!0;e=e.sibling}return!1}function vr(e){for(e=e.return;e!==null;){if(e.tag===3||e.tag===5||e.tag===27)return e;e=e.return}return null}function A_(e){var t=!1;for(e=e.return;e!==null&&(e.tag===4&&(t=!0),!(e.tag===3||e.tag===5||e.tag===27));)e=e.return;return t}function px(e){var t=[null,null],n=vr(e);return n===null||mx(t,e,n.child,{foundSelf:!1}),t}function mx(e,t,n,i){for(;n!==null;){if(n===t)i.foundSelf=!0;else if(n.tag===5||n.tag===27||n.tag===6){if(i.foundSelf)return e[1]=n,!0;e[0]=n}else if((n.tag!==22||n.memoizedState===null)&&mx(e,t,n.child,i))return!0;n=n.sibling}return!1}function nn(e){switch(e.tag){case 5:case 27:case 6:return e.stateNode;case 3:return e.stateNode.containerInfo;default:throw Error(J(559))}}var jr=null,fm=null;function XT(e,t,n){return e===n?!0:e===t?(jr=e,!0):!1}function WT(e,t,n){return e===n?(fm=e,!1):e===t?(fm!==null&&(jr=e),!0):!1}function w_(e){if(e===null)return null;do e=e===null?null:e.return;while(e&&e.tag!==5&&e.tag!==27&&e.tag!==3);return e||null}function hm(e,t,n){for(var i=0,a=e;a;a=n(a))i++;a=0;for(var s=t;s;s=n(s))a++;for(;0<i-a;)e=n(e),i--;for(;0<a-i;)t=n(t),a--;for(;i--;){if(e===t||t!==null&&e===t.alternate)return e;e=n(e),t=n(t)}return null}var we=Object.assign,qT=Symbol.for("react.element"),Nu=Symbol.for("react.transitional.element"),zl=Symbol.for("react.portal"),$r=Symbol.for("react.fragment"),gx=Symbol.for("react.strict_mode"),dm=Symbol.for("react.profiler"),vx=Symbol.for("react.consumer"),ea=Symbol.for("react.context"),bg=Symbol.for("react.forward_ref"),pm=Symbol.for("react.suspense"),mm=Symbol.for("react.suspense_list"),Mg=Symbol.for("react.memo"),ss=Symbol.for("react.lazy"),gm=Symbol.for("react.activity"),YT=Symbol.for("react.legacy_hidden"),ZT=Symbol.for("react.memo_cache_sentinel"),vm=Symbol.for("react.view_transition"),JT=Symbol.for("react.recoverable"),C_=Symbol.iterator;function Ul(e){return e===null||typeof e!="object"?null:(e=C_&&e[C_]||e["@@iterator"],typeof e=="function"?e:null)}var KT=Symbol.for("react.client.reference");function _m(e){if(e==null)return null;if(typeof e=="function")return e.$$typeof===KT?null:e.displayName||e.name||null;if(typeof e=="string")return e;switch(e){case $r:return"Fragment";case dm:return"Profiler";case gx:return"StrictMode";case pm:return"Suspense";case mm:return"SuspenseList";case gm:return"Activity";case vm:return"ViewTransition"}if(typeof e=="object")switch(e.$$typeof){case zl:return"Portal";case ea:return e.displayName||"Context";case vx:return(e._context.displayName||"Context")+".Consumer";case bg:var t=e.render;return e=e.displayName,e||(e=t.displayName||t.name||"",e=e!==""?"ForwardRef("+e+")":"ForwardRef"),e;case Mg:return t=e.displayName||null,t!==null?t:_m(e.type)||"Memo";case ss:t=e._payload,e=e._init;try{return _m(e(t))}catch{}}return null}var Fl=Array.isArray,Bt=cx.__CLIENT_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE,he=GT.__DOM_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE,ar={pending:!1,data:null,method:null,action:null},ym=[],to=-1;function la(e){return{current:e}}function Sn(e){0>to||(e.current=ym[to],ym[to]=null,to--)}function Le(e,t){to++,ym[to]=e.current,e.current=t}var sa=la(null),nc=la(null),ps=la(null),_f=la(null);function yf(e,t){switch(Le(ps,t),Le(nc,e),Le(sa,null),t.nodeType){case 9:case 11:e=(e=t.documentElement)&&(e=e.namespaceURI)?Gy(e):0;break;default:if(e=t.tagName,t=t.namespaceURI)t=Gy(t),e=Fb(t,e);else switch(e){case"svg":e=1;break;case"math":e=2;break;default:e=0}}Sn(sa),Le(sa,e)}function So(){Sn(sa),Sn(nc),Sn(ps)}function xm(e){var t=e.memoizedState;t!==null&&(No._currentValue=t.memoizedState,Le(_f,e)),t=sa.current;var n=Fb(t,e.type);t!==n&&(Le(nc,e),Le(sa,n))}function xf(e){nc.current===e&&(Sn(sa),Sn(nc)),_f.current===e&&(Sn(_f),No._currentValue=ar)}var Ip,R_;function is(e){if(Ip===void 0)try{throw Error()}catch(n){var t=n.stack.trim().match(/\n( *(at )?)/);Ip=t&&t[1]||"",R_=-1<n.stack.indexOf(`
    at`)?" (<anonymous>)":-1<n.stack.indexOf("@")?"@unknown:0:0":""}return`
`+Ip+e+R_}var Op=!1;function Pp(e,t){if(!e||Op)return"";Op=!0;var n=Error.prepareStackTrace;Error.prepareStackTrace=void 0;try{var i={DetermineComponentFrameRoot:function(){try{if(t){var p=function(){throw Error()};if(Object.defineProperty(p.prototype,"props",{set:function(){throw Error()}}),typeof Reflect=="object"&&Reflect.construct){try{Reflect.construct(p,[])}catch(_){var u=_}Reflect.construct(e,[],p)}else{try{p.call()}catch(_){u=_}p=!1;try{var d=Object.getOwnPropertyDescriptor(e.prototype,"props");Object.defineProperty(e.prototype,"props",{configurable:!0,set:function(){throw Error()}}),p=!0,new e}finally{p&&(d!==void 0?Object.defineProperty(e.prototype,"props",d):delete e.prototype.props)}}}else{try{throw Error()}catch(_){u=_}(p=e())&&typeof p.catch=="function"&&p.catch(function(){})}}catch(_){if(_&&u&&typeof _.stack=="string")return[_.stack,u.stack]}return[null,null]}};i.DetermineComponentFrameRoot.displayName="DetermineComponentFrameRoot";var a=Object.getOwnPropertyDescriptor(i.DetermineComponentFrameRoot,"name");a&&a.configurable&&Object.defineProperty(i.DetermineComponentFrameRoot,"name",{value:"DetermineComponentFrameRoot"});var s=i.DetermineComponentFrameRoot(),r=s[0],o=s[1];if(r&&o){var l=r.split(`
`),c=o.split(`
`);for(a=i=0;i<l.length&&!l[i].includes("DetermineComponentFrameRoot");)i++;for(;a<c.length&&!c[a].includes("DetermineComponentFrameRoot");)a++;if(i===l.length||a===c.length)for(i=l.length-1,a=c.length-1;1<=i&&0<=a&&l[i]!==c[a];)a--;for(;1<=i&&0<=a;i--,a--)if(l[i]!==c[a]){if(i!==1||a!==1)do if(i--,a--,0>a||l[i]!==c[a]){var f=`
`+l[i].replace(" at new "," at ");return e.displayName&&f.includes("<anonymous>")&&(f=f.replace("<anonymous>",e.displayName)),f}while(1<=i&&0<=a);break}}}finally{Op=!1,Error.prepareStackTrace=n}return(n=e?e.displayName||e.name:"")?is(n):""}function QT(e,t){switch(e.tag){case 26:case 27:case 5:return is(e.type);case 16:return is("Lazy");case 13:return e.child!==t&&t!==null?is("Suspense Fallback"):is("Suspense");case 19:return is("SuspenseList");case 0:case 15:return Pp(e.type,!1);case 11:return Pp(e.type.render,!1);case 1:return Pp(e.type,!0);case 31:return is("Activity");case 30:return is("ViewTransition");default:return""}}function D_(e){try{var t="",n=null;do t+=QT(e,n),n=e,e=e.return;while(e);return t}catch(i){return`
Error generating stack: `+i.message+`
`+i.stack}}var Sm=Object.prototype.hasOwnProperty,Tg=an.unstable_scheduleCallback,Bp=an.unstable_cancelCallback,jT=an.unstable_shouldYield,$T=an.unstable_requestPaint,li=an.unstable_now,tE=an.unstable_getCurrentPriorityLevel,_x=an.unstable_ImmediatePriority,yx=an.unstable_UserBlockingPriority,Sf=an.unstable_NormalPriority,eE=an.unstable_LowPriority,xx=an.unstable_IdlePriority,nE=an.log,iE=an.unstable_setDisableYieldValue,yc=null,ci=null;function ls(e){if(typeof nE=="function"&&iE(e),ci&&typeof ci.setStrictMode=="function")try{ci.setStrictMode(yc,e)}catch{}}var ui=Math.clz32?Math.clz32:rE,aE=Math.log,sE=Math.LN2;function rE(e){return e>>>=0,e===0?32:31-(aE(e)/sE|0)|0}var Uu=256,Lu=262144,Iu=4194304;function $s(e){var t=e&42;if(t!==0)return t;switch(e&-e){case 1:return 1;case 2:return 2;case 4:return 4;case 8:return 8;case 16:return 16;case 32:return 32;case 64:return 64;case 128:return 128;case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:return e&-e;case 262144:case 524288:case 1048576:case 2097152:return e&3932160;case 4194304:case 8388608:case 16777216:case 33554432:return e&62914560;case 67108864:return 67108864;case 134217728:return 134217728;case 268435456:return 268435456;case 536870912:return 536870912;case 1073741824:return 0;default:return e}}function Yf(e,t,n){var i=e.pendingLanes;if(i===0)return 0;var a=0,s=e.suspendedLanes,r=e.pingedLanes;e=e.warmLanes;var o=i&134217727;return o!==0?(i=o&~s,i!==0?a=$s(i):(r&=o,r!==0?a=$s(r):n||(n=o&~e,n!==0&&(a=$s(n))))):(o=i&~s,o!==0?a=$s(o):r!==0?a=$s(r):n||(n=i&~e,n!==0&&(a=$s(n)))),a===0?0:t!==0&&t!==a&&(t&s)===0&&(s=a&-a,n=t&-t,s>=n||s===32&&(n&4194048)!==0)?t:a}function xc(e,t){return(e.pendingLanes&~(e.suspendedLanes&~e.pingedLanes)&t)===0}function Sx(e,t){(t&8)!==0&&(t|=t&32);var n=e.entangledLanes;if(n!==0)for(e=e.entanglements,n&=t;0<n;){var i=31-ui(n),a=1<<i;t|=e[i],n&=~a}return t}function oE(e,t){switch(e){case 1:case 2:case 4:case 8:case 64:return t+250;case 16:case 32:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:return t+5e3;case 4194304:case 8388608:case 16777216:case 33554432:return-1;case 67108864:case 134217728:case 268435456:case 536870912:case 1073741824:return-1;default:return-1}}function bx(){var e=Iu;return Iu<<=1,(Iu&62914560)===0&&(Iu=4194304),e}function zp(e){for(var t=[],n=0;31>n;n++)t.push(e);return t}function Sc(e,t){e.pendingLanes|=t,t!==268435456&&(e.suspendedLanes=0,e.pingedLanes=0,e.warmLanes=0)}function lE(e,t,n,i,a,s){var r=e.pendingLanes;e.pendingLanes=n,e.suspendedLanes=0,e.pingedLanes=0,e.warmLanes=0,e.expiredLanes&=n,e.entangledLanes&=n,e.errorRecoveryDisabledLanes&=n,e.shellSuspendCounter=0;var o=e.entanglements,l=e.expirationTimes,c=e.hiddenUpdates;for(n=r&~n;0<n;){var f=31-ui(n),p=1<<f;o[f]=0,l[f]=-1;var u=c[f];if(u!==null)for(c[f]=null,f=0;f<u.length;f++){var d=u[f];d!==null&&(d.lane&=-536870913)}n&=~p}i!==0&&Mx(e,i,0),s!==0&&a===0&&e.tag!==0&&(e.suspendedLanes|=s&~(r&~t))}function Mx(e,t,n){e.pendingLanes|=t,e.suspendedLanes&=~t;var i=31-ui(t);e.entangledLanes|=t,e.entanglements[i]=e.entanglements[i]|1073741824|n&261930}function Tx(e,t){var n=e.entangledLanes|=t;for(e=e.entanglements;n;){var i=31-ui(n),a=1<<i;a&t|e[i]&t&&(e[i]|=t),n&=~a}}function Ex(e,t){var n=t&-t;return n=(n&42)!==0?1:Eg(n),(n&(e.suspendedLanes|t))!==0?0:n}function Eg(e){switch(e){case 2:e=1;break;case 8:e=4;break;case 32:e=16;break;case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:case 4194304:case 8388608:case 16777216:case 33554432:e=128;break;case 268435456:e=134217728;break;default:e=0}return e}function Ag(e){return e&=-e,2<e?8<e?(e&134217727)!==0?32:268435456:8:2}function Ax(){var e=he.p;return e!==0?e:(e=window.event,e===void 0?32:oM(e.type))}function N_(e,t){var n=he.p;try{return he.p=e,t()}finally{he.p=n}}var Pa=Math.random().toString(36).slice(2),yn="__reactFiber$"+Pa,ti="__reactProps$"+Pa,Io="__reactContainer$"+Pa,U_="__reactEvents$"+Pa,cE="__reactListeners$"+Pa,uE="__reactHandles$"+Pa,L_="__reactResources$"+Pa,bc="__reactMarker$"+Pa,bf="__reactLoad$"+Pa;function Zf(e){delete e[yn],delete e[ti],delete e[cE],delete e[uE]}function nr(e){var t;if(t=e[yn])return t;for(var n=e.parentNode;n;){if(t=n[Io]||n[yn]){if(n=t.alternate,t.child!==null||n!==null&&n.child!==null)for(e=Ky(e);e!==null;){if(n=e[yn])return n;e=Ky(e)}return t}e=n,n=e.parentNode}return null}function Oo(e){if(e=e[yn]||e[Io]){var t=e.tag;if(t===5||t===6||t===13||t===31||t===26||t===27||t===3)return e}return null}function Vl(e){var t=e.tag;if(t===5||t===26||t===27||t===6)return e.stateNode;throw Error(J(33))}function uo(e){var t=e[L_];return t||(t=e[L_]={hoistableStyles:new Map,hoistableScripts:new Map}),t}function fn(e){e[bc]=!0}function wx(e){e[bf]=void 0}var Cx=new Set,Rx={};function _r(e,t){bo(e,t),bo(e+"Capture",t)}function bo(e,t){for(Rx[e]=t,e=0;e<t.length;e++)Cx.add(t[e])}var fE=RegExp("^[:A-Z_a-z\\u00C0-\\u00D6\\u00D8-\\u00F6\\u00F8-\\u02FF\\u0370-\\u037D\\u037F-\\u1FFF\\u200C-\\u200D\\u2070-\\u218F\\u2C00-\\u2FEF\\u3001-\\uD7FF\\uF900-\\uFDCF\\uFDF0-\\uFFFD][:A-Z_a-z\\u00C0-\\u00D6\\u00D8-\\u00F6\\u00F8-\\u02FF\\u0370-\\u037D\\u037F-\\u1FFF\\u200C-\\u200D\\u2070-\\u218F\\u2C00-\\u2FEF\\u3001-\\uD7FF\\uF900-\\uFDCF\\uFDF0-\\uFFFD\\-.0-9\\u00B7\\u0300-\\u036F\\u203F-\\u2040]*$"),I_={},O_={};function hE(e){return Sm.call(O_,e)?!0:Sm.call(I_,e)?!1:fE.test(e)?O_[e]=!0:(I_[e]=!0,!1)}var ce=!1;function P_(){var e=ce;return ce=!1,e}function ju(e,t,n){if(hE(t))if(n===null)e.removeAttribute(t);else{switch(typeof n){case"undefined":case"function":case"symbol":e.removeAttribute(t);return;case"boolean":var i=t.toLowerCase().slice(0,5);if(i!=="data-"&&i!=="aria-"){e.removeAttribute(t);return}}e.setAttribute(t,n)}}function Ou(e,t,n){if(n===null)e.removeAttribute(t);else{switch(typeof n){case"undefined":case"function":case"symbol":case"boolean":e.removeAttribute(t);return}e.setAttribute(t,n)}}function Ta(e,t,n,i){if(i===null)e.removeAttribute(n);else{switch(typeof i){case"undefined":case"function":case"symbol":case"boolean":e.removeAttribute(n);return}e.setAttributeNS(t,n,i)}}function ai(e){switch(typeof e){case"bigint":case"boolean":case"number":case"string":case"undefined":return e;case"object":return e;default:return""}}function Dx(e){var t=e.type;return(e=e.nodeName)&&e.toLowerCase()==="input"&&(t==="checkbox"||t==="radio")}function dE(e,t,n){var i=Object.getOwnPropertyDescriptor(e.constructor.prototype,t);if(!e.hasOwnProperty(t)&&typeof i<"u"&&typeof i.get=="function"&&typeof i.set=="function"){var a=i.get,s=i.set;return Object.defineProperty(e,t,{configurable:!0,get:function(){return a.call(this)},set:function(r){n=""+r,s.call(this,r)}}),Object.defineProperty(e,t,{enumerable:i.enumerable}),{getValue:function(){return n},setValue:function(r){n=""+r},stopTracking:function(){e._valueTracker=null,delete e[t]}}}}function bm(e){if(!e._valueTracker){var t=Dx(e)?"checked":"value";e._valueTracker=dE(e,t,""+e[t])}}function Nx(e){if(!e)return!1;var t=e._valueTracker;if(!t)return!0;var n=t.getValue(),i="";return e&&(i=Dx(e)?e.checked?"true":"false":e.value),e=i,e!==n?(t.setValue(e),!0):!1}var pE=/[\n"\\]/g;function Ei(e){return e.replace(pE,function(t){return"\\"+t.charCodeAt(0).toString(16)+" "})}function Mm(e,t,n,i,a,s,r,o){e.name="",r!=null&&typeof r!="function"&&typeof r!="symbol"&&typeof r!="boolean"?e.type=r:e.removeAttribute("type"),t!=null?r==="number"?(t===0&&e.value===""||e.value!=t)&&(e.value=""+ai(t)):e.value!==""+ai(t)&&(e.value=""+ai(t)):r!=="submit"&&r!=="reset"||e.removeAttribute("value"),t!=null?r==="number"&&e.value==t?Fp(e,ai(e.value)):Fp(e,ai(t)):n!=null?Fp(e,ai(n)):i!=null&&e.removeAttribute("value"),a==null&&s!=null&&(e.defaultChecked=!!s),a!=null&&(e.checked=a&&typeof a!="function"&&typeof a!="symbol"),o!=null&&typeof o!="function"&&typeof o!="symbol"&&typeof o!="boolean"?e.name=""+ai(o):e.removeAttribute("name")}function Ux(e,t,n,i,a,s,r,o){if(s!=null&&typeof s!="function"&&typeof s!="symbol"&&typeof s!="boolean"&&(e.type=s),t!=null||n!=null){if(!(s!=="submit"&&s!=="reset"||t!=null)){bm(e);return}n=n!=null?""+ai(n):"",t=t!=null?""+ai(t):n,o||t===e.value||(e.value=t),e.defaultValue=t}i=i??a,i=typeof i!="function"&&typeof i!="symbol"&&!!i,e.checked=o?e.checked:!!i,e.defaultChecked=!!i,r!=null&&typeof r!="function"&&typeof r!="symbol"&&typeof r!="boolean"&&(e.name=r),bm(e)}function Fp(e,t){e.defaultValue!==""+t&&(e.defaultValue=""+t)}function fo(e,t,n,i){if(e=e.options,t){t={};for(var a=0;a<n.length;a++)t["$"+n[a]]=!0;for(n=0;n<e.length;n++)a=t.hasOwnProperty("$"+e[n].value),e[n].selected!==a&&(e[n].selected=a),a&&i&&(e[n].defaultSelected=!0)}else{for(n=""+ai(n),t=null,a=0;a<e.length;a++){if(e[a].value===n){e[a].selected=!0,i&&(e[a].defaultSelected=!0);return}t!==null||e[a].disabled||(t=e[a])}t!==null&&(t.selected=!0)}}function Lx(e,t,n){if(t!=null&&(t=""+ai(t),t!==e.value&&(e.value=t),n==null)){e.defaultValue!==t&&(e.defaultValue=t);return}e.defaultValue=n!=null?""+ai(n):""}function Ix(e,t,n,i){if(t==null){if(i!=null){if(n!=null)throw Error(J(92));if(Fl(i)){if(1<i.length)throw Error(J(93));i=i[0]}n=i}n==null&&(n=""),t=n}n=ai(t),e.defaultValue=n,i=e.textContent,i===n&&i!==""&&i!==null&&(e.value=i),bm(e)}function Mo(e,t){if(t){var n=e.firstChild;if(n&&n===e.lastChild&&n.nodeType===3){n.nodeValue=t;return}}e.textContent=t}var mE=new Set("animationIterationCount aspectRatio borderImageOutset borderImageSlice borderImageWidth boxFlex boxFlexGroup boxOrdinalGroup columnCount columns flex flexGrow flexPositive flexShrink flexNegative flexOrder gridArea gridRow gridRowEnd gridRowSpan gridRowStart gridColumn gridColumnEnd gridColumnSpan gridColumnStart fontWeight lineClamp lineHeight opacity order orphans scale tabSize widows zIndex zoom fillOpacity floodOpacity stopOpacity strokeDasharray strokeDashoffset strokeMiterlimit strokeOpacity strokeWidth MozAnimationIterationCount MozBoxFlex MozBoxFlexGroup MozLineClamp msAnimationIterationCount msFlex msZoom msFlexGrow msFlexNegative msFlexOrder msFlexPositive msFlexShrink msGridColumn msGridColumnSpan msGridRow msGridRowSpan WebkitAnimationIterationCount WebkitBoxFlex WebKitBoxFlexGroup WebkitBoxOrdinalGroup WebkitColumnCount WebkitColumns WebkitFlex WebkitFlexGrow WebkitFlexPositive WebkitFlexShrink WebkitLineClamp".split(" "));function B_(e,t,n){var i=t.indexOf("--")===0;n==null||typeof n=="boolean"||n===""?i?e.setProperty(t,""):t==="float"?e.cssFloat="":e[t]="":i?e.setProperty(t,n):typeof n!="number"||n===0||mE.has(t)?t==="float"?e.cssFloat=n:e[t]=(""+n).trim():e[t]=n+"px"}function Ox(e,t,n){if(t!=null&&typeof t!="object")throw Error(J(62));if(e=e.style,n!=null){for(var i in n)!n.hasOwnProperty(i)||t!=null&&t.hasOwnProperty(i)||(i.indexOf("--")===0?e.setProperty(i,""):i==="float"?e.cssFloat="":e[i]="",ce=!0);for(var a in t)i=t[a],t.hasOwnProperty(a)&&n[a]!==i&&(B_(e,a,i),ce=!0)}else for(var s in t)t.hasOwnProperty(s)&&B_(e,s,t[s])}function wg(e){if(e.indexOf("-")===-1)return!1;switch(e){case"annotation-xml":case"color-profile":case"font-face":case"font-face-src":case"font-face-uri":case"font-face-format":case"font-face-name":case"missing-glyph":return!1;default:return!0}}var gE=new Map([["acceptCharset","accept-charset"],["htmlFor","for"],["httpEquiv","http-equiv"],["crossOrigin","crossorigin"],["accentHeight","accent-height"],["alignmentBaseline","alignment-baseline"],["arabicForm","arabic-form"],["baselineShift","baseline-shift"],["capHeight","cap-height"],["clipPath","clip-path"],["clipRule","clip-rule"],["colorInterpolation","color-interpolation"],["colorInterpolationFilters","color-interpolation-filters"],["colorProfile","color-profile"],["colorRendering","color-rendering"],["dominantBaseline","dominant-baseline"],["enableBackground","enable-background"],["fillOpacity","fill-opacity"],["fillRule","fill-rule"],["floodColor","flood-color"],["floodOpacity","flood-opacity"],["fontFamily","font-family"],["fontSize","font-size"],["fontSizeAdjust","font-size-adjust"],["fontStretch","font-stretch"],["fontStyle","font-style"],["fontVariant","font-variant"],["fontWeight","font-weight"],["glyphName","glyph-name"],["glyphOrientationHorizontal","glyph-orientation-horizontal"],["glyphOrientationVertical","glyph-orientation-vertical"],["horizAdvX","horiz-adv-x"],["horizOriginX","horiz-origin-x"],["imageRendering","image-rendering"],["letterSpacing","letter-spacing"],["lightingColor","lighting-color"],["markerEnd","marker-end"],["markerMid","marker-mid"],["markerStart","marker-start"],["maskType","mask-type"],["overlinePosition","overline-position"],["overlineThickness","overline-thickness"],["paintOrder","paint-order"],["panose-1","panose-1"],["pointerEvents","pointer-events"],["renderingIntent","rendering-intent"],["shapeRendering","shape-rendering"],["stopColor","stop-color"],["stopOpacity","stop-opacity"],["strikethroughPosition","strikethrough-position"],["strikethroughThickness","strikethrough-thickness"],["strokeDasharray","stroke-dasharray"],["strokeDashoffset","stroke-dashoffset"],["strokeLinecap","stroke-linecap"],["strokeLinejoin","stroke-linejoin"],["strokeMiterlimit","stroke-miterlimit"],["strokeOpacity","stroke-opacity"],["strokeWidth","stroke-width"],["textAnchor","text-anchor"],["textDecoration","text-decoration"],["textRendering","text-rendering"],["transformOrigin","transform-origin"],["underlinePosition","underline-position"],["underlineThickness","underline-thickness"],["unicodeBidi","unicode-bidi"],["unicodeRange","unicode-range"],["unitsPerEm","units-per-em"],["vAlphabetic","v-alphabetic"],["vHanging","v-hanging"],["vIdeographic","v-ideographic"],["vMathematical","v-mathematical"],["vectorEffect","vector-effect"],["vertAdvY","vert-adv-y"],["vertOriginX","vert-origin-x"],["vertOriginY","vert-origin-y"],["wordSpacing","word-spacing"],["writingMode","writing-mode"],["xmlnsXlink","xmlns:xlink"],["xHeight","x-height"]]),vE=/^[\u0000-\u001F ]*j[\r\n\t]*a[\r\n\t]*v[\r\n\t]*a[\r\n\t]*s[\r\n\t]*c[\r\n\t]*r[\r\n\t]*i[\r\n\t]*p[\r\n\t]*t[\r\n\t]*:/i;function $u(e){return vE.test(""+e)?"javascript:throw new Error('React has blocked a javascript: URL as a security precaution.')":e}function na(){}var Tm=null;function Cg(e){return e=e.target||e.srcElement||window,e.correspondingUseElement&&(e=e.correspondingUseElement),e.nodeType===3?e.parentNode:e}var eo=null,ho=null;function z_(e){var t=Oo(e);if(t&&(e=t.stateNode)){var n=e[ti]||null;t:switch(e=t.stateNode,t.type){case"input":if(Mm(e,n.value,n.defaultValue,n.defaultValue,n.checked,n.defaultChecked,n.type,n.name),t=n.name,n.type==="radio"&&t!=null){for(n=e;n.parentNode;)n=n.parentNode;for(n=n.querySelectorAll('input[name="'+Ei(""+t)+'"][type="radio"]'),t=0;t<n.length;t++){var i=n[t];if(i!==e&&i.form===e.form){var a=i[ti]||null;if(!a)throw Error(J(90));Mm(i,a.value,a.defaultValue,a.defaultValue,a.checked,a.defaultChecked,a.type,a.name)}}for(t=0;t<n.length;t++)i=n[t],i.form===e.form&&Nx(i)}break t;case"textarea":Lx(e,n.value,n.defaultValue);break t;case"select":t=n.value,t!=null&&fo(e,!!n.multiple,t,!1)}}}var Vp=!1;function Px(e,t,n){if(Vp)return e(t,n);Vp=!0;try{var i=e(t);return i}finally{if(Vp=!1,(eo!==null||ho!==null)&&(lh(),eo&&(t=eo,e=ho,ho=eo=null,z_(t),e)))for(t=0;t<e.length;t++)z_(e[t])}}function ic(e,t){var n=e.stateNode;if(n===null)return null;var i=n[ti]||null;if(i===null)return null;n=i[t];t:switch(t){case"onClick":case"onClickCapture":case"onDoubleClick":case"onDoubleClickCapture":case"onMouseDown":case"onMouseDownCapture":case"onMouseMove":case"onMouseMoveCapture":case"onMouseUp":case"onMouseUpCapture":case"onMouseEnter":(i=!i.disabled)||(e=e.type,i=!(e==="button"||e==="input"||e==="select"||e==="textarea")),e=!i;break t;default:e=!1}if(e)return null;if(n&&typeof n!="function")throw Error(J(231,t,typeof n));return n}var Da=!(typeof window>"u"||typeof window.document>"u"||typeof window.document.createElement>"u"),Em=!1;if(Da)try{qr={},Object.defineProperty(qr,"passive",{get:function(){Em=!0}}),window.addEventListener("test",qr,qr),window.removeEventListener("test",qr,qr)}catch{Em=!1}var qr,cs=null,Rg=null,tf=null;function Bx(){if(tf)return tf;var e,t=Rg,n=t.length,i,a="value"in cs?cs.value:cs.textContent,s=a.length;for(e=0;e<n&&t[e]===a[e];e++);var r=n-e;for(i=1;i<=r&&t[n-i]===a[s-i];i++);return tf=a.slice(e,1<i?1-i:void 0)}function ef(e){var t=e.keyCode;return"charCode"in e?(e=e.charCode,e===0&&t===13&&(e=13)):e=t,e===10&&(e=13),32<=e||e===13?e:0}function Pu(){return!0}function F_(){return!1}function Xn(e){function t(n,i,a,s,r){this._reactName=n,this._targetInst=a,this.type=i,this.nativeEvent=s,this.target=r,this.currentTarget=null;for(var o in e)e.hasOwnProperty(o)&&(n=e[o],this[o]=n?n(s):s[o]);return this.isDefaultPrevented=(s.defaultPrevented!=null?s.defaultPrevented:s.returnValue===!1)?Pu:F_,this.isPropagationStopped=F_,this}return we(t.prototype,{preventDefault:function(){this.defaultPrevented=!0;var n=this.nativeEvent;n&&(n.preventDefault?n.preventDefault():typeof n.returnValue!="unknown"&&(n.returnValue=!1),this.isDefaultPrevented=Pu)},stopPropagation:function(){var n=this.nativeEvent;n&&(n.stopPropagation?n.stopPropagation():typeof n.cancelBubble!="unknown"&&(n.cancelBubble=!0),this.isPropagationStopped=Pu)},persist:function(){},isPersistent:Pu}),t}var Rs={eventPhase:0,bubbles:0,cancelable:0,timeStamp:function(e){return e.timeStamp||Date.now()},defaultPrevented:0,isTrusted:0},Jf=Xn(Rs),Mc=we({},Rs,{view:0,detail:0}),_E=Xn(Mc),Hp,Gp,Ll,Kf=we({},Mc,{screenX:0,screenY:0,clientX:0,clientY:0,pageX:0,pageY:0,ctrlKey:0,shiftKey:0,altKey:0,metaKey:0,getModifierState:Dg,button:0,buttons:0,relatedTarget:function(e){return e.relatedTarget===void 0?e.fromElement===e.srcElement?e.toElement:e.fromElement:e.relatedTarget},movementX:function(e){return"movementX"in e?e.movementX:(e!==Ll&&(Ll&&e.type==="mousemove"?(Hp=e.screenX-Ll.screenX,Gp=e.screenY-Ll.screenY):Gp=Hp=0,Ll=e),Hp)},movementY:function(e){return"movementY"in e?e.movementY:Gp}}),V_=Xn(Kf),yE=we({},Kf,{dataTransfer:0}),xE=Xn(yE),SE=we({},Mc,{relatedTarget:0}),kp=Xn(SE),bE=we({},Rs,{animationName:0,elapsedTime:0,pseudoElement:0}),ME=Xn(bE),TE=we({},Rs,{clipboardData:function(e){return"clipboardData"in e?e.clipboardData:window.clipboardData}}),EE=Xn(TE),AE=we({},Rs,{data:0}),H_=Xn(AE),wE={Esc:"Escape",Spacebar:" ",Left:"ArrowLeft",Up:"ArrowUp",Right:"ArrowRight",Down:"ArrowDown",Del:"Delete",Win:"OS",Menu:"ContextMenu",Apps:"ContextMenu",Scroll:"ScrollLock",MozPrintableKey:"Unidentified"},CE={8:"Backspace",9:"Tab",12:"Clear",13:"Enter",16:"Shift",17:"Control",18:"Alt",19:"Pause",20:"CapsLock",27:"Escape",32:" ",33:"PageUp",34:"PageDown",35:"End",36:"Home",37:"ArrowLeft",38:"ArrowUp",39:"ArrowRight",40:"ArrowDown",45:"Insert",46:"Delete",112:"F1",113:"F2",114:"F3",115:"F4",116:"F5",117:"F6",118:"F7",119:"F8",120:"F9",121:"F10",122:"F11",123:"F12",144:"NumLock",145:"ScrollLock",224:"Meta"},RE={Alt:"altKey",Control:"ctrlKey",Meta:"metaKey",Shift:"shiftKey"};function DE(e){var t=this.nativeEvent;return t.getModifierState?t.getModifierState(e):(e=RE[e])?!!t[e]:!1}function Dg(){return DE}var NE=we({},Mc,{key:function(e){if(e.key){var t=wE[e.key]||e.key;if(t!=="Unidentified")return t}return e.type==="keypress"?(e=ef(e),e===13?"Enter":String.fromCharCode(e)):e.type==="keydown"||e.type==="keyup"?CE[e.keyCode]||"Unidentified":""},code:0,location:0,ctrlKey:0,shiftKey:0,altKey:0,metaKey:0,repeat:0,locale:0,getModifierState:Dg,charCode:function(e){return e.type==="keypress"?ef(e):0},keyCode:function(e){return e.type==="keydown"||e.type==="keyup"?e.keyCode:0},which:function(e){return e.type==="keypress"?ef(e):e.type==="keydown"||e.type==="keyup"?e.keyCode:0}}),UE=Xn(NE),LE=we({},Kf,{pointerId:0,width:0,height:0,pressure:0,tangentialPressure:0,tiltX:0,tiltY:0,twist:0,pointerType:0,isPrimary:0}),G_=Xn(LE),IE=we({},Rs,{submitter:0}),OE=Xn(IE),PE=we({},Mc,{touches:0,targetTouches:0,changedTouches:0,altKey:0,metaKey:0,ctrlKey:0,shiftKey:0,getModifierState:Dg}),BE=Xn(PE),zE=we({},Rs,{propertyName:0,elapsedTime:0,pseudoElement:0}),FE=Xn(zE),VE=we({},Kf,{deltaX:function(e){return"deltaX"in e?e.deltaX:"wheelDeltaX"in e?-e.wheelDeltaX:0},deltaY:function(e){return"deltaY"in e?e.deltaY:"wheelDeltaY"in e?-e.wheelDeltaY:"wheelDelta"in e?-e.wheelDelta:0},deltaZ:0,deltaMode:0}),HE=Xn(VE),GE=we({},Rs,{newState:0,oldState:0,source:0}),kE=Xn(GE),XE=[9,13,27,32],Ng=Da&&"CompositionEvent"in window,kl=null;Da&&"documentMode"in document&&(kl=document.documentMode);var WE=Da&&"TextEvent"in window&&!kl,zx=Da&&(!Ng||kl&&8<kl&&11>=kl),k_=" ",X_=!1;function Fx(e,t){switch(e){case"keyup":return XE.indexOf(t.keyCode)!==-1;case"keydown":return t.keyCode!==229;case"keypress":case"mousedown":case"focusout":return!0;default:return!1}}function Vx(e){return e=e.detail,typeof e=="object"&&"data"in e?e.data:null}var no=!1;function qE(e,t){switch(e){case"compositionend":return Vx(t);case"keypress":return t.which!==32?null:(X_=!0,k_);case"textInput":return e=t.data,e===k_&&X_?null:e;default:return null}}function YE(e,t){if(no)return e==="compositionend"||!Ng&&Fx(e,t)?(e=Bx(),tf=Rg=cs=null,no=!1,e):null;switch(e){case"paste":return null;case"keypress":if(!(t.ctrlKey||t.altKey||t.metaKey)||t.ctrlKey&&t.altKey){if(t.char&&1<t.char.length)return t.char;if(t.which)return String.fromCharCode(t.which)}return null;case"compositionend":return zx&&t.locale!=="ko"?null:t.data;default:return null}}var ZE={color:!0,date:!0,datetime:!0,"datetime-local":!0,email:!0,month:!0,number:!0,password:!0,range:!0,search:!0,tel:!0,text:!0,time:!0,url:!0,week:!0};function W_(e){var t=e&&e.nodeName&&e.nodeName.toLowerCase();return t==="input"?!!ZE[e.type]:t==="textarea"}function Hx(e,t,n,i){eo?ho?ho.push(i):ho=[i]:eo=i,t=Xf(t,"onChange"),0<t.length&&(n=new Jf("onChange","change",null,n,i),e.push({event:n,listeners:t}))}var Xl=null,ac=null;function JE(e){Pb(e,0)}function Qf(e){var t=Vl(e);if(Nx(t))return e}function q_(e,t){if(e==="change")return t}var Gx=!1;Da&&(Da?(zu="oninput"in document,zu||(Xp=document.createElement("div"),Xp.setAttribute("oninput","return;"),zu=typeof Xp.oninput=="function"),Bu=zu):Bu=!1,Gx=Bu&&(!document.documentMode||9<document.documentMode));var Bu,zu,Xp;function Y_(){Xl&&(Xl.detachEvent("onpropertychange",kx),ac=Xl=null)}function kx(e){if(e.propertyName==="value"&&Qf(ac)){var t=[];Hx(t,ac,e,Cg(e)),Px(JE,t)}}function KE(e,t,n){e==="focusin"?(Y_(),Xl=t,ac=n,Xl.attachEvent("onpropertychange",kx)):e==="focusout"&&Y_()}function QE(e){if(e==="selectionchange"||e==="keyup"||e==="keydown")return Qf(ac)}function jE(e,t){if(e==="click")return Qf(t)}function $E(e,t){if(e==="input"||e==="change")return Qf(t)}function tA(e,t){return e===t&&(e!==0||1/e===1/t)||e!==e&&t!==t}var hi=typeof Object.is=="function"?Object.is:tA;function sc(e,t){if(hi(e,t))return!0;if(typeof e!="object"||e===null||typeof t!="object"||t===null)return!1;var n=Object.keys(e),i=Object.keys(t);if(n.length!==i.length)return!1;for(i=0;i<n.length;i++){var a=n[i];if(!Sm.call(t,a)||!hi(e[a],t[a]))return!1}return!0}function Am(e){if(e=e||(typeof document<"u"?document:void 0),typeof e>"u")return null;try{return e.activeElement||e.body}catch{return e.body}}function Z_(e){for(;e&&e.firstChild;)e=e.firstChild;return e}function J_(e,t){var n=Z_(e);e=0;for(var i;n;){if(n.nodeType===3){if(i=e+n.textContent.length,e<=t&&i>=t)return{node:n,offset:t-e};e=i}t:{for(;n;){if(n.nextSibling){n=n.nextSibling;break t}n=n.parentNode}n=void 0}n=Z_(n)}}function Xx(e,t){return e&&t?e===t?!0:e&&e.nodeType===3?!1:t&&t.nodeType===3?Xx(e,t.parentNode):"contains"in e?e.contains(t):e.compareDocumentPosition?!!(e.compareDocumentPosition(t)&16):!1:!1}function Wx(e){e=e!=null&&e.ownerDocument!=null&&e.ownerDocument.defaultView!=null?e.ownerDocument.defaultView:window;for(var t=Am(e.document);t instanceof e.HTMLIFrameElement;){try{var n=typeof t.contentWindow.location.href=="string"}catch{n=!1}if(n)e=t.contentWindow;else break;t=Am(e.document)}return t}function Ug(e){var t=e&&e.nodeName&&e.nodeName.toLowerCase();return t&&(t==="input"&&(e.type==="text"||e.type==="search"||e.type==="tel"||e.type==="url"||e.type==="password")||t==="textarea"||e.contentEditable==="true")}var eA=Da&&"documentMode"in document&&11>=document.documentMode,io=null,wm=null,Wl=null,Cm=!1;function K_(e,t,n){var i=n.window===n?n.document:n.nodeType===9?n:n.ownerDocument;Cm||io==null||io!==Am(i)||(i=io,"selectionStart"in i&&Ug(i)?i={start:i.selectionStart,end:i.selectionEnd}:(i=(i.ownerDocument&&i.ownerDocument.defaultView||window).getSelection(),i={anchorNode:i.anchorNode,anchorOffset:i.anchorOffset,focusNode:i.focusNode,focusOffset:i.focusOffset}),Wl&&sc(Wl,i)||(Wl=i,i=Xf(wm,"onSelect"),0<i.length&&(t=new Jf("onSelect","select",null,t,n),e.push({event:t,listeners:i}),t.target=io)))}function Qs(e,t){var n={};return n[e.toLowerCase()]=t.toLowerCase(),n["Webkit"+e]="webkit"+t,n["Moz"+e]="moz"+t,n}var ao={animationend:Qs("Animation","AnimationEnd"),animationiteration:Qs("Animation","AnimationIteration"),animationstart:Qs("Animation","AnimationStart"),transitionrun:Qs("Transition","TransitionRun"),transitionstart:Qs("Transition","TransitionStart"),transitioncancel:Qs("Transition","TransitionCancel"),transitionend:Qs("Transition","TransitionEnd")},Wp={},qx={};Da&&(qx=document.createElement("div").style,"AnimationEvent"in window||(delete ao.animationend.animation,delete ao.animationiteration.animation,delete ao.animationstart.animation),"TransitionEvent"in window||delete ao.transitionend.transition);function yr(e){if(Wp[e])return Wp[e];if(!ao[e])return e;var t=ao[e],n;for(n in t)if(t.hasOwnProperty(n)&&n in qx)return Wp[e]=t[n];return e}var Yx=yr("animationend"),Zx=yr("animationiteration"),Jx=yr("animationstart"),nA=yr("transitionrun"),iA=yr("transitionstart"),aA=yr("transitioncancel"),Kx=yr("transitionend"),Qx=new Map,Rm="abort auxClick beforeToggle cancel canPlay canPlayThrough click close contextMenu copy cut drag dragEnd dragEnter dragExit dragLeave dragOver dragStart drop durationChange emptied encrypted ended error fullscreenChange fullscreenError gotPointerCapture input invalid keyDown keyPress keyUp load loadedData loadedMetadata loadStart lostPointerCapture mouseDown mouseMove mouseOut mouseOver mouseUp paste pause play playing pointerCancel pointerDown pointerMove pointerOut pointerOver pointerUp progress rateChange reset resize seeked seeking stalled submit suspend timeUpdate touchCancel touchEnd touchStart volumeChange scroll toggle touchMove waiting wheel".split(" ");Rm.push("scrollEnd");function Fi(e,t){Qx.set(e,t),_r(t,[e])}var sA=0;function Na(e,t){if(e.name!=null&&e.name!=="auto")return e.name;if(t.autoName!==null)return t.autoName;e=zi.identifierPrefix;var n=sA++;return e="_"+e+"t_"+n.toString(32)+"_",t.autoName=e}function Q_(e){if(e==null||typeof e=="string")return e;var t=null,n=xo;if(n!==null)for(var i=0;i<n.length;i++){var a=e[n[i]];if(a!=null){if(a==="none")return"none";t=t==null?a:t+(" "+a)}}return t??e.default}function Ba(e,t){return e=Q_(e),t=Q_(t),t==null?e==="auto"?null:e:t==="auto"?null:t}var Mf=typeof reportError=="function"?reportError:function(e){if(typeof window=="object"&&typeof window.ErrorEvent=="function"){var t=new window.ErrorEvent("error",{bubbles:!0,cancelable:!0,message:typeof e=="object"&&e!==null&&typeof e.message=="string"?String(e.message):String(e),error:e});if(!window.dispatchEvent(t))return}else if(typeof process=="object"&&typeof process.emit=="function"){process.emit("uncaughtException",e);return}console.error(e)},bi=[],so=0,Lg=0;function jf(){for(var e=so,t=Lg=so=0;t<e;){var n=bi[t];bi[t++]=null;var i=bi[t];bi[t++]=null;var a=bi[t];bi[t++]=null;var s=bi[t];if(bi[t++]=null,i!==null&&a!==null){var r=i.pending;r===null?a.next=a:(a.next=r.next,r.next=a),i.pending=a}s!==0&&jx(n,a,s)}}function $f(e,t,n,i){bi[so++]=e,bi[so++]=t,bi[so++]=n,bi[so++]=i,Lg|=i,e.lanes|=i,e=e.alternate,e!==null&&(e.lanes|=i)}function Ig(e,t,n,i){return $f(e,t,n,i),Tf(e)}function xr(e,t){return $f(e,null,null,t),Tf(e)}function jx(e,t,n){e.lanes|=n;var i=e.alternate;i!==null&&(i.lanes|=n);for(var a=!1,s=e.return;s!==null;)s.childLanes|=n,i=s.alternate,i!==null&&(i.childLanes|=n),s.tag===22&&(e=s.stateNode,e===null||e._visibility&1||(a=!0)),e=s,s=s.return;return e.tag===3?(s=e.stateNode,a&&t!==null&&(a=31-ui(n),e=s.hiddenUpdates,i=e[a],i===null?e[a]=[t]:i.push(t),t.lane=n|536870912),s):null}function Tf(e){if(50<ec)throw ec=0,hf=null,Error(J(185));for(var t=e.return;t!==null;)e=t,t=e.return;return e.tag===3?e.stateNode:null}var ro={};function rA(e,t,n,i){this.tag=e,this.key=n,this.sibling=this.child=this.return=this.stateNode=this.type=this.elementType=null,this.index=0,this.refCleanup=this.ref=null,this.pendingProps=t,this.dependencies=this.memoizedState=this.updateQueue=this.memoizedProps=null,this.mode=i,this.subtreeFlags=this.flags=0,this.deletions=null,this.childLanes=this.lanes=0,this.alternate=null}function Qn(e,t,n,i){return new rA(e,t,n,i)}function Og(e){return e=e.prototype,!(!e||!e.isReactComponent)}function Ca(e,t){var n=e.alternate;return n===null?(n=Qn(e.tag,t,e.key,e.mode),n.elementType=e.elementType,n.type=e.type,n.stateNode=e.stateNode,n.alternate=e,e.alternate=n):(n.pendingProps=t,n.type=e.type,n.flags=0,n.subtreeFlags=0,n.deletions=null),n.flags=e.flags&1206910976,n.childLanes=e.childLanes,n.lanes=e.lanes,n.child=e.child,n.memoizedProps=e.memoizedProps,n.memoizedState=e.memoizedState,n.updateQueue=e.updateQueue,t=e.dependencies,n.dependencies=t===null?null:{lanes:t.lanes,firstContext:t.firstContext},n.sibling=e.sibling,n.index=e.index,n.ref=e.ref,n.refCleanup=e.refCleanup,n}function $x(e,t){e.flags&=1206910978;var n=e.alternate;return n===null?(e.childLanes=0,e.lanes=t,e.child=null,e.subtreeFlags=0,e.memoizedProps=null,e.memoizedState=null,e.updateQueue=null,e.dependencies=null,e.stateNode=null):(e.childLanes=n.childLanes,e.lanes=n.lanes,e.child=n.child,e.subtreeFlags=0,e.deletions=null,e.memoizedProps=n.memoizedProps,e.memoizedState=n.memoizedState,e.updateQueue=n.updateQueue,e.type=n.type,t=n.dependencies,e.dependencies=t===null?null:{lanes:t.lanes,firstContext:t.firstContext}),e}function nf(e,t,n,i,a,s){var r=0;if(i=e,typeof i=="function")Og(i)&&(r=1);else if(typeof i=="string")r=Uw(e,n,sa.current)?26:e==="html"||e==="head"||e==="body"?27:5;else t:switch(i){case gm:return e=Qn(31,n,t,a),e.elementType=gm,e.lanes=s,e;case $r:return sr(n.children,a,s,t);case gx:r=8,a|=24;break;case dm:return e=Qn(12,n,t,a|2),e.elementType=dm,e.lanes=s,e;case pm:return e=Qn(13,n,t,a),e.elementType=pm,e.lanes=s,e;case mm:return e=Qn(19,n,t,a),e.elementType=mm,e.lanes=s,e;case YT:case vm:return e=a|32,e=Qn(30,n,t,e),e.elementType=vm,e.lanes=s,e.stateNode={autoName:null,paired:null,clones:null,ref:null},e;default:if(typeof i=="object"&&i!==null)switch(i.$$typeof){case ea:r=10;break t;case vx:r=9;break t;case bg:r=11;break t;case Mg:r=14;break t;case ss:r=16,i=null;break t}r=29,n=Error(J(130,e===null?"null":typeof e,"")),i=null}return t=Qn(r,n,t,a),t.elementType=e,t.type=i,t.lanes=s,t}function sr(e,t,n,i){return e=Qn(7,e,i,t),e.lanes=n,e}function qp(e,t,n){return e=Qn(6,e,null,t),e.lanes=n,e}function tS(e){var t=Qn(18,null,null,0);return t.stateNode=e,t}function Yp(e,t,n){return t=Qn(4,e.children!==null?e.children:[],e.key,t),t.lanes=n,t.stateNode={containerInfo:e.containerInfo,pendingChildren:null,implementation:e.implementation},t}var j_=new WeakMap;function Ai(e,t){if(typeof e=="object"&&e!==null){var n=j_.get(e);return n!==void 0?n:(t={value:e,source:t,stack:D_(t)},j_.set(e,t),t)}return{value:e,source:t,stack:D_(t)}}var oo=[],lo=0,Ef=null,rc=0,Mi=[],Ti=0,Ts=null,ia=1,aa="";function Aa(e,t){oo[lo++]=rc,oo[lo++]=Ef,Ef=e,rc=t}function eS(e,t,n){Mi[Ti++]=ia,Mi[Ti++]=aa,Mi[Ti++]=Ts,Ts=e;var i=ia;e=aa;var a=32-ui(i)-1;i&=~(1<<a),n+=1;var s=32-ui(t)+a;if(30<s){var r=a-a%5;s=(i&(1<<r)-1).toString(32),i>>=r,a-=r,ia=1<<32-ui(t)+a|n<<a|i,aa=s+e}else ia=1<<s|n<<a|i,aa=e}function th(e){e.return!==null&&(Aa(e,1),eS(e,1,0))}function Pg(e){for(;e===Ef;)Ef=oo[--lo],oo[lo]=null,rc=oo[--lo],oo[lo]=null;for(;e===Ts;)Ts=Mi[--Ti],Mi[Ti]=null,aa=Mi[--Ti],Mi[Ti]=null,ia=Mi[--Ti],Mi[Ti]=null}function nS(e,t){Mi[Ti++]=ia,Mi[Ti++]=aa,Mi[Ti++]=Ts,ia=t.id,aa=t.overflow,Ts=e}var hn=null,Ue=null,Qt=!1,ms=null,wi=!1,Dm=Error(J(519));function Es(e){var t=Error(J(418,1<arguments.length&&arguments[1]!==void 0&&arguments[1]?"text":"HTML",""));throw oc(Ai(t,e)),Dm}function $_(e){var t=e.stateNode,n=e.type,i=e.memoizedProps;switch(t[yn]=e,t[ti]=i,n){case"dialog":jt("cancel",t),jt("close",t);break;case"iframe":case"object":case"embed":jt("load",t);break;case"video":case"audio":for(n=0;n<fc.length;n++)jt(fc[n],t);break;case"source":jt("error",t);break;case"img":case"image":case"link":jt("error",t),jt("load",t);break;case"details":jt("toggle",t);break;case"input":jt("invalid",t),Ux(t,i.value,i.defaultValue,i.checked,i.defaultChecked,i.type,i.name,!0);break;case"select":jt("invalid",t);break;case"textarea":jt("invalid",t),Ix(t,i.value,i.defaultValue,i.children)}n=i.children,typeof n!="string"&&typeof n!="number"&&typeof n!="bigint"||t.textContent===""+n||i.suppressHydrationWarning===!0||zb(t.textContent,n)?(i.popover!=null&&(jt("beforetoggle",t),jt("toggle",t)),i.onScroll!=null&&jt("scroll",t),i.onScrollEnd!=null&&jt("scrollend",t),i.onClick!=null&&(t.onclick=na),t=!0):t=!1,t||Es(e,!0)}function Af(e){for(hn=e.return;hn;)switch(hn.tag){case 5:case 31:case 13:wi=!1;return;case 27:case 3:wi=!0;return;default:hn=hn.return}}function Yr(e){if(e!==hn)return!1;if(!Qt)return Af(e),Qt=!0,!1;var t=e.tag,n;if((n=t!==3&&t!==27)&&((n=t===5)&&(n=e.type,n=!(n!=="form"&&n!=="button")||pg(e.type,e.memoizedProps)),n=!n),n&&Ue&&Es(e),Af(e),t===13){if(e=e.memoizedState,e=e!==null?e.dehydrated:null,!e)throw Error(J(317));Ue=Jy(e)}else if(t===31){if(e=e.memoizedState,e=e!==null?e.dehydrated:null,!e)throw Error(J(317));Ue=Jy(e)}else t===27?(t=Ue,Ds(e.type)?(e=_g,_g=null,Ue=e):Ue=t):Ue=hn?Ci(e.stateNode.nextSibling):null;return!0}function cr(){Ue=hn=null,Qt=!1}function Zp(){var e=ms;return e!==null&&(Jn===null?Jn=e:Jn.push.apply(Jn,e),ms=null),e}function oc(e){ms===null?ms=[e]:ms.push(e)}var Nm=la(null),Sr=null,wa=null;function us(e,t,n){Le(Nm,t._currentValue),t._currentValue=n}function Ra(e){e._currentValue=Nm.current,Sn(Nm)}function af(e,t,n){for(;e!==null;){var i=e.alternate;if((e.childLanes&t)!==t?(e.childLanes|=t,i!==null&&(i.childLanes|=t)):i!==null&&(i.childLanes&t)!==t&&(i.childLanes|=t),e===n)break;e=e.return}}function Um(e,t,n,i){var a=e.child;for(a!==null&&(a.return=e);a!==null;){var s=a.dependencies;if(s!==null){var r=a.child;s=s.firstContext;t:for(;s!==null;){var o=s;s=a;for(var l=0;l<t.length;l++)if(o.context===t[l]){s.lanes|=n,o=s.alternate,o!==null&&(o.lanes|=n),af(s.return,n,e),i||(r=null);break t}s=o.next}}else if(a.tag===18){if(r=a.return,r===null)throw Error(J(341));r.lanes|=n,s=r.alternate,s!==null&&(s.lanes|=n),af(r,n,e),r=null}else a.tag===13&&a.memoizedState!==null&&a.memoizedState.dehydrated===null?(a.lanes|=n,r=a.alternate,r!==null&&(r.lanes|=n),af(a.return,n,e),r=a.child,r=r!==null?r.sibling:null):r=a.child;if(r!==null)r.return=a;else for(r=a;r!==null;){if(r===e){r=null;break}if(a=r.sibling,a!==null){a.return=r.return,r=a;break}r=r.return}a=r}}function ur(e,t,n,i){e=null;for(var a=t,s=!1;a!==null;){if(!s){if((a.flags&524288)!==0)s=!0;else if((a.flags&262144)!==0)break}if(a.tag===10){var r=a.alternate;if(r===null)throw Error(J(387));if(r=r.memoizedProps,r!==null){var o=a.type;hi(a.pendingProps.value,r.value)||(e!==null?e.push(o):e=[o])}}else if(a===_f.current){if(r=a.alternate,r===null)throw Error(J(387));r.memoizedState.memoizedState!==a.memoizedState.memoizedState&&(e!==null?e.push(No):e=[No])}a=a.return}return e!==null&&Um(t,e,n,i),t.flags|=262144,e!==null}function wf(e){for(e=e.firstContext;e!==null;){if(!hi(e.context._currentValue,e.memoizedValue))return!0;e=e.next}return!1}function fr(e){Sr=e,wa=null,e=e.dependencies,e!==null&&(e.firstContext=null)}function xn(e){return iS(Sr,e)}function Fu(e,t){return Sr===null&&fr(e),iS(e,t)}function iS(e,t){var n=t._currentValue;if(t={context:t,memoizedValue:n,next:null},wa===null){if(e===null)throw Error(J(308));wa=t,e.dependencies={lanes:0,firstContext:t},e.flags|=524288}else wa=wa.next=t;return n}var oA=typeof AbortController<"u"?AbortController:function(){var e=[],t=this.signal={aborted:!1,addEventListener:function(n,i){e.push(i)}};this.abort=function(){t.aborted=!0,e.forEach(function(n){return n()})}},lA=an.unstable_scheduleCallback,cA=an.unstable_NormalPriority,Qe={$$typeof:ea,Consumer:null,Provider:null,_currentValue:null,_currentValue2:null,_threadCount:0};function Bg(){return{controller:new oA,data:new Map,refCount:0}}function Tc(e){e.refCount--,e.refCount===0&&lA(cA,function(){e.controller.abort()})}function ty(e,t){if((e.pendingLanes&4194048)!==0){var n=e.transitionTypes;for(n===null&&(n=e.transitionTypes=[]),e=0;e<t.length;e++){var i=t[e];n.indexOf(i)===-1&&n.push(i)}}}var Hl=null;function uA(e){var t=e.transitionTypes;return e.transitionTypes=null,t}var ql=null,Lm=0,hr=0,po=null;function fA(e,t){if(ql===null){var n=ql=[];Lm=0,hr=f0(),po={status:"pending",value:void 0,then:function(i){n.push(i)}}}return Lm++,t.then(ey,ey),t}function ey(){if(--Lm===0&&(Hl=null,ql!==null)){po!==null&&(po.status="fulfilled");var e=ql;ql=null,hr=0,po=null;for(var t=0;t<e.length;t++)(0,e[t])()}}function hA(e,t){var n=[],i={status:"pending",value:null,reason:null,then:function(a){n.push(a)}};return e.then(function(){i.status="fulfilled",i.value=t;for(var a=0;a<n.length;a++)(0,n[a])(t)},function(a){for(i.status="rejected",i.reason=a,a=0;a<n.length;a++)(0,n[a])(void 0)}),i}var ny=Bt.S;Bt.S=function(e,t){if(Sb=li(),typeof t=="object"&&t!==null&&typeof t.then=="function"&&fA(e,t),Hl!==null)for(var n=Co;n!==null;)ty(n,Hl),n=n.next;if(n=e.types,n!==null){for(var i=Co;i!==null;)ty(i,n),i=i.next;if(hr!==0){i=Hl,i===null&&(i=Hl=[]);for(var a=0;a<n.length;a++){var s=n[a];i.indexOf(s)===-1&&i.push(s)}}}ny!==null&&ny(e,t)};var rr=la(null);function zg(){var e=rr.current;return e!==null?e:Ae.pooledCache}function sf(e,t){t===null?Le(rr,rr.current):Le(rr,t.pool)}function aS(){var e=zg();return e===null?null:{parent:Qe._currentValue,pool:e}}var Po=Error(J(460)),Fg=Error(J(474)),eh=Error(J(542)),Cf={then:function(){}};function iy(e){return e=e.status,e==="fulfilled"||e==="rejected"}function sS(e,t,n){switch(n=e[n],n===void 0?e.push(t):n!==t&&(t.then(na,na),t=n),t.status){case"fulfilled":return t.value;case"rejected":throw e=t.reason,sy(e),e===void 0&&!("reason"in t)?Error(J(600)):e;default:if(typeof t.status=="string")t.then(na,na);else{if(e=Ae,e!==null&&100<e.shellSuspendCounter)throw Error(J(482));e=t,e.status="pending",e.then(function(i){if(t.status==="pending"){var a=t;a.status="fulfilled",a.value=i}},function(i){if(t.status==="pending"){var a=t;a.status="rejected",a.reason=i}})}switch(t.status){case"fulfilled":return t.value;case"rejected":throw e=t.reason,sy(e),e}throw or=t,Po}}function tr(e){try{var t=e._init;return t(e._payload)}catch(n){throw n!==null&&typeof n=="object"&&typeof n.then=="function"?(or=n,Po):n}}var or=null;function ay(){if(or===null)throw Error(J(459));var e=or;return or=null,e}function sy(e){if(e===Po||e===eh)throw Error(J(483))}var mo=null,lc=0;function Vu(e){var t=lc;return lc+=1,mo===null&&(mo=[]),sS(mo,e,t)}function ns(e,t){t=t.props.ref,e.ref=t!==void 0?t:null}function Hu(e,t){throw t.$$typeof===qT?Error(J(525)):(e=Object.prototype.toString.call(t),Error(J(31,e==="[object Object]"?"object with keys {"+Object.keys(t).join(", ")+"}":e)))}function rS(e){function t(h,m){if(e){var g=h.deletions;g===null?(h.deletions=[m],h.flags|=16):g.push(m)}}function n(h,m){if(!e)return null;for(;m!==null;)t(h,m),m=m.sibling;return null}function i(h){for(var m=new Map;h!==null;)h.key===null?m.set(h.index,h):m.set(h.key,h),h=h.sibling;return m}function a(h,m){return h=Ca(h,m),h.index=0,h.sibling=null,h}function s(h,m,g){return h.index=g,e?(g=h.alternate,g!==null?(g=g.index,g<m?(h.flags|=2,m):g):(h.flags|=134217730,m)):(h.flags|=1048576,m)}function r(h){return e&&h.alternate===null&&(h.flags|=134217730),h}function o(h,m,g,x){return m===null||m.tag!==6?(m=qp(g,h.mode,x),m.return=h,m):(m=a(m,g),m.return=h,m)}function l(h,m,g,x){var E=g.type;return E===$r?(h=f(h,m,g.props.children,x,g.key),ns(h,g),h):m!==null&&(m.elementType===E||typeof E=="object"&&E!==null&&E.$$typeof===ss&&tr(E)===m.type)?(m=a(m,g.props),ns(m,g),m.return=h,m):(m=nf(g.type,g.key,g.props,null,h.mode,x),ns(m,g),m.return=h,m)}function c(h,m,g,x){return m===null||m.tag!==4||m.stateNode.containerInfo!==g.containerInfo||m.stateNode.implementation!==g.implementation?(m=Yp(g,h.mode,x),m.return=h,m):(m=a(m,g.children||[]),m.return=h,m)}function f(h,m,g,x,E){return m===null||m.tag!==7?(m=sr(g,h.mode,x,E),m.return=h,m):(m=a(m,g),m.return=h,m)}function p(h,m,g){if(typeof m=="string"&&m!==""||typeof m=="number"||typeof m=="bigint")return m=qp(""+m,h.mode,g),m.return=h,m;if(typeof m=="object"&&m!==null){switch(m.$$typeof){case Nu:return g=nf(m.type,m.key,m.props,null,h.mode,g),ns(g,m),g.return=h,g;case zl:return m=Yp(m,h.mode,g),m.return=h,m;case ss:return m=tr(m),p(h,m,g)}if(Fl(m)||Ul(m))return m=sr(m,h.mode,g,null),m.return=h,m;if(typeof m.then=="function")return p(h,Vu(m),g);if(m.$$typeof===ea)return p(h,Fu(h,m),g);Hu(h,m)}return null}function u(h,m,g,x){var E=m!==null?m.key:null;if(typeof g=="string"&&g!==""||typeof g=="number"||typeof g=="bigint")return E!==null?null:o(h,m,""+g,x);if(typeof g=="object"&&g!==null){switch(g.$$typeof){case Nu:return g.key===E?l(h,m,g,x):null;case zl:return g.key===E?c(h,m,g,x):null;case ss:return g=tr(g),u(h,m,g,x)}if(Fl(g)||Ul(g))return E!==null?null:f(h,m,g,x,null);if(typeof g.then=="function")return u(h,m,Vu(g),x);if(g.$$typeof===ea)return u(h,m,Fu(h,g),x);Hu(h,g)}return null}function d(h,m,g,x,E){if(typeof x=="string"&&x!==""||typeof x=="number"||typeof x=="bigint")return h=h.get(g)||null,o(m,h,""+x,E);if(typeof x=="object"&&x!==null){switch(x.$$typeof){case Nu:return h=h.get(x.key===null?g:x.key)||null,l(m,h,x,E);case zl:return h=h.get(x.key===null?g:x.key)||null,c(m,h,x,E);case ss:return x=tr(x),d(h,m,g,x,E)}if(Fl(x)||Ul(x))return h=h.get(g)||null,f(m,h,x,E,null);if(typeof x.then=="function")return d(h,m,g,Vu(x),E);if(x.$$typeof===ea)return d(h,m,g,Fu(m,x),E);Hu(m,x)}return null}function _(h,m,g,x){for(var E=null,A=null,w=m,y=m=0,T=null;w!==null&&y<g.length;y++){w.index>y?(T=w,w=null):T=w.sibling;var L=u(h,w,g[y],x);if(L===null){w===null&&(w=T);break}e&&w&&L.alternate===null&&t(h,w),m=s(L,m,y),A===null?E=L:A.sibling=L,A=L,w=T}if(y===g.length)return n(h,w),Qt&&Aa(h,y),E;if(w===null){for(;y<g.length;y++)w=p(h,g[y],x),w!==null&&(m=s(w,m,y),A===null?E=w:A.sibling=w,A=w);return Qt&&Aa(h,y),E}for(w=i(w);y<g.length;y++)T=d(w,h,y,g[y],x),T!==null&&(e&&(L=T.alternate,L!==null&&w.delete(L.key===null?y:L.key)),m=s(T,m,y),A===null?E=T:A.sibling=T,A=T);return e&&w.forEach(function(R){return t(h,R)}),Qt&&Aa(h,y),E}function b(h,m,g,x){if(g==null)throw Error(J(151));for(var E=null,A=null,w=m,y=m=0,T=null,L=g.next();w!==null&&!L.done;y++,L=g.next()){w.index>y?(T=w,w=null):T=w.sibling;var R=u(h,w,L.value,x);if(R===null){w===null&&(w=T);break}e&&w&&R.alternate===null&&t(h,w),m=s(R,m,y),A===null?E=R:A.sibling=R,A=R,w=T}if(L.done)return n(h,w),Qt&&Aa(h,y),E;if(w===null){for(;!L.done;y++,L=g.next())L=p(h,L.value,x),L!==null&&(m=s(L,m,y),A===null?E=L:A.sibling=L,A=L);return Qt&&Aa(h,y),E}for(w=i(w);!L.done;y++,L=g.next())L=d(w,h,y,L.value,x),L!==null&&(e&&(T=L.alternate,T!==null&&w.delete(T.key===null?y:T.key)),m=s(L,m,y),A===null?E=L:A.sibling=L,A=L);return e&&w.forEach(function(N){return t(h,N)}),Qt&&Aa(h,y),E}function v(h,m,g,x){if(typeof g=="object"&&g!==null&&g.type===$r&&g.key===null&&g.props.ref===void 0&&(g=g.props.children),typeof g=="object"&&g!==null){switch(g.$$typeof){case Nu:t:{for(var E=g.key;m!==null;){if(m.key===E){if(E=g.type,E===$r){if(m.tag===7){n(h,m.sibling),x=a(m,g.props.children),ns(x,g),x.return=h,h=x;break t}}else if(m.elementType===E||typeof E=="object"&&E!==null&&E.$$typeof===ss&&tr(E)===m.type){n(h,m.sibling),x=a(m,g.props),ns(x,g),x.return=h,h=x;break t}n(h,m);break}else t(h,m);m=m.sibling}g.type===$r?(x=sr(g.props.children,h.mode,x,g.key),ns(x,g),x.return=h,h=x):(x=nf(g.type,g.key,g.props,null,h.mode,x),ns(x,g),x.return=h,h=x)}return r(h);case zl:t:{for(E=g.key;m!==null;){if(m.key===E)if(m.tag===4&&m.stateNode.containerInfo===g.containerInfo&&m.stateNode.implementation===g.implementation){n(h,m.sibling),x=a(m,g.children||[]),x.return=h,h=x;break t}else{n(h,m);break}else t(h,m);m=m.sibling}x=Yp(g,h.mode,x),x.return=h,h=x}return r(h);case ss:return g=tr(g),v(h,m,g,x)}if(Fl(g))return _(h,m,g,x);if(Ul(g)){if(E=Ul(g),typeof E!="function")throw Error(J(150));return g=E.call(g),b(h,m,g,x)}if(typeof g.then=="function")return v(h,m,Vu(g),x);if(g.$$typeof===ea)return v(h,m,Fu(h,g),x);Hu(h,g)}return typeof g=="string"&&g!==""||typeof g=="number"||typeof g=="bigint"?(g=""+g,m!==null&&m.tag===6?(n(h,m.sibling),x=a(m,g),x.return=h,h=x):(n(h,m),x=qp(g,h.mode,x),x.return=h,h=x),r(h)):n(h,m)}return function(h,m,g,x){try{lc=0;var E=v(h,m,g,x);return mo=null,E}catch(w){if(w===Po||w===eh)throw w;var A=Qn(29,w,null,h.mode);return A.lanes=x,A.return=h,A}}}var dr=rS(!0),oS=rS(!1),rs=!1;function Vg(e){e.updateQueue={baseState:e.memoizedState,firstBaseUpdate:null,lastBaseUpdate:null,shared:{pending:null,lanes:0,hiddenCallbacks:null},callbacks:null}}function Im(e,t){e=e.updateQueue,t.updateQueue===e&&(t.updateQueue={baseState:e.baseState,firstBaseUpdate:e.firstBaseUpdate,lastBaseUpdate:e.lastBaseUpdate,shared:e.shared,callbacks:null})}function gs(e){return{lane:e,tag:0,payload:null,callback:null,next:null}}function vs(e,t,n){var i=e.updateQueue;if(i===null)return null;if(i=i.shared,(fe&2)!==0){var a=i.pending;return a===null?t.next=t:(t.next=a.next,a.next=t),i.pending=t,t=Tf(e),jx(e,null,n),t}return $f(e,i,t,n),Tf(e)}function Yl(e,t,n){if(t=t.updateQueue,t!==null&&(t=t.shared,(n&4194048)!==0)){var i=t.lanes;i&=e.pendingLanes,n|=i,t.lanes=n,Tx(e,n)}}function Jp(e,t){var n=e.updateQueue,i=e.alternate;if(i!==null&&(i=i.updateQueue,n===i)){var a=null,s=null;if(n=n.firstBaseUpdate,n!==null){do{var r={lane:n.lane,tag:n.tag,payload:n.payload,callback:null,next:null};s===null?a=s=r:s=s.next=r,n=n.next}while(n!==null);s===null?a=s=t:s=s.next=t}else a=s=t;n={baseState:i.baseState,firstBaseUpdate:a,lastBaseUpdate:s,shared:i.shared,callbacks:i.callbacks},e.updateQueue=n;return}e=n.lastBaseUpdate,e===null?n.firstBaseUpdate=t:e.next=t,n.lastBaseUpdate=t}var Om=!1;function Zl(){if(Om){var e=po;if(e!==null)throw e}}function Jl(e,t,n,i){Om=!1;var a=e.updateQueue;rs=!1;var s=a.firstBaseUpdate,r=a.lastBaseUpdate,o=a.shared.pending;if(o!==null){a.shared.pending=null;var l=o,c=l.next;l.next=null,r===null?s=c:r.next=c,r=l;var f=e.alternate;f!==null&&(f=f.updateQueue,o=f.lastBaseUpdate,o!==r&&(o===null?f.firstBaseUpdate=c:o.next=c,f.lastBaseUpdate=l))}if(s!==null){var p=a.baseState;r=0,f=c=l=null,o=s;do{var u=o.lane&-536870913,d=u!==o.lane;if(d?(ee&u)===u:(i&u)===u){u!==0&&u===hr&&(Om=!0),f!==null&&(f=f.next={lane:0,tag:o.tag,payload:o.payload,callback:null,next:null});t:{var _=e,b=o;u=t;var v=n;switch(b.tag){case 1:if(_=b.payload,typeof _=="function"){p=_.call(v,p,u);break t}p=_;break t;case 3:_.flags=_.flags&-65537|128;case 0:if(_=b.payload,u=typeof _=="function"?_.call(v,p,u):_,u==null)break t;p=we({},p,u);break t;case 2:rs=!0}}u=o.callback,u!==null&&(e.flags|=64,d&&(e.flags|=8192),d=a.callbacks,d===null?a.callbacks=[u]:d.push(u))}else d={lane:u,tag:o.tag,payload:o.payload,callback:o.callback,next:null},f===null?(c=f=d,l=p):f=f.next=d,r|=u;if(o=o.next,o===null){if(o=a.shared.pending,o===null)break;d=o,o=d.next,d.next=null,a.lastBaseUpdate=d,a.shared.pending=null}}while(!0);f===null&&(l=p),a.baseState=l,a.firstBaseUpdate=c,a.lastBaseUpdate=f,s===null&&(a.shared.lanes=0),Cs|=r,e.lanes=r,e.memoizedState=p}}function lS(e,t){if(typeof e!="function")throw Error(J(191,e));e.call(t)}function cS(e,t){var n=e.callbacks;if(n!==null)for(e.callbacks=null,e=0;e<n.length;e++)lS(n[e],t)}var As=la(null),Rf=la(0);function ry(e,t){e=Oa,Le(Rf,e),Le(As,t),Oa=e|t.baseLanes}function Pm(){Le(Rf,Oa),Le(As,As.current)}function Hg(){Oa=Rf.current,Sn(As),Sn(Rf)}var Tn=la(null),Un=null;function _s(e){var t=e.alternate;Le(bn,bn.current&1),Le(Tn,e),Un===null&&(t===null||As.current!==null||t.memoizedState!==null)&&(Un=e)}function Bm(e){Le(bn,bn.current),Le(Tn,e),Un===null&&(Un=e)}function uS(e){e.tag===22?(Le(bn,bn.current),Le(Tn,e),Un===null&&(Un=e)):ys()}function ys(){Le(bn,bn.current),Le(Tn,Tn.current)}function si(e){Sn(Tn),Un===e&&(Un=null),Sn(bn)}var bn=la(0);function cc(e,t){Le(Tn,Tn.current),Le(bn,t)}function Gg(e){Sn(bn),Sn(Tn),Un===e&&(Un=null)}function Df(e){for(var t=e;t!==null;){if(t.tag===13){var n=t.memoizedState;if(n!==null&&(n=n.dehydrated,n===null||vg(n)||m0(n)))return t}else if(t.tag===19&&t.memoizedProps.revealOrder!=="independent"){if((t.flags&128)!==0)return t}else if(t.child!==null){t.child.return=t,t=t.child;continue}if(t===e)break;for(;t.sibling===null;){if(t.return===null||t.return===e)return null;t=t.return}t.sibling.return=t.return,t=t.sibling}return null}var Ua=0,Yt=null,Me=null,Ke=null,Nf=!1,go=!1,pr=!1,Uf=0,uc=0,vo=null,dA=0;function We(){throw Error(J(321))}function kg(e,t){if(t===null)return!1;for(var n=0;n<t.length&&n<e.length;n++)if(!hi(e[n],t[n]))return!1;return!0}function Xg(e,t,n,i,a,s){return Ua=s,Yt=t,t.memoizedState=null,t.updateQueue=null,t.lanes=0,Bt.H=e===null||e.memoizedState===null?HS:GS,pr=!1,s=n(i,a),pr=!1,go&&(s=hS(t,n,i,a)),fS(e),s}function fS(e){Bt.H=Lf;var t=Me!==null&&Me.next!==null;if(Ua=0,Ke=Me=Yt=null,Nf=!1,uc=0,vo=null,t)throw Error(J(300));e===null||je||(e=e.dependencies,e!==null&&wf(e)&&(je=!0))}function hS(e,t,n,i){Yt=e;var a=0;do{if(go&&(vo=null),uc=0,go=!1,25<=a)throw Error(J(301));if(a+=1,Ke=Me=null,e.updateQueue!=null){var s=e.updateQueue;s.lastEffect=null,s.events=null,s.stores=null,s.memoCache!=null&&(s.memoCache.index=0)}Bt.H=SA,s=t(n,i)}while(go);return s}function pA(){var e=Bt.H,t=e.useState()[0];return t=typeof t.then=="function"?Ec(t):t,e=e.useState()[0],(Me!==null?Me.memoizedState:null)!==e&&(Yt.flags|=1024),t}function Wg(){var e=Uf!==0;return Uf=0,e}function qg(e,t,n){t.updateQueue=e.updateQueue,t.flags&=-2053,e.lanes&=~n}function Yg(e){if(Nf){for(e=e.memoizedState;e!==null;){var t=e.queue;t!==null&&(t.pending=null),e=e.next}Nf=!1}Ua=0,Ke=Me=Yt=null,go=!1,uc=Uf=0,vo=null}function kn(){var e={memoizedState:null,baseState:null,baseQueue:null,queue:null,next:null};return Ke===null?Yt.memoizedState=Ke=e:Ke=Ke.next=e,Ke}function Je(){if(Me===null){var e=Yt.alternate;e=e!==null?e.memoizedState:null}else e=Me.next;var t=Ke===null?Yt.memoizedState:Ke.next;if(t!==null)Ke=t,Me=e;else{if(e===null)throw Yt.alternate===null?Error(J(467)):Error(J(310));Me=e,e={memoizedState:Me.memoizedState,baseState:Me.baseState,baseQueue:Me.baseQueue,queue:Me.queue,next:null},Ke===null?Yt.memoizedState=Ke=e:Ke=Ke.next=e}return Ke}function nh(){return{lastEffect:null,events:null,stores:null,memoCache:null}}function Ec(e){var t=uc;return uc+=1,vo===null&&(vo=[]),e=sS(vo,e,t),t=Yt,(Ke===null?t.memoizedState:Ke.next)===null&&(t=t.alternate,Bt.H=t===null||t.memoizedState===null?HS:GS),e}function ih(e){if(e!==null&&typeof e=="object"){if(typeof e.then=="function")return Ec(e);if(e.$$typeof===JT)return;if(e.$$typeof===ea)return xn(e)}throw Error(J(438,String(e)))}function Zg(e){var t=null,n=Yt.updateQueue;if(n!==null&&(t=n.memoCache),t==null){var i=Yt.alternate;i!==null&&(i=i.updateQueue,i!==null&&(i=i.memoCache,i!=null&&(t={data:i.data.map(function(a){return a.slice()}),index:0})))}if(t==null&&(t={data:[],index:0}),n===null&&(n=nh(),Yt.updateQueue=n),n.memoCache=t,n=t.data[t.index],n===void 0)for(n=t.data[t.index]=Array(e),i=0;i<e;i++)n[i]=ZT;return t.index++,n}function La(e,t){return typeof t=="function"?t(e):t}function rf(e){var t=Je();return Jg(t,Me,e)}function Jg(e,t,n){var i=e.queue;if(i===null)throw Error(J(311));i.lastRenderedReducer=n;var a=e.baseQueue,s=i.pending;if(s!==null){if(a!==null){var r=a.next;a.next=s.next,s.next=r}t.baseQueue=a=s,i.pending=null}if(s=e.baseState,a===null)e.memoizedState=s;else{t=a.next;var o=r=null,l=null,c=t,f=!1;do{var p=c.lane&-536870913;if(p!==c.lane?(ee&p)===p:(Ua&p)===p){var u=c.revertLane;if(u===0)l!==null&&(l=l.next={lane:0,revertLane:0,gesture:null,action:c.action,hasEagerState:c.hasEagerState,eagerState:c.eagerState,next:null}),p===hr&&(f=!0);else if((Ua&u)===u){c=c.next,u===hr&&(f=!0);continue}else p={lane:0,revertLane:c.revertLane,gesture:null,action:c.action,hasEagerState:c.hasEagerState,eagerState:c.eagerState,next:null},l===null?(o=l=p,r=s):l=l.next=p,Yt.lanes|=u,Cs|=u;p=c.action,pr&&n(s,p),s=c.hasEagerState?c.eagerState:n(s,p)}else u={lane:p,revertLane:c.revertLane,gesture:c.gesture,action:c.action,hasEagerState:c.hasEagerState,eagerState:c.eagerState,next:null},l===null?(o=l=u,r=s):l=l.next=u,Yt.lanes|=p,Cs|=p;c=c.next}while(c!==null&&c!==t);if(l===null?r=s:l.next=o,!hi(s,e.memoizedState)&&(je=!0,f&&(n=po,n!==null)))throw n;e.memoizedState=s,e.baseState=r,e.baseQueue=l,i.lastRenderedState=s}return a===null&&(i.lanes=0),[e.memoizedState,i.dispatch]}function Kp(e){var t=Je(),n=t.queue;if(n===null)throw Error(J(311));n.lastRenderedReducer=e;var i=n.dispatch,a=n.pending,s=t.memoizedState;if(a!==null){n.pending=null;var r=a=a.next;do s=e(s,r.action),r=r.next;while(r!==a);hi(s,t.memoizedState)||(je=!0),t.memoizedState=s,t.baseQueue===null&&(t.baseState=s),n.lastRenderedState=s}return[s,i]}function dS(e,t,n){var i=Yt,a=Je(),s=Qt;if(s){if(n===void 0)throw Error(J(407));n=n()}else n=t();var r=!hi((Me||a).memoizedState,n);if(r&&(a.memoizedState=n,je=!0),a=a.queue,Kg(gS.bind(null,i,a,e),[e]),e=a.getSnapshot!==t||r||Ke!==null&&(Ke.memoizedState.tag&1)!==0,To(e?9:8,{destroy:void 0},mS.bind(null,i,a,n,t),null),e){if(i.flags|=2048,Ae===null)throw Error(J(349));s||(Ua&127)!==0||pS(i,t,n)}return n}function pS(e,t,n){e.flags|=16384,e={getSnapshot:t,value:n},t=Yt.updateQueue,t===null?(t=nh(),Yt.updateQueue=t,t.stores=[e]):(n=t.stores,n===null?t.stores=[e]:n.push(e))}function mS(e,t,n,i){t.value=n,t.getSnapshot=i,vS(t)&&_S(e)}function gS(e,t,n){return n(function(){vS(t)&&_S(e)})}function vS(e){var t=e.getSnapshot;e=e.value;try{var n=t();return!hi(e,n)}catch{return!0}}function _S(e){var t=xr(e,2);t!==null&&jn(t,e,2)}function zm(e){var t=kn();if(typeof e=="function"){var n=e;if(e=n(),pr){ls(!0);try{n()}finally{ls(!1)}}}return t.memoizedState=t.baseState=e,t.queue={pending:null,lanes:0,dispatch:null,lastRenderedReducer:La,lastRenderedState:e},t}function yS(e,t,n,i){return e.baseState=n,Jg(e,Me,typeof i=="function"?i:La)}function mA(e,t,n,i,a){if(sh(e))throw Error(J(485));if(e=t.action,e!==null){var s={payload:a,action:e,next:null,isTransition:!0,status:"pending",value:null,reason:null,listeners:[],then:function(r){s.listeners.push(r)}};Bt.T!==null?n(!0):s.isTransition=!1,i(s),n=t.pending,n===null?(s.next=t.pending=s,xS(t,s)):(s.next=n.next,t.pending=n.next=s)}}function xS(e,t){var n=t.action,i=t.payload,a=e.state;if(t.isTransition){var s=Bt.T,r={};r.types=s!==null?s.types:null,Bt.T=r;try{var o=n(a,i),l=Bt.S;l!==null&&l(r,o),oy(e,t,o)}catch(c){Fm(e,t,c)}finally{s!==null&&r.types!==null&&(s.types=r.types),Bt.T=s}}else try{s=n(a,i),oy(e,t,s)}catch(c){Fm(e,t,c)}}function oy(e,t,n){n!==null&&typeof n=="object"&&typeof n.then=="function"?n.then(function(i){ly(e,t,i)},function(i){return Fm(e,t,i)}):ly(e,t,n)}function ly(e,t,n){t.status="fulfilled",t.value=n,SS(t),e.state=n,t=e.pending,t!==null&&(n=t.next,n===t?e.pending=null:(n=n.next,t.next=n,xS(e,n)))}function Fm(e,t,n){var i=e.pending;if(e.pending=null,i!==null){i=i.next;do t.status="rejected",t.reason=n,SS(t),t=t.next;while(t!==i)}e.action=null}function SS(e){e=e.listeners;for(var t=0;t<e.length;t++)(0,e[t])()}function bS(e,t){return t}function cy(e,t){if(Qt){var n=Ae.formState;if(n!==null){t:{var i=Yt;if(Qt){if(Ue){e:{for(var a=Ue,s=wi;a.nodeType!==8;){if(!s){a=null;break e}if(a=Ci(a.nextSibling),a===null){a=null;break e}}s=a.data,a=s==="F!"||s==="F"?a:null}if(a){Ue=Ci(a.nextSibling),i=a.data==="F!";break t}}Es(i)}i=!1}i&&(t=n[0])}}return n=kn(),n.memoizedState=n.baseState=t,i={pending:null,lanes:0,dispatch:null,lastRenderedReducer:bS,lastRenderedState:t},n.queue=i,n=zS.bind(null,Yt,i),i.dispatch=n,i=zm(!1),s=t0.bind(null,Yt,!1,i.queue),i=kn(),a={state:t,dispatch:null,action:e,pending:null},i.queue=a,n=mA.bind(null,Yt,a,s,n),a.dispatch=n,i.memoizedState=e,[t,n,!1]}function uy(e){var t=Je();return MS(t,Me,e)}function MS(e,t,n){if(t=Jg(e,t,bS)[0],e=rf(La)[0],typeof t=="object"&&t!==null&&typeof t.then=="function")try{var i=Ec(t)}catch(r){throw r===Po?eh:r}else i=t;t=Je();var a=t.queue,s=a.dispatch;return n!==t.memoizedState&&(Yt.flags|=2048,To(9,{destroy:void 0},gA.bind(null,a,n),null)),[i,s,e]}function gA(e,t){e.action=t}function fy(e){var t=Je(),n=Me;if(n!==null)return MS(t,n,e);Je(),t=t.memoizedState,n=Je();var i=n.queue.dispatch;return n.memoizedState=e,[t,i,!1]}function To(e,t,n,i){return e={tag:e,create:n,deps:i,inst:t,next:null},t=Yt.updateQueue,t===null&&(t=nh(),Yt.updateQueue=t),n=t.lastEffect,n===null?t.lastEffect=e.next=e:(i=n.next,n.next=e,e.next=i,t.lastEffect=e),e}function TS(){return Je().memoizedState}function of(e,t,n,i){var a=kn();Yt.flags|=e,a.memoizedState=To(1|t,{destroy:void 0},n,i===void 0?null:i)}function ah(e,t,n,i){var a=Je();i=i===void 0?null:i;var s=a.memoizedState.inst;Me!==null&&i!==null&&kg(i,Me.memoizedState.deps)?a.memoizedState=To(t,s,n,i):(Yt.flags|=e,a.memoizedState=To(1|t,s,n,i))}function hy(e,t){of(8390656,8,e,t)}function Kg(e,t){ah(2048,8,e,t)}function vA(e){Yt.flags|=4;var t=Yt.updateQueue;if(t===null)t=nh(),Yt.updateQueue=t,t.events=[e];else{var n=t.events;n===null?t.events=[e]:n.push(e)}}function ES(e){var t=Je().memoizedState;return vA({ref:t,nextImpl:e}),function(){if((fe&2)!==0)throw Error(J(440));return t.impl.apply(void 0,arguments)}}function AS(e,t){return ah(4,2,e,t)}function wS(e,t){return ah(4,4,e,t)}function CS(e,t){if(typeof t=="function"){e=e();var n=t(e);return function(){typeof n=="function"?n():t(null)}}if(t!=null)return e=e(),t.current=e,function(){t.current=null}}function RS(e,t,n){n=n!=null?n.concat([e]):null,ah(4,4,CS.bind(null,t,e),n)}function Qg(){}function DS(e,t){var n=Je();t=t===void 0?null:t;var i=n.memoizedState;return t!==null&&kg(t,i[1])?i[0]:(n.memoizedState=[e,t],e)}function NS(e,t){var n=Je();t=t===void 0?null:t;var i=n.memoizedState;if(t!==null&&kg(t,i[1]))return i[0];if(i=e(),pr){ls(!0);try{e()}finally{ls(!1)}}return n.memoizedState=[i,t],i}function jg(e,t,n){return n===void 0||(Ua&1073741824)!==0&&(ee&261930)===0?e.memoizedState=t:(e.memoizedState=n,e=Mb(),Yt.lanes|=e,Cs|=e,n)}function US(e,t,n,i){return hi(n,t)?n:As.current!==null?(e=jg(e,n,i),hi(e,t)||(je=!0),e):(Ua&106)===0||(Ua&1073741824)!==0&&(ee&261930)===0?(je=!0,e.memoizedState=n):(e=Mb(),Yt.lanes|=e,Cs|=e,t)}function LS(e,t,n,i,a){var s=he.p;he.p=s!==0&&8>s?s:8;var r=Bt.T,o={};o.types=r!==null?r.types:null,Bt.T=o,t0(e,!1,t,n);try{var l=a(),c=Bt.S;if(c!==null&&c(o,l),l!==null&&typeof l=="object"&&typeof l.then=="function"){var f=hA(l,i);Kl(e,t,f,fi(e))}else Kl(e,t,i,fi(e))}catch(p){Kl(e,t,{then:function(){},status:"rejected",reason:p},fi())}finally{he.p=s,r!==null&&o.types!==null&&(r.types=o.types),Bt.T=r}}function _A(){}function Vm(e,t,n,i){if(e.tag!==5)throw Error(J(476));var a=IS(e).queue;LS(e,a,t,ar,n===null?_A:function(){return OS(e),n(i)})}function IS(e){var t=e.memoizedState;if(t!==null)return t;t={memoizedState:ar,baseState:ar,baseQueue:null,queue:{pending:null,lanes:0,dispatch:null,lastRenderedReducer:La,lastRenderedState:ar},next:null};var n={};return t.next={memoizedState:n,baseState:n,baseQueue:null,queue:{pending:null,lanes:0,dispatch:null,lastRenderedReducer:La,lastRenderedState:n},next:null},e.memoizedState=t,e=e.alternate,e!==null&&(e.memoizedState=t),t}function OS(e){var t=IS(e);t.next===null&&(t=e.alternate.memoizedState),Kl(e,t.next.queue,{},fi())}function $g(){return xn(No)}function PS(){return Je().memoizedState}function BS(){return Je().memoizedState}function yA(e){for(var t=e.return;t!==null;){switch(t.tag){case 24:case 3:var n=fi();e=gs(n);var i=vs(t,e,n);i!==null&&(jn(i,t,n),Yl(i,t,n)),t={cache:Bg()},e.payload=t;return}t=t.return}}function xA(e,t,n){var i=fi();n={lane:i,revertLane:0,gesture:null,action:n,hasEagerState:!1,eagerState:null,next:null},sh(e)?FS(t,n):(n=Ig(e,t,n,i),n!==null&&(jn(n,e,i),VS(n,t,i)))}function zS(e,t,n){var i=fi();Kl(e,t,n,i)}function Kl(e,t,n,i){var a={lane:i,revertLane:0,gesture:null,action:n,hasEagerState:!1,eagerState:null,next:null};if(sh(e))FS(t,a);else{var s=e.alternate;if(e.lanes===0&&(s===null||s.lanes===0)&&(s=t.lastRenderedReducer,s!==null))try{var r=t.lastRenderedState,o=s(r,n);if(a.hasEagerState=!0,a.eagerState=o,hi(o,r))return $f(e,t,a,0),Ae===null&&jf(),!1}catch{}if(n=Ig(e,t,a,i),n!==null)return jn(n,e,i),VS(n,t,i),!0}return!1}function t0(e,t,n,i){if(i={lane:2,revertLane:f0(),gesture:null,action:i,hasEagerState:!1,eagerState:null,next:null},sh(e)){if(t)throw Error(J(479))}else t=Ig(e,n,i,2),t!==null&&jn(t,e,2)}function sh(e){var t=e.alternate;return e===Yt||t!==null&&t===Yt}function FS(e,t){go=Nf=!0;var n=e.pending;n===null?t.next=t:(t.next=n.next,n.next=t),e.pending=t}function VS(e,t,n){if((n&4194048)!==0){var i=t.lanes;i&=e.pendingLanes,n|=i,t.lanes=n,Tx(e,n)}}var Lf={readContext:xn,use:ih,useCallback:We,useContext:We,useEffect:We,useImperativeHandle:We,useLayoutEffect:We,useInsertionEffect:We,useMemo:We,useReducer:We,useRef:We,useState:We,useDebugValue:We,useDeferredValue:We,useTransition:We,useSyncExternalStore:We,useId:We,useHostTransitionStatus:We,useFormState:We,useActionState:We,useOptimistic:We,useMemoCache:We,useCacheRefresh:We,useEffectEvent:We},HS={readContext:xn,use:ih,useCallback:function(e,t){return kn().memoizedState=[e,t===void 0?null:t],e},useContext:xn,useEffect:hy,useImperativeHandle:function(e,t,n){n=n!=null?n.concat([e]):null,of(4194308,4,CS.bind(null,t,e),n)},useLayoutEffect:function(e,t){return of(4194308,4,e,t)},useInsertionEffect:function(e,t){of(4,2,e,t)},useMemo:function(e,t){var n=kn();t=t===void 0?null:t;var i=e();if(pr){ls(!0);try{e()}finally{ls(!1)}}return n.memoizedState=[i,t],i},useReducer:function(e,t,n){var i=kn();if(n!==void 0){var a=n(t);if(pr){ls(!0);try{n(t)}finally{ls(!1)}}}else a=t;return i.memoizedState=i.baseState=a,e={pending:null,lanes:0,dispatch:null,lastRenderedReducer:e,lastRenderedState:a},i.queue=e,e=e.dispatch=xA.bind(null,Yt,e),[i.memoizedState,e]},useRef:function(e){var t=kn();return e={current:e},t.memoizedState=e},useState:function(e){e=zm(e);var t=e.queue,n=zS.bind(null,Yt,t);return t.dispatch=n,[e.memoizedState,n]},useDebugValue:Qg,useDeferredValue:function(e,t){var n=kn();return jg(n,e,t)},useTransition:function(){var e=zm(!1);return e=LS.bind(null,Yt,e.queue,!0,!1),kn().memoizedState=e,[!1,e]},useSyncExternalStore:function(e,t,n){var i=Yt,a=kn();if(Qt){if(n===void 0)throw Error(J(407));n=n()}else{if(n=t(),Ae===null)throw Error(J(349));(ee&127)!==0||pS(i,t,n)}a.memoizedState=n;var s={value:n,getSnapshot:t};return a.queue=s,hy(gS.bind(null,i,s,e),[e]),i.flags|=2048,To(9,{destroy:void 0},mS.bind(null,i,s,n,t),null),n},useId:function(){var e=kn(),t=Ae.identifierPrefix;if(Qt){var n=aa,i=ia;n=(i&~(1<<32-ui(i)-1)).toString(32)+n,t="_"+t+"R_"+n,n=Uf++,0<n&&(t+="H"+n.toString(32)),t+="_"}else n=dA++,t="_"+t+"r_"+n.toString(32)+"_";return e.memoizedState=t},useHostTransitionStatus:$g,useFormState:cy,useActionState:cy,useOptimistic:function(e){var t=kn();t.memoizedState=t.baseState=e;var n={pending:null,lanes:0,dispatch:null,lastRenderedReducer:null,lastRenderedState:null};return t.queue=n,t=t0.bind(null,Yt,!0,n),n.dispatch=t,[e,t]},useMemoCache:Zg,useCacheRefresh:function(){return kn().memoizedState=yA.bind(null,Yt)},useEffectEvent:function(e){var t=kn(),n={impl:e};return t.memoizedState=n,function(){if((fe&2)!==0)throw Error(J(440));return n.impl.apply(void 0,arguments)}}},GS={readContext:xn,use:ih,useCallback:DS,useContext:xn,useEffect:Kg,useImperativeHandle:RS,useInsertionEffect:AS,useLayoutEffect:wS,useMemo:NS,useReducer:rf,useRef:TS,useState:function(){return rf(La)},useDebugValue:Qg,useDeferredValue:function(e,t){var n=Je();return US(n,Me.memoizedState,e,t)},useTransition:function(){var e=rf(La)[0],t=Je().memoizedState;return[typeof e=="boolean"?e:Ec(e),t]},useSyncExternalStore:dS,useId:PS,useHostTransitionStatus:$g,useFormState:uy,useActionState:uy,useOptimistic:function(e,t){var n=Je();return yS(n,Me,e,t)},useMemoCache:Zg,useCacheRefresh:BS,useEffectEvent:ES},SA={readContext:xn,use:ih,useCallback:DS,useContext:xn,useEffect:Kg,useImperativeHandle:RS,useInsertionEffect:AS,useLayoutEffect:wS,useMemo:NS,useReducer:Kp,useRef:TS,useState:function(){return Kp(La)},useDebugValue:Qg,useDeferredValue:function(e,t){var n=Je();return Me===null?jg(n,e,t):US(n,Me.memoizedState,e,t)},useTransition:function(){var e=Kp(La)[0],t=Je().memoizedState;return[typeof e=="boolean"?e:Ec(e),t]},useSyncExternalStore:dS,useId:PS,useHostTransitionStatus:$g,useFormState:fy,useActionState:fy,useOptimistic:function(e,t){var n=Je();return Me!==null?yS(n,Me,e,t):(n.baseState=e,[e,n.queue.dispatch])},useMemoCache:Zg,useCacheRefresh:BS,useEffectEvent:ES};function Qp(e,t,n,i){t=e.memoizedState,n=n(i,t),n=n==null?t:we({},t,n),e.memoizedState=n,e.lanes===0&&(e.updateQueue.baseState=n)}var Hm={enqueueSetState:function(e,t,n){e=e._reactInternals;var i=fi(),a=gs(i);a.payload=t,n!=null&&(a.callback=n),t=vs(e,a,i),t!==null&&(jn(t,e,i),Yl(t,e,i))},enqueueReplaceState:function(e,t,n){e=e._reactInternals;var i=fi(),a=gs(i);a.tag=1,a.payload=t,n!=null&&(a.callback=n),t=vs(e,a,i),t!==null&&(jn(t,e,i),Yl(t,e,i))},enqueueForceUpdate:function(e,t){e=e._reactInternals;var n=fi(),i=gs(n);i.tag=2,t!=null&&(i.callback=t),t=vs(e,i,n),t!==null&&(jn(t,e,n),Yl(t,e,n))}};function dy(e,t,n,i,a,s,r){return e=e.stateNode,typeof e.shouldComponentUpdate=="function"?e.shouldComponentUpdate(i,s,r):t.prototype&&t.prototype.isPureReactComponent?!sc(n,i)||!sc(a,s):!0}function py(e,t,n,i){e=t.state,typeof t.componentWillReceiveProps=="function"&&t.componentWillReceiveProps(n,i),typeof t.UNSAFE_componentWillReceiveProps=="function"&&t.UNSAFE_componentWillReceiveProps(n,i),t.state!==e&&Hm.enqueueReplaceState(t,t.state,null)}function mr(e,t){var n=t;if("ref"in t){n={};for(var i in t)i!=="ref"&&(n[i]=t[i])}if(e=e.defaultProps){n===t&&(n=we({},n));for(var a in e)n[a]===void 0&&(n[a]=e[a])}return n}function kS(e){Mf(e)}function XS(e){console.error(e)}function WS(e){Mf(e)}function If(e,t){try{var n=e.onUncaughtError;n(t.value,{componentStack:t.stack})}catch(i){setTimeout(function(){throw i})}}function my(e,t,n){try{var i=e.onCaughtError;i(n.value,{componentStack:n.stack,errorBoundary:t.tag===1?t.stateNode:null})}catch(a){setTimeout(function(){throw a})}}function Gm(e,t,n){return n=gs(n),n.tag=3,n.payload={element:null},n.callback=function(){If(e,t)},n}function qS(e){return e=gs(e),e.tag=3,e}function YS(e,t,n,i){var a=n.type.getDerivedStateFromError;if(typeof a=="function"){var s=i.value;e.payload=function(){return a(s)},e.callback=function(){my(t,n,i)}}var r=n.stateNode;r!==null&&typeof r.componentDidCatch=="function"&&(e.callback=function(){my(t,n,i),typeof a!="function"&&(xs===null?xs=new Set([this]):xs.add(this));var o=i.stack;this.componentDidCatch(i.value,{componentStack:o!==null?o:""})})}function bA(e,t,n,i,a){if(n.flags|=32768,i!==null&&typeof i=="object"&&typeof i.then=="function"){if(t=n.alternate,t!==null&&ur(t,n,a,!0),n=Tn.current,n!==null){switch(n.tag){case 31:case 13:case 19:return Un===null?Gf():n.alternate===null&&qe===0&&(qe=3),n.flags&=-257,n.flags|=65536,n.lanes=a,i===Cf?n.flags|=16384:(t=n.updateQueue,t===null?n.updateQueue=new Set([i]):t.add(i),am(e,i,a)),!1;case 22:return n.flags|=65536,i===Cf?n.flags|=16384:(t=n.updateQueue,t===null?(t={transitions:null,markerInstances:null,retryQueue:new Set([i])},n.updateQueue=t):(n=t.retryQueue,n===null?t.retryQueue=new Set([i]):n.add(i)),am(e,i,a)),!1}throw Error(J(435,n.tag))}return am(e,i,a),Gf(),!1}if(Qt)return t=Tn.current,t!==null?((t.flags&65536)===0&&(t.flags|=256),t.flags|=65536,t.lanes=a,i!==Dm&&(e=Error(J(422),{cause:i}),oc(Ai(e,n)))):(i!==Dm&&(t=Error(J(423),{cause:i}),oc(Ai(t,n))),e=e.current.alternate,e.flags|=65536,a&=-a,e.lanes|=a,i=Ai(i,n),a=Gm(e.stateNode,i,a),Jp(e,a),qe!==4&&(qe=2)),!1;var s=Error(J(520),{cause:i});if(s=Ai(s,n),tc===null?tc=[s]:tc.push(s),qe!==4&&(qe=2),t===null)return!0;i=Ai(i,n),n=t;do{switch(n.tag){case 3:return n.flags|=65536,e=a&-a,n.lanes|=e,e=Gm(n.stateNode,i,e),Jp(n,e),!1;case 1:if(t=n.type,s=n.stateNode,(n.flags&128)===0&&(typeof t.getDerivedStateFromError=="function"||s!==null&&typeof s.componentDidCatch=="function"&&(xs===null||!xs.has(s))))return n.flags|=65536,a&=-a,n.lanes|=a,a=qS(a),YS(a,e,n,i),Jp(n,a),!1;break;case 22:if(n.memoizedState!==null)return n.flags|=65536,!1}n=n.return}while(n!==null);return!1}var e0=Error(J(461)),je=!1;function en(e,t,n,i){t.child=e===null?oS(t,null,n,i):dr(t,e.child,n,i)}function gy(e,t,n,i,a){n=n.render;var s=t.ref;if("ref"in i){var r={};for(var o in i)o!=="ref"&&(r[o]=i[o])}else r=i;return fr(t),i=Xg(e,t,n,r,s,a),o=Wg(),e!==null&&!je?(qg(e,t,a),Ia(e,t,a)):(Qt&&o&&th(t),t.flags|=1,en(e,t,i,a),t.child)}function vy(e,t,n,i,a){if(e===null){var s=n.type;return typeof s=="function"&&!Og(s)&&s.defaultProps===void 0&&n.compare===null?(t.tag=15,t.type=s,ZS(e,t,s,i,a)):(e=nf(n.type,null,i,t,t.mode,a),e.ref=t.ref,e.return=t,t.child=e)}if(s=e.child,!i0(e,a)){var r=s.memoizedProps;if(n=n.compare,n=n!==null?n:sc,n(r,i)&&e.ref===t.ref)return Ia(e,t,a)}return t.flags|=1,e=Ca(s,i),e.ref=t.ref,e.return=t,t.child=e}function ZS(e,t,n,i,a){if(e!==null){var s=e.memoizedProps;if(sc(s,i)&&e.ref===t.ref)if(je=!1,t.pendingProps=i=s,i0(e,a))(e.flags&131072)!==0&&(je=!0);else return t.lanes=e.lanes,Ia(e,t,a)}return km(e,t,n,i,a)}function JS(e,t,n,i){var a=i.children,s=e!==null?e.memoizedState:null;if(e===null&&t.stateNode===null&&(t.stateNode={_visibility:1,_pendingMarkers:null,_retryCache:null,_transitions:null}),i.mode==="hidden"){if((t.flags&128)!==0){if(s=s!==null?s.baseLanes|n:n,e!==null){for(i=t.child=e.child,a=0;i!==null;)a=a|i.lanes|i.childLanes,i=i.sibling;i=a&~s}else i=0,t.child=null;return _y(e,t,s,n,i)}if((n&536870912)!==0)t.memoizedState={baseLanes:0,cachePool:null},e!==null&&sf(t,s!==null?s.cachePool:null),s!==null?ry(t,s):Pm(),uS(t);else return i=t.lanes=536870912,_y(e,t,s!==null?s.baseLanes|n:n,n,i)}else s!==null?(sf(t,s.cachePool),ry(t,s),ys(),t.memoizedState=null):(e!==null&&sf(t,null),Pm(),ys());return en(e,t,a,n),t.child}function Ql(e,t){return e!==null&&e.tag===22||t.stateNode!==null||(t.stateNode={_visibility:1,_pendingMarkers:null,_retryCache:null,_transitions:null}),t.sibling}function _y(e,t,n,i,a){var s=zg();return s=s===null?null:{parent:Qe._currentValue,pool:s},t.memoizedState={baseLanes:n,cachePool:s},e!==null&&sf(t,null),Pm(),uS(t),e!==null&&ur(e,t,i,!0),t.childLanes=a,null}function lf(e,t){return t=rh({mode:t.mode,children:t.children},e.mode),t.ref=e.ref,e.child=t,t.return=e,t}function yy(e,t,n){return dr(t,e.child,null,n),e=lf(t,t.pendingProps),e.flags|=2,si(t),t.memoizedState=null,e}function MA(e,t,n){var i=t.pendingProps,a=(t.flags&128)!==0;if(t.flags&=-129,e===null){if(Qt){if(i.mode==="hidden")return e=lf(t,i),t.lanes=536870912,e.memoizedState={baseLanes:0,cachePool:null},Ql(null,e);if(Bm(t),(e=Ue)?(e=Jb(e,wi),e=e!==null&&e.data==="&"?e:null,e!==null&&(t.memoizedState={dehydrated:e,treeContext:Ts!==null?{id:ia,overflow:aa}:null,retryLane:536870912,hydrationErrors:null},n=tS(e),n.return=t,t.child=n,hn=t,Ue=null)):e=null,e===null)throw Es(t);return t.lanes=536870912,null}return lf(t,i)}var s=e.memoizedState;if(s!==null){var r=s.dehydrated;if(Bm(t),a)if(t.flags&256)t.flags&=-257,t=yy(e,t,n);else if(t.memoizedState!==null)t.child=e.child,t.flags|=128,t=null;else throw Error(J(558));else if(je||ur(e,t,n,!1),a=(n&e.childLanes)!==0,je||a){if(As.current===null){if(i=Ae,i!==null&&(r=Ex(i,n),r!==0&&r!==s.retryLane))throw s.retryLane=r,xr(e,r),jn(i,e,r),e0;Gf()}t=yy(e,t,n)}else e=s.treeContext,Ue=Ci(r.nextSibling),hn=t,Qt=!0,ms=null,wi=!1,e!==null&&nS(t,e),t=lf(t,i),t.flags|=134221824;return t}return e=Ca(e.child,{mode:i.mode,children:i.children}),e.ref=t.ref,t.child=e,e.return=t,e}function Jr(e,t){var n=t.ref;if(n===null)e!==null&&e.ref!==null&&(t.flags|=4194816);else{if(typeof n!="function"&&typeof n!="object")throw Error(J(284));(e===null||e.ref!==n)&&(t.flags|=4194816)}}function km(e,t,n,i,a){return fr(t),n=Xg(e,t,n,i,void 0,a),i=Wg(),e!==null&&!je?(qg(e,t,a),Ia(e,t,a)):(Qt&&i&&th(t),t.flags|=1,en(e,t,n,a),t.child)}function xy(e,t,n,i,a,s){return fr(t),t.updateQueue=null,n=hS(t,i,n,a),fS(e),i=Wg(),e!==null&&!je?(qg(e,t,s),Ia(e,t,s)):(Qt&&i&&th(t),t.flags|=1,en(e,t,n,s),t.child)}function Sy(e,t,n,i,a){if(fr(t),t.stateNode===null){var s=ro,r=n.contextType;typeof r=="object"&&r!==null&&(s=xn(r)),s=new n(i,s),t.memoizedState=s.state!==null&&s.state!==void 0?s.state:null,s.updater=Hm,t.stateNode=s,s._reactInternals=t,s=t.stateNode,s.props=i,s.state=t.memoizedState,s.refs={},Vg(t),r=n.contextType,s.context=typeof r=="object"&&r!==null?xn(r):ro,s.state=t.memoizedState,r=n.getDerivedStateFromProps,typeof r=="function"&&(Qp(t,n,r,i),s.state=t.memoizedState),typeof n.getDerivedStateFromProps=="function"||typeof s.getSnapshotBeforeUpdate=="function"||typeof s.UNSAFE_componentWillMount!="function"&&typeof s.componentWillMount!="function"||(r=s.state,typeof s.componentWillMount=="function"&&s.componentWillMount(),typeof s.UNSAFE_componentWillMount=="function"&&s.UNSAFE_componentWillMount(),r!==s.state&&Hm.enqueueReplaceState(s,s.state,null),Jl(t,i,s,a),Zl(),s.state=t.memoizedState),typeof s.componentDidMount=="function"&&(t.flags|=4194308),i=!0}else if(e===null){s=t.stateNode;var o=t.memoizedProps,l=mr(n,o);s.props=l;var c=s.context,f=n.contextType;r=ro,typeof f=="object"&&f!==null&&(r=xn(f));var p=n.getDerivedStateFromProps;f=typeof p=="function"||typeof s.getSnapshotBeforeUpdate=="function",o=t.pendingProps!==o,f||typeof s.UNSAFE_componentWillReceiveProps!="function"&&typeof s.componentWillReceiveProps!="function"||(o||c!==r)&&py(t,s,i,r),rs=!1;var u=t.memoizedState;s.state=u,Jl(t,i,s,a),Zl(),c=t.memoizedState,o||u!==c||rs?(typeof p=="function"&&(Qp(t,n,p,i),c=t.memoizedState),(l=rs||dy(t,n,l,i,u,c,r))?(f||typeof s.UNSAFE_componentWillMount!="function"&&typeof s.componentWillMount!="function"||(typeof s.componentWillMount=="function"&&s.componentWillMount(),typeof s.UNSAFE_componentWillMount=="function"&&s.UNSAFE_componentWillMount()),typeof s.componentDidMount=="function"&&(t.flags|=4194308)):(typeof s.componentDidMount=="function"&&(t.flags|=4194308),t.memoizedProps=i,t.memoizedState=c),s.props=i,s.state=c,s.context=r,i=l):(typeof s.componentDidMount=="function"&&(t.flags|=4194308),i=!1)}else{s=t.stateNode,Im(e,t),r=t.memoizedProps,f=mr(n,r),s.props=f,p=t.pendingProps,u=s.context,c=n.contextType,l=ro,typeof c=="object"&&c!==null&&(l=xn(c)),o=n.getDerivedStateFromProps,(c=typeof o=="function"||typeof s.getSnapshotBeforeUpdate=="function")||typeof s.UNSAFE_componentWillReceiveProps!="function"&&typeof s.componentWillReceiveProps!="function"||(r!==p||u!==l)&&py(t,s,i,l),rs=!1,u=t.memoizedState,s.state=u,Jl(t,i,s,a),Zl();var d=t.memoizedState;r!==p||u!==d||rs||e!==null&&e.dependencies!==null&&wf(e.dependencies)?(typeof o=="function"&&(Qp(t,n,o,i),d=t.memoizedState),(f=rs||dy(t,n,f,i,u,d,l)||e!==null&&e.dependencies!==null&&wf(e.dependencies))?(c||typeof s.UNSAFE_componentWillUpdate!="function"&&typeof s.componentWillUpdate!="function"||(typeof s.componentWillUpdate=="function"&&s.componentWillUpdate(i,d,l),typeof s.UNSAFE_componentWillUpdate=="function"&&s.UNSAFE_componentWillUpdate(i,d,l)),typeof s.componentDidUpdate=="function"&&(t.flags|=4),typeof s.getSnapshotBeforeUpdate=="function"&&(t.flags|=1024)):(typeof s.componentDidUpdate!="function"||r===e.memoizedProps&&u===e.memoizedState||(t.flags|=4),typeof s.getSnapshotBeforeUpdate!="function"||r===e.memoizedProps&&u===e.memoizedState||(t.flags|=1024),t.memoizedProps=i,t.memoizedState=d),s.props=i,s.state=d,s.context=l,i=f):(typeof s.componentDidUpdate!="function"||r===e.memoizedProps&&u===e.memoizedState||(t.flags|=4),typeof s.getSnapshotBeforeUpdate!="function"||r===e.memoizedProps&&u===e.memoizedState||(t.flags|=1024),i=!1)}return s=i,Jr(e,t),i=(t.flags&128)!==0,s||i?(s=t.stateNode,n=i&&typeof n.getDerivedStateFromError!="function"?null:s.render(),t.flags|=1,e!==null&&i?(t.child=dr(t,e.child,null,a),t.child=dr(t,null,n,a)):en(e,t,n,a),t.memoizedState=s.state,e=t.child):e=Ia(e,t,a),e}function by(e,t,n,i){return cr(),t.flags|=256,en(e,t,n,i),t.child}var Xm={dehydrated:null,treeContext:null,retryLane:0,hydrationErrors:null};function Wm(e){return{baseLanes:e,cachePool:aS()}}function qm(e,t,n){return e=e!==null?e.childLanes&~n:0,t&&(e|=oi),e}function KS(e,t,n){var i=t.pendingProps,a=!1,s=(t.flags&128)!==0,r;if((r=s)||(r=e!==null&&e.memoizedState===null?!1:(bn.current&2)!==0),r&&(a=!0,t.flags&=-129),r=(t.flags&32)!==0,t.flags&=-33,e===null){if(Qt){if(a?_s(t):ys(),(e=Ue)?(e=Jb(e,wi),e=e!==null&&e.data!=="&"?e:null,e!==null&&(t.memoizedState={dehydrated:e,treeContext:Ts!==null?{id:ia,overflow:aa}:null,retryLane:536870912,hydrationErrors:null},n=tS(e),n.return=t,t.child=n,hn=t,Ue=null)):e=null,e===null)throw Es(t);return m0(e)?t.lanes=32:t.lanes=536870912,null}return s=i.children,i=i.fallback,a?(ys(),a=t.mode,s=rh({mode:"hidden",children:s},a),i=sr(i,a,n,null),s.return=t,i.return=t,s.sibling=i,t.child=s,i=t.child,i.memoizedState=Wm(n),i.childLanes=qm(e,r,n),t.memoizedState=Xm,Ql(null,i)):(_s(t),n0(t,s))}var o=e.memoizedState;if(o!==null){var l=o.dehydrated;if(l!==null)return TA(e,t,s,r,i,l,o,n)}return a?(ys(),a=i.fallback,s=t.mode,o=e.child,l=o.sibling,i=Ca(o,{mode:"hidden",children:i.children}),i.subtreeFlags=o.subtreeFlags&1206910976,l!==null?a=Ca(l,a):(a=sr(a,s,n,null),a.flags|=2),a.return=t,i.return=t,i.sibling=a,t.child=i,Ql(null,i),i=t.child,a=e.child.memoizedState,a===null?a=Wm(n):(s=a.cachePool,s!==null?(o=Qe._currentValue,s=s.parent!==o?{parent:o,pool:o}:s):s=aS(),a={baseLanes:a.baseLanes|n,cachePool:s}),i.memoizedState=a,i.childLanes=qm(e,r,n),t.memoizedState=Xm,Ql(e.child,i)):(_s(t),n=e.child,e=n.sibling,n=Ca(n,{mode:"visible",children:i.children}),n.return=t,n.sibling=null,e!==null&&(r=t.deletions,r===null?(t.deletions=[e],t.flags|=16):r.push(e)),t.child=n,t.memoizedState=null,n)}function n0(e,t){return t=rh({mode:"visible",children:t},e.mode),t.return=e,e.child=t}function rh(e,t){return e=Qn(22,e,null,t),e.lanes=0,e}function Gu(e,t,n){return dr(t,e.child,null,n),e=n0(t,t.pendingProps.children),e.flags|=2,t.memoizedState=null,e}function TA(e,t,n,i,a,s,r,o){if(n)return t.flags&256?(_s(t),t.flags&=-257,Gu(e,t,o)):t.memoizedState!==null?(ys(),t.child=e.child,t.flags|=128,null):(ys(),s=a.fallback,r=t.mode,a=rh({mode:"visible",children:a.children},r),s=sr(s,r,o,null),s.flags|=2,a.return=t,s.return=t,a.sibling=s,t.child=a,dr(t,e.child,null,o),a=t.child,a.memoizedState=Wm(o),a.childLanes=qm(e,i,o),t.memoizedState=Xm,Ql(null,a));if(_s(t),m0(s)){if(i=s.nextSibling&&s.nextSibling.dataset,i)var l=i.dgst;return i=l,i!==""&&(a=Error(J(419)),a.stack="",a.digest=i,oc({value:a,source:null,stack:null})),Gu(e,t,o)}if(je||ur(e,t,o,!1),i=(o&e.childLanes)!==0,je||i){if(As.current!==null)return Gu(e,t,o);if(i=Ae,i!==null&&(a=Ex(i,o),a!==0&&a!==r.retryLane))throw r.retryLane=a,xr(e,a),jn(i,e,a),e0;return vg(s)||Gf(),Gu(e,t,o)}return vg(s)?(t.flags|=192,t.child=e.child,null):(e=r.treeContext,Ue=Ci(s.nextSibling),hn=t,Qt=!0,ms=null,wi=!1,e!==null&&nS(t,e),t=n0(t,a.children),t.flags|=134221824,t)}function My(e,t,n){e.lanes|=t;var i=e.alternate;i!==null&&(i.lanes|=t),af(e.return,t,n)}function Ty(e){for(var t=null;e!==null;){var n=e.alternate;n!==null&&Df(n)===null&&(t=e),e=e.sibling}return t}function ku(e,t,n,i,a,s){var r=e.memoizedState;r===null?e.memoizedState={isBackwards:t,rendering:null,renderingStartTime:0,last:i,tail:n,tailMode:a,treeForkCount:s}:(r.isBackwards=t,r.rendering=null,r.renderingStartTime=0,r.last=i,r.tail=n,r.tailMode=a,r.treeForkCount=s)}function jp(e){var t=e.child;for(e.child=null;t!==null;){var n=t.sibling;t.sibling=e.child,e.child=t,t=n}}function Ym(e,t,n){var i=t.pendingProps,a=i.revealOrder,s=i.tail;i=i.children;var r=bn.current;if(t.flags&128)return cc(t,r),null;var o=(r&2)!==0;if(o?(r=r&1|2,t.flags|=128):r&=1,cc(t,r),a==="backwards"&&e!==null?(jp(e),en(e,t,i,n),jp(e)):en(e,t,i,n),i=Qt?rc:0,!o&&e!==null&&(e.flags&128)!==0)t:for(e=t.child;e!==null;){if(e.tag===13)e.memoizedState!==null&&My(e,n,t);else if(e.tag===19)My(e,n,t);else if(e.child!==null){e.child.return=e,e=e.child;continue}if(e===t)break t;for(;e.sibling===null;){if(e.return===null||e.return===t)break t;e=e.return}e.sibling.return=e.return,e=e.sibling}switch(a){case"backwards":n=Ty(t.child),n===null?(a=t.child,t.child=null):(a=n.sibling,n.sibling=null,jp(t)),ku(t,!0,a,null,s,i);break;case"unstable_legacy-backwards":for(n=null,a=t.child,t.child=null;a!==null;){if(e=a.alternate,e!==null&&Df(e)===null){t.child=a;break}e=a.sibling,a.sibling=n,n=a,a=e}ku(t,!0,n,null,s,i);break;case"together":ku(t,!1,null,null,void 0,i);break;case"independent":t.memoizedState=null;break;default:n=Ty(t.child),n===null?(a=t.child,t.child=null):(a=n.sibling,n.sibling=null),ku(t,!1,a,n,s,i)}return t.child}function Ey(e,t,n){var i=t.pendingProps;return us(t,t.type,i.value),en(e,t,i.children,n),t.child}function Ia(e,t,n){if(e!==null&&(t.dependencies=e.dependencies),Cs|=t.lanes,(n&t.childLanes)===0)if(e!==null){if(ur(e,t,n,!1),(n&t.childLanes)===0)return null}else return null;if(e!==null&&t.child!==e.child)throw Error(J(153));if(t.child!==null){for(e=t.child,n=Ca(e,e.pendingProps),t.child=n,n.return=t;e.sibling!==null;)e=e.sibling,n=n.sibling=Ca(e,e.pendingProps),n.return=t;n.sibling=null}return t.child}function i0(e,t){return(e.lanes&t)!==0?!0:(e=e.dependencies,!!(e!==null&&wf(e)))}function EA(e,t,n){switch(t.tag){case 3:yf(t,t.stateNode.containerInfo),us(t,Qe,e.memoizedState.cache),cr();break;case 27:case 5:xm(t);break;case 4:yf(t,t.stateNode.containerInfo);break;case 10:us(t,t.type,t.memoizedProps.value);break;case 31:if(t.memoizedState!==null)return t.flags|=128,Bm(t),null;break;case 13:var i=t.memoizedState;if(i!==null){if(i.dehydrated!==null)return _s(t),t.flags|=128,null;i=ur(e,t,n,!1);var a=t.child.childLanes;return i||(n&a)!==0?KS(e,t,n):(_s(t),e=Ia(e,t,n),e!==null?e.sibling:null)}_s(t);break;case 19:if(t.flags&128)return Ym(e,t,n);if(a=(e.flags&128)!==0,i=(n&t.childLanes)!==0,i||(ur(e,t,n,!1),i=(n&t.childLanes)!==0),a){if(i)return Ym(e,t,n);t.flags|=128}if(a=t.memoizedState,a!==null&&(a.rendering=null,a.tail=null,a.lastEffect=null),cc(t,bn.current),i)break;return null;case 22:return t.lanes=0,JS(e,t,n,t.pendingProps);case 24:us(t,Qe,e.memoizedState.cache)}return Ia(e,t,n)}function QS(e,t,n){if(e!==null)if(e.memoizedProps!==t.pendingProps)je=!0;else{if(!i0(e,n)&&(t.flags&128)===0)return je=!1,EA(e,t,n);je=(e.flags&131072)!==0}else je=!1,Qt&&(t.flags&1048576)!==0&&eS(t,rc,t.index);switch(t.lanes=0,t.tag){case 16:t:{var i=t.pendingProps;if(e=tr(t.elementType),t.type=e,typeof e=="function")Og(e)?(i=mr(e,i),t.tag=1,t=Sy(null,t,e,i,n)):(t.tag=0,t=km(null,t,e,i,n));else{if(e!=null){var a=e.$$typeof;if(a===bg){t.tag=11,t=gy(null,t,e,i,n);break t}else if(a===Mg){t.tag=14,t=vy(null,t,e,i,n);break t}else if(a===ea){t.tag=10,t.type=e,t=Ey(null,t,n);break t}}throw t=_m(e)||e,Error(J(306,t,""))}}return t;case 0:return km(e,t,t.type,t.pendingProps,n);case 1:return i=t.type,a=mr(i,t.pendingProps),Sy(e,t,i,a,n);case 3:t:{if(yf(t,t.stateNode.containerInfo),e===null)throw Error(J(387));i=t.pendingProps;var s=t.memoizedState;a=s.element,Im(e,t),Jl(t,i,null,n);var r=t.memoizedState;if(i=r.cache,us(t,Qe,i),i!==s.cache&&Um(t,[Qe],n,!0),Zl(),i=r.element,s.isDehydrated)if(s={element:i,isDehydrated:!1,cache:r.cache},t.updateQueue.baseState=s,t.memoizedState=s,t.flags&256){t=by(e,t,i,n);break t}else if(i!==a){a=Ai(Error(J(424)),t),oc(a),t=by(e,t,i,n);break t}else for(e=t.stateNode.containerInfo,e.nodeType===9?e=e.body:e=e.nodeName==="HTML"?e.ownerDocument.body:e,Ue=Ci(e.firstChild),hn=t,Qt=!0,ms=null,wi=!0,n=oS(t,null,i,n),t.child=n;n;)n.flags=n.flags&-3|134221824,n=n.sibling;else{if(cr(),i===a){t=Ia(e,t,n);break t}en(e,t,i,n)}t=t.child}return t;case 26:return Jr(e,t),e===null?(n=jy(t.type,null,t.pendingProps,null))?t.memoizedState=n:Qt||(t.stateNode=Vb(t.type,t.pendingProps,ps.current,t)):t.memoizedState=jy(t.type,e.memoizedProps,t.pendingProps,e.memoizedState),null;case 27:return xm(t),e===null&&Qt&&(i=t.stateNode=Kb(t.type,t.pendingProps,ps.current),hn=t,wi=!0,a=Ue,Ds(t.type)?(_g=a,Ue=Ci(i.firstChild)):Ue=a),en(e,t,t.pendingProps.children,n),Jr(e,t),e===null&&(t.flags|=4194304),t.child;case 5:return e===null&&Qt&&((a=i=Ue)&&(i=vw(i,t.type,t.pendingProps,wi),i!==null?(t.stateNode=i,hn=t,Ue=Ci(i.firstChild),wi=!1,a=!0):a=!1),a||Es(t)),xm(t),a=t.type,s=t.pendingProps,r=e!==null?e.memoizedProps:null,i=s.children,pg(a,s)?i=null:r!==null&&pg(a,r)&&(t.flags|=32),t.memoizedState!==null&&(a=Xg(e,t,pA,null,null,n),No._currentValue=a),Jr(e,t),en(e,t,i,n),t.child;case 6:return e===null&&Qt&&((e=n=Ue)&&(n=_w(n,t.pendingProps,wi),n!==null?(t.stateNode=n,hn=t,Ue=null,e=!0):e=!1),e||Es(t)),null;case 13:return KS(e,t,n);case 4:return yf(t,t.stateNode.containerInfo),i=t.pendingProps,e===null?t.child=dr(t,null,i,n):en(e,t,i,n),t.child;case 11:return gy(e,t,t.type,t.pendingProps,n);case 7:return i=t.pendingProps,Jr(e,t),en(e,t,i,n),t.child;case 8:return en(e,t,t.pendingProps.children,n),t.child;case 12:return en(e,t,t.pendingProps.children,n),t.child;case 10:return Ey(e,t,n);case 9:return a=t.type._context,i=t.pendingProps.children,fr(t),a=xn(a),i=i(a),t.flags|=1,en(e,t,i,n),t.child;case 14:return vy(e,t,t.type,t.pendingProps,n);case 15:return ZS(e,t,t.type,t.pendingProps,n);case 19:return Ym(e,t,n);case 31:return MA(e,t,n);case 22:return JS(e,t,n,t.pendingProps);case 24:return fr(t),i=xn(Qe),e===null?(a=zg(),a===null&&(a=Ae,s=Bg(),a.pooledCache=s,s.refCount++,s!==null&&(a.pooledCacheLanes|=n),a=s),t.memoizedState={parent:i,cache:a},Vg(t),us(t,Qe,a)):((e.lanes&n)!==0&&(Im(e,t),Jl(t,null,null,n),Zl()),a=e.memoizedState,s=t.memoizedState,a.parent!==i?(a={parent:i,cache:i},t.memoizedState=a,t.lanes===0&&(t.memoizedState=t.updateQueue.baseState=a),us(t,Qe,i)):(i=s.cache,us(t,Qe,i),i!==a.cache&&Um(t,[Qe],n,!0))),en(e,t,t.pendingProps.children,n),t.child;case 30:return t.stateNode===null&&(t.stateNode={autoName:null,paired:null,clones:null,ref:null}),i=t.pendingProps,i.name!=null&&i.name!=="auto"?t.flags|=e===null?18882560:18874368:Qt&&th(t),e!==null&&e.memoizedProps.name!==i.name?t.flags|=4194816:Jr(e,t),en(e,t,i.children,n),t.child;case 29:throw t.pendingProps}throw Error(J(156,t.tag))}function Ea(e){e.flags|=4}function $p(e,t,n,i,a){var s;if((s=(e.mode&32)!==0)&&(s=n===null?ex(t,i):ex(t,i)&&(i.src!==n.src||i.srcSet!==n.srcSet)),s){if(e.flags|=16777216,(a&335544128)===a)if(e.stateNode.complete)e.flags|=8192;else if(Ab())e.flags|=8192;else throw or=Cf,Fg}else e.flags&=-16777217}function Ay(e,t){if(t.type!=="stylesheet"||(t.state.loading&4)!==0)e.flags&=-16777217;else if(e.flags|=16777216,!tM(t))if(Ab())e.flags|=8192;else throw or=Cf,Fg}function Xu(e,t){t!==null&&(e.flags|=4),e.flags&16384&&(t=e.tag!==22?bx():536870912,e.lanes|=t,Eo|=t)}function Il(e,t){if(!Qt)switch(e.tailMode){case"visible":break;case"collapsed":for(var n=e.tail,i=null;n!==null;)n.alternate!==null&&(i=n),n=n.sibling;i===null?t||e.tail===null?e.tail=null:e.tail.sibling=null:i.sibling=null;break;default:for(t=e.tail,n=null;t!==null;)t.alternate!==null&&(n=t),t=t.sibling;n===null?e.tail=null:n.sibling=null}}function Ne(e){var t=e.alternate!==null&&e.alternate.child===e.child,n=0,i=0;if(t)for(var a=e.child;a!==null;)n|=a.lanes|a.childLanes,i|=a.subtreeFlags&1206910976,i|=a.flags&1206910976,a.return=e,a=a.sibling;else for(a=e.child;a!==null;)n|=a.lanes|a.childLanes,i|=a.subtreeFlags,i|=a.flags,a.return=e,a=a.sibling;return e.subtreeFlags|=i,e.childLanes=n,t}function AA(e,t,n){var i=t.pendingProps;switch(Pg(t),t.tag){case 16:case 15:case 0:case 11:case 7:case 8:case 12:case 9:case 14:return Ne(t),null;case 1:return Ne(t),null;case 3:return n=t.stateNode,i=null,e!==null&&(i=e.memoizedState.cache),t.memoizedState.cache!==i&&(t.flags|=2048),Ra(Qe),So(),n.pendingContext&&(n.context=n.pendingContext,n.pendingContext=null),(e===null||e.child===null)&&(Yr(t)?Ea(t):e===null||e.memoizedState.isDehydrated&&(t.flags&256)===0||(t.flags|=1024,Zp())),Ne(t),null;case 26:var a=t.type,s=t.memoizedState;return e===null?(Ea(t),s!==null?(Ne(t),Ay(t,s)):(Ne(t),$p(t,a,null,i,n))):s?s!==e.memoizedState?(Ea(t),Ne(t),Ay(t,s)):(Ne(t),t.flags&=-16777217):(e=e.memoizedProps,e!==i&&Ea(t),Ne(t),$p(t,a,e,i,n)),null;case 27:if(xf(t),n=ps.current,a=t.type,e!==null&&t.stateNode!=null)e.memoizedProps!==i&&Ea(t);else{if(!i){if(t.stateNode===null)throw Error(J(166));return Ne(t),t.subtreeFlags&=-33554433,null}e=sa.current,Yr(t)?$_(t,e):(e=Kb(a,i,n),t.stateNode=e,Ea(t))}return Ne(t),t.subtreeFlags&=-33554433,null;case 5:if(xf(t),a=t.type,e!==null&&t.stateNode!=null)e.memoizedProps!==i&&Ea(t);else{if(!i){if(t.stateNode===null)throw Error(J(166));return Ne(t),t.subtreeFlags&=-33554433,null}if(s=sa.current,Yr(t))$_(t,s);else{var r=dc(ps.current);switch(s){case 1:s=r.createElementNS("http://www.w3.org/2000/svg",a);break;case 2:s=r.createElementNS("http://www.w3.org/1998/Math/MathML",a);break;default:switch(a){case"svg":s=r.createElementNS("http://www.w3.org/2000/svg",a);break;case"math":s=r.createElementNS("http://www.w3.org/1998/Math/MathML",a);break;case"script":s=r.createElement("div"),s.innerHTML="<script><\/script>",s=s.removeChild(s.firstChild);break;case"select":s=typeof i.is=="string"?r.createElement("select",{is:i.is}):r.createElement("select"),i.multiple?s.multiple=!0:i.size&&(s.size=i.size);break;default:s=typeof i.is=="string"?r.createElement(a,{is:i.is}):r.createElement(a)}}s[yn]=t,s[ti]=i;t:for(r=t.child;r!==null;){if(r.tag===5||r.tag===6)s.appendChild(r.stateNode);else if(r.tag!==4&&r.tag!==27&&r.child!==null){r.child.return=r,r=r.child;continue}if(r===t)break t;for(;r.sibling===null;){if(r.return===null||r.return===t)break t;r=r.return}r.sibling.return=r.return,r=r.sibling}t.stateNode=s;t:switch(Mn(s,a,i),a){case"button":case"input":case"select":case"textarea":i=!!i.autoFocus;break t;case"img":i=!0;break t;default:i=!1}i&&Ea(t)}}return Ne(t),t.subtreeFlags&=-33554433,$p(t,t.type,e===null?null:e.memoizedProps,t.pendingProps,n),null;case 6:if(e&&t.stateNode!=null)e.memoizedProps!==i&&Ea(t);else{if(typeof i!="string"&&t.stateNode===null)throw Error(J(166));if(e=ps.current,Yr(t)){if(e=t.stateNode,n=t.memoizedProps,i=null,a=hn,a!==null)switch(a.tag){case 27:case 5:i=a.memoizedProps}e[yn]=t,e=!!(e.nodeValue===n||i!==null&&i.suppressHydrationWarning===!0||zb(e.nodeValue,n)),e||Es(t,!0)}else e=dc(e).createTextNode(i),e[yn]=t,t.stateNode=e}return Ne(t),null;case 31:if(n=t.memoizedState,e===null||e.memoizedState!==null){if(i=Yr(t),n!==null){if(e===null){if(!i)throw Error(J(318));if(e=t.memoizedState,e=e!==null?e.dehydrated:null,!e)throw Error(J(557));e[yn]=t}else cr(),(t.flags&128)===0&&(t.memoizedState=null),t.flags|=4;Ne(t),e=!1}else n=Zp(),e!==null&&e.memoizedState!==null&&(e.memoizedState.hydrationErrors=n),e=!0;if(!e)return t.flags&256?(si(t),t):(si(t),null);if((t.flags&128)!==0)throw Error(J(558))}return Ne(t),null;case 13:if(i=t.memoizedState,e===null||e.memoizedState!==null&&e.memoizedState.dehydrated!==null){if(a=Yr(t),i!==null&&i.dehydrated!==null){if(e===null){if(!a)throw Error(J(318));if(a=t.memoizedState,a=a!==null?a.dehydrated:null,!a)throw Error(J(317));a[yn]=t}else cr(),(t.flags&128)===0&&(t.memoizedState=null),t.flags|=4;Ne(t),a=!1}else a=Zp(),e!==null&&e.memoizedState!==null&&(e.memoizedState.hydrationErrors=a),a=!0;if(!a)return t.flags&256?(si(t),t):(si(t),null)}return si(t),(t.flags&128)!==0?(t.lanes=n,t):(n=i!==null,e=e!==null&&e.memoizedState!==null,n&&(i=t.child,a=null,i.alternate!==null&&i.alternate.memoizedState!==null&&i.alternate.memoizedState.cachePool!==null&&(a=i.alternate.memoizedState.cachePool.pool),s=null,i.memoizedState!==null&&i.memoizedState.cachePool!==null&&(s=i.memoizedState.cachePool.pool),s!==a&&(i.flags|=2048)),n!==e&&n&&(t.child.flags|=8192),Xu(t,t.updateQueue),Ne(t),null);case 4:return So(),e===null&&h0(t.stateNode.containerInfo),t.flags|=67108864,Ne(t),null;case 10:return Ra(t.type),Ne(t),null;case 19:if(Gg(t),i=t.memoizedState,i===null)return Ne(t),null;if(a=(t.flags&128)!==0,s=i.rendering,s===null)if(a)Il(i,!1);else{if(qe!==0||e!==null&&(e.flags&128)!==0)for(e=t.child;e!==null;){if(s=Df(e),s!==null){for(t.flags|=128,Il(i,!1),e=s.updateQueue,t.updateQueue=e,Xu(t,e),t.subtreeFlags=0,e=n,n=t.child;n!==null;)$x(n,e),n=n.sibling;return cc(t,bn.current&1|2),Qt&&Aa(t,i.treeForkCount),t.child}e=e.sibling}i.tail!==null&&li()>Vf&&(t.flags|=128,a=!0,Il(i,!1),t.lanes=4194304)}else{if(!a)if(e=Df(s),e!==null){if(t.flags|=128,a=!0,e=e.updateQueue,t.updateQueue=e,Xu(t,e),Il(i,!0),i.tail===null&&i.tailMode!=="collapsed"&&i.tailMode!=="visible"&&!s.alternate&&!Qt)return Ne(t),null}else 2*li()-i.renderingStartTime>Vf&&n!==536870912&&(t.flags|=128,a=!0,Il(i,!1),t.lanes=4194304);i.isBackwards?(s.sibling=t.child,t.child=s):(e=i.last,e!==null?e.sibling=s:t.child=s,i.last=s)}if(i.tail!==null){e=i.tail;t:{for(n=e;n!==null;){if(n.alternate!==null){n=!1;break t}n=n.sibling}n=!0}return i.rendering=e,i.tail=e.sibling,i.renderingStartTime=li(),e.sibling=null,s=bn.current,s=a?s&1|2:s&1,i.tailMode==="visible"||i.tailMode==="collapsed"||!n||Qt?cc(t,s):(n=s,Le(Tn,t),Le(bn,n),Un===null&&(Un=t)),Qt&&Aa(t,i.treeForkCount),e}return Ne(t),null;case 22:case 23:return si(t),Hg(),i=t.memoizedState!==null,e!==null?e.memoizedState!==null!==i&&(t.flags|=8192):i&&(t.flags|=8192),i?(n&536870912)!==0&&(t.flags&128)===0&&(Ne(t),t.subtreeFlags&6&&(t.flags|=8192)):Ne(t),n=t.updateQueue,n!==null&&Xu(t,n.retryQueue),n=null,e!==null&&e.memoizedState!==null&&e.memoizedState.cachePool!==null&&(n=e.memoizedState.cachePool.pool),i=null,t.memoizedState!==null&&t.memoizedState.cachePool!==null&&(i=t.memoizedState.cachePool.pool),i!==n&&(t.flags|=2048),e!==null&&Sn(rr),null;case 24:return n=null,e!==null&&(n=e.memoizedState.cache),t.memoizedState.cache!==n&&(t.flags|=2048),Ra(Qe),Ne(t),null;case 25:return null;case 30:return t.flags|=33554432,Ne(t),null}throw Error(J(156,t.tag))}function wA(e,t){switch(Pg(t),t.tag){case 1:return e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 3:return Ra(Qe),So(),e=t.flags,(e&65536)!==0&&(e&128)===0?(t.flags=e&-65537|128,t):null;case 26:case 27:case 5:return xf(t),null;case 31:if(t.memoizedState!==null){if(si(t),t.alternate===null)throw Error(J(340));cr()}return e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 13:if(si(t),e=t.memoizedState,e!==null&&e.dehydrated!==null){if(t.alternate===null)throw Error(J(340));cr()}return e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 19:return Gg(t),e=t.flags,e&65536?(t.flags=e&-65537|128,e=t.memoizedState,e!==null&&(e.rendering=null,e.tail=null),t.flags|=4,t):null;case 4:return So(),null;case 10:return Ra(t.type),null;case 22:case 23:return si(t),Hg(),e!==null&&Sn(rr),e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 24:return Ra(Qe),null;case 25:return null;default:return null}}function jS(e,t){switch(Pg(t),t.tag){case 3:Ra(Qe),So();break;case 26:case 27:case 5:xf(t);break;case 4:So();break;case 31:t.memoizedState!==null&&si(t);break;case 13:si(t);break;case 19:Gg(t);break;case 10:Ra(t.type);break;case 22:case 23:si(t),Hg(),e!==null&&Sn(rr);break;case 24:Ra(Qe)}}function Ac(e,t){try{var n=t.updateQueue,i=n!==null?n.lastEffect:null;if(i!==null){var a=i.next;n=a;do{if((n.tag&e)===e){i=void 0;var s=n.create,r=n.inst;i=s(),r.destroy=i}n=n.next}while(n!==a)}}catch(o){ye(t,t.return,o)}}function ws(e,t,n){try{var i=t.updateQueue,a=i!==null?i.lastEffect:null;if(a!==null){var s=a.next;i=s;do{if((i.tag&e)===e){var r=i.inst,o=r.destroy;if(o!==void 0){r.destroy=void 0,a=t;var l=n,c=o;try{c()}catch(f){ye(a,l,f)}}}i=i.next}while(i!==s)}}catch(f){ye(t,t.return,f)}}function $S(e){var t=e.updateQueue;if(t!==null){var n=e.stateNode;try{cS(t,n)}catch(i){ye(e,e.return,i)}}}function tb(e,t,n){n.props=mr(e.type,e.memoizedProps),n.state=e.memoizedState;try{n.componentWillUnmount()}catch(i){ye(e,t,i)}}function $i(e,t){try{var n=e.ref;if(n!==null){switch(e.tag){case 26:case 27:case 5:var i=e.stateNode;break;case 30:var a=e.stateNode,s=Na(e.memoizedProps,a);(a.ref===null||a.ref.name!==s)&&(a.ref=Xb(s)),i=a.ref;break;case 7:if(e.stateNode===null){var r=new di(e);$n(e.child,!1,mw,r,void 0,void 0),e.stateNode=r}i=e.stateNode;break;default:i=e.stateNode}typeof n=="function"?e.refCleanup=n(i):n.current=i}}catch(o){ye(e,t,o)}}function _n(e,t){var n=e.ref,i=e.refCleanup;if(n!==null)if(typeof i=="function")try{i()}catch(a){ye(e,t,a)}finally{e.refCleanup=null,e=e.alternate,e!=null&&(e.refCleanup=null)}else if(typeof n=="function")try{n(null)}catch(a){ye(e,t,a)}else n.current=null}function Of(e,t){if((e.tag===5||e.tag===27||e.tag===6)&&e.alternate===null&&t!==null)for(var n=0;n<t.length;n++)Zb(e.stateNode,t[n])}function wy(e){for(var t=e.return;t!==null&&(s0(t)&&Zb(e.stateNode,t.stateNode),!a0(t));)t=t.return}function jl(e){for(var t=e.return;t!==null&&(s0(t)&&gw(e.stateNode,t.stateNode),!a0(t));)t=t.return}function a0(e){return e.tag===5||e.tag===3||e.tag===27}function s0(e){return e&&e.tag===7&&e.stateNode!==null}function Zm(e){var t=e.type,n=e.memoizedProps,i=e.stateNode;try{t:switch(t){case"button":case"input":case"select":case"textarea":n.autoFocus&&i.focus();break t;case"img":n.src?i.src=n.src:n.srcSet&&(i.srcset=n.srcSet)}}catch(a){ye(e,e.return,a)}}function tm(e,t,n){try{var i=e.stateNode;QA(i,e.type,n,t),i[ti]=t}catch(a){ye(e,e.return,a)}}function eb(e){return e.tag===5||e.tag===3||e.tag===26||e.tag===27&&Ds(e.type)||e.tag===4}function em(e){t:for(;;){for(;e.sibling===null;){if(e.return===null||eb(e.return))return null;e=e.return}for(e.sibling.return=e.return,e=e.sibling;e.tag!==5&&e.tag!==6&&e.tag!==18;){if(e.tag===27&&Ds(e.type)||e.flags&2||e.child===null||e.tag===4)continue t;e.child.return=e,e=e.child}if(!(e.flags&2))return e.stateNode}}function Jm(e,t,n,i){var a=e.tag;if(a===5||a===6)a=e.stateNode,t?(n.nodeType===9?n.body:n.nodeName==="HTML"?n.ownerDocument.body:n).insertBefore(a,t):(t=n.nodeType===9?n.body:n.nodeName==="HTML"?n.ownerDocument.body:n,t.appendChild(a),n=n._reactRootContainer,n!=null||t.onclick!==null||(t.onclick=na)),Of(e,i),ce=!0;else if(a!==4&&(a===27&&(Of(e,i),i=null,Ds(e.type)&&(n=e.stateNode,t=null)),e=e.child,e!==null))for(Jm(e,t,n,i),e=e.sibling;e!==null;)Jm(e,t,n,i),e=e.sibling}function Pf(e,t,n,i){var a=e.tag;if(a===5||a===6)a=e.stateNode,t?n.insertBefore(a,t):n.appendChild(a),Of(e,i),ce=!0;else if(a!==4&&(a===27&&(Of(e,i),i=null,Ds(e.type)&&(n=e.stateNode)),e=e.child,e!==null))for(Pf(e,t,n,i),e=e.sibling;e!==null;)Pf(e,t,n,i),e=e.sibling}function nb(e){var t=e.stateNode,n=e.memoizedProps;try{for(var i=e.type,a=t.attributes;a.length;)t.removeAttributeNode(a[0]);Mn(t,i,n),t[yn]=e,t[ti]=n}catch(s){ye(e,e.return,s)}}var Bf=!1,ri=null;function Cy(e){(e.tag===30||(e.subtreeFlags&33554432)!==0)&&(Bf=!0)}var ta=null;function Ry(){var e=ta;return ta=null,e}var Kn=0;function Bo(e,t,n,i,a){return Kn=0,ib(e.child,t,n,i,a)}function ib(e,t,n,i,a){for(var s=!1;e!==null;){if(e.tag===5){var r=e.stateNode;if(i!==null){var o=mg(r);i.push(o),o.view&&(s=!0)}else s||mg(r).view&&(s=!0);Bf=!0,Hb(r,Kn===0?t:t+"_"+Kn,n),Kn++}else(e.tag!==22||e.memoizedState===null)&&(e.tag===30&&a||ib(e.child,t,n,i,a)&&(s=!0));e=e.sibling}return s}function oa(e,t){for(;e!==null;)e.tag===5?Gb(e.stateNode,e.memoizedProps):(e.tag!==22||e.memoizedState===null)&&(e.tag===30&&t||oa(e.child,t)),e=e.sibling}function cf(e){if((e.subtreeFlags&18874368)!==0)for(e=e.child;e!==null;){if((e.tag!==22||e.memoizedState===null)&&(cf(e),e.tag===30&&(e.flags&18874368)!==0&&e.stateNode.paired)){var t=e.memoizedProps;if(t.name==null||t.name==="auto")throw Error(J(544));var n=t.name;t=Ba(t.default,t.share),t!=="none"&&(Bo(e,n,t,null,!1)||oa(e.child,!1))}e=e.sibling}}function Km(e,t){if(e.tag===30){var n=e.stateNode,i=e.memoizedProps,a=Na(i,n),s=Ba(i.default,n.paired?i.share:i.enter);s!=="none"?Bo(e,a,s,null,!1)?(cf(e),n.paired||t||Ao(e,i.onEnter)):oa(e.child,!1):cf(e)}else if((e.subtreeFlags&33554432)!==0)for(e=e.child;e!==null;)Km(e,t),e=e.sibling;else cf(e)}function Qm(e){if(ri!==null&&ri.size!==0){var t=ri;if((e.subtreeFlags&18874368)!==0)for(e=e.child;e!==null;){if(e.tag!==22||e.memoizedState===null){if(e.tag===30&&(e.flags&18874368)!==0){var n=e.memoizedProps,i=n.name;if(i!=null&&i!=="auto"){var a=t.get(i);if(a!==void 0){var s=Ba(n.default,n.share);if(s!=="none"&&(Bo(e,i,s,null,!1)?(s=e.stateNode,a.paired=s,s.paired=a,Ao(e,n.onShare)):oa(e.child,!1)),t.delete(i),t.size===0)break}}}Qm(e)}e=e.sibling}}}function jm(e){if(e.tag===30){var t=e.memoizedProps,n=Na(t,e.stateNode),i=ri!==null?ri.get(n):void 0,a=Ba(t.default,i!==void 0?t.share:t.exit);a!=="none"&&(Bo(e,n,a,null,!1)?i!==void 0?(a=e.stateNode,i.paired=a,a.paired=i,ri.delete(n),Ao(e,t.onShare)):Ao(e,t.onExit):oa(e.child,!1)),ri!==null&&Qm(e)}else if((e.subtreeFlags&33554432)!==0)for(e=e.child;e!==null;)jm(e),e=e.sibling;else ri!==null&&Qm(e)}function ab(e){for(e=e.child;e!==null;){if(e.tag===30){var t=e.memoizedProps,n=Na(t,e.stateNode);t=Ba(t.default,t.update),e.flags&=-5,t!=="none"&&Bo(e,n,t,e.memoizedState=[],!1)}else(e.subtreeFlags&33554432)!==0&&ab(e);e=e.sibling}}function $m(e){if((e.subtreeFlags&18874368)!==0)for(e=e.child;e!==null;){if(e.tag!==22||e.memoizedState===null){if(e.tag===30&&(e.flags&18874368)!==0){var t=e.stateNode;t.paired!==null&&(t.paired=null,oa(e.child,!1))}$m(e)}e=e.sibling}}function uf(e){if(e.tag===30)e.stateNode.paired=null,oa(e.child,!1),$m(e);else if((e.subtreeFlags&33554432)!==0)for(e=e.child;e!==null;)uf(e),e=e.sibling;else $m(e)}function sb(e){for(e=e.child;e!==null;)e.tag===30?oa(e.child,!1):(e.subtreeFlags&33554432)!==0&&sb(e),e=e.sibling}function r0(e,t,n,i,a,s,r){for(var o=!1;t!==null;){if(t.tag===5){var l=t.stateNode;if(s!==null&&Kn<s.length){var c=s[Kn],f=mg(l);(c.view||f.view)&&(o=!0);var p;if(p=(e.flags&4)===0)if(f.clip)p=!0;else{p=c.rect;var u=f.rect;p=p.y!==u.y||p.x!==u.x||p.height!==u.height||p.width!==u.width}p&&(e.flags|=4),f.abs?f=!c.abs:(c=c.rect,f=f.rect,f=c.height!==f.height||c.width!==f.width),f&&(e.flags|=32)}else e.flags|=32;(e.flags&4)!==0&&Hb(l,Kn===0?n:n+"_"+Kn,a),o&&(e.flags&4)!==0||(ta===null&&(ta=[]),ta.push(l,Kn===0?i:i+"_"+Kn,t.memoizedProps)),Kn++}else(t.tag!==22||t.memoizedState===null)&&(t.tag===30&&r?e.flags|=t.flags&32:r0(e,t.child,n,i,a,s,r)&&(o=!0));t=t.sibling}return o}function rb(e,t){for(e=e.child;e!==null;){if(e.tag===30){var n=e.memoizedProps,i=e.stateNode,a=Na(n,i),s=Ba(n.default,n.update);if(t){i=i.clones;var r=i===null?null:i.map(iw)}else r=e.memoizedState,e.memoizedState=null;i=e;var o=e.child;Kn=0,a=r0(i,o,a,a,s,r,!1),(e.flags&4)!==0&&a&&(t||Ao(e,n.onUpdate))}else(e.subtreeFlags&33554432)!==0&&rb(e,t);e=e.sibling}}var cn=!1,pe=!1,Ki=!1,nm=!1,Dy=typeof WeakSet=="function"?WeakSet:Set,un=null,Qi=!1,Gl=!1,zf=!1,tg=!1;function CA(e,t,n){if(e=e.containerInfo,hg=Uo,e=Wx(e),Ug(e)){if("selectionStart"in e)var i={start:e.selectionStart,end:e.selectionEnd};else t:{i=(i=e.ownerDocument)&&i.defaultView||window;var a=i.getSelection&&i.getSelection();if(a&&a.rangeCount!==0){i=a.anchorNode;var s=a.anchorOffset,r=a.focusNode;a=a.focusOffset;try{i.nodeType,r.nodeType}catch{i=null;break t}var o=0,l=-1,c=-1,f=0,p=0,u=e,d=null;e:for(;;){for(var _;u!==i||s!==0&&u.nodeType!==3||(l=o+s),u!==r||a!==0&&u.nodeType!==3||(c=o+a),u.nodeType===3&&(o+=u.nodeValue.length),(_=u.firstChild)!==null;)d=u,u=_;for(;;){if(u===e)break e;if(d===i&&++f===s&&(l=o),d===r&&++p===a&&(c=o),(_=u.nextSibling)!==null)break;u=d,d=u.parentNode}u=_}i=l===-1||c===-1?null:{start:l,end:c}}else i=null}i=i||{start:0,end:0}}else i=null;for(dg={focusedElem:e,selectionRange:i},Uo=!1,n=(n&335544064)===n,un=t,t=n?9270:1024;un!==null;){if(e=un,n&&(i=e.deletions,i!==null))for(s=0;s<i.length;s++)n&&jm(i[s]);if(e.alternate===null&&(e.flags&2)!==0)n&&Cy(e),Wu(n);else{if(e.tag===22){if(i=e.alternate,e.memoizedState!==null){i!==null&&i.memoizedState===null&&n&&jm(i),Wu(n);continue}else if(i!==null&&i.memoizedState!==null){n&&Cy(e),Wu(n);continue}}i=e.child,(e.subtreeFlags&t)!==0&&i!==null?(i.return=e,un=i):(n&&ab(e),Wu(n))}}ri=null}function Wu(e){for(;un!==null;){var t=un,n=e,i=t.alternate,a=t.flags;switch(t.tag){case 0:case 11:case 15:break;case 1:if((a&1024)!==0&&i!==null){n=void 0,a=i.memoizedProps,i=i.memoizedState;var s=t.stateNode;try{var r=mr(t.type,a);n=s.getSnapshotBeforeUpdate(r,i),s.__reactInternalSnapshotBeforeUpdate=n}catch(o){ye(t,t.return,o)}}break;case 3:if((a&1024)!==0){if(i=t.stateNode.containerInfo,n=i.nodeType,n===9)gg(i);else if(n===1)switch(i.nodeName){case"HEAD":case"HTML":case"BODY":gg(i);break;default:i.textContent=""}}break;case 5:case 26:case 27:case 6:case 4:case 17:break;case 30:n&&i!==null&&(n=Na(i.memoizedProps,i.stateNode),a=t.memoizedProps,a=Ba(a.default,a.update),a!=="none"&&Bo(i,n,a,i.memoizedState=[],!0));break;default:if((a&1024)!==0)throw Error(J(163))}if(i=t.sibling,i!==null){i.return=t.return,un=i;break}un=t.return}}function ob(e,t,n){var i=n.flags;switch(n.tag){case 0:case 11:case 15:ji(e,n),i&4&&Ac(5,n);break;case 1:if(ji(e,n),i&4)if(e=n.stateNode,t===null)try{e.componentDidMount()}catch(r){ye(n,n.return,r)}else{var a=mr(n.type,t.memoizedProps);t=t.memoizedState;try{e.componentDidUpdate(a,t,e.__reactInternalSnapshotBeforeUpdate)}catch(r){ye(n,n.return,r)}}i&64&&$S(n),i&512&&$i(n,n.return);break;case 3:if(ji(e,n),i&64&&(e=n.updateQueue,e!==null)){if(t=null,n.child!==null)switch(n.child.tag){case 27:case 5:t=n.child.stateNode;break;case 1:t=n.child.stateNode}try{cS(e,t)}catch(r){ye(n,n.return,r)}}break;case 27:t===null&&i&4&&nb(n);case 26:case 5:ji(e,n),t===null&&i&4&&Zm(n),i&512&&$i(n,n.return);break;case 12:ji(e,n);break;case 31:ji(e,n),i&4&&fb(e,n);break;case 13:ji(e,n),i&4&&hb(e,n),i&64&&(e=n.memoizedState,e!==null&&(e=e.dehydrated,e!==null&&(n=VA.bind(null,n),yw(e,n))));break;case 22:if(i=n.memoizedState!==null||cn,!i){var s=t!==null&&t.memoizedState!==null||pe;t=cn,a=pe,cn=i,(pe=s)&&!a?(i=2,(n.subtreeFlags&8772)!==0&&(i|=1),Oi(e,n,i)):ji(e,n),cn=t,pe=a}break;case 30:ji(e,n),i&512&&$i(n,n.return);break;case 7:i&512&&$i(n,n.return);default:ji(e,n)}}function eg(e,t){for(e=e.child;e!==null;)lb(e,t),e=e.sibling}function lb(e,t){switch(e.tag){case 5:case 26:try{var n=e.stateNode;if(t){var i=n.style;typeof i.setProperty=="function"?i.setProperty("display","none","important"):i.display="none"}else{var a=e.stateNode,s=e.memoizedProps.style,r=s!=null&&s.hasOwnProperty("display")?s.display:null;a.style.display=r==null||typeof r=="boolean"?"":(""+r).trim()}}catch(l){ye(e,e.return,l)}ng(e,t);break;case 6:try{e.stateNode.nodeValue=t?"":e.memoizedProps,ce=!0}catch(l){ye(e,e.return,l)}break;case 18:try{var o=e.stateNode;t?qy(o,!0):qy(e.stateNode,!1)}catch(l){ye(e,e.return,l)}break;case 22:case 23:e.memoizedState===null&&eg(e,t);break;default:eg(e,t)}}function ng(e,t){if(e.subtreeFlags&67108864)for(e=e.child;e!==null;){t:{var n=e,i=t;switch(n.tag){case 4:lb(n,i);break t;case 22:n.memoizedState===null&&ng(n,i);break t;default:ng(n,i)}}e=e.sibling}}function cb(e){var t=e.alternate;t!==null&&(e.alternate=null,cb(t)),e.child=null,e.deletions=null,e.sibling=null,e.tag===5&&(t=e.stateNode,t!==null&&Zf(t)),e.stateNode=null,e.return=null,e.dependencies=null,e.memoizedProps=null,e.memoizedState=null,e.pendingProps=null,e.stateNode=null,e.updateQueue=null}var Ve=null,Zn=!1;function Ii(e,t,n){for(n=n.child;n!==null;)ub(e,t,n),n=n.sibling}function ub(e,t,n){if(ci&&typeof ci.onCommitFiberUnmount=="function")try{ci.onCommitFiberUnmount(yc,n)}catch{}switch(n.tag){case 26:pe||_n(n,t),Ii(e,t,n),n.memoizedState?n.memoizedState.count--:n.stateNode&&!pe&&(n=n.stateNode,n.parentNode.removeChild(n));break;case 27:pe||_n(n,t),jl(n);var i=Ve,a=Zn;Ds(n.type)&&(Ve=n.stateNode,Zn=!1),Ii(e,t,n),Qb(n.stateNode,n.type,n.memoizedProps),Ve=i,Zn=a;break;case 5:pe||_n(n,t),jl(n);case 6:if(n.tag===6&&jl(n),i=Ve,a=Zn,Ve=null,Ii(e,t,n),Ve=i,Zn=a,Ve!==null)if(Zn)try{(Ve.nodeType===9?Ve.body:Ve.nodeName==="HTML"?Ve.ownerDocument.body:Ve).removeChild(n.stateNode),ce=!0}catch(s){ye(n,t,s)}else try{Ve.removeChild(n.stateNode),ce=!0}catch(s){ye(n,t,s)}break;case 18:Ve!==null&&(Zn?(e=Ve,Wy(e.nodeType===9?e.body:e.nodeName==="HTML"?e.ownerDocument.body:e,n.stateNode),Lo(e)):Wy(Ve,n.stateNode));break;case 4:i=Ve,a=Zn,Ve=n.stateNode.containerInfo,Zn=!0,Ii(e,t,n),Ve=i,Zn=a;break;case 0:case 11:case 14:case 15:ws(2,n,t),pe||ws(4,n,t),Ii(e,t,n);break;case 1:pe||(_n(n,t),i=n.stateNode,typeof i.componentWillUnmount=="function"&&tb(n,t,i)),Ii(e,t,n);break;case 21:Ii(e,t,n);break;case 22:pe=(i=pe)||n.memoizedState!==null,Ii(e,t,n),pe=i;break;case 30:_n(n,t),Ii(e,t,n);break;case 7:pe||_n(n,t),Ii(e,t,n);break;default:Ii(e,t,n)}}function fb(e,t){if(t.memoizedState===null&&(e=t.alternate,e!==null&&(e=e.memoizedState,e!==null))){e=e.dehydrated;try{Lo(e)}catch(n){ye(t,t.return,n)}}}function hb(e,t){if(t.memoizedState===null&&(e=t.alternate,e!==null&&(e=e.memoizedState,e!==null&&(e=e.dehydrated,e!==null))))try{Lo(e)}catch(n){ye(t,t.return,n)}}function RA(e){switch(e.tag){case 31:case 13:case 19:var t=e.stateNode;return t===null&&(t=e.stateNode=new Dy),t;case 22:return e=e.stateNode,t=e._retryCache,t===null&&(t=e._retryCache=new Dy),t;default:throw Error(J(435,e.tag))}}function qu(e,t){var n=RA(e);t.forEach(function(i){if(!n.has(i)){n.add(i);var a=HA.bind(null,e,i);i.then(a,a)}})}function Hn(e,t,n){var i=t.deletions;if(i!==null)for(var a=0;a<i.length;a++){var s=i[a],r=e,o=t,l=o;t:for(;l!==null;){switch(l.tag){case 27:if(Ds(l.type)){Ve=l.stateNode,Zn=!1;break t}break;case 5:Ve=l.stateNode,Zn=!1;break t;case 3:case 4:Ve=l.stateNode.containerInfo,Zn=!0;break t}l=l.return}if(Ve===null)throw Error(J(160));ub(r,o,s),Ve=null,Zn=!1,r=s.alternate,r!==null&&(r.return=null),s.return=null}if(t.subtreeFlags&13886)for(t=t.child;t!==null;)db(t,e,n),t=t.sibling}var Pi=null;function db(e,t,n){var i=e.alternate,a=e.flags;switch(e.tag){case 0:case 11:case 14:case 15:if(a&4&&(i=e.updateQueue,i=i!==null?i.events:null,i!==null))for(var s=0;s<i.length;s++){var r=i[s];r.ref.impl=r.nextImpl}Hn(t,e,n),Gn(e),a&4&&(ws(3,e,e.return),Ac(3,e),ws(5,e,e.return));break;case 1:Hn(t,e,n),Gn(e),a&512&&(pe||i===null||_n(i,i.return)),a&64&&cn&&(e=e.updateQueue,e!==null&&(t=e.callbacks,t!==null&&(n=e.shared.hiddenCallbacks,e.shared.hiddenCallbacks=n===null?t:n.concat(t))));break;case 26:if(s=Pi,Hn(t,e,n),Gn(e),a&512&&(pe||i===null||_n(i,i.return)),a&4)if(a=i!==null?i.memoizedState:null,n=e.memoizedState,i===null)if(n===null)if(e.stateNode===null)if(cn)e.stateNode=Vb(e.type,e.memoizedProps,t.containerInfo,e);else{t:{t=e.type,n=e.memoizedProps,a=s.ownerDocument||s;e:switch(t){case"title":i=a.getElementsByTagName("title")[0],(!i||i[bc]||i[yn]||i.namespaceURI==="http://www.w3.org/2000/svg"||i.hasAttribute("itemprop"))&&(i=a.createElement(t),a.head.insertBefore(i,a.querySelector("head > title"))),Mn(i,t,n),i[yn]=e,fn(i),t=i;break t;case"link":if(s=tx("link","href",a).get(t+(n.href||""))){for(r=0;r<s.length;r++)if(i=s[r],i.getAttribute("href")===(n.href==null||n.href===""?null:n.href)&&i.getAttribute("rel")===(n.rel==null?null:n.rel)&&i.getAttribute("title")===(n.title==null?null:n.title)&&i.getAttribute("crossorigin")===(n.crossOrigin==null?null:n.crossOrigin)){s.splice(r,1);break e}}i=a.createElement(t),Mn(i,t,n),a.head.appendChild(i);break;case"meta":if(s=tx("meta","content",a).get(t+(n.content||""))){for(r=0;r<s.length;r++)if(i=s[r],i.getAttribute("content")===(n.content==null?null:""+n.content)&&i.getAttribute("name")===(n.name==null?null:n.name)&&i.getAttribute("property")===(n.property==null?null:n.property)&&i.getAttribute("http-equiv")===(n.httpEquiv==null?null:n.httpEquiv)&&i.getAttribute("charset")===(n.charSet==null?null:n.charSet)){s.splice(r,1);break e}}i=a.createElement(t),Mn(i,t,n),a.head.appendChild(i);break;default:throw Error(J(468,t))}i[yn]=e,fn(i),t=i}e.stateNode=t}else cn||yg(s,e.type,e.stateNode);else e.stateNode=$y(s,n,e.memoizedProps);else a!==n?(a===null?(t=i.stateNode,t===null||pe||t.parentNode.removeChild(t)):a.count--,n===null?cn||yg(s,e.type,e.stateNode):$y(s,n,e.memoizedProps)):n===null&&e.stateNode!==null&&tm(e,e.memoizedProps,i.memoizedProps);break;case 27:Hn(t,e,n),Gn(e),a&512&&(pe||i===null||_n(i,i.return)),i!==null&&a&4&&tm(e,e.memoizedProps,i.memoizedProps);break;case 5:if(s=Ki,Ki=!1,Hn(t,e,n),Ki=s,Gn(e),a&512&&(pe||i===null||_n(i,i.return)),e.flags&32){t=e.stateNode;try{Mo(t,""),ce=!0}catch(f){ye(e,e.return,f)}}a&4&&e.stateNode!=null&&(t=e.memoizedProps,tm(e,t,i!==null?i.memoizedProps:t)),a&1024&&(nm=!0);break;case 6:if(Hn(t,e,n),Gn(e),a&4){if(e.stateNode===null)throw Error(J(162));t=e.memoizedProps,n=e.stateNode;try{n.nodeValue=t,ce=!0}catch(f){ye(e,e.return,f)}}break;case 3:if(ce=!1,pf=null,s=Pi,Pi=pc(t.containerInfo),Hn(t,e,n),Pi=s,Gn(e),a&4&&i!==null&&i.memoizedState.isDehydrated)try{Lo(t.containerInfo)}catch(f){ye(e,e.return,f)}nm&&(nm=!1,pb(e)),ce=!1;break;case 4:a=Ki,Ki=cn,i=P_(),s=Pi,Pi=pc(e.stateNode.containerInfo),Hn(t,e,n),Gn(e),Pi=s,ce&&Gl&&(zf=!0),ce=i,Ki=a;break;case 12:Hn(t,e,n),Gn(e);break;case 31:Hn(t,e,n),Gn(e),a&4&&(t=e.updateQueue,t!==null&&(e.updateQueue=null,qu(e,t)));break;case 13:Hn(t,e,n),Gn(e),e.child.flags&8192&&e.memoizedState!==null!=(i!==null&&i.memoizedState!==null)&&(oh=li()),a&4&&(t=e.updateQueue,t!==null&&(e.updateQueue=null,qu(e,t)));break;case 22:s=e.memoizedState!==null,r=i!==null&&i.memoizedState!==null;var o=cn,l=pe,c=Ki;cn=o||s,Ki=c||s,pe=l||r,Hn(t,e,n),pe=l,Ki=c,cn=o,Gn(e),a&8192&&(t=e.stateNode,t._visibility=s?t._visibility&-2:t._visibility|1,!s||i===null||r||cn||pe||(t=r||pe,n=cn,i=pe,cn=s||cn,pe=t,as(e,2),cn=n,pe=i),!s&&Ki||eg(e,s)),a&4&&(t=e.updateQueue,t!==null&&(n=t.retryQueue,n!==null&&(t.retryQueue=null,qu(e,n))));break;case 19:Hn(t,e,n),Gn(e),a&4&&(t=e.updateQueue,t!==null&&(e.updateQueue=null,qu(e,t)));break;case 30:a&512&&(pe||i===null||_n(i,i.return)),a=P_(),s=Gl,r=(n&335544064)===n,o=e.memoizedProps,Gl=r&&Ba(o.default,o.update)!=="none",Hn(t,e,n),Gn(e),r&&i!==null&&ce&&(e.flags|=4),Gl=s,ce=a;break;case 21:break;case 7:a&512&&(pe||i===null||_n(i,i.return)),i&&i.stateNode!==null&&(i.stateNode._fragmentFiber=e);default:Hn(t,e,n),Gn(e)}}function Gn(e){var t=e.flags;if(t&2){try{for(var n,i=e.return;i!==null;){if(eb(i)){n=i;break}i=i.return}i=null;for(var a=e.return;a!==null;){if(s0(a)){var s=a.stateNode;i===null?i=[s]:i.push(s)}if(a0(a))break;a=a.return}var r=i;if(n==null)throw Error(J(160));switch(n.tag){case 27:var o=n.stateNode,l=em(e);Pf(e,l,o,r);break;case 5:var c=n.stateNode;n.flags&32&&(Mo(c,""),n.flags&=-33);var f=em(e);Pf(e,f,c,r);break;case 3:case 4:var p=n.stateNode.containerInfo,u=em(e);Jm(e,u,p,r);break;default:throw Error(J(161))}}catch(d){ye(e,e.return,d)}e.flags&=-3}t&4096&&(e.flags&=-4097)}function pb(e){if(e.subtreeFlags&1024)for(e=e.child;e!==null;){var t=e;pb(t),t.tag===5&&t.flags&1024&&(t=t.stateNode,Uo=!0,t.reset(),Uo=!1),e=e.sibling}}function Zr(e,t){if(t.subtreeFlags&9270)for(t=t.child;t!==null;)mb(t,e),t=t.sibling;else rb(t,!1)}function mb(e,t){var n=e.alternate;if(n===null)Km(e,!1);else switch(e.tag){case 3:if(tg=Qi=!1,Ry(),Zr(t,e),!Qi&&!zf){if(e=ta,e!==null)for(var i=0;i<e.length;i+=3){n=e[i];var a=e[i+1];Gb(n,e[i+2]),n=n.ownerDocument.documentElement,n!==null&&n.animate({opacity:[0,0],pointerEvents:["none","none"]},{duration:0,fill:"forwards",pseudoElement:"::view-transition-group("+a+")"})}e=t.containerInfo,e=e.nodeType===9?e.documentElement:e.ownerDocument.documentElement,e!==null&&e.style.viewTransitionName===""&&(e.style.viewTransitionName="none",e.animate({opacity:[0,0],pointerEvents:["none","none"]},{duration:0,fill:"forwards",pseudoElement:"::view-transition-group(root)"}),e.animate({width:[0,0],height:[0,0]},{duration:0,fill:"forwards",pseudoElement:"::view-transition"})),tg=!0}ta=null;break;case 5:Zr(t,e);break;case 4:i=Qi,Qi=!1,Zr(t,e),Qi&&(zf=!0),Qi=i;break;case 22:e.memoizedState===null&&(n.memoizedState!==null?Km(e,!1):Zr(t,e));break;case 30:i=Qi,a=Ry(),Qi=!1,Zr(t,e),Qi&&(e.flags|=4);var s=e.memoizedProps,r=e.stateNode;t=Na(s,r),r=Na(n.memoizedProps,r);var o=Ba(s.default,s.update);o==="none"?t=!1:(s=n.memoizedState,n.memoizedState=null,n=e.child,Kn=0,t=r0(e,n,t,r,o,s,!0),Kn!==(s===null?0:s.length)&&(e.flags|=32)),(e.flags&4)!==0&&t?(Ao(e,e.memoizedProps.onUpdate),ta=a):a!==null&&(a.push.apply(a,ta),ta=a),Qi=(e.flags&32)!==0?!0:i;break;default:Zr(t,e)}}function ji(e,t){if(t.subtreeFlags&8772)for(t=t.child;t!==null;)ob(e,t.alternate,t),t=t.sibling}function as(e,t){for(e=e.child;e!==null;){var n=e,i=t;switch(n.tag){case 0:case 11:case 14:case 15:ws(4,n,n.return),as(n,i);break;case 1:_n(n,n.return);var a=n.stateNode;typeof a.componentWillUnmount=="function"&&tb(n,n.return,a),as(n,i);break;case 27:(i&2)!==0&&Qb(n.stateNode,n.type,n.memoizedProps);case 5:_n(n,n.return),n.tag!==5&&n.tag!==27||jl(n),as(n,i);break;case 6:jl(n);break;case 26:_n(n,n.return),a=n.stateNode,n.memoizedState!==null||a===null||pe||a.parentNode.removeChild(a),as(n,i);break;case 22:n.memoizedState===null&&as(n,i);break;case 30:_n(n,n.return),as(n,i);break;case 7:_n(n,n.return);default:as(n,i)}e=e.sibling}}function Oi(e,t,n){for(n=(t.subtreeFlags&8772)!==0?n:n&-2,t=t.child;t!==null;){var i=t.alternate,a=e,s=t,r=s.flags,o=(n&1)!==0;switch(s.tag){case 0:case 11:case 15:Oi(a,s,n),Ac(4,s);break;case 1:if(Oi(a,s,n),i=s,a=i.stateNode,typeof a.componentDidMount=="function")try{a.componentDidMount()}catch(f){ye(i,i.return,f)}if(i=s,a=i.updateQueue,a!==null){var l=i.stateNode;try{var c=a.shared.hiddenCallbacks;if(c!==null)for(a.shared.hiddenCallbacks=null,a=0;a<c.length;a++)lS(c[a],l)}catch(f){ye(i,i.return,f)}}o&&r&64&&$S(s),$i(s,s.return);break;case 27:(n&2)!==0&&nb(s);case 5:s.tag!==5&&s.tag!==27||wy(s),Oi(a,s,n),o&&i===null&&r&4&&Zm(s),$i(s,s.return);break;case 6:wy(s);break;case 26:l=s.stateNode,s.memoizedState!==null||l===null||cn||yg(pc(l.ownerDocument),s.type,l),Oi(a,s,n),o&&i===null&&r&4&&Zm(s),$i(s,s.return);break;case 12:Oi(a,s,n);break;case 31:Oi(a,s,n),o&&r&4&&fb(a,s);break;case 13:Oi(a,s,n),o&&r&4&&hb(a,s);break;case 22:s.memoizedState===null&&Oi(a,s,n),$i(s,s.return);break;case 30:Oi(a,s,n),$i(s,s.return);break;case 7:$i(s,s.return);default:Oi(a,s,n)}t=t.sibling}}function o0(e,t){var n=null;e!==null&&e.memoizedState!==null&&e.memoizedState.cachePool!==null&&(n=e.memoizedState.cachePool.pool),e=null,t.memoizedState!==null&&t.memoizedState.cachePool!==null&&(e=t.memoizedState.cachePool.pool),e!==n&&(e!=null&&e.refCount++,n!=null&&Tc(n))}function l0(e,t){e=null,t.alternate!==null&&(e=t.alternate.memoizedState.cache),t=t.memoizedState.cache,t!==e&&(t.refCount++,e!=null&&Tc(e))}function Si(e,t,n,i){var a=(n&335544064)===n;if(t.subtreeFlags&(a?10262:10256))for(t=t.child;t!==null;)gb(e,t,n,i),t=t.sibling;else a&&sb(t)}function gb(e,t,n,i){var a=(n&335544064)===n;a&&t.alternate===null&&t.return!==null&&t.return.alternate!==null&&uf(t);var s=t.flags;switch(t.tag){case 0:case 11:case 15:Si(e,t,n,i),s&2048&&Ac(9,t);break;case 1:Si(e,t,n,i);break;case 3:Si(e,t,n,i),a&&tg&&(e=e.containerInfo,e=e.nodeType===9?e.body:e.nodeName==="HTML"?e.ownerDocument.body:e,e.style.viewTransitionName==="root"&&(e.style.viewTransitionName=""),e=e.ownerDocument.documentElement,e!==null&&e.style.viewTransitionName==="none"&&(e.style.viewTransitionName="")),s&2048&&(s=null,t.alternate!==null&&(s=t.alternate.memoizedState.cache),t=t.memoizedState.cache,t!==s&&(t.refCount++,s!=null&&Tc(s)));break;case 12:if(s&2048){Si(e,t,n,i),s=t.stateNode;try{var r=t.memoizedProps,o=r.id,l=r.onPostCommit;typeof l=="function"&&l(o,t.alternate===null?"mount":"update",s.passiveEffectDuration,-0)}catch(c){ye(t,t.return,c)}}else Si(e,t,n,i);break;case 31:Si(e,t,n,i);break;case 13:Si(e,t,n,i);break;case 23:break;case 22:r=t.stateNode,o=t.alternate,t.memoizedState!==null?(a&&o!==null&&o.memoizedState===null&&uf(o),r._visibility&2?Si(e,t,n,i):$l(e,t)):(a&&o!==null&&o.memoizedState!==null&&uf(t),r._visibility&2?Si(e,t,n,i):(r._visibility|=2,Kr(e,t,n,i,(t.subtreeFlags&10256)!==0||!1))),s&2048&&o0(o,t);break;case 24:Si(e,t,n,i),s&2048&&l0(t.alternate,t);break;case 30:a&&(s=t.alternate,s!==null&&(oa(s.child,!0),oa(t.child,!0))),Si(e,t,n,i);break;default:Si(e,t,n,i)}}function Kr(e,t,n,i,a){for(a=a&&((t.subtreeFlags&10256)!==0||!1),t=t.child;t!==null;){var s=e,r=t,o=n,l=i,c=r.flags;switch(r.tag){case 0:case 11:case 15:Kr(s,r,o,l,a),Ac(8,r);break;case 23:break;case 22:var f=r.stateNode;r.memoizedState!==null?f._visibility&2?Kr(s,r,o,l,a):$l(s,r):(f._visibility|=2,Kr(s,r,o,l,a)),a&&c&2048&&o0(r.alternate,r);break;case 24:Kr(s,r,o,l,a),a&&c&2048&&l0(r.alternate,r);break;default:Kr(s,r,o,l,a)}t=t.sibling}}function $l(e,t){if(t.subtreeFlags&10256)for(t=t.child;t!==null;){var n=e,i=t,a=i.flags;switch(i.tag){case 22:$l(n,i),a&2048&&o0(i.alternate,i);break;case 24:$l(n,i),a&2048&&l0(i.alternate,i);break;default:$l(n,i)}t=t.sibling}}var er=8192;function js(e,t,n){if(e.subtreeFlags&er)for(e=e.child;e!==null;)vb(e,t,n),e=e.sibling}function vb(e,t,n){switch(e.tag){case 26:js(e,t,n),e.flags&er&&(e.memoizedState!==null?Lw(n,Pi,e.memoizedState,e.memoizedProps):(e=e.stateNode,(t&335544128)===t&&nx(n,e)));break;case 5:js(e,t,n),e.flags&er&&(e=e.stateNode,(t&335544128)===t&&nx(n,e));break;case 3:case 4:var i=Pi;Pi=pc(e.stateNode.containerInfo),js(e,t,n),Pi=i;break;case 22:e.memoizedState===null&&(i=e.alternate,i!==null&&i.memoizedState!==null?(i=er,er=16777216,js(e,t,n),er=i):js(e,t,n));break;case 30:if((e.flags&er)!==0&&(i=e.memoizedProps.name,i!=null&&i!=="auto")){var a=e.stateNode;a.paired=null,ri===null&&(ri=new Map),ri.set(i,a)}js(e,t,n);break;default:js(e,t,n)}}function _b(e){var t=e.alternate;if(t!==null&&(e=t.child,e!==null)){t.child=null;do t=e.sibling,e.sibling=null,e=t;while(e!==null)}}function Ol(e){var t=e.deletions;if((e.flags&16)!==0){if(t!==null)for(var n=0;n<t.length;n++){var i=t[n];un=i,xb(i,e)}_b(e)}if(e.subtreeFlags&10256)for(e=e.child;e!==null;)yb(e),e=e.sibling}function yb(e){switch(e.tag){case 0:case 11:case 15:Ol(e),e.flags&2048&&ws(9,e,e.return);break;case 3:Ol(e);break;case 12:Ol(e);break;case 22:var t=e.stateNode;e.memoizedState!==null&&t._visibility&2&&(e.return===null||e.return.tag!==13)?(t._visibility&=-3,ff(e)):Ol(e);break;default:Ol(e)}}function ff(e){var t=e.deletions;if((e.flags&16)!==0){if(t!==null)for(var n=0;n<t.length;n++){var i=t[n];un=i,xb(i,e)}_b(e)}for(e=e.child;e!==null;){switch(t=e,t.tag){case 0:case 11:case 15:ws(8,t,t.return),ff(t);break;case 22:n=t.stateNode,n._visibility&2&&(n._visibility&=-3,ff(t));break;default:ff(t)}e=e.sibling}}function xb(e,t){for(;un!==null;){var n=un;switch(n.tag){case 0:case 11:case 15:ws(8,n,t);break;case 23:case 22:if(n.memoizedState!==null&&n.memoizedState.cachePool!==null){var i=n.memoizedState.cachePool.pool;i!=null&&i.refCount++}break;case 24:Tc(n.memoizedState.cache)}if(i=n.child,i!==null)i.return=n,un=i;else t:for(n=e;un!==null;){i=un;var a=i.sibling,s=i.return;if(cb(i),i===n){un=null;break t}if(a!==null){a.return=s,un=a;break t}un=s}}}var DA={getCacheForType:function(e){var t=xn(Qe),n=t.data.get(e);return n===void 0&&(n=e(),t.data.set(e,n)),n},cacheSignal:function(){return xn(Qe).controller.signal}},NA=typeof WeakMap=="function"?WeakMap:Map,fe=0,Ae=null,$t=null,ee=0,ve=0,ii=null,fs=!1,zo=!1,c0=!1,Oa=0,qe=0,Cs=0,lr=0,Ff=0,oi=0,Eo=0,tc=null,Jn=null,ig=!1,oh=0,Sb=0,Vf=1/0,Hf=null,xs=null,He=0,zi=null,gr=null,ra=0,ag=0,sg=null,bb=null,_o=null,yo=null,xo=null,ec=0,hf=null;function fi(){return(fe&2)!==0&&ee!==0?ee&-ee:Bt.T!==null?f0():Ax()}function Mb(){if(oi===0)if((ee&536870912)===0||Qt){var e=Lu;Lu<<=1,(Lu&3932160)===0&&(Lu=262144),oi=e}else oi=536870912;return e=Tn.current,e!==null&&(e.flags|=32),oi}function Ao(e,t){if(t!=null){var n=e.stateNode,i=n.ref;i===null&&(i=n.ref=Xb(Na(e.memoizedProps,n))),yo===null&&(yo=[]),yo.push(t.bind(null,i))}}function jn(e,t,n){(e===Ae&&(ve===2||ve===9)||e.cancelPendingCommit!==null)&&(wo(e,0),hs(e,ee,oi,!1)),Sc(e,n),((fe&2)===0||e!==Ae)&&(e===Ae&&((fe&2)===0&&(lr|=n),qe===4&&hs(e,ee,oi,!1)),ca(e))}function Tb(e,t,n){if((fe&6)!==0)throw Error(J(327));var i=!n&&(t&127)===0&&(t&e.expiredLanes)===0||xc(e,t),a=i?IA(e,t):im(e,t,!0),s=i;do{if(a===0){zo&&!i&&hs(e,t,0,!1);break}else{if(n=e.current.alternate,s&&!UA(n)){a=im(e,t,!1),s=!1;continue}if(a===2){if(s=t,e.errorRecoveryDisabledLanes&s)var r=0;else r=e.pendingLanes&-536870913,r=r!==0?r:r&536870912?536870912:0;if(r!==0){t=r;t:{var o=e;a=tc;var l=o.current.memoizedState.isDehydrated;if(l&&(wo(o,r).flags|=256),r=im(o,r,!1),r!==2&&r!==6){if(c0&&!l){o.errorRecoveryDisabledLanes|=s,lr|=s,a=4;break t}s=Jn,Jn=a,s!==null&&(Jn===null?Jn=s:Jn.push.apply(Jn,s))}a=r}if(s=!1,a!==2)continue}}if(a===1){wo(e,0),hs(e,t,0,!0);break}t:{switch(i=e,s=a,s){case 0:case 1:throw Error(J(345));case 4:if((t&4194048)!==t&&(t&62914560)!==t)break;case 6:hs(i,t,oi,!fs);break t;case 2:Jn=null;break;case 3:case 5:break;default:throw Error(J(329))}if((t&62914560)===t&&(a=oh+300-li(),10<a)){if(hs(i,t,oi,!fs),Yf(i,0,!0)!==0)break t;ra=t,i.timeoutHandle=d0(Ny.bind(null,i,n,Jn,Hf,ig,t,oi,lr,Eo,fs,s,"Throttled",-0,0),a);break t}Ny(i,n,Jn,Hf,ig,t,oi,lr,Eo,fs,s,null,-0,0)}}break}while(!0);ca(e)}function Ny(e,t,n,i,a,s,r,o,l,c,f,p,u,d){e.timeoutHandle=-1;var _=t.subtreeFlags,b=(s&335544064)===s;if(p=null,(b||_&8192||(_&16785408)===16785408)&&(p={stylesheets:null,count:0,imgCount:0,imgBytes:0,suspenseyImages:[],waitingForImages:!0,waitingForViewTransition:!1,unsuspend:na},ri=null,vb(t,s,p),b&&(_=p,b=e.containerInfo,b=(b.nodeType===9?b:b.ownerDocument).__reactViewTransition,b!=null&&(_.count++,_.waitingForViewTransition=!0,_=mc.bind(_),b.finished.then(_,_))),_=(s&62914560)===s?oh-li():(s&4194048)===s?Sb-li():0,_=Iw(p,_),_!==null)){ra=s,e.cancelPendingCommit=_(Ly.bind(null,e,t,s,n,i,a,r,o,l,c,f,p,null,u,d)),hs(e,s,r,!c);return}Ly(e,t,s,n,i,a,r,o,l,c,f,p)}function UA(e){for(var t=e;;){var n=t.tag;if((n===0||n===11||n===15)&&t.flags&16384&&(n=t.updateQueue,n!==null&&(n=n.stores,n!==null)))for(var i=0;i<n.length;i++){var a=n[i],s=a.getSnapshot;a=a.value;try{if(!hi(s(),a))return!1}catch{return!1}}if(n=t.child,t.subtreeFlags&16384&&n!==null)n.return=t,t=n;else{if(t===e)break;for(;t.sibling===null;){if(t.return===null||t.return===e)return!0;t=t.return}t.sibling.return=t.return,t=t.sibling}}return!0}function hs(e,t,n,i){t=Sx(e,t),t&=~Ff,t&=~lr,e.suspendedLanes|=t,e.pingedLanes&=~t,i&&(e.warmLanes|=t),i=e.expirationTimes;for(var a=t;0<a;){var s=31-ui(a),r=1<<s;i[s]=-1,a&=~r}n!==0&&Mx(e,n,t)}function lh(){return(fe&6)===0?(wc(0,!1),!1):!0}function u0(){if($t!==null){if(ve===0)var e=$t.return;else e=$t,wa=Sr=null,Yg(e),mo=null,lc=0,e=$t;for(;e!==null;)jS(e.alternate,e),e=e.return;$t=null}}function wo(e,t){var n=e.timeoutHandle;return n!==-1&&(e.timeoutHandle=-1,tw(n)),n=e.cancelPendingCommit,n!==null&&(e.cancelPendingCommit=null,n()),ra=0,u0(),Ae=e,$t=n=Ca(e.current,null),ee=t,ve=0,ii=null,fs=!1,zo=xc(e,t),c0=!1,Eo=oi=Ff=lr=Cs=qe=0,Jn=tc=null,ig=!1,Oa=Sx(e,t),jf(),n}function Eb(e,t){Yt=null,Bt.H=Lf,t===Po||t===eh?(t=ay(),ve=3):t===Fg?(t=ay(),ve=4):ve=t===e0?8:t!==null&&typeof t=="object"&&typeof t.then=="function"?6:1,ii=t,$t===null&&(qe=1,If(e,Ai(t,e.current)))}function Ab(){var e=Tn.current;return e===null?!0:(ee&4194048)===ee?Un===null:(ee&62914560)===ee||(ee&536870912)!==0?e===Un:!1}function wb(){var e=Bt.H;return Bt.H=Lf,e===null?Lf:e}function Cb(){var e=Bt.A;return Bt.A=DA,e}function Gf(){qe=4,fs||(ee&4194048)!==ee&&Tn.current!==null||(zo=!0),(Cs&134217727)===0&&(lr&134217727)===0||Ae===null||hs(Ae,ee,oi,!1)}function im(e,t,n){var i=fe;fe|=2;var a=wb(),s=Cb();(Ae!==e||ee!==t)&&(Hf=null,wo(e,t)),t=!1;var r=qe;t:do try{if(ve!==0&&$t!==null){var o=$t,l=ii;switch(ve){case 8:u0(),r=6;break t;case 3:case 2:case 9:case 6:Tn.current===null&&(t=!0);var c=ve;if(ve=0,ii=null,co(e,o,l,c),n&&zo){r=0;break t}break;default:c=ve,ve=0,ii=null,co(e,o,l,c)}}LA(),r=qe;break}catch(f){Eb(e,f)}while(!0);return t&&e.shellSuspendCounter++,wa=Sr=null,fe=i,Bt.H=a,Bt.A=s,$t===null&&(Ae=null,ee=0,jf()),r}function LA(){for(;$t!==null;)Rb($t)}function IA(e,t){var n=fe;fe|=2;var i=wb(),a=Cb();Ae!==e||ee!==t?(Hf=null,Vf=li()+500,wo(e,t)):zo=xc(e,t);t:do try{if(ve!==0&&$t!==null){t=$t;var s=ii;e:switch(ve){case 1:ve=0,ii=null,co(e,t,s,1);break;case 2:case 9:if(iy(s)){ve=0,ii=null,Uy(t);break}t=function(){ve!==2&&ve!==9||Ae!==e||(ve=7),ca(e)},s.then(t,t);break t;case 3:ve=7;break t;case 4:ve=5;break t;case 7:iy(s)?(ve=0,ii=null,Uy(t)):(ve=0,ii=null,co(e,t,s,7));break;case 5:var r=null;switch($t.tag){case 26:r=$t.memoizedState;case 5:case 27:var o=$t;if(r?tM(r):o.stateNode.complete){ve=0,ii=null;var l=o.sibling;if(l!==null)$t=l;else{var c=o.return;c!==null?($t=c,ch(c)):$t=null}break e}}ve=0,ii=null,co(e,t,s,5);break;case 6:ve=0,ii=null,co(e,t,s,6);break;case 8:u0(),qe=6;break t;default:throw Error(J(462))}}OA();break}catch(f){Eb(e,f)}while(!0);return wa=Sr=null,Bt.H=i,Bt.A=a,fe=n,$t!==null?0:(Ae=null,ee=0,jf(),qe)}function OA(){for(;$t!==null&&!jT();)Rb($t)}function Rb(e){var t=QS(e.alternate,e,Oa);e.memoizedProps=e.pendingProps,t===null?ch(e):$t=t}function Uy(e){var t=e,n=t.alternate;switch(t.tag){case 15:case 0:t=xy(n,t,t.pendingProps,t.type,void 0,ee);break;case 11:t=xy(n,t,t.pendingProps,t.type.render,t.ref,ee);break;case 5:Yg(t);var i=t;i===hn&&(Qt?(Af(i),i.tag===5&&i.stateNode!=null&&(Ue=i.stateNode)):(Af(i),Qt=!0));default:jS(n,t),t=$t=$x(t,Oa),t=QS(n,t,Oa)}e.memoizedProps=e.pendingProps,t===null?ch(e):$t=t}function co(e,t,n,i){wa=Sr=null,Yg(t),mo=null,lc=0;var a=t.return;try{if(bA(e,a,t,n,ee)){qe=1,If(e,Ai(n,e.current)),$t=null;return}}catch(s){if(a!==null)throw $t=a,s;qe=1,If(e,Ai(n,e.current)),$t=null;return}t.flags&32768?(Qt||i===1?e=!0:zo||(ee&536870912)!==0?e=!1:(fs=e=!0,(i===2||i===9||i===3||i===6)&&(i=Tn.current,i!==null&&i.tag===13&&(i.flags|=16384))),Db(t,e)):ch(t)}function ch(e){var t=e;do{if((t.flags&32768)!==0){Db(t,fs);return}e=t.return;var n=AA(t.alternate,t,Oa);if(n!==null){$t=n;return}if(t=t.sibling,t!==null){$t=t;return}$t=t=e}while(t!==null);qe===0&&(qe=5)}function Db(e,t){do{var n=wA(e.alternate,e);if(n!==null){n.flags&=32767,$t=n;return}if(n=e.return,n!==null&&(n.flags|=32768,n.subtreeFlags=0,n.deletions=null),!t&&(e=e.sibling,e!==null)){$t=e;return}$t=e=n}while(e!==null);qe=6,$t=null}function Ly(e,t,n,i,a,s,r,o,l,c,f,p){e.cancelPendingCommit=null;do uh();while(He!==0);if((fe&6)!==0)throw Error(J(327));if(t!==null){if(t===e.current)throw Error(J(177));e===Ae&&($t=Ae=null,ee=0),gr=t,zi=e,ra=n,sg=a,bb=i,PA(e,t,n,r,o,l,p)}}function PA(e,t,n,i,a,s,r){var o=t.lanes|t.childLanes;if(ag=o,o|=Lg,lE(e,n,o,i,a,s),yo=null,(n&335544064)===n?(xo=uA(e),i=10262):(xo=null,i=10256),(t.subtreeFlags&i)!==0||(t.flags&i)!==0?(e.callbackNode=null,e.callbackPriority=0,GA(Sf,function(){return cg(),null})):(e.callbackNode=null,e.callbackPriority=0),Bf=!1,i=(t.flags&13878)!==0,(t.subtreeFlags&13878)!==0||i){i=Bt.T,Bt.T=null,a=he.p,he.p=2,s=fe,fe|=4;try{CA(e,t,n)}finally{fe=s,he.p=a,Bt.T=i}}He=1,Bf?_o=rw(r,e.containerInfo,xo,rg,og,zA,lg,cg,BA,null,null):(rg(),og(),lg())}function BA(e){if(He!==0){var t=zi.onRecoverableError;t(e,{componentStack:null})}}function zA(){He===3&&(He=0,mb(gr,zi),He=4)}function rg(){if(He===1){He=0;var e=zi,t=gr,n=ra,i=(t.flags&13878)!==0;if((t.subtreeFlags&13878)!==0||i){i=Bt.T,Bt.T=null;var a=he.p;he.p=2;var s=fe;fe|=4;try{Gl=zf=!1,db(t,e,n),n=dg;var r=Wx(e.containerInfo),o=n.focusedElem,l=n.selectionRange;if(r!==o&&o&&o.ownerDocument&&Xx(o.ownerDocument.documentElement,o)){if(l!==null&&Ug(o)){var c=l.start,f=l.end;if(f===void 0&&(f=c),"selectionStart"in o)o.selectionStart=c,o.selectionEnd=Math.min(f,o.value.length);else{var p=o.ownerDocument||document,u=p&&p.defaultView||window;if(u.getSelection){var d=u.getSelection(),_=o.textContent.length,b=Math.min(l.start,_),v=l.end===void 0?b:Math.min(l.end,_);!d.extend&&b>v&&(r=v,v=b,b=r);var h=J_(o,b),m=J_(o,v);if(h&&m&&(d.rangeCount!==1||d.anchorNode!==h.node||d.anchorOffset!==h.offset||d.focusNode!==m.node||d.focusOffset!==m.offset)){var g=p.createRange();g.setStart(h.node,h.offset),d.removeAllRanges(),b>v?(d.addRange(g),d.extend(m.node,m.offset)):(g.setEnd(m.node,m.offset),d.addRange(g))}}}}for(p=[],d=o;d=d.parentNode;)d.nodeType===1&&p.push({element:d,left:d.scrollLeft,top:d.scrollTop});for(typeof o.focus=="function"&&o.focus(),o=0;o<p.length;o++){var x=p[o];x.element.scrollLeft=x.left,x.element.scrollTop=x.top}}Uo=!!hg,dg=hg=null}finally{fe=s,he.p=a,Bt.T=i}}e.current=t,He=2}}function og(){if(He===2){He=0;var e=zi,t=gr,n=(t.flags&8772)!==0;if((t.subtreeFlags&8772)!==0||n){n=Bt.T,Bt.T=null;var i=he.p;he.p=2;var a=fe;fe|=4;try{ob(e,t.alternate,t)}finally{fe=a,he.p=i,Bt.T=n}}He=3}}function lg(){if(He===4||He===3){He=0;var e=_o;_o=null,$T();var t=zi,n=gr,i=ra,a=bb,s=(i&335544064)===i?10262:10256;if((n.subtreeFlags&s)!==0||(n.flags&s)!==0?He=5:(He=0,gr=zi=null,Nb(t,t.pendingLanes)),s=t.pendingLanes,s===0&&(xs=null),Ag(i),n=n.stateNode,ci&&typeof ci.onCommitFiberRoot=="function")try{ci.onCommitFiberRoot(yc,n,void 0,(n.current.flags&128)===128)}catch{}if(a!==null){n=Bt.T,s=he.p,he.p=2,Bt.T=null;try{for(var r=t.onRecoverableError,o=0;o<a.length;o++){var l=a[o];r(l.value,{componentStack:l.stack})}}finally{Bt.T=n,he.p=s}}if(a=yo,r=xo,xo=null,a!==null&&(yo=null,r===null&&(r=[]),e!==null))for(l=0;l<a.length;l++)n=(0,a[l])(r),n!==void 0&&e.finished.finally(n);(ra&3)!==0&&uh(),ca(t),s=t.pendingLanes,(i&261930)!==0&&(s&42)!==0?t===hf?ec++:(ec=0,hf=t):(ec=0,hf=null),wc(0,!1)}}function Nb(e,t){(e.pooledCacheLanes&=t)===0&&(t=e.pooledCache,t!=null&&(e.pooledCache=null,Tc(t)))}function uh(){return _o!==null&&(_o.skipTransition(),_o=null),rg(),og(),lg(),cg()}function cg(){if(He!==5)return!1;var e=zi,t=ag;ag=0;var n=Ag(ra),i=Bt.T,a=he.p;try{he.p=32>n?32:n,Bt.T=null,n=sg,sg=null;var s=zi,r=ra;if(He=0,gr=zi=null,ra=0,(fe&6)!==0)throw Error(J(331));var o=fe;if(fe|=4,yb(s.current),gb(s,s.current,r,n),fe=o,wc(0,!1),ci&&typeof ci.onPostCommitFiberRoot=="function")try{ci.onPostCommitFiberRoot(yc,s)}catch{}return!0}finally{he.p=a,Bt.T=i,Nb(e,t)}}function Iy(e,t,n){t=Ai(n,t),t=Gm(e.stateNode,t,2),e=vs(e,t,2),e!==null&&(Sc(e,2),ca(e))}function ye(e,t,n){if(e.tag===3)Iy(e,e,n);else for(;t!==null;){if(t.tag===3){Iy(t,e,n);break}else if(t.tag===1){var i=t.stateNode;if(typeof t.type.getDerivedStateFromError=="function"||typeof i.componentDidCatch=="function"&&(xs===null||!xs.has(i))){e=Ai(n,e),n=qS(2),i=vs(t,n,2),i!==null&&(YS(n,i,t,e),Sc(i,2),ca(i));break}}t=t.return}}function am(e,t,n){var i=e.pingCache;if(i===null){i=e.pingCache=new NA;var a=new Set;i.set(t,a)}else a=i.get(t),a===void 0&&(a=new Set,i.set(t,a));a.has(n)||(c0=!0,a.add(n),e=FA.bind(null,e,t,n),t.then(e,e))}function FA(e,t,n){var i=e.pingCache;i!==null&&i.delete(t),e.pingedLanes|=e.suspendedLanes&n,e.warmLanes&=~n,Ae===e&&(ee&n)===n&&((qe===4||qe===3&&(ee&62914560)===ee&&300>li()-oh)&&(fe&2)===0?wo(e,0):Ff|=n,Eo===ee&&(Eo=0)),ca(e)}function Ub(e,t){t===0&&(t=bx()),e=xr(e,t),e!==null&&(Sc(e,t),ca(e))}function VA(e){var t=e.memoizedState,n=0;t!==null&&(n=t.retryLane),Ub(e,n)}function HA(e,t){var n=0;switch(e.tag){case 31:case 13:var i=e.stateNode,a=e.memoizedState;a!==null&&(n=a.retryLane);break;case 19:i=e.stateNode;break;case 22:i=e.stateNode._retryCache;break;default:throw Error(J(314))}i!==null&&i.delete(t),Ub(e,n)}function GA(e,t){return Tg(e,t)}var Co=null,Qr=null,ug=!1,kf=!1,sm=!1,ds=0;function ca(e){e!==Qr&&e.next===null&&(Qr===null?Co=Qr=e:Qr=Qr.next=e),kf=!0,ug||(ug=!0,XA())}function wc(e,t){if(!sm&&kf){sm=!0;do for(var n=!1,i=Co;i!==null;){if(!t)if(e!==0){var a=i.pendingLanes;if(a===0)var s=0;else{var r=i.suspendedLanes,o=i.pingedLanes;s=(1<<31-ui(42|e)+1)-1,s&=a&~(r&~o),s=s&201326741?s&201326741|1:s?s|2:0}s!==0&&(n=!0,Oy(i,s))}else s=ee,s=Yf(i,i===Ae?s:0,i.cancelPendingCommit!==null||i.timeoutHandle!==-1),(s&3)===0||xc(i,s)||(n=!0,Oy(i,s));i=i.next}while(n);sm=!1}}function kA(){Lb()}function Lb(){kf=ug=!1;var e=0;ds!==0&&$A()&&(e=ds);for(var t=li(),n=null,i=Co;i!==null;){var a=i.next,s=Ib(i,t);s===0?(i.next=null,n===null?Co=a:n.next=a,a===null&&(Qr=n)):(n=i,(e!==0||(s&3)!==0)&&(kf=!0)),i=a}He!==0&&He!==5||wc(e,!1),ds!==0&&(ds=0)}function Ib(e,t){for(var n=e.suspendedLanes,i=e.pingedLanes,a=e.expirationTimes,s=e.pendingLanes&-62914561;0<s;){var r=31-ui(s),o=1<<r,l=a[r];l===-1?((o&n)===0||(o&i)!==0)&&(a[r]=oE(o,t)):l<=t&&(e.expiredLanes|=o),s&=~o}if(t=Ae,n=ee,n=Yf(e,e===t?n:0,e.cancelPendingCommit!==null||e.timeoutHandle!==-1),i=e.callbackNode,n===0||e===t&&(ve===2||ve===9)||e.cancelPendingCommit!==null)return i!==null&&i!==null&&Bp(i),e.callbackNode=null,e.callbackPriority=0;if((n&3)===0||xc(e,n)){if(t=n&-n,t===e.callbackPriority)return t;switch(i!==null&&Bp(i),Ag(n)){case 2:case 8:n=yx;break;case 32:n=Sf;break;case 268435456:n=xx;break;default:n=Sf}return i=Ob.bind(null,e),n=Tg(n,i),e.callbackPriority=t,e.callbackNode=n,t}return i!==null&&i!==null&&Bp(i),e.callbackPriority=2,e.callbackNode=null,2}function Ob(e,t){if(He!==0&&He!==5)return e.callbackNode=null,e.callbackPriority=0,null;var n=e.callbackNode;if(uh()&&e.callbackNode!==n)return null;var i=ee;return i=Yf(e,e===Ae?i:0,e.cancelPendingCommit!==null||e.timeoutHandle!==-1),i===0?null:(Tb(e,i,t),Ib(e,li()),e.callbackNode!=null&&e.callbackNode===n?Ob.bind(null,e):null)}function Oy(e,t){if(uh())return null;Tb(e,t,!0)}function XA(){ew(function(){(fe&6)!==0?Tg(_x,kA):Lb()})}function f0(){if(ds===0){var e=hr;e===0&&(e=Uu,Uu<<=1,(Uu&261888)===0&&(Uu=256)),ds=e}return ds}function Py(e){return e==null||typeof e=="symbol"||typeof e=="boolean"?null:typeof e=="function"?e:$u(e)}function WA(e,t,n,i,a){if(t==="submit"&&n&&n.stateNode===a){var s=Py((a[ti]||null).action),r=i.submitter;r&&(t=(t=r[ti]||null)?Py(t.formAction):r.getAttribute("formAction"),t!==null&&(s=t,r=null));var o=new Jf("action","action",null,i,a);e.push({event:o,listeners:[{instance:null,listener:function(){if(i.defaultPrevented){if(ds!==0){var l=new FormData(a,r);Vm(n,{pending:!0,data:l,method:a.method,action:s},null,l)}}else typeof s=="function"&&(o.preventDefault(),l=new FormData(a,r),Vm(n,{pending:!0,data:l,method:a.method,action:s},s,l))},currentTarget:a}]})}}for(Yu=0;Yu<Rm.length;Yu++)Zu=Rm[Yu],By=Zu.toLowerCase(),zy=Zu[0].toUpperCase()+Zu.slice(1),Fi(By,"on"+zy);var Zu,By,zy,Yu;Fi(Yx,"onAnimationEnd");Fi(Zx,"onAnimationIteration");Fi(Jx,"onAnimationStart");Fi("dblclick","onDoubleClick");Fi("focusin","onFocus");Fi("focusout","onBlur");Fi(nA,"onTransitionRun");Fi(iA,"onTransitionStart");Fi(aA,"onTransitionCancel");Fi(Kx,"onTransitionEnd");bo("onMouseEnter",["mouseout","mouseover"]);bo("onMouseLeave",["mouseout","mouseover"]);bo("onPointerEnter",["pointerout","pointerover"]);bo("onPointerLeave",["pointerout","pointerover"]);_r("onChange","change click focusin focusout input keydown keyup selectionchange".split(" "));_r("onSelect","focusout contextmenu dragend focusin keydown keyup mousedown mouseup selectionchange".split(" "));_r("onBeforeInput",["compositionend","keypress","textInput","paste"]);_r("onCompositionEnd","compositionend focusout keydown keypress keyup mousedown".split(" "));_r("onCompositionStart","compositionstart focusout keydown keypress keyup mousedown".split(" "));_r("onCompositionUpdate","compositionupdate focusout keydown keypress keyup mousedown".split(" "));var fc="abort canplay canplaythrough durationchange emptied encrypted ended error loadeddata loadedmetadata loadstart pause play playing progress ratechange resize seeked seeking stalled suspend timeupdate volumechange waiting".split(" "),qA=new Set("beforetoggle cancel close invalid load scroll scrollend toggle".split(" ").concat(fc));function Pb(e,t){t=(t&4)!==0;for(var n=0;n<e.length;n++){var i=e[n],a=i.event;i=i.listeners;t:{var s=void 0;if(t)for(var r=i.length-1;0<=r;r--){var o=i[r],l=o.instance,c=o.currentTarget;if(o=o.listener,l!==s&&a.isPropagationStopped())break t;s=o,a.currentTarget=c;try{s(a)}catch(f){Mf(f)}a.currentTarget=null,s=l}else for(r=0;r<i.length;r++){if(o=i[r],l=o.instance,c=o.currentTarget,o=o.listener,l!==s&&a.isPropagationStopped())break t;s=o,a.currentTarget=c;try{s(a)}catch(f){Mf(f)}a.currentTarget=null,s=l}}}}function jt(e,t){var n=t[U_];n===void 0&&(n=t[U_]=new Set);var i=e+"__bubble";n.has(i)||(Bb(t,e,2,!1),n.add(i))}function rm(e,t,n){var i=0;t&&(i|=4),Bb(n,e,i,t)}var Ju="_reactListening"+Math.random().toString(36).slice(2);function h0(e){if(!e[Ju]){e[Ju]=!0,Cx.forEach(function(n){n!=="selectionchange"&&(qA.has(n)||rm(n,!1,e),rm(n,!0,e))});var t=e.nodeType===9?e:e.ownerDocument;t===null||t[Ju]||(t[Ju]=!0,rm("selectionchange",!1,t))}}function Bb(e,t,n,i){switch(oM(t)){case 2:var a=zw;break;case 8:a=Fw;break;default:a=y0}n=a.bind(null,t,n,e),a=void 0,!Em||t!=="touchstart"&&t!=="touchmove"&&t!=="wheel"||(a=!0),i?a!==void 0?e.addEventListener(t,n,{capture:!0,passive:a}):e.addEventListener(t,n,!0):a!==void 0?e.addEventListener(t,n,{passive:a}):e.addEventListener(t,n,!1)}function om(e,t,n,i,a){var s=i;if((t&1)===0&&(t&2)===0&&i!==null)t:for(;;){if(i===null)return;var r=i.tag;if(r===3||r===4){var o=i.stateNode.containerInfo;if(o===a)break;if(r===4)for(r=i.return;r!==null;){var l=r.tag;if((l===3||l===4)&&r.stateNode.containerInfo===a)return;r=r.return}for(;o!==null;){if(r=nr(o),r===null)return;if(l=r.tag,l===5||l===6||l===26||l===27){i=s=r;continue t}o=o.parentNode}}i=i.return}Px(function(){var c=s,f=Cg(n),p=[];t:{var u=Qx.get(e);if(u!==void 0){var d=Jf,_=e;switch(e){case"keypress":if(ef(n)===0)break t;case"keydown":case"keyup":d=UE;break;case"focusin":_="focus",d=kp;break;case"focusout":_="blur",d=kp;break;case"beforeblur":case"afterblur":d=kp;break;case"click":if(n.button===2)break t;case"auxclick":case"dblclick":case"mousedown":case"mousemove":case"mouseup":case"mouseout":case"mouseover":case"contextmenu":d=V_;break;case"drag":case"dragend":case"dragenter":case"dragexit":case"dragleave":case"dragover":case"dragstart":case"drop":d=xE;break;case"touchcancel":case"touchend":case"touchmove":case"touchstart":d=BE;break;case Yx:case Zx:case Jx:d=ME;break;case Kx:d=FE;break;case"scroll":case"scrollend":d=_E;break;case"wheel":d=HE;break;case"copy":case"cut":case"paste":d=EE;break;case"gotpointercapture":case"lostpointercapture":case"pointercancel":case"pointerdown":case"pointermove":case"pointerout":case"pointerover":case"pointerup":d=G_;break;case"submit":d=OE;break;case"toggle":case"beforetoggle":d=kE}var b=(t&4)!==0,v=!b&&(e==="scroll"||e==="scrollend"),h=b?u!==null?u+"Capture":null:u;b=[];for(var m=c,g;m!==null;){var x=m;if(g=x.stateNode,x=x.tag,x!==5&&x!==26&&x!==27||g===null||h===null||(x=ic(m,h),x!=null&&b.push(hc(m,x,g))),v)break;m=m.return}0<b.length&&(u=new d(u,_,null,n,f),p.push({event:u,listeners:b}))}}if((t&7)===0){t:{if(d=e==="mouseover"||e==="pointerover",u=e==="mouseout"||e==="pointerout",d&&n!==Tm&&(_=n.relatedTarget||n.fromElement)&&(nr(_)||_[Io]))break t;(u||d)&&(_=f.window===f?f:(d=f.ownerDocument)?d.defaultView||d.parentWindow:window,u?(d=n.relatedTarget||n.toElement,u=c,d=d?nr(d):null,d!==null&&(v=_c(d),b=d.tag,d!==v||b!==5&&b!==27&&b!==6)&&(d=null)):(u=null,d=c),u!==d&&(b=V_,x="onMouseLeave",h="onMouseEnter",m="mouse",(e==="pointerout"||e==="pointerover")&&(b=G_,x="onPointerLeave",h="onPointerEnter",m="pointer"),v=u==null?_:Vl(u),g=d==null?_:Vl(d),_=new b(x,m+"leave",u,n,f),_.target=v,_.relatedTarget=g,x=null,nr(f)===c&&(b=new b(h,m+"enter",d,n,f),b.target=g,b.relatedTarget=v,x=b),v=x,b=u&&d?hm(u,d,YA):null,u!==null&&Fy(p,_,u,b,!1),d!==null&&v!==null&&Fy(p,v,d,b,!0)))}t:{if(u=c?Vl(c):window,d=u.nodeName&&u.nodeName.toLowerCase(),d==="select"||d==="input"&&u.type==="file")var E=q_;else if(W_(u))if(Gx)E=$E;else{E=QE;var A=KE}else d=u.nodeName,!d||d.toLowerCase()!=="input"||u.type!=="checkbox"&&u.type!=="radio"?c&&wg(c.elementType)&&(E=q_):E=jE;if(E&&(E=E(e,c))){Hx(p,E,n,f);break t}A&&A(e,u,c)}switch(A=c?Vl(c):window,e){case"focusin":(W_(A)||A.contentEditable==="true")&&(io=A,wm=c,Wl=null);break;case"focusout":Wl=wm=io=null;break;case"mousedown":Cm=!0;break;case"contextmenu":case"mouseup":case"dragend":Cm=!1,K_(p,n,f);break;case"selectionchange":if(eA)break;case"keydown":case"keyup":K_(p,n,f)}var w;if(Ng)t:{switch(e){case"compositionstart":var y="onCompositionStart";break t;case"compositionend":y="onCompositionEnd";break t;case"compositionupdate":y="onCompositionUpdate";break t}y=void 0}else no?Fx(e,n)&&(y="onCompositionEnd"):e==="keydown"&&n.keyCode===229&&(y="onCompositionStart");y&&(zx&&n.locale!=="ko"&&(no||y!=="onCompositionStart"?y==="onCompositionEnd"&&no&&(w=Bx()):(cs=f,Rg="value"in cs?cs.value:cs.textContent,no=!0)),A=Xf(c,y),0<A.length&&(y=new H_(y,e,null,n,f),p.push({event:y,listeners:A}),w?y.data=w:(w=Vx(n),w!==null&&(y.data=w)))),(w=WE?qE(e,n):YE(e,n))&&(y=Xf(c,"onBeforeInput"),0<y.length&&(A=new H_("onBeforeInput","beforeinput",null,n,f),p.push({event:A,listeners:y}),A.data=w)),WA(p,e,c,n,f)}Pb(p,t)})}function hc(e,t,n){return{instance:e,listener:t,currentTarget:n}}function Xf(e,t){for(var n=t+"Capture",i=[];e!==null;){var a=e,s=a.stateNode;if(a=a.tag,a!==5&&a!==26&&a!==27||s===null||(a=ic(e,n),a!=null&&i.unshift(hc(e,a,s)),a=ic(e,t),a!=null&&i.push(hc(e,a,s))),e.tag===3)return i;e=e.return}return[]}function YA(e){if(e===null)return null;do e=e.return;while(e&&e.tag!==5&&e.tag!==27);return e||null}function Fy(e,t,n,i,a){for(var s=t._reactName,r=[];n!==null&&n!==i;){var o=n,l=o.alternate,c=o.stateNode;if(o=o.tag,l!==null&&l===i)break;o!==5&&o!==26&&o!==27||c===null||(l=c,a?(c=ic(n,s),c!=null&&r.unshift(hc(n,c,l))):a||(c=ic(n,s),c!=null&&r.push(hc(n,c,l)))),n=n.return}r.length!==0&&e.push({event:t,listeners:r})}var ZA=/\r\n?/g,JA=/\u0000|\uFFFD/g;function Vy(e){return(typeof e=="string"?e:""+e).replace(ZA,`
`).replace(JA,"")}function zb(e,t){return t=Vy(t),Vy(e)===t}function _e(e,t,n,i,a,s){switch(n){case"children":if(typeof i=="string")t==="body"||t==="textarea"&&i===""||Mo(e,i);else if(typeof i=="number"||typeof i=="bigint")t!=="body"&&Mo(e,""+i);else return;break;case"className":Ou(e,"class",i);break;case"tabIndex":Ou(e,"tabindex",i);break;case"dir":case"role":case"viewBox":case"width":case"height":Ou(e,n,i);break;case"style":Ox(e,i,s);return;case"data":if(t!=="object"){Ou(e,"data",i);break}case"src":case"href":if(i===""&&(t!=="a"||n!=="href")){e.removeAttribute(n);break}if(i==null||typeof i=="function"||typeof i=="symbol"||typeof i=="boolean"){e.removeAttribute(n);break}i=$u(i),e.setAttribute(n,i);break;case"action":case"formAction":if(typeof i=="function"){e.setAttribute(n,"javascript:throw new Error('A React form was unexpectedly submitted. If you called form.submit() manually, consider using form.requestSubmit() instead. If you\\'re trying to use event.stopPropagation() in a submit event handler, consider also calling event.preventDefault().')");break}else typeof s=="function"&&(n==="formAction"?(t!=="input"&&_e(e,t,"name",a.name,a,null),_e(e,t,"formEncType",a.formEncType,a,null),_e(e,t,"formMethod",a.formMethod,a,null),_e(e,t,"formTarget",a.formTarget,a,null)):(_e(e,t,"encType",a.encType,a,null),_e(e,t,"method",a.method,a,null),_e(e,t,"target",a.target,a,null)));if(i==null||typeof i=="symbol"||typeof i=="boolean"){e.removeAttribute(n);break}i=$u(i),e.setAttribute(n,i);break;case"onClick":i!=null&&(e.onclick=na);return;case"onScroll":i!=null&&jt("scroll",e);return;case"onScrollEnd":i!=null&&jt("scrollend",e);return;case"dangerouslySetInnerHTML":if(i!=null){if(typeof i!="object"||!("__html"in i))throw Error(J(61));if(n=i.__html,n!=null){if(a.children!=null)throw Error(J(60));s?.__html!==n&&(e.innerHTML=n)}}break;case"multiple":e.multiple=i&&typeof i!="function"&&typeof i!="symbol";break;case"muted":e.muted=i&&typeof i!="function"&&typeof i!="symbol";break;case"suppressContentEditableWarning":case"suppressHydrationWarning":case"defaultValue":case"defaultChecked":case"innerHTML":case"ref":break;case"autoFocus":break;case"xlinkHref":if(i==null||typeof i=="function"||typeof i=="boolean"||typeof i=="symbol"){e.removeAttribute("xlink:href");break}n=$u(i),e.setAttributeNS("http://www.w3.org/1999/xlink","xlink:href",n);break;case"contentEditable":case"spellCheck":case"draggable":case"value":case"autoReverse":case"externalResourcesRequired":case"focusable":case"preserveAlpha":i!=null&&typeof i!="function"&&typeof i!="symbol"?e.setAttribute(n,i):e.removeAttribute(n);break;case"inert":case"allowFullScreen":case"async":case"autoPlay":case"controls":case"credentialless":case"default":case"defer":case"disabled":case"disablePictureInPicture":case"disableRemotePlayback":case"formNoValidate":case"hidden":case"loop":case"noModule":case"noValidate":case"open":case"playsInline":case"readOnly":case"required":case"reversed":case"scoped":case"seamless":case"itemScope":i&&typeof i!="function"&&typeof i!="symbol"?e.setAttribute(n,""):e.removeAttribute(n);break;case"capture":case"download":i===!0?e.setAttribute(n,""):i!==!1&&i!=null&&typeof i!="function"&&typeof i!="symbol"?e.setAttribute(n,i):e.removeAttribute(n);break;case"cols":case"rows":case"size":case"span":i!=null&&typeof i!="function"&&typeof i!="symbol"&&!isNaN(i)&&1<=i?e.setAttribute(n,i):e.removeAttribute(n);break;case"rowSpan":case"start":i==null||typeof i=="function"||typeof i=="symbol"||isNaN(i)?e.removeAttribute(n):e.setAttribute(n,i);break;case"popover":jt("beforetoggle",e),jt("toggle",e),ju(e,"popover",i);break;case"xlinkActuate":Ta(e,"http://www.w3.org/1999/xlink","xlink:actuate",i);break;case"xlinkArcrole":Ta(e,"http://www.w3.org/1999/xlink","xlink:arcrole",i);break;case"xlinkRole":Ta(e,"http://www.w3.org/1999/xlink","xlink:role",i);break;case"xlinkShow":Ta(e,"http://www.w3.org/1999/xlink","xlink:show",i);break;case"xlinkTitle":Ta(e,"http://www.w3.org/1999/xlink","xlink:title",i);break;case"xlinkType":Ta(e,"http://www.w3.org/1999/xlink","xlink:type",i);break;case"xmlBase":Ta(e,"http://www.w3.org/XML/1998/namespace","xml:base",i);break;case"xmlLang":Ta(e,"http://www.w3.org/XML/1998/namespace","xml:lang",i);break;case"xmlSpace":Ta(e,"http://www.w3.org/XML/1998/namespace","xml:space",i);break;case"is":ju(e,"is",i);break;case"innerText":case"textContent":return;default:if(!(2<n.length)||n[0]!=="o"&&n[0]!=="O"||n[1]!=="n"&&n[1]!=="N")n=gE.get(n)||n,ju(e,n,i);else return}ce=!0}function fg(e,t,n,i,a,s){switch(n){case"style":Ox(e,i,s);return;case"dangerouslySetInnerHTML":if(i!=null){if(typeof i!="object"||!("__html"in i))throw Error(J(61));if(n=i.__html,n!=null){if(a.children!=null)throw Error(J(60));s?.__html!==n&&(e.innerHTML=n)}}break;case"children":if(typeof i=="string")Mo(e,i);else if(typeof i=="number"||typeof i=="bigint")Mo(e,""+i);else return;break;case"onScroll":i!=null&&jt("scroll",e);return;case"onScrollEnd":i!=null&&jt("scrollend",e);return;case"onClick":i!=null&&(e.onclick=na);return;case"suppressContentEditableWarning":case"suppressHydrationWarning":case"innerHTML":case"ref":return;case"innerText":case"textContent":return;default:if(!Rx.hasOwnProperty(n))t:{if(n[0]==="o"&&n[1]==="n"&&(a=n.endsWith("Capture"),s=n.slice(2,a?n.length-7:void 0),t=e[ti]||null,t=t!=null?t[n]:null,typeof t=="function"&&e.removeEventListener(s,t,a),typeof i=="function")){typeof t!="function"&&t!==null&&(n in e?e[n]=null:e.hasAttribute(n)&&e.removeAttribute(n)),e.addEventListener(s,i,a);break t}ce=!0,n in e?e[n]=i:i===!0?e.setAttribute(n,""):ju(e,n,i)}return}ce=!0}function Mn(e,t,n){switch(t){case"div":case"span":case"svg":case"path":case"a":case"g":case"p":case"li":break;case"img":jt("error",e),jt("load",e);var i=!1,a=!1,s;for(s in n)if(n.hasOwnProperty(s)){var r=n[s];if(r!=null)switch(s){case"src":i=!0;break;case"srcSet":a=!0;break;case"children":case"dangerouslySetInnerHTML":throw Error(J(137,t));default:_e(e,t,s,r,n,null)}}a&&_e(e,t,"srcSet",n.srcSet,n,null),i&&_e(e,t,"src",n.src,n,null);return;case"input":jt("invalid",e);var o=s=r=a=null,l=null,c=null;for(i in n)if(n.hasOwnProperty(i)){var f=n[i];if(f!=null)switch(i){case"name":a=f;break;case"type":r=f;break;case"checked":l=f;break;case"defaultChecked":c=f;break;case"value":s=f;break;case"defaultValue":o=f;break;case"children":case"dangerouslySetInnerHTML":if(f!=null)throw Error(J(137,t));break;default:_e(e,t,i,f,n,null)}}Ux(e,s,o,l,c,r,a,!1);return;case"select":jt("invalid",e),i=r=s=null;for(a in n)if(n.hasOwnProperty(a)&&(o=n[a],o!=null))switch(a){case"value":s=o;break;case"defaultValue":r=o;break;case"multiple":i=o;default:_e(e,t,a,o,n,null)}t=s,n=r,e.multiple=!!i,t!=null?fo(e,!!i,t,!1):n!=null&&fo(e,!!i,n,!0);return;case"textarea":jt("invalid",e),s=a=i=null;for(r in n)if(n.hasOwnProperty(r)&&(o=n[r],o!=null))switch(r){case"value":i=o;break;case"defaultValue":a=o;break;case"children":s=o;break;case"dangerouslySetInnerHTML":if(o!=null)throw Error(J(91));break;default:_e(e,t,r,o,n,null)}Ix(e,i,a,s);return;case"option":for(l in n)n.hasOwnProperty(l)&&(i=n[l],i!=null)&&(l==="selected"?e.selected=i&&typeof i!="function"&&typeof i!="symbol":_e(e,t,l,i,n,null));return;case"dialog":jt("beforetoggle",e),jt("toggle",e),jt("cancel",e),jt("close",e);break;case"iframe":case"object":jt("load",e);break;case"video":case"audio":for(i=0;i<fc.length;i++)jt(fc[i],e);break;case"image":jt("error",e),jt("load",e);break;case"details":jt("toggle",e);break;case"embed":case"source":case"link":jt("error",e),jt("load",e);case"area":case"base":case"br":case"col":case"hr":case"keygen":case"meta":case"param":case"track":case"wbr":case"menuitem":for(c in n)if(n.hasOwnProperty(c)&&(i=n[c],i!=null))switch(c){case"children":case"dangerouslySetInnerHTML":throw Error(J(137,t));default:_e(e,t,c,i,n,null)}return;default:if(wg(t)){for(f in n)n.hasOwnProperty(f)&&(i=n[f],i!==void 0&&fg(e,t,f,i,n,void 0));return}}for(o in n)n.hasOwnProperty(o)&&(i=n[o],i!=null&&_e(e,t,o,i,n,null))}var KA={};function QA(e,t,n,i){switch(t){case"div":case"span":case"svg":case"path":case"a":case"g":case"p":case"li":break;case"input":var a=null,s=null,r=null,o=null,l=null,c=null,f=null;for(d in n){var p=n[d];if(n.hasOwnProperty(d)&&p!=null)switch(d){case"checked":break;case"value":break;case"defaultValue":l=p;default:i.hasOwnProperty(d)||_e(e,t,d,null,i,p)}}for(var u in i){var d=i[u];if(p=n[u],i.hasOwnProperty(u)&&(d!=null||p!=null))switch(u){case"type":d!==p&&(ce=!0),s=d;break;case"name":d!==p&&(ce=!0),a=d;break;case"checked":d!==p&&(ce=!0),c=d;break;case"defaultChecked":d!==p&&(ce=!0),f=d;break;case"value":d!==p&&(ce=!0),r=d;break;case"defaultValue":d!==p&&(ce=!0),o=d;break;case"children":case"dangerouslySetInnerHTML":if(d!=null)throw Error(J(137,t));break;default:d!==p&&_e(e,t,u,d,i,p)}}Mm(e,r,o,l,c,f,s,a);return;case"select":d=r=o=u=null;for(s in n)if(l=n[s],n.hasOwnProperty(s)&&l!=null)switch(s){case"value":break;case"multiple":d=l;default:i.hasOwnProperty(s)||_e(e,t,s,null,i,l)}for(a in i)if(s=i[a],l=n[a],i.hasOwnProperty(a)&&(s!=null||l!=null))switch(a){case"value":s!==l&&(ce=!0),u=s;break;case"defaultValue":s!==l&&(ce=!0),o=s;break;case"multiple":s!==l&&(ce=!0),r=s;default:s!==l&&_e(e,t,a,s,i,l)}t=o,n=r,i=d,u!=null?fo(e,!!n,u,!1):!!i!=!!n&&(t!=null?fo(e,!!n,t,!0):fo(e,!!n,n?[]:"",!1));return;case"textarea":d=u=null;for(o in n)if(a=n[o],n.hasOwnProperty(o)&&a!=null&&!i.hasOwnProperty(o))switch(o){case"value":break;case"children":break;default:_e(e,t,o,null,i,a)}for(r in i)if(a=i[r],s=n[r],i.hasOwnProperty(r)&&(a!=null||s!=null))switch(r){case"value":a!==s&&(ce=!0),u=a;break;case"defaultValue":a!==s&&(ce=!0),d=a;break;case"children":break;case"dangerouslySetInnerHTML":if(a!=null)throw Error(J(91));break;default:a!==s&&_e(e,t,r,a,i,s)}Lx(e,u,d);return;case"option":for(var _ in n)u=n[_],n.hasOwnProperty(_)&&u!=null&&!i.hasOwnProperty(_)&&(_==="selected"?e.selected=!1:_e(e,t,_,null,i,u));for(l in i)u=i[l],d=n[l],i.hasOwnProperty(l)&&u!==d&&(u!=null||d!=null)&&(l==="selected"?(u!==d&&(ce=!0),e.selected=u&&typeof u!="function"&&typeof u!="symbol"):_e(e,t,l,u,i,d));return;case"img":case"link":case"area":case"base":case"br":case"col":case"embed":case"hr":case"keygen":case"meta":case"param":case"source":case"track":case"wbr":case"menuitem":for(var b in n)u=n[b],n.hasOwnProperty(b)&&u!=null&&!i.hasOwnProperty(b)&&_e(e,t,b,null,i,u);for(c in i)if(u=i[c],d=n[c],i.hasOwnProperty(c)&&u!==d&&(u!=null||d!=null))switch(c){case"children":case"dangerouslySetInnerHTML":if(u!=null)throw Error(J(137,t));break;default:_e(e,t,c,u,i,d)}return;default:if(wg(t)){for(var v in n)u=n[v],n.hasOwnProperty(v)&&u!==void 0&&!i.hasOwnProperty(v)&&fg(e,t,v,void 0,i,u);for(f in i)u=i[f],d=n[f],!i.hasOwnProperty(f)||u===d||u===void 0&&d===void 0||fg(e,t,f,u,i,d);return}}for(var h in n)u=n[h],n.hasOwnProperty(h)&&u!=null&&!i.hasOwnProperty(h)&&_e(e,t,h,null,i,u);for(p in i)u=i[p],d=n[p],!i.hasOwnProperty(p)||u===d||u==null&&d==null||_e(e,t,p,u,i,d)}function Hy(e){switch(e){case"css":case"script":case"font":case"img":case"image":case"input":case"link":return!0;default:return!1}}function jA(){if(typeof performance.getEntriesByType=="function"){for(var e=0,t=0,n=performance.getEntriesByType("resource"),i=0;i<n.length;i++){var a=n[i],s=a.transferSize,r=a.initiatorType,o=a.duration;if(s&&o&&Hy(r)){for(r=0,o=a.responseEnd,i+=1;i<n.length;i++){var l=n[i],c=l.startTime;if(c>o)break;var f=l.transferSize,p=l.initiatorType;f&&Hy(p)&&(l=l.responseEnd,r+=f*(l<o?1:(o-c)/(l-c)))}if(--i,t+=8*(s+r)/(a.duration/1e3),e++,10<e)break}}if(0<e)return t/e/1e6}return navigator.connection&&(e=navigator.connection.downlink,typeof e=="number")?e:5}var hg=null,dg=null;function dc(e){return e.nodeType===9?e:e.ownerDocument}function Gy(e){switch(e){case"http://www.w3.org/2000/svg":return 1;case"http://www.w3.org/1998/Math/MathML":return 2;default:return 0}}function Fb(e,t){if(e===0)switch(t){case"svg":return 1;case"math":return 2;default:return 0}return e===1&&t==="foreignObject"?0:e}function Vb(e,t,n,i){return n=dc(n).createElement(e),n[yn]=i,n[ti]=t,Mn(n,e,t),fn(n),n}function pg(e,t){return e==="textarea"||e==="noscript"||typeof t.children=="string"||typeof t.children=="number"||typeof t.children=="bigint"||typeof t.dangerouslySetInnerHTML=="object"&&t.dangerouslySetInnerHTML!==null&&t.dangerouslySetInnerHTML.__html!=null}var lm=null;function $A(){var e=window.event;return e&&e.type==="popstate"?e===lm?!1:(lm=e,!0):(lm=null,!1)}var d0=typeof setTimeout=="function"?setTimeout:void 0,tw=typeof clearTimeout=="function"?clearTimeout:void 0,ky=typeof Promise=="function"?Promise:void 0,Xy=typeof requestAnimationFrame=="function"?requestAnimationFrame:d0,ew=typeof queueMicrotask=="function"?queueMicrotask:typeof ky<"u"?function(e){return ky.resolve(null).then(e).catch(nw)}:d0;function nw(e){setTimeout(function(){throw e})}function Ds(e){return e==="head"}function Wy(e,t){var n=t,i=0;do{var a=n.nextSibling;if(e.removeChild(n),a&&a.nodeType===8)if(n=a.data,n==="/$"||n==="/&"){if(i===0){e.removeChild(a),Lo(t);return}i--}else if(n==="$"||n==="$?"||n==="$~"||n==="$!"||n==="&")i++;else if(n==="html")um(e.ownerDocument.documentElement);else if(n==="head"){n=e.ownerDocument.head,um(n);for(var s=n.firstChild;s;){var r=s.nextSibling,o=s.nodeName;s[bc]||o==="SCRIPT"||o==="STYLE"||o==="LINK"&&s.rel.toLowerCase()==="stylesheet"||n.removeChild(s),s=r}}else n==="body"&&um(e.ownerDocument.body);n=a}while(n);Lo(t)}function qy(e,t){var n=e;e=0;do{var i=n.nextSibling;if(n.nodeType===1?t?(n._stashedDisplay=n.style.display,n.style.display="none"):(n.style.display=n._stashedDisplay||"",n.getAttribute("style")===""&&n.removeAttribute("style")):n.nodeType===3&&(t?(n._stashedText=n.nodeValue,n.nodeValue=""):n.nodeValue=n._stashedText||""),i&&i.nodeType===8)if(n=i.data,n==="/$"){if(e===0)break;e--}else n!=="$"&&n!=="$?"&&n!=="$~"&&n!=="$!"||e++;n=i}while(n)}function Hb(e,t,n){if(t=CSS.escape(t)!==t?"r-"+btoa(t).replace(/=/g,""):t,e.style.viewTransitionName=t,n!=null&&(e.style.viewTransitionClass=n),n=getComputedStyle(e),n.display==="inline"){if(t=e.getClientRects(),t.length===1)var i=1;else for(var a=i=0;a<t.length;a++){var s=t[a];0<s.width&&0<s.height&&i++}i===1&&(e=e.style,e.display=t.length===1?"inline-block":"block",e.marginTop="-"+n.paddingTop,e.marginBottom="-"+n.paddingBottom)}}function Gb(e,t){e=e.style,t=t.style;var n=t!=null?t.hasOwnProperty("viewTransitionName")?t.viewTransitionName:t.hasOwnProperty("view-transition-name")?t["view-transition-name"]:null:null;e.viewTransitionName=n==null||typeof n=="boolean"?"":(""+n).trim(),n=t!=null?t.hasOwnProperty("viewTransitionClass")?t.viewTransitionClass:t.hasOwnProperty("view-transition-class")?t["view-transition-class"]:null:null,e.viewTransitionClass=n==null||typeof n=="boolean"?"":(""+n).trim(),e.display==="inline-block"&&(t==null?e.display=e.margin="":(n=t.display,e.display=n==null||typeof n=="boolean"?"":n,n=t.margin,n!=null?e.margin=n:(n=t.hasOwnProperty("marginTop")?t.marginTop:t["margin-top"],e.marginTop=n==null||typeof n=="boolean"?"":n,t=t.hasOwnProperty("marginBottom")?t.marginBottom:t["margin-bottom"],e.marginBottom=t==null||typeof t=="boolean"?"":t)))}function kb(e,t,n){return n=n.ownerDocument.defaultView,{rect:e,abs:t.position==="absolute"||t.position==="fixed",clip:t.clipPath!=="none"||t.overflow!=="visible"||t.filter!=="none"||t.mask!=="none"||t.mask!=="none"||t.borderRadius!=="0px",view:0<=e.bottom&&0<=e.right&&e.top<=n.innerHeight&&e.left<=n.innerWidth}}function mg(e){var t=e.getBoundingClientRect(),n=getComputedStyle(e);return kb(t,n,e)}function iw(e){var t=e.getBoundingClientRect();t=new DOMRect(t.x+2e4,t.y+2e4,t.width,t.height);var n=getComputedStyle(e);return kb(t,n,e)}function aw(e){return e.documentElement.clientHeight}function sw(e){this.addEventListener("load",e),this.addEventListener("error",e)}function rw(e,t,n,i,a,s,r,o,l){var c=t.nodeType===9?t:t.ownerDocument;try{var f=c.startViewTransition({update:function(){var u=c.defaultView,d=u.navigation&&u.navigation.transition,_=c.fonts.status;i();var b=[];if(_==="loaded"&&(aw(c),c.fonts.status==="loading"&&b.push(c.fonts.ready)),_=b.length,e!==null)for(var v=e.suspenseyImages,h=0,m=0;m<v.length;m++){var g=v[m];if(!g.complete){var x=g.getBoundingClientRect();if(0<x.bottom&&0<x.right&&x.top<u.innerHeight&&x.left<u.innerWidth){if(h+=eM(g),h>mf){b.length=_;break}g=new Promise(sw.bind(g)),b.push(g)}}}if(0<b.length)return u=Promise.race([Promise.all(b),new Promise(function(E){return setTimeout(E,500)})]).then(a,a),(d?Promise.allSettled([d.finished,u]):u).then(s,s);if(a(),d)return d.finished.then(s,s);s()},types:n});c.__reactViewTransition=f;var p=[];return f.ready.then(function(){for(var u=c.documentElement.getAnimations({subtree:!0}),d=0;d<u.length;d++){var _=u[d],b=_.effect,v=b.pseudoElement;if(v!=null&&v.startsWith("::view-transition")){p.push(_),_=b.getKeyframes();for(var h=v=void 0,m=!0,g=0;g<_.length;g++){var x=_[g],E=x.width;if(v===void 0)v=E;else if(v!==E){m=!1;break}if(E=x.height,h===void 0)h=E;else if(h!==E){m=!1;break}delete x.width,delete x.height,x.transform==="none"&&delete x.transform}m&&v!==void 0&&h!==void 0&&(b.setKeyframes(_),m=getComputedStyle(b.target,b.pseudoElement),m.width!==v||m.height!==h)&&(m=_[0],m.width=v,m.height=h,m=_[_.length-1],m.width=v,m.height=h,b.setKeyframes(_))}}r()},function(u){c.__reactViewTransition===f&&(c.__reactViewTransition=null);try{typeof u=="object"&&u!==null&&u.name==="InvalidStateError"&&(u.message==="View transition was skipped because document visibility state is hidden."||u.message==="Skipping view transition because document visibility state has become hidden."||u.message==="Skipping view transition because viewport size changed."||u.message==="Transition was aborted because of invalid state")&&(u=null),u!==null&&l(u)}finally{i(),a(),r()}}),f.finished.finally(function(){for(var u=0;u<p.length;u++)p[u].cancel();c.__reactViewTransition===f&&(c.__reactViewTransition=null),o()}),f}catch{return i(),a(),r(),null}}function ir(e,t){this._scope=document.documentElement,this._selector="::view-transition-"+e+"("+t+")"}ir.prototype.animate=function(e,t){return t=typeof t=="number"?{duration:t}:we({},t),t.pseudoElement=this._selector,this._scope.animate(e,t)};ir.prototype.getAnimations=function(){for(var e=this._scope,t=this._selector,n=e.getAnimations({subtree:!0}),i=[],a=0;a<n.length;a++){var s=n[a].effect;s!==null&&s.target===e&&s.pseudoElement===t&&i.push(n[a])}return i};ir.prototype.getComputedStyle=function(){return getComputedStyle(this._scope,this._selector)};function Xb(e){return{name:e,group:new ir("group",e),imagePair:new ir("image-pair",e),old:new ir("old",e),new:new ir("new",e)}}function di(e){this._fragmentFiber=e,this._observers=this._eventListeners=null}di.prototype.addEventListener=function(e,t,n){var i=null,a=null;if(!(n!=null&&typeof n!="boolean"&&(i=n.signal||null,i!==null&&i.aborted))){this._eventListeners===null&&(this._eventListeners=[]);var s=this._eventListeners;if(Wb(s,e,t,n)===-1){var r=this,o=t;n!=null&&typeof n!="boolean"&&n.once===!0&&(o=function(l){r.removeEventListener(e,t,n),typeof t=="function"?t.call(this,l):t.handleEvent(l)}),i!==null&&(a=r.removeEventListener.bind(r,e,t,n),i.addEventListener("abort",a,{once:!0}),a=i.removeEventListener.bind(i,"abort",a)),i=Ro(n),s.push({type:e,listener:t,optionsOrUseCapture:n,attachedListener:o,cleanup:a}),$n(this._fragmentFiber.child,!1,ow,e,o,i)}this._eventListeners=s}};function ow(e,t,n,i){return nn(e).addEventListener(t,n,i),!1}di.prototype.removeEventListener=function(e,t,n){var i=this._eventListeners;if(i!==null&&(t=Wb(i,e,t,n),t!==-1)){var a=i[t];n=a.attachedListener;var s=a.cleanup;a=Ro(a.optionsOrUseCapture),$n(this._fragmentFiber.child,!1,lw,e,n,a),i.splice(t,1),s!==null&&s()}};function lw(e,t,n,i){return nn(e).removeEventListener(t,n,i),!1}function Ro(e){return e!=null&&typeof e!="boolean"&&(e.once===!0||e.signal instanceof AbortSignal)?{capture:e.capture,passive:e.passive}:e}function Yy(e){return e==null?"c=0":typeof e=="boolean"?"c="+(e?"1":"0"):"c="+(e.capture?"1":"0")}function Wb(e,t,n,i){if(e.length===0)return-1;i=Yy(i);for(var a=0;a<e.length;a++){var s=e[a];if(s.type===t&&s.listener===n&&Yy(s.optionsOrUseCapture)===i)return a}return-1}di.prototype.dispatchEvent=function(e){var t=vr(this._fragmentFiber);if(t===null)return!0;t=nn(t);var n=this._eventListeners;if(n!==null&&0<n.length||!e.bubbles){var i=t.nodeType===9?t.createComment(""):document.createTextNode("");if(n)for(var a=0;a<n.length;a++){var s=n[a];i.addEventListener(s.type,s.attachedListener,Ro(s.optionsOrUseCapture))}if(t.appendChild(i),e=i.dispatchEvent(e),n)for(a=0;a<n.length;a++)s=n[a],i.removeEventListener(s.type,s.attachedListener,Ro(s.optionsOrUseCapture));return t.removeChild(i),e}return t.dispatchEvent(e)};di.prototype.focus=function(e){$n(this._fragmentFiber.child,!0,qb,e,void 0,void 0)};function qb(e,t){return e.tag===6?!1:(e=nn(e),xw(e,t))}di.prototype.focusLast=function(e){var t=[];$n(this._fragmentFiber.child,!0,p0,t,void 0,void 0);for(var n=t.length-1;0<=n&&!qb(t[n],e);n--);};function p0(e,t){return t.push(e),!1}di.prototype.blur=function(){var e=vr(this._fragmentFiber);e!==null&&(e=nn(e),e=dc(e).activeElement,e!==null&&$n(this._fragmentFiber.child,!1,cw,e,void 0,void 0))};function cw(e,t){return e.tag===6?!1:(e=nn(e),e===t||e.contains(t)?(t.blur(),!0):!1)}di.prototype.observeUsing=function(e){this._observers===null&&(this._observers=new Set),this._observers.add(e),$n(this._fragmentFiber.child,!1,uw,e,void 0,void 0)};function uw(e,t){return e.tag===6||(e=nn(e),t.observe(e)),!1}di.prototype.unobserveUsing=function(e){var t=this._observers;if(t!==null&&t.has(e)){t.delete(e),$n(this._fragmentFiber.child,!1,fw,e,void 0,void 0);for(var n=t=0;n<Bi.length;n++){var i=Bi[n];i.fragmentInstance===this&&i.observer===e?e.unobserve(i.instance):Bi[t++]=i}Bi.length=t}};function fw(e,t){return e.tag===6||(e=nn(e),t.unobserve(e)),!1}var Bi=[],cm=!1;function hw(e,t,n){Bi.push({fragmentInstance:e,observer:t,instance:n}),cm||(cm=!0,Sw(function(){cm=!1;var i=Bi;Bi=[];for(var a=0;a<i.length;a++){var s=i[a];s.observer.unobserve(s.instance)}}))}di.prototype.getClientRects=function(){var e=[];return $n(this._fragmentFiber.child,!1,dw,e,void 0,void 0),e};function dw(e,t){if(e.tag===6){e=e.stateNode;var n=e.ownerDocument.createRange();n.selectNodeContents(e),t.push.apply(t,n.getClientRects())}else e=nn(e),t.push.apply(t,e.getClientRects());return!1}di.prototype.getRootNode=function(e){var t=vr(this._fragmentFiber);return t===null?this:nn(t).getRootNode(e)};di.prototype.compareDocumentPosition=function(e){var t=vr(this._fragmentFiber);if(t===null)return Node.DOCUMENT_POSITION_DISCONNECTED;var n=[];$n(this._fragmentFiber.child,!1,p0,n,void 0,void 0);var i=nn(t);if(n.length===0){if(n=i,A_(this._fragmentFiber)){t:{for(t=this._fragmentFiber.return;t!==null;){if(t.tag===4){t=t.stateNode.containerInfo;break t}if(t.tag===3||t.tag===5||t.tag===27)break;t=t.return}t=null}t!=null&&(n=t)}t=this._fragmentFiber;var a=i=n.compareDocumentPosition(e);return n===e?a=Node.DOCUMENT_POSITION_CONTAINS:i&Node.DOCUMENT_POSITION_CONTAINED_BY&&(n=px(t)[1],n===null?a=Node.DOCUMENT_POSITION_PRECEDING:(e=nn(n).compareDocumentPosition(e),a=e===0||e&Node.DOCUMENT_POSITION_FOLLOWING?Node.DOCUMENT_POSITION_FOLLOWING:Node.DOCUMENT_POSITION_PRECEDING)),a|=Node.DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC}t=nn(n[0]),a=nn(n[n.length-1]);var s=A_(this._fragmentFiber)?t.parentElement:i;if(s==null)return Node.DOCUMENT_POSITION_DISCONNECTED;i=s.compareDocumentPosition(t)&Node.DOCUMENT_POSITION_CONTAINED_BY,s=s.compareDocumentPosition(a)&Node.DOCUMENT_POSITION_CONTAINED_BY;var r=t.compareDocumentPosition(e),o=a.compareDocumentPosition(e),l=r&Node.DOCUMENT_POSITION_CONTAINED_BY||o&Node.DOCUMENT_POSITION_CONTAINED_BY;return o=i&&s&&r&Node.DOCUMENT_POSITION_FOLLOWING&&o&Node.DOCUMENT_POSITION_PRECEDING,t=i&&t===e||s&&a===e||l||o?Node.DOCUMENT_POSITION_CONTAINED_BY:!i&&t===e||!s&&a===e?Node.DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC:r,t&Node.DOCUMENT_POSITION_DISCONNECTED||t&Node.DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC||pw(t,this._fragmentFiber,n[0],n[n.length-1],e)?t:Node.DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC};function pw(e,t,n,i,a){var s=nr(a);if(e&Node.DOCUMENT_POSITION_CONTAINED_BY){if(n=!!s)t:{for(;s!==null;){if(s.tag===7&&(s===t||s.alternate===t)){n=!0;break t}s=s.return}n=!1}return n}if(e&Node.DOCUMENT_POSITION_CONTAINS){if(s===null)return s=a.ownerDocument,a===s||a===s.documentElement||a===s.body;t:{for(s=t,t=vr(t);s!==null;){if(!(s.tag!==5&&s.tag!==3&&s.tag!==27||s!==t&&s.alternate!==t)){s=!0;break t}s=s.return}s=!1}return s}return e&Node.DOCUMENT_POSITION_PRECEDING?((t=!!s)&&!(t=s===n)&&(t=hm(n,s,w_),t===null?t=!1:($n(t,!0,XT,s,n),s=jr,jr=null,t=s!==null)),t):e&Node.DOCUMENT_POSITION_FOLLOWING?((t=!!s)&&!(t=s===i)&&(t=hm(i,s,w_),t===null?t=!1:($n(t,!0,WT,s,i),s=jr,fm=jr=null,t=s!==null)),t):!1}function Zy(e,t){var n=e.ownerDocument.createRange();n.selectNodeContents(e),e=n.getBoundingClientRect(),window.scrollTo(window.scrollX+e.left,t?window.scrollY+e.top:window.scrollY+e.bottom-window.innerHeight)}di.prototype.scrollIntoView=function(e){if(typeof e=="object")throw Error(J(566));var t=[];$n(this._fragmentFiber.child,!1,p0,t,void 0,void 0);var n=e!==!1;if(t.length===0){var i=px(this._fragmentFiber);if(i=n?i[1]||i[0]||vr(this._fragmentFiber):i[0]||i[1],i===null)return;if(i.tag===6){e=nn(i),Zy(e,n);return}if(i=nn(i),i.nodeType!==9){if(i.nodeType===11){n="host"in i?i.host:null,n!==null&&n.scrollIntoView(e);return}i.scrollIntoView(e)}}for(i=n?t.length-1:0;i!==(n?-1:t.length);){var a=t[i];a.tag===6?(a=nn(a),Zy(a,n)):nn(a).scrollIntoView(e),i+=n?-1:1}};function mw(e,t){return e=nn(e),Yb(e,t),!1}function Yb(e,t){e.reactFragments==null&&(e.reactFragments=new Set),e.reactFragments.add(t)}function Zb(e,t){var n=t._eventListeners;if(n!==null)for(var i=0;i<n.length;i++){var a=n[i];e.addEventListener(a.type,a.attachedListener,Ro(a.optionsOrUseCapture))}e.nodeType!==3&&(n=t._observers,n!==null&&n.forEach(function(s){for(var r=0,o=0;o<Bi.length;o++){var l=Bi[o];(l.fragmentInstance!==t||l.observer!==s||l.instance!==e)&&(Bi[r++]=l)}Bi.length=r,s.observe(e)}),Yb(e,t))}function gw(e,t){var n=t._eventListeners;if(n!==null)for(var i=0;i<n.length;i++){var a=n[i];e.removeEventListener(a.type,a.attachedListener,Ro(a.optionsOrUseCapture))}e.nodeType!==3&&(n=t._observers,n!==null&&n.forEach(function(s){typeof s.rootMargin=="string"?hw(t,s,e):s.unobserve(e)}),e.reactFragments!=null&&e.reactFragments.delete(t))}function gg(e){var t=e.firstChild;for(t&&t.nodeType===10&&(t=t.nextSibling);t;){var n=t;switch(t=t.nextSibling,n.nodeName){case"HTML":case"HEAD":case"BODY":gg(n),Zf(n);continue;case"SCRIPT":case"STYLE":continue;case"LINK":if(n.rel.toLowerCase()==="stylesheet")continue}e.removeChild(n)}}function vw(e,t,n,i){for(;e.nodeType===1;){var a=n;if(e.nodeName.toLowerCase()!==t.toLowerCase()){if(!i&&(e.nodeName!=="INPUT"||e.type!=="hidden"))break}else if(i){if(!e[bc])switch(t){case"meta":if(!e.hasAttribute("itemprop"))break;return e;case"link":if(s=e.getAttribute("rel"),s==="stylesheet"&&e.hasAttribute("data-precedence"))break;if(s!==a.rel||e.getAttribute("href")!==(a.href==null||a.href===""?null:a.href)||e.getAttribute("crossorigin")!==(a.crossOrigin==null?null:a.crossOrigin)||e.getAttribute("title")!==(a.title==null?null:a.title))break;return e;case"style":if(e.hasAttribute("data-precedence"))break;return e;case"script":if(s=e.getAttribute("src"),(s!==(a.src==null?null:a.src)||e.getAttribute("type")!==(a.type==null?null:a.type)||e.getAttribute("crossorigin")!==(a.crossOrigin==null?null:a.crossOrigin))&&s&&e.hasAttribute("async")&&!e.hasAttribute("itemprop"))break;return e;default:return e}}else if(t==="input"&&e.type==="hidden"){var s=a.name==null?null:""+a.name;if(a.type==="hidden"&&e.getAttribute("name")===s)return e}else return e;if(e=Ci(e.nextSibling),e===null)break}return null}function _w(e,t,n){if(t==="")return null;for(;e.nodeType!==3;)if((e.nodeType!==1||e.nodeName!=="INPUT"||e.type!=="hidden")&&!n||(e=Ci(e.nextSibling),e===null))return null;return e}function Jb(e,t){for(;e.nodeType!==8;)if((e.nodeType!==1||e.nodeName!=="INPUT"||e.type!=="hidden")&&!t||(e=Ci(e.nextSibling),e===null))return null;return e}function vg(e){return e.data==="$?"||e.data==="$~"}function m0(e){return e.data==="$!"||e.data==="$?"&&e.ownerDocument.readyState!=="loading"}function yw(e,t){var n=e.ownerDocument;if(e.data==="$~")e._reactRetry=t;else if(e.data!=="$?"||n.readyState!=="loading")t();else{var i=function(){t(),n.removeEventListener("DOMContentLoaded",i)};n.addEventListener("DOMContentLoaded",i),e._reactRetry=i}}function Ci(e){for(;e!=null;e=e.nextSibling){var t=e.nodeType;if(t===1||t===3)break;if(t===8){if(t=e.data,t==="$"||t==="$!"||t==="$?"||t==="$~"||t==="&"||t==="F!"||t==="F")break;if(t==="/$"||t==="/&")return null}}return e}var _g=null;function Jy(e){e=e.nextSibling;for(var t=0;e;){if(e.nodeType===8){var n=e.data;if(n==="/$"||n==="/&"){if(t===0)return Ci(e.nextSibling);t--}else n!=="$"&&n!=="$!"&&n!=="$?"&&n!=="$~"&&n!=="&"||t++}e=e.nextSibling}return null}function Ky(e){e=e.previousSibling;for(var t=0;e;){if(e.nodeType===8){var n=e.data;if(n==="$"||n==="$!"||n==="$?"||n==="$~"||n==="&"){if(t===0)return e;t--}else n!=="/$"&&n!=="/&"||t++}e=e.previousSibling}return null}function xw(e,t){function n(){i=!0}if(e.ownerDocument.activeElement===e)return!0;var i=!1;try{e.ownerDocument.addEventListener("focus",n,!0),(e.focus||HTMLElement.prototype.focus).call(e,t)}finally{e.ownerDocument.removeEventListener("focus",n,!0)}return i}function Sw(e){Xy(function(){Xy(function(t){return e(t)})})}function Kb(e,t,n){switch(t=dc(n),e){case"html":if(e=t.documentElement,!e)throw Error(J(452));return e;case"head":if(e=t.head,!e)throw Error(J(453));return e;case"body":if(e=t.body,!e)throw Error(J(454));return e;default:throw Error(J(451))}}function Qb(e,t,n){for(var i in n){var a=n[i];n.hasOwnProperty(i)&&a!=null&&_e(e,t,i,null,KA,a)}n.dangerouslySetInnerHTML!=null&&(e.textContent=""),e.onclick===na&&(e.onclick=null),Zf(e)}function um(e){for(var t=e.attributes;t.length;)e.removeAttributeNode(t[0]);Zf(e)}var Ri=new Map,Qy=new Set;function pc(e){if(typeof e.getRootNode=="function"){var t=e.getRootNode();if(t.nodeType===9||t.nodeType===11)return t}return e.nodeType===9?e:e.ownerDocument}var za=he.d;he.d={f:bw,r:Mw,D:Tw,C:Ew,L:Aw,m:ww,X:Rw,S:Cw,M:Dw};function bw(){var e=za.f(),t=lh();return e||t}function Mw(e){var t=Oo(e);t!==null&&t.tag===5&&t.type==="form"?OS(t):za.r(e)}var Fo=typeof document>"u"?null:document;function jb(e,t,n){var i=Fo;if(i&&typeof t=="string"&&t){var a=Ei(t);a='link[rel="'+e+'"][href="'+a+'"]',typeof n=="string"&&(a+='[crossorigin="'+n+'"]'),Qy.has(a)||(Qy.add(a),e={rel:e,crossOrigin:n,href:t},i.querySelector(a)===null&&(t=i.createElement("link"),Mn(t,"link",e),fn(t),i.head.appendChild(t)))}}function Tw(e){za.D(e),jb("dns-prefetch",e,null)}function Ew(e,t){za.C(e,t),jb("preconnect",e,t)}function Aw(e,t,n){za.L(e,t,n);var i=Fo;if(i&&e&&t){var a='link[rel="preload"][as="'+Ei(t)+'"]';t==="image"&&n&&n.imageSrcSet?(a+='[imagesrcset="'+Ei(n.imageSrcSet)+'"]',typeof n.imageSizes=="string"&&(a+='[imagesizes="'+Ei(n.imageSizes)+'"]')):a+='[href="'+Ei(e)+'"]';var s=a;switch(t){case"style":s=Do(e);break;case"script":s=Vo(e)}if(!(Ri.has(s)||(e=we({rel:"preload",href:t==="image"&&n&&n.imageSrcSet?void 0:e,as:t},n),Ri.set(s,e),i.querySelector(a)!==null||t==="style"&&i.querySelector(Cc(s))||t==="script"&&i.querySelector(Rc(s))))){var r=i.createElement("link");Mn(r,"link",e),t==="style"&&(r[bf]=!0,r.onload=r.onerror=function(){wx(r)}),fn(r),i.head.appendChild(r)}}}function ww(e,t){za.m(e,t);var n=Fo;if(n&&e){var i=t&&typeof t.as=="string"?t.as:"script",a='link[rel="modulepreload"][as="'+Ei(i)+'"][href="'+Ei(e)+'"]',s=a;switch(i){case"audioworklet":case"paintworklet":case"serviceworker":case"sharedworker":case"worker":case"script":s=Vo(e)}if(!Ri.has(s)&&(e=we({rel:"modulepreload",href:e},t),Ri.set(s,e),n.querySelector(a)===null)){switch(i){case"audioworklet":case"paintworklet":case"serviceworker":case"sharedworker":case"worker":case"script":if(n.querySelector(Rc(s)))return}i=n.createElement("link"),Mn(i,"link",e),fn(i),n.head.appendChild(i)}}}function Cw(e,t,n){za.S(e,t,n);var i=Fo;if(i&&e){var a=uo(i).hoistableStyles,s=Do(e);t=t||"default";var r=a.get(s);if(!r){var o={loading:0,preload:null};if(r=i.querySelector(Cc(s)))o.loading=5;else{e=we({rel:"stylesheet",href:e,"data-precedence":t},n),(n=Ri.get(s))&&g0(e,n);var l=r=i.createElement("link");fn(l),Mn(l,"link",e),l._p=new Promise(function(c,f){l.onload=c,l.onerror=f}),l.addEventListener("load",function(){o.loading|=1}),l.addEventListener("error",function(){o.loading|=2}),o.loading|=4,df(r,t,i)}r={type:"stylesheet",instance:r,count:1,state:o},a.set(s,r)}}}function Rw(e,t){za.X(e,t);var n=Fo;if(n&&e){var i=uo(n).hoistableScripts,a=Vo(e),s=i.get(a);s||(s=n.querySelector(Rc(a)),s||(e=we({src:e,async:!0},t),(t=Ri.get(a))&&v0(e,t),s=n.createElement("script"),fn(s),Mn(s,"link",e),n.head.appendChild(s)),s={type:"script",instance:s,count:1,state:null},i.set(a,s))}}function Dw(e,t){za.M(e,t);var n=Fo;if(n&&e){var i=uo(n).hoistableScripts,a=Vo(e),s=i.get(a);s||(s=n.querySelector(Rc(a)),s||(e=we({src:e,async:!0,type:"module"},t),(t=Ri.get(a))&&v0(e,t),s=n.createElement("script"),fn(s),Mn(s,"link",e),n.head.appendChild(s)),s={type:"script",instance:s,count:1,state:null},i.set(a,s))}}function jy(e,t,n,i){var a=(a=ps.current)?pc(a):null;if(!a)throw Error(J(446));switch(e){case"meta":case"title":return null;case"style":return typeof n.precedence=="string"&&typeof n.href=="string"?(n=Do(n.href),t=uo(a).hoistableStyles,i=t.get(n),i||(i={type:"style",instance:null,count:0,state:null},t.set(n,i)),i):{type:"void",instance:null,count:0,state:null};case"link":if(n.rel==="stylesheet"&&typeof n.href=="string"&&typeof n.precedence=="string"){e=Do(n.href);var s=uo(a).hoistableStyles,r=s.get(e);if(r||(a=a.ownerDocument||a,r={type:"stylesheet",instance:null,count:0,state:{loading:0,preload:null}},s.set(e,r),(s=a.querySelector(Cc(e)))?s._p||(r.instance=s,r.state.loading=5):(s=Ri.get(e),s||(s={rel:"preload",as:"style",href:n.href,crossOrigin:n.crossOrigin,integrity:n.integrity,media:n.media,hrefLang:n.hrefLang,referrerPolicy:n.referrerPolicy},Ri.set(e,s)),Nw(a,e,s,r.state))),t&&i===null)throw Error(J(528,""));return r}if(t&&i!==null)throw Error(J(529,""));return null;case"script":return t=n.async,n=n.src,typeof n=="string"&&t&&typeof t!="function"&&typeof t!="symbol"?(n=Vo(n),t=uo(a).hoistableScripts,i=t.get(n),i||(i={type:"script",instance:null,count:0,state:null},t.set(n,i)),i):{type:"void",instance:null,count:0,state:null};default:throw Error(J(444,e))}}function Do(e){return'href="'+Ei(e)+'"'}function Cc(e){return'link[rel="stylesheet"]['+e+"]"}function $b(e){return we({},e,{"data-precedence":e.precedence,precedence:null})}function Nw(e,t,n,i){if(t=e.querySelector('link[rel="preload"][as="style"]['+t+"]")){if(t[bf]!==!0){i.loading=1;return}}else t=e.createElement("link"),t[bf]=!0,t.onload=t.onerror=wx.bind(null,t),Mn(t,"link",n),fn(t),e.head.appendChild(t);i.preload=t,t.addEventListener("load",function(){return i.loading|=1}),t.addEventListener("error",function(){return i.loading|=2})}function Vo(e){return'[src="'+Ei(e)+'"]'}function Rc(e){return"script[async]"+e}function $y(e,t,n){if(t.count++,t.instance===null)switch(t.type){case"style":var i=e.querySelector('style[data-href~="'+Ei(n.href)+'"]');if(i)return t.instance=i,fn(i),i;var a=we({},n,{"data-href":n.href,"data-precedence":n.precedence,href:null,precedence:null});return i=(e.ownerDocument||e).createElement("style"),fn(i),Mn(i,"style",a),df(i,n.precedence,e),t.instance=i;case"stylesheet":a=Do(n.href);var s=e.querySelector(Cc(a));if(s)return t.state.loading|=4,t.instance=s,fn(s),s;i=$b(n),(a=Ri.get(a))&&g0(i,a),s=(e.ownerDocument||e).createElement("link"),fn(s);var r=s;return r._p=new Promise(function(o,l){r.onload=o,r.onerror=l}),Mn(s,"link",i),t.state.loading|=4,df(s,n.precedence,e),t.instance=s;case"script":return s=Vo(n.src),(a=e.querySelector(Rc(s)))?(t.instance=a,fn(a),a):(i=n,(a=Ri.get(s))&&(i=we({},n),v0(i,a)),e=e.ownerDocument||e,a=e.createElement("script"),fn(a),Mn(a,"link",i),e.head.appendChild(a),t.instance=a);case"void":return null;default:throw Error(J(443,t.type))}else t.type==="stylesheet"&&(t.state.loading&4)===0&&(i=t.instance,t.state.loading|=4,df(i,n.precedence,e));return t.instance}function df(e,t,n){for(var i=n.querySelectorAll('link[rel="stylesheet"][data-precedence],style[data-precedence]'),a=i.length?i[i.length-1]:null,s=a,r=0;r<i.length;r++){var o=i[r];if(o.dataset.precedence===t)s=o;else if(s!==a)break}s?s.parentNode.insertBefore(e,s.nextSibling):(t=n.nodeType===9?n.head:n,t.insertBefore(e,t.firstChild))}function g0(e,t){e.crossOrigin==null&&(e.crossOrigin=t.crossOrigin),e.referrerPolicy==null&&(e.referrerPolicy=t.referrerPolicy),e.title==null&&(e.title=t.title)}function v0(e,t){e.crossOrigin==null&&(e.crossOrigin=t.crossOrigin),e.referrerPolicy==null&&(e.referrerPolicy=t.referrerPolicy),e.integrity==null&&(e.integrity=t.integrity)}var pf=null;function tx(e,t,n){if(pf===null){var i=new Map,a=pf=new Map;a.set(n,i)}else a=pf,i=a.get(n),i||(i=new Map,a.set(n,i));if(i.has(e))return i;for(i.set(e,null),n=n.getElementsByTagName(e),a=0;a<n.length;a++){var s=n[a];if(!(s[bc]||s[yn]||e==="link"&&s.getAttribute("rel")==="stylesheet")&&s.namespaceURI!=="http://www.w3.org/2000/svg"){var r=s.getAttribute(t)||"";r=e+r;var o=i.get(r);o?o.push(s):i.set(r,[s])}}return i}function yg(e,t,n){e=e.ownerDocument||e,e.head.insertBefore(n,t==="title"?e.querySelector("head > title"):null)}function Uw(e,t,n){if(n===1||t.itemProp!=null)return!1;switch(e){case"meta":case"title":return!0;case"style":if(typeof t.precedence!="string"||typeof t.href!="string"||t.href==="")break;return!0;case"link":if(typeof t.rel!="string"||typeof t.href!="string"||t.href===""||t.onLoad||t.onError)break;return t.rel==="stylesheet"?(e=t.disabled,typeof t.precedence=="string"&&e==null):!0;case"script":if(t.async&&typeof t.async!="function"&&typeof t.async!="symbol"&&!t.onLoad&&!t.onError&&t.src&&typeof t.src=="string")return!0}return!1}function ex(e,t){return e==="img"&&t.src!=null&&t.src!==""&&t.onLoad==null&&t.loading!=="lazy"}function tM(e){return!(e.type==="stylesheet"&&(e.state.loading&3)===0)}function eM(e){return(e.width||100)*(e.height||100)*(typeof devicePixelRatio=="number"?devicePixelRatio:1)*.25}function nx(e,t){typeof t.decode=="function"&&(e.imgCount++,t.complete||(e.imgBytes+=eM(t),e.suspenseyImages.push(t)),e=Ow.bind(e),t.decode().then(e,e))}function Lw(e,t,n,i){if(n.type==="stylesheet"&&(typeof i.media!="string"||matchMedia(i.media).matches!==!1)&&(n.state.loading&4)===0){if(n.instance===null){var a=Do(i.href),s=t.querySelector(Cc(a));if(s){t=s._p,t!==null&&typeof t=="object"&&typeof t.then=="function"&&(e.count++,e=mc.bind(e),t.then(e,e)),n.state.loading|=4,n.instance=s,fn(s);return}s=t.ownerDocument||t,i=$b(i),(a=Ri.get(a))&&g0(i,a),s=s.createElement("link"),fn(s);var r=s;r._p=new Promise(function(o,l){r.onload=o,r.onerror=l}),Mn(s,"link",i),n.instance=s}e.stylesheets===null&&(e.stylesheets=new Map),e.stylesheets.set(n,t),(t=n.state.preload)&&(n.state.loading&3)===0&&(e.count++,n=mc.bind(e),t.addEventListener("load",n),t.addEventListener("error",n))}}var mf=0;function Iw(e,t){return e.stylesheets&&e.count===0&&gf(e,e.stylesheets),0<e.count||0<e.imgCount?function(n){var i=setTimeout(function(){if(e.stylesheets&&gf(e,e.stylesheets),e.unsuspend){var s=e.unsuspend;e.unsuspend=null,s()}},6e4+t);0<e.imgBytes&&mf===0&&(mf=62500*jA());var a=setTimeout(function(){if(e.waitingForImages=!1,e.count===0&&(e.stylesheets&&gf(e,e.stylesheets),e.unsuspend)){var s=e.unsuspend;e.unsuspend=null,s()}},(e.imgBytes>mf?50:800)+t);return e.unsuspend=n,function(){e.unsuspend=null,clearTimeout(i),clearTimeout(a)}}:null}function nM(e){if(e.count===0&&(e.imgCount===0||!e.waitingForImages)){if(e.stylesheets)gf(e,e.stylesheets);else if(e.unsuspend){var t=e.unsuspend;e.unsuspend=null,t()}}}function mc(){this.count--,nM(this)}function Ow(){this.imgCount--,nM(this)}var Wf=null;function gf(e,t){e.stylesheets=null,e.unsuspend!==null&&(e.count++,Wf=new Map,t.forEach(Pw,e),Wf=null,mc.call(e))}function Pw(e,t){if(!(t.state.loading&4)){var n=Wf.get(e);if(n)var i=n.get(null);else{n=new Map,Wf.set(e,n);for(var a=e.querySelectorAll("link[data-precedence],style[data-precedence]"),s=0;s<a.length;s++){var r=a[s];(r.nodeName==="LINK"||r.getAttribute("media")!=="not all")&&(n.set(r.dataset.precedence,r),i=r)}i&&n.set(null,i)}a=t.instance,r=a.getAttribute("data-precedence"),s=n.get(r)||i,s===i&&n.set(null,a),n.set(r,a),this.count++,i=mc.bind(this),a.addEventListener("load",i),a.addEventListener("error",i),s?s.parentNode.insertBefore(a,s.nextSibling):(e=e.nodeType===9?e.head:e,e.insertBefore(a,e.firstChild)),t.state.loading|=4}}var No={$$typeof:ea,Provider:null,Consumer:null,_currentValue:ar,_currentValue2:ar,_threadCount:0};function Bw(e,t,n,i,a,s,r,o,l){this.tag=1,this.containerInfo=e,this.pingCache=this.current=this.pendingChildren=null,this.timeoutHandle=-1,this.callbackNode=this.next=this.pendingContext=this.context=this.cancelPendingCommit=null,this.callbackPriority=0,this.expirationTimes=zp(-1),this.entangledLanes=this.shellSuspendCounter=this.errorRecoveryDisabledLanes=this.expiredLanes=this.warmLanes=this.pingedLanes=this.suspendedLanes=this.pendingLanes=0,this.entanglements=zp(0),this.hiddenUpdates=zp(null),this.identifierPrefix=i,this.onUncaughtError=a,this.onCaughtError=s,this.onRecoverableError=r,this.pooledCache=null,this.pooledCacheLanes=0,this.formState=l,this.transitionTypes=null,this.incompleteTransitions=new Map}function iM(e,t,n,i,a,s,r,o,l,c,f,p){return e=new Bw(e,t,n,r,l,c,f,p,o),t=1,s===!0&&(t|=24),s=Qn(3,null,null,t),e.current=s,s.stateNode=e,t=Bg(),t.refCount++,e.pooledCache=t,t.refCount++,s.memoizedState={element:i,isDehydrated:n,cache:t},Vg(s),e}function aM(e){return e?(e=ro,e):ro}function sM(e,t,n,i,a,s){a=aM(a),i.context===null?i.context=a:i.pendingContext=a,i=gs(t),i.payload={element:n},s=s===void 0?null:s,s!==null&&(i.callback=s),n=vs(e,i,t),n!==null&&(jn(n,e,t),Yl(n,e,t))}function ix(e,t){if(e=e.memoizedState,e!==null&&e.dehydrated!==null){var n=e.retryLane;e.retryLane=n!==0&&n<t?n:t}}function _0(e,t){ix(e,t),(e=e.alternate)&&ix(e,t)}function rM(e){if(e.tag===13||e.tag===31){var t=xr(e,67108864);t!==null&&jn(t,e,67108864),_0(e,67108864)}}function ax(e){if(e.tag===13||e.tag===31){var t=fi();t=Eg(t);var n=xr(e,t);n!==null&&jn(n,e,t),_0(e,t)}}var Uo=!0;function zw(e,t,n,i){var a=Bt.T;Bt.T=null;var s=he.p;try{he.p=2,y0(e,t,n,i)}finally{he.p=s,Bt.T=a}}function Fw(e,t,n,i){var a=Bt.T;Bt.T=null;var s=he.p;try{he.p=8,y0(e,t,n,i)}finally{he.p=s,Bt.T=a}}function y0(e,t,n,i){if(Uo){var a=xg(i);if(a===null)om(e,t,i,qf,n),sx(e,i);else if(Hw(a,e,t,n,i))i.stopPropagation();else if(sx(e,i),t&4&&-1<Vw.indexOf(e)){for(;a!==null;){var s=Oo(a);if(s!==null)switch(s.tag){case 3:if(s=s.stateNode,s.current.memoizedState.isDehydrated){var r=$s(s.pendingLanes);if(r!==0){var o=s;for(o.pendingLanes|=2,o.entangledLanes|=2;r;){var l=1<<31-ui(r);o.entanglements[1]|=l,r&=~l}ca(s),(fe&6)===0&&(Vf=li()+500,wc(0,!1))}}break;case 31:case 13:o=xr(s,2),o!==null&&jn(o,s,2),lh(),_0(s,2)}if(s=xg(i),s===null&&om(e,t,i,qf,n),s===a)break;a=s}a!==null&&i.stopPropagation()}else om(e,t,i,null,n)}}function xg(e){return e=Cg(e),x0(e)}var qf=null;function x0(e){if(qf=null,e=nr(e),e!==null){var t=_c(e);if(t===null)e=null;else{var n=t.tag;if(n===13){if(e=fx(t),e!==null)return e;e=null}else if(n===31){if(e=hx(t),e!==null)return e;e=null}else if(n===3){if(t.stateNode.current.memoizedState.isDehydrated)return t.tag===3?t.stateNode.containerInfo:null;e=null}else t!==e&&(e=null)}}return qf=e,null}function oM(e){switch(e){case"beforetoggle":case"cancel":case"click":case"close":case"contextmenu":case"copy":case"cut":case"auxclick":case"dblclick":case"dragend":case"dragstart":case"drop":case"focusin":case"focusout":case"input":case"invalid":case"keydown":case"keypress":case"keyup":case"mousedown":case"mouseup":case"paste":case"pause":case"play":case"pointercancel":case"pointerdown":case"pointerup":case"ratechange":case"reset":case"seeked":case"submit":case"toggle":case"touchcancel":case"touchend":case"touchstart":case"volumechange":case"change":case"selectionchange":case"textInput":case"compositionstart":case"compositionend":case"compositionupdate":case"beforeblur":case"afterblur":case"beforeinput":case"blur":case"fullscreenchange":case"fullscreenerror":case"focus":case"hashchange":case"popstate":case"select":case"selectstart":return 2;case"drag":case"dragenter":case"dragexit":case"dragleave":case"dragover":case"mousemove":case"mouseout":case"mouseover":case"pointermove":case"pointerout":case"pointerover":case"resize":case"scroll":case"touchmove":case"wheel":case"mouseenter":case"mouseleave":case"pointerenter":case"pointerleave":return 8;case"message":switch(tE()){case _x:return 2;case yx:return 8;case Sf:case eE:return 32;case xx:return 268435456;default:return 32}default:return 32}}var Sg=!1,Ss=null,bs=null,Ms=null,gc=new Map,vc=new Map,os=[],Vw="mousedown mouseup touchcancel touchend touchstart auxclick dblclick pointercancel pointerdown pointerup dragend dragstart drop compositionend compositionstart keydown keypress keyup input textInput copy cut paste click change contextmenu reset".split(" ");function sx(e,t){switch(e){case"focusin":case"focusout":Ss=null;break;case"dragenter":case"dragleave":bs=null;break;case"mouseover":case"mouseout":Ms=null;break;case"pointerover":case"pointerout":gc.delete(t.pointerId);break;case"gotpointercapture":case"lostpointercapture":vc.delete(t.pointerId)}}function Pl(e,t,n,i,a,s){return e===null||e.nativeEvent!==s?(e={blockedOn:t,domEventName:n,eventSystemFlags:i,nativeEvent:s,targetContainers:[a]},t!==null&&(t=Oo(t),t!==null&&rM(t)),e):(e.eventSystemFlags|=i,t=e.targetContainers,a!==null&&t.indexOf(a)===-1&&t.push(a),e)}function Hw(e,t,n,i,a){switch(t){case"focusin":return Ss=Pl(Ss,e,t,n,i,a),!0;case"dragenter":return bs=Pl(bs,e,t,n,i,a),!0;case"mouseover":return Ms=Pl(Ms,e,t,n,i,a),!0;case"pointerover":var s=a.pointerId;return gc.set(s,Pl(gc.get(s)||null,e,t,n,i,a)),!0;case"gotpointercapture":return s=a.pointerId,vc.set(s,Pl(vc.get(s)||null,e,t,n,i,a)),!0}return!1}function lM(e){var t=nr(e.target);if(t!==null){var n=_c(t);if(n!==null){if(t=n.tag,t===13){if(t=fx(n),t!==null){e.blockedOn=t,N_(e.priority,function(){ax(n)});return}}else if(t===31){if(t=hx(n),t!==null){e.blockedOn=t,N_(e.priority,function(){ax(n)});return}}else if(t===3&&n.stateNode.current.memoizedState.isDehydrated){e.blockedOn=n.tag===3?n.stateNode.containerInfo:null;return}}}e.blockedOn=null}function vf(e){if(e.blockedOn!==null)return!1;for(var t=e.targetContainers;0<t.length;){var n=xg(e.nativeEvent);if(n===null){n=e.nativeEvent;var i=new n.constructor(n.type,n);Tm=i,n.target.dispatchEvent(i),Tm=null}else return t=Oo(n),t!==null&&rM(t),e.blockedOn=n,!1;t.shift()}return!0}function rx(e,t,n){vf(e)&&n.delete(t)}function Gw(){Sg=!1,Ss!==null&&vf(Ss)&&(Ss=null),bs!==null&&vf(bs)&&(bs=null),Ms!==null&&vf(Ms)&&(Ms=null),gc.forEach(rx),vc.forEach(rx)}function Ku(e,t){e.blockedOn===t&&(e.blockedOn=null,Sg||(Sg=!0,an.unstable_scheduleCallback(an.unstable_NormalPriority,Gw)))}var Qu=null;function ox(e){Qu!==e&&(Qu=e,an.unstable_scheduleCallback(an.unstable_NormalPriority,function(){Qu===e&&(Qu=null);for(var t=0;t<e.length;t+=3){var n=e[t],i=e[t+1],a=e[t+2];if(typeof i!="function"){if(x0(i||n)===null)continue;break}var s=Oo(n);s!==null&&(e.splice(t,3),t-=3,Vm(s,{pending:!0,data:a,method:n.method,action:i},i,a))}}))}function Lo(e){function t(l){return Ku(l,e)}Ss!==null&&Ku(Ss,e),bs!==null&&Ku(bs,e),Ms!==null&&Ku(Ms,e),gc.forEach(t),vc.forEach(t);for(var n=0;n<os.length;n++){var i=os[n];i.blockedOn===e&&(i.blockedOn=null)}for(;0<os.length&&(n=os[0],n.blockedOn===null);)lM(n),n.blockedOn===null&&os.shift();if(n=(e.ownerDocument||e).$$reactFormReplay,n!=null)for(i=0;i<n.length;i+=3){var a=n[i],s=n[i+1],r=a[ti]||null;if(typeof s=="function")r||ox(n);else if(r){var o=null;if(s&&s.hasAttribute("formAction")){if(a=s,r=s[ti]||null)o=r.formAction;else if(x0(a)!==null)continue}else o=r.action;typeof o=="function"?n[i+1]=o:(n.splice(i,3),i-=3),ox(n)}}}function cM(){function e(s){s.canIntercept&&s.info==="react-transition"&&s.intercept({handler:function(){return new Promise(function(r){return a=r})},focusReset:"manual",scroll:"manual"})}function t(){a!==null&&(a(),a=null),i||setTimeout(n,20)}function n(){if(!i&&!navigation.transition){var s=navigation.currentEntry;s&&s.url!=null&&navigation.navigate(s.url,{state:s.getState(),info:"react-transition",history:"replace"})}}if(typeof navigation=="object"){var i=!1,a=null;return navigation.addEventListener("navigate",e),navigation.addEventListener("navigatesuccess",t),navigation.addEventListener("navigateerror",t),setTimeout(n,100),function(){i=!0,navigation.removeEventListener("navigate",e),navigation.removeEventListener("navigatesuccess",t),navigation.removeEventListener("navigateerror",t),a!==null&&(a(),a=null)}}}function S0(e){this._internalRoot=e}fh.prototype.render=S0.prototype.render=function(e){var t=this._internalRoot;if(t===null)throw Error(J(409));var n=t.current,i=fi();sM(n,i,e,t,null,null)};fh.prototype.unmount=S0.prototype.unmount=function(){var e=this._internalRoot;if(e!==null){this._internalRoot=null;var t=e.containerInfo;sM(e.current,2,null,e,null,null),lh(),t[Io]=null}};function fh(e){this._internalRoot=e}fh.prototype.unstable_scheduleHydration=function(e){if(e){var t=Ax();e={blockedOn:null,target:e,priority:t};for(var n=0;n<os.length&&t!==0&&t<os[n].priority;n++);os.splice(n,0,e),n===0&&lM(e)}};var lx=cx.version;if(lx!=="19.3.1")throw Error(J(527,lx,"19.3.1"));he.findDOMNode=function(e){var t=e._reactInternals;if(t===void 0)throw typeof e.render=="function"?Error(J(188)):(e=Object.keys(e).join(","),Error(J(268,e)));return e=kT(t),e=e!==null?dx(e):null,e=e===null?null:e.stateNode,e};var kw={bundleType:0,version:"19.3.1",rendererPackageName:"react-dom",currentDispatcherRef:Bt,reconcilerVersion:"19.3.1"};if(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__<"u"&&(Bl=__REACT_DEVTOOLS_GLOBAL_HOOK__,!Bl.isDisabled&&Bl.supportsFiber))try{yc=Bl.inject(kw),ci=Bl}catch{}var Bl;hh.createRoot=function(e,t){if(!ux(e))throw Error(J(299));var n=!1,i="",a=kS,s=XS,r=WS;return t!=null&&(t.unstable_strictMode===!0&&(n=!0),t.identifierPrefix!==void 0&&(i=t.identifierPrefix),t.onUncaughtError!==void 0&&(a=t.onUncaughtError),t.onCaughtError!==void 0&&(s=t.onCaughtError),t.onRecoverableError!==void 0&&(r=t.onRecoverableError)),t=iM(e,1,!1,null,null,n,i,null,a,s,r,cM),e[Io]=t.current,h0(e),new S0(t)};hh.hydrateRoot=function(e,t,n){if(!ux(e))throw Error(J(299));var i=!1,a="",s=kS,r=XS,o=WS,l=null;return n!=null&&(n.unstable_strictMode===!0&&(i=!0),n.identifierPrefix!==void 0&&(a=n.identifierPrefix),n.onUncaughtError!==void 0&&(s=n.onUncaughtError),n.onCaughtError!==void 0&&(r=n.onCaughtError),n.onRecoverableError!==void 0&&(o=n.onRecoverableError),n.formState!==void 0&&(l=n.formState)),t=iM(e,1,!0,t,n??null,i,a,l,s,r,o,cM),t.context=aM(null),n=t.current,i=fi(),i=Eg(i),a=gs(i),a.callback=null,vs(n,a,i),n=i,t.current.lanes=n,Sc(t,n),ca(t),e[Io]=t.current,h0(e),new fh(t)};hh.version="19.3.1"});var dM=Zi((GN,hM)=>{"use strict";function fM(){if(!(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__>"u"||typeof __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE!="function"))try{__REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE(fM)}catch(e){console.error(e)}}fM(),hM.exports=uM()});var cT=Zi(vp=>{"use strict";var RN=Symbol.for("react.transitional.element"),DN=Symbol.for("react.fragment");function lT(e,t,n){var i=null;if(n!==void 0&&(i=""+n),t.key!==void 0&&(i=""+t.key),"key"in t){n={};for(var a in t)a!=="key"&&(n[a]=t[a])}else n=t;return t=n.ref,{$$typeof:RN,type:e,key:i,ref:t!==void 0?t:null,props:n}}vp.Fragment=DN;vp.jsx=lT;vp.jsxs=lT});var Xv=Zi((vO,uT)=>{"use strict";uT.exports=cT()});var Cn=Al(wl(),1),dT=Al(dM(),1);var ft=Al(wl(),1);var VM=0,tv=1,HM=2;var fu=1,GM=2,dl=3,Wa=0,wn=1,ga=2,va=0,wr=1,ev=2,nv=3,iv=4,kM=5;var zs=100,XM=101,WM=102,qM=103,YM=104,ZM=200,JM=201,KM=202,QM=203,Bh=204,zh=205,jM=206,$M=207,t1=208,e1=209,n1=210,i1=211,a1=212,s1=213,r1=214,Fh=0,Vh=1,Hh=2,Cr=3,Gh=4,kh=5,Xh=6,Wh=7,av=0,o1=1,l1=2,Wi=0,sv=1,rv=2,ov=3,hu=4,lv=5,cv=6,uv=7;var fv=300,Xs=301,Ur=302,_d=303,yd=304,du=306,qh=1e3,ha=1001,Yh=1002,pn=1003,c1=1004;var pu=1005;var En=1006,xd=1007;var Ws=1008;var ei=1009,hv=1010,dv=1011,pl=1012,Sd=1013,qi=1014,Ni=1015,_a=1016,bd=1017,Md=1018,ml=1020,pv=35902,mv=35899,gv=1021,vv=1022,Ui=1023,da=1026,qs=1027,Td=1028,Ed=1029,Lr=1030,Ad=1031;var wd=1033,mu=33776,gu=33777,vu=33778,_u=33779,Cd=35840,Rd=35841,Dd=35842,Nd=35843,Ud=36196,Ld=37492,Id=37496,Od=37488,Pd=37489,Bd=37490,zd=37491,Fd=37808,Vd=37809,Hd=37810,Gd=37811,kd=37812,Xd=37813,Wd=37814,qd=37815,Yd=37816,Zd=37817,Jd=37818,Kd=37819,Qd=37820,jd=37821,$d=36492,tp=36494,ep=36495,np=36283,ip=36284,ap=36285,sp=36286;var Fc=2300,Zh=2301,Ph=2302,W0=2303,q0=2400,Y0=2401,Z0=2402;var u1=3200;var _v=0,f1=1,Ja="",On="srgb",Rr="srgb-linear",Vc="linear",me="srgb";var Ar=7680;var J0=519,h1=512,d1=513,p1=514,rp=515,m1=516,g1=517,op=518,v1=519,K0=35044;var yv="300 es",ki=2e3,nl=2001;function Xw(e){for(let t=e.length-1;t>=0;--t)if(e[t]>=65535)return!0;return!1}function Ww(e){return ArrayBuffer.isView(e)&&!(e instanceof DataView)}function Hc(e){return document.createElementNS("http://www.w3.org/1999/xhtml",e)}function _1(){let e=Hc("canvas");return e.style.display="block",e}var pM={},il=null;function xv(...e){let t="THREE."+e.shift();il?il("log",t,...e):console.log(t,...e)}function y1(e){let t=e[0];if(typeof t=="string"&&t.startsWith("TSL:")){let n=e[1];n&&n.isStackTrace?e[0]+=" "+n.getLocation():e[1]='Stack trace not available. Enable "THREE.Node.captureStackTrace" to capture stack traces.'}return e}function Ot(...e){e=y1(e);let t="THREE."+e.shift();if(il)il("warn",t,...e);else{let n=e[0];n&&n.isStackTrace?console.warn(n.getError(t)):console.warn(t,...e)}}function zt(...e){e=y1(e);let t="THREE."+e.shift();if(il)il("error",t,...e);else{let n=e[0];n&&n.isStackTrace?console.error(n.getError(t)):console.error(t,...e)}}function Gc(...e){let t=e.join(" ");t in pM||(pM[t]=!0,Ot(...e))}function x1(e,t,n){return new Promise(function(i,a){function s(){switch(e.clientWaitSync(t,e.SYNC_FLUSH_COMMANDS_BIT,0)){case e.WAIT_FAILED:a();break;case e.TIMEOUT_EXPIRED:setTimeout(s,n);break;default:i()}}setTimeout(s,n)})}var S1={[Fh]:Vh,[Hh]:Xh,[Gh]:Wh,[Cr]:kh,[Vh]:Fh,[Xh]:Hh,[Wh]:Gh,[kh]:Cr},qa=class{addEventListener(t,n){this._listeners===void 0&&(this._listeners={});let i=this._listeners;i[t]===void 0&&(i[t]=[]),i[t].indexOf(n)===-1&&i[t].push(n)}hasEventListener(t,n){let i=this._listeners;return i===void 0?!1:i[t]!==void 0&&i[t].indexOf(n)!==-1}removeEventListener(t,n){let i=this._listeners;if(i===void 0)return;let a=i[t];if(a!==void 0){let s=a.indexOf(n);s!==-1&&a.splice(s,1)}}dispatchEvent(t){let n=this._listeners;if(n===void 0)return;let i=n[t.type];if(i!==void 0){t.target=this;let a=i.slice(0);for(let s=0,r=a.length;s<r;s++)a[s].call(this,t);t.target=null}}},Ln=["00","01","02","03","04","05","06","07","08","09","0a","0b","0c","0d","0e","0f","10","11","12","13","14","15","16","17","18","19","1a","1b","1c","1d","1e","1f","20","21","22","23","24","25","26","27","28","29","2a","2b","2c","2d","2e","2f","30","31","32","33","34","35","36","37","38","39","3a","3b","3c","3d","3e","3f","40","41","42","43","44","45","46","47","48","49","4a","4b","4c","4d","4e","4f","50","51","52","53","54","55","56","57","58","59","5a","5b","5c","5d","5e","5f","60","61","62","63","64","65","66","67","68","69","6a","6b","6c","6d","6e","6f","70","71","72","73","74","75","76","77","78","79","7a","7b","7c","7d","7e","7f","80","81","82","83","84","85","86","87","88","89","8a","8b","8c","8d","8e","8f","90","91","92","93","94","95","96","97","98","99","9a","9b","9c","9d","9e","9f","a0","a1","a2","a3","a4","a5","a6","a7","a8","a9","aa","ab","ac","ad","ae","af","b0","b1","b2","b3","b4","b5","b6","b7","b8","b9","ba","bb","bc","bd","be","bf","c0","c1","c2","c3","c4","c5","c6","c7","c8","c9","ca","cb","cc","cd","ce","cf","d0","d1","d2","d3","d4","d5","d6","d7","d8","d9","da","db","dc","dd","de","df","e0","e1","e2","e3","e4","e5","e6","e7","e8","e9","ea","eb","ec","ed","ee","ef","f0","f1","f2","f3","f4","f5","f6","f7","f8","f9","fa","fb","fc","fd","fe","ff"],mM=1234567,Bc=Math.PI/180,al=180/Math.PI;function gl(){let e=Math.random()*4294967295|0,t=Math.random()*4294967295|0,n=Math.random()*4294967295|0,i=Math.random()*4294967295|0;return(Ln[e&255]+Ln[e>>8&255]+Ln[e>>16&255]+Ln[e>>24&255]+"-"+Ln[t&255]+Ln[t>>8&255]+"-"+Ln[t>>16&15|64]+Ln[t>>24&255]+"-"+Ln[n&63|128]+Ln[n>>8&255]+"-"+Ln[n>>16&255]+Ln[n>>24&255]+Ln[i&255]+Ln[i>>8&255]+Ln[i>>16&255]+Ln[i>>24&255]).toLowerCase()}function ne(e,t,n){return Math.max(t,Math.min(n,e))}function Sv(e,t){return(e%t+t)%t}function qw(e,t,n,i,a){return i+(e-t)*(a-i)/(n-t)}function Yw(e,t,n){return e!==t?(n-e)/(t-e):0}function zc(e,t,n){return(1-n)*e+n*t}function Zw(e,t,n,i){return zc(e,t,1-Math.exp(-n*i))}function Jw(e,t=1){return t-Math.abs(Sv(e,t*2)-t)}function Kw(e,t,n){return e<=t?0:e>=n?1:(e=(e-t)/(n-t),e*e*(3-2*e))}function Qw(e,t,n){return e<=t?0:e>=n?1:(e=(e-t)/(n-t),e*e*e*(e*(e*6-15)+10))}function jw(e,t){return e+Math.floor(Math.random()*(t-e+1))}function $w(e,t){return e+Math.random()*(t-e)}function tC(e){return e*(.5-Math.random())}function eC(e){e!==void 0&&(mM=e);let t=mM+=1831565813;return t=Math.imul(t^t>>>15,t|1),t^=t+Math.imul(t^t>>>7,t|61),((t^t>>>14)>>>0)/4294967296}function nC(e){return e*Bc}function iC(e){return e*al}function aC(e){return(e&e-1)===0&&e!==0}function sC(e){return Math.pow(2,Math.ceil(Math.log(e)/Math.LN2))}function rC(e){return Math.pow(2,Math.floor(Math.log(e)/Math.LN2))}function oC(e,t,n,i,a){let s=Math.cos,r=Math.sin,o=s(n/2),l=r(n/2),c=s((t+i)/2),f=r((t+i)/2),p=s((t-i)/2),u=r((t-i)/2),d=s((i-t)/2),_=r((i-t)/2);switch(a){case"XYX":e.set(o*f,l*p,l*u,o*c);break;case"YZY":e.set(l*u,o*f,l*p,o*c);break;case"ZXZ":e.set(l*p,l*u,o*f,o*c);break;case"XZX":e.set(o*f,l*_,l*d,o*c);break;case"YXY":e.set(l*d,o*f,l*_,o*c);break;case"ZYZ":e.set(l*_,l*d,o*f,o*c);break;default:Ot("MathUtils: .setQuaternionFromProperEuler() encountered an unknown order: "+a)}}function tl(e,t){switch(t.constructor){case Float32Array:return e;case Uint32Array:return e/4294967295;case Uint16Array:return e/65535;case Uint8Array:return e/255;case Int32Array:return Math.max(e/2147483647,-1);case Int16Array:return Math.max(e/32767,-1);case Int8Array:return Math.max(e/127,-1);default:throw new Error("Invalid component type.")}}function Wn(e,t){switch(t.constructor){case Float32Array:return e;case Uint32Array:return Math.round(e*4294967295);case Uint16Array:return Math.round(e*65535);case Uint8Array:return Math.round(e*255);case Int32Array:return Math.round(e*2147483647);case Int16Array:return Math.round(e*32767);case Int8Array:return Math.round(e*127);default:throw new Error("Invalid component type.")}}var vl={DEG2RAD:Bc,RAD2DEG:al,generateUUID:gl,clamp:ne,euclideanModulo:Sv,mapLinear:qw,inverseLerp:Yw,lerp:zc,damp:Zw,pingpong:Jw,smoothstep:Kw,smootherstep:Qw,randInt:jw,randFloat:$w,randFloatSpread:tC,seededRandom:eC,degToRad:nC,radToDeg:iC,isPowerOfTwo:aC,ceilPowerOfTwo:sC,floorPowerOfTwo:rC,setQuaternionFromProperEuler:oC,normalize:Wn,denormalize:tl},te=class e{constructor(t=0,n=0){e.prototype.isVector2=!0,this.x=t,this.y=n}get width(){return this.x}set width(t){this.x=t}get height(){return this.y}set height(t){this.y=t}set(t,n){return this.x=t,this.y=n,this}setScalar(t){return this.x=t,this.y=t,this}setX(t){return this.x=t,this}setY(t){return this.y=t,this}setComponent(t,n){switch(t){case 0:this.x=n;break;case 1:this.y=n;break;default:throw new Error("index is out of range: "+t)}return this}getComponent(t){switch(t){case 0:return this.x;case 1:return this.y;default:throw new Error("index is out of range: "+t)}}clone(){return new this.constructor(this.x,this.y)}copy(t){return this.x=t.x,this.y=t.y,this}add(t){return this.x+=t.x,this.y+=t.y,this}addScalar(t){return this.x+=t,this.y+=t,this}addVectors(t,n){return this.x=t.x+n.x,this.y=t.y+n.y,this}addScaledVector(t,n){return this.x+=t.x*n,this.y+=t.y*n,this}sub(t){return this.x-=t.x,this.y-=t.y,this}subScalar(t){return this.x-=t,this.y-=t,this}subVectors(t,n){return this.x=t.x-n.x,this.y=t.y-n.y,this}multiply(t){return this.x*=t.x,this.y*=t.y,this}multiplyScalar(t){return this.x*=t,this.y*=t,this}divide(t){return this.x/=t.x,this.y/=t.y,this}divideScalar(t){return this.multiplyScalar(1/t)}applyMatrix3(t){let n=this.x,i=this.y,a=t.elements;return this.x=a[0]*n+a[3]*i+a[6],this.y=a[1]*n+a[4]*i+a[7],this}min(t){return this.x=Math.min(this.x,t.x),this.y=Math.min(this.y,t.y),this}max(t){return this.x=Math.max(this.x,t.x),this.y=Math.max(this.y,t.y),this}clamp(t,n){return this.x=ne(this.x,t.x,n.x),this.y=ne(this.y,t.y,n.y),this}clampScalar(t,n){return this.x=ne(this.x,t,n),this.y=ne(this.y,t,n),this}clampLength(t,n){let i=this.length();return this.divideScalar(i||1).multiplyScalar(ne(i,t,n))}floor(){return this.x=Math.floor(this.x),this.y=Math.floor(this.y),this}ceil(){return this.x=Math.ceil(this.x),this.y=Math.ceil(this.y),this}round(){return this.x=Math.round(this.x),this.y=Math.round(this.y),this}roundToZero(){return this.x=Math.trunc(this.x),this.y=Math.trunc(this.y),this}negate(){return this.x=-this.x,this.y=-this.y,this}dot(t){return this.x*t.x+this.y*t.y}cross(t){return this.x*t.y-this.y*t.x}lengthSq(){return this.x*this.x+this.y*this.y}length(){return Math.sqrt(this.x*this.x+this.y*this.y)}manhattanLength(){return Math.abs(this.x)+Math.abs(this.y)}normalize(){return this.divideScalar(this.length()||1)}angle(){return Math.atan2(-this.y,-this.x)+Math.PI}angleTo(t){let n=Math.sqrt(this.lengthSq()*t.lengthSq());if(n===0)return Math.PI/2;let i=this.dot(t)/n;return Math.acos(ne(i,-1,1))}distanceTo(t){return Math.sqrt(this.distanceToSquared(t))}distanceToSquared(t){let n=this.x-t.x,i=this.y-t.y;return n*n+i*i}manhattanDistanceTo(t){return Math.abs(this.x-t.x)+Math.abs(this.y-t.y)}setLength(t){return this.normalize().multiplyScalar(t)}lerp(t,n){return this.x+=(t.x-this.x)*n,this.y+=(t.y-this.y)*n,this}lerpVectors(t,n,i){return this.x=t.x+(n.x-t.x)*i,this.y=t.y+(n.y-t.y)*i,this}equals(t){return t.x===this.x&&t.y===this.y}fromArray(t,n=0){return this.x=t[n],this.y=t[n+1],this}toArray(t=[],n=0){return t[n]=this.x,t[n+1]=this.y,t}fromBufferAttribute(t,n){return this.x=t.getX(n),this.y=t.getY(n),this}rotateAround(t,n){let i=Math.cos(n),a=Math.sin(n),s=this.x-t.x,r=this.y-t.y;return this.x=s*i-r*a+t.x,this.y=s*a+r*i+t.y,this}random(){return this.x=Math.random(),this.y=Math.random(),this}*[Symbol.iterator](){yield this.x,yield this.y}},pa=class{constructor(t=0,n=0,i=0,a=1){this.isQuaternion=!0,this._x=t,this._y=n,this._z=i,this._w=a}static slerpFlat(t,n,i,a,s,r,o){let l=i[a+0],c=i[a+1],f=i[a+2],p=i[a+3],u=s[r+0],d=s[r+1],_=s[r+2],b=s[r+3];if(p!==b||l!==u||c!==d||f!==_){let v=l*u+c*d+f*_+p*b;v<0&&(u=-u,d=-d,_=-_,b=-b,v=-v);let h=1-o;if(v<.9995){let m=Math.acos(v),g=Math.sin(m);h=Math.sin(h*m)/g,o=Math.sin(o*m)/g,l=l*h+u*o,c=c*h+d*o,f=f*h+_*o,p=p*h+b*o}else{l=l*h+u*o,c=c*h+d*o,f=f*h+_*o,p=p*h+b*o;let m=1/Math.sqrt(l*l+c*c+f*f+p*p);l*=m,c*=m,f*=m,p*=m}}t[n]=l,t[n+1]=c,t[n+2]=f,t[n+3]=p}static multiplyQuaternionsFlat(t,n,i,a,s,r){let o=i[a],l=i[a+1],c=i[a+2],f=i[a+3],p=s[r],u=s[r+1],d=s[r+2],_=s[r+3];return t[n]=o*_+f*p+l*d-c*u,t[n+1]=l*_+f*u+c*p-o*d,t[n+2]=c*_+f*d+o*u-l*p,t[n+3]=f*_-o*p-l*u-c*d,t}get x(){return this._x}set x(t){this._x=t,this._onChangeCallback()}get y(){return this._y}set y(t){this._y=t,this._onChangeCallback()}get z(){return this._z}set z(t){this._z=t,this._onChangeCallback()}get w(){return this._w}set w(t){this._w=t,this._onChangeCallback()}set(t,n,i,a){return this._x=t,this._y=n,this._z=i,this._w=a,this._onChangeCallback(),this}clone(){return new this.constructor(this._x,this._y,this._z,this._w)}copy(t){return this._x=t.x,this._y=t.y,this._z=t.z,this._w=t.w,this._onChangeCallback(),this}setFromEuler(t,n=!0){let i=t._x,a=t._y,s=t._z,r=t._order,o=Math.cos,l=Math.sin,c=o(i/2),f=o(a/2),p=o(s/2),u=l(i/2),d=l(a/2),_=l(s/2);switch(r){case"XYZ":this._x=u*f*p+c*d*_,this._y=c*d*p-u*f*_,this._z=c*f*_+u*d*p,this._w=c*f*p-u*d*_;break;case"YXZ":this._x=u*f*p+c*d*_,this._y=c*d*p-u*f*_,this._z=c*f*_-u*d*p,this._w=c*f*p+u*d*_;break;case"ZXY":this._x=u*f*p-c*d*_,this._y=c*d*p+u*f*_,this._z=c*f*_+u*d*p,this._w=c*f*p-u*d*_;break;case"ZYX":this._x=u*f*p-c*d*_,this._y=c*d*p+u*f*_,this._z=c*f*_-u*d*p,this._w=c*f*p+u*d*_;break;case"YZX":this._x=u*f*p+c*d*_,this._y=c*d*p+u*f*_,this._z=c*f*_-u*d*p,this._w=c*f*p-u*d*_;break;case"XZY":this._x=u*f*p-c*d*_,this._y=c*d*p-u*f*_,this._z=c*f*_+u*d*p,this._w=c*f*p+u*d*_;break;default:Ot("Quaternion: .setFromEuler() encountered an unknown order: "+r)}return n===!0&&this._onChangeCallback(),this}setFromAxisAngle(t,n){let i=n/2,a=Math.sin(i);return this._x=t.x*a,this._y=t.y*a,this._z=t.z*a,this._w=Math.cos(i),this._onChangeCallback(),this}setFromRotationMatrix(t){let n=t.elements,i=n[0],a=n[4],s=n[8],r=n[1],o=n[5],l=n[9],c=n[2],f=n[6],p=n[10],u=i+o+p;if(u>0){let d=.5/Math.sqrt(u+1);this._w=.25/d,this._x=(f-l)*d,this._y=(s-c)*d,this._z=(r-a)*d}else if(i>o&&i>p){let d=2*Math.sqrt(1+i-o-p);this._w=(f-l)/d,this._x=.25*d,this._y=(a+r)/d,this._z=(s+c)/d}else if(o>p){let d=2*Math.sqrt(1+o-i-p);this._w=(s-c)/d,this._x=(a+r)/d,this._y=.25*d,this._z=(l+f)/d}else{let d=2*Math.sqrt(1+p-i-o);this._w=(r-a)/d,this._x=(s+c)/d,this._y=(l+f)/d,this._z=.25*d}return this._onChangeCallback(),this}setFromUnitVectors(t,n){let i=t.dot(n)+1;return i<1e-8?(i=0,Math.abs(t.x)>Math.abs(t.z)?(this._x=-t.y,this._y=t.x,this._z=0,this._w=i):(this._x=0,this._y=-t.z,this._z=t.y,this._w=i)):(this._x=t.y*n.z-t.z*n.y,this._y=t.z*n.x-t.x*n.z,this._z=t.x*n.y-t.y*n.x,this._w=i),this.normalize()}angleTo(t){return 2*Math.acos(Math.abs(ne(this.dot(t),-1,1)))}rotateTowards(t,n){let i=this.angleTo(t);if(i===0)return this;let a=Math.min(1,n/i);return this.slerp(t,a),this}identity(){return this.set(0,0,0,1)}invert(){return this.conjugate()}conjugate(){return this._x*=-1,this._y*=-1,this._z*=-1,this._onChangeCallback(),this}dot(t){return this._x*t._x+this._y*t._y+this._z*t._z+this._w*t._w}lengthSq(){return this._x*this._x+this._y*this._y+this._z*this._z+this._w*this._w}length(){return Math.sqrt(this._x*this._x+this._y*this._y+this._z*this._z+this._w*this._w)}normalize(){let t=this.length();return t===0?(this._x=0,this._y=0,this._z=0,this._w=1):(t=1/t,this._x=this._x*t,this._y=this._y*t,this._z=this._z*t,this._w=this._w*t),this._onChangeCallback(),this}multiply(t){return this.multiplyQuaternions(this,t)}premultiply(t){return this.multiplyQuaternions(t,this)}multiplyQuaternions(t,n){let i=t._x,a=t._y,s=t._z,r=t._w,o=n._x,l=n._y,c=n._z,f=n._w;return this._x=i*f+r*o+a*c-s*l,this._y=a*f+r*l+s*o-i*c,this._z=s*f+r*c+i*l-a*o,this._w=r*f-i*o-a*l-s*c,this._onChangeCallback(),this}slerp(t,n){let i=t._x,a=t._y,s=t._z,r=t._w,o=this.dot(t);o<0&&(i=-i,a=-a,s=-s,r=-r,o=-o);let l=1-n;if(o<.9995){let c=Math.acos(o),f=Math.sin(c);l=Math.sin(l*c)/f,n=Math.sin(n*c)/f,this._x=this._x*l+i*n,this._y=this._y*l+a*n,this._z=this._z*l+s*n,this._w=this._w*l+r*n,this._onChangeCallback()}else this._x=this._x*l+i*n,this._y=this._y*l+a*n,this._z=this._z*l+s*n,this._w=this._w*l+r*n,this.normalize();return this}slerpQuaternions(t,n,i){return this.copy(t).slerp(n,i)}random(){let t=2*Math.PI*Math.random(),n=2*Math.PI*Math.random(),i=Math.random(),a=Math.sqrt(1-i),s=Math.sqrt(i);return this.set(a*Math.sin(t),a*Math.cos(t),s*Math.sin(n),s*Math.cos(n))}equals(t){return t._x===this._x&&t._y===this._y&&t._z===this._z&&t._w===this._w}fromArray(t,n=0){return this._x=t[n],this._y=t[n+1],this._z=t[n+2],this._w=t[n+3],this._onChangeCallback(),this}toArray(t=[],n=0){return t[n]=this._x,t[n+1]=this._y,t[n+2]=this._z,t[n+3]=this._w,t}fromBufferAttribute(t,n){return this._x=t.getX(n),this._y=t.getY(n),this._z=t.getZ(n),this._w=t.getW(n),this._onChangeCallback(),this}toJSON(){return this.toArray()}_onChange(t){return this._onChangeCallback=t,this}_onChangeCallback(){}*[Symbol.iterator](){yield this._x,yield this._y,yield this._z,yield this._w}},B=class e{constructor(t=0,n=0,i=0){e.prototype.isVector3=!0,this.x=t,this.y=n,this.z=i}set(t,n,i){return i===void 0&&(i=this.z),this.x=t,this.y=n,this.z=i,this}setScalar(t){return this.x=t,this.y=t,this.z=t,this}setX(t){return this.x=t,this}setY(t){return this.y=t,this}setZ(t){return this.z=t,this}setComponent(t,n){switch(t){case 0:this.x=n;break;case 1:this.y=n;break;case 2:this.z=n;break;default:throw new Error("index is out of range: "+t)}return this}getComponent(t){switch(t){case 0:return this.x;case 1:return this.y;case 2:return this.z;default:throw new Error("index is out of range: "+t)}}clone(){return new this.constructor(this.x,this.y,this.z)}copy(t){return this.x=t.x,this.y=t.y,this.z=t.z,this}add(t){return this.x+=t.x,this.y+=t.y,this.z+=t.z,this}addScalar(t){return this.x+=t,this.y+=t,this.z+=t,this}addVectors(t,n){return this.x=t.x+n.x,this.y=t.y+n.y,this.z=t.z+n.z,this}addScaledVector(t,n){return this.x+=t.x*n,this.y+=t.y*n,this.z+=t.z*n,this}sub(t){return this.x-=t.x,this.y-=t.y,this.z-=t.z,this}subScalar(t){return this.x-=t,this.y-=t,this.z-=t,this}subVectors(t,n){return this.x=t.x-n.x,this.y=t.y-n.y,this.z=t.z-n.z,this}multiply(t){return this.x*=t.x,this.y*=t.y,this.z*=t.z,this}multiplyScalar(t){return this.x*=t,this.y*=t,this.z*=t,this}multiplyVectors(t,n){return this.x=t.x*n.x,this.y=t.y*n.y,this.z=t.z*n.z,this}applyEuler(t){return this.applyQuaternion(gM.setFromEuler(t))}applyAxisAngle(t,n){return this.applyQuaternion(gM.setFromAxisAngle(t,n))}applyMatrix3(t){let n=this.x,i=this.y,a=this.z,s=t.elements;return this.x=s[0]*n+s[3]*i+s[6]*a,this.y=s[1]*n+s[4]*i+s[7]*a,this.z=s[2]*n+s[5]*i+s[8]*a,this}applyNormalMatrix(t){return this.applyMatrix3(t).normalize()}applyMatrix4(t){let n=this.x,i=this.y,a=this.z,s=t.elements,r=1/(s[3]*n+s[7]*i+s[11]*a+s[15]);return this.x=(s[0]*n+s[4]*i+s[8]*a+s[12])*r,this.y=(s[1]*n+s[5]*i+s[9]*a+s[13])*r,this.z=(s[2]*n+s[6]*i+s[10]*a+s[14])*r,this}applyQuaternion(t){let n=this.x,i=this.y,a=this.z,s=t.x,r=t.y,o=t.z,l=t.w,c=2*(r*a-o*i),f=2*(o*n-s*a),p=2*(s*i-r*n);return this.x=n+l*c+r*p-o*f,this.y=i+l*f+o*c-s*p,this.z=a+l*p+s*f-r*c,this}project(t){return this.applyMatrix4(t.matrixWorldInverse).applyMatrix4(t.projectionMatrix)}unproject(t){return this.applyMatrix4(t.projectionMatrixInverse).applyMatrix4(t.matrixWorld)}transformDirection(t){let n=this.x,i=this.y,a=this.z,s=t.elements;return this.x=s[0]*n+s[4]*i+s[8]*a,this.y=s[1]*n+s[5]*i+s[9]*a,this.z=s[2]*n+s[6]*i+s[10]*a,this.normalize()}divide(t){return this.x/=t.x,this.y/=t.y,this.z/=t.z,this}divideScalar(t){return this.multiplyScalar(1/t)}min(t){return this.x=Math.min(this.x,t.x),this.y=Math.min(this.y,t.y),this.z=Math.min(this.z,t.z),this}max(t){return this.x=Math.max(this.x,t.x),this.y=Math.max(this.y,t.y),this.z=Math.max(this.z,t.z),this}clamp(t,n){return this.x=ne(this.x,t.x,n.x),this.y=ne(this.y,t.y,n.y),this.z=ne(this.z,t.z,n.z),this}clampScalar(t,n){return this.x=ne(this.x,t,n),this.y=ne(this.y,t,n),this.z=ne(this.z,t,n),this}clampLength(t,n){let i=this.length();return this.divideScalar(i||1).multiplyScalar(ne(i,t,n))}floor(){return this.x=Math.floor(this.x),this.y=Math.floor(this.y),this.z=Math.floor(this.z),this}ceil(){return this.x=Math.ceil(this.x),this.y=Math.ceil(this.y),this.z=Math.ceil(this.z),this}round(){return this.x=Math.round(this.x),this.y=Math.round(this.y),this.z=Math.round(this.z),this}roundToZero(){return this.x=Math.trunc(this.x),this.y=Math.trunc(this.y),this.z=Math.trunc(this.z),this}negate(){return this.x=-this.x,this.y=-this.y,this.z=-this.z,this}dot(t){return this.x*t.x+this.y*t.y+this.z*t.z}lengthSq(){return this.x*this.x+this.y*this.y+this.z*this.z}length(){return Math.sqrt(this.x*this.x+this.y*this.y+this.z*this.z)}manhattanLength(){return Math.abs(this.x)+Math.abs(this.y)+Math.abs(this.z)}normalize(){return this.divideScalar(this.length()||1)}setLength(t){return this.normalize().multiplyScalar(t)}lerp(t,n){return this.x+=(t.x-this.x)*n,this.y+=(t.y-this.y)*n,this.z+=(t.z-this.z)*n,this}lerpVectors(t,n,i){return this.x=t.x+(n.x-t.x)*i,this.y=t.y+(n.y-t.y)*i,this.z=t.z+(n.z-t.z)*i,this}cross(t){return this.crossVectors(this,t)}crossVectors(t,n){let i=t.x,a=t.y,s=t.z,r=n.x,o=n.y,l=n.z;return this.x=a*l-s*o,this.y=s*r-i*l,this.z=i*o-a*r,this}projectOnVector(t){let n=t.lengthSq();if(n===0)return this.set(0,0,0);let i=t.dot(this)/n;return this.copy(t).multiplyScalar(i)}projectOnPlane(t){return b0.copy(this).projectOnVector(t),this.sub(b0)}reflect(t){return this.sub(b0.copy(t).multiplyScalar(2*this.dot(t)))}angleTo(t){let n=Math.sqrt(this.lengthSq()*t.lengthSq());if(n===0)return Math.PI/2;let i=this.dot(t)/n;return Math.acos(ne(i,-1,1))}distanceTo(t){return Math.sqrt(this.distanceToSquared(t))}distanceToSquared(t){let n=this.x-t.x,i=this.y-t.y,a=this.z-t.z;return n*n+i*i+a*a}manhattanDistanceTo(t){return Math.abs(this.x-t.x)+Math.abs(this.y-t.y)+Math.abs(this.z-t.z)}setFromSpherical(t){return this.setFromSphericalCoords(t.radius,t.phi,t.theta)}setFromSphericalCoords(t,n,i){let a=Math.sin(n)*t;return this.x=a*Math.sin(i),this.y=Math.cos(n)*t,this.z=a*Math.cos(i),this}setFromCylindrical(t){return this.setFromCylindricalCoords(t.radius,t.theta,t.y)}setFromCylindricalCoords(t,n,i){return this.x=t*Math.sin(n),this.y=i,this.z=t*Math.cos(n),this}setFromMatrixPosition(t){let n=t.elements;return this.x=n[12],this.y=n[13],this.z=n[14],this}setFromMatrixScale(t){let n=this.setFromMatrixColumn(t,0).length(),i=this.setFromMatrixColumn(t,1).length(),a=this.setFromMatrixColumn(t,2).length();return this.x=n,this.y=i,this.z=a,this}setFromMatrixColumn(t,n){return this.fromArray(t.elements,n*4)}setFromMatrix3Column(t,n){return this.fromArray(t.elements,n*3)}setFromEuler(t){return this.x=t._x,this.y=t._y,this.z=t._z,this}setFromColor(t){return this.x=t.r,this.y=t.g,this.z=t.b,this}equals(t){return t.x===this.x&&t.y===this.y&&t.z===this.z}fromArray(t,n=0){return this.x=t[n],this.y=t[n+1],this.z=t[n+2],this}toArray(t=[],n=0){return t[n]=this.x,t[n+1]=this.y,t[n+2]=this.z,t}fromBufferAttribute(t,n){return this.x=t.getX(n),this.y=t.getY(n),this.z=t.getZ(n),this}random(){return this.x=Math.random(),this.y=Math.random(),this.z=Math.random(),this}randomDirection(){let t=Math.random()*Math.PI*2,n=Math.random()*2-1,i=Math.sqrt(1-n*n);return this.x=i*Math.cos(t),this.y=n,this.z=i*Math.sin(t),this}*[Symbol.iterator](){yield this.x,yield this.y,yield this.z}},b0=new B,gM=new pa,Zt=class e{constructor(t,n,i,a,s,r,o,l,c){e.prototype.isMatrix3=!0,this.elements=[1,0,0,0,1,0,0,0,1],t!==void 0&&this.set(t,n,i,a,s,r,o,l,c)}set(t,n,i,a,s,r,o,l,c){let f=this.elements;return f[0]=t,f[1]=a,f[2]=o,f[3]=n,f[4]=s,f[5]=l,f[6]=i,f[7]=r,f[8]=c,this}identity(){return this.set(1,0,0,0,1,0,0,0,1),this}copy(t){let n=this.elements,i=t.elements;return n[0]=i[0],n[1]=i[1],n[2]=i[2],n[3]=i[3],n[4]=i[4],n[5]=i[5],n[6]=i[6],n[7]=i[7],n[8]=i[8],this}extractBasis(t,n,i){return t.setFromMatrix3Column(this,0),n.setFromMatrix3Column(this,1),i.setFromMatrix3Column(this,2),this}setFromMatrix4(t){let n=t.elements;return this.set(n[0],n[4],n[8],n[1],n[5],n[9],n[2],n[6],n[10]),this}multiply(t){return this.multiplyMatrices(this,t)}premultiply(t){return this.multiplyMatrices(t,this)}multiplyMatrices(t,n){let i=t.elements,a=n.elements,s=this.elements,r=i[0],o=i[3],l=i[6],c=i[1],f=i[4],p=i[7],u=i[2],d=i[5],_=i[8],b=a[0],v=a[3],h=a[6],m=a[1],g=a[4],x=a[7],E=a[2],A=a[5],w=a[8];return s[0]=r*b+o*m+l*E,s[3]=r*v+o*g+l*A,s[6]=r*h+o*x+l*w,s[1]=c*b+f*m+p*E,s[4]=c*v+f*g+p*A,s[7]=c*h+f*x+p*w,s[2]=u*b+d*m+_*E,s[5]=u*v+d*g+_*A,s[8]=u*h+d*x+_*w,this}multiplyScalar(t){let n=this.elements;return n[0]*=t,n[3]*=t,n[6]*=t,n[1]*=t,n[4]*=t,n[7]*=t,n[2]*=t,n[5]*=t,n[8]*=t,this}determinant(){let t=this.elements,n=t[0],i=t[1],a=t[2],s=t[3],r=t[4],o=t[5],l=t[6],c=t[7],f=t[8];return n*r*f-n*o*c-i*s*f+i*o*l+a*s*c-a*r*l}invert(){let t=this.elements,n=t[0],i=t[1],a=t[2],s=t[3],r=t[4],o=t[5],l=t[6],c=t[7],f=t[8],p=f*r-o*c,u=o*l-f*s,d=c*s-r*l,_=n*p+i*u+a*d;if(_===0)return this.set(0,0,0,0,0,0,0,0,0);let b=1/_;return t[0]=p*b,t[1]=(a*c-f*i)*b,t[2]=(o*i-a*r)*b,t[3]=u*b,t[4]=(f*n-a*l)*b,t[5]=(a*s-o*n)*b,t[6]=d*b,t[7]=(i*l-c*n)*b,t[8]=(r*n-i*s)*b,this}transpose(){let t,n=this.elements;return t=n[1],n[1]=n[3],n[3]=t,t=n[2],n[2]=n[6],n[6]=t,t=n[5],n[5]=n[7],n[7]=t,this}getNormalMatrix(t){return this.setFromMatrix4(t).invert().transpose()}transposeIntoArray(t){let n=this.elements;return t[0]=n[0],t[1]=n[3],t[2]=n[6],t[3]=n[1],t[4]=n[4],t[5]=n[7],t[6]=n[2],t[7]=n[5],t[8]=n[8],this}setUvTransform(t,n,i,a,s,r,o){let l=Math.cos(s),c=Math.sin(s);return this.set(i*l,i*c,-i*(l*r+c*o)+r+t,-a*c,a*l,-a*(-c*r+l*o)+o+n,0,0,1),this}scale(t,n){return this.premultiply(M0.makeScale(t,n)),this}rotate(t){return this.premultiply(M0.makeRotation(-t)),this}translate(t,n){return this.premultiply(M0.makeTranslation(t,n)),this}makeTranslation(t,n){return t.isVector2?this.set(1,0,t.x,0,1,t.y,0,0,1):this.set(1,0,t,0,1,n,0,0,1),this}makeRotation(t){let n=Math.cos(t),i=Math.sin(t);return this.set(n,-i,0,i,n,0,0,0,1),this}makeScale(t,n){return this.set(t,0,0,0,n,0,0,0,1),this}equals(t){let n=this.elements,i=t.elements;for(let a=0;a<9;a++)if(n[a]!==i[a])return!1;return!0}fromArray(t,n=0){for(let i=0;i<9;i++)this.elements[i]=t[i+n];return this}toArray(t=[],n=0){let i=this.elements;return t[n]=i[0],t[n+1]=i[1],t[n+2]=i[2],t[n+3]=i[3],t[n+4]=i[4],t[n+5]=i[5],t[n+6]=i[6],t[n+7]=i[7],t[n+8]=i[8],t}clone(){return new this.constructor().fromArray(this.elements)}},M0=new Zt,vM=new Zt().set(.4123908,.3575843,.1804808,.212639,.7151687,.0721923,.0193308,.1191948,.9505322),_M=new Zt().set(3.2409699,-1.5373832,-.4986108,-.9692436,1.8759675,.0415551,.0556301,-.203977,1.0569715);function lC(){let e={enabled:!0,workingColorSpace:Rr,spaces:{},convert:function(a,s,r){return this.enabled===!1||s===r||!s||!r||(this.spaces[s].transfer===me&&(a.r=Xa(a.r),a.g=Xa(a.g),a.b=Xa(a.b)),this.spaces[s].primaries!==this.spaces[r].primaries&&(a.applyMatrix3(this.spaces[s].toXYZ),a.applyMatrix3(this.spaces[r].fromXYZ)),this.spaces[r].transfer===me&&(a.r=el(a.r),a.g=el(a.g),a.b=el(a.b))),a},workingToColorSpace:function(a,s){return this.convert(a,this.workingColorSpace,s)},colorSpaceToWorking:function(a,s){return this.convert(a,s,this.workingColorSpace)},getPrimaries:function(a){return this.spaces[a].primaries},getTransfer:function(a){return a===Ja?Vc:this.spaces[a].transfer},getToneMappingMode:function(a){return this.spaces[a].outputColorSpaceConfig.toneMappingMode||"standard"},getLuminanceCoefficients:function(a,s=this.workingColorSpace){return a.fromArray(this.spaces[s].luminanceCoefficients)},define:function(a){Object.assign(this.spaces,a)},_getMatrix:function(a,s,r){return a.copy(this.spaces[s].toXYZ).multiply(this.spaces[r].fromXYZ)},_getDrawingBufferColorSpace:function(a){return this.spaces[a].outputColorSpaceConfig.drawingBufferColorSpace},_getUnpackColorSpace:function(a=this.workingColorSpace){return this.spaces[a].workingColorSpaceConfig.unpackColorSpace},fromWorkingColorSpace:function(a,s){return Gc("ColorManagement: .fromWorkingColorSpace() has been renamed to .workingToColorSpace()."),e.workingToColorSpace(a,s)},toWorkingColorSpace:function(a,s){return Gc("ColorManagement: .toWorkingColorSpace() has been renamed to .colorSpaceToWorking()."),e.colorSpaceToWorking(a,s)}},t=[.64,.33,.3,.6,.15,.06],n=[.2126,.7152,.0722],i=[.3127,.329];return e.define({[Rr]:{primaries:t,whitePoint:i,transfer:Vc,toXYZ:vM,fromXYZ:_M,luminanceCoefficients:n,workingColorSpaceConfig:{unpackColorSpace:On},outputColorSpaceConfig:{drawingBufferColorSpace:On}},[On]:{primaries:t,whitePoint:i,transfer:me,toXYZ:vM,fromXYZ:_M,luminanceCoefficients:n,outputColorSpaceConfig:{drawingBufferColorSpace:On}}}),e}var se=lC();function Xa(e){return e<.04045?e*.0773993808:Math.pow(e*.9478672986+.0521327014,2.4)}function el(e){return e<.0031308?e*12.92:1.055*Math.pow(e,.41666)-.055}var Ho,Jh=class{static getDataURL(t,n="image/png"){if(/^data:/i.test(t.src)||typeof HTMLCanvasElement>"u")return t.src;let i;if(t instanceof HTMLCanvasElement)i=t;else{Ho===void 0&&(Ho=Hc("canvas")),Ho.width=t.width,Ho.height=t.height;let a=Ho.getContext("2d");t instanceof ImageData?a.putImageData(t,0,0):a.drawImage(t,0,0,t.width,t.height),i=Ho}return i.toDataURL(n)}static sRGBToLinear(t){if(typeof HTMLImageElement<"u"&&t instanceof HTMLImageElement||typeof HTMLCanvasElement<"u"&&t instanceof HTMLCanvasElement||typeof ImageBitmap<"u"&&t instanceof ImageBitmap){let n=Hc("canvas");n.width=t.width,n.height=t.height;let i=n.getContext("2d");i.drawImage(t,0,0,t.width,t.height);let a=i.getImageData(0,0,t.width,t.height),s=a.data;for(let r=0;r<s.length;r++)s[r]=Xa(s[r]/255)*255;return i.putImageData(a,0,0),n}else if(t.data){let n=t.data.slice(0);for(let i=0;i<n.length;i++)n instanceof Uint8Array||n instanceof Uint8ClampedArray?n[i]=Math.floor(Xa(n[i]/255)*255):n[i]=Xa(n[i]);return{data:n,width:t.width,height:t.height}}else return Ot("ImageUtils.sRGBToLinear(): Unsupported image type. No color space conversion applied."),t}},cC=0,sl=class{constructor(t=null){this.isSource=!0,Object.defineProperty(this,"id",{value:cC++}),this.uuid=gl(),this.data=t,this.dataReady=!0,this.version=0}getSize(t){let n=this.data;return typeof HTMLVideoElement<"u"&&n instanceof HTMLVideoElement?t.set(n.videoWidth,n.videoHeight,0):typeof VideoFrame<"u"&&n instanceof VideoFrame?t.set(n.displayHeight,n.displayWidth,0):n!==null?t.set(n.width,n.height,n.depth||0):t.set(0,0,0),t}set needsUpdate(t){t===!0&&this.version++}toJSON(t){let n=t===void 0||typeof t=="string";if(!n&&t.images[this.uuid]!==void 0)return t.images[this.uuid];let i={uuid:this.uuid,url:""},a=this.data;if(a!==null){let s;if(Array.isArray(a)){s=[];for(let r=0,o=a.length;r<o;r++)a[r].isDataTexture?s.push(T0(a[r].image)):s.push(T0(a[r]))}else s=T0(a);i.url=s}return n||(t.images[this.uuid]=i),i}};function T0(e){return typeof HTMLImageElement<"u"&&e instanceof HTMLImageElement||typeof HTMLCanvasElement<"u"&&e instanceof HTMLCanvasElement||typeof ImageBitmap<"u"&&e instanceof ImageBitmap?Jh.getDataURL(e):e.data?{data:Array.from(e.data),width:e.width,height:e.height,type:e.data.constructor.name}:(Ot("Texture: Unable to serialize Texture."),{})}var uC=0,E0=new B,qn=class e extends qa{constructor(t=e.DEFAULT_IMAGE,n=e.DEFAULT_MAPPING,i=ha,a=ha,s=En,r=Ws,o=Ui,l=ei,c=e.DEFAULT_ANISOTROPY,f=Ja){super(),this.isTexture=!0,Object.defineProperty(this,"id",{value:uC++}),this.uuid=gl(),this.name="",this.source=new sl(t),this.mipmaps=[],this.mapping=n,this.channel=0,this.wrapS=i,this.wrapT=a,this.magFilter=s,this.minFilter=r,this.anisotropy=c,this.format=o,this.internalFormat=null,this.type=l,this.offset=new te(0,0),this.repeat=new te(1,1),this.center=new te(0,0),this.rotation=0,this.matrixAutoUpdate=!0,this.matrix=new Zt,this.generateMipmaps=!0,this.premultiplyAlpha=!1,this.flipY=!0,this.unpackAlignment=4,this.colorSpace=f,this.userData={},this.updateRanges=[],this.version=0,this.onUpdate=null,this.renderTarget=null,this.isRenderTargetTexture=!1,this.isArrayTexture=!!(t&&t.depth&&t.depth>1),this.pmremVersion=0}get width(){return this.source.getSize(E0).x}get height(){return this.source.getSize(E0).y}get depth(){return this.source.getSize(E0).z}get image(){return this.source.data}set image(t=null){this.source.data=t}updateMatrix(){this.matrix.setUvTransform(this.offset.x,this.offset.y,this.repeat.x,this.repeat.y,this.rotation,this.center.x,this.center.y)}addUpdateRange(t,n){this.updateRanges.push({start:t,count:n})}clearUpdateRanges(){this.updateRanges.length=0}clone(){return new this.constructor().copy(this)}copy(t){return this.name=t.name,this.source=t.source,this.mipmaps=t.mipmaps.slice(0),this.mapping=t.mapping,this.channel=t.channel,this.wrapS=t.wrapS,this.wrapT=t.wrapT,this.magFilter=t.magFilter,this.minFilter=t.minFilter,this.anisotropy=t.anisotropy,this.format=t.format,this.internalFormat=t.internalFormat,this.type=t.type,this.offset.copy(t.offset),this.repeat.copy(t.repeat),this.center.copy(t.center),this.rotation=t.rotation,this.matrixAutoUpdate=t.matrixAutoUpdate,this.matrix.copy(t.matrix),this.generateMipmaps=t.generateMipmaps,this.premultiplyAlpha=t.premultiplyAlpha,this.flipY=t.flipY,this.unpackAlignment=t.unpackAlignment,this.colorSpace=t.colorSpace,this.renderTarget=t.renderTarget,this.isRenderTargetTexture=t.isRenderTargetTexture,this.isArrayTexture=t.isArrayTexture,this.userData=JSON.parse(JSON.stringify(t.userData)),this.needsUpdate=!0,this}setValues(t){for(let n in t){let i=t[n];if(i===void 0){Ot(`Texture.setValues(): parameter '${n}' has value of undefined.`);continue}let a=this[n];if(a===void 0){Ot(`Texture.setValues(): property '${n}' does not exist.`);continue}a&&i&&a.isVector2&&i.isVector2||a&&i&&a.isVector3&&i.isVector3||a&&i&&a.isMatrix3&&i.isMatrix3?a.copy(i):this[n]=i}}toJSON(t){let n=t===void 0||typeof t=="string";if(!n&&t.textures[this.uuid]!==void 0)return t.textures[this.uuid];let i={metadata:{version:4.7,type:"Texture",generator:"Texture.toJSON"},uuid:this.uuid,name:this.name,image:this.source.toJSON(t).uuid,mapping:this.mapping,channel:this.channel,repeat:[this.repeat.x,this.repeat.y],offset:[this.offset.x,this.offset.y],center:[this.center.x,this.center.y],rotation:this.rotation,wrap:[this.wrapS,this.wrapT],format:this.format,internalFormat:this.internalFormat,type:this.type,colorSpace:this.colorSpace,minFilter:this.minFilter,magFilter:this.magFilter,anisotropy:this.anisotropy,flipY:this.flipY,generateMipmaps:this.generateMipmaps,premultiplyAlpha:this.premultiplyAlpha,unpackAlignment:this.unpackAlignment};return Object.keys(this.userData).length>0&&(i.userData=this.userData),n||(t.textures[this.uuid]=i),i}dispose(){this.dispatchEvent({type:"dispose"})}transformUv(t){if(this.mapping!==fv)return t;if(t.applyMatrix3(this.matrix),t.x<0||t.x>1)switch(this.wrapS){case qh:t.x=t.x-Math.floor(t.x);break;case ha:t.x=t.x<0?0:1;break;case Yh:Math.abs(Math.floor(t.x)%2)===1?t.x=Math.ceil(t.x)-t.x:t.x=t.x-Math.floor(t.x);break}if(t.y<0||t.y>1)switch(this.wrapT){case qh:t.y=t.y-Math.floor(t.y);break;case ha:t.y=t.y<0?0:1;break;case Yh:Math.abs(Math.floor(t.y)%2)===1?t.y=Math.ceil(t.y)-t.y:t.y=t.y-Math.floor(t.y);break}return this.flipY&&(t.y=1-t.y),t}set needsUpdate(t){t===!0&&(this.version++,this.source.needsUpdate=!0)}set needsPMREMUpdate(t){t===!0&&this.pmremVersion++}};qn.DEFAULT_IMAGE=null;qn.DEFAULT_MAPPING=fv;qn.DEFAULT_ANISOTROPY=1;var Ge=class e{constructor(t=0,n=0,i=0,a=1){e.prototype.isVector4=!0,this.x=t,this.y=n,this.z=i,this.w=a}get width(){return this.z}set width(t){this.z=t}get height(){return this.w}set height(t){this.w=t}set(t,n,i,a){return this.x=t,this.y=n,this.z=i,this.w=a,this}setScalar(t){return this.x=t,this.y=t,this.z=t,this.w=t,this}setX(t){return this.x=t,this}setY(t){return this.y=t,this}setZ(t){return this.z=t,this}setW(t){return this.w=t,this}setComponent(t,n){switch(t){case 0:this.x=n;break;case 1:this.y=n;break;case 2:this.z=n;break;case 3:this.w=n;break;default:throw new Error("index is out of range: "+t)}return this}getComponent(t){switch(t){case 0:return this.x;case 1:return this.y;case 2:return this.z;case 3:return this.w;default:throw new Error("index is out of range: "+t)}}clone(){return new this.constructor(this.x,this.y,this.z,this.w)}copy(t){return this.x=t.x,this.y=t.y,this.z=t.z,this.w=t.w!==void 0?t.w:1,this}add(t){return this.x+=t.x,this.y+=t.y,this.z+=t.z,this.w+=t.w,this}addScalar(t){return this.x+=t,this.y+=t,this.z+=t,this.w+=t,this}addVectors(t,n){return this.x=t.x+n.x,this.y=t.y+n.y,this.z=t.z+n.z,this.w=t.w+n.w,this}addScaledVector(t,n){return this.x+=t.x*n,this.y+=t.y*n,this.z+=t.z*n,this.w+=t.w*n,this}sub(t){return this.x-=t.x,this.y-=t.y,this.z-=t.z,this.w-=t.w,this}subScalar(t){return this.x-=t,this.y-=t,this.z-=t,this.w-=t,this}subVectors(t,n){return this.x=t.x-n.x,this.y=t.y-n.y,this.z=t.z-n.z,this.w=t.w-n.w,this}multiply(t){return this.x*=t.x,this.y*=t.y,this.z*=t.z,this.w*=t.w,this}multiplyScalar(t){return this.x*=t,this.y*=t,this.z*=t,this.w*=t,this}applyMatrix4(t){let n=this.x,i=this.y,a=this.z,s=this.w,r=t.elements;return this.x=r[0]*n+r[4]*i+r[8]*a+r[12]*s,this.y=r[1]*n+r[5]*i+r[9]*a+r[13]*s,this.z=r[2]*n+r[6]*i+r[10]*a+r[14]*s,this.w=r[3]*n+r[7]*i+r[11]*a+r[15]*s,this}divide(t){return this.x/=t.x,this.y/=t.y,this.z/=t.z,this.w/=t.w,this}divideScalar(t){return this.multiplyScalar(1/t)}setAxisAngleFromQuaternion(t){this.w=2*Math.acos(t.w);let n=Math.sqrt(1-t.w*t.w);return n<1e-4?(this.x=1,this.y=0,this.z=0):(this.x=t.x/n,this.y=t.y/n,this.z=t.z/n),this}setAxisAngleFromRotationMatrix(t){let n,i,a,s,l=t.elements,c=l[0],f=l[4],p=l[8],u=l[1],d=l[5],_=l[9],b=l[2],v=l[6],h=l[10];if(Math.abs(f-u)<.01&&Math.abs(p-b)<.01&&Math.abs(_-v)<.01){if(Math.abs(f+u)<.1&&Math.abs(p+b)<.1&&Math.abs(_+v)<.1&&Math.abs(c+d+h-3)<.1)return this.set(1,0,0,0),this;n=Math.PI;let g=(c+1)/2,x=(d+1)/2,E=(h+1)/2,A=(f+u)/4,w=(p+b)/4,y=(_+v)/4;return g>x&&g>E?g<.01?(i=0,a=.707106781,s=.707106781):(i=Math.sqrt(g),a=A/i,s=w/i):x>E?x<.01?(i=.707106781,a=0,s=.707106781):(a=Math.sqrt(x),i=A/a,s=y/a):E<.01?(i=.707106781,a=.707106781,s=0):(s=Math.sqrt(E),i=w/s,a=y/s),this.set(i,a,s,n),this}let m=Math.sqrt((v-_)*(v-_)+(p-b)*(p-b)+(u-f)*(u-f));return Math.abs(m)<.001&&(m=1),this.x=(v-_)/m,this.y=(p-b)/m,this.z=(u-f)/m,this.w=Math.acos((c+d+h-1)/2),this}setFromMatrixPosition(t){let n=t.elements;return this.x=n[12],this.y=n[13],this.z=n[14],this.w=n[15],this}min(t){return this.x=Math.min(this.x,t.x),this.y=Math.min(this.y,t.y),this.z=Math.min(this.z,t.z),this.w=Math.min(this.w,t.w),this}max(t){return this.x=Math.max(this.x,t.x),this.y=Math.max(this.y,t.y),this.z=Math.max(this.z,t.z),this.w=Math.max(this.w,t.w),this}clamp(t,n){return this.x=ne(this.x,t.x,n.x),this.y=ne(this.y,t.y,n.y),this.z=ne(this.z,t.z,n.z),this.w=ne(this.w,t.w,n.w),this}clampScalar(t,n){return this.x=ne(this.x,t,n),this.y=ne(this.y,t,n),this.z=ne(this.z,t,n),this.w=ne(this.w,t,n),this}clampLength(t,n){let i=this.length();return this.divideScalar(i||1).multiplyScalar(ne(i,t,n))}floor(){return this.x=Math.floor(this.x),this.y=Math.floor(this.y),this.z=Math.floor(this.z),this.w=Math.floor(this.w),this}ceil(){return this.x=Math.ceil(this.x),this.y=Math.ceil(this.y),this.z=Math.ceil(this.z),this.w=Math.ceil(this.w),this}round(){return this.x=Math.round(this.x),this.y=Math.round(this.y),this.z=Math.round(this.z),this.w=Math.round(this.w),this}roundToZero(){return this.x=Math.trunc(this.x),this.y=Math.trunc(this.y),this.z=Math.trunc(this.z),this.w=Math.trunc(this.w),this}negate(){return this.x=-this.x,this.y=-this.y,this.z=-this.z,this.w=-this.w,this}dot(t){return this.x*t.x+this.y*t.y+this.z*t.z+this.w*t.w}lengthSq(){return this.x*this.x+this.y*this.y+this.z*this.z+this.w*this.w}length(){return Math.sqrt(this.x*this.x+this.y*this.y+this.z*this.z+this.w*this.w)}manhattanLength(){return Math.abs(this.x)+Math.abs(this.y)+Math.abs(this.z)+Math.abs(this.w)}normalize(){return this.divideScalar(this.length()||1)}setLength(t){return this.normalize().multiplyScalar(t)}lerp(t,n){return this.x+=(t.x-this.x)*n,this.y+=(t.y-this.y)*n,this.z+=(t.z-this.z)*n,this.w+=(t.w-this.w)*n,this}lerpVectors(t,n,i){return this.x=t.x+(n.x-t.x)*i,this.y=t.y+(n.y-t.y)*i,this.z=t.z+(n.z-t.z)*i,this.w=t.w+(n.w-t.w)*i,this}equals(t){return t.x===this.x&&t.y===this.y&&t.z===this.z&&t.w===this.w}fromArray(t,n=0){return this.x=t[n],this.y=t[n+1],this.z=t[n+2],this.w=t[n+3],this}toArray(t=[],n=0){return t[n]=this.x,t[n+1]=this.y,t[n+2]=this.z,t[n+3]=this.w,t}fromBufferAttribute(t,n){return this.x=t.getX(n),this.y=t.getY(n),this.z=t.getZ(n),this.w=t.getW(n),this}random(){return this.x=Math.random(),this.y=Math.random(),this.z=Math.random(),this.w=Math.random(),this}*[Symbol.iterator](){yield this.x,yield this.y,yield this.z,yield this.w}},Kh=class extends qa{constructor(t=1,n=1,i={}){super(),i=Object.assign({generateMipmaps:!1,internalFormat:null,minFilter:En,depthBuffer:!0,stencilBuffer:!1,resolveDepthBuffer:!0,resolveStencilBuffer:!0,depthTexture:null,samples:0,count:1,depth:1,multiview:!1},i),this.isRenderTarget=!0,this.width=t,this.height=n,this.depth=i.depth,this.scissor=new Ge(0,0,t,n),this.scissorTest=!1,this.viewport=new Ge(0,0,t,n),this.textures=[];let a={width:t,height:n,depth:i.depth},s=new qn(a),r=i.count;for(let o=0;o<r;o++)this.textures[o]=s.clone(),this.textures[o].isRenderTargetTexture=!0,this.textures[o].renderTarget=this;this._setTextureOptions(i),this.depthBuffer=i.depthBuffer,this.stencilBuffer=i.stencilBuffer,this.resolveDepthBuffer=i.resolveDepthBuffer,this.resolveStencilBuffer=i.resolveStencilBuffer,this._depthTexture=null,this.depthTexture=i.depthTexture,this.samples=i.samples,this.multiview=i.multiview}_setTextureOptions(t={}){let n={minFilter:En,generateMipmaps:!1,flipY:!1,internalFormat:null};t.mapping!==void 0&&(n.mapping=t.mapping),t.wrapS!==void 0&&(n.wrapS=t.wrapS),t.wrapT!==void 0&&(n.wrapT=t.wrapT),t.wrapR!==void 0&&(n.wrapR=t.wrapR),t.magFilter!==void 0&&(n.magFilter=t.magFilter),t.minFilter!==void 0&&(n.minFilter=t.minFilter),t.format!==void 0&&(n.format=t.format),t.type!==void 0&&(n.type=t.type),t.anisotropy!==void 0&&(n.anisotropy=t.anisotropy),t.colorSpace!==void 0&&(n.colorSpace=t.colorSpace),t.flipY!==void 0&&(n.flipY=t.flipY),t.generateMipmaps!==void 0&&(n.generateMipmaps=t.generateMipmaps),t.internalFormat!==void 0&&(n.internalFormat=t.internalFormat);for(let i=0;i<this.textures.length;i++)this.textures[i].setValues(n)}get texture(){return this.textures[0]}set texture(t){this.textures[0]=t}set depthTexture(t){this._depthTexture!==null&&(this._depthTexture.renderTarget=null),t!==null&&(t.renderTarget=this),this._depthTexture=t}get depthTexture(){return this._depthTexture}setSize(t,n,i=1){if(this.width!==t||this.height!==n||this.depth!==i){this.width=t,this.height=n,this.depth=i;for(let a=0,s=this.textures.length;a<s;a++)this.textures[a].image.width=t,this.textures[a].image.height=n,this.textures[a].image.depth=i,this.textures[a].isData3DTexture!==!0&&(this.textures[a].isArrayTexture=this.textures[a].image.depth>1);this.dispose()}this.viewport.set(0,0,t,n),this.scissor.set(0,0,t,n)}clone(){return new this.constructor().copy(this)}copy(t){this.width=t.width,this.height=t.height,this.depth=t.depth,this.scissor.copy(t.scissor),this.scissorTest=t.scissorTest,this.viewport.copy(t.viewport),this.textures.length=0;for(let n=0,i=t.textures.length;n<i;n++){this.textures[n]=t.textures[n].clone(),this.textures[n].isRenderTargetTexture=!0,this.textures[n].renderTarget=this;let a=Object.assign({},t.textures[n].image);this.textures[n].source=new sl(a)}return this.depthBuffer=t.depthBuffer,this.stencilBuffer=t.stencilBuffer,this.resolveDepthBuffer=t.resolveDepthBuffer,this.resolveStencilBuffer=t.resolveStencilBuffer,t.depthTexture!==null&&(this.depthTexture=t.depthTexture.clone()),this.samples=t.samples,this}dispose(){this.dispatchEvent({type:"dispose"})}},gi=class extends Kh{constructor(t=1,n=1,i={}){super(t,n,i),this.isWebGLRenderTarget=!0}},kc=class extends qn{constructor(t=null,n=1,i=1,a=1){super(null),this.isDataArrayTexture=!0,this.image={data:t,width:n,height:i,depth:a},this.magFilter=pn,this.minFilter=pn,this.wrapR=ha,this.generateMipmaps=!1,this.flipY=!1,this.unpackAlignment=1,this.layerUpdates=new Set}addLayerUpdate(t){this.layerUpdates.add(t)}clearLayerUpdates(){this.layerUpdates.clear()}};var Qh=class extends qn{constructor(t=null,n=1,i=1,a=1){super(null),this.isData3DTexture=!0,this.image={data:t,width:n,height:i,depth:a},this.magFilter=pn,this.minFilter=pn,this.wrapR=ha,this.generateMipmaps=!1,this.flipY=!1,this.unpackAlignment=1}};var Ie=class e{constructor(t,n,i,a,s,r,o,l,c,f,p,u,d,_,b,v){e.prototype.isMatrix4=!0,this.elements=[1,0,0,0,0,1,0,0,0,0,1,0,0,0,0,1],t!==void 0&&this.set(t,n,i,a,s,r,o,l,c,f,p,u,d,_,b,v)}set(t,n,i,a,s,r,o,l,c,f,p,u,d,_,b,v){let h=this.elements;return h[0]=t,h[4]=n,h[8]=i,h[12]=a,h[1]=s,h[5]=r,h[9]=o,h[13]=l,h[2]=c,h[6]=f,h[10]=p,h[14]=u,h[3]=d,h[7]=_,h[11]=b,h[15]=v,this}identity(){return this.set(1,0,0,0,0,1,0,0,0,0,1,0,0,0,0,1),this}clone(){return new e().fromArray(this.elements)}copy(t){let n=this.elements,i=t.elements;return n[0]=i[0],n[1]=i[1],n[2]=i[2],n[3]=i[3],n[4]=i[4],n[5]=i[5],n[6]=i[6],n[7]=i[7],n[8]=i[8],n[9]=i[9],n[10]=i[10],n[11]=i[11],n[12]=i[12],n[13]=i[13],n[14]=i[14],n[15]=i[15],this}copyPosition(t){let n=this.elements,i=t.elements;return n[12]=i[12],n[13]=i[13],n[14]=i[14],this}setFromMatrix3(t){let n=t.elements;return this.set(n[0],n[3],n[6],0,n[1],n[4],n[7],0,n[2],n[5],n[8],0,0,0,0,1),this}extractBasis(t,n,i){return this.determinant()===0?(t.set(1,0,0),n.set(0,1,0),i.set(0,0,1),this):(t.setFromMatrixColumn(this,0),n.setFromMatrixColumn(this,1),i.setFromMatrixColumn(this,2),this)}makeBasis(t,n,i){return this.set(t.x,n.x,i.x,0,t.y,n.y,i.y,0,t.z,n.z,i.z,0,0,0,0,1),this}extractRotation(t){if(t.determinant()===0)return this.identity();let n=this.elements,i=t.elements,a=1/Go.setFromMatrixColumn(t,0).length(),s=1/Go.setFromMatrixColumn(t,1).length(),r=1/Go.setFromMatrixColumn(t,2).length();return n[0]=i[0]*a,n[1]=i[1]*a,n[2]=i[2]*a,n[3]=0,n[4]=i[4]*s,n[5]=i[5]*s,n[6]=i[6]*s,n[7]=0,n[8]=i[8]*r,n[9]=i[9]*r,n[10]=i[10]*r,n[11]=0,n[12]=0,n[13]=0,n[14]=0,n[15]=1,this}makeRotationFromEuler(t){let n=this.elements,i=t.x,a=t.y,s=t.z,r=Math.cos(i),o=Math.sin(i),l=Math.cos(a),c=Math.sin(a),f=Math.cos(s),p=Math.sin(s);if(t.order==="XYZ"){let u=r*f,d=r*p,_=o*f,b=o*p;n[0]=l*f,n[4]=-l*p,n[8]=c,n[1]=d+_*c,n[5]=u-b*c,n[9]=-o*l,n[2]=b-u*c,n[6]=_+d*c,n[10]=r*l}else if(t.order==="YXZ"){let u=l*f,d=l*p,_=c*f,b=c*p;n[0]=u+b*o,n[4]=_*o-d,n[8]=r*c,n[1]=r*p,n[5]=r*f,n[9]=-o,n[2]=d*o-_,n[6]=b+u*o,n[10]=r*l}else if(t.order==="ZXY"){let u=l*f,d=l*p,_=c*f,b=c*p;n[0]=u-b*o,n[4]=-r*p,n[8]=_+d*o,n[1]=d+_*o,n[5]=r*f,n[9]=b-u*o,n[2]=-r*c,n[6]=o,n[10]=r*l}else if(t.order==="ZYX"){let u=r*f,d=r*p,_=o*f,b=o*p;n[0]=l*f,n[4]=_*c-d,n[8]=u*c+b,n[1]=l*p,n[5]=b*c+u,n[9]=d*c-_,n[2]=-c,n[6]=o*l,n[10]=r*l}else if(t.order==="YZX"){let u=r*l,d=r*c,_=o*l,b=o*c;n[0]=l*f,n[4]=b-u*p,n[8]=_*p+d,n[1]=p,n[5]=r*f,n[9]=-o*f,n[2]=-c*f,n[6]=d*p+_,n[10]=u-b*p}else if(t.order==="XZY"){let u=r*l,d=r*c,_=o*l,b=o*c;n[0]=l*f,n[4]=-p,n[8]=c*f,n[1]=u*p+b,n[5]=r*f,n[9]=d*p-_,n[2]=_*p-d,n[6]=o*f,n[10]=b*p+u}return n[3]=0,n[7]=0,n[11]=0,n[12]=0,n[13]=0,n[14]=0,n[15]=1,this}makeRotationFromQuaternion(t){return this.compose(fC,t,hC)}lookAt(t,n,i){let a=this.elements;return pi.subVectors(t,n),pi.lengthSq()===0&&(pi.z=1),pi.normalize(),Ns.crossVectors(i,pi),Ns.lengthSq()===0&&(Math.abs(i.z)===1?pi.x+=1e-4:pi.z+=1e-4,pi.normalize(),Ns.crossVectors(i,pi)),Ns.normalize(),dh.crossVectors(pi,Ns),a[0]=Ns.x,a[4]=dh.x,a[8]=pi.x,a[1]=Ns.y,a[5]=dh.y,a[9]=pi.y,a[2]=Ns.z,a[6]=dh.z,a[10]=pi.z,this}multiply(t){return this.multiplyMatrices(this,t)}premultiply(t){return this.multiplyMatrices(t,this)}multiplyMatrices(t,n){let i=t.elements,a=n.elements,s=this.elements,r=i[0],o=i[4],l=i[8],c=i[12],f=i[1],p=i[5],u=i[9],d=i[13],_=i[2],b=i[6],v=i[10],h=i[14],m=i[3],g=i[7],x=i[11],E=i[15],A=a[0],w=a[4],y=a[8],T=a[12],L=a[1],R=a[5],N=a[9],P=a[13],z=a[2],F=a[6],V=a[10],X=a[14],at=a[3],et=a[7],dt=a[11],St=a[15];return s[0]=r*A+o*L+l*z+c*at,s[4]=r*w+o*R+l*F+c*et,s[8]=r*y+o*N+l*V+c*dt,s[12]=r*T+o*P+l*X+c*St,s[1]=f*A+p*L+u*z+d*at,s[5]=f*w+p*R+u*F+d*et,s[9]=f*y+p*N+u*V+d*dt,s[13]=f*T+p*P+u*X+d*St,s[2]=_*A+b*L+v*z+h*at,s[6]=_*w+b*R+v*F+h*et,s[10]=_*y+b*N+v*V+h*dt,s[14]=_*T+b*P+v*X+h*St,s[3]=m*A+g*L+x*z+E*at,s[7]=m*w+g*R+x*F+E*et,s[11]=m*y+g*N+x*V+E*dt,s[15]=m*T+g*P+x*X+E*St,this}multiplyScalar(t){let n=this.elements;return n[0]*=t,n[4]*=t,n[8]*=t,n[12]*=t,n[1]*=t,n[5]*=t,n[9]*=t,n[13]*=t,n[2]*=t,n[6]*=t,n[10]*=t,n[14]*=t,n[3]*=t,n[7]*=t,n[11]*=t,n[15]*=t,this}determinant(){let t=this.elements,n=t[0],i=t[4],a=t[8],s=t[12],r=t[1],o=t[5],l=t[9],c=t[13],f=t[2],p=t[6],u=t[10],d=t[14],_=t[3],b=t[7],v=t[11],h=t[15],m=l*d-c*u,g=o*d-c*p,x=o*u-l*p,E=r*d-c*f,A=r*u-l*f,w=r*p-o*f;return n*(b*m-v*g+h*x)-i*(_*m-v*E+h*A)+a*(_*g-b*E+h*w)-s*(_*x-b*A+v*w)}transpose(){let t=this.elements,n;return n=t[1],t[1]=t[4],t[4]=n,n=t[2],t[2]=t[8],t[8]=n,n=t[6],t[6]=t[9],t[9]=n,n=t[3],t[3]=t[12],t[12]=n,n=t[7],t[7]=t[13],t[13]=n,n=t[11],t[11]=t[14],t[14]=n,this}setPosition(t,n,i){let a=this.elements;return t.isVector3?(a[12]=t.x,a[13]=t.y,a[14]=t.z):(a[12]=t,a[13]=n,a[14]=i),this}invert(){let t=this.elements,n=t[0],i=t[1],a=t[2],s=t[3],r=t[4],o=t[5],l=t[6],c=t[7],f=t[8],p=t[9],u=t[10],d=t[11],_=t[12],b=t[13],v=t[14],h=t[15],m=n*o-i*r,g=n*l-a*r,x=n*c-s*r,E=i*l-a*o,A=i*c-s*o,w=a*c-s*l,y=f*b-p*_,T=f*v-u*_,L=f*h-d*_,R=p*v-u*b,N=p*h-d*b,P=u*h-d*v,z=m*P-g*N+x*R+E*L-A*T+w*y;if(z===0)return this.set(0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0);let F=1/z;return t[0]=(o*P-l*N+c*R)*F,t[1]=(a*N-i*P-s*R)*F,t[2]=(b*w-v*A+h*E)*F,t[3]=(u*A-p*w-d*E)*F,t[4]=(l*L-r*P-c*T)*F,t[5]=(n*P-a*L+s*T)*F,t[6]=(v*x-_*w-h*g)*F,t[7]=(f*w-u*x+d*g)*F,t[8]=(r*N-o*L+c*y)*F,t[9]=(i*L-n*N-s*y)*F,t[10]=(_*A-b*x+h*m)*F,t[11]=(p*x-f*A-d*m)*F,t[12]=(o*T-r*R-l*y)*F,t[13]=(n*R-i*T+a*y)*F,t[14]=(b*g-_*E-v*m)*F,t[15]=(f*E-p*g+u*m)*F,this}scale(t){let n=this.elements,i=t.x,a=t.y,s=t.z;return n[0]*=i,n[4]*=a,n[8]*=s,n[1]*=i,n[5]*=a,n[9]*=s,n[2]*=i,n[6]*=a,n[10]*=s,n[3]*=i,n[7]*=a,n[11]*=s,this}getMaxScaleOnAxis(){let t=this.elements,n=t[0]*t[0]+t[1]*t[1]+t[2]*t[2],i=t[4]*t[4]+t[5]*t[5]+t[6]*t[6],a=t[8]*t[8]+t[9]*t[9]+t[10]*t[10];return Math.sqrt(Math.max(n,i,a))}makeTranslation(t,n,i){return t.isVector3?this.set(1,0,0,t.x,0,1,0,t.y,0,0,1,t.z,0,0,0,1):this.set(1,0,0,t,0,1,0,n,0,0,1,i,0,0,0,1),this}makeRotationX(t){let n=Math.cos(t),i=Math.sin(t);return this.set(1,0,0,0,0,n,-i,0,0,i,n,0,0,0,0,1),this}makeRotationY(t){let n=Math.cos(t),i=Math.sin(t);return this.set(n,0,i,0,0,1,0,0,-i,0,n,0,0,0,0,1),this}makeRotationZ(t){let n=Math.cos(t),i=Math.sin(t);return this.set(n,-i,0,0,i,n,0,0,0,0,1,0,0,0,0,1),this}makeRotationAxis(t,n){let i=Math.cos(n),a=Math.sin(n),s=1-i,r=t.x,o=t.y,l=t.z,c=s*r,f=s*o;return this.set(c*r+i,c*o-a*l,c*l+a*o,0,c*o+a*l,f*o+i,f*l-a*r,0,c*l-a*o,f*l+a*r,s*l*l+i,0,0,0,0,1),this}makeScale(t,n,i){return this.set(t,0,0,0,0,n,0,0,0,0,i,0,0,0,0,1),this}makeShear(t,n,i,a,s,r){return this.set(1,i,s,0,t,1,r,0,n,a,1,0,0,0,0,1),this}compose(t,n,i){let a=this.elements,s=n._x,r=n._y,o=n._z,l=n._w,c=s+s,f=r+r,p=o+o,u=s*c,d=s*f,_=s*p,b=r*f,v=r*p,h=o*p,m=l*c,g=l*f,x=l*p,E=i.x,A=i.y,w=i.z;return a[0]=(1-(b+h))*E,a[1]=(d+x)*E,a[2]=(_-g)*E,a[3]=0,a[4]=(d-x)*A,a[5]=(1-(u+h))*A,a[6]=(v+m)*A,a[7]=0,a[8]=(_+g)*w,a[9]=(v-m)*w,a[10]=(1-(u+b))*w,a[11]=0,a[12]=t.x,a[13]=t.y,a[14]=t.z,a[15]=1,this}decompose(t,n,i){let a=this.elements;t.x=a[12],t.y=a[13],t.z=a[14];let s=this.determinant();if(s===0)return i.set(1,1,1),n.identity(),this;let r=Go.set(a[0],a[1],a[2]).length(),o=Go.set(a[4],a[5],a[6]).length(),l=Go.set(a[8],a[9],a[10]).length();s<0&&(r=-r),Vi.copy(this);let c=1/r,f=1/o,p=1/l;return Vi.elements[0]*=c,Vi.elements[1]*=c,Vi.elements[2]*=c,Vi.elements[4]*=f,Vi.elements[5]*=f,Vi.elements[6]*=f,Vi.elements[8]*=p,Vi.elements[9]*=p,Vi.elements[10]*=p,n.setFromRotationMatrix(Vi),i.x=r,i.y=o,i.z=l,this}makePerspective(t,n,i,a,s,r,o=ki,l=!1){let c=this.elements,f=2*s/(n-t),p=2*s/(i-a),u=(n+t)/(n-t),d=(i+a)/(i-a),_,b;if(l)_=s/(r-s),b=r*s/(r-s);else if(o===ki)_=-(r+s)/(r-s),b=-2*r*s/(r-s);else if(o===nl)_=-r/(r-s),b=-r*s/(r-s);else throw new Error("THREE.Matrix4.makePerspective(): Invalid coordinate system: "+o);return c[0]=f,c[4]=0,c[8]=u,c[12]=0,c[1]=0,c[5]=p,c[9]=d,c[13]=0,c[2]=0,c[6]=0,c[10]=_,c[14]=b,c[3]=0,c[7]=0,c[11]=-1,c[15]=0,this}makeOrthographic(t,n,i,a,s,r,o=ki,l=!1){let c=this.elements,f=2/(n-t),p=2/(i-a),u=-(n+t)/(n-t),d=-(i+a)/(i-a),_,b;if(l)_=1/(r-s),b=r/(r-s);else if(o===ki)_=-2/(r-s),b=-(r+s)/(r-s);else if(o===nl)_=-1/(r-s),b=-s/(r-s);else throw new Error("THREE.Matrix4.makeOrthographic(): Invalid coordinate system: "+o);return c[0]=f,c[4]=0,c[8]=0,c[12]=u,c[1]=0,c[5]=p,c[9]=0,c[13]=d,c[2]=0,c[6]=0,c[10]=_,c[14]=b,c[3]=0,c[7]=0,c[11]=0,c[15]=1,this}equals(t){let n=this.elements,i=t.elements;for(let a=0;a<16;a++)if(n[a]!==i[a])return!1;return!0}fromArray(t,n=0){for(let i=0;i<16;i++)this.elements[i]=t[i+n];return this}toArray(t=[],n=0){let i=this.elements;return t[n]=i[0],t[n+1]=i[1],t[n+2]=i[2],t[n+3]=i[3],t[n+4]=i[4],t[n+5]=i[5],t[n+6]=i[6],t[n+7]=i[7],t[n+8]=i[8],t[n+9]=i[9],t[n+10]=i[10],t[n+11]=i[11],t[n+12]=i[12],t[n+13]=i[13],t[n+14]=i[14],t[n+15]=i[15],t}},Go=new B,Vi=new Ie,fC=new B(0,0,0),hC=new B(1,1,1),Ns=new B,dh=new B,pi=new B,yM=new Ie,xM=new pa,Xi=class e{constructor(t=0,n=0,i=0,a=e.DEFAULT_ORDER){this.isEuler=!0,this._x=t,this._y=n,this._z=i,this._order=a}get x(){return this._x}set x(t){this._x=t,this._onChangeCallback()}get y(){return this._y}set y(t){this._y=t,this._onChangeCallback()}get z(){return this._z}set z(t){this._z=t,this._onChangeCallback()}get order(){return this._order}set order(t){this._order=t,this._onChangeCallback()}set(t,n,i,a=this._order){return this._x=t,this._y=n,this._z=i,this._order=a,this._onChangeCallback(),this}clone(){return new this.constructor(this._x,this._y,this._z,this._order)}copy(t){return this._x=t._x,this._y=t._y,this._z=t._z,this._order=t._order,this._onChangeCallback(),this}setFromRotationMatrix(t,n=this._order,i=!0){let a=t.elements,s=a[0],r=a[4],o=a[8],l=a[1],c=a[5],f=a[9],p=a[2],u=a[6],d=a[10];switch(n){case"XYZ":this._y=Math.asin(ne(o,-1,1)),Math.abs(o)<.9999999?(this._x=Math.atan2(-f,d),this._z=Math.atan2(-r,s)):(this._x=Math.atan2(u,c),this._z=0);break;case"YXZ":this._x=Math.asin(-ne(f,-1,1)),Math.abs(f)<.9999999?(this._y=Math.atan2(o,d),this._z=Math.atan2(l,c)):(this._y=Math.atan2(-p,s),this._z=0);break;case"ZXY":this._x=Math.asin(ne(u,-1,1)),Math.abs(u)<.9999999?(this._y=Math.atan2(-p,d),this._z=Math.atan2(-r,c)):(this._y=0,this._z=Math.atan2(l,s));break;case"ZYX":this._y=Math.asin(-ne(p,-1,1)),Math.abs(p)<.9999999?(this._x=Math.atan2(u,d),this._z=Math.atan2(l,s)):(this._x=0,this._z=Math.atan2(-r,c));break;case"YZX":this._z=Math.asin(ne(l,-1,1)),Math.abs(l)<.9999999?(this._x=Math.atan2(-f,c),this._y=Math.atan2(-p,s)):(this._x=0,this._y=Math.atan2(o,d));break;case"XZY":this._z=Math.asin(-ne(r,-1,1)),Math.abs(r)<.9999999?(this._x=Math.atan2(u,c),this._y=Math.atan2(o,s)):(this._x=Math.atan2(-f,d),this._y=0);break;default:Ot("Euler: .setFromRotationMatrix() encountered an unknown order: "+n)}return this._order=n,i===!0&&this._onChangeCallback(),this}setFromQuaternion(t,n,i){return yM.makeRotationFromQuaternion(t),this.setFromRotationMatrix(yM,n,i)}setFromVector3(t,n=this._order){return this.set(t.x,t.y,t.z,n)}reorder(t){return xM.setFromEuler(this),this.setFromQuaternion(xM,t)}equals(t){return t._x===this._x&&t._y===this._y&&t._z===this._z&&t._order===this._order}fromArray(t){return this._x=t[0],this._y=t[1],this._z=t[2],t[3]!==void 0&&(this._order=t[3]),this._onChangeCallback(),this}toArray(t=[],n=0){return t[n]=this._x,t[n+1]=this._y,t[n+2]=this._z,t[n+3]=this._order,t}_onChange(t){return this._onChangeCallback=t,this}_onChangeCallback(){}*[Symbol.iterator](){yield this._x,yield this._y,yield this._z,yield this._order}};Xi.DEFAULT_ORDER="XYZ";var Xc=class{constructor(){this.mask=1}set(t){this.mask=(1<<t|0)>>>0}enable(t){this.mask|=1<<t|0}enableAll(){this.mask=-1}toggle(t){this.mask^=1<<t|0}disable(t){this.mask&=~(1<<t|0)}disableAll(){this.mask=0}test(t){return(this.mask&t.mask)!==0}isEnabled(t){return(this.mask&(1<<t|0))!==0}},dC=0,SM=new B,ko=new pa,Fa=new Ie,ph=new B,Dc=new B,pC=new B,mC=new pa,bM=new B(1,0,0),MM=new B(0,1,0),TM=new B(0,0,1),EM={type:"added"},gC={type:"removed"},Xo={type:"childadded",child:null},A0={type:"childremoved",child:null},tn=class e extends qa{constructor(){super(),this.isObject3D=!0,Object.defineProperty(this,"id",{value:dC++}),this.uuid=gl(),this.name="",this.type="Object3D",this.parent=null,this.children=[],this.up=e.DEFAULT_UP.clone();let t=new B,n=new Xi,i=new pa,a=new B(1,1,1);function s(){i.setFromEuler(n,!1)}function r(){n.setFromQuaternion(i,void 0,!1)}n._onChange(s),i._onChange(r),Object.defineProperties(this,{position:{configurable:!0,enumerable:!0,value:t},rotation:{configurable:!0,enumerable:!0,value:n},quaternion:{configurable:!0,enumerable:!0,value:i},scale:{configurable:!0,enumerable:!0,value:a},modelViewMatrix:{value:new Ie},normalMatrix:{value:new Zt}}),this.matrix=new Ie,this.matrixWorld=new Ie,this.matrixAutoUpdate=e.DEFAULT_MATRIX_AUTO_UPDATE,this.matrixWorldAutoUpdate=e.DEFAULT_MATRIX_WORLD_AUTO_UPDATE,this.matrixWorldNeedsUpdate=!1,this.layers=new Xc,this.visible=!0,this.castShadow=!1,this.receiveShadow=!1,this.frustumCulled=!0,this.renderOrder=0,this.animations=[],this.customDepthMaterial=void 0,this.customDistanceMaterial=void 0,this.static=!1,this.userData={},this.pivot=null}onBeforeShadow(){}onAfterShadow(){}onBeforeRender(){}onAfterRender(){}applyMatrix4(t){this.matrixAutoUpdate&&this.updateMatrix(),this.matrix.premultiply(t),this.matrix.decompose(this.position,this.quaternion,this.scale)}applyQuaternion(t){return this.quaternion.premultiply(t),this}setRotationFromAxisAngle(t,n){this.quaternion.setFromAxisAngle(t,n)}setRotationFromEuler(t){this.quaternion.setFromEuler(t,!0)}setRotationFromMatrix(t){this.quaternion.setFromRotationMatrix(t)}setRotationFromQuaternion(t){this.quaternion.copy(t)}rotateOnAxis(t,n){return ko.setFromAxisAngle(t,n),this.quaternion.multiply(ko),this}rotateOnWorldAxis(t,n){return ko.setFromAxisAngle(t,n),this.quaternion.premultiply(ko),this}rotateX(t){return this.rotateOnAxis(bM,t)}rotateY(t){return this.rotateOnAxis(MM,t)}rotateZ(t){return this.rotateOnAxis(TM,t)}translateOnAxis(t,n){return SM.copy(t).applyQuaternion(this.quaternion),this.position.add(SM.multiplyScalar(n)),this}translateX(t){return this.translateOnAxis(bM,t)}translateY(t){return this.translateOnAxis(MM,t)}translateZ(t){return this.translateOnAxis(TM,t)}localToWorld(t){return this.updateWorldMatrix(!0,!1),t.applyMatrix4(this.matrixWorld)}worldToLocal(t){return this.updateWorldMatrix(!0,!1),t.applyMatrix4(Fa.copy(this.matrixWorld).invert())}lookAt(t,n,i){t.isVector3?ph.copy(t):ph.set(t,n,i);let a=this.parent;this.updateWorldMatrix(!0,!1),Dc.setFromMatrixPosition(this.matrixWorld),this.isCamera||this.isLight?Fa.lookAt(Dc,ph,this.up):Fa.lookAt(ph,Dc,this.up),this.quaternion.setFromRotationMatrix(Fa),a&&(Fa.extractRotation(a.matrixWorld),ko.setFromRotationMatrix(Fa),this.quaternion.premultiply(ko.invert()))}add(t){if(arguments.length>1){for(let n=0;n<arguments.length;n++)this.add(arguments[n]);return this}return t===this?(zt("Object3D.add: object can't be added as a child of itself.",t),this):(t&&t.isObject3D?(t.removeFromParent(),t.parent=this,this.children.push(t),t.dispatchEvent(EM),Xo.child=t,this.dispatchEvent(Xo),Xo.child=null):zt("Object3D.add: object not an instance of THREE.Object3D.",t),this)}remove(t){if(arguments.length>1){for(let i=0;i<arguments.length;i++)this.remove(arguments[i]);return this}let n=this.children.indexOf(t);return n!==-1&&(t.parent=null,this.children.splice(n,1),t.dispatchEvent(gC),A0.child=t,this.dispatchEvent(A0),A0.child=null),this}removeFromParent(){let t=this.parent;return t!==null&&t.remove(this),this}clear(){return this.remove(...this.children)}attach(t){return this.updateWorldMatrix(!0,!1),Fa.copy(this.matrixWorld).invert(),t.parent!==null&&(t.parent.updateWorldMatrix(!0,!1),Fa.multiply(t.parent.matrixWorld)),t.applyMatrix4(Fa),t.removeFromParent(),t.parent=this,this.children.push(t),t.updateWorldMatrix(!1,!0),t.dispatchEvent(EM),Xo.child=t,this.dispatchEvent(Xo),Xo.child=null,this}getObjectById(t){return this.getObjectByProperty("id",t)}getObjectByName(t){return this.getObjectByProperty("name",t)}getObjectByProperty(t,n){if(this[t]===n)return this;for(let i=0,a=this.children.length;i<a;i++){let r=this.children[i].getObjectByProperty(t,n);if(r!==void 0)return r}}getObjectsByProperty(t,n,i=[]){this[t]===n&&i.push(this);let a=this.children;for(let s=0,r=a.length;s<r;s++)a[s].getObjectsByProperty(t,n,i);return i}getWorldPosition(t){return this.updateWorldMatrix(!0,!1),t.setFromMatrixPosition(this.matrixWorld)}getWorldQuaternion(t){return this.updateWorldMatrix(!0,!1),this.matrixWorld.decompose(Dc,t,pC),t}getWorldScale(t){return this.updateWorldMatrix(!0,!1),this.matrixWorld.decompose(Dc,mC,t),t}getWorldDirection(t){this.updateWorldMatrix(!0,!1);let n=this.matrixWorld.elements;return t.set(n[8],n[9],n[10]).normalize()}raycast(){}traverse(t){t(this);let n=this.children;for(let i=0,a=n.length;i<a;i++)n[i].traverse(t)}traverseVisible(t){if(this.visible===!1)return;t(this);let n=this.children;for(let i=0,a=n.length;i<a;i++)n[i].traverseVisible(t)}traverseAncestors(t){let n=this.parent;n!==null&&(t(n),n.traverseAncestors(t))}updateMatrix(){this.matrix.compose(this.position,this.quaternion,this.scale);let t=this.pivot;if(t!==null){let n=t.x,i=t.y,a=t.z,s=this.matrix.elements;s[12]+=n-s[0]*n-s[4]*i-s[8]*a,s[13]+=i-s[1]*n-s[5]*i-s[9]*a,s[14]+=a-s[2]*n-s[6]*i-s[10]*a}this.matrixWorldNeedsUpdate=!0}updateMatrixWorld(t){this.matrixAutoUpdate&&this.updateMatrix(),(this.matrixWorldNeedsUpdate||t)&&(this.matrixWorldAutoUpdate===!0&&(this.parent===null?this.matrixWorld.copy(this.matrix):this.matrixWorld.multiplyMatrices(this.parent.matrixWorld,this.matrix)),this.matrixWorldNeedsUpdate=!1,t=!0);let n=this.children;for(let i=0,a=n.length;i<a;i++)n[i].updateMatrixWorld(t)}updateWorldMatrix(t,n){let i=this.parent;if(t===!0&&i!==null&&i.updateWorldMatrix(!0,!1),this.matrixAutoUpdate&&this.updateMatrix(),this.matrixWorldAutoUpdate===!0&&(this.parent===null?this.matrixWorld.copy(this.matrix):this.matrixWorld.multiplyMatrices(this.parent.matrixWorld,this.matrix)),n===!0){let a=this.children;for(let s=0,r=a.length;s<r;s++)a[s].updateWorldMatrix(!1,!0)}}toJSON(t){let n=t===void 0||typeof t=="string",i={};n&&(t={geometries:{},materials:{},textures:{},images:{},shapes:{},skeletons:{},animations:{},nodes:{}},i.metadata={version:4.7,type:"Object",generator:"Object3D.toJSON"});let a={};a.uuid=this.uuid,a.type=this.type,this.name!==""&&(a.name=this.name),this.castShadow===!0&&(a.castShadow=!0),this.receiveShadow===!0&&(a.receiveShadow=!0),this.visible===!1&&(a.visible=!1),this.frustumCulled===!1&&(a.frustumCulled=!1),this.renderOrder!==0&&(a.renderOrder=this.renderOrder),this.static!==!1&&(a.static=this.static),Object.keys(this.userData).length>0&&(a.userData=this.userData),a.layers=this.layers.mask,a.matrix=this.matrix.toArray(),a.up=this.up.toArray(),this.pivot!==null&&(a.pivot=this.pivot.toArray()),this.matrixAutoUpdate===!1&&(a.matrixAutoUpdate=!1),this.morphTargetDictionary!==void 0&&(a.morphTargetDictionary=Object.assign({},this.morphTargetDictionary)),this.morphTargetInfluences!==void 0&&(a.morphTargetInfluences=this.morphTargetInfluences.slice()),this.isInstancedMesh&&(a.type="InstancedMesh",a.count=this.count,a.instanceMatrix=this.instanceMatrix.toJSON(),this.instanceColor!==null&&(a.instanceColor=this.instanceColor.toJSON())),this.isBatchedMesh&&(a.type="BatchedMesh",a.perObjectFrustumCulled=this.perObjectFrustumCulled,a.sortObjects=this.sortObjects,a.drawRanges=this._drawRanges,a.reservedRanges=this._reservedRanges,a.geometryInfo=this._geometryInfo.map(o=>({...o,boundingBox:o.boundingBox?o.boundingBox.toJSON():void 0,boundingSphere:o.boundingSphere?o.boundingSphere.toJSON():void 0})),a.instanceInfo=this._instanceInfo.map(o=>({...o})),a.availableInstanceIds=this._availableInstanceIds.slice(),a.availableGeometryIds=this._availableGeometryIds.slice(),a.nextIndexStart=this._nextIndexStart,a.nextVertexStart=this._nextVertexStart,a.geometryCount=this._geometryCount,a.maxInstanceCount=this._maxInstanceCount,a.maxVertexCount=this._maxVertexCount,a.maxIndexCount=this._maxIndexCount,a.geometryInitialized=this._geometryInitialized,a.matricesTexture=this._matricesTexture.toJSON(t),a.indirectTexture=this._indirectTexture.toJSON(t),this._colorsTexture!==null&&(a.colorsTexture=this._colorsTexture.toJSON(t)),this.boundingSphere!==null&&(a.boundingSphere=this.boundingSphere.toJSON()),this.boundingBox!==null&&(a.boundingBox=this.boundingBox.toJSON()));function s(o,l){return o[l.uuid]===void 0&&(o[l.uuid]=l.toJSON(t)),l.uuid}if(this.isScene)this.background&&(this.background.isColor?a.background=this.background.toJSON():this.background.isTexture&&(a.background=this.background.toJSON(t).uuid)),this.environment&&this.environment.isTexture&&this.environment.isRenderTargetTexture!==!0&&(a.environment=this.environment.toJSON(t).uuid);else if(this.isMesh||this.isLine||this.isPoints){a.geometry=s(t.geometries,this.geometry);let o=this.geometry.parameters;if(o!==void 0&&o.shapes!==void 0){let l=o.shapes;if(Array.isArray(l))for(let c=0,f=l.length;c<f;c++){let p=l[c];s(t.shapes,p)}else s(t.shapes,l)}}if(this.isSkinnedMesh&&(a.bindMode=this.bindMode,a.bindMatrix=this.bindMatrix.toArray(),this.skeleton!==void 0&&(s(t.skeletons,this.skeleton),a.skeleton=this.skeleton.uuid)),this.material!==void 0)if(Array.isArray(this.material)){let o=[];for(let l=0,c=this.material.length;l<c;l++)o.push(s(t.materials,this.material[l]));a.material=o}else a.material=s(t.materials,this.material);if(this.children.length>0){a.children=[];for(let o=0;o<this.children.length;o++)a.children.push(this.children[o].toJSON(t).object)}if(this.animations.length>0){a.animations=[];for(let o=0;o<this.animations.length;o++){let l=this.animations[o];a.animations.push(s(t.animations,l))}}if(n){let o=r(t.geometries),l=r(t.materials),c=r(t.textures),f=r(t.images),p=r(t.shapes),u=r(t.skeletons),d=r(t.animations),_=r(t.nodes);o.length>0&&(i.geometries=o),l.length>0&&(i.materials=l),c.length>0&&(i.textures=c),f.length>0&&(i.images=f),p.length>0&&(i.shapes=p),u.length>0&&(i.skeletons=u),d.length>0&&(i.animations=d),_.length>0&&(i.nodes=_)}return i.object=a,i;function r(o){let l=[];for(let c in o){let f=o[c];delete f.metadata,l.push(f)}return l}}clone(t){return new this.constructor().copy(this,t)}copy(t,n=!0){if(this.name=t.name,this.up.copy(t.up),this.position.copy(t.position),this.rotation.order=t.rotation.order,this.quaternion.copy(t.quaternion),this.scale.copy(t.scale),t.pivot!==null&&(this.pivot=t.pivot.clone()),this.matrix.copy(t.matrix),this.matrixWorld.copy(t.matrixWorld),this.matrixAutoUpdate=t.matrixAutoUpdate,this.matrixWorldAutoUpdate=t.matrixWorldAutoUpdate,this.matrixWorldNeedsUpdate=t.matrixWorldNeedsUpdate,this.layers.mask=t.layers.mask,this.visible=t.visible,this.castShadow=t.castShadow,this.receiveShadow=t.receiveShadow,this.frustumCulled=t.frustumCulled,this.renderOrder=t.renderOrder,this.static=t.static,this.animations=t.animations.slice(),this.userData=JSON.parse(JSON.stringify(t.userData)),n===!0)for(let i=0;i<t.children.length;i++){let a=t.children[i];this.add(a.clone())}return this}};tn.DEFAULT_UP=new B(0,1,0);tn.DEFAULT_MATRIX_AUTO_UPDATE=!0;tn.DEFAULT_MATRIX_WORLD_AUTO_UPDATE=!0;var sn=class extends tn{constructor(){super(),this.isGroup=!0,this.type="Group"}},vC={type:"move"},rl=class{constructor(){this._targetRay=null,this._grip=null,this._hand=null}getHandSpace(){return this._hand===null&&(this._hand=new sn,this._hand.matrixAutoUpdate=!1,this._hand.visible=!1,this._hand.joints={},this._hand.inputState={pinching:!1}),this._hand}getTargetRaySpace(){return this._targetRay===null&&(this._targetRay=new sn,this._targetRay.matrixAutoUpdate=!1,this._targetRay.visible=!1,this._targetRay.hasLinearVelocity=!1,this._targetRay.linearVelocity=new B,this._targetRay.hasAngularVelocity=!1,this._targetRay.angularVelocity=new B),this._targetRay}getGripSpace(){return this._grip===null&&(this._grip=new sn,this._grip.matrixAutoUpdate=!1,this._grip.visible=!1,this._grip.hasLinearVelocity=!1,this._grip.linearVelocity=new B,this._grip.hasAngularVelocity=!1,this._grip.angularVelocity=new B),this._grip}dispatchEvent(t){return this._targetRay!==null&&this._targetRay.dispatchEvent(t),this._grip!==null&&this._grip.dispatchEvent(t),this._hand!==null&&this._hand.dispatchEvent(t),this}connect(t){if(t&&t.hand){let n=this._hand;if(n)for(let i of t.hand.values())this._getHandJoint(n,i)}return this.dispatchEvent({type:"connected",data:t}),this}disconnect(t){return this.dispatchEvent({type:"disconnected",data:t}),this._targetRay!==null&&(this._targetRay.visible=!1),this._grip!==null&&(this._grip.visible=!1),this._hand!==null&&(this._hand.visible=!1),this}update(t,n,i){let a=null,s=null,r=null,o=this._targetRay,l=this._grip,c=this._hand;if(t&&n.session.visibilityState!=="visible-blurred"){if(c&&t.hand){r=!0;for(let b of t.hand.values()){let v=n.getJointPose(b,i),h=this._getHandJoint(c,b);v!==null&&(h.matrix.fromArray(v.transform.matrix),h.matrix.decompose(h.position,h.rotation,h.scale),h.matrixWorldNeedsUpdate=!0,h.jointRadius=v.radius),h.visible=v!==null}let f=c.joints["index-finger-tip"],p=c.joints["thumb-tip"],u=f.position.distanceTo(p.position),d=.02,_=.005;c.inputState.pinching&&u>d+_?(c.inputState.pinching=!1,this.dispatchEvent({type:"pinchend",handedness:t.handedness,target:this})):!c.inputState.pinching&&u<=d-_&&(c.inputState.pinching=!0,this.dispatchEvent({type:"pinchstart",handedness:t.handedness,target:this}))}else l!==null&&t.gripSpace&&(s=n.getPose(t.gripSpace,i),s!==null&&(l.matrix.fromArray(s.transform.matrix),l.matrix.decompose(l.position,l.rotation,l.scale),l.matrixWorldNeedsUpdate=!0,s.linearVelocity?(l.hasLinearVelocity=!0,l.linearVelocity.copy(s.linearVelocity)):l.hasLinearVelocity=!1,s.angularVelocity?(l.hasAngularVelocity=!0,l.angularVelocity.copy(s.angularVelocity)):l.hasAngularVelocity=!1));o!==null&&(a=n.getPose(t.targetRaySpace,i),a===null&&s!==null&&(a=s),a!==null&&(o.matrix.fromArray(a.transform.matrix),o.matrix.decompose(o.position,o.rotation,o.scale),o.matrixWorldNeedsUpdate=!0,a.linearVelocity?(o.hasLinearVelocity=!0,o.linearVelocity.copy(a.linearVelocity)):o.hasLinearVelocity=!1,a.angularVelocity?(o.hasAngularVelocity=!0,o.angularVelocity.copy(a.angularVelocity)):o.hasAngularVelocity=!1,this.dispatchEvent(vC)))}return o!==null&&(o.visible=a!==null),l!==null&&(l.visible=s!==null),c!==null&&(c.visible=r!==null),this}_getHandJoint(t,n){if(t.joints[n.jointName]===void 0){let i=new sn;i.matrixAutoUpdate=!1,i.visible=!1,t.joints[n.jointName]=i,t.add(i)}return t.joints[n.jointName]}},b1={aliceblue:15792383,antiquewhite:16444375,aqua:65535,aquamarine:8388564,azure:15794175,beige:16119260,bisque:16770244,black:0,blanchedalmond:16772045,blue:255,blueviolet:9055202,brown:10824234,burlywood:14596231,cadetblue:6266528,chartreuse:8388352,chocolate:13789470,coral:16744272,cornflowerblue:6591981,cornsilk:16775388,crimson:14423100,cyan:65535,darkblue:139,darkcyan:35723,darkgoldenrod:12092939,darkgray:11119017,darkgreen:25600,darkgrey:11119017,darkkhaki:12433259,darkmagenta:9109643,darkolivegreen:5597999,darkorange:16747520,darkorchid:10040012,darkred:9109504,darksalmon:15308410,darkseagreen:9419919,darkslateblue:4734347,darkslategray:3100495,darkslategrey:3100495,darkturquoise:52945,darkviolet:9699539,deeppink:16716947,deepskyblue:49151,dimgray:6908265,dimgrey:6908265,dodgerblue:2003199,firebrick:11674146,floralwhite:16775920,forestgreen:2263842,fuchsia:16711935,gainsboro:14474460,ghostwhite:16316671,gold:16766720,goldenrod:14329120,gray:8421504,green:32768,greenyellow:11403055,grey:8421504,honeydew:15794160,hotpink:16738740,indianred:13458524,indigo:4915330,ivory:16777200,khaki:15787660,lavender:15132410,lavenderblush:16773365,lawngreen:8190976,lemonchiffon:16775885,lightblue:11393254,lightcoral:15761536,lightcyan:14745599,lightgoldenrodyellow:16448210,lightgray:13882323,lightgreen:9498256,lightgrey:13882323,lightpink:16758465,lightsalmon:16752762,lightseagreen:2142890,lightskyblue:8900346,lightslategray:7833753,lightslategrey:7833753,lightsteelblue:11584734,lightyellow:16777184,lime:65280,limegreen:3329330,linen:16445670,magenta:16711935,maroon:8388608,mediumaquamarine:6737322,mediumblue:205,mediumorchid:12211667,mediumpurple:9662683,mediumseagreen:3978097,mediumslateblue:8087790,mediumspringgreen:64154,mediumturquoise:4772300,mediumvioletred:13047173,midnightblue:1644912,mintcream:16121850,mistyrose:16770273,moccasin:16770229,navajowhite:16768685,navy:128,oldlace:16643558,olive:8421376,olivedrab:7048739,orange:16753920,orangered:16729344,orchid:14315734,palegoldenrod:15657130,palegreen:10025880,paleturquoise:11529966,palevioletred:14381203,papayawhip:16773077,peachpuff:16767673,peru:13468991,pink:16761035,plum:14524637,powderblue:11591910,purple:8388736,rebeccapurple:6697881,red:16711680,rosybrown:12357519,royalblue:4286945,saddlebrown:9127187,salmon:16416882,sandybrown:16032864,seagreen:3050327,seashell:16774638,sienna:10506797,silver:12632256,skyblue:8900331,slateblue:6970061,slategray:7372944,slategrey:7372944,snow:16775930,springgreen:65407,steelblue:4620980,tan:13808780,teal:32896,thistle:14204888,tomato:16737095,turquoise:4251856,violet:15631086,wheat:16113331,white:16777215,whitesmoke:16119285,yellow:16776960,yellowgreen:10145074},Us={h:0,s:0,l:0},mh={h:0,s:0,l:0};function w0(e,t,n){return n<0&&(n+=1),n>1&&(n-=1),n<1/6?e+(t-e)*6*n:n<1/2?t:n<2/3?e+(t-e)*6*(2/3-n):e}var Ft=class{constructor(t,n,i){return this.isColor=!0,this.r=1,this.g=1,this.b=1,this.set(t,n,i)}set(t,n,i){if(n===void 0&&i===void 0){let a=t;a&&a.isColor?this.copy(a):typeof a=="number"?this.setHex(a):typeof a=="string"&&this.setStyle(a)}else this.setRGB(t,n,i);return this}setScalar(t){return this.r=t,this.g=t,this.b=t,this}setHex(t,n=On){return t=Math.floor(t),this.r=(t>>16&255)/255,this.g=(t>>8&255)/255,this.b=(t&255)/255,se.colorSpaceToWorking(this,n),this}setRGB(t,n,i,a=se.workingColorSpace){return this.r=t,this.g=n,this.b=i,se.colorSpaceToWorking(this,a),this}setHSL(t,n,i,a=se.workingColorSpace){if(t=Sv(t,1),n=ne(n,0,1),i=ne(i,0,1),n===0)this.r=this.g=this.b=i;else{let s=i<=.5?i*(1+n):i+n-i*n,r=2*i-s;this.r=w0(r,s,t+1/3),this.g=w0(r,s,t),this.b=w0(r,s,t-1/3)}return se.colorSpaceToWorking(this,a),this}setStyle(t,n=On){function i(s){s!==void 0&&parseFloat(s)<1&&Ot("Color: Alpha component of "+t+" will be ignored.")}let a;if(a=/^(\w+)\(([^\)]*)\)/.exec(t)){let s,r=a[1],o=a[2];switch(r){case"rgb":case"rgba":if(s=/^\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(?:,\s*(\d*\.?\d+)\s*)?$/.exec(o))return i(s[4]),this.setRGB(Math.min(255,parseInt(s[1],10))/255,Math.min(255,parseInt(s[2],10))/255,Math.min(255,parseInt(s[3],10))/255,n);if(s=/^\s*(\d+)\%\s*,\s*(\d+)\%\s*,\s*(\d+)\%\s*(?:,\s*(\d*\.?\d+)\s*)?$/.exec(o))return i(s[4]),this.setRGB(Math.min(100,parseInt(s[1],10))/100,Math.min(100,parseInt(s[2],10))/100,Math.min(100,parseInt(s[3],10))/100,n);break;case"hsl":case"hsla":if(s=/^\s*(\d*\.?\d+)\s*,\s*(\d*\.?\d+)\%\s*,\s*(\d*\.?\d+)\%\s*(?:,\s*(\d*\.?\d+)\s*)?$/.exec(o))return i(s[4]),this.setHSL(parseFloat(s[1])/360,parseFloat(s[2])/100,parseFloat(s[3])/100,n);break;default:Ot("Color: Unknown color model "+t)}}else if(a=/^\#([A-Fa-f\d]+)$/.exec(t)){let s=a[1],r=s.length;if(r===3)return this.setRGB(parseInt(s.charAt(0),16)/15,parseInt(s.charAt(1),16)/15,parseInt(s.charAt(2),16)/15,n);if(r===6)return this.setHex(parseInt(s,16),n);Ot("Color: Invalid hex color "+t)}else if(t&&t.length>0)return this.setColorName(t,n);return this}setColorName(t,n=On){let i=b1[t.toLowerCase()];return i!==void 0?this.setHex(i,n):Ot("Color: Unknown color "+t),this}clone(){return new this.constructor(this.r,this.g,this.b)}copy(t){return this.r=t.r,this.g=t.g,this.b=t.b,this}copySRGBToLinear(t){return this.r=Xa(t.r),this.g=Xa(t.g),this.b=Xa(t.b),this}copyLinearToSRGB(t){return this.r=el(t.r),this.g=el(t.g),this.b=el(t.b),this}convertSRGBToLinear(){return this.copySRGBToLinear(this),this}convertLinearToSRGB(){return this.copyLinearToSRGB(this),this}getHex(t=On){return se.workingToColorSpace(In.copy(this),t),Math.round(ne(In.r*255,0,255))*65536+Math.round(ne(In.g*255,0,255))*256+Math.round(ne(In.b*255,0,255))}getHexString(t=On){return("000000"+this.getHex(t).toString(16)).slice(-6)}getHSL(t,n=se.workingColorSpace){se.workingToColorSpace(In.copy(this),n);let i=In.r,a=In.g,s=In.b,r=Math.max(i,a,s),o=Math.min(i,a,s),l,c,f=(o+r)/2;if(o===r)l=0,c=0;else{let p=r-o;switch(c=f<=.5?p/(r+o):p/(2-r-o),r){case i:l=(a-s)/p+(a<s?6:0);break;case a:l=(s-i)/p+2;break;case s:l=(i-a)/p+4;break}l/=6}return t.h=l,t.s=c,t.l=f,t}getRGB(t,n=se.workingColorSpace){return se.workingToColorSpace(In.copy(this),n),t.r=In.r,t.g=In.g,t.b=In.b,t}getStyle(t=On){se.workingToColorSpace(In.copy(this),t);let n=In.r,i=In.g,a=In.b;return t!==On?`color(${t} ${n.toFixed(3)} ${i.toFixed(3)} ${a.toFixed(3)})`:`rgb(${Math.round(n*255)},${Math.round(i*255)},${Math.round(a*255)})`}offsetHSL(t,n,i){return this.getHSL(Us),this.setHSL(Us.h+t,Us.s+n,Us.l+i)}add(t){return this.r+=t.r,this.g+=t.g,this.b+=t.b,this}addColors(t,n){return this.r=t.r+n.r,this.g=t.g+n.g,this.b=t.b+n.b,this}addScalar(t){return this.r+=t,this.g+=t,this.b+=t,this}sub(t){return this.r=Math.max(0,this.r-t.r),this.g=Math.max(0,this.g-t.g),this.b=Math.max(0,this.b-t.b),this}multiply(t){return this.r*=t.r,this.g*=t.g,this.b*=t.b,this}multiplyScalar(t){return this.r*=t,this.g*=t,this.b*=t,this}lerp(t,n){return this.r+=(t.r-this.r)*n,this.g+=(t.g-this.g)*n,this.b+=(t.b-this.b)*n,this}lerpColors(t,n,i){return this.r=t.r+(n.r-t.r)*i,this.g=t.g+(n.g-t.g)*i,this.b=t.b+(n.b-t.b)*i,this}lerpHSL(t,n){this.getHSL(Us),t.getHSL(mh);let i=zc(Us.h,mh.h,n),a=zc(Us.s,mh.s,n),s=zc(Us.l,mh.l,n);return this.setHSL(i,a,s),this}setFromVector3(t){return this.r=t.x,this.g=t.y,this.b=t.z,this}applyMatrix3(t){let n=this.r,i=this.g,a=this.b,s=t.elements;return this.r=s[0]*n+s[3]*i+s[6]*a,this.g=s[1]*n+s[4]*i+s[7]*a,this.b=s[2]*n+s[5]*i+s[8]*a,this}equals(t){return t.r===this.r&&t.g===this.g&&t.b===this.b}fromArray(t,n=0){return this.r=t[n],this.g=t[n+1],this.b=t[n+2],this}toArray(t=[],n=0){return t[n]=this.r,t[n+1]=this.g,t[n+2]=this.b,t}fromBufferAttribute(t,n){return this.r=t.getX(n),this.g=t.getY(n),this.b=t.getZ(n),this}toJSON(){return this.getHex()}*[Symbol.iterator](){yield this.r,yield this.g,yield this.b}},In=new Ft;Ft.NAMES=b1;var Wc=class e{constructor(t,n=1,i=1e3){this.isFog=!0,this.name="",this.color=new Ft(t),this.near=n,this.far=i}clone(){return new e(this.color,this.near,this.far)}toJSON(){return{type:"Fog",name:this.name,color:this.color.getHex(),near:this.near,far:this.far}}},qc=class extends tn{constructor(){super(),this.isScene=!0,this.type="Scene",this.background=null,this.environment=null,this.fog=null,this.backgroundBlurriness=0,this.backgroundIntensity=1,this.backgroundRotation=new Xi,this.environmentIntensity=1,this.environmentRotation=new Xi,this.overrideMaterial=null,typeof __THREE_DEVTOOLS__<"u"&&__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("observe",{detail:this}))}copy(t,n){return super.copy(t,n),t.background!==null&&(this.background=t.background.clone()),t.environment!==null&&(this.environment=t.environment.clone()),t.fog!==null&&(this.fog=t.fog.clone()),this.backgroundBlurriness=t.backgroundBlurriness,this.backgroundIntensity=t.backgroundIntensity,this.backgroundRotation.copy(t.backgroundRotation),this.environmentIntensity=t.environmentIntensity,this.environmentRotation.copy(t.environmentRotation),t.overrideMaterial!==null&&(this.overrideMaterial=t.overrideMaterial.clone()),this.matrixAutoUpdate=t.matrixAutoUpdate,this}toJSON(t){let n=super.toJSON(t);return this.fog!==null&&(n.object.fog=this.fog.toJSON()),this.backgroundBlurriness>0&&(n.object.backgroundBlurriness=this.backgroundBlurriness),this.backgroundIntensity!==1&&(n.object.backgroundIntensity=this.backgroundIntensity),n.object.backgroundRotation=this.backgroundRotation.toArray(),this.environmentIntensity!==1&&(n.object.environmentIntensity=this.environmentIntensity),n.object.environmentRotation=this.environmentRotation.toArray(),n}},Hi=new B,Va=new B,C0=new B,Ha=new B,Wo=new B,qo=new B,AM=new B,R0=new B,D0=new B,N0=new B,U0=new Ge,L0=new Ge,I0=new Ge,Bs=class e{constructor(t=new B,n=new B,i=new B){this.a=t,this.b=n,this.c=i}static getNormal(t,n,i,a){a.subVectors(i,n),Hi.subVectors(t,n),a.cross(Hi);let s=a.lengthSq();return s>0?a.multiplyScalar(1/Math.sqrt(s)):a.set(0,0,0)}static getBarycoord(t,n,i,a,s){Hi.subVectors(a,n),Va.subVectors(i,n),C0.subVectors(t,n);let r=Hi.dot(Hi),o=Hi.dot(Va),l=Hi.dot(C0),c=Va.dot(Va),f=Va.dot(C0),p=r*c-o*o;if(p===0)return s.set(0,0,0),null;let u=1/p,d=(c*l-o*f)*u,_=(r*f-o*l)*u;return s.set(1-d-_,_,d)}static containsPoint(t,n,i,a){return this.getBarycoord(t,n,i,a,Ha)===null?!1:Ha.x>=0&&Ha.y>=0&&Ha.x+Ha.y<=1}static getInterpolation(t,n,i,a,s,r,o,l){return this.getBarycoord(t,n,i,a,Ha)===null?(l.x=0,l.y=0,"z"in l&&(l.z=0),"w"in l&&(l.w=0),null):(l.setScalar(0),l.addScaledVector(s,Ha.x),l.addScaledVector(r,Ha.y),l.addScaledVector(o,Ha.z),l)}static getInterpolatedAttribute(t,n,i,a,s,r){return U0.setScalar(0),L0.setScalar(0),I0.setScalar(0),U0.fromBufferAttribute(t,n),L0.fromBufferAttribute(t,i),I0.fromBufferAttribute(t,a),r.setScalar(0),r.addScaledVector(U0,s.x),r.addScaledVector(L0,s.y),r.addScaledVector(I0,s.z),r}static isFrontFacing(t,n,i,a){return Hi.subVectors(i,n),Va.subVectors(t,n),Hi.cross(Va).dot(a)<0}set(t,n,i){return this.a.copy(t),this.b.copy(n),this.c.copy(i),this}setFromPointsAndIndices(t,n,i,a){return this.a.copy(t[n]),this.b.copy(t[i]),this.c.copy(t[a]),this}setFromAttributeAndIndices(t,n,i,a){return this.a.fromBufferAttribute(t,n),this.b.fromBufferAttribute(t,i),this.c.fromBufferAttribute(t,a),this}clone(){return new this.constructor().copy(this)}copy(t){return this.a.copy(t.a),this.b.copy(t.b),this.c.copy(t.c),this}getArea(){return Hi.subVectors(this.c,this.b),Va.subVectors(this.a,this.b),Hi.cross(Va).length()*.5}getMidpoint(t){return t.addVectors(this.a,this.b).add(this.c).multiplyScalar(1/3)}getNormal(t){return e.getNormal(this.a,this.b,this.c,t)}getPlane(t){return t.setFromCoplanarPoints(this.a,this.b,this.c)}getBarycoord(t,n){return e.getBarycoord(t,this.a,this.b,this.c,n)}getInterpolation(t,n,i,a,s){return e.getInterpolation(t,this.a,this.b,this.c,n,i,a,s)}containsPoint(t){return e.containsPoint(t,this.a,this.b,this.c)}isFrontFacing(t){return e.isFrontFacing(this.a,this.b,this.c,t)}intersectsBox(t){return t.intersectsTriangle(this)}closestPointToPoint(t,n){let i=this.a,a=this.b,s=this.c,r,o;Wo.subVectors(a,i),qo.subVectors(s,i),R0.subVectors(t,i);let l=Wo.dot(R0),c=qo.dot(R0);if(l<=0&&c<=0)return n.copy(i);D0.subVectors(t,a);let f=Wo.dot(D0),p=qo.dot(D0);if(f>=0&&p<=f)return n.copy(a);let u=l*p-f*c;if(u<=0&&l>=0&&f<=0)return r=l/(l-f),n.copy(i).addScaledVector(Wo,r);N0.subVectors(t,s);let d=Wo.dot(N0),_=qo.dot(N0);if(_>=0&&d<=_)return n.copy(s);let b=d*c-l*_;if(b<=0&&c>=0&&_<=0)return o=c/(c-_),n.copy(i).addScaledVector(qo,o);let v=f*_-d*p;if(v<=0&&p-f>=0&&d-_>=0)return AM.subVectors(s,a),o=(p-f)/(p-f+(d-_)),n.copy(a).addScaledVector(AM,o);let h=1/(v+b+u);return r=b*h,o=u*h,n.copy(i).addScaledVector(Wo,r).addScaledVector(qo,o)}equals(t){return t.a.equals(this.a)&&t.b.equals(this.b)&&t.c.equals(this.c)}},ma=class{constructor(t=new B(1/0,1/0,1/0),n=new B(-1/0,-1/0,-1/0)){this.isBox3=!0,this.min=t,this.max=n}set(t,n){return this.min.copy(t),this.max.copy(n),this}setFromArray(t){this.makeEmpty();for(let n=0,i=t.length;n<i;n+=3)this.expandByPoint(Gi.fromArray(t,n));return this}setFromBufferAttribute(t){this.makeEmpty();for(let n=0,i=t.count;n<i;n++)this.expandByPoint(Gi.fromBufferAttribute(t,n));return this}setFromPoints(t){this.makeEmpty();for(let n=0,i=t.length;n<i;n++)this.expandByPoint(t[n]);return this}setFromCenterAndSize(t,n){let i=Gi.copy(n).multiplyScalar(.5);return this.min.copy(t).sub(i),this.max.copy(t).add(i),this}setFromObject(t,n=!1){return this.makeEmpty(),this.expandByObject(t,n)}clone(){return new this.constructor().copy(this)}copy(t){return this.min.copy(t.min),this.max.copy(t.max),this}makeEmpty(){return this.min.x=this.min.y=this.min.z=1/0,this.max.x=this.max.y=this.max.z=-1/0,this}isEmpty(){return this.max.x<this.min.x||this.max.y<this.min.y||this.max.z<this.min.z}getCenter(t){return this.isEmpty()?t.set(0,0,0):t.addVectors(this.min,this.max).multiplyScalar(.5)}getSize(t){return this.isEmpty()?t.set(0,0,0):t.subVectors(this.max,this.min)}expandByPoint(t){return this.min.min(t),this.max.max(t),this}expandByVector(t){return this.min.sub(t),this.max.add(t),this}expandByScalar(t){return this.min.addScalar(-t),this.max.addScalar(t),this}expandByObject(t,n=!1){t.updateWorldMatrix(!1,!1);let i=t.geometry;if(i!==void 0){let s=i.getAttribute("position");if(n===!0&&s!==void 0&&t.isInstancedMesh!==!0)for(let r=0,o=s.count;r<o;r++)t.isMesh===!0?t.getVertexPosition(r,Gi):Gi.fromBufferAttribute(s,r),Gi.applyMatrix4(t.matrixWorld),this.expandByPoint(Gi);else t.boundingBox!==void 0?(t.boundingBox===null&&t.computeBoundingBox(),gh.copy(t.boundingBox)):(i.boundingBox===null&&i.computeBoundingBox(),gh.copy(i.boundingBox)),gh.applyMatrix4(t.matrixWorld),this.union(gh)}let a=t.children;for(let s=0,r=a.length;s<r;s++)this.expandByObject(a[s],n);return this}containsPoint(t){return t.x>=this.min.x&&t.x<=this.max.x&&t.y>=this.min.y&&t.y<=this.max.y&&t.z>=this.min.z&&t.z<=this.max.z}containsBox(t){return this.min.x<=t.min.x&&t.max.x<=this.max.x&&this.min.y<=t.min.y&&t.max.y<=this.max.y&&this.min.z<=t.min.z&&t.max.z<=this.max.z}getParameter(t,n){return n.set((t.x-this.min.x)/(this.max.x-this.min.x),(t.y-this.min.y)/(this.max.y-this.min.y),(t.z-this.min.z)/(this.max.z-this.min.z))}intersectsBox(t){return t.max.x>=this.min.x&&t.min.x<=this.max.x&&t.max.y>=this.min.y&&t.min.y<=this.max.y&&t.max.z>=this.min.z&&t.min.z<=this.max.z}intersectsSphere(t){return this.clampPoint(t.center,Gi),Gi.distanceToSquared(t.center)<=t.radius*t.radius}intersectsPlane(t){let n,i;return t.normal.x>0?(n=t.normal.x*this.min.x,i=t.normal.x*this.max.x):(n=t.normal.x*this.max.x,i=t.normal.x*this.min.x),t.normal.y>0?(n+=t.normal.y*this.min.y,i+=t.normal.y*this.max.y):(n+=t.normal.y*this.max.y,i+=t.normal.y*this.min.y),t.normal.z>0?(n+=t.normal.z*this.min.z,i+=t.normal.z*this.max.z):(n+=t.normal.z*this.max.z,i+=t.normal.z*this.min.z),n<=-t.constant&&i>=-t.constant}intersectsTriangle(t){if(this.isEmpty())return!1;this.getCenter(Nc),vh.subVectors(this.max,Nc),Yo.subVectors(t.a,Nc),Zo.subVectors(t.b,Nc),Jo.subVectors(t.c,Nc),Ls.subVectors(Zo,Yo),Is.subVectors(Jo,Zo),br.subVectors(Yo,Jo);let n=[0,-Ls.z,Ls.y,0,-Is.z,Is.y,0,-br.z,br.y,Ls.z,0,-Ls.x,Is.z,0,-Is.x,br.z,0,-br.x,-Ls.y,Ls.x,0,-Is.y,Is.x,0,-br.y,br.x,0];return!O0(n,Yo,Zo,Jo,vh)||(n=[1,0,0,0,1,0,0,0,1],!O0(n,Yo,Zo,Jo,vh))?!1:(_h.crossVectors(Ls,Is),n=[_h.x,_h.y,_h.z],O0(n,Yo,Zo,Jo,vh))}clampPoint(t,n){return n.copy(t).clamp(this.min,this.max)}distanceToPoint(t){return this.clampPoint(t,Gi).distanceTo(t)}getBoundingSphere(t){return this.isEmpty()?t.makeEmpty():(this.getCenter(t.center),t.radius=this.getSize(Gi).length()*.5),t}intersect(t){return this.min.max(t.min),this.max.min(t.max),this.isEmpty()&&this.makeEmpty(),this}union(t){return this.min.min(t.min),this.max.max(t.max),this}applyMatrix4(t){return this.isEmpty()?this:(Ga[0].set(this.min.x,this.min.y,this.min.z).applyMatrix4(t),Ga[1].set(this.min.x,this.min.y,this.max.z).applyMatrix4(t),Ga[2].set(this.min.x,this.max.y,this.min.z).applyMatrix4(t),Ga[3].set(this.min.x,this.max.y,this.max.z).applyMatrix4(t),Ga[4].set(this.max.x,this.min.y,this.min.z).applyMatrix4(t),Ga[5].set(this.max.x,this.min.y,this.max.z).applyMatrix4(t),Ga[6].set(this.max.x,this.max.y,this.min.z).applyMatrix4(t),Ga[7].set(this.max.x,this.max.y,this.max.z).applyMatrix4(t),this.setFromPoints(Ga),this)}translate(t){return this.min.add(t),this.max.add(t),this}equals(t){return t.min.equals(this.min)&&t.max.equals(this.max)}toJSON(){return{min:this.min.toArray(),max:this.max.toArray()}}fromJSON(t){return this.min.fromArray(t.min),this.max.fromArray(t.max),this}},Ga=[new B,new B,new B,new B,new B,new B,new B,new B],Gi=new B,gh=new ma,Yo=new B,Zo=new B,Jo=new B,Ls=new B,Is=new B,br=new B,Nc=new B,vh=new B,_h=new B,Mr=new B;function O0(e,t,n,i,a){for(let s=0,r=e.length-3;s<=r;s+=3){Mr.fromArray(e,s);let o=a.x*Math.abs(Mr.x)+a.y*Math.abs(Mr.y)+a.z*Math.abs(Mr.z),l=t.dot(Mr),c=n.dot(Mr),f=i.dot(Mr);if(Math.max(-Math.max(l,c,f),Math.min(l,c,f))>o)return!1}return!0}var $e=new B,yh=new te,_C=0,Bn=class{constructor(t,n,i=!1){if(Array.isArray(t))throw new TypeError("THREE.BufferAttribute: array should be a Typed Array.");this.isBufferAttribute=!0,Object.defineProperty(this,"id",{value:_C++}),this.name="",this.array=t,this.itemSize=n,this.count=t!==void 0?t.length/n:0,this.normalized=i,this.usage=K0,this.updateRanges=[],this.gpuType=Ni,this.version=0}onUploadCallback(){}set needsUpdate(t){t===!0&&this.version++}setUsage(t){return this.usage=t,this}addUpdateRange(t,n){this.updateRanges.push({start:t,count:n})}clearUpdateRanges(){this.updateRanges.length=0}copy(t){return this.name=t.name,this.array=new t.array.constructor(t.array),this.itemSize=t.itemSize,this.count=t.count,this.normalized=t.normalized,this.usage=t.usage,this.gpuType=t.gpuType,this}copyAt(t,n,i){t*=this.itemSize,i*=n.itemSize;for(let a=0,s=this.itemSize;a<s;a++)this.array[t+a]=n.array[i+a];return this}copyArray(t){return this.array.set(t),this}applyMatrix3(t){if(this.itemSize===2)for(let n=0,i=this.count;n<i;n++)yh.fromBufferAttribute(this,n),yh.applyMatrix3(t),this.setXY(n,yh.x,yh.y);else if(this.itemSize===3)for(let n=0,i=this.count;n<i;n++)$e.fromBufferAttribute(this,n),$e.applyMatrix3(t),this.setXYZ(n,$e.x,$e.y,$e.z);return this}applyMatrix4(t){for(let n=0,i=this.count;n<i;n++)$e.fromBufferAttribute(this,n),$e.applyMatrix4(t),this.setXYZ(n,$e.x,$e.y,$e.z);return this}applyNormalMatrix(t){for(let n=0,i=this.count;n<i;n++)$e.fromBufferAttribute(this,n),$e.applyNormalMatrix(t),this.setXYZ(n,$e.x,$e.y,$e.z);return this}transformDirection(t){for(let n=0,i=this.count;n<i;n++)$e.fromBufferAttribute(this,n),$e.transformDirection(t),this.setXYZ(n,$e.x,$e.y,$e.z);return this}set(t,n=0){return this.array.set(t,n),this}getComponent(t,n){let i=this.array[t*this.itemSize+n];return this.normalized&&(i=tl(i,this.array)),i}setComponent(t,n,i){return this.normalized&&(i=Wn(i,this.array)),this.array[t*this.itemSize+n]=i,this}getX(t){let n=this.array[t*this.itemSize];return this.normalized&&(n=tl(n,this.array)),n}setX(t,n){return this.normalized&&(n=Wn(n,this.array)),this.array[t*this.itemSize]=n,this}getY(t){let n=this.array[t*this.itemSize+1];return this.normalized&&(n=tl(n,this.array)),n}setY(t,n){return this.normalized&&(n=Wn(n,this.array)),this.array[t*this.itemSize+1]=n,this}getZ(t){let n=this.array[t*this.itemSize+2];return this.normalized&&(n=tl(n,this.array)),n}setZ(t,n){return this.normalized&&(n=Wn(n,this.array)),this.array[t*this.itemSize+2]=n,this}getW(t){let n=this.array[t*this.itemSize+3];return this.normalized&&(n=tl(n,this.array)),n}setW(t,n){return this.normalized&&(n=Wn(n,this.array)),this.array[t*this.itemSize+3]=n,this}setXY(t,n,i){return t*=this.itemSize,this.normalized&&(n=Wn(n,this.array),i=Wn(i,this.array)),this.array[t+0]=n,this.array[t+1]=i,this}setXYZ(t,n,i,a){return t*=this.itemSize,this.normalized&&(n=Wn(n,this.array),i=Wn(i,this.array),a=Wn(a,this.array)),this.array[t+0]=n,this.array[t+1]=i,this.array[t+2]=a,this}setXYZW(t,n,i,a,s){return t*=this.itemSize,this.normalized&&(n=Wn(n,this.array),i=Wn(i,this.array),a=Wn(a,this.array),s=Wn(s,this.array)),this.array[t+0]=n,this.array[t+1]=i,this.array[t+2]=a,this.array[t+3]=s,this}onUpload(t){return this.onUploadCallback=t,this}clone(){return new this.constructor(this.array,this.itemSize).copy(this)}toJSON(){let t={itemSize:this.itemSize,type:this.array.constructor.name,array:Array.from(this.array),normalized:this.normalized};return this.name!==""&&(t.name=this.name),this.usage!==K0&&(t.usage=this.usage),t}};var Yc=class extends Bn{constructor(t,n,i){super(new Uint16Array(t),n,i)}};var Zc=class extends Bn{constructor(t,n,i){super(new Uint32Array(t),n,i)}};var Oe=class extends Bn{constructor(t,n,i){super(new Float32Array(t),n,i)}},yC=new ma,Uc=new B,P0=new B,Ya=class{constructor(t=new B,n=-1){this.isSphere=!0,this.center=t,this.radius=n}set(t,n){return this.center.copy(t),this.radius=n,this}setFromPoints(t,n){let i=this.center;n!==void 0?i.copy(n):yC.setFromPoints(t).getCenter(i);let a=0;for(let s=0,r=t.length;s<r;s++)a=Math.max(a,i.distanceToSquared(t[s]));return this.radius=Math.sqrt(a),this}copy(t){return this.center.copy(t.center),this.radius=t.radius,this}isEmpty(){return this.radius<0}makeEmpty(){return this.center.set(0,0,0),this.radius=-1,this}containsPoint(t){return t.distanceToSquared(this.center)<=this.radius*this.radius}distanceToPoint(t){return t.distanceTo(this.center)-this.radius}intersectsSphere(t){let n=this.radius+t.radius;return t.center.distanceToSquared(this.center)<=n*n}intersectsBox(t){return t.intersectsSphere(this)}intersectsPlane(t){return Math.abs(t.distanceToPoint(this.center))<=this.radius}clampPoint(t,n){let i=this.center.distanceToSquared(t);return n.copy(t),i>this.radius*this.radius&&(n.sub(this.center).normalize(),n.multiplyScalar(this.radius).add(this.center)),n}getBoundingBox(t){return this.isEmpty()?(t.makeEmpty(),t):(t.set(this.center,this.center),t.expandByScalar(this.radius),t)}applyMatrix4(t){return this.center.applyMatrix4(t),this.radius=this.radius*t.getMaxScaleOnAxis(),this}translate(t){return this.center.add(t),this}expandByPoint(t){if(this.isEmpty())return this.center.copy(t),this.radius=0,this;Uc.subVectors(t,this.center);let n=Uc.lengthSq();if(n>this.radius*this.radius){let i=Math.sqrt(n),a=(i-this.radius)*.5;this.center.addScaledVector(Uc,a/i),this.radius+=a}return this}union(t){return t.isEmpty()?this:this.isEmpty()?(this.copy(t),this):(this.center.equals(t.center)===!0?this.radius=Math.max(this.radius,t.radius):(P0.subVectors(t.center,this.center).setLength(t.radius),this.expandByPoint(Uc.copy(t.center).add(P0)),this.expandByPoint(Uc.copy(t.center).sub(P0))),this)}equals(t){return t.center.equals(this.center)&&t.radius===this.radius}clone(){return new this.constructor().copy(this)}toJSON(){return{radius:this.radius,center:this.center.toArray()}}fromJSON(t){return this.radius=t.radius,this.center.fromArray(t.center),this}},xC=0,Di=new Ie,B0=new tn,Ko=new B,mi=new ma,Lc=new ma,dn=new B,An=class e extends qa{constructor(){super(),this.isBufferGeometry=!0,Object.defineProperty(this,"id",{value:xC++}),this.uuid=gl(),this.name="",this.type="BufferGeometry",this.index=null,this.indirect=null,this.indirectOffset=0,this.attributes={},this.morphAttributes={},this.morphTargetsRelative=!1,this.groups=[],this.boundingBox=null,this.boundingSphere=null,this.drawRange={start:0,count:1/0},this.userData={}}getIndex(){return this.index}setIndex(t){return Array.isArray(t)?this.index=new(Xw(t)?Zc:Yc)(t,1):this.index=t,this}setIndirect(t,n=0){return this.indirect=t,this.indirectOffset=n,this}getIndirect(){return this.indirect}getAttribute(t){return this.attributes[t]}setAttribute(t,n){return this.attributes[t]=n,this}deleteAttribute(t){return delete this.attributes[t],this}hasAttribute(t){return this.attributes[t]!==void 0}addGroup(t,n,i=0){this.groups.push({start:t,count:n,materialIndex:i})}clearGroups(){this.groups=[]}setDrawRange(t,n){this.drawRange.start=t,this.drawRange.count=n}applyMatrix4(t){let n=this.attributes.position;n!==void 0&&(n.applyMatrix4(t),n.needsUpdate=!0);let i=this.attributes.normal;if(i!==void 0){let s=new Zt().getNormalMatrix(t);i.applyNormalMatrix(s),i.needsUpdate=!0}let a=this.attributes.tangent;return a!==void 0&&(a.transformDirection(t),a.needsUpdate=!0),this.boundingBox!==null&&this.computeBoundingBox(),this.boundingSphere!==null&&this.computeBoundingSphere(),this}applyQuaternion(t){return Di.makeRotationFromQuaternion(t),this.applyMatrix4(Di),this}rotateX(t){return Di.makeRotationX(t),this.applyMatrix4(Di),this}rotateY(t){return Di.makeRotationY(t),this.applyMatrix4(Di),this}rotateZ(t){return Di.makeRotationZ(t),this.applyMatrix4(Di),this}translate(t,n,i){return Di.makeTranslation(t,n,i),this.applyMatrix4(Di),this}scale(t,n,i){return Di.makeScale(t,n,i),this.applyMatrix4(Di),this}lookAt(t){return B0.lookAt(t),B0.updateMatrix(),this.applyMatrix4(B0.matrix),this}center(){return this.computeBoundingBox(),this.boundingBox.getCenter(Ko).negate(),this.translate(Ko.x,Ko.y,Ko.z),this}setFromPoints(t){let n=this.getAttribute("position");if(n===void 0){let i=[];for(let a=0,s=t.length;a<s;a++){let r=t[a];i.push(r.x,r.y,r.z||0)}this.setAttribute("position",new Oe(i,3))}else{let i=Math.min(t.length,n.count);for(let a=0;a<i;a++){let s=t[a];n.setXYZ(a,s.x,s.y,s.z||0)}t.length>n.count&&Ot("BufferGeometry: Buffer size too small for points data. Use .dispose() and create a new geometry."),n.needsUpdate=!0}return this}computeBoundingBox(){this.boundingBox===null&&(this.boundingBox=new ma);let t=this.attributes.position,n=this.morphAttributes.position;if(t&&t.isGLBufferAttribute){zt("BufferGeometry.computeBoundingBox(): GLBufferAttribute requires a manual bounding box.",this),this.boundingBox.set(new B(-1/0,-1/0,-1/0),new B(1/0,1/0,1/0));return}if(t!==void 0){if(this.boundingBox.setFromBufferAttribute(t),n)for(let i=0,a=n.length;i<a;i++){let s=n[i];mi.setFromBufferAttribute(s),this.morphTargetsRelative?(dn.addVectors(this.boundingBox.min,mi.min),this.boundingBox.expandByPoint(dn),dn.addVectors(this.boundingBox.max,mi.max),this.boundingBox.expandByPoint(dn)):(this.boundingBox.expandByPoint(mi.min),this.boundingBox.expandByPoint(mi.max))}}else this.boundingBox.makeEmpty();(isNaN(this.boundingBox.min.x)||isNaN(this.boundingBox.min.y)||isNaN(this.boundingBox.min.z))&&zt('BufferGeometry.computeBoundingBox(): Computed min/max have NaN values. The "position" attribute is likely to have NaN values.',this)}computeBoundingSphere(){this.boundingSphere===null&&(this.boundingSphere=new Ya);let t=this.attributes.position,n=this.morphAttributes.position;if(t&&t.isGLBufferAttribute){zt("BufferGeometry.computeBoundingSphere(): GLBufferAttribute requires a manual bounding sphere.",this),this.boundingSphere.set(new B,1/0);return}if(t){let i=this.boundingSphere.center;if(mi.setFromBufferAttribute(t),n)for(let s=0,r=n.length;s<r;s++){let o=n[s];Lc.setFromBufferAttribute(o),this.morphTargetsRelative?(dn.addVectors(mi.min,Lc.min),mi.expandByPoint(dn),dn.addVectors(mi.max,Lc.max),mi.expandByPoint(dn)):(mi.expandByPoint(Lc.min),mi.expandByPoint(Lc.max))}mi.getCenter(i);let a=0;for(let s=0,r=t.count;s<r;s++)dn.fromBufferAttribute(t,s),a=Math.max(a,i.distanceToSquared(dn));if(n)for(let s=0,r=n.length;s<r;s++){let o=n[s],l=this.morphTargetsRelative;for(let c=0,f=o.count;c<f;c++)dn.fromBufferAttribute(o,c),l&&(Ko.fromBufferAttribute(t,c),dn.add(Ko)),a=Math.max(a,i.distanceToSquared(dn))}this.boundingSphere.radius=Math.sqrt(a),isNaN(this.boundingSphere.radius)&&zt('BufferGeometry.computeBoundingSphere(): Computed radius is NaN. The "position" attribute is likely to have NaN values.',this)}}computeTangents(){let t=this.index,n=this.attributes;if(t===null||n.position===void 0||n.normal===void 0||n.uv===void 0){zt("BufferGeometry: .computeTangents() failed. Missing required attributes (index, position, normal or uv)");return}let i=n.position,a=n.normal,s=n.uv;this.hasAttribute("tangent")===!1&&this.setAttribute("tangent",new Bn(new Float32Array(4*i.count),4));let r=this.getAttribute("tangent"),o=[],l=[];for(let y=0;y<i.count;y++)o[y]=new B,l[y]=new B;let c=new B,f=new B,p=new B,u=new te,d=new te,_=new te,b=new B,v=new B;function h(y,T,L){c.fromBufferAttribute(i,y),f.fromBufferAttribute(i,T),p.fromBufferAttribute(i,L),u.fromBufferAttribute(s,y),d.fromBufferAttribute(s,T),_.fromBufferAttribute(s,L),f.sub(c),p.sub(c),d.sub(u),_.sub(u);let R=1/(d.x*_.y-_.x*d.y);isFinite(R)&&(b.copy(f).multiplyScalar(_.y).addScaledVector(p,-d.y).multiplyScalar(R),v.copy(p).multiplyScalar(d.x).addScaledVector(f,-_.x).multiplyScalar(R),o[y].add(b),o[T].add(b),o[L].add(b),l[y].add(v),l[T].add(v),l[L].add(v))}let m=this.groups;m.length===0&&(m=[{start:0,count:t.count}]);for(let y=0,T=m.length;y<T;++y){let L=m[y],R=L.start,N=L.count;for(let P=R,z=R+N;P<z;P+=3)h(t.getX(P+0),t.getX(P+1),t.getX(P+2))}let g=new B,x=new B,E=new B,A=new B;function w(y){E.fromBufferAttribute(a,y),A.copy(E);let T=o[y];g.copy(T),g.sub(E.multiplyScalar(E.dot(T))).normalize(),x.crossVectors(A,T);let R=x.dot(l[y])<0?-1:1;r.setXYZW(y,g.x,g.y,g.z,R)}for(let y=0,T=m.length;y<T;++y){let L=m[y],R=L.start,N=L.count;for(let P=R,z=R+N;P<z;P+=3)w(t.getX(P+0)),w(t.getX(P+1)),w(t.getX(P+2))}}computeVertexNormals(){let t=this.index,n=this.getAttribute("position");if(n!==void 0){let i=this.getAttribute("normal");if(i===void 0)i=new Bn(new Float32Array(n.count*3),3),this.setAttribute("normal",i);else for(let u=0,d=i.count;u<d;u++)i.setXYZ(u,0,0,0);let a=new B,s=new B,r=new B,o=new B,l=new B,c=new B,f=new B,p=new B;if(t)for(let u=0,d=t.count;u<d;u+=3){let _=t.getX(u+0),b=t.getX(u+1),v=t.getX(u+2);a.fromBufferAttribute(n,_),s.fromBufferAttribute(n,b),r.fromBufferAttribute(n,v),f.subVectors(r,s),p.subVectors(a,s),f.cross(p),o.fromBufferAttribute(i,_),l.fromBufferAttribute(i,b),c.fromBufferAttribute(i,v),o.add(f),l.add(f),c.add(f),i.setXYZ(_,o.x,o.y,o.z),i.setXYZ(b,l.x,l.y,l.z),i.setXYZ(v,c.x,c.y,c.z)}else for(let u=0,d=n.count;u<d;u+=3)a.fromBufferAttribute(n,u+0),s.fromBufferAttribute(n,u+1),r.fromBufferAttribute(n,u+2),f.subVectors(r,s),p.subVectors(a,s),f.cross(p),i.setXYZ(u+0,f.x,f.y,f.z),i.setXYZ(u+1,f.x,f.y,f.z),i.setXYZ(u+2,f.x,f.y,f.z);this.normalizeNormals(),i.needsUpdate=!0}}normalizeNormals(){let t=this.attributes.normal;for(let n=0,i=t.count;n<i;n++)dn.fromBufferAttribute(t,n),dn.normalize(),t.setXYZ(n,dn.x,dn.y,dn.z)}toNonIndexed(){function t(o,l){let c=o.array,f=o.itemSize,p=o.normalized,u=new c.constructor(l.length*f),d=0,_=0;for(let b=0,v=l.length;b<v;b++){o.isInterleavedBufferAttribute?d=l[b]*o.data.stride+o.offset:d=l[b]*f;for(let h=0;h<f;h++)u[_++]=c[d++]}return new Bn(u,f,p)}if(this.index===null)return Ot("BufferGeometry.toNonIndexed(): BufferGeometry is already non-indexed."),this;let n=new e,i=this.index.array,a=this.attributes;for(let o in a){let l=a[o],c=t(l,i);n.setAttribute(o,c)}let s=this.morphAttributes;for(let o in s){let l=[],c=s[o];for(let f=0,p=c.length;f<p;f++){let u=c[f],d=t(u,i);l.push(d)}n.morphAttributes[o]=l}n.morphTargetsRelative=this.morphTargetsRelative;let r=this.groups;for(let o=0,l=r.length;o<l;o++){let c=r[o];n.addGroup(c.start,c.count,c.materialIndex)}return n}toJSON(){let t={metadata:{version:4.7,type:"BufferGeometry",generator:"BufferGeometry.toJSON"}};if(t.uuid=this.uuid,t.type=this.type,this.name!==""&&(t.name=this.name),Object.keys(this.userData).length>0&&(t.userData=this.userData),this.parameters!==void 0){let l=this.parameters;for(let c in l)l[c]!==void 0&&(t[c]=l[c]);return t}t.data={attributes:{}};let n=this.index;n!==null&&(t.data.index={type:n.array.constructor.name,array:Array.prototype.slice.call(n.array)});let i=this.attributes;for(let l in i){let c=i[l];t.data.attributes[l]=c.toJSON(t.data)}let a={},s=!1;for(let l in this.morphAttributes){let c=this.morphAttributes[l],f=[];for(let p=0,u=c.length;p<u;p++){let d=c[p];f.push(d.toJSON(t.data))}f.length>0&&(a[l]=f,s=!0)}s&&(t.data.morphAttributes=a,t.data.morphTargetsRelative=this.morphTargetsRelative);let r=this.groups;r.length>0&&(t.data.groups=JSON.parse(JSON.stringify(r)));let o=this.boundingSphere;return o!==null&&(t.data.boundingSphere=o.toJSON()),t}clone(){return new this.constructor().copy(this)}copy(t){this.index=null,this.attributes={},this.morphAttributes={},this.groups=[],this.boundingBox=null,this.boundingSphere=null;let n={};this.name=t.name;let i=t.index;i!==null&&this.setIndex(i.clone());let a=t.attributes;for(let c in a){let f=a[c];this.setAttribute(c,f.clone(n))}let s=t.morphAttributes;for(let c in s){let f=[],p=s[c];for(let u=0,d=p.length;u<d;u++)f.push(p[u].clone(n));this.morphAttributes[c]=f}this.morphTargetsRelative=t.morphTargetsRelative;let r=t.groups;for(let c=0,f=r.length;c<f;c++){let p=r[c];this.addGroup(p.start,p.count,p.materialIndex)}let o=t.boundingBox;o!==null&&(this.boundingBox=o.clone());let l=t.boundingSphere;return l!==null&&(this.boundingSphere=l.clone()),this.drawRange.start=t.drawRange.start,this.drawRange.count=t.drawRange.count,this.userData=t.userData,this}dispose(){this.dispatchEvent({type:"dispose"})}};var SC=0,Za=class extends qa{constructor(){super(),this.isMaterial=!0,Object.defineProperty(this,"id",{value:SC++}),this.uuid=gl(),this.name="",this.type="Material",this.blending=wr,this.side=Wa,this.vertexColors=!1,this.opacity=1,this.transparent=!1,this.alphaHash=!1,this.blendSrc=Bh,this.blendDst=zh,this.blendEquation=zs,this.blendSrcAlpha=null,this.blendDstAlpha=null,this.blendEquationAlpha=null,this.blendColor=new Ft(0,0,0),this.blendAlpha=0,this.depthFunc=Cr,this.depthTest=!0,this.depthWrite=!0,this.stencilWriteMask=255,this.stencilFunc=J0,this.stencilRef=0,this.stencilFuncMask=255,this.stencilFail=Ar,this.stencilZFail=Ar,this.stencilZPass=Ar,this.stencilWrite=!1,this.clippingPlanes=null,this.clipIntersection=!1,this.clipShadows=!1,this.shadowSide=null,this.colorWrite=!0,this.precision=null,this.polygonOffset=!1,this.polygonOffsetFactor=0,this.polygonOffsetUnits=0,this.dithering=!1,this.alphaToCoverage=!1,this.premultipliedAlpha=!1,this.forceSinglePass=!1,this.allowOverride=!0,this.visible=!0,this.toneMapped=!0,this.userData={},this.version=0,this._alphaTest=0}get alphaTest(){return this._alphaTest}set alphaTest(t){this._alphaTest>0!=t>0&&this.version++,this._alphaTest=t}onBeforeRender(){}onBeforeCompile(){}customProgramCacheKey(){return this.onBeforeCompile.toString()}setValues(t){if(t!==void 0)for(let n in t){let i=t[n];if(i===void 0){Ot(`Material: parameter '${n}' has value of undefined.`);continue}let a=this[n];if(a===void 0){Ot(`Material: '${n}' is not a property of THREE.${this.type}.`);continue}a&&a.isColor?a.set(i):a&&a.isVector3&&i&&i.isVector3?a.copy(i):this[n]=i}}toJSON(t){let n=t===void 0||typeof t=="string";n&&(t={textures:{},images:{}});let i={metadata:{version:4.7,type:"Material",generator:"Material.toJSON"}};i.uuid=this.uuid,i.type=this.type,this.name!==""&&(i.name=this.name),this.color&&this.color.isColor&&(i.color=this.color.getHex()),this.roughness!==void 0&&(i.roughness=this.roughness),this.metalness!==void 0&&(i.metalness=this.metalness),this.sheen!==void 0&&(i.sheen=this.sheen),this.sheenColor&&this.sheenColor.isColor&&(i.sheenColor=this.sheenColor.getHex()),this.sheenRoughness!==void 0&&(i.sheenRoughness=this.sheenRoughness),this.emissive&&this.emissive.isColor&&(i.emissive=this.emissive.getHex()),this.emissiveIntensity!==void 0&&this.emissiveIntensity!==1&&(i.emissiveIntensity=this.emissiveIntensity),this.specular&&this.specular.isColor&&(i.specular=this.specular.getHex()),this.specularIntensity!==void 0&&(i.specularIntensity=this.specularIntensity),this.specularColor&&this.specularColor.isColor&&(i.specularColor=this.specularColor.getHex()),this.shininess!==void 0&&(i.shininess=this.shininess),this.clearcoat!==void 0&&(i.clearcoat=this.clearcoat),this.clearcoatRoughness!==void 0&&(i.clearcoatRoughness=this.clearcoatRoughness),this.clearcoatMap&&this.clearcoatMap.isTexture&&(i.clearcoatMap=this.clearcoatMap.toJSON(t).uuid),this.clearcoatRoughnessMap&&this.clearcoatRoughnessMap.isTexture&&(i.clearcoatRoughnessMap=this.clearcoatRoughnessMap.toJSON(t).uuid),this.clearcoatNormalMap&&this.clearcoatNormalMap.isTexture&&(i.clearcoatNormalMap=this.clearcoatNormalMap.toJSON(t).uuid,i.clearcoatNormalScale=this.clearcoatNormalScale.toArray()),this.sheenColorMap&&this.sheenColorMap.isTexture&&(i.sheenColorMap=this.sheenColorMap.toJSON(t).uuid),this.sheenRoughnessMap&&this.sheenRoughnessMap.isTexture&&(i.sheenRoughnessMap=this.sheenRoughnessMap.toJSON(t).uuid),this.dispersion!==void 0&&(i.dispersion=this.dispersion),this.iridescence!==void 0&&(i.iridescence=this.iridescence),this.iridescenceIOR!==void 0&&(i.iridescenceIOR=this.iridescenceIOR),this.iridescenceThicknessRange!==void 0&&(i.iridescenceThicknessRange=this.iridescenceThicknessRange),this.iridescenceMap&&this.iridescenceMap.isTexture&&(i.iridescenceMap=this.iridescenceMap.toJSON(t).uuid),this.iridescenceThicknessMap&&this.iridescenceThicknessMap.isTexture&&(i.iridescenceThicknessMap=this.iridescenceThicknessMap.toJSON(t).uuid),this.anisotropy!==void 0&&(i.anisotropy=this.anisotropy),this.anisotropyRotation!==void 0&&(i.anisotropyRotation=this.anisotropyRotation),this.anisotropyMap&&this.anisotropyMap.isTexture&&(i.anisotropyMap=this.anisotropyMap.toJSON(t).uuid),this.map&&this.map.isTexture&&(i.map=this.map.toJSON(t).uuid),this.matcap&&this.matcap.isTexture&&(i.matcap=this.matcap.toJSON(t).uuid),this.alphaMap&&this.alphaMap.isTexture&&(i.alphaMap=this.alphaMap.toJSON(t).uuid),this.lightMap&&this.lightMap.isTexture&&(i.lightMap=this.lightMap.toJSON(t).uuid,i.lightMapIntensity=this.lightMapIntensity),this.aoMap&&this.aoMap.isTexture&&(i.aoMap=this.aoMap.toJSON(t).uuid,i.aoMapIntensity=this.aoMapIntensity),this.bumpMap&&this.bumpMap.isTexture&&(i.bumpMap=this.bumpMap.toJSON(t).uuid,i.bumpScale=this.bumpScale),this.normalMap&&this.normalMap.isTexture&&(i.normalMap=this.normalMap.toJSON(t).uuid,i.normalMapType=this.normalMapType,i.normalScale=this.normalScale.toArray()),this.displacementMap&&this.displacementMap.isTexture&&(i.displacementMap=this.displacementMap.toJSON(t).uuid,i.displacementScale=this.displacementScale,i.displacementBias=this.displacementBias),this.roughnessMap&&this.roughnessMap.isTexture&&(i.roughnessMap=this.roughnessMap.toJSON(t).uuid),this.metalnessMap&&this.metalnessMap.isTexture&&(i.metalnessMap=this.metalnessMap.toJSON(t).uuid),this.emissiveMap&&this.emissiveMap.isTexture&&(i.emissiveMap=this.emissiveMap.toJSON(t).uuid),this.specularMap&&this.specularMap.isTexture&&(i.specularMap=this.specularMap.toJSON(t).uuid),this.specularIntensityMap&&this.specularIntensityMap.isTexture&&(i.specularIntensityMap=this.specularIntensityMap.toJSON(t).uuid),this.specularColorMap&&this.specularColorMap.isTexture&&(i.specularColorMap=this.specularColorMap.toJSON(t).uuid),this.envMap&&this.envMap.isTexture&&(i.envMap=this.envMap.toJSON(t).uuid,this.combine!==void 0&&(i.combine=this.combine)),this.envMapRotation!==void 0&&(i.envMapRotation=this.envMapRotation.toArray()),this.envMapIntensity!==void 0&&(i.envMapIntensity=this.envMapIntensity),this.reflectivity!==void 0&&(i.reflectivity=this.reflectivity),this.refractionRatio!==void 0&&(i.refractionRatio=this.refractionRatio),this.gradientMap&&this.gradientMap.isTexture&&(i.gradientMap=this.gradientMap.toJSON(t).uuid),this.transmission!==void 0&&(i.transmission=this.transmission),this.transmissionMap&&this.transmissionMap.isTexture&&(i.transmissionMap=this.transmissionMap.toJSON(t).uuid),this.thickness!==void 0&&(i.thickness=this.thickness),this.thicknessMap&&this.thicknessMap.isTexture&&(i.thicknessMap=this.thicknessMap.toJSON(t).uuid),this.attenuationDistance!==void 0&&this.attenuationDistance!==1/0&&(i.attenuationDistance=this.attenuationDistance),this.attenuationColor!==void 0&&(i.attenuationColor=this.attenuationColor.getHex()),this.size!==void 0&&(i.size=this.size),this.shadowSide!==null&&(i.shadowSide=this.shadowSide),this.sizeAttenuation!==void 0&&(i.sizeAttenuation=this.sizeAttenuation),this.blending!==wr&&(i.blending=this.blending),this.side!==Wa&&(i.side=this.side),this.vertexColors===!0&&(i.vertexColors=!0),this.opacity<1&&(i.opacity=this.opacity),this.transparent===!0&&(i.transparent=!0),this.blendSrc!==Bh&&(i.blendSrc=this.blendSrc),this.blendDst!==zh&&(i.blendDst=this.blendDst),this.blendEquation!==zs&&(i.blendEquation=this.blendEquation),this.blendSrcAlpha!==null&&(i.blendSrcAlpha=this.blendSrcAlpha),this.blendDstAlpha!==null&&(i.blendDstAlpha=this.blendDstAlpha),this.blendEquationAlpha!==null&&(i.blendEquationAlpha=this.blendEquationAlpha),this.blendColor&&this.blendColor.isColor&&(i.blendColor=this.blendColor.getHex()),this.blendAlpha!==0&&(i.blendAlpha=this.blendAlpha),this.depthFunc!==Cr&&(i.depthFunc=this.depthFunc),this.depthTest===!1&&(i.depthTest=this.depthTest),this.depthWrite===!1&&(i.depthWrite=this.depthWrite),this.colorWrite===!1&&(i.colorWrite=this.colorWrite),this.stencilWriteMask!==255&&(i.stencilWriteMask=this.stencilWriteMask),this.stencilFunc!==J0&&(i.stencilFunc=this.stencilFunc),this.stencilRef!==0&&(i.stencilRef=this.stencilRef),this.stencilFuncMask!==255&&(i.stencilFuncMask=this.stencilFuncMask),this.stencilFail!==Ar&&(i.stencilFail=this.stencilFail),this.stencilZFail!==Ar&&(i.stencilZFail=this.stencilZFail),this.stencilZPass!==Ar&&(i.stencilZPass=this.stencilZPass),this.stencilWrite===!0&&(i.stencilWrite=this.stencilWrite),this.rotation!==void 0&&this.rotation!==0&&(i.rotation=this.rotation),this.polygonOffset===!0&&(i.polygonOffset=!0),this.polygonOffsetFactor!==0&&(i.polygonOffsetFactor=this.polygonOffsetFactor),this.polygonOffsetUnits!==0&&(i.polygonOffsetUnits=this.polygonOffsetUnits),this.linewidth!==void 0&&this.linewidth!==1&&(i.linewidth=this.linewidth),this.dashSize!==void 0&&(i.dashSize=this.dashSize),this.gapSize!==void 0&&(i.gapSize=this.gapSize),this.scale!==void 0&&(i.scale=this.scale),this.dithering===!0&&(i.dithering=!0),this.alphaTest>0&&(i.alphaTest=this.alphaTest),this.alphaHash===!0&&(i.alphaHash=!0),this.alphaToCoverage===!0&&(i.alphaToCoverage=!0),this.premultipliedAlpha===!0&&(i.premultipliedAlpha=!0),this.forceSinglePass===!0&&(i.forceSinglePass=!0),this.allowOverride===!1&&(i.allowOverride=!1),this.wireframe===!0&&(i.wireframe=!0),this.wireframeLinewidth>1&&(i.wireframeLinewidth=this.wireframeLinewidth),this.wireframeLinecap!=="round"&&(i.wireframeLinecap=this.wireframeLinecap),this.wireframeLinejoin!=="round"&&(i.wireframeLinejoin=this.wireframeLinejoin),this.flatShading===!0&&(i.flatShading=!0),this.visible===!1&&(i.visible=!1),this.toneMapped===!1&&(i.toneMapped=!1),this.fog===!1&&(i.fog=!1),Object.keys(this.userData).length>0&&(i.userData=this.userData);function a(s){let r=[];for(let o in s){let l=s[o];delete l.metadata,r.push(l)}return r}if(n){let s=a(t.textures),r=a(t.images);s.length>0&&(i.textures=s),r.length>0&&(i.images=r)}return i}clone(){return new this.constructor().copy(this)}copy(t){this.name=t.name,this.blending=t.blending,this.side=t.side,this.vertexColors=t.vertexColors,this.opacity=t.opacity,this.transparent=t.transparent,this.blendSrc=t.blendSrc,this.blendDst=t.blendDst,this.blendEquation=t.blendEquation,this.blendSrcAlpha=t.blendSrcAlpha,this.blendDstAlpha=t.blendDstAlpha,this.blendEquationAlpha=t.blendEquationAlpha,this.blendColor.copy(t.blendColor),this.blendAlpha=t.blendAlpha,this.depthFunc=t.depthFunc,this.depthTest=t.depthTest,this.depthWrite=t.depthWrite,this.stencilWriteMask=t.stencilWriteMask,this.stencilFunc=t.stencilFunc,this.stencilRef=t.stencilRef,this.stencilFuncMask=t.stencilFuncMask,this.stencilFail=t.stencilFail,this.stencilZFail=t.stencilZFail,this.stencilZPass=t.stencilZPass,this.stencilWrite=t.stencilWrite;let n=t.clippingPlanes,i=null;if(n!==null){let a=n.length;i=new Array(a);for(let s=0;s!==a;++s)i[s]=n[s].clone()}return this.clippingPlanes=i,this.clipIntersection=t.clipIntersection,this.clipShadows=t.clipShadows,this.shadowSide=t.shadowSide,this.colorWrite=t.colorWrite,this.precision=t.precision,this.polygonOffset=t.polygonOffset,this.polygonOffsetFactor=t.polygonOffsetFactor,this.polygonOffsetUnits=t.polygonOffsetUnits,this.dithering=t.dithering,this.alphaTest=t.alphaTest,this.alphaHash=t.alphaHash,this.alphaToCoverage=t.alphaToCoverage,this.premultipliedAlpha=t.premultipliedAlpha,this.forceSinglePass=t.forceSinglePass,this.allowOverride=t.allowOverride,this.visible=t.visible,this.toneMapped=t.toneMapped,this.userData=JSON.parse(JSON.stringify(t.userData)),this}dispose(){this.dispatchEvent({type:"dispose"})}set needsUpdate(t){t===!0&&this.version++}};var ka=new B,z0=new B,xh=new B,Os=new B,F0=new B,Sh=new B,V0=new B,Jc=class{constructor(t=new B,n=new B(0,0,-1)){this.origin=t,this.direction=n}set(t,n){return this.origin.copy(t),this.direction.copy(n),this}copy(t){return this.origin.copy(t.origin),this.direction.copy(t.direction),this}at(t,n){return n.copy(this.origin).addScaledVector(this.direction,t)}lookAt(t){return this.direction.copy(t).sub(this.origin).normalize(),this}recast(t){return this.origin.copy(this.at(t,ka)),this}closestPointToPoint(t,n){n.subVectors(t,this.origin);let i=n.dot(this.direction);return i<0?n.copy(this.origin):n.copy(this.origin).addScaledVector(this.direction,i)}distanceToPoint(t){return Math.sqrt(this.distanceSqToPoint(t))}distanceSqToPoint(t){let n=ka.subVectors(t,this.origin).dot(this.direction);return n<0?this.origin.distanceToSquared(t):(ka.copy(this.origin).addScaledVector(this.direction,n),ka.distanceToSquared(t))}distanceSqToSegment(t,n,i,a){z0.copy(t).add(n).multiplyScalar(.5),xh.copy(n).sub(t).normalize(),Os.copy(this.origin).sub(z0);let s=t.distanceTo(n)*.5,r=-this.direction.dot(xh),o=Os.dot(this.direction),l=-Os.dot(xh),c=Os.lengthSq(),f=Math.abs(1-r*r),p,u,d,_;if(f>0)if(p=r*l-o,u=r*o-l,_=s*f,p>=0)if(u>=-_)if(u<=_){let b=1/f;p*=b,u*=b,d=p*(p+r*u+2*o)+u*(r*p+u+2*l)+c}else u=s,p=Math.max(0,-(r*u+o)),d=-p*p+u*(u+2*l)+c;else u=-s,p=Math.max(0,-(r*u+o)),d=-p*p+u*(u+2*l)+c;else u<=-_?(p=Math.max(0,-(-r*s+o)),u=p>0?-s:Math.min(Math.max(-s,-l),s),d=-p*p+u*(u+2*l)+c):u<=_?(p=0,u=Math.min(Math.max(-s,-l),s),d=u*(u+2*l)+c):(p=Math.max(0,-(r*s+o)),u=p>0?s:Math.min(Math.max(-s,-l),s),d=-p*p+u*(u+2*l)+c);else u=r>0?-s:s,p=Math.max(0,-(r*u+o)),d=-p*p+u*(u+2*l)+c;return i&&i.copy(this.origin).addScaledVector(this.direction,p),a&&a.copy(z0).addScaledVector(xh,u),d}intersectSphere(t,n){ka.subVectors(t.center,this.origin);let i=ka.dot(this.direction),a=ka.dot(ka)-i*i,s=t.radius*t.radius;if(a>s)return null;let r=Math.sqrt(s-a),o=i-r,l=i+r;return l<0?null:o<0?this.at(l,n):this.at(o,n)}intersectsSphere(t){return t.radius<0?!1:this.distanceSqToPoint(t.center)<=t.radius*t.radius}distanceToPlane(t){let n=t.normal.dot(this.direction);if(n===0)return t.distanceToPoint(this.origin)===0?0:null;let i=-(this.origin.dot(t.normal)+t.constant)/n;return i>=0?i:null}intersectPlane(t,n){let i=this.distanceToPlane(t);return i===null?null:this.at(i,n)}intersectsPlane(t){let n=t.distanceToPoint(this.origin);return n===0||t.normal.dot(this.direction)*n<0}intersectBox(t,n){let i,a,s,r,o,l,c=1/this.direction.x,f=1/this.direction.y,p=1/this.direction.z,u=this.origin;return c>=0?(i=(t.min.x-u.x)*c,a=(t.max.x-u.x)*c):(i=(t.max.x-u.x)*c,a=(t.min.x-u.x)*c),f>=0?(s=(t.min.y-u.y)*f,r=(t.max.y-u.y)*f):(s=(t.max.y-u.y)*f,r=(t.min.y-u.y)*f),i>r||s>a||((s>i||isNaN(i))&&(i=s),(r<a||isNaN(a))&&(a=r),p>=0?(o=(t.min.z-u.z)*p,l=(t.max.z-u.z)*p):(o=(t.max.z-u.z)*p,l=(t.min.z-u.z)*p),i>l||o>a)||((o>i||i!==i)&&(i=o),(l<a||a!==a)&&(a=l),a<0)?null:this.at(i>=0?i:a,n)}intersectsBox(t){return this.intersectBox(t,ka)!==null}intersectTriangle(t,n,i,a,s){F0.subVectors(n,t),Sh.subVectors(i,t),V0.crossVectors(F0,Sh);let r=this.direction.dot(V0),o;if(r>0){if(a)return null;o=1}else if(r<0)o=-1,r=-r;else return null;Os.subVectors(this.origin,t);let l=o*this.direction.dot(Sh.crossVectors(Os,Sh));if(l<0)return null;let c=o*this.direction.dot(F0.cross(Os));if(c<0||l+c>r)return null;let f=-o*Os.dot(V0);return f<0?null:this.at(f/r,s)}applyMatrix4(t){return this.origin.applyMatrix4(t),this.direction.transformDirection(t),this}equals(t){return t.origin.equals(this.origin)&&t.direction.equals(this.direction)}clone(){return new this.constructor().copy(this)}},Dr=class extends Za{constructor(t){super(),this.isMeshBasicMaterial=!0,this.type="MeshBasicMaterial",this.color=new Ft(16777215),this.map=null,this.lightMap=null,this.lightMapIntensity=1,this.aoMap=null,this.aoMapIntensity=1,this.specularMap=null,this.alphaMap=null,this.envMap=null,this.envMapRotation=new Xi,this.combine=av,this.reflectivity=1,this.refractionRatio=.98,this.wireframe=!1,this.wireframeLinewidth=1,this.wireframeLinecap="round",this.wireframeLinejoin="round",this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.map=t.map,this.lightMap=t.lightMap,this.lightMapIntensity=t.lightMapIntensity,this.aoMap=t.aoMap,this.aoMapIntensity=t.aoMapIntensity,this.specularMap=t.specularMap,this.alphaMap=t.alphaMap,this.envMap=t.envMap,this.envMapRotation.copy(t.envMapRotation),this.combine=t.combine,this.reflectivity=t.reflectivity,this.refractionRatio=t.refractionRatio,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.wireframeLinecap=t.wireframeLinecap,this.wireframeLinejoin=t.wireframeLinejoin,this.fog=t.fog,this}},wM=new Ie,Tr=new Jc,bh=new Ya,CM=new B,Mh=new B,Th=new B,Eh=new B,H0=new B,Ah=new B,RM=new B,wh=new B,mn=class extends tn{constructor(t=new An,n=new Dr){super(),this.isMesh=!0,this.type="Mesh",this.geometry=t,this.material=n,this.morphTargetDictionary=void 0,this.morphTargetInfluences=void 0,this.count=1,this.updateMorphTargets()}copy(t,n){return super.copy(t,n),t.morphTargetInfluences!==void 0&&(this.morphTargetInfluences=t.morphTargetInfluences.slice()),t.morphTargetDictionary!==void 0&&(this.morphTargetDictionary=Object.assign({},t.morphTargetDictionary)),this.material=Array.isArray(t.material)?t.material.slice():t.material,this.geometry=t.geometry,this}updateMorphTargets(){let n=this.geometry.morphAttributes,i=Object.keys(n);if(i.length>0){let a=n[i[0]];if(a!==void 0){this.morphTargetInfluences=[],this.morphTargetDictionary={};for(let s=0,r=a.length;s<r;s++){let o=a[s].name||String(s);this.morphTargetInfluences.push(0),this.morphTargetDictionary[o]=s}}}}getVertexPosition(t,n){let i=this.geometry,a=i.attributes.position,s=i.morphAttributes.position,r=i.morphTargetsRelative;n.fromBufferAttribute(a,t);let o=this.morphTargetInfluences;if(s&&o){Ah.set(0,0,0);for(let l=0,c=s.length;l<c;l++){let f=o[l],p=s[l];f!==0&&(H0.fromBufferAttribute(p,t),r?Ah.addScaledVector(H0,f):Ah.addScaledVector(H0.sub(n),f))}n.add(Ah)}return n}raycast(t,n){let i=this.geometry,a=this.material,s=this.matrixWorld;a!==void 0&&(i.boundingSphere===null&&i.computeBoundingSphere(),bh.copy(i.boundingSphere),bh.applyMatrix4(s),Tr.copy(t.ray).recast(t.near),!(bh.containsPoint(Tr.origin)===!1&&(Tr.intersectSphere(bh,CM)===null||Tr.origin.distanceToSquared(CM)>(t.far-t.near)**2))&&(wM.copy(s).invert(),Tr.copy(t.ray).applyMatrix4(wM),!(i.boundingBox!==null&&Tr.intersectsBox(i.boundingBox)===!1)&&this._computeIntersections(t,n,Tr)))}_computeIntersections(t,n,i){let a,s=this.geometry,r=this.material,o=s.index,l=s.attributes.position,c=s.attributes.uv,f=s.attributes.uv1,p=s.attributes.normal,u=s.groups,d=s.drawRange;if(o!==null)if(Array.isArray(r))for(let _=0,b=u.length;_<b;_++){let v=u[_],h=r[v.materialIndex],m=Math.max(v.start,d.start),g=Math.min(o.count,Math.min(v.start+v.count,d.start+d.count));for(let x=m,E=g;x<E;x+=3){let A=o.getX(x),w=o.getX(x+1),y=o.getX(x+2);a=Ch(this,h,t,i,c,f,p,A,w,y),a&&(a.faceIndex=Math.floor(x/3),a.face.materialIndex=v.materialIndex,n.push(a))}}else{let _=Math.max(0,d.start),b=Math.min(o.count,d.start+d.count);for(let v=_,h=b;v<h;v+=3){let m=o.getX(v),g=o.getX(v+1),x=o.getX(v+2);a=Ch(this,r,t,i,c,f,p,m,g,x),a&&(a.faceIndex=Math.floor(v/3),n.push(a))}}else if(l!==void 0)if(Array.isArray(r))for(let _=0,b=u.length;_<b;_++){let v=u[_],h=r[v.materialIndex],m=Math.max(v.start,d.start),g=Math.min(l.count,Math.min(v.start+v.count,d.start+d.count));for(let x=m,E=g;x<E;x+=3){let A=x,w=x+1,y=x+2;a=Ch(this,h,t,i,c,f,p,A,w,y),a&&(a.faceIndex=Math.floor(x/3),a.face.materialIndex=v.materialIndex,n.push(a))}}else{let _=Math.max(0,d.start),b=Math.min(l.count,d.start+d.count);for(let v=_,h=b;v<h;v+=3){let m=v,g=v+1,x=v+2;a=Ch(this,r,t,i,c,f,p,m,g,x),a&&(a.faceIndex=Math.floor(v/3),n.push(a))}}}};function bC(e,t,n,i,a,s,r,o){let l;if(t.side===wn?l=i.intersectTriangle(r,s,a,!0,o):l=i.intersectTriangle(a,s,r,t.side===Wa,o),l===null)return null;wh.copy(o),wh.applyMatrix4(e.matrixWorld);let c=n.ray.origin.distanceTo(wh);return c<n.near||c>n.far?null:{distance:c,point:wh.clone(),object:e}}function Ch(e,t,n,i,a,s,r,o,l,c){e.getVertexPosition(o,Mh),e.getVertexPosition(l,Th),e.getVertexPosition(c,Eh);let f=bC(e,t,n,i,Mh,Th,Eh,RM);if(f){let p=new B;Bs.getBarycoord(RM,Mh,Th,Eh,p),a&&(f.uv=Bs.getInterpolatedAttribute(a,o,l,c,p,new te)),s&&(f.uv1=Bs.getInterpolatedAttribute(s,o,l,c,p,new te)),r&&(f.normal=Bs.getInterpolatedAttribute(r,o,l,c,p,new B),f.normal.dot(i.direction)>0&&f.normal.multiplyScalar(-1));let u={a:o,b:l,c,normal:new B,materialIndex:0};Bs.getNormal(Mh,Th,Eh,u.normal),f.face=u,f.barycoord=p}return f}var Kc=class extends qn{constructor(t=null,n=1,i=1,a,s,r,o,l,c=pn,f=pn,p,u){super(null,r,o,l,c,f,a,s,p,u),this.isDataTexture=!0,this.image={data:t,width:n,height:i},this.generateMipmaps=!1,this.flipY=!1,this.unpackAlignment=1}};var Qc=class extends Bn{constructor(t,n,i,a=1){super(t,n,i),this.isInstancedBufferAttribute=!0,this.meshPerAttribute=a}copy(t){return super.copy(t),this.meshPerAttribute=t.meshPerAttribute,this}toJSON(){let t=super.toJSON();return t.meshPerAttribute=this.meshPerAttribute,t.isInstancedBufferAttribute=!0,t}},Qo=new Ie,DM=new Ie,Rh=[],NM=new ma,MC=new Ie,Ic=new mn,Oc=new Ya,jc=class extends mn{constructor(t,n,i){super(t,n),this.isInstancedMesh=!0,this.instanceMatrix=new Qc(new Float32Array(i*16),16),this.previousInstanceMatrix=null,this.instanceColor=null,this.morphTexture=null,this.count=i,this.boundingBox=null,this.boundingSphere=null;for(let a=0;a<i;a++)this.setMatrixAt(a,MC)}computeBoundingBox(){let t=this.geometry,n=this.count;this.boundingBox===null&&(this.boundingBox=new ma),t.boundingBox===null&&t.computeBoundingBox(),this.boundingBox.makeEmpty();for(let i=0;i<n;i++)this.getMatrixAt(i,Qo),NM.copy(t.boundingBox).applyMatrix4(Qo),this.boundingBox.union(NM)}computeBoundingSphere(){let t=this.geometry,n=this.count;this.boundingSphere===null&&(this.boundingSphere=new Ya),t.boundingSphere===null&&t.computeBoundingSphere(),this.boundingSphere.makeEmpty();for(let i=0;i<n;i++)this.getMatrixAt(i,Qo),Oc.copy(t.boundingSphere).applyMatrix4(Qo),this.boundingSphere.union(Oc)}copy(t,n){return super.copy(t,n),this.instanceMatrix.copy(t.instanceMatrix),t.previousInstanceMatrix!==null&&(this.previousInstanceMatrix=t.previousInstanceMatrix.clone()),t.morphTexture!==null&&(this.morphTexture=t.morphTexture.clone()),t.instanceColor!==null&&(this.instanceColor=t.instanceColor.clone()),this.count=t.count,t.boundingBox!==null&&(this.boundingBox=t.boundingBox.clone()),t.boundingSphere!==null&&(this.boundingSphere=t.boundingSphere.clone()),this}getColorAt(t,n){n.fromArray(this.instanceColor.array,t*3)}getMatrixAt(t,n){n.fromArray(this.instanceMatrix.array,t*16)}getMorphAt(t,n){let i=n.morphTargetInfluences,a=this.morphTexture.source.data.data,s=i.length+1,r=t*s+1;for(let o=0;o<i.length;o++)i[o]=a[r+o]}raycast(t,n){let i=this.matrixWorld,a=this.count;if(Ic.geometry=this.geometry,Ic.material=this.material,Ic.material!==void 0&&(this.boundingSphere===null&&this.computeBoundingSphere(),Oc.copy(this.boundingSphere),Oc.applyMatrix4(i),t.ray.intersectsSphere(Oc)!==!1))for(let s=0;s<a;s++){this.getMatrixAt(s,Qo),DM.multiplyMatrices(i,Qo),Ic.matrixWorld=DM,Ic.raycast(t,Rh);for(let r=0,o=Rh.length;r<o;r++){let l=Rh[r];l.instanceId=s,l.object=this,n.push(l)}Rh.length=0}}setColorAt(t,n){this.instanceColor===null&&(this.instanceColor=new Qc(new Float32Array(this.instanceMatrix.count*3).fill(1),3)),n.toArray(this.instanceColor.array,t*3)}setMatrixAt(t,n){n.toArray(this.instanceMatrix.array,t*16)}setMorphAt(t,n){let i=n.morphTargetInfluences,a=i.length+1;this.morphTexture===null&&(this.morphTexture=new Kc(new Float32Array(a*this.count),a,this.count,Td,Ni));let s=this.morphTexture.source.data.data,r=0;for(let c=0;c<i.length;c++)r+=i[c];let o=this.geometry.morphTargetsRelative?1:1-r,l=a*t;s[l]=o,s.set(i,l+1)}updateMorphTargets(){}dispose(){this.dispatchEvent({type:"dispose"}),this.morphTexture!==null&&(this.morphTexture.dispose(),this.morphTexture=null)}},G0=new B,TC=new B,EC=new Zt,fa=class{constructor(t=new B(1,0,0),n=0){this.isPlane=!0,this.normal=t,this.constant=n}set(t,n){return this.normal.copy(t),this.constant=n,this}setComponents(t,n,i,a){return this.normal.set(t,n,i),this.constant=a,this}setFromNormalAndCoplanarPoint(t,n){return this.normal.copy(t),this.constant=-n.dot(this.normal),this}setFromCoplanarPoints(t,n,i){let a=G0.subVectors(i,n).cross(TC.subVectors(t,n)).normalize();return this.setFromNormalAndCoplanarPoint(a,t),this}copy(t){return this.normal.copy(t.normal),this.constant=t.constant,this}normalize(){let t=1/this.normal.length();return this.normal.multiplyScalar(t),this.constant*=t,this}negate(){return this.constant*=-1,this.normal.negate(),this}distanceToPoint(t){return this.normal.dot(t)+this.constant}distanceToSphere(t){return this.distanceToPoint(t.center)-t.radius}projectPoint(t,n){return n.copy(t).addScaledVector(this.normal,-this.distanceToPoint(t))}intersectLine(t,n){let i=t.delta(G0),a=this.normal.dot(i);if(a===0)return this.distanceToPoint(t.start)===0?n.copy(t.start):null;let s=-(t.start.dot(this.normal)+this.constant)/a;return s<0||s>1?null:n.copy(t.start).addScaledVector(i,s)}intersectsLine(t){let n=this.distanceToPoint(t.start),i=this.distanceToPoint(t.end);return n<0&&i>0||i<0&&n>0}intersectsBox(t){return t.intersectsPlane(this)}intersectsSphere(t){return t.intersectsPlane(this)}coplanarPoint(t){return t.copy(this.normal).multiplyScalar(-this.constant)}applyMatrix4(t,n){let i=n||EC.getNormalMatrix(t),a=this.coplanarPoint(G0).applyMatrix4(t),s=this.normal.applyMatrix3(i).normalize();return this.constant=-a.dot(s),this}translate(t){return this.constant-=t.dot(this.normal),this}equals(t){return t.normal.equals(this.normal)&&t.constant===this.constant}clone(){return new this.constructor().copy(this)}},Er=new Ya,AC=new te(.5,.5),Dh=new B,ol=class{constructor(t=new fa,n=new fa,i=new fa,a=new fa,s=new fa,r=new fa){this.planes=[t,n,i,a,s,r]}set(t,n,i,a,s,r){let o=this.planes;return o[0].copy(t),o[1].copy(n),o[2].copy(i),o[3].copy(a),o[4].copy(s),o[5].copy(r),this}copy(t){let n=this.planes;for(let i=0;i<6;i++)n[i].copy(t.planes[i]);return this}setFromProjectionMatrix(t,n=ki,i=!1){let a=this.planes,s=t.elements,r=s[0],o=s[1],l=s[2],c=s[3],f=s[4],p=s[5],u=s[6],d=s[7],_=s[8],b=s[9],v=s[10],h=s[11],m=s[12],g=s[13],x=s[14],E=s[15];if(a[0].setComponents(c-r,d-f,h-_,E-m).normalize(),a[1].setComponents(c+r,d+f,h+_,E+m).normalize(),a[2].setComponents(c+o,d+p,h+b,E+g).normalize(),a[3].setComponents(c-o,d-p,h-b,E-g).normalize(),i)a[4].setComponents(l,u,v,x).normalize(),a[5].setComponents(c-l,d-u,h-v,E-x).normalize();else if(a[4].setComponents(c-l,d-u,h-v,E-x).normalize(),n===ki)a[5].setComponents(c+l,d+u,h+v,E+x).normalize();else if(n===nl)a[5].setComponents(l,u,v,x).normalize();else throw new Error("THREE.Frustum.setFromProjectionMatrix(): Invalid coordinate system: "+n);return this}intersectsObject(t){if(t.boundingSphere!==void 0)t.boundingSphere===null&&t.computeBoundingSphere(),Er.copy(t.boundingSphere).applyMatrix4(t.matrixWorld);else{let n=t.geometry;n.boundingSphere===null&&n.computeBoundingSphere(),Er.copy(n.boundingSphere).applyMatrix4(t.matrixWorld)}return this.intersectsSphere(Er)}intersectsSprite(t){Er.center.set(0,0,0);let n=AC.distanceTo(t.center);return Er.radius=.7071067811865476+n,Er.applyMatrix4(t.matrixWorld),this.intersectsSphere(Er)}intersectsSphere(t){let n=this.planes,i=t.center,a=-t.radius;for(let s=0;s<6;s++)if(n[s].distanceToPoint(i)<a)return!1;return!0}intersectsBox(t){let n=this.planes;for(let i=0;i<6;i++){let a=n[i];if(Dh.x=a.normal.x>0?t.max.x:t.min.x,Dh.y=a.normal.y>0?t.max.y:t.min.y,Dh.z=a.normal.z>0?t.max.z:t.min.z,a.distanceToPoint(Dh)<0)return!1}return!0}containsPoint(t){let n=this.planes;for(let i=0;i<6;i++)if(n[i].distanceToPoint(t)<0)return!1;return!0}clone(){return new this.constructor().copy(this)}};var ll=class extends Za{constructor(t){super(),this.isLineBasicMaterial=!0,this.type="LineBasicMaterial",this.color=new Ft(16777215),this.map=null,this.linewidth=1,this.linecap="round",this.linejoin="round",this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.map=t.map,this.linewidth=t.linewidth,this.linecap=t.linecap,this.linejoin=t.linejoin,this.fog=t.fog,this}},jh=new B,$h=new B,UM=new Ie,Pc=new Jc,Nh=new Ya,k0=new B,LM=new B,td=class extends tn{constructor(t=new An,n=new ll){super(),this.isLine=!0,this.type="Line",this.geometry=t,this.material=n,this.morphTargetDictionary=void 0,this.morphTargetInfluences=void 0,this.updateMorphTargets()}copy(t,n){return super.copy(t,n),this.material=Array.isArray(t.material)?t.material.slice():t.material,this.geometry=t.geometry,this}computeLineDistances(){let t=this.geometry;if(t.index===null){let n=t.attributes.position,i=[0];for(let a=1,s=n.count;a<s;a++)jh.fromBufferAttribute(n,a-1),$h.fromBufferAttribute(n,a),i[a]=i[a-1],i[a]+=jh.distanceTo($h);t.setAttribute("lineDistance",new Oe(i,1))}else Ot("Line.computeLineDistances(): Computation only possible with non-indexed BufferGeometry.");return this}raycast(t,n){let i=this.geometry,a=this.matrixWorld,s=t.params.Line.threshold,r=i.drawRange;if(i.boundingSphere===null&&i.computeBoundingSphere(),Nh.copy(i.boundingSphere),Nh.applyMatrix4(a),Nh.radius+=s,t.ray.intersectsSphere(Nh)===!1)return;UM.copy(a).invert(),Pc.copy(t.ray).applyMatrix4(UM);let o=s/((this.scale.x+this.scale.y+this.scale.z)/3),l=o*o,c=this.isLineSegments?2:1,f=i.index,u=i.attributes.position;if(f!==null){let d=Math.max(0,r.start),_=Math.min(f.count,r.start+r.count);for(let b=d,v=_-1;b<v;b+=c){let h=f.getX(b),m=f.getX(b+1),g=Uh(this,t,Pc,l,h,m,b);g&&n.push(g)}if(this.isLineLoop){let b=f.getX(_-1),v=f.getX(d),h=Uh(this,t,Pc,l,b,v,_-1);h&&n.push(h)}}else{let d=Math.max(0,r.start),_=Math.min(u.count,r.start+r.count);for(let b=d,v=_-1;b<v;b+=c){let h=Uh(this,t,Pc,l,b,b+1,b);h&&n.push(h)}if(this.isLineLoop){let b=Uh(this,t,Pc,l,_-1,d,_-1);b&&n.push(b)}}}updateMorphTargets(){let n=this.geometry.morphAttributes,i=Object.keys(n);if(i.length>0){let a=n[i[0]];if(a!==void 0){this.morphTargetInfluences=[],this.morphTargetDictionary={};for(let s=0,r=a.length;s<r;s++){let o=a[s].name||String(s);this.morphTargetInfluences.push(0),this.morphTargetDictionary[o]=s}}}}};function Uh(e,t,n,i,a,s,r){let o=e.geometry.attributes.position;if(jh.fromBufferAttribute(o,a),$h.fromBufferAttribute(o,s),n.distanceSqToSegment(jh,$h,k0,LM)>i)return;k0.applyMatrix4(e.matrixWorld);let c=t.ray.origin.distanceTo(k0);if(!(c<t.near||c>t.far))return{distance:c,point:LM.clone().applyMatrix4(e.matrixWorld),index:r,face:null,faceIndex:null,barycoord:null,object:e}}var IM=new B,OM=new B,$c=class extends td{constructor(t,n){super(t,n),this.isLineSegments=!0,this.type="LineSegments"}computeLineDistances(){let t=this.geometry;if(t.index===null){let n=t.attributes.position,i=[];for(let a=0,s=n.count;a<s;a+=2)IM.fromBufferAttribute(n,a),OM.fromBufferAttribute(n,a+1),i[a]=a===0?0:i[a-1],i[a+1]=i[a]+IM.distanceTo(OM);t.setAttribute("lineDistance",new Oe(i,1))}else Ot("LineSegments.computeLineDistances(): Computation only possible with non-indexed BufferGeometry.");return this}};var tu=class extends qn{constructor(t=[],n=Xs,i,a,s,r,o,l,c,f){super(t,n,i,a,s,r,o,l,c,f),this.isCubeTexture=!0,this.flipY=!1}get images(){return this.image}set images(t){this.image=t}},eu=class extends qn{constructor(t,n,i,a,s,r,o,l,c){super(t,n,i,a,s,r,o,l,c),this.isCanvasTexture=!0,this.needsUpdate=!0}},Fs=class extends qn{constructor(t,n,i=qi,a,s,r,o=pn,l=pn,c,f=da,p=1){if(f!==da&&f!==qs)throw new Error("DepthTexture format must be either THREE.DepthFormat or THREE.DepthStencilFormat");let u={width:t,height:n,depth:p};super(u,a,s,r,o,l,f,i,c),this.isDepthTexture=!0,this.flipY=!1,this.generateMipmaps=!1,this.compareFunction=null}copy(t){return super.copy(t),this.source=new sl(Object.assign({},t.image)),this.compareFunction=t.compareFunction,this}toJSON(t){let n=super.toJSON(t);return this.compareFunction!==null&&(n.compareFunction=this.compareFunction),n}},ed=class extends Fs{constructor(t,n=qi,i=Xs,a,s,r=pn,o=pn,l,c=da){let f={width:t,height:t,depth:1},p=[f,f,f,f,f,f];super(t,t,n,i,a,s,r,o,l,c),this.image=p,this.isCubeDepthTexture=!0,this.isCubeTexture=!0}get images(){return this.image}set images(t){this.image=t}},nu=class extends qn{constructor(t=null){super(),this.sourceTexture=t,this.isExternalTexture=!0}copy(t){return super.copy(t),this.sourceTexture=t.sourceTexture,this}},Vs=class e extends An{constructor(t=1,n=1,i=1,a=1,s=1,r=1){super(),this.type="BoxGeometry",this.parameters={width:t,height:n,depth:i,widthSegments:a,heightSegments:s,depthSegments:r};let o=this;a=Math.floor(a),s=Math.floor(s),r=Math.floor(r);let l=[],c=[],f=[],p=[],u=0,d=0;_("z","y","x",-1,-1,i,n,t,r,s,0),_("z","y","x",1,-1,i,n,-t,r,s,1),_("x","z","y",1,1,t,i,n,a,r,2),_("x","z","y",1,-1,t,i,-n,a,r,3),_("x","y","z",1,-1,t,n,i,a,s,4),_("x","y","z",-1,-1,t,n,-i,a,s,5),this.setIndex(l),this.setAttribute("position",new Oe(c,3)),this.setAttribute("normal",new Oe(f,3)),this.setAttribute("uv",new Oe(p,2));function _(b,v,h,m,g,x,E,A,w,y,T){let L=x/w,R=E/y,N=x/2,P=E/2,z=A/2,F=w+1,V=y+1,X=0,at=0,et=new B;for(let dt=0;dt<V;dt++){let St=dt*R-P;for(let ht=0;ht<F;ht++){let Vt=ht*L-N;et[b]=Vt*m,et[v]=St*g,et[h]=z,c.push(et.x,et.y,et.z),et[b]=0,et[v]=0,et[h]=A>0?1:-1,f.push(et.x,et.y,et.z),p.push(ht/w),p.push(1-dt/y),X+=1}}for(let dt=0;dt<y;dt++)for(let St=0;St<w;St++){let ht=u+St+F*dt,Vt=u+St+F*(dt+1),Ce=u+(St+1)+F*(dt+1),xe=u+(St+1)+F*dt;l.push(ht,Vt,xe),l.push(Vt,Ce,xe),at+=6}o.addGroup(d,at,T),d+=at,u+=X}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new e(t.width,t.height,t.depth,t.widthSegments,t.heightSegments,t.depthSegments)}};var iu=class e extends An{constructor(t=1,n=32,i=0,a=Math.PI*2){super(),this.type="CircleGeometry",this.parameters={radius:t,segments:n,thetaStart:i,thetaLength:a},n=Math.max(3,n);let s=[],r=[],o=[],l=[],c=new B,f=new te;r.push(0,0,0),o.push(0,0,1),l.push(.5,.5);for(let p=0,u=3;p<=n;p++,u+=3){let d=i+p/n*a;c.x=t*Math.cos(d),c.y=t*Math.sin(d),r.push(c.x,c.y,c.z),o.push(0,0,1),f.x=(r[u]/t+1)/2,f.y=(r[u+1]/t+1)/2,l.push(f.x,f.y)}for(let p=1;p<=n;p++)s.push(p,p+1,0);this.setIndex(s),this.setAttribute("position",new Oe(r,3)),this.setAttribute("normal",new Oe(o,3)),this.setAttribute("uv",new Oe(l,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new e(t.radius,t.segments,t.thetaStart,t.thetaLength)}},cl=class e extends An{constructor(t=1,n=1,i=1,a=32,s=1,r=!1,o=0,l=Math.PI*2){super(),this.type="CylinderGeometry",this.parameters={radiusTop:t,radiusBottom:n,height:i,radialSegments:a,heightSegments:s,openEnded:r,thetaStart:o,thetaLength:l};let c=this;a=Math.floor(a),s=Math.floor(s);let f=[],p=[],u=[],d=[],_=0,b=[],v=i/2,h=0;m(),r===!1&&(t>0&&g(!0),n>0&&g(!1)),this.setIndex(f),this.setAttribute("position",new Oe(p,3)),this.setAttribute("normal",new Oe(u,3)),this.setAttribute("uv",new Oe(d,2));function m(){let x=new B,E=new B,A=0,w=(n-t)/i;for(let y=0;y<=s;y++){let T=[],L=y/s,R=L*(n-t)+t;for(let N=0;N<=a;N++){let P=N/a,z=P*l+o,F=Math.sin(z),V=Math.cos(z);E.x=R*F,E.y=-L*i+v,E.z=R*V,p.push(E.x,E.y,E.z),x.set(F,w,V).normalize(),u.push(x.x,x.y,x.z),d.push(P,1-L),T.push(_++)}b.push(T)}for(let y=0;y<a;y++)for(let T=0;T<s;T++){let L=b[T][y],R=b[T+1][y],N=b[T+1][y+1],P=b[T][y+1];(t>0||T!==0)&&(f.push(L,R,P),A+=3),(n>0||T!==s-1)&&(f.push(R,N,P),A+=3)}c.addGroup(h,A,0),h+=A}function g(x){let E=_,A=new te,w=new B,y=0,T=x===!0?t:n,L=x===!0?1:-1;for(let N=1;N<=a;N++)p.push(0,v*L,0),u.push(0,L,0),d.push(.5,.5),_++;let R=_;for(let N=0;N<=a;N++){let z=N/a*l+o,F=Math.cos(z),V=Math.sin(z);w.x=T*V,w.y=v*L,w.z=T*F,p.push(w.x,w.y,w.z),u.push(0,L,0),A.x=F*.5+.5,A.y=V*.5*L+.5,d.push(A.x,A.y),_++}for(let N=0;N<a;N++){let P=E+N,z=R+N;x===!0?f.push(z,z+1,P):f.push(z+1,z,P),y+=3}c.addGroup(h,y,x===!0?1:2),h+=y}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new e(t.radiusTop,t.radiusBottom,t.height,t.radialSegments,t.heightSegments,t.openEnded,t.thetaStart,t.thetaLength)}},ul=class e extends cl{constructor(t=1,n=1,i=32,a=1,s=!1,r=0,o=Math.PI*2){super(0,t,n,i,a,s,r,o),this.type="ConeGeometry",this.parameters={radius:t,height:n,radialSegments:i,heightSegments:a,openEnded:s,thetaStart:r,thetaLength:o}}static fromJSON(t){return new e(t.radius,t.height,t.radialSegments,t.heightSegments,t.openEnded,t.thetaStart,t.thetaLength)}},nd=class e extends An{constructor(t=[],n=[],i=1,a=0){super(),this.type="PolyhedronGeometry",this.parameters={vertices:t,indices:n,radius:i,detail:a};let s=[],r=[];o(a),c(i),f(),this.setAttribute("position",new Oe(s,3)),this.setAttribute("normal",new Oe(s.slice(),3)),this.setAttribute("uv",new Oe(r,2)),a===0?this.computeVertexNormals():this.normalizeNormals();function o(m){let g=new B,x=new B,E=new B;for(let A=0;A<n.length;A+=3)d(n[A+0],g),d(n[A+1],x),d(n[A+2],E),l(g,x,E,m)}function l(m,g,x,E){let A=E+1,w=[];for(let y=0;y<=A;y++){w[y]=[];let T=m.clone().lerp(x,y/A),L=g.clone().lerp(x,y/A),R=A-y;for(let N=0;N<=R;N++)N===0&&y===A?w[y][N]=T:w[y][N]=T.clone().lerp(L,N/R)}for(let y=0;y<A;y++)for(let T=0;T<2*(A-y)-1;T++){let L=Math.floor(T/2);T%2===0?(u(w[y][L+1]),u(w[y+1][L]),u(w[y][L])):(u(w[y][L+1]),u(w[y+1][L+1]),u(w[y+1][L]))}}function c(m){let g=new B;for(let x=0;x<s.length;x+=3)g.x=s[x+0],g.y=s[x+1],g.z=s[x+2],g.normalize().multiplyScalar(m),s[x+0]=g.x,s[x+1]=g.y,s[x+2]=g.z}function f(){let m=new B;for(let g=0;g<s.length;g+=3){m.x=s[g+0],m.y=s[g+1],m.z=s[g+2];let x=v(m)/2/Math.PI+.5,E=h(m)/Math.PI+.5;r.push(x,1-E)}_(),p()}function p(){for(let m=0;m<r.length;m+=6){let g=r[m+0],x=r[m+2],E=r[m+4],A=Math.max(g,x,E),w=Math.min(g,x,E);A>.9&&w<.1&&(g<.2&&(r[m+0]+=1),x<.2&&(r[m+2]+=1),E<.2&&(r[m+4]+=1))}}function u(m){s.push(m.x,m.y,m.z)}function d(m,g){let x=m*3;g.x=t[x+0],g.y=t[x+1],g.z=t[x+2]}function _(){let m=new B,g=new B,x=new B,E=new B,A=new te,w=new te,y=new te;for(let T=0,L=0;T<s.length;T+=9,L+=6){m.set(s[T+0],s[T+1],s[T+2]),g.set(s[T+3],s[T+4],s[T+5]),x.set(s[T+6],s[T+7],s[T+8]),A.set(r[L+0],r[L+1]),w.set(r[L+2],r[L+3]),y.set(r[L+4],r[L+5]),E.copy(m).add(g).add(x).divideScalar(3);let R=v(E);b(A,L+0,m,R),b(w,L+2,g,R),b(y,L+4,x,R)}}function b(m,g,x,E){E<0&&m.x===1&&(r[g]=m.x-1),x.x===0&&x.z===0&&(r[g]=E/2/Math.PI+.5)}function v(m){return Math.atan2(m.z,-m.x)}function h(m){return Math.atan2(-m.y,Math.sqrt(m.x*m.x+m.z*m.z))}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new e(t.vertices,t.indices,t.radius,t.detail)}};var fl=class e extends nd{constructor(t=1,n=0){let i=(1+Math.sqrt(5))/2,a=[-1,i,0,1,i,0,-1,-i,0,1,-i,0,0,-1,i,0,1,i,0,-1,-i,0,1,-i,i,0,-1,i,0,1,-i,0,-1,-i,0,1],s=[0,11,5,0,5,1,0,1,7,0,7,10,0,10,11,1,5,9,5,11,4,11,10,2,10,7,6,7,1,8,3,9,4,3,4,2,3,2,6,3,6,8,3,8,9,4,9,5,2,4,11,6,2,10,8,6,7,9,8,1];super(a,s,t,n),this.type="IcosahedronGeometry",this.parameters={radius:t,detail:n}}static fromJSON(t){return new e(t.radius,t.detail)}};var au=class e extends An{constructor(t=1,n=1,i=1,a=1){super(),this.type="PlaneGeometry",this.parameters={width:t,height:n,widthSegments:i,heightSegments:a};let s=t/2,r=n/2,o=Math.floor(i),l=Math.floor(a),c=o+1,f=l+1,p=t/o,u=n/l,d=[],_=[],b=[],v=[];for(let h=0;h<f;h++){let m=h*u-r;for(let g=0;g<c;g++){let x=g*p-s;_.push(x,-m,0),b.push(0,0,1),v.push(g/o),v.push(1-h/l)}}for(let h=0;h<l;h++)for(let m=0;m<o;m++){let g=m+c*h,x=m+c*(h+1),E=m+1+c*(h+1),A=m+1+c*h;d.push(g,x,A),d.push(x,E,A)}this.setIndex(d),this.setAttribute("position",new Oe(_,3)),this.setAttribute("normal",new Oe(b,3)),this.setAttribute("uv",new Oe(v,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new e(t.width,t.height,t.widthSegments,t.heightSegments)}};var su=class e extends An{constructor(t=1,n=32,i=16,a=0,s=Math.PI*2,r=0,o=Math.PI){super(),this.type="SphereGeometry",this.parameters={radius:t,widthSegments:n,heightSegments:i,phiStart:a,phiLength:s,thetaStart:r,thetaLength:o},n=Math.max(3,Math.floor(n)),i=Math.max(2,Math.floor(i));let l=Math.min(r+o,Math.PI),c=0,f=[],p=new B,u=new B,d=[],_=[],b=[],v=[];for(let h=0;h<=i;h++){let m=[],g=h/i,x=0;h===0&&r===0?x=.5/n:h===i&&l===Math.PI&&(x=-.5/n);for(let E=0;E<=n;E++){let A=E/n;p.x=-t*Math.cos(a+A*s)*Math.sin(r+g*o),p.y=t*Math.cos(r+g*o),p.z=t*Math.sin(a+A*s)*Math.sin(r+g*o),_.push(p.x,p.y,p.z),u.copy(p).normalize(),b.push(u.x,u.y,u.z),v.push(A+x,1-g),m.push(c++)}f.push(m)}for(let h=0;h<i;h++)for(let m=0;m<n;m++){let g=f[h][m+1],x=f[h][m],E=f[h+1][m],A=f[h+1][m+1];(h!==0||r>0)&&d.push(g,x,A),(h!==i-1||l<Math.PI)&&d.push(x,E,A)}this.setIndex(d),this.setAttribute("position",new Oe(_,3)),this.setAttribute("normal",new Oe(b,3)),this.setAttribute("uv",new Oe(v,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new e(t.radius,t.widthSegments,t.heightSegments,t.phiStart,t.phiLength,t.thetaStart,t.thetaLength)}};function Ir(e){let t={};for(let n in e){t[n]={};for(let i in e[n]){let a=e[n][i];a&&(a.isColor||a.isMatrix3||a.isMatrix4||a.isVector2||a.isVector3||a.isVector4||a.isTexture||a.isQuaternion)?a.isRenderTargetTexture?(Ot("UniformsUtils: Textures of render targets cannot be cloned via cloneUniforms() or mergeUniforms()."),t[n][i]=null):t[n][i]=a.clone():Array.isArray(a)?t[n][i]=a.slice():t[n][i]=a}}return t}function zn(e){let t={};for(let n=0;n<e.length;n++){let i=Ir(e[n]);for(let a in i)t[a]=i[a]}return t}function wC(e){let t=[];for(let n=0;n<e.length;n++)t.push(e[n].clone());return t}function bv(e){let t=e.getRenderTarget();return t===null?e.outputColorSpace:t.isXRRenderTarget===!0?t.texture.colorSpace:se.workingColorSpace}var M1={clone:Ir,merge:zn},CC=`void main() {
	gl_Position = projectionMatrix * modelViewMatrix * vec4( position, 1.0 );
}`,RC=`void main() {
	gl_FragColor = vec4( 1.0, 0.0, 0.0, 1.0 );
}`,Yn=class extends Za{constructor(t){super(),this.isShaderMaterial=!0,this.type="ShaderMaterial",this.defines={},this.uniforms={},this.uniformsGroups=[],this.vertexShader=CC,this.fragmentShader=RC,this.linewidth=1,this.wireframe=!1,this.wireframeLinewidth=1,this.fog=!1,this.lights=!1,this.clipping=!1,this.forceSinglePass=!0,this.extensions={clipCullDistance:!1,multiDraw:!1},this.defaultAttributeValues={color:[1,1,1],uv:[0,0],uv1:[0,0]},this.index0AttributeName=void 0,this.uniformsNeedUpdate=!1,this.glslVersion=null,t!==void 0&&this.setValues(t)}copy(t){return super.copy(t),this.fragmentShader=t.fragmentShader,this.vertexShader=t.vertexShader,this.uniforms=Ir(t.uniforms),this.uniformsGroups=wC(t.uniformsGroups),this.defines=Object.assign({},t.defines),this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.fog=t.fog,this.lights=t.lights,this.clipping=t.clipping,this.extensions=Object.assign({},t.extensions),this.glslVersion=t.glslVersion,this.defaultAttributeValues=Object.assign({},t.defaultAttributeValues),this.index0AttributeName=t.index0AttributeName,this.uniformsNeedUpdate=t.uniformsNeedUpdate,this}toJSON(t){let n=super.toJSON(t);n.glslVersion=this.glslVersion,n.uniforms={};for(let a in this.uniforms){let r=this.uniforms[a].value;r&&r.isTexture?n.uniforms[a]={type:"t",value:r.toJSON(t).uuid}:r&&r.isColor?n.uniforms[a]={type:"c",value:r.getHex()}:r&&r.isVector2?n.uniforms[a]={type:"v2",value:r.toArray()}:r&&r.isVector3?n.uniforms[a]={type:"v3",value:r.toArray()}:r&&r.isVector4?n.uniforms[a]={type:"v4",value:r.toArray()}:r&&r.isMatrix3?n.uniforms[a]={type:"m3",value:r.toArray()}:r&&r.isMatrix4?n.uniforms[a]={type:"m4",value:r.toArray()}:n.uniforms[a]={value:r}}Object.keys(this.defines).length>0&&(n.defines=this.defines),n.vertexShader=this.vertexShader,n.fragmentShader=this.fragmentShader,n.lights=this.lights,n.clipping=this.clipping;let i={};for(let a in this.extensions)this.extensions[a]===!0&&(i[a]=!0);return Object.keys(i).length>0&&(n.extensions=i),n}},id=class extends Yn{constructor(t){super(t),this.isRawShaderMaterial=!0,this.type="RawShaderMaterial"}},Nr=class extends Za{constructor(t){super(),this.isMeshStandardMaterial=!0,this.type="MeshStandardMaterial",this.defines={STANDARD:""},this.color=new Ft(16777215),this.roughness=1,this.metalness=0,this.map=null,this.lightMap=null,this.lightMapIntensity=1,this.aoMap=null,this.aoMapIntensity=1,this.emissive=new Ft(0),this.emissiveIntensity=1,this.emissiveMap=null,this.bumpMap=null,this.bumpScale=1,this.normalMap=null,this.normalMapType=_v,this.normalScale=new te(1,1),this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.roughnessMap=null,this.metalnessMap=null,this.alphaMap=null,this.envMap=null,this.envMapRotation=new Xi,this.envMapIntensity=1,this.wireframe=!1,this.wireframeLinewidth=1,this.wireframeLinecap="round",this.wireframeLinejoin="round",this.flatShading=!1,this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.defines={STANDARD:""},this.color.copy(t.color),this.roughness=t.roughness,this.metalness=t.metalness,this.map=t.map,this.lightMap=t.lightMap,this.lightMapIntensity=t.lightMapIntensity,this.aoMap=t.aoMap,this.aoMapIntensity=t.aoMapIntensity,this.emissive.copy(t.emissive),this.emissiveMap=t.emissiveMap,this.emissiveIntensity=t.emissiveIntensity,this.bumpMap=t.bumpMap,this.bumpScale=t.bumpScale,this.normalMap=t.normalMap,this.normalMapType=t.normalMapType,this.normalScale.copy(t.normalScale),this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.roughnessMap=t.roughnessMap,this.metalnessMap=t.metalnessMap,this.alphaMap=t.alphaMap,this.envMap=t.envMap,this.envMapRotation.copy(t.envMapRotation),this.envMapIntensity=t.envMapIntensity,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.wireframeLinecap=t.wireframeLinecap,this.wireframeLinejoin=t.wireframeLinejoin,this.flatShading=t.flatShading,this.fog=t.fog,this}};var ad=class extends Za{constructor(t){super(),this.isMeshDepthMaterial=!0,this.type="MeshDepthMaterial",this.depthPacking=u1,this.map=null,this.alphaMap=null,this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.wireframe=!1,this.wireframeLinewidth=1,this.setValues(t)}copy(t){return super.copy(t),this.depthPacking=t.depthPacking,this.map=t.map,this.alphaMap=t.alphaMap,this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this}},sd=class extends Za{constructor(t){super(),this.isMeshDistanceMaterial=!0,this.type="MeshDistanceMaterial",this.map=null,this.alphaMap=null,this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.setValues(t)}copy(t){return super.copy(t),this.map=t.map,this.alphaMap=t.alphaMap,this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this}};function Lh(e,t){return!e||e.constructor===t?e:typeof t.BYTES_PER_ELEMENT=="number"?new t(e):Array.prototype.slice.call(e)}var Hs=class{constructor(t,n,i,a){this.parameterPositions=t,this._cachedIndex=0,this.resultBuffer=a!==void 0?a:new n.constructor(i),this.sampleValues=n,this.valueSize=i,this.settings=null,this.DefaultSettings_={}}evaluate(t){let n=this.parameterPositions,i=this._cachedIndex,a=n[i],s=n[i-1];t:{e:{let r;n:{i:if(!(t<a)){for(let o=i+2;;){if(a===void 0){if(t<s)break i;return i=n.length,this._cachedIndex=i,this.copySampleValue_(i-1)}if(i===o)break;if(s=a,a=n[++i],t<a)break e}r=n.length;break n}if(!(t>=s)){let o=n[1];t<o&&(i=2,s=o);for(let l=i-2;;){if(s===void 0)return this._cachedIndex=0,this.copySampleValue_(0);if(i===l)break;if(a=s,s=n[--i-1],t>=s)break e}r=i,i=0;break n}break t}for(;i<r;){let o=i+r>>>1;t<n[o]?r=o:i=o+1}if(a=n[i],s=n[i-1],s===void 0)return this._cachedIndex=0,this.copySampleValue_(0);if(a===void 0)return i=n.length,this._cachedIndex=i,this.copySampleValue_(i-1)}this._cachedIndex=i,this.intervalChanged_(i,s,a)}return this.interpolate_(i,s,t,a)}getSettings_(){return this.settings||this.DefaultSettings_}copySampleValue_(t){let n=this.resultBuffer,i=this.sampleValues,a=this.valueSize,s=t*a;for(let r=0;r!==a;++r)n[r]=i[s+r];return n}interpolate_(){throw new Error("call to abstract method")}intervalChanged_(){}},rd=class extends Hs{constructor(t,n,i,a){super(t,n,i,a),this._weightPrev=-0,this._offsetPrev=-0,this._weightNext=-0,this._offsetNext=-0,this.DefaultSettings_={endingStart:q0,endingEnd:q0}}intervalChanged_(t,n,i){let a=this.parameterPositions,s=t-2,r=t+1,o=a[s],l=a[r];if(o===void 0)switch(this.getSettings_().endingStart){case Y0:s=t,o=2*n-i;break;case Z0:s=a.length-2,o=n+a[s]-a[s+1];break;default:s=t,o=i}if(l===void 0)switch(this.getSettings_().endingEnd){case Y0:r=t,l=2*i-n;break;case Z0:r=1,l=i+a[1]-a[0];break;default:r=t-1,l=n}let c=(i-n)*.5,f=this.valueSize;this._weightPrev=c/(n-o),this._weightNext=c/(l-i),this._offsetPrev=s*f,this._offsetNext=r*f}interpolate_(t,n,i,a){let s=this.resultBuffer,r=this.sampleValues,o=this.valueSize,l=t*o,c=l-o,f=this._offsetPrev,p=this._offsetNext,u=this._weightPrev,d=this._weightNext,_=(i-n)/(a-n),b=_*_,v=b*_,h=-u*v+2*u*b-u*_,m=(1+u)*v+(-1.5-2*u)*b+(-.5+u)*_+1,g=(-1-d)*v+(1.5+d)*b+.5*_,x=d*v-d*b;for(let E=0;E!==o;++E)s[E]=h*r[f+E]+m*r[c+E]+g*r[l+E]+x*r[p+E];return s}},od=class extends Hs{constructor(t,n,i,a){super(t,n,i,a)}interpolate_(t,n,i,a){let s=this.resultBuffer,r=this.sampleValues,o=this.valueSize,l=t*o,c=l-o,f=(i-n)/(a-n),p=1-f;for(let u=0;u!==o;++u)s[u]=r[c+u]*p+r[l+u]*f;return s}},ld=class extends Hs{constructor(t,n,i,a){super(t,n,i,a)}interpolate_(t){return this.copySampleValue_(t-1)}},cd=class extends Hs{interpolate_(t,n,i,a){let s=this.resultBuffer,r=this.sampleValues,o=this.valueSize,l=t*o,c=l-o,f=this.settings||this.DefaultSettings_,p=f.inTangents,u=f.outTangents;if(!p||!u){let b=(i-n)/(a-n),v=1-b;for(let h=0;h!==o;++h)s[h]=r[c+h]*v+r[l+h]*b;return s}let d=o*2,_=t-1;for(let b=0;b!==o;++b){let v=r[c+b],h=r[l+b],m=_*d+b*2,g=u[m],x=u[m+1],E=t*d+b*2,A=p[E],w=p[E+1],y=(i-n)/(a-n),T,L,R,N,P;for(let z=0;z<8;z++){T=y*y,L=T*y,R=1-y,N=R*R,P=N*R;let V=P*n+3*N*y*g+3*R*T*A+L*a-i;if(Math.abs(V)<1e-10)break;let X=3*N*(g-n)+6*R*y*(A-g)+3*T*(a-A);if(Math.abs(X)<1e-10)break;y=y-V/X,y=Math.max(0,Math.min(1,y))}s[b]=P*v+3*N*y*x+3*R*T*w+L*h}return s}},vi=class{constructor(t,n,i,a){if(t===void 0)throw new Error("THREE.KeyframeTrack: track name is undefined");if(n===void 0||n.length===0)throw new Error("THREE.KeyframeTrack: no keyframes in track named "+t);this.name=t,this.times=Lh(n,this.TimeBufferType),this.values=Lh(i,this.ValueBufferType),this.setInterpolation(a||this.DefaultInterpolation)}static toJSON(t){let n=t.constructor,i;if(n.toJSON!==this.toJSON)i=n.toJSON(t);else{i={name:t.name,times:Lh(t.times,Array),values:Lh(t.values,Array)};let a=t.getInterpolation();a!==t.DefaultInterpolation&&(i.interpolation=a)}return i.type=t.ValueTypeName,i}InterpolantFactoryMethodDiscrete(t){return new ld(this.times,this.values,this.getValueSize(),t)}InterpolantFactoryMethodLinear(t){return new od(this.times,this.values,this.getValueSize(),t)}InterpolantFactoryMethodSmooth(t){return new rd(this.times,this.values,this.getValueSize(),t)}InterpolantFactoryMethodBezier(t){let n=new cd(this.times,this.values,this.getValueSize(),t);return this.settings&&(n.settings=this.settings),n}setInterpolation(t){let n;switch(t){case Fc:n=this.InterpolantFactoryMethodDiscrete;break;case Zh:n=this.InterpolantFactoryMethodLinear;break;case Ph:n=this.InterpolantFactoryMethodSmooth;break;case W0:n=this.InterpolantFactoryMethodBezier;break}if(n===void 0){let i="unsupported interpolation for "+this.ValueTypeName+" keyframe track named "+this.name;if(this.createInterpolant===void 0)if(t!==this.DefaultInterpolation)this.setInterpolation(this.DefaultInterpolation);else throw new Error(i);return Ot("KeyframeTrack:",i),this}return this.createInterpolant=n,this}getInterpolation(){switch(this.createInterpolant){case this.InterpolantFactoryMethodDiscrete:return Fc;case this.InterpolantFactoryMethodLinear:return Zh;case this.InterpolantFactoryMethodSmooth:return Ph;case this.InterpolantFactoryMethodBezier:return W0}}getValueSize(){return this.values.length/this.times.length}shift(t){if(t!==0){let n=this.times;for(let i=0,a=n.length;i!==a;++i)n[i]+=t}return this}scale(t){if(t!==1){let n=this.times;for(let i=0,a=n.length;i!==a;++i)n[i]*=t}return this}trim(t,n){let i=this.times,a=i.length,s=0,r=a-1;for(;s!==a&&i[s]<t;)++s;for(;r!==-1&&i[r]>n;)--r;if(++r,s!==0||r!==a){s>=r&&(r=Math.max(r,1),s=r-1);let o=this.getValueSize();this.times=i.slice(s,r),this.values=this.values.slice(s*o,r*o)}return this}validate(){let t=!0,n=this.getValueSize();n-Math.floor(n)!==0&&(zt("KeyframeTrack: Invalid value size in track.",this),t=!1);let i=this.times,a=this.values,s=i.length;s===0&&(zt("KeyframeTrack: Track is empty.",this),t=!1);let r=null;for(let o=0;o!==s;o++){let l=i[o];if(typeof l=="number"&&isNaN(l)){zt("KeyframeTrack: Time is not a valid number.",this,o,l),t=!1;break}if(r!==null&&r>l){zt("KeyframeTrack: Out of order keys.",this,o,l,r),t=!1;break}r=l}if(a!==void 0&&Ww(a))for(let o=0,l=a.length;o!==l;++o){let c=a[o];if(isNaN(c)){zt("KeyframeTrack: Value is not a valid number.",this,o,c),t=!1;break}}return t}optimize(){let t=this.times.slice(),n=this.values.slice(),i=this.getValueSize(),a=this.getInterpolation()===Ph,s=t.length-1,r=1;for(let o=1;o<s;++o){let l=!1,c=t[o],f=t[o+1];if(c!==f&&(o!==1||c!==t[0]))if(a)l=!0;else{let p=o*i,u=p-i,d=p+i;for(let _=0;_!==i;++_){let b=n[p+_];if(b!==n[u+_]||b!==n[d+_]){l=!0;break}}}if(l){if(o!==r){t[r]=t[o];let p=o*i,u=r*i;for(let d=0;d!==i;++d)n[u+d]=n[p+d]}++r}}if(s>0){t[r]=t[s];for(let o=s*i,l=r*i,c=0;c!==i;++c)n[l+c]=n[o+c];++r}return r!==t.length?(this.times=t.slice(0,r),this.values=n.slice(0,r*i)):(this.times=t,this.values=n),this}clone(){let t=this.times.slice(),n=this.values.slice(),i=this.constructor,a=new i(this.name,t,n);return a.createInterpolant=this.createInterpolant,a}};vi.prototype.ValueTypeName="";vi.prototype.TimeBufferType=Float32Array;vi.prototype.ValueBufferType=Float32Array;vi.prototype.DefaultInterpolation=Zh;var Gs=class extends vi{constructor(t,n,i){super(t,n,i)}};Gs.prototype.ValueTypeName="bool";Gs.prototype.ValueBufferType=Array;Gs.prototype.DefaultInterpolation=Fc;Gs.prototype.InterpolantFactoryMethodLinear=void 0;Gs.prototype.InterpolantFactoryMethodSmooth=void 0;var ud=class extends vi{constructor(t,n,i,a){super(t,n,i,a)}};ud.prototype.ValueTypeName="color";var fd=class extends vi{constructor(t,n,i,a){super(t,n,i,a)}};fd.prototype.ValueTypeName="number";var hd=class extends Hs{constructor(t,n,i,a){super(t,n,i,a)}interpolate_(t,n,i,a){let s=this.resultBuffer,r=this.sampleValues,o=this.valueSize,l=(i-n)/(a-n),c=t*o;for(let f=c+o;c!==f;c+=4)pa.slerpFlat(s,0,r,c-o,r,c,l);return s}},ru=class extends vi{constructor(t,n,i,a){super(t,n,i,a)}InterpolantFactoryMethodLinear(t){return new hd(this.times,this.values,this.getValueSize(),t)}};ru.prototype.ValueTypeName="quaternion";ru.prototype.InterpolantFactoryMethodSmooth=void 0;var ks=class extends vi{constructor(t,n,i){super(t,n,i)}};ks.prototype.ValueTypeName="string";ks.prototype.ValueBufferType=Array;ks.prototype.DefaultInterpolation=Fc;ks.prototype.InterpolantFactoryMethodLinear=void 0;ks.prototype.InterpolantFactoryMethodSmooth=void 0;var dd=class extends vi{constructor(t,n,i,a){super(t,n,i,a)}};dd.prototype.ValueTypeName="vector";var pd=class{constructor(t,n,i){let a=this,s=!1,r=0,o=0,l,c=[];this.onStart=void 0,this.onLoad=t,this.onProgress=n,this.onError=i,this._abortController=null,this.itemStart=function(f){o++,s===!1&&a.onStart!==void 0&&a.onStart(f,r,o),s=!0},this.itemEnd=function(f){r++,a.onProgress!==void 0&&a.onProgress(f,r,o),r===o&&(s=!1,a.onLoad!==void 0&&a.onLoad())},this.itemError=function(f){a.onError!==void 0&&a.onError(f)},this.resolveURL=function(f){return l?l(f):f},this.setURLModifier=function(f){return l=f,this},this.addHandler=function(f,p){return c.push(f,p),this},this.removeHandler=function(f){let p=c.indexOf(f);return p!==-1&&c.splice(p,2),this},this.getHandler=function(f){for(let p=0,u=c.length;p<u;p+=2){let d=c[p],_=c[p+1];if(d.global&&(d.lastIndex=0),d.test(f))return _}return null},this.abort=function(){return this.abortController.abort(),this._abortController=null,this}}get abortController(){return this._abortController||(this._abortController=new AbortController),this._abortController}},T1=new pd,md=class{constructor(t){this.manager=t!==void 0?t:T1,this.crossOrigin="anonymous",this.withCredentials=!1,this.path="",this.resourcePath="",this.requestHeader={},typeof __THREE_DEVTOOLS__<"u"&&__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("observe",{detail:this}))}load(){}loadAsync(t,n){let i=this;return new Promise(function(a,s){i.load(t,a,n,s)})}parse(){}setCrossOrigin(t){return this.crossOrigin=t,this}setWithCredentials(t){return this.withCredentials=t,this}setPath(t){return this.path=t,this}setResourcePath(t){return this.resourcePath=t,this}setRequestHeader(t){return this.requestHeader=t,this}abort(){return this}};md.DEFAULT_MATERIAL_NAME="__DEFAULT";var ou=class extends tn{constructor(t,n=1){super(),this.isLight=!0,this.type="Light",this.color=new Ft(t),this.intensity=n}dispose(){this.dispatchEvent({type:"dispose"})}copy(t,n){return super.copy(t,n),this.color.copy(t.color),this.intensity=t.intensity,this}toJSON(t){let n=super.toJSON(t);return n.object.color=this.color.getHex(),n.object.intensity=this.intensity,n}},lu=class extends ou{constructor(t,n,i){super(t,i),this.isHemisphereLight=!0,this.type="HemisphereLight",this.position.copy(tn.DEFAULT_UP),this.updateMatrix(),this.groundColor=new Ft(n)}copy(t,n){return super.copy(t,n),this.groundColor.copy(t.groundColor),this}toJSON(t){let n=super.toJSON(t);return n.object.groundColor=this.groundColor.getHex(),n}},X0=new Ie,PM=new B,BM=new B,Q0=class{constructor(t){this.camera=t,this.intensity=1,this.bias=0,this.biasNode=null,this.normalBias=0,this.radius=1,this.blurSamples=8,this.mapSize=new te(512,512),this.mapType=ei,this.map=null,this.mapPass=null,this.matrix=new Ie,this.autoUpdate=!0,this.needsUpdate=!1,this._frustum=new ol,this._frameExtents=new te(1,1),this._viewportCount=1,this._viewports=[new Ge(0,0,1,1)]}getViewportCount(){return this._viewportCount}getFrustum(){return this._frustum}updateMatrices(t){let n=this.camera,i=this.matrix;PM.setFromMatrixPosition(t.matrixWorld),n.position.copy(PM),BM.setFromMatrixPosition(t.target.matrixWorld),n.lookAt(BM),n.updateMatrixWorld(),X0.multiplyMatrices(n.projectionMatrix,n.matrixWorldInverse),this._frustum.setFromProjectionMatrix(X0,n.coordinateSystem,n.reversedDepth),n.coordinateSystem===nl||n.reversedDepth?i.set(.5,0,0,.5,0,.5,0,.5,0,0,1,0,0,0,0,1):i.set(.5,0,0,.5,0,.5,0,.5,0,0,.5,.5,0,0,0,1),i.multiply(X0)}getViewport(t){return this._viewports[t]}getFrameExtents(){return this._frameExtents}dispose(){this.map&&this.map.dispose(),this.mapPass&&this.mapPass.dispose()}copy(t){return this.camera=t.camera.clone(),this.intensity=t.intensity,this.bias=t.bias,this.radius=t.radius,this.autoUpdate=t.autoUpdate,this.needsUpdate=t.needsUpdate,this.normalBias=t.normalBias,this.blurSamples=t.blurSamples,this.mapSize.copy(t.mapSize),this.biasNode=t.biasNode,this}clone(){return new this.constructor().copy(this)}toJSON(){let t={};return this.intensity!==1&&(t.intensity=this.intensity),this.bias!==0&&(t.bias=this.bias),this.normalBias!==0&&(t.normalBias=this.normalBias),this.radius!==1&&(t.radius=this.radius),(this.mapSize.x!==512||this.mapSize.y!==512)&&(t.mapSize=this.mapSize.toArray()),t.camera=this.camera.toJSON(!1).object,delete t.camera.matrix,t}},Ih=new B,Oh=new pa,ua=new B,cu=class extends tn{constructor(){super(),this.isCamera=!0,this.type="Camera",this.matrixWorldInverse=new Ie,this.projectionMatrix=new Ie,this.projectionMatrixInverse=new Ie,this.coordinateSystem=ki,this._reversedDepth=!1}get reversedDepth(){return this._reversedDepth}copy(t,n){return super.copy(t,n),this.matrixWorldInverse.copy(t.matrixWorldInverse),this.projectionMatrix.copy(t.projectionMatrix),this.projectionMatrixInverse.copy(t.projectionMatrixInverse),this.coordinateSystem=t.coordinateSystem,this}getWorldDirection(t){return super.getWorldDirection(t).negate()}updateMatrixWorld(t){super.updateMatrixWorld(t),this.matrixWorld.decompose(Ih,Oh,ua),ua.x===1&&ua.y===1&&ua.z===1?this.matrixWorldInverse.copy(this.matrixWorld).invert():this.matrixWorldInverse.compose(Ih,Oh,ua.set(1,1,1)).invert()}updateWorldMatrix(t,n){super.updateWorldMatrix(t,n),this.matrixWorld.decompose(Ih,Oh,ua),ua.x===1&&ua.y===1&&ua.z===1?this.matrixWorldInverse.copy(this.matrixWorld).invert():this.matrixWorldInverse.compose(Ih,Oh,ua.set(1,1,1)).invert()}clone(){return new this.constructor().copy(this)}},Ps=new B,zM=new te,FM=new te,Pn=class extends cu{constructor(t=50,n=1,i=.1,a=2e3){super(),this.isPerspectiveCamera=!0,this.type="PerspectiveCamera",this.fov=t,this.zoom=1,this.near=i,this.far=a,this.focus=10,this.aspect=n,this.view=null,this.filmGauge=35,this.filmOffset=0,this.updateProjectionMatrix()}copy(t,n){return super.copy(t,n),this.fov=t.fov,this.zoom=t.zoom,this.near=t.near,this.far=t.far,this.focus=t.focus,this.aspect=t.aspect,this.view=t.view===null?null:Object.assign({},t.view),this.filmGauge=t.filmGauge,this.filmOffset=t.filmOffset,this}setFocalLength(t){let n=.5*this.getFilmHeight()/t;this.fov=al*2*Math.atan(n),this.updateProjectionMatrix()}getFocalLength(){let t=Math.tan(Bc*.5*this.fov);return .5*this.getFilmHeight()/t}getEffectiveFOV(){return al*2*Math.atan(Math.tan(Bc*.5*this.fov)/this.zoom)}getFilmWidth(){return this.filmGauge*Math.min(this.aspect,1)}getFilmHeight(){return this.filmGauge/Math.max(this.aspect,1)}getViewBounds(t,n,i){Ps.set(-1,-1,.5).applyMatrix4(this.projectionMatrixInverse),n.set(Ps.x,Ps.y).multiplyScalar(-t/Ps.z),Ps.set(1,1,.5).applyMatrix4(this.projectionMatrixInverse),i.set(Ps.x,Ps.y).multiplyScalar(-t/Ps.z)}getViewSize(t,n){return this.getViewBounds(t,zM,FM),n.subVectors(FM,zM)}setViewOffset(t,n,i,a,s,r){this.aspect=t/n,this.view===null&&(this.view={enabled:!0,fullWidth:1,fullHeight:1,offsetX:0,offsetY:0,width:1,height:1}),this.view.enabled=!0,this.view.fullWidth=t,this.view.fullHeight=n,this.view.offsetX=i,this.view.offsetY=a,this.view.width=s,this.view.height=r,this.updateProjectionMatrix()}clearViewOffset(){this.view!==null&&(this.view.enabled=!1),this.updateProjectionMatrix()}updateProjectionMatrix(){let t=this.near,n=t*Math.tan(Bc*.5*this.fov)/this.zoom,i=2*n,a=this.aspect*i,s=-.5*a,r=this.view;if(this.view!==null&&this.view.enabled){let l=r.fullWidth,c=r.fullHeight;s+=r.offsetX*a/l,n-=r.offsetY*i/c,a*=r.width/l,i*=r.height/c}let o=this.filmOffset;o!==0&&(s+=t*o/this.getFilmWidth()),this.projectionMatrix.makePerspective(s,s+a,n,n-i,t,this.far,this.coordinateSystem,this.reversedDepth),this.projectionMatrixInverse.copy(this.projectionMatrix).invert()}toJSON(t){let n=super.toJSON(t);return n.object.fov=this.fov,n.object.zoom=this.zoom,n.object.near=this.near,n.object.far=this.far,n.object.focus=this.focus,n.object.aspect=this.aspect,this.view!==null&&(n.object.view=Object.assign({},this.view)),n.object.filmGauge=this.filmGauge,n.object.filmOffset=this.filmOffset,n}};var hl=class extends cu{constructor(t=-1,n=1,i=1,a=-1,s=.1,r=2e3){super(),this.isOrthographicCamera=!0,this.type="OrthographicCamera",this.zoom=1,this.view=null,this.left=t,this.right=n,this.top=i,this.bottom=a,this.near=s,this.far=r,this.updateProjectionMatrix()}copy(t,n){return super.copy(t,n),this.left=t.left,this.right=t.right,this.top=t.top,this.bottom=t.bottom,this.near=t.near,this.far=t.far,this.zoom=t.zoom,this.view=t.view===null?null:Object.assign({},t.view),this}setViewOffset(t,n,i,a,s,r){this.view===null&&(this.view={enabled:!0,fullWidth:1,fullHeight:1,offsetX:0,offsetY:0,width:1,height:1}),this.view.enabled=!0,this.view.fullWidth=t,this.view.fullHeight=n,this.view.offsetX=i,this.view.offsetY=a,this.view.width=s,this.view.height=r,this.updateProjectionMatrix()}clearViewOffset(){this.view!==null&&(this.view.enabled=!1),this.updateProjectionMatrix()}updateProjectionMatrix(){let t=(this.right-this.left)/(2*this.zoom),n=(this.top-this.bottom)/(2*this.zoom),i=(this.right+this.left)/2,a=(this.top+this.bottom)/2,s=i-t,r=i+t,o=a+n,l=a-n;if(this.view!==null&&this.view.enabled){let c=(this.right-this.left)/this.view.fullWidth/this.zoom,f=(this.top-this.bottom)/this.view.fullHeight/this.zoom;s+=c*this.view.offsetX,r=s+c*this.view.width,o-=f*this.view.offsetY,l=o-f*this.view.height}this.projectionMatrix.makeOrthographic(s,r,o,l,this.near,this.far,this.coordinateSystem,this.reversedDepth),this.projectionMatrixInverse.copy(this.projectionMatrix).invert()}toJSON(t){let n=super.toJSON(t);return n.object.zoom=this.zoom,n.object.left=this.left,n.object.right=this.right,n.object.top=this.top,n.object.bottom=this.bottom,n.object.near=this.near,n.object.far=this.far,this.view!==null&&(n.object.view=Object.assign({},this.view)),n}},j0=class extends Q0{constructor(){super(new hl(-5,5,5,-5,.5,500)),this.isDirectionalLightShadow=!0}},uu=class extends ou{constructor(t,n){super(t,n),this.isDirectionalLight=!0,this.type="DirectionalLight",this.position.copy(tn.DEFAULT_UP),this.updateMatrix(),this.target=new tn,this.shadow=new j0}dispose(){super.dispose(),this.shadow.dispose()}copy(t){return super.copy(t),this.target=t.target.clone(),this.shadow=t.shadow.clone(),this}toJSON(t){let n=super.toJSON(t);return n.object.shadow=this.shadow.toJSON(),n.object.target=this.target.uuid,n}};var jo=-90,$o=1,gd=class extends tn{constructor(t,n,i){super(),this.type="CubeCamera",this.renderTarget=i,this.coordinateSystem=null,this.activeMipmapLevel=0;let a=new Pn(jo,$o,t,n);a.layers=this.layers,this.add(a);let s=new Pn(jo,$o,t,n);s.layers=this.layers,this.add(s);let r=new Pn(jo,$o,t,n);r.layers=this.layers,this.add(r);let o=new Pn(jo,$o,t,n);o.layers=this.layers,this.add(o);let l=new Pn(jo,$o,t,n);l.layers=this.layers,this.add(l);let c=new Pn(jo,$o,t,n);c.layers=this.layers,this.add(c)}updateCoordinateSystem(){let t=this.coordinateSystem,n=this.children.concat(),[i,a,s,r,o,l]=n;for(let c of n)this.remove(c);if(t===ki)i.up.set(0,1,0),i.lookAt(1,0,0),a.up.set(0,1,0),a.lookAt(-1,0,0),s.up.set(0,0,-1),s.lookAt(0,1,0),r.up.set(0,0,1),r.lookAt(0,-1,0),o.up.set(0,1,0),o.lookAt(0,0,1),l.up.set(0,1,0),l.lookAt(0,0,-1);else if(t===nl)i.up.set(0,-1,0),i.lookAt(-1,0,0),a.up.set(0,-1,0),a.lookAt(1,0,0),s.up.set(0,0,1),s.lookAt(0,1,0),r.up.set(0,0,-1),r.lookAt(0,-1,0),o.up.set(0,-1,0),o.lookAt(0,0,1),l.up.set(0,-1,0),l.lookAt(0,0,-1);else throw new Error("THREE.CubeCamera.updateCoordinateSystem(): Invalid coordinate system: "+t);for(let c of n)this.add(c),c.updateMatrixWorld()}update(t,n){this.parent===null&&this.updateMatrixWorld();let{renderTarget:i,activeMipmapLevel:a}=this;this.coordinateSystem!==t.coordinateSystem&&(this.coordinateSystem=t.coordinateSystem,this.updateCoordinateSystem());let[s,r,o,l,c,f]=this.children,p=t.getRenderTarget(),u=t.getActiveCubeFace(),d=t.getActiveMipmapLevel(),_=t.xr.enabled;t.xr.enabled=!1;let b=i.texture.generateMipmaps;i.texture.generateMipmaps=!1;let v=!1;t.isWebGLRenderer===!0?v=t.state.buffers.depth.getReversed():v=t.reversedDepthBuffer,t.setRenderTarget(i,0,a),v&&t.autoClear===!1&&t.clearDepth(),t.render(n,s),t.setRenderTarget(i,1,a),v&&t.autoClear===!1&&t.clearDepth(),t.render(n,r),t.setRenderTarget(i,2,a),v&&t.autoClear===!1&&t.clearDepth(),t.render(n,o),t.setRenderTarget(i,3,a),v&&t.autoClear===!1&&t.clearDepth(),t.render(n,l),t.setRenderTarget(i,4,a),v&&t.autoClear===!1&&t.clearDepth(),t.render(n,c),i.texture.generateMipmaps=b,t.setRenderTarget(i,5,a),v&&t.autoClear===!1&&t.clearDepth(),t.render(n,f),t.setRenderTarget(p,u,d),t.xr.enabled=_,i.texture.needsPMREMUpdate=!0}},vd=class extends Pn{constructor(t=[]){super(),this.isArrayCamera=!0,this.isMultiViewCamera=!1,this.cameras=t}};var Mv="\\[\\]\\.:\\/",DC=new RegExp("["+Mv+"]","g"),Tv="[^"+Mv+"]",NC="[^"+Mv.replace("\\.","")+"]",UC=/((?:WC+[\/:])*)/.source.replace("WC",Tv),LC=/(WCOD+)?/.source.replace("WCOD",NC),IC=/(?:\.(WC+)(?:\[(.+)\])?)?/.source.replace("WC",Tv),OC=/\.(WC+)(?:\[(.+)\])?/.source.replace("WC",Tv),PC=new RegExp("^"+UC+LC+IC+OC+"$"),BC=["material","materials","bones","map"],$0=class{constructor(t,n,i){let a=i||Be.parseTrackName(n);this._targetGroup=t,this._bindings=t.subscribe_(n,a)}getValue(t,n){this.bind();let i=this._targetGroup.nCachedObjects_,a=this._bindings[i];a!==void 0&&a.getValue(t,n)}setValue(t,n){let i=this._bindings;for(let a=this._targetGroup.nCachedObjects_,s=i.length;a!==s;++a)i[a].setValue(t,n)}bind(){let t=this._bindings;for(let n=this._targetGroup.nCachedObjects_,i=t.length;n!==i;++n)t[n].bind()}unbind(){let t=this._bindings;for(let n=this._targetGroup.nCachedObjects_,i=t.length;n!==i;++n)t[n].unbind()}},Be=class e{constructor(t,n,i){this.path=n,this.parsedPath=i||e.parseTrackName(n),this.node=e.findNode(t,this.parsedPath.nodeName),this.rootNode=t,this.getValue=this._getValue_unbound,this.setValue=this._setValue_unbound}static create(t,n,i){return t&&t.isAnimationObjectGroup?new e.Composite(t,n,i):new e(t,n,i)}static sanitizeNodeName(t){return t.replace(/\s/g,"_").replace(DC,"")}static parseTrackName(t){let n=PC.exec(t);if(n===null)throw new Error("PropertyBinding: Cannot parse trackName: "+t);let i={nodeName:n[2],objectName:n[3],objectIndex:n[4],propertyName:n[5],propertyIndex:n[6]},a=i.nodeName&&i.nodeName.lastIndexOf(".");if(a!==void 0&&a!==-1){let s=i.nodeName.substring(a+1);BC.indexOf(s)!==-1&&(i.nodeName=i.nodeName.substring(0,a),i.objectName=s)}if(i.propertyName===null||i.propertyName.length===0)throw new Error("PropertyBinding: can not parse propertyName from trackName: "+t);return i}static findNode(t,n){if(n===void 0||n===""||n==="."||n===-1||n===t.name||n===t.uuid)return t;if(t.skeleton){let i=t.skeleton.getBoneByName(n);if(i!==void 0)return i}if(t.children){let i=function(s){for(let r=0;r<s.length;r++){let o=s[r];if(o.name===n||o.uuid===n)return o;let l=i(o.children);if(l)return l}return null},a=i(t.children);if(a)return a}return null}_getValue_unavailable(){}_setValue_unavailable(){}_getValue_direct(t,n){t[n]=this.targetObject[this.propertyName]}_getValue_array(t,n){let i=this.resolvedProperty;for(let a=0,s=i.length;a!==s;++a)t[n++]=i[a]}_getValue_arrayElement(t,n){t[n]=this.resolvedProperty[this.propertyIndex]}_getValue_toArray(t,n){this.resolvedProperty.toArray(t,n)}_setValue_direct(t,n){this.targetObject[this.propertyName]=t[n]}_setValue_direct_setNeedsUpdate(t,n){this.targetObject[this.propertyName]=t[n],this.targetObject.needsUpdate=!0}_setValue_direct_setMatrixWorldNeedsUpdate(t,n){this.targetObject[this.propertyName]=t[n],this.targetObject.matrixWorldNeedsUpdate=!0}_setValue_array(t,n){let i=this.resolvedProperty;for(let a=0,s=i.length;a!==s;++a)i[a]=t[n++]}_setValue_array_setNeedsUpdate(t,n){let i=this.resolvedProperty;for(let a=0,s=i.length;a!==s;++a)i[a]=t[n++];this.targetObject.needsUpdate=!0}_setValue_array_setMatrixWorldNeedsUpdate(t,n){let i=this.resolvedProperty;for(let a=0,s=i.length;a!==s;++a)i[a]=t[n++];this.targetObject.matrixWorldNeedsUpdate=!0}_setValue_arrayElement(t,n){this.resolvedProperty[this.propertyIndex]=t[n]}_setValue_arrayElement_setNeedsUpdate(t,n){this.resolvedProperty[this.propertyIndex]=t[n],this.targetObject.needsUpdate=!0}_setValue_arrayElement_setMatrixWorldNeedsUpdate(t,n){this.resolvedProperty[this.propertyIndex]=t[n],this.targetObject.matrixWorldNeedsUpdate=!0}_setValue_fromArray(t,n){this.resolvedProperty.fromArray(t,n)}_setValue_fromArray_setNeedsUpdate(t,n){this.resolvedProperty.fromArray(t,n),this.targetObject.needsUpdate=!0}_setValue_fromArray_setMatrixWorldNeedsUpdate(t,n){this.resolvedProperty.fromArray(t,n),this.targetObject.matrixWorldNeedsUpdate=!0}_getValue_unbound(t,n){this.bind(),this.getValue(t,n)}_setValue_unbound(t,n){this.bind(),this.setValue(t,n)}bind(){let t=this.node,n=this.parsedPath,i=n.objectName,a=n.propertyName,s=n.propertyIndex;if(t||(t=e.findNode(this.rootNode,n.nodeName),this.node=t),this.getValue=this._getValue_unavailable,this.setValue=this._setValue_unavailable,!t){Ot("PropertyBinding: No target node found for track: "+this.path+".");return}if(i){let c=n.objectIndex;switch(i){case"materials":if(!t.material){zt("PropertyBinding: Can not bind to material as node does not have a material.",this);return}if(!t.material.materials){zt("PropertyBinding: Can not bind to material.materials as node.material does not have a materials array.",this);return}t=t.material.materials;break;case"bones":if(!t.skeleton){zt("PropertyBinding: Can not bind to bones as node does not have a skeleton.",this);return}t=t.skeleton.bones;for(let f=0;f<t.length;f++)if(t[f].name===c){c=f;break}break;case"map":if("map"in t){t=t.map;break}if(!t.material){zt("PropertyBinding: Can not bind to material as node does not have a material.",this);return}if(!t.material.map){zt("PropertyBinding: Can not bind to material.map as node.material does not have a map.",this);return}t=t.material.map;break;default:if(t[i]===void 0){zt("PropertyBinding: Can not bind to objectName of node undefined.",this);return}t=t[i]}if(c!==void 0){if(t[c]===void 0){zt("PropertyBinding: Trying to bind to objectIndex of objectName, but is undefined.",this,t);return}t=t[c]}}let r=t[a];if(r===void 0){let c=n.nodeName;zt("PropertyBinding: Trying to update property for track: "+c+"."+a+" but it wasn't found.",t);return}let o=this.Versioning.None;this.targetObject=t,t.isMaterial===!0?o=this.Versioning.NeedsUpdate:t.isObject3D===!0&&(o=this.Versioning.MatrixWorldNeedsUpdate);let l=this.BindingType.Direct;if(s!==void 0){if(a==="morphTargetInfluences"){if(!t.geometry){zt("PropertyBinding: Can not bind to morphTargetInfluences because node does not have a geometry.",this);return}if(!t.geometry.morphAttributes){zt("PropertyBinding: Can not bind to morphTargetInfluences because node does not have a geometry.morphAttributes.",this);return}t.morphTargetDictionary[s]!==void 0&&(s=t.morphTargetDictionary[s])}l=this.BindingType.ArrayElement,this.resolvedProperty=r,this.propertyIndex=s}else r.fromArray!==void 0&&r.toArray!==void 0?(l=this.BindingType.HasFromToArray,this.resolvedProperty=r):Array.isArray(r)?(l=this.BindingType.EntireArray,this.resolvedProperty=r):this.propertyName=a;this.getValue=this.GetterByBindingType[l],this.setValue=this.SetterByBindingTypeAndVersioning[l][o]}unbind(){this.node=null,this.getValue=this._getValue_unbound,this.setValue=this._setValue_unbound}};Be.Composite=$0;Be.prototype.BindingType={Direct:0,EntireArray:1,ArrayElement:2,HasFromToArray:3};Be.prototype.Versioning={None:0,NeedsUpdate:1,MatrixWorldNeedsUpdate:2};Be.prototype.GetterByBindingType=[Be.prototype._getValue_direct,Be.prototype._getValue_array,Be.prototype._getValue_arrayElement,Be.prototype._getValue_toArray];Be.prototype.SetterByBindingTypeAndVersioning=[[Be.prototype._setValue_direct,Be.prototype._setValue_direct_setNeedsUpdate,Be.prototype._setValue_direct_setMatrixWorldNeedsUpdate],[Be.prototype._setValue_array,Be.prototype._setValue_array_setNeedsUpdate,Be.prototype._setValue_array_setMatrixWorldNeedsUpdate],[Be.prototype._setValue_arrayElement,Be.prototype._setValue_arrayElement_setNeedsUpdate,Be.prototype._setValue_arrayElement_setMatrixWorldNeedsUpdate],[Be.prototype._setValue_fromArray,Be.prototype._setValue_fromArray_setNeedsUpdate,Be.prototype._setValue_fromArray_setMatrixWorldNeedsUpdate]];var kN=new Float32Array(1);function Ev(e,t,n,i){let a=zC(i);switch(n){case gv:return e*t;case Td:return e*t/a.components*a.byteLength;case Ed:return e*t/a.components*a.byteLength;case Lr:return e*t*2/a.components*a.byteLength;case Ad:return e*t*2/a.components*a.byteLength;case vv:return e*t*3/a.components*a.byteLength;case Ui:return e*t*4/a.components*a.byteLength;case wd:return e*t*4/a.components*a.byteLength;case mu:case gu:return Math.floor((e+3)/4)*Math.floor((t+3)/4)*8;case vu:case _u:return Math.floor((e+3)/4)*Math.floor((t+3)/4)*16;case Rd:case Nd:return Math.max(e,16)*Math.max(t,8)/4;case Cd:case Dd:return Math.max(e,8)*Math.max(t,8)/2;case Ud:case Ld:case Od:case Pd:return Math.floor((e+3)/4)*Math.floor((t+3)/4)*8;case Id:case Bd:case zd:return Math.floor((e+3)/4)*Math.floor((t+3)/4)*16;case Fd:return Math.floor((e+3)/4)*Math.floor((t+3)/4)*16;case Vd:return Math.floor((e+4)/5)*Math.floor((t+3)/4)*16;case Hd:return Math.floor((e+4)/5)*Math.floor((t+4)/5)*16;case Gd:return Math.floor((e+5)/6)*Math.floor((t+4)/5)*16;case kd:return Math.floor((e+5)/6)*Math.floor((t+5)/6)*16;case Xd:return Math.floor((e+7)/8)*Math.floor((t+4)/5)*16;case Wd:return Math.floor((e+7)/8)*Math.floor((t+5)/6)*16;case qd:return Math.floor((e+7)/8)*Math.floor((t+7)/8)*16;case Yd:return Math.floor((e+9)/10)*Math.floor((t+4)/5)*16;case Zd:return Math.floor((e+9)/10)*Math.floor((t+5)/6)*16;case Jd:return Math.floor((e+9)/10)*Math.floor((t+7)/8)*16;case Kd:return Math.floor((e+9)/10)*Math.floor((t+9)/10)*16;case Qd:return Math.floor((e+11)/12)*Math.floor((t+9)/10)*16;case jd:return Math.floor((e+11)/12)*Math.floor((t+11)/12)*16;case $d:case tp:case ep:return Math.ceil(e/4)*Math.ceil(t/4)*16;case np:case ip:return Math.ceil(e/4)*Math.ceil(t/4)*8;case ap:case sp:return Math.ceil(e/4)*Math.ceil(t/4)*16}throw new Error(`Unable to determine texture byte length for ${n} format.`)}function zC(e){switch(e){case ei:case hv:return{byteLength:1,components:1};case pl:case dv:case _a:return{byteLength:2,components:1};case bd:case Md:return{byteLength:2,components:4};case qi:case Sd:case Ni:return{byteLength:4,components:1};case pv:case mv:return{byteLength:4,components:3}}throw new Error(`Unknown texture type ${e}.`)}typeof __THREE_DEVTOOLS__<"u"&&__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("register",{detail:{revision:"183"}}));typeof window<"u"&&(window.__THREE__?Ot("WARNING: Multiple instances of Three.js being imported."):window.__THREE__="183");function Y1(){let e=null,t=!1,n=null,i=null;function a(s,r){n(s,r),i=e.requestAnimationFrame(a)}return{start:function(){t!==!0&&n!==null&&(i=e.requestAnimationFrame(a),t=!0)},stop:function(){e.cancelAnimationFrame(i),t=!1},setAnimationLoop:function(s){n=s},setContext:function(s){e=s}}}function VC(e){let t=new WeakMap;function n(o,l){let c=o.array,f=o.usage,p=c.byteLength,u=e.createBuffer();e.bindBuffer(l,u),e.bufferData(l,c,f),o.onUploadCallback();let d;if(c instanceof Float32Array)d=e.FLOAT;else if(typeof Float16Array<"u"&&c instanceof Float16Array)d=e.HALF_FLOAT;else if(c instanceof Uint16Array)o.isFloat16BufferAttribute?d=e.HALF_FLOAT:d=e.UNSIGNED_SHORT;else if(c instanceof Int16Array)d=e.SHORT;else if(c instanceof Uint32Array)d=e.UNSIGNED_INT;else if(c instanceof Int32Array)d=e.INT;else if(c instanceof Int8Array)d=e.BYTE;else if(c instanceof Uint8Array)d=e.UNSIGNED_BYTE;else if(c instanceof Uint8ClampedArray)d=e.UNSIGNED_BYTE;else throw new Error("THREE.WebGLAttributes: Unsupported buffer data format: "+c);return{buffer:u,type:d,bytesPerElement:c.BYTES_PER_ELEMENT,version:o.version,size:p}}function i(o,l,c){let f=l.array,p=l.updateRanges;if(e.bindBuffer(c,o),p.length===0)e.bufferSubData(c,0,f);else{p.sort((d,_)=>d.start-_.start);let u=0;for(let d=1;d<p.length;d++){let _=p[u],b=p[d];b.start<=_.start+_.count+1?_.count=Math.max(_.count,b.start+b.count-_.start):(++u,p[u]=b)}p.length=u+1;for(let d=0,_=p.length;d<_;d++){let b=p[d];e.bufferSubData(c,b.start*f.BYTES_PER_ELEMENT,f,b.start,b.count)}l.clearUpdateRanges()}l.onUploadCallback()}function a(o){return o.isInterleavedBufferAttribute&&(o=o.data),t.get(o)}function s(o){o.isInterleavedBufferAttribute&&(o=o.data);let l=t.get(o);l&&(e.deleteBuffer(l.buffer),t.delete(o))}function r(o,l){if(o.isInterleavedBufferAttribute&&(o=o.data),o.isGLBufferAttribute){let f=t.get(o);(!f||f.version<o.version)&&t.set(o,{buffer:o.buffer,type:o.type,bytesPerElement:o.elementSize,version:o.version});return}let c=t.get(o);if(c===void 0)t.set(o,n(o,l));else if(c.version<o.version){if(c.size!==o.array.byteLength)throw new Error("THREE.WebGLAttributes: The size of the buffer attribute's array buffer does not match the original size. Resizing buffer attributes is not supported.");i(c.buffer,o,l),c.version=o.version}}return{get:a,remove:s,update:r}}var HC=`#ifdef USE_ALPHAHASH
	if ( diffuseColor.a < getAlphaHashThreshold( vPosition ) ) discard;
#endif`,GC=`#ifdef USE_ALPHAHASH
	const float ALPHA_HASH_SCALE = 0.05;
	float hash2D( vec2 value ) {
		return fract( 1.0e4 * sin( 17.0 * value.x + 0.1 * value.y ) * ( 0.1 + abs( sin( 13.0 * value.y + value.x ) ) ) );
	}
	float hash3D( vec3 value ) {
		return hash2D( vec2( hash2D( value.xy ), value.z ) );
	}
	float getAlphaHashThreshold( vec3 position ) {
		float maxDeriv = max(
			length( dFdx( position.xyz ) ),
			length( dFdy( position.xyz ) )
		);
		float pixScale = 1.0 / ( ALPHA_HASH_SCALE * maxDeriv );
		vec2 pixScales = vec2(
			exp2( floor( log2( pixScale ) ) ),
			exp2( ceil( log2( pixScale ) ) )
		);
		vec2 alpha = vec2(
			hash3D( floor( pixScales.x * position.xyz ) ),
			hash3D( floor( pixScales.y * position.xyz ) )
		);
		float lerpFactor = fract( log2( pixScale ) );
		float x = ( 1.0 - lerpFactor ) * alpha.x + lerpFactor * alpha.y;
		float a = min( lerpFactor, 1.0 - lerpFactor );
		vec3 cases = vec3(
			x * x / ( 2.0 * a * ( 1.0 - a ) ),
			( x - 0.5 * a ) / ( 1.0 - a ),
			1.0 - ( ( 1.0 - x ) * ( 1.0 - x ) / ( 2.0 * a * ( 1.0 - a ) ) )
		);
		float threshold = ( x < ( 1.0 - a ) )
			? ( ( x < a ) ? cases.x : cases.y )
			: cases.z;
		return clamp( threshold , 1.0e-6, 1.0 );
	}
#endif`,kC=`#ifdef USE_ALPHAMAP
	diffuseColor.a *= texture2D( alphaMap, vAlphaMapUv ).g;
#endif`,XC=`#ifdef USE_ALPHAMAP
	uniform sampler2D alphaMap;
#endif`,WC=`#ifdef USE_ALPHATEST
	#ifdef ALPHA_TO_COVERAGE
	diffuseColor.a = smoothstep( alphaTest, alphaTest + fwidth( diffuseColor.a ), diffuseColor.a );
	if ( diffuseColor.a == 0.0 ) discard;
	#else
	if ( diffuseColor.a < alphaTest ) discard;
	#endif
#endif`,qC=`#ifdef USE_ALPHATEST
	uniform float alphaTest;
#endif`,YC=`#ifdef USE_AOMAP
	float ambientOcclusion = ( texture2D( aoMap, vAoMapUv ).r - 1.0 ) * aoMapIntensity + 1.0;
	reflectedLight.indirectDiffuse *= ambientOcclusion;
	#if defined( USE_CLEARCOAT ) 
		clearcoatSpecularIndirect *= ambientOcclusion;
	#endif
	#if defined( USE_SHEEN ) 
		sheenSpecularIndirect *= ambientOcclusion;
	#endif
	#if defined( USE_ENVMAP ) && defined( STANDARD )
		float dotNV = saturate( dot( geometryNormal, geometryViewDir ) );
		reflectedLight.indirectSpecular *= computeSpecularOcclusion( dotNV, ambientOcclusion, material.roughness );
	#endif
#endif`,ZC=`#ifdef USE_AOMAP
	uniform sampler2D aoMap;
	uniform float aoMapIntensity;
#endif`,JC=`#ifdef USE_BATCHING
	#if ! defined( GL_ANGLE_multi_draw )
	#define gl_DrawID _gl_DrawID
	uniform int _gl_DrawID;
	#endif
	uniform highp sampler2D batchingTexture;
	uniform highp usampler2D batchingIdTexture;
	mat4 getBatchingMatrix( const in float i ) {
		int size = textureSize( batchingTexture, 0 ).x;
		int j = int( i ) * 4;
		int x = j % size;
		int y = j / size;
		vec4 v1 = texelFetch( batchingTexture, ivec2( x, y ), 0 );
		vec4 v2 = texelFetch( batchingTexture, ivec2( x + 1, y ), 0 );
		vec4 v3 = texelFetch( batchingTexture, ivec2( x + 2, y ), 0 );
		vec4 v4 = texelFetch( batchingTexture, ivec2( x + 3, y ), 0 );
		return mat4( v1, v2, v3, v4 );
	}
	float getIndirectIndex( const in int i ) {
		int size = textureSize( batchingIdTexture, 0 ).x;
		int x = i % size;
		int y = i / size;
		return float( texelFetch( batchingIdTexture, ivec2( x, y ), 0 ).r );
	}
#endif
#ifdef USE_BATCHING_COLOR
	uniform sampler2D batchingColorTexture;
	vec4 getBatchingColor( const in float i ) {
		int size = textureSize( batchingColorTexture, 0 ).x;
		int j = int( i );
		int x = j % size;
		int y = j / size;
		return texelFetch( batchingColorTexture, ivec2( x, y ), 0 );
	}
#endif`,KC=`#ifdef USE_BATCHING
	mat4 batchingMatrix = getBatchingMatrix( getIndirectIndex( gl_DrawID ) );
#endif`,QC=`vec3 transformed = vec3( position );
#ifdef USE_ALPHAHASH
	vPosition = vec3( position );
#endif`,jC=`vec3 objectNormal = vec3( normal );
#ifdef USE_TANGENT
	vec3 objectTangent = vec3( tangent.xyz );
#endif`,$C=`float G_BlinnPhong_Implicit( ) {
	return 0.25;
}
float D_BlinnPhong( const in float shininess, const in float dotNH ) {
	return RECIPROCAL_PI * ( shininess * 0.5 + 1.0 ) * pow( dotNH, shininess );
}
vec3 BRDF_BlinnPhong( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, const in vec3 specularColor, const in float shininess ) {
	vec3 halfDir = normalize( lightDir + viewDir );
	float dotNH = saturate( dot( normal, halfDir ) );
	float dotVH = saturate( dot( viewDir, halfDir ) );
	vec3 F = F_Schlick( specularColor, 1.0, dotVH );
	float G = G_BlinnPhong_Implicit( );
	float D = D_BlinnPhong( shininess, dotNH );
	return F * ( G * D );
} // validated`,t3=`#ifdef USE_IRIDESCENCE
	const mat3 XYZ_TO_REC709 = mat3(
		 3.2404542, -0.9692660,  0.0556434,
		-1.5371385,  1.8760108, -0.2040259,
		-0.4985314,  0.0415560,  1.0572252
	);
	vec3 Fresnel0ToIor( vec3 fresnel0 ) {
		vec3 sqrtF0 = sqrt( fresnel0 );
		return ( vec3( 1.0 ) + sqrtF0 ) / ( vec3( 1.0 ) - sqrtF0 );
	}
	vec3 IorToFresnel0( vec3 transmittedIor, float incidentIor ) {
		return pow2( ( transmittedIor - vec3( incidentIor ) ) / ( transmittedIor + vec3( incidentIor ) ) );
	}
	float IorToFresnel0( float transmittedIor, float incidentIor ) {
		return pow2( ( transmittedIor - incidentIor ) / ( transmittedIor + incidentIor ));
	}
	vec3 evalSensitivity( float OPD, vec3 shift ) {
		float phase = 2.0 * PI * OPD * 1.0e-9;
		vec3 val = vec3( 5.4856e-13, 4.4201e-13, 5.2481e-13 );
		vec3 pos = vec3( 1.6810e+06, 1.7953e+06, 2.2084e+06 );
		vec3 var = vec3( 4.3278e+09, 9.3046e+09, 6.6121e+09 );
		vec3 xyz = val * sqrt( 2.0 * PI * var ) * cos( pos * phase + shift ) * exp( - pow2( phase ) * var );
		xyz.x += 9.7470e-14 * sqrt( 2.0 * PI * 4.5282e+09 ) * cos( 2.2399e+06 * phase + shift[ 0 ] ) * exp( - 4.5282e+09 * pow2( phase ) );
		xyz /= 1.0685e-7;
		vec3 rgb = XYZ_TO_REC709 * xyz;
		return rgb;
	}
	vec3 evalIridescence( float outsideIOR, float eta2, float cosTheta1, float thinFilmThickness, vec3 baseF0 ) {
		vec3 I;
		float iridescenceIOR = mix( outsideIOR, eta2, smoothstep( 0.0, 0.03, thinFilmThickness ) );
		float sinTheta2Sq = pow2( outsideIOR / iridescenceIOR ) * ( 1.0 - pow2( cosTheta1 ) );
		float cosTheta2Sq = 1.0 - sinTheta2Sq;
		if ( cosTheta2Sq < 0.0 ) {
			return vec3( 1.0 );
		}
		float cosTheta2 = sqrt( cosTheta2Sq );
		float R0 = IorToFresnel0( iridescenceIOR, outsideIOR );
		float R12 = F_Schlick( R0, 1.0, cosTheta1 );
		float T121 = 1.0 - R12;
		float phi12 = 0.0;
		if ( iridescenceIOR < outsideIOR ) phi12 = PI;
		float phi21 = PI - phi12;
		vec3 baseIOR = Fresnel0ToIor( clamp( baseF0, 0.0, 0.9999 ) );		vec3 R1 = IorToFresnel0( baseIOR, iridescenceIOR );
		vec3 R23 = F_Schlick( R1, 1.0, cosTheta2 );
		vec3 phi23 = vec3( 0.0 );
		if ( baseIOR[ 0 ] < iridescenceIOR ) phi23[ 0 ] = PI;
		if ( baseIOR[ 1 ] < iridescenceIOR ) phi23[ 1 ] = PI;
		if ( baseIOR[ 2 ] < iridescenceIOR ) phi23[ 2 ] = PI;
		float OPD = 2.0 * iridescenceIOR * thinFilmThickness * cosTheta2;
		vec3 phi = vec3( phi21 ) + phi23;
		vec3 R123 = clamp( R12 * R23, 1e-5, 0.9999 );
		vec3 r123 = sqrt( R123 );
		vec3 Rs = pow2( T121 ) * R23 / ( vec3( 1.0 ) - R123 );
		vec3 C0 = R12 + Rs;
		I = C0;
		vec3 Cm = Rs - T121;
		for ( int m = 1; m <= 2; ++ m ) {
			Cm *= r123;
			vec3 Sm = 2.0 * evalSensitivity( float( m ) * OPD, float( m ) * phi );
			I += Cm * Sm;
		}
		return max( I, vec3( 0.0 ) );
	}
#endif`,e3=`#ifdef USE_BUMPMAP
	uniform sampler2D bumpMap;
	uniform float bumpScale;
	vec2 dHdxy_fwd() {
		vec2 dSTdx = dFdx( vBumpMapUv );
		vec2 dSTdy = dFdy( vBumpMapUv );
		float Hll = bumpScale * texture2D( bumpMap, vBumpMapUv ).x;
		float dBx = bumpScale * texture2D( bumpMap, vBumpMapUv + dSTdx ).x - Hll;
		float dBy = bumpScale * texture2D( bumpMap, vBumpMapUv + dSTdy ).x - Hll;
		return vec2( dBx, dBy );
	}
	vec3 perturbNormalArb( vec3 surf_pos, vec3 surf_norm, vec2 dHdxy, float faceDirection ) {
		vec3 vSigmaX = normalize( dFdx( surf_pos.xyz ) );
		vec3 vSigmaY = normalize( dFdy( surf_pos.xyz ) );
		vec3 vN = surf_norm;
		vec3 R1 = cross( vSigmaY, vN );
		vec3 R2 = cross( vN, vSigmaX );
		float fDet = dot( vSigmaX, R1 ) * faceDirection;
		vec3 vGrad = sign( fDet ) * ( dHdxy.x * R1 + dHdxy.y * R2 );
		return normalize( abs( fDet ) * surf_norm - vGrad );
	}
#endif`,n3=`#if NUM_CLIPPING_PLANES > 0
	vec4 plane;
	#ifdef ALPHA_TO_COVERAGE
		float distanceToPlane, distanceGradient;
		float clipOpacity = 1.0;
		#pragma unroll_loop_start
		for ( int i = 0; i < UNION_CLIPPING_PLANES; i ++ ) {
			plane = clippingPlanes[ i ];
			distanceToPlane = - dot( vClipPosition, plane.xyz ) + plane.w;
			distanceGradient = fwidth( distanceToPlane ) / 2.0;
			clipOpacity *= smoothstep( - distanceGradient, distanceGradient, distanceToPlane );
			if ( clipOpacity == 0.0 ) discard;
		}
		#pragma unroll_loop_end
		#if UNION_CLIPPING_PLANES < NUM_CLIPPING_PLANES
			float unionClipOpacity = 1.0;
			#pragma unroll_loop_start
			for ( int i = UNION_CLIPPING_PLANES; i < NUM_CLIPPING_PLANES; i ++ ) {
				plane = clippingPlanes[ i ];
				distanceToPlane = - dot( vClipPosition, plane.xyz ) + plane.w;
				distanceGradient = fwidth( distanceToPlane ) / 2.0;
				unionClipOpacity *= 1.0 - smoothstep( - distanceGradient, distanceGradient, distanceToPlane );
			}
			#pragma unroll_loop_end
			clipOpacity *= 1.0 - unionClipOpacity;
		#endif
		diffuseColor.a *= clipOpacity;
		if ( diffuseColor.a == 0.0 ) discard;
	#else
		#pragma unroll_loop_start
		for ( int i = 0; i < UNION_CLIPPING_PLANES; i ++ ) {
			plane = clippingPlanes[ i ];
			if ( dot( vClipPosition, plane.xyz ) > plane.w ) discard;
		}
		#pragma unroll_loop_end
		#if UNION_CLIPPING_PLANES < NUM_CLIPPING_PLANES
			bool clipped = true;
			#pragma unroll_loop_start
			for ( int i = UNION_CLIPPING_PLANES; i < NUM_CLIPPING_PLANES; i ++ ) {
				plane = clippingPlanes[ i ];
				clipped = ( dot( vClipPosition, plane.xyz ) > plane.w ) && clipped;
			}
			#pragma unroll_loop_end
			if ( clipped ) discard;
		#endif
	#endif
#endif`,i3=`#if NUM_CLIPPING_PLANES > 0
	varying vec3 vClipPosition;
	uniform vec4 clippingPlanes[ NUM_CLIPPING_PLANES ];
#endif`,a3=`#if NUM_CLIPPING_PLANES > 0
	varying vec3 vClipPosition;
#endif`,s3=`#if NUM_CLIPPING_PLANES > 0
	vClipPosition = - mvPosition.xyz;
#endif`,r3=`#if defined( USE_COLOR ) || defined( USE_COLOR_ALPHA )
	diffuseColor *= vColor;
#endif`,o3=`#if defined( USE_COLOR ) || defined( USE_COLOR_ALPHA )
	varying vec4 vColor;
#endif`,l3=`#if defined( USE_COLOR ) || defined( USE_COLOR_ALPHA ) || defined( USE_INSTANCING_COLOR ) || defined( USE_BATCHING_COLOR )
	varying vec4 vColor;
#endif`,c3=`#if defined( USE_COLOR ) || defined( USE_COLOR_ALPHA ) || defined( USE_INSTANCING_COLOR ) || defined( USE_BATCHING_COLOR )
	vColor = vec4( 1.0 );
#endif
#ifdef USE_COLOR_ALPHA
	vColor *= color;
#elif defined( USE_COLOR )
	vColor.rgb *= color;
#endif
#ifdef USE_INSTANCING_COLOR
	vColor.rgb *= instanceColor.rgb;
#endif
#ifdef USE_BATCHING_COLOR
	vColor *= getBatchingColor( getIndirectIndex( gl_DrawID ) );
#endif`,u3=`#define PI 3.141592653589793
#define PI2 6.283185307179586
#define PI_HALF 1.5707963267948966
#define RECIPROCAL_PI 0.3183098861837907
#define RECIPROCAL_PI2 0.15915494309189535
#define EPSILON 1e-6
#ifndef saturate
#define saturate( a ) clamp( a, 0.0, 1.0 )
#endif
#define whiteComplement( a ) ( 1.0 - saturate( a ) )
float pow2( const in float x ) { return x*x; }
vec3 pow2( const in vec3 x ) { return x*x; }
float pow3( const in float x ) { return x*x*x; }
float pow4( const in float x ) { float x2 = x*x; return x2*x2; }
float max3( const in vec3 v ) { return max( max( v.x, v.y ), v.z ); }
float average( const in vec3 v ) { return dot( v, vec3( 0.3333333 ) ); }
highp float rand( const in vec2 uv ) {
	const highp float a = 12.9898, b = 78.233, c = 43758.5453;
	highp float dt = dot( uv.xy, vec2( a,b ) ), sn = mod( dt, PI );
	return fract( sin( sn ) * c );
}
#ifdef HIGH_PRECISION
	float precisionSafeLength( vec3 v ) { return length( v ); }
#else
	float precisionSafeLength( vec3 v ) {
		float maxComponent = max3( abs( v ) );
		return length( v / maxComponent ) * maxComponent;
	}
#endif
struct IncidentLight {
	vec3 color;
	vec3 direction;
	bool visible;
};
struct ReflectedLight {
	vec3 directDiffuse;
	vec3 directSpecular;
	vec3 indirectDiffuse;
	vec3 indirectSpecular;
};
#ifdef USE_ALPHAHASH
	varying vec3 vPosition;
#endif
vec3 transformDirection( in vec3 dir, in mat4 matrix ) {
	return normalize( ( matrix * vec4( dir, 0.0 ) ).xyz );
}
vec3 inverseTransformDirection( in vec3 dir, in mat4 matrix ) {
	return normalize( ( vec4( dir, 0.0 ) * matrix ).xyz );
}
bool isPerspectiveMatrix( mat4 m ) {
	return m[ 2 ][ 3 ] == - 1.0;
}
vec2 equirectUv( in vec3 dir ) {
	float u = atan( dir.z, dir.x ) * RECIPROCAL_PI2 + 0.5;
	float v = asin( clamp( dir.y, - 1.0, 1.0 ) ) * RECIPROCAL_PI + 0.5;
	return vec2( u, v );
}
vec3 BRDF_Lambert( const in vec3 diffuseColor ) {
	return RECIPROCAL_PI * diffuseColor;
}
vec3 F_Schlick( const in vec3 f0, const in float f90, const in float dotVH ) {
	float fresnel = exp2( ( - 5.55473 * dotVH - 6.98316 ) * dotVH );
	return f0 * ( 1.0 - fresnel ) + ( f90 * fresnel );
}
float F_Schlick( const in float f0, const in float f90, const in float dotVH ) {
	float fresnel = exp2( ( - 5.55473 * dotVH - 6.98316 ) * dotVH );
	return f0 * ( 1.0 - fresnel ) + ( f90 * fresnel );
} // validated`,f3=`#ifdef ENVMAP_TYPE_CUBE_UV
	#define cubeUV_minMipLevel 4.0
	#define cubeUV_minTileSize 16.0
	float getFace( vec3 direction ) {
		vec3 absDirection = abs( direction );
		float face = - 1.0;
		if ( absDirection.x > absDirection.z ) {
			if ( absDirection.x > absDirection.y )
				face = direction.x > 0.0 ? 0.0 : 3.0;
			else
				face = direction.y > 0.0 ? 1.0 : 4.0;
		} else {
			if ( absDirection.z > absDirection.y )
				face = direction.z > 0.0 ? 2.0 : 5.0;
			else
				face = direction.y > 0.0 ? 1.0 : 4.0;
		}
		return face;
	}
	vec2 getUV( vec3 direction, float face ) {
		vec2 uv;
		if ( face == 0.0 ) {
			uv = vec2( direction.z, direction.y ) / abs( direction.x );
		} else if ( face == 1.0 ) {
			uv = vec2( - direction.x, - direction.z ) / abs( direction.y );
		} else if ( face == 2.0 ) {
			uv = vec2( - direction.x, direction.y ) / abs( direction.z );
		} else if ( face == 3.0 ) {
			uv = vec2( - direction.z, direction.y ) / abs( direction.x );
		} else if ( face == 4.0 ) {
			uv = vec2( - direction.x, direction.z ) / abs( direction.y );
		} else {
			uv = vec2( direction.x, direction.y ) / abs( direction.z );
		}
		return 0.5 * ( uv + 1.0 );
	}
	vec3 bilinearCubeUV( sampler2D envMap, vec3 direction, float mipInt ) {
		float face = getFace( direction );
		float filterInt = max( cubeUV_minMipLevel - mipInt, 0.0 );
		mipInt = max( mipInt, cubeUV_minMipLevel );
		float faceSize = exp2( mipInt );
		highp vec2 uv = getUV( direction, face ) * ( faceSize - 2.0 ) + 1.0;
		if ( face > 2.0 ) {
			uv.y += faceSize;
			face -= 3.0;
		}
		uv.x += face * faceSize;
		uv.x += filterInt * 3.0 * cubeUV_minTileSize;
		uv.y += 4.0 * ( exp2( CUBEUV_MAX_MIP ) - faceSize );
		uv.x *= CUBEUV_TEXEL_WIDTH;
		uv.y *= CUBEUV_TEXEL_HEIGHT;
		#ifdef texture2DGradEXT
			return texture2DGradEXT( envMap, uv, vec2( 0.0 ), vec2( 0.0 ) ).rgb;
		#else
			return texture2D( envMap, uv ).rgb;
		#endif
	}
	#define cubeUV_r0 1.0
	#define cubeUV_m0 - 2.0
	#define cubeUV_r1 0.8
	#define cubeUV_m1 - 1.0
	#define cubeUV_r4 0.4
	#define cubeUV_m4 2.0
	#define cubeUV_r5 0.305
	#define cubeUV_m5 3.0
	#define cubeUV_r6 0.21
	#define cubeUV_m6 4.0
	float roughnessToMip( float roughness ) {
		float mip = 0.0;
		if ( roughness >= cubeUV_r1 ) {
			mip = ( cubeUV_r0 - roughness ) * ( cubeUV_m1 - cubeUV_m0 ) / ( cubeUV_r0 - cubeUV_r1 ) + cubeUV_m0;
		} else if ( roughness >= cubeUV_r4 ) {
			mip = ( cubeUV_r1 - roughness ) * ( cubeUV_m4 - cubeUV_m1 ) / ( cubeUV_r1 - cubeUV_r4 ) + cubeUV_m1;
		} else if ( roughness >= cubeUV_r5 ) {
			mip = ( cubeUV_r4 - roughness ) * ( cubeUV_m5 - cubeUV_m4 ) / ( cubeUV_r4 - cubeUV_r5 ) + cubeUV_m4;
		} else if ( roughness >= cubeUV_r6 ) {
			mip = ( cubeUV_r5 - roughness ) * ( cubeUV_m6 - cubeUV_m5 ) / ( cubeUV_r5 - cubeUV_r6 ) + cubeUV_m5;
		} else {
			mip = - 2.0 * log2( 1.16 * roughness );		}
		return mip;
	}
	vec4 textureCubeUV( sampler2D envMap, vec3 sampleDir, float roughness ) {
		float mip = clamp( roughnessToMip( roughness ), cubeUV_m0, CUBEUV_MAX_MIP );
		float mipF = fract( mip );
		float mipInt = floor( mip );
		vec3 color0 = bilinearCubeUV( envMap, sampleDir, mipInt );
		if ( mipF == 0.0 ) {
			return vec4( color0, 1.0 );
		} else {
			vec3 color1 = bilinearCubeUV( envMap, sampleDir, mipInt + 1.0 );
			return vec4( mix( color0, color1, mipF ), 1.0 );
		}
	}
#endif`,h3=`vec3 transformedNormal = objectNormal;
#ifdef USE_TANGENT
	vec3 transformedTangent = objectTangent;
#endif
#ifdef USE_BATCHING
	mat3 bm = mat3( batchingMatrix );
	transformedNormal /= vec3( dot( bm[ 0 ], bm[ 0 ] ), dot( bm[ 1 ], bm[ 1 ] ), dot( bm[ 2 ], bm[ 2 ] ) );
	transformedNormal = bm * transformedNormal;
	#ifdef USE_TANGENT
		transformedTangent = bm * transformedTangent;
	#endif
#endif
#ifdef USE_INSTANCING
	mat3 im = mat3( instanceMatrix );
	transformedNormal /= vec3( dot( im[ 0 ], im[ 0 ] ), dot( im[ 1 ], im[ 1 ] ), dot( im[ 2 ], im[ 2 ] ) );
	transformedNormal = im * transformedNormal;
	#ifdef USE_TANGENT
		transformedTangent = im * transformedTangent;
	#endif
#endif
transformedNormal = normalMatrix * transformedNormal;
#ifdef FLIP_SIDED
	transformedNormal = - transformedNormal;
#endif
#ifdef USE_TANGENT
	transformedTangent = ( modelViewMatrix * vec4( transformedTangent, 0.0 ) ).xyz;
	#ifdef FLIP_SIDED
		transformedTangent = - transformedTangent;
	#endif
#endif`,d3=`#ifdef USE_DISPLACEMENTMAP
	uniform sampler2D displacementMap;
	uniform float displacementScale;
	uniform float displacementBias;
#endif`,p3=`#ifdef USE_DISPLACEMENTMAP
	transformed += normalize( objectNormal ) * ( texture2D( displacementMap, vDisplacementMapUv ).x * displacementScale + displacementBias );
#endif`,m3=`#ifdef USE_EMISSIVEMAP
	vec4 emissiveColor = texture2D( emissiveMap, vEmissiveMapUv );
	#ifdef DECODE_VIDEO_TEXTURE_EMISSIVE
		emissiveColor = sRGBTransferEOTF( emissiveColor );
	#endif
	totalEmissiveRadiance *= emissiveColor.rgb;
#endif`,g3=`#ifdef USE_EMISSIVEMAP
	uniform sampler2D emissiveMap;
#endif`,v3="gl_FragColor = linearToOutputTexel( gl_FragColor );",_3=`vec4 LinearTransferOETF( in vec4 value ) {
	return value;
}
vec4 sRGBTransferEOTF( in vec4 value ) {
	return vec4( mix( pow( value.rgb * 0.9478672986 + vec3( 0.0521327014 ), vec3( 2.4 ) ), value.rgb * 0.0773993808, vec3( lessThanEqual( value.rgb, vec3( 0.04045 ) ) ) ), value.a );
}
vec4 sRGBTransferOETF( in vec4 value ) {
	return vec4( mix( pow( value.rgb, vec3( 0.41666 ) ) * 1.055 - vec3( 0.055 ), value.rgb * 12.92, vec3( lessThanEqual( value.rgb, vec3( 0.0031308 ) ) ) ), value.a );
}`,y3=`#ifdef USE_ENVMAP
	#ifdef ENV_WORLDPOS
		vec3 cameraToFrag;
		if ( isOrthographic ) {
			cameraToFrag = normalize( vec3( - viewMatrix[ 0 ][ 2 ], - viewMatrix[ 1 ][ 2 ], - viewMatrix[ 2 ][ 2 ] ) );
		} else {
			cameraToFrag = normalize( vWorldPosition - cameraPosition );
		}
		vec3 worldNormal = inverseTransformDirection( normal, viewMatrix );
		#ifdef ENVMAP_MODE_REFLECTION
			vec3 reflectVec = reflect( cameraToFrag, worldNormal );
		#else
			vec3 reflectVec = refract( cameraToFrag, worldNormal, refractionRatio );
		#endif
	#else
		vec3 reflectVec = vReflect;
	#endif
	#ifdef ENVMAP_TYPE_CUBE
		vec4 envColor = textureCube( envMap, envMapRotation * vec3( flipEnvMap * reflectVec.x, reflectVec.yz ) );
		#ifdef ENVMAP_BLENDING_MULTIPLY
			outgoingLight = mix( outgoingLight, outgoingLight * envColor.xyz, specularStrength * reflectivity );
		#elif defined( ENVMAP_BLENDING_MIX )
			outgoingLight = mix( outgoingLight, envColor.xyz, specularStrength * reflectivity );
		#elif defined( ENVMAP_BLENDING_ADD )
			outgoingLight += envColor.xyz * specularStrength * reflectivity;
		#endif
	#endif
#endif`,x3=`#ifdef USE_ENVMAP
	uniform float envMapIntensity;
	uniform float flipEnvMap;
	uniform mat3 envMapRotation;
	#ifdef ENVMAP_TYPE_CUBE
		uniform samplerCube envMap;
	#else
		uniform sampler2D envMap;
	#endif
#endif`,S3=`#ifdef USE_ENVMAP
	uniform float reflectivity;
	#if defined( USE_BUMPMAP ) || defined( USE_NORMALMAP ) || defined( PHONG ) || defined( LAMBERT )
		#define ENV_WORLDPOS
	#endif
	#ifdef ENV_WORLDPOS
		varying vec3 vWorldPosition;
		uniform float refractionRatio;
	#else
		varying vec3 vReflect;
	#endif
#endif`,b3=`#ifdef USE_ENVMAP
	#if defined( USE_BUMPMAP ) || defined( USE_NORMALMAP ) || defined( PHONG ) || defined( LAMBERT )
		#define ENV_WORLDPOS
	#endif
	#ifdef ENV_WORLDPOS
		
		varying vec3 vWorldPosition;
	#else
		varying vec3 vReflect;
		uniform float refractionRatio;
	#endif
#endif`,M3=`#ifdef USE_ENVMAP
	#ifdef ENV_WORLDPOS
		vWorldPosition = worldPosition.xyz;
	#else
		vec3 cameraToVertex;
		if ( isOrthographic ) {
			cameraToVertex = normalize( vec3( - viewMatrix[ 0 ][ 2 ], - viewMatrix[ 1 ][ 2 ], - viewMatrix[ 2 ][ 2 ] ) );
		} else {
			cameraToVertex = normalize( worldPosition.xyz - cameraPosition );
		}
		vec3 worldNormal = inverseTransformDirection( transformedNormal, viewMatrix );
		#ifdef ENVMAP_MODE_REFLECTION
			vReflect = reflect( cameraToVertex, worldNormal );
		#else
			vReflect = refract( cameraToVertex, worldNormal, refractionRatio );
		#endif
	#endif
#endif`,T3=`#ifdef USE_FOG
	vFogDepth = - mvPosition.z;
#endif`,E3=`#ifdef USE_FOG
	varying float vFogDepth;
#endif`,A3=`#ifdef USE_FOG
	#ifdef FOG_EXP2
		float fogFactor = 1.0 - exp( - fogDensity * fogDensity * vFogDepth * vFogDepth );
	#else
		float fogFactor = smoothstep( fogNear, fogFar, vFogDepth );
	#endif
	gl_FragColor.rgb = mix( gl_FragColor.rgb, fogColor, fogFactor );
#endif`,w3=`#ifdef USE_FOG
	uniform vec3 fogColor;
	varying float vFogDepth;
	#ifdef FOG_EXP2
		uniform float fogDensity;
	#else
		uniform float fogNear;
		uniform float fogFar;
	#endif
#endif`,C3=`#ifdef USE_GRADIENTMAP
	uniform sampler2D gradientMap;
#endif
vec3 getGradientIrradiance( vec3 normal, vec3 lightDirection ) {
	float dotNL = dot( normal, lightDirection );
	vec2 coord = vec2( dotNL * 0.5 + 0.5, 0.0 );
	#ifdef USE_GRADIENTMAP
		return vec3( texture2D( gradientMap, coord ).r );
	#else
		vec2 fw = fwidth( coord ) * 0.5;
		return mix( vec3( 0.7 ), vec3( 1.0 ), smoothstep( 0.7 - fw.x, 0.7 + fw.x, coord.x ) );
	#endif
}`,R3=`#ifdef USE_LIGHTMAP
	uniform sampler2D lightMap;
	uniform float lightMapIntensity;
#endif`,D3=`LambertMaterial material;
material.diffuseColor = diffuseColor.rgb;
material.specularStrength = specularStrength;`,N3=`varying vec3 vViewPosition;
struct LambertMaterial {
	vec3 diffuseColor;
	float specularStrength;
};
void RE_Direct_Lambert( const in IncidentLight directLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in LambertMaterial material, inout ReflectedLight reflectedLight ) {
	float dotNL = saturate( dot( geometryNormal, directLight.direction ) );
	vec3 irradiance = dotNL * directLight.color;
	reflectedLight.directDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
void RE_IndirectDiffuse_Lambert( const in vec3 irradiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in LambertMaterial material, inout ReflectedLight reflectedLight ) {
	reflectedLight.indirectDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
#define RE_Direct				RE_Direct_Lambert
#define RE_IndirectDiffuse		RE_IndirectDiffuse_Lambert`,U3=`uniform bool receiveShadow;
uniform vec3 ambientLightColor;
#if defined( USE_LIGHT_PROBES )
	uniform vec3 lightProbe[ 9 ];
#endif
vec3 shGetIrradianceAt( in vec3 normal, in vec3 shCoefficients[ 9 ] ) {
	float x = normal.x, y = normal.y, z = normal.z;
	vec3 result = shCoefficients[ 0 ] * 0.886227;
	result += shCoefficients[ 1 ] * 2.0 * 0.511664 * y;
	result += shCoefficients[ 2 ] * 2.0 * 0.511664 * z;
	result += shCoefficients[ 3 ] * 2.0 * 0.511664 * x;
	result += shCoefficients[ 4 ] * 2.0 * 0.429043 * x * y;
	result += shCoefficients[ 5 ] * 2.0 * 0.429043 * y * z;
	result += shCoefficients[ 6 ] * ( 0.743125 * z * z - 0.247708 );
	result += shCoefficients[ 7 ] * 2.0 * 0.429043 * x * z;
	result += shCoefficients[ 8 ] * 0.429043 * ( x * x - y * y );
	return result;
}
vec3 getLightProbeIrradiance( const in vec3 lightProbe[ 9 ], const in vec3 normal ) {
	vec3 worldNormal = inverseTransformDirection( normal, viewMatrix );
	vec3 irradiance = shGetIrradianceAt( worldNormal, lightProbe );
	return irradiance;
}
vec3 getAmbientLightIrradiance( const in vec3 ambientLightColor ) {
	vec3 irradiance = ambientLightColor;
	return irradiance;
}
float getDistanceAttenuation( const in float lightDistance, const in float cutoffDistance, const in float decayExponent ) {
	float distanceFalloff = 1.0 / max( pow( lightDistance, decayExponent ), 0.01 );
	if ( cutoffDistance > 0.0 ) {
		distanceFalloff *= pow2( saturate( 1.0 - pow4( lightDistance / cutoffDistance ) ) );
	}
	return distanceFalloff;
}
float getSpotAttenuation( const in float coneCosine, const in float penumbraCosine, const in float angleCosine ) {
	return smoothstep( coneCosine, penumbraCosine, angleCosine );
}
#if NUM_DIR_LIGHTS > 0
	struct DirectionalLight {
		vec3 direction;
		vec3 color;
	};
	uniform DirectionalLight directionalLights[ NUM_DIR_LIGHTS ];
	void getDirectionalLightInfo( const in DirectionalLight directionalLight, out IncidentLight light ) {
		light.color = directionalLight.color;
		light.direction = directionalLight.direction;
		light.visible = true;
	}
#endif
#if NUM_POINT_LIGHTS > 0
	struct PointLight {
		vec3 position;
		vec3 color;
		float distance;
		float decay;
	};
	uniform PointLight pointLights[ NUM_POINT_LIGHTS ];
	void getPointLightInfo( const in PointLight pointLight, const in vec3 geometryPosition, out IncidentLight light ) {
		vec3 lVector = pointLight.position - geometryPosition;
		light.direction = normalize( lVector );
		float lightDistance = length( lVector );
		light.color = pointLight.color;
		light.color *= getDistanceAttenuation( lightDistance, pointLight.distance, pointLight.decay );
		light.visible = ( light.color != vec3( 0.0 ) );
	}
#endif
#if NUM_SPOT_LIGHTS > 0
	struct SpotLight {
		vec3 position;
		vec3 direction;
		vec3 color;
		float distance;
		float decay;
		float coneCos;
		float penumbraCos;
	};
	uniform SpotLight spotLights[ NUM_SPOT_LIGHTS ];
	void getSpotLightInfo( const in SpotLight spotLight, const in vec3 geometryPosition, out IncidentLight light ) {
		vec3 lVector = spotLight.position - geometryPosition;
		light.direction = normalize( lVector );
		float angleCos = dot( light.direction, spotLight.direction );
		float spotAttenuation = getSpotAttenuation( spotLight.coneCos, spotLight.penumbraCos, angleCos );
		if ( spotAttenuation > 0.0 ) {
			float lightDistance = length( lVector );
			light.color = spotLight.color * spotAttenuation;
			light.color *= getDistanceAttenuation( lightDistance, spotLight.distance, spotLight.decay );
			light.visible = ( light.color != vec3( 0.0 ) );
		} else {
			light.color = vec3( 0.0 );
			light.visible = false;
		}
	}
#endif
#if NUM_RECT_AREA_LIGHTS > 0
	struct RectAreaLight {
		vec3 color;
		vec3 position;
		vec3 halfWidth;
		vec3 halfHeight;
	};
	uniform sampler2D ltc_1;	uniform sampler2D ltc_2;
	uniform RectAreaLight rectAreaLights[ NUM_RECT_AREA_LIGHTS ];
#endif
#if NUM_HEMI_LIGHTS > 0
	struct HemisphereLight {
		vec3 direction;
		vec3 skyColor;
		vec3 groundColor;
	};
	uniform HemisphereLight hemisphereLights[ NUM_HEMI_LIGHTS ];
	vec3 getHemisphereLightIrradiance( const in HemisphereLight hemiLight, const in vec3 normal ) {
		float dotNL = dot( normal, hemiLight.direction );
		float hemiDiffuseWeight = 0.5 * dotNL + 0.5;
		vec3 irradiance = mix( hemiLight.groundColor, hemiLight.skyColor, hemiDiffuseWeight );
		return irradiance;
	}
#endif`,L3=`#ifdef USE_ENVMAP
	vec3 getIBLIrradiance( const in vec3 normal ) {
		#ifdef ENVMAP_TYPE_CUBE_UV
			vec3 worldNormal = inverseTransformDirection( normal, viewMatrix );
			vec4 envMapColor = textureCubeUV( envMap, envMapRotation * worldNormal, 1.0 );
			return PI * envMapColor.rgb * envMapIntensity;
		#else
			return vec3( 0.0 );
		#endif
	}
	vec3 getIBLRadiance( const in vec3 viewDir, const in vec3 normal, const in float roughness ) {
		#ifdef ENVMAP_TYPE_CUBE_UV
			vec3 reflectVec = reflect( - viewDir, normal );
			reflectVec = normalize( mix( reflectVec, normal, pow4( roughness ) ) );
			reflectVec = inverseTransformDirection( reflectVec, viewMatrix );
			vec4 envMapColor = textureCubeUV( envMap, envMapRotation * reflectVec, roughness );
			return envMapColor.rgb * envMapIntensity;
		#else
			return vec3( 0.0 );
		#endif
	}
	#ifdef USE_ANISOTROPY
		vec3 getIBLAnisotropyRadiance( const in vec3 viewDir, const in vec3 normal, const in float roughness, const in vec3 bitangent, const in float anisotropy ) {
			#ifdef ENVMAP_TYPE_CUBE_UV
				vec3 bentNormal = cross( bitangent, viewDir );
				bentNormal = normalize( cross( bentNormal, bitangent ) );
				bentNormal = normalize( mix( bentNormal, normal, pow2( pow2( 1.0 - anisotropy * ( 1.0 - roughness ) ) ) ) );
				return getIBLRadiance( viewDir, bentNormal, roughness );
			#else
				return vec3( 0.0 );
			#endif
		}
	#endif
#endif`,I3=`ToonMaterial material;
material.diffuseColor = diffuseColor.rgb;`,O3=`varying vec3 vViewPosition;
struct ToonMaterial {
	vec3 diffuseColor;
};
void RE_Direct_Toon( const in IncidentLight directLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in ToonMaterial material, inout ReflectedLight reflectedLight ) {
	vec3 irradiance = getGradientIrradiance( geometryNormal, directLight.direction ) * directLight.color;
	reflectedLight.directDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
void RE_IndirectDiffuse_Toon( const in vec3 irradiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in ToonMaterial material, inout ReflectedLight reflectedLight ) {
	reflectedLight.indirectDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
#define RE_Direct				RE_Direct_Toon
#define RE_IndirectDiffuse		RE_IndirectDiffuse_Toon`,P3=`BlinnPhongMaterial material;
material.diffuseColor = diffuseColor.rgb;
material.specularColor = specular;
material.specularShininess = shininess;
material.specularStrength = specularStrength;`,B3=`varying vec3 vViewPosition;
struct BlinnPhongMaterial {
	vec3 diffuseColor;
	vec3 specularColor;
	float specularShininess;
	float specularStrength;
};
void RE_Direct_BlinnPhong( const in IncidentLight directLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in BlinnPhongMaterial material, inout ReflectedLight reflectedLight ) {
	float dotNL = saturate( dot( geometryNormal, directLight.direction ) );
	vec3 irradiance = dotNL * directLight.color;
	reflectedLight.directDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
	reflectedLight.directSpecular += irradiance * BRDF_BlinnPhong( directLight.direction, geometryViewDir, geometryNormal, material.specularColor, material.specularShininess ) * material.specularStrength;
}
void RE_IndirectDiffuse_BlinnPhong( const in vec3 irradiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in BlinnPhongMaterial material, inout ReflectedLight reflectedLight ) {
	reflectedLight.indirectDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
#define RE_Direct				RE_Direct_BlinnPhong
#define RE_IndirectDiffuse		RE_IndirectDiffuse_BlinnPhong`,z3=`PhysicalMaterial material;
material.diffuseColor = diffuseColor.rgb;
material.diffuseContribution = diffuseColor.rgb * ( 1.0 - metalnessFactor );
material.metalness = metalnessFactor;
vec3 dxy = max( abs( dFdx( nonPerturbedNormal ) ), abs( dFdy( nonPerturbedNormal ) ) );
float geometryRoughness = max( max( dxy.x, dxy.y ), dxy.z );
material.roughness = max( roughnessFactor, 0.0525 );material.roughness += geometryRoughness;
material.roughness = min( material.roughness, 1.0 );
#ifdef IOR
	material.ior = ior;
	#ifdef USE_SPECULAR
		float specularIntensityFactor = specularIntensity;
		vec3 specularColorFactor = specularColor;
		#ifdef USE_SPECULAR_COLORMAP
			specularColorFactor *= texture2D( specularColorMap, vSpecularColorMapUv ).rgb;
		#endif
		#ifdef USE_SPECULAR_INTENSITYMAP
			specularIntensityFactor *= texture2D( specularIntensityMap, vSpecularIntensityMapUv ).a;
		#endif
		material.specularF90 = mix( specularIntensityFactor, 1.0, metalnessFactor );
	#else
		float specularIntensityFactor = 1.0;
		vec3 specularColorFactor = vec3( 1.0 );
		material.specularF90 = 1.0;
	#endif
	material.specularColor = min( pow2( ( material.ior - 1.0 ) / ( material.ior + 1.0 ) ) * specularColorFactor, vec3( 1.0 ) ) * specularIntensityFactor;
	material.specularColorBlended = mix( material.specularColor, diffuseColor.rgb, metalnessFactor );
#else
	material.specularColor = vec3( 0.04 );
	material.specularColorBlended = mix( material.specularColor, diffuseColor.rgb, metalnessFactor );
	material.specularF90 = 1.0;
#endif
#ifdef USE_CLEARCOAT
	material.clearcoat = clearcoat;
	material.clearcoatRoughness = clearcoatRoughness;
	material.clearcoatF0 = vec3( 0.04 );
	material.clearcoatF90 = 1.0;
	#ifdef USE_CLEARCOATMAP
		material.clearcoat *= texture2D( clearcoatMap, vClearcoatMapUv ).x;
	#endif
	#ifdef USE_CLEARCOAT_ROUGHNESSMAP
		material.clearcoatRoughness *= texture2D( clearcoatRoughnessMap, vClearcoatRoughnessMapUv ).y;
	#endif
	material.clearcoat = saturate( material.clearcoat );	material.clearcoatRoughness = max( material.clearcoatRoughness, 0.0525 );
	material.clearcoatRoughness += geometryRoughness;
	material.clearcoatRoughness = min( material.clearcoatRoughness, 1.0 );
#endif
#ifdef USE_DISPERSION
	material.dispersion = dispersion;
#endif
#ifdef USE_IRIDESCENCE
	material.iridescence = iridescence;
	material.iridescenceIOR = iridescenceIOR;
	#ifdef USE_IRIDESCENCEMAP
		material.iridescence *= texture2D( iridescenceMap, vIridescenceMapUv ).r;
	#endif
	#ifdef USE_IRIDESCENCE_THICKNESSMAP
		material.iridescenceThickness = (iridescenceThicknessMaximum - iridescenceThicknessMinimum) * texture2D( iridescenceThicknessMap, vIridescenceThicknessMapUv ).g + iridescenceThicknessMinimum;
	#else
		material.iridescenceThickness = iridescenceThicknessMaximum;
	#endif
#endif
#ifdef USE_SHEEN
	material.sheenColor = sheenColor;
	#ifdef USE_SHEEN_COLORMAP
		material.sheenColor *= texture2D( sheenColorMap, vSheenColorMapUv ).rgb;
	#endif
	material.sheenRoughness = clamp( sheenRoughness, 0.0001, 1.0 );
	#ifdef USE_SHEEN_ROUGHNESSMAP
		material.sheenRoughness *= texture2D( sheenRoughnessMap, vSheenRoughnessMapUv ).a;
	#endif
#endif
#ifdef USE_ANISOTROPY
	#ifdef USE_ANISOTROPYMAP
		mat2 anisotropyMat = mat2( anisotropyVector.x, anisotropyVector.y, - anisotropyVector.y, anisotropyVector.x );
		vec3 anisotropyPolar = texture2D( anisotropyMap, vAnisotropyMapUv ).rgb;
		vec2 anisotropyV = anisotropyMat * normalize( 2.0 * anisotropyPolar.rg - vec2( 1.0 ) ) * anisotropyPolar.b;
	#else
		vec2 anisotropyV = anisotropyVector;
	#endif
	material.anisotropy = length( anisotropyV );
	if( material.anisotropy == 0.0 ) {
		anisotropyV = vec2( 1.0, 0.0 );
	} else {
		anisotropyV /= material.anisotropy;
		material.anisotropy = saturate( material.anisotropy );
	}
	material.alphaT = mix( pow2( material.roughness ), 1.0, pow2( material.anisotropy ) );
	material.anisotropyT = tbn[ 0 ] * anisotropyV.x + tbn[ 1 ] * anisotropyV.y;
	material.anisotropyB = tbn[ 1 ] * anisotropyV.x - tbn[ 0 ] * anisotropyV.y;
#endif`,F3=`uniform sampler2D dfgLUT;
struct PhysicalMaterial {
	vec3 diffuseColor;
	vec3 diffuseContribution;
	vec3 specularColor;
	vec3 specularColorBlended;
	float roughness;
	float metalness;
	float specularF90;
	float dispersion;
	#ifdef USE_CLEARCOAT
		float clearcoat;
		float clearcoatRoughness;
		vec3 clearcoatF0;
		float clearcoatF90;
	#endif
	#ifdef USE_IRIDESCENCE
		float iridescence;
		float iridescenceIOR;
		float iridescenceThickness;
		vec3 iridescenceFresnel;
		vec3 iridescenceF0;
		vec3 iridescenceFresnelDielectric;
		vec3 iridescenceFresnelMetallic;
	#endif
	#ifdef USE_SHEEN
		vec3 sheenColor;
		float sheenRoughness;
	#endif
	#ifdef IOR
		float ior;
	#endif
	#ifdef USE_TRANSMISSION
		float transmission;
		float transmissionAlpha;
		float thickness;
		float attenuationDistance;
		vec3 attenuationColor;
	#endif
	#ifdef USE_ANISOTROPY
		float anisotropy;
		float alphaT;
		vec3 anisotropyT;
		vec3 anisotropyB;
	#endif
};
vec3 clearcoatSpecularDirect = vec3( 0.0 );
vec3 clearcoatSpecularIndirect = vec3( 0.0 );
vec3 sheenSpecularDirect = vec3( 0.0 );
vec3 sheenSpecularIndirect = vec3(0.0 );
vec3 Schlick_to_F0( const in vec3 f, const in float f90, const in float dotVH ) {
    float x = clamp( 1.0 - dotVH, 0.0, 1.0 );
    float x2 = x * x;
    float x5 = clamp( x * x2 * x2, 0.0, 0.9999 );
    return ( f - vec3( f90 ) * x5 ) / ( 1.0 - x5 );
}
float V_GGX_SmithCorrelated( const in float alpha, const in float dotNL, const in float dotNV ) {
	float a2 = pow2( alpha );
	float gv = dotNL * sqrt( a2 + ( 1.0 - a2 ) * pow2( dotNV ) );
	float gl = dotNV * sqrt( a2 + ( 1.0 - a2 ) * pow2( dotNL ) );
	return 0.5 / max( gv + gl, EPSILON );
}
float D_GGX( const in float alpha, const in float dotNH ) {
	float a2 = pow2( alpha );
	float denom = pow2( dotNH ) * ( a2 - 1.0 ) + 1.0;
	return RECIPROCAL_PI * a2 / pow2( denom );
}
#ifdef USE_ANISOTROPY
	float V_GGX_SmithCorrelated_Anisotropic( const in float alphaT, const in float alphaB, const in float dotTV, const in float dotBV, const in float dotTL, const in float dotBL, const in float dotNV, const in float dotNL ) {
		float gv = dotNL * length( vec3( alphaT * dotTV, alphaB * dotBV, dotNV ) );
		float gl = dotNV * length( vec3( alphaT * dotTL, alphaB * dotBL, dotNL ) );
		float v = 0.5 / ( gv + gl );
		return v;
	}
	float D_GGX_Anisotropic( const in float alphaT, const in float alphaB, const in float dotNH, const in float dotTH, const in float dotBH ) {
		float a2 = alphaT * alphaB;
		highp vec3 v = vec3( alphaB * dotTH, alphaT * dotBH, a2 * dotNH );
		highp float v2 = dot( v, v );
		float w2 = a2 / v2;
		return RECIPROCAL_PI * a2 * pow2 ( w2 );
	}
#endif
#ifdef USE_CLEARCOAT
	vec3 BRDF_GGX_Clearcoat( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, const in PhysicalMaterial material) {
		vec3 f0 = material.clearcoatF0;
		float f90 = material.clearcoatF90;
		float roughness = material.clearcoatRoughness;
		float alpha = pow2( roughness );
		vec3 halfDir = normalize( lightDir + viewDir );
		float dotNL = saturate( dot( normal, lightDir ) );
		float dotNV = saturate( dot( normal, viewDir ) );
		float dotNH = saturate( dot( normal, halfDir ) );
		float dotVH = saturate( dot( viewDir, halfDir ) );
		vec3 F = F_Schlick( f0, f90, dotVH );
		float V = V_GGX_SmithCorrelated( alpha, dotNL, dotNV );
		float D = D_GGX( alpha, dotNH );
		return F * ( V * D );
	}
#endif
vec3 BRDF_GGX( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, const in PhysicalMaterial material ) {
	vec3 f0 = material.specularColorBlended;
	float f90 = material.specularF90;
	float roughness = material.roughness;
	float alpha = pow2( roughness );
	vec3 halfDir = normalize( lightDir + viewDir );
	float dotNL = saturate( dot( normal, lightDir ) );
	float dotNV = saturate( dot( normal, viewDir ) );
	float dotNH = saturate( dot( normal, halfDir ) );
	float dotVH = saturate( dot( viewDir, halfDir ) );
	vec3 F = F_Schlick( f0, f90, dotVH );
	#ifdef USE_IRIDESCENCE
		F = mix( F, material.iridescenceFresnel, material.iridescence );
	#endif
	#ifdef USE_ANISOTROPY
		float dotTL = dot( material.anisotropyT, lightDir );
		float dotTV = dot( material.anisotropyT, viewDir );
		float dotTH = dot( material.anisotropyT, halfDir );
		float dotBL = dot( material.anisotropyB, lightDir );
		float dotBV = dot( material.anisotropyB, viewDir );
		float dotBH = dot( material.anisotropyB, halfDir );
		float V = V_GGX_SmithCorrelated_Anisotropic( material.alphaT, alpha, dotTV, dotBV, dotTL, dotBL, dotNV, dotNL );
		float D = D_GGX_Anisotropic( material.alphaT, alpha, dotNH, dotTH, dotBH );
	#else
		float V = V_GGX_SmithCorrelated( alpha, dotNL, dotNV );
		float D = D_GGX( alpha, dotNH );
	#endif
	return F * ( V * D );
}
vec2 LTC_Uv( const in vec3 N, const in vec3 V, const in float roughness ) {
	const float LUT_SIZE = 64.0;
	const float LUT_SCALE = ( LUT_SIZE - 1.0 ) / LUT_SIZE;
	const float LUT_BIAS = 0.5 / LUT_SIZE;
	float dotNV = saturate( dot( N, V ) );
	vec2 uv = vec2( roughness, sqrt( 1.0 - dotNV ) );
	uv = uv * LUT_SCALE + LUT_BIAS;
	return uv;
}
float LTC_ClippedSphereFormFactor( const in vec3 f ) {
	float l = length( f );
	return max( ( l * l + f.z ) / ( l + 1.0 ), 0.0 );
}
vec3 LTC_EdgeVectorFormFactor( const in vec3 v1, const in vec3 v2 ) {
	float x = dot( v1, v2 );
	float y = abs( x );
	float a = 0.8543985 + ( 0.4965155 + 0.0145206 * y ) * y;
	float b = 3.4175940 + ( 4.1616724 + y ) * y;
	float v = a / b;
	float theta_sintheta = ( x > 0.0 ) ? v : 0.5 * inversesqrt( max( 1.0 - x * x, 1e-7 ) ) - v;
	return cross( v1, v2 ) * theta_sintheta;
}
vec3 LTC_Evaluate( const in vec3 N, const in vec3 V, const in vec3 P, const in mat3 mInv, const in vec3 rectCoords[ 4 ] ) {
	vec3 v1 = rectCoords[ 1 ] - rectCoords[ 0 ];
	vec3 v2 = rectCoords[ 3 ] - rectCoords[ 0 ];
	vec3 lightNormal = cross( v1, v2 );
	if( dot( lightNormal, P - rectCoords[ 0 ] ) < 0.0 ) return vec3( 0.0 );
	vec3 T1, T2;
	T1 = normalize( V - N * dot( V, N ) );
	T2 = - cross( N, T1 );
	mat3 mat = mInv * transpose( mat3( T1, T2, N ) );
	vec3 coords[ 4 ];
	coords[ 0 ] = mat * ( rectCoords[ 0 ] - P );
	coords[ 1 ] = mat * ( rectCoords[ 1 ] - P );
	coords[ 2 ] = mat * ( rectCoords[ 2 ] - P );
	coords[ 3 ] = mat * ( rectCoords[ 3 ] - P );
	coords[ 0 ] = normalize( coords[ 0 ] );
	coords[ 1 ] = normalize( coords[ 1 ] );
	coords[ 2 ] = normalize( coords[ 2 ] );
	coords[ 3 ] = normalize( coords[ 3 ] );
	vec3 vectorFormFactor = vec3( 0.0 );
	vectorFormFactor += LTC_EdgeVectorFormFactor( coords[ 0 ], coords[ 1 ] );
	vectorFormFactor += LTC_EdgeVectorFormFactor( coords[ 1 ], coords[ 2 ] );
	vectorFormFactor += LTC_EdgeVectorFormFactor( coords[ 2 ], coords[ 3 ] );
	vectorFormFactor += LTC_EdgeVectorFormFactor( coords[ 3 ], coords[ 0 ] );
	float result = LTC_ClippedSphereFormFactor( vectorFormFactor );
	return vec3( result );
}
#if defined( USE_SHEEN )
float D_Charlie( float roughness, float dotNH ) {
	float alpha = pow2( roughness );
	float invAlpha = 1.0 / alpha;
	float cos2h = dotNH * dotNH;
	float sin2h = max( 1.0 - cos2h, 0.0078125 );
	return ( 2.0 + invAlpha ) * pow( sin2h, invAlpha * 0.5 ) / ( 2.0 * PI );
}
float V_Neubelt( float dotNV, float dotNL ) {
	return saturate( 1.0 / ( 4.0 * ( dotNL + dotNV - dotNL * dotNV ) ) );
}
vec3 BRDF_Sheen( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, vec3 sheenColor, const in float sheenRoughness ) {
	vec3 halfDir = normalize( lightDir + viewDir );
	float dotNL = saturate( dot( normal, lightDir ) );
	float dotNV = saturate( dot( normal, viewDir ) );
	float dotNH = saturate( dot( normal, halfDir ) );
	float D = D_Charlie( sheenRoughness, dotNH );
	float V = V_Neubelt( dotNV, dotNL );
	return sheenColor * ( D * V );
}
#endif
float IBLSheenBRDF( const in vec3 normal, const in vec3 viewDir, const in float roughness ) {
	float dotNV = saturate( dot( normal, viewDir ) );
	float r2 = roughness * roughness;
	float rInv = 1.0 / ( roughness + 0.1 );
	float a = -1.9362 + 1.0678 * roughness + 0.4573 * r2 - 0.8469 * rInv;
	float b = -0.6014 + 0.5538 * roughness - 0.4670 * r2 - 0.1255 * rInv;
	float DG = exp( a * dotNV + b );
	return saturate( DG );
}
vec3 EnvironmentBRDF( const in vec3 normal, const in vec3 viewDir, const in vec3 specularColor, const in float specularF90, const in float roughness ) {
	float dotNV = saturate( dot( normal, viewDir ) );
	vec2 fab = texture2D( dfgLUT, vec2( roughness, dotNV ) ).rg;
	return specularColor * fab.x + specularF90 * fab.y;
}
#ifdef USE_IRIDESCENCE
void computeMultiscatteringIridescence( const in vec3 normal, const in vec3 viewDir, const in vec3 specularColor, const in float specularF90, const in float iridescence, const in vec3 iridescenceF0, const in float roughness, inout vec3 singleScatter, inout vec3 multiScatter ) {
#else
void computeMultiscattering( const in vec3 normal, const in vec3 viewDir, const in vec3 specularColor, const in float specularF90, const in float roughness, inout vec3 singleScatter, inout vec3 multiScatter ) {
#endif
	float dotNV = saturate( dot( normal, viewDir ) );
	vec2 fab = texture2D( dfgLUT, vec2( roughness, dotNV ) ).rg;
	#ifdef USE_IRIDESCENCE
		vec3 Fr = mix( specularColor, iridescenceF0, iridescence );
	#else
		vec3 Fr = specularColor;
	#endif
	vec3 FssEss = Fr * fab.x + specularF90 * fab.y;
	float Ess = fab.x + fab.y;
	float Ems = 1.0 - Ess;
	vec3 Favg = Fr + ( 1.0 - Fr ) * 0.047619;	vec3 Fms = FssEss * Favg / ( 1.0 - Ems * Favg );
	singleScatter += FssEss;
	multiScatter += Fms * Ems;
}
vec3 BRDF_GGX_Multiscatter( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, const in PhysicalMaterial material ) {
	vec3 singleScatter = BRDF_GGX( lightDir, viewDir, normal, material );
	float dotNL = saturate( dot( normal, lightDir ) );
	float dotNV = saturate( dot( normal, viewDir ) );
	vec2 dfgV = texture2D( dfgLUT, vec2( material.roughness, dotNV ) ).rg;
	vec2 dfgL = texture2D( dfgLUT, vec2( material.roughness, dotNL ) ).rg;
	vec3 FssEss_V = material.specularColorBlended * dfgV.x + material.specularF90 * dfgV.y;
	vec3 FssEss_L = material.specularColorBlended * dfgL.x + material.specularF90 * dfgL.y;
	float Ess_V = dfgV.x + dfgV.y;
	float Ess_L = dfgL.x + dfgL.y;
	float Ems_V = 1.0 - Ess_V;
	float Ems_L = 1.0 - Ess_L;
	vec3 Favg = material.specularColorBlended + ( 1.0 - material.specularColorBlended ) * 0.047619;
	vec3 Fms = FssEss_V * FssEss_L * Favg / ( 1.0 - Ems_V * Ems_L * Favg + EPSILON );
	float compensationFactor = Ems_V * Ems_L;
	vec3 multiScatter = Fms * compensationFactor;
	return singleScatter + multiScatter;
}
#if NUM_RECT_AREA_LIGHTS > 0
	void RE_Direct_RectArea_Physical( const in RectAreaLight rectAreaLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in PhysicalMaterial material, inout ReflectedLight reflectedLight ) {
		vec3 normal = geometryNormal;
		vec3 viewDir = geometryViewDir;
		vec3 position = geometryPosition;
		vec3 lightPos = rectAreaLight.position;
		vec3 halfWidth = rectAreaLight.halfWidth;
		vec3 halfHeight = rectAreaLight.halfHeight;
		vec3 lightColor = rectAreaLight.color;
		float roughness = material.roughness;
		vec3 rectCoords[ 4 ];
		rectCoords[ 0 ] = lightPos + halfWidth - halfHeight;		rectCoords[ 1 ] = lightPos - halfWidth - halfHeight;
		rectCoords[ 2 ] = lightPos - halfWidth + halfHeight;
		rectCoords[ 3 ] = lightPos + halfWidth + halfHeight;
		vec2 uv = LTC_Uv( normal, viewDir, roughness );
		vec4 t1 = texture2D( ltc_1, uv );
		vec4 t2 = texture2D( ltc_2, uv );
		mat3 mInv = mat3(
			vec3( t1.x, 0, t1.y ),
			vec3(    0, 1,    0 ),
			vec3( t1.z, 0, t1.w )
		);
		vec3 fresnel = ( material.specularColorBlended * t2.x + ( material.specularF90 - material.specularColorBlended ) * t2.y );
		reflectedLight.directSpecular += lightColor * fresnel * LTC_Evaluate( normal, viewDir, position, mInv, rectCoords );
		reflectedLight.directDiffuse += lightColor * material.diffuseContribution * LTC_Evaluate( normal, viewDir, position, mat3( 1.0 ), rectCoords );
		#ifdef USE_CLEARCOAT
			vec3 Ncc = geometryClearcoatNormal;
			vec2 uvClearcoat = LTC_Uv( Ncc, viewDir, material.clearcoatRoughness );
			vec4 t1Clearcoat = texture2D( ltc_1, uvClearcoat );
			vec4 t2Clearcoat = texture2D( ltc_2, uvClearcoat );
			mat3 mInvClearcoat = mat3(
				vec3( t1Clearcoat.x, 0, t1Clearcoat.y ),
				vec3(             0, 1,             0 ),
				vec3( t1Clearcoat.z, 0, t1Clearcoat.w )
			);
			vec3 fresnelClearcoat = material.clearcoatF0 * t2Clearcoat.x + ( material.clearcoatF90 - material.clearcoatF0 ) * t2Clearcoat.y;
			clearcoatSpecularDirect += lightColor * fresnelClearcoat * LTC_Evaluate( Ncc, viewDir, position, mInvClearcoat, rectCoords );
		#endif
	}
#endif
void RE_Direct_Physical( const in IncidentLight directLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in PhysicalMaterial material, inout ReflectedLight reflectedLight ) {
	float dotNL = saturate( dot( geometryNormal, directLight.direction ) );
	vec3 irradiance = dotNL * directLight.color;
	#ifdef USE_CLEARCOAT
		float dotNLcc = saturate( dot( geometryClearcoatNormal, directLight.direction ) );
		vec3 ccIrradiance = dotNLcc * directLight.color;
		clearcoatSpecularDirect += ccIrradiance * BRDF_GGX_Clearcoat( directLight.direction, geometryViewDir, geometryClearcoatNormal, material );
	#endif
	#ifdef USE_SHEEN
 
 		sheenSpecularDirect += irradiance * BRDF_Sheen( directLight.direction, geometryViewDir, geometryNormal, material.sheenColor, material.sheenRoughness );
 
 		float sheenAlbedoV = IBLSheenBRDF( geometryNormal, geometryViewDir, material.sheenRoughness );
 		float sheenAlbedoL = IBLSheenBRDF( geometryNormal, directLight.direction, material.sheenRoughness );
 
 		float sheenEnergyComp = 1.0 - max3( material.sheenColor ) * max( sheenAlbedoV, sheenAlbedoL );
 
 		irradiance *= sheenEnergyComp;
 
 	#endif
	reflectedLight.directSpecular += irradiance * BRDF_GGX_Multiscatter( directLight.direction, geometryViewDir, geometryNormal, material );
	reflectedLight.directDiffuse += irradiance * BRDF_Lambert( material.diffuseContribution );
}
void RE_IndirectDiffuse_Physical( const in vec3 irradiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in PhysicalMaterial material, inout ReflectedLight reflectedLight ) {
	vec3 diffuse = irradiance * BRDF_Lambert( material.diffuseContribution );
	#ifdef USE_SHEEN
		float sheenAlbedo = IBLSheenBRDF( geometryNormal, geometryViewDir, material.sheenRoughness );
		float sheenEnergyComp = 1.0 - max3( material.sheenColor ) * sheenAlbedo;
		diffuse *= sheenEnergyComp;
	#endif
	reflectedLight.indirectDiffuse += diffuse;
}
void RE_IndirectSpecular_Physical( const in vec3 radiance, const in vec3 irradiance, const in vec3 clearcoatRadiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in PhysicalMaterial material, inout ReflectedLight reflectedLight) {
	#ifdef USE_CLEARCOAT
		clearcoatSpecularIndirect += clearcoatRadiance * EnvironmentBRDF( geometryClearcoatNormal, geometryViewDir, material.clearcoatF0, material.clearcoatF90, material.clearcoatRoughness );
	#endif
	#ifdef USE_SHEEN
		sheenSpecularIndirect += irradiance * material.sheenColor * IBLSheenBRDF( geometryNormal, geometryViewDir, material.sheenRoughness ) * RECIPROCAL_PI;
 	#endif
	vec3 singleScatteringDielectric = vec3( 0.0 );
	vec3 multiScatteringDielectric = vec3( 0.0 );
	vec3 singleScatteringMetallic = vec3( 0.0 );
	vec3 multiScatteringMetallic = vec3( 0.0 );
	#ifdef USE_IRIDESCENCE
		computeMultiscatteringIridescence( geometryNormal, geometryViewDir, material.specularColor, material.specularF90, material.iridescence, material.iridescenceFresnelDielectric, material.roughness, singleScatteringDielectric, multiScatteringDielectric );
		computeMultiscatteringIridescence( geometryNormal, geometryViewDir, material.diffuseColor, material.specularF90, material.iridescence, material.iridescenceFresnelMetallic, material.roughness, singleScatteringMetallic, multiScatteringMetallic );
	#else
		computeMultiscattering( geometryNormal, geometryViewDir, material.specularColor, material.specularF90, material.roughness, singleScatteringDielectric, multiScatteringDielectric );
		computeMultiscattering( geometryNormal, geometryViewDir, material.diffuseColor, material.specularF90, material.roughness, singleScatteringMetallic, multiScatteringMetallic );
	#endif
	vec3 singleScattering = mix( singleScatteringDielectric, singleScatteringMetallic, material.metalness );
	vec3 multiScattering = mix( multiScatteringDielectric, multiScatteringMetallic, material.metalness );
	vec3 totalScatteringDielectric = singleScatteringDielectric + multiScatteringDielectric;
	vec3 diffuse = material.diffuseContribution * ( 1.0 - totalScatteringDielectric );
	vec3 cosineWeightedIrradiance = irradiance * RECIPROCAL_PI;
	vec3 indirectSpecular = radiance * singleScattering;
	indirectSpecular += multiScattering * cosineWeightedIrradiance;
	vec3 indirectDiffuse = diffuse * cosineWeightedIrradiance;
	#ifdef USE_SHEEN
		float sheenAlbedo = IBLSheenBRDF( geometryNormal, geometryViewDir, material.sheenRoughness );
		float sheenEnergyComp = 1.0 - max3( material.sheenColor ) * sheenAlbedo;
		indirectSpecular *= sheenEnergyComp;
		indirectDiffuse *= sheenEnergyComp;
	#endif
	reflectedLight.indirectSpecular += indirectSpecular;
	reflectedLight.indirectDiffuse += indirectDiffuse;
}
#define RE_Direct				RE_Direct_Physical
#define RE_Direct_RectArea		RE_Direct_RectArea_Physical
#define RE_IndirectDiffuse		RE_IndirectDiffuse_Physical
#define RE_IndirectSpecular		RE_IndirectSpecular_Physical
float computeSpecularOcclusion( const in float dotNV, const in float ambientOcclusion, const in float roughness ) {
	return saturate( pow( dotNV + ambientOcclusion, exp2( - 16.0 * roughness - 1.0 ) ) - 1.0 + ambientOcclusion );
}`,V3=`
vec3 geometryPosition = - vViewPosition;
vec3 geometryNormal = normal;
vec3 geometryViewDir = ( isOrthographic ) ? vec3( 0, 0, 1 ) : normalize( vViewPosition );
vec3 geometryClearcoatNormal = vec3( 0.0 );
#ifdef USE_CLEARCOAT
	geometryClearcoatNormal = clearcoatNormal;
#endif
#ifdef USE_IRIDESCENCE
	float dotNVi = saturate( dot( normal, geometryViewDir ) );
	if ( material.iridescenceThickness == 0.0 ) {
		material.iridescence = 0.0;
	} else {
		material.iridescence = saturate( material.iridescence );
	}
	if ( material.iridescence > 0.0 ) {
		material.iridescenceFresnelDielectric = evalIridescence( 1.0, material.iridescenceIOR, dotNVi, material.iridescenceThickness, material.specularColor );
		material.iridescenceFresnelMetallic = evalIridescence( 1.0, material.iridescenceIOR, dotNVi, material.iridescenceThickness, material.diffuseColor );
		material.iridescenceFresnel = mix( material.iridescenceFresnelDielectric, material.iridescenceFresnelMetallic, material.metalness );
		material.iridescenceF0 = Schlick_to_F0( material.iridescenceFresnel, 1.0, dotNVi );
	}
#endif
IncidentLight directLight;
#if ( NUM_POINT_LIGHTS > 0 ) && defined( RE_Direct )
	PointLight pointLight;
	#if defined( USE_SHADOWMAP ) && NUM_POINT_LIGHT_SHADOWS > 0
	PointLightShadow pointLightShadow;
	#endif
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_POINT_LIGHTS; i ++ ) {
		pointLight = pointLights[ i ];
		getPointLightInfo( pointLight, geometryPosition, directLight );
		#if defined( USE_SHADOWMAP ) && ( UNROLLED_LOOP_INDEX < NUM_POINT_LIGHT_SHADOWS ) && ( defined( SHADOWMAP_TYPE_PCF ) || defined( SHADOWMAP_TYPE_BASIC ) )
		pointLightShadow = pointLightShadows[ i ];
		directLight.color *= ( directLight.visible && receiveShadow ) ? getPointShadow( pointShadowMap[ i ], pointLightShadow.shadowMapSize, pointLightShadow.shadowIntensity, pointLightShadow.shadowBias, pointLightShadow.shadowRadius, vPointShadowCoord[ i ], pointLightShadow.shadowCameraNear, pointLightShadow.shadowCameraFar ) : 1.0;
		#endif
		RE_Direct( directLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if ( NUM_SPOT_LIGHTS > 0 ) && defined( RE_Direct )
	SpotLight spotLight;
	vec4 spotColor;
	vec3 spotLightCoord;
	bool inSpotLightMap;
	#if defined( USE_SHADOWMAP ) && NUM_SPOT_LIGHT_SHADOWS > 0
	SpotLightShadow spotLightShadow;
	#endif
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_SPOT_LIGHTS; i ++ ) {
		spotLight = spotLights[ i ];
		getSpotLightInfo( spotLight, geometryPosition, directLight );
		#if ( UNROLLED_LOOP_INDEX < NUM_SPOT_LIGHT_SHADOWS_WITH_MAPS )
		#define SPOT_LIGHT_MAP_INDEX UNROLLED_LOOP_INDEX
		#elif ( UNROLLED_LOOP_INDEX < NUM_SPOT_LIGHT_SHADOWS )
		#define SPOT_LIGHT_MAP_INDEX NUM_SPOT_LIGHT_MAPS
		#else
		#define SPOT_LIGHT_MAP_INDEX ( UNROLLED_LOOP_INDEX - NUM_SPOT_LIGHT_SHADOWS + NUM_SPOT_LIGHT_SHADOWS_WITH_MAPS )
		#endif
		#if ( SPOT_LIGHT_MAP_INDEX < NUM_SPOT_LIGHT_MAPS )
			spotLightCoord = vSpotLightCoord[ i ].xyz / vSpotLightCoord[ i ].w;
			inSpotLightMap = all( lessThan( abs( spotLightCoord * 2. - 1. ), vec3( 1.0 ) ) );
			spotColor = texture2D( spotLightMap[ SPOT_LIGHT_MAP_INDEX ], spotLightCoord.xy );
			directLight.color = inSpotLightMap ? directLight.color * spotColor.rgb : directLight.color;
		#endif
		#undef SPOT_LIGHT_MAP_INDEX
		#if defined( USE_SHADOWMAP ) && ( UNROLLED_LOOP_INDEX < NUM_SPOT_LIGHT_SHADOWS )
		spotLightShadow = spotLightShadows[ i ];
		directLight.color *= ( directLight.visible && receiveShadow ) ? getShadow( spotShadowMap[ i ], spotLightShadow.shadowMapSize, spotLightShadow.shadowIntensity, spotLightShadow.shadowBias, spotLightShadow.shadowRadius, vSpotLightCoord[ i ] ) : 1.0;
		#endif
		RE_Direct( directLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if ( NUM_DIR_LIGHTS > 0 ) && defined( RE_Direct )
	DirectionalLight directionalLight;
	#if defined( USE_SHADOWMAP ) && NUM_DIR_LIGHT_SHADOWS > 0
	DirectionalLightShadow directionalLightShadow;
	#endif
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_DIR_LIGHTS; i ++ ) {
		directionalLight = directionalLights[ i ];
		getDirectionalLightInfo( directionalLight, directLight );
		#if defined( USE_SHADOWMAP ) && ( UNROLLED_LOOP_INDEX < NUM_DIR_LIGHT_SHADOWS )
		directionalLightShadow = directionalLightShadows[ i ];
		directLight.color *= ( directLight.visible && receiveShadow ) ? getShadow( directionalShadowMap[ i ], directionalLightShadow.shadowMapSize, directionalLightShadow.shadowIntensity, directionalLightShadow.shadowBias, directionalLightShadow.shadowRadius, vDirectionalShadowCoord[ i ] ) : 1.0;
		#endif
		RE_Direct( directLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if ( NUM_RECT_AREA_LIGHTS > 0 ) && defined( RE_Direct_RectArea )
	RectAreaLight rectAreaLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_RECT_AREA_LIGHTS; i ++ ) {
		rectAreaLight = rectAreaLights[ i ];
		RE_Direct_RectArea( rectAreaLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if defined( RE_IndirectDiffuse )
	vec3 iblIrradiance = vec3( 0.0 );
	vec3 irradiance = getAmbientLightIrradiance( ambientLightColor );
	#if defined( USE_LIGHT_PROBES )
		irradiance += getLightProbeIrradiance( lightProbe, geometryNormal );
	#endif
	#if ( NUM_HEMI_LIGHTS > 0 )
		#pragma unroll_loop_start
		for ( int i = 0; i < NUM_HEMI_LIGHTS; i ++ ) {
			irradiance += getHemisphereLightIrradiance( hemisphereLights[ i ], geometryNormal );
		}
		#pragma unroll_loop_end
	#endif
#endif
#if defined( RE_IndirectSpecular )
	vec3 radiance = vec3( 0.0 );
	vec3 clearcoatRadiance = vec3( 0.0 );
#endif`,H3=`#if defined( RE_IndirectDiffuse )
	#ifdef USE_LIGHTMAP
		vec4 lightMapTexel = texture2D( lightMap, vLightMapUv );
		vec3 lightMapIrradiance = lightMapTexel.rgb * lightMapIntensity;
		irradiance += lightMapIrradiance;
	#endif
	#if defined( USE_ENVMAP ) && defined( ENVMAP_TYPE_CUBE_UV )
		#if defined( STANDARD ) || defined( LAMBERT ) || defined( PHONG )
			iblIrradiance += getIBLIrradiance( geometryNormal );
		#endif
	#endif
#endif
#if defined( USE_ENVMAP ) && defined( RE_IndirectSpecular )
	#ifdef USE_ANISOTROPY
		radiance += getIBLAnisotropyRadiance( geometryViewDir, geometryNormal, material.roughness, material.anisotropyB, material.anisotropy );
	#else
		radiance += getIBLRadiance( geometryViewDir, geometryNormal, material.roughness );
	#endif
	#ifdef USE_CLEARCOAT
		clearcoatRadiance += getIBLRadiance( geometryViewDir, geometryClearcoatNormal, material.clearcoatRoughness );
	#endif
#endif`,G3=`#if defined( RE_IndirectDiffuse )
	#if defined( LAMBERT ) || defined( PHONG )
		irradiance += iblIrradiance;
	#endif
	RE_IndirectDiffuse( irradiance, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
#endif
#if defined( RE_IndirectSpecular )
	RE_IndirectSpecular( radiance, iblIrradiance, clearcoatRadiance, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
#endif`,k3=`#if defined( USE_LOGARITHMIC_DEPTH_BUFFER )
	gl_FragDepth = vIsPerspective == 0.0 ? gl_FragCoord.z : log2( vFragDepth ) * logDepthBufFC * 0.5;
#endif`,X3=`#if defined( USE_LOGARITHMIC_DEPTH_BUFFER )
	uniform float logDepthBufFC;
	varying float vFragDepth;
	varying float vIsPerspective;
#endif`,W3=`#ifdef USE_LOGARITHMIC_DEPTH_BUFFER
	varying float vFragDepth;
	varying float vIsPerspective;
#endif`,q3=`#ifdef USE_LOGARITHMIC_DEPTH_BUFFER
	vFragDepth = 1.0 + gl_Position.w;
	vIsPerspective = float( isPerspectiveMatrix( projectionMatrix ) );
#endif`,Y3=`#ifdef USE_MAP
	vec4 sampledDiffuseColor = texture2D( map, vMapUv );
	#ifdef DECODE_VIDEO_TEXTURE
		sampledDiffuseColor = sRGBTransferEOTF( sampledDiffuseColor );
	#endif
	diffuseColor *= sampledDiffuseColor;
#endif`,Z3=`#ifdef USE_MAP
	uniform sampler2D map;
#endif`,J3=`#if defined( USE_MAP ) || defined( USE_ALPHAMAP )
	#if defined( USE_POINTS_UV )
		vec2 uv = vUv;
	#else
		vec2 uv = ( uvTransform * vec3( gl_PointCoord.x, 1.0 - gl_PointCoord.y, 1 ) ).xy;
	#endif
#endif
#ifdef USE_MAP
	diffuseColor *= texture2D( map, uv );
#endif
#ifdef USE_ALPHAMAP
	diffuseColor.a *= texture2D( alphaMap, uv ).g;
#endif`,K3=`#if defined( USE_POINTS_UV )
	varying vec2 vUv;
#else
	#if defined( USE_MAP ) || defined( USE_ALPHAMAP )
		uniform mat3 uvTransform;
	#endif
#endif
#ifdef USE_MAP
	uniform sampler2D map;
#endif
#ifdef USE_ALPHAMAP
	uniform sampler2D alphaMap;
#endif`,Q3=`float metalnessFactor = metalness;
#ifdef USE_METALNESSMAP
	vec4 texelMetalness = texture2D( metalnessMap, vMetalnessMapUv );
	metalnessFactor *= texelMetalness.b;
#endif`,j3=`#ifdef USE_METALNESSMAP
	uniform sampler2D metalnessMap;
#endif`,$3=`#ifdef USE_INSTANCING_MORPH
	float morphTargetInfluences[ MORPHTARGETS_COUNT ];
	float morphTargetBaseInfluence = texelFetch( morphTexture, ivec2( 0, gl_InstanceID ), 0 ).r;
	for ( int i = 0; i < MORPHTARGETS_COUNT; i ++ ) {
		morphTargetInfluences[i] =  texelFetch( morphTexture, ivec2( i + 1, gl_InstanceID ), 0 ).r;
	}
#endif`,t2=`#if defined( USE_MORPHCOLORS )
	vColor *= morphTargetBaseInfluence;
	for ( int i = 0; i < MORPHTARGETS_COUNT; i ++ ) {
		#if defined( USE_COLOR_ALPHA )
			if ( morphTargetInfluences[ i ] != 0.0 ) vColor += getMorph( gl_VertexID, i, 2 ) * morphTargetInfluences[ i ];
		#elif defined( USE_COLOR )
			if ( morphTargetInfluences[ i ] != 0.0 ) vColor += getMorph( gl_VertexID, i, 2 ).rgb * morphTargetInfluences[ i ];
		#endif
	}
#endif`,e2=`#ifdef USE_MORPHNORMALS
	objectNormal *= morphTargetBaseInfluence;
	for ( int i = 0; i < MORPHTARGETS_COUNT; i ++ ) {
		if ( morphTargetInfluences[ i ] != 0.0 ) objectNormal += getMorph( gl_VertexID, i, 1 ).xyz * morphTargetInfluences[ i ];
	}
#endif`,n2=`#ifdef USE_MORPHTARGETS
	#ifndef USE_INSTANCING_MORPH
		uniform float morphTargetBaseInfluence;
		uniform float morphTargetInfluences[ MORPHTARGETS_COUNT ];
	#endif
	uniform sampler2DArray morphTargetsTexture;
	uniform ivec2 morphTargetsTextureSize;
	vec4 getMorph( const in int vertexIndex, const in int morphTargetIndex, const in int offset ) {
		int texelIndex = vertexIndex * MORPHTARGETS_TEXTURE_STRIDE + offset;
		int y = texelIndex / morphTargetsTextureSize.x;
		int x = texelIndex - y * morphTargetsTextureSize.x;
		ivec3 morphUV = ivec3( x, y, morphTargetIndex );
		return texelFetch( morphTargetsTexture, morphUV, 0 );
	}
#endif`,i2=`#ifdef USE_MORPHTARGETS
	transformed *= morphTargetBaseInfluence;
	for ( int i = 0; i < MORPHTARGETS_COUNT; i ++ ) {
		if ( morphTargetInfluences[ i ] != 0.0 ) transformed += getMorph( gl_VertexID, i, 0 ).xyz * morphTargetInfluences[ i ];
	}
#endif`,a2=`float faceDirection = gl_FrontFacing ? 1.0 : - 1.0;
#ifdef FLAT_SHADED
	vec3 fdx = dFdx( vViewPosition );
	vec3 fdy = dFdy( vViewPosition );
	vec3 normal = normalize( cross( fdx, fdy ) );
#else
	vec3 normal = normalize( vNormal );
	#ifdef DOUBLE_SIDED
		normal *= faceDirection;
	#endif
#endif
#if defined( USE_NORMALMAP_TANGENTSPACE ) || defined( USE_CLEARCOAT_NORMALMAP ) || defined( USE_ANISOTROPY )
	#ifdef USE_TANGENT
		mat3 tbn = mat3( normalize( vTangent ), normalize( vBitangent ), normal );
	#else
		mat3 tbn = getTangentFrame( - vViewPosition, normal,
		#if defined( USE_NORMALMAP )
			vNormalMapUv
		#elif defined( USE_CLEARCOAT_NORMALMAP )
			vClearcoatNormalMapUv
		#else
			vUv
		#endif
		);
	#endif
	#if defined( DOUBLE_SIDED ) && ! defined( FLAT_SHADED )
		tbn[0] *= faceDirection;
		tbn[1] *= faceDirection;
	#endif
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	#ifdef USE_TANGENT
		mat3 tbn2 = mat3( normalize( vTangent ), normalize( vBitangent ), normal );
	#else
		mat3 tbn2 = getTangentFrame( - vViewPosition, normal, vClearcoatNormalMapUv );
	#endif
	#if defined( DOUBLE_SIDED ) && ! defined( FLAT_SHADED )
		tbn2[0] *= faceDirection;
		tbn2[1] *= faceDirection;
	#endif
#endif
vec3 nonPerturbedNormal = normal;`,s2=`#ifdef USE_NORMALMAP_OBJECTSPACE
	normal = texture2D( normalMap, vNormalMapUv ).xyz * 2.0 - 1.0;
	#ifdef FLIP_SIDED
		normal = - normal;
	#endif
	#ifdef DOUBLE_SIDED
		normal = normal * faceDirection;
	#endif
	normal = normalize( normalMatrix * normal );
#elif defined( USE_NORMALMAP_TANGENTSPACE )
	vec3 mapN = texture2D( normalMap, vNormalMapUv ).xyz * 2.0 - 1.0;
	mapN.xy *= normalScale;
	normal = normalize( tbn * mapN );
#elif defined( USE_BUMPMAP )
	normal = perturbNormalArb( - vViewPosition, normal, dHdxy_fwd(), faceDirection );
#endif`,r2=`#ifndef FLAT_SHADED
	varying vec3 vNormal;
	#ifdef USE_TANGENT
		varying vec3 vTangent;
		varying vec3 vBitangent;
	#endif
#endif`,o2=`#ifndef FLAT_SHADED
	varying vec3 vNormal;
	#ifdef USE_TANGENT
		varying vec3 vTangent;
		varying vec3 vBitangent;
	#endif
#endif`,l2=`#ifndef FLAT_SHADED
	vNormal = normalize( transformedNormal );
	#ifdef USE_TANGENT
		vTangent = normalize( transformedTangent );
		vBitangent = normalize( cross( vNormal, vTangent ) * tangent.w );
	#endif
#endif`,c2=`#ifdef USE_NORMALMAP
	uniform sampler2D normalMap;
	uniform vec2 normalScale;
#endif
#ifdef USE_NORMALMAP_OBJECTSPACE
	uniform mat3 normalMatrix;
#endif
#if ! defined ( USE_TANGENT ) && ( defined ( USE_NORMALMAP_TANGENTSPACE ) || defined ( USE_CLEARCOAT_NORMALMAP ) || defined( USE_ANISOTROPY ) )
	mat3 getTangentFrame( vec3 eye_pos, vec3 surf_norm, vec2 uv ) {
		vec3 q0 = dFdx( eye_pos.xyz );
		vec3 q1 = dFdy( eye_pos.xyz );
		vec2 st0 = dFdx( uv.st );
		vec2 st1 = dFdy( uv.st );
		vec3 N = surf_norm;
		vec3 q1perp = cross( q1, N );
		vec3 q0perp = cross( N, q0 );
		vec3 T = q1perp * st0.x + q0perp * st1.x;
		vec3 B = q1perp * st0.y + q0perp * st1.y;
		float det = max( dot( T, T ), dot( B, B ) );
		float scale = ( det == 0.0 ) ? 0.0 : inversesqrt( det );
		return mat3( T * scale, B * scale, N );
	}
#endif`,u2=`#ifdef USE_CLEARCOAT
	vec3 clearcoatNormal = nonPerturbedNormal;
#endif`,f2=`#ifdef USE_CLEARCOAT_NORMALMAP
	vec3 clearcoatMapN = texture2D( clearcoatNormalMap, vClearcoatNormalMapUv ).xyz * 2.0 - 1.0;
	clearcoatMapN.xy *= clearcoatNormalScale;
	clearcoatNormal = normalize( tbn2 * clearcoatMapN );
#endif`,h2=`#ifdef USE_CLEARCOATMAP
	uniform sampler2D clearcoatMap;
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	uniform sampler2D clearcoatNormalMap;
	uniform vec2 clearcoatNormalScale;
#endif
#ifdef USE_CLEARCOAT_ROUGHNESSMAP
	uniform sampler2D clearcoatRoughnessMap;
#endif`,d2=`#ifdef USE_IRIDESCENCEMAP
	uniform sampler2D iridescenceMap;
#endif
#ifdef USE_IRIDESCENCE_THICKNESSMAP
	uniform sampler2D iridescenceThicknessMap;
#endif`,p2=`#ifdef OPAQUE
diffuseColor.a = 1.0;
#endif
#ifdef USE_TRANSMISSION
diffuseColor.a *= material.transmissionAlpha;
#endif
gl_FragColor = vec4( outgoingLight, diffuseColor.a );`,m2=`vec3 packNormalToRGB( const in vec3 normal ) {
	return normalize( normal ) * 0.5 + 0.5;
}
vec3 unpackRGBToNormal( const in vec3 rgb ) {
	return 2.0 * rgb.xyz - 1.0;
}
const float PackUpscale = 256. / 255.;const float UnpackDownscale = 255. / 256.;const float ShiftRight8 = 1. / 256.;
const float Inv255 = 1. / 255.;
const vec4 PackFactors = vec4( 1.0, 256.0, 256.0 * 256.0, 256.0 * 256.0 * 256.0 );
const vec2 UnpackFactors2 = vec2( UnpackDownscale, 1.0 / PackFactors.g );
const vec3 UnpackFactors3 = vec3( UnpackDownscale / PackFactors.rg, 1.0 / PackFactors.b );
const vec4 UnpackFactors4 = vec4( UnpackDownscale / PackFactors.rgb, 1.0 / PackFactors.a );
vec4 packDepthToRGBA( const in float v ) {
	if( v <= 0.0 )
		return vec4( 0., 0., 0., 0. );
	if( v >= 1.0 )
		return vec4( 1., 1., 1., 1. );
	float vuf;
	float af = modf( v * PackFactors.a, vuf );
	float bf = modf( vuf * ShiftRight8, vuf );
	float gf = modf( vuf * ShiftRight8, vuf );
	return vec4( vuf * Inv255, gf * PackUpscale, bf * PackUpscale, af );
}
vec3 packDepthToRGB( const in float v ) {
	if( v <= 0.0 )
		return vec3( 0., 0., 0. );
	if( v >= 1.0 )
		return vec3( 1., 1., 1. );
	float vuf;
	float bf = modf( v * PackFactors.b, vuf );
	float gf = modf( vuf * ShiftRight8, vuf );
	return vec3( vuf * Inv255, gf * PackUpscale, bf );
}
vec2 packDepthToRG( const in float v ) {
	if( v <= 0.0 )
		return vec2( 0., 0. );
	if( v >= 1.0 )
		return vec2( 1., 1. );
	float vuf;
	float gf = modf( v * 256., vuf );
	return vec2( vuf * Inv255, gf );
}
float unpackRGBAToDepth( const in vec4 v ) {
	return dot( v, UnpackFactors4 );
}
float unpackRGBToDepth( const in vec3 v ) {
	return dot( v, UnpackFactors3 );
}
float unpackRGToDepth( const in vec2 v ) {
	return v.r * UnpackFactors2.r + v.g * UnpackFactors2.g;
}
vec4 pack2HalfToRGBA( const in vec2 v ) {
	vec4 r = vec4( v.x, fract( v.x * 255.0 ), v.y, fract( v.y * 255.0 ) );
	return vec4( r.x - r.y / 255.0, r.y, r.z - r.w / 255.0, r.w );
}
vec2 unpackRGBATo2Half( const in vec4 v ) {
	return vec2( v.x + ( v.y / 255.0 ), v.z + ( v.w / 255.0 ) );
}
float viewZToOrthographicDepth( const in float viewZ, const in float near, const in float far ) {
	return ( viewZ + near ) / ( near - far );
}
float orthographicDepthToViewZ( const in float depth, const in float near, const in float far ) {
	#ifdef USE_REVERSED_DEPTH_BUFFER
	
		return depth * ( far - near ) - far;
	#else
		return depth * ( near - far ) - near;
	#endif
}
float viewZToPerspectiveDepth( const in float viewZ, const in float near, const in float far ) {
	return ( ( near + viewZ ) * far ) / ( ( far - near ) * viewZ );
}
float perspectiveDepthToViewZ( const in float depth, const in float near, const in float far ) {
	
	#ifdef USE_REVERSED_DEPTH_BUFFER
		return ( near * far ) / ( ( near - far ) * depth - near );
	#else
		return ( near * far ) / ( ( far - near ) * depth - far );
	#endif
}`,g2=`#ifdef PREMULTIPLIED_ALPHA
	gl_FragColor.rgb *= gl_FragColor.a;
#endif`,v2=`vec4 mvPosition = vec4( transformed, 1.0 );
#ifdef USE_BATCHING
	mvPosition = batchingMatrix * mvPosition;
#endif
#ifdef USE_INSTANCING
	mvPosition = instanceMatrix * mvPosition;
#endif
mvPosition = modelViewMatrix * mvPosition;
gl_Position = projectionMatrix * mvPosition;`,_2=`#ifdef DITHERING
	gl_FragColor.rgb = dithering( gl_FragColor.rgb );
#endif`,y2=`#ifdef DITHERING
	vec3 dithering( vec3 color ) {
		float grid_position = rand( gl_FragCoord.xy );
		vec3 dither_shift_RGB = vec3( 0.25 / 255.0, -0.25 / 255.0, 0.25 / 255.0 );
		dither_shift_RGB = mix( 2.0 * dither_shift_RGB, -2.0 * dither_shift_RGB, grid_position );
		return color + dither_shift_RGB;
	}
#endif`,x2=`float roughnessFactor = roughness;
#ifdef USE_ROUGHNESSMAP
	vec4 texelRoughness = texture2D( roughnessMap, vRoughnessMapUv );
	roughnessFactor *= texelRoughness.g;
#endif`,S2=`#ifdef USE_ROUGHNESSMAP
	uniform sampler2D roughnessMap;
#endif`,b2=`#if NUM_SPOT_LIGHT_COORDS > 0
	varying vec4 vSpotLightCoord[ NUM_SPOT_LIGHT_COORDS ];
#endif
#if NUM_SPOT_LIGHT_MAPS > 0
	uniform sampler2D spotLightMap[ NUM_SPOT_LIGHT_MAPS ];
#endif
#ifdef USE_SHADOWMAP
	#if NUM_DIR_LIGHT_SHADOWS > 0
		#if defined( SHADOWMAP_TYPE_PCF )
			uniform sampler2DShadow directionalShadowMap[ NUM_DIR_LIGHT_SHADOWS ];
		#else
			uniform sampler2D directionalShadowMap[ NUM_DIR_LIGHT_SHADOWS ];
		#endif
		varying vec4 vDirectionalShadowCoord[ NUM_DIR_LIGHT_SHADOWS ];
		struct DirectionalLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform DirectionalLightShadow directionalLightShadows[ NUM_DIR_LIGHT_SHADOWS ];
	#endif
	#if NUM_SPOT_LIGHT_SHADOWS > 0
		#if defined( SHADOWMAP_TYPE_PCF )
			uniform sampler2DShadow spotShadowMap[ NUM_SPOT_LIGHT_SHADOWS ];
		#else
			uniform sampler2D spotShadowMap[ NUM_SPOT_LIGHT_SHADOWS ];
		#endif
		struct SpotLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform SpotLightShadow spotLightShadows[ NUM_SPOT_LIGHT_SHADOWS ];
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0
		#if defined( SHADOWMAP_TYPE_PCF )
			uniform samplerCubeShadow pointShadowMap[ NUM_POINT_LIGHT_SHADOWS ];
		#elif defined( SHADOWMAP_TYPE_BASIC )
			uniform samplerCube pointShadowMap[ NUM_POINT_LIGHT_SHADOWS ];
		#endif
		varying vec4 vPointShadowCoord[ NUM_POINT_LIGHT_SHADOWS ];
		struct PointLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
			float shadowCameraNear;
			float shadowCameraFar;
		};
		uniform PointLightShadow pointLightShadows[ NUM_POINT_LIGHT_SHADOWS ];
	#endif
	#if defined( SHADOWMAP_TYPE_PCF )
		float interleavedGradientNoise( vec2 position ) {
			return fract( 52.9829189 * fract( dot( position, vec2( 0.06711056, 0.00583715 ) ) ) );
		}
		vec2 vogelDiskSample( int sampleIndex, int samplesCount, float phi ) {
			const float goldenAngle = 2.399963229728653;
			float r = sqrt( ( float( sampleIndex ) + 0.5 ) / float( samplesCount ) );
			float theta = float( sampleIndex ) * goldenAngle + phi;
			return vec2( cos( theta ), sin( theta ) ) * r;
		}
	#endif
	#if defined( SHADOWMAP_TYPE_PCF )
		float getShadow( sampler2DShadow shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord ) {
			float shadow = 1.0;
			shadowCoord.xyz /= shadowCoord.w;
			shadowCoord.z += shadowBias;
			bool inFrustum = shadowCoord.x >= 0.0 && shadowCoord.x <= 1.0 && shadowCoord.y >= 0.0 && shadowCoord.y <= 1.0;
			bool frustumTest = inFrustum && shadowCoord.z <= 1.0;
			if ( frustumTest ) {
				vec2 texelSize = vec2( 1.0 ) / shadowMapSize;
				float radius = shadowRadius * texelSize.x;
				float phi = interleavedGradientNoise( gl_FragCoord.xy ) * PI2;
				shadow = (
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 0, 5, phi ) * radius, shadowCoord.z ) ) +
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 1, 5, phi ) * radius, shadowCoord.z ) ) +
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 2, 5, phi ) * radius, shadowCoord.z ) ) +
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 3, 5, phi ) * radius, shadowCoord.z ) ) +
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 4, 5, phi ) * radius, shadowCoord.z ) )
				) * 0.2;
			}
			return mix( 1.0, shadow, shadowIntensity );
		}
	#elif defined( SHADOWMAP_TYPE_VSM )
		float getShadow( sampler2D shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord ) {
			float shadow = 1.0;
			shadowCoord.xyz /= shadowCoord.w;
			#ifdef USE_REVERSED_DEPTH_BUFFER
				shadowCoord.z -= shadowBias;
			#else
				shadowCoord.z += shadowBias;
			#endif
			bool inFrustum = shadowCoord.x >= 0.0 && shadowCoord.x <= 1.0 && shadowCoord.y >= 0.0 && shadowCoord.y <= 1.0;
			bool frustumTest = inFrustum && shadowCoord.z <= 1.0;
			if ( frustumTest ) {
				vec2 distribution = texture2D( shadowMap, shadowCoord.xy ).rg;
				float mean = distribution.x;
				float variance = distribution.y * distribution.y;
				#ifdef USE_REVERSED_DEPTH_BUFFER
					float hard_shadow = step( mean, shadowCoord.z );
				#else
					float hard_shadow = step( shadowCoord.z, mean );
				#endif
				
				if ( hard_shadow == 1.0 ) {
					shadow = 1.0;
				} else {
					variance = max( variance, 0.0000001 );
					float d = shadowCoord.z - mean;
					float p_max = variance / ( variance + d * d );
					p_max = clamp( ( p_max - 0.3 ) / 0.65, 0.0, 1.0 );
					shadow = max( hard_shadow, p_max );
				}
			}
			return mix( 1.0, shadow, shadowIntensity );
		}
	#else
		float getShadow( sampler2D shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord ) {
			float shadow = 1.0;
			shadowCoord.xyz /= shadowCoord.w;
			#ifdef USE_REVERSED_DEPTH_BUFFER
				shadowCoord.z -= shadowBias;
			#else
				shadowCoord.z += shadowBias;
			#endif
			bool inFrustum = shadowCoord.x >= 0.0 && shadowCoord.x <= 1.0 && shadowCoord.y >= 0.0 && shadowCoord.y <= 1.0;
			bool frustumTest = inFrustum && shadowCoord.z <= 1.0;
			if ( frustumTest ) {
				float depth = texture2D( shadowMap, shadowCoord.xy ).r;
				#ifdef USE_REVERSED_DEPTH_BUFFER
					shadow = step( depth, shadowCoord.z );
				#else
					shadow = step( shadowCoord.z, depth );
				#endif
			}
			return mix( 1.0, shadow, shadowIntensity );
		}
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0
	#if defined( SHADOWMAP_TYPE_PCF )
	float getPointShadow( samplerCubeShadow shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord, float shadowCameraNear, float shadowCameraFar ) {
		float shadow = 1.0;
		vec3 lightToPosition = shadowCoord.xyz;
		vec3 bd3D = normalize( lightToPosition );
		vec3 absVec = abs( lightToPosition );
		float viewSpaceZ = max( max( absVec.x, absVec.y ), absVec.z );
		if ( viewSpaceZ - shadowCameraFar <= 0.0 && viewSpaceZ - shadowCameraNear >= 0.0 ) {
			#ifdef USE_REVERSED_DEPTH_BUFFER
				float dp = ( shadowCameraNear * ( shadowCameraFar - viewSpaceZ ) ) / ( viewSpaceZ * ( shadowCameraFar - shadowCameraNear ) );
				dp -= shadowBias;
			#else
				float dp = ( shadowCameraFar * ( viewSpaceZ - shadowCameraNear ) ) / ( viewSpaceZ * ( shadowCameraFar - shadowCameraNear ) );
				dp += shadowBias;
			#endif
			float texelSize = shadowRadius / shadowMapSize.x;
			vec3 absDir = abs( bd3D );
			vec3 tangent = absDir.x > absDir.z ? vec3( 0.0, 1.0, 0.0 ) : vec3( 1.0, 0.0, 0.0 );
			tangent = normalize( cross( bd3D, tangent ) );
			vec3 bitangent = cross( bd3D, tangent );
			float phi = interleavedGradientNoise( gl_FragCoord.xy ) * PI2;
			vec2 sample0 = vogelDiskSample( 0, 5, phi );
			vec2 sample1 = vogelDiskSample( 1, 5, phi );
			vec2 sample2 = vogelDiskSample( 2, 5, phi );
			vec2 sample3 = vogelDiskSample( 3, 5, phi );
			vec2 sample4 = vogelDiskSample( 4, 5, phi );
			shadow = (
				texture( shadowMap, vec4( bd3D + ( tangent * sample0.x + bitangent * sample0.y ) * texelSize, dp ) ) +
				texture( shadowMap, vec4( bd3D + ( tangent * sample1.x + bitangent * sample1.y ) * texelSize, dp ) ) +
				texture( shadowMap, vec4( bd3D + ( tangent * sample2.x + bitangent * sample2.y ) * texelSize, dp ) ) +
				texture( shadowMap, vec4( bd3D + ( tangent * sample3.x + bitangent * sample3.y ) * texelSize, dp ) ) +
				texture( shadowMap, vec4( bd3D + ( tangent * sample4.x + bitangent * sample4.y ) * texelSize, dp ) )
			) * 0.2;
		}
		return mix( 1.0, shadow, shadowIntensity );
	}
	#elif defined( SHADOWMAP_TYPE_BASIC )
	float getPointShadow( samplerCube shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord, float shadowCameraNear, float shadowCameraFar ) {
		float shadow = 1.0;
		vec3 lightToPosition = shadowCoord.xyz;
		vec3 absVec = abs( lightToPosition );
		float viewSpaceZ = max( max( absVec.x, absVec.y ), absVec.z );
		if ( viewSpaceZ - shadowCameraFar <= 0.0 && viewSpaceZ - shadowCameraNear >= 0.0 ) {
			float dp = ( shadowCameraFar * ( viewSpaceZ - shadowCameraNear ) ) / ( viewSpaceZ * ( shadowCameraFar - shadowCameraNear ) );
			dp += shadowBias;
			vec3 bd3D = normalize( lightToPosition );
			float depth = textureCube( shadowMap, bd3D ).r;
			#ifdef USE_REVERSED_DEPTH_BUFFER
				depth = 1.0 - depth;
			#endif
			shadow = step( dp, depth );
		}
		return mix( 1.0, shadow, shadowIntensity );
	}
	#endif
	#endif
#endif`,M2=`#if NUM_SPOT_LIGHT_COORDS > 0
	uniform mat4 spotLightMatrix[ NUM_SPOT_LIGHT_COORDS ];
	varying vec4 vSpotLightCoord[ NUM_SPOT_LIGHT_COORDS ];
#endif
#ifdef USE_SHADOWMAP
	#if NUM_DIR_LIGHT_SHADOWS > 0
		uniform mat4 directionalShadowMatrix[ NUM_DIR_LIGHT_SHADOWS ];
		varying vec4 vDirectionalShadowCoord[ NUM_DIR_LIGHT_SHADOWS ];
		struct DirectionalLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform DirectionalLightShadow directionalLightShadows[ NUM_DIR_LIGHT_SHADOWS ];
	#endif
	#if NUM_SPOT_LIGHT_SHADOWS > 0
		struct SpotLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform SpotLightShadow spotLightShadows[ NUM_SPOT_LIGHT_SHADOWS ];
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0
		uniform mat4 pointShadowMatrix[ NUM_POINT_LIGHT_SHADOWS ];
		varying vec4 vPointShadowCoord[ NUM_POINT_LIGHT_SHADOWS ];
		struct PointLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
			float shadowCameraNear;
			float shadowCameraFar;
		};
		uniform PointLightShadow pointLightShadows[ NUM_POINT_LIGHT_SHADOWS ];
	#endif
#endif`,T2=`#if ( defined( USE_SHADOWMAP ) && ( NUM_DIR_LIGHT_SHADOWS > 0 || NUM_POINT_LIGHT_SHADOWS > 0 ) ) || ( NUM_SPOT_LIGHT_COORDS > 0 )
	vec3 shadowWorldNormal = inverseTransformDirection( transformedNormal, viewMatrix );
	vec4 shadowWorldPosition;
#endif
#if defined( USE_SHADOWMAP )
	#if NUM_DIR_LIGHT_SHADOWS > 0
		#pragma unroll_loop_start
		for ( int i = 0; i < NUM_DIR_LIGHT_SHADOWS; i ++ ) {
			shadowWorldPosition = worldPosition + vec4( shadowWorldNormal * directionalLightShadows[ i ].shadowNormalBias, 0 );
			vDirectionalShadowCoord[ i ] = directionalShadowMatrix[ i ] * shadowWorldPosition;
		}
		#pragma unroll_loop_end
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0
		#pragma unroll_loop_start
		for ( int i = 0; i < NUM_POINT_LIGHT_SHADOWS; i ++ ) {
			shadowWorldPosition = worldPosition + vec4( shadowWorldNormal * pointLightShadows[ i ].shadowNormalBias, 0 );
			vPointShadowCoord[ i ] = pointShadowMatrix[ i ] * shadowWorldPosition;
		}
		#pragma unroll_loop_end
	#endif
#endif
#if NUM_SPOT_LIGHT_COORDS > 0
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_SPOT_LIGHT_COORDS; i ++ ) {
		shadowWorldPosition = worldPosition;
		#if ( defined( USE_SHADOWMAP ) && UNROLLED_LOOP_INDEX < NUM_SPOT_LIGHT_SHADOWS )
			shadowWorldPosition.xyz += shadowWorldNormal * spotLightShadows[ i ].shadowNormalBias;
		#endif
		vSpotLightCoord[ i ] = spotLightMatrix[ i ] * shadowWorldPosition;
	}
	#pragma unroll_loop_end
#endif`,E2=`float getShadowMask() {
	float shadow = 1.0;
	#ifdef USE_SHADOWMAP
	#if NUM_DIR_LIGHT_SHADOWS > 0
	DirectionalLightShadow directionalLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_DIR_LIGHT_SHADOWS; i ++ ) {
		directionalLight = directionalLightShadows[ i ];
		shadow *= receiveShadow ? getShadow( directionalShadowMap[ i ], directionalLight.shadowMapSize, directionalLight.shadowIntensity, directionalLight.shadowBias, directionalLight.shadowRadius, vDirectionalShadowCoord[ i ] ) : 1.0;
	}
	#pragma unroll_loop_end
	#endif
	#if NUM_SPOT_LIGHT_SHADOWS > 0
	SpotLightShadow spotLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_SPOT_LIGHT_SHADOWS; i ++ ) {
		spotLight = spotLightShadows[ i ];
		shadow *= receiveShadow ? getShadow( spotShadowMap[ i ], spotLight.shadowMapSize, spotLight.shadowIntensity, spotLight.shadowBias, spotLight.shadowRadius, vSpotLightCoord[ i ] ) : 1.0;
	}
	#pragma unroll_loop_end
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0 && ( defined( SHADOWMAP_TYPE_PCF ) || defined( SHADOWMAP_TYPE_BASIC ) )
	PointLightShadow pointLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_POINT_LIGHT_SHADOWS; i ++ ) {
		pointLight = pointLightShadows[ i ];
		shadow *= receiveShadow ? getPointShadow( pointShadowMap[ i ], pointLight.shadowMapSize, pointLight.shadowIntensity, pointLight.shadowBias, pointLight.shadowRadius, vPointShadowCoord[ i ], pointLight.shadowCameraNear, pointLight.shadowCameraFar ) : 1.0;
	}
	#pragma unroll_loop_end
	#endif
	#endif
	return shadow;
}`,A2=`#ifdef USE_SKINNING
	mat4 boneMatX = getBoneMatrix( skinIndex.x );
	mat4 boneMatY = getBoneMatrix( skinIndex.y );
	mat4 boneMatZ = getBoneMatrix( skinIndex.z );
	mat4 boneMatW = getBoneMatrix( skinIndex.w );
#endif`,w2=`#ifdef USE_SKINNING
	uniform mat4 bindMatrix;
	uniform mat4 bindMatrixInverse;
	uniform highp sampler2D boneTexture;
	mat4 getBoneMatrix( const in float i ) {
		int size = textureSize( boneTexture, 0 ).x;
		int j = int( i ) * 4;
		int x = j % size;
		int y = j / size;
		vec4 v1 = texelFetch( boneTexture, ivec2( x, y ), 0 );
		vec4 v2 = texelFetch( boneTexture, ivec2( x + 1, y ), 0 );
		vec4 v3 = texelFetch( boneTexture, ivec2( x + 2, y ), 0 );
		vec4 v4 = texelFetch( boneTexture, ivec2( x + 3, y ), 0 );
		return mat4( v1, v2, v3, v4 );
	}
#endif`,C2=`#ifdef USE_SKINNING
	vec4 skinVertex = bindMatrix * vec4( transformed, 1.0 );
	vec4 skinned = vec4( 0.0 );
	skinned += boneMatX * skinVertex * skinWeight.x;
	skinned += boneMatY * skinVertex * skinWeight.y;
	skinned += boneMatZ * skinVertex * skinWeight.z;
	skinned += boneMatW * skinVertex * skinWeight.w;
	transformed = ( bindMatrixInverse * skinned ).xyz;
#endif`,R2=`#ifdef USE_SKINNING
	mat4 skinMatrix = mat4( 0.0 );
	skinMatrix += skinWeight.x * boneMatX;
	skinMatrix += skinWeight.y * boneMatY;
	skinMatrix += skinWeight.z * boneMatZ;
	skinMatrix += skinWeight.w * boneMatW;
	skinMatrix = bindMatrixInverse * skinMatrix * bindMatrix;
	objectNormal = vec4( skinMatrix * vec4( objectNormal, 0.0 ) ).xyz;
	#ifdef USE_TANGENT
		objectTangent = vec4( skinMatrix * vec4( objectTangent, 0.0 ) ).xyz;
	#endif
#endif`,D2=`float specularStrength;
#ifdef USE_SPECULARMAP
	vec4 texelSpecular = texture2D( specularMap, vSpecularMapUv );
	specularStrength = texelSpecular.r;
#else
	specularStrength = 1.0;
#endif`,N2=`#ifdef USE_SPECULARMAP
	uniform sampler2D specularMap;
#endif`,U2=`#if defined( TONE_MAPPING )
	gl_FragColor.rgb = toneMapping( gl_FragColor.rgb );
#endif`,L2=`#ifndef saturate
#define saturate( a ) clamp( a, 0.0, 1.0 )
#endif
uniform float toneMappingExposure;
vec3 LinearToneMapping( vec3 color ) {
	return saturate( toneMappingExposure * color );
}
vec3 ReinhardToneMapping( vec3 color ) {
	color *= toneMappingExposure;
	return saturate( color / ( vec3( 1.0 ) + color ) );
}
vec3 CineonToneMapping( vec3 color ) {
	color *= toneMappingExposure;
	color = max( vec3( 0.0 ), color - 0.004 );
	return pow( ( color * ( 6.2 * color + 0.5 ) ) / ( color * ( 6.2 * color + 1.7 ) + 0.06 ), vec3( 2.2 ) );
}
vec3 RRTAndODTFit( vec3 v ) {
	vec3 a = v * ( v + 0.0245786 ) - 0.000090537;
	vec3 b = v * ( 0.983729 * v + 0.4329510 ) + 0.238081;
	return a / b;
}
vec3 ACESFilmicToneMapping( vec3 color ) {
	const mat3 ACESInputMat = mat3(
		vec3( 0.59719, 0.07600, 0.02840 ),		vec3( 0.35458, 0.90834, 0.13383 ),
		vec3( 0.04823, 0.01566, 0.83777 )
	);
	const mat3 ACESOutputMat = mat3(
		vec3(  1.60475, -0.10208, -0.00327 ),		vec3( -0.53108,  1.10813, -0.07276 ),
		vec3( -0.07367, -0.00605,  1.07602 )
	);
	color *= toneMappingExposure / 0.6;
	color = ACESInputMat * color;
	color = RRTAndODTFit( color );
	color = ACESOutputMat * color;
	return saturate( color );
}
const mat3 LINEAR_REC2020_TO_LINEAR_SRGB = mat3(
	vec3( 1.6605, - 0.1246, - 0.0182 ),
	vec3( - 0.5876, 1.1329, - 0.1006 ),
	vec3( - 0.0728, - 0.0083, 1.1187 )
);
const mat3 LINEAR_SRGB_TO_LINEAR_REC2020 = mat3(
	vec3( 0.6274, 0.0691, 0.0164 ),
	vec3( 0.3293, 0.9195, 0.0880 ),
	vec3( 0.0433, 0.0113, 0.8956 )
);
vec3 agxDefaultContrastApprox( vec3 x ) {
	vec3 x2 = x * x;
	vec3 x4 = x2 * x2;
	return + 15.5 * x4 * x2
		- 40.14 * x4 * x
		+ 31.96 * x4
		- 6.868 * x2 * x
		+ 0.4298 * x2
		+ 0.1191 * x
		- 0.00232;
}
vec3 AgXToneMapping( vec3 color ) {
	const mat3 AgXInsetMatrix = mat3(
		vec3( 0.856627153315983, 0.137318972929847, 0.11189821299995 ),
		vec3( 0.0951212405381588, 0.761241990602591, 0.0767994186031903 ),
		vec3( 0.0482516061458583, 0.101439036467562, 0.811302368396859 )
	);
	const mat3 AgXOutsetMatrix = mat3(
		vec3( 1.1271005818144368, - 0.1413297634984383, - 0.14132976349843826 ),
		vec3( - 0.11060664309660323, 1.157823702216272, - 0.11060664309660294 ),
		vec3( - 0.016493938717834573, - 0.016493938717834257, 1.2519364065950405 )
	);
	const float AgxMinEv = - 12.47393;	const float AgxMaxEv = 4.026069;
	color *= toneMappingExposure;
	color = LINEAR_SRGB_TO_LINEAR_REC2020 * color;
	color = AgXInsetMatrix * color;
	color = max( color, 1e-10 );	color = log2( color );
	color = ( color - AgxMinEv ) / ( AgxMaxEv - AgxMinEv );
	color = clamp( color, 0.0, 1.0 );
	color = agxDefaultContrastApprox( color );
	color = AgXOutsetMatrix * color;
	color = pow( max( vec3( 0.0 ), color ), vec3( 2.2 ) );
	color = LINEAR_REC2020_TO_LINEAR_SRGB * color;
	color = clamp( color, 0.0, 1.0 );
	return color;
}
vec3 NeutralToneMapping( vec3 color ) {
	const float StartCompression = 0.8 - 0.04;
	const float Desaturation = 0.15;
	color *= toneMappingExposure;
	float x = min( color.r, min( color.g, color.b ) );
	float offset = x < 0.08 ? x - 6.25 * x * x : 0.04;
	color -= offset;
	float peak = max( color.r, max( color.g, color.b ) );
	if ( peak < StartCompression ) return color;
	float d = 1. - StartCompression;
	float newPeak = 1. - d * d / ( peak + d - StartCompression );
	color *= newPeak / peak;
	float g = 1. - 1. / ( Desaturation * ( peak - newPeak ) + 1. );
	return mix( color, vec3( newPeak ), g );
}
vec3 CustomToneMapping( vec3 color ) { return color; }`,I2=`#ifdef USE_TRANSMISSION
	material.transmission = transmission;
	material.transmissionAlpha = 1.0;
	material.thickness = thickness;
	material.attenuationDistance = attenuationDistance;
	material.attenuationColor = attenuationColor;
	#ifdef USE_TRANSMISSIONMAP
		material.transmission *= texture2D( transmissionMap, vTransmissionMapUv ).r;
	#endif
	#ifdef USE_THICKNESSMAP
		material.thickness *= texture2D( thicknessMap, vThicknessMapUv ).g;
	#endif
	vec3 pos = vWorldPosition;
	vec3 v = normalize( cameraPosition - pos );
	vec3 n = inverseTransformDirection( normal, viewMatrix );
	vec4 transmitted = getIBLVolumeRefraction(
		n, v, material.roughness, material.diffuseContribution, material.specularColorBlended, material.specularF90,
		pos, modelMatrix, viewMatrix, projectionMatrix, material.dispersion, material.ior, material.thickness,
		material.attenuationColor, material.attenuationDistance );
	material.transmissionAlpha = mix( material.transmissionAlpha, transmitted.a, material.transmission );
	totalDiffuse = mix( totalDiffuse, transmitted.rgb, material.transmission );
#endif`,O2=`#ifdef USE_TRANSMISSION
	uniform float transmission;
	uniform float thickness;
	uniform float attenuationDistance;
	uniform vec3 attenuationColor;
	#ifdef USE_TRANSMISSIONMAP
		uniform sampler2D transmissionMap;
	#endif
	#ifdef USE_THICKNESSMAP
		uniform sampler2D thicknessMap;
	#endif
	uniform vec2 transmissionSamplerSize;
	uniform sampler2D transmissionSamplerMap;
	uniform mat4 modelMatrix;
	uniform mat4 projectionMatrix;
	varying vec3 vWorldPosition;
	float w0( float a ) {
		return ( 1.0 / 6.0 ) * ( a * ( a * ( - a + 3.0 ) - 3.0 ) + 1.0 );
	}
	float w1( float a ) {
		return ( 1.0 / 6.0 ) * ( a *  a * ( 3.0 * a - 6.0 ) + 4.0 );
	}
	float w2( float a ){
		return ( 1.0 / 6.0 ) * ( a * ( a * ( - 3.0 * a + 3.0 ) + 3.0 ) + 1.0 );
	}
	float w3( float a ) {
		return ( 1.0 / 6.0 ) * ( a * a * a );
	}
	float g0( float a ) {
		return w0( a ) + w1( a );
	}
	float g1( float a ) {
		return w2( a ) + w3( a );
	}
	float h0( float a ) {
		return - 1.0 + w1( a ) / ( w0( a ) + w1( a ) );
	}
	float h1( float a ) {
		return 1.0 + w3( a ) / ( w2( a ) + w3( a ) );
	}
	vec4 bicubic( sampler2D tex, vec2 uv, vec4 texelSize, float lod ) {
		uv = uv * texelSize.zw + 0.5;
		vec2 iuv = floor( uv );
		vec2 fuv = fract( uv );
		float g0x = g0( fuv.x );
		float g1x = g1( fuv.x );
		float h0x = h0( fuv.x );
		float h1x = h1( fuv.x );
		float h0y = h0( fuv.y );
		float h1y = h1( fuv.y );
		vec2 p0 = ( vec2( iuv.x + h0x, iuv.y + h0y ) - 0.5 ) * texelSize.xy;
		vec2 p1 = ( vec2( iuv.x + h1x, iuv.y + h0y ) - 0.5 ) * texelSize.xy;
		vec2 p2 = ( vec2( iuv.x + h0x, iuv.y + h1y ) - 0.5 ) * texelSize.xy;
		vec2 p3 = ( vec2( iuv.x + h1x, iuv.y + h1y ) - 0.5 ) * texelSize.xy;
		return g0( fuv.y ) * ( g0x * textureLod( tex, p0, lod ) + g1x * textureLod( tex, p1, lod ) ) +
			g1( fuv.y ) * ( g0x * textureLod( tex, p2, lod ) + g1x * textureLod( tex, p3, lod ) );
	}
	vec4 textureBicubic( sampler2D sampler, vec2 uv, float lod ) {
		vec2 fLodSize = vec2( textureSize( sampler, int( lod ) ) );
		vec2 cLodSize = vec2( textureSize( sampler, int( lod + 1.0 ) ) );
		vec2 fLodSizeInv = 1.0 / fLodSize;
		vec2 cLodSizeInv = 1.0 / cLodSize;
		vec4 fSample = bicubic( sampler, uv, vec4( fLodSizeInv, fLodSize ), floor( lod ) );
		vec4 cSample = bicubic( sampler, uv, vec4( cLodSizeInv, cLodSize ), ceil( lod ) );
		return mix( fSample, cSample, fract( lod ) );
	}
	vec3 getVolumeTransmissionRay( const in vec3 n, const in vec3 v, const in float thickness, const in float ior, const in mat4 modelMatrix ) {
		vec3 refractionVector = refract( - v, normalize( n ), 1.0 / ior );
		vec3 modelScale;
		modelScale.x = length( vec3( modelMatrix[ 0 ].xyz ) );
		modelScale.y = length( vec3( modelMatrix[ 1 ].xyz ) );
		modelScale.z = length( vec3( modelMatrix[ 2 ].xyz ) );
		return normalize( refractionVector ) * thickness * modelScale;
	}
	float applyIorToRoughness( const in float roughness, const in float ior ) {
		return roughness * clamp( ior * 2.0 - 2.0, 0.0, 1.0 );
	}
	vec4 getTransmissionSample( const in vec2 fragCoord, const in float roughness, const in float ior ) {
		float lod = log2( transmissionSamplerSize.x ) * applyIorToRoughness( roughness, ior );
		return textureBicubic( transmissionSamplerMap, fragCoord.xy, lod );
	}
	vec3 volumeAttenuation( const in float transmissionDistance, const in vec3 attenuationColor, const in float attenuationDistance ) {
		if ( isinf( attenuationDistance ) ) {
			return vec3( 1.0 );
		} else {
			vec3 attenuationCoefficient = -log( attenuationColor ) / attenuationDistance;
			vec3 transmittance = exp( - attenuationCoefficient * transmissionDistance );			return transmittance;
		}
	}
	vec4 getIBLVolumeRefraction( const in vec3 n, const in vec3 v, const in float roughness, const in vec3 diffuseColor,
		const in vec3 specularColor, const in float specularF90, const in vec3 position, const in mat4 modelMatrix,
		const in mat4 viewMatrix, const in mat4 projMatrix, const in float dispersion, const in float ior, const in float thickness,
		const in vec3 attenuationColor, const in float attenuationDistance ) {
		vec4 transmittedLight;
		vec3 transmittance;
		#ifdef USE_DISPERSION
			float halfSpread = ( ior - 1.0 ) * 0.025 * dispersion;
			vec3 iors = vec3( ior - halfSpread, ior, ior + halfSpread );
			for ( int i = 0; i < 3; i ++ ) {
				vec3 transmissionRay = getVolumeTransmissionRay( n, v, thickness, iors[ i ], modelMatrix );
				vec3 refractedRayExit = position + transmissionRay;
				vec4 ndcPos = projMatrix * viewMatrix * vec4( refractedRayExit, 1.0 );
				vec2 refractionCoords = ndcPos.xy / ndcPos.w;
				refractionCoords += 1.0;
				refractionCoords /= 2.0;
				vec4 transmissionSample = getTransmissionSample( refractionCoords, roughness, iors[ i ] );
				transmittedLight[ i ] = transmissionSample[ i ];
				transmittedLight.a += transmissionSample.a;
				transmittance[ i ] = diffuseColor[ i ] * volumeAttenuation( length( transmissionRay ), attenuationColor, attenuationDistance )[ i ];
			}
			transmittedLight.a /= 3.0;
		#else
			vec3 transmissionRay = getVolumeTransmissionRay( n, v, thickness, ior, modelMatrix );
			vec3 refractedRayExit = position + transmissionRay;
			vec4 ndcPos = projMatrix * viewMatrix * vec4( refractedRayExit, 1.0 );
			vec2 refractionCoords = ndcPos.xy / ndcPos.w;
			refractionCoords += 1.0;
			refractionCoords /= 2.0;
			transmittedLight = getTransmissionSample( refractionCoords, roughness, ior );
			transmittance = diffuseColor * volumeAttenuation( length( transmissionRay ), attenuationColor, attenuationDistance );
		#endif
		vec3 attenuatedColor = transmittance * transmittedLight.rgb;
		vec3 F = EnvironmentBRDF( n, v, specularColor, specularF90, roughness );
		float transmittanceFactor = ( transmittance.r + transmittance.g + transmittance.b ) / 3.0;
		return vec4( ( 1.0 - F ) * attenuatedColor, 1.0 - ( 1.0 - transmittedLight.a ) * transmittanceFactor );
	}
#endif`,P2=`#if defined( USE_UV ) || defined( USE_ANISOTROPY )
	varying vec2 vUv;
#endif
#ifdef USE_MAP
	varying vec2 vMapUv;
#endif
#ifdef USE_ALPHAMAP
	varying vec2 vAlphaMapUv;
#endif
#ifdef USE_LIGHTMAP
	varying vec2 vLightMapUv;
#endif
#ifdef USE_AOMAP
	varying vec2 vAoMapUv;
#endif
#ifdef USE_BUMPMAP
	varying vec2 vBumpMapUv;
#endif
#ifdef USE_NORMALMAP
	varying vec2 vNormalMapUv;
#endif
#ifdef USE_EMISSIVEMAP
	varying vec2 vEmissiveMapUv;
#endif
#ifdef USE_METALNESSMAP
	varying vec2 vMetalnessMapUv;
#endif
#ifdef USE_ROUGHNESSMAP
	varying vec2 vRoughnessMapUv;
#endif
#ifdef USE_ANISOTROPYMAP
	varying vec2 vAnisotropyMapUv;
#endif
#ifdef USE_CLEARCOATMAP
	varying vec2 vClearcoatMapUv;
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	varying vec2 vClearcoatNormalMapUv;
#endif
#ifdef USE_CLEARCOAT_ROUGHNESSMAP
	varying vec2 vClearcoatRoughnessMapUv;
#endif
#ifdef USE_IRIDESCENCEMAP
	varying vec2 vIridescenceMapUv;
#endif
#ifdef USE_IRIDESCENCE_THICKNESSMAP
	varying vec2 vIridescenceThicknessMapUv;
#endif
#ifdef USE_SHEEN_COLORMAP
	varying vec2 vSheenColorMapUv;
#endif
#ifdef USE_SHEEN_ROUGHNESSMAP
	varying vec2 vSheenRoughnessMapUv;
#endif
#ifdef USE_SPECULARMAP
	varying vec2 vSpecularMapUv;
#endif
#ifdef USE_SPECULAR_COLORMAP
	varying vec2 vSpecularColorMapUv;
#endif
#ifdef USE_SPECULAR_INTENSITYMAP
	varying vec2 vSpecularIntensityMapUv;
#endif
#ifdef USE_TRANSMISSIONMAP
	uniform mat3 transmissionMapTransform;
	varying vec2 vTransmissionMapUv;
#endif
#ifdef USE_THICKNESSMAP
	uniform mat3 thicknessMapTransform;
	varying vec2 vThicknessMapUv;
#endif`,B2=`#if defined( USE_UV ) || defined( USE_ANISOTROPY )
	varying vec2 vUv;
#endif
#ifdef USE_MAP
	uniform mat3 mapTransform;
	varying vec2 vMapUv;
#endif
#ifdef USE_ALPHAMAP
	uniform mat3 alphaMapTransform;
	varying vec2 vAlphaMapUv;
#endif
#ifdef USE_LIGHTMAP
	uniform mat3 lightMapTransform;
	varying vec2 vLightMapUv;
#endif
#ifdef USE_AOMAP
	uniform mat3 aoMapTransform;
	varying vec2 vAoMapUv;
#endif
#ifdef USE_BUMPMAP
	uniform mat3 bumpMapTransform;
	varying vec2 vBumpMapUv;
#endif
#ifdef USE_NORMALMAP
	uniform mat3 normalMapTransform;
	varying vec2 vNormalMapUv;
#endif
#ifdef USE_DISPLACEMENTMAP
	uniform mat3 displacementMapTransform;
	varying vec2 vDisplacementMapUv;
#endif
#ifdef USE_EMISSIVEMAP
	uniform mat3 emissiveMapTransform;
	varying vec2 vEmissiveMapUv;
#endif
#ifdef USE_METALNESSMAP
	uniform mat3 metalnessMapTransform;
	varying vec2 vMetalnessMapUv;
#endif
#ifdef USE_ROUGHNESSMAP
	uniform mat3 roughnessMapTransform;
	varying vec2 vRoughnessMapUv;
#endif
#ifdef USE_ANISOTROPYMAP
	uniform mat3 anisotropyMapTransform;
	varying vec2 vAnisotropyMapUv;
#endif
#ifdef USE_CLEARCOATMAP
	uniform mat3 clearcoatMapTransform;
	varying vec2 vClearcoatMapUv;
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	uniform mat3 clearcoatNormalMapTransform;
	varying vec2 vClearcoatNormalMapUv;
#endif
#ifdef USE_CLEARCOAT_ROUGHNESSMAP
	uniform mat3 clearcoatRoughnessMapTransform;
	varying vec2 vClearcoatRoughnessMapUv;
#endif
#ifdef USE_SHEEN_COLORMAP
	uniform mat3 sheenColorMapTransform;
	varying vec2 vSheenColorMapUv;
#endif
#ifdef USE_SHEEN_ROUGHNESSMAP
	uniform mat3 sheenRoughnessMapTransform;
	varying vec2 vSheenRoughnessMapUv;
#endif
#ifdef USE_IRIDESCENCEMAP
	uniform mat3 iridescenceMapTransform;
	varying vec2 vIridescenceMapUv;
#endif
#ifdef USE_IRIDESCENCE_THICKNESSMAP
	uniform mat3 iridescenceThicknessMapTransform;
	varying vec2 vIridescenceThicknessMapUv;
#endif
#ifdef USE_SPECULARMAP
	uniform mat3 specularMapTransform;
	varying vec2 vSpecularMapUv;
#endif
#ifdef USE_SPECULAR_COLORMAP
	uniform mat3 specularColorMapTransform;
	varying vec2 vSpecularColorMapUv;
#endif
#ifdef USE_SPECULAR_INTENSITYMAP
	uniform mat3 specularIntensityMapTransform;
	varying vec2 vSpecularIntensityMapUv;
#endif
#ifdef USE_TRANSMISSIONMAP
	uniform mat3 transmissionMapTransform;
	varying vec2 vTransmissionMapUv;
#endif
#ifdef USE_THICKNESSMAP
	uniform mat3 thicknessMapTransform;
	varying vec2 vThicknessMapUv;
#endif`,z2=`#if defined( USE_UV ) || defined( USE_ANISOTROPY )
	vUv = vec3( uv, 1 ).xy;
#endif
#ifdef USE_MAP
	vMapUv = ( mapTransform * vec3( MAP_UV, 1 ) ).xy;
#endif
#ifdef USE_ALPHAMAP
	vAlphaMapUv = ( alphaMapTransform * vec3( ALPHAMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_LIGHTMAP
	vLightMapUv = ( lightMapTransform * vec3( LIGHTMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_AOMAP
	vAoMapUv = ( aoMapTransform * vec3( AOMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_BUMPMAP
	vBumpMapUv = ( bumpMapTransform * vec3( BUMPMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_NORMALMAP
	vNormalMapUv = ( normalMapTransform * vec3( NORMALMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_DISPLACEMENTMAP
	vDisplacementMapUv = ( displacementMapTransform * vec3( DISPLACEMENTMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_EMISSIVEMAP
	vEmissiveMapUv = ( emissiveMapTransform * vec3( EMISSIVEMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_METALNESSMAP
	vMetalnessMapUv = ( metalnessMapTransform * vec3( METALNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_ROUGHNESSMAP
	vRoughnessMapUv = ( roughnessMapTransform * vec3( ROUGHNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_ANISOTROPYMAP
	vAnisotropyMapUv = ( anisotropyMapTransform * vec3( ANISOTROPYMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_CLEARCOATMAP
	vClearcoatMapUv = ( clearcoatMapTransform * vec3( CLEARCOATMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	vClearcoatNormalMapUv = ( clearcoatNormalMapTransform * vec3( CLEARCOAT_NORMALMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_CLEARCOAT_ROUGHNESSMAP
	vClearcoatRoughnessMapUv = ( clearcoatRoughnessMapTransform * vec3( CLEARCOAT_ROUGHNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_IRIDESCENCEMAP
	vIridescenceMapUv = ( iridescenceMapTransform * vec3( IRIDESCENCEMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_IRIDESCENCE_THICKNESSMAP
	vIridescenceThicknessMapUv = ( iridescenceThicknessMapTransform * vec3( IRIDESCENCE_THICKNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SHEEN_COLORMAP
	vSheenColorMapUv = ( sheenColorMapTransform * vec3( SHEEN_COLORMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SHEEN_ROUGHNESSMAP
	vSheenRoughnessMapUv = ( sheenRoughnessMapTransform * vec3( SHEEN_ROUGHNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SPECULARMAP
	vSpecularMapUv = ( specularMapTransform * vec3( SPECULARMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SPECULAR_COLORMAP
	vSpecularColorMapUv = ( specularColorMapTransform * vec3( SPECULAR_COLORMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SPECULAR_INTENSITYMAP
	vSpecularIntensityMapUv = ( specularIntensityMapTransform * vec3( SPECULAR_INTENSITYMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_TRANSMISSIONMAP
	vTransmissionMapUv = ( transmissionMapTransform * vec3( TRANSMISSIONMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_THICKNESSMAP
	vThicknessMapUv = ( thicknessMapTransform * vec3( THICKNESSMAP_UV, 1 ) ).xy;
#endif`,F2=`#if defined( USE_ENVMAP ) || defined( DISTANCE ) || defined ( USE_SHADOWMAP ) || defined ( USE_TRANSMISSION ) || NUM_SPOT_LIGHT_COORDS > 0
	vec4 worldPosition = vec4( transformed, 1.0 );
	#ifdef USE_BATCHING
		worldPosition = batchingMatrix * worldPosition;
	#endif
	#ifdef USE_INSTANCING
		worldPosition = instanceMatrix * worldPosition;
	#endif
	worldPosition = modelMatrix * worldPosition;
#endif`,V2=`varying vec2 vUv;
uniform mat3 uvTransform;
void main() {
	vUv = ( uvTransform * vec3( uv, 1 ) ).xy;
	gl_Position = vec4( position.xy, 1.0, 1.0 );
}`,H2=`uniform sampler2D t2D;
uniform float backgroundIntensity;
varying vec2 vUv;
void main() {
	vec4 texColor = texture2D( t2D, vUv );
	#ifdef DECODE_VIDEO_TEXTURE
		texColor = vec4( mix( pow( texColor.rgb * 0.9478672986 + vec3( 0.0521327014 ), vec3( 2.4 ) ), texColor.rgb * 0.0773993808, vec3( lessThanEqual( texColor.rgb, vec3( 0.04045 ) ) ) ), texColor.w );
	#endif
	texColor.rgb *= backgroundIntensity;
	gl_FragColor = texColor;
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
}`,G2=`varying vec3 vWorldDirection;
#include <common>
void main() {
	vWorldDirection = transformDirection( position, modelMatrix );
	#include <begin_vertex>
	#include <project_vertex>
	gl_Position.z = gl_Position.w;
}`,k2=`#ifdef ENVMAP_TYPE_CUBE
	uniform samplerCube envMap;
#elif defined( ENVMAP_TYPE_CUBE_UV )
	uniform sampler2D envMap;
#endif
uniform float flipEnvMap;
uniform float backgroundBlurriness;
uniform float backgroundIntensity;
uniform mat3 backgroundRotation;
varying vec3 vWorldDirection;
#include <cube_uv_reflection_fragment>
void main() {
	#ifdef ENVMAP_TYPE_CUBE
		vec4 texColor = textureCube( envMap, backgroundRotation * vec3( flipEnvMap * vWorldDirection.x, vWorldDirection.yz ) );
	#elif defined( ENVMAP_TYPE_CUBE_UV )
		vec4 texColor = textureCubeUV( envMap, backgroundRotation * vWorldDirection, backgroundBlurriness );
	#else
		vec4 texColor = vec4( 0.0, 0.0, 0.0, 1.0 );
	#endif
	texColor.rgb *= backgroundIntensity;
	gl_FragColor = texColor;
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
}`,X2=`varying vec3 vWorldDirection;
#include <common>
void main() {
	vWorldDirection = transformDirection( position, modelMatrix );
	#include <begin_vertex>
	#include <project_vertex>
	gl_Position.z = gl_Position.w;
}`,W2=`uniform samplerCube tCube;
uniform float tFlip;
uniform float opacity;
varying vec3 vWorldDirection;
void main() {
	vec4 texColor = textureCube( tCube, vec3( tFlip * vWorldDirection.x, vWorldDirection.yz ) );
	gl_FragColor = texColor;
	gl_FragColor.a *= opacity;
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
}`,q2=`#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
varying vec2 vHighPrecisionZW;
void main() {
	#include <uv_vertex>
	#include <batching_vertex>
	#include <skinbase_vertex>
	#include <morphinstance_vertex>
	#ifdef USE_DISPLACEMENTMAP
		#include <beginnormal_vertex>
		#include <morphnormal_vertex>
		#include <skinnormal_vertex>
	#endif
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vHighPrecisionZW = gl_Position.zw;
}`,Y2=`#if DEPTH_PACKING == 3200
	uniform float opacity;
#endif
#include <common>
#include <packing>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
varying vec2 vHighPrecisionZW;
void main() {
	vec4 diffuseColor = vec4( 1.0 );
	#include <clipping_planes_fragment>
	#if DEPTH_PACKING == 3200
		diffuseColor.a = opacity;
	#endif
	#include <map_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <logdepthbuf_fragment>
	#ifdef USE_REVERSED_DEPTH_BUFFER
		float fragCoordZ = vHighPrecisionZW[ 0 ] / vHighPrecisionZW[ 1 ];
	#else
		float fragCoordZ = 0.5 * vHighPrecisionZW[ 0 ] / vHighPrecisionZW[ 1 ] + 0.5;
	#endif
	#if DEPTH_PACKING == 3200
		gl_FragColor = vec4( vec3( 1.0 - fragCoordZ ), opacity );
	#elif DEPTH_PACKING == 3201
		gl_FragColor = packDepthToRGBA( fragCoordZ );
	#elif DEPTH_PACKING == 3202
		gl_FragColor = vec4( packDepthToRGB( fragCoordZ ), 1.0 );
	#elif DEPTH_PACKING == 3203
		gl_FragColor = vec4( packDepthToRG( fragCoordZ ), 0.0, 1.0 );
	#endif
}`,Z2=`#define DISTANCE
varying vec3 vWorldPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <batching_vertex>
	#include <skinbase_vertex>
	#include <morphinstance_vertex>
	#ifdef USE_DISPLACEMENTMAP
		#include <beginnormal_vertex>
		#include <morphnormal_vertex>
		#include <skinnormal_vertex>
	#endif
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <worldpos_vertex>
	#include <clipping_planes_vertex>
	vWorldPosition = worldPosition.xyz;
}`,J2=`#define DISTANCE
uniform vec3 referencePosition;
uniform float nearDistance;
uniform float farDistance;
varying vec3 vWorldPosition;
#include <common>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <clipping_planes_pars_fragment>
void main () {
	vec4 diffuseColor = vec4( 1.0 );
	#include <clipping_planes_fragment>
	#include <map_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	float dist = length( vWorldPosition - referencePosition );
	dist = ( dist - nearDistance ) / ( farDistance - nearDistance );
	dist = saturate( dist );
	gl_FragColor = vec4( dist, 0.0, 0.0, 1.0 );
}`,K2=`varying vec3 vWorldDirection;
#include <common>
void main() {
	vWorldDirection = transformDirection( position, modelMatrix );
	#include <begin_vertex>
	#include <project_vertex>
}`,Q2=`uniform sampler2D tEquirect;
varying vec3 vWorldDirection;
#include <common>
void main() {
	vec3 direction = normalize( vWorldDirection );
	vec2 sampleUV = equirectUv( direction );
	gl_FragColor = texture2D( tEquirect, sampleUV );
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
}`,j2=`uniform float scale;
attribute float lineDistance;
varying float vLineDistance;
#include <common>
#include <uv_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <morphtarget_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	vLineDistance = scale * lineDistance;
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <fog_vertex>
}`,$2=`uniform vec3 diffuse;
uniform float opacity;
uniform float dashSize;
uniform float totalSize;
varying float vLineDistance;
#include <common>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <fog_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	if ( mod( vLineDistance, totalSize ) > dashSize ) {
		discard;
	}
	vec3 outgoingLight = vec3( 0.0 );
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	outgoingLight = diffuseColor.rgb;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
}`,tR=`#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <envmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#if defined ( USE_ENVMAP ) || defined ( USE_SKINNING )
		#include <beginnormal_vertex>
		#include <morphnormal_vertex>
		#include <skinbase_vertex>
		#include <skinnormal_vertex>
		#include <defaultnormal_vertex>
	#endif
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <worldpos_vertex>
	#include <envmap_vertex>
	#include <fog_vertex>
}`,eR=`uniform vec3 diffuse;
uniform float opacity;
#ifndef FLAT_SHADED
	varying vec3 vNormal;
#endif
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <envmap_common_pars_fragment>
#include <envmap_pars_fragment>
#include <fog_pars_fragment>
#include <specularmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <specularmap_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	#ifdef USE_LIGHTMAP
		vec4 lightMapTexel = texture2D( lightMap, vLightMapUv );
		reflectedLight.indirectDiffuse += lightMapTexel.rgb * lightMapIntensity * RECIPROCAL_PI;
	#else
		reflectedLight.indirectDiffuse += vec3( 1.0 );
	#endif
	#include <aomap_fragment>
	reflectedLight.indirectDiffuse *= diffuseColor.rgb;
	vec3 outgoingLight = reflectedLight.indirectDiffuse;
	#include <envmap_fragment>
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,nR=`#define LAMBERT
varying vec3 vViewPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <envmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <shadowmap_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vViewPosition = - mvPosition.xyz;
	#include <worldpos_vertex>
	#include <envmap_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
}`,iR=`#define LAMBERT
uniform vec3 diffuse;
uniform vec3 emissive;
uniform float opacity;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <emissivemap_pars_fragment>
#include <cube_uv_reflection_fragment>
#include <envmap_common_pars_fragment>
#include <envmap_pars_fragment>
#include <envmap_physical_pars_fragment>
#include <fog_pars_fragment>
#include <bsdfs>
#include <lights_pars_begin>
#include <normal_pars_fragment>
#include <lights_lambert_pars_fragment>
#include <shadowmap_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <specularmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	vec3 totalEmissiveRadiance = emissive;
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <specularmap_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	#include <emissivemap_fragment>
	#include <lights_lambert_fragment>
	#include <lights_fragment_begin>
	#include <lights_fragment_maps>
	#include <lights_fragment_end>
	#include <aomap_fragment>
	vec3 outgoingLight = reflectedLight.directDiffuse + reflectedLight.indirectDiffuse + totalEmissiveRadiance;
	#include <envmap_fragment>
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,aR=`#define MATCAP
varying vec3 vViewPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <color_pars_vertex>
#include <displacementmap_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <fog_vertex>
	vViewPosition = - mvPosition.xyz;
}`,sR=`#define MATCAP
uniform vec3 diffuse;
uniform float opacity;
uniform sampler2D matcap;
varying vec3 vViewPosition;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <fog_pars_fragment>
#include <normal_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	vec3 viewDir = normalize( vViewPosition );
	vec3 x = normalize( vec3( viewDir.z, 0.0, - viewDir.x ) );
	vec3 y = cross( viewDir, x );
	vec2 uv = vec2( dot( x, normal ), dot( y, normal ) ) * 0.495 + 0.5;
	#ifdef USE_MATCAP
		vec4 matcapColor = texture2D( matcap, uv );
	#else
		vec4 matcapColor = vec4( vec3( mix( 0.2, 0.8, uv.y ) ), 1.0 );
	#endif
	vec3 outgoingLight = diffuseColor.rgb * matcapColor.rgb;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,rR=`#define NORMAL
#if defined( FLAT_SHADED ) || defined( USE_BUMPMAP ) || defined( USE_NORMALMAP_TANGENTSPACE )
	varying vec3 vViewPosition;
#endif
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphinstance_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
#if defined( FLAT_SHADED ) || defined( USE_BUMPMAP ) || defined( USE_NORMALMAP_TANGENTSPACE )
	vViewPosition = - mvPosition.xyz;
#endif
}`,oR=`#define NORMAL
uniform float opacity;
#if defined( FLAT_SHADED ) || defined( USE_BUMPMAP ) || defined( USE_NORMALMAP_TANGENTSPACE )
	varying vec3 vViewPosition;
#endif
#include <uv_pars_fragment>
#include <normal_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( 0.0, 0.0, 0.0, opacity );
	#include <clipping_planes_fragment>
	#include <logdepthbuf_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	gl_FragColor = vec4( normalize( normal ) * 0.5 + 0.5, diffuseColor.a );
	#ifdef OPAQUE
		gl_FragColor.a = 1.0;
	#endif
}`,lR=`#define PHONG
varying vec3 vViewPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <envmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <shadowmap_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphinstance_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vViewPosition = - mvPosition.xyz;
	#include <worldpos_vertex>
	#include <envmap_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
}`,cR=`#define PHONG
uniform vec3 diffuse;
uniform vec3 emissive;
uniform vec3 specular;
uniform float shininess;
uniform float opacity;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <emissivemap_pars_fragment>
#include <cube_uv_reflection_fragment>
#include <envmap_common_pars_fragment>
#include <envmap_pars_fragment>
#include <envmap_physical_pars_fragment>
#include <fog_pars_fragment>
#include <bsdfs>
#include <lights_pars_begin>
#include <normal_pars_fragment>
#include <lights_phong_pars_fragment>
#include <shadowmap_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <specularmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	vec3 totalEmissiveRadiance = emissive;
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <specularmap_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	#include <emissivemap_fragment>
	#include <lights_phong_fragment>
	#include <lights_fragment_begin>
	#include <lights_fragment_maps>
	#include <lights_fragment_end>
	#include <aomap_fragment>
	vec3 outgoingLight = reflectedLight.directDiffuse + reflectedLight.indirectDiffuse + reflectedLight.directSpecular + reflectedLight.indirectSpecular + totalEmissiveRadiance;
	#include <envmap_fragment>
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,uR=`#define STANDARD
varying vec3 vViewPosition;
#ifdef USE_TRANSMISSION
	varying vec3 vWorldPosition;
#endif
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <shadowmap_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vViewPosition = - mvPosition.xyz;
	#include <worldpos_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
#ifdef USE_TRANSMISSION
	vWorldPosition = worldPosition.xyz;
#endif
}`,fR=`#define STANDARD
#ifdef PHYSICAL
	#define IOR
	#define USE_SPECULAR
#endif
uniform vec3 diffuse;
uniform vec3 emissive;
uniform float roughness;
uniform float metalness;
uniform float opacity;
#ifdef IOR
	uniform float ior;
#endif
#ifdef USE_SPECULAR
	uniform float specularIntensity;
	uniform vec3 specularColor;
	#ifdef USE_SPECULAR_COLORMAP
		uniform sampler2D specularColorMap;
	#endif
	#ifdef USE_SPECULAR_INTENSITYMAP
		uniform sampler2D specularIntensityMap;
	#endif
#endif
#ifdef USE_CLEARCOAT
	uniform float clearcoat;
	uniform float clearcoatRoughness;
#endif
#ifdef USE_DISPERSION
	uniform float dispersion;
#endif
#ifdef USE_IRIDESCENCE
	uniform float iridescence;
	uniform float iridescenceIOR;
	uniform float iridescenceThicknessMinimum;
	uniform float iridescenceThicknessMaximum;
#endif
#ifdef USE_SHEEN
	uniform vec3 sheenColor;
	uniform float sheenRoughness;
	#ifdef USE_SHEEN_COLORMAP
		uniform sampler2D sheenColorMap;
	#endif
	#ifdef USE_SHEEN_ROUGHNESSMAP
		uniform sampler2D sheenRoughnessMap;
	#endif
#endif
#ifdef USE_ANISOTROPY
	uniform vec2 anisotropyVector;
	#ifdef USE_ANISOTROPYMAP
		uniform sampler2D anisotropyMap;
	#endif
#endif
varying vec3 vViewPosition;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <emissivemap_pars_fragment>
#include <iridescence_fragment>
#include <cube_uv_reflection_fragment>
#include <envmap_common_pars_fragment>
#include <envmap_physical_pars_fragment>
#include <fog_pars_fragment>
#include <lights_pars_begin>
#include <normal_pars_fragment>
#include <lights_physical_pars_fragment>
#include <transmission_pars_fragment>
#include <shadowmap_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <clearcoat_pars_fragment>
#include <iridescence_pars_fragment>
#include <roughnessmap_pars_fragment>
#include <metalnessmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	vec3 totalEmissiveRadiance = emissive;
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <roughnessmap_fragment>
	#include <metalnessmap_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	#include <clearcoat_normal_fragment_begin>
	#include <clearcoat_normal_fragment_maps>
	#include <emissivemap_fragment>
	#include <lights_physical_fragment>
	#include <lights_fragment_begin>
	#include <lights_fragment_maps>
	#include <lights_fragment_end>
	#include <aomap_fragment>
	vec3 totalDiffuse = reflectedLight.directDiffuse + reflectedLight.indirectDiffuse;
	vec3 totalSpecular = reflectedLight.directSpecular + reflectedLight.indirectSpecular;
	#include <transmission_fragment>
	vec3 outgoingLight = totalDiffuse + totalSpecular + totalEmissiveRadiance;
	#ifdef USE_SHEEN
 
		outgoingLight = outgoingLight + sheenSpecularDirect + sheenSpecularIndirect;
 
 	#endif
	#ifdef USE_CLEARCOAT
		float dotNVcc = saturate( dot( geometryClearcoatNormal, geometryViewDir ) );
		vec3 Fcc = F_Schlick( material.clearcoatF0, material.clearcoatF90, dotNVcc );
		outgoingLight = outgoingLight * ( 1.0 - material.clearcoat * Fcc ) + ( clearcoatSpecularDirect + clearcoatSpecularIndirect ) * material.clearcoat;
	#endif
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,hR=`#define TOON
varying vec3 vViewPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <shadowmap_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vViewPosition = - mvPosition.xyz;
	#include <worldpos_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
}`,dR=`#define TOON
uniform vec3 diffuse;
uniform vec3 emissive;
uniform float opacity;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <emissivemap_pars_fragment>
#include <gradientmap_pars_fragment>
#include <fog_pars_fragment>
#include <bsdfs>
#include <lights_pars_begin>
#include <normal_pars_fragment>
#include <lights_toon_pars_fragment>
#include <shadowmap_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	vec3 totalEmissiveRadiance = emissive;
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	#include <emissivemap_fragment>
	#include <lights_toon_fragment>
	#include <lights_fragment_begin>
	#include <lights_fragment_maps>
	#include <lights_fragment_end>
	#include <aomap_fragment>
	vec3 outgoingLight = reflectedLight.directDiffuse + reflectedLight.indirectDiffuse + totalEmissiveRadiance;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,pR=`uniform float size;
uniform float scale;
#include <common>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <morphtarget_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
#ifdef USE_POINTS_UV
	varying vec2 vUv;
	uniform mat3 uvTransform;
#endif
void main() {
	#ifdef USE_POINTS_UV
		vUv = ( uvTransform * vec3( uv, 1 ) ).xy;
	#endif
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <project_vertex>
	gl_PointSize = size;
	#ifdef USE_SIZEATTENUATION
		bool isPerspective = isPerspectiveMatrix( projectionMatrix );
		if ( isPerspective ) gl_PointSize *= ( scale / - mvPosition.z );
	#endif
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <worldpos_vertex>
	#include <fog_vertex>
}`,mR=`uniform vec3 diffuse;
uniform float opacity;
#include <common>
#include <color_pars_fragment>
#include <map_particle_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <fog_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	vec3 outgoingLight = vec3( 0.0 );
	#include <logdepthbuf_fragment>
	#include <map_particle_fragment>
	#include <color_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	outgoingLight = diffuseColor.rgb;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
}`,gR=`#include <common>
#include <batching_pars_vertex>
#include <fog_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <shadowmap_pars_vertex>
void main() {
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphinstance_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <worldpos_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
}`,vR=`uniform vec3 color;
uniform float opacity;
#include <common>
#include <fog_pars_fragment>
#include <bsdfs>
#include <lights_pars_begin>
#include <logdepthbuf_pars_fragment>
#include <shadowmap_pars_fragment>
#include <shadowmask_pars_fragment>
void main() {
	#include <logdepthbuf_fragment>
	gl_FragColor = vec4( color, opacity * ( 1.0 - getShadowMask() ) );
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
}`,_R=`uniform float rotation;
uniform vec2 center;
#include <common>
#include <uv_pars_vertex>
#include <fog_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	vec4 mvPosition = modelViewMatrix[ 3 ];
	vec2 scale = vec2( length( modelMatrix[ 0 ].xyz ), length( modelMatrix[ 1 ].xyz ) );
	#ifndef USE_SIZEATTENUATION
		bool isPerspective = isPerspectiveMatrix( projectionMatrix );
		if ( isPerspective ) scale *= - mvPosition.z;
	#endif
	vec2 alignedPosition = ( position.xy - ( center - vec2( 0.5 ) ) ) * scale;
	vec2 rotatedPosition;
	rotatedPosition.x = cos( rotation ) * alignedPosition.x - sin( rotation ) * alignedPosition.y;
	rotatedPosition.y = sin( rotation ) * alignedPosition.x + cos( rotation ) * alignedPosition.y;
	mvPosition.xy += rotatedPosition;
	gl_Position = projectionMatrix * mvPosition;
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <fog_vertex>
}`,yR=`uniform vec3 diffuse;
uniform float opacity;
#include <common>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <fog_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	vec3 outgoingLight = vec3( 0.0 );
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	outgoingLight = diffuseColor.rgb;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
}`,Kt={alphahash_fragment:HC,alphahash_pars_fragment:GC,alphamap_fragment:kC,alphamap_pars_fragment:XC,alphatest_fragment:WC,alphatest_pars_fragment:qC,aomap_fragment:YC,aomap_pars_fragment:ZC,batching_pars_vertex:JC,batching_vertex:KC,begin_vertex:QC,beginnormal_vertex:jC,bsdfs:$C,iridescence_fragment:t3,bumpmap_pars_fragment:e3,clipping_planes_fragment:n3,clipping_planes_pars_fragment:i3,clipping_planes_pars_vertex:a3,clipping_planes_vertex:s3,color_fragment:r3,color_pars_fragment:o3,color_pars_vertex:l3,color_vertex:c3,common:u3,cube_uv_reflection_fragment:f3,defaultnormal_vertex:h3,displacementmap_pars_vertex:d3,displacementmap_vertex:p3,emissivemap_fragment:m3,emissivemap_pars_fragment:g3,colorspace_fragment:v3,colorspace_pars_fragment:_3,envmap_fragment:y3,envmap_common_pars_fragment:x3,envmap_pars_fragment:S3,envmap_pars_vertex:b3,envmap_physical_pars_fragment:L3,envmap_vertex:M3,fog_vertex:T3,fog_pars_vertex:E3,fog_fragment:A3,fog_pars_fragment:w3,gradientmap_pars_fragment:C3,lightmap_pars_fragment:R3,lights_lambert_fragment:D3,lights_lambert_pars_fragment:N3,lights_pars_begin:U3,lights_toon_fragment:I3,lights_toon_pars_fragment:O3,lights_phong_fragment:P3,lights_phong_pars_fragment:B3,lights_physical_fragment:z3,lights_physical_pars_fragment:F3,lights_fragment_begin:V3,lights_fragment_maps:H3,lights_fragment_end:G3,logdepthbuf_fragment:k3,logdepthbuf_pars_fragment:X3,logdepthbuf_pars_vertex:W3,logdepthbuf_vertex:q3,map_fragment:Y3,map_pars_fragment:Z3,map_particle_fragment:J3,map_particle_pars_fragment:K3,metalnessmap_fragment:Q3,metalnessmap_pars_fragment:j3,morphinstance_vertex:$3,morphcolor_vertex:t2,morphnormal_vertex:e2,morphtarget_pars_vertex:n2,morphtarget_vertex:i2,normal_fragment_begin:a2,normal_fragment_maps:s2,normal_pars_fragment:r2,normal_pars_vertex:o2,normal_vertex:l2,normalmap_pars_fragment:c2,clearcoat_normal_fragment_begin:u2,clearcoat_normal_fragment_maps:f2,clearcoat_pars_fragment:h2,iridescence_pars_fragment:d2,opaque_fragment:p2,packing:m2,premultiplied_alpha_fragment:g2,project_vertex:v2,dithering_fragment:_2,dithering_pars_fragment:y2,roughnessmap_fragment:x2,roughnessmap_pars_fragment:S2,shadowmap_pars_fragment:b2,shadowmap_pars_vertex:M2,shadowmap_vertex:T2,shadowmask_pars_fragment:E2,skinbase_vertex:A2,skinning_pars_vertex:w2,skinning_vertex:C2,skinnormal_vertex:R2,specularmap_fragment:D2,specularmap_pars_fragment:N2,tonemapping_fragment:U2,tonemapping_pars_fragment:L2,transmission_fragment:I2,transmission_pars_fragment:O2,uv_pars_fragment:P2,uv_pars_vertex:B2,uv_vertex:z2,worldpos_vertex:F2,background_vert:V2,background_frag:H2,backgroundCube_vert:G2,backgroundCube_frag:k2,cube_vert:X2,cube_frag:W2,depth_vert:q2,depth_frag:Y2,distance_vert:Z2,distance_frag:J2,equirect_vert:K2,equirect_frag:Q2,linedashed_vert:j2,linedashed_frag:$2,meshbasic_vert:tR,meshbasic_frag:eR,meshlambert_vert:nR,meshlambert_frag:iR,meshmatcap_vert:aR,meshmatcap_frag:sR,meshnormal_vert:rR,meshnormal_frag:oR,meshphong_vert:lR,meshphong_frag:cR,meshphysical_vert:uR,meshphysical_frag:fR,meshtoon_vert:hR,meshtoon_frag:dR,points_vert:pR,points_frag:mR,shadow_vert:gR,shadow_frag:vR,sprite_vert:_R,sprite_frag:yR},mt={common:{diffuse:{value:new Ft(16777215)},opacity:{value:1},map:{value:null},mapTransform:{value:new Zt},alphaMap:{value:null},alphaMapTransform:{value:new Zt},alphaTest:{value:0}},specularmap:{specularMap:{value:null},specularMapTransform:{value:new Zt}},envmap:{envMap:{value:null},envMapRotation:{value:new Zt},flipEnvMap:{value:-1},reflectivity:{value:1},ior:{value:1.5},refractionRatio:{value:.98},dfgLUT:{value:null}},aomap:{aoMap:{value:null},aoMapIntensity:{value:1},aoMapTransform:{value:new Zt}},lightmap:{lightMap:{value:null},lightMapIntensity:{value:1},lightMapTransform:{value:new Zt}},bumpmap:{bumpMap:{value:null},bumpMapTransform:{value:new Zt},bumpScale:{value:1}},normalmap:{normalMap:{value:null},normalMapTransform:{value:new Zt},normalScale:{value:new te(1,1)}},displacementmap:{displacementMap:{value:null},displacementMapTransform:{value:new Zt},displacementScale:{value:1},displacementBias:{value:0}},emissivemap:{emissiveMap:{value:null},emissiveMapTransform:{value:new Zt}},metalnessmap:{metalnessMap:{value:null},metalnessMapTransform:{value:new Zt}},roughnessmap:{roughnessMap:{value:null},roughnessMapTransform:{value:new Zt}},gradientmap:{gradientMap:{value:null}},fog:{fogDensity:{value:25e-5},fogNear:{value:1},fogFar:{value:2e3},fogColor:{value:new Ft(16777215)}},lights:{ambientLightColor:{value:[]},lightProbe:{value:[]},directionalLights:{value:[],properties:{direction:{},color:{}}},directionalLightShadows:{value:[],properties:{shadowIntensity:1,shadowBias:{},shadowNormalBias:{},shadowRadius:{},shadowMapSize:{}}},directionalShadowMatrix:{value:[]},spotLights:{value:[],properties:{color:{},position:{},direction:{},distance:{},coneCos:{},penumbraCos:{},decay:{}}},spotLightShadows:{value:[],properties:{shadowIntensity:1,shadowBias:{},shadowNormalBias:{},shadowRadius:{},shadowMapSize:{}}},spotLightMap:{value:[]},spotLightMatrix:{value:[]},pointLights:{value:[],properties:{color:{},position:{},decay:{},distance:{}}},pointLightShadows:{value:[],properties:{shadowIntensity:1,shadowBias:{},shadowNormalBias:{},shadowRadius:{},shadowMapSize:{},shadowCameraNear:{},shadowCameraFar:{}}},pointShadowMatrix:{value:[]},hemisphereLights:{value:[],properties:{direction:{},skyColor:{},groundColor:{}}},rectAreaLights:{value:[],properties:{color:{},position:{},width:{},height:{}}},ltc_1:{value:null},ltc_2:{value:null}},points:{diffuse:{value:new Ft(16777215)},opacity:{value:1},size:{value:1},scale:{value:1},map:{value:null},alphaMap:{value:null},alphaMapTransform:{value:new Zt},alphaTest:{value:0},uvTransform:{value:new Zt}},sprite:{diffuse:{value:new Ft(16777215)},opacity:{value:1},center:{value:new te(.5,.5)},rotation:{value:0},map:{value:null},mapTransform:{value:new Zt},alphaMap:{value:null},alphaMapTransform:{value:new Zt},alphaTest:{value:0}}},xa={basic:{uniforms:zn([mt.common,mt.specularmap,mt.envmap,mt.aomap,mt.lightmap,mt.fog]),vertexShader:Kt.meshbasic_vert,fragmentShader:Kt.meshbasic_frag},lambert:{uniforms:zn([mt.common,mt.specularmap,mt.envmap,mt.aomap,mt.lightmap,mt.emissivemap,mt.bumpmap,mt.normalmap,mt.displacementmap,mt.fog,mt.lights,{emissive:{value:new Ft(0)},envMapIntensity:{value:1}}]),vertexShader:Kt.meshlambert_vert,fragmentShader:Kt.meshlambert_frag},phong:{uniforms:zn([mt.common,mt.specularmap,mt.envmap,mt.aomap,mt.lightmap,mt.emissivemap,mt.bumpmap,mt.normalmap,mt.displacementmap,mt.fog,mt.lights,{emissive:{value:new Ft(0)},specular:{value:new Ft(1118481)},shininess:{value:30},envMapIntensity:{value:1}}]),vertexShader:Kt.meshphong_vert,fragmentShader:Kt.meshphong_frag},standard:{uniforms:zn([mt.common,mt.envmap,mt.aomap,mt.lightmap,mt.emissivemap,mt.bumpmap,mt.normalmap,mt.displacementmap,mt.roughnessmap,mt.metalnessmap,mt.fog,mt.lights,{emissive:{value:new Ft(0)},roughness:{value:1},metalness:{value:0},envMapIntensity:{value:1}}]),vertexShader:Kt.meshphysical_vert,fragmentShader:Kt.meshphysical_frag},toon:{uniforms:zn([mt.common,mt.aomap,mt.lightmap,mt.emissivemap,mt.bumpmap,mt.normalmap,mt.displacementmap,mt.gradientmap,mt.fog,mt.lights,{emissive:{value:new Ft(0)}}]),vertexShader:Kt.meshtoon_vert,fragmentShader:Kt.meshtoon_frag},matcap:{uniforms:zn([mt.common,mt.bumpmap,mt.normalmap,mt.displacementmap,mt.fog,{matcap:{value:null}}]),vertexShader:Kt.meshmatcap_vert,fragmentShader:Kt.meshmatcap_frag},points:{uniforms:zn([mt.points,mt.fog]),vertexShader:Kt.points_vert,fragmentShader:Kt.points_frag},dashed:{uniforms:zn([mt.common,mt.fog,{scale:{value:1},dashSize:{value:1},totalSize:{value:2}}]),vertexShader:Kt.linedashed_vert,fragmentShader:Kt.linedashed_frag},depth:{uniforms:zn([mt.common,mt.displacementmap]),vertexShader:Kt.depth_vert,fragmentShader:Kt.depth_frag},normal:{uniforms:zn([mt.common,mt.bumpmap,mt.normalmap,mt.displacementmap,{opacity:{value:1}}]),vertexShader:Kt.meshnormal_vert,fragmentShader:Kt.meshnormal_frag},sprite:{uniforms:zn([mt.sprite,mt.fog]),vertexShader:Kt.sprite_vert,fragmentShader:Kt.sprite_frag},background:{uniforms:{uvTransform:{value:new Zt},t2D:{value:null},backgroundIntensity:{value:1}},vertexShader:Kt.background_vert,fragmentShader:Kt.background_frag},backgroundCube:{uniforms:{envMap:{value:null},flipEnvMap:{value:-1},backgroundBlurriness:{value:0},backgroundIntensity:{value:1},backgroundRotation:{value:new Zt}},vertexShader:Kt.backgroundCube_vert,fragmentShader:Kt.backgroundCube_frag},cube:{uniforms:{tCube:{value:null},tFlip:{value:-1},opacity:{value:1}},vertexShader:Kt.cube_vert,fragmentShader:Kt.cube_frag},equirect:{uniforms:{tEquirect:{value:null}},vertexShader:Kt.equirect_vert,fragmentShader:Kt.equirect_frag},distance:{uniforms:zn([mt.common,mt.displacementmap,{referencePosition:{value:new B},nearDistance:{value:1},farDistance:{value:1e3}}]),vertexShader:Kt.distance_vert,fragmentShader:Kt.distance_frag},shadow:{uniforms:zn([mt.lights,mt.fog,{color:{value:new Ft(0)},opacity:{value:1}}]),vertexShader:Kt.shadow_vert,fragmentShader:Kt.shadow_frag}};xa.physical={uniforms:zn([xa.standard.uniforms,{clearcoat:{value:0},clearcoatMap:{value:null},clearcoatMapTransform:{value:new Zt},clearcoatNormalMap:{value:null},clearcoatNormalMapTransform:{value:new Zt},clearcoatNormalScale:{value:new te(1,1)},clearcoatRoughness:{value:0},clearcoatRoughnessMap:{value:null},clearcoatRoughnessMapTransform:{value:new Zt},dispersion:{value:0},iridescence:{value:0},iridescenceMap:{value:null},iridescenceMapTransform:{value:new Zt},iridescenceIOR:{value:1.3},iridescenceThicknessMinimum:{value:100},iridescenceThicknessMaximum:{value:400},iridescenceThicknessMap:{value:null},iridescenceThicknessMapTransform:{value:new Zt},sheen:{value:0},sheenColor:{value:new Ft(0)},sheenColorMap:{value:null},sheenColorMapTransform:{value:new Zt},sheenRoughness:{value:1},sheenRoughnessMap:{value:null},sheenRoughnessMapTransform:{value:new Zt},transmission:{value:0},transmissionMap:{value:null},transmissionMapTransform:{value:new Zt},transmissionSamplerSize:{value:new te},transmissionSamplerMap:{value:null},thickness:{value:0},thicknessMap:{value:null},thicknessMapTransform:{value:new Zt},attenuationDistance:{value:0},attenuationColor:{value:new Ft(0)},specularColor:{value:new Ft(1,1,1)},specularColorMap:{value:null},specularColorMapTransform:{value:new Zt},specularIntensity:{value:1},specularIntensityMap:{value:null},specularIntensityMapTransform:{value:new Zt},anisotropyVector:{value:new te},anisotropyMap:{value:null},anisotropyMapTransform:{value:new Zt}}]),vertexShader:Kt.meshphysical_vert,fragmentShader:Kt.meshphysical_frag};var lp={r:0,b:0,g:0},Or=new Xi,xR=new Ie;function SR(e,t,n,i,a,s){let r=new Ft(0),o=a===!0?0:1,l,c,f=null,p=0,u=null;function d(m){let g=m.isScene===!0?m.background:null;if(g&&g.isTexture){let x=m.backgroundBlurriness>0;g=t.get(g,x)}return g}function _(m){let g=!1,x=d(m);x===null?v(r,o):x&&x.isColor&&(v(x,1),g=!0);let E=e.xr.getEnvironmentBlendMode();E==="additive"?n.buffers.color.setClear(0,0,0,1,s):E==="alpha-blend"&&n.buffers.color.setClear(0,0,0,0,s),(e.autoClear||g)&&(n.buffers.depth.setTest(!0),n.buffers.depth.setMask(!0),n.buffers.color.setMask(!0),e.clear(e.autoClearColor,e.autoClearDepth,e.autoClearStencil))}function b(m,g){let x=d(g);x&&(x.isCubeTexture||x.mapping===du)?(c===void 0&&(c=new mn(new Vs(1,1,1),new Yn({name:"BackgroundCubeMaterial",uniforms:Ir(xa.backgroundCube.uniforms),vertexShader:xa.backgroundCube.vertexShader,fragmentShader:xa.backgroundCube.fragmentShader,side:wn,depthTest:!1,depthWrite:!1,fog:!1,allowOverride:!1})),c.geometry.deleteAttribute("normal"),c.geometry.deleteAttribute("uv"),c.onBeforeRender=function(E,A,w){this.matrixWorld.copyPosition(w.matrixWorld)},Object.defineProperty(c.material,"envMap",{get:function(){return this.uniforms.envMap.value}}),i.update(c)),Or.copy(g.backgroundRotation),Or.x*=-1,Or.y*=-1,Or.z*=-1,x.isCubeTexture&&x.isRenderTargetTexture===!1&&(Or.y*=-1,Or.z*=-1),c.material.uniforms.envMap.value=x,c.material.uniforms.flipEnvMap.value=x.isCubeTexture&&x.isRenderTargetTexture===!1?-1:1,c.material.uniforms.backgroundBlurriness.value=g.backgroundBlurriness,c.material.uniforms.backgroundIntensity.value=g.backgroundIntensity,c.material.uniforms.backgroundRotation.value.setFromMatrix4(xR.makeRotationFromEuler(Or)),c.material.toneMapped=se.getTransfer(x.colorSpace)!==me,(f!==x||p!==x.version||u!==e.toneMapping)&&(c.material.needsUpdate=!0,f=x,p=x.version,u=e.toneMapping),c.layers.enableAll(),m.unshift(c,c.geometry,c.material,0,0,null)):x&&x.isTexture&&(l===void 0&&(l=new mn(new au(2,2),new Yn({name:"BackgroundMaterial",uniforms:Ir(xa.background.uniforms),vertexShader:xa.background.vertexShader,fragmentShader:xa.background.fragmentShader,side:Wa,depthTest:!1,depthWrite:!1,fog:!1,allowOverride:!1})),l.geometry.deleteAttribute("normal"),Object.defineProperty(l.material,"map",{get:function(){return this.uniforms.t2D.value}}),i.update(l)),l.material.uniforms.t2D.value=x,l.material.uniforms.backgroundIntensity.value=g.backgroundIntensity,l.material.toneMapped=se.getTransfer(x.colorSpace)!==me,x.matrixAutoUpdate===!0&&x.updateMatrix(),l.material.uniforms.uvTransform.value.copy(x.matrix),(f!==x||p!==x.version||u!==e.toneMapping)&&(l.material.needsUpdate=!0,f=x,p=x.version,u=e.toneMapping),l.layers.enableAll(),m.unshift(l,l.geometry,l.material,0,0,null))}function v(m,g){m.getRGB(lp,bv(e)),n.buffers.color.setClear(lp.r,lp.g,lp.b,g,s)}function h(){c!==void 0&&(c.geometry.dispose(),c.material.dispose(),c=void 0),l!==void 0&&(l.geometry.dispose(),l.material.dispose(),l=void 0)}return{getClearColor:function(){return r},setClearColor:function(m,g=1){r.set(m),o=g,v(r,o)},getClearAlpha:function(){return o},setClearAlpha:function(m){o=m,v(r,o)},render:_,addToRenderList:b,dispose:h}}function bR(e,t){let n=e.getParameter(e.MAX_VERTEX_ATTRIBS),i={},a=u(null),s=a,r=!1;function o(R,N,P,z,F){let V=!1,X=p(R,z,P,N);s!==X&&(s=X,c(s.object)),V=d(R,z,P,F),V&&_(R,z,P,F),F!==null&&t.update(F,e.ELEMENT_ARRAY_BUFFER),(V||r)&&(r=!1,x(R,N,P,z),F!==null&&e.bindBuffer(e.ELEMENT_ARRAY_BUFFER,t.get(F).buffer))}function l(){return e.createVertexArray()}function c(R){return e.bindVertexArray(R)}function f(R){return e.deleteVertexArray(R)}function p(R,N,P,z){let F=z.wireframe===!0,V=i[N.id];V===void 0&&(V={},i[N.id]=V);let X=R.isInstancedMesh===!0?R.id:0,at=V[X];at===void 0&&(at={},V[X]=at);let et=at[P.id];et===void 0&&(et={},at[P.id]=et);let dt=et[F];return dt===void 0&&(dt=u(l()),et[F]=dt),dt}function u(R){let N=[],P=[],z=[];for(let F=0;F<n;F++)N[F]=0,P[F]=0,z[F]=0;return{geometry:null,program:null,wireframe:!1,newAttributes:N,enabledAttributes:P,attributeDivisors:z,object:R,attributes:{},index:null}}function d(R,N,P,z){let F=s.attributes,V=N.attributes,X=0,at=P.getAttributes();for(let et in at)if(at[et].location>=0){let St=F[et],ht=V[et];if(ht===void 0&&(et==="instanceMatrix"&&R.instanceMatrix&&(ht=R.instanceMatrix),et==="instanceColor"&&R.instanceColor&&(ht=R.instanceColor)),St===void 0||St.attribute!==ht||ht&&St.data!==ht.data)return!0;X++}return s.attributesNum!==X||s.index!==z}function _(R,N,P,z){let F={},V=N.attributes,X=0,at=P.getAttributes();for(let et in at)if(at[et].location>=0){let St=V[et];St===void 0&&(et==="instanceMatrix"&&R.instanceMatrix&&(St=R.instanceMatrix),et==="instanceColor"&&R.instanceColor&&(St=R.instanceColor));let ht={};ht.attribute=St,St&&St.data&&(ht.data=St.data),F[et]=ht,X++}s.attributes=F,s.attributesNum=X,s.index=z}function b(){let R=s.newAttributes;for(let N=0,P=R.length;N<P;N++)R[N]=0}function v(R){h(R,0)}function h(R,N){let P=s.newAttributes,z=s.enabledAttributes,F=s.attributeDivisors;P[R]=1,z[R]===0&&(e.enableVertexAttribArray(R),z[R]=1),F[R]!==N&&(e.vertexAttribDivisor(R,N),F[R]=N)}function m(){let R=s.newAttributes,N=s.enabledAttributes;for(let P=0,z=N.length;P<z;P++)N[P]!==R[P]&&(e.disableVertexAttribArray(P),N[P]=0)}function g(R,N,P,z,F,V,X){X===!0?e.vertexAttribIPointer(R,N,P,F,V):e.vertexAttribPointer(R,N,P,z,F,V)}function x(R,N,P,z){b();let F=z.attributes,V=P.getAttributes(),X=N.defaultAttributeValues;for(let at in V){let et=V[at];if(et.location>=0){let dt=F[at];if(dt===void 0&&(at==="instanceMatrix"&&R.instanceMatrix&&(dt=R.instanceMatrix),at==="instanceColor"&&R.instanceColor&&(dt=R.instanceColor)),dt!==void 0){let St=dt.normalized,ht=dt.itemSize,Vt=t.get(dt);if(Vt===void 0)continue;let Ce=Vt.buffer,xe=Vt.type,K=Vt.bytesPerElement,j=xe===e.INT||xe===e.UNSIGNED_INT||dt.gpuType===Sd;if(dt.isInterleavedBufferAttribute){let ct=dt.data,It=ct.stride,Lt=dt.offset;if(ct.isInstancedInterleavedBuffer){for(let Pt=0;Pt<et.locationSize;Pt++)h(et.location+Pt,ct.meshPerAttribute);R.isInstancedMesh!==!0&&z._maxInstanceCount===void 0&&(z._maxInstanceCount=ct.meshPerAttribute*ct.count)}else for(let Pt=0;Pt<et.locationSize;Pt++)v(et.location+Pt);e.bindBuffer(e.ARRAY_BUFFER,Ce);for(let Pt=0;Pt<et.locationSize;Pt++)g(et.location+Pt,ht/et.locationSize,xe,St,It*K,(Lt+ht/et.locationSize*Pt)*K,j)}else{if(dt.isInstancedBufferAttribute){for(let ct=0;ct<et.locationSize;ct++)h(et.location+ct,dt.meshPerAttribute);R.isInstancedMesh!==!0&&z._maxInstanceCount===void 0&&(z._maxInstanceCount=dt.meshPerAttribute*dt.count)}else for(let ct=0;ct<et.locationSize;ct++)v(et.location+ct);e.bindBuffer(e.ARRAY_BUFFER,Ce);for(let ct=0;ct<et.locationSize;ct++)g(et.location+ct,ht/et.locationSize,xe,St,ht*K,ht/et.locationSize*ct*K,j)}}else if(X!==void 0){let St=X[at];if(St!==void 0)switch(St.length){case 2:e.vertexAttrib2fv(et.location,St);break;case 3:e.vertexAttrib3fv(et.location,St);break;case 4:e.vertexAttrib4fv(et.location,St);break;default:e.vertexAttrib1fv(et.location,St)}}}}m()}function E(){T();for(let R in i){let N=i[R];for(let P in N){let z=N[P];for(let F in z){let V=z[F];for(let X in V)f(V[X].object),delete V[X];delete z[F]}}delete i[R]}}function A(R){if(i[R.id]===void 0)return;let N=i[R.id];for(let P in N){let z=N[P];for(let F in z){let V=z[F];for(let X in V)f(V[X].object),delete V[X];delete z[F]}}delete i[R.id]}function w(R){for(let N in i){let P=i[N];for(let z in P){let F=P[z];if(F[R.id]===void 0)continue;let V=F[R.id];for(let X in V)f(V[X].object),delete V[X];delete F[R.id]}}}function y(R){for(let N in i){let P=i[N],z=R.isInstancedMesh===!0?R.id:0,F=P[z];if(F!==void 0){for(let V in F){let X=F[V];for(let at in X)f(X[at].object),delete X[at];delete F[V]}delete P[z],Object.keys(P).length===0&&delete i[N]}}}function T(){L(),r=!0,s!==a&&(s=a,c(s.object))}function L(){a.geometry=null,a.program=null,a.wireframe=!1}return{setup:o,reset:T,resetDefaultState:L,dispose:E,releaseStatesOfGeometry:A,releaseStatesOfObject:y,releaseStatesOfProgram:w,initAttributes:b,enableAttribute:v,disableUnusedAttributes:m}}function MR(e,t,n){let i;function a(c){i=c}function s(c,f){e.drawArrays(i,c,f),n.update(f,i,1)}function r(c,f,p){p!==0&&(e.drawArraysInstanced(i,c,f,p),n.update(f,i,p))}function o(c,f,p){if(p===0)return;t.get("WEBGL_multi_draw").multiDrawArraysWEBGL(i,c,0,f,0,p);let d=0;for(let _=0;_<p;_++)d+=f[_];n.update(d,i,1)}function l(c,f,p,u){if(p===0)return;let d=t.get("WEBGL_multi_draw");if(d===null)for(let _=0;_<c.length;_++)r(c[_],f[_],u[_]);else{d.multiDrawArraysInstancedWEBGL(i,c,0,f,0,u,0,p);let _=0;for(let b=0;b<p;b++)_+=f[b]*u[b];n.update(_,i,1)}}this.setMode=a,this.render=s,this.renderInstances=r,this.renderMultiDraw=o,this.renderMultiDrawInstances=l}function TR(e,t,n,i){let a;function s(){if(a!==void 0)return a;if(t.has("EXT_texture_filter_anisotropic")===!0){let w=t.get("EXT_texture_filter_anisotropic");a=e.getParameter(w.MAX_TEXTURE_MAX_ANISOTROPY_EXT)}else a=0;return a}function r(w){return!(w!==Ui&&i.convert(w)!==e.getParameter(e.IMPLEMENTATION_COLOR_READ_FORMAT))}function o(w){let y=w===_a&&(t.has("EXT_color_buffer_half_float")||t.has("EXT_color_buffer_float"));return!(w!==ei&&i.convert(w)!==e.getParameter(e.IMPLEMENTATION_COLOR_READ_TYPE)&&w!==Ni&&!y)}function l(w){if(w==="highp"){if(e.getShaderPrecisionFormat(e.VERTEX_SHADER,e.HIGH_FLOAT).precision>0&&e.getShaderPrecisionFormat(e.FRAGMENT_SHADER,e.HIGH_FLOAT).precision>0)return"highp";w="mediump"}return w==="mediump"&&e.getShaderPrecisionFormat(e.VERTEX_SHADER,e.MEDIUM_FLOAT).precision>0&&e.getShaderPrecisionFormat(e.FRAGMENT_SHADER,e.MEDIUM_FLOAT).precision>0?"mediump":"lowp"}let c=n.precision!==void 0?n.precision:"highp",f=l(c);f!==c&&(Ot("WebGLRenderer:",c,"not supported, using",f,"instead."),c=f);let p=n.logarithmicDepthBuffer===!0,u=n.reversedDepthBuffer===!0&&t.has("EXT_clip_control"),d=e.getParameter(e.MAX_TEXTURE_IMAGE_UNITS),_=e.getParameter(e.MAX_VERTEX_TEXTURE_IMAGE_UNITS),b=e.getParameter(e.MAX_TEXTURE_SIZE),v=e.getParameter(e.MAX_CUBE_MAP_TEXTURE_SIZE),h=e.getParameter(e.MAX_VERTEX_ATTRIBS),m=e.getParameter(e.MAX_VERTEX_UNIFORM_VECTORS),g=e.getParameter(e.MAX_VARYING_VECTORS),x=e.getParameter(e.MAX_FRAGMENT_UNIFORM_VECTORS),E=e.getParameter(e.MAX_SAMPLES),A=e.getParameter(e.SAMPLES);return{isWebGL2:!0,getMaxAnisotropy:s,getMaxPrecision:l,textureFormatReadable:r,textureTypeReadable:o,precision:c,logarithmicDepthBuffer:p,reversedDepthBuffer:u,maxTextures:d,maxVertexTextures:_,maxTextureSize:b,maxCubemapSize:v,maxAttributes:h,maxVertexUniforms:m,maxVaryings:g,maxFragmentUniforms:x,maxSamples:E,samples:A}}function ER(e){let t=this,n=null,i=0,a=!1,s=!1,r=new fa,o=new Zt,l={value:null,needsUpdate:!1};this.uniform=l,this.numPlanes=0,this.numIntersection=0,this.init=function(p,u){let d=p.length!==0||u||i!==0||a;return a=u,i=p.length,d},this.beginShadows=function(){s=!0,f(null)},this.endShadows=function(){s=!1},this.setGlobalState=function(p,u){n=f(p,u,0)},this.setState=function(p,u,d){let _=p.clippingPlanes,b=p.clipIntersection,v=p.clipShadows,h=e.get(p);if(!a||_===null||_.length===0||s&&!v)s?f(null):c();else{let m=s?0:i,g=m*4,x=h.clippingState||null;l.value=x,x=f(_,u,g,d);for(let E=0;E!==g;++E)x[E]=n[E];h.clippingState=x,this.numIntersection=b?this.numPlanes:0,this.numPlanes+=m}};function c(){l.value!==n&&(l.value=n,l.needsUpdate=i>0),t.numPlanes=i,t.numIntersection=0}function f(p,u,d,_){let b=p!==null?p.length:0,v=null;if(b!==0){if(v=l.value,_!==!0||v===null){let h=d+b*4,m=u.matrixWorldInverse;o.getNormalMatrix(m),(v===null||v.length<h)&&(v=new Float32Array(h));for(let g=0,x=d;g!==b;++g,x+=4)r.copy(p[g]).applyMatrix4(m,o),r.normal.toArray(v,x),v[x+3]=r.constant}l.value=v,l.needsUpdate=!0}return t.numPlanes=b,t.numIntersection=0,v}}var Ys=4,E1=[.125,.215,.35,.446,.526,.582],Br=20,AR=256,yu=new hl,A1=new Ft,Av=null,wv=0,Cv=0,Rv=!1,wR=new B,up=class{constructor(t){this._renderer=t,this._pingPongRenderTarget=null,this._lodMax=0,this._cubeSize=0,this._sizeLods=[],this._sigmas=[],this._lodMeshes=[],this._backgroundBox=null,this._cubemapMaterial=null,this._equirectMaterial=null,this._blurMaterial=null,this._ggxMaterial=null}fromScene(t,n=0,i=.1,a=100,s={}){let{size:r=256,position:o=wR}=s;Av=this._renderer.getRenderTarget(),wv=this._renderer.getActiveCubeFace(),Cv=this._renderer.getActiveMipmapLevel(),Rv=this._renderer.xr.enabled,this._renderer.xr.enabled=!1,this._setSize(r);let l=this._allocateTargets();return l.depthBuffer=!0,this._sceneToCubeUV(t,i,a,l,o),n>0&&this._blur(l,0,0,n),this._applyPMREM(l),this._cleanup(l),l}fromEquirectangular(t,n=null){return this._fromTexture(t,n)}fromCubemap(t,n=null){return this._fromTexture(t,n)}compileCubemapShader(){this._cubemapMaterial===null&&(this._cubemapMaterial=R1(),this._compileMaterial(this._cubemapMaterial))}compileEquirectangularShader(){this._equirectMaterial===null&&(this._equirectMaterial=C1(),this._compileMaterial(this._equirectMaterial))}dispose(){this._dispose(),this._cubemapMaterial!==null&&this._cubemapMaterial.dispose(),this._equirectMaterial!==null&&this._equirectMaterial.dispose(),this._backgroundBox!==null&&(this._backgroundBox.geometry.dispose(),this._backgroundBox.material.dispose())}_setSize(t){this._lodMax=Math.floor(Math.log2(t)),this._cubeSize=Math.pow(2,this._lodMax)}_dispose(){this._blurMaterial!==null&&this._blurMaterial.dispose(),this._ggxMaterial!==null&&this._ggxMaterial.dispose(),this._pingPongRenderTarget!==null&&this._pingPongRenderTarget.dispose();for(let t=0;t<this._lodMeshes.length;t++)this._lodMeshes[t].geometry.dispose()}_cleanup(t){this._renderer.setRenderTarget(Av,wv,Cv),this._renderer.xr.enabled=Rv,t.scissorTest=!1,_l(t,0,0,t.width,t.height)}_fromTexture(t,n){t.mapping===Xs||t.mapping===Ur?this._setSize(t.image.length===0?16:t.image[0].width||t.image[0].image.width):this._setSize(t.image.width/4),Av=this._renderer.getRenderTarget(),wv=this._renderer.getActiveCubeFace(),Cv=this._renderer.getActiveMipmapLevel(),Rv=this._renderer.xr.enabled,this._renderer.xr.enabled=!1;let i=n||this._allocateTargets();return this._textureToCubeUV(t,i),this._applyPMREM(i),this._cleanup(i),i}_allocateTargets(){let t=3*Math.max(this._cubeSize,112),n=4*this._cubeSize,i={magFilter:En,minFilter:En,generateMipmaps:!1,type:_a,format:Ui,colorSpace:Rr,depthBuffer:!1},a=w1(t,n,i);if(this._pingPongRenderTarget===null||this._pingPongRenderTarget.width!==t||this._pingPongRenderTarget.height!==n){this._pingPongRenderTarget!==null&&this._dispose(),this._pingPongRenderTarget=w1(t,n,i);let{_lodMax:s}=this;({lodMeshes:this._lodMeshes,sizeLods:this._sizeLods,sigmas:this._sigmas}=CR(s)),this._blurMaterial=DR(s,t,n),this._ggxMaterial=RR(s,t,n)}return a}_compileMaterial(t){let n=new mn(new An,t);this._renderer.compile(n,yu)}_sceneToCubeUV(t,n,i,a,s){let l=new Pn(90,1,n,i),c=[1,-1,1,1,1,1],f=[1,1,1,-1,-1,-1],p=this._renderer,u=p.autoClear,d=p.toneMapping;p.getClearColor(A1),p.toneMapping=Wi,p.autoClear=!1,p.state.buffers.depth.getReversed()&&(p.setRenderTarget(a),p.clearDepth(),p.setRenderTarget(null)),this._backgroundBox===null&&(this._backgroundBox=new mn(new Vs,new Dr({name:"PMREM.Background",side:wn,depthWrite:!1,depthTest:!1})));let b=this._backgroundBox,v=b.material,h=!1,m=t.background;m?m.isColor&&(v.color.copy(m),t.background=null,h=!0):(v.color.copy(A1),h=!0);for(let g=0;g<6;g++){let x=g%3;x===0?(l.up.set(0,c[g],0),l.position.set(s.x,s.y,s.z),l.lookAt(s.x+f[g],s.y,s.z)):x===1?(l.up.set(0,0,c[g]),l.position.set(s.x,s.y,s.z),l.lookAt(s.x,s.y+f[g],s.z)):(l.up.set(0,c[g],0),l.position.set(s.x,s.y,s.z),l.lookAt(s.x,s.y,s.z+f[g]));let E=this._cubeSize;_l(a,x*E,g>2?E:0,E,E),p.setRenderTarget(a),h&&p.render(b,l),p.render(t,l)}p.toneMapping=d,p.autoClear=u,t.background=m}_textureToCubeUV(t,n){let i=this._renderer,a=t.mapping===Xs||t.mapping===Ur;a?(this._cubemapMaterial===null&&(this._cubemapMaterial=R1()),this._cubemapMaterial.uniforms.flipEnvMap.value=t.isRenderTargetTexture===!1?-1:1):this._equirectMaterial===null&&(this._equirectMaterial=C1());let s=a?this._cubemapMaterial:this._equirectMaterial,r=this._lodMeshes[0];r.material=s;let o=s.uniforms;o.envMap.value=t;let l=this._cubeSize;_l(n,0,0,3*l,2*l),i.setRenderTarget(n),i.render(r,yu)}_applyPMREM(t){let n=this._renderer,i=n.autoClear;n.autoClear=!1;let a=this._lodMeshes.length;for(let s=1;s<a;s++)this._applyGGXFilter(t,s-1,s);n.autoClear=i}_applyGGXFilter(t,n,i){let a=this._renderer,s=this._pingPongRenderTarget,r=this._ggxMaterial,o=this._lodMeshes[i];o.material=r;let l=r.uniforms,c=i/(this._lodMeshes.length-1),f=n/(this._lodMeshes.length-1),p=Math.sqrt(c*c-f*f),u=0+c*1.25,d=p*u,{_lodMax:_}=this,b=this._sizeLods[i],v=3*b*(i>_-Ys?i-_+Ys:0),h=4*(this._cubeSize-b);l.envMap.value=t.texture,l.roughness.value=d,l.mipInt.value=_-n,_l(s,v,h,3*b,2*b),a.setRenderTarget(s),a.render(o,yu),l.envMap.value=s.texture,l.roughness.value=0,l.mipInt.value=_-i,_l(t,v,h,3*b,2*b),a.setRenderTarget(t),a.render(o,yu)}_blur(t,n,i,a,s){let r=this._pingPongRenderTarget;this._halfBlur(t,r,n,i,a,"latitudinal",s),this._halfBlur(r,t,i,i,a,"longitudinal",s)}_halfBlur(t,n,i,a,s,r,o){let l=this._renderer,c=this._blurMaterial;r!=="latitudinal"&&r!=="longitudinal"&&zt("blur direction must be either latitudinal or longitudinal!");let f=3,p=this._lodMeshes[a];p.material=c;let u=c.uniforms,d=this._sizeLods[i]-1,_=isFinite(s)?Math.PI/(2*d):2*Math.PI/(2*Br-1),b=s/_,v=isFinite(s)?1+Math.floor(f*b):Br;v>Br&&Ot(`sigmaRadians, ${s}, is too large and will clip, as it requested ${v} samples when the maximum is set to ${Br}`);let h=[],m=0;for(let w=0;w<Br;++w){let y=w/b,T=Math.exp(-y*y/2);h.push(T),w===0?m+=T:w<v&&(m+=2*T)}for(let w=0;w<h.length;w++)h[w]=h[w]/m;u.envMap.value=t.texture,u.samples.value=v,u.weights.value=h,u.latitudinal.value=r==="latitudinal",o&&(u.poleAxis.value=o);let{_lodMax:g}=this;u.dTheta.value=_,u.mipInt.value=g-i;let x=this._sizeLods[a],E=3*x*(a>g-Ys?a-g+Ys:0),A=4*(this._cubeSize-x);_l(n,E,A,3*x,2*x),l.setRenderTarget(n),l.render(p,yu)}};function CR(e){let t=[],n=[],i=[],a=e,s=e-Ys+1+E1.length;for(let r=0;r<s;r++){let o=Math.pow(2,a);t.push(o);let l=1/o;r>e-Ys?l=E1[r-e+Ys-1]:r===0&&(l=0),n.push(l);let c=1/(o-2),f=-c,p=1+c,u=[f,f,p,f,p,p,f,f,p,p,f,p],d=6,_=6,b=3,v=2,h=1,m=new Float32Array(b*_*d),g=new Float32Array(v*_*d),x=new Float32Array(h*_*d);for(let A=0;A<d;A++){let w=A%3*2/3-1,y=A>2?0:-1,T=[w,y,0,w+2/3,y,0,w+2/3,y+1,0,w,y,0,w+2/3,y+1,0,w,y+1,0];m.set(T,b*_*A),g.set(u,v*_*A);let L=[A,A,A,A,A,A];x.set(L,h*_*A)}let E=new An;E.setAttribute("position",new Bn(m,b)),E.setAttribute("uv",new Bn(g,v)),E.setAttribute("faceIndex",new Bn(x,h)),i.push(new mn(E,null)),a>Ys&&a--}return{lodMeshes:i,sizeLods:t,sigmas:n}}function w1(e,t,n){let i=new gi(e,t,n);return i.texture.mapping=du,i.texture.name="PMREM.cubeUv",i.scissorTest=!0,i}function _l(e,t,n,i,a){e.viewport.set(t,n,i,a),e.scissor.set(t,n,i,a)}function RR(e,t,n){return new Yn({name:"PMREMGGXConvolution",defines:{GGX_SAMPLES:AR,CUBEUV_TEXEL_WIDTH:1/t,CUBEUV_TEXEL_HEIGHT:1/n,CUBEUV_MAX_MIP:`${e}.0`},uniforms:{envMap:{value:null},roughness:{value:0},mipInt:{value:0}},vertexShader:dp(),fragmentShader:`

			precision highp float;
			precision highp int;

			varying vec3 vOutputDirection;

			uniform sampler2D envMap;
			uniform float roughness;
			uniform float mipInt;

			#define ENVMAP_TYPE_CUBE_UV
			#include <cube_uv_reflection_fragment>

			#define PI 3.14159265359

			// Van der Corput radical inverse
			float radicalInverse_VdC(uint bits) {
				bits = (bits << 16u) | (bits >> 16u);
				bits = ((bits & 0x55555555u) << 1u) | ((bits & 0xAAAAAAAAu) >> 1u);
				bits = ((bits & 0x33333333u) << 2u) | ((bits & 0xCCCCCCCCu) >> 2u);
				bits = ((bits & 0x0F0F0F0Fu) << 4u) | ((bits & 0xF0F0F0F0u) >> 4u);
				bits = ((bits & 0x00FF00FFu) << 8u) | ((bits & 0xFF00FF00u) >> 8u);
				return float(bits) * 2.3283064365386963e-10; // / 0x100000000
			}

			// Hammersley sequence
			vec2 hammersley(uint i, uint N) {
				return vec2(float(i) / float(N), radicalInverse_VdC(i));
			}

			// GGX VNDF importance sampling (Eric Heitz 2018)
			// "Sampling the GGX Distribution of Visible Normals"
			// https://jcgt.org/published/0007/04/01/
			vec3 importanceSampleGGX_VNDF(vec2 Xi, vec3 V, float roughness) {
				float alpha = roughness * roughness;

				// Section 4.1: Orthonormal basis
				vec3 T1 = vec3(1.0, 0.0, 0.0);
				vec3 T2 = cross(V, T1);

				// Section 4.2: Parameterization of projected area
				float r = sqrt(Xi.x);
				float phi = 2.0 * PI * Xi.y;
				float t1 = r * cos(phi);
				float t2 = r * sin(phi);
				float s = 0.5 * (1.0 + V.z);
				t2 = (1.0 - s) * sqrt(1.0 - t1 * t1) + s * t2;

				// Section 4.3: Reprojection onto hemisphere
				vec3 Nh = t1 * T1 + t2 * T2 + sqrt(max(0.0, 1.0 - t1 * t1 - t2 * t2)) * V;

				// Section 3.4: Transform back to ellipsoid configuration
				return normalize(vec3(alpha * Nh.x, alpha * Nh.y, max(0.0, Nh.z)));
			}

			void main() {
				vec3 N = normalize(vOutputDirection);
				vec3 V = N; // Assume view direction equals normal for pre-filtering

				vec3 prefilteredColor = vec3(0.0);
				float totalWeight = 0.0;

				// For very low roughness, just sample the environment directly
				if (roughness < 0.001) {
					gl_FragColor = vec4(bilinearCubeUV(envMap, N, mipInt), 1.0);
					return;
				}

				// Tangent space basis for VNDF sampling
				vec3 up = abs(N.z) < 0.999 ? vec3(0.0, 0.0, 1.0) : vec3(1.0, 0.0, 0.0);
				vec3 tangent = normalize(cross(up, N));
				vec3 bitangent = cross(N, tangent);

				for(uint i = 0u; i < uint(GGX_SAMPLES); i++) {
					vec2 Xi = hammersley(i, uint(GGX_SAMPLES));

					// For PMREM, V = N, so in tangent space V is always (0, 0, 1)
					vec3 H_tangent = importanceSampleGGX_VNDF(Xi, vec3(0.0, 0.0, 1.0), roughness);

					// Transform H back to world space
					vec3 H = normalize(tangent * H_tangent.x + bitangent * H_tangent.y + N * H_tangent.z);
					vec3 L = normalize(2.0 * dot(V, H) * H - V);

					float NdotL = max(dot(N, L), 0.0);

					if(NdotL > 0.0) {
						// Sample environment at fixed mip level
						// VNDF importance sampling handles the distribution filtering
						vec3 sampleColor = bilinearCubeUV(envMap, L, mipInt);

						// Weight by NdotL for the split-sum approximation
						// VNDF PDF naturally accounts for the visible microfacet distribution
						prefilteredColor += sampleColor * NdotL;
						totalWeight += NdotL;
					}
				}

				if (totalWeight > 0.0) {
					prefilteredColor = prefilteredColor / totalWeight;
				}

				gl_FragColor = vec4(prefilteredColor, 1.0);
			}
		`,blending:va,depthTest:!1,depthWrite:!1})}function DR(e,t,n){let i=new Float32Array(Br),a=new B(0,1,0);return new Yn({name:"SphericalGaussianBlur",defines:{n:Br,CUBEUV_TEXEL_WIDTH:1/t,CUBEUV_TEXEL_HEIGHT:1/n,CUBEUV_MAX_MIP:`${e}.0`},uniforms:{envMap:{value:null},samples:{value:1},weights:{value:i},latitudinal:{value:!1},dTheta:{value:0},mipInt:{value:0},poleAxis:{value:a}},vertexShader:dp(),fragmentShader:`

			precision mediump float;
			precision mediump int;

			varying vec3 vOutputDirection;

			uniform sampler2D envMap;
			uniform int samples;
			uniform float weights[ n ];
			uniform bool latitudinal;
			uniform float dTheta;
			uniform float mipInt;
			uniform vec3 poleAxis;

			#define ENVMAP_TYPE_CUBE_UV
			#include <cube_uv_reflection_fragment>

			vec3 getSample( float theta, vec3 axis ) {

				float cosTheta = cos( theta );
				// Rodrigues' axis-angle rotation
				vec3 sampleDirection = vOutputDirection * cosTheta
					+ cross( axis, vOutputDirection ) * sin( theta )
					+ axis * dot( axis, vOutputDirection ) * ( 1.0 - cosTheta );

				return bilinearCubeUV( envMap, sampleDirection, mipInt );

			}

			void main() {

				vec3 axis = latitudinal ? poleAxis : cross( poleAxis, vOutputDirection );

				if ( all( equal( axis, vec3( 0.0 ) ) ) ) {

					axis = vec3( vOutputDirection.z, 0.0, - vOutputDirection.x );

				}

				axis = normalize( axis );

				gl_FragColor = vec4( 0.0, 0.0, 0.0, 1.0 );
				gl_FragColor.rgb += weights[ 0 ] * getSample( 0.0, axis );

				for ( int i = 1; i < n; i++ ) {

					if ( i >= samples ) {

						break;

					}

					float theta = dTheta * float( i );
					gl_FragColor.rgb += weights[ i ] * getSample( -1.0 * theta, axis );
					gl_FragColor.rgb += weights[ i ] * getSample( theta, axis );

				}

			}
		`,blending:va,depthTest:!1,depthWrite:!1})}function C1(){return new Yn({name:"EquirectangularToCubeUV",uniforms:{envMap:{value:null}},vertexShader:dp(),fragmentShader:`

			precision mediump float;
			precision mediump int;

			varying vec3 vOutputDirection;

			uniform sampler2D envMap;

			#include <common>

			void main() {

				vec3 outputDirection = normalize( vOutputDirection );
				vec2 uv = equirectUv( outputDirection );

				gl_FragColor = vec4( texture2D ( envMap, uv ).rgb, 1.0 );

			}
		`,blending:va,depthTest:!1,depthWrite:!1})}function R1(){return new Yn({name:"CubemapToCubeUV",uniforms:{envMap:{value:null},flipEnvMap:{value:-1}},vertexShader:dp(),fragmentShader:`

			precision mediump float;
			precision mediump int;

			uniform float flipEnvMap;

			varying vec3 vOutputDirection;

			uniform samplerCube envMap;

			void main() {

				gl_FragColor = textureCube( envMap, vec3( flipEnvMap * vOutputDirection.x, vOutputDirection.yz ) );

			}
		`,blending:va,depthTest:!1,depthWrite:!1})}function dp(){return`

		precision mediump float;
		precision mediump int;

		attribute float faceIndex;

		varying vec3 vOutputDirection;

		// RH coordinate system; PMREM face-indexing convention
		vec3 getDirection( vec2 uv, float face ) {

			uv = 2.0 * uv - 1.0;

			vec3 direction = vec3( uv, 1.0 );

			if ( face == 0.0 ) {

				direction = direction.zyx; // ( 1, v, u ) pos x

			} else if ( face == 1.0 ) {

				direction = direction.xzy;
				direction.xz *= -1.0; // ( -u, 1, -v ) pos y

			} else if ( face == 2.0 ) {

				direction.x *= -1.0; // ( -u, v, 1 ) pos z

			} else if ( face == 3.0 ) {

				direction = direction.zyx;
				direction.xz *= -1.0; // ( -1, v, -u ) neg x

			} else if ( face == 4.0 ) {

				direction = direction.xzy;
				direction.xy *= -1.0; // ( -u, -1, v ) neg y

			} else if ( face == 5.0 ) {

				direction.z *= -1.0; // ( u, v, -1 ) neg z

			}

			return direction;

		}

		void main() {

			vOutputDirection = getDirection( uv, faceIndex );
			gl_Position = vec4( position, 1.0 );

		}
	`}var fp=class extends gi{constructor(t=1,n={}){super(t,t,n),this.isWebGLCubeRenderTarget=!0;let i={width:t,height:t,depth:1},a=[i,i,i,i,i,i];this.texture=new tu(a),this._setTextureOptions(n),this.texture.isRenderTargetTexture=!0}fromEquirectangularTexture(t,n){this.texture.type=n.type,this.texture.colorSpace=n.colorSpace,this.texture.generateMipmaps=n.generateMipmaps,this.texture.minFilter=n.minFilter,this.texture.magFilter=n.magFilter;let i={uniforms:{tEquirect:{value:null}},vertexShader:`

				varying vec3 vWorldDirection;

				vec3 transformDirection( in vec3 dir, in mat4 matrix ) {

					return normalize( ( matrix * vec4( dir, 0.0 ) ).xyz );

				}

				void main() {

					vWorldDirection = transformDirection( position, modelMatrix );

					#include <begin_vertex>
					#include <project_vertex>

				}
			`,fragmentShader:`

				uniform sampler2D tEquirect;

				varying vec3 vWorldDirection;

				#include <common>

				void main() {

					vec3 direction = normalize( vWorldDirection );

					vec2 sampleUV = equirectUv( direction );

					gl_FragColor = texture2D( tEquirect, sampleUV );

				}
			`},a=new Vs(5,5,5),s=new Yn({name:"CubemapFromEquirect",uniforms:Ir(i.uniforms),vertexShader:i.vertexShader,fragmentShader:i.fragmentShader,side:wn,blending:va});s.uniforms.tEquirect.value=n;let r=new mn(a,s),o=n.minFilter;return n.minFilter===Ws&&(n.minFilter=En),new gd(1,10,this).update(t,r),n.minFilter=o,r.geometry.dispose(),r.material.dispose(),this}clear(t,n=!0,i=!0,a=!0){let s=t.getRenderTarget();for(let r=0;r<6;r++)t.setRenderTarget(this,r),t.clear(n,i,a);t.setRenderTarget(s)}};function NR(e){let t=new WeakMap,n=new WeakMap,i=null;function a(u,d=!1){return u==null?null:d?r(u):s(u)}function s(u){if(u&&u.isTexture){let d=u.mapping;if(d===_d||d===yd)if(t.has(u)){let _=t.get(u).texture;return o(_,u.mapping)}else{let _=u.image;if(_&&_.height>0){let b=new fp(_.height);return b.fromEquirectangularTexture(e,u),t.set(u,b),u.addEventListener("dispose",c),o(b.texture,u.mapping)}else return null}}return u}function r(u){if(u&&u.isTexture){let d=u.mapping,_=d===_d||d===yd,b=d===Xs||d===Ur;if(_||b){let v=n.get(u),h=v!==void 0?v.texture.pmremVersion:0;if(u.isRenderTargetTexture&&u.pmremVersion!==h)return i===null&&(i=new up(e)),v=_?i.fromEquirectangular(u,v):i.fromCubemap(u,v),v.texture.pmremVersion=u.pmremVersion,n.set(u,v),v.texture;if(v!==void 0)return v.texture;{let m=u.image;return _&&m&&m.height>0||b&&m&&l(m)?(i===null&&(i=new up(e)),v=_?i.fromEquirectangular(u):i.fromCubemap(u),v.texture.pmremVersion=u.pmremVersion,n.set(u,v),u.addEventListener("dispose",f),v.texture):null}}}return u}function o(u,d){return d===_d?u.mapping=Xs:d===yd&&(u.mapping=Ur),u}function l(u){let d=0,_=6;for(let b=0;b<_;b++)u[b]!==void 0&&d++;return d===_}function c(u){let d=u.target;d.removeEventListener("dispose",c);let _=t.get(d);_!==void 0&&(t.delete(d),_.dispose())}function f(u){let d=u.target;d.removeEventListener("dispose",f);let _=n.get(d);_!==void 0&&(n.delete(d),_.dispose())}function p(){t=new WeakMap,n=new WeakMap,i!==null&&(i.dispose(),i=null)}return{get:a,dispose:p}}function UR(e){let t={};function n(i){if(t[i]!==void 0)return t[i];let a=e.getExtension(i);return t[i]=a,a}return{has:function(i){return n(i)!==null},init:function(){n("EXT_color_buffer_float"),n("WEBGL_clip_cull_distance"),n("OES_texture_float_linear"),n("EXT_color_buffer_half_float"),n("WEBGL_multisampled_render_to_texture"),n("WEBGL_render_shared_exponent")},get:function(i){let a=n(i);return a===null&&Gc("WebGLRenderer: "+i+" extension not supported."),a}}}function LR(e,t,n,i){let a={},s=new WeakMap;function r(p){let u=p.target;u.index!==null&&t.remove(u.index);for(let _ in u.attributes)t.remove(u.attributes[_]);u.removeEventListener("dispose",r),delete a[u.id];let d=s.get(u);d&&(t.remove(d),s.delete(u)),i.releaseStatesOfGeometry(u),u.isInstancedBufferGeometry===!0&&delete u._maxInstanceCount,n.memory.geometries--}function o(p,u){return a[u.id]===!0||(u.addEventListener("dispose",r),a[u.id]=!0,n.memory.geometries++),u}function l(p){let u=p.attributes;for(let d in u)t.update(u[d],e.ARRAY_BUFFER)}function c(p){let u=[],d=p.index,_=p.attributes.position,b=0;if(_===void 0)return;if(d!==null){let m=d.array;b=d.version;for(let g=0,x=m.length;g<x;g+=3){let E=m[g+0],A=m[g+1],w=m[g+2];u.push(E,A,A,w,w,E)}}else{let m=_.array;b=_.version;for(let g=0,x=m.length/3-1;g<x;g+=3){let E=g+0,A=g+1,w=g+2;u.push(E,A,A,w,w,E)}}let v=new(_.count>=65535?Zc:Yc)(u,1);v.version=b;let h=s.get(p);h&&t.remove(h),s.set(p,v)}function f(p){let u=s.get(p);if(u){let d=p.index;d!==null&&u.version<d.version&&c(p)}else c(p);return s.get(p)}return{get:o,update:l,getWireframeAttribute:f}}function IR(e,t,n){let i;function a(u){i=u}let s,r;function o(u){s=u.type,r=u.bytesPerElement}function l(u,d){e.drawElements(i,d,s,u*r),n.update(d,i,1)}function c(u,d,_){_!==0&&(e.drawElementsInstanced(i,d,s,u*r,_),n.update(d,i,_))}function f(u,d,_){if(_===0)return;t.get("WEBGL_multi_draw").multiDrawElementsWEBGL(i,d,0,s,u,0,_);let v=0;for(let h=0;h<_;h++)v+=d[h];n.update(v,i,1)}function p(u,d,_,b){if(_===0)return;let v=t.get("WEBGL_multi_draw");if(v===null)for(let h=0;h<u.length;h++)c(u[h]/r,d[h],b[h]);else{v.multiDrawElementsInstancedWEBGL(i,d,0,s,u,0,b,0,_);let h=0;for(let m=0;m<_;m++)h+=d[m]*b[m];n.update(h,i,1)}}this.setMode=a,this.setIndex=o,this.render=l,this.renderInstances=c,this.renderMultiDraw=f,this.renderMultiDrawInstances=p}function OR(e){let t={geometries:0,textures:0},n={frame:0,calls:0,triangles:0,points:0,lines:0};function i(s,r,o){switch(n.calls++,r){case e.TRIANGLES:n.triangles+=o*(s/3);break;case e.LINES:n.lines+=o*(s/2);break;case e.LINE_STRIP:n.lines+=o*(s-1);break;case e.LINE_LOOP:n.lines+=o*s;break;case e.POINTS:n.points+=o*s;break;default:zt("WebGLInfo: Unknown draw mode:",r);break}}function a(){n.calls=0,n.triangles=0,n.points=0,n.lines=0}return{memory:t,render:n,programs:null,autoReset:!0,reset:a,update:i}}function PR(e,t,n){let i=new WeakMap,a=new Ge;function s(r,o,l){let c=r.morphTargetInfluences,f=o.morphAttributes.position||o.morphAttributes.normal||o.morphAttributes.color,p=f!==void 0?f.length:0,u=i.get(o);if(u===void 0||u.count!==p){let T=function(){w.dispose(),i.delete(o),o.removeEventListener("dispose",T)};u!==void 0&&u.texture.dispose();let d=o.morphAttributes.position!==void 0,_=o.morphAttributes.normal!==void 0,b=o.morphAttributes.color!==void 0,v=o.morphAttributes.position||[],h=o.morphAttributes.normal||[],m=o.morphAttributes.color||[],g=0;d===!0&&(g=1),_===!0&&(g=2),b===!0&&(g=3);let x=o.attributes.position.count*g,E=1;x>t.maxTextureSize&&(E=Math.ceil(x/t.maxTextureSize),x=t.maxTextureSize);let A=new Float32Array(x*E*4*p),w=new kc(A,x,E,p);w.type=Ni,w.needsUpdate=!0;let y=g*4;for(let L=0;L<p;L++){let R=v[L],N=h[L],P=m[L],z=x*E*4*L;for(let F=0;F<R.count;F++){let V=F*y;d===!0&&(a.fromBufferAttribute(R,F),A[z+V+0]=a.x,A[z+V+1]=a.y,A[z+V+2]=a.z,A[z+V+3]=0),_===!0&&(a.fromBufferAttribute(N,F),A[z+V+4]=a.x,A[z+V+5]=a.y,A[z+V+6]=a.z,A[z+V+7]=0),b===!0&&(a.fromBufferAttribute(P,F),A[z+V+8]=a.x,A[z+V+9]=a.y,A[z+V+10]=a.z,A[z+V+11]=P.itemSize===4?a.w:1)}}u={count:p,texture:w,size:new te(x,E)},i.set(o,u),o.addEventListener("dispose",T)}if(r.isInstancedMesh===!0&&r.morphTexture!==null)l.getUniforms().setValue(e,"morphTexture",r.morphTexture,n);else{let d=0;for(let b=0;b<c.length;b++)d+=c[b];let _=o.morphTargetsRelative?1:1-d;l.getUniforms().setValue(e,"morphTargetBaseInfluence",_),l.getUniforms().setValue(e,"morphTargetInfluences",c)}l.getUniforms().setValue(e,"morphTargetsTexture",u.texture,n),l.getUniforms().setValue(e,"morphTargetsTextureSize",u.size)}return{update:s}}function BR(e,t,n,i,a){let s=new WeakMap;function r(c){let f=a.render.frame,p=c.geometry,u=t.get(c,p);if(s.get(u)!==f&&(t.update(u),s.set(u,f)),c.isInstancedMesh&&(c.hasEventListener("dispose",l)===!1&&c.addEventListener("dispose",l),s.get(c)!==f&&(n.update(c.instanceMatrix,e.ARRAY_BUFFER),c.instanceColor!==null&&n.update(c.instanceColor,e.ARRAY_BUFFER),s.set(c,f))),c.isSkinnedMesh){let d=c.skeleton;s.get(d)!==f&&(d.update(),s.set(d,f))}return u}function o(){s=new WeakMap}function l(c){let f=c.target;f.removeEventListener("dispose",l),i.releaseStatesOfObject(f),n.remove(f.instanceMatrix),f.instanceColor!==null&&n.remove(f.instanceColor)}return{update:r,dispose:o}}var zR={[sv]:"LINEAR_TONE_MAPPING",[rv]:"REINHARD_TONE_MAPPING",[ov]:"CINEON_TONE_MAPPING",[hu]:"ACES_FILMIC_TONE_MAPPING",[cv]:"AGX_TONE_MAPPING",[uv]:"NEUTRAL_TONE_MAPPING",[lv]:"CUSTOM_TONE_MAPPING"};function FR(e,t,n,i,a){let s=new gi(t,n,{type:e,depthBuffer:i,stencilBuffer:a}),r=new gi(t,n,{type:_a,depthBuffer:!1,stencilBuffer:!1}),o=new An;o.setAttribute("position",new Oe([-1,3,0,-1,-1,0,3,-1,0],3)),o.setAttribute("uv",new Oe([0,2,0,0,2,0],2));let l=new id({uniforms:{tDiffuse:{value:null}},vertexShader:`
			precision highp float;

			uniform mat4 modelViewMatrix;
			uniform mat4 projectionMatrix;

			attribute vec3 position;
			attribute vec2 uv;

			varying vec2 vUv;

			void main() {
				vUv = uv;
				gl_Position = projectionMatrix * modelViewMatrix * vec4( position, 1.0 );
			}`,fragmentShader:`
			precision highp float;

			uniform sampler2D tDiffuse;

			varying vec2 vUv;

			#include <tonemapping_pars_fragment>
			#include <colorspace_pars_fragment>

			void main() {
				gl_FragColor = texture2D( tDiffuse, vUv );

				#ifdef LINEAR_TONE_MAPPING
					gl_FragColor.rgb = LinearToneMapping( gl_FragColor.rgb );
				#elif defined( REINHARD_TONE_MAPPING )
					gl_FragColor.rgb = ReinhardToneMapping( gl_FragColor.rgb );
				#elif defined( CINEON_TONE_MAPPING )
					gl_FragColor.rgb = CineonToneMapping( gl_FragColor.rgb );
				#elif defined( ACES_FILMIC_TONE_MAPPING )
					gl_FragColor.rgb = ACESFilmicToneMapping( gl_FragColor.rgb );
				#elif defined( AGX_TONE_MAPPING )
					gl_FragColor.rgb = AgXToneMapping( gl_FragColor.rgb );
				#elif defined( NEUTRAL_TONE_MAPPING )
					gl_FragColor.rgb = NeutralToneMapping( gl_FragColor.rgb );
				#elif defined( CUSTOM_TONE_MAPPING )
					gl_FragColor.rgb = CustomToneMapping( gl_FragColor.rgb );
				#endif

				#ifdef SRGB_TRANSFER
					gl_FragColor = sRGBTransferOETF( gl_FragColor );
				#endif
			}`,depthTest:!1,depthWrite:!1}),c=new mn(o,l),f=new hl(-1,1,1,-1,0,1),p=null,u=null,d=!1,_,b=null,v=[],h=!1;this.setSize=function(m,g){s.setSize(m,g),r.setSize(m,g);for(let x=0;x<v.length;x++){let E=v[x];E.setSize&&E.setSize(m,g)}},this.setEffects=function(m){v=m,h=v.length>0&&v[0].isRenderPass===!0;let g=s.width,x=s.height;for(let E=0;E<v.length;E++){let A=v[E];A.setSize&&A.setSize(g,x)}},this.begin=function(m,g){if(d||m.toneMapping===Wi&&v.length===0)return!1;if(b=g,g!==null){let x=g.width,E=g.height;(s.width!==x||s.height!==E)&&this.setSize(x,E)}return h===!1&&m.setRenderTarget(s),_=m.toneMapping,m.toneMapping=Wi,!0},this.hasRenderPass=function(){return h},this.end=function(m,g){m.toneMapping=_,d=!0;let x=s,E=r;for(let A=0;A<v.length;A++){let w=v[A];if(w.enabled!==!1&&(w.render(m,E,x,g),w.needsSwap!==!1)){let y=x;x=E,E=y}}if(p!==m.outputColorSpace||u!==m.toneMapping){p=m.outputColorSpace,u=m.toneMapping,l.defines={},se.getTransfer(p)===me&&(l.defines.SRGB_TRANSFER="");let A=zR[u];A&&(l.defines[A]=""),l.needsUpdate=!0}l.uniforms.tDiffuse.value=x.texture,m.setRenderTarget(b),m.render(c,f),b=null,d=!1},this.isCompositing=function(){return d},this.dispose=function(){s.dispose(),r.dispose(),o.dispose(),l.dispose()}}var Z1=new qn,Uv=new Fs(1,1),J1=new kc,K1=new Qh,Q1=new tu,D1=[],N1=[],U1=new Float32Array(16),L1=new Float32Array(9),I1=new Float32Array(4);function xl(e,t,n){let i=e[0];if(i<=0||i>0)return e;let a=t*n,s=D1[a];if(s===void 0&&(s=new Float32Array(a),D1[a]=s),t!==0){i.toArray(s,0);for(let r=1,o=0;r!==t;++r)o+=n,e[r].toArray(s,o)}return s}function rn(e,t){if(e.length!==t.length)return!1;for(let n=0,i=e.length;n<i;n++)if(e[n]!==t[n])return!1;return!0}function on(e,t){for(let n=0,i=t.length;n<i;n++)e[n]=t[n]}function pp(e,t){let n=N1[t];n===void 0&&(n=new Int32Array(t),N1[t]=n);for(let i=0;i!==t;++i)n[i]=e.allocateTextureUnit();return n}function VR(e,t){let n=this.cache;n[0]!==t&&(e.uniform1f(this.addr,t),n[0]=t)}function HR(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y)&&(e.uniform2f(this.addr,t.x,t.y),n[0]=t.x,n[1]=t.y);else{if(rn(n,t))return;e.uniform2fv(this.addr,t),on(n,t)}}function GR(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y||n[2]!==t.z)&&(e.uniform3f(this.addr,t.x,t.y,t.z),n[0]=t.x,n[1]=t.y,n[2]=t.z);else if(t.r!==void 0)(n[0]!==t.r||n[1]!==t.g||n[2]!==t.b)&&(e.uniform3f(this.addr,t.r,t.g,t.b),n[0]=t.r,n[1]=t.g,n[2]=t.b);else{if(rn(n,t))return;e.uniform3fv(this.addr,t),on(n,t)}}function kR(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y||n[2]!==t.z||n[3]!==t.w)&&(e.uniform4f(this.addr,t.x,t.y,t.z,t.w),n[0]=t.x,n[1]=t.y,n[2]=t.z,n[3]=t.w);else{if(rn(n,t))return;e.uniform4fv(this.addr,t),on(n,t)}}function XR(e,t){let n=this.cache,i=t.elements;if(i===void 0){if(rn(n,t))return;e.uniformMatrix2fv(this.addr,!1,t),on(n,t)}else{if(rn(n,i))return;I1.set(i),e.uniformMatrix2fv(this.addr,!1,I1),on(n,i)}}function WR(e,t){let n=this.cache,i=t.elements;if(i===void 0){if(rn(n,t))return;e.uniformMatrix3fv(this.addr,!1,t),on(n,t)}else{if(rn(n,i))return;L1.set(i),e.uniformMatrix3fv(this.addr,!1,L1),on(n,i)}}function qR(e,t){let n=this.cache,i=t.elements;if(i===void 0){if(rn(n,t))return;e.uniformMatrix4fv(this.addr,!1,t),on(n,t)}else{if(rn(n,i))return;U1.set(i),e.uniformMatrix4fv(this.addr,!1,U1),on(n,i)}}function YR(e,t){let n=this.cache;n[0]!==t&&(e.uniform1i(this.addr,t),n[0]=t)}function ZR(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y)&&(e.uniform2i(this.addr,t.x,t.y),n[0]=t.x,n[1]=t.y);else{if(rn(n,t))return;e.uniform2iv(this.addr,t),on(n,t)}}function JR(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y||n[2]!==t.z)&&(e.uniform3i(this.addr,t.x,t.y,t.z),n[0]=t.x,n[1]=t.y,n[2]=t.z);else{if(rn(n,t))return;e.uniform3iv(this.addr,t),on(n,t)}}function KR(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y||n[2]!==t.z||n[3]!==t.w)&&(e.uniform4i(this.addr,t.x,t.y,t.z,t.w),n[0]=t.x,n[1]=t.y,n[2]=t.z,n[3]=t.w);else{if(rn(n,t))return;e.uniform4iv(this.addr,t),on(n,t)}}function QR(e,t){let n=this.cache;n[0]!==t&&(e.uniform1ui(this.addr,t),n[0]=t)}function jR(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y)&&(e.uniform2ui(this.addr,t.x,t.y),n[0]=t.x,n[1]=t.y);else{if(rn(n,t))return;e.uniform2uiv(this.addr,t),on(n,t)}}function $R(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y||n[2]!==t.z)&&(e.uniform3ui(this.addr,t.x,t.y,t.z),n[0]=t.x,n[1]=t.y,n[2]=t.z);else{if(rn(n,t))return;e.uniform3uiv(this.addr,t),on(n,t)}}function tD(e,t){let n=this.cache;if(t.x!==void 0)(n[0]!==t.x||n[1]!==t.y||n[2]!==t.z||n[3]!==t.w)&&(e.uniform4ui(this.addr,t.x,t.y,t.z,t.w),n[0]=t.x,n[1]=t.y,n[2]=t.z,n[3]=t.w);else{if(rn(n,t))return;e.uniform4uiv(this.addr,t),on(n,t)}}function eD(e,t,n){let i=this.cache,a=n.allocateTextureUnit();i[0]!==a&&(e.uniform1i(this.addr,a),i[0]=a);let s;this.type===e.SAMPLER_2D_SHADOW?(Uv.compareFunction=n.isReversedDepthBuffer()?op:rp,s=Uv):s=Z1,n.setTexture2D(t||s,a)}function nD(e,t,n){let i=this.cache,a=n.allocateTextureUnit();i[0]!==a&&(e.uniform1i(this.addr,a),i[0]=a),n.setTexture3D(t||K1,a)}function iD(e,t,n){let i=this.cache,a=n.allocateTextureUnit();i[0]!==a&&(e.uniform1i(this.addr,a),i[0]=a),n.setTextureCube(t||Q1,a)}function aD(e,t,n){let i=this.cache,a=n.allocateTextureUnit();i[0]!==a&&(e.uniform1i(this.addr,a),i[0]=a),n.setTexture2DArray(t||J1,a)}function sD(e){switch(e){case 5126:return VR;case 35664:return HR;case 35665:return GR;case 35666:return kR;case 35674:return XR;case 35675:return WR;case 35676:return qR;case 5124:case 35670:return YR;case 35667:case 35671:return ZR;case 35668:case 35672:return JR;case 35669:case 35673:return KR;case 5125:return QR;case 36294:return jR;case 36295:return $R;case 36296:return tD;case 35678:case 36198:case 36298:case 36306:case 35682:return eD;case 35679:case 36299:case 36307:return nD;case 35680:case 36300:case 36308:case 36293:return iD;case 36289:case 36303:case 36311:case 36292:return aD}}function rD(e,t){e.uniform1fv(this.addr,t)}function oD(e,t){let n=xl(t,this.size,2);e.uniform2fv(this.addr,n)}function lD(e,t){let n=xl(t,this.size,3);e.uniform3fv(this.addr,n)}function cD(e,t){let n=xl(t,this.size,4);e.uniform4fv(this.addr,n)}function uD(e,t){let n=xl(t,this.size,4);e.uniformMatrix2fv(this.addr,!1,n)}function fD(e,t){let n=xl(t,this.size,9);e.uniformMatrix3fv(this.addr,!1,n)}function hD(e,t){let n=xl(t,this.size,16);e.uniformMatrix4fv(this.addr,!1,n)}function dD(e,t){e.uniform1iv(this.addr,t)}function pD(e,t){e.uniform2iv(this.addr,t)}function mD(e,t){e.uniform3iv(this.addr,t)}function gD(e,t){e.uniform4iv(this.addr,t)}function vD(e,t){e.uniform1uiv(this.addr,t)}function _D(e,t){e.uniform2uiv(this.addr,t)}function yD(e,t){e.uniform3uiv(this.addr,t)}function xD(e,t){e.uniform4uiv(this.addr,t)}function SD(e,t,n){let i=this.cache,a=t.length,s=pp(n,a);rn(i,s)||(e.uniform1iv(this.addr,s),on(i,s));let r;this.type===e.SAMPLER_2D_SHADOW?r=Uv:r=Z1;for(let o=0;o!==a;++o)n.setTexture2D(t[o]||r,s[o])}function bD(e,t,n){let i=this.cache,a=t.length,s=pp(n,a);rn(i,s)||(e.uniform1iv(this.addr,s),on(i,s));for(let r=0;r!==a;++r)n.setTexture3D(t[r]||K1,s[r])}function MD(e,t,n){let i=this.cache,a=t.length,s=pp(n,a);rn(i,s)||(e.uniform1iv(this.addr,s),on(i,s));for(let r=0;r!==a;++r)n.setTextureCube(t[r]||Q1,s[r])}function TD(e,t,n){let i=this.cache,a=t.length,s=pp(n,a);rn(i,s)||(e.uniform1iv(this.addr,s),on(i,s));for(let r=0;r!==a;++r)n.setTexture2DArray(t[r]||J1,s[r])}function ED(e){switch(e){case 5126:return rD;case 35664:return oD;case 35665:return lD;case 35666:return cD;case 35674:return uD;case 35675:return fD;case 35676:return hD;case 5124:case 35670:return dD;case 35667:case 35671:return pD;case 35668:case 35672:return mD;case 35669:case 35673:return gD;case 5125:return vD;case 36294:return _D;case 36295:return yD;case 36296:return xD;case 35678:case 36198:case 36298:case 36306:case 35682:return SD;case 35679:case 36299:case 36307:return bD;case 35680:case 36300:case 36308:case 36293:return MD;case 36289:case 36303:case 36311:case 36292:return TD}}var Lv=class{constructor(t,n,i){this.id=t,this.addr=i,this.cache=[],this.type=n.type,this.setValue=sD(n.type)}},Iv=class{constructor(t,n,i){this.id=t,this.addr=i,this.cache=[],this.type=n.type,this.size=n.size,this.setValue=ED(n.type)}},Ov=class{constructor(t){this.id=t,this.seq=[],this.map={}}setValue(t,n,i){let a=this.seq;for(let s=0,r=a.length;s!==r;++s){let o=a[s];o.setValue(t,n[o.id],i)}}},Dv=/(\w+)(\])?(\[|\.)?/g;function O1(e,t){e.seq.push(t),e.map[t.id]=t}function AD(e,t,n){let i=e.name,a=i.length;for(Dv.lastIndex=0;;){let s=Dv.exec(i),r=Dv.lastIndex,o=s[1],l=s[2]==="]",c=s[3];if(l&&(o=o|0),c===void 0||c==="["&&r+2===a){O1(n,c===void 0?new Lv(o,e,t):new Iv(o,e,t));break}else{let p=n.map[o];p===void 0&&(p=new Ov(o),O1(n,p)),n=p}}}var yl=class{constructor(t,n){this.seq=[],this.map={};let i=t.getProgramParameter(n,t.ACTIVE_UNIFORMS);for(let r=0;r<i;++r){let o=t.getActiveUniform(n,r),l=t.getUniformLocation(n,o.name);AD(o,l,this)}let a=[],s=[];for(let r of this.seq)r.type===t.SAMPLER_2D_SHADOW||r.type===t.SAMPLER_CUBE_SHADOW||r.type===t.SAMPLER_2D_ARRAY_SHADOW?a.push(r):s.push(r);a.length>0&&(this.seq=a.concat(s))}setValue(t,n,i,a){let s=this.map[n];s!==void 0&&s.setValue(t,i,a)}setOptional(t,n,i){let a=n[i];a!==void 0&&this.setValue(t,i,a)}static upload(t,n,i,a){for(let s=0,r=n.length;s!==r;++s){let o=n[s],l=i[o.id];l.needsUpdate!==!1&&o.setValue(t,l.value,a)}}static seqWithValue(t,n){let i=[];for(let a=0,s=t.length;a!==s;++a){let r=t[a];r.id in n&&i.push(r)}return i}};function P1(e,t,n){let i=e.createShader(t);return e.shaderSource(i,n),e.compileShader(i),i}var wD=37297,CD=0;function RD(e,t){let n=e.split(`
`),i=[],a=Math.max(t-6,0),s=Math.min(t+6,n.length);for(let r=a;r<s;r++){let o=r+1;i.push(`${o===t?">":" "} ${o}: ${n[r]}`)}return i.join(`
`)}var B1=new Zt;function DD(e){se._getMatrix(B1,se.workingColorSpace,e);let t=`mat3( ${B1.elements.map(n=>n.toFixed(4))} )`;switch(se.getTransfer(e)){case Vc:return[t,"LinearTransferOETF"];case me:return[t,"sRGBTransferOETF"];default:return Ot("WebGLProgram: Unsupported color space: ",e),[t,"LinearTransferOETF"]}}function z1(e,t,n){let i=e.getShaderParameter(t,e.COMPILE_STATUS),s=(e.getShaderInfoLog(t)||"").trim();if(i&&s==="")return"";let r=/ERROR: 0:(\d+)/.exec(s);if(r){let o=parseInt(r[1]);return n.toUpperCase()+`

`+s+`

`+RD(e.getShaderSource(t),o)}else return s}function ND(e,t){let n=DD(t);return[`vec4 ${e}( vec4 value ) {`,`	return ${n[1]}( vec4( value.rgb * ${n[0]}, value.a ) );`,"}"].join(`
`)}var UD={[sv]:"Linear",[rv]:"Reinhard",[ov]:"Cineon",[hu]:"ACESFilmic",[cv]:"AgX",[uv]:"Neutral",[lv]:"Custom"};function LD(e,t){let n=UD[t];return n===void 0?(Ot("WebGLProgram: Unsupported toneMapping:",t),"vec3 "+e+"( vec3 color ) { return LinearToneMapping( color ); }"):"vec3 "+e+"( vec3 color ) { return "+n+"ToneMapping( color ); }"}var cp=new B;function ID(){se.getLuminanceCoefficients(cp);let e=cp.x.toFixed(4),t=cp.y.toFixed(4),n=cp.z.toFixed(4);return["float luminance( const in vec3 rgb ) {",`	const vec3 weights = vec3( ${e}, ${t}, ${n} );`,"	return dot( weights, rgb );","}"].join(`
`)}function OD(e){return[e.extensionClipCullDistance?"#extension GL_ANGLE_clip_cull_distance : require":"",e.extensionMultiDraw?"#extension GL_ANGLE_multi_draw : require":""].filter(Su).join(`
`)}function PD(e){let t=[];for(let n in e){let i=e[n];i!==!1&&t.push("#define "+n+" "+i)}return t.join(`
`)}function BD(e,t){let n={},i=e.getProgramParameter(t,e.ACTIVE_ATTRIBUTES);for(let a=0;a<i;a++){let s=e.getActiveAttrib(t,a),r=s.name,o=1;s.type===e.FLOAT_MAT2&&(o=2),s.type===e.FLOAT_MAT3&&(o=3),s.type===e.FLOAT_MAT4&&(o=4),n[r]={type:s.type,location:e.getAttribLocation(t,r),locationSize:o}}return n}function Su(e){return e!==""}function F1(e,t){let n=t.numSpotLightShadows+t.numSpotLightMaps-t.numSpotLightShadowsWithMaps;return e.replace(/NUM_DIR_LIGHTS/g,t.numDirLights).replace(/NUM_SPOT_LIGHTS/g,t.numSpotLights).replace(/NUM_SPOT_LIGHT_MAPS/g,t.numSpotLightMaps).replace(/NUM_SPOT_LIGHT_COORDS/g,n).replace(/NUM_RECT_AREA_LIGHTS/g,t.numRectAreaLights).replace(/NUM_POINT_LIGHTS/g,t.numPointLights).replace(/NUM_HEMI_LIGHTS/g,t.numHemiLights).replace(/NUM_DIR_LIGHT_SHADOWS/g,t.numDirLightShadows).replace(/NUM_SPOT_LIGHT_SHADOWS_WITH_MAPS/g,t.numSpotLightShadowsWithMaps).replace(/NUM_SPOT_LIGHT_SHADOWS/g,t.numSpotLightShadows).replace(/NUM_POINT_LIGHT_SHADOWS/g,t.numPointLightShadows)}function V1(e,t){return e.replace(/NUM_CLIPPING_PLANES/g,t.numClippingPlanes).replace(/UNION_CLIPPING_PLANES/g,t.numClippingPlanes-t.numClipIntersection)}var zD=/^[ \t]*#include +<([\w\d./]+)>/gm;function Pv(e){return e.replace(zD,VD)}var FD=new Map;function VD(e,t){let n=Kt[t];if(n===void 0){let i=FD.get(t);if(i!==void 0)n=Kt[i],Ot('WebGLRenderer: Shader chunk "%s" has been deprecated. Use "%s" instead.',t,i);else throw new Error("Can not resolve #include <"+t+">")}return Pv(n)}var HD=/#pragma unroll_loop_start\s+for\s*\(\s*int\s+i\s*=\s*(\d+)\s*;\s*i\s*<\s*(\d+)\s*;\s*i\s*\+\+\s*\)\s*{([\s\S]+?)}\s+#pragma unroll_loop_end/g;function H1(e){return e.replace(HD,GD)}function GD(e,t,n,i){let a="";for(let s=parseInt(t);s<parseInt(n);s++)a+=i.replace(/\[\s*i\s*\]/g,"[ "+s+" ]").replace(/UNROLLED_LOOP_INDEX/g,s);return a}function G1(e){let t=`precision ${e.precision} float;
	precision ${e.precision} int;
	precision ${e.precision} sampler2D;
	precision ${e.precision} samplerCube;
	precision ${e.precision} sampler3D;
	precision ${e.precision} sampler2DArray;
	precision ${e.precision} sampler2DShadow;
	precision ${e.precision} samplerCubeShadow;
	precision ${e.precision} sampler2DArrayShadow;
	precision ${e.precision} isampler2D;
	precision ${e.precision} isampler3D;
	precision ${e.precision} isamplerCube;
	precision ${e.precision} isampler2DArray;
	precision ${e.precision} usampler2D;
	precision ${e.precision} usampler3D;
	precision ${e.precision} usamplerCube;
	precision ${e.precision} usampler2DArray;
	`;return e.precision==="highp"?t+=`
#define HIGH_PRECISION`:e.precision==="mediump"?t+=`
#define MEDIUM_PRECISION`:e.precision==="lowp"&&(t+=`
#define LOW_PRECISION`),t}var kD={[fu]:"SHADOWMAP_TYPE_PCF",[dl]:"SHADOWMAP_TYPE_VSM"};function XD(e){return kD[e.shadowMapType]||"SHADOWMAP_TYPE_BASIC"}var WD={[Xs]:"ENVMAP_TYPE_CUBE",[Ur]:"ENVMAP_TYPE_CUBE",[du]:"ENVMAP_TYPE_CUBE_UV"};function qD(e){return e.envMap===!1?"ENVMAP_TYPE_CUBE":WD[e.envMapMode]||"ENVMAP_TYPE_CUBE"}var YD={[Ur]:"ENVMAP_MODE_REFRACTION"};function ZD(e){return e.envMap===!1?"ENVMAP_MODE_REFLECTION":YD[e.envMapMode]||"ENVMAP_MODE_REFLECTION"}var JD={[av]:"ENVMAP_BLENDING_MULTIPLY",[o1]:"ENVMAP_BLENDING_MIX",[l1]:"ENVMAP_BLENDING_ADD"};function KD(e){return e.envMap===!1?"ENVMAP_BLENDING_NONE":JD[e.combine]||"ENVMAP_BLENDING_NONE"}function QD(e){let t=e.envMapCubeUVHeight;if(t===null)return null;let n=Math.log2(t)-2,i=1/t;return{texelWidth:1/(3*Math.max(Math.pow(2,n),112)),texelHeight:i,maxMip:n}}function jD(e,t,n,i){let a=e.getContext(),s=n.defines,r=n.vertexShader,o=n.fragmentShader,l=XD(n),c=qD(n),f=ZD(n),p=KD(n),u=QD(n),d=OD(n),_=PD(s),b=a.createProgram(),v,h,m=n.glslVersion?"#version "+n.glslVersion+`
`:"";n.isRawShaderMaterial?(v=["#define SHADER_TYPE "+n.shaderType,"#define SHADER_NAME "+n.shaderName,_].filter(Su).join(`
`),v.length>0&&(v+=`
`),h=["#define SHADER_TYPE "+n.shaderType,"#define SHADER_NAME "+n.shaderName,_].filter(Su).join(`
`),h.length>0&&(h+=`
`)):(v=[G1(n),"#define SHADER_TYPE "+n.shaderType,"#define SHADER_NAME "+n.shaderName,_,n.extensionClipCullDistance?"#define USE_CLIP_DISTANCE":"",n.batching?"#define USE_BATCHING":"",n.batchingColor?"#define USE_BATCHING_COLOR":"",n.instancing?"#define USE_INSTANCING":"",n.instancingColor?"#define USE_INSTANCING_COLOR":"",n.instancingMorph?"#define USE_INSTANCING_MORPH":"",n.useFog&&n.fog?"#define USE_FOG":"",n.useFog&&n.fogExp2?"#define FOG_EXP2":"",n.map?"#define USE_MAP":"",n.envMap?"#define USE_ENVMAP":"",n.envMap?"#define "+f:"",n.lightMap?"#define USE_LIGHTMAP":"",n.aoMap?"#define USE_AOMAP":"",n.bumpMap?"#define USE_BUMPMAP":"",n.normalMap?"#define USE_NORMALMAP":"",n.normalMapObjectSpace?"#define USE_NORMALMAP_OBJECTSPACE":"",n.normalMapTangentSpace?"#define USE_NORMALMAP_TANGENTSPACE":"",n.displacementMap?"#define USE_DISPLACEMENTMAP":"",n.emissiveMap?"#define USE_EMISSIVEMAP":"",n.anisotropy?"#define USE_ANISOTROPY":"",n.anisotropyMap?"#define USE_ANISOTROPYMAP":"",n.clearcoatMap?"#define USE_CLEARCOATMAP":"",n.clearcoatRoughnessMap?"#define USE_CLEARCOAT_ROUGHNESSMAP":"",n.clearcoatNormalMap?"#define USE_CLEARCOAT_NORMALMAP":"",n.iridescenceMap?"#define USE_IRIDESCENCEMAP":"",n.iridescenceThicknessMap?"#define USE_IRIDESCENCE_THICKNESSMAP":"",n.specularMap?"#define USE_SPECULARMAP":"",n.specularColorMap?"#define USE_SPECULAR_COLORMAP":"",n.specularIntensityMap?"#define USE_SPECULAR_INTENSITYMAP":"",n.roughnessMap?"#define USE_ROUGHNESSMAP":"",n.metalnessMap?"#define USE_METALNESSMAP":"",n.alphaMap?"#define USE_ALPHAMAP":"",n.alphaHash?"#define USE_ALPHAHASH":"",n.transmission?"#define USE_TRANSMISSION":"",n.transmissionMap?"#define USE_TRANSMISSIONMAP":"",n.thicknessMap?"#define USE_THICKNESSMAP":"",n.sheenColorMap?"#define USE_SHEEN_COLORMAP":"",n.sheenRoughnessMap?"#define USE_SHEEN_ROUGHNESSMAP":"",n.mapUv?"#define MAP_UV "+n.mapUv:"",n.alphaMapUv?"#define ALPHAMAP_UV "+n.alphaMapUv:"",n.lightMapUv?"#define LIGHTMAP_UV "+n.lightMapUv:"",n.aoMapUv?"#define AOMAP_UV "+n.aoMapUv:"",n.emissiveMapUv?"#define EMISSIVEMAP_UV "+n.emissiveMapUv:"",n.bumpMapUv?"#define BUMPMAP_UV "+n.bumpMapUv:"",n.normalMapUv?"#define NORMALMAP_UV "+n.normalMapUv:"",n.displacementMapUv?"#define DISPLACEMENTMAP_UV "+n.displacementMapUv:"",n.metalnessMapUv?"#define METALNESSMAP_UV "+n.metalnessMapUv:"",n.roughnessMapUv?"#define ROUGHNESSMAP_UV "+n.roughnessMapUv:"",n.anisotropyMapUv?"#define ANISOTROPYMAP_UV "+n.anisotropyMapUv:"",n.clearcoatMapUv?"#define CLEARCOATMAP_UV "+n.clearcoatMapUv:"",n.clearcoatNormalMapUv?"#define CLEARCOAT_NORMALMAP_UV "+n.clearcoatNormalMapUv:"",n.clearcoatRoughnessMapUv?"#define CLEARCOAT_ROUGHNESSMAP_UV "+n.clearcoatRoughnessMapUv:"",n.iridescenceMapUv?"#define IRIDESCENCEMAP_UV "+n.iridescenceMapUv:"",n.iridescenceThicknessMapUv?"#define IRIDESCENCE_THICKNESSMAP_UV "+n.iridescenceThicknessMapUv:"",n.sheenColorMapUv?"#define SHEEN_COLORMAP_UV "+n.sheenColorMapUv:"",n.sheenRoughnessMapUv?"#define SHEEN_ROUGHNESSMAP_UV "+n.sheenRoughnessMapUv:"",n.specularMapUv?"#define SPECULARMAP_UV "+n.specularMapUv:"",n.specularColorMapUv?"#define SPECULAR_COLORMAP_UV "+n.specularColorMapUv:"",n.specularIntensityMapUv?"#define SPECULAR_INTENSITYMAP_UV "+n.specularIntensityMapUv:"",n.transmissionMapUv?"#define TRANSMISSIONMAP_UV "+n.transmissionMapUv:"",n.thicknessMapUv?"#define THICKNESSMAP_UV "+n.thicknessMapUv:"",n.vertexTangents&&n.flatShading===!1?"#define USE_TANGENT":"",n.vertexColors?"#define USE_COLOR":"",n.vertexAlphas?"#define USE_COLOR_ALPHA":"",n.vertexUv1s?"#define USE_UV1":"",n.vertexUv2s?"#define USE_UV2":"",n.vertexUv3s?"#define USE_UV3":"",n.pointsUvs?"#define USE_POINTS_UV":"",n.flatShading?"#define FLAT_SHADED":"",n.skinning?"#define USE_SKINNING":"",n.morphTargets?"#define USE_MORPHTARGETS":"",n.morphNormals&&n.flatShading===!1?"#define USE_MORPHNORMALS":"",n.morphColors?"#define USE_MORPHCOLORS":"",n.morphTargetsCount>0?"#define MORPHTARGETS_TEXTURE_STRIDE "+n.morphTextureStride:"",n.morphTargetsCount>0?"#define MORPHTARGETS_COUNT "+n.morphTargetsCount:"",n.doubleSided?"#define DOUBLE_SIDED":"",n.flipSided?"#define FLIP_SIDED":"",n.shadowMapEnabled?"#define USE_SHADOWMAP":"",n.shadowMapEnabled?"#define "+l:"",n.sizeAttenuation?"#define USE_SIZEATTENUATION":"",n.numLightProbes>0?"#define USE_LIGHT_PROBES":"",n.logarithmicDepthBuffer?"#define USE_LOGARITHMIC_DEPTH_BUFFER":"",n.reversedDepthBuffer?"#define USE_REVERSED_DEPTH_BUFFER":"","uniform mat4 modelMatrix;","uniform mat4 modelViewMatrix;","uniform mat4 projectionMatrix;","uniform mat4 viewMatrix;","uniform mat3 normalMatrix;","uniform vec3 cameraPosition;","uniform bool isOrthographic;","#ifdef USE_INSTANCING","	attribute mat4 instanceMatrix;","#endif","#ifdef USE_INSTANCING_COLOR","	attribute vec3 instanceColor;","#endif","#ifdef USE_INSTANCING_MORPH","	uniform sampler2D morphTexture;","#endif","attribute vec3 position;","attribute vec3 normal;","attribute vec2 uv;","#ifdef USE_UV1","	attribute vec2 uv1;","#endif","#ifdef USE_UV2","	attribute vec2 uv2;","#endif","#ifdef USE_UV3","	attribute vec2 uv3;","#endif","#ifdef USE_TANGENT","	attribute vec4 tangent;","#endif","#if defined( USE_COLOR_ALPHA )","	attribute vec4 color;","#elif defined( USE_COLOR )","	attribute vec3 color;","#endif","#ifdef USE_SKINNING","	attribute vec4 skinIndex;","	attribute vec4 skinWeight;","#endif",`
`].filter(Su).join(`
`),h=[G1(n),"#define SHADER_TYPE "+n.shaderType,"#define SHADER_NAME "+n.shaderName,_,n.useFog&&n.fog?"#define USE_FOG":"",n.useFog&&n.fogExp2?"#define FOG_EXP2":"",n.alphaToCoverage?"#define ALPHA_TO_COVERAGE":"",n.map?"#define USE_MAP":"",n.matcap?"#define USE_MATCAP":"",n.envMap?"#define USE_ENVMAP":"",n.envMap?"#define "+c:"",n.envMap?"#define "+f:"",n.envMap?"#define "+p:"",u?"#define CUBEUV_TEXEL_WIDTH "+u.texelWidth:"",u?"#define CUBEUV_TEXEL_HEIGHT "+u.texelHeight:"",u?"#define CUBEUV_MAX_MIP "+u.maxMip+".0":"",n.lightMap?"#define USE_LIGHTMAP":"",n.aoMap?"#define USE_AOMAP":"",n.bumpMap?"#define USE_BUMPMAP":"",n.normalMap?"#define USE_NORMALMAP":"",n.normalMapObjectSpace?"#define USE_NORMALMAP_OBJECTSPACE":"",n.normalMapTangentSpace?"#define USE_NORMALMAP_TANGENTSPACE":"",n.emissiveMap?"#define USE_EMISSIVEMAP":"",n.anisotropy?"#define USE_ANISOTROPY":"",n.anisotropyMap?"#define USE_ANISOTROPYMAP":"",n.clearcoat?"#define USE_CLEARCOAT":"",n.clearcoatMap?"#define USE_CLEARCOATMAP":"",n.clearcoatRoughnessMap?"#define USE_CLEARCOAT_ROUGHNESSMAP":"",n.clearcoatNormalMap?"#define USE_CLEARCOAT_NORMALMAP":"",n.dispersion?"#define USE_DISPERSION":"",n.iridescence?"#define USE_IRIDESCENCE":"",n.iridescenceMap?"#define USE_IRIDESCENCEMAP":"",n.iridescenceThicknessMap?"#define USE_IRIDESCENCE_THICKNESSMAP":"",n.specularMap?"#define USE_SPECULARMAP":"",n.specularColorMap?"#define USE_SPECULAR_COLORMAP":"",n.specularIntensityMap?"#define USE_SPECULAR_INTENSITYMAP":"",n.roughnessMap?"#define USE_ROUGHNESSMAP":"",n.metalnessMap?"#define USE_METALNESSMAP":"",n.alphaMap?"#define USE_ALPHAMAP":"",n.alphaTest?"#define USE_ALPHATEST":"",n.alphaHash?"#define USE_ALPHAHASH":"",n.sheen?"#define USE_SHEEN":"",n.sheenColorMap?"#define USE_SHEEN_COLORMAP":"",n.sheenRoughnessMap?"#define USE_SHEEN_ROUGHNESSMAP":"",n.transmission?"#define USE_TRANSMISSION":"",n.transmissionMap?"#define USE_TRANSMISSIONMAP":"",n.thicknessMap?"#define USE_THICKNESSMAP":"",n.vertexTangents&&n.flatShading===!1?"#define USE_TANGENT":"",n.vertexColors||n.instancingColor?"#define USE_COLOR":"",n.vertexAlphas||n.batchingColor?"#define USE_COLOR_ALPHA":"",n.vertexUv1s?"#define USE_UV1":"",n.vertexUv2s?"#define USE_UV2":"",n.vertexUv3s?"#define USE_UV3":"",n.pointsUvs?"#define USE_POINTS_UV":"",n.gradientMap?"#define USE_GRADIENTMAP":"",n.flatShading?"#define FLAT_SHADED":"",n.doubleSided?"#define DOUBLE_SIDED":"",n.flipSided?"#define FLIP_SIDED":"",n.shadowMapEnabled?"#define USE_SHADOWMAP":"",n.shadowMapEnabled?"#define "+l:"",n.premultipliedAlpha?"#define PREMULTIPLIED_ALPHA":"",n.numLightProbes>0?"#define USE_LIGHT_PROBES":"",n.decodeVideoTexture?"#define DECODE_VIDEO_TEXTURE":"",n.decodeVideoTextureEmissive?"#define DECODE_VIDEO_TEXTURE_EMISSIVE":"",n.logarithmicDepthBuffer?"#define USE_LOGARITHMIC_DEPTH_BUFFER":"",n.reversedDepthBuffer?"#define USE_REVERSED_DEPTH_BUFFER":"","uniform mat4 viewMatrix;","uniform vec3 cameraPosition;","uniform bool isOrthographic;",n.toneMapping!==Wi?"#define TONE_MAPPING":"",n.toneMapping!==Wi?Kt.tonemapping_pars_fragment:"",n.toneMapping!==Wi?LD("toneMapping",n.toneMapping):"",n.dithering?"#define DITHERING":"",n.opaque?"#define OPAQUE":"",Kt.colorspace_pars_fragment,ND("linearToOutputTexel",n.outputColorSpace),ID(),n.useDepthPacking?"#define DEPTH_PACKING "+n.depthPacking:"",`
`].filter(Su).join(`
`)),r=Pv(r),r=F1(r,n),r=V1(r,n),o=Pv(o),o=F1(o,n),o=V1(o,n),r=H1(r),o=H1(o),n.isRawShaderMaterial!==!0&&(m=`#version 300 es
`,v=[d,"#define attribute in","#define varying out","#define texture2D texture"].join(`
`)+`
`+v,h=["#define varying in",n.glslVersion===yv?"":"layout(location = 0) out highp vec4 pc_fragColor;",n.glslVersion===yv?"":"#define gl_FragColor pc_fragColor","#define gl_FragDepthEXT gl_FragDepth","#define texture2D texture","#define textureCube texture","#define texture2DProj textureProj","#define texture2DLodEXT textureLod","#define texture2DProjLodEXT textureProjLod","#define textureCubeLodEXT textureLod","#define texture2DGradEXT textureGrad","#define texture2DProjGradEXT textureProjGrad","#define textureCubeGradEXT textureGrad"].join(`
`)+`
`+h);let g=m+v+r,x=m+h+o,E=P1(a,a.VERTEX_SHADER,g),A=P1(a,a.FRAGMENT_SHADER,x);a.attachShader(b,E),a.attachShader(b,A),n.index0AttributeName!==void 0?a.bindAttribLocation(b,0,n.index0AttributeName):n.morphTargets===!0&&a.bindAttribLocation(b,0,"position"),a.linkProgram(b);function w(R){if(e.debug.checkShaderErrors){let N=a.getProgramInfoLog(b)||"",P=a.getShaderInfoLog(E)||"",z=a.getShaderInfoLog(A)||"",F=N.trim(),V=P.trim(),X=z.trim(),at=!0,et=!0;if(a.getProgramParameter(b,a.LINK_STATUS)===!1)if(at=!1,typeof e.debug.onShaderError=="function")e.debug.onShaderError(a,b,E,A);else{let dt=z1(a,E,"vertex"),St=z1(a,A,"fragment");zt("THREE.WebGLProgram: Shader Error "+a.getError()+" - VALIDATE_STATUS "+a.getProgramParameter(b,a.VALIDATE_STATUS)+`

Material Name: `+R.name+`
Material Type: `+R.type+`

Program Info Log: `+F+`
`+dt+`
`+St)}else F!==""?Ot("WebGLProgram: Program Info Log:",F):(V===""||X==="")&&(et=!1);et&&(R.diagnostics={runnable:at,programLog:F,vertexShader:{log:V,prefix:v},fragmentShader:{log:X,prefix:h}})}a.deleteShader(E),a.deleteShader(A),y=new yl(a,b),T=BD(a,b)}let y;this.getUniforms=function(){return y===void 0&&w(this),y};let T;this.getAttributes=function(){return T===void 0&&w(this),T};let L=n.rendererExtensionParallelShaderCompile===!1;return this.isReady=function(){return L===!1&&(L=a.getProgramParameter(b,wD)),L},this.destroy=function(){i.releaseStatesOfProgram(this),a.deleteProgram(b),this.program=void 0},this.type=n.shaderType,this.name=n.shaderName,this.id=CD++,this.cacheKey=t,this.usedTimes=1,this.program=b,this.vertexShader=E,this.fragmentShader=A,this}var $D=0,Bv=class{constructor(){this.shaderCache=new Map,this.materialCache=new Map}update(t){let n=t.vertexShader,i=t.fragmentShader,a=this._getShaderStage(n),s=this._getShaderStage(i),r=this._getShaderCacheForMaterial(t);return r.has(a)===!1&&(r.add(a),a.usedTimes++),r.has(s)===!1&&(r.add(s),s.usedTimes++),this}remove(t){let n=this.materialCache.get(t);for(let i of n)i.usedTimes--,i.usedTimes===0&&this.shaderCache.delete(i.code);return this.materialCache.delete(t),this}getVertexShaderID(t){return this._getShaderStage(t.vertexShader).id}getFragmentShaderID(t){return this._getShaderStage(t.fragmentShader).id}dispose(){this.shaderCache.clear(),this.materialCache.clear()}_getShaderCacheForMaterial(t){let n=this.materialCache,i=n.get(t);return i===void 0&&(i=new Set,n.set(t,i)),i}_getShaderStage(t){let n=this.shaderCache,i=n.get(t);return i===void 0&&(i=new zv(t),n.set(t,i)),i}},zv=class{constructor(t){this.id=$D++,this.code=t,this.usedTimes=0}};function tN(e,t,n,i,a,s){let r=new Xc,o=new Bv,l=new Set,c=[],f=new Map,p=i.logarithmicDepthBuffer,u=i.precision,d={MeshDepthMaterial:"depth",MeshDistanceMaterial:"distance",MeshNormalMaterial:"normal",MeshBasicMaterial:"basic",MeshLambertMaterial:"lambert",MeshPhongMaterial:"phong",MeshToonMaterial:"toon",MeshStandardMaterial:"physical",MeshPhysicalMaterial:"physical",MeshMatcapMaterial:"matcap",LineBasicMaterial:"basic",LineDashedMaterial:"dashed",PointsMaterial:"points",ShadowMaterial:"shadow",SpriteMaterial:"sprite"};function _(y){return l.add(y),y===0?"uv":`uv${y}`}function b(y,T,L,R,N){let P=R.fog,z=N.geometry,F=y.isMeshStandardMaterial||y.isMeshLambertMaterial||y.isMeshPhongMaterial?R.environment:null,V=y.isMeshStandardMaterial||y.isMeshLambertMaterial&&!y.envMap||y.isMeshPhongMaterial&&!y.envMap,X=t.get(y.envMap||F,V),at=X&&X.mapping===du?X.image.height:null,et=d[y.type];y.precision!==null&&(u=i.getMaxPrecision(y.precision),u!==y.precision&&Ot("WebGLProgram.getParameters:",y.precision,"not supported, using",u,"instead."));let dt=z.morphAttributes.position||z.morphAttributes.normal||z.morphAttributes.color,St=dt!==void 0?dt.length:0,ht=0;z.morphAttributes.position!==void 0&&(ht=1),z.morphAttributes.normal!==void 0&&(ht=2),z.morphAttributes.color!==void 0&&(ht=3);let Vt,Ce,xe,K;if(et){let oe=xa[et];Vt=oe.vertexShader,Ce=oe.fragmentShader}else Vt=y.vertexShader,Ce=y.fragmentShader,o.update(y),xe=o.getVertexShaderID(y),K=o.getFragmentShaderID(y);let j=e.getRenderTarget(),ct=e.state.buffers.depth.getReversed(),It=N.isInstancedMesh===!0,Lt=N.isBatchedMesh===!0,Pt=!!y.map,Ye=!!y.matcap,ie=!!X,Jt=!!y.aoMap,de=!!y.lightMap,Xt=!!y.bumpMap,Te=!!y.normalMap,D=!!y.displacementMap,ze=!!y.emissiveMap,re=!!y.metalnessMap,ue=!!y.roughnessMap,At=y.anisotropy>0,C=y.clearcoat>0,S=y.dispersion>0,O=y.iridescence>0,Z=y.sheen>0,Q=y.transmission>0,Y=At&&!!y.anisotropyMap,Tt=C&&!!y.clearcoatMap,lt=C&&!!y.clearcoatNormalMap,Rt=C&&!!y.clearcoatRoughnessMap,Ut=O&&!!y.iridescenceMap,tt=O&&!!y.iridescenceThicknessMap,rt=Z&&!!y.sheenColorMap,Mt=Z&&!!y.sheenRoughnessMap,bt=!!y.specularMap,gt=!!y.specularColorMap,Wt=!!y.specularIntensityMap,U=Q&&!!y.transmissionMap,ut=Q&&!!y.thicknessMap,ot=!!y.gradientMap,xt=!!y.alphaMap,nt=y.alphaTest>0,q=!!y.alphaHash,Et=!!y.extensions,Ht=Wi;y.toneMapped&&(j===null||j.isXRRenderTarget===!0)&&(Ht=e.toneMapping);let ge={shaderID:et,shaderType:y.type,shaderName:y.name,vertexShader:Vt,fragmentShader:Ce,defines:y.defines,customVertexShaderID:xe,customFragmentShaderID:K,isRawShaderMaterial:y.isRawShaderMaterial===!0,glslVersion:y.glslVersion,precision:u,batching:Lt,batchingColor:Lt&&N._colorsTexture!==null,instancing:It,instancingColor:It&&N.instanceColor!==null,instancingMorph:It&&N.morphTexture!==null,outputColorSpace:j===null?e.outputColorSpace:j.isXRRenderTarget===!0?j.texture.colorSpace:Rr,alphaToCoverage:!!y.alphaToCoverage,map:Pt,matcap:Ye,envMap:ie,envMapMode:ie&&X.mapping,envMapCubeUVHeight:at,aoMap:Jt,lightMap:de,bumpMap:Xt,normalMap:Te,displacementMap:D,emissiveMap:ze,normalMapObjectSpace:Te&&y.normalMapType===f1,normalMapTangentSpace:Te&&y.normalMapType===_v,metalnessMap:re,roughnessMap:ue,anisotropy:At,anisotropyMap:Y,clearcoat:C,clearcoatMap:Tt,clearcoatNormalMap:lt,clearcoatRoughnessMap:Rt,dispersion:S,iridescence:O,iridescenceMap:Ut,iridescenceThicknessMap:tt,sheen:Z,sheenColorMap:rt,sheenRoughnessMap:Mt,specularMap:bt,specularColorMap:gt,specularIntensityMap:Wt,transmission:Q,transmissionMap:U,thicknessMap:ut,gradientMap:ot,opaque:y.transparent===!1&&y.blending===wr&&y.alphaToCoverage===!1,alphaMap:xt,alphaTest:nt,alphaHash:q,combine:y.combine,mapUv:Pt&&_(y.map.channel),aoMapUv:Jt&&_(y.aoMap.channel),lightMapUv:de&&_(y.lightMap.channel),bumpMapUv:Xt&&_(y.bumpMap.channel),normalMapUv:Te&&_(y.normalMap.channel),displacementMapUv:D&&_(y.displacementMap.channel),emissiveMapUv:ze&&_(y.emissiveMap.channel),metalnessMapUv:re&&_(y.metalnessMap.channel),roughnessMapUv:ue&&_(y.roughnessMap.channel),anisotropyMapUv:Y&&_(y.anisotropyMap.channel),clearcoatMapUv:Tt&&_(y.clearcoatMap.channel),clearcoatNormalMapUv:lt&&_(y.clearcoatNormalMap.channel),clearcoatRoughnessMapUv:Rt&&_(y.clearcoatRoughnessMap.channel),iridescenceMapUv:Ut&&_(y.iridescenceMap.channel),iridescenceThicknessMapUv:tt&&_(y.iridescenceThicknessMap.channel),sheenColorMapUv:rt&&_(y.sheenColorMap.channel),sheenRoughnessMapUv:Mt&&_(y.sheenRoughnessMap.channel),specularMapUv:bt&&_(y.specularMap.channel),specularColorMapUv:gt&&_(y.specularColorMap.channel),specularIntensityMapUv:Wt&&_(y.specularIntensityMap.channel),transmissionMapUv:U&&_(y.transmissionMap.channel),thicknessMapUv:ut&&_(y.thicknessMap.channel),alphaMapUv:xt&&_(y.alphaMap.channel),vertexTangents:!!z.attributes.tangent&&(Te||At),vertexColors:y.vertexColors,vertexAlphas:y.vertexColors===!0&&!!z.attributes.color&&z.attributes.color.itemSize===4,pointsUvs:N.isPoints===!0&&!!z.attributes.uv&&(Pt||xt),fog:!!P,useFog:y.fog===!0,fogExp2:!!P&&P.isFogExp2,flatShading:y.wireframe===!1&&(y.flatShading===!0||z.attributes.normal===void 0&&Te===!1&&(y.isMeshLambertMaterial||y.isMeshPhongMaterial||y.isMeshStandardMaterial||y.isMeshPhysicalMaterial)),sizeAttenuation:y.sizeAttenuation===!0,logarithmicDepthBuffer:p,reversedDepthBuffer:ct,skinning:N.isSkinnedMesh===!0,morphTargets:z.morphAttributes.position!==void 0,morphNormals:z.morphAttributes.normal!==void 0,morphColors:z.morphAttributes.color!==void 0,morphTargetsCount:St,morphTextureStride:ht,numDirLights:T.directional.length,numPointLights:T.point.length,numSpotLights:T.spot.length,numSpotLightMaps:T.spotLightMap.length,numRectAreaLights:T.rectArea.length,numHemiLights:T.hemi.length,numDirLightShadows:T.directionalShadowMap.length,numPointLightShadows:T.pointShadowMap.length,numSpotLightShadows:T.spotShadowMap.length,numSpotLightShadowsWithMaps:T.numSpotLightShadowsWithMaps,numLightProbes:T.numLightProbes,numClippingPlanes:s.numPlanes,numClipIntersection:s.numIntersection,dithering:y.dithering,shadowMapEnabled:e.shadowMap.enabled&&L.length>0,shadowMapType:e.shadowMap.type,toneMapping:Ht,decodeVideoTexture:Pt&&y.map.isVideoTexture===!0&&se.getTransfer(y.map.colorSpace)===me,decodeVideoTextureEmissive:ze&&y.emissiveMap.isVideoTexture===!0&&se.getTransfer(y.emissiveMap.colorSpace)===me,premultipliedAlpha:y.premultipliedAlpha,doubleSided:y.side===ga,flipSided:y.side===wn,useDepthPacking:y.depthPacking>=0,depthPacking:y.depthPacking||0,index0AttributeName:y.index0AttributeName,extensionClipCullDistance:Et&&y.extensions.clipCullDistance===!0&&n.has("WEBGL_clip_cull_distance"),extensionMultiDraw:(Et&&y.extensions.multiDraw===!0||Lt)&&n.has("WEBGL_multi_draw"),rendererExtensionParallelShaderCompile:n.has("KHR_parallel_shader_compile"),customProgramCacheKey:y.customProgramCacheKey()};return ge.vertexUv1s=l.has(1),ge.vertexUv2s=l.has(2),ge.vertexUv3s=l.has(3),l.clear(),ge}function v(y){let T=[];if(y.shaderID?T.push(y.shaderID):(T.push(y.customVertexShaderID),T.push(y.customFragmentShaderID)),y.defines!==void 0)for(let L in y.defines)T.push(L),T.push(y.defines[L]);return y.isRawShaderMaterial===!1&&(h(T,y),m(T,y),T.push(e.outputColorSpace)),T.push(y.customProgramCacheKey),T.join()}function h(y,T){y.push(T.precision),y.push(T.outputColorSpace),y.push(T.envMapMode),y.push(T.envMapCubeUVHeight),y.push(T.mapUv),y.push(T.alphaMapUv),y.push(T.lightMapUv),y.push(T.aoMapUv),y.push(T.bumpMapUv),y.push(T.normalMapUv),y.push(T.displacementMapUv),y.push(T.emissiveMapUv),y.push(T.metalnessMapUv),y.push(T.roughnessMapUv),y.push(T.anisotropyMapUv),y.push(T.clearcoatMapUv),y.push(T.clearcoatNormalMapUv),y.push(T.clearcoatRoughnessMapUv),y.push(T.iridescenceMapUv),y.push(T.iridescenceThicknessMapUv),y.push(T.sheenColorMapUv),y.push(T.sheenRoughnessMapUv),y.push(T.specularMapUv),y.push(T.specularColorMapUv),y.push(T.specularIntensityMapUv),y.push(T.transmissionMapUv),y.push(T.thicknessMapUv),y.push(T.combine),y.push(T.fogExp2),y.push(T.sizeAttenuation),y.push(T.morphTargetsCount),y.push(T.morphAttributeCount),y.push(T.numDirLights),y.push(T.numPointLights),y.push(T.numSpotLights),y.push(T.numSpotLightMaps),y.push(T.numHemiLights),y.push(T.numRectAreaLights),y.push(T.numDirLightShadows),y.push(T.numPointLightShadows),y.push(T.numSpotLightShadows),y.push(T.numSpotLightShadowsWithMaps),y.push(T.numLightProbes),y.push(T.shadowMapType),y.push(T.toneMapping),y.push(T.numClippingPlanes),y.push(T.numClipIntersection),y.push(T.depthPacking)}function m(y,T){r.disableAll(),T.instancing&&r.enable(0),T.instancingColor&&r.enable(1),T.instancingMorph&&r.enable(2),T.matcap&&r.enable(3),T.envMap&&r.enable(4),T.normalMapObjectSpace&&r.enable(5),T.normalMapTangentSpace&&r.enable(6),T.clearcoat&&r.enable(7),T.iridescence&&r.enable(8),T.alphaTest&&r.enable(9),T.vertexColors&&r.enable(10),T.vertexAlphas&&r.enable(11),T.vertexUv1s&&r.enable(12),T.vertexUv2s&&r.enable(13),T.vertexUv3s&&r.enable(14),T.vertexTangents&&r.enable(15),T.anisotropy&&r.enable(16),T.alphaHash&&r.enable(17),T.batching&&r.enable(18),T.dispersion&&r.enable(19),T.batchingColor&&r.enable(20),T.gradientMap&&r.enable(21),y.push(r.mask),r.disableAll(),T.fog&&r.enable(0),T.useFog&&r.enable(1),T.flatShading&&r.enable(2),T.logarithmicDepthBuffer&&r.enable(3),T.reversedDepthBuffer&&r.enable(4),T.skinning&&r.enable(5),T.morphTargets&&r.enable(6),T.morphNormals&&r.enable(7),T.morphColors&&r.enable(8),T.premultipliedAlpha&&r.enable(9),T.shadowMapEnabled&&r.enable(10),T.doubleSided&&r.enable(11),T.flipSided&&r.enable(12),T.useDepthPacking&&r.enable(13),T.dithering&&r.enable(14),T.transmission&&r.enable(15),T.sheen&&r.enable(16),T.opaque&&r.enable(17),T.pointsUvs&&r.enable(18),T.decodeVideoTexture&&r.enable(19),T.decodeVideoTextureEmissive&&r.enable(20),T.alphaToCoverage&&r.enable(21),y.push(r.mask)}function g(y){let T=d[y.type],L;if(T){let R=xa[T];L=M1.clone(R.uniforms)}else L=y.uniforms;return L}function x(y,T){let L=f.get(T);return L!==void 0?++L.usedTimes:(L=new jD(e,T,y,a),c.push(L),f.set(T,L)),L}function E(y){if(--y.usedTimes===0){let T=c.indexOf(y);c[T]=c[c.length-1],c.pop(),f.delete(y.cacheKey),y.destroy()}}function A(y){o.remove(y)}function w(){o.dispose()}return{getParameters:b,getProgramCacheKey:v,getUniforms:g,acquireProgram:x,releaseProgram:E,releaseShaderCache:A,programs:c,dispose:w}}function eN(){let e=new WeakMap;function t(r){return e.has(r)}function n(r){let o=e.get(r);return o===void 0&&(o={},e.set(r,o)),o}function i(r){e.delete(r)}function a(r,o,l){e.get(r)[o]=l}function s(){e=new WeakMap}return{has:t,get:n,remove:i,update:a,dispose:s}}function nN(e,t){return e.groupOrder!==t.groupOrder?e.groupOrder-t.groupOrder:e.renderOrder!==t.renderOrder?e.renderOrder-t.renderOrder:e.material.id!==t.material.id?e.material.id-t.material.id:e.materialVariant!==t.materialVariant?e.materialVariant-t.materialVariant:e.z!==t.z?e.z-t.z:e.id-t.id}function k1(e,t){return e.groupOrder!==t.groupOrder?e.groupOrder-t.groupOrder:e.renderOrder!==t.renderOrder?e.renderOrder-t.renderOrder:e.z!==t.z?t.z-e.z:e.id-t.id}function X1(){let e=[],t=0,n=[],i=[],a=[];function s(){t=0,n.length=0,i.length=0,a.length=0}function r(u){let d=0;return u.isInstancedMesh&&(d+=2),u.isSkinnedMesh&&(d+=1),d}function o(u,d,_,b,v,h){let m=e[t];return m===void 0?(m={id:u.id,object:u,geometry:d,material:_,materialVariant:r(u),groupOrder:b,renderOrder:u.renderOrder,z:v,group:h},e[t]=m):(m.id=u.id,m.object=u,m.geometry=d,m.material=_,m.materialVariant=r(u),m.groupOrder=b,m.renderOrder=u.renderOrder,m.z=v,m.group=h),t++,m}function l(u,d,_,b,v,h){let m=o(u,d,_,b,v,h);_.transmission>0?i.push(m):_.transparent===!0?a.push(m):n.push(m)}function c(u,d,_,b,v,h){let m=o(u,d,_,b,v,h);_.transmission>0?i.unshift(m):_.transparent===!0?a.unshift(m):n.unshift(m)}function f(u,d){n.length>1&&n.sort(u||nN),i.length>1&&i.sort(d||k1),a.length>1&&a.sort(d||k1)}function p(){for(let u=t,d=e.length;u<d;u++){let _=e[u];if(_.id===null)break;_.id=null,_.object=null,_.geometry=null,_.material=null,_.group=null}}return{opaque:n,transmissive:i,transparent:a,init:s,push:l,unshift:c,finish:p,sort:f}}function iN(){let e=new WeakMap;function t(i,a){let s=e.get(i),r;return s===void 0?(r=new X1,e.set(i,[r])):a>=s.length?(r=new X1,s.push(r)):r=s[a],r}function n(){e=new WeakMap}return{get:t,dispose:n}}function aN(){let e={};return{get:function(t){if(e[t.id]!==void 0)return e[t.id];let n;switch(t.type){case"DirectionalLight":n={direction:new B,color:new Ft};break;case"SpotLight":n={position:new B,direction:new B,color:new Ft,distance:0,coneCos:0,penumbraCos:0,decay:0};break;case"PointLight":n={position:new B,color:new Ft,distance:0,decay:0};break;case"HemisphereLight":n={direction:new B,skyColor:new Ft,groundColor:new Ft};break;case"RectAreaLight":n={color:new Ft,position:new B,halfWidth:new B,halfHeight:new B};break}return e[t.id]=n,n}}}function sN(){let e={};return{get:function(t){if(e[t.id]!==void 0)return e[t.id];let n;switch(t.type){case"DirectionalLight":n={shadowIntensity:1,shadowBias:0,shadowNormalBias:0,shadowRadius:1,shadowMapSize:new te};break;case"SpotLight":n={shadowIntensity:1,shadowBias:0,shadowNormalBias:0,shadowRadius:1,shadowMapSize:new te};break;case"PointLight":n={shadowIntensity:1,shadowBias:0,shadowNormalBias:0,shadowRadius:1,shadowMapSize:new te,shadowCameraNear:1,shadowCameraFar:1e3};break}return e[t.id]=n,n}}}var rN=0;function oN(e,t){return(t.castShadow?2:0)-(e.castShadow?2:0)+(t.map?1:0)-(e.map?1:0)}function lN(e){let t=new aN,n=sN(),i={version:0,hash:{directionalLength:-1,pointLength:-1,spotLength:-1,rectAreaLength:-1,hemiLength:-1,numDirectionalShadows:-1,numPointShadows:-1,numSpotShadows:-1,numSpotMaps:-1,numLightProbes:-1},ambient:[0,0,0],probe:[],directional:[],directionalShadow:[],directionalShadowMap:[],directionalShadowMatrix:[],spot:[],spotLightMap:[],spotShadow:[],spotShadowMap:[],spotLightMatrix:[],rectArea:[],rectAreaLTC1:null,rectAreaLTC2:null,point:[],pointShadow:[],pointShadowMap:[],pointShadowMatrix:[],hemi:[],numSpotLightShadowsWithMaps:0,numLightProbes:0};for(let c=0;c<9;c++)i.probe.push(new B);let a=new B,s=new Ie,r=new Ie;function o(c){let f=0,p=0,u=0;for(let T=0;T<9;T++)i.probe[T].set(0,0,0);let d=0,_=0,b=0,v=0,h=0,m=0,g=0,x=0,E=0,A=0,w=0;c.sort(oN);for(let T=0,L=c.length;T<L;T++){let R=c[T],N=R.color,P=R.intensity,z=R.distance,F=null;if(R.shadow&&R.shadow.map&&(R.shadow.map.texture.format===Lr?F=R.shadow.map.texture:F=R.shadow.map.depthTexture||R.shadow.map.texture),R.isAmbientLight)f+=N.r*P,p+=N.g*P,u+=N.b*P;else if(R.isLightProbe){for(let V=0;V<9;V++)i.probe[V].addScaledVector(R.sh.coefficients[V],P);w++}else if(R.isDirectionalLight){let V=t.get(R);if(V.color.copy(R.color).multiplyScalar(R.intensity),R.castShadow){let X=R.shadow,at=n.get(R);at.shadowIntensity=X.intensity,at.shadowBias=X.bias,at.shadowNormalBias=X.normalBias,at.shadowRadius=X.radius,at.shadowMapSize=X.mapSize,i.directionalShadow[d]=at,i.directionalShadowMap[d]=F,i.directionalShadowMatrix[d]=R.shadow.matrix,m++}i.directional[d]=V,d++}else if(R.isSpotLight){let V=t.get(R);V.position.setFromMatrixPosition(R.matrixWorld),V.color.copy(N).multiplyScalar(P),V.distance=z,V.coneCos=Math.cos(R.angle),V.penumbraCos=Math.cos(R.angle*(1-R.penumbra)),V.decay=R.decay,i.spot[b]=V;let X=R.shadow;if(R.map&&(i.spotLightMap[E]=R.map,E++,X.updateMatrices(R),R.castShadow&&A++),i.spotLightMatrix[b]=X.matrix,R.castShadow){let at=n.get(R);at.shadowIntensity=X.intensity,at.shadowBias=X.bias,at.shadowNormalBias=X.normalBias,at.shadowRadius=X.radius,at.shadowMapSize=X.mapSize,i.spotShadow[b]=at,i.spotShadowMap[b]=F,x++}b++}else if(R.isRectAreaLight){let V=t.get(R);V.color.copy(N).multiplyScalar(P),V.halfWidth.set(R.width*.5,0,0),V.halfHeight.set(0,R.height*.5,0),i.rectArea[v]=V,v++}else if(R.isPointLight){let V=t.get(R);if(V.color.copy(R.color).multiplyScalar(R.intensity),V.distance=R.distance,V.decay=R.decay,R.castShadow){let X=R.shadow,at=n.get(R);at.shadowIntensity=X.intensity,at.shadowBias=X.bias,at.shadowNormalBias=X.normalBias,at.shadowRadius=X.radius,at.shadowMapSize=X.mapSize,at.shadowCameraNear=X.camera.near,at.shadowCameraFar=X.camera.far,i.pointShadow[_]=at,i.pointShadowMap[_]=F,i.pointShadowMatrix[_]=R.shadow.matrix,g++}i.point[_]=V,_++}else if(R.isHemisphereLight){let V=t.get(R);V.skyColor.copy(R.color).multiplyScalar(P),V.groundColor.copy(R.groundColor).multiplyScalar(P),i.hemi[h]=V,h++}}v>0&&(e.has("OES_texture_float_linear")===!0?(i.rectAreaLTC1=mt.LTC_FLOAT_1,i.rectAreaLTC2=mt.LTC_FLOAT_2):(i.rectAreaLTC1=mt.LTC_HALF_1,i.rectAreaLTC2=mt.LTC_HALF_2)),i.ambient[0]=f,i.ambient[1]=p,i.ambient[2]=u;let y=i.hash;(y.directionalLength!==d||y.pointLength!==_||y.spotLength!==b||y.rectAreaLength!==v||y.hemiLength!==h||y.numDirectionalShadows!==m||y.numPointShadows!==g||y.numSpotShadows!==x||y.numSpotMaps!==E||y.numLightProbes!==w)&&(i.directional.length=d,i.spot.length=b,i.rectArea.length=v,i.point.length=_,i.hemi.length=h,i.directionalShadow.length=m,i.directionalShadowMap.length=m,i.pointShadow.length=g,i.pointShadowMap.length=g,i.spotShadow.length=x,i.spotShadowMap.length=x,i.directionalShadowMatrix.length=m,i.pointShadowMatrix.length=g,i.spotLightMatrix.length=x+E-A,i.spotLightMap.length=E,i.numSpotLightShadowsWithMaps=A,i.numLightProbes=w,y.directionalLength=d,y.pointLength=_,y.spotLength=b,y.rectAreaLength=v,y.hemiLength=h,y.numDirectionalShadows=m,y.numPointShadows=g,y.numSpotShadows=x,y.numSpotMaps=E,y.numLightProbes=w,i.version=rN++)}function l(c,f){let p=0,u=0,d=0,_=0,b=0,v=f.matrixWorldInverse;for(let h=0,m=c.length;h<m;h++){let g=c[h];if(g.isDirectionalLight){let x=i.directional[p];x.direction.setFromMatrixPosition(g.matrixWorld),a.setFromMatrixPosition(g.target.matrixWorld),x.direction.sub(a),x.direction.transformDirection(v),p++}else if(g.isSpotLight){let x=i.spot[d];x.position.setFromMatrixPosition(g.matrixWorld),x.position.applyMatrix4(v),x.direction.setFromMatrixPosition(g.matrixWorld),a.setFromMatrixPosition(g.target.matrixWorld),x.direction.sub(a),x.direction.transformDirection(v),d++}else if(g.isRectAreaLight){let x=i.rectArea[_];x.position.setFromMatrixPosition(g.matrixWorld),x.position.applyMatrix4(v),r.identity(),s.copy(g.matrixWorld),s.premultiply(v),r.extractRotation(s),x.halfWidth.set(g.width*.5,0,0),x.halfHeight.set(0,g.height*.5,0),x.halfWidth.applyMatrix4(r),x.halfHeight.applyMatrix4(r),_++}else if(g.isPointLight){let x=i.point[u];x.position.setFromMatrixPosition(g.matrixWorld),x.position.applyMatrix4(v),u++}else if(g.isHemisphereLight){let x=i.hemi[b];x.direction.setFromMatrixPosition(g.matrixWorld),x.direction.transformDirection(v),b++}}}return{setup:o,setupView:l,state:i}}function W1(e){let t=new lN(e),n=[],i=[];function a(f){c.camera=f,n.length=0,i.length=0}function s(f){n.push(f)}function r(f){i.push(f)}function o(){t.setup(n)}function l(f){t.setupView(n,f)}let c={lightsArray:n,shadowsArray:i,camera:null,lights:t,transmissionRenderTarget:{}};return{init:a,state:c,setupLights:o,setupLightsView:l,pushLight:s,pushShadow:r}}function cN(e){let t=new WeakMap;function n(a,s=0){let r=t.get(a),o;return r===void 0?(o=new W1(e),t.set(a,[o])):s>=r.length?(o=new W1(e),r.push(o)):o=r[s],o}function i(){t=new WeakMap}return{get:n,dispose:i}}var uN=`void main() {
	gl_Position = vec4( position, 1.0 );
}`,fN=`uniform sampler2D shadow_pass;
uniform vec2 resolution;
uniform float radius;
void main() {
	const float samples = float( VSM_SAMPLES );
	float mean = 0.0;
	float squared_mean = 0.0;
	float uvStride = samples <= 1.0 ? 0.0 : 2.0 / ( samples - 1.0 );
	float uvStart = samples <= 1.0 ? 0.0 : - 1.0;
	for ( float i = 0.0; i < samples; i ++ ) {
		float uvOffset = uvStart + i * uvStride;
		#ifdef HORIZONTAL_PASS
			vec2 distribution = texture2D( shadow_pass, ( gl_FragCoord.xy + vec2( uvOffset, 0.0 ) * radius ) / resolution ).rg;
			mean += distribution.x;
			squared_mean += distribution.y * distribution.y + distribution.x * distribution.x;
		#else
			float depth = texture2D( shadow_pass, ( gl_FragCoord.xy + vec2( 0.0, uvOffset ) * radius ) / resolution ).r;
			mean += depth;
			squared_mean += depth * depth;
		#endif
	}
	mean = mean / samples;
	squared_mean = squared_mean / samples;
	float std_dev = sqrt( max( 0.0, squared_mean - mean * mean ) );
	gl_FragColor = vec4( mean, std_dev, 0.0, 1.0 );
}`,hN=[new B(1,0,0),new B(-1,0,0),new B(0,1,0),new B(0,-1,0),new B(0,0,1),new B(0,0,-1)],dN=[new B(0,-1,0),new B(0,-1,0),new B(0,0,1),new B(0,0,-1),new B(0,-1,0),new B(0,-1,0)],q1=new Ie,xu=new B,Nv=new B;function pN(e,t,n){let i=new ol,a=new te,s=new te,r=new Ge,o=new ad,l=new sd,c={},f=n.maxTextureSize,p={[Wa]:wn,[wn]:Wa,[ga]:ga},u=new Yn({defines:{VSM_SAMPLES:8},uniforms:{shadow_pass:{value:null},resolution:{value:new te},radius:{value:4}},vertexShader:uN,fragmentShader:fN}),d=u.clone();d.defines.HORIZONTAL_PASS=1;let _=new An;_.setAttribute("position",new Bn(new Float32Array([-1,-1,.5,3,-1,.5,-1,3,.5]),3));let b=new mn(_,u),v=this;this.enabled=!1,this.autoUpdate=!0,this.needsUpdate=!1,this.type=fu;let h=this.type;this.render=function(A,w,y){if(v.enabled===!1||v.autoUpdate===!1&&v.needsUpdate===!1||A.length===0)return;this.type===GM&&(Ot("WebGLShadowMap: PCFSoftShadowMap has been deprecated. Using PCFShadowMap instead."),this.type=fu);let T=e.getRenderTarget(),L=e.getActiveCubeFace(),R=e.getActiveMipmapLevel(),N=e.state;N.setBlending(va),N.buffers.depth.getReversed()===!0?N.buffers.color.setClear(0,0,0,0):N.buffers.color.setClear(1,1,1,1),N.buffers.depth.setTest(!0),N.setScissorTest(!1);let P=h!==this.type;P&&w.traverse(function(z){z.material&&(Array.isArray(z.material)?z.material.forEach(F=>F.needsUpdate=!0):z.material.needsUpdate=!0)});for(let z=0,F=A.length;z<F;z++){let V=A[z],X=V.shadow;if(X===void 0){Ot("WebGLShadowMap:",V,"has no shadow.");continue}if(X.autoUpdate===!1&&X.needsUpdate===!1)continue;a.copy(X.mapSize);let at=X.getFrameExtents();a.multiply(at),s.copy(X.mapSize),(a.x>f||a.y>f)&&(a.x>f&&(s.x=Math.floor(f/at.x),a.x=s.x*at.x,X.mapSize.x=s.x),a.y>f&&(s.y=Math.floor(f/at.y),a.y=s.y*at.y,X.mapSize.y=s.y));let et=e.state.buffers.depth.getReversed();if(X.camera._reversedDepth=et,X.map===null||P===!0){if(X.map!==null&&(X.map.depthTexture!==null&&(X.map.depthTexture.dispose(),X.map.depthTexture=null),X.map.dispose()),this.type===dl){if(V.isPointLight){Ot("WebGLShadowMap: VSM shadow maps are not supported for PointLights. Use PCF or BasicShadowMap instead.");continue}X.map=new gi(a.x,a.y,{format:Lr,type:_a,minFilter:En,magFilter:En,generateMipmaps:!1}),X.map.texture.name=V.name+".shadowMap",X.map.depthTexture=new Fs(a.x,a.y,Ni),X.map.depthTexture.name=V.name+".shadowMapDepth",X.map.depthTexture.format=da,X.map.depthTexture.compareFunction=null,X.map.depthTexture.minFilter=pn,X.map.depthTexture.magFilter=pn}else V.isPointLight?(X.map=new fp(a.x),X.map.depthTexture=new ed(a.x,qi)):(X.map=new gi(a.x,a.y),X.map.depthTexture=new Fs(a.x,a.y,qi)),X.map.depthTexture.name=V.name+".shadowMap",X.map.depthTexture.format=da,this.type===fu?(X.map.depthTexture.compareFunction=et?op:rp,X.map.depthTexture.minFilter=En,X.map.depthTexture.magFilter=En):(X.map.depthTexture.compareFunction=null,X.map.depthTexture.minFilter=pn,X.map.depthTexture.magFilter=pn);X.camera.updateProjectionMatrix()}let dt=X.map.isWebGLCubeRenderTarget?6:1;for(let St=0;St<dt;St++){if(X.map.isWebGLCubeRenderTarget)e.setRenderTarget(X.map,St),e.clear();else{St===0&&(e.setRenderTarget(X.map),e.clear());let ht=X.getViewport(St);r.set(s.x*ht.x,s.y*ht.y,s.x*ht.z,s.y*ht.w),N.viewport(r)}if(V.isPointLight){let ht=X.camera,Vt=X.matrix,Ce=V.distance||ht.far;Ce!==ht.far&&(ht.far=Ce,ht.updateProjectionMatrix()),xu.setFromMatrixPosition(V.matrixWorld),ht.position.copy(xu),Nv.copy(ht.position),Nv.add(hN[St]),ht.up.copy(dN[St]),ht.lookAt(Nv),ht.updateMatrixWorld(),Vt.makeTranslation(-xu.x,-xu.y,-xu.z),q1.multiplyMatrices(ht.projectionMatrix,ht.matrixWorldInverse),X._frustum.setFromProjectionMatrix(q1,ht.coordinateSystem,ht.reversedDepth)}else X.updateMatrices(V);i=X.getFrustum(),x(w,y,X.camera,V,this.type)}X.isPointLightShadow!==!0&&this.type===dl&&m(X,y),X.needsUpdate=!1}h=this.type,v.needsUpdate=!1,e.setRenderTarget(T,L,R)};function m(A,w){let y=t.update(b);u.defines.VSM_SAMPLES!==A.blurSamples&&(u.defines.VSM_SAMPLES=A.blurSamples,d.defines.VSM_SAMPLES=A.blurSamples,u.needsUpdate=!0,d.needsUpdate=!0),A.mapPass===null&&(A.mapPass=new gi(a.x,a.y,{format:Lr,type:_a})),u.uniforms.shadow_pass.value=A.map.depthTexture,u.uniforms.resolution.value=A.mapSize,u.uniforms.radius.value=A.radius,e.setRenderTarget(A.mapPass),e.clear(),e.renderBufferDirect(w,null,y,u,b,null),d.uniforms.shadow_pass.value=A.mapPass.texture,d.uniforms.resolution.value=A.mapSize,d.uniforms.radius.value=A.radius,e.setRenderTarget(A.map),e.clear(),e.renderBufferDirect(w,null,y,d,b,null)}function g(A,w,y,T){let L=null,R=y.isPointLight===!0?A.customDistanceMaterial:A.customDepthMaterial;if(R!==void 0)L=R;else if(L=y.isPointLight===!0?l:o,e.localClippingEnabled&&w.clipShadows===!0&&Array.isArray(w.clippingPlanes)&&w.clippingPlanes.length!==0||w.displacementMap&&w.displacementScale!==0||w.alphaMap&&w.alphaTest>0||w.map&&w.alphaTest>0||w.alphaToCoverage===!0){let N=L.uuid,P=w.uuid,z=c[N];z===void 0&&(z={},c[N]=z);let F=z[P];F===void 0&&(F=L.clone(),z[P]=F,w.addEventListener("dispose",E)),L=F}if(L.visible=w.visible,L.wireframe=w.wireframe,T===dl?L.side=w.shadowSide!==null?w.shadowSide:w.side:L.side=w.shadowSide!==null?w.shadowSide:p[w.side],L.alphaMap=w.alphaMap,L.alphaTest=w.alphaToCoverage===!0?.5:w.alphaTest,L.map=w.map,L.clipShadows=w.clipShadows,L.clippingPlanes=w.clippingPlanes,L.clipIntersection=w.clipIntersection,L.displacementMap=w.displacementMap,L.displacementScale=w.displacementScale,L.displacementBias=w.displacementBias,L.wireframeLinewidth=w.wireframeLinewidth,L.linewidth=w.linewidth,y.isPointLight===!0&&L.isMeshDistanceMaterial===!0){let N=e.properties.get(L);N.light=y}return L}function x(A,w,y,T,L){if(A.visible===!1)return;if(A.layers.test(w.layers)&&(A.isMesh||A.isLine||A.isPoints)&&(A.castShadow||A.receiveShadow&&L===dl)&&(!A.frustumCulled||i.intersectsObject(A))){A.modelViewMatrix.multiplyMatrices(y.matrixWorldInverse,A.matrixWorld);let P=t.update(A),z=A.material;if(Array.isArray(z)){let F=P.groups;for(let V=0,X=F.length;V<X;V++){let at=F[V],et=z[at.materialIndex];if(et&&et.visible){let dt=g(A,et,T,L);A.onBeforeShadow(e,A,w,y,P,dt,at),e.renderBufferDirect(y,null,P,dt,A,at),A.onAfterShadow(e,A,w,y,P,dt,at)}}}else if(z.visible){let F=g(A,z,T,L);A.onBeforeShadow(e,A,w,y,P,F,null),e.renderBufferDirect(y,null,P,F,A,null),A.onAfterShadow(e,A,w,y,P,F,null)}}let N=A.children;for(let P=0,z=N.length;P<z;P++)x(N[P],w,y,T,L)}function E(A){A.target.removeEventListener("dispose",E);for(let y in c){let T=c[y],L=A.target.uuid;L in T&&(T[L].dispose(),delete T[L])}}}function mN(e,t){function n(){let U=!1,ut=new Ge,ot=null,xt=new Ge(0,0,0,0);return{setMask:function(nt){ot!==nt&&!U&&(e.colorMask(nt,nt,nt,nt),ot=nt)},setLocked:function(nt){U=nt},setClear:function(nt,q,Et,Ht,ge){ge===!0&&(nt*=Ht,q*=Ht,Et*=Ht),ut.set(nt,q,Et,Ht),xt.equals(ut)===!1&&(e.clearColor(nt,q,Et,Ht),xt.copy(ut))},reset:function(){U=!1,ot=null,xt.set(-1,0,0,0)}}}function i(){let U=!1,ut=!1,ot=null,xt=null,nt=null;return{setReversed:function(q){if(ut!==q){let Et=t.get("EXT_clip_control");q?Et.clipControlEXT(Et.LOWER_LEFT_EXT,Et.ZERO_TO_ONE_EXT):Et.clipControlEXT(Et.LOWER_LEFT_EXT,Et.NEGATIVE_ONE_TO_ONE_EXT),ut=q;let Ht=nt;nt=null,this.setClear(Ht)}},getReversed:function(){return ut},setTest:function(q){q?j(e.DEPTH_TEST):ct(e.DEPTH_TEST)},setMask:function(q){ot!==q&&!U&&(e.depthMask(q),ot=q)},setFunc:function(q){if(ut&&(q=S1[q]),xt!==q){switch(q){case Fh:e.depthFunc(e.NEVER);break;case Vh:e.depthFunc(e.ALWAYS);break;case Hh:e.depthFunc(e.LESS);break;case Cr:e.depthFunc(e.LEQUAL);break;case Gh:e.depthFunc(e.EQUAL);break;case kh:e.depthFunc(e.GEQUAL);break;case Xh:e.depthFunc(e.GREATER);break;case Wh:e.depthFunc(e.NOTEQUAL);break;default:e.depthFunc(e.LEQUAL)}xt=q}},setLocked:function(q){U=q},setClear:function(q){nt!==q&&(nt=q,ut&&(q=1-q),e.clearDepth(q))},reset:function(){U=!1,ot=null,xt=null,nt=null,ut=!1}}}function a(){let U=!1,ut=null,ot=null,xt=null,nt=null,q=null,Et=null,Ht=null,ge=null;return{setTest:function(oe){U||(oe?j(e.STENCIL_TEST):ct(e.STENCIL_TEST))},setMask:function(oe){ut!==oe&&!U&&(e.stencilMask(oe),ut=oe)},setFunc:function(oe,Fn,yi){(ot!==oe||xt!==Fn||nt!==yi)&&(e.stencilFunc(oe,Fn,yi),ot=oe,xt=Fn,nt=yi)},setOp:function(oe,Fn,yi){(q!==oe||Et!==Fn||Ht!==yi)&&(e.stencilOp(oe,Fn,yi),q=oe,Et=Fn,Ht=yi)},setLocked:function(oe){U=oe},setClear:function(oe){ge!==oe&&(e.clearStencil(oe),ge=oe)},reset:function(){U=!1,ut=null,ot=null,xt=null,nt=null,q=null,Et=null,Ht=null,ge=null}}}let s=new n,r=new i,o=new a,l=new WeakMap,c=new WeakMap,f={},p={},u=new WeakMap,d=[],_=null,b=!1,v=null,h=null,m=null,g=null,x=null,E=null,A=null,w=new Ft(0,0,0),y=0,T=!1,L=null,R=null,N=null,P=null,z=null,F=e.getParameter(e.MAX_COMBINED_TEXTURE_IMAGE_UNITS),V=!1,X=0,at=e.getParameter(e.VERSION);at.indexOf("WebGL")!==-1?(X=parseFloat(/^WebGL (\d)/.exec(at)[1]),V=X>=1):at.indexOf("OpenGL ES")!==-1&&(X=parseFloat(/^OpenGL ES (\d)/.exec(at)[1]),V=X>=2);let et=null,dt={},St=e.getParameter(e.SCISSOR_BOX),ht=e.getParameter(e.VIEWPORT),Vt=new Ge().fromArray(St),Ce=new Ge().fromArray(ht);function xe(U,ut,ot,xt){let nt=new Uint8Array(4),q=e.createTexture();e.bindTexture(U,q),e.texParameteri(U,e.TEXTURE_MIN_FILTER,e.NEAREST),e.texParameteri(U,e.TEXTURE_MAG_FILTER,e.NEAREST);for(let Et=0;Et<ot;Et++)U===e.TEXTURE_3D||U===e.TEXTURE_2D_ARRAY?e.texImage3D(ut,0,e.RGBA,1,1,xt,0,e.RGBA,e.UNSIGNED_BYTE,nt):e.texImage2D(ut+Et,0,e.RGBA,1,1,0,e.RGBA,e.UNSIGNED_BYTE,nt);return q}let K={};K[e.TEXTURE_2D]=xe(e.TEXTURE_2D,e.TEXTURE_2D,1),K[e.TEXTURE_CUBE_MAP]=xe(e.TEXTURE_CUBE_MAP,e.TEXTURE_CUBE_MAP_POSITIVE_X,6),K[e.TEXTURE_2D_ARRAY]=xe(e.TEXTURE_2D_ARRAY,e.TEXTURE_2D_ARRAY,1,1),K[e.TEXTURE_3D]=xe(e.TEXTURE_3D,e.TEXTURE_3D,1,1),s.setClear(0,0,0,1),r.setClear(1),o.setClear(0),j(e.DEPTH_TEST),r.setFunc(Cr),Xt(!1),Te(tv),j(e.CULL_FACE),Jt(va);function j(U){f[U]!==!0&&(e.enable(U),f[U]=!0)}function ct(U){f[U]!==!1&&(e.disable(U),f[U]=!1)}function It(U,ut){return p[U]!==ut?(e.bindFramebuffer(U,ut),p[U]=ut,U===e.DRAW_FRAMEBUFFER&&(p[e.FRAMEBUFFER]=ut),U===e.FRAMEBUFFER&&(p[e.DRAW_FRAMEBUFFER]=ut),!0):!1}function Lt(U,ut){let ot=d,xt=!1;if(U){ot=u.get(ut),ot===void 0&&(ot=[],u.set(ut,ot));let nt=U.textures;if(ot.length!==nt.length||ot[0]!==e.COLOR_ATTACHMENT0){for(let q=0,Et=nt.length;q<Et;q++)ot[q]=e.COLOR_ATTACHMENT0+q;ot.length=nt.length,xt=!0}}else ot[0]!==e.BACK&&(ot[0]=e.BACK,xt=!0);xt&&e.drawBuffers(ot)}function Pt(U){return _!==U?(e.useProgram(U),_=U,!0):!1}let Ye={[zs]:e.FUNC_ADD,[XM]:e.FUNC_SUBTRACT,[WM]:e.FUNC_REVERSE_SUBTRACT};Ye[qM]=e.MIN,Ye[YM]=e.MAX;let ie={[ZM]:e.ZERO,[JM]:e.ONE,[KM]:e.SRC_COLOR,[Bh]:e.SRC_ALPHA,[n1]:e.SRC_ALPHA_SATURATE,[t1]:e.DST_COLOR,[jM]:e.DST_ALPHA,[QM]:e.ONE_MINUS_SRC_COLOR,[zh]:e.ONE_MINUS_SRC_ALPHA,[e1]:e.ONE_MINUS_DST_COLOR,[$M]:e.ONE_MINUS_DST_ALPHA,[i1]:e.CONSTANT_COLOR,[a1]:e.ONE_MINUS_CONSTANT_COLOR,[s1]:e.CONSTANT_ALPHA,[r1]:e.ONE_MINUS_CONSTANT_ALPHA};function Jt(U,ut,ot,xt,nt,q,Et,Ht,ge,oe){if(U===va){b===!0&&(ct(e.BLEND),b=!1);return}if(b===!1&&(j(e.BLEND),b=!0),U!==kM){if(U!==v||oe!==T){if((h!==zs||x!==zs)&&(e.blendEquation(e.FUNC_ADD),h=zs,x=zs),oe)switch(U){case wr:e.blendFuncSeparate(e.ONE,e.ONE_MINUS_SRC_ALPHA,e.ONE,e.ONE_MINUS_SRC_ALPHA);break;case ev:e.blendFunc(e.ONE,e.ONE);break;case nv:e.blendFuncSeparate(e.ZERO,e.ONE_MINUS_SRC_COLOR,e.ZERO,e.ONE);break;case iv:e.blendFuncSeparate(e.DST_COLOR,e.ONE_MINUS_SRC_ALPHA,e.ZERO,e.ONE);break;default:zt("WebGLState: Invalid blending: ",U);break}else switch(U){case wr:e.blendFuncSeparate(e.SRC_ALPHA,e.ONE_MINUS_SRC_ALPHA,e.ONE,e.ONE_MINUS_SRC_ALPHA);break;case ev:e.blendFuncSeparate(e.SRC_ALPHA,e.ONE,e.ONE,e.ONE);break;case nv:zt("WebGLState: SubtractiveBlending requires material.premultipliedAlpha = true");break;case iv:zt("WebGLState: MultiplyBlending requires material.premultipliedAlpha = true");break;default:zt("WebGLState: Invalid blending: ",U);break}m=null,g=null,E=null,A=null,w.set(0,0,0),y=0,v=U,T=oe}return}nt=nt||ut,q=q||ot,Et=Et||xt,(ut!==h||nt!==x)&&(e.blendEquationSeparate(Ye[ut],Ye[nt]),h=ut,x=nt),(ot!==m||xt!==g||q!==E||Et!==A)&&(e.blendFuncSeparate(ie[ot],ie[xt],ie[q],ie[Et]),m=ot,g=xt,E=q,A=Et),(Ht.equals(w)===!1||ge!==y)&&(e.blendColor(Ht.r,Ht.g,Ht.b,ge),w.copy(Ht),y=ge),v=U,T=!1}function de(U,ut){U.side===ga?ct(e.CULL_FACE):j(e.CULL_FACE);let ot=U.side===wn;ut&&(ot=!ot),Xt(ot),U.blending===wr&&U.transparent===!1?Jt(va):Jt(U.blending,U.blendEquation,U.blendSrc,U.blendDst,U.blendEquationAlpha,U.blendSrcAlpha,U.blendDstAlpha,U.blendColor,U.blendAlpha,U.premultipliedAlpha),r.setFunc(U.depthFunc),r.setTest(U.depthTest),r.setMask(U.depthWrite),s.setMask(U.colorWrite);let xt=U.stencilWrite;o.setTest(xt),xt&&(o.setMask(U.stencilWriteMask),o.setFunc(U.stencilFunc,U.stencilRef,U.stencilFuncMask),o.setOp(U.stencilFail,U.stencilZFail,U.stencilZPass)),ze(U.polygonOffset,U.polygonOffsetFactor,U.polygonOffsetUnits),U.alphaToCoverage===!0?j(e.SAMPLE_ALPHA_TO_COVERAGE):ct(e.SAMPLE_ALPHA_TO_COVERAGE)}function Xt(U){L!==U&&(U?e.frontFace(e.CW):e.frontFace(e.CCW),L=U)}function Te(U){U!==VM?(j(e.CULL_FACE),U!==R&&(U===tv?e.cullFace(e.BACK):U===HM?e.cullFace(e.FRONT):e.cullFace(e.FRONT_AND_BACK))):ct(e.CULL_FACE),R=U}function D(U){U!==N&&(V&&e.lineWidth(U),N=U)}function ze(U,ut,ot){U?(j(e.POLYGON_OFFSET_FILL),(P!==ut||z!==ot)&&(P=ut,z=ot,r.getReversed()&&(ut=-ut),e.polygonOffset(ut,ot))):ct(e.POLYGON_OFFSET_FILL)}function re(U){U?j(e.SCISSOR_TEST):ct(e.SCISSOR_TEST)}function ue(U){U===void 0&&(U=e.TEXTURE0+F-1),et!==U&&(e.activeTexture(U),et=U)}function At(U,ut,ot){ot===void 0&&(et===null?ot=e.TEXTURE0+F-1:ot=et);let xt=dt[ot];xt===void 0&&(xt={type:void 0,texture:void 0},dt[ot]=xt),(xt.type!==U||xt.texture!==ut)&&(et!==ot&&(e.activeTexture(ot),et=ot),e.bindTexture(U,ut||K[U]),xt.type=U,xt.texture=ut)}function C(){let U=dt[et];U!==void 0&&U.type!==void 0&&(e.bindTexture(U.type,null),U.type=void 0,U.texture=void 0)}function S(){try{e.compressedTexImage2D(...arguments)}catch(U){zt("WebGLState:",U)}}function O(){try{e.compressedTexImage3D(...arguments)}catch(U){zt("WebGLState:",U)}}function Z(){try{e.texSubImage2D(...arguments)}catch(U){zt("WebGLState:",U)}}function Q(){try{e.texSubImage3D(...arguments)}catch(U){zt("WebGLState:",U)}}function Y(){try{e.compressedTexSubImage2D(...arguments)}catch(U){zt("WebGLState:",U)}}function Tt(){try{e.compressedTexSubImage3D(...arguments)}catch(U){zt("WebGLState:",U)}}function lt(){try{e.texStorage2D(...arguments)}catch(U){zt("WebGLState:",U)}}function Rt(){try{e.texStorage3D(...arguments)}catch(U){zt("WebGLState:",U)}}function Ut(){try{e.texImage2D(...arguments)}catch(U){zt("WebGLState:",U)}}function tt(){try{e.texImage3D(...arguments)}catch(U){zt("WebGLState:",U)}}function rt(U){Vt.equals(U)===!1&&(e.scissor(U.x,U.y,U.z,U.w),Vt.copy(U))}function Mt(U){Ce.equals(U)===!1&&(e.viewport(U.x,U.y,U.z,U.w),Ce.copy(U))}function bt(U,ut){let ot=c.get(ut);ot===void 0&&(ot=new WeakMap,c.set(ut,ot));let xt=ot.get(U);xt===void 0&&(xt=e.getUniformBlockIndex(ut,U.name),ot.set(U,xt))}function gt(U,ut){let xt=c.get(ut).get(U);l.get(ut)!==xt&&(e.uniformBlockBinding(ut,xt,U.__bindingPointIndex),l.set(ut,xt))}function Wt(){e.disable(e.BLEND),e.disable(e.CULL_FACE),e.disable(e.DEPTH_TEST),e.disable(e.POLYGON_OFFSET_FILL),e.disable(e.SCISSOR_TEST),e.disable(e.STENCIL_TEST),e.disable(e.SAMPLE_ALPHA_TO_COVERAGE),e.blendEquation(e.FUNC_ADD),e.blendFunc(e.ONE,e.ZERO),e.blendFuncSeparate(e.ONE,e.ZERO,e.ONE,e.ZERO),e.blendColor(0,0,0,0),e.colorMask(!0,!0,!0,!0),e.clearColor(0,0,0,0),e.depthMask(!0),e.depthFunc(e.LESS),r.setReversed(!1),e.clearDepth(1),e.stencilMask(4294967295),e.stencilFunc(e.ALWAYS,0,4294967295),e.stencilOp(e.KEEP,e.KEEP,e.KEEP),e.clearStencil(0),e.cullFace(e.BACK),e.frontFace(e.CCW),e.polygonOffset(0,0),e.activeTexture(e.TEXTURE0),e.bindFramebuffer(e.FRAMEBUFFER,null),e.bindFramebuffer(e.DRAW_FRAMEBUFFER,null),e.bindFramebuffer(e.READ_FRAMEBUFFER,null),e.useProgram(null),e.lineWidth(1),e.scissor(0,0,e.canvas.width,e.canvas.height),e.viewport(0,0,e.canvas.width,e.canvas.height),f={},et=null,dt={},p={},u=new WeakMap,d=[],_=null,b=!1,v=null,h=null,m=null,g=null,x=null,E=null,A=null,w=new Ft(0,0,0),y=0,T=!1,L=null,R=null,N=null,P=null,z=null,Vt.set(0,0,e.canvas.width,e.canvas.height),Ce.set(0,0,e.canvas.width,e.canvas.height),s.reset(),r.reset(),o.reset()}return{buffers:{color:s,depth:r,stencil:o},enable:j,disable:ct,bindFramebuffer:It,drawBuffers:Lt,useProgram:Pt,setBlending:Jt,setMaterial:de,setFlipSided:Xt,setCullFace:Te,setLineWidth:D,setPolygonOffset:ze,setScissorTest:re,activeTexture:ue,bindTexture:At,unbindTexture:C,compressedTexImage2D:S,compressedTexImage3D:O,texImage2D:Ut,texImage3D:tt,updateUBOMapping:bt,uniformBlockBinding:gt,texStorage2D:lt,texStorage3D:Rt,texSubImage2D:Z,texSubImage3D:Q,compressedTexSubImage2D:Y,compressedTexSubImage3D:Tt,scissor:rt,viewport:Mt,reset:Wt}}function gN(e,t,n,i,a,s,r){let o=t.has("WEBGL_multisampled_render_to_texture")?t.get("WEBGL_multisampled_render_to_texture"):null,l=typeof navigator>"u"?!1:/OculusBrowser/g.test(navigator.userAgent),c=new te,f=new WeakMap,p,u=new WeakMap,d=!1;try{d=typeof OffscreenCanvas<"u"&&new OffscreenCanvas(1,1).getContext("2d")!==null}catch{}function _(C,S){return d?new OffscreenCanvas(C,S):Hc("canvas")}function b(C,S,O){let Z=1,Q=At(C);if((Q.width>O||Q.height>O)&&(Z=O/Math.max(Q.width,Q.height)),Z<1)if(typeof HTMLImageElement<"u"&&C instanceof HTMLImageElement||typeof HTMLCanvasElement<"u"&&C instanceof HTMLCanvasElement||typeof ImageBitmap<"u"&&C instanceof ImageBitmap||typeof VideoFrame<"u"&&C instanceof VideoFrame){let Y=Math.floor(Z*Q.width),Tt=Math.floor(Z*Q.height);p===void 0&&(p=_(Y,Tt));let lt=S?_(Y,Tt):p;return lt.width=Y,lt.height=Tt,lt.getContext("2d").drawImage(C,0,0,Y,Tt),Ot("WebGLRenderer: Texture has been resized from ("+Q.width+"x"+Q.height+") to ("+Y+"x"+Tt+")."),lt}else return"data"in C&&Ot("WebGLRenderer: Image in DataTexture is too big ("+Q.width+"x"+Q.height+")."),C;return C}function v(C){return C.generateMipmaps}function h(C){e.generateMipmap(C)}function m(C){return C.isWebGLCubeRenderTarget?e.TEXTURE_CUBE_MAP:C.isWebGL3DRenderTarget?e.TEXTURE_3D:C.isWebGLArrayRenderTarget||C.isCompressedArrayTexture?e.TEXTURE_2D_ARRAY:e.TEXTURE_2D}function g(C,S,O,Z,Q=!1){if(C!==null){if(e[C]!==void 0)return e[C];Ot("WebGLRenderer: Attempt to use non-existing WebGL internal format '"+C+"'")}let Y=S;if(S===e.RED&&(O===e.FLOAT&&(Y=e.R32F),O===e.HALF_FLOAT&&(Y=e.R16F),O===e.UNSIGNED_BYTE&&(Y=e.R8)),S===e.RED_INTEGER&&(O===e.UNSIGNED_BYTE&&(Y=e.R8UI),O===e.UNSIGNED_SHORT&&(Y=e.R16UI),O===e.UNSIGNED_INT&&(Y=e.R32UI),O===e.BYTE&&(Y=e.R8I),O===e.SHORT&&(Y=e.R16I),O===e.INT&&(Y=e.R32I)),S===e.RG&&(O===e.FLOAT&&(Y=e.RG32F),O===e.HALF_FLOAT&&(Y=e.RG16F),O===e.UNSIGNED_BYTE&&(Y=e.RG8)),S===e.RG_INTEGER&&(O===e.UNSIGNED_BYTE&&(Y=e.RG8UI),O===e.UNSIGNED_SHORT&&(Y=e.RG16UI),O===e.UNSIGNED_INT&&(Y=e.RG32UI),O===e.BYTE&&(Y=e.RG8I),O===e.SHORT&&(Y=e.RG16I),O===e.INT&&(Y=e.RG32I)),S===e.RGB_INTEGER&&(O===e.UNSIGNED_BYTE&&(Y=e.RGB8UI),O===e.UNSIGNED_SHORT&&(Y=e.RGB16UI),O===e.UNSIGNED_INT&&(Y=e.RGB32UI),O===e.BYTE&&(Y=e.RGB8I),O===e.SHORT&&(Y=e.RGB16I),O===e.INT&&(Y=e.RGB32I)),S===e.RGBA_INTEGER&&(O===e.UNSIGNED_BYTE&&(Y=e.RGBA8UI),O===e.UNSIGNED_SHORT&&(Y=e.RGBA16UI),O===e.UNSIGNED_INT&&(Y=e.RGBA32UI),O===e.BYTE&&(Y=e.RGBA8I),O===e.SHORT&&(Y=e.RGBA16I),O===e.INT&&(Y=e.RGBA32I)),S===e.RGB&&(O===e.UNSIGNED_INT_5_9_9_9_REV&&(Y=e.RGB9_E5),O===e.UNSIGNED_INT_10F_11F_11F_REV&&(Y=e.R11F_G11F_B10F)),S===e.RGBA){let Tt=Q?Vc:se.getTransfer(Z);O===e.FLOAT&&(Y=e.RGBA32F),O===e.HALF_FLOAT&&(Y=e.RGBA16F),O===e.UNSIGNED_BYTE&&(Y=Tt===me?e.SRGB8_ALPHA8:e.RGBA8),O===e.UNSIGNED_SHORT_4_4_4_4&&(Y=e.RGBA4),O===e.UNSIGNED_SHORT_5_5_5_1&&(Y=e.RGB5_A1)}return(Y===e.R16F||Y===e.R32F||Y===e.RG16F||Y===e.RG32F||Y===e.RGBA16F||Y===e.RGBA32F)&&t.get("EXT_color_buffer_float"),Y}function x(C,S){let O;return C?S===null||S===qi||S===ml?O=e.DEPTH24_STENCIL8:S===Ni?O=e.DEPTH32F_STENCIL8:S===pl&&(O=e.DEPTH24_STENCIL8,Ot("DepthTexture: 16 bit depth attachment is not supported with stencil. Using 24-bit attachment.")):S===null||S===qi||S===ml?O=e.DEPTH_COMPONENT24:S===Ni?O=e.DEPTH_COMPONENT32F:S===pl&&(O=e.DEPTH_COMPONENT16),O}function E(C,S){return v(C)===!0||C.isFramebufferTexture&&C.minFilter!==pn&&C.minFilter!==En?Math.log2(Math.max(S.width,S.height))+1:C.mipmaps!==void 0&&C.mipmaps.length>0?C.mipmaps.length:C.isCompressedTexture&&Array.isArray(C.image)?S.mipmaps.length:1}function A(C){let S=C.target;S.removeEventListener("dispose",A),y(S),S.isVideoTexture&&f.delete(S)}function w(C){let S=C.target;S.removeEventListener("dispose",w),L(S)}function y(C){let S=i.get(C);if(S.__webglInit===void 0)return;let O=C.source,Z=u.get(O);if(Z){let Q=Z[S.__cacheKey];Q.usedTimes--,Q.usedTimes===0&&T(C),Object.keys(Z).length===0&&u.delete(O)}i.remove(C)}function T(C){let S=i.get(C);e.deleteTexture(S.__webglTexture);let O=C.source,Z=u.get(O);delete Z[S.__cacheKey],r.memory.textures--}function L(C){let S=i.get(C);if(C.depthTexture&&(C.depthTexture.dispose(),i.remove(C.depthTexture)),C.isWebGLCubeRenderTarget)for(let Z=0;Z<6;Z++){if(Array.isArray(S.__webglFramebuffer[Z]))for(let Q=0;Q<S.__webglFramebuffer[Z].length;Q++)e.deleteFramebuffer(S.__webglFramebuffer[Z][Q]);else e.deleteFramebuffer(S.__webglFramebuffer[Z]);S.__webglDepthbuffer&&e.deleteRenderbuffer(S.__webglDepthbuffer[Z])}else{if(Array.isArray(S.__webglFramebuffer))for(let Z=0;Z<S.__webglFramebuffer.length;Z++)e.deleteFramebuffer(S.__webglFramebuffer[Z]);else e.deleteFramebuffer(S.__webglFramebuffer);if(S.__webglDepthbuffer&&e.deleteRenderbuffer(S.__webglDepthbuffer),S.__webglMultisampledFramebuffer&&e.deleteFramebuffer(S.__webglMultisampledFramebuffer),S.__webglColorRenderbuffer)for(let Z=0;Z<S.__webglColorRenderbuffer.length;Z++)S.__webglColorRenderbuffer[Z]&&e.deleteRenderbuffer(S.__webglColorRenderbuffer[Z]);S.__webglDepthRenderbuffer&&e.deleteRenderbuffer(S.__webglDepthRenderbuffer)}let O=C.textures;for(let Z=0,Q=O.length;Z<Q;Z++){let Y=i.get(O[Z]);Y.__webglTexture&&(e.deleteTexture(Y.__webglTexture),r.memory.textures--),i.remove(O[Z])}i.remove(C)}let R=0;function N(){R=0}function P(){let C=R;return C>=a.maxTextures&&Ot("WebGLTextures: Trying to use "+C+" texture units while this GPU supports only "+a.maxTextures),R+=1,C}function z(C){let S=[];return S.push(C.wrapS),S.push(C.wrapT),S.push(C.wrapR||0),S.push(C.magFilter),S.push(C.minFilter),S.push(C.anisotropy),S.push(C.internalFormat),S.push(C.format),S.push(C.type),S.push(C.generateMipmaps),S.push(C.premultiplyAlpha),S.push(C.flipY),S.push(C.unpackAlignment),S.push(C.colorSpace),S.join()}function F(C,S){let O=i.get(C);if(C.isVideoTexture&&re(C),C.isRenderTargetTexture===!1&&C.isExternalTexture!==!0&&C.version>0&&O.__version!==C.version){let Z=C.image;if(Z===null)Ot("WebGLRenderer: Texture marked for update but no image data found.");else if(Z.complete===!1)Ot("WebGLRenderer: Texture marked for update but image is incomplete");else{K(O,C,S);return}}else C.isExternalTexture&&(O.__webglTexture=C.sourceTexture?C.sourceTexture:null);n.bindTexture(e.TEXTURE_2D,O.__webglTexture,e.TEXTURE0+S)}function V(C,S){let O=i.get(C);if(C.isRenderTargetTexture===!1&&C.version>0&&O.__version!==C.version){K(O,C,S);return}else C.isExternalTexture&&(O.__webglTexture=C.sourceTexture?C.sourceTexture:null);n.bindTexture(e.TEXTURE_2D_ARRAY,O.__webglTexture,e.TEXTURE0+S)}function X(C,S){let O=i.get(C);if(C.isRenderTargetTexture===!1&&C.version>0&&O.__version!==C.version){K(O,C,S);return}n.bindTexture(e.TEXTURE_3D,O.__webglTexture,e.TEXTURE0+S)}function at(C,S){let O=i.get(C);if(C.isCubeDepthTexture!==!0&&C.version>0&&O.__version!==C.version){j(O,C,S);return}n.bindTexture(e.TEXTURE_CUBE_MAP,O.__webglTexture,e.TEXTURE0+S)}let et={[qh]:e.REPEAT,[ha]:e.CLAMP_TO_EDGE,[Yh]:e.MIRRORED_REPEAT},dt={[pn]:e.NEAREST,[c1]:e.NEAREST_MIPMAP_NEAREST,[pu]:e.NEAREST_MIPMAP_LINEAR,[En]:e.LINEAR,[xd]:e.LINEAR_MIPMAP_NEAREST,[Ws]:e.LINEAR_MIPMAP_LINEAR},St={[h1]:e.NEVER,[v1]:e.ALWAYS,[d1]:e.LESS,[rp]:e.LEQUAL,[p1]:e.EQUAL,[op]:e.GEQUAL,[m1]:e.GREATER,[g1]:e.NOTEQUAL};function ht(C,S){if(S.type===Ni&&t.has("OES_texture_float_linear")===!1&&(S.magFilter===En||S.magFilter===xd||S.magFilter===pu||S.magFilter===Ws||S.minFilter===En||S.minFilter===xd||S.minFilter===pu||S.minFilter===Ws)&&Ot("WebGLRenderer: Unable to use linear filtering with floating point textures. OES_texture_float_linear not supported on this device."),e.texParameteri(C,e.TEXTURE_WRAP_S,et[S.wrapS]),e.texParameteri(C,e.TEXTURE_WRAP_T,et[S.wrapT]),(C===e.TEXTURE_3D||C===e.TEXTURE_2D_ARRAY)&&e.texParameteri(C,e.TEXTURE_WRAP_R,et[S.wrapR]),e.texParameteri(C,e.TEXTURE_MAG_FILTER,dt[S.magFilter]),e.texParameteri(C,e.TEXTURE_MIN_FILTER,dt[S.minFilter]),S.compareFunction&&(e.texParameteri(C,e.TEXTURE_COMPARE_MODE,e.COMPARE_REF_TO_TEXTURE),e.texParameteri(C,e.TEXTURE_COMPARE_FUNC,St[S.compareFunction])),t.has("EXT_texture_filter_anisotropic")===!0){if(S.magFilter===pn||S.minFilter!==pu&&S.minFilter!==Ws||S.type===Ni&&t.has("OES_texture_float_linear")===!1)return;if(S.anisotropy>1||i.get(S).__currentAnisotropy){let O=t.get("EXT_texture_filter_anisotropic");e.texParameterf(C,O.TEXTURE_MAX_ANISOTROPY_EXT,Math.min(S.anisotropy,a.getMaxAnisotropy())),i.get(S).__currentAnisotropy=S.anisotropy}}}function Vt(C,S){let O=!1;C.__webglInit===void 0&&(C.__webglInit=!0,S.addEventListener("dispose",A));let Z=S.source,Q=u.get(Z);Q===void 0&&(Q={},u.set(Z,Q));let Y=z(S);if(Y!==C.__cacheKey){Q[Y]===void 0&&(Q[Y]={texture:e.createTexture(),usedTimes:0},r.memory.textures++,O=!0),Q[Y].usedTimes++;let Tt=Q[C.__cacheKey];Tt!==void 0&&(Q[C.__cacheKey].usedTimes--,Tt.usedTimes===0&&T(S)),C.__cacheKey=Y,C.__webglTexture=Q[Y].texture}return O}function Ce(C,S,O){return Math.floor(Math.floor(C/O)/S)}function xe(C,S,O,Z){let Y=C.updateRanges;if(Y.length===0)n.texSubImage2D(e.TEXTURE_2D,0,0,0,S.width,S.height,O,Z,S.data);else{Y.sort((tt,rt)=>tt.start-rt.start);let Tt=0;for(let tt=1;tt<Y.length;tt++){let rt=Y[Tt],Mt=Y[tt],bt=rt.start+rt.count,gt=Ce(Mt.start,S.width,4),Wt=Ce(rt.start,S.width,4);Mt.start<=bt+1&&gt===Wt&&Ce(Mt.start+Mt.count-1,S.width,4)===gt?rt.count=Math.max(rt.count,Mt.start+Mt.count-rt.start):(++Tt,Y[Tt]=Mt)}Y.length=Tt+1;let lt=e.getParameter(e.UNPACK_ROW_LENGTH),Rt=e.getParameter(e.UNPACK_SKIP_PIXELS),Ut=e.getParameter(e.UNPACK_SKIP_ROWS);e.pixelStorei(e.UNPACK_ROW_LENGTH,S.width);for(let tt=0,rt=Y.length;tt<rt;tt++){let Mt=Y[tt],bt=Math.floor(Mt.start/4),gt=Math.ceil(Mt.count/4),Wt=bt%S.width,U=Math.floor(bt/S.width),ut=gt,ot=1;e.pixelStorei(e.UNPACK_SKIP_PIXELS,Wt),e.pixelStorei(e.UNPACK_SKIP_ROWS,U),n.texSubImage2D(e.TEXTURE_2D,0,Wt,U,ut,ot,O,Z,S.data)}C.clearUpdateRanges(),e.pixelStorei(e.UNPACK_ROW_LENGTH,lt),e.pixelStorei(e.UNPACK_SKIP_PIXELS,Rt),e.pixelStorei(e.UNPACK_SKIP_ROWS,Ut)}}function K(C,S,O){let Z=e.TEXTURE_2D;(S.isDataArrayTexture||S.isCompressedArrayTexture)&&(Z=e.TEXTURE_2D_ARRAY),S.isData3DTexture&&(Z=e.TEXTURE_3D);let Q=Vt(C,S),Y=S.source;n.bindTexture(Z,C.__webglTexture,e.TEXTURE0+O);let Tt=i.get(Y);if(Y.version!==Tt.__version||Q===!0){n.activeTexture(e.TEXTURE0+O);let lt=se.getPrimaries(se.workingColorSpace),Rt=S.colorSpace===Ja?null:se.getPrimaries(S.colorSpace),Ut=S.colorSpace===Ja||lt===Rt?e.NONE:e.BROWSER_DEFAULT_WEBGL;e.pixelStorei(e.UNPACK_FLIP_Y_WEBGL,S.flipY),e.pixelStorei(e.UNPACK_PREMULTIPLY_ALPHA_WEBGL,S.premultiplyAlpha),e.pixelStorei(e.UNPACK_ALIGNMENT,S.unpackAlignment),e.pixelStorei(e.UNPACK_COLORSPACE_CONVERSION_WEBGL,Ut);let tt=b(S.image,!1,a.maxTextureSize);tt=ue(S,tt);let rt=s.convert(S.format,S.colorSpace),Mt=s.convert(S.type),bt=g(S.internalFormat,rt,Mt,S.colorSpace,S.isVideoTexture);ht(Z,S);let gt,Wt=S.mipmaps,U=S.isVideoTexture!==!0,ut=Tt.__version===void 0||Q===!0,ot=Y.dataReady,xt=E(S,tt);if(S.isDepthTexture)bt=x(S.format===qs,S.type),ut&&(U?n.texStorage2D(e.TEXTURE_2D,1,bt,tt.width,tt.height):n.texImage2D(e.TEXTURE_2D,0,bt,tt.width,tt.height,0,rt,Mt,null));else if(S.isDataTexture)if(Wt.length>0){U&&ut&&n.texStorage2D(e.TEXTURE_2D,xt,bt,Wt[0].width,Wt[0].height);for(let nt=0,q=Wt.length;nt<q;nt++)gt=Wt[nt],U?ot&&n.texSubImage2D(e.TEXTURE_2D,nt,0,0,gt.width,gt.height,rt,Mt,gt.data):n.texImage2D(e.TEXTURE_2D,nt,bt,gt.width,gt.height,0,rt,Mt,gt.data);S.generateMipmaps=!1}else U?(ut&&n.texStorage2D(e.TEXTURE_2D,xt,bt,tt.width,tt.height),ot&&xe(S,tt,rt,Mt)):n.texImage2D(e.TEXTURE_2D,0,bt,tt.width,tt.height,0,rt,Mt,tt.data);else if(S.isCompressedTexture)if(S.isCompressedArrayTexture){U&&ut&&n.texStorage3D(e.TEXTURE_2D_ARRAY,xt,bt,Wt[0].width,Wt[0].height,tt.depth);for(let nt=0,q=Wt.length;nt<q;nt++)if(gt=Wt[nt],S.format!==Ui)if(rt!==null)if(U){if(ot)if(S.layerUpdates.size>0){let Et=Ev(gt.width,gt.height,S.format,S.type);for(let Ht of S.layerUpdates){let ge=gt.data.subarray(Ht*Et/gt.data.BYTES_PER_ELEMENT,(Ht+1)*Et/gt.data.BYTES_PER_ELEMENT);n.compressedTexSubImage3D(e.TEXTURE_2D_ARRAY,nt,0,0,Ht,gt.width,gt.height,1,rt,ge)}S.clearLayerUpdates()}else n.compressedTexSubImage3D(e.TEXTURE_2D_ARRAY,nt,0,0,0,gt.width,gt.height,tt.depth,rt,gt.data)}else n.compressedTexImage3D(e.TEXTURE_2D_ARRAY,nt,bt,gt.width,gt.height,tt.depth,0,gt.data,0,0);else Ot("WebGLRenderer: Attempt to load unsupported compressed texture format in .uploadTexture()");else U?ot&&n.texSubImage3D(e.TEXTURE_2D_ARRAY,nt,0,0,0,gt.width,gt.height,tt.depth,rt,Mt,gt.data):n.texImage3D(e.TEXTURE_2D_ARRAY,nt,bt,gt.width,gt.height,tt.depth,0,rt,Mt,gt.data)}else{U&&ut&&n.texStorage2D(e.TEXTURE_2D,xt,bt,Wt[0].width,Wt[0].height);for(let nt=0,q=Wt.length;nt<q;nt++)gt=Wt[nt],S.format!==Ui?rt!==null?U?ot&&n.compressedTexSubImage2D(e.TEXTURE_2D,nt,0,0,gt.width,gt.height,rt,gt.data):n.compressedTexImage2D(e.TEXTURE_2D,nt,bt,gt.width,gt.height,0,gt.data):Ot("WebGLRenderer: Attempt to load unsupported compressed texture format in .uploadTexture()"):U?ot&&n.texSubImage2D(e.TEXTURE_2D,nt,0,0,gt.width,gt.height,rt,Mt,gt.data):n.texImage2D(e.TEXTURE_2D,nt,bt,gt.width,gt.height,0,rt,Mt,gt.data)}else if(S.isDataArrayTexture)if(U){if(ut&&n.texStorage3D(e.TEXTURE_2D_ARRAY,xt,bt,tt.width,tt.height,tt.depth),ot)if(S.layerUpdates.size>0){let nt=Ev(tt.width,tt.height,S.format,S.type);for(let q of S.layerUpdates){let Et=tt.data.subarray(q*nt/tt.data.BYTES_PER_ELEMENT,(q+1)*nt/tt.data.BYTES_PER_ELEMENT);n.texSubImage3D(e.TEXTURE_2D_ARRAY,0,0,0,q,tt.width,tt.height,1,rt,Mt,Et)}S.clearLayerUpdates()}else n.texSubImage3D(e.TEXTURE_2D_ARRAY,0,0,0,0,tt.width,tt.height,tt.depth,rt,Mt,tt.data)}else n.texImage3D(e.TEXTURE_2D_ARRAY,0,bt,tt.width,tt.height,tt.depth,0,rt,Mt,tt.data);else if(S.isData3DTexture)U?(ut&&n.texStorage3D(e.TEXTURE_3D,xt,bt,tt.width,tt.height,tt.depth),ot&&n.texSubImage3D(e.TEXTURE_3D,0,0,0,0,tt.width,tt.height,tt.depth,rt,Mt,tt.data)):n.texImage3D(e.TEXTURE_3D,0,bt,tt.width,tt.height,tt.depth,0,rt,Mt,tt.data);else if(S.isFramebufferTexture){if(ut)if(U)n.texStorage2D(e.TEXTURE_2D,xt,bt,tt.width,tt.height);else{let nt=tt.width,q=tt.height;for(let Et=0;Et<xt;Et++)n.texImage2D(e.TEXTURE_2D,Et,bt,nt,q,0,rt,Mt,null),nt>>=1,q>>=1}}else if(Wt.length>0){if(U&&ut){let nt=At(Wt[0]);n.texStorage2D(e.TEXTURE_2D,xt,bt,nt.width,nt.height)}for(let nt=0,q=Wt.length;nt<q;nt++)gt=Wt[nt],U?ot&&n.texSubImage2D(e.TEXTURE_2D,nt,0,0,rt,Mt,gt):n.texImage2D(e.TEXTURE_2D,nt,bt,rt,Mt,gt);S.generateMipmaps=!1}else if(U){if(ut){let nt=At(tt);n.texStorage2D(e.TEXTURE_2D,xt,bt,nt.width,nt.height)}ot&&n.texSubImage2D(e.TEXTURE_2D,0,0,0,rt,Mt,tt)}else n.texImage2D(e.TEXTURE_2D,0,bt,rt,Mt,tt);v(S)&&h(Z),Tt.__version=Y.version,S.onUpdate&&S.onUpdate(S)}C.__version=S.version}function j(C,S,O){if(S.image.length!==6)return;let Z=Vt(C,S),Q=S.source;n.bindTexture(e.TEXTURE_CUBE_MAP,C.__webglTexture,e.TEXTURE0+O);let Y=i.get(Q);if(Q.version!==Y.__version||Z===!0){n.activeTexture(e.TEXTURE0+O);let Tt=se.getPrimaries(se.workingColorSpace),lt=S.colorSpace===Ja?null:se.getPrimaries(S.colorSpace),Rt=S.colorSpace===Ja||Tt===lt?e.NONE:e.BROWSER_DEFAULT_WEBGL;e.pixelStorei(e.UNPACK_FLIP_Y_WEBGL,S.flipY),e.pixelStorei(e.UNPACK_PREMULTIPLY_ALPHA_WEBGL,S.premultiplyAlpha),e.pixelStorei(e.UNPACK_ALIGNMENT,S.unpackAlignment),e.pixelStorei(e.UNPACK_COLORSPACE_CONVERSION_WEBGL,Rt);let Ut=S.isCompressedTexture||S.image[0].isCompressedTexture,tt=S.image[0]&&S.image[0].isDataTexture,rt=[];for(let q=0;q<6;q++)!Ut&&!tt?rt[q]=b(S.image[q],!0,a.maxCubemapSize):rt[q]=tt?S.image[q].image:S.image[q],rt[q]=ue(S,rt[q]);let Mt=rt[0],bt=s.convert(S.format,S.colorSpace),gt=s.convert(S.type),Wt=g(S.internalFormat,bt,gt,S.colorSpace),U=S.isVideoTexture!==!0,ut=Y.__version===void 0||Z===!0,ot=Q.dataReady,xt=E(S,Mt);ht(e.TEXTURE_CUBE_MAP,S);let nt;if(Ut){U&&ut&&n.texStorage2D(e.TEXTURE_CUBE_MAP,xt,Wt,Mt.width,Mt.height);for(let q=0;q<6;q++){nt=rt[q].mipmaps;for(let Et=0;Et<nt.length;Et++){let Ht=nt[Et];S.format!==Ui?bt!==null?U?ot&&n.compressedTexSubImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,Et,0,0,Ht.width,Ht.height,bt,Ht.data):n.compressedTexImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,Et,Wt,Ht.width,Ht.height,0,Ht.data):Ot("WebGLRenderer: Attempt to load unsupported compressed texture format in .setTextureCube()"):U?ot&&n.texSubImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,Et,0,0,Ht.width,Ht.height,bt,gt,Ht.data):n.texImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,Et,Wt,Ht.width,Ht.height,0,bt,gt,Ht.data)}}}else{if(nt=S.mipmaps,U&&ut){nt.length>0&&xt++;let q=At(rt[0]);n.texStorage2D(e.TEXTURE_CUBE_MAP,xt,Wt,q.width,q.height)}for(let q=0;q<6;q++)if(tt){U?ot&&n.texSubImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,0,0,0,rt[q].width,rt[q].height,bt,gt,rt[q].data):n.texImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,0,Wt,rt[q].width,rt[q].height,0,bt,gt,rt[q].data);for(let Et=0;Et<nt.length;Et++){let ge=nt[Et].image[q].image;U?ot&&n.texSubImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,Et+1,0,0,ge.width,ge.height,bt,gt,ge.data):n.texImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,Et+1,Wt,ge.width,ge.height,0,bt,gt,ge.data)}}else{U?ot&&n.texSubImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,0,0,0,bt,gt,rt[q]):n.texImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,0,Wt,bt,gt,rt[q]);for(let Et=0;Et<nt.length;Et++){let Ht=nt[Et];U?ot&&n.texSubImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,Et+1,0,0,bt,gt,Ht.image[q]):n.texImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+q,Et+1,Wt,bt,gt,Ht.image[q])}}}v(S)&&h(e.TEXTURE_CUBE_MAP),Y.__version=Q.version,S.onUpdate&&S.onUpdate(S)}C.__version=S.version}function ct(C,S,O,Z,Q,Y){let Tt=s.convert(O.format,O.colorSpace),lt=s.convert(O.type),Rt=g(O.internalFormat,Tt,lt,O.colorSpace),Ut=i.get(S),tt=i.get(O);if(tt.__renderTarget=S,!Ut.__hasExternalTextures){let rt=Math.max(1,S.width>>Y),Mt=Math.max(1,S.height>>Y);Q===e.TEXTURE_3D||Q===e.TEXTURE_2D_ARRAY?n.texImage3D(Q,Y,Rt,rt,Mt,S.depth,0,Tt,lt,null):n.texImage2D(Q,Y,Rt,rt,Mt,0,Tt,lt,null)}n.bindFramebuffer(e.FRAMEBUFFER,C),ze(S)?o.framebufferTexture2DMultisampleEXT(e.FRAMEBUFFER,Z,Q,tt.__webglTexture,0,D(S)):(Q===e.TEXTURE_2D||Q>=e.TEXTURE_CUBE_MAP_POSITIVE_X&&Q<=e.TEXTURE_CUBE_MAP_NEGATIVE_Z)&&e.framebufferTexture2D(e.FRAMEBUFFER,Z,Q,tt.__webglTexture,Y),n.bindFramebuffer(e.FRAMEBUFFER,null)}function It(C,S,O){if(e.bindRenderbuffer(e.RENDERBUFFER,C),S.depthBuffer){let Z=S.depthTexture,Q=Z&&Z.isDepthTexture?Z.type:null,Y=x(S.stencilBuffer,Q),Tt=S.stencilBuffer?e.DEPTH_STENCIL_ATTACHMENT:e.DEPTH_ATTACHMENT;ze(S)?o.renderbufferStorageMultisampleEXT(e.RENDERBUFFER,D(S),Y,S.width,S.height):O?e.renderbufferStorageMultisample(e.RENDERBUFFER,D(S),Y,S.width,S.height):e.renderbufferStorage(e.RENDERBUFFER,Y,S.width,S.height),e.framebufferRenderbuffer(e.FRAMEBUFFER,Tt,e.RENDERBUFFER,C)}else{let Z=S.textures;for(let Q=0;Q<Z.length;Q++){let Y=Z[Q],Tt=s.convert(Y.format,Y.colorSpace),lt=s.convert(Y.type),Rt=g(Y.internalFormat,Tt,lt,Y.colorSpace);ze(S)?o.renderbufferStorageMultisampleEXT(e.RENDERBUFFER,D(S),Rt,S.width,S.height):O?e.renderbufferStorageMultisample(e.RENDERBUFFER,D(S),Rt,S.width,S.height):e.renderbufferStorage(e.RENDERBUFFER,Rt,S.width,S.height)}}e.bindRenderbuffer(e.RENDERBUFFER,null)}function Lt(C,S,O){let Z=S.isWebGLCubeRenderTarget===!0;if(n.bindFramebuffer(e.FRAMEBUFFER,C),!(S.depthTexture&&S.depthTexture.isDepthTexture))throw new Error("renderTarget.depthTexture must be an instance of THREE.DepthTexture");let Q=i.get(S.depthTexture);if(Q.__renderTarget=S,(!Q.__webglTexture||S.depthTexture.image.width!==S.width||S.depthTexture.image.height!==S.height)&&(S.depthTexture.image.width=S.width,S.depthTexture.image.height=S.height,S.depthTexture.needsUpdate=!0),Z){if(Q.__webglInit===void 0&&(Q.__webglInit=!0,S.depthTexture.addEventListener("dispose",A)),Q.__webglTexture===void 0){Q.__webglTexture=e.createTexture(),n.bindTexture(e.TEXTURE_CUBE_MAP,Q.__webglTexture),ht(e.TEXTURE_CUBE_MAP,S.depthTexture);let Ut=s.convert(S.depthTexture.format),tt=s.convert(S.depthTexture.type),rt;S.depthTexture.format===da?rt=e.DEPTH_COMPONENT24:S.depthTexture.format===qs&&(rt=e.DEPTH24_STENCIL8);for(let Mt=0;Mt<6;Mt++)e.texImage2D(e.TEXTURE_CUBE_MAP_POSITIVE_X+Mt,0,rt,S.width,S.height,0,Ut,tt,null)}}else F(S.depthTexture,0);let Y=Q.__webglTexture,Tt=D(S),lt=Z?e.TEXTURE_CUBE_MAP_POSITIVE_X+O:e.TEXTURE_2D,Rt=S.depthTexture.format===qs?e.DEPTH_STENCIL_ATTACHMENT:e.DEPTH_ATTACHMENT;if(S.depthTexture.format===da)ze(S)?o.framebufferTexture2DMultisampleEXT(e.FRAMEBUFFER,Rt,lt,Y,0,Tt):e.framebufferTexture2D(e.FRAMEBUFFER,Rt,lt,Y,0);else if(S.depthTexture.format===qs)ze(S)?o.framebufferTexture2DMultisampleEXT(e.FRAMEBUFFER,Rt,lt,Y,0,Tt):e.framebufferTexture2D(e.FRAMEBUFFER,Rt,lt,Y,0);else throw new Error("Unknown depthTexture format")}function Pt(C){let S=i.get(C),O=C.isWebGLCubeRenderTarget===!0;if(S.__boundDepthTexture!==C.depthTexture){let Z=C.depthTexture;if(S.__depthDisposeCallback&&S.__depthDisposeCallback(),Z){let Q=()=>{delete S.__boundDepthTexture,delete S.__depthDisposeCallback,Z.removeEventListener("dispose",Q)};Z.addEventListener("dispose",Q),S.__depthDisposeCallback=Q}S.__boundDepthTexture=Z}if(C.depthTexture&&!S.__autoAllocateDepthBuffer)if(O)for(let Z=0;Z<6;Z++)Lt(S.__webglFramebuffer[Z],C,Z);else{let Z=C.texture.mipmaps;Z&&Z.length>0?Lt(S.__webglFramebuffer[0],C,0):Lt(S.__webglFramebuffer,C,0)}else if(O){S.__webglDepthbuffer=[];for(let Z=0;Z<6;Z++)if(n.bindFramebuffer(e.FRAMEBUFFER,S.__webglFramebuffer[Z]),S.__webglDepthbuffer[Z]===void 0)S.__webglDepthbuffer[Z]=e.createRenderbuffer(),It(S.__webglDepthbuffer[Z],C,!1);else{let Q=C.stencilBuffer?e.DEPTH_STENCIL_ATTACHMENT:e.DEPTH_ATTACHMENT,Y=S.__webglDepthbuffer[Z];e.bindRenderbuffer(e.RENDERBUFFER,Y),e.framebufferRenderbuffer(e.FRAMEBUFFER,Q,e.RENDERBUFFER,Y)}}else{let Z=C.texture.mipmaps;if(Z&&Z.length>0?n.bindFramebuffer(e.FRAMEBUFFER,S.__webglFramebuffer[0]):n.bindFramebuffer(e.FRAMEBUFFER,S.__webglFramebuffer),S.__webglDepthbuffer===void 0)S.__webglDepthbuffer=e.createRenderbuffer(),It(S.__webglDepthbuffer,C,!1);else{let Q=C.stencilBuffer?e.DEPTH_STENCIL_ATTACHMENT:e.DEPTH_ATTACHMENT,Y=S.__webglDepthbuffer;e.bindRenderbuffer(e.RENDERBUFFER,Y),e.framebufferRenderbuffer(e.FRAMEBUFFER,Q,e.RENDERBUFFER,Y)}}n.bindFramebuffer(e.FRAMEBUFFER,null)}function Ye(C,S,O){let Z=i.get(C);S!==void 0&&ct(Z.__webglFramebuffer,C,C.texture,e.COLOR_ATTACHMENT0,e.TEXTURE_2D,0),O!==void 0&&Pt(C)}function ie(C){let S=C.texture,O=i.get(C),Z=i.get(S);C.addEventListener("dispose",w);let Q=C.textures,Y=C.isWebGLCubeRenderTarget===!0,Tt=Q.length>1;if(Tt||(Z.__webglTexture===void 0&&(Z.__webglTexture=e.createTexture()),Z.__version=S.version,r.memory.textures++),Y){O.__webglFramebuffer=[];for(let lt=0;lt<6;lt++)if(S.mipmaps&&S.mipmaps.length>0){O.__webglFramebuffer[lt]=[];for(let Rt=0;Rt<S.mipmaps.length;Rt++)O.__webglFramebuffer[lt][Rt]=e.createFramebuffer()}else O.__webglFramebuffer[lt]=e.createFramebuffer()}else{if(S.mipmaps&&S.mipmaps.length>0){O.__webglFramebuffer=[];for(let lt=0;lt<S.mipmaps.length;lt++)O.__webglFramebuffer[lt]=e.createFramebuffer()}else O.__webglFramebuffer=e.createFramebuffer();if(Tt)for(let lt=0,Rt=Q.length;lt<Rt;lt++){let Ut=i.get(Q[lt]);Ut.__webglTexture===void 0&&(Ut.__webglTexture=e.createTexture(),r.memory.textures++)}if(C.samples>0&&ze(C)===!1){O.__webglMultisampledFramebuffer=e.createFramebuffer(),O.__webglColorRenderbuffer=[],n.bindFramebuffer(e.FRAMEBUFFER,O.__webglMultisampledFramebuffer);for(let lt=0;lt<Q.length;lt++){let Rt=Q[lt];O.__webglColorRenderbuffer[lt]=e.createRenderbuffer(),e.bindRenderbuffer(e.RENDERBUFFER,O.__webglColorRenderbuffer[lt]);let Ut=s.convert(Rt.format,Rt.colorSpace),tt=s.convert(Rt.type),rt=g(Rt.internalFormat,Ut,tt,Rt.colorSpace,C.isXRRenderTarget===!0),Mt=D(C);e.renderbufferStorageMultisample(e.RENDERBUFFER,Mt,rt,C.width,C.height),e.framebufferRenderbuffer(e.FRAMEBUFFER,e.COLOR_ATTACHMENT0+lt,e.RENDERBUFFER,O.__webglColorRenderbuffer[lt])}e.bindRenderbuffer(e.RENDERBUFFER,null),C.depthBuffer&&(O.__webglDepthRenderbuffer=e.createRenderbuffer(),It(O.__webglDepthRenderbuffer,C,!0)),n.bindFramebuffer(e.FRAMEBUFFER,null)}}if(Y){n.bindTexture(e.TEXTURE_CUBE_MAP,Z.__webglTexture),ht(e.TEXTURE_CUBE_MAP,S);for(let lt=0;lt<6;lt++)if(S.mipmaps&&S.mipmaps.length>0)for(let Rt=0;Rt<S.mipmaps.length;Rt++)ct(O.__webglFramebuffer[lt][Rt],C,S,e.COLOR_ATTACHMENT0,e.TEXTURE_CUBE_MAP_POSITIVE_X+lt,Rt);else ct(O.__webglFramebuffer[lt],C,S,e.COLOR_ATTACHMENT0,e.TEXTURE_CUBE_MAP_POSITIVE_X+lt,0);v(S)&&h(e.TEXTURE_CUBE_MAP),n.unbindTexture()}else if(Tt){for(let lt=0,Rt=Q.length;lt<Rt;lt++){let Ut=Q[lt],tt=i.get(Ut),rt=e.TEXTURE_2D;(C.isWebGL3DRenderTarget||C.isWebGLArrayRenderTarget)&&(rt=C.isWebGL3DRenderTarget?e.TEXTURE_3D:e.TEXTURE_2D_ARRAY),n.bindTexture(rt,tt.__webglTexture),ht(rt,Ut),ct(O.__webglFramebuffer,C,Ut,e.COLOR_ATTACHMENT0+lt,rt,0),v(Ut)&&h(rt)}n.unbindTexture()}else{let lt=e.TEXTURE_2D;if((C.isWebGL3DRenderTarget||C.isWebGLArrayRenderTarget)&&(lt=C.isWebGL3DRenderTarget?e.TEXTURE_3D:e.TEXTURE_2D_ARRAY),n.bindTexture(lt,Z.__webglTexture),ht(lt,S),S.mipmaps&&S.mipmaps.length>0)for(let Rt=0;Rt<S.mipmaps.length;Rt++)ct(O.__webglFramebuffer[Rt],C,S,e.COLOR_ATTACHMENT0,lt,Rt);else ct(O.__webglFramebuffer,C,S,e.COLOR_ATTACHMENT0,lt,0);v(S)&&h(lt),n.unbindTexture()}C.depthBuffer&&Pt(C)}function Jt(C){let S=C.textures;for(let O=0,Z=S.length;O<Z;O++){let Q=S[O];if(v(Q)){let Y=m(C),Tt=i.get(Q).__webglTexture;n.bindTexture(Y,Tt),h(Y),n.unbindTexture()}}}let de=[],Xt=[];function Te(C){if(C.samples>0){if(ze(C)===!1){let S=C.textures,O=C.width,Z=C.height,Q=e.COLOR_BUFFER_BIT,Y=C.stencilBuffer?e.DEPTH_STENCIL_ATTACHMENT:e.DEPTH_ATTACHMENT,Tt=i.get(C),lt=S.length>1;if(lt)for(let Ut=0;Ut<S.length;Ut++)n.bindFramebuffer(e.FRAMEBUFFER,Tt.__webglMultisampledFramebuffer),e.framebufferRenderbuffer(e.FRAMEBUFFER,e.COLOR_ATTACHMENT0+Ut,e.RENDERBUFFER,null),n.bindFramebuffer(e.FRAMEBUFFER,Tt.__webglFramebuffer),e.framebufferTexture2D(e.DRAW_FRAMEBUFFER,e.COLOR_ATTACHMENT0+Ut,e.TEXTURE_2D,null,0);n.bindFramebuffer(e.READ_FRAMEBUFFER,Tt.__webglMultisampledFramebuffer);let Rt=C.texture.mipmaps;Rt&&Rt.length>0?n.bindFramebuffer(e.DRAW_FRAMEBUFFER,Tt.__webglFramebuffer[0]):n.bindFramebuffer(e.DRAW_FRAMEBUFFER,Tt.__webglFramebuffer);for(let Ut=0;Ut<S.length;Ut++){if(C.resolveDepthBuffer&&(C.depthBuffer&&(Q|=e.DEPTH_BUFFER_BIT),C.stencilBuffer&&C.resolveStencilBuffer&&(Q|=e.STENCIL_BUFFER_BIT)),lt){e.framebufferRenderbuffer(e.READ_FRAMEBUFFER,e.COLOR_ATTACHMENT0,e.RENDERBUFFER,Tt.__webglColorRenderbuffer[Ut]);let tt=i.get(S[Ut]).__webglTexture;e.framebufferTexture2D(e.DRAW_FRAMEBUFFER,e.COLOR_ATTACHMENT0,e.TEXTURE_2D,tt,0)}e.blitFramebuffer(0,0,O,Z,0,0,O,Z,Q,e.NEAREST),l===!0&&(de.length=0,Xt.length=0,de.push(e.COLOR_ATTACHMENT0+Ut),C.depthBuffer&&C.resolveDepthBuffer===!1&&(de.push(Y),Xt.push(Y),e.invalidateFramebuffer(e.DRAW_FRAMEBUFFER,Xt)),e.invalidateFramebuffer(e.READ_FRAMEBUFFER,de))}if(n.bindFramebuffer(e.READ_FRAMEBUFFER,null),n.bindFramebuffer(e.DRAW_FRAMEBUFFER,null),lt)for(let Ut=0;Ut<S.length;Ut++){n.bindFramebuffer(e.FRAMEBUFFER,Tt.__webglMultisampledFramebuffer),e.framebufferRenderbuffer(e.FRAMEBUFFER,e.COLOR_ATTACHMENT0+Ut,e.RENDERBUFFER,Tt.__webglColorRenderbuffer[Ut]);let tt=i.get(S[Ut]).__webglTexture;n.bindFramebuffer(e.FRAMEBUFFER,Tt.__webglFramebuffer),e.framebufferTexture2D(e.DRAW_FRAMEBUFFER,e.COLOR_ATTACHMENT0+Ut,e.TEXTURE_2D,tt,0)}n.bindFramebuffer(e.DRAW_FRAMEBUFFER,Tt.__webglMultisampledFramebuffer)}else if(C.depthBuffer&&C.resolveDepthBuffer===!1&&l){let S=C.stencilBuffer?e.DEPTH_STENCIL_ATTACHMENT:e.DEPTH_ATTACHMENT;e.invalidateFramebuffer(e.DRAW_FRAMEBUFFER,[S])}}}function D(C){return Math.min(a.maxSamples,C.samples)}function ze(C){let S=i.get(C);return C.samples>0&&t.has("WEBGL_multisampled_render_to_texture")===!0&&S.__useRenderToTexture!==!1}function re(C){let S=r.render.frame;f.get(C)!==S&&(f.set(C,S),C.update())}function ue(C,S){let O=C.colorSpace,Z=C.format,Q=C.type;return C.isCompressedTexture===!0||C.isVideoTexture===!0||O!==Rr&&O!==Ja&&(se.getTransfer(O)===me?(Z!==Ui||Q!==ei)&&Ot("WebGLTextures: sRGB encoded textures have to use RGBAFormat and UnsignedByteType."):zt("WebGLTextures: Unsupported texture color space:",O)),S}function At(C){return typeof HTMLImageElement<"u"&&C instanceof HTMLImageElement?(c.width=C.naturalWidth||C.width,c.height=C.naturalHeight||C.height):typeof VideoFrame<"u"&&C instanceof VideoFrame?(c.width=C.displayWidth,c.height=C.displayHeight):(c.width=C.width,c.height=C.height),c}this.allocateTextureUnit=P,this.resetTextureUnits=N,this.setTexture2D=F,this.setTexture2DArray=V,this.setTexture3D=X,this.setTextureCube=at,this.rebindTextures=Ye,this.setupRenderTarget=ie,this.updateRenderTargetMipmap=Jt,this.updateMultisampleRenderTarget=Te,this.setupDepthRenderbuffer=Pt,this.setupFrameBufferTexture=ct,this.useMultisampledRTT=ze,this.isReversedDepthBuffer=function(){return n.buffers.depth.getReversed()}}function vN(e,t){function n(i,a=Ja){let s,r=se.getTransfer(a);if(i===ei)return e.UNSIGNED_BYTE;if(i===bd)return e.UNSIGNED_SHORT_4_4_4_4;if(i===Md)return e.UNSIGNED_SHORT_5_5_5_1;if(i===pv)return e.UNSIGNED_INT_5_9_9_9_REV;if(i===mv)return e.UNSIGNED_INT_10F_11F_11F_REV;if(i===hv)return e.BYTE;if(i===dv)return e.SHORT;if(i===pl)return e.UNSIGNED_SHORT;if(i===Sd)return e.INT;if(i===qi)return e.UNSIGNED_INT;if(i===Ni)return e.FLOAT;if(i===_a)return e.HALF_FLOAT;if(i===gv)return e.ALPHA;if(i===vv)return e.RGB;if(i===Ui)return e.RGBA;if(i===da)return e.DEPTH_COMPONENT;if(i===qs)return e.DEPTH_STENCIL;if(i===Td)return e.RED;if(i===Ed)return e.RED_INTEGER;if(i===Lr)return e.RG;if(i===Ad)return e.RG_INTEGER;if(i===wd)return e.RGBA_INTEGER;if(i===mu||i===gu||i===vu||i===_u)if(r===me)if(s=t.get("WEBGL_compressed_texture_s3tc_srgb"),s!==null){if(i===mu)return s.COMPRESSED_SRGB_S3TC_DXT1_EXT;if(i===gu)return s.COMPRESSED_SRGB_ALPHA_S3TC_DXT1_EXT;if(i===vu)return s.COMPRESSED_SRGB_ALPHA_S3TC_DXT3_EXT;if(i===_u)return s.COMPRESSED_SRGB_ALPHA_S3TC_DXT5_EXT}else return null;else if(s=t.get("WEBGL_compressed_texture_s3tc"),s!==null){if(i===mu)return s.COMPRESSED_RGB_S3TC_DXT1_EXT;if(i===gu)return s.COMPRESSED_RGBA_S3TC_DXT1_EXT;if(i===vu)return s.COMPRESSED_RGBA_S3TC_DXT3_EXT;if(i===_u)return s.COMPRESSED_RGBA_S3TC_DXT5_EXT}else return null;if(i===Cd||i===Rd||i===Dd||i===Nd)if(s=t.get("WEBGL_compressed_texture_pvrtc"),s!==null){if(i===Cd)return s.COMPRESSED_RGB_PVRTC_4BPPV1_IMG;if(i===Rd)return s.COMPRESSED_RGB_PVRTC_2BPPV1_IMG;if(i===Dd)return s.COMPRESSED_RGBA_PVRTC_4BPPV1_IMG;if(i===Nd)return s.COMPRESSED_RGBA_PVRTC_2BPPV1_IMG}else return null;if(i===Ud||i===Ld||i===Id||i===Od||i===Pd||i===Bd||i===zd)if(s=t.get("WEBGL_compressed_texture_etc"),s!==null){if(i===Ud||i===Ld)return r===me?s.COMPRESSED_SRGB8_ETC2:s.COMPRESSED_RGB8_ETC2;if(i===Id)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ETC2_EAC:s.COMPRESSED_RGBA8_ETC2_EAC;if(i===Od)return s.COMPRESSED_R11_EAC;if(i===Pd)return s.COMPRESSED_SIGNED_R11_EAC;if(i===Bd)return s.COMPRESSED_RG11_EAC;if(i===zd)return s.COMPRESSED_SIGNED_RG11_EAC}else return null;if(i===Fd||i===Vd||i===Hd||i===Gd||i===kd||i===Xd||i===Wd||i===qd||i===Yd||i===Zd||i===Jd||i===Kd||i===Qd||i===jd)if(s=t.get("WEBGL_compressed_texture_astc"),s!==null){if(i===Fd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_4x4_KHR:s.COMPRESSED_RGBA_ASTC_4x4_KHR;if(i===Vd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_5x4_KHR:s.COMPRESSED_RGBA_ASTC_5x4_KHR;if(i===Hd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_5x5_KHR:s.COMPRESSED_RGBA_ASTC_5x5_KHR;if(i===Gd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_6x5_KHR:s.COMPRESSED_RGBA_ASTC_6x5_KHR;if(i===kd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_6x6_KHR:s.COMPRESSED_RGBA_ASTC_6x6_KHR;if(i===Xd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_8x5_KHR:s.COMPRESSED_RGBA_ASTC_8x5_KHR;if(i===Wd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_8x6_KHR:s.COMPRESSED_RGBA_ASTC_8x6_KHR;if(i===qd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_8x8_KHR:s.COMPRESSED_RGBA_ASTC_8x8_KHR;if(i===Yd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_10x5_KHR:s.COMPRESSED_RGBA_ASTC_10x5_KHR;if(i===Zd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_10x6_KHR:s.COMPRESSED_RGBA_ASTC_10x6_KHR;if(i===Jd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_10x8_KHR:s.COMPRESSED_RGBA_ASTC_10x8_KHR;if(i===Kd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_10x10_KHR:s.COMPRESSED_RGBA_ASTC_10x10_KHR;if(i===Qd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_12x10_KHR:s.COMPRESSED_RGBA_ASTC_12x10_KHR;if(i===jd)return r===me?s.COMPRESSED_SRGB8_ALPHA8_ASTC_12x12_KHR:s.COMPRESSED_RGBA_ASTC_12x12_KHR}else return null;if(i===$d||i===tp||i===ep)if(s=t.get("EXT_texture_compression_bptc"),s!==null){if(i===$d)return r===me?s.COMPRESSED_SRGB_ALPHA_BPTC_UNORM_EXT:s.COMPRESSED_RGBA_BPTC_UNORM_EXT;if(i===tp)return s.COMPRESSED_RGB_BPTC_SIGNED_FLOAT_EXT;if(i===ep)return s.COMPRESSED_RGB_BPTC_UNSIGNED_FLOAT_EXT}else return null;if(i===np||i===ip||i===ap||i===sp)if(s=t.get("EXT_texture_compression_rgtc"),s!==null){if(i===np)return s.COMPRESSED_RED_RGTC1_EXT;if(i===ip)return s.COMPRESSED_SIGNED_RED_RGTC1_EXT;if(i===ap)return s.COMPRESSED_RED_GREEN_RGTC2_EXT;if(i===sp)return s.COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT}else return null;return i===ml?e.UNSIGNED_INT_24_8:e[i]!==void 0?e[i]:null}return{convert:n}}var _N=`
void main() {

	gl_Position = vec4( position, 1.0 );

}`,yN=`
uniform sampler2DArray depthColor;
uniform float depthWidth;
uniform float depthHeight;

void main() {

	vec2 coord = vec2( gl_FragCoord.x / depthWidth, gl_FragCoord.y / depthHeight );

	if ( coord.x >= 1.0 ) {

		gl_FragDepth = texture( depthColor, vec3( coord.x - 1.0, coord.y, 1 ) ).r;

	} else {

		gl_FragDepth = texture( depthColor, vec3( coord.x, coord.y, 0 ) ).r;

	}

}`,Fv=class{constructor(){this.texture=null,this.mesh=null,this.depthNear=0,this.depthFar=0}init(t,n){if(this.texture===null){let i=new nu(t.texture);(t.depthNear!==n.depthNear||t.depthFar!==n.depthFar)&&(this.depthNear=t.depthNear,this.depthFar=t.depthFar),this.texture=i}}getMesh(t){if(this.texture!==null&&this.mesh===null){let n=t.cameras[0].viewport,i=new Yn({vertexShader:_N,fragmentShader:yN,uniforms:{depthColor:{value:this.texture},depthWidth:{value:n.z},depthHeight:{value:n.w}}});this.mesh=new mn(new au(20,20),i)}return this.mesh}reset(){this.texture=null,this.mesh=null}getDepthTexture(){return this.texture}},Vv=class extends qa{constructor(t,n){super();let i=this,a=null,s=1,r=null,o="local-floor",l=1,c=null,f=null,p=null,u=null,d=null,_=null,b=typeof XRWebGLBinding<"u",v=new Fv,h={},m=n.getContextAttributes(),g=null,x=null,E=[],A=[],w=new te,y=null,T=new Pn;T.viewport=new Ge;let L=new Pn;L.viewport=new Ge;let R=[T,L],N=new vd,P=null,z=null;this.cameraAutoUpdate=!0,this.enabled=!1,this.isPresenting=!1,this.getController=function(K){let j=E[K];return j===void 0&&(j=new rl,E[K]=j),j.getTargetRaySpace()},this.getControllerGrip=function(K){let j=E[K];return j===void 0&&(j=new rl,E[K]=j),j.getGripSpace()},this.getHand=function(K){let j=E[K];return j===void 0&&(j=new rl,E[K]=j),j.getHandSpace()};function F(K){let j=A.indexOf(K.inputSource);if(j===-1)return;let ct=E[j];ct!==void 0&&(ct.update(K.inputSource,K.frame,c||r),ct.dispatchEvent({type:K.type,data:K.inputSource}))}function V(){a.removeEventListener("select",F),a.removeEventListener("selectstart",F),a.removeEventListener("selectend",F),a.removeEventListener("squeeze",F),a.removeEventListener("squeezestart",F),a.removeEventListener("squeezeend",F),a.removeEventListener("end",V),a.removeEventListener("inputsourceschange",X);for(let K=0;K<E.length;K++){let j=A[K];j!==null&&(A[K]=null,E[K].disconnect(j))}P=null,z=null,v.reset();for(let K in h)delete h[K];t.setRenderTarget(g),d=null,u=null,p=null,a=null,x=null,xe.stop(),i.isPresenting=!1,t.setPixelRatio(y),t.setSize(w.width,w.height,!1),i.dispatchEvent({type:"sessionend"})}this.setFramebufferScaleFactor=function(K){s=K,i.isPresenting===!0&&Ot("WebXRManager: Cannot change framebuffer scale while presenting.")},this.setReferenceSpaceType=function(K){o=K,i.isPresenting===!0&&Ot("WebXRManager: Cannot change reference space type while presenting.")},this.getReferenceSpace=function(){return c||r},this.setReferenceSpace=function(K){c=K},this.getBaseLayer=function(){return u!==null?u:d},this.getBinding=function(){return p===null&&b&&(p=new XRWebGLBinding(a,n)),p},this.getFrame=function(){return _},this.getSession=function(){return a},this.setSession=async function(K){if(a=K,a!==null){if(g=t.getRenderTarget(),a.addEventListener("select",F),a.addEventListener("selectstart",F),a.addEventListener("selectend",F),a.addEventListener("squeeze",F),a.addEventListener("squeezestart",F),a.addEventListener("squeezeend",F),a.addEventListener("end",V),a.addEventListener("inputsourceschange",X),m.xrCompatible!==!0&&await n.makeXRCompatible(),y=t.getPixelRatio(),t.getSize(w),b&&"createProjectionLayer"in XRWebGLBinding.prototype){let ct=null,It=null,Lt=null;m.depth&&(Lt=m.stencil?n.DEPTH24_STENCIL8:n.DEPTH_COMPONENT24,ct=m.stencil?qs:da,It=m.stencil?ml:qi);let Pt={colorFormat:n.RGBA8,depthFormat:Lt,scaleFactor:s};p=this.getBinding(),u=p.createProjectionLayer(Pt),a.updateRenderState({layers:[u]}),t.setPixelRatio(1),t.setSize(u.textureWidth,u.textureHeight,!1),x=new gi(u.textureWidth,u.textureHeight,{format:Ui,type:ei,depthTexture:new Fs(u.textureWidth,u.textureHeight,It,void 0,void 0,void 0,void 0,void 0,void 0,ct),stencilBuffer:m.stencil,colorSpace:t.outputColorSpace,samples:m.antialias?4:0,resolveDepthBuffer:u.ignoreDepthValues===!1,resolveStencilBuffer:u.ignoreDepthValues===!1})}else{let ct={antialias:m.antialias,alpha:!0,depth:m.depth,stencil:m.stencil,framebufferScaleFactor:s};d=new XRWebGLLayer(a,n,ct),a.updateRenderState({baseLayer:d}),t.setPixelRatio(1),t.setSize(d.framebufferWidth,d.framebufferHeight,!1),x=new gi(d.framebufferWidth,d.framebufferHeight,{format:Ui,type:ei,colorSpace:t.outputColorSpace,stencilBuffer:m.stencil,resolveDepthBuffer:d.ignoreDepthValues===!1,resolveStencilBuffer:d.ignoreDepthValues===!1})}x.isXRRenderTarget=!0,this.setFoveation(l),c=null,r=await a.requestReferenceSpace(o),xe.setContext(a),xe.start(),i.isPresenting=!0,i.dispatchEvent({type:"sessionstart"})}},this.getEnvironmentBlendMode=function(){if(a!==null)return a.environmentBlendMode},this.getDepthTexture=function(){return v.getDepthTexture()};function X(K){for(let j=0;j<K.removed.length;j++){let ct=K.removed[j],It=A.indexOf(ct);It>=0&&(A[It]=null,E[It].disconnect(ct))}for(let j=0;j<K.added.length;j++){let ct=K.added[j],It=A.indexOf(ct);if(It===-1){for(let Pt=0;Pt<E.length;Pt++)if(Pt>=A.length){A.push(ct),It=Pt;break}else if(A[Pt]===null){A[Pt]=ct,It=Pt;break}if(It===-1)break}let Lt=E[It];Lt&&Lt.connect(ct)}}let at=new B,et=new B;function dt(K,j,ct){at.setFromMatrixPosition(j.matrixWorld),et.setFromMatrixPosition(ct.matrixWorld);let It=at.distanceTo(et),Lt=j.projectionMatrix.elements,Pt=ct.projectionMatrix.elements,Ye=Lt[14]/(Lt[10]-1),ie=Lt[14]/(Lt[10]+1),Jt=(Lt[9]+1)/Lt[5],de=(Lt[9]-1)/Lt[5],Xt=(Lt[8]-1)/Lt[0],Te=(Pt[8]+1)/Pt[0],D=Ye*Xt,ze=Ye*Te,re=It/(-Xt+Te),ue=re*-Xt;if(j.matrixWorld.decompose(K.position,K.quaternion,K.scale),K.translateX(ue),K.translateZ(re),K.matrixWorld.compose(K.position,K.quaternion,K.scale),K.matrixWorldInverse.copy(K.matrixWorld).invert(),Lt[10]===-1)K.projectionMatrix.copy(j.projectionMatrix),K.projectionMatrixInverse.copy(j.projectionMatrixInverse);else{let At=Ye+re,C=ie+re,S=D-ue,O=ze+(It-ue),Z=Jt*ie/C*At,Q=de*ie/C*At;K.projectionMatrix.makePerspective(S,O,Z,Q,At,C),K.projectionMatrixInverse.copy(K.projectionMatrix).invert()}}function St(K,j){j===null?K.matrixWorld.copy(K.matrix):K.matrixWorld.multiplyMatrices(j.matrixWorld,K.matrix),K.matrixWorldInverse.copy(K.matrixWorld).invert()}this.updateCamera=function(K){if(a===null)return;let j=K.near,ct=K.far;v.texture!==null&&(v.depthNear>0&&(j=v.depthNear),v.depthFar>0&&(ct=v.depthFar)),N.near=L.near=T.near=j,N.far=L.far=T.far=ct,(P!==N.near||z!==N.far)&&(a.updateRenderState({depthNear:N.near,depthFar:N.far}),P=N.near,z=N.far),N.layers.mask=K.layers.mask|6,T.layers.mask=N.layers.mask&-5,L.layers.mask=N.layers.mask&-3;let It=K.parent,Lt=N.cameras;St(N,It);for(let Pt=0;Pt<Lt.length;Pt++)St(Lt[Pt],It);Lt.length===2?dt(N,T,L):N.projectionMatrix.copy(T.projectionMatrix),ht(K,N,It)};function ht(K,j,ct){ct===null?K.matrix.copy(j.matrixWorld):(K.matrix.copy(ct.matrixWorld),K.matrix.invert(),K.matrix.multiply(j.matrixWorld)),K.matrix.decompose(K.position,K.quaternion,K.scale),K.updateMatrixWorld(!0),K.projectionMatrix.copy(j.projectionMatrix),K.projectionMatrixInverse.copy(j.projectionMatrixInverse),K.isPerspectiveCamera&&(K.fov=al*2*Math.atan(1/K.projectionMatrix.elements[5]),K.zoom=1)}this.getCamera=function(){return N},this.getFoveation=function(){if(!(u===null&&d===null))return l},this.setFoveation=function(K){l=K,u!==null&&(u.fixedFoveation=K),d!==null&&d.fixedFoveation!==void 0&&(d.fixedFoveation=K)},this.hasDepthSensing=function(){return v.texture!==null},this.getDepthSensingMesh=function(){return v.getMesh(N)},this.getCameraTexture=function(K){return h[K]};let Vt=null;function Ce(K,j){if(f=j.getViewerPose(c||r),_=j,f!==null){let ct=f.views;d!==null&&(t.setRenderTargetFramebuffer(x,d.framebuffer),t.setRenderTarget(x));let It=!1;ct.length!==N.cameras.length&&(N.cameras.length=0,It=!0);for(let ie=0;ie<ct.length;ie++){let Jt=ct[ie],de=null;if(d!==null)de=d.getViewport(Jt);else{let Te=p.getViewSubImage(u,Jt);de=Te.viewport,ie===0&&(t.setRenderTargetTextures(x,Te.colorTexture,Te.depthStencilTexture),t.setRenderTarget(x))}let Xt=R[ie];Xt===void 0&&(Xt=new Pn,Xt.layers.enable(ie),Xt.viewport=new Ge,R[ie]=Xt),Xt.matrix.fromArray(Jt.transform.matrix),Xt.matrix.decompose(Xt.position,Xt.quaternion,Xt.scale),Xt.projectionMatrix.fromArray(Jt.projectionMatrix),Xt.projectionMatrixInverse.copy(Xt.projectionMatrix).invert(),Xt.viewport.set(de.x,de.y,de.width,de.height),ie===0&&(N.matrix.copy(Xt.matrix),N.matrix.decompose(N.position,N.quaternion,N.scale)),It===!0&&N.cameras.push(Xt)}let Lt=a.enabledFeatures;if(Lt&&Lt.includes("depth-sensing")&&a.depthUsage=="gpu-optimized"&&b){p=i.getBinding();let ie=p.getDepthInformation(ct[0]);ie&&ie.isValid&&ie.texture&&v.init(ie,a.renderState)}if(Lt&&Lt.includes("camera-access")&&b){t.state.unbindTexture(),p=i.getBinding();for(let ie=0;ie<ct.length;ie++){let Jt=ct[ie].camera;if(Jt){let de=h[Jt];de||(de=new nu,h[Jt]=de);let Xt=p.getCameraImage(Jt);de.sourceTexture=Xt}}}}for(let ct=0;ct<E.length;ct++){let It=A[ct],Lt=E[ct];It!==null&&Lt!==void 0&&Lt.update(It,j,c||r)}Vt&&Vt(K,j),j.detectedPlanes&&i.dispatchEvent({type:"planesdetected",data:j}),_=null}let xe=new Y1;xe.setAnimationLoop(Ce),this.setAnimationLoop=function(K){Vt=K},this.dispose=function(){}}},Pr=new Xi,xN=new Ie;function SN(e,t){function n(v,h){v.matrixAutoUpdate===!0&&v.updateMatrix(),h.value.copy(v.matrix)}function i(v,h){h.color.getRGB(v.fogColor.value,bv(e)),h.isFog?(v.fogNear.value=h.near,v.fogFar.value=h.far):h.isFogExp2&&(v.fogDensity.value=h.density)}function a(v,h,m,g,x){h.isMeshBasicMaterial?s(v,h):h.isMeshLambertMaterial?(s(v,h),h.envMap&&(v.envMapIntensity.value=h.envMapIntensity)):h.isMeshToonMaterial?(s(v,h),p(v,h)):h.isMeshPhongMaterial?(s(v,h),f(v,h),h.envMap&&(v.envMapIntensity.value=h.envMapIntensity)):h.isMeshStandardMaterial?(s(v,h),u(v,h),h.isMeshPhysicalMaterial&&d(v,h,x)):h.isMeshMatcapMaterial?(s(v,h),_(v,h)):h.isMeshDepthMaterial?s(v,h):h.isMeshDistanceMaterial?(s(v,h),b(v,h)):h.isMeshNormalMaterial?s(v,h):h.isLineBasicMaterial?(r(v,h),h.isLineDashedMaterial&&o(v,h)):h.isPointsMaterial?l(v,h,m,g):h.isSpriteMaterial?c(v,h):h.isShadowMaterial?(v.color.value.copy(h.color),v.opacity.value=h.opacity):h.isShaderMaterial&&(h.uniformsNeedUpdate=!1)}function s(v,h){v.opacity.value=h.opacity,h.color&&v.diffuse.value.copy(h.color),h.emissive&&v.emissive.value.copy(h.emissive).multiplyScalar(h.emissiveIntensity),h.map&&(v.map.value=h.map,n(h.map,v.mapTransform)),h.alphaMap&&(v.alphaMap.value=h.alphaMap,n(h.alphaMap,v.alphaMapTransform)),h.bumpMap&&(v.bumpMap.value=h.bumpMap,n(h.bumpMap,v.bumpMapTransform),v.bumpScale.value=h.bumpScale,h.side===wn&&(v.bumpScale.value*=-1)),h.normalMap&&(v.normalMap.value=h.normalMap,n(h.normalMap,v.normalMapTransform),v.normalScale.value.copy(h.normalScale),h.side===wn&&v.normalScale.value.negate()),h.displacementMap&&(v.displacementMap.value=h.displacementMap,n(h.displacementMap,v.displacementMapTransform),v.displacementScale.value=h.displacementScale,v.displacementBias.value=h.displacementBias),h.emissiveMap&&(v.emissiveMap.value=h.emissiveMap,n(h.emissiveMap,v.emissiveMapTransform)),h.specularMap&&(v.specularMap.value=h.specularMap,n(h.specularMap,v.specularMapTransform)),h.alphaTest>0&&(v.alphaTest.value=h.alphaTest);let m=t.get(h),g=m.envMap,x=m.envMapRotation;g&&(v.envMap.value=g,Pr.copy(x),Pr.x*=-1,Pr.y*=-1,Pr.z*=-1,g.isCubeTexture&&g.isRenderTargetTexture===!1&&(Pr.y*=-1,Pr.z*=-1),v.envMapRotation.value.setFromMatrix4(xN.makeRotationFromEuler(Pr)),v.flipEnvMap.value=g.isCubeTexture&&g.isRenderTargetTexture===!1?-1:1,v.reflectivity.value=h.reflectivity,v.ior.value=h.ior,v.refractionRatio.value=h.refractionRatio),h.lightMap&&(v.lightMap.value=h.lightMap,v.lightMapIntensity.value=h.lightMapIntensity,n(h.lightMap,v.lightMapTransform)),h.aoMap&&(v.aoMap.value=h.aoMap,v.aoMapIntensity.value=h.aoMapIntensity,n(h.aoMap,v.aoMapTransform))}function r(v,h){v.diffuse.value.copy(h.color),v.opacity.value=h.opacity,h.map&&(v.map.value=h.map,n(h.map,v.mapTransform))}function o(v,h){v.dashSize.value=h.dashSize,v.totalSize.value=h.dashSize+h.gapSize,v.scale.value=h.scale}function l(v,h,m,g){v.diffuse.value.copy(h.color),v.opacity.value=h.opacity,v.size.value=h.size*m,v.scale.value=g*.5,h.map&&(v.map.value=h.map,n(h.map,v.uvTransform)),h.alphaMap&&(v.alphaMap.value=h.alphaMap,n(h.alphaMap,v.alphaMapTransform)),h.alphaTest>0&&(v.alphaTest.value=h.alphaTest)}function c(v,h){v.diffuse.value.copy(h.color),v.opacity.value=h.opacity,v.rotation.value=h.rotation,h.map&&(v.map.value=h.map,n(h.map,v.mapTransform)),h.alphaMap&&(v.alphaMap.value=h.alphaMap,n(h.alphaMap,v.alphaMapTransform)),h.alphaTest>0&&(v.alphaTest.value=h.alphaTest)}function f(v,h){v.specular.value.copy(h.specular),v.shininess.value=Math.max(h.shininess,1e-4)}function p(v,h){h.gradientMap&&(v.gradientMap.value=h.gradientMap)}function u(v,h){v.metalness.value=h.metalness,h.metalnessMap&&(v.metalnessMap.value=h.metalnessMap,n(h.metalnessMap,v.metalnessMapTransform)),v.roughness.value=h.roughness,h.roughnessMap&&(v.roughnessMap.value=h.roughnessMap,n(h.roughnessMap,v.roughnessMapTransform)),h.envMap&&(v.envMapIntensity.value=h.envMapIntensity)}function d(v,h,m){v.ior.value=h.ior,h.sheen>0&&(v.sheenColor.value.copy(h.sheenColor).multiplyScalar(h.sheen),v.sheenRoughness.value=h.sheenRoughness,h.sheenColorMap&&(v.sheenColorMap.value=h.sheenColorMap,n(h.sheenColorMap,v.sheenColorMapTransform)),h.sheenRoughnessMap&&(v.sheenRoughnessMap.value=h.sheenRoughnessMap,n(h.sheenRoughnessMap,v.sheenRoughnessMapTransform))),h.clearcoat>0&&(v.clearcoat.value=h.clearcoat,v.clearcoatRoughness.value=h.clearcoatRoughness,h.clearcoatMap&&(v.clearcoatMap.value=h.clearcoatMap,n(h.clearcoatMap,v.clearcoatMapTransform)),h.clearcoatRoughnessMap&&(v.clearcoatRoughnessMap.value=h.clearcoatRoughnessMap,n(h.clearcoatRoughnessMap,v.clearcoatRoughnessMapTransform)),h.clearcoatNormalMap&&(v.clearcoatNormalMap.value=h.clearcoatNormalMap,n(h.clearcoatNormalMap,v.clearcoatNormalMapTransform),v.clearcoatNormalScale.value.copy(h.clearcoatNormalScale),h.side===wn&&v.clearcoatNormalScale.value.negate())),h.dispersion>0&&(v.dispersion.value=h.dispersion),h.iridescence>0&&(v.iridescence.value=h.iridescence,v.iridescenceIOR.value=h.iridescenceIOR,v.iridescenceThicknessMinimum.value=h.iridescenceThicknessRange[0],v.iridescenceThicknessMaximum.value=h.iridescenceThicknessRange[1],h.iridescenceMap&&(v.iridescenceMap.value=h.iridescenceMap,n(h.iridescenceMap,v.iridescenceMapTransform)),h.iridescenceThicknessMap&&(v.iridescenceThicknessMap.value=h.iridescenceThicknessMap,n(h.iridescenceThicknessMap,v.iridescenceThicknessMapTransform))),h.transmission>0&&(v.transmission.value=h.transmission,v.transmissionSamplerMap.value=m.texture,v.transmissionSamplerSize.value.set(m.width,m.height),h.transmissionMap&&(v.transmissionMap.value=h.transmissionMap,n(h.transmissionMap,v.transmissionMapTransform)),v.thickness.value=h.thickness,h.thicknessMap&&(v.thicknessMap.value=h.thicknessMap,n(h.thicknessMap,v.thicknessMapTransform)),v.attenuationDistance.value=h.attenuationDistance,v.attenuationColor.value.copy(h.attenuationColor)),h.anisotropy>0&&(v.anisotropyVector.value.set(h.anisotropy*Math.cos(h.anisotropyRotation),h.anisotropy*Math.sin(h.anisotropyRotation)),h.anisotropyMap&&(v.anisotropyMap.value=h.anisotropyMap,n(h.anisotropyMap,v.anisotropyMapTransform))),v.specularIntensity.value=h.specularIntensity,v.specularColor.value.copy(h.specularColor),h.specularColorMap&&(v.specularColorMap.value=h.specularColorMap,n(h.specularColorMap,v.specularColorMapTransform)),h.specularIntensityMap&&(v.specularIntensityMap.value=h.specularIntensityMap,n(h.specularIntensityMap,v.specularIntensityMapTransform))}function _(v,h){h.matcap&&(v.matcap.value=h.matcap)}function b(v,h){let m=t.get(h).light;v.referencePosition.value.setFromMatrixPosition(m.matrixWorld),v.nearDistance.value=m.shadow.camera.near,v.farDistance.value=m.shadow.camera.far}return{refreshFogUniforms:i,refreshMaterialUniforms:a}}function bN(e,t,n,i){let a={},s={},r=[],o=e.getParameter(e.MAX_UNIFORM_BUFFER_BINDINGS);function l(m,g){let x=g.program;i.uniformBlockBinding(m,x)}function c(m,g){let x=a[m.id];x===void 0&&(_(m),x=f(m),a[m.id]=x,m.addEventListener("dispose",v));let E=g.program;i.updateUBOMapping(m,E);let A=t.render.frame;s[m.id]!==A&&(u(m),s[m.id]=A)}function f(m){let g=p();m.__bindingPointIndex=g;let x=e.createBuffer(),E=m.__size,A=m.usage;return e.bindBuffer(e.UNIFORM_BUFFER,x),e.bufferData(e.UNIFORM_BUFFER,E,A),e.bindBuffer(e.UNIFORM_BUFFER,null),e.bindBufferBase(e.UNIFORM_BUFFER,g,x),x}function p(){for(let m=0;m<o;m++)if(r.indexOf(m)===-1)return r.push(m),m;return zt("WebGLRenderer: Maximum number of simultaneously usable uniforms groups reached."),0}function u(m){let g=a[m.id],x=m.uniforms,E=m.__cache;e.bindBuffer(e.UNIFORM_BUFFER,g);for(let A=0,w=x.length;A<w;A++){let y=Array.isArray(x[A])?x[A]:[x[A]];for(let T=0,L=y.length;T<L;T++){let R=y[T];if(d(R,A,T,E)===!0){let N=R.__offset,P=Array.isArray(R.value)?R.value:[R.value],z=0;for(let F=0;F<P.length;F++){let V=P[F],X=b(V);typeof V=="number"||typeof V=="boolean"?(R.__data[0]=V,e.bufferSubData(e.UNIFORM_BUFFER,N+z,R.__data)):V.isMatrix3?(R.__data[0]=V.elements[0],R.__data[1]=V.elements[1],R.__data[2]=V.elements[2],R.__data[3]=0,R.__data[4]=V.elements[3],R.__data[5]=V.elements[4],R.__data[6]=V.elements[5],R.__data[7]=0,R.__data[8]=V.elements[6],R.__data[9]=V.elements[7],R.__data[10]=V.elements[8],R.__data[11]=0):(V.toArray(R.__data,z),z+=X.storage/Float32Array.BYTES_PER_ELEMENT)}e.bufferSubData(e.UNIFORM_BUFFER,N,R.__data)}}}e.bindBuffer(e.UNIFORM_BUFFER,null)}function d(m,g,x,E){let A=m.value,w=g+"_"+x;if(E[w]===void 0)return typeof A=="number"||typeof A=="boolean"?E[w]=A:E[w]=A.clone(),!0;{let y=E[w];if(typeof A=="number"||typeof A=="boolean"){if(y!==A)return E[w]=A,!0}else if(y.equals(A)===!1)return y.copy(A),!0}return!1}function _(m){let g=m.uniforms,x=0,E=16;for(let w=0,y=g.length;w<y;w++){let T=Array.isArray(g[w])?g[w]:[g[w]];for(let L=0,R=T.length;L<R;L++){let N=T[L],P=Array.isArray(N.value)?N.value:[N.value];for(let z=0,F=P.length;z<F;z++){let V=P[z],X=b(V),at=x%E,et=at%X.boundary,dt=at+et;x+=et,dt!==0&&E-dt<X.storage&&(x+=E-dt),N.__data=new Float32Array(X.storage/Float32Array.BYTES_PER_ELEMENT),N.__offset=x,x+=X.storage}}}let A=x%E;return A>0&&(x+=E-A),m.__size=x,m.__cache={},this}function b(m){let g={boundary:0,storage:0};return typeof m=="number"||typeof m=="boolean"?(g.boundary=4,g.storage=4):m.isVector2?(g.boundary=8,g.storage=8):m.isVector3||m.isColor?(g.boundary=16,g.storage=12):m.isVector4?(g.boundary=16,g.storage=16):m.isMatrix3?(g.boundary=48,g.storage=48):m.isMatrix4?(g.boundary=64,g.storage=64):m.isTexture?Ot("WebGLRenderer: Texture samplers can not be part of an uniforms group."):Ot("WebGLRenderer: Unsupported uniform value type.",m),g}function v(m){let g=m.target;g.removeEventListener("dispose",v);let x=r.indexOf(g.__bindingPointIndex);r.splice(x,1),e.deleteBuffer(a[g.id]),delete a[g.id],delete s[g.id]}function h(){for(let m in a)e.deleteBuffer(a[m]);r=[],a={},s={}}return{bind:l,update:c,dispose:h}}var MN=new Uint16Array([12469,15057,12620,14925,13266,14620,13807,14376,14323,13990,14545,13625,14713,13328,14840,12882,14931,12528,14996,12233,15039,11829,15066,11525,15080,11295,15085,10976,15082,10705,15073,10495,13880,14564,13898,14542,13977,14430,14158,14124,14393,13732,14556,13410,14702,12996,14814,12596,14891,12291,14937,11834,14957,11489,14958,11194,14943,10803,14921,10506,14893,10278,14858,9960,14484,14039,14487,14025,14499,13941,14524,13740,14574,13468,14654,13106,14743,12678,14818,12344,14867,11893,14889,11509,14893,11180,14881,10751,14852,10428,14812,10128,14765,9754,14712,9466,14764,13480,14764,13475,14766,13440,14766,13347,14769,13070,14786,12713,14816,12387,14844,11957,14860,11549,14868,11215,14855,10751,14825,10403,14782,10044,14729,9651,14666,9352,14599,9029,14967,12835,14966,12831,14963,12804,14954,12723,14936,12564,14917,12347,14900,11958,14886,11569,14878,11247,14859,10765,14828,10401,14784,10011,14727,9600,14660,9289,14586,8893,14508,8533,15111,12234,15110,12234,15104,12216,15092,12156,15067,12010,15028,11776,14981,11500,14942,11205,14902,10752,14861,10393,14812,9991,14752,9570,14682,9252,14603,8808,14519,8445,14431,8145,15209,11449,15208,11451,15202,11451,15190,11438,15163,11384,15117,11274,15055,10979,14994,10648,14932,10343,14871,9936,14803,9532,14729,9218,14645,8742,14556,8381,14461,8020,14365,7603,15273,10603,15272,10607,15267,10619,15256,10631,15231,10614,15182,10535,15118,10389,15042,10167,14963,9787,14883,9447,14800,9115,14710,8665,14615,8318,14514,7911,14411,7507,14279,7198,15314,9675,15313,9683,15309,9712,15298,9759,15277,9797,15229,9773,15166,9668,15084,9487,14995,9274,14898,8910,14800,8539,14697,8234,14590,7790,14479,7409,14367,7067,14178,6621,15337,8619,15337,8631,15333,8677,15325,8769,15305,8871,15264,8940,15202,8909,15119,8775,15022,8565,14916,8328,14804,8009,14688,7614,14569,7287,14448,6888,14321,6483,14088,6171,15350,7402,15350,7419,15347,7480,15340,7613,15322,7804,15287,7973,15229,8057,15148,8012,15046,7846,14933,7611,14810,7357,14682,7069,14552,6656,14421,6316,14251,5948,14007,5528,15356,5942,15356,5977,15353,6119,15348,6294,15332,6551,15302,6824,15249,7044,15171,7122,15070,7050,14949,6861,14818,6611,14679,6349,14538,6067,14398,5651,14189,5311,13935,4958,15359,4123,15359,4153,15356,4296,15353,4646,15338,5160,15311,5508,15263,5829,15188,6042,15088,6094,14966,6001,14826,5796,14678,5543,14527,5287,14377,4985,14133,4586,13869,4257,15360,1563,15360,1642,15358,2076,15354,2636,15341,3350,15317,4019,15273,4429,15203,4732,15105,4911,14981,4932,14836,4818,14679,4621,14517,4386,14359,4156,14083,3795,13808,3437,15360,122,15360,137,15358,285,15355,636,15344,1274,15322,2177,15281,2765,15215,3223,15120,3451,14995,3569,14846,3567,14681,3466,14511,3305,14344,3121,14037,2800,13753,2467,15360,0,15360,1,15359,21,15355,89,15346,253,15325,479,15287,796,15225,1148,15133,1492,15008,1749,14856,1882,14685,1886,14506,1783,14324,1608,13996,1398,13702,1183]),ya=null;function TN(){return ya===null&&(ya=new Kc(MN,16,16,Lr,_a),ya.name="DFG_LUT",ya.minFilter=En,ya.magFilter=En,ya.wrapS=ha,ya.wrapT=ha,ya.generateMipmaps=!1,ya.needsUpdate=!0),ya}var hp=class{constructor(t={}){let{canvas:n=_1(),context:i=null,depth:a=!0,stencil:s=!1,alpha:r=!1,antialias:o=!1,premultipliedAlpha:l=!0,preserveDrawingBuffer:c=!1,powerPreference:f="default",failIfMajorPerformanceCaveat:p=!1,reversedDepthBuffer:u=!1,outputBufferType:d=ei}=t;this.isWebGLRenderer=!0;let _;if(i!==null){if(typeof WebGLRenderingContext<"u"&&i instanceof WebGLRenderingContext)throw new Error("THREE.WebGLRenderer: WebGL 1 is not supported since r163.");_=i.getContextAttributes().alpha}else _=r;let b=d,v=new Set([wd,Ad,Ed]),h=new Set([ei,qi,pl,ml,bd,Md]),m=new Uint32Array(4),g=new Int32Array(4),x=null,E=null,A=[],w=[],y=null;this.domElement=n,this.debug={checkShaderErrors:!0,onShaderError:null},this.autoClear=!0,this.autoClearColor=!0,this.autoClearDepth=!0,this.autoClearStencil=!0,this.sortObjects=!0,this.clippingPlanes=[],this.localClippingEnabled=!1,this.toneMapping=Wi,this.toneMappingExposure=1,this.transmissionResolutionScale=1;let T=this,L=!1;this._outputColorSpace=On;let R=0,N=0,P=null,z=-1,F=null,V=new Ge,X=new Ge,at=null,et=new Ft(0),dt=0,St=n.width,ht=n.height,Vt=1,Ce=null,xe=null,K=new Ge(0,0,St,ht),j=new Ge(0,0,St,ht),ct=!1,It=new ol,Lt=!1,Pt=!1,Ye=new Ie,ie=new B,Jt=new Ge,de={background:null,fog:null,environment:null,overrideMaterial:null,isScene:!0},Xt=!1;function Te(){return P===null?Vt:1}let D=i;function ze(M,I){return n.getContext(M,I)}try{let M={alpha:!0,depth:a,stencil:s,antialias:o,premultipliedAlpha:l,preserveDrawingBuffer:c,powerPreference:f,failIfMajorPerformanceCaveat:p};if("setAttribute"in n&&n.setAttribute("data-engine",`three.js r${"183"}`),n.addEventListener("webglcontextlost",Et,!1),n.addEventListener("webglcontextrestored",Ht,!1),n.addEventListener("webglcontextcreationerror",ge,!1),D===null){let I="webgl2";if(D=ze(I,M),D===null)throw ze(I)?new Error("Error creating WebGL context with your selected attributes."):new Error("Error creating WebGL context.")}}catch(M){throw zt("WebGLRenderer: "+M.message),M}let re,ue,At,C,S,O,Z,Q,Y,Tt,lt,Rt,Ut,tt,rt,Mt,bt,gt,Wt,U,ut,ot,xt;function nt(){re=new UR(D),re.init(),ut=new vN(D,re),ue=new TR(D,re,t,ut),At=new mN(D,re),ue.reversedDepthBuffer&&u&&At.buffers.depth.setReversed(!0),C=new OR(D),S=new eN,O=new gN(D,re,At,S,ue,ut,C),Z=new NR(T),Q=new VC(D),ot=new bR(D,Q),Y=new LR(D,Q,C,ot),Tt=new BR(D,Y,Q,ot,C),gt=new PR(D,ue,O),rt=new ER(S),lt=new tN(T,Z,re,ue,ot,rt),Rt=new SN(T,S),Ut=new iN,tt=new cN(re),bt=new SR(T,Z,At,Tt,_,l),Mt=new pN(T,Tt,ue),xt=new bN(D,C,ue,At),Wt=new MR(D,re,C),U=new IR(D,re,C),C.programs=lt.programs,T.capabilities=ue,T.extensions=re,T.properties=S,T.renderLists=Ut,T.shadowMap=Mt,T.state=At,T.info=C}nt(),b!==ei&&(y=new FR(b,n.width,n.height,a,s));let q=new Vv(T,D);this.xr=q,this.getContext=function(){return D},this.getContextAttributes=function(){return D.getContextAttributes()},this.forceContextLoss=function(){let M=re.get("WEBGL_lose_context");M&&M.loseContext()},this.forceContextRestore=function(){let M=re.get("WEBGL_lose_context");M&&M.restoreContext()},this.getPixelRatio=function(){return Vt},this.setPixelRatio=function(M){M!==void 0&&(Vt=M,this.setSize(St,ht,!1))},this.getSize=function(M){return M.set(St,ht)},this.setSize=function(M,I,W=!0){if(q.isPresenting){Ot("WebGLRenderer: Can't change size while VR device is presenting.");return}St=M,ht=I,n.width=Math.floor(M*Vt),n.height=Math.floor(I*Vt),W===!0&&(n.style.width=M+"px",n.style.height=I+"px"),y!==null&&y.setSize(n.width,n.height),this.setViewport(0,0,M,I)},this.getDrawingBufferSize=function(M){return M.set(St*Vt,ht*Vt).floor()},this.setDrawingBufferSize=function(M,I,W){St=M,ht=I,Vt=W,n.width=Math.floor(M*W),n.height=Math.floor(I*W),this.setViewport(0,0,M,I)},this.setEffects=function(M){if(b===ei){console.error("THREE.WebGLRenderer: setEffects() requires outputBufferType set to HalfFloatType or FloatType.");return}if(M){for(let I=0;I<M.length;I++)if(M[I].isOutputPass===!0){console.warn("THREE.WebGLRenderer: OutputPass is not needed in setEffects(). Tone mapping and color space conversion are applied automatically.");break}}y.setEffects(M||[])},this.getCurrentViewport=function(M){return M.copy(V)},this.getViewport=function(M){return M.copy(K)},this.setViewport=function(M,I,W,k){M.isVector4?K.set(M.x,M.y,M.z,M.w):K.set(M,I,W,k),At.viewport(V.copy(K).multiplyScalar(Vt).round())},this.getScissor=function(M){return M.copy(j)},this.setScissor=function(M,I,W,k){M.isVector4?j.set(M.x,M.y,M.z,M.w):j.set(M,I,W,k),At.scissor(X.copy(j).multiplyScalar(Vt).round())},this.getScissorTest=function(){return ct},this.setScissorTest=function(M){At.setScissorTest(ct=M)},this.setOpaqueSort=function(M){Ce=M},this.setTransparentSort=function(M){xe=M},this.getClearColor=function(M){return M.copy(bt.getClearColor())},this.setClearColor=function(){bt.setClearColor(...arguments)},this.getClearAlpha=function(){return bt.getClearAlpha()},this.setClearAlpha=function(){bt.setClearAlpha(...arguments)},this.clear=function(M=!0,I=!0,W=!0){let k=0;if(M){let G=!1;if(P!==null){let it=P.texture.format;G=v.has(it)}if(G){let it=P.texture.type,_t=h.has(it),pt=bt.getClearColor(),wt=bt.getClearAlpha(),Ct=pt.r,Gt=pt.g,qt=pt.b;_t?(m[0]=Ct,m[1]=Gt,m[2]=qt,m[3]=wt,D.clearBufferuiv(D.COLOR,0,m)):(g[0]=Ct,g[1]=Gt,g[2]=qt,g[3]=wt,D.clearBufferiv(D.COLOR,0,g))}else k|=D.COLOR_BUFFER_BIT}I&&(k|=D.DEPTH_BUFFER_BIT),W&&(k|=D.STENCIL_BUFFER_BIT,this.state.buffers.stencil.setMask(4294967295)),k!==0&&D.clear(k)},this.clearColor=function(){this.clear(!0,!1,!1)},this.clearDepth=function(){this.clear(!1,!0,!1)},this.clearStencil=function(){this.clear(!1,!1,!0)},this.dispose=function(){n.removeEventListener("webglcontextlost",Et,!1),n.removeEventListener("webglcontextrestored",Ht,!1),n.removeEventListener("webglcontextcreationerror",ge,!1),bt.dispose(),Ut.dispose(),tt.dispose(),S.dispose(),Z.dispose(),Tt.dispose(),ot.dispose(),xt.dispose(),lt.dispose(),q.dispose(),q.removeEventListener("sessionstart",Mu),q.removeEventListener("sessionend",Tu),Sa.stop()};function Et(M){M.preventDefault(),xv("WebGLRenderer: Context Lost."),L=!0}function Ht(){xv("WebGLRenderer: Context Restored."),L=!1;let M=C.autoReset,I=Mt.enabled,W=Mt.autoUpdate,k=Mt.needsUpdate,G=Mt.type;nt(),C.autoReset=M,Mt.enabled=I,Mt.autoUpdate=W,Mt.needsUpdate=k,Mt.type=G}function ge(M){zt("WebGLRenderer: A WebGL context could not be created. Reason: ",M.statusMessage)}function oe(M){let I=M.target;I.removeEventListener("dispose",oe),Fn(I)}function Fn(M){yi(M),S.remove(M)}function yi(M){let I=S.get(M).programs;I!==void 0&&(I.forEach(function(W){lt.releaseProgram(W)}),M.isShaderMaterial&&lt.releaseShaderCache(M))}this.renderBufferDirect=function(M,I,W,k,G,it){I===null&&(I=de);let _t=G.isMesh&&G.matrixWorld.determinant()<0,pt=st(M,I,W,k,G);At.setMaterial(k,_t);let wt=W.index,Ct=1;if(k.wireframe===!0){if(wt=Y.getWireframeAttribute(W),wt===void 0)return;Ct=2}let Gt=W.drawRange,qt=W.attributes.position,Nt=Gt.start*Ct,Se=(Gt.start+Gt.count)*Ct;it!==null&&(Nt=Math.max(Nt,it.start*Ct),Se=Math.min(Se,(it.start+it.count)*Ct)),wt!==null?(Nt=Math.max(Nt,0),Se=Math.min(Se,wt.count)):qt!=null&&(Nt=Math.max(Nt,0),Se=Math.min(Se,qt.count));let Ze=Se-Nt;if(Ze<0||Ze===1/0)return;ot.setup(G,k,pt,W,wt);let Xe,be=Wt;if(wt!==null&&(Xe=Q.get(wt),be=U,be.setIndex(Xe)),G.isMesh)k.wireframe===!0?(At.setLineWidth(k.wireframeLinewidth*Te()),be.setMode(D.LINES)):be.setMode(D.TRIANGLES);else if(G.isLine){let Rn=k.linewidth;Rn===void 0&&(Rn=1),At.setLineWidth(Rn*Te()),G.isLineSegments?be.setMode(D.LINES):G.isLineLoop?be.setMode(D.LINE_LOOP):be.setMode(D.LINE_STRIP)}else G.isPoints?be.setMode(D.POINTS):G.isSprite&&be.setMode(D.TRIANGLES);if(G.isBatchedMesh)if(G._multiDrawInstances!==null)Gc("WebGLRenderer: renderMultiDrawInstances has been deprecated and will be removed in r184. Append to renderMultiDraw arguments and use indirection."),be.renderMultiDrawInstances(G._multiDrawStarts,G._multiDrawCounts,G._multiDrawCount,G._multiDrawInstances);else if(re.get("WEBGL_multi_draw"))be.renderMultiDraw(G._multiDrawStarts,G._multiDrawCounts,G._multiDrawCount);else{let Rn=G._multiDrawStarts,Dt=G._multiDrawCounts,ni=G._multiDrawCount,le=wt?Q.get(wt).bytesPerElement:1,Li=S.get(k).currentProgram.getUniforms();for(let Yi=0;Yi<ni;Yi++)Li.setValue(D,"_gl_DrawID",Yi),be.render(Rn[Yi]/le,Dt[Yi])}else if(G.isInstancedMesh)be.renderInstances(Nt,Ze,G.count);else if(W.isInstancedBufferGeometry){let Rn=W._maxInstanceCount!==void 0?W._maxInstanceCount:1/0,Dt=Math.min(W.instanceCount,Rn);be.renderInstances(Nt,Ze,Dt)}else be.render(Nt,Ze)};function Ml(M,I,W){M.transparent===!0&&M.side===ga&&M.forceSinglePass===!1?(M.side=wn,M.needsUpdate=!0,Js(M,I,W),M.side=Wa,M.needsUpdate=!0,Js(M,I,W),M.side=ga):Js(M,I,W)}this.compile=function(M,I,W=null){W===null&&(W=M),E=tt.get(W),E.init(I),w.push(E),W.traverseVisible(function(G){G.isLight&&G.layers.test(I.layers)&&(E.pushLight(G),G.castShadow&&E.pushShadow(G))}),M!==W&&M.traverseVisible(function(G){G.isLight&&G.layers.test(I.layers)&&(E.pushLight(G),G.castShadow&&E.pushShadow(G))}),E.setupLights();let k=new Set;return M.traverse(function(G){if(!(G.isMesh||G.isPoints||G.isLine||G.isSprite))return;let it=G.material;if(it)if(Array.isArray(it))for(let _t=0;_t<it.length;_t++){let pt=it[_t];Ml(pt,W,G),k.add(pt)}else Ml(it,W,G),k.add(it)}),E=w.pop(),k},this.compileAsync=function(M,I,W=null){let k=this.compile(M,I,W);return new Promise(G=>{function it(){if(k.forEach(function(_t){S.get(_t).currentProgram.isReady()&&k.delete(_t)}),k.size===0){G(M);return}setTimeout(it,10)}re.get("KHR_parallel_shader_compile")!==null?it():setTimeout(it,10)})};let Qa=null;function _p(M){Qa&&Qa(M)}function Mu(){Sa.stop()}function Tu(){Sa.start()}let Sa=new Y1;Sa.setAnimationLoop(_p),typeof self<"u"&&Sa.setContext(self),this.setAnimationLoop=function(M){Qa=M,q.setAnimationLoop(M),M===null?Sa.stop():Sa.start()},q.addEventListener("sessionstart",Mu),q.addEventListener("sessionend",Tu),this.render=function(M,I){if(I!==void 0&&I.isCamera!==!0){zt("WebGLRenderer.render: camera is not an instance of THREE.Camera.");return}if(L===!0)return;let W=q.enabled===!0&&q.isPresenting===!0,k=y!==null&&(P===null||W)&&y.begin(T,P);if(M.matrixWorldAutoUpdate===!0&&M.updateMatrixWorld(),I.parent===null&&I.matrixWorldAutoUpdate===!0&&I.updateMatrixWorld(),q.enabled===!0&&q.isPresenting===!0&&(y===null||y.isCompositing()===!1)&&(q.cameraAutoUpdate===!0&&q.updateCamera(I),I=q.getCamera()),M.isScene===!0&&M.onBeforeRender(T,M,I,P),E=tt.get(M,w.length),E.init(I),w.push(E),Ye.multiplyMatrices(I.projectionMatrix,I.matrixWorldInverse),It.setFromProjectionMatrix(Ye,ki,I.reversedDepth),Pt=this.localClippingEnabled,Lt=rt.init(this.clippingPlanes,Pt),x=Ut.get(M,A.length),x.init(),A.push(x),q.enabled===!0&&q.isPresenting===!0){let _t=T.xr.getDepthSensingMesh();_t!==null&&Tl(_t,I,-1/0,T.sortObjects)}Tl(M,I,0,T.sortObjects),x.finish(),T.sortObjects===!0&&x.sort(Ce,xe),Xt=q.enabled===!1||q.isPresenting===!1||q.hasDepthSensing()===!1,Xt&&bt.addToRenderList(x,M),this.info.render.frame++,Lt===!0&&rt.beginShadows();let G=E.state.shadowsArray;if(Mt.render(G,M,I),Lt===!0&&rt.endShadows(),this.info.autoReset===!0&&this.info.reset(),(k&&y.hasRenderPass())===!1){let _t=x.opaque,pt=x.transmissive;if(E.setupLights(),I.isArrayCamera){let wt=I.cameras;if(pt.length>0)for(let Ct=0,Gt=wt.length;Ct<Gt;Ct++){let qt=wt[Ct];El(_t,pt,M,qt)}Xt&&bt.render(M);for(let Ct=0,Gt=wt.length;Ct<Gt;Ct++){let qt=wt[Ct];Eu(x,M,qt,qt.viewport)}}else pt.length>0&&El(_t,pt,M,I),Xt&&bt.render(M),Eu(x,M,I)}P!==null&&N===0&&(O.updateMultisampleRenderTarget(P),O.updateRenderTargetMipmap(P)),k&&y.end(T),M.isScene===!0&&M.onAfterRender(T,M,I),ot.resetDefaultState(),z=-1,F=null,w.pop(),w.length>0?(E=w[w.length-1],Lt===!0&&rt.setGlobalState(T.clippingPlanes,E.state.camera)):E=null,A.pop(),A.length>0?x=A[A.length-1]:x=null};function Tl(M,I,W,k){if(M.visible===!1)return;if(M.layers.test(I.layers)){if(M.isGroup)W=M.renderOrder;else if(M.isLOD)M.autoUpdate===!0&&M.update(I);else if(M.isLight)E.pushLight(M),M.castShadow&&E.pushShadow(M);else if(M.isSprite){if(!M.frustumCulled||It.intersectsSprite(M)){k&&Jt.setFromMatrixPosition(M.matrixWorld).applyMatrix4(Ye);let _t=Tt.update(M),pt=M.material;pt.visible&&x.push(M,_t,pt,W,Jt.z,null)}}else if((M.isMesh||M.isLine||M.isPoints)&&(!M.frustumCulled||It.intersectsObject(M))){let _t=Tt.update(M),pt=M.material;if(k&&(M.boundingSphere!==void 0?(M.boundingSphere===null&&M.computeBoundingSphere(),Jt.copy(M.boundingSphere.center)):(_t.boundingSphere===null&&_t.computeBoundingSphere(),Jt.copy(_t.boundingSphere.center)),Jt.applyMatrix4(M.matrixWorld).applyMatrix4(Ye)),Array.isArray(pt)){let wt=_t.groups;for(let Ct=0,Gt=wt.length;Ct<Gt;Ct++){let qt=wt[Ct],Nt=pt[qt.materialIndex];Nt&&Nt.visible&&x.push(M,_t,Nt,W,Jt.z,qt)}}else pt.visible&&x.push(M,_t,pt,W,Jt.z,null)}}let it=M.children;for(let _t=0,pt=it.length;_t<pt;_t++)Tl(it[_t],I,W,k)}function Eu(M,I,W,k){let{opaque:G,transmissive:it,transparent:_t}=M;E.setupLightsView(W),Lt===!0&&rt.setGlobalState(T.clippingPlanes,W),k&&At.viewport(V.copy(k)),G.length>0&&ba(G,I,W),it.length>0&&ba(it,I,W),_t.length>0&&ba(_t,I,W),At.buffers.depth.setTest(!0),At.buffers.depth.setMask(!0),At.buffers.color.setMask(!0),At.setPolygonOffset(!1)}function El(M,I,W,k){if((W.isScene===!0?W.overrideMaterial:null)!==null)return;if(E.state.transmissionRenderTarget[k.id]===void 0){let Nt=re.has("EXT_color_buffer_half_float")||re.has("EXT_color_buffer_float");E.state.transmissionRenderTarget[k.id]=new gi(1,1,{generateMipmaps:!0,type:Nt?_a:ei,minFilter:Ws,samples:Math.max(4,ue.samples),stencilBuffer:s,resolveDepthBuffer:!1,resolveStencilBuffer:!1,colorSpace:se.workingColorSpace})}let it=E.state.transmissionRenderTarget[k.id],_t=k.viewport||V;it.setSize(_t.z*T.transmissionResolutionScale,_t.w*T.transmissionResolutionScale);let pt=T.getRenderTarget(),wt=T.getActiveCubeFace(),Ct=T.getActiveMipmapLevel();T.setRenderTarget(it),T.getClearColor(et),dt=T.getClearAlpha(),dt<1&&T.setClearColor(16777215,.5),T.clear(),Xt&&bt.render(W);let Gt=T.toneMapping;T.toneMapping=Wi;let qt=k.viewport;if(k.viewport!==void 0&&(k.viewport=void 0),E.setupLightsView(k),Lt===!0&&rt.setGlobalState(T.clippingPlanes,k),ba(M,W,k),O.updateMultisampleRenderTarget(it),O.updateRenderTargetMipmap(it),re.has("WEBGL_multisampled_render_to_texture")===!1){let Nt=!1;for(let Se=0,Ze=I.length;Se<Ze;Se++){let Xe=I[Se],{object:be,geometry:Rn,material:Dt,group:ni}=Xe;if(Dt.side===ga&&be.layers.test(k.layers)){let le=Dt.side;Dt.side=wn,Dt.needsUpdate=!0,Fr(be,W,k,Rn,Dt,ni),Dt.side=le,Dt.needsUpdate=!0,Nt=!0}}Nt===!0&&(O.updateMultisampleRenderTarget(it),O.updateRenderTargetMipmap(it))}T.setRenderTarget(pt,wt,Ct),T.setClearColor(et,dt),qt!==void 0&&(k.viewport=qt),T.toneMapping=Gt}function ba(M,I,W){let k=I.isScene===!0?I.overrideMaterial:null;for(let G=0,it=M.length;G<it;G++){let _t=M[G],{object:pt,geometry:wt,group:Ct}=_t,Gt=_t.material;Gt.allowOverride===!0&&k!==null&&(Gt=k),pt.layers.test(W.layers)&&Fr(pt,I,W,wt,Gt,Ct)}}function Fr(M,I,W,k,G,it){M.onBeforeRender(T,I,W,k,G,it),M.modelViewMatrix.multiplyMatrices(W.matrixWorldInverse,M.matrixWorld),M.normalMatrix.getNormalMatrix(M.modelViewMatrix),G.onBeforeRender(T,I,W,k,M,it),G.transparent===!0&&G.side===ga&&G.forceSinglePass===!1?(G.side=wn,G.needsUpdate=!0,T.renderBufferDirect(W,I,k,G,M,it),G.side=Wa,G.needsUpdate=!0,T.renderBufferDirect(W,I,k,G,M,it),G.side=ga):T.renderBufferDirect(W,I,k,G,M,it),M.onAfterRender(T,I,W,k,G,it)}function Js(M,I,W){I.isScene!==!0&&(I=de);let k=S.get(M),G=E.state.lights,it=E.state.shadowsArray,_t=G.state.version,pt=lt.getParameters(M,G.state,it,I,W),wt=lt.getProgramCacheKey(pt),Ct=k.programs;k.environment=M.isMeshStandardMaterial||M.isMeshLambertMaterial||M.isMeshPhongMaterial?I.environment:null,k.fog=I.fog;let Gt=M.isMeshStandardMaterial||M.isMeshLambertMaterial&&!M.envMap||M.isMeshPhongMaterial&&!M.envMap;k.envMap=Z.get(M.envMap||k.environment,Gt),k.envMapRotation=k.environment!==null&&M.envMap===null?I.environmentRotation:M.envMapRotation,Ct===void 0&&(M.addEventListener("dispose",oe),Ct=new Map,k.programs=Ct);let qt=Ct.get(wt);if(qt!==void 0){if(k.currentProgram===qt&&k.lightsStateVersion===_t)return $(M,pt),qt}else pt.uniforms=lt.getUniforms(M),M.onBeforeCompile(pt,T),qt=lt.acquireProgram(pt,wt),Ct.set(wt,qt),k.uniforms=pt.uniforms;let Nt=k.uniforms;return(!M.isShaderMaterial&&!M.isRawShaderMaterial||M.clipping===!0)&&(Nt.clippingPlanes=rt.uniform),$(M,pt),k.needsLights=ae(M),k.lightsStateVersion=_t,k.needsLights&&(Nt.ambientLightColor.value=G.state.ambient,Nt.lightProbe.value=G.state.probe,Nt.directionalLights.value=G.state.directional,Nt.directionalLightShadows.value=G.state.directionalShadow,Nt.spotLights.value=G.state.spot,Nt.spotLightShadows.value=G.state.spotShadow,Nt.rectAreaLights.value=G.state.rectArea,Nt.ltc_1.value=G.state.rectAreaLTC1,Nt.ltc_2.value=G.state.rectAreaLTC2,Nt.pointLights.value=G.state.point,Nt.pointLightShadows.value=G.state.pointShadow,Nt.hemisphereLights.value=G.state.hemi,Nt.directionalShadowMatrix.value=G.state.directionalShadowMatrix,Nt.spotLightMatrix.value=G.state.spotLightMatrix,Nt.spotLightMap.value=G.state.spotLightMap,Nt.pointShadowMatrix.value=G.state.pointShadowMatrix),k.currentProgram=qt,k.uniformsList=null,qt}function H(M){if(M.uniformsList===null){let I=M.currentProgram.getUniforms();M.uniformsList=yl.seqWithValue(I.seq,M.uniforms)}return M.uniformsList}function $(M,I){let W=S.get(M);W.outputColorSpace=I.outputColorSpace,W.batching=I.batching,W.batchingColor=I.batchingColor,W.instancing=I.instancing,W.instancingColor=I.instancingColor,W.instancingMorph=I.instancingMorph,W.skinning=I.skinning,W.morphTargets=I.morphTargets,W.morphNormals=I.morphNormals,W.morphColors=I.morphColors,W.morphTargetsCount=I.morphTargetsCount,W.numClippingPlanes=I.numClippingPlanes,W.numIntersection=I.numClipIntersection,W.vertexAlphas=I.vertexAlphas,W.vertexTangents=I.vertexTangents,W.toneMapping=I.toneMapping}function st(M,I,W,k,G){I.isScene!==!0&&(I=de),O.resetTextureUnits();let it=I.fog,_t=k.isMeshStandardMaterial||k.isMeshLambertMaterial||k.isMeshPhongMaterial?I.environment:null,pt=P===null?T.outputColorSpace:P.isXRRenderTarget===!0?P.texture.colorSpace:Rr,wt=k.isMeshStandardMaterial||k.isMeshLambertMaterial&&!k.envMap||k.isMeshPhongMaterial&&!k.envMap,Ct=Z.get(k.envMap||_t,wt),Gt=k.vertexColors===!0&&!!W.attributes.color&&W.attributes.color.itemSize===4,qt=!!W.attributes.tangent&&(!!k.normalMap||k.anisotropy>0),Nt=!!W.morphAttributes.position,Se=!!W.morphAttributes.normal,Ze=!!W.morphAttributes.color,Xe=Wi;k.toneMapped&&(P===null||P.isXRRenderTarget===!0)&&(Xe=T.toneMapping);let be=W.morphAttributes.position||W.morphAttributes.normal||W.morphAttributes.color,Rn=be!==void 0?be.length:0,Dt=S.get(k),ni=E.state.lights;if(Lt===!0&&(Pt===!0||M!==F)){let ln=M===F&&k.id===z;rt.setState(k,M,ln)}let le=!1;k.version===Dt.__version?(Dt.needsLights&&Dt.lightsStateVersion!==ni.state.version||Dt.outputColorSpace!==pt||G.isBatchedMesh&&Dt.batching===!1||!G.isBatchedMesh&&Dt.batching===!0||G.isBatchedMesh&&Dt.batchingColor===!0&&G.colorTexture===null||G.isBatchedMesh&&Dt.batchingColor===!1&&G.colorTexture!==null||G.isInstancedMesh&&Dt.instancing===!1||!G.isInstancedMesh&&Dt.instancing===!0||G.isSkinnedMesh&&Dt.skinning===!1||!G.isSkinnedMesh&&Dt.skinning===!0||G.isInstancedMesh&&Dt.instancingColor===!0&&G.instanceColor===null||G.isInstancedMesh&&Dt.instancingColor===!1&&G.instanceColor!==null||G.isInstancedMesh&&Dt.instancingMorph===!0&&G.morphTexture===null||G.isInstancedMesh&&Dt.instancingMorph===!1&&G.morphTexture!==null||Dt.envMap!==Ct||k.fog===!0&&Dt.fog!==it||Dt.numClippingPlanes!==void 0&&(Dt.numClippingPlanes!==rt.numPlanes||Dt.numIntersection!==rt.numIntersection)||Dt.vertexAlphas!==Gt||Dt.vertexTangents!==qt||Dt.morphTargets!==Nt||Dt.morphNormals!==Se||Dt.morphColors!==Ze||Dt.toneMapping!==Xe||Dt.morphTargetsCount!==Rn)&&(le=!0):(le=!0,Dt.__version=k.version);let Li=Dt.currentProgram;le===!0&&(Li=Js(k,I,G));let Yi=!1,Ks=!1,Vr=!1,Ee=Li.getUniforms(),vn=Dt.uniforms;if(At.useProgram(Li.program)&&(Yi=!0,Ks=!0,Vr=!0),k.id!==z&&(z=k.id,Ks=!0),Yi||F!==M){At.buffers.depth.getReversed()&&M.reversedDepth!==!0&&(M._reversedDepth=!0,M.updateProjectionMatrix()),Ee.setValue(D,"projectionMatrix",M.projectionMatrix),Ee.setValue(D,"viewMatrix",M.matrixWorldInverse);let $a=Ee.map.cameraPosition;$a!==void 0&&$a.setValue(D,ie.setFromMatrixPosition(M.matrixWorld)),ue.logarithmicDepthBuffer&&Ee.setValue(D,"logDepthBufFC",2/(Math.log(M.far+1)/Math.LN2)),(k.isMeshPhongMaterial||k.isMeshToonMaterial||k.isMeshLambertMaterial||k.isMeshBasicMaterial||k.isMeshStandardMaterial||k.isShaderMaterial)&&Ee.setValue(D,"isOrthographic",M.isOrthographicCamera===!0),F!==M&&(F=M,Ks=!0,Vr=!0)}if(Dt.needsLights&&(ni.state.directionalShadowMap.length>0&&Ee.setValue(D,"directionalShadowMap",ni.state.directionalShadowMap,O),ni.state.spotShadowMap.length>0&&Ee.setValue(D,"spotShadowMap",ni.state.spotShadowMap,O),ni.state.pointShadowMap.length>0&&Ee.setValue(D,"pointShadowMap",ni.state.pointShadowMap,O)),G.isSkinnedMesh){Ee.setOptional(D,G,"bindMatrix"),Ee.setOptional(D,G,"bindMatrixInverse");let ln=G.skeleton;ln&&(ln.boneTexture===null&&ln.computeBoneTexture(),Ee.setValue(D,"boneTexture",ln.boneTexture,O))}G.isBatchedMesh&&(Ee.setOptional(D,G,"batchingTexture"),Ee.setValue(D,"batchingTexture",G._matricesTexture,O),Ee.setOptional(D,G,"batchingIdTexture"),Ee.setValue(D,"batchingIdTexture",G._indirectTexture,O),Ee.setOptional(D,G,"batchingColorTexture"),G._colorsTexture!==null&&Ee.setValue(D,"batchingColorTexture",G._colorsTexture,O));let ja=W.morphAttributes;if((ja.position!==void 0||ja.normal!==void 0||ja.color!==void 0)&&gt.update(G,W,Li),(Ks||Dt.receiveShadow!==G.receiveShadow)&&(Dt.receiveShadow=G.receiveShadow,Ee.setValue(D,"receiveShadow",G.receiveShadow)),(k.isMeshStandardMaterial||k.isMeshLambertMaterial||k.isMeshPhongMaterial)&&k.envMap===null&&I.environment!==null&&(vn.envMapIntensity.value=I.environmentIntensity),vn.dfgLUT!==void 0&&(vn.dfgLUT.value=TN()),Ks&&(Ee.setValue(D,"toneMappingExposure",T.toneMappingExposure),Dt.needsLights&&vt(vn,Vr),it&&k.fog===!0&&Rt.refreshFogUniforms(vn,it),Rt.refreshMaterialUniforms(vn,k,Vt,ht,E.state.transmissionRenderTarget[M.id]),yl.upload(D,H(Dt),vn,O)),k.isShaderMaterial&&k.uniformsNeedUpdate===!0&&(yl.upload(D,H(Dt),vn,O),k.uniformsNeedUpdate=!1),k.isSpriteMaterial&&Ee.setValue(D,"center",G.center),Ee.setValue(D,"modelViewMatrix",G.modelViewMatrix),Ee.setValue(D,"normalMatrix",G.normalMatrix),Ee.setValue(D,"modelMatrix",G.matrixWorld),k.isShaderMaterial||k.isRawShaderMaterial){let ln=k.uniformsGroups;for(let $a=0,Hr=ln.length;$a<Hr;$a++){let Wv=ln[$a];xt.update(Wv,Li),xt.bind(Wv,Li)}}return Li}function vt(M,I){M.ambientLightColor.needsUpdate=I,M.lightProbe.needsUpdate=I,M.directionalLights.needsUpdate=I,M.directionalLightShadows.needsUpdate=I,M.pointLights.needsUpdate=I,M.pointLightShadows.needsUpdate=I,M.spotLights.needsUpdate=I,M.spotLightShadows.needsUpdate=I,M.rectAreaLights.needsUpdate=I,M.hemisphereLights.needsUpdate=I}function ae(M){return M.isMeshLambertMaterial||M.isMeshToonMaterial||M.isMeshPhongMaterial||M.isMeshStandardMaterial||M.isShadowMaterial||M.isShaderMaterial&&M.lights===!0}this.getActiveCubeFace=function(){return R},this.getActiveMipmapLevel=function(){return N},this.getRenderTarget=function(){return P},this.setRenderTargetTextures=function(M,I,W){let k=S.get(M);k.__autoAllocateDepthBuffer=M.resolveDepthBuffer===!1,k.__autoAllocateDepthBuffer===!1&&(k.__useRenderToTexture=!1),S.get(M.texture).__webglTexture=I,S.get(M.depthTexture).__webglTexture=k.__autoAllocateDepthBuffer?void 0:W,k.__hasExternalTextures=!0},this.setRenderTargetFramebuffer=function(M,I){let W=S.get(M);W.__webglFramebuffer=I,W.__useDefaultFramebuffer=I===void 0};let Pe=D.createFramebuffer();this.setRenderTarget=function(M,I=0,W=0){P=M,R=I,N=W;let k=null,G=!1,it=!1;if(M){let pt=S.get(M);if(pt.__useDefaultFramebuffer!==void 0){At.bindFramebuffer(D.FRAMEBUFFER,pt.__webglFramebuffer),V.copy(M.viewport),X.copy(M.scissor),at=M.scissorTest,At.viewport(V),At.scissor(X),At.setScissorTest(at),z=-1;return}else if(pt.__webglFramebuffer===void 0)O.setupRenderTarget(M);else if(pt.__hasExternalTextures)O.rebindTextures(M,S.get(M.texture).__webglTexture,S.get(M.depthTexture).__webglTexture);else if(M.depthBuffer){let Gt=M.depthTexture;if(pt.__boundDepthTexture!==Gt){if(Gt!==null&&S.has(Gt)&&(M.width!==Gt.image.width||M.height!==Gt.image.height))throw new Error("WebGLRenderTarget: Attached DepthTexture is initialized to the incorrect size.");O.setupDepthRenderbuffer(M)}}let wt=M.texture;(wt.isData3DTexture||wt.isDataArrayTexture||wt.isCompressedArrayTexture)&&(it=!0);let Ct=S.get(M).__webglFramebuffer;M.isWebGLCubeRenderTarget?(Array.isArray(Ct[I])?k=Ct[I][W]:k=Ct[I],G=!0):M.samples>0&&O.useMultisampledRTT(M)===!1?k=S.get(M).__webglMultisampledFramebuffer:Array.isArray(Ct)?k=Ct[W]:k=Ct,V.copy(M.viewport),X.copy(M.scissor),at=M.scissorTest}else V.copy(K).multiplyScalar(Vt).floor(),X.copy(j).multiplyScalar(Vt).floor(),at=ct;if(W!==0&&(k=Pe),At.bindFramebuffer(D.FRAMEBUFFER,k)&&At.drawBuffers(M,k),At.viewport(V),At.scissor(X),At.setScissorTest(at),G){let pt=S.get(M.texture);D.framebufferTexture2D(D.FRAMEBUFFER,D.COLOR_ATTACHMENT0,D.TEXTURE_CUBE_MAP_POSITIVE_X+I,pt.__webglTexture,W)}else if(it){let pt=I;for(let wt=0;wt<M.textures.length;wt++){let Ct=S.get(M.textures[wt]);D.framebufferTextureLayer(D.FRAMEBUFFER,D.COLOR_ATTACHMENT0+wt,Ct.__webglTexture,W,pt)}}else if(M!==null&&W!==0){let pt=S.get(M.texture);D.framebufferTexture2D(D.FRAMEBUFFER,D.COLOR_ATTACHMENT0,D.TEXTURE_2D,pt.__webglTexture,W)}z=-1},this.readRenderTargetPixels=function(M,I,W,k,G,it,_t,pt=0){if(!(M&&M.isWebGLRenderTarget)){zt("WebGLRenderer.readRenderTargetPixels: renderTarget is not THREE.WebGLRenderTarget.");return}let wt=S.get(M).__webglFramebuffer;if(M.isWebGLCubeRenderTarget&&_t!==void 0&&(wt=wt[_t]),wt){At.bindFramebuffer(D.FRAMEBUFFER,wt);try{let Ct=M.textures[pt],Gt=Ct.format,qt=Ct.type;if(M.textures.length>1&&D.readBuffer(D.COLOR_ATTACHMENT0+pt),!ue.textureFormatReadable(Gt)){zt("WebGLRenderer.readRenderTargetPixels: renderTarget is not in RGBA or implementation defined format.");return}if(!ue.textureTypeReadable(qt)){zt("WebGLRenderer.readRenderTargetPixels: renderTarget is not in UnsignedByteType or implementation defined type.");return}I>=0&&I<=M.width-k&&W>=0&&W<=M.height-G&&D.readPixels(I,W,k,G,ut.convert(Gt),ut.convert(qt),it)}finally{let Ct=P!==null?S.get(P).__webglFramebuffer:null;At.bindFramebuffer(D.FRAMEBUFFER,Ct)}}},this.readRenderTargetPixelsAsync=async function(M,I,W,k,G,it,_t,pt=0){if(!(M&&M.isWebGLRenderTarget))throw new Error("THREE.WebGLRenderer.readRenderTargetPixels: renderTarget is not THREE.WebGLRenderTarget.");let wt=S.get(M).__webglFramebuffer;if(M.isWebGLCubeRenderTarget&&_t!==void 0&&(wt=wt[_t]),wt)if(I>=0&&I<=M.width-k&&W>=0&&W<=M.height-G){At.bindFramebuffer(D.FRAMEBUFFER,wt);let Ct=M.textures[pt],Gt=Ct.format,qt=Ct.type;if(M.textures.length>1&&D.readBuffer(D.COLOR_ATTACHMENT0+pt),!ue.textureFormatReadable(Gt))throw new Error("THREE.WebGLRenderer.readRenderTargetPixelsAsync: renderTarget is not in RGBA or implementation defined format.");if(!ue.textureTypeReadable(qt))throw new Error("THREE.WebGLRenderer.readRenderTargetPixelsAsync: renderTarget is not in UnsignedByteType or implementation defined type.");let Nt=D.createBuffer();D.bindBuffer(D.PIXEL_PACK_BUFFER,Nt),D.bufferData(D.PIXEL_PACK_BUFFER,it.byteLength,D.STREAM_READ),D.readPixels(I,W,k,G,ut.convert(Gt),ut.convert(qt),0);let Se=P!==null?S.get(P).__webglFramebuffer:null;At.bindFramebuffer(D.FRAMEBUFFER,Se);let Ze=D.fenceSync(D.SYNC_GPU_COMMANDS_COMPLETE,0);return D.flush(),await x1(D,Ze,4),D.bindBuffer(D.PIXEL_PACK_BUFFER,Nt),D.getBufferSubData(D.PIXEL_PACK_BUFFER,0,it),D.deleteBuffer(Nt),D.deleteSync(Ze),it}else throw new Error("THREE.WebGLRenderer.readRenderTargetPixelsAsync: requested read bounds are out of range.")},this.copyFramebufferToTexture=function(M,I=null,W=0){let k=Math.pow(2,-W),G=Math.floor(M.image.width*k),it=Math.floor(M.image.height*k),_t=I!==null?I.x:0,pt=I!==null?I.y:0;O.setTexture2D(M,0),D.copyTexSubImage2D(D.TEXTURE_2D,W,0,0,_t,pt,G,it),At.unbindTexture()};let ke=D.createFramebuffer(),Re=D.createFramebuffer();this.copyTextureToTexture=function(M,I,W=null,k=null,G=0,it=0){let _t,pt,wt,Ct,Gt,qt,Nt,Se,Ze,Xe=M.isCompressedTexture?M.mipmaps[it]:M.image;if(W!==null)_t=W.max.x-W.min.x,pt=W.max.y-W.min.y,wt=W.isBox3?W.max.z-W.min.z:1,Ct=W.min.x,Gt=W.min.y,qt=W.isBox3?W.min.z:0;else{let vn=Math.pow(2,-G);_t=Math.floor(Xe.width*vn),pt=Math.floor(Xe.height*vn),M.isDataArrayTexture?wt=Xe.depth:M.isData3DTexture?wt=Math.floor(Xe.depth*vn):wt=1,Ct=0,Gt=0,qt=0}k!==null?(Nt=k.x,Se=k.y,Ze=k.z):(Nt=0,Se=0,Ze=0);let be=ut.convert(I.format),Rn=ut.convert(I.type),Dt;I.isData3DTexture?(O.setTexture3D(I,0),Dt=D.TEXTURE_3D):I.isDataArrayTexture||I.isCompressedArrayTexture?(O.setTexture2DArray(I,0),Dt=D.TEXTURE_2D_ARRAY):(O.setTexture2D(I,0),Dt=D.TEXTURE_2D),D.pixelStorei(D.UNPACK_FLIP_Y_WEBGL,I.flipY),D.pixelStorei(D.UNPACK_PREMULTIPLY_ALPHA_WEBGL,I.premultiplyAlpha),D.pixelStorei(D.UNPACK_ALIGNMENT,I.unpackAlignment);let ni=D.getParameter(D.UNPACK_ROW_LENGTH),le=D.getParameter(D.UNPACK_IMAGE_HEIGHT),Li=D.getParameter(D.UNPACK_SKIP_PIXELS),Yi=D.getParameter(D.UNPACK_SKIP_ROWS),Ks=D.getParameter(D.UNPACK_SKIP_IMAGES);D.pixelStorei(D.UNPACK_ROW_LENGTH,Xe.width),D.pixelStorei(D.UNPACK_IMAGE_HEIGHT,Xe.height),D.pixelStorei(D.UNPACK_SKIP_PIXELS,Ct),D.pixelStorei(D.UNPACK_SKIP_ROWS,Gt),D.pixelStorei(D.UNPACK_SKIP_IMAGES,qt);let Vr=M.isDataArrayTexture||M.isData3DTexture,Ee=I.isDataArrayTexture||I.isData3DTexture;if(M.isDepthTexture){let vn=S.get(M),ja=S.get(I),ln=S.get(vn.__renderTarget),$a=S.get(ja.__renderTarget);At.bindFramebuffer(D.READ_FRAMEBUFFER,ln.__webglFramebuffer),At.bindFramebuffer(D.DRAW_FRAMEBUFFER,$a.__webglFramebuffer);for(let Hr=0;Hr<wt;Hr++)Vr&&(D.framebufferTextureLayer(D.READ_FRAMEBUFFER,D.COLOR_ATTACHMENT0,S.get(M).__webglTexture,G,qt+Hr),D.framebufferTextureLayer(D.DRAW_FRAMEBUFFER,D.COLOR_ATTACHMENT0,S.get(I).__webglTexture,it,Ze+Hr)),D.blitFramebuffer(Ct,Gt,_t,pt,Nt,Se,_t,pt,D.DEPTH_BUFFER_BIT,D.NEAREST);At.bindFramebuffer(D.READ_FRAMEBUFFER,null),At.bindFramebuffer(D.DRAW_FRAMEBUFFER,null)}else if(G!==0||M.isRenderTargetTexture||S.has(M)){let vn=S.get(M),ja=S.get(I);At.bindFramebuffer(D.READ_FRAMEBUFFER,ke),At.bindFramebuffer(D.DRAW_FRAMEBUFFER,Re);for(let ln=0;ln<wt;ln++)Vr?D.framebufferTextureLayer(D.READ_FRAMEBUFFER,D.COLOR_ATTACHMENT0,vn.__webglTexture,G,qt+ln):D.framebufferTexture2D(D.READ_FRAMEBUFFER,D.COLOR_ATTACHMENT0,D.TEXTURE_2D,vn.__webglTexture,G),Ee?D.framebufferTextureLayer(D.DRAW_FRAMEBUFFER,D.COLOR_ATTACHMENT0,ja.__webglTexture,it,Ze+ln):D.framebufferTexture2D(D.DRAW_FRAMEBUFFER,D.COLOR_ATTACHMENT0,D.TEXTURE_2D,ja.__webglTexture,it),G!==0?D.blitFramebuffer(Ct,Gt,_t,pt,Nt,Se,_t,pt,D.COLOR_BUFFER_BIT,D.NEAREST):Ee?D.copyTexSubImage3D(Dt,it,Nt,Se,Ze+ln,Ct,Gt,_t,pt):D.copyTexSubImage2D(Dt,it,Nt,Se,Ct,Gt,_t,pt);At.bindFramebuffer(D.READ_FRAMEBUFFER,null),At.bindFramebuffer(D.DRAW_FRAMEBUFFER,null)}else Ee?M.isDataTexture||M.isData3DTexture?D.texSubImage3D(Dt,it,Nt,Se,Ze,_t,pt,wt,be,Rn,Xe.data):I.isCompressedArrayTexture?D.compressedTexSubImage3D(Dt,it,Nt,Se,Ze,_t,pt,wt,be,Xe.data):D.texSubImage3D(Dt,it,Nt,Se,Ze,_t,pt,wt,be,Rn,Xe):M.isDataTexture?D.texSubImage2D(D.TEXTURE_2D,it,Nt,Se,_t,pt,be,Rn,Xe.data):M.isCompressedTexture?D.compressedTexSubImage2D(D.TEXTURE_2D,it,Nt,Se,Xe.width,Xe.height,be,Xe.data):D.texSubImage2D(D.TEXTURE_2D,it,Nt,Se,_t,pt,be,Rn,Xe);D.pixelStorei(D.UNPACK_ROW_LENGTH,ni),D.pixelStorei(D.UNPACK_IMAGE_HEIGHT,le),D.pixelStorei(D.UNPACK_SKIP_PIXELS,Li),D.pixelStorei(D.UNPACK_SKIP_ROWS,Yi),D.pixelStorei(D.UNPACK_SKIP_IMAGES,Ks),it===0&&I.generateMipmaps&&D.generateMipmap(Dt),At.unbindTexture()},this.initRenderTarget=function(M){S.get(M).__webglFramebuffer===void 0&&O.setupRenderTarget(M)},this.initTexture=function(M){M.isCubeTexture?O.setTextureCube(M,0):M.isData3DTexture?O.setTexture3D(M,0):M.isDataArrayTexture||M.isCompressedArrayTexture?O.setTexture2DArray(M,0):O.setTexture2D(M,0),At.unbindTexture()},this.resetState=function(){R=0,N=0,P=null,At.reset(),ot.reset()},typeof __THREE_DEVTOOLS__<"u"&&__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("observe",{detail:this}))}get coordinateSystem(){return ki}get outputColorSpace(){return this._outputColorSpace}set outputColorSpace(t){this._outputColorSpace=t;let n=this.getContext();n.drawingBufferColorSpace=se._getDrawingBufferColorSpace(t),n.unpackColorSpace=se._getUnpackColorSpace()}};var AN=[["#f9d392","#f2a754","#809f62","#356b52"],["#9aaac4","#fbeca9","#626d96","#333f6d"],["#f4cecd","#ffedb5","#96b870","#4c845c"],["#b2d4d3","#fbdf9d","#6aab9d","#356d74"]],zr=["#e1b07e","#91b0b5","#bf9fbd","#c8cd83"];function Hv(e,t,n,i){let[a,s,r,o]=AN[t.style];e.fillStyle=a,e.fillRect(0,0,n,i);let l=(c,f,p,u,d)=>{e.fillStyle=c,e.beginPath(),e.ellipse(f*n,p*i,u*n,d*i,0,0,2*Math.PI),e.fill()};if(l(s,t.style===1?.26:.7,.25,t.style===1?.075:.1,.14),l(r,.24,.83,.68,.36),l(o,.8,1.02,.68,.48),t.kind==="build"){let c=(t.level||1)+2;if(t.style===1)for(let f=0;f<8;f++){let p=Math.PI+f*Math.PI/7;l(zr[1],.5+Math.cos(p)*.25,.75+Math.sin(p)*.47,.095,.12)}else if(t.style===2)e.fillStyle=zr[0],e.fillRect(n*.36,i*(.75-c*.095),n*.28,i*c*.095),e.fillStyle=zr[2],e.fillRect(n*.31,i*(.72-c*.095),n*.38,i*.12);else for(let f=0;f<c;f++)l(zr[t.style],.52+Math.sin(f*1.9)*.03,.76-f*.1,.25-f*.032,.072);return}if(t.style===1||t.style===3){l(t.style===1?"#839ab0":"#b4d9c7",.46,.76,.34,.13);for(let c=0;c<4;c++)l(s,.34+c*.037,.71+c*.028,.058,.006);if(t.style===1)for(let c=0;c<12;c++)l("#faf0cf",.13+c*.173%.78,.08+c*.267%.35,.007,.009);else for(let c=0;c<6;c++)l("#e0eee2",.08+c*.11,.33+Math.sin(c*.7)*.05,.13,.041)}else for(let c=0;c<17;c++){let f=.06+c*.369%.87,p=.63+c*.173%.31;t.style===2?(l(o,f,p+.024,.006,.036),l(c%2===0?"#d76f8d":s,f,p,.02,.022)):l(c%3===0?s:r,f,p,.012,.015)}}function j1(e,t,n,i){let a=new hp({antialias:!0,alpha:!1,powerPreference:"low-power"});a.setPixelRatio(Math.min(devicePixelRatio,1.5)),a.outputColorSpace=On,a.toneMapping=hu,e.append(a.domElement);let s=a.domElement;s.tabIndex=0,s.setAttribute("aria-label","Clover Hollow 3D view. Focus to walk with WASD or arrow keys.");let r=new qc,o=new Ft("#c7e9ce"),l=new Wc(o,25,110);r.background=o,r.fog=l;let c=new Pn(55,1,.07,150);c.rotation.order="YXZ";let f=new lu("#ecfce8","#837b57",2),p=new uu("#ffebba",2.7);p.position.set(-11,24,12),r.add(f,p);let u=new Ft("#80c2e2"),d=new su(127,24,12),_=new Yn({side:wn,depthWrite:!1,uniforms:{upper:{value:u},horizon:{value:o}},vertexShader:`varying float height;
      void main(){height=position.y/127.0;gl_Position=projectionMatrix*modelViewMatrix*vec4(position,1.0);}`,fragmentShader:`uniform vec3 upper;uniform vec3 horizon;varying float height;
      void main(){gl_FragColor=vec4(mix(horizon,upper,smoothstep(-0.06,0.58,height)),1.0);
      #include <colorspace_fragment>
      }`}),b=new mn(d,_);b.renderOrder=-20,r.add(b);let v={globe:new fl(1,2),rock:new fl(1,0),branch:new cl(.55,1,1,8),block:new Vs(1,1,1),pine:new ul(1,1,8),plane:new iu(1,48)},h=new Map,m=new Set,g=(H,$,st,vt,ae,Pe,ke,Re,M=r)=>{let I=h.get($);I||(I=new Nr({color:$,roughness:.92,flatShading:H==="rock"||H==="pine"}),h.set($,I));let W=new mn(v[H],I);return W.position.set(st,vt,ae),W.scale.set(Pe,ke,Re),M.add(W),W},x=(H,$,st,vt,ae,Pe,ke=r)=>{let Re=g("plane",H,$,st,vt,ae,Pe,1,ke);return Re.rotation.x=-Math.PI/2,Re},E=(H,$,st,vt,ae=r)=>x("#526f35",H,.016,$,st,vt,ae),A=Number(BigInt(t().seed)%2147483647n)||7,w=()=>(A=A*48271%2147483647)/2147483647;g("block","#7fa64c",0,-.49,0,200,.95,200);for(let H=0;H<23;H++){let $=H*2.399,st=20+w()*39;x(H%2?"#84ab55":"#79a04c",Math.cos($)*st,.001+H*6e-5,Math.sin($)*st,3+w()*7,2+w()*5)}for(let H=0;H<21;H++){let $=H/20;x(H%3?"#b2b275":"#c0bb85",.8+8*$+Math.sin($*7)*2,.009+H*1e-4,9-17*$,.54+w()*.29,.72)}let y=(H,$,st,vt)=>{if(E(H+.35,$+.25,1.55*st,.95*st),g("branch","#66553d",H,1.4*st,$,.23*st,2.8*st,.22*st),vt){let ae=["#356749","#427754","#5b8b55"];for(let Pe=0;Pe<3;Pe++){let ke=(2.1-Pe*.46)*st,Re=g("pine",ae[Pe],H,(2.8+Pe*.93)*st,$,ke,2.65*st,ke);Re.rotation.y=Pe*1.12}return}g("globe","#4b8247",H,3.3*st,$,1.8*st,1.7*st,1.72*st),g("globe","#659851",H-.87*st,3.7*st,$-.25*st,1.18*st,1.12*st,1.17*st),g("globe","#87ad5b",H+.43*st,4.28*st,$+.12*st,1.15*st,1.2*st,1.02*st)};for(let H=0;H<52;H++){let $=(H+.45)*Math.PI*2/52,st=23+w()*20;y(Math.cos($)*st,Math.sin($)*st,.76+w()*.63,H%3!==0)}[[-10,-1,1],[-10,-13,1.2],[13,-12,1.1],[15,6,.83],[-4,-12,.74],[5,-12,.88],[-9,11,1],[-19,-1,1.15]].forEach(([H,$,st],vt)=>y(H,$,st,vt%3===0));for(let H=0;H<20;H++){let $=H*.38,st=59+w()*22,vt=g("pine",H%2?"#769c8e":"#8eae9a",Math.cos($)*st,2,Math.sin($)*st,11+w()*8,15+w()*18,12);vt.rotation.y=$}let T=new sn;r.add(T);for(let H=0;H<15;H++){let $=H*2.399,st=36+w()*40,vt=Math.cos($)*st,ae=Math.sin($)*st,Pe=14+w()*10;for(let ke=0;ke<4;ke++)g("globe","#f9f9de",vt+ke*2.15,Pe+Math.sin(ke*1.5),ae,3.4,1.2,1.9,T)}x("#5a916c",-15,.012,-9,5.07,4.67),x("#70c5bb",-15,.025,-9,4.48,4.03);let L=new sn;r.add(L);for(let H=0;H<10;H++){let $=H*2.399,st=Math.sqrt(w())*3.4;x("#a7dbb9",-15+Math.cos($)*st,.031,-9+Math.sin($)*st,.54,.045,L)}for(let H=0;H<12;H++){let $=H*.51;g("rock","#a0a38a",-15+Math.cos($)*4.6,.12,-9+Math.sin($)*4.25,.34,.25,.27)}let R=new ul(.06,.28,3),N=new Nr({color:"#7da34e"}),P=new jc(R,N,1150),z=new tn,F=new Ft,V=0;for(let H=0;H<430;H++){let $=(w()-.5)*51,st=(w()-.5)*44;if(!(Math.hypot($+15,st+9)<5.2||Math.hypot($-9,st+8)<3.4))for(let vt=0;vt<3&&V<1150;vt++){let ae=.45+w()*.65;z.position.set($+vt*.075,.12*ae,st+(w()-.5)*.12),z.scale.setScalar(ae),z.rotation.set(0,w()*6.28,(vt-1)*.34),z.updateMatrix(),P.setMatrixAt(V,z.matrix),P.setColorAt(V++,F.set(vt%2?"#d9ef9e":"#b6d97b"))}}P.count=V,r.add(P);let X=["#e5cd62","#f6efd2","#d78992","#a696c2"];for(let H=0;H<130;H++){let $=(w()-.5)*48,st=(w()-.5)*41;if(Math.hypot($+15,st+9)<5.3||Math.hypot($-9,st+8)<3.4)continue;let vt=.08+w()*.09,ae=X[H%X.length];g("branch","#4a7b3e",$,.1,st,.026,.24,.026),g("globe",ae,$,.21,st,vt,.068,vt)}E(9,-8.9,2.9,2.2),g("globe","#af845b",9,1.15,-9,2.42,1.83,1.87),g("globe","#5d8649",9,2.26,-9.22,2.56,.99,1.92),g("globe","#364837",9,.89,-7.43,.89,.91,.33),g("globe","#645644",9,.51,-7.13,.88,.32,.31),g("globe","#dbb27f",9,.17,-7.08,.91,.18,.63),g("block","#eccb87",9.08,.26,-7.01,1.19,.14,.72),g("globe","#e8cb76",7.49,1.23,-7.73,.31,.32,.22),g("globe","#e8cb76",10.56,1.17,-7.73,.31,.32,.22);for(let H=0;H<7;H++)g("rock","#83a158",7.2+H*.57,2.72+Math.sin(H*.7)*.33,-8.82,.27,.16,.31);let at=g("globe","#e9b05b",-6,.34,5,.34,.34,.34),et=g("globe","#d96653",-6,.34,5,.348,.08,.348);E(-6,5,.46,.38);let dt=new sn;dt.position.set(7,0,3.55),r.add(dt),E(0,0,1.1,.58,dt),g("block","#92764c",-.52,.85,0,.12,1.73,.13,dt).rotation.z=-.22,g("block","#92764c",.52,.85,0,.12,1.73,.13,dt).rotation.z=.22,g("block","#92764c",0,1.1,.02,1.3,.11,.23,dt),g("block","#e1c391",0,1.67,.04,1.32,1.05,.15,dt);let St=document.createElement("canvas");St.width=192,St.height=144;let ht=St.getContext("2d"),Vt=new eu(St);m.add(Vt),Vt.colorSpace=On;let Ce=g("block","#f8eed5",0,1.67,.125,1.19,.89,.021,dt),xe=new Nr({map:Vt,roughness:1});Ce.material=xe,h.set("artMaterial",xe),ht.fillStyle="#f8eed5",ht.fillRect(0,0,192,144),Vt.needsUpdate=!0;let K="",j=new sn,ct=new sn;r.add(j,ct),E(0,0,1.03,.69,ct);let It=new sn;j.add(It);let Lt="#d97d47",Pt="#fff0cc",Ye="#392b27",ie=g("globe",Lt,0,.62,-.14,.66,.67,.84,It);g("globe",Pt,0,.58,.5,.46,.48,.28,It);let Jt=new sn;Jt.position.set(0,1.21,.29),It.add(Jt),g("globe",Lt,0,0,0,.71,.6,.61,Jt),g("globe",Pt,0,-.24,.4,.55,.29,.28,Jt);let de=[],Xt=[];for(let H of[-1,1]){let $=new sn;Jt.add($),$.position.set(H*.48,.37,-.035),$.rotation.z=H*-.25,g("globe","#b85c35",0,.23,0,.2,.4,.19,$),g("globe","#efb58c",0,.23,.14,.12,.27,.07,$),de.push($),Xt.push(g("globe",Ye,H*.25,.035,.55,.073,.1,.047,Jt)),g("globe","#fff8e1",H*.25-.016,.075,.585,.023,.029,.015,Jt),g("globe","#eea277",H*.43,-.16,.48,.12,.065,.036,Jt)}g("globe",Ye,0,-.12,.683,.105,.066,.045,Jt);let Te=new sn;Te.position.set(0,.6,-.67),It.add(Te),g("globe","#b8653c",0,.24,-.37,.27,.32,.57,Te).rotation.x=.4,g("globe",Pt,0,.43,-.73,.23,.24,.25,Te);let D=[[-.44,.17,.36],[.44,.17,.36],[-.4,.16,-.54],[.4,.16,-.54]].map(([H,$,st])=>g("globe","#b3623b",H,$,st,.21,.24,.32,It)),ze=new sn;r.add(ze);let re=new Dr({color:"#fff2c8"});h.set("stars",re);for(let H=0;H<50;H++){let $=H*2.399,st=14+w()*58,vt=new mn(v.rock,re);vt.position.set(Math.cos($)*58,st,Math.sin($)*58),vt.scale.setScalar(.07+w()*.13),ze.add(vt)}let ue=new An,At=new Float32Array(672);ue.setAttribute("position",new Bn(At,3));let C=new ll({color:"#bde4da",transparent:!0,opacity:.49});h.set("rain",C);let S=new $c(ue,C);S.frustumCulled=!1,r.add(S);let O=Array.from({length:112},()=>[w()*25-12.5,w()*13,w()*25-12.5]),Z=new Map,Q=H=>{let $=new sn;if($.position.set(H.x,0,H.z),r.add($),H.kind==="build"){E(0,0,1.28,.89,$);let st=(H.level||1)+2;if(H.style===1)for(let vt=0;vt<8;vt++){let ae=Math.PI-vt*Math.PI/7;g("rock",zr[vt%2?1:3],Math.cos(ae)*.96,.24+Math.sin(ae)*(.98+st*.12),0,.34,.39,.36,$).rotation.z=ae}else if(H.style===2){for(let vt=0;vt<st;vt++)g("block","#ac794e",0,.19+vt*.37,0,.91-vt*.08,.32,.69-vt*.05,$).rotation.y=vt%2*.19;g("pine",zr[2],0,.43+st*.36,0,.61,.54,.61,$)}else for(let vt=0;vt<st;vt++){let ae=.91-vt*.13,Pe=zr[(H.style+(vt===0?1:0))%4];g("rock",Pe,Math.sin(vt*1.8)*.13,.27+vt*.4,0,ae,.3,ae*.78,$).rotation.y=vt*.65}}else{E(0,0,.43,.32,$);let st=H.kind==="berry"?"#a74e87":H.kind==="mushroom"?"#db785d":"#edab50";if(g("branch","#e6d3a0",0,.21,0,.12,.38,.12,$),H.kind==="berry")for(let vt=0;vt<4;vt++)g("globe",st,Math.cos(vt*1.7)*.15,.24+Math.sin(vt)*.11,Math.sin(vt*1.7)*.15,.16,.17,.15,$);else g("globe",st,0,.42,0,.28,.17,.28,$);H.kind!=="mushroom"&&(g("globe","#b3c66e",.13,.49,-.07,.18,.07,.09,$).rotation.z=.45)}return $},Y=-1,Tt=H=>{let $=[...H.food,...H.creations.filter(Re=>Re.kind==="build")],st=new Set($.map(Re=>Re.id));for(let[Re,M]of Z)st.has(Re)||(r.remove(M.group),Z.delete(Re));for(let Re of $){let M=Re.kind+Re.style+Re.level,I=Z.get(Re.id);I?.key!==M&&(I&&r.remove(I.group),Z.set(Re.id,{group:Q(Re),key:M}))}let vt=[...H.creations].reverse().find(Re=>Re.kind==="art"),ae=H.action==="draw"&&H.until>0,Pe=ae?Math.max(.12,1-(H.until-H.tick)/16):1,ke=ae?"draft-"+H.tick:String(vt?.id??"");ke!==K&&(ht.fillStyle="#f8eed5",ht.fillRect(0,0,192,144),(ae||vt)&&(ht.save(),ht.beginPath(),ht.rect(0,0,192*Pe,144),ht.clip(),Hv(ht,ae?{kind:"art",style:H.focusStyle,id:0,x:0,z:0}:vt,192,144),ht.restore()),Vt.needsUpdate=!0,K=ke)},lt=!0,Rt=0,Ut=0,tt=null,rt=0,Mt=0,bt=new Set,gt=new Set,Wt={w:"up",a:"left",s:"down",d:"right",ArrowUp:"up",ArrowLeft:"left",ArrowDown:"down",ArrowRight:"right"},U=H=>{lt!==H&&(lt=H,n(H))},ut=H=>{let $=Wt[H.key.length===1?H.key.toLowerCase():H.key];!$||H.altKey||H.ctrlKey||H.metaKey||(H.preventDefault(),bt.add($),U(!1))},ot=H=>bt.delete(Wt[H.key.length===1?H.key.toLowerCase():H.key]),xt=()=>{bt.clear(),gt.clear()},nt=()=>bt.clear(),q=H=>{tt!==null||H.button!==0||(tt=H.pointerId,rt=H.clientX,Mt=H.clientY,s.focus({preventScroll:!0}),s.setPointerCapture(H.pointerId))},Et=H=>{if(H.pointerId!==tt)return;let $=H.clientX-rt,st=H.clientY-Mt;!$&&!st||(U(!1),Rt-=$*.006,Ut=vl.clamp(Ut-st*.005,-1.18,1.05),rt=H.clientX,Mt=H.clientY)},Ht=H=>{H.pointerId===tt&&(tt=null)},ge=H=>{H.preventDefault(),i(),a.setAnimationLoop(null)};s.addEventListener("pointerdown",q),s.addEventListener("pointermove",Et),s.addEventListener("lostpointercapture",Ht),s.addEventListener("keydown",ut),s.addEventListener("keyup",ot),s.addEventListener("blur",nt),s.addEventListener("webglcontextlost",ge),window.addEventListener("blur",xt);let oe=new ResizeObserver(()=>{let H=e.clientWidth,$=e.clientHeight;!H||!$||(c.aspect=H/$,c.fov=c.aspect<1?63:55,c.updateProjectionMatrix(),a.setSize(H,$,!1))});oe.observe(e);let Fn=new B,yi=new B,Ml=new B,Qa=new B,_p=new Ft("#d4ecd8"),Mu=new Ft("#202d4b"),Tu=new Ft("#819d9b"),Sa=new Ft("#8ec9df"),Tl=new Ft("#101a33"),Eu=new Ft("#526c80"),El=(H,$,st)=>{Qa.set(H,0,$),Qa.lengthSq()&&(Qa.normalize().applyAxisAngle(tn.DEFAULT_UP,Rt).multiplyScalar(st),c.position.add(Qa),c.position.x=vl.clamp(c.position.x,-23,23),c.position.z=vl.clamp(c.position.z,-21,22))},ba=t().position;j.position.set(ba.x,0,ba.z),c.position.set(ba.x+2.6,1.74,ba.z+5.6);let Fr=0,Js=0;return a.setAnimationLoop(H=>{if(document.hidden||H-Js<30)return;Js=H;let $=Fr?Math.min(.08,(H-Fr)/1e3):.033;Fr=H;let st=H/1e3,vt=t();vt.tick!==Y&&(Tt(vt),Y=vt.tick);let ae=vl.smoothstep(Math.sin((vt.minute-345)/1440*Math.PI*2),-.18,.36),Pe=vt.weather==="Rain";o.copy(Mu).lerp(Pe?Tu:_p,ae),u.copy(Tl).lerp(Pe?Eu:Sa,ae),l.color.copy(o),f.intensity=.78+ae*1.39,p.intensity=.16+ae*(Pe?.8:2.35),ze.visible=ae<.27&&!Pe,T.visible=ae>.16,T.position.x=Math.sin(st*.022)*6,L.position.x=Math.sin(st*.43)*.15,Fn.set(vt.position.x,0,vt.position.z);let ke=j.position.distanceTo(Fn)>.035,Re=1-Math.exp(-$*(ke?2.65:7)),M=ke?Math.atan2(Fn.x-j.position.x,Fn.z-j.position.z):lt?Math.atan2(c.position.x-j.position.x,c.position.z-j.position.z):j.rotation.y;j.rotation.y+=Math.atan2(Math.sin(M-j.rotation.y),Math.cos(M-j.rotation.y))*(1-Math.exp(-$*5)),j.position.lerp(Fn,Re),ct.position.copy(j.position);let I=vt.mood==="Dreaming",W=vt.until>0;It.scale.y+=((I?.54:1)-It.scale.y)*(1-Math.exp(-$*5)),It.position.y=I?.04:Math.abs(Math.sin(st*(ke?9:3.1)))*(ke?.15:.035),ie.rotation.x=I?.15:0,Te.rotation.y=Math.sin(st*(vt.happiness>76?8:4))*(I?.08:.4),Jt.rotation.z=Math.sin(st*2)*.08,Jt.rotation.x=W&&(vt.action==="draw"||vt.action==="eat")?Math.sin(st*5)*.17:0;let k=I||Math.sin(st*1.42)>.996;for(let it=0;it<2;it++)Xt[it].scale.y=k?.021:.1,de[it].rotation.x=Math.sin(st*2.3+it*2)*.12;for(let it=0;it<4;it++)D[it].position.y=.17+(ke?Math.max(0,Math.sin(st*9+it*Math.PI/2))*.18:0);let G=vt.action==="play"&&W;if(at.position.y=.34+(G?Math.abs(Math.sin(st*5))*.75:0),et.position.y=at.position.y,lt)yi.set(j.position.x+2.45,1.71,j.position.z+5.15),c.position.lerp(yi,1-Math.exp(-$*3.3)),Ml.set(j.position.x,.99,j.position.z),c.lookAt(Ml),Rt=c.rotation.y,Ut=c.rotation.x;else{c.rotation.set(Ut,Rt,0,"YXZ");let it=_t=>+(bt.has(_t)||gt.has(_t));El(it("right")-it("left"),it("down")-it("up"),$*4.8)}S.visible=Pe,Pe&&(O.forEach(([it,_t,pt],wt)=>{let Ct=it+c.position.x,Gt=(_t-st*8%13+13)%13,qt=pt+c.position.z;At.set([Ct,Gt,qt,Ct+.11,Gt-.57,qt],wt*6)}),ue.attributes.position.needsUpdate=!0),b.position.copy(c.position),a.render(r,c)}),{follow:()=>U(!0),move:(H,$)=>{$?(gt.has(H)||El(+(H==="right")-+(H==="left"),+(H==="down")-+(H==="up"),.28),gt.add(H),U(!1)):gt.delete(H)},dispose:()=>{oe.disconnect(),a.setAnimationLoop(null),window.removeEventListener("blur",xt);for(let H of Object.values(v))H.dispose();ue.dispose(),R.dispose(),N.dispose(),d.dispose(),_.dispose();for(let H of h.values())H.dispose();for(let H of m)H.dispose();a.dispose(),s.removeEventListener("webglcontextlost",ge),a.forceContextLoss(),s.remove()}}}var Sl=e=>t=>(0,ft.createElement)(e,{className:t.variant},t.children),$1=Sl("main"),_i=Sl("div"),bl=Sl("div"),Ka=Sl("section"),gn=Sl("p"),tT=Sl("span"),Zs=e=>(0,ft.createElement)("h"+(e.level||2),{},e.children),eT=e=>e.children?(0,ft.createElement)("p",{className:"notice",role:"status"},e.children):null,nT=e=>(0,ft.createElement)("img",{src:e.src,alt:e.alt,width:37,height:37}),Gv=e=>(0,ft.createElement)("a",{href:e.href,className:e.variant,"data-workspace-apps-trigger":e.variant==="apps"?"true":void 0},e.children),mp=e=>(0,ft.createElement)("button",{type:"button",className:e.variant,disabled:e.disabled,onClick:e.onClick},e.children),gp=(0,ft.createContext)({name:"",setName:e=>{},save:async e=>{}}),iT=e=>{let{name:t,setName:n}=(0,ft.useContext)(gp);return(0,ft.createElement)("label",{},e.label,(0,ft.createElement)("input",{value:t,maxLength:e.limit,placeholder:e.placeholder,onChange:i=>n(i.target.value)}))};function kv({mood:e}){return(0,ft.createElement)("div",{className:"portrait","aria-hidden":!0},(0,ft.createElement)("i",{className:"ear left"}),(0,ft.createElement)("i",{className:"ear right"}),(0,ft.createElement)("div",{className:e==="Dreaming"?"face sleeping":"face"},(0,ft.createElement)("i"),(0,ft.createElement)("i"),(0,ft.createElement)("b")))}function aT({state:e}){let t=[["Happiness","\u2665",e.happiness,"happy"],["Energy","\u2726",e.energy,"energy"],["Fullness","\u25CF",e.fullness,"fullness"],["Curiosity","\u25C8",e.curiosity,"curiosity"]];return(0,ft.createElement)("div",{className:"meters"},t.map(([n,i,a,s])=>(0,ft.createElement)("div",{key:s,className:"meter "+s},(0,ft.createElement)("div",{className:"meterLabel"},(0,ft.createElement)("span",{"aria-hidden":!0},i),(0,ft.createElement)("label",{htmlFor:"need-"+s},n),(0,ft.createElement)("span",{},Math.round(a)+"%")),(0,ft.createElement)("progress",{id:"need-"+s,value:a,max:100}))))}var wN=["Sun over Clover Hill","Moonlit Pond","Wildflower Study","After the Rain"],CN=["Stone cairn","Meadow arch","Woodland tower","Clover stack"];function sT({item:e}){let t=(0,ft.useRef)(null);(0,ft.useEffect)(()=>{let i=t.current?.getContext("2d");i&&Hv(i,e,156,108)},[e.id,e.kind,e.style,e.level]);let n=(e.kind==="art"?wN:CN)[e.style];return(0,ft.createElement)("figure",{className:"thumb"},(0,ft.createElement)("canvas",{ref:t,width:156,height:108,"aria-hidden":!0}),(0,ft.createElement)("figcaption",{},n))}function rT({state:e}){let t=(0,ft.useRef)(null),n=(0,ft.useRef)(null),i=(0,ft.useRef)(e);i.current=e;let[a,s]=(0,ft.useState)(!1),[r,o]=(0,ft.useState)(!0),[l,c]=(0,ft.useState)(()=>navigator.maxTouchPoints>0||matchMedia("(pointer: coarse)").matches);return(0,ft.useEffect)(()=>{let f=matchMedia("(pointer: coarse)"),p=()=>c(navigator.maxTouchPoints>0||f.matches);return f.addEventListener("change",p),()=>f.removeEventListener("change",p)},[]),(0,ft.useEffect)(()=>{try{let f=j1(t.current,()=>i.current,o,()=>s(!0));return n.current=f,()=>{n.current=null,f.dispose()}}catch{s(!0)}},[]),(0,ft.createElement)("div",{className:"viewport"+(a?" failed":"")},(0,ft.createElement)("div",{ref:t,className:"surface"}),a?(0,ft.createElement)("div",{className:"viewFallback"},(0,ft.createElement)(kv,{mood:e.mood}),(0,ft.createElement)("p",{},"3D view needs WebGL 2. Your pet is still active."),(0,ft.createElement)("p",{},"Follow along in Field notes.")):(0,ft.createElement)("div",{className:"overlay"},(0,ft.createElement)("div",{className:"viewHeading"},(0,ft.createElement)("span",{className:"liveDot"}),r?"Watching "+e.name:"First-person view"),(0,ft.createElement)("button",{type:"button",className:"focusPet","aria-pressed":r,onClick:()=>n.current?.follow()},"\u25CE  Find "+e.name),l?(0,ft.createElement)("div",{className:"dpad",role:"group","aria-label":"Walk"},...[["up","\u2191","Move forward"],["left","\u2190","Move left"],["down","\u2193","Move back"],["right","\u2192","Move right"]].map(([f,p,u])=>(0,ft.createElement)("button",{key:f,type:"button",className:f,"aria-label":u,onPointerDown:d=>{d.preventDefault(),d.currentTarget.setPointerCapture(d.pointerId),n.current?.move(f,!0)},onLostPointerCapture:()=>n.current?.move(f,!1),onPointerUp:()=>n.current?.move(f,!1),onKeyDown:d=>{(d.key===" "||d.key==="Enter")&&(d.preventDefault(),n.current?.move(f,!0))},onKeyUp:()=>n.current?.move(f,!1),onBlur:()=>n.current?.move(f,!1)},p))):null,(0,ft.createElement)("span",{className:"viewHelp"},l?"Drag to look":"Drag to look \xB7 WASD to walk")))}function oT({state:e}){let{save:t}=(0,ft.useContext)(gp),[n,i]=(0,ft.useState)(e.mode),[a,s]=(0,ft.useState)(e.localModel),[r,o]=(0,ft.useState)(""),[l,c]=(0,ft.useState)(!1),[f,p]=(0,ft.useState)(""),[u,d]=(0,ft.useState)(!1),_=n==="local"||n==="both",b=n==="service"||n==="both",v=async h=>{h.preventDefault(),d(!0);try{await t({mode:n,model:a,key:r,clearKey:l}),o(""),c(!1),p("Decision settings saved.")}catch(m){p(m.message)}finally{d(!1)}};return(0,ft.createElement)("details",{},(0,ft.createElement)("summary",{},"Decision settings"),(0,ft.createElement)("form",{onSubmit:v},(0,ft.createElement)("p",{className:"muted"},"Models can suggest activities. Instincts handle food, sleep and outages."),(0,ft.createElement)("label",{},"Decision mode",(0,ft.createElement)("select",{value:n,onChange:h=>i(h.target.value)},(0,ft.createElement)("option",{value:"instincts"},"Instincts"),(0,ft.createElement)("option",{value:"local"},"Local Ollama model"),(0,ft.createElement)("option",{value:"service",disabled:!e.hasService},"Workspace AI Service"),(0,ft.createElement)("option",{value:"both",disabled:!e.hasService},"Local model + AI Service"))),_&&(0,ft.createElement)("label",{},"Installed Ollama model",(0,ft.createElement)("input",{value:a,placeholder:"qwen3.5:9b",maxLength:100,required:!0,onChange:h=>s(h.target.value)})),b&&(0,ft.createElement)("label",{},e.paired?"Replace the Pet app key":"Pet app key",(0,ft.createElement)("input",{value:r,type:"password",autoComplete:"new-password",maxLength:200,required:!e.paired,onChange:h=>o(h.target.value)})),b&&(0,ft.createElement)("p",{className:"muted"},"In AI Service, create a local app key with the ID pet. Paste it here."),!b&&e.paired&&(0,ft.createElement)("label",{className:"check"},(0,ft.createElement)("input",{type:"checkbox",checked:l,onChange:h=>c(h.target.checked)}),"Forget the saved AI Service key"),(0,ft.createElement)("button",{className:"secondary",type:"submit",disabled:u},"Save settings"),(f||e.modelStatus)&&(0,ft.createElement)("p",{role:"status",className:"muted"},[f,e.modelStatus].filter(Boolean).join(" "))))}var yt=Al(Xv(),1);function fT({data:e,onAction:t}){return(0,yt.jsxs)($1,{children:[(0,yt.jsxs)(_i,{variant:"topbar",children:[(0,yt.jsxs)(_i,{variant:"wordmark",children:[(0,yt.jsx)(nT,{src:"./icon.svg",alt:""}),(0,yt.jsx)(gn,{children:"Pet"}),(0,yt.jsx)(gn,{variant:"workspaceLabel",children:"One Workspace"})]}),(0,yt.jsxs)(_i,{variant:"nav",children:[(0,yt.jsx)(Gv,{variant:"apps",href:e.workspace,children:"Apps"}),e.session&&(0,yt.jsx)(mp,{variant:"quiet",onClick:()=>t("logout"),children:"Sign out"})]})]}),(0,yt.jsx)(eT,{children:e.notice}),e.loading&&(0,yt.jsx)(Ka,{variant:"welcome",children:(0,yt.jsx)(Zs,{level:1,children:"Opening the hollow\u2026"})}),e.signedOut&&(0,yt.jsxs)(Ka,{variant:"welcome",children:[(0,yt.jsx)(Zs,{level:1,children:"Clover Hollow"}),(0,yt.jsx)(gn,{children:"Visit a little woodland world and see what your pet is up to."}),e.oidc&&(0,yt.jsx)(Gv,{variant:"primary",href:"./auth/login",children:"Continue with MyAccounts"})]}),e.unadopted&&(0,yt.jsxs)(Ka,{variant:"welcome",children:[(0,yt.jsx)(Zs,{level:1,children:"Make yourself at home"}),e.owner?(0,yt.jsxs)(yt.Fragment,{children:[(0,yt.jsx)(gn,{children:"Name the pet who will live in Clover Hollow."}),(0,yt.jsx)(iT,{label:"Pet name",name:"name",limit:22,placeholder:"Pip"}),(0,yt.jsx)(mp,{variant:"primary",onClick:()=>t("adopt"),disabled:e.busy,children:"Welcome your pet"})]}):(0,yt.jsx)(gn,{children:"The Workspace owner can welcome the first pet."})]}),e.world&&(0,yt.jsxs)(yt.Fragment,{children:[(0,yt.jsxs)(_i,{variant:"titlebar",children:[(0,yt.jsxs)(bl,{children:[(0,yt.jsx)(Zs,{level:1,children:"Clover Hollow"}),(0,yt.jsx)(gn,{variant:"location",children:e.weather})]}),(0,yt.jsxs)(_i,{variant:"worldclock",children:[(0,yt.jsx)(gn,{variant:"clockIcon",children:e.sun}),(0,yt.jsxs)(bl,{children:[(0,yt.jsx)(gn,{variant:"clockValue",children:e.clock}),(0,yt.jsx)(gn,{variant:"clockDay",children:e.day})]})]})]}),(0,yt.jsxs)(_i,{variant:"layout",children:[(0,yt.jsxs)(bl,{variant:"mainColumn",children:[(0,yt.jsx)(Ka,{variant:"worldCard",children:(0,yt.jsx)(rT,{state:e.world})}),(0,yt.jsx)(_i,{variant:"viewNote",children:(0,yt.jsx)(gn,{variant:"active",children:e.status})}),(0,yt.jsxs)(Ka,{variant:"journal",children:[(0,yt.jsxs)(_i,{variant:"sectionTop",children:[(0,yt.jsx)(Zs,{level:2,children:"Field notes"}),(0,yt.jsx)(gn,{variant:"smallLabel",children:e.journalDay})]}),e.events.map((n,i)=>(0,yt.jsxs)(_i,{variant:"event",children:[(0,yt.jsx)(gn,{variant:"eventIcon",children:n.icon}),(0,yt.jsx)(gn,{children:n.text}),(0,yt.jsx)(gn,{variant:"eventTime",children:n.time})]},i))]})]}),(0,yt.jsxs)(bl,{variant:"sideColumn",children:[(0,yt.jsxs)(Ka,{variant:"petCard",children:[(0,yt.jsxs)(_i,{variant:"petIdentity",children:[(0,yt.jsx)(kv,{mood:e.world.mood}),(0,yt.jsxs)(bl,{children:[(0,yt.jsx)(Zs,{level:2,children:e.world.name}),(0,yt.jsx)(tT,{variant:"mood",children:e.world.mood})]})]}),(0,yt.jsx)(gn,{variant:"activity",children:e.activity}),(0,yt.jsx)(aT,{state:e.world}),(0,yt.jsx)(mp,{variant:"treat",onClick:()=>t("treat"),disabled:e.treatDisabled,children:e.treatLabel})]}),(0,yt.jsxs)(Ka,{variant:"gallery",children:[(0,yt.jsxs)(_i,{variant:"sectionTop",children:[(0,yt.jsx)(Zs,{level:2,children:"Made in the hollow"}),(0,yt.jsx)(gn,{variant:"smallLabel",children:e.galleryCount})]}),e.emptyGallery&&(0,yt.jsx)(gn,{variant:"muted",children:"Paintings and sculptures will appear here as your pet makes them."}),(0,yt.jsx)(_i,{variant:"galleryItems",children:e.creations.map((n,i)=>(0,yt.jsx)(sT,{item:n},i))})]}),e.owner&&(0,yt.jsx)(Ka,{variant:"brainCard",children:(0,yt.jsx)(oT,{state:e.settings})})]})]})]})]})}var hT=e=>String(Math.floor(e/60)).padStart(2,"0")+":"+String(e%60).padStart(2,"0"),NN={explore:"Exploring the meadow",sleep:"At the burrow",eat:"Foraging for a snack",play:"Playing in the meadow",draw:"Working on a painting",build:"Building a sculpture"},UN={explore:"\u2727",home:"\u2302",treat:"\u273F",sleep:"\u263E",play:"\u25C8",art:"\u25D2",build:"\u25A6",eat:"\u25CF"};async function bu(e,t){let n=await fetch(e,{cache:"no-store",...t}),i=await n.json();if(!n.ok)throw Object.assign(new Error(i.error||"Pet could not finish the request."),{status:n.status});return i}function LN(){let[e,t]=(0,Cn.useState)(),[n,i]=(0,Cn.useState)(),[a,s]=(0,Cn.useState)(),[r,o]=(0,Cn.useState)(!0),[l,c]=(0,Cn.useState)(!1),[f,p]=(0,Cn.useState)(""),[u,d]=(0,Cn.useState)(!1),[_,b]=(0,Cn.useState)(!1),[v,h]=(0,Cn.useState)("");(0,Cn.useEffect)(()=>{let N=!0,P=new AbortController;async function z(){try{let F=await bu("./api/bootstrap",{signal:P.signal});N&&t(F);let V=await bu("./api/me",{signal:P.signal});N&&i(V)}catch(F){N&&F.status!==401&&p(F.message)}finally{N&&o(!1)}}return z(),()=>{N=!1,P.abort()}},[]);let m=n?.user.id;(0,Cn.useEffect)(()=>{if(!m)return;let N=new AbortController,P,z=async()=>{try{let F=await bu("./api/world",{signal:N.signal,headers:{"X-Workspace-Expected-Account":m}});if(N.signal.aborted)return;s(V=>V?.world&&(!F.world||F.world.tick<V.world.tick||F.world.tick===V.world.tick&&F.world.treatAfter<V.world.treatAfter)?{...F,world:V.world}:F),d(!1)}catch(F){if(N.signal.aborted)return;if(F.status===401||F.status===409){i(void 0),s(void 0),c(!0),p("Your session changed. Continue with MyAccounts to return to the hollow.");return}d(!0)}P=window.setTimeout(z,document.hidden?8e3:1800)};return z(),()=>{N.abort(),clearTimeout(P)}},[m]);let g=async(N,P,z)=>{if(n)try{return await bu(N,{method:P,body:z===void 0?void 0:JSON.stringify(z),headers:{"Content-Type":"application/json","X-CSRF-Token":n.csrf,"X-Workspace-Expected-Account":n.user.id}})}catch(F){throw F.status===401&&(i(void 0),s(void 0),c(!0)),F}},x=a&&{mode:a.mode,modelStatus:a.modelStatus,localModel:a.localModel||"",paired:!!a.paired,hasService:!!a.hasService},E=async N=>{await g("./api/decisions","PUT",N);let P=await bu("./api/world",{headers:{"X-Workspace-Expected-Account":m}});s(P)},A=a?.world,w=Math.max(0,(A?.treatAfter||0)-(A?.saved||0)),y=async N=>{if(!_){b(!0),p("");try{if(N==="logout"&&(await g("./auth/logout","POST"),s(void 0),i(void 0)),N==="adopt"||N==="treat"){let P=await g("./api/"+N,"POST",N==="adopt"?{name:v.trim()||"Pip"}:void 0);s(z=>z&&(!z.world||z.world.tick<=P.tick)?{...z,world:P}:z),N==="treat"&&p("A clover treat for "+P.name+"!")}}catch(P){p(P.message)}finally{b(!1)}}},T=A?.minute??510,L=T<360||T>=1200,R=[...A?.creations??[]].sort((N,P)=>P.id-N.id);return(0,Cn.createElement)(gp.Provider,{value:{name:v,setName:h,save:E}},(0,Cn.createElement)(fT,{data:{workspace:e?.workspace||"/",oidc:e?.oidc||l,session:!!n,loading:r||!!n&&!a&&!u,signedOut:!r&&!n,unadopted:a&&!A,owner:n?.user.role==="owner",notice:f||(u?"Connection interrupted. Reconnecting to the hollow\u2026":""),busy:_,world:A,weather:A?.weather==="Rain"?"Rain in the meadow":L?"A clear night":"Clear skies",sun:L?"\u263E":A?.weather==="Rain"?"\u2602":"\u2600",clock:hT(T),day:"DAY "+(A?.day||1),status:u?"Reconnecting":"Habitat live",journalDay:"Latest activity",activity:A?A.mood==="Dreaming"?"Asleep in the burrow":NN[A.action]:"",events:(A?.journal??[]).slice(-6).reverse().map(N=>({...N,icon:UN[N.kind]||"\xB7",time:hT((N.tick*2+510)%1440)})),creations:R,galleryCount:String(R.length),emptyGallery:!R.length,treatDisabled:_||u||w>0,treatLabel:w?"Next treat in "+w+"s":"Give a clover treat",settings:x},onAction:y}))}(0,dT.createRoot)(document.querySelector("#root")).render((0,Cn.createElement)(LN));
/*! Bundled license information:

react/cjs/react.production.js:
  (**
   * @license React
   * react.production.js
   *
   * Copyright (c) Meta Platforms, Inc. and affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   *)

scheduler/cjs/scheduler.production.js:
  (**
   * @license React
   * scheduler.production.js
   *
   * Copyright (c) Meta Platforms, Inc. and affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   *)

react-dom/cjs/react-dom.production.js:
  (**
   * @license React
   * react-dom.production.js
   *
   * Copyright (c) Meta Platforms, Inc. and affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   *)

react-dom/cjs/react-dom-client.production.js:
  (**
   * @license React
   * react-dom-client.production.js
   *
   * Copyright (c) Meta Platforms, Inc. and affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   *)

react/cjs/react-jsx-runtime.production.js:
  (**
   * @license React
   * react-jsx-runtime.production.js
   *
   * Copyright (c) Meta Platforms, Inc. and affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   *)

three/build/three.core.js:
three/build/three.module.js:
  (**
   * @license
   * Copyright 2010-2026 Three.js Authors
   * SPDX-License-Identifier: MIT
   *)
*/
