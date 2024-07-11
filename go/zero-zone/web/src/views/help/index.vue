<template>
    <h1>Hello App!</h1>
    <p>
        <!-- 使用 router-link 组件进行导航 -->
        <!-- 通过传递 `to` 来指定链接 -->
        <!-- `<router-link>` 将呈现一个带有正确 `href` 属性的 `<a>` 标签 -->
        <router-link to="/">Go to Home</router-link>
        <router-link to="/test">Go to Test</router-link>
        <router-link to="/login">login</router-link>

        <br>
    <h1>{{ $filters.sexName(1) }} </h1>
    <button @click="handleClick">click</button>
    <br>


    <h2>pinia</h2>
    <h1>{{ $store.test.useTestStore().count }}</h1>
    <button @click="$store.test.useTestStore().add">click</button>
    <br>

    <h2>demo api request</h2>
    <button @click="handleDemoClick">click</button>
    <h1>{{ res }}</h1>
    <br>

    <h2>全局测试</h2>
    <base-no-data></base-no-data>
    <base-no-data>888</base-no-data>

    <el-row class="mb-4">
        <el-button>Default</el-button>
        <el-button type="primary">Primary</el-button>
        <el-button type="success">Success</el-button>
        <el-button type="info">Info</el-button>
        <el-button type="warning">Warning</el-button>
        <el-button type="danger">Danger</el-button>
    </el-row>
    <el-icon :size="100" color="red">
        <Edit />
    </el-icon>
    </p>
    <!-- 路由出口 -->
    <!-- 路由匹配到的组件将渲染在这里 -->
    <router-view></router-view>
</template>

<script setup>
  import { getCurrentInstance, proxyRefs } from 'vue';
  const {proxy} = getCurrentInstance();
  
  async function handleClick() {
    proxy.submitOk("submitOk");
    proxy.submitFail("submitFail");
  }
  
  let res = $ref(null)
  
  async function handleDemoClick() {
    res = await proxy.$api.demo.time();
  }
  </script>