import request from "@/utils/request";

const BASE_API = "/web/api/system/role";

export default {
  listPage(query, headers) {
    return request({
      //url: BASE_API + '/listPage',
      url: "/admin/sys/role/list",
      method: "get",
      params: query,
      headers,
    });
  },
  list(query) {
    return request({
      url: BASE_API + "/list",
      method: "get",
      params: query,
    });
  },
  detail(id) {
    return request({
      url: BASE_API + "/detail",
      method: "get",
      params: { roleId: id },
    });
  },
  add(data) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/role/add",
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/role/update",
      method: "post",
      data,
    });
  },
  delete(id) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/role/delete",
      method: "post",
      //params: { roleId: id },
      data: { id: id },
    });
  },
  updateStatus(id, status) {
    return request({
      url: BASE_API + "/updateStatus",
      method: "post",
      data: { roleId: id, status: status },
    });
  },
};
