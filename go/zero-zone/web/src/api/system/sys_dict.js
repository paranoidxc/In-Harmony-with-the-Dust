import request from "@/utils/request";

const BASE_API = "/web/api/system/dict";

export default {
  list(query) {
    return request({
      url: "/admin/config/dict/list",
      method: "get",
      params: query,
    });
  },

  listPage(query, headers) {
    return request({
      url: "/admin/config/dict/data/page",
      method: "get",
      params: query,
      headers,
    });
  },

  add(data) {
    return request({
      url: "/admin/config/dict/add",
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: "/admin/config/dict/update",
      method: "post",
      data,
    });
  },
  delete(data) {
    return request({
      url: "/admin/config/dict/delete",
      method: "post",
      data,
    });
  },
};
