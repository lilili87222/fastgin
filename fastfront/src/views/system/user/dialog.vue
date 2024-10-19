<template>
  <el-drawer
    v-model="drawer"
    :direction="direction"
    :size="'50%'"
    @close="drawerClose"
  >
    <template #header>
      <h4>{{ dialogType === "create" ? "新增用户" : "编辑用户" }}</h4>
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
          <template v-if="item.prop === 'Status'">
            <el-select
              v-model="formData[item.prop]"
              placeholder="请选择状态"
              style="width: 100%"
            >
              <el-option label="正常" :value="1" />
              <el-option label="禁用" :value="2" />
            </el-select>
          </template>
          <template v-else-if="item.prop === 'RoleIds'">
            <el-select
              v-model="formData[item.prop]"
              multiple
              placeholder="请选择角色"
              style="width: 100%"
            >
              <el-option
                v-for="item in roleList"
                :key="item.Id"
                :label="item.Name"
                :value="item.Id"
              />
            </el-select>
          </template>
          <template v-else-if="item.prop === 'Password'">
            <el-input
              v-model="formData[item.prop]"
              autocomplete="off"
              :type="passwordType"
              :placeholder="dialogType === 'create' ? '新密码' : '重置密码'"
            />
            <span class="show-pwd" @click="showPwd">
              <svg-icon
                :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'"
              />
            </span>
          </template>
          <template v-else-if="item.prop === 'Introduction'">
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
import { computed, reactive, ref } from "vue";
import type { DrawerProps, FormRules } from "element-plus";
import JSEncrypt from "jsencrypt";

import { createUser, updateUserById } from "@/api/system/user";
import type { TUserFormData } from "@/types/system/user";
import type { TRoleTableData } from "@/types/system/role";

import { publicKey } from "./const";

var checkPhone = (rule, value, callback) => {
  if (!value) {
    return callback(new Error("手机号不能为空"));
  } else {
    const reg = /^1[3-9][0-9]\d{8}$/;
    if (reg.test(value)) {
      callback();
    } else {
      return callback(new Error("请输入正确的手机号"));
    }
  }
};
const drawer = ref(false);
const direction = ref<DrawerProps["direction"]>("rtl");
const emits = defineEmits(["getUserData"]);

const rules = reactive<FormRules<TUserFormData>>({
  UserName: [
    { required: true, message: "请输入用户名", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" },
  ],
  Password: [
    { required: true, message: "请输入", trigger: "blur" },
    { min: 6, max: 30, message: "长度在 6 到 30 个字符", trigger: "blur" },
  ],
  RoleIds: [{ required: true, message: "请选择角色", trigger: "change" }],
  NickName: [
    { required: false, message: "请输入昵称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" },
  ],
  Mobile: [{ required: true, validator: checkPhone, trigger: "blur" }],
  Status: [{ required: true, message: "请选择状态", trigger: "change" }],
  Introduction: [
    { required: false, message: "说明", trigger: "blur" },
    { min: 0, max: 100, message: "长度在 0 到 100 个字符", trigger: "blur" },
  ],
});

const formData = ref<TUserFormData>({
  UserName: "",
  Password: "",
  NickName: "",
  Status: 1,
  Mobile: "",
  Avatar: "",
  Introduction: "",
  RoleIds: "",
});

const dialogType = ref("");

const fromCol = computed(() => [
  { prop: "UserName", label: "用户名", placeholder: "用户名" },
  {
    prop: "Password",
    label: dialogType.value === "create" ? "新密码" : "重置密码",
    placeholder: dialogType.value === "create" ? "新密码" : "重置密码",
  },
  { prop: "RoleIds", label: "角色", placeholder: "角色" },
  { prop: "Status", label: "状态", placeholder: "状态" },
  { prop: "NickName", label: "昵称", placeholder: "昵称" },
  { prop: "Mobile", label: "手机号", placeholder: "手机号" },
  { prop: "Introduction", label: "说明", placeholder: "说明" },
]);

const roleList = ref<TRoleTableData[]>([]);
//打开
const openDrawer = (
  row: TUserFormData,
  type: string,
  roleLists: TRoleTableData[]
) => {
  formData.value = row;
  dialogType.value = type;
  roleList.value = roleLists;
  drawer.value = true;
};

defineExpose({
  openDrawer,
});

const passwordType = ref("password");
const showPwd = () => {
  if (passwordType.value === "password") {
    passwordType.value = "";
  } else {
    passwordType.value = "password";
  }
};

//关闭
const drawerClose = () => {
  formData.value = {
    UserName: "",
    Password: "",
    NickName: "",
    Status: 1,
    Mobile: "",
    Avatar: "",
    Introduction: "",
    RoleIds: "",
  };
  formRef.value.resetFields();
};

//取消
const cancelClick = () => {
  drawer.value = false;
};

const formRef = ref();
//提交
const submitForm = () => {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      let data = { ...formData.value };
      if (formData.value.Password !== "") {
        // // 密码RSA加密处理
        // const encryptor = new JSEncrypt();
        // // 设置公钥
        // encryptor.setPublicKey(publicKey);
        // // 加密密码
        // //const encPassword = encryptor.encrypt(formData.value.Password)
        // const encPassword = formData.value.Password;
        // data.password = encPassword;
      }
      if (dialogType.value === "create") {
        createUser(data).then((res) => {
          ElMessage.success(res.Message);
          formRef.value.resetFields();
          emits("getUserData");
          drawer.value = false;
        });
      } else {
        if (!data.Id) return;
        updateUserById(data.Id, data).then((res) => {
          ElMessage.success(res.Message);
          formRef.value.resetFields();
          emits("getUserData");
          drawer.value = false;
        });
      }
    }
  });
};
</script>
<style>
.show-pwd {
  position: absolute;
  right: 10px;
  top: 3px;
  font-size: 16px;
  color: #889aa4;
  cursor: pointer;
  user-select: none;
}
.el-drawer {
  max-width: 90%; /* 设置最大宽度为 90% */
  min-width: 300px; /* 设置最小宽度为 300px */
}
</style>
