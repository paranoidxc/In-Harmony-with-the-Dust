import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + '/saasCooperateAuth/list',
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + '/saasCooperateAuth/page',
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + '/saasCooperateAuth/create',
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + '/saasCooperateAuth/update',
      method: "post",
      data,
    });
  },
  syncAuth(id) {
    return request({
      url: BASE_API + '/saasCooperateAuth/syncAuth',
      method: "get",
      params: { id: id },
    });
  },
  detail(id) {
      return request({
        url: BASE_API + '/saasCooperateAuth/detail',
        method: "get",
        params: { id: id },
      });
  },
  delete(id) {
    return request({
      url: BASE_API + '/saasCooperateAuth/delete',
      method: "post",
      data: { id: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + '/saasCooperateAuth/deletes',
      method: "post",
      data: { id: ids },
    });
  },
};