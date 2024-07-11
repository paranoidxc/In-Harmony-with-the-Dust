<template>
  <!-- 路由出口 -->
  <!-- 路由匹配到的组件将渲染在这里 -->
  <p>
    <!-- 使用 router-link 组件进行导航 -->
    <!-- 通过传递 `to` 来指定链接 -->
    <!-- `<router-link>` 将呈现一个带有正确 `href` 属性的 `<a>` 标签 -->
    <!--<router-link to="/">Go to Home</router-link>
    <router-link to="/test">Go to Test</router-link>
    <router-link to="/login">login</router-link>
    -->
  </p>
  <layout />
</template>

<script setup>
import { getCurrentInstance, proxyRefs } from 'vue';
const {proxy} = getCurrentInstance();
import layout from '@/layout/index.vue';
</script>
