<template>
  <!-- {{ route.meta }} -->
  <div class="app flex-between-center p-x-10">
    <div class="flex-center-center">
      <div
        class="m-r-10"
        style="cursor: pointer"
        @click="proxy.$store.settings.useSettingsStore().update"
      >
        <el-icon :size="12">
          <component
            :is="
              proxy.$store.settings.useSettingsStore().isShowMenu
                ? 'DArrowRight'
                : 'DArrowLeft'
            "
          />
        </el-icon>
      </div>

      <el-breadcrumb :separator-icon="ArrowRight">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item
          v-for="item in route.meta.breadcrumbItemList"
          :key="item"
        >
          <span class="text-color-grey">{{ item }}</span>
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <!--
            <wx-mp-account />
            <quick-mark />
            <quick-search />
            -->
    <el-dropdown
      class="avatar-container right-menu-item hover-effect"
      trigger="click"
    >
      <div class="flex-center-center" style="cursor: pointer">
        <el-avatar class="" :size="32" :src="userInfo.avatar" />
        <div class="flex-center-center">
          <span class="m-l-6"> {{ userInfo.username }} </span>
          <el-icon :size="16" style="width: 16px">
            <ArrowDown />
          </el-icon>
        </div>
      </div>
      <template #dropdown>
        <el-dropdown-menu style="width: 160px">
          <el-dropdown-item icon="avatar" @click="ucenter"
            >个人中心
          </el-dropdown-item>
          <el-dropdown-item divided icon="reading-lamp" @click="myLogout"
            >退出
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>
<script setup>
import { useRoute } from "vue-router";

const route = useRoute();
import { ArrowRight } from "@element-plus/icons-vue";
import { getCurrentInstance, toRefs } from "vue";
import store from "@/store/index.js";

const { proxy } = getCurrentInstance();

let useUserStore = proxy.$store.user.useUserStore();
let { logout } = useUserStore;
let { userObj, userInfo } = toRefs(useUserStore);

const ucenter = () => {
  proxy.$router.push({ path: "/ucenter" });
};

const myLogout = () => {
  store.user.useUserStore().$reset();
  store.settings.useSettingsStore().$reset();
  window.localStorage.clear();
  window.sessionStorage.clear();
  location.reload(); // 强制刷新页面
};
</script>
<style lang="scss" scoped></style>
