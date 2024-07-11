import request from "@/utils/request";

const BASE_API = "/web/api/system/user";

export default {
  // 获取验证码
  getCaptcha() {
    return request({
      url: "/captcha?t=" + new Date().getTime().toString(),
      method: "get",
    });
  },
  // 登录
  login(data) {
    return request({
      url: "/admin/user/login",
      method: "post",
      data,
      // headers: {
      //   // 客户端信息Base64明文：web:123456
      //   Authorization: 'Basic d2ViOjEyMzQ1Ng==',
      // },
    });
  },
  // 注销
  logout() {
    return request({
      //url: '/auth/logout',
      //method: "delete",
      url: "/admin/user/logout",
      method: "post",
    });
  },
  // 获取用户权限
  getUserPerm() {
    return request({
      //url: '/web/api/system/perm/getUserPerm',
      url: "/admin/user/permmenu",
      method: "get",
      // params: { systemSource: 0 }
    });
  },
  getUserInfo() {
    return request({
      url: "/admin/user/info",
      method: "get",
    });
  },
  getProfileInfo() {
    return request({
      url: "/admin/user/profile/info",
      method: "get",
    });
  },
  //修改个人资料
  updateProfileInfo(data) {
    return request({
      url: "/admin/user/profile/update",
      method: "post",
      data,
    });
  },
  // 修改密码
  updatePassword(data) {
    return request({
      url: "/admin/user/password/update",
      method: "post",
      data,
    });
  },
  updateAvatar() {
    return request({
      url: "/admin/user/avatar/generate",
      method: "get",
    });
  },

  listPage(query, headers) {
    return request({
      //url: BASE_API + "/",
      url: "/admin/sys/user/page",
      method: "get",
      params: query,
      headers,
    });
  },
  add(data) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/user/add",
      method: "post",
      data,
    });
  },
  detail(id) {
    return request({
      url: "/admin/sys/user/detail",
      method: "get",
      params: { id: id },
    });
  },
  update(data) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/user/update",
      method: "post",
      data,
    });
  },
  delete(data) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/user/delete",
      method: "post",
      data,
    });
  },
  updateStatus(id, status) {
    return request({
      url: BASE_API + "/updateStatus",
      method: "post",
      data: { userId: id, status: status },
    });
  },
  resetPassword(id) {
    return request({
      url: "/admin/sys/user/password/update",
      method: "post",
      data: { id: id },
    });
  },
  getUserInfoById(userId) {
    return request({
      url: BASE_API + "",
      method: "get",
      params: {
        userId: userId,
      },
    });
  },
  // 保存用户角色
  saveRoleIds(data) {
    return request({
      url: BASE_API + "/saveRoleIds",
      method: "post",
      data: data,
    });
  },
};
