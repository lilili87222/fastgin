import { createApp } from "vue";

import App from "./App.vue";
import router from "./router";
import { setupStore } from "./store";

import "@/styles/index.scss";
import SvgIcon from "./icons"; // icon
import "./permission"; // permission control
import vPermission from "./directive/permission/index"; // permission control
import { checkEnableLogs } from "./utils/error-log"; // error log
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

const app = createApp(App);
setupStore(app);
Object.entries(ElementPlusIconsVue).forEach(([key, component]) => {
  app.component(key, component);
});
app.use(router);
app.component("svg-icon", SvgIcon);
app.directive("permission", vPermission);
checkEnableLogs(app);

app.mount("#app");
