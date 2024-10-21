<template>
  <div class="app-container">
    <el-card class="container-card" shadow="always">
      <SearchForm
        :searchColumn="searchColumn"
        :searchAction="searchAction"
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
            <template v-if="item.prop === 'Status'">
              <el-tag
                :type="statusTagFilter(scope.row.Status)"
                disable-transitions
                >{{ scope.row.Status }}</el-tag
              >
            </template>
            <template v-else-if="item.prop === 'TimeCost'">
              <el-tag
                :type="timeCostTagFilter(scope.row.TimeCost)"
                disable-transitions
                >{{ scope.row.TimeCost }}</el-tag
              >
            </template>
            <template v-else-if="item.prop === 'StartTime'">
              {{ parseGoTime(scope.row.StartTime) }}
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
          min-width="100"
        >
          <template #default="scope">
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

      <Pagination
        v-show="total > 0"
        :total="total"
        v-model:page="params.PageNum"
        v-model:limit="params.PageSize"
        @pagination="onPaginaion"
      ></Pagination>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import {
  getOperationLogs,
  batchDeleteOperationLogByIds,
} from "@/api/log/operationLog.js";
import SearchForm from "@/components/SearchForm/index.vue";
import Pagination from "@/components/Pagination/index.vue";
import type { TLogs, TLogsQuery } from "@/types/log/operation-log";
import { parseGoTime } from "@/utils/index";

const searchColumn = [
  { prop: "user_name", label: "请求人", placeholder: "请求人" },
  { prop: "ip", label: "IP地址", placeholder: "IP地址" },
  { prop: "path", label: "请求路径", placeholder: "请求路径" },
  { prop: "status", label: "请求状态", placeholder: "请求状态" },
];

const tableColumn = [
  { prop: "UserName", label: "请求人", minWidth: 95 },
  { prop: "Ip", label: "IP地址", minWidth: 105 },
  { prop: "Path", label: "请求路径", minWidth: 105 },
  { prop: "Status", label: "请求状态", minWidth: 105 },
  { prop: "StartTime", label: "发起时间", minWidth: 105 },
  { prop: "TimeCost", label: "请求耗时", minWidth: 105 },
  { prop: "Desc", label: "说明", minWidth: 80 },
];

// 查询参数
const params = ref<TLogsQuery>({
  PageNum: 1,
  PageSize: 10,
});

// 表格数据
const tableData = ref<TLogs[]>([]);
const total = ref(0);
const loading = ref(false);

onMounted(() => {
  getTableData();
});

// 获取表格数据
const getTableData = () => {
  loading.value = true;
  getOperationLogs(params.value)
    .then((res) => {
      const { data: Data } = res;
      tableData.value = Data.Data;
      total.value = Data.Total;
    })
    .finally(() => {
      loading.value = false;
    });
};

// 表格多选
const multipleSelection = ref<TLogs[]>([]);

const searchAction = computed(() => [
  { label: "查询", event: "search", type: "primary" },
  {
    label: "批量删除",
    event: "delete",
    type: "danger",
    disable: multipleSelection.value.length === 0,
  },
]);

//切换页面
const onPaginaion = (val: any) => {
  params.value.PageNum = val.page;
  params.value.PageSize = val.limit;
  getTableData();
};

//清空
const onClear = (form: TLogsQuery) => {
  params.value = form;
  params.value.PageNum = 1;
  params.value.PageSize = 10;
  getTableData();
};

//搜索
const onSearch = (form: TLogsQuery) => {
  params.value = form;
  params.value.PageNum = 1;
  params.value.PageSize = 10;
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
        Ids.push(x.Id);
      });
      batchDeleteOperationLogByIds({
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

const statusTagFilter = (val: number) => {
  switch (val) {
    case 200:
      return "success";
    case 400:
    case 401:
    case 403:
    case 500:
      return "danger";
    default:
      return "info";
  }
};
const timeCostTagFilter = (val: number) => {
  if (val <= 200) {
    return "success";
  } else if (val > 200 && val <= 1000) {
    return "";
  } else if (val > 1000 && val <= 2000) {
    return "warning";
  } else {
    return "danger";
  }
};

// 表格多选
const handleSelectionChange = (val: TLogs[]) => {
  multipleSelection.value = val;
};

// 单个删除
const singleDelete = (Id: number) => {
  loading.value = true;
  batchDeleteOperationLogByIds({
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
.container-card {
  margin: 10px;
}

.delete-popover {
  margin-left: 10px;
}
</style>
