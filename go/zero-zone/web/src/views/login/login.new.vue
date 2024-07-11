<template>
  <!-- Sign In -->
  <div class="app">
    <!--登录表单-->
    <div v-if="!isSignUp" class="wrapper">
      <div
        style="
          width: 50%;
          display: flex;
          justify-content: center;
          align-items: center;
        "
      >
        <div style="width: 86%">
          <transition
            name="animate__animated animate__bounce"
            enter-active-class="animate__fadeInUp"
            leave-active-class="animate__zoomOut"
            appear
          >
            <p
              style="
                text-align: center;
                font-weight: bold;
                font-size: 40px;
                margin-bottom: 20px;
              "
            >
             后台登录
            </p>
          </transition>
          <transition
            name="animate__animated animate__bounce"
            enter-active-class="animate__fadeInUp"
            leave-active-class="animate__zoomOut"
            appear
          >
            <el-form ref="loginFormRef" :model="loginForm">
              <el-form-item
                prop="username"
                :rules="[
                  { required: true, trigger: 'blur', message: '请输入账号' },
                ]"
              >
                <el-input
                  size="large"
                  prop="username"
                  v-model="loginForm.username"
                  prefix-icon="User"
                  placeholder="请输入账号"
                  clearable
                  maxlength="30"
                />
              </el-form-item>
              <el-form-item
                prop="password"
                :rules="[
                  { required: true, trigger: 'blur', message: '请输入密码' },
                ]"
              >
                <el-input
                  size="large"
                  v-model="loginForm.password"
                  prefix-icon="Lock"
                  placeholder="请输入密码"
                  show-password
                  maxlength="50"
                />
              </el-form-item>
              <div style="text-align: center">
                <el-button
                  size="large"
                  type="primary"
                  style="width: 60%"
                  class="m-t-10"
                  @click="handleLogin"
                  >登 录
                </el-button>
              </div>
            </el-form>
          </transition>
        </div>
      </div>

      <div
        class="handle-wrapper"
        style="
          width: 50%;
          display: flex;
          justify-content: center;
          align-items: center;
        "
      >
        <el-button
          class=""
          @click="handleSwitchToSignUp"
          size="large"
          type="primary"
          style="width: 50%"
          >注 册
        </el-button>
      </div>
    </div>

    <!-- Sign Up -->
    <div v-if="isSignUp" class="wrapper">
      <div
        class="handle-wrapper2"
        style="
          width: 50%;
          display: flex;
          justify-content: center;
          align-items: center;
        "
      >
        <el-button
          @click="handleSwitchToSignIn"
          size="large"
          type="primary"
          style="width: 50%"
          >登 录
        </el-button>
      </div>
      <div
        style="
          width: 50%;
          display: flex;
          justify-content: center;
          align-items: center;
        "
      >
        <div style="width: 86%">
          <transition
            name="animate__animated animate__bounce"
            enter-active-class="animate__fadeInUp"
            leave-active-class="animate__zoomOut"
            appear
          >
            <p
              style="
                text-align: center;
                font-weight: bold;
                font-size: 40px;
                margin-bottom: 20px;
              "
            >
              注册账号
            </p>
          </transition>
          <transition
            name="animate__animated animate__bounce"
            enter-active-class="animate__fadeInUp"
            leave-active-class="animate__zoomOut"
            appear
          >
            <el-form>
              <el-form-item
                prop="password"
                :rules="[
                  { required: true, message: 'age is required' },
                  { type: 'number', message: 'age must be a number' },
                ]"
              >
                <el-input
                  size="large"
                  prop="username"
                  v-model="loginForm.username"
                  prefix-icon="User"
                  placeholder="请输入账号"
                  clearable
                  maxlength="30"
                />
              </el-form-item>
              <el-form-item prop="password">
                <el-input
                  size="large"
                  v-model="loginForm.password"
                  prefix-icon="Lock"
                  placeholder="请输入密码"
                  show-password
                  maxlength="50"
                />
              </el-form-item>
              <el-form-item prop="password">
                <el-input
                  size="large"
                  v-model="loginForm.password"
                  prefix-icon="Lock"
                  placeholder="请输入密码"
                  show-password
                  maxlength="50"
                />
              </el-form-item>
              <div style="text-align: center">
                <el-button
                  size="large"
                  type="primary"
                  style="width: 60%"
                  class="m-t-10"
                  @click="handleToSignUp"
                  >注 册
                </el-button>
              </div>
            </el-form>
          </transition>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.app {
  height: 100%;
  width: 100%;
  background: #e5e7ee;
  display: flex;
  justify-content: center;
  align-items: center;
  background: url("https://res.cloudinary.com/dci1eujqw/image/upload/v1616769558/Codepen/waldemar-brandt-aThdSdgx0YM-unsplash_cnq4sb.jpg");
  background-attachment: fixed;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}

.handle-wrapper {
  background: url("https://res.cloudinary.com/dci1eujqw/image/upload/v1616769558/Codepen/waldemar-brandt-aThdSdgx0YM-unsplash_cnq4sb.jpg");
  background-attachment: fixed;
  background-position: right;
  background-repeat: no-repeat;
  background-size: cover;
  border-bottom-rigth-radius: 10px;
  border-top-right-radius: 10px;
  transform: translateX(0);
  transition: transform 0.6s ease-in-out;
  width: 200%;
}

.handle-wrapper2 {
  background: url("https://res.cloudinary.com/dci1eujqw/image/upload/v1616769558/Codepen/waldemar-brandt-aThdSdgx0YM-unsplash_cnq4sb.jpg");
  background-attachment: fixed;
  background-position: left;
  background-repeat: no-repeat;
  background-size: cover;
  border-bottom-left-radius: 10px;
  border-top-left-radius: 10px;
  transform: translateX(0);
  transition: transform 0.6s ease-in-out;
  width: 200%;
}

.wrapper {
  background: #f9f9f9;
  border-radius: 10px;
  width: 50%;
  height: 400px;
  display: flex;
  flex-direction: row;
}
</style>

<script setup>
import { getCurrentInstance, proxyRefs } from "vue";
import { ElMessage } from "element-plus";
import "animate.css";

const { proxy } = getCurrentInstance();
const { login } = proxy.$store.user.useUserStore();
const loginForm = $ref({});

let isSignUp = $ref(false);
let goRight = $ref("");

const handleSwitchToSignUp = () => {
  isSignUp = true;
  //goRight = "animate__slideOutRight animate__bounce animate__delay-2s";
};
const handleSwitchToSignIn = () => {
  isSignUp = false;
};

const handleToSignUp = () => {
  ElMessage({
    message: h("p", { style: "line-height: 1; font-size: 14px" }, [
      h("span", null, "该功能未开放"),
    ]),
  });
};

function handleLogin() {
  proxy.$refs.loginFormRef.validate((valid) => {
    if (valid) {
      login(loginForm).then((res) => {
        if (res.code === 200) {
          let fullPath = proxy.$route.fullPath;
          if (fullPath.startsWith("/login?redirect=")) {
            let lastPath = fullPath.replace("/login?redirect=", "");
            // 跳转到上次退出的页面
            proxy.$router.push({ path: lastPath });
          } else {
            // 跳转到首页
            proxy.$router.push({ path: "/" });
          }
        }
      });
    }
  });
}
</script>
