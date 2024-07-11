import request from "@/utils/request";

const BASE_API = "/web/api/system/user";

export default {
  // 获取用户权限
  listPage(query) {
    return request({
      url: "/admin/sys/dept/list",
      method: "get",
      params: query,
    });
  },
  search(query) {
    return request({
      url: "/admin/sys/dept/search",
      method: "get",
      params: query,
    });
  },
  add(data) {
    return request({
      url: "/admin/sys/dept/add",
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: "/admin/sys/dept/update",
      method: "post",
      data,
    });
  },
  delete(data) {
    return request({
      url: "/admin/sys/dept/delete",
      method: "post",
      data,
    });
  },
};
