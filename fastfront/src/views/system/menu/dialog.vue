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
          <template v-if="item.prop === 'sort'">
            <el-input-number
              v-model="formData[item.prop]"
              controls-position="right"
              :min="1"
              :max="999"
            />
          </template>

          <template v-else-if="item.prop === 'icon'">
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
            v-else-if="['status', 'hidden', 'noCache'].includes(item.prop)"
          >
            <el-radio-group v-model="formData[item.prop]">
              <el-radio-button label="是" :value="2" />
              <el-radio-button label="否" :value="1" />
            </el-radio-group>
          </template>

          <template v-else-if="item.prop === 'parent_id'">
            <el-tree-select
              ref="treeSelectRef"
              v-model="formData[item.prop]"
              :data="treeselect"
              check-strictly
              filterable
              show-checkbox
              node-key="id"
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
import type { TMenuForm } from "@/types/system/menu";

const drawer = ref(false);
const direction = ref<DrawerProps["direction"]>("rtl");
const emits = defineEmits(["getMenuData"]);

const defaultProps = ref({
  children: "children",
  label: "title",
});

const rules = reactive<FormRules>({
  title: [
    { required: true, message: "请输入标题", trigger: "blur" },
    { min: 1, max: 50, message: "长度在 1 到 50 个字符", trigger: "blur" },
  ],
  name: [
    { required: true, message: "请输入名称", trigger: "blur" },
    { min: 1, max: 100, message: "长度在 1 到 100 个字符", trigger: "blur" },
  ],
  path: [
    { required: true, message: "请输入访问路径", trigger: "blur" },
    { min: 1, max: 100, message: "长度在 1 到 100 个字符", trigger: "blur" },
  ],
  icon: [{ required: true, message: "请选择图标", trigger: "blur" }],
  sort: [{ required: true, message: "请输入排序", trigger: "blur" }],
  component: [
    { required: true, message: "请输入组件路径", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
  redirect: [
    { required: false, message: "请输入重定向", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
  activeMenu: [
    { required: false, message: "请输入高亮菜单", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
  parent_id: [{ required: true, message: "请选择上级目录", trigger: "change" }],
});

const fromCol = [
  { prop: "title", label: "菜单标题", placeholder: "菜单标题" },
  { prop: "name", label: "名称", placeholder: "名称" },
  { prop: "sort", label: "排序", placeholder: "排序" },
  { prop: "icon", label: "图标", placeholder: "图标" },
  { prop: "path", label: "路由地址", placeholder: "路由地址" },
  { prop: "component", label: "组件路径", placeholder: "组件路径" },
  { prop: "redirect", label: "重定向", placeholder: "重定向" },
  { prop: "status", label: "禁用", placeholder: "禁用" }, //1 否 2是
  { prop: "hidden", label: "隐藏", placeholder: "隐藏" }, //2否 1 是
  { prop: "noCache", label: "缓存", placeholder: "缓存" }, //2 是 1 否
  { prop: "activeMenu", label: "高亮菜单", placeholder: "高亮菜单" },
  { prop: "parent_id", label: "上级目录", placeholder: "上级目录" },
];

const dialogType = ref("");
const treeselect = ref([]);
//打开
const openDrawer = (row: TMenuForm, type: string, treeselectData: any) => {
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
  formData.value.parent_id = data.id;
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
  formData.value.icon = iconName;
};

const handleIconName = (iconName: string) => {
  formData.value.icon = "";
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
        updateRoleById(data.id, data).then((res) => {
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
