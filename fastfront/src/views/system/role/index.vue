<template>
  <div class="app-container">
    <el-card shadow="always">
      <SearchForm
        :searchColumn="searchColumn"
        @onClear="onClear"
        :searchAction="searchAction"
        @onAdd="onAdd"
        @onDelete="onDelete"
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
          min-width="130"
        >
          <template #default="scope">
            <el-button
              @click="update(scope.row)"
              class="custom-btn"
              type="primary"
              >编辑</el-button
            >
            <el-button
              @click="openPermission(scope.row)"
              class="custom-btn"
              type="warning"
              >权限</el-button
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
      <Dialog ref="DrawerRef" @getRoleData="getRoleData"></Dialog>
      <Permission ref="PermissionRef" @getRoleData="getRoleData"></Permission>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { batchDeleteRoleByIds, getRoles } from "@/api/system/role";
import SearchForm from "@/components/SearchForm/index.vue";
import Pagination from "@/components/Pagination/index.vue";
import type { TRoleQuery, TRoleTable } from "@/types/system/role";

import Dialog from "./dialog.vue";
import Permission from "./permission.vue";

const searchColumn = [
  { prop: "name", label: "角色名称", placeholder: "用户名" },
  { prop: "keyword", label: "关键字", placeholder: "昵称" },
  {
    prop: "status",
    label: "角色状态",
    placeholder: "角色状态",
    type: "select",
    options: [
      { label: "正常", value: 1 },
      { label: "禁用", value: 2 },
    ],
  },
];

// 查询参数
const params = ref<TRoleQuery>({
  page_num: 1,
  page_size: 10,
});

// 表格数据
const tableData = ref<TRoleTable[]>([]);
const total = ref(0);
const loading = ref(false);

onMounted(() => {
  getTableData();
});

// 获取表格数据
const getTableData = () => {
  loading.value = true;
  getRoles(params.value)
    .then((res) => {
      tableData.value = res.data.data;
      total.value = res.data.total;
    })
    .finally(() => {
      loading.value = false;
    });
};

const DrawerRef = ref();
//新增
const onAdd = () => {
  DrawerRef.value.openDrawer({}, "create");
};

//编辑
const update = (row: TRoleTable) => {
  DrawerRef.value.openDrawer({ ...row }, "update");
};

const PermissionRef = ref();
//权限
const openPermission = (row: TRoleTable) => {
  PermissionRef.value.openDrawer({ ...row });
};

//清空
const onClear = (form: TRoleQuery) => {
  params.value = form;
  params.value.page_num = 1;
  params.value.page_size = 10;
  getTableData();
};

const getRoleData = () => {
  getTableData();
};

// 表格多选
const multipleSelection = ref<TRoleTable[]>([]);

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

const tableColumn = [
  { prop: "name", label: "角色名称", minWidth: 105 },
  { prop: "keyword", label: "关键字", minWidth: 95 },
  { prop: "sort", label: "等级", minWidth: 80 },
  { prop: "status", label: "角色状态", minWidth: 105 },
  { prop: "creator", label: "创建人", minWidth: 95 },
  { prop: "des", label: "说明", minWidth: 80 },
];

//搜索
const onSearch = (form: TRoleQuery) => {
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
      batchDeleteRoleByIds({
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
const handleSelectionChange = (val: TRoleTable[]) => {
  multipleSelection.value = val;
};

// 单个删除
const singleDelete = (Id: number) => {
  loading.value = true;
  batchDeleteRoleByIds({
    Ids: [Id],
  })
    .then((res) => {
      getTableData();
      ElMessage.success(res.message);
    })
    .finally(() => (loading.value = false));
};
</script>

<style scoped>
.delete-popover {
  margin-left: 10px;
}
</style>
