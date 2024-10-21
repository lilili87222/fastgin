<template>
  <div class="app-container">
    <el-card shadow="always">
      <SearchForm
        :searchColumn="searchColumn"
        :searchAction="searchAction"
        @onClear="onClear"
        @onDelete="onDelete"
        @onAdd="onAdd"
        @onSearch="onSearch"
      ></SearchForm>
      <el-table
        border
        v-loading="loading"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        stripe
        :data="tableData"
      >
        <el-table-column type="selection" min-width="55" align="center" />
        <el-table-column
          v-for="item in tableColumn"
          :prop="item.prop"
          sortable
          :min-width="item.minWidth"
          :label="item.label"
          align="center"
        >
          <template #default="scope">
            <template v-if="item.prop === 'status'">
              <el-tag
                size="small"
                :type="scope.row.status === 1 ? 'success' : 'danger'"
                disable-transitions
                >{{ scope.row.status === 1 ? "正常" : "禁用" }}</el-tag
              >
            </template>

            <template v-else>
              {{ scope.row[item.prop] }}
            </template>
          </template>
        </el-table-column>
        <el-table-column
          fixed="right"
          label="操作"
          align="center"
          min-width="150"
        >
          <template #default="scope">
            <el-button
              @click="update(scope.row)"
              type="primary"
              class="custom-btn"
              >编辑</el-button
            >
            <el-popconfirm
              title="确定删除吗？"
              @confirm="singleDelete(scope.row.id)"
            >
              <template #reference>
                <el-button type="danger" class="custom-btn">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <Pagination
        v-show="total > 0"
        :total="total"
        v-model:page="params.page_num"
        v-model:limit="params.page_size"
        @pagination="onPaginaion"
      ></Pagination>
      <Dialog ref="DrawerRef" @getUserData="getUserData"></Dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { batchDeleteUserByIds, getUsers } from "@/api/system/user";
import { getRoles } from "@/api/system/role";
import SearchForm from "@/components/SearchForm/index.vue";
import Pagination from "@/components/Pagination/index.vue";
import type { TUserQuery, TUserTable } from "@/types/system/user";
import type { TRoleTable } from "@/types/system/role";

import Dialog from "./dialog.vue";

const searchColumn = [
  { prop: "user_name", label: "用户名", placeholder: "用户名" },
  { prop: "nick_name", label: "昵称", placeholder: "昵称" },
  {
    prop: "status",
    label: "状态",
    placeholder: "状态",
    type: "select",
    options: [
      { label: "正常", value: 1 },
      { label: "禁用", value: 2 },
    ],
  },
  { prop: "mobile", label: "手机号", placeholder: "手机号" },
];

const tableColumn = [
  { prop: "user_name", label: "用户名", minWidth: 95 },
  { prop: "nick_name", label: "昵称", minWidth: 80 },
  { prop: "status", label: "状态", minWidth: 80 },
  { prop: "mobile", label: "手机号", minWidth: 95 },
  { prop: "creator", label: "创建人", minWidth: 95 },
  { prop: "des", label: "说明", minWidth: 80 },
];

// 查询参数
const params = ref<TUserQuery>({
  page_num: 1,
  page_size: 10,
});
// 表格数据
const tableData = ref<TUserTable[]>([]);
const total = ref(0);
const loading = ref(false);

onMounted(() => {
  getTableData();
  getRole();
});

const roleList = ref<TRoleTable[]>([]);
const getRole = () => {
  getRoles().then((res) => {
    roleList.value = res.data.data;
  });
};

// 获取表格数据
const getTableData = () => {
  loading.value = true;
  getUsers(params.value)
    .then((res) => {
      tableData.value = res.data.data;
      total.value = res.data.total;
    })
    .finally(() => {
      loading.value = false;
    });
};

// 表格多选
const multipleSelection = ref<TUserTable[]>([]);

const searchAction = computed(() => [
  { label: "查询", event: "search", type: "primary" },
  { label: "新增", event: "add", type: "warning" },
  {
    label: "批量删除",
    event: "delete",
    type: "danger",
    disable: multipleSelection.value.length === 0,
  },
]);

//分页
const onPaginaion = (val: any) => {
  params.value.page_num = val.page;
  params.value.page_size = val.limit;
  getTableData();
};

const getUserData = () => {
  getTableData();
};

// 搜索
const onSearch = (form: TUserQuery) => {
  params.value = form;
  params.value.page_num = 1;
  params.value.page_size = 10;
  getTableData();
};

//批量删除
const onDelete = () => {
  ElMessageBox.confirm("此操作将永久删除, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      loading.value = true;
      const Ids: number[] = [];
      multipleSelection.value.forEach((x: any) => {
        Ids.push(x.id);
      });
      batchDeleteUserByIds({
        Ids,
      })
        .then((res) => {
          getTableData();
          ElMessage.success(res.message);
        })
        .finally(() => {
          loading.value = false;
        });
    })
    .catch(() => {
      ElMessage.info("已取消删除");
    });
};

// 表格多选
const handleSelectionChange = (val: TUserTable[]) => {
  multipleSelection.value = val;
};

const DrawerRef = ref();
// 新增弹窗
const onAdd = () => {
  DrawerRef.value.openDrawer({}, "create", roleList.value);
};

// 编辑弹窗
const update = (row: TUserTable) => {
  DrawerRef.value.openDrawer({ ...row }, "update", roleList.value);
};

//清空
const onClear = (form: TUserQuery) => {
  params.value = form;
  params.value.page_num = 1;
  params.value.page_size = 10;
  getTableData();
};

// 单个删除
const singleDelete = (Id: number) => {
  loading.value = true;
  batchDeleteUserByIds({
    Ids: [Id],
  }).then((res) => {
    getTableData();
    ElMessage.success(res.message);
  });
};
</script>

<style scoped>
.delete-popover {
  margin-left: 10px;
}
</style>
