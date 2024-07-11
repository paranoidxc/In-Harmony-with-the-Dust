<template>
  <div class="app">
    <el-tabs tab-position="left" v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="基本资料" name="first">
        <el-form ref="infoFormRef" :model="infoForm">
          <el-form-item label="头像" :label-width="140">
            <el-avatar
              :src="infoForm.avatar"
              :size="80"
              style="margin-right: 10px"
            />
            <el-input
              class="dsN"
              v-model="infoForm.avatar"
              autocomplete="off"
            />
            <el-button type="primary" @click="handleAvatar">
              生成头像
            </el-button>
          </el-form-item>

          <el-form-item
            label="姓名"
            :label-width="140"
            prop="username"
            :rules="[
              { required: true, trigger: 'blur', message: '请输入新密码' },
            ]"
          >
            <el-input v-model="infoForm.username" autocomplete="off" />
          </el-form-item>
          <el-form-item label="昵称" :label-width="140">
            <el-input v-model="infoForm.nickname" autocomplete="off" />
          </el-form-item>
          <el-form-item label="电子邮箱" :label-width="140">
            <el-input v-model="infoForm.email" autocomplete="off" />
          </el-form-item>
          <el-form-item label="手机号" :label-width="140">
            <el-input v-model="infoForm.mobile" autocomplete="off" />
          </el-form-item>
          <el-form-item label="性别" :label-width="140">
            <el-select v-model.number="infoForm.gender" placeholder="Select">
              <el-option
                v-for="item in genderOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="" :label-width="140">
            <el-button type="primary" @click="handleInfoConfirm">
              更新资料
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="安全设置" name="second">
        <el-form ref="pwdFormRef" :model="pwdForm">
          <el-form-item
            label="旧密码"
            :label-width="140"
            prop="oldPassword"
            :rules="[
              { required: true, trigger: 'blur', message: '请输入旧密码' },
            ]"
          >
            <el-input v-model="pwdForm.oldPassword" autocomplete="off" />
          </el-form-item>
          <el-form-item
            label="新密码"
            :label-width="140"
            prop="newPassword"
            :rules="[
              { required: true, trigger: 'blur', message: '请输入新密码' },
            ]"
          >
            <el-input v-model="pwdForm.newPassword" autocomplete="off" />
          </el-form-item>
          <el-form-item
            label="确认密码"
            :label-width="140"
            prop="repeatNewPassword"
            :rules="[
              {
                required: true,
                trigger: 'blur',
                message: '请再次输入新密码',
              },
            ]"
          >
            <el-input v-model="pwdForm.repeatNewPassword" autocomplete="off" />
          </el-form-item>
          <el-form-item label="" :label-width="140">
            <el-button type="primary" @click="handlePwdConfirm">
              更新密码
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import md5 from "js-md5";

const activeName = $ref("first");

import { getCurrentInstance, proxyRefs } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

const { proxy } = getCurrentInstance();

let infoForm = $ref({});
let pwdForm = $ref({});

import tableDataApi from "@/api/system/sys_user.js";

const getProfileInfo = async () => {
  let res = await tableDataApi.getProfileInfo();
  if (res.code === 200) {
    infoForm = res.data;
  }
};
getProfileInfo();

const handleInfoConfirm = async () => {
  proxy.$refs.infoFormRef.validate((valid) => {
    if (valid) {
      tableDataApi.updateProfileInfo(infoForm).then((res) => {
        if (res.code === 200) {
          ElMessage({
            message: "操作成功",
            type: "success",
            plain: true,
          });
        }
      });
    }
  });
};

const handlePwdConfirm = async () => {
  proxy.$refs.pwdFormRef.validate((valid) => {
    if (valid) {
      let copyPwdForm = {
        oldPassword: md5(pwdForm.oldPassword),
        newPassword: md5(pwdForm.newPassword),
        repeatNewPassword: md5(pwdForm.repeatNewPassword),
      };
      if (copyPwdForm.newPassword != copyPwdForm.repeatNewPassword) {
        ElMessage({
          message: "新密码 和 确认密码不一致",
          type: "error",
          plain: true,
        });
        return;
      }

      tableDataApi.updatePassword(copyPwdForm).then((res) => {
        if (res.code === 200) {
          ElMessage({
            message: "操作成功",
            type: "success",
            plain: true,
          });
          pwdForm = {};
        }
      });
    }
  });
};

const handleAvatar = async () => {
  tableDataApi.updateAvatar().then((res) => {
    if (res.code === 200) {
      infoForm.avatar = res.data.avatarUrl;
    }
  });
};
const handleClick = (tab, event) => {
  console.log(tab, event);
};

const genderOptions = [
  {
    value: 0,
    label: "保密",
  },
  {
    value: 1,
    label: "女",
  },
  {
    value: 2,
    label: "男",
  },
];
const getGenderLabel = (idx) => {
  const index = genderOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return genderOptions[index].label;
  }
  return "";
};
</script>

<style scoped lang="scss">
.app {
  margin: 10px;
}

demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
</style>
