import { createRouter, createWebHistory } from "vue-router"; // createWebHashHistory, createWebHistory
import type { Router, RouteRecordRaw, RouteComponent } from "vue-router";

/* Layout */
const Layout = (): RouteComponent => import("@/layout/index.vue");

/**
 * 静态路由
 * path: 路由地址
 * component: 组件
 * name: 路由名称
 * meta: 路由元信息
 * redirect: 重定向地址
 * hidden: 是否隐藏路由
 * alwaysShow: 是否一直显示根路由
 * icon: 路由图标
 * noCache: 是否缓存路由
 * affix: 是否固定标签
 */
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: "/login",
    component: () => import("@/views/login/index.vue"),
    meta: { hidden: true },
  },
  {
    path: "/404",
    component: () => import("@/views/error-page/404.vue"),
    meta: { hidden: true },
  },
  {
    path: "/401",
    component: () => import("@/views/error-page/401.vue"),
    meta: { hidden: true },
  },
  // {
  //   path: "/theme",
  //   component: Layout,
  //   children: [
  //     {
  //       path: "index",
  //       component: () => import("@/views/theme/index.vue"),
  //       name: "Theme",
  //       meta: { title: "主题", icon: "theme", hidden: true },
  //     },
  //   ],
  // },
  // {
  //   path: "/icons",
  //   component: Layout,
  //   children: [
  //     {
  //       path: "index",
  //       component: () => import("@/views/icons/index.vue"),
  //       name: "icons",
  //       meta: { title: "icons", icon: "icons" },
  //     },
  //   ],
  // },
  {
    path: "/",
    component: Layout,
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        component: () => import("@/views/dashboard/index.vue"),
        name: "Dashboard",
        meta: { title: "首页", icon: "dashboard", affix: true },
      },
    ],
  },

  {
    path: "/profile",
    component: Layout,
    redirect: "/profile/index",
    meta: { hidden: true },
    children: [
      {
        path: "index",
        component: () => import("@/views/profile/index.vue"),
        name: "Profile",
        meta: { title: "个人中心", icon: "user", noCache: true },
      },
    ],
  },
  {
    path: "/:pathMatch(.*)*", // 通配符路径匹配所有未匹配的路径
    component: () => import("@/views/error-page/404.vue"), // 重定向到404页面
    meta: { hidden: true },
  },
];

const createTheRouter = (): Router =>
  createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    scrollBehavior: () => ({ top: 0 }),
    routes: constantRoutes,
  });

interface RouterPro extends Router {
  matcher: unknown;
}

const router = createTheRouter() as RouterPro;

export function resetRouter() {
  const newRouter = createTheRouter() as RouterPro;
  router.matcher = newRouter.matcher; // reset router
}

export default router;
