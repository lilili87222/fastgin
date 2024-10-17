<template>
  <div>
    <el-form ref="formRef" :inline="true" :model="form">
      <el-form-item
        v-for="item in searchColumn"
        :key="item.prop"
        :label="item.label"
        :prop="item.prop"
      >
        <template v-if="item.type === 'select'">
          <el-select
            style="width: 200px"
            @clear="clear(item.prop)"
            clearable
            v-model="form[item.prop]"
            :placeholder="item.placeholder"
          >
            <el-option
              v-for="option in item.options"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            ></el-option>
          </el-select>
        </template>
        <template v-else>
          <el-input
            clearable
            @clear="clear(item.prop)"
            v-model="form[item.prop]"
            :placeholder="item.placeholder"
          ></el-input>
        </template>
      </el-form-item>

      <el-form-item>
        <el-button
          v-for="action in props.searchAction"
          :key="action.event"
          :disabled="action.disable"
          class="custom-btn"
          :type="action.type || 'primary'"
          @click="handleAction(action.event)"
        >
          {{ action.label }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

// 定义 props 的类型
const props = defineProps({
  searchColumn: {
    type: Array,
    default: () => [],
  } as any,
  searchAction: {
    type: Array,
    default: () => [],
  } as any,
});

const emits = defineEmits([
  "onSearch",
  "onClear",
  "onDelete",
  "onAdd",
  "onRestart",
  "onStop",
]);

const form = ref({});
const clear = (i) => {
  form.value[i] = null;
  emits("onClear", form.value);
};

const formRef = ref(null);

const handleAction = (event) => {
  if (event === "search") {
    emits("onSearch", form.value);
  } else if (event === "delete") {
    emits("onDelete", form.value);
  } else if (event === "add") {
    emits("onAdd", form.value);
  } else if (event === "restart") {
    emits("onRestart", form.value);
  } else if (event === "stop") {
    emits("onStop", form.value);
  }
};
</script>
<style scoped>
:deep(.el-input.el-input--default.el-input--suffix) {
  width: 200px !important;
}
</style>
