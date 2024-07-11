<template>
  <el-scrollbar style="border-bottom: 1px solid #ebeef5">
    <base-right-click class="flex">
      <div
        v-for="item in tabsList"
        :key="item"
        class="item"
        :class="{ active: $route.meta.fullPath === item.meta.fullPath }"
        style="display: inline-block; white-space: nowrap"
      >
        <div
          class="flex-between-center"
          style="height: 20px; padding: 18px 10px; margin: 0; border: none"
          @click.right="handleRightClick(item, $event)"
        >
          <el-icon :size="8" style="margin-right: 6px" class="tab-star">
          </el-icon>
          <router-link :to="item.meta.fullPath" @click="activeTabs(item)">
            <span class="m-r-4">{{ item.meta.title }}</span>
          </router-link>
          <el-icon
            style="cursor: pointer"
            class="tab-close"
            v-if="item.meta.fullPath !== '/'"
            :size="12"
            @click="handleClose(item)"
          >
            <Close />
          </el-icon>
        </div>
      </div>

      <template #right-show="{ isShow }">
        <div class="right-menu flex-column b-rd-4 bg-color-white">
          <div class="option" @click="handleCloseCurrent">
            <el-icon :size="10">
              <Close />
            </el-icon>
            <span> 关闭当前</span>
          </div>
          <div class="option" @click="handleCloseLeft">
            <el-icon :size="10">
              <Close />
            </el-icon>
            <span> 关闭左侧</span>
          </div>
          <div class="option" @click="handleCloseRight">
            <el-icon :size="10">
              <Close />
            </el-icon>
            <span> 关闭右侧</span>
          </div>
          <div class="option" @click="handleCloseAll">
            <el-icon :size="10">
              <Close />
            </el-icon>
            <span> 关闭所有</span>
          </div>
        </div>
      </template>
    </base-right-click>
  </el-scrollbar>
</template>
<script setup>
const { proxy } = getCurrentInstance();
let useSettingsStore = proxy.$store.settings.useSettingsStore();
let { keepAliveList, routerMap } = toRefs(proxy.$store.user.useUserStore());
let { tabsList } = toRefs(useSettingsStore);
let { curRouter } = toRefs(useSettingsStore);
let { activeTabs } = useSettingsStore;
let chooseItem = $ref(null);
let activeName = $ref("");

import { useRouter } from "vue-router";

const router = useRouter();
//let xxx = $ref(router.currentRoute);

let left = $ref(0);
let top = $ref(0);
let isCollapse = $ref(false);
let rightActive = $ref("");
let contextMenuVisible = $ref(false);
const openContextMenu = (e) => {
  let id = "";
  if (e.srcElement.nodeName === "SPAN") {
    id = e.srcElement.offsetParent.id;
  } else {
    id = e.srcElement.id;
  }
  if (id) {
    contextMenuVisible = true;
    let width;
    if (isCollapse) {
      width = 54;
    } else {
      width = 220;
    }

    left = e.clientX - width;
    //top = e.clientY + 10;
    top = 40;
    rightActive = id.substring(4);
  }
};

watch(
  () => contextMenuVisible,
  () => {
    if (contextMenuVisible) {
      document.body.addEventListener("click", () => {
        contextMenuVisible = false;
      });
    } else {
      document.body.removeEventListener("click", () => {
        contextMenuVisible = false;
      });
    }
  },
);

const initActiveName = () => {
  let activeRouter = $ref({});
  for (const element of tabsList.value) {
    if (element.meta.title == tab.props.name) {
      activeRouter = element;
      activeName = element.meta.title;
      //console.log("ACTIVE", element);
      break;
    }
  }

  proxy.$router.push({ path: activeRouter.meta.fullPath });
};
const handleClick = (tab, event) => {
  //console.log(tab, event);
  //console.log(tab.props.name);
  //console.log(tabsList);
  //console.log(tabsList.value);
  let activeRouter = $ref({});
  for (const element of tabsList.value) {
    if (element.meta.title == tab.props.name) {
      activeRouter = element;
      activeName = element.meta.title;
      //console.log("ACTIVE", element);
      break;
    }
  }

  proxy.$router.push({ path: activeRouter.meta.fullPath });
};

