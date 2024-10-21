import { defineStore } from "pinia";
import { constantRoutes } from "@/router";
import type { RouteRecordRaw } from "vue-router";
import { getUserMenuTreeByUserId } from "@/api/system/menu";
import router from "@/router"; // 确保正确导入 router 实例
import store from "@/store";
interface IPermissionState {
  routes: Array<RouteRecordRaw>;
  addRoutes: Array<RouteRecordRaw>;
}

// 定义 Layout 组件的导入
const Layout = () => import("@/layout/index.vue");

interface Menu {
  path: string;
  name: string;
  component?: string;
  icon?: string;
  hidden?: number;
  title?: string;
  redirect?: string;
  alwaysShow?: number;
  noCache?: number;
  breadcrumb?: number;
  activeMenu?: string;
  children?: Menu[];
}
const modules = import.meta.glob("../../views/**/**.vue");

// 递归转换函数
export const getRoutesFromMenuTree = (menuTree: Menu[]): RouteRecordRaw[] => {
  return menuTree.map((menu) => {
    let component;
    if (menu.component === "Layout") {
      component = Layout;
    } else if (menu.component) {
      // 确保路径正确，并且使用箭头函数包裹动态导入
      const componentPath = `../../views/${menu.component.replace(".vue", "")}`;
      component =
        modules[componentPath] ||
        (() => import(/* @vite-ignore */ `../../views/${menu.component}`));
    }

    const route: any = {
      path: menu.path,
      name: menu.name,
      component, // 这里 component 应该是一个函数
      // hidden: menu.hidden === 1,
      redirect: menu.redirect,
      // alwaysShow: menu.alwaysShow === 1,
      children: menu.children
        ? getRoutesFromMenuTree(menu.children)
        : undefined,
      meta: {
        title: menu.title,
        icon: menu.icon,
        noCache: menu.noCache === 1,
        breadcrumb: menu.breadcrumb === 1,
        activeMenu: menu.activeMenu,
      },
    };
    return route;
  });
};

export default defineStore({
  id: "permission",
  state: (): IPermissionState => ({
    routes: [],
    addRoutes: [],
  }),
  getters: {},
  actions: {
    setRoutes(routes: RouteRecordRaw[]) {
      this.addRoutes = routes;
      this.routes = constantRoutes.concat(routes);
      routes.forEach((route) => {
        router.addRoute(route); // 动态添加路由
      });
    },
    async generateRoutes() {
      try {
        // 获取用户菜单树
        const id = store.user().userId;
        const res = (await getUserMenuTreeByUserId(id)) as any;
        const menuTree = res.data; // 确保这里的数据结构正确

        // 将菜单树转换为路由配置
        const accessedRoutes = getRoutesFromMenuTree(menuTree);

        // 设置路由
        this.setRoutes(accessedRoutes);
        return accessedRoutes;
      } catch (error) {
        console.error("Failed to generate routes:", error);
      }
    },
  },
});
