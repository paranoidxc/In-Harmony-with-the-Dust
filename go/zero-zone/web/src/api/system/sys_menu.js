import request from "@/utils/request";

const BASE_API = "/web/api/system/menu";

export default {
  listPage(query) {
    return request({
      url: "/admin/sys/perm/menu/list",
      method: "get",
      params: query,
    });
  },
  add(data) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/perm/menu/add",
      method: "post",
      data,
    });
  },
  update(data) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/perm/menu/update",
      method: "post",
      data,
    });
  },
  delete(id) {
    return request({
      //url: BASE_API,
      url: "/admin/sys/perm/menu/delete",
      method: "post",
      params: {
        id: id,
      },
    });
  },

  menuList(params) {
    return request({
      //url: BASE_API + "/menuList",
      url: "/admin/sys/perm/menu/list",
      method: "get",
      params: params,
    });
  },
  menuTree(params) {
    return request({
      url: BASE_API + "/menuTree",
      method: "get",
      params: params,
    });
  },

  /**
   * 获取菜单关联按钮权限信息
   */
  getPermListByMenuId(id) {
    return request({
      url: BASE_API + "/getPermListByMenuId",
      method: "get",
      params: {
        menuId: id,
      },
    });
  },
  /**
   * 删除菜单关联按钮权限
   */
  deleteMenuReBtnPerm(id) {
    return request({
      url: BASE_API + "/deleteMenuReBtnPerm",
      method: "delete",
      params: {
        id: id,
      },
    });
  },
  /**
   * 新增菜单关联按钮权限
   */
  addMenuReBtnPerm(data) {
    return request({
      url: BASE_API + "/addMenuReBtnPerm",
      method: "post",
      data,
    });
  },
  /**
   * 更新菜单关联按钮权限
   */
  updateMenuReBtnPerm(data) {
    return request({
      url: BASE_API + "/updateMenuReBtnPerm",
      method: "put",
      data,
    });
  },
};
