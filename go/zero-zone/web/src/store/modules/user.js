import { defineStore } from "pinia";
import sysUserApi from "@/api/system/sys_user.js";
// 动态导入拿到所有页面 eg: {/src/views/test/index.vue: () => import("/src/views/test/index.vue")}
const views = import.meta.glob("@/views/**/**.vue");
import { useRoute, useRouter } from "vue-router";
import store from "@/store";

export const useUserStore = defineStore("user", () => {
  const route = useRoute();
  const router = useRouter();
  let isLogin = ref(false);
  let tokenObj = ref({});
  let userObj = ref({});
  let userInfo = ref({});
  let routerMap = ref({}); // 全路径'/system/user' -> 路由信息
  let keepAliveList = [];

  // 登录
  async function login(loginObj) {
    if (isLogin.value) {
      return;
    }
    let result = await sysUserApi.login({
      username: loginObj.username.trim(),
      account: loginObj.username.trim(),
      password: loginObj.password.trim(),
      captchaId: "",
      verifyCode: "",
    });
    if (result.code === 200) {
      isLogin.value = true;
      tokenObj.value = result.data;
      getUserInfo();
    }
    return result;
  }

  // 退出登录
  function logout() {
    // 清空pinia存储的数据
    this.$reset();
    this.$reset();

    store.settings.useSettingsStore().$reset();

    // window.localStorage.setItem('user2', 'hello');
    // window.localStorage.removeItem('user2');
    // tips: pinia持久化的无法通过这种方式清空数据，只能删除同样方式存储的值 eg: window.localStorage.setItem('user2', 'hello');
    window.localStorage.clear();
    window.sessionStorage.clear();

    // 跳转登录页
    router.push(`/login?redirect=${route.fullPath}`);
    // window.location.href = '/login';
    location.reload(); // 强制刷新页面
  }

  // 获取用户 & 权限数据
  async function getUserInfo() {
    let result = await sysUserApi.getUserPerm();
    if (result.code === 200) {
      userObj.value = result.data;

      let info = await sysUserApi.getUserInfo();
      //console.log("INFO", info);
      userInfo.value = info.data;

      //console.log("INFO username", userInfo.value.username);
      // 初始化系统设置数据
      store.system.useSystemStore().init();
    }
  }

  const routerList = computed(() => {
    // 拿到后台的权限数据
    return generateRouterList({}, userObj.value.permissionTreeList);
  });

  // 生成侧边栏菜单 & 权限路由数据
  function generateRouterList(parentObj, permList) {
    //console.log(11111)
    //console.log(permList)
    let result = [];
    if (!permList || permList.length === 0) {
      return result;
    }
    // console.log("=====views", views, "permList", permList);

    for (let index = 0; index < permList.length; index++) {
      //console.log(30001)
      let permItem = permList[index];
      //console.log("permItem", permItem);

      // 填充字段数据
      if (!permItem.meta) {
        permItem.meta = {};
      }

      if (permItem.keepAlive) {
        let view_file_name = "";
        let view_file_names = permItem.newTest.split("/");
        if (view_file_names.length) {
          view_file_name = view_file_names[view_file_names.length - 1];
        }
        if (view_file_name.length) {
          if (keepAliveList.indexOf(view_file_name) == -1) {
            keepAliveList.push(view_file_name);
          }
        }
      }

      //permItem.meta.keepAlive = false;
      //console.log(30003)
      if (!permItem.meta.isParentView) {
        permItem.meta.isParentView = false;
      }
      //console.log(30004)
      if (!permItem.meta.sort) {
        permItem.meta.sort = 10000;
      }
      permItem.meta.isShow = permItem.isShow;
      permItem.meta.icon = permItem.icon;
      //console.log(30005)

      //let title = permItem.meta.title;
      let title = permItem.title;
      //console.log(30006)
      if (title) {
        //console.log(30007)
        if (parentObj.meta) {
          //console.log(30008);
          // [子级]
          // 面包屑数据
          //console.log(
          // "面包屑数据",
          //parentObj.meta.breadcrumbItemList.concat([title]),
          //);
          permItem.meta.breadcrumbItemList =
            parentObj.meta.breadcrumbItemList.concat([title]);
          // 全路径
          //permItem.meta.fullPath = parentObj.meta.fullPath + "/" + permItem.path;
          //console.log("全路径", permItem.path);
        } else {
          //console.log(30009)
          // [顶级]
          //permItem.meta.breadcrumbItemList = [title];
          //permItem.meta.fullPath = permItem.path;
          //permItem.meta.title = title;
          permItem.meta.breadcrumbItemList = [title];
        }
        permItem.meta.fullPath = permItem.path;
        permItem.meta.title = title;
        //permItem.meta.keepAlive = false;
      }

      //console.log(30010)
      // 组件页面显示处理
      //console.log("map views", permItem);

      //permItem.component = views[`/src/views/${permItem.component}.vue`];
      permItem.component = views[`/src/views/${permItem.newTest}.vue`];
      // permItem.component = views[`/src/views/test/index.vue`];
      //console.log("map views", "`/src/views/${permItem.component}.vue`")
      // console.log("30011 permItem.meta.fullPath", permItem.meta.fullPath);
      //console.log("30011 permItem.component", permItem.component);
      //console.log("30011 permItem.newTest", permItem.newTest);
      //console.log("30012 permItem", permItem);

      routerMap.value[permItem.meta.fullPath] = permItem;
      // 递归处理
      if (permItem.children.length > 0) {
        //console.log(30013);
        permItem.children = generateRouterList(permItem, permItem.children);
      }

      //console.log(30014)
      result.push(permItem);
    }

    //console.log("routerMap", routerMap);

    // 从小到大 升序
    result = result.sort((a, b) => {
      return a.meta.sort - b.meta.sort;
    });

    //console.log("result", result);
    return result;
  }

  return {
    isLogin,
    login,
    logout,
    tokenObj,
    userObj,
    userInfo,
    getUserInfo,
    routerList,
    routerMap,
    keepAliveList,
  };
});
