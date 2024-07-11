import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + '/hxOrder/list',
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + '/hxOrder/page',
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + '/hxOrder/create',
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + '/hxOrder/update',
      method: "post",
      data,
    });
  },
  detail(id) {
      return request({
        url: BASE_API + '/hxOrder/detail',
        method: "get",
        params: { id: id },
      });
  },
  delete(id) {
    return request({
      url: BASE_API + '/hxOrder/delete',
      method: "post",
      data: { id: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + '/hxOrder/deletes',
      method: "post",
      data: { id: ids },
    });
  },
};