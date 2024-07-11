<template>
  <div>
    <iframe ref="authIframe" :src="authUrl" @load="onIframeLoad"></iframe>
  </div>
</template>
<style>
iframe {
  width: 100%;
  height: 960px;
  border: none; /* 可选，去掉边框 */
}
</style>
<script>
export default {
  data() {
    return {
      authUrl: '',
    };
  },
  created() {
    this.authUrl = decodeURIComponent(this.$route.query.authUrl);
    console.log("获得url：" + this.authUrl);
  },
  methods: {
    onIframeLoad() {
      window.addEventListener('message', this.receiveMessage, false);
    },
    receiveMessage(event) {
      // if (event.origin !== 'https://third-party-authorization-page.com') return;
      // 根据第三方页面发送的数据结构来处理
      console.log('Received message:', event.data);
      // 跳转网址
      window.location.href = 'http://localhost:8003/callback/meituan/auth?auth_code='+event.data.auth_code+'&message=456&state='+event.data.state+'&third_id=5';

      // 处理授权信息
    },
  },
  beforeDestroy() {
    // 清除监听器
    window.removeEventListener('message', this.receiveMessage);
  },
};
</script>
