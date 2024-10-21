import { defineStore } from "pinia";
import {
  // login as apiLogin,
  // logout as apiLogout,
  getInfo as apiGetInfo,
} from "@/api/system/user";
import { login as apiLogin, logout as apiLogout } from "@/api/system/base";
import { getToken, setToken, removeToken } from "@/utils/auth";
import router, { resetRouter } from "@/router";
import tagsViewStore from "./tagsView";
import permissionStore from "./permission";

export interface IUserState {
  token: string;
  userId: string;
  name: string;
  avatar: string;
  introduction: string;
  roles: string[];
}

export default defineStore({
  id: "user",
  state: (): IUserState => ({
    token: getToken(),
    userId: "",
    name: "",
    avatar: "",
    introduction: "",
    roles: [],
  }),
  getters: {},
  actions: {
    // user login
    login(userInfo): Promise<void> {
      const { username, password } = userInfo;
      return new Promise((resolve, reject) => {
        apiLogin({ user_name: username.trim(), password: password })
          .then((response) => {
            const { data: Data } = response;
            this.token = Data.Token;
            setToken(Data.Token);
            resolve();
          })
          .catch((error) => {
            reject(error);
          });
      });
    },

    // get user info
    getInfo() {
      return new Promise((resolve, reject) => {
        apiGetInfo(this.token)
          .then((response) => {
            const { data: Data } = response;

            if (!Data) {
              reject("Verification failed, please Login again.");
            }

            const { Id, Roles, UserName, Avatar, Introduction } = Data;

            // roles must be a non-empty array
            if (!Roles || Roles.length <= 0) {
              reject("getInfo: roles must be a non-null array!");
            }
            this.userId = Id;
            this.roles = Roles;
            this.name = UserName;
            this.avatar = Avatar;
            this.introduction = Introduction;

            resolve(Data);
          })
          .catch((error) => {
            reject(error);
          });
      });
    },

    // user logout
    logout(): Promise<void> {
      return new Promise((resolve, reject) => {
        apiLogout()
          .then(() => {
            this.token = "";
            this.roles = [];
            removeToken();
            resetRouter();

            // reset visited views and cached views
            // to fixed https://github.com/PanJiaChen/vue-element-admin/issues/2485
            tagsViewStore().delAllViews();

            resolve();
          })
          .catch((error) => {
            reject(error);
          });
      });
    },

    // remove token
    resetToken() {
      this.token = "";
      this.roles = [];
      removeToken();
    },

    // dynamically modify permissions
    async changeRoles(role) {
      const token = role + "-token";

      this.token = token;
      setToken(token);

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const infoRes = (await this.getInfo()) as any;

      let roles = [];
      if (infoRes.roles) {
        roles = infoRes.roles;
      }

      resetRouter();

      // generate accessible routes map based on roles
      const accessRoutes = await permissionStore().generateRoutes();
      // dynamically add accessible routes
      // router.addRoutes(accessRoutes);

      accessRoutes?.forEach((item) => {
        router.addRoute(item);
      });

      // reset visited views and cached views
      tagsViewStore().delAllViews();
    },
  },
});
