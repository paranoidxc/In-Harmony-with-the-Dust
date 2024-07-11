import request from "@/utils/request";

const BASE_API = "/web/api/system/curd";

export default {
  // 获取用户权限
  all(query) {
    return request({
      url: "/admin/sys/autocurd/all",
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: "/admin/sys/autocurd/create",
      method: "post",
      data,
    });
  },
};
