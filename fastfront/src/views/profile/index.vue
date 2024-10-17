<template>
  <div class="app-container">
    <div v-if="user">
      <el-row :gutter="20">
        <el-col :span="6" :xs="24">
          <user-card :user="user" />
        </el-col>

        <el-col :span="18" :xs="24">
          <el-card>
            <el-tabs v-model="activeTab">
              <el-tab-pane label="账户" name="account">
                <account :user="user" />
              </el-tab-pane>
            </el-tabs>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import UserCard from "./components/UserCard.vue";
import Account from "./components/Account.vue";
import { onMounted, ref } from "vue";
import store from "@/store";

const user = ref({});
const activeTab = ref("account");

onMounted(() => {
  getUser();
});

const getUser = () => {
  user.value = {
    name: store.user().name,
    role: store
      .user()
      .roles.map((item: any) => item.Name)
      .join(" | "),
    email: "admin@test.com",
    avatar: store.user().avatar,
  };
};
</script>
