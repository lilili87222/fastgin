<template>
  <div class="app-container">
    <el-card shadow="always">
      <SearchForm
        :searchColumn="searchColumn"
        :searchAction="searchAction"
        @onAdd="onAdd"
        @onClear="onClear"
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
          :min-width="item.minWidth"
          sortable
          :label="item.label"
          align="center"
        >
          <template #default="scope">
            <template v-if="item.prop === 'method'">
              <el-tag
                :type="methodsTagFilter(scope.row.method)"
                disable-transitions
                >{{ scope.row.method }}
              </el-tag>
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
          min-width="120"
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
      <Dialog ref="DrawerRef" @getApiData="getApiData"></Dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import { batchDeleteApiByIds, getApis } from "@/api/system/api";
import SearchForm from "@/components/SearchForm/index.vue";
import Pagination from "@/components/Pagination/index.vue";
import type { TApiQuery, TApiTable } from "@/types/system/api";

import Dialog from "./dialog.vue";

const searchColumn = [
  { prop: "path", label: "访问路径", placeholder: "访问路径" },
  { prop: "category", label: "所属类别", placeholder: "所属类别" },
  {
    prop: "method",
    label: "请求方式",
    placeholder: "请求方式",
    type: "select",
    options: [
      { label: "GET[获取资源]", value: "GET" },
      { label: "POST[新增资源]", value: "POST" },
      { label: "PUT[全部更新]", value: "PUT" },
      { label: "PATCH[增量更新]", value: "PATCH" },
      { label: "DELETE[删除资源]", value: "DELETE" },
    ],
  },
  { prop: "creator", label: "创建人", placeholder: "创建人" },
];

const tableColumn = [
  { prop: "path", label: "访问路径", minWidth: 115 },
  { prop: "category", label: "所属类别", minWidth: 115 },
  { prop: "method", label: "请求方式", minWidth: 115 },
  { prop: "creator", label: "创建人", minWidth: 95 },
  { prop: "des", label: "说明", minWidth: 110 },
];

// 查询参数
const params = ref<TApiQuery>({
  page_num: 1,
  page_size: 10,
});
// 表格数据
const tableData = ref<TApiTable[]>([]);
const total = ref(0);
const loading = ref(false);

onMounted(() => {
  getTableData();
});

// 获取表格数据
const getTableData = () => {
  loading.value = true;
  getApis(params.value)
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
const update = (row: TApiTable) => {
  DrawerRef.value.openDrawer({ ...row }, "update");
};

//清空
const onClear = (form: TApiQuery) => {
  params.value = form;
  params.value.page_num = 1;
  params.value.page_size = 10;
  getTableData();
};

const getApiData = () => {
  getTableData();
};

// 表格多选
const multipleSelection = ref<TApiTable[]>([]);
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

//搜索
const onSearch = (form: TApiQuery) => {
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
    .then(async (res) => {
      loading.value = true;
      const Ids: number[] = [];
      multipleSelection.value.forEach((x: any) => {
        Ids.push(x.id);
      });
      batchDeleteApiByIds({
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

const methodsTagFilter = (val: string) => {
  if (val === "GET") {
    return "primary";
  } else if (val === "POST") {
    return "success";
  } else if (val === "PUT") {
    return "info";
  } else if (val === "PATCH") {
    return "warning";
  } else if (val === "DELETE") {
    return "danger";
  } else {
    return "info";
  }
};

// 表格多选
const handleSelectionChange = (val: TApiTable[]) => {
  multipleSelection.value = val;
};

// 单个删除
const singleDelete = (Id) => {
  loading.value = true;
  batchDeleteApiByIds({
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
