<template>
  <div class="app">
    <el-tabs
      v-model="activeName"
      class="demo-tabs"
      type="card"
      @contextmenu.prevent="openContextMenu($event)"
      @tab-click="handleClick"
      @edit="handleTabsEdit"
    >
      <el-tab-pane
        v-for="item in tabsList"
        :label="item.meta.title"
        :name="item.meta.title"
        closable
        :tab="item"
      >
        <template #label>
          <span class="custom-tabs-label">
            <el-icon><CollectionTag /></el-icon>
            <span class="dot">
              {{ item.meta.title }}
            </span>
          </span>
        </template>
      </el-tab-pane>
    </el-tabs>

    <!--自定义右键菜单html代码-->
    <div
      v-show="contextMenuVisible"
      :style="{ left: left + 'px', top: top + 'px' }"
      class="contextmenu"
    >
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
    </div>
  </div>
  <!--{{ $route.meta.fullPath }}
                                                                                                                                            {{ $route.meta.title }}-->
  <div>
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
            style="height: 20px; padding: 16px 10px; margin: 0; border: none"
            @click.right="handleRightClick(item, $event)"
          >
            <el-icon :size="16" style="margin-right: 4px">
              <StarFilled />
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
  </div>
</template>
<script setup>
const { proxy } = getCurrentInstance();
let useSettingsStore = proxy.$store.settings.useSettingsStore();
let { tabsList } = toRefs(useSettingsStore);
let { curRouter } = toRefs(useSettingsStore);
let { activeTabs } = useSettingsStore;
let chooseItem = $ref(null);
let activeName = $ref("");

import { useRouter } from "vue-router";

const router = useRouter();
let xxx = $ref(router.currentRoute);
console.log("router", router.currentRoute);
console.log("xxxx", xxx);

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
  tabsList.value.splice(tabsList.value.indexOf(item), 1);
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

  console.log("length", tabsList.value.length);
  console.log("idx", idx);
  console.log("chooseItem", chooseItem);

  tabsList.value.splice(idx, tabsList.value.length - idx);
}
</script>
<style lang="scss" scoped>
.app {
  position: relative;
  border-top: 1px solid #ebeef5;
  //box-shadow: 0 1px 5px rgba(0, 0, 0, 0.2);

  .item {
    margin-top: -1px;
    border-top: 1px solid #ebeef5;
    border-right: 1px solid #ebeef5;

    span {
      color: #333;
    }

    &.active {
      color: blue;
      border-top: 1px solid blue;
      font-weight: bold;
    }

    &.active span {
      color: blue;
    }

    .tab-close:hover {
      background: #ccc;
      border-radius: 8px;
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
  /*
                                                                                                                                                                                                                                                                                                                                                            border: 1px solid #ccc;
                                                                                                                                                                                                                                                                                                                                                            background: #fff;
                                                                                                                                                                                                                                                                                                                                                            z-index: 3000;
                                                                                                                                                                                                                                                                                                                                                            box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.2);
                                                                                                                                                                                                                                                                                                                                                            */
  position: absolute;
  list-style-type: none;
  padding: 5px 0;
  font-size: 14px;
  color: #333;
}
</style>
