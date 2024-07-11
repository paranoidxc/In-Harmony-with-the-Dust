import request from "@/utils/request";

const BASE_API = "/admin/feat";

export default {
  all(query) {
    return request({
      url: BASE_API + '/{{ .LowerCaseName }}/list',
      method: "get",
      params: query,
    });
  },
  page(query) {
    return request({
      url: BASE_API + '/{{ .LowerCaseName }}/page',
      method: "get",
      params: query,
    });
  },
  create(data) {
    return request({
      url: BASE_API + '/{{ .LowerCaseName }}/create',
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      url: BASE_API + '/{{ .LowerCaseName }}/update',
      method: "post",
      data,
    });
  },
  detail(id) {
      return request({
        url: BASE_API + '/{{ .LowerCaseName }}/detail',
        method: "get",
        params: { {{ .PrimaryKeyJson }}: id },
      });
  },
  delete(id) {
    return request({
      url: BASE_API + '/{{ .LowerCaseName }}/delete',
      method: "post",
      data: { {{ .PrimaryKeyJson }}: id },
    });
  },
  deletes(ids) {
    return request({
      url: BASE_API + '/{{ .LowerCaseName }}/deletes',
      method: "post",
      data: { {{ .PrimaryKeyJson }}: ids },
    });
  },
};