import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  list(query) {
    return request({
      url: BASE_API + "/redis/list",
      method: "get",
      params: query,
    });
  },
  detail(id) {
    return request({
      url: BASE_API + "/redis/detail",
      method: "get",
      params: { id: id },
    });
  },
  delete(key) {
    return request({
      url: BASE_API + "/redis/delete",
      method: "post",
      data: { key: key },
    });
  },
  deletes(keys) {
    return request({
      url: BASE_API + "/redis/deletes",
      method: "post",
      data: { key: keys },
    });
  },
};
