<template>
  <div class="login-container">
    <div class="lefts" style="flex: 1">
      <div class="top-row">
        <img
          class="logo"
          src="https://appuploader.net/statics/images/appuploader_logo.png"
          alt="logo"
        />
        <span>Admin</span>
      </div>
      <div class="bg-container">
        <img class="bg" alt="bg" src="../../assets/login/login-box-bg.svg" />
        <h3>后台管理系统</h3>
      </div>
    </div>
    <div class="rights">
      <LoginForm
        v-if="showForm === 'login'"
        @changeShowForm="changeShowForm"
        key="login"
      ></LoginForm>

      <SignForm
        v-if="showForm === 'sign'"
        @changeShowForm="changeShowForm"
        key="register"
      ></SignForm>
      <ForgetForm
        v-if="showForm === 'forget'"
        @changeShowForm="changeShowForm"
        key="forget"
      ></ForgetForm>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import LoginForm from "./components/LoginForm.vue";
import SignForm from "./components/SignForm.vue";
import ForgetForm from "./components/ForgetForm.vue";

const showForm = ref("login");

const changeShowForm = (value: string) => {
  showForm.value = value;
};
</script>

<style lang="scss">
@import "./login.scss";
</style>
