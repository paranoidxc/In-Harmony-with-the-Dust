import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + '/uhxOrder/list',
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + '/uhxOrder/page',
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + '/uhxOrder/create',
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + '/uhxOrder/update',
      method: "post",
      data,
    });
  },
  detail(id) {
      return request({
        url: BASE_API + '/uhxOrder/detail',
        method: "get",
        params: { id: id },
      });
  },
  delete(id) {
    return request({
      url: BASE_API + '/uhxOrder/delete',
      method: "post",
      data: { id: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + '/uhxOrder/deletes',
      method: "post",
      data: { id: ids },
    });
  },
};