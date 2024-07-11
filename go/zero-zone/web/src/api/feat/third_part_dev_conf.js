import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + "/thirdPartDevConf/list",
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + "/thirdPartDevConf/page",
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + "/thirdPartDevConf/create",
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + "/thirdPartDevConf/update",
      method: "post",
      data,
    });
  },
  authUrl(id) {
    return request({
      url: BASE_API + "/thirdPartDevConf/authUrl",
      method: "get",
      params: { id: id },
    });
  },
  qrCode(id) {
    return request({
      url: BASE_API + "/thirdPartDevConf/qrCode",
      method: "get",
      params: { id: id },
    });
  },
  detail(id) {
    return request({
      url: BASE_API + "/thirdPartDevConf/detail",
      method: "get",
      params: { id: id },
    });
  },
  refreshClientToken(id) {
    return request({
      url: BASE_API + "/thirdPartDevConf/refreshClientToken",
      method: "get",
      params: { id: id },
    });
  },
  delete(id) {
    return request({
      url: BASE_API + "/thirdPartDevConf/delete",
      method: "post",
      data: { id: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + "/thirdPartDevConf/deletes",
      method: "post",
      data: { id: ids },
    });
  },
};
