(window.webpackJsonp_name_chunkhash_4_=window.webpackJsonp_name_chunkhash_4_||[]).push([[34],{"+bLY":function(e,t,n){"use strict";var i=n("q1tI"),a=n.n(i),r=n("a/bO"),o=n("wEEd"),c=n("ANGJ"),s=n("Yd4P"),l=n.n(s),d=n("DyER"),u=n("Vzdd"),m=n("nKUr");const b={};for(let e=0;e<25;e+=1)b[`&:nth-of-type(${e})`]={zIndex:25-e};const j=Object(c.a)(e=>({container:{position:"relative",marginTop:e.spacing(-1),top:e.spacing(),paddingBottom:e.spacing(),overflow:"hidden"},banner:{position:"relative",...b},bannerInner:{position:"absolute",left:0,right:0,bottom:0}})),h=()=>{const e=Object(r.d)().filter(e=>e.props.type===u.a).slice(-20).reverse();return Object(m.jsx)(v,{elements:e})},v=({elements:e})=>{const t=j(),[n,i]=a.a.useState({}),r=a.a.useCallback((e,t)=>{i(n=>t.height===n[e]?n:{...n,[e]:t.height})},[i]),c=a.a.useCallback(e=>{i(t=>(void 0===t[e]||delete t[e],t))},[i]),s=Object(o.c)(e,e=>e.id,{from:{opacity:0,height:0},leave:{opacity:0,height:0},update:({id:e})=>({opacity:1,height:n[e]||0}),config:{tension:210,friction:20}});return Object(m.jsx)("div",{className:t.container,children:s.map(({item:{id:e,props:n,remove:i},props:a,key:s})=>Object(m.jsx)(o.a.div,{style:a,className:t.banner,children:Object(m.jsxs)("div",{className:t.bannerInner,children:[Object(m.jsx)(l.a,{onResize:t=>r(e,t)}),Object(m.jsx)(d.a,{onUnmount:()=>c(e)}),n.render({id:e,remove:i})]})},s))})};t.a=a.a.memo(h)},"0YP5":function(e,t,n){"use strict";var i=n("q1tI"),a=n.n(i);const r=Object.freeze({getItem:()=>null,setItem:()=>null});function o(){return(new Date).getTime()}t.a=function(e,{delay:t=864e5,localStorage:n=(()=>{var e;return null===(e=window)||void 0===e?void 0:e.localStorage})()||r,sessionStorage:i=(()=>{var e;return null===(e=window)||void 0===e?void 0:e.sessionStorage})()||r,reminderStorage:c}={}){const s=a.a.useMemo(()=>{const t=[e,"veteran"].join(".");return{key:t,val:"true"===n.getItem(t),set:e=>n.setItem(t,e.toString())}},[e,n]),l=a.a.useMemo(()=>{const t=[e,"remindme"].join("."),a="session"===c?i:n;return{key:t,val:parseInt(a.getItem(t)||"0",10),set:e=>a.setItem(t,(o()+e).toString())}},[e,c,n,i]),[d,u]=a.a.useState(()=>!s.val&&l.val<o());return a.a.useMemo(()=>[d,{close:()=>{s.set(!0),u(!1)},remind:e=>{const n=e||t;l.set(n),u(!1)}}],[d,u,s,l,t])}},"4OOL":function(e,t,n){"use strict";n.d(t,"b",(function(){return i}));const i={normal:{active:1e4,inactive:3e5},infrequent:{active:3e4,inactive:6e5},long:{active:12e4,inactive:18e5}},a={licenseExpiryReminder:null,alwaysShowLocalCluster:!1,linkPolicy:{allowList:!1,URLs:[]},catalog:{disabled:!1,URL:"https://catalog.sensu.io",releaseVersion:"version"},preferences:{pageSize:25,theme:"sensu",pollInterval:i.normal.active,serializationFormat:null},pagePreferences:[]};t.a=a},"7TpP":function(e,t,n){"use strict";var i=n("q1tI"),a=n.n(i),r=n("4OOL");let o;!function(e){e.json="json",e.yaml="yaml"}(o||(o={})),t.a=a.a.createContext(r.a)},ANGJ:function(e,t,n){"use strict";n("CSxS");var i=n("fySL"),a=n("UnQg");const r=Object(i.createMakeStyles)({useTheme:a.a});t.a=(...e)=>{const t=r.makeStyles()(...e);return(...e)=>t(...e).classes}},DyER:function(e,t,n){"use strict";var i,a,r,o=n("q1tI"),c=n.n(o),s=n("17x9"),l=n.n(s);class d extends c.a.PureComponent{componentWillUnmount(){this.props.onUnmount()}render(){return null}}i=d,a="propTypes",r={onUnmount:l.a.func.isRequired},a in i?Object.defineProperty(i,a,{value:r,enumerable:!0,configurable:!0,writable:!0}):i[a]=r,t.a=d},H0Zs:function(e,t,n){"use strict";let i;n.d(t,"a",(function(){return i})),function(e){e.redirect="redirect-to",e.filters="filters",e.q="q",e.order="order",e.tab="tab",e.version="version",e.limit="limit",e.offset="offset"}(i||(i={}))},IeoD:function(e,t,n){"use strict";n.r(t);var i=n("q1tI"),a=n.n(i),r=n("Ty5D"),o=n("wAjv"),c=n("imBb"),s=n.n(c),l=n("nKUr");const d=e=>({trap:new s.a(e),bindings:[]}),u=a.a.createContext(d());var m=e=>{const t=Object(i.useContext)(u);Object(i.useEffect)(()=>(((e,t)=>{e.bindings.push(t),e.trap.bind(t.keys,()=>{const n=e.bindings.reverse().find(e=>e.keys===t.keys);n&&n.callback()})})(t,e),()=>((e,t)=>{e.bindings=e.bindings.reduce((e,n)=>n.id===t.id?e:[...e,n],[])})(t,e)))};const b={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"mutation",name:{kind:"Name",value:"ToggleSwitcher"},variableDefinitions:[],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"toggleModal"},arguments:[{kind:"Argument",name:{kind:"Name",value:"modal"},value:{kind:"EnumValue",value:"PREFERENCES_MODAL"}}],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:83,source:{body:"\n  mutation ToggleSwitcher {\n    toggleModal(modal: PREFERENCES_MODAL) @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}};var j=()=>{const e=Object(o.a)();return m({id:"preferences-toggle",keys:"ctrl+,",callback:()=>e.mutate({mutation:b})}),Object(l.jsx)(a.a.Fragment,{})};const h={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"mutation",name:{kind:"Name",value:"ToggleSwitcher"},variableDefinitions:[],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"toggleModal"},arguments:[{kind:"Argument",name:{kind:"Name",value:"modal"},value:{kind:"EnumValue",value:"SYS_STATUS_MODAL"}}],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:82,source:{body:"\n  mutation ToggleSwitcher {\n    toggleModal(modal: SYS_STATUS_MODAL) @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}},v=()=>{const e=Object(o.a)();return m({id:"system-status-toggle",keys:"ctrl+.",callback:()=>e.mutate({mutation:h})}),Object(l.jsx)(a.a.Fragment,{})};var k=a.a.memo(v),p=n("Vz6H"),f=n("pCzp");const O={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"query",name:{kind:"Name",value:"SwitcherIsOpenQuery"},variableDefinitions:[],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"modalStack"},arguments:[],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:58,source:{body:"\n  query SwitcherIsOpenQuery {\n    modalStack @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}},g={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"mutation",name:{kind:"Name",value:"CloseSwitcher"},variableDefinitions:[{kind:"VariableDefinition",variable:{kind:"Variable",name:{kind:"Name",value:"modal"}},type:{kind:"NamedType",name:{kind:"Name",value:"Modal"}},directives:[]}],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"clearModal"},arguments:[{kind:"Argument",name:{kind:"Name",value:"modal"},value:{kind:"Variable",name:{kind:"Name",value:"modal"}}}],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:85,source:{body:"\n  mutation CloseSwitcher($modal: Modal) {\n    clearModal(modal: $modal) @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}};var x=e=>{var t;const n=null===(t=Object(f.a)({query:O}).data)||void 0===t?void 0:t.modalStack,i=(n[n.length-1]||"")===e,r=Object(o.a)();return[i,a.a.useCallback(()=>r.mutate({mutation:g,variables:{modal:e}}),[r,e])]},y=n("Ks7j"),S=n("mvca"),N=n("NZDO"),C=n("umvS"),D=n("TFnf"),w=n("GVSF"),T=n("akmu"),M=n("b7jQ"),E=n("OGDC"),I=n("Gqia"),q=n("n7JV");const P=Object(S.a)(T.a)(({theme:e})=>({marginBottom:e.spacing(2)}));var A=({open:e,onClose:t,loading:n=!1,title:i,desc:a,message:r,fields:o,actions:c,...s})=>{const d=Object(y.a)("sm","lt");let u=o;return"function"!=typeof u||n||(u=u()),Object(l.jsx)(C.a,{open:e,fullScreen:d,...s,children:!n&&Object(l.jsxs)(l.Fragment,{children:[Object(l.jsx)(M.a,{sx:{px:d?.5:void 0},children:Object(l.jsxs)(N.a,{display:"flex",flexDirection:"row",children:[Object(l.jsx)(N.a,{flexGrow:1,alignItems:"center",children:Object(l.jsx)(I.a,{variant:"h6",children:i})}),Object(l.jsx)(N.a,{mt:-1,mb:-1,children:t&&Object(l.jsx)(E.a,{"aria-label":"close",onClick:t,edge:"end",size:"large",children:Object(l.jsx)(q.a,{of:"close"})})})]})}),Object(l.jsxs)(w.a,{sx:{px:d?.5:void 0},children:[r&&r(),a&&Object(l.jsx)(P,{children:a}),Object(l.jsx)("div",{children:u})]}),c&&Object(l.jsx)(D.a,{children:c})]})})},L=n("hGmu"),F=n("H9le"),Q=n("ZvkB"),R=n("a6xD"),G=n("T4Ez"),z=n("DFFc"),V=n("BkAX"),U=n("mkGA"),_=n("YM+J"),B=n("4enW"),W=n("zK79"),$=n("LutX"),K=n("brhb");const J={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"query",name:{kind:"Name",value:"GetIdentity"},variableDefinitions:[],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"auth"},arguments:[],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"accessToken"},arguments:[],directives:[]}]}}]}}],loc:{start:0,end:70,source:{body:"\n  query GetIdentity {\n    auth @client {\n      accessToken\n    }\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}};var Y=()=>{const e=(Object(f.a)({query:J}).data.auth.accessToken||"").split(".")[1]||"";return JSON.parse(window.atob(e)||"{}")},H=n("uDTK"),Z=n("WDlM");const X={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"mutation",name:{kind:"Name",value:"PresentModal"},variableDefinitions:[{kind:"VariableDefinition",variable:{kind:"Variable",name:{kind:"Name",value:"modal"}},type:{kind:"NamedType",name:{kind:"Name",value:"Modal"}},directives:[]}],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"presentModal"},arguments:[{kind:"Argument",name:{kind:"Name",value:"modal"},value:{kind:"Variable",name:{kind:"Name",value:"modal"}}}],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:86,source:{body:"\n  mutation PresentModal($modal: Modal) {\n    presentModal(modal: $modal) @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}},ee={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"mutation",name:{kind:"Name",value:"ToggleModal"},variableDefinitions:[{kind:"VariableDefinition",variable:{kind:"Variable",name:{kind:"Name",value:"modal"}},type:{kind:"NamedType",name:{kind:"Name",value:"Modal"}},directives:[]}],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"toggleModal"},arguments:[{kind:"Argument",name:{kind:"Name",value:"modal"},value:{kind:"Variable",name:{kind:"Name",value:"modal"}}}],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:84,source:{body:"\n  mutation ToggleModal($modal: Modal) {\n    toggleModal(modal: $modal) @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}};var te=()=>{const e=Object(o.a)();return a.a.useCallback((t,{mode:n}={})=>e.mutate({mutation:"present"===n?X:ee,variables:{modal:t}}),[e])},ne=n("V4jd"),ie=n("KQEc");const ae={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"mutation",name:{kind:"Name",value:"SetThemeMudation"},variableDefinitions:[{kind:"VariableDefinition",variable:{kind:"Variable",name:{kind:"Name",value:"theme"}},type:{kind:"NonNullType",type:{kind:"NamedType",name:{kind:"Name",value:"String"}}},directives:[]}],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"setTheme"},arguments:[{kind:"Argument",name:{kind:"Name",value:"theme"},value:{kind:"Variable",name:{kind:"Name",value:"theme"}}}],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:88,source:{body:"\n  mutation SetThemeMudation($theme: String!) {\n    setTheme(theme: $theme) @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}},re={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"mutation",name:{kind:"Name",value:"ToggleDarkMutation"},variableDefinitions:[{kind:"VariableDefinition",variable:{kind:"Variable",name:{kind:"Name",value:"value"}},type:{kind:"NonNullType",type:{kind:"NamedType",name:{kind:"Name",value:"ThemeDarkModeValue"}}},directives:[]}],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"setDarkMode"},arguments:[{kind:"Argument",name:{kind:"Name",value:"value"},value:{kind:"Variable",name:{kind:"Name",value:"value"}}}],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:105,source:{body:"\n  mutation ToggleDarkMutation($value: ThemeDarkModeValue!) {\n    setDarkMode(value: $value) @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}},oe=(e,{sep:t=" · ",max:n=5}={})=>{const i=e.slice(0).sort().slice(0,n);return e.length>i.length&&i.push(`and ${e.length-i.length} others`),i.join(t)};var ce=()=>{var e;const t=Object(o.a)(),n=Object(H.a)(),i=Y(),r=te(),c=":"+Object(Z.a)(),s=":"+Object(Z.a)(),d=a.a.useCallback(()=>t.mutate({mutation:re,variables:{value:n.dark?"LIGHT":"DARK"}}),[t,n.dark]),u=a.a.useCallback(()=>t.mutate({mutation:re,variables:{value:n.usingSystemColourScheme?"LIGHT":"UNSET"}}),[t,n.usingSystemColourScheme]),[m,b]=a.a.useState(null),j=a.a.useCallback(e=>b(e.currentTarget),[b]),h=e=>()=>{t.mutate({mutation:ae,variables:{theme:e}}),b(null)};return Object(l.jsxs)(l.Fragment,{children:[Object(l.jsx)(F.a,{children:Object(l.jsxs)(Q.a,{children:[Object(l.jsx)(R.a,{children:Object(l.jsx)(ie.a,{username:i.sub})}),Object(l.jsx)(G.a,{primary:i.sub,secondary:oe(i.groups)}),Object(l.jsxs)(z.a,{children:["basic"===(null===(e=i.provider)||void 0===e?void 0:e.provider_id)&&Object(l.jsx)(K.a,{title:"Change Password",children:Object(l.jsx)(E.a,{onClick:()=>r("PASSWORD_UPDATE_MODAL",{mode:"present"}),edge:"end",size:"large",children:Object(l.jsx)(q.a,{of:"lock"})})}),Object(l.jsx)(K.a,{title:"Sign-out",children:Object(l.jsx)(E.a,{onClick:()=>Object(ne.a)(t),edge:"end",size:"large",children:Object(l.jsx)(q.a,{of:"signout"})})})]})]})}),Object(l.jsxs)(F.a,{subheader:Object(l.jsx)(V.a,{children:"Brightness"}),children:[Object(l.jsxs)(Q.a,{children:[Object(l.jsx)(U.a,{children:Object(l.jsx)(q.a,{of:"moon"})}),Object(l.jsx)(G.a,{id:c,primary:"Dark mode"}),Object(l.jsx)(z.a,{children:Object(l.jsx)(_.a,{edge:"end",onChange:d,checked:n.dark,inputProps:{"aria-labelledby":c}})})]}),Object(l.jsxs)(Q.a,{children:[Object(l.jsx)(U.a,{children:Object(l.jsx)(q.a,{of:"bulb"})}),Object(l.jsx)(G.a,{id:s,primary:"Use system settings",secondary:"Set dark mode to use the light or dark selection located in your system settings."}),Object(l.jsx)(z.a,{children:Object(l.jsx)(_.a,{edge:"end",onChange:u,checked:n.usingSystemColourScheme,inputProps:{"aria-labelledby":s}})})]})]}),Object(l.jsx)(F.a,{subheader:Object(l.jsx)(V.a,{children:"Appearance"}),children:Object(l.jsxs)(Q.a,{component:"li",button:!0,onClick:j,children:[Object(l.jsx)(U.a,{children:Object(l.jsx)(q.a,{of:"eye"})}),Object(l.jsx)(G.a,{primary:"Theme",secondary:n.value})]})}),Object(l.jsx)(B.a,{id:"theme-menu",anchorEl:m,open:Boolean(m),onClose:()=>b(null),children:Object(l.jsxs)(W.a,{children:[Object(l.jsx)($.a,{selected:"sensu"===n.value,onClick:h("sensu"),children:Object(l.jsx)(G.a,{primary:"Sensu Go"})}),Object(l.jsx)($.a,{selected:"uchiwa"===n.value,onClick:h("uchiwa"),children:Object(l.jsx)(G.a,{primary:"Uchiwa"})}),Object(l.jsx)($.a,{selected:"highcontrast"===n.value,onClick:h("highcontrast"),children:Object(l.jsx)(G.a,{primary:"High Contrast"})}),Object(l.jsx)($.a,{selected:"classic"===n.value,onClick:h("classic"),children:Object(l.jsx)(G.a,{primary:"Classic"})}),Object(l.jsx)($.a,{selected:"deuteranopia"===n.value,onClick:h("deuteranopia"),children:Object(l.jsx)(G.a,{primary:"Deuteranopia"})}),Object(l.jsx)($.a,{selected:"tritanopia"===n.value,onClick:h("tritanopia"),children:Object(l.jsx)(G.a,{primary:"Tritanopia"})})]})})]})};var se=()=>{const[e,t]=x("PREFERENCES_MODAL");return Object(l.jsx)(A,{open:e,onClose:t,TransitionComponent:L.a,maxWidth:"xs",title:"Preferences",fields:Object(l.jsx)(ce,{})})},le=n("kUGQ"),de=n("0YP5"),ue=n("InSj"),me=n("ZwXh"),be=n("5I82");const je=({onClickSignin:e,onClickDismiss:t,...n})=>Object(l.jsx)(me.a,{message:"You are not authenticated and as a result some functionality may be unavailable.",variant:"info",actions:Object(l.jsxs)(l.Fragment,{children:[Object(l.jsx)(be.a,{variant:"outlined",color:"inherit",onClick:e,children:"Sign-in"}),Object(l.jsx)(be.a,{color:"inherit",onClick:t,sx:{ml:1},children:"Dismiss"})]}),...n});var he=a.a.memo(je),ve=n("aP3S"),ke=n("H0Zs");const pe={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"query",name:{kind:"Name",value:"SigninAlertQuery"},variableDefinitions:[],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"auth"},arguments:[],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"accessToken"},arguments:[],directives:[]},{kind:"Field",name:{kind:"Name",value:"invalid"},arguments:[],directives:[]}]}}]}}],loc:{start:0,end:89,source:{body:"\n  query SigninAlertQuery {\n    auth @client {\n      accessToken\n      invalid\n    }\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}};const fe=()=>{var e;const t=Object(f.a)({query:pe}),[n,{remind:i}]=Object(de.a)("io.alert.signin"),o=a.a.useCallback(()=>i(5*ve.c),[i]),{replace:c,location:s}=Object(r.h)(),d=a.a.useCallback(()=>{c(function(e){let t=e.search||"?";if("/signin"!==e.pathname){const n=new URLSearchParams(t);n.set(ke.a.redirect,e.pathname+e.search),t="?"+n.toString()}return"/signin"+t}(s)),i(1*ve.c)},[i,c,s]),u=t.networkStatus,m=null===(e=t.data)||void 0===e?void 0:e.auth,b=!!m.accessToken&&!m.invalid;return!n||u<6||b?null:Object(l.jsx)(ue.a,{children:Object(l.jsx)(he,{onClickDismiss:o,onClickSignin:d})})};var Oe=a.a.memo(fe),ge=n("nVmw"),xe=n("wEEd"),ye=n("vzfW");const Se=224,Ne=64,Ce=48,De=56,we=48;var Te=n("CjiS");var Me=({icon:e})=>Object(l.jsx)(N.a,{display:"flex",alignItems:"center",justifyContent:"center",flexShrink:0,width:48,height:48,children:e||Object(l.jsx)(N.a,{})});const Ee=({accountId:e,toolbarItems:t=[],isOpen:n})=>{const i=Object(xe.b)({height:n?De:0,opacity:n?1:0});return Object(l.jsxs)(N.a,{display:"flex",alignItems:"center",height:n?De:"auto",flexDirection:n?"row":"column",paddingTop:.5,paddingBottom:.5,children:[Object(l.jsx)(K.a,{title:Object(l.jsxs)("span",{children:["Authenticated as ",Object(l.jsx)("em",{children:e})]}),open:!n&&void 0,placement:"right",children:Object(l.jsx)(N.a,{display:e?"flex":"none",flexDirection:"row",alignItems:"center",width:n?void 0:we,children:Object(l.jsx)(Me,{icon:Object(l.jsx)(ie.a,{username:e||""})})})}),Object(l.jsx)(I.a,{component:xe.a.p,variant:"body1",color:"inherit",noWrap:!0,style:i,sx:{display:"flex",flexGrow:1,alignItems:"center"},children:e}),t.map(({id:e,hint:t,icon:n,onClick:i})=>Object(l.jsx)(K.a,{title:t||"",open:!!t&&void 0,placement:"right",children:Object(l.jsx)(E.a,{color:"inherit",onClick:i,size:"large",children:n})},e))]})};var Ie=Object(i.memo)(Ee),qe=n("UnQg"),Pe=n("2Qr1");const Ae=Object(i.memo)(({color:e})=>{const t=Object(qe.a)(),n=e||t.palette.text.primary,i=Object(Pe.a)(n,.125);return Object(l.jsx)(N.a,{component:"hr",border:"0",margin:"0",marginTop:"-1px",height:"1px",style:{background:`linear-gradient(\n          to right,\n          rgba(0, 0, 0, 0),\n          ${i} 10%,\n          ${Object(Pe.a)(n,.5)},\n          ${i} 90%,\n          rgba(0, 0, 0, 0)\n        )`}})});Ae.displayName="HorizontalRule";var Le=Ae,Fe=n("XikB"),Qe=n("TSYQ"),Re=n.n(Qe),Ge=n("55Ip"),ze=n("ANGJ");const Ve=Object(ze.a)(e=>({root:{paddingTop:0,paddingBottom:0,color:"inherit"},active:{backgroundColor:e.palette.action.hover,fontWeight:600}})),Ue=({active:e,collapsed:t,children:n,disabled:i,to:a,onClick:r})=>{const o=Object(ye.a)(),c=Ve();return Object(l.jsx)(Q.a,{className:Re()(c.root,{[c.active]:e}),button:!0,component:a?Ge.b:"button",disabled:i,disableGutters:!0,to:a,onClick:r,sx:{color:e&&t?o.palette.primary.contrastText:o.palette.text.primary,borderRadius:.25},children:n})},_e=({adornment:e,collapsed:t=!1,contents:n,disabled:i=!1,hint:a,href:o,icon:c,onClick:s})=>{const d=(e=>{const t=Object(r.i)();if(!e)return!1;let n="";try{n=new URL(e,"https://contoso.org").pathname}catch(t){n=e}return!!Object(r.g)(t.pathname,{path:n,exact:!1,strict:!1})})(o||""),u=Object(l.jsx)(N.a,{component:"li",display:"flex",justifyContent:"left",height:Ce,children:Object(l.jsxs)(Ue,{active:d,collapsed:t,disabled:i||!o&&!s,to:o,onClick:s,children:[Object(l.jsx)(Me,{icon:c}),Object(l.jsx)(I.a,{variant:"body1",color:"inherit",noWrap:!0,sx:{alignItems:"center",marginLeft:1,flexGrow:1,minWidth:0},children:Object(l.jsx)(N.a,{component:"span",fontWeight:d?600:"inherit",children:n})}),e&&Object(l.jsx)(N.a,{display:"flex",alignItems:"center",justifyContent:"center",width:48,height:48,children:e})]})});return t?Object(l.jsx)(K.a,{title:a||n,placement:"right",children:u}):u},Be=({icon:e,collapsed:t,contents:n,links:i,expanded:o,onExpand:c,onClick:s,...d})=>{const{pathname:u}=Object(r.i)(),m=a.a.useMemo(()=>i.find(e=>Object(r.g)(u,{path:e.href})),[i,u]),b=!t&&void 0!==m||o,j=Object(xe.b)({height:b?i.length*Ce:0,opacity:b?1:0}),h=Object(xe.b)({transform:b?"rotateX(180deg) translateY(5%)":"rotateX(0deg) translateY(0px)"});return Object(l.jsxs)(a.a.Fragment,{children:[Object(l.jsx)(_e,{...d,icon:e,adornment:Object(l.jsx)(xe.a.span,{style:h,children:Object(l.jsx)(q.a,{of:"keyboardDown"})}),collapsed:t,contents:n,onClick:c}),Object(l.jsx)("li",{children:Object(l.jsx)(N.a,{component:xe.a.ul,overflow:"hidden",style:j,children:i.map(({id:e,onClick:n,...i})=>Object(l.jsx)(_e,{disabled:!b,...i,onClick:()=>{n&&n(),s()},collapsed:t,icon:void 0},e))})})]})};var We=e=>void 0!==e.links?Object(l.jsx)(Be,{...e}):Object(l.jsx)(_e,{...e}),$e=n("Jk64");var Ke=({title:e,links:t=[],toolbarItems:n=[],accountId:r,variant:o="mini",expanded:c=!1,onClose:s,onToggle:d,onTempExpand:u})=>{const m=!!c||"full"===o,{topBarHeight:b}=Object(i.useContext)(Fe.a),[j,h]=Object(i.useState)(null),v=Object(i.useCallback)(e=>{m?h(j!==e?e:null):(u&&u(),h(e))},[h,u,j,m]),k=Object(ye.a)(),p=m?k.palette.text.primary:k.palette.primary.contrastText,f=m?Se:"mini"===o?Ne:0,O=b>0?window.innerHeight-b:"100vh",g=m?k.palette.background.default:k.palette.primary.main,x=Object(xe.b)({color:p,width:f,backgroundColor:g,outline:0,paddingTop:"env(safe-area-inset-top)",paddingBottom:"env(safe-area-inset-bottom)",borderColor:k.palette.divider,borderStyle:"solid",borderRightWidth:("full"===o?1:0)+"px"}),y=Object(l.jsxs)(N.a,{component:xe.a.div,position:"fixed",style:{...x,height:"100vh"},flexDirection:"column",paddingLeft:1,paddingRight:1,display:"flex",overflow:"hidden",zIndex:"1",children:[Object(l.jsxs)(N.a,{display:"flex",alignItems:"center",height:De,children:[Object(l.jsx)(K.a,{title:"Main Menu",children:Object(l.jsx)(E.a,{color:"inherit",onClick:d,"aria-expanded":m,size:"large",children:Object(l.jsx)(q.a,{of:"menu"})})})," ",Object(l.jsx)(I.a,{variant:"h6",marginLeft:1,children:e||Object(l.jsx)($e.a,{of:"sensu-wordmark"})})]}),Object(l.jsx)(Le,{color:p}),Object(l.jsx)(N.a,{component:"ul",display:"flex",flexGrow:1,flexDirection:"column",style:{overflowY:"auto",overflowX:"hidden"},children:t.map(e=>Object(i.createElement)(We,{...e,key:e.id,collapsed:!m,expanded:m&&j===e.id,onExpand:()=>v(e.id),onClick:()=>{const t=e.onClick;void 0!==t&&t(),s&&s()}}))}),Object(l.jsx)(Le,{color:p}),Object(l.jsx)(Ie,{accountId:r,isOpen:m,toolbarItems:n})]});return"mini"===o&&c?Object(l.jsxs)(a.a.Fragment,{children:[Object(l.jsx)(N.a,{width:Ne,flex:"0 0 auto",position:"relative"}),Object(l.jsx)(Te.a,{disableAutoFocus:!0,onClose:s,open:!0,children:y})]}):"hidden"===o?Object(l.jsx)(Te.a,{disableAutoFocus:!0,keepMounted:!0,onClose:s,open:c,children:y}):Object(l.jsxs)(a.a.Fragment,{children:[Object(l.jsx)(N.a,{style:{width:f},position:"relative",flex:"0 0 auto"}),a.a.cloneElement(y,{style:{...y.props.style,height:O}})]})};var Je=({content:e,...t})=>{const n=Object(y.a)("sm","lt");return Object(l.jsx)(ge.a,{...t,content:e,drawer:Object(l.jsx)(Ke,{variant:"mini"}),mobile:n,topBar:!1})};const Ye=a.a.lazy(()=>Promise.all([n.e(1),n.e(11),n.e(2),n.e(5),n.e(10),n.e(86)]).then(n.bind(null,"i36c"))),He=a.a.lazy(()=>Promise.all([n.e(1),n.e(11),n.e(38),n.e(2),n.e(5),n.e(31),n.e(10),n.e(79)]).then(n.bind(null,"XUE8")));t.default=()=>Object(l.jsxs)(a.a.Suspense,{fallback:Object(l.jsx)(a.a.Fragment,{}),children:[Object(l.jsxs)(r.e,{children:[Object(l.jsx)(r.c,{path:"/c/:cluster/n/:namespace",render:({match:e})=>Object(l.jsx)(Je,{content:Object(l.jsxs)(r.e,{children:[Object(l.jsx)(r.c,{exact:!0,path:e.path+"/catalog",component:Ye}),Object(l.jsx)(r.c,{path:e.path+"/catalog/:catalogNamespace/:integration",component:He})]})})}),Object(l.jsx)(r.b,{to:"/c/~/n/default/catalog"})]}),Object(l.jsx)(Oe,{}),Object(l.jsx)(le.a,{}),Object(l.jsx)(j,{}),Object(l.jsx)(k,{}),Object(l.jsx)(se,{}),Object(l.jsx)(p.a,{})]})},InSj:function(e,t,n){"use strict";var i=n("q1tI"),a=n.n(i),r=n("a/bO"),o=n("Vzdd"),c=n("nKUr");const s=({children:e})=>{const t={children:{render:()=>e,type:o.a}};return Object(c.jsx)(r.b,{...t})};t.a=a.a.memo(s)},J4nZ:function(e,t,n){"use strict";function i(e,...t){const n=Object.assign({},e);let i=0;for(const a in e){if(i===t.length)break;t.includes(a)&&(delete n[a],i++)}return n}function a(e,...t){const n={};let i=0;for(const a in e){if(i===t.length)break;t.includes(a)&&(n[a]=e[a],i++)}return n}function r(e,t){const n=Object.assign({},e);for(const i in e)t(i,e[i])||delete n[i];return n}function o(e,t){return Object.keys(t).reduce((e,n)=>e[n]&&0!==e[n]&&""!==e[n]?e:{...e,[n]:t[n]},e)}function c(e,t){if(e===t)return!0;if(Object.keys(e).length!==Object.keys(t).length)return!1;for(const n of Object.keys(e))if(e[n]!==t[n])return!1;return!0}function s(e,t){const[n,...i]=t.split(".");if(void 0!==e[n])return 0===i.length?e[n]:s(e[n],i.join("."))}n.d(t,"c",(function(){return i})),n.d(t,"d",(function(){return a})),n.d(t,"a",(function(){return r})),n.d(t,"e",(function(){return o})),n.d(t,"f",(function(){return c})),n.d(t,"b",(function(){return s}))},KQEc:function(e,t,n){"use strict";n.d(t,"a",(function(){return s}));var i=n("q1tI"),a=n.n(i),r=n("+uPJ"),o=n("nKUr");const c=a.a.lazy(()=>n.e(88).then(n.bind(null,"NH6r")));var s=e=>Object(o.jsx)(a.a.Suspense,{fallback:Object(o.jsx)(r.a,{variant:"icon"}),children:Object(o.jsx)(c,{...e})})},M0mT:function(e,t,n){"use strict";var i=n("q1tI"),a=n.n(i),r=n("Yd4P"),o=n.n(r),c=n("wEEd"),s=n("ANGJ"),l=n("DyER"),d=n("a/bO"),u=n("Vzdd"),m=n("nKUr");const b=Object(s.a)(e=>({toast:{position:"relative",left:0,right:0},toastPadding:{[e.breakpoints.up("md")]:{paddingBottom:10,paddingRight:10}}}));t.a=()=>{const e=b(),t=Object(d.d)().filter(e=>e.props.type===u.c).slice(-20),[n,i]=a.a.useState({}),r=a.a.useCallback((e,t)=>{i(n=>t.height===n[e]?n:{...n,[e]:t.height})},[i]),s=a.a.useCallback(e=>{i(t=>(void 0===t[e]||delete t[e],t))},[i]),j=Object(c.c)(t,e=>e.id,{from:{opacity:0,height:0},update:({id:e})=>({opacity:1,height:n[e]||0}),leave:{opacity:0,height:0},config:{tension:210,friction:20}}).map(({item:{id:t,props:n,remove:i},props:a,key:d})=>Object(m.jsx)(c.a.div,{style:a,className:e.toast,children:Object(m.jsxs)("div",{style:{position:"relative"},children:[Object(m.jsx)(o.a,{onResize:e=>r(t,e)}),Object(m.jsx)(l.a,{onUnmount:()=>s(t)}),Object(m.jsx)("div",{className:e.toastPadding,children:n.render({id:t,remove:i})})]})},d));return Object(m.jsx)(m.Fragment,{children:j})}},Q2fi:function(e,t,n){"use strict";var i=n("q1tI"),a=n.n(i);t.a=a.a.createContext(void 0)},Tzv7:function(e,t,n){"use strict";var i=n("q1tI"),a=n("7TpP");t.a=function(){return Object(i.useContext)(a.a)}},V4jd:function(e,t,n){"use strict";const i={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"mutation",name:{kind:"Name",value:"InvalidateTokensMutation"},variableDefinitions:[],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"invalidateTokens"},arguments:[],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}]}]}}],loc:{start:0,end:72,source:{body:"\n  mutation InvalidateTokensMutation {\n    invalidateTokens @client\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}};t.a=e=>e.mutate({mutation:i}).catch(e=>{throw e.networkError||e})},Vzdd:function(e,t,n){"use strict";n.d(t,"c",(function(){return i})),n.d(t,"a",(function(){return a})),n.d(t,"b",(function(){return r}));const i="TOAST",a="BANNER",r="MODAL"},WDlM:function(e,t,n){"use strict";var i=n("q1tI"),a=n.n(i),r=n("3ebn");t.a=()=>a.a.useMemo(r.a,[])},Zrqx:function(e,t,n){"use strict";n.d(t,"e",(function(){return i})),n.d(t,"c",(function(){return a})),n.d(t,"d",(function(){return r})),n.d(t,"a",(function(){return o})),n.d(t,"b",(function(){return c})),n.d(t,"f",(function(){return s}));const i=(e,t)=>{if(e===t)return!0;if(e.length!==t.length)return!1;for(let n=0;n<e.length;n+=1)if(e[n]!==t[n])return!1;return!0},a=(e,t,n)=>e.slice(0,t).concat([{...e[t],...n}]).concat(e.slice(t+1)),r=(e,t)=>e.slice(0,t).concat(e.slice(t+1)),o=(e,t)=>e.some(e=>t.includes(e)),c=(e,t)=>e.filter(e=>t.includes(e)),s=e=>Array.isArray(e)?e:void 0===e?[]:[e]},"a/bO":function(e,t,n){"use strict";n.d(t,"a",(function(){return u})),n.d(t,"b",(function(){return m})),n.d(t,"d",(function(){return b})),n.d(t,"c",(function(){return j}));var i=n("q1tI"),a=n.n(i),r=n("3ebn"),o=n("Zrqx"),c=n("WDlM"),s=n("nKUr");const l=a.a.createContext([]),d=a.a.createContext({setChild:()=>{},removeChild:()=>{}}),u=({children:e})=>{const[t,n]=a.a.useState([]),i=a.a.useMemo(()=>({removeChild:e=>{n(t=>{const n=t.findIndex(t=>t.id===e);return-1===n?t:Object(o.d)(t,n)})},setChild:(e=Object(r.a)(),t)=>(n(n=>{const a=n.findIndex(t=>t.id===e);if(-1===a){const a={id:e,props:t,update:t=>i.setChild(e,t),remove:()=>i.removeChild(e)};return n.concat([a])}return Object(o.c)(n,a,{props:t})}),e)}),[]);return Object(s.jsx)(d.Provider,{value:i,children:Object(s.jsx)(l.Provider,{value:t,children:e})})},m=({children:e})=>((e=>{const t=Object(c.a)(),{setChild:n,removeChild:i}=a.a.useContext(d);a.a.useEffect(()=>{n(t,e)},[e,t,n]),a.a.useEffect(()=>()=>i(t),[t,i])})(e),null),b=()=>a.a.useContext(l),j=()=>a.a.useContext(d)},eqBg:function(e,t,n){"use strict";var i=n("Tzv7"),a=n("J4nZ");t.a=function(e){const t=Object(i.a)();if(e){const n=t.pagePreferences.find(t=>t.page===e);if(n)return{...t,preferences:Object(a.e)(n,t.preferences)}}return{...t,preferences:Object(a.e)({pageSize:0,order:"",selector:""},t.preferences)}}},pCzp:function(e,t,n){"use strict";var i=n("q1tI"),a=n("6koa"),r=n.n(a),o=n("dMq0"),c=n("jNO+"),s=n("wAjv"),l=n("prjk");function d(e){const{data:t,loading:n,networkStatus:i}=e.currentResult();return{observable:e,data:t,loading:n,networkStatus:i,refetch:e.refetch,fetchMore:e.fetchMore,updateQuery:e.updateQuery,startPolling:e.startPolling,stopPolling:e.stopPolling,subscribeToMore:e.subscribeToMore}}function u(e){return 1===e.networkStatus?"initial":e.networkStatus<=6?"update":8===e.networkStatus||e.aborted?"aborted":void 0}const m={onError:e=>{if(null==e||!e.networkError||!Object(l.a)(e.networkError)){if(null!==e&&"object"==typeof e)throw e;throw new Error(e)}}};function b(e){const t=i.useRef(void 0);t.current={client:Object(s.a)(),...m,...e};const n=i.useRef(null),a=i.useRef(null),l=i.useRef(null);var b,j;null!==l.current&&null!==n.current&&(b=t.current,j=l.current,b.client!==j.client||b.query!==j.query?n.current=null:function(e,t){return e.pollInterval!==t.pollInterval||e.fetchPolicy!==t.fetchPolicy||e.errorPolicy!==t.errorPolicy||e.fetchResults!==t.fetchResults||e.notifyOnNetworkStatusChange!==t.notifyOnNetworkStatusChange||!r()(e.variables,t.variables)}(t.current,l.current)&&n.current.setOptions(t.current).catch(()=>null));const h=i.useCallback(()=>{a.current&&a.current.unsubscribe();const e=t.current.client.watchQuery(t.current);n.current=e,function(e){e.refetch=e.refetch.bind(e),e.fetchMore=e.fetchMore.bind(e),e.updateQuery=e.updateQuery.bind(e),e.startPolling=e.startPolling.bind(e),e.stopPolling=e.stopPolling.bind(e),e.subscribeToMore=e.subscribeToMore.bind(e)}(e),a.current=e.subscribe({next:e=>{const{data:n,errors:i,loading:a,networkStatus:r}=e;let c=null;if(i&&i.length>0&&(c=new o.a({graphQLErrors:i})),k(e=>({...e,aborted:!1,error:c,data:n,loading:a,networkStatus:r})),c){if(!t.current.onError)throw c;t.current.onError(c)}},error:e=>{if(e.networkError instanceof c.a)k(e=>({...e,aborted:!0,error:null}));else{if(k(t=>({...t,error:e})),!t.current.onError)throw e;t.current.onError(e)}}})},[]);null===n.current&&h();const[v,k]=i.useState(()=>({...d(n.current),error:null,aborted:!1,resubscribe:h}));return i.useEffect(()=>()=>{a.current.unsubscribe()},[]),l.current=t.current,{...v,loadingState:u(v)}}const j={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"query",name:{kind:"Name",value:"LocalNetworkStatusQuery"},variableDefinitions:[],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"localNetwork"},arguments:[],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"offline"},arguments:[],directives:[]},{kind:"Field",name:{kind:"Name",value:"retry"},arguments:[],directives:[]}]}}]}}],loc:{start:0,end:98,source:{body:"\n  query LocalNetworkStatusQuery {\n    localNetwork @client {\n      offline\n      retry\n    }\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}};t.a=function(e){const t=b({query:j,onError:e.onError}),n=i.useRef(!1),{data:{localNetwork:{offline:a=!1,retry:r=!1}={}}={}}=t,o=b(e),c=o.resubscribe,s=o.error;return i.useEffect(()=>{!n.current||a&&!r||c()},[a,r,c,s]),n.current=a,o}},uDTK:function(e,t,n){"use strict";var i=n("eqBg"),a=n("pCzp"),r=n("q1tI");var o=()=>{const e=window.matchMedia("(prefers-color-scheme: dark)"),[t,n]=Object(r.useState)(e.matches),i=e=>n(e.matches);return Object(r.useEffect)(()=>(e.addListener(i),()=>e.removeListener(i))),t};const c={kind:"Document",definitions:[{kind:"OperationDefinition",operation:"query",name:{kind:"Name",value:"ThemeProviderQuery"},variableDefinitions:[],directives:[],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"theme"},arguments:[],directives:[{kind:"Directive",name:{kind:"Name",value:"client"},arguments:[]}],selectionSet:{kind:"SelectionSet",selections:[{kind:"Field",name:{kind:"Name",value:"value"},arguments:[],directives:[]},{kind:"Field",name:{kind:"Name",value:"dark"},arguments:[],directives:[]}]}}]}}],loc:{start:0,end:83,source:{body:"\n  query ThemeProviderQuery {\n    theme @client {\n      value\n      dark\n    }\n  }\n",name:"GraphQL request",locationOffset:{line:1,column:1}}}};t.a=()=>{const e=Object(a.a)({query:c}),t=o(),n=Object(i.a)(),r=e.data;if(!r.theme)throw new Error("Unable to get stored theme preferences from the store. May indicate that the store is not availabe in the current context.");return{value:r.theme.value||n.preferences.theme,dark:"UNSET"!==r.theme.dark?"DARK"===r.theme.dark:t,usingSystemColourScheme:"UNSET"===r.theme.dark}}},vzfW:function(e,t,n){"use strict";n("CSxS");var i=n("UnQg");t.a=()=>Object(i.a)()},wAjv:function(e,t,n){"use strict";var i=n("q1tI"),a=n("Q2fi");t.a=function(){const e=i.useContext(a.a);if(!e||!e.client)throw new Error("useApolloClient must be used within ApolloProvider");return e.client}}}]);
//# sourceMappingURL=34_edf5.js.map