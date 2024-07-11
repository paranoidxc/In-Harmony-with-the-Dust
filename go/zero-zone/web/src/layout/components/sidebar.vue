<template>
  <div
    style="display: flex; padding: 0; background: black; flex-direction: column"
  >
    <h2
      style="text-align: center; color: #fff; display: block; padding: 10px 0"
    >
      ZERO-ZONE
    </h2>
    <el-menu
      class="el-menu-my"
      router
      :default-active="$route.meta.fullPath"
      :collapse="false"
      :unique-opened="false"
      @select="handleSelect"
    >
      <el-scrollbar>
        <sidebar-item :router-list="routerList" />
      </el-scrollbar>
    </el-menu>
  </div>
</template>

<script setup>
import sidebarItem from "./sidebar-item.vue";
import { getCurrentInstance, toRefs } from "vue";

const { proxy } = getCurrentInstance();
let { routerList, routerMap } = toRefs(proxy.$store.user.useUserStore());
let { activeTabs } = proxy.$store.settings.useSettingsStore();
let isCollapse = $ref(true);
let { keepAliveList } = toRefs(proxy.$store.user.useUserStore());

//console.log("routerList");
//console.log(routerList);
//console.log("routerMap");
//console.log(routerMap);

/**
 * 选中菜单时触发
 * @param index 选中菜单项的 index  eg: /system/role （router 以 index 作为 path 进行路由跳转，或 router 属性直接跳转）
 * @param indexPath 选中菜单项的 index path eg: ['/system', '/system/role']
 * @param item 选中菜单项
 * @param routeResult vue-router 的返回值（如果 router 为 true）
 */
function handleSelect(index, indexPath, item, routeResult) {
  // console.log(index, indexPath, item, routeResult);
  // proxy.$router.push(index);
  let router = routerMap.value[index];
  let path = router.newTest.replace("feat/", "");
  if (router.meta.keepAlive) {
    let pos = keepAliveList.value.indexOf(path);
    if (pos == -1) {
      keepAliveList.value.push(path);
    }
  }
  activeTabs(routerMap.value[index]);
}
</script>

<style lang="scss" scoped></style>
