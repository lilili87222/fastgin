<template>
  <el-dialog
    v-model="dialogVisible"
    title="安全验证"
    width="350px"
    :before-close="handleClose"
  >
    <div id="rotation">
      <div class="check">
        <div class="title">请拖动滑块，使得图片的角度调整为正</div>
        <div class="img-con">
          <img :src="showImg" :style="{ transform: imgAngle }" v-if="showImg" />

          <div class="check-state" v-if="currentState">
            {{ currentState.message }}
          </div>
        </div>
        <div
          ref="sliderCon"
          class="slider-con"
          :class="{ 'err-anim': showError }"
        >
          <div
            ref="slider"
            class="slider"
            :class="{ sliding }"
            :style="{ '--move': `${slidMove}px` }"
            @mousedown="onMouseDown"
            @touchstart="onTouchStart"
            @mouseup="onMouseUp"
            @mousemove="onMouseMove"
            @mouseleave="onMouseUp"
            @touchmove="onTouchMove"
            @touchend="onTouchEnd"
          ></div>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import Cookies from "js-cookie";

import { getCode } from "@/api/system/base";
import store from "@/store";

//验证码三种状态
const showError = ref(false);
const showSuccess = ref(false);
const checking = ref(false);

const sliding = ref(false);
const slidMove = ref(0);
const showImg = ref("");
const imgList = ref<any>({});
const dialogVisible = ref(false);
const sliderConWidth = ref(0);
const sliderWidth = ref(0);
let sliderLeft = ref(0);

const route = useRoute();
const router = useRouter();

const currentState = computed(() => {
  if (showError.value) {
    return { message: "错误" };
  } else if (showSuccess.value) {
    return { message: "正确" };
  } else if (checking.value) {
    return { message: "验证中" };
  }
  return null;
});

//旋转角度
const angle = computed(() => {
  const sliderConWidthValue = sliderConWidth.value ?? 0;
  const sliderWidthValue = sliderWidth.value ?? 0;
  let ratio = slidMove.value / (sliderConWidthValue - sliderWidthValue);
  ratio = Math.max(0, Math.min(1, ratio));
  return 360 * ratio;
});

const imgAngle = computed(() => {
  return `rotate(${angle.value}deg)`;
});

const redirect = ref();
interface QueryType {
  [propname: string]: string;
}
const getOtherQuery = (query: QueryType) => {
  return Object.keys(query).reduce((acc: QueryType, cur) => {
    if (cur !== "redirect") {
      acc[cur] = query[cur];
    }
    return acc;
  }, {});
};

const otherQuery = ref({});
watch(
  () => route.query,
  (query) => {
    if (query) {
      redirect.value = query.redirect || "";
      otherQuery.value = getOtherQuery(query as QueryType);
    }
  },
  { immediate: true }
);

const getData = () => {
  let time = new Date().getTime();
  let params = { time };
  getCode(params).then((res) => {
    imgList.value = res.data;
    showImg.value = "data:image/png;base64," + imgList.value.image;
  });
};

const type = ref("");
//打开
const open = (types: string) => {
  type.value = types;
  getData();
  dialogVisible.value = true;
};

const handleClose = () => {
  resetSlider();
  dialogVisible.value = false;
};

//重置
const resetSlider = () => {
  sliding.value = false;
  slidMove.value = 0;
  showImg.value = "";
  checking.value = false;
  showSuccess.value = false;
  showError.value = false;
  sliderCon.value.style.backgroundSize = "0% 100%";
};

const sliderCon = ref();
const slider = ref();
const onMouseDown = (event) => {
  if (checking.value) return;
  sliding.value = true;
  sliderLeft.value = event.clientX; // 记录鼠标按下时的x位置
  sliderConWidth.value = sliderCon.value.clientWidth; // 记录滑槽的宽度
  sliderWidth.value = slider.value.clientWidth; // 记录滑块的宽度

  // 绑定全局鼠标移动和松开事件
  document.addEventListener("mousemove", onMouseMove);
  document.addEventListener("mouseup", onMouseUp);
};

const emits = defineEmits(["handle"]);
const onMouseUp = () => {
  if (sliding.value && !checking.value) {
    checking.value = true;
    const actions = {
      login: handleLogin,
      reset: handleReset,
      sign: handleRegister,
    };
    actions[type.value]?.(); // 调用对应的处理函数
  }
  sliding.value = false;

  // 移除全局鼠标移动和松开事件
  document.removeEventListener("mousemove", onMouseMove);
  document.removeEventListener("mouseup", onMouseUp);
};

