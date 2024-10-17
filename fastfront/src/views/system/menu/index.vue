<template>
  <div class="app-container">
    <el-card shadow="always">
      <SearchForm
        :searchColumn="searchColumn"
        :searchAction="searchAction"
        @onAdd="onAdd"
        @onDelete="onDelete"
        @onSearch="onSearch"
      ></SearchForm>
      <el-table
        border
        :tree-props="{ children: 'Children' }"
        row-key="Id"
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
          :min-width="item.minWidth"
          :label="item.label"
          align="center"
        >
          <template #default="scope">
            <template v-if="item.prop === 'Status'">
              <el-tag
                size="small"
                :type="scope.row.Status === 1 ? 'success' : 'danger'"
                >{{ scope.row.Status === 1 ? "否" : "是" }}</el-tag
              >
            </template>
            <template v-else-if="item.prop === 'Hidden'">
              <el-tag
                size="small"
                :type="scope.row.Hidden === 1 ? 'danger' : 'success'"
                >{{ scope.row.Hidden === 1 ? "是" : "否" }}</el-tag
              >
            </template>
            <template v-else-if="item.prop === 'NoCache'">
              <el-tag
                size="small"
                :type="scope.row.NoCache === 1 ? 'danger' : 'success'"
                >{{ scope.row.NoCache === 1 ? "否" : "是" }}</el-tag
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
              @confirm="singleDelete(scope.row.Id)"
            >
              <template #reference>
                <el-button type="danger" class="custom-btn">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <Dialog ref="DrawerRef" @getMenuData="getMenuData"></Dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import { batchDeleteMenuByIds, getMenuTree } from "@/api/system/menu";
import SearchForm from "@/components/SearchForm/index.vue";
import type { TMenuQuery, TMenuTableData } from "@/types/system/menu";

import Dialog from "./dialog.vue";

const searchColumn = [];

// 表格多选
const multipleSelection = ref<TMenuTableData[]>([]);

const searchAction = computed(() => [
  { label: "新增", event: "add", type: "warning" },
  {
    label: "批量删除",
    event: "delete",
    type: "danger",
    disable: multipleSelection.value.length === 0,
  },
]);

const getMenuData = () => {
  getTableData();
};

const tableColumn = [
  { prop: "Title", label: "菜单标题", minWidth: 105 },
  { prop: "Name", label: "名称", minWidth: 80 },
  { prop: "Icon", label: "图标", minWidth: 80 },
  { prop: "Path", label: "路由地址", minWidth: 105 },
  { prop: "Component", label: "组件路径", minWidth: 105 },
  { prop: "Redirect", label: "重定向", minWidth: 105 },
  { prop: "Sort", label: "排序", minWidth: 80 },
  { prop: "Status", label: "禁用", minWidth: 80 },
  { prop: "Hidden", label: "隐藏", minWidth: 80 },
  { prop: "NoCache", label: "缓存", minWidth: 80 },
  { prop: "ActiveMenu", label: "高亮菜单", minWidth: 105 },
];

// 查询参数
const params = ref<TMenuQuery>({
  PageNum: 1,
  PageSize: 10,
});
// 表格数据
const tableData = ref<TMenuTableData[]>([]);
const total = ref(0);
const loading = ref(false);

const treeselectData = ref<any>([]);
onMounted(() => {
  getTableData();
});

// 获取表格数据
const getTableData = () => {
  loading.value = true;
  getMenuTree()
    .then((res) => {
      const { Data } = res;
      tableData.value = Data;
      treeselectData.value = [{ Id: 0, Title: "顶级类目", Children: Data }];
      total.value = Data.Total;
    })
    .finally(() => {
      loading.value = false;
    });
};

//搜索
const onSearch = (form: TMenuQuery) => {
  params.value = form;
  params.value.PageNum = 1;
  params.value.PageSize = 10;
  getTableData();
};

const DrawerRef = ref();
// 新增
const onAdd = () => {
  DrawerRef.value.openDrawer({}, "create", treeselectData.value);
};
// 编辑
const update = (row: TMenuTableData) => {
  DrawerRef.value.openDrawer({ ...row }, "update", treeselectData.value);
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
        Ids.push(x.Id);
      });
      batchDeleteMenuByIds({
        Ids,
      })
        .then((res) => {
          getTableData();
          ElMessage.success(res.Message);
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
const handleSelectionChange = (val: TMenuTableData[]) => {
  multipleSelection.value = val;
};

// 单个删除
const singleDelete = (Id: number) => {
  loading.value = true;
  batchDeleteMenuByIds({
    Ids: [Id],
  })
    .then((res) => {
      getTableData();

      ElMessage.success(res.Message);
    })
    .finally(() => (loading.value = false));
};
</script>

<style scoped>
.delete-popover {
  margin-left: 10px;
}
</style>
