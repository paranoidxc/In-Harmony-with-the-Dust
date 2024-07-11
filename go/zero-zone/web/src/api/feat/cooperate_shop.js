import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + "/cooperateShop/list",
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + "/cooperateShop/page",
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + "/cooperateShop/create",
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + "/cooperateShop/update",
      method: "post",
      data,
    });
  },
  detail(id) {
    return request({
      url: BASE_API + "/cooperateShop/detail",
      method: "get",
      params: { id: id },
    });
  },
  delete(id) {
    return request({
      url: BASE_API + "/cooperateShop/delete",
      method: "post",
      data: { id: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + "/cooperateShop/deletes",
      method: "post",
      data: { id: ids },
    });
  },
};