const onMouseMove = (event) => {
  if (sliding.value && !checking.value) {
    let m = event.clientX - sliderLeft.value;
    if (m < 0) {
      m = 0;
    } else if (m > sliderConWidth.value - sliderWidth.value) {
      m = sliderConWidth.value - sliderWidth.value;
    }
    slidMove.value = m;
    // 计算滑动比例并更新背景
    const percentage =
      ((slidMove.value + sliderWidth.value) / sliderConWidth.value) * 100;
    sliderCon.value.style.backgroundSize = `${percentage}% 100%`;
  }
};

const handleLogin = () => {
  emits("handle", {
    captcha_id: imgList.value.captcha_id,
    captcha_code: angle.value,
  });
};

const clear = () => {
  showError.value = true;
  resetSlider();
  getData();
};

const handleReset = () => {
  console.log("handleReset");
};

const handleRegister = () => {
  emits("handle", {
    captcha_id: imgList.value.captcha_id,
    captcha_code: angle.value,
  });
};

const closeDia = () => {
  dialogVisible.value = false;
};
defineExpose({
  open,
  clear,
  closeDia,
});

// 移动端
// 新增触摸事件处理函数
const onTouchStart = (event) => {
  if (checking.value) return;
  sliding.value = true;
  sliderLeft.value = event.touches[0].clientX; // 获取触摸点的X坐标
  sliderConWidth.value = sliderCon.value.clientWidth;
  sliderWidth.value = slider.value.clientWidth;
  document.addEventListener("touchmove", onTouchMove);
  document.addEventListener("touchend", onTouchEnd);
};

const onTouchMove = (event) => {
  if (sliding.value && !checking.value) {
    let moveDistance = event.touches[0].clientX - sliderLeft.value;
    if (moveDistance < 0) {
      moveDistance = 0;
    } else if (moveDistance > sliderConWidth.value - sliderWidth.value) {
      moveDistance = sliderConWidth.value - sliderWidth.value;
    }
    slidMove.value = moveDistance;
    // 计算滑动比例并更新背景
    const percentage =
      ((slidMove.value + sliderWidth.value) / sliderConWidth.value) * 100;
    sliderCon.value.style.backgroundSize = `${percentage}% 100%`;
  }
};

const onTouchEnd = () => {
  sliding.value = false;
  document.removeEventListener("touchmove", onTouchMove);
  document.removeEventListener("touchend", onTouchEnd);
  if (!checking.value) {
    checking.value = true;
    const actions = {
      login: handleLogin,
      reset: handleReset,
      register: handleRegister,
    };
    actions[type.value]?.(); // 调用对应的处理函数
  }
};
</script>

<style lang="scss" scoped>
#rotation {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.check {
  --slider-size: 42px;
  width: 330px;
  background: white;
  border-radius: 16px;
  padding: 20px 0;
  display: flex;
  flex-direction: column;
  align-items: center;

  .title {
    color: #464b52;
    font-size: 16px;
    font-weight: 700;
    margin-bottom: 16px;
  }
}

.check .img-con {
  position: relative;
  overflow: hidden;
  pointer-events: none; /* 禁止任何指针事件，包括拖拽 */
  user-select: none; /* 禁止文本选中 */
  width: 256px;
  height: 256px;
  border-radius: 50%;
}

.check .img-con img {
  width: 100%;
  height: 100%;
  user-select: none;
}

.check .slider-con {
  width: 80%;
  height: var(--slider-size);
  border-radius: 42px;
  margin-top: 1rem;
  position: relative;
  background: linear-gradient(
    to right,
    #8c9eff 0%,
    /* 深灰色 */ #8c9eff 50%,
    /* 在50%位置保持深灰色 */ #8c9eff 100% /* 在100%位置变为浅灰色 */
  );

  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2); /* 增加阴影 */
  background-size: 0% 100%; /* 初始背景宽度为 0，表示没有滑动 */
  background-repeat: no-repeat;
}
.slider-con .slider {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.2);
  cursor: move;
  --move: 0px;
  transform: translateX(var(--move));
  background: url("@/assets/login/code.svg") no-repeat;
  background-size: cover;
}
.slider-con .slider.sliding {
  background: url("@/assets/login/code.svg") no-repeat;
  background-size: cover;
}
.check-state {
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  position: absolute;
  top: 0;
  left: 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

body {
  padding: 0;
  margin: 0;
  background: #fef5e0;
}

.slider-con.err-anim {
  animation: jitter 0.5s;
}

.slider-con.err-anim .slider {
  background: #ff4e4e;
}

@keyframes jitter {
  20% {
    transform: translateX(-5px);
  }

  40% {
    transform: translateX(10px);
  }

  60% {
    transform: translateX(-5px);
  }

  80% {
    transform: translateX(10px);
  }

  100% {
    transform: translateX(0);
  }
}
</style>
