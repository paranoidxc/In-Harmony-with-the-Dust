import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  index(query) {
    return request({
      url: BASE_API + "/dashboard/index",
      method: "get",
      params: query,
    });
  },
};
