import { createApp } from "vue";
import App from "./App.vue";

const app = createApp(App);

// 路由
import router from "@/router";

app.use(router);

// element-plus
import ElementPlus from "element-plus";
import zhCn from "element-plus/es/locale/lang/zh-cn";
//import "element-plus/dist/index.css";

app.use(ElementPlus, {
  locale: zhCn,
});

// 自定义样式
import "@/styles/index.scss";

// 注册所有图标
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

// pinia
import { createPinia } from "pinia";

const pinia = createPinia();
// 持久化存储
import { createPersistedState } from "pinia-plugin-persistedstate";

pinia.use(
  createPersistedState({
    auto: true, // 启用所有 Store 默认持久化
  }),
);
// 重写 $reset 方法 => 解决组合式api中无法使用问题
pinia.use(({ store }) => {
  const initialState = JSON.parse(JSON.stringify(store.$state));
  store.$reset = () => {
    store.$patch(initialState);
  };
});
app.use(pinia);
// store
import store from "@/store";

app.config.globalProperties.$store = store;

// 配置全局api
import api from "@/api";

app.config.globalProperties.$api = api;

// 全局组件注册
import myComponent from "@/components/index";

Object.keys(myComponent).forEach((key) => {
  app.component(key, myComponent[key]);
});

// 混入
import mixin from "@/utils/mixin";

app.mixin(mixin);

// 全局过滤器
import { filters } from "@/utils/filters.js";

app.config.globalProperties.$filters = filters;

// 动态路由权限
import "@/router/permission.js";

app.mount("#app");