const handleTabsEdit = (targetName, action) => {
  if (action === "add") {
  } else if (action === "remove") {
    console.log("remove");
    console.log(targetName);
    let activeRouter = $ref({});
    for (let i = 0; i < tabsList.value.length; i++) {
      let element = tabsList.value[i];
      if (element.meta.title == targetName) {
        activeRouter = element;
        break;
      }
    }
    //console.log(activeRouter);
    tabsList.value.splice(tabsList.value.indexOf(activeRouter), 1);
  }
};

// 保留首页
watch(
  tabsList,
  (newValue) => {
    if (newValue.length === 0) {
      tabsList.value.push({ meta: { title: "首页", fullPath: "/" } });
    }
  },
  { immediate: true, deep: true },
);

function handleClose(item) {
  let path = item.newTest.replace("feat/", "");
  let pos = keepAliveList.value.indexOf(path);
  if (pos > -1) {
    keepAliveList.value.splice(pos, 1);
  }
  let findPos = tabsList.value.indexOf(item);
  tabsList.value.splice(findPos, 1);
  //tabsList.value.splice(tabsList.value.indexOf(item), 1);
  if (item.meta.fullPath == proxy.$route.meta.fullPath) {
    let prev = findPos - 1;
    if (prev > -1) {
      let ret = tabsList.value[prev];
      proxy.$router.push({ path: ret.path });
    }
  }
}

function handleRightClick(item) {
  chooseItem = item;
}

function handleCloseCurrent() {
  handleClose(chooseItem);
}

function handleCloseAll() {
  tabsList.value = [];
}

function handleCloseLeft() {
  tabsList.value.splice(0, tabsList.value.indexOf(chooseItem));
}

function handleCloseRight() {
  let idx = tabsList.value.indexOf(chooseItem);

  //console.log("length", tabsList.value.length);
  //console.log("idx", idx);
  //console.log("chooseItem", chooseItem);

  tabsList.value.splice(idx, tabsList.value.length - idx);
}
</script>
<style lang="scss" scoped>
.app {
  position: relative;
  border-top: 1px solid #ebeef5;

  .item {
    margin-top: -1px;
    border-top: 1px solid #ebeef5;
    border-right: 1px solid #ebeef5;

    .tab-star {
      background: #dddddd;
      border-radius: 20px;
    }

    .tab-star,
    .tab-close,
    span {
      color: #303133;
      font-size: 14px;
    }

    &.active {
      color: #4d70ff;
      color: var(--el-color-primary-dark-2);
      border-top: 2px solid #4d70ff;
      border-top: 2px solid var(--el-color-primary-dark-2);
      font-weight: bold;
      background: #f5f9fe;
      background: var(--el-color-primary-light-9);
    }

    &.active span,
    &.active .tab-close,
    &.active .tab-star {
      /*color: var(--el-color-primary);*/
      color: #4d70ff;
      color: var(--el-color-primary-dark-2);
    }

    &.active .tab-star {
      background: #4d70ff;
      background: var(--el-color-primary-dark-2);
    }

    .tab-close:hover {
      background: #ccc;
      border-radius: 20px;
    }
  }

  .item:last-child {
    border: 1px solid red !important;
  }

  .right-menu {
    .option {
      text-align: center;
      padding: 5px 10px;
      cursor: pointer;

      &:hover {
        background: #eee;
      }
    }
  }
}

a {
  text-decoration: none; // 去掉下换线
  color: black; //文字颜色更改
}

.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}

.demo-tabs .custom-tabs-label .el-icon {
  vertical-align: middle;
}

.demo-tabs .custom-tabs-label span {
  vertical-align: middle;
  margin-left: 4px;
}

.contextmenu {
  width: 100px;
  margin: 0;
  border: 1px solid #ccc;
  z-index: 3000;
  position: absolute;
  list-style-type: none;
  padding: 5px 0;
  font-size: 14px;
  color: #333;
}
</style>
