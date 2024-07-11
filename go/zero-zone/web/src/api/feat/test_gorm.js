import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + '/testGorm/all',
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + '/testGorm/page',
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + '/testGorm/create',
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + '/testGorm/update',
      method: "post",
      data,
    });
  },
  detail(id) {
      return request({
        url: BASE_API + '/testGorm/detail',
        method: "get",
        params: { id: id },
      });
  },
  delete(id) {
    return request({
      url: BASE_API + '/testGorm/delete',
      method: "post",
      data: { id: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + '/testGorm/deletes',
      method: "post",
      data: { id: ids },
    });
  },
};