<template>
  <el-drawer
    v-model="drawer"
    :direction="direction"
    :size="'50%'"
    @close="drawerClose"
  >
    <template #header>
      <h4>{{ dialogType === "create" ? "新增接口" : "编辑接口" }}</h4>
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
          <template v-if="item.prop === 'method'">
            <el-select
              v-model="formData[item.prop]"
              placeholder="请选择请求方式"
            >
              <el-option label="GET[获取资源]" value="GET" />
              <el-option label="POST[新增资源]" value="POST" />
              <el-option label="PUT[全部更新]" value="PUT" />
              <el-option label="PATCH[增量更新]" value="PATCH" />
              <el-option label="DELETE[删除资源]" value="DELETE" />
            </el-select>
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
import { ref } from "vue";
import { ElMessage } from "element-plus";
import type { DrawerProps } from "element-plus";
import { createApi, updateApiById } from "@/api/system/api";
import type { TApiForm } from "@/types/system/api";

const drawer = ref(false);
const direction = ref<DrawerProps["direction"]>("rtl");

const emits = defineEmits(["getApiData"]);

const rules = {
  path: [
    { required: true, message: "请输入访问路径", trigger: "blur" },
    { min: 1, max: 100, message: "长度在 1 到 100 个字符", trigger: "blur" },
  ],
  category: [
    { required: true, message: "请输入所属类别", trigger: "blur" },
    { min: 1, max: 50, message: "长度在 1 到 50 个字符", trigger: "blur" },
  ],
  method: [{ required: true, message: "请选择请求方式", trigger: "change" }],
  des: [
    { required: false, message: "说明", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
};

const fromCol = [
  { prop: "path", label: "访问路径", placeholder: "访问路径" },
  { prop: "category", label: "所属类别", placeholder: "所属类别" },
  { prop: "method", label: "请求方式", placeholder: "请求方式" },
  { prop: "des", label: "说明", placeholder: "说明" },
];

const dialogType = ref("");
//打开
const openDrawer = (row: TApiForm, type: string) => {
  formData.value = row;
  dialogType.value = type;
  drawer.value = true;
};

defineExpose({
  openDrawer,
});

const formData = ref<TApiForm>({
  path: "",
  category: "",
  method: "",
  des: "",
});

const formRef = ref();
//关闭
const drawerClose = () => {
  formData.value = {
    path: "",
    category: "",
    method: "",
    des: "",
  };
  formRef.value.resetFields();
};
//取消
const cancelClick = () => {
  drawer.value = false;
};

//提交
const submitForm = () => {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      let data = { ...formData.value };
      if (dialogType.value === "create") {
        createApi(data).then((res) => {
          ElMessage.success(res.message);
          formRef.value.resetFields();
          emits("getApiData");
          drawer.value = false;
        });
      } else {
        if (!data.id) return;
        updateApiById(data.id, data).then((res) => {
          ElMessage.success(res.message);
          formRef.value.resetFields();
          emits("getApiData");
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
