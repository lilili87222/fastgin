<template>
  <el-drawer
    v-model="drawer"
    :direction="direction"
    :size="'50%'"
    @close="drawerClose"
  >
    <template #header>
      <h4>{{ dialogType === "create" ? "新增角色" : "编辑角色" }}</h4>
    </template>
    <template #default>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item
          v-for="item in fromCol"
          :key="item.prop"
          :label="item.label"
          :prop="item.prop"
        >
          <template v-if="item.prop === 'status'">
            <el-select
              v-model="formData[item.prop]"
              placeholder="请选择状态"
              style="width: 100%"
            >
              <el-option label="正常" :value="1" />
              <el-option label="禁用" :value="2" />
            </el-select>
          </template>

          <template v-else-if="item.prop === 'sort'">
            <el-input-number
              v-model="formData[item.prop]"
              controls-position="right"
              :min="1"
              :max="999"
            />
          </template>
          <template v-else-if="item.prop === 'des'">
            <el-input
              v-model="formData[item.prop]"
              type="textarea"
              placeholder="说明"
              show-word-limit
              maxlength="100"
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
import type { TRoleForm } from "@/types/system/role";

const drawer = ref(false);
const direction = ref<DrawerProps["direction"]>("rtl");
const emits = defineEmits(["getRoleData"]);

const rules = reactive<FormRules<TRoleForm>>({
  name: [
    { required: true, message: "请输入角色名称", trigger: "blur" },
    { min: 1, max: 20, message: "长度在 1 到 20 个字符", trigger: "blur" },
  ],
  keyword: [
    { required: true, message: "请输入关键字", trigger: "blur" },
    { min: 1, max: 20, message: "长度在 1 到 20 个字符", trigger: "blur" },
  ],
  status: [{ required: true, message: "请选择角色状态", trigger: "change" }],
  des: [
    { required: false, message: "说明", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
});

const formData = ref<TRoleForm>({
  name: "",
  keyword: "",
  status: 1,
  sort: 999,
  des: "",
});

const fromCol = [
  { prop: "name", label: "角色名称", placeholder: "角色名称" },

  { prop: "keyword", label: "关键字", placeholder: "关键字" },
  {
    prop: "status",
    label: "角色状态",
    placeholder: "角色状态",
  },
  {
    prop: "sort",
    label: "等级(1最高)",
    placeholder: "等级(1最高)",
  },
  {
    prop: "des",
    label: "说明",
    placeholder: "说明",
  },
];

const dialogType = ref("");
//打开
const openDrawer = (row: TRoleForm, type: string) => {
  formData.value = row;
  dialogType.value = type;
  drawer.value = true;
};

defineExpose({
  openDrawer,
});

//关闭
const drawerClose = () => {
  formData.value = {
    name: "",
    keyword: "",
    status: 1,
    sort: 999,
    des: "",
  };
  formRef.value.resetFields();
};

//取消按钮
const cancelClick = () => {
  drawer.value = false;
};

const formRef = ref();
const submitForm = () => {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      let data = { ...formData.value };
      if (dialogType.value === "create") {
        createRole(data).then((res) => {
          ElMessage.success(res.message);
          formRef.value.resetFields();
          emits("getRoleData");
          drawer.value = false;
        });
      } else {
        if (!data.id) return;
        updateRoleById(data.id, data).then((res) => {
          ElMessage.success(res.message);
          formRef.value.resetFields();
          emits("getRoleData");
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
