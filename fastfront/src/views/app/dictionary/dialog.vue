<template>
  <el-drawer
    v-model="drawer"
    :direction="direction"
    :size="'50%'"
    @close="drawerClose"
  >
    <template #header>
      <h4>{{ dialogType === "create" ? "新增字典" : "编辑字典" }}</h4>
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
          <template v-if="item.prop === 'des'">
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
import type { TDictionaryForm } from "@/types/app/dictionary";
import {
  createDictionary,
  getDictionaryDetail,
  updateDictionaryById,
} from "@/api/app/dictionary";

const drawer = ref(false);
const direction = ref<DrawerProps["direction"]>("rtl");

const emits = defineEmits(["getDictionaryData"]);

const rules = {
  key: [
    { required: true, message: "请输入配置的Key", trigger: "blur" },
    { min: 1, max: 20, message: "长度在 1 到 20 个字符", trigger: "blur" },
  ],
  value: [
    { required: true, message: "请输入配置的值", trigger: "blur" },
    { min: 1, max: 512, message: "长度在 1 到 512 个字符", trigger: "blur" },
  ],
  des: [
    { required: false, message: "说明", trigger: "blur" },
    { min: 1, max: 512, message: "长度在 1 到 512 个字符", trigger: "blur" },
  ],
};

const fromCol = [
  { prop: "key", label: "配置的Key", placeholder: "配置的Key" },
  { prop: "value", label: "配置的值", placeholder: "配置的值" },
  { prop: "des", label: "说明", placeholder: "说明" },
];

const dialogType = ref("");
//打开
const openDrawer = (type: string, Id?: number) => {
  dialogType.value = type;
  if (type === "update" && Id) {
    getDictionaryDetail(Id).then((res) => {
      formData.value = res.data;
    });
  }
  drawer.value = true;
};

defineExpose({
  openDrawer,
});

const formData = ref<TDictionaryForm>({
  value: "",
  key: "",
  des: "",
});

const formRef = ref();
//关闭
const drawerClose = () => {
  formData.value = { value: "", key: "", des: "" };
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
        createDictionary(data).then((res) => {
          ElMessage.success(res.message);
          formRef.value.resetFields();
          emits("getDictionaryData");
          drawer.value = false;
        });
      } else {
        if (!data.id) return;
        updateDictionaryById(data.id, data).then((res) => {
          ElMessage.success(res.message);
          formRef.value.resetFields();
          emits("getDictionaryData");
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
