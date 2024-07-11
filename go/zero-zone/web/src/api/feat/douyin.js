import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + '/douyin/all',
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + '/douyin/page',
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + '/douyin/create',
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + '/douyin/update',
      method: "post",
      data,
    });
  },
  detail(id) {
      return request({
        url: BASE_API + '/douyin/detail',
        method: "get",
        params: { id: id },
      });
  },
  delete(id) {
    return request({
      url: BASE_API + '/douyin/delete',
      method: "post",
      data: { id: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + '/douyin/deletes',
      method: "post",
      data: { id: ids },
    });
  },
};