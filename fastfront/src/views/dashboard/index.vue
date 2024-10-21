<template>
  <div class="app-container">
    <SearchForm
      :searchColumn="searchColumn"
      :searchAction="searchAction"
      @onRestart="onRestart"
      @onStop="onStop"
    ></SearchForm>
    <div>
      <el-row :gutter="20">
        <template v-for="item in cards" :key="item.key">
          <el-col :xs="24" :sm="12">
            <el-card class="card-item">
              <h3>{{ item.title }}</h3>
              <el-divider />
              <div style="height: 23vh; overflow-y: auto; overflow-x: hidden">
                <!-- RunTime -->
                <template v-if="item.key === 'os'">
                  <el-row
                    :gutter="20"
                    v-for="(value, key) in systemInfo.os"
                    :key="key"
                    style="margin: 20px 0"
                  >
                    <el-col :span="12">{{ key }}:</el-col>
                    <el-col
                      :span="12"
                      :title="value"
                      class="ellipsis"
                      @click="copyToClipboard(value)"
                      >{{ value }}</el-col
                    >
                  </el-row>
                </template>
                <!-- Disk -->
                <template v-else-if="item.key === 'disk'">
                  <el-row :gutter="20">
                    <el-col :span="24" :sm="12">
                      <el-row
                        v-for="(value, key) in systemInfo.disk"
                        :key="key"
                        style="margin: 20px 10px"
                      >
                        <el-col :span="12">{{ key }}:</el-col>
                        <el-col
                          :span="12"
                          class="ellipsis"
                          :title="value"
                          @click="copyToClipboard(value)"
                        >
                          {{ value }}
                        </el-col>
                      </el-row>
                    </el-col>
                    <el-col :span="24" :sm="12">
                      <el-progress
                        type="circle"
                        :percentage="diskPercent"
                        :color="customColorMethod"
                        style="
                          display: flex;
                          justify-content: center;
                          min-height: 200px;
                        "
                      />
                    </el-col>
                  </el-row>
                </template>
                <!-- CPU -->
                <template v-else-if="item.key === 'cpu'">
                  <el-row
                    :gutter="20"
                    v-for="(value, key) in systemInfo.cpu"
                    :key="key"
                    style="margin: 20px 0"
                  >
                    <el-col :span="12">
                      <template v-if="key !== 'states'"> {{ key }}:</template>
                    </el-col>
                    <el-col
                      :span="12"
                      class="ellipsis"
                      :title="
                        key === 'cpu_percent'
                          ? value[0]?.toFixed(2) + '%'
                          : value
                      "
                      @click="
                        copyToClipboard(
                          key === 'cpu_percent'
                            ? value[0]?.toFixed(2) + '%'
                            : value
                        )
                      "
                    >
                      <template v-if="key === 'count'">
                        <span>{{ value }}</span>
                      </template>
                      <template v-else-if="key === 'cpu_percent'">
                        <span> {{ value[0]?.toFixed(2) }} % </span>
                      </template>
                    </el-col>
                  </el-row>
                </template>
                <!-- Ram -->
                <template v-else>
                  <el-row :gutter="20">
                    <el-col :span="24" :sm="12">
                      <el-row
                        v-for="(value, key) in systemInfo.memory"
                        :key="key"
                        style="margin: 20px 10px"
                      >
                        <el-col :span="12">{{ key }}:</el-col>
                        <el-col
                          :span="12"
                          class="ellipsis"
                          :title="value"
                          @click="copyToClipboard(value)"
                        >
                          {{ value }}
                        </el-col>
                      </el-row>
                    </el-col>
                    <el-col :span="24" :sm="12">
                      <el-progress
                        type="circle"
                        :percentage="ramPercent"
                        :color="customColorMethod"
                        style="
                          display: flex;
                          justify-content: center;
                          min-height: 200px;
                        "
                      />
                    </el-col>
                  </el-row>
                </template>
              </div>
            </el-card>
          </el-col>
        </template>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from "vue";
import { getSystemInfo, restartServer, stopServer } from "@/api/server";
import SearchForm from "@/components/SearchForm/index.vue";
import type { TSystemInfo } from "@/types/dashboard";

const searchColumn = ref([]);
const searchAction = computed(() => [
  { label: "重启服务", event: "restart", type: "warning" },
  { label: "停止服务", event: "stop", type: "danger" },
]);
const cards = [
  { title: "RunTime", key: "os" },
  { title: "Disk", key: "disk" },
  { title: "CPU", key: "cpu" },
  { title: "Ram", key: "memory" },
];
const systemInfo = reactive<TSystemInfo>({
  cpu: {
    count: 0,
    cpu_percent: [],
    states: [],
  },
  memory: {
    available: 0,
    total: 0,
    used: 0,
    used_percent: 0,
  },
  disk: {
    partition: "",
    total: 0,
    used: 0,
    free: 0,
    used_percent: 0,
  },
  os: {
    compiler: "",
    go_version: "",
    num_cpu: 0,
    num_goroutine: 0,
    os: "",
  },
});

const diskPercent = ref(0);
const ramPercent = ref(0);
const fetchData = () => {
  getSystemInfo().then((res) => {
    const { data } = res.data;
    systemInfo.cpu = data.cpu;
    systemInfo.memory = data.mem;
    ramPercent.value = Number(data.mem.used_percent.toFixed(2));
    delete systemInfo.memory.used_percent;
    systemInfo.disk = data.disk;
    diskPercent.value = Number(data.disk.used_percent.toFixed(2));
    delete systemInfo.disk.used_percent;
    systemInfo.os = data.os;
  });
};

onMounted(() => {
  fetchData(); // Initial fetch
  // const interval = setInterval(fetchData, 1000 * 10); // Fetch every second
  // onBeforeUnmount(() => clearInterval(interval)); // Clear interval on unmount
});

const customColorMethod = (percentage: number) => {
  if (percentage > 90) return "#ff4949"; // 红色
  if (percentage >= 60) return "#e6a23c"; // 橙色
  if (percentage >= 30) return "#20a0ff"; // 蓝色
  return "#13ce66"; // 绿色
};

//点击复制
const copyToClipboard = async (text: any) => {
  try {
    await navigator.clipboard.writeText(text);
    ElMessage.success("复制成功");
  } catch (err) {
    ElMessage.error("复制失败");
  }
};

//重启
const onRestart = () => {
  restartServer().then((res) => {
    ElMessage.success(res.message || "重启成功");
  });
};

//停止
const onStop = () => {
  stopServer().then((res) => {
    ElMessage.success(res.message || "停止成功");
  });
};
</script>

<style lang="scss" scoped>
.card-item {
  margin: 10px 0; /* 设置上下间距为10px */
  display: flex; /* 使用 Flexbox 布局 */
  flex-direction: column; /* 垂直排列子元素 */
}

.card-item h2 {
  margin-top: 0; /* 去掉 h2 顶部间距 */
}
.ellipsis {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>
