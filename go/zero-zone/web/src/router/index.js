import { createRouter, createWebHashHistory } from "vue-router";

// 本地静态路由
export const constantRoutes = [
  /*
                                 {
                                   path: '/',
                                   component: () => import('@/views/dashboard/index.vue'),
                                 },
                                   */

  {
    path: "/ucenter",
    component: () => import("@/views/system/ucenter.vue"),
    meta: {
      isParentView: false,
    },
  },
  {
    path: "/login",
    component: () => import("@/views/login/login.vue"),
    meta: {
      isParentView: true,
    },
  },
  {
    path: "/search",
    component: () => import("@/views/help/search.vue"),
  },
  {
    path: "/curd",
    component: () => import("@/views/help/curd.vue"),
  },
  {
    path: "/test",
    component: () => import("@/views/help/index.vue"),
  },
  {
    path: "/open-auth",
    component: () => import("@/components/feat/open_auth.vue"),
    meta: {
      isParentView: true,
    }
  },
  {
    path: "/test-layout",
    component: () => import("@/views/help/index.vue"),
    // component: () => import('@/layout/index.vue'),
    // component: () => import('@/layout/parentView.vue'),
    meta: {
      isParentView: true,
      xxx: true,
    },
    children: [
      {
        path: "xxx", // 加斜杠 全路径，不加的话会拼接父path eg：/test-layout/xxx
        component: () => import("@/views/help/index.vue"),
      },
    ],
  },
  {
    // path: '/404',
    path: "/:pathMatch(.*)*", // 防止浏览器刷新时路由未找到警告提示: vue-router.mjs:35 [Vue Router warn]: No match found for location with path "/xxx"
    component: () => import("@/views/error-page/404.vue"),
  },
];

// 创建路由
const router = createRouter({
  history: createWebHashHistory(),
  routes: constantRoutes,
});

export default router;
