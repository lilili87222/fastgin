<template>
  <el-drawer v-model="drawer" :direction="direction" :size="'20%'">
    <template #header>
      <h4>系统设置</h4>
    </template>
    <template #default>
      <div class="drawer-container">
        <div>
          <div class="drawer-item">
            <span>开启 Tags-View</span>
            <el-switch v-model="tagsView" class="drawer-switch" />
          </div>

          <div class="drawer-item">
            <span>固定 Header</span>
            <el-switch v-model="fixedHeader" class="drawer-switch" />
          </div>

          <div class="drawer-item">
            <span>侧边栏 Logo</span>
            <el-switch v-model="sidebarLogo" class="drawer-switch" />
          </div>

          <div class="drawer-item">
            <span>侧边栏子目录气泡显示</span>
            <el-switch v-model="secondMenuPopup" class="drawer-switch" />
          </div>
        </div>
      </div>
    </template>
    <template #footer>
      <div style="flex: auto">
        <el-button @click="cancelClick">关闭</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";
import type { DrawerProps } from "element-plus";
import store from "@/store";

const drawer = ref(false);
const direction = ref<DrawerProps["direction"]>("rtl");

const emits = defineEmits(["getDictData"]);

const tagsView = computed({
  get: () => store.settings().tagsView,
  set: (val) => {
    store.settings().changeSetting({
      key: "tagsView",
      value: val,
    });
  },
});

const secondMenuPopup = computed({
  get: () => store.settings().secondMenuPopup,
  set: (val) => {
    store.settings().changeSetting({
      key: "secondMenuPopup",
      value: val,
    });
  },
});
const sidebarLogo = computed({
  get: () => store.settings().sidebarLogo,
  set: (val) => {
    store.settings().changeSetting({
      key: "sidebarLogo",
      value: val,
    });
  },
});
const fixedHeader = computed({
  get: () => store.settings().fixedHeader,
  set: (val) => {
    store.settings().changeSetting({
      key: "fixedHeader",
      value: val,
    });
  },
});

//打开
const openDrawer = (type: string, Id?: number) => {
  drawer.value = true;
};

defineExpose({
  openDrawer,
});

//取消
const cancelClick = () => {
  drawer.value = false;
};
</script>
<style lang="scss">
.el-drawer {
  max-width: 90%; /* 设置最大宽度为 90% */
  min-width: 300px; /* 设置最小宽度为 300px */
}
.drawer-container {
  padding: 24px;
  font-size: 14px;
  line-height: 1.5;
  word-wrap: break-word;

  .drawer-title {
    margin-bottom: 12px;
    color: rgba(0, 0, 0, 0.85);
    font-size: 14px;
    line-height: 22px;
  }

  .drawer-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: rgba(0, 0, 0, 0.65);
    font-size: 14px;
    padding: 12px 0;
  }

  .drawer-switch {
    float: right;
  }
}
</style>
