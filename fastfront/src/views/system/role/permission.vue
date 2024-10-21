<template>
  <el-drawer v-model="drawer" :direction="direction" :size="'50%'">
    <template #header>
      <h4>修改权限</h4>
    </template>
    <template #default>
      <el-tabs v-model="activeName" class="demo-tabs">
        <el-tab-pane label="角色菜单" name="menu">
          <el-tree
            ref="roleMenuTree"
            v-loading="menuTreeLoading"
            :props="{ children: 'children', label: 'title' }"
            :data="menuData"
            show-checkbox
            node-key="id"
            check-strictly
            :default-checked-keys="defaultCheckedRoleMenu"
          />
        </el-tab-pane>
        <el-tab-pane label="角色接口" name="api">
          <el-tree
            ref="roleApiTree"
            v-loading="apiTreeLoading"
            :props="{ children: 'children', label: 'des' }"
            :data="apiData"
            show-checkbox
            node-key="id"
            :default-checked-keys="defaultCheckedRoleApi"
          />
        </el-tab-pane>
      </el-tabs>
    </template>
    <template #footer>
      <div style="flex: auto">
        <el-button @click="cancelClick">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessage } from "element-plus";
import type { DrawerProps } from "element-plus";

import {
  getRoleApisById,
  getRoleMenusById,
  updateRoleApisById,
  updateRoleMenusById,
} from "@/api/system/role";
import { getApiTree } from "@/api/system/api";
import { getMenuTree } from "@/api/system/menu";
import type { TRoleForm } from "@/types/system/role";
import type { TMenuTable } from "@/types/system/menu";
import type { TApiTable } from "@/types/system/api";

const drawer = ref(false);
const direction = ref<DrawerProps["direction"]>("rtl");
const emits = defineEmits(["getRoleData"]);

const activeName = ref("menu");

//获取角色接口
const apiData = ref<TApiTable[]>([]);
const apiTreeLoading = ref(false);
const getApiData = () => {
  apiTreeLoading.value = true;
  getApiTree().then((res) => {
    apiData.value = res.data;
    apiTreeLoading.value = false;
  });
};

//获取角色菜单
const menuData = ref<TMenuTable[]>([]);
const menuTreeLoading = ref(false);
const getMenuData = () => {
  menuTreeLoading.value = true;
  getMenuTree().then((res) => {
    menuData.value = res.data;
    menuTreeLoading.value = false;
  });
};

const roleMenuTree = ref();
const roleApiTree = ref();
const defaultCheckedRoleMenu = ref<number[]>([]);
const defaultCheckedRoleApi = ref<number[]>([]);
//获取当前角色拥有的菜单和接口
const getTree = () => {
  getRoleMenusById(roleId.value).then((res) => {
    const menus = res.data;
    const menuIds: number[] = [];
    menus.forEach((x: { id: any }) => {
      menuIds.push(x.id);
    });
    defaultCheckedRoleMenu.value = menuIds;
    roleMenuTree.value.setCheckedKeys(defaultCheckedRoleMenu.value);
  });
  getRoleApisById(roleId.value).then((res) => {
    const apis = res.data;
    const apiIds: number[] = [];
    apis.forEach((x: { id: any }) => {
      apiIds.push(x.id);
    });
    defaultCheckedRoleApi.value = apiIds;
    roleApiTree.value.setCheckedKeys(defaultCheckedRoleApi.value);
  });
};

const roleId = ref(0);
//打开
const openDrawer = (row: TRoleForm) => {
  roleId.value = row.id || 0;
  getTree();
  getMenuData();
  getApiData();
  drawer.value = true;
};

defineExpose({
  openDrawer,
});

//取消
const cancelClick = () => {
  drawer.value = false;
};

//提交
const submitForm = async () => {
  let menuIds = roleMenuTree.value.getCheckedKeys();
  const idsHalf = roleMenuTree.value.getHalfCheckedKeys();
  menuIds = menuIds.concat(idsHalf);
  menuIds = [...new Set(menuIds)];
  const apiIds = roleApiTree.value.getCheckedKeys(true);
  await updateRoleMenusById(roleId.value, { menuIds }).then((res) => {
    ElMessage.success(res.message);
  });
  await updateRoleApisById(roleId.value, { apiIds }).then((res) => {
    ElMessage.success(res.message);
  });
  emits("getRoleData");
  drawer.value = false;
};
</script>
<style>
.el-drawer {
  max-width: 90%; /* 设置最大宽度为 90% */
  min-width: 300px; /* 设置最小宽度为 300px */
}
</style>
