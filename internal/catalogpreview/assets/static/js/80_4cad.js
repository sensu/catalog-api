(window.webpackJsonp_name_chunkhash_4_=window.webpackJsonp_name_chunkhash_4_||[]).push([[80],{"/xD5":function(e,t,n){"use strict";n("q1tI");var r=n("NZDO"),i=n("umvS"),a=n("TFnf"),s=n("GVSF"),l=n("b7jQ"),o=n("OGDC"),c=n("Gqia"),d=n("n7JV"),u=n("Ks7j"),p=n("nKUr");t.a=({open:e=!0,maxWidth:t="xs",onClose:n,title:m,titleEl:b,contents:h,actions:j,fullScreen:g="sm",ContentProps:f={},...x})=>{let v=!1,O="sm";"always"===g?v=!0:"never"!==g&&(O=g);return Object(u.a)(O,"lt")&&(v=!0),Object(p.jsxs)(i.a,{open:e,onClose:n,title:m,maxWidth:t,fullScreen:v,...x,children:[Object(p.jsx)(l.a,{children:Object(p.jsxs)(r.a,{display:"flex",flexDirection:"row",children:[Object(p.jsx)(r.a,{flexGrow:1,alignItems:"center",children:Object(p.jsx)(c.a,{variant:"h6",children:b||m})}),Object(p.jsx)(r.a,{mt:-1,mb:-1,children:n&&Object(p.jsx)(o.a,{"aria-label":"close",onClick:n,edge:"end",size:"large",children:Object(p.jsx)(d.a,{of:"close"})})})]})}),Object(p.jsx)(s.a,{...f,children:h}),j&&Object(p.jsx)(a.a,{children:j})]})}},"7BXk":function(e,t,n){"use strict";n.d(t,"b",(function(){return j}));var r=n("q1tI"),i=n.n(r),a=n("ANGJ"),s=n("/DBP"),l=n.n(s),o=n("24wR"),c=n.n(o),d=n("LutX"),u=n("G43+"),p=n("tAME"),m=n("1h/R"),b=n.n(m),h=n("nKUr");function j(e,{query:t,isHighlighted:n}){const r=l()(e,t),a=c()(e,r);return Object(h.jsx)(d.a,{selected:n,component:"div",onMouseDown:e=>e.preventDefault(),children:Object(h.jsx)("div",{children:a.map((e,t)=>e.highlight?Object(h.jsx)("span",{style:{fontWeight:800},children:e.text},String(t)):Object(h.jsx)(i.a.Fragment,{children:e.text},String(t)))})})}const g=({containerProps:{className:e,...t},children:n,inputRef:r})=>{var i;return Object(h.jsx)(p.a,{...t,anchorEl:()=>null==r?void 0:r.current,placement:"bottom-start",keepMounted:!0,style:{zIndex:1e4,width:(null==r||null===(i=r.current)||void 0===i?void 0:i.clientWidth)||null},open:!!n,children:Object(h.jsx)(u.a,{className:e,children:n})})},f=Object(a.a)(e=>({suggestionsContainer:{maxHeight:280,overflowY:"auto",overflowX:"hidden",zIndex:e.zIndex.tooltip},suggestionsContainerOpen:{},suggestion:{display:"block"},suggestionsList:{margin:0,padding:0,listStyleType:"none"},container:{width:"100%"}}));t.a=({onSuggestionsClearRequested:e=(()=>{}),renderSuggestion:t=j,renderSuggestionsContainer:n=g,inputProps:r,...a})=>{const s=f(),l=i.a.useRef(null),o=r||{};return Object(h.jsx)(b.a,{theme:s,onSuggestionsClearRequested:e,renderSuggestion:t,renderSuggestionsContainer:e=>n({...e,inputRef:l}),inputProps:{...o,ref:l},...a})}},"SA+o":function(e,t,n){"use strict";n.r(t);var r=n("q1tI"),i=n.n(r),a=n("SYuh"),s=n("aipG"),l=n("fQ5S"),o=n("1MYv"),c=n("NZDO"),d=n("n7JV"),u=n("a20y"),p=n("0p0Q"),m=n("/xD5"),b=n("auwq"),h=n("nP3w"),j=n("86yx"),g=n("7BXk"),f=n("5I82"),x=n("8cYg"),v=n("MGIy"),O=n("nKUr");function y({value:e,onChange:t,inputValue:n,onInputChange:r,clearOnBlur:i,clearOnEscape:a,clearIcon:s,clearText:l,disableClearable:o,disabledItemsFocusable:c,limitTags:d,renderTags:u,size:p,...m},b){return Object(O.jsx)(x.a,{multiple:!0,options:[],value:e,onChange:t,inputValue:n,onInputChange:r,clearOnBlur:i,clearIcon:s,clearText:l,disableClearable:o,disabledItemsFocusable:c,limitTags:d,renderTags:u,size:p,freeSolo:!0,ChipProps:{variant:"filled"===m.variant?"outlined":"filled"},renderInput:e=>{var t,n,r,i;return Object(O.jsx)(v.a,{...m,...e,InputProps:{...m.InputProps||{},...e.InputProps,startAdornment:Object(O.jsxs)(O.Fragment,{children:[null===(t=m.InputProps)||void 0===t?void 0:t.startAdornment,null===(n=e.InputProps)||void 0===n?void 0:n.startAdornment]}),endAdornment:Object(O.jsxs)(O.Fragment,{children:[null===(r=e.InputProps)||void 0===r?void 0:r.endAdornment,null===(i=m.InputProps)||void 0===i?void 0:i.endAdornment]})},inputProps:{...m.inputProps||{},...e.inputProps},ref:b})}})}var S=i.a.forwardRef(y),k=n("awfS"),C=n("Fg+5"),w=n("sbgx"),$=n("BkAX"),q=n("brhb"),F=n("Gqia"),E=n("ANGJ");function B(e){return I(e,"toUpperCase")}function I(e,t){const n=e[0]||"",r=e.slice(1);return n[t]()+r}var D=i.a.createContext({autosuggest:{subscribe:()=>({request:()=>{},cancel:()=>{}})},schema:{},groups:[],hideDesc:!1,omitDesc:!1,loading:!1});const P=()=>{},L=(e,t)=>{const{autosuggest:n}=i.a.useContext(D),[r,a]=i.a.useState([]),s=i.a.useRef(P),l=i.a.useRef(P);i.a.useEffect(()=>{if(!n)return;let r=((e,t)=>{if(!t)return t;const n=e.split("/",3);let r=n[2];if(r)return"check_config"===r&&(r="check"),`${r}${t}`;console.warn("convertFilter()","unable to parse ref",{ref:e,refParts:n})})(e,t);r&&(r="fieldSelector:"+r);const i=n.subscribe(e=>a(e));return s.current=(t,n)=>i.request(t,{...n,ref:e,filter:r}),l.current=(t,i)=>{n.selected&&n.selected(t,{...i,ref:e,filter:r})},()=>{i.cancel(),s.current=P,l.current=P}},[n,a,e,t]);return[r,Object.freeze({request:(...e)=>s.current(...e),selected:(...e)=>l.current(...e)})]};var T=n("J3Jp"),_=n.n(T),A=n("aP3S");const R=e=>{const[t]=(e=>{let t,n=e.trim();if(/^(TZ|CRON_TZ)=[^\s]+\s+/.test(n)){const e=n.indexOf(" "),r=n.indexOf("=");t=n.slice(r+1,e),n=n.slice(e).trim()}return[n,t]})(e);if(/^@(yearly|annually)$/.test(t))return"0 0 0 1 1 *";if(/^@monthly$/.test(t))return"0 0 0 1 * *";if(/^@weekly$/.test(t))return"0 0 0 * * 0";if(/^@(daily|midnight)$/.test(t))return"0 0 0 * * *";if(/^@hourly$/.test(t))return"0 0 * * * *";const n=/^@every +(.*)$/.exec(t);if(n){const e=Math.round(Object(A.d)(n[1])/A.g)*A.g;if(e>0){if(e%A.a==0)return`0 0 0 */${e/A.a} * *`;if(e%A.b==0)return`0 0 */${e/A.b} * * *`;if(e%A.c==0)return`0 */${e/A.c} * * * *`}const t=e/A.g;return t>0?`*/${t} * * * * *`:"* * * * * *"}return t};var N=n("Zrqx"),V=n("oKJj");const M="Please fill out this field.",W=(e,t,n)=>{const r=e?e.trim():"";if(n&&""===r)return M;if(void 0!==t.minLength&&r.length<t.minLength)return`Value may not be shorter than ${t.minLength.toLocaleString()}.`;var i;if(void 0!==t.maxLength&&r.length>t.maxLength)return`Value may not be longer than ${t.maxLength.toLocaleString()}.`;if(""!==r&&t.pattern){if(!new RegExp(t.pattern,"g").test(r))return(i=t.patternTitle)?`Please match the requested format: ${i}.`:"Unexpected input."}const a=Object(N.f)(t.format);return a.includes("cron")&&!(e=>{try{const t=R(e);return _.a.toString(t),!0}catch(e){return!1}})(r)?"Please provide valid cron expression.":a.includes("io.sensu.selector")&&!(e=>{try{const t=new V.Parser(e);if("BinaryExpr"===t.parse().kind)return!0}catch(e){}return!1})(r)?"Please provide valid selector expression.":void 0},z=(e,t)=>{switch(e.type){case"string":return n=>W(n,e,t);case"number":case"integer":return n=>((e,t,n)=>void 0===e?n?M:void 0:void 0===t.minimum||e>=t.minimum?void 0===t.maximum||e<=t.maximum?void 0===t.exclusiveMinimum||e>t.exclusiveMinimum?void 0===t.exclusiveMaximum||e<t.exclusiveMaximum?void 0!==t.multipleOf&&e%t.multipleOf!=0?`Value must be multiple of ${t.multipleOf.toLocaleString()}.`:void 0:`Value must be less than ${t.exclusiveMaximum.toLocaleString()}.`:`Value must be greater than ${t.exclusiveMinimum.toLocaleString()}.`:`Value must be less than or equal to ${t.maximum.toLocaleString()}.`:`Value must be greater than or equal to ${t.minimum.toLocaleString()}.`)(n,e,t);case"array":return n=>((e,t,n)=>{const r=(e||[]).length;if(n&&0===r)return M;const i=t.minItems||0,a=t.maxItems||4294967295;return r<i||r>a?((e,t)=>`Must have between ${e.toLocaleString()}...${t.toLocaleString()} elements.`)(i,a):Array.isArray(e)&&"string"===t.items.type?e.find(e=>W(e,t.items,!1)):void 0})(n,e,t)}},G=(e,t=".")=>e.filter(e=>!!e).join(t),J=(e,t,n="")=>{const{dependencies:r}=e;return!r||r.every(e=>{const r=G(n.split(".").slice(0,-1)),i=G([r,e.name]);let a=Object(j.d)(t,i);return void 0===a&&(a=Object(j.d)(t.__virtual,i)),e.condition.includes(a)})},U=e=>{if(void 0===e.oneOf)return;let t;for(const n in e.oneOf){const r=e.oneOf[n].properties||{},i=Object.keys(r).reduce((e,t)=>{const n=r[t];switch(n.type){case"boolean":case"integer":case"number":case"string":if(n.const)return[...e,t]}return e},[]);if(t=void 0===t?i:Object(N.b)(t,i),0===t.length)break}return t&&1===t.length?[t[0],e.oneOf[0].properties[t[0]]]:void 0},K=(e,t=!1)=>i.a.useMemo(()=>z(e,t),[e,t]),X=(e,t,n)=>{const r=i.a.useContext(D),a=K(e,n),{input:s,meta:l}=Object(h.b)(t,{validate:a}),o=Object(h.c)(),c=o.getFieldState.bind(null,t);return{input:s,meta:{...l,displayErr:(l.visited||l.submitFailed)&&!l.valid,disabled:void 0!==l.initial&&e.readOnly||r.loading||o.getState().submitting},getFieldState:c}},Z=Object(E.a)({group:{boxShadow:"none"},arrayButtonContainer:{display:"block",textAlign:"right"},arrayButton:{marginTop:"5px"},arrayField:{},chipHelperText:{marginBottom:"initial"},chipContainer:{marginBottom:0},chipInputRoot:{paddingTop:"20px !important"},chipInput:{marginBottom:"0px !important"},chipLabel:{top:"0px !important"}}),H=({name:e,path:t,field:n,required:r,omitDesc:a=!1,hideDesc:s=!1})=>{if("string"!==n.items.type)throw new Error("component called with incompatible type "+n.items.type);const{input:l,meta:o,getFieldState:c}=X(n,t,r),[d,u]=i.a.useState(""),[p,m]=L(n.items.ref||"",n.items.refFilter),b=l.value?l.value:[];let h=n.description;var j;(o.displayErr&&(h=o.error),n.items.examples)&&(h=Object(O.jsxs)("div",{style:{overflow:"scroll"},children:[Object(O.jsx)("span",{children:h}),Object(O.jsxs)("div",{children:["Example(s): ",null===(j=n.items.examples)||void 0===j?void 0:j.toString()]})]}));const f=Boolean(o.initial&&n.readOnly);let x=Object(O.jsx)(S,{...l,inputValue:d,onInputChange:(e,t,n)=>{("input"===n||"reset"===n&&null!==e)&&u(t)},clearOnBlur:!0,disabled:f,error:o.displayErr,fullWidth:!0,helperText:Object(O.jsx)(k.a,{in:o.displayErr||!a&&(!s||o.active),children:h}),label:n.title||B(e),InputProps:{name:t,autoComplete:"off"},onBlur:l.onBlur,onFocus:l.onFocus,onChange:(e,t)=>{const r=c();if(!r)return;let i=t.map(e=>e.trim());n.uniqueItems&&(i=Array.from(new Set(t))),i=i.filter(e=>!W(e||"",n.items,!0)),r.change(i)},required:Boolean(b.length?void 0:r),value:b,variant:"filled"});if(n.items.ref){const e=i.a.cloneElement(x);x=Object(O.jsx)(g.a,{focusInputOnSuggestionClick:!0,getSuggestionValue:e=>e,shouldRenderSuggestions:()=>!0,onSuggestionsFetchRequested:({value:e,reason:t})=>{["input-focused","input-changed"].includes(t)&&void 0!==e&&m.request(e,{})},onSuggestionSelected:(e,{suggestion:t})=>{const r=c();t&&r&&(m.selected(t,{}),((e,t,n)=>{let r=[...t.value||new Array,e];n.uniqueItems&&(r=Array.from(new Set(r))),t.change(r),t.error})(t.trim(),r,n),u(""))},inputProps:{onBlur:e=>l.onBlur(e),onChange:(e,t)=>u(t.newValue||""),onFocus:l.onFocus,value:d},renderInputComponent:({value:t,onBlur:n,onChange:r,onFocus:a,onKeyDown:s,ref:l,...o})=>i.a.cloneElement(e,{inputValue:t,onBlur:n,onFocus:a,onKeyDown:s,onInputChange:(e,t,n)=>{e&&n&&["reset","clear"].includes(n)?u(""):e&&r&&r(e)},ref:l,...o}),renderSectionTitle:e=>0===p[e].values.length?null:Object(O.jsx)($.a,{children:p[e].title}),multiSection:!0,suggestions:p.map((e,t)=>t),getSectionSuggestions:e=>p[e].values})}return Object(O.jsx)(C.a,{fullWidth:!0,margin:"dense",children:i.a.cloneElement(x,{key:t})})},Q=({name:e,path:t,field:n,required:r,component:a,omitDesc:s})=>{var l;const o=Z(),{input:u,meta:p,getFieldState:m}=X(n,t,r),b=null!==(l=u.value)&&void 0!==l?l:[],h=b.length?b.length:1,[j,g]=i.a.useState(h),x=[];for(let e=0;e<j;e+=1)x.push(Object(O.jsx)(c.a,{bgcolor:e%2==0?"divider":void 0,borderRadius:"8px",children:Object(O.jsx)(a,{name:`${t}[${e}]`,path:`${t}[${e}]`,schema:n.items,omitDesc:!0,hideDesc:!0},`${t}[${e}]`)}));const v=Array.isArray(p.error)?void 0:p.error;return Object(O.jsxs)(c.a,{marginTop:.5,marginBottom:.5,children:[!s&&Object(O.jsxs)(O.Fragment,{children:[Object(O.jsx)(F.a,{variant:"body1",children:n.title||B(e)}),Object(O.jsx)(F.a,{variant:"body2",color:"textSecondary",paragraph:!0,children:n.description})]}),Object(O.jsxs)("div",{className:o.arrayField,children:[x,Object(O.jsxs)("div",{className:o.arrayButtonContainer,children:[Object(O.jsx)(q.a,{title:"Add",children:Object(O.jsx)(f.a,{color:"primary",className:o.arrayButton,onClick:()=>g(j+1),disabled:p.disabled,children:Object(O.jsx)(d.a,{of:"new"})})}),Object(O.jsx)(q.a,{title:"Remove",children:Object(O.jsx)(f.a,{color:"primary",className:o.arrayButton,onClick:()=>{const e=m();if(!e)return;const t=j>0?j-1:0;g(t),e.change(b.slice(0,t))},disabled:p.disabled,children:Object(O.jsx)(d.a,{of:"remove"})})})]})]}),p.displayErr&&Object(O.jsx)(w.a,{error:!0,children:v})]},t)};var Y=({field:e,...t})=>"string"===e.items.type?Object(O.jsx)(H,{field:e,...t}):Object(O.jsx)(Q,{field:e,...t}),ee=n("kQF4"),te=n("T4Ez"),ne=n("LutX"),re=n("rT0G"),ie=n("mvca");const ae=Object(ie.a)(ne.a,{name:"WizardFormBooleanField",slot:"MenuItem"})({whiteSpace:"normal"});var se=({name:e,path:t,field:n,required:r,omitDesc:a=!1,hideDesc:s=!1})=>{const{input:l,meta:o,getFieldState:c}=X(n,t,r),d=i.a.useCallback(e=>{const t=c();if(!t)return;const n="1"===e.target.value;t.value!==n&&t.change(n)},[c]),u=(n.enum||[!0,!1]).map(e=>{const r=(n.enumLocale||{})[e.toString()]||{};return Object(O.jsx)(ae,{value:e?"1":"0",children:Object(O.jsx)(te.a,{primary:r.title||e.toString(),secondary:r.description})},`${t}-${e.toString()}`)});return Object(O.jsxs)(C.a,{error:o.displayErr,fullWidth:!0,margin:"dense",variant:"filled",children:[Object(O.jsx)(ee.a,{htmlFor:t,children:n.title||B(e)}),Object(O.jsx)(re.a,{...l,disabled:o.disabled,required:r,value:l.value?"1":"0",renderValue:e=>{const t="1"===e?"true":"false",r=(n.enumLocale||{})[t];return(null==r?void 0:r.title)||B(t)},onChange:d,children:u}),Object(O.jsx)(w.a,{error:o.displayErr,children:Object(O.jsx)(k.a,{in:o.displayErr||!a&&(!s||o.active),children:o.displayErr?o.error:n.description})})]},t)};var le=({path:e,field:t,required:n,component:r})=>{const a=K(t,n),{input:s}=Object(h.b)(e,{validate:a}),l=t.oneOf;if(void 0===l)throw new Error("<WizardForm.DiscriminatingUnionField /> called with field that does not appear to be discrimating union");const[o,c]=i.a.useMemo(()=>{const[e,n]=U(t)||[];if(!e||void 0===t)throw new Error("<DiscriminatingUnionField /> called with field that does not appear to be discrimating union");const r=l.reduce((t,n)=>{var r,i;const a=(n.properties||{})[e],s=a.const||"";return{...t,enum:[...t.enum||[],s],enumLocale:{...t.enumLocale||{},[s.toString()]:{title:(null===(r=a.constLocale)||void 0===r?void 0:r.title)||s.toLocaleString(),description:null===(i=a.constLocale)||void 0===i?void 0:i.description}}}},{...n,const:void 0,readOnly:t.readOnly,enum:[],enumLocale:{}});return[e,r]},[t,l]),d=(s.value||{})[o],u=i.a.useMemo(()=>{const e=l.find(e=>(e.properties||{})[o].const===d);let t={...e,properties:{}};for(const n in null==e?void 0:e.properties)n!==o&&(t={...t,properties:{...t.properties,[n]:null==e?void 0:e.properties[n]}});return t},[o,d,l]);return Object(O.jsxs)(O.Fragment,{children:[Object(O.jsx)(r,{name:o,path:`${e}.${o}`,schema:c,required:n},`${e}.${o}`),Object.keys(u.properties).map(t=>{var n;return Object(O.jsx)(r,{name:t,path:`${e}.${t}`,schema:u.properties[t],required:null===(n=u.required)||void 0===n?void 0:n.includes(t)},`${e}.${t}`)})]})},oe=n("EfzM");const ce=({name:e,path:t,field:n,required:a,omitDesc:s=!1,hideDesc:l=!1})=>{const{input:o,meta:c,getFieldState:d}=X(n,t,a);let u=n.description;c.displayErr&&(u=c.error),n.examples&&(u=Object(O.jsxs)("div",{children:[Object(O.jsx)("span",{children:u}),Object(O.jsxs)("div",{children:["Example(s): ",n.examples.toString()]})]}));const p=Object(N.f)(n.format),m=Boolean(c.initial&&n.readOnly),b=i.a.useCallback(e=>{let t;const n=d(),r=e.target.value.trim();""!==r&&(t=parseInt(r,10)),void 0===t&&""!==r||n&&n.change(t)},[d]);let h;return p.includes("duration")?h=Object(O.jsx)(oe.a,{position:"end",children:"sec."}):p.includes("bytes")&&(h=Object(O.jsx)(oe.a,{position:"end",children:"byte(s)"})),Object(O.jsx)(C.a,{fullWidth:!0,margin:"dense",children:Object(r.createElement)(v.a,{...o,id:t,key:t,label:n.title||B(e),disabled:m,variant:"filled",value:void 0!==o.value?o.value:" ",type:"number",required:a,onChange:b,error:c.displayErr,helperText:Object(O.jsx)(k.a,{in:c.displayErr||!s&&(!l||c.active),children:u}),autoComplete:"off",InputProps:{endAdornment:h},inputProps:{min:n.minimum,max:n.maximum}})},t+"-parent")};var de=i.a.memo(ce);const ue=Object(E.a)(()=>({root:{display:"inline-grid",width:"100%",boxShadow:"none"}})),pe=({name:e,path:t,field:n,required:r,component:i,omitDesc:a=!1,hideDesc:s=!1})=>{const l=ue(),o=Object(N.f)(n.required),d=K(n,r);Object(h.b)(t,{validate:d});const u=Object(h.c)(),{values:p}=u.getState();let m={};return n.properties&&(m=n.properties),n.properties&&0!==Object.keys(n.properties).length?Object(O.jsxs)(c.a,{className:l.root,children:[!a&&Object(O.jsxs)(O.Fragment,{children:[Object(O.jsx)(F.a,{variant:"body1",color:"textPrimary",children:n.title||e?B(e):B(t)}),Object(O.jsx)(F.a,{variant:"body2",color:"textSecondary",paragraph:!0,children:n.description})]}),Object.keys(m).map(e=>{const n=m[e];return J(n,p,t)?Object(O.jsx)(i,{name:e,path:`${t}.${e}`,schema:n,required:o.includes(e),hideDesc:s,omitDesc:!1},`${t}.${e}`):null})]},t):null};var me=i.a.memo(pe);const be=Object(ie.a)(ne.a,{name:"WizardFormSelectField",slot:"item"})({whiteSpace:"normal"}),he=({name:e,path:t,field:n,required:r,hideDesc:i=!1,omitDesc:a=!1})=>{const{input:s,meta:l}=X(n,t,r);if(!n.enum)throw new Error("<SelectBoxField /> with field that does not contain any enum values");const o=n.enum.map(e=>{const t=(n.enumLocale||{})[e.toString()]||{};return{title:t.title||e,description:t.description,value:e,selected:e===s.value}}),c=o.find(e=>e.selected);return Object(O.jsxs)(C.a,{error:l.displayErr,variant:"filled",fullWidth:!0,margin:"dense",children:[Object(O.jsx)(ee.a,{htmlFor:t,children:n.title||B(e)}),Object(O.jsx)(re.a,{...s,required:r,value:void 0===s.value?"":s.value,renderValue:()=>null==c?void 0:c.title,disabled:l.disabled,children:o.map(e=>Object(O.jsx)(be,{value:e.value,children:Object(O.jsx)(te.a,{primary:e.title,secondary:e.description})},`${t}-${e.title}`))}),Object(O.jsx)(w.a,{component:"div",children:Object(O.jsx)(k.a,{in:l.displayErr||!!n.description&&!a&&(!i||l.active),children:l.displayErr?l.error:n.description})})]},t)};var je=i.a.memo(he);const ge=({name:e,path:t,field:n,required:a,hideDesc:s=!1,omitDesc:l=!1})=>{const{input:o,meta:c,getFieldState:d}=X(n,t,a),[u,p]=L(n.ref||"",n.refFilter),m=Object(N.f)(n.format);let b=n.description;if(c.displayErr)b=c.error;else if(o.value&&m.includes("cron")){const e=(e=>{try{const t=R(e);return _.a.toString(t)}catch(e){return""}})(o.value);e&&(b=e)}n.examples&&(b=Object(O.jsxs)("div",{style:{overflow:"auto"},children:[Object(O.jsx)("span",{children:b}),Object(O.jsxs)("div",{children:["Example(s): ",n.examples.toString()]})]}));let h="text";const j=Object(N.b)(["tel","email","url"],m);j[0]&&(h=j[0]);let f=!1;Object(N.a)(["ecmascript-5.1","sh"],m)&&(f=!0);let x=Object(r.createElement)(v.a,{...o,id:t,key:t,variant:"filled",label:n.title||B(e),disabled:c.disabled,value:o.value||"",type:h,required:a,error:c.displayErr,fullWidth:!0,helperText:Object(O.jsx)(k.a,{in:c.displayErr||!l&&(!s||c.active),children:b}),multiline:f,autoComplete:"off",rows:f?4:void 0,inputProps:{minLength:n.minLength,maxLength:n.maxLength,pattern:n.pattern,title:n.patternTitle}});if(n.ref){const e=i.a.cloneElement(x);x=Object(O.jsx)(g.a,{focusInputOnSuggestionClick:!0,getSuggestionValue:e=>e,onSuggestionsFetchRequested:({value:e,reason:t})=>{["input-focused","input-changed"].includes(t)&&p.request(e,{})},onSuggestionSelected:(e,{suggestion:t})=>{const n=d();t&&n&&(p.selected(t,{}),n.change(t))},shouldRenderSuggestions:()=>!0,inputProps:{...o,value:o.value||""},renderInputComponent:t=>i.a.cloneElement(e,t),renderSectionTitle:e=>Object(O.jsx)($.a,{children:u[e].title}),multiSection:!0,suggestions:u.map((e,t)=>t),getSectionSuggestions:e=>u[e].values})}return Object(O.jsx)(C.a,{fullWidth:!0,margin:"dense",children:x},t+"-parent")};var fe=i.a.memo(ge),xe=n("J4nZ");const ve=Object(E.a)(e=>({root:{display:"inline-grid"},pair:{display:"flex",marginBottom:e.spacing(1)},item:{marginRight:e.spacing(.5),flex:"1 1 auto"},removeContainer:{textAlign:"right"},arrayButton:{}})),Oe=e=>e?Object.keys(e).reduce((t,n)=>(t.push({key:n,value:e[n]}),t),[]):[{key:void 0,value:void 0}],ye=e=>e.reduce((e,t)=>(e[t.key]=t.value,e),{}),Se=(e,t,n,r,i=!1)=>{i?r[n].key=e:r[n].value=e,t.change(ye(r))},ke=({path:e,name:t,field:n,required:r,omitDesc:a})=>{const s=ve(),{meta:l,getFieldState:o}=X(n,e,r),c=Oe(l.initial),[u,p]=i.a.useState(c),m=i.a.useRef(Object(xe.c)(l,"touched"));i.a.useEffect(()=>{!l.touched&&m.current.touched&&p(Oe(l.initial)),m.current=Object(xe.c)(l,"touched")},[l,p]);const b=[];for(let t=0;t<u.length;t+=1)b.push(Object(O.jsxs)("div",{className:s.pair,children:[Object(O.jsx)(v.a,{id:`${e}[${t}].key`,label:"key",value:u[t].key,onChange:e=>{const n=o();n&&Se(e.target.value,n,t,u,!0)},variant:"filled",className:s.item,disabled:l.disabled},`${e}[${t}].key`),Object(O.jsx)(v.a,{id:`${e}[${t}].value`,label:"value",value:u[t].value,onChange:e=>{const n=o();n&&Se(e.target.value,n,t,u)},multiline:!0,variant:"filled",className:s.item,disabled:l.disabled},`${e}[${t}].value`),Object(O.jsx)(q.a,{title:"Remove",children:Object(O.jsx)(f.a,{color:"primary",className:s.arrayButton,onClick:()=>{const e=o(),n=u.slice(0);n.splice(t,1),p(n),e&&e.change(ye(n))},disabled:l.disabled,children:Object(O.jsx)(d.a,{of:"remove"})})})]},`div.${t}.pair`));return Object(O.jsxs)(C.a,{fullWidth:!0,margin:"dense",children:[Object(O.jsx)(F.a,{variant:"body1",color:"textPrimary",children:!a&&(n.title||B(t))},e),Object(O.jsx)(F.a,{variant:"body2",color:"textSecondary",paragraph:!a,children:!a&&n.description}),b,Object(O.jsx)("div",{className:s.removeContainer,children:Object(O.jsx)(q.a,{title:"Add",children:Object(O.jsx)(f.a,{color:"primary",className:s.arrayButton,onClick:()=>{const e=[...u,...Oe()];p(e)},disabled:l.disabled,children:Object(O.jsx)(d.a,{of:"new"})})})})]},e+"-parent")};var Ce=i.a.memo(ke);const we=({name:e,path:t="",schema:n,required:r=!1,omitDesc:i=!1,hideDesc:a=!1})=>{let s=t;n.virtual&&(s="__virtual."+s);const l=Object(h.c)(),o=Object(j.d)(l.getState().initialValues,s);if(n.deprecated&&!o)return null;const c={name:e,path:s,required:r,component:we,omitDesc:i,hideDesc:a};switch(n.type){case"string":case"number":case"integer":case"boolean":if(n.enum)return Object(O.jsx)(je,{field:n,...c})}switch(n.type){case"boolean":return Object(O.jsx)(se,{field:n,...c});case"string":return Object(O.jsx)(fe,{field:n,...c});case"number":case"integer":return Object(O.jsx)(de,{field:n,...c});case"array":return Object(O.jsx)(Y,{field:n,...c});case"object":return n.additionalProperties?Object(O.jsx)(Ce,{field:n,...c}):U(n)?Object(O.jsx)(le,{field:n,...c}):Object(O.jsx)(me,{field:n,...c});default:return null}};var $e=i.a.memo(we),qe=n("aUsF"),Fe=n.n(qe);const Ee=(e,t,n="")=>{e&&Object.keys(e).forEach(r=>{const i=e[r];t(n,r,i),"object"===i.type&&i.properties&&Ee(i.properties,t,G([n,r]))})},Be=({onSubmit:e,resource:t={},schema:n,...r})=>{const a=i.a.useMemo(()=>{let e={};const r={};let i=Object.assign({},t);return Ee(n.properties,(n,a,s)=>{const l=G([n,a]);s.default&&(s.virtual?e=Object(j.e)(e,l,s.default):void 0===Object(j.d)(t,l)&&(i=Object(j.e)(i,l,s.default))),s.dependencies&&s.dependencies.forEach(e=>{const t=G([n,e.name]),i={path:l,field:a,ref:e};r[t]=r[t]||[],r[t].push(i)})}),Object.keys(r).forEach(t=>{for(let n=0;n<r[t].length;n++){const a=r[t][n],s=Object(j.d)(i,a.path);if(s&&0!==s){e=Object(j.e)(e,t,a.field);break}}}),{...i,__virtual:e}},[t,n]),s=i.a.useCallback((t,r)=>{let i={...t,__virtual:void 0};return Ee(n.properties,(e,n,r)=>{const a=G([e,n]);J(r,t,a)||(i=Object(j.e)(i,a,void 0))}),e(i).then(()=>{r.restart(t)})},[n.properties,e]),l=i.a.useRef();l.current||(l.current=Object(j.a)({initialValues:a,debug:void 0,onSubmit:s,validateOnBlur:!0}),l.current.pauseValidation());const o=l.current;i.a.useEffect(()=>{Fe()(a,o.getState().initialValues)||(o.setConfig("keepDirtyOnReinitialize",!0),o.initialize(a),o.setConfig("keepDirtyOnReinitialize",!1))},[a]);const c=i.a.useContext(D);return{form:o,config:{...c,resource:t,schema:n,...r}}},Ie=({form:e,config:t,children:n})=>Object(O.jsx)(O.Fragment,{children:Object(O.jsx)(h.a,{form:e,children:()=>Object(O.jsx)(D.Provider,{value:t,children:n})})});var De=n("Ty5D");const Pe=[{type:"section",title:"Complete"},{type:"markdown",body:"Your response has been record."}],Le=[{type:"section",title:"Failed"},{type:"markdown",body:"An unrecoverable error occurred; consult the logs for more."}],Te=[{type:"section",title:"Summary"},{type:"custom",Component:e=>Object(O.jsx)(u.a,{md:`\nPlease verify that the details collected are corrrect.\n\n\`\`\`json\n${JSON.stringify(e.variables,null,"\t")}\n\`\`\`\n      `})}],_e="Are you sure you'd like to close the wizard? Installation of this integration will be canceled and any changes will be lost.",Ae=e=>(e.length>0&&"section"!==e[0].type&&(e=[{type:"section"},...e]),(e=>e.reduce((e,t)=>"section"!==t.type?e:[...e,t],[]))(e).reduce((t,n)=>{const r=((e,t)=>{let n=e.indexOf(t)+1;-1===n&&(n=e.length-1);const r=e.slice(n).findIndex(e=>"section"===e.type);return e.slice(n,-1!==r?r+n:void 0)})(e,n);return[...t,{title:n.title,blocks:r}]},[]));function Re(e){return Object(xe.c)(e,"errors","dirty","submitting","submitSucceeded","submitFailed")}var Ne=({submitLabel:e="Submit",nextLabel:t="Next",prevLabel:n="Back",endLabel:r="Finish",title:a="Questionnaire",PromptComponent:s=De.a,open:l=!0,...o})=>{const c=(({blocks:e,completionBlock:t=Pe,failureBlock:n=Le,summaryBlock:r=Te,exitTitle:a=_e,disableConfirmExit:s,suggester:l,onSubmit:o,onClose:c})=>{const[d,u]=i.a.useState({chapter:"collect",page:0,dirty:!1,submitting:!1,submitSuccess:!1,submitFailed:!1});r&&(e=[...e,...r]);const p=i.a.useMemo(()=>({collect:Ae(e),complete:Ae(t),failed:Ae(n)}),[e,t,n]),m=i.a.useMemo(()=>e.reduce((e,t)=>{if("question"!==t.type)return e;const n={...e,properties:{...e.properties||{},[t.name]:t.input}};return t.required&&(n.required=[...n.required||[],t.name]),n},{}),[e]),b=i.a.useCallback(e=>o(e).then(()=>u(e=>({...e,chapter:"complete",page:0}))).catch(()=>u(e=>({...e,chapter:"failed",page:0}))),[o,u]),h=Be({autosuggest:l,schema:m,hideDesc:!0,onSubmit:b}),j=d.page>0,g=i.a.useCallback(()=>{j&&u(e=>({...e,page:e.page-1}))},[j,u]),f=!d.submitting&&0===Object.keys(d.errors||{}).length,x=i.a.useCallback(()=>{f&&(d.page===p[d.chapter].length-1?"collect"===d.chapter?h.form.submit():c():u(e=>({...e,page:e.page+1})))},[f,d,u,h.form,p,c]);i.a.useEffect(()=>h.form.subscribe(e=>{const t=Re(d),n=Re(e);Fe()(t,n)||u(e=>({...e,...n}))},{errors:!0,dirty:!0,submitting:!0,submitSucceeded:!0,submitFailed:!0}),[h.form,d,u]);const v=Object(xe.b)(h.form.getState().values,"__virtual"),O="collect"!==d.chapter||!d.dirty&&0===Object.keys(v).length,y=i.a.useCallback(()=>{(O||s||confirm(a))&&c()},[O,a,s,c]);return{...d,values:v,chapters:p,formState:h,canGoBack:j,handleBack:g,handleClose:y,canContinue:f,canSafelyClose:O,handleContinue:x}})(o),d=c.chapters[c.chapter],h=d[c.page];return Object(O.jsxs)(Ie,{...c.formState,children:[Object(O.jsx)(m.a,{title:h.title||a,open:l,onClose:c.handleClose,maxWidth:"sm",PaperProps:{sx:{minHeight:480,maxHeight:"min(720px, 100vw)",minWidth:e=>`min(${e.breakpoints.values.sm}px, 100vw)`}},ContentProps:{dividers:!0,sx:{px:3,py:2.5}},contents:h.blocks.reduce((e,t,n)=>{let r=null;switch(t.type){case"markdown":r=Object(O.jsx)(u.a,{md:t.body},`${c.page}-markdown-${n}`);break;case"question":r=Object(O.jsx)($e,{name:t.name,path:t.name,schema:t.input,required:t.required,hideDesc:!0},`${c.page}-question-${n}`);break;case"custom":r=Object(O.jsx)(t.Component,{variables:c.values},`${c.page}-custom-${n}`)}return[...e,r]},[]),actions:Object(O.jsx)(b.a,{steps:d.length,activeStep:c.page,position:"static",variant:"dots",sx:{flexGrow:1},nextButton:Object(O.jsx)(p.a,{disabled:!c.canContinue,loading:c.submitting,onClick:c.handleContinue,children:d.length-1===c.page?"collect"===c.chapter?e:r:t}),backButton:Object(O.jsx)(p.a,{disabled:!c.canGoBack||c.submitting,onClick:c.handleBack,children:n})})}),Object(O.jsx)(s,{when:!o.disableConfirmExit&&!c.canSafelyClose,message:o.exitTitle||_e})]})},Ve=n("uGfC");t.default=({collectNamespace:e="if-required",patcher:t,postInstallBlock:n,prompts:r=[],submitLabel:p="Apply",summaryBlock:m,title:b="Wizard",...h})=>{let j=Object.assign([],r||{});const g=Object(a.a)();let f;return("always"===e||!g.namespace&&"never"!==e)&&(j[0]&&"section"!==j[0].type&&(j=[{type:"section",title:"Installer"},...j]),j=[{type:"section",title:"Select Namespace"},{type:"markdown",body:"Select the namespace and cluster in-which you would like to _install_ the integration."},{type:"question",name:"__hidden.namespace",input:{type:"string",title:"Namespace",default:g.namespace,pattern:"^[\\w.\\-:]+$",patternTitle:"may only contain letters, numbers, periods, colons, and dashes"},required:!0},{type:"question",name:"__hidden.cluster",input:{type:"string",title:"Cluster",default:g.cluster}},...j]),n&&(f=[{type:"section",title:"Next Steps"},...n]),m=i.a.useMemo(()=>m||(t?[{type:"section",title:"Summary"},{type:"custom",Component:e=>Object(O.jsxs)(O.Fragment,{children:[Object(O.jsx)(u.a,{md:"The following resources will be applied..."}),Object(O.jsx)(c.a,{children:t(e.variables).map((e,t)=>Object(O.jsxs)(s.a,{children:[Object(O.jsxs)(o.a,{expandIcon:Object(O.jsx)(d.a,{of:"expand"}),"aria-controls":`panel-${t}-content`,id:`panel-${t}-header`,children:[Object(O.jsx)(F.a,{sx:{width:"60%",flexShrink:0},children:e.metadata.name}),Object(O.jsxs)(F.a,{sx:{color:"text.secondary"},children:[e.api_version," ",e.type]})]}),Object(O.jsx)(l.a,{children:Object(O.jsx)(Ve.a,{data:e})})]},e.metadata.name))})]})}]:void 0),[m,t]),Object(O.jsx)(Ne,{...h,blocks:j,completionBlock:f,submitLabel:p,summaryBlock:m,title:b})}}}]);
//# sourceMappingURL=80_4cad.js.map