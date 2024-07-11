import request from "@/utils/request";

const BASE_API = "/web/api/system/user";

export default {
  listPage(query, headers) {
    return request({
      url: "/admin/log/login/page",
      method: "get",
      params: query,
      headers,
    });
  },
};
