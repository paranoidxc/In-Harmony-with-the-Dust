<template>
  <el-icon color="#409efc" @click="handleRefresh">
    <RefreshRight />
  </el-icon>

  <el-icon color="#409efc" @click="handleClick">
    <FullScreen />
  </el-icon>
</template>

<script setup>
const handleRefresh = () => {
  window.location.reload();
};
const handleClick = () => {};
</script>
<style lang="scss" scoped></style>
