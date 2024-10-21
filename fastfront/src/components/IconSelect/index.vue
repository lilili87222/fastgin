<template>
  <el-tabs v-model="activeName" @tab-click="handleClick">
    <el-tab-pane label="SVG" name="SVG">
      <div class="icon-body">
        <div class="icon-list">
          <div
            v-for="(item, index) in nameArr"
            :key="index"
            @click="selectedIcon(item)"
            :class="{ highlight: props.iconVal && item === props.iconVal }"
          >
            <svg-icon :icon-class="item" />

            <span class="name">{{ item }}</span>
          </div>
        </div>
      </div>
    </el-tab-pane>
    <el-tab-pane label="Element Icons" name="ElementIcons">
      <div class="icon-body">
        <div class="icon-list">
          <div
            v-for="item of elementIcons"
            :key="item"
            @click="selectedIcon(item)"
            :class="{ highlight: props.iconVal && item === props.iconVal }"
          >
            <el-icon style="width: 30px; height: 20px">
              <component :is="item" class="svg-icon disabled"></component>
            </el-icon>
            <span class="name">{{ item }}</span>
          </div>
        </div>
      </div>
    </el-tab-pane>
  </el-tabs>
</template>

<script setup lang="ts">
import ElementPlusIconsVue from "@/views/icons/element-icons";
import nameArr from "./requireIcons";
import { ref } from "vue";

const emits = defineEmits(["selectIcon", "selectName"]);

const props = defineProps({
  iconVal: {
    type: String,
    default: " ",
  },
});

const activeName = ref("SVG");
const elementIcons = Object.keys(ElementPlusIconsVue);
// 点击图标时触发
const selectedIcon = (name: string) => {
  emits("selectIcon", name);
};

const handleClick = (tab: any) => {
  emits("selectName", tab.paneName);
};
</script>

<style lang="scss" scoped>
.icon-body {
  width: 100%;
  padding: 10px;

  .icon-list {
    height: 200px;
    overflow-y: scroll;
    display: flex;
    flex-wrap: wrap; // 使用 flex 布局，自动适应移动设备
    justify-content: space-between;

    @media (min-width: 1400px) {
      div {
        width: 25%; // 大屏幕：每行显示 4 个图标
      }
    }

    @media (min-width: 768px) and (max-width: 1399px) {
      div {
        width: 25%; // 中屏幕：每行显示 4 个图标
      }
      .name {
        display: none;
      }
    }

    @media (max-width: 767px) {
      div {
        width: 33.33%; // 小屏幕：每行显示 3 个图标
      }
      .name {
        display: none;
      }
    }

    div {
      height: 50px; // 设置高度以便于居中
      margin-bottom: 10px;
      cursor: pointer;
      transition: background-color 0.3s;
      display: flex; // 设置为 flex 布局
      align-items: center; // 垂直居中
      justify-content: center; // 水平居中
    }

    .svg-icon {
      height: 30px;
      margin-right: 5px; // 图标和文本之间的间距
    }

    span {
      display: inline-block;
      vertical-align: middle; // 中垂直对齐
      fill: currentColor;
      overflow: hidden;
    }
  }

  .highlight {
    background-color: #e0f7fa;
    border-radius: 5px;
  }
}
</style>
