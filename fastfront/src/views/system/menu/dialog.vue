<template>
  <el-drawer
    v-model="drawer"
    :direction="direction"
    :size="'50%'"
    @close="drawerClose"
  >
    <template #header>
      <h4>{{ dialogType === "create" ? "新增菜单" : "编辑菜单" }}</h4>
    </template>
    <template #default>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item
          v-for="item in fromCol"
          :key="item.prop"
          :label="item.label"
          :prop="item.prop"
        >
          <template v-if="item.prop === 'Sort'">
            <el-input-number
              v-model="formData[item.prop]"
              controls-position="right"
              :min="1"
              :max="999"
            />
          </template>

          <template v-else-if="item.prop === 'Icon'">
            <el-input
              v-model="formData[item.prop]"
              readonly
              :placeholder="item.placeholder"
            >
              <template #prefix>
                <template v-if="activeName === 'SVG'">
                  <svg-icon :icon-class="formData[item.prop] || ''" />
                </template>
                <template v-else>
                  <el-icon>
                    <component
                      :is="formData[item.prop]"
                      class="svg-icon disabled"
                    ></component>
                  </el-icon>
                </template>
              </template>
            </el-input>

            <IconSelect
              :iconVal="formData[item.prop] || ' '"
              @selectIcon="handleIconSelect"
              @select-name="handleIconName"
            />
          </template>
          <template
            v-else-if="['Status', 'Hidden', 'NoCache'].includes(item.prop)"
          >
            <el-radio-group v-model="formData[item.prop]">
              <el-radio-button label="是" :value="2" />
              <el-radio-button label="否" :value="1" />
            </el-radio-group>
          </template>
          <template v-else-if="item.prop === 'Hidden'">
            <el-radio-group v-model="formData[item.prop]">
              <el-radio-button label="是" :value="1" />
              <el-radio-button label="否" :value="2" />
            </el-radio-group>
          </template>
          <template v-else-if="item.prop === 'ParentId'">
            <el-tree-select
              ref="treeSelectRef"
              v-model="formData[item.prop]"
              :data="treeselect"
              check-strictly
              filterable
              show-checkbox
              node-key="Id"
              :props="defaultProps"
              style="width: 100%"
              @node-click="handleNodeClick"
            />
          </template>
          <template v-else>
            <el-input
              v-model="formData[item.prop]"
              :placeholder="item.placeholder"
            ></el-input>
          </template>
        </el-form-item>
      </el-form>
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
import { reactive, ref } from "vue";
import type { DrawerProps, FormRules } from "element-plus";

import { createRole, updateRoleById } from "@/api/system/role";
import IconSelect from "@/components/IconSelect/index.vue";
import type { TMenuFormData } from "@/types/system/menu";

const drawer = ref(false);
const direction = ref<DrawerProps["direction"]>("rtl");
const emits = defineEmits(["getMenuData"]);

const defaultProps = ref({
  children: "Children",
  label: "Title",
});

const rules = reactive<FormRules>({
  Title: [
    { required: true, message: "请输入标题", trigger: "blur" },
    { min: 1, max: 50, message: "长度在 1 到 50 个字符", trigger: "blur" },
  ],
  Name: [
    { required: true, message: "请输入名称", trigger: "blur" },
    { min: 1, max: 100, message: "长度在 1 到 100 个字符", trigger: "blur" },
  ],
  Path: [
    { required: true, message: "请输入访问路径", trigger: "blur" },
    { min: 1, max: 100, message: "长度在 1 到 100 个字符", trigger: "blur" },
  ],
  Icon: [{ required: true, message: "请选择图标", trigger: "blur" }],
  Sort: [{ required: true, message: "请输入排序", trigger: "blur" }],
  Component: [
    { required: true, message: "请输入组件路径", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
  Redirect: [
    { required: false, message: "请输入重定向", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
  ActiveMenu: [
    { required: false, message: "请输入高亮菜单", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
  ParentId: [{ required: true, message: "请选择上级目录", trigger: "change" }],
});

const fromCol = [
  { prop: "Title", label: "菜单标题", placeholder: "菜单标题" },
  { prop: "Name", label: "名称", placeholder: "名称" },
  { prop: "Sort", label: "排序", placeholder: "排序" },
  { prop: "Icon", label: "图标", placeholder: "图标" },
  { prop: "Path", label: "路由地址", placeholder: "路由地址" },
  { prop: "Component", label: "组件路径", placeholder: "组件路径" },
  { prop: "Redirect", label: "重定向", placeholder: "重定向" },
  { prop: "Status", label: "禁用", placeholder: "禁用" }, //1 否 2是
  { prop: "Hidden", label: "隐藏", placeholder: "隐藏" }, //2否 1 是
  { prop: "NoCache", label: "缓存", placeholder: "缓存" }, //2 是 1 否
  { prop: "ActiveMenu", label: "高亮菜单", placeholder: "高亮菜单" },
  { prop: "ParentId", label: "上级目录", placeholder: "上级目录" },
];

const dialogType = ref("");
const treeselect = ref([]);
//打开
const openDrawer = (row: TMenuFormData, type: string, treeselectData: any) => {
  formData.value = row;
  dialogType.value = type;
  treeselect.value = treeselectData;
  drawer.value = true;
};

defineExpose({
  openDrawer,
});

const formData = ref<any>({});

const handleNodeClick = (data: any) => {
  formData.value.ParentId = data.Id;
};

const treeSelectRef = ref();
//关闭
const drawerClose = () => {
  formData.value = {};
  formRef.value.resetFields();
};

//取消
const cancelClick = () => {
  drawer.value = false;
};

const formRef = ref();

//图标选择
const activeName = ref("SVG");
const handleIconSelect = (iconName: string) => {
  formData.value.Icon = iconName;
};

const handleIconName = (iconName: string) => {
  formData.value.Icon = "";
  activeName.value = iconName;
};

//提交
const submitForm = () => {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      let data = { ...formData.value };

      if (dialogType.value === "create") {
        createRole(data).then((res) => {
          ElMessage.success(res.message);
          formRef.value.resetFields();
          emits("getMenuData");
          drawer.value = false;
        });
      } else {
        updateRoleById(data.Id, data).then((res) => {
          ElMessage.success(res.message);
          formRef.value.resetFields();
          emits("getMenuData");
          drawer.value = false;
        });
      }
    }
  });
};
</script>
<style>
.el-drawer {
  max-width: 90%; /* 设置最大宽度为 90% */
  min-width: 300px; /* 设置最小宽度为 300px */
}
</style>
