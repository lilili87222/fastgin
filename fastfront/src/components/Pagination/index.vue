<template>
  <div :class="{ hidden: hidden }" class="pagination-container">
    <el-pagination
      :background="background"
      v-model:current-page="currentPage"
      v-model:page-size="PageSize"
      :layout="layout"
      :page-sizes="PageSizes"
      :total="total"
      v-bind="$attrs"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :pager-count="5"
    />
  </div>
</template>

<script>
import { defineComponent } from "vue";
import { scrollTo } from "@/utils/scroll-to";

export default defineComponent({
  name: "Pagination",
  props: {
    total: {
      required: true,
      type: Number,
    },
    page: {
      type: Number,
      default: 1,
    },
    limit: {
      type: Number,
      default: 20,
    },
    PageSizes: {
      type: Array,
      default() {
        return [10, 20, 30, 50];
      },
    },
    layout: {
      type: String,
      default: "total, sizes, prev, pager, next, jumper",
    },
    background: {
      type: Boolean,
      default: true,
    },
    autoScroll: {
      type: Boolean,
      default: true,
    },
    hidden: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    currentPage: {
      get() {
        return this.page;
      },
      set(val) {
        this.$emit("update:page", val);
      },
    },
    PageSize: {
      get() {
        return this.limit;
      },
      set(val) {
        this.$emit("update:limit", val);
        // 当页面大小改变时，重置为第一页
        this.$emit("update:page", 1);
      },
    },
  },
  methods: {
    handleSizeChange(val) {
      this.$emit("pagination", { page: 1, limit: val }); //当页面大小改变时，重置为第一页
      if (this.autoScroll) {
        scrollTo(0, 800);
      }
    },
    handleCurrentChange(val) {
      this.$emit("pagination", { page: val, limit: this.PageSize });
      if (this.autoScroll) {
        scrollTo(0, 800);
      }
    },
  },
});
</script>

<style scoped>
.pagination-container {
  display: flex;
  justify-content: flex-end;
  padding: 22px 6px;
  overflow-x: auto; /* 允许在 X 轴上滚动 */
}

.pagination-container.hidden {
  display: none;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .pagination-container {
    padding: 16px 8px; /* 减少内边距 */
  }

  .el-pagination {
    font-size: 12px;
  }

  /* 隐藏不必要的分页组件部分 */
  :deep(.el-pagination__sizes),
  :deep(.el-pagination__total),
  :deep(.el-pagination__jump) {
    display: none !important; /* 使用 !important 确保样式生效 */
  }

  .el-pagination__pager {
    margin: 0 5px;
  }

  .el-pagination .el-pager li {
    padding: 0 8px;
  }

  .el-button {
    padding: 10px; /* 调整按钮的内边距 */
    margin: 0 4px; /* 调整按钮的间距 */
  }
}
</style>
