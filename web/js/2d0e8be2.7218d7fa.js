(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["2d0e8be2"],{"8b24":function(e,t,n){"use strict";n.r(t);var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("q-page",{staticClass:"q-pa-lg"},[n("h5",{staticClass:"q-mt-none"},[e._v("Startpage")])])},c=[],o=(n("8e6e"),n("8a81"),n("ac6a"),n("cadf"),n("06db"),n("456d"),n("c47a")),a=n.n(o),s=n("2f62");function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(n,!0).forEach((function(t){a()(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(n).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}var p={name:"PageIndex",data:function(){return{leftDrawerOpen:!1}},created:function(){this.$store.dispatch("vault/loadSample"),this.$store.dispatch("vault/authorize")},computed:l({},Object(s["c"])({myLines:function(e){return e.vault.lines.all}}),{},Object(s["b"])({myOddLines:"vault/oddLines"}),{},Object(s["b"])({myProjects:"vault/getAllProjects"}))},u=p,b=n("2877"),f=n("fe09"),O=Object(b["a"])(u,r,c,!1,null,null,null);t["default"]=O.exports;O.options.components=Object.assign({QPage:f["r"]},O.options.components||{})}}]);
