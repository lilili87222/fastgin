<template>
  <el-form :model="form" label-width="120px" :rules="rules" ref="formRef">
    <el-form-item label="用户名">
      {{ userName }}
    </el-form-item>
    <el-form-item label="旧密码" prop="OldPassword">
      <el-input
        v-model.trim="form.OldPassword"
        type="password"
        placeholder="请输入旧密码"
      />
    </el-form-item>
    <el-form-item label="新密码" prop="newPassword">
      <el-input
        v-model.trim="form.newPassword"
        type="password"
        placeholder="请输入新密码"
      />
    </el-form-item>
    <el-form-item label="重复密码" prop="confirmPassword">
      <el-input
        v-model.trim="form.confirmPassword"
        type="password"
        placeholder="请再次输入密码"
      />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submitForm">修改密码</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import store from "@/store";
import { changePwd } from "@/api/system/user";

const formRef = ref();
const userName = ref(store.user().name); // 模拟获取用户名
const form = ref({
  OldPassword: "",
  confirmPassword: "",
  newPassword: "",
});

const rules = ref({
  OldPassword: [
    { required: true, message: "请输入旧密码", trigger: "blur" },
    { min: 6, message: "密码不能少于6位", trigger: "blur" },
  ],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, message: "密码不能少于6位", trigger: "blur" },
  ],
  confirmPassword: [
    { required: true, message: "请再次输入密码", trigger: "blur" },
    {
      validator: (rule, value, callback) => {
        if (value === "") {
          callback(new Error("请再次输入密码"));
        } else if (value !== form.value.newPassword) {
          callback(new Error("两次输入的密码不一致"));
        } else {
          callback();
        }
      },
      trigger: "blur",
    },
  ],
});

const submitForm = () => {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      changePwd(form.value).then((res) => {
        console.log(res);
        ElMessage({
          message: "用户信息已成功更新",
          type: "success",
          duration: 5000,
        });
        setTimeout(async () => {
          await store.user().logout();
          location.reload(); // 为了重新实例化vue-router对象 避免bug
        }, 1500);
      });
    }
  });
};
</script>
