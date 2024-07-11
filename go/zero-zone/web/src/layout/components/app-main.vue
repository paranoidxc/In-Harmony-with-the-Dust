<template>
  <el-scrollbar>
    <div :style="{ height: height, width: width }">
      <!-- 路由视图 -->
      <router-view v-slot="{ Component }">
        <keep-alive :include="keepAliveList">
          <component :is="Component" />
        </keep-alive>
      </router-view>
    </div>
  </el-scrollbar>
</template>
<script setup>
defineProps({
  height: {
    type: String,
    default: "",
  },
  width: {
    type: String,
    default: "",
  },
});
import { getCurrentInstance, toRefs } from "vue";

const { proxy } = getCurrentInstance();
let { keepAliveList } = toRefs(proxy.$store.user.useUserStore());
</script>
<style lang="scss" scoped></style>
