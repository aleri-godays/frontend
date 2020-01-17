(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["app"],{0:function(n,e,o){n.exports=o("2f39")},"2f39":function(n,e,o){"use strict";o.r(e);var t={};o.r(t),o.d(t,"oddLines",(function(){return S})),o.d(t,"getAllProjects",(function(){return L})),o.d(t,"getAllEntries",(function(){return P})),o.d(t,"getUser",(function(){return T}));var r={};o.r(r),o.d(r,"setLines",(function(){return E})),o.d(r,"setUser",(function(){return D})),o.d(r,"setProjects",(function(){return I})),o.d(r,"setEntries",(function(){return V}));var c={};o.r(c),o.d(c,"loadSample",(function(){return A})),o.d(c,"login",(function(){return J})),o.d(c,"authorize",(function(){return U})),o.d(c,"getProjects",(function(){return C})),o.d(c,"getEntries",(function(){return N})),o.d(c,"logout",(function(){return O})),o.d(c,"deleteTime",(function(){return _})),o.d(c,"addTime",(function(){return q})),o.d(c,"editTime",(function(){return $}));var i=o("967e"),l=o.n(i),u=(o("96cf"),o("7d6e"),o("e54f"),o("985d"),o("31cd"),o("2b0e")),a=o("b05d"),s=o("ca78"),f=o("7f67");u["a"].use(a["a"],{config:{},components:{QTime:s["a"],ClosePopup:f["a"]}});var d=function(){var n=this,e=n.$createElement,o=n._self._c||e;return o("div",{attrs:{id:"q-app"}},[o("router-view")],1)},p=[],g={name:"App"},m=g,h=o("2877"),v=Object(h["a"])(m,d,p,!1,null,null,null),w=v.exports,b=o("2f62"),j=o("bfa9"),y={lines:{all:[]},user:{name:"",token:"",loginState:!1},projects:{all:[{name:"",ID:0,entries:[]}]},entries:{all:[{id:0,date:0,project_id:0,comment:"",duration:0}]}};o("ac6a"),o("cadf"),o("06db");function S(n){return n.lines.all.filter((function(n){return n.ID%2===0}))}function L(n){return n.projects.all}function P(n){return n.entries.all}function T(n){return n.user}o("7f7f");function E(n,e){null===e&&(e=[]),n.lines.all=e}function D(n,e){null===e&&(e=""),n.user.name=e.login,n.user.token=e.jwt,n.user.loginState=!0}function I(n,e){null===e&&(e=[]),console.log("project items: ",e),n.projects.all=e,console.log("project state: ",n.projects.all)}function V(n,e){null===e&&(e=[]),n.entries.all=e}var k=o("bc3a"),x=o.n(k);function A(n){var e=n.commit,o=[{ID:1,Title:"Line1",Value:1},{ID:2,Title:"Line2",Value:2},{ID:3,Title:"Line3",Value:3}];e("setLines",o)}function J(n){var e=n.commit;console.log("calling /me"),x.a.get("/api/v1/frontend/me").then((function(n){e("setUser",n.data),console.log("response: ",n)})).catch((function(n){console.log("error, redirecting to /login",n),window.location="/login"}))}function U(n){var e=n.commit;console.log("calling /me"),x.a.get("/api/v1/frontend/me").then((function(n){e("setUser",n.data),console.log("response: ",n)})).catch((function(n){console.log("error ",n)}))}function C(n){var e=n.commit;x.a.get("/api/v1/frontend/project").then((function(n){e("setProjects",n.data),console.log("project response: ",n)})).catch((function(n){return console.log(n)}))}function N(n){var e=n.commit;x.a.get("/api/v1/frontend/entry").then((function(n){e("setEntries",n.data),console.log("entries response: ",n)})).catch((function(n){return console.log(n)}))}function O(n){n.commit;x.a.get("/logout").then((function(n){localStorage.clear(),console.log("response: ",n)})).then(window.location="/").catch((function(n){return console.log("error",n)}),window.location="/logout",localStorage.clear())}function _(n,e){n.commit;x.a.delete("/api/v1/frontend/entry/"+e).then((function(n){console.log("response: ",n)})).then(window.location.reload()).catch((function(n){console.log("error, could not log out",n)}))}function q(n,e){n.commit;console.log("timeToSubmit: "+JSON.stringify(e)),x.a.post("/api/v1/frontend/entry",e).then((function(n){console.log("response: ",n)})).then(window.location.reload()).catch((function(n){console.log("error, could not log out",n)}))}function $(n,e){n.commit;console.log("timeToSubmit: "+JSON.stringify(e)),x.a.put("/api/v1/frontend/entry/"+e.id,e).then((function(n){console.log("response: ",n)})).then(window.location.reload()).catch((function(n){console.log("error, could not log out",n)}))}var z={namespaced:!0,getters:t,mutations:r,actions:c,state:y},B=o("85ff"),M=o.n(B),Q=!0,F={isEnabled:!0,logLevel:Q?"error":"debug",stringifyArguments:!1,showLogLevel:!0,showMethodName:!0,separator:"|",showConsoleColors:!0};u["a"].use(M.a,F),u["a"].use(b["a"]);var G=function(){var n=new b["a"].Store({modules:{vault:z},plugins:[(new j["a"]).plugin],strict:!1});return n},H=o("8c4f"),K=[{path:"/",component:function(){return Promise.all([o.e("559d2ce5"),o.e("2d22c0ff")]).then(o.bind(null,"f241"))},children:[{path:"",component:function(){return Promise.all([o.e("559d2ce5"),o.e("2d0e8be2")]).then(o.bind(null,"8b24"))}},{path:"/projects",component:function(){return Promise.all([o.e("559d2ce5"),o.e("2f9afefb")]).then(o.bind(null,"7f1d"))}}]}];K.push({path:"*",component:function(){return Promise.all([o.e("559d2ce5"),o.e("4b47640d")]).then(o.bind(null,"e51e"))}});var R=K;u["a"].use(H["a"]);var W=function(){var n=new H["a"]({scrollBehavior:function(){return{x:0,y:0}},routes:R,mode:"hash",base:""});return n},X=function(){var n="function"===typeof G?G({Vue:u["a"]}):G,e="function"===typeof W?W({Vue:u["a"],store:n}):W;n.$router=e;var o={el:"#q-app",router:e,store:n,render:function(n){return n(w)}};return{app:o,store:n,router:e}},Y=X(),Z=Y.app;Y.store,Y.router;function nn(){return l.a.async((function(n){while(1)switch(n.prev=n.next){case 0:new u["a"](Z);case 1:case"end":return n.stop()}}))}nn()},"31cd":function(n,e,o){}},[[0,"runtime","vendor"]]]);