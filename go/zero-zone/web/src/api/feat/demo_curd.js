import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + '/demoCurd/all',
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + '/demoCurd/page',
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + '/demoCurd/create',
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + '/demoCurd/update',
      method: "post",
      data,
    });
  },
  detail(id) {
      return request({
        url: BASE_API + '/demoCurd/detail',
        method: "get",
        params: { firmId: id },
      });
  },
  delete(id) {
    return request({
      url: BASE_API + '/demoCurd/delete',
      method: "post",
      data: { firmId: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + '/demoCurd/deletes',
      method: "post",
      data: { firmId: ids },
    });
  },
};