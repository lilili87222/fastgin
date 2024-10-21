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
            <template v-if="item.prop === 'Method'"> </template>
            <template
              v-else-if="['CreatedAt', 'UpdatedAt'].includes(item.prop)"
            >
              {{ parseGoTime(scope.row[item.prop]) }}
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
              @click="update(scope.row.ID)"
              type="primary"
              class="custom-btn"
              >编辑</el-button
            >
            <el-popconfirm
              title="确定删除吗？"
              @confirm="singleDelete(scope.row.ID)"
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
      <Dialog ref="DrawerRef" @getDictData="getDictData"></Dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import {
  batchDeleteDictionaryById,
  batchDeleteDictionaryByIds,
  getDictionary,
} from "@/api/app/dictionary";
import SearchForm from "@/components/SearchForm/index.vue";
import Pagination from "@/components/Pagination/index.vue";
import type {
  TDictionaryTableData,
  TDictionaryQuery,
} from "@/types/app/dictionary";
import { parseGoTime } from "@/utils/index";

import Dialog from "./dialog.vue";

const searchColumn = [
  { prop: "Value", label: "字典名称", placeholder: "字典名称" },
  { prop: "Key", label: "字典类型", placeholder: "字典类型" },
  { prop: "Desc", label: "说明", placeholder: "说明" },
];

const tableColumn = [
  { prop: "Value", label: "字典名称", minWidth: 110 },
  { prop: "Key", label: "字典类型", minWidth: 110 },
  { prop: "Desc", label: "说明", minWidth: 85 },
  { prop: "CreatedAt", label: "创建时间", minWidth: 110 },
  { prop: "UpdatedAt", label: "更新时间", minWidth: 110 },
];

// 查询参数
const params = ref<TDictionaryQuery>({
  PageNum: 1,
  PageSize: 10,
});
// 表格数据
const tableData = ref<TDictionaryTableData[]>([]);
const total = ref(0);
const loading = ref(false);

onMounted(() => {
  getTableData();
});

// 获取表格数据
const getTableData = () => {
  loading.value = true;
  getDictionary(params.value)
    .then((res) => {
      tableData.value = res.data.Data;
      total.value = res.data.Total;
    })
    .finally(() => {
      loading.value = false;
    });
};

const DrawerRef = ref();
//新增
const onAdd = () => {
  DrawerRef.value.openDrawer("create");
};
//编辑
const update = (Id: number) => {
  DrawerRef.value.openDrawer("update", Id);
};

//清空
const onClear = (form: TDictionaryQuery) => {
  params.value = form;
  params.value.PageNum = 1;
  params.value.PageSize = 10;
  getTableData();
};

const getDictData = () => {
  getTableData();
};

// 表格多选
const multipleSelection = ref<TDictionaryTableData[]>([]);
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
  params.value.PageNum = val.page;
  params.value.PageSize = val.limit;
  getTableData();
};

//搜索
const onSearch = (form: TDictionaryQuery) => {
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
    .then(async (res) => {
      loading.value = true;
      const Ids: number[] = [];
      multipleSelection.value.forEach((x: any) => {
        Ids.push(x.ID);
      });
      batchDeleteDictionaryByIds({
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
const handleSelectionChange = (val: TDictionaryTableData[]) => {
  multipleSelection.value = val;
};

// 单个删除
const singleDelete = (Id) => {
  loading.value = true;
  batchDeleteDictionaryById(Id)
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
