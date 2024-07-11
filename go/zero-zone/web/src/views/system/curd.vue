<template>
  <div style="margin: 20px">
    <div class=""></div>

    <!-- query -->
    <table class="demo-typo-size">
      <tbody>
        <tr>
          <td width="240px">生成CURD</td>
          <td>
            <el-text class="mx-1" type="primary"> 以下全部</el-text>
          </td>
        </tr>
        <tr>
          <td width="240px">生成Api文件</td>
          <td>
            <el-text class="mx-1" type="primary">
              只会生成 goctl 的api文件, api文件加入总的core.api文件
            </el-text>
          </td>
        </tr>
        <tr>
          <td width="">生成Handle/Login文件(重写)</td>
          <td>
            <el-text class="mx-1" type="primary">
              只会利用 goctl api 命令生成对应的文件并重写 logic 文件
            </el-text>
          </td>
        </tr>
        <tr>
          <td width="">生成Handle/Login文件(不重写)</td>
          <td>
            <el-text class="mx-1" type="primary">
              只会利用 goctl api 命令生成对应的文件 logic 文件
            </el-text>
          </td>
        </tr>

        <tr>
          <td width="">生成Model文件</td>
          <td>
            <el-text class="mx-1" type="primary">
              只会利用 goctl model 命令生成对应的文件的模型文件
            </el-text>
          </td>
        </tr>
        <tr>
          <td width="">生成Vue文件</td>
          <td>
            <el-text class="mx-1" type="primary">
              只会覆盖 web/src/xxx.api 和 web/views/feat/xxx.vue 文件
            </el-text>
          </td>
        </tr>
      </tbody>
    </table>
    <br />

    <el-table :data="tableData" height="460" stripe style="width: 100%" border>
      <el-table-column fixed prop="name" label="模型名称" width="200" />
      <el-table-column fixed="right" label="操作" width="">
        <template #default="scope">
          <el-popconfirm
            title="请仔细，确定要执行该操作么?"
            @confirm="handleCreate(scope.row, 1, 0, 0, 0, 0, 0)"
          >
            <template #reference>
              <el-button size="small"> 生成CURD</el-button>
            </template>
          </el-popconfirm>

          <!--
          <el-button size="small" @click="handleCreate(scope.row, 1,0, 0, 0, 0)">
            生成CURD
          </el-button>
          -->
          <el-popconfirm
            title="请仔细，确定要执行该操作么?"
            @confirm="handleCreate(scope.row, 0, 1, 0, 0, 0, 0)"
          >
            <template #reference>
              <el-button size="small"> 生成Api文件</el-button>
            </template>
          </el-popconfirm>
          <!--
          <el-button size="small" @click="handleCreate(scope.row, 0,1,0, 0, 0, 0)">
            生成Api文件
          </el-button>
          -->

          <el-popconfirm
            title="请仔细，确定要执行该操作么?"
            @confirm="handleCreate(scope.row, 0, 0, 1, 0, 0, 1)"
          >
            <template #reference>
              <el-button size="small"> 生成Handle/Logic (重写)</el-button>
            </template>
          </el-popconfirm>

          <el-popconfirm
            title="请仔细，确定要执行该操作么?"
            @confirm="handleCreate(scope.row, 0, 0, 1, 0, 0, 0)"
          >
            <template #reference>
              <el-button size="small"> 生成Handle/Logic (不重写)</el-button>
            </template>
          </el-popconfirm>

          <!--
          <el-button size="small" @click="handleCreate(scope.row, 0, 0,1,0, 0)">
            生成Handle/Logic
          </el-button>
          -->

          <el-popconfirm
            title="请仔细，确定要执行该操作么?"
            @confirm="handleCreate(scope.row, 0, 0, 0, 1, 0, 0)"
          >
            <template #reference>
              <el-button size="small"> 生成Model</el-button>
            </template>
          </el-popconfirm>

          <!--
          <el-button size="small" @click="handleCreate(scope.row, 0, 0, 0,1, 0)">
            生成Model
          </el-button>
          -->

          <el-popconfirm
            title="请仔细，确定要执行该操作么?"
            @confirm="handleCreate(scope.row, 0, 0, 0, 0, 1, 0)"
          >
            <template #reference>
              <el-button size="small"> 生成Vue</el-button>
            </template>
          </el-popconfirm>

          <el-button size="small" @click="handleCreate(scope.row, 0, 0, 0,0, 0, 0, 1)">
            生成菜单权限
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

const { proxy } = getCurrentInstance();

import sysTableApi from "@/api/system/sys_curd.js";

let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableForm = $ref({
  type: 0,
  status: 1,
  parentId: 0,
});
let dialogFormVisible = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(10);
let total = $ref(0);
let curPage = $ref(1);

const small = ref(false);
const background = ref(false);

const getTableDataList = async (cur = 1, limit = 10) => {
  let result = await sysTableApi.all({ page: cur, limit: limit });
  if (result.code == 200) {
    tableData = result.data.list;
  }
};
getTableDataList();
/* 请求分页 */

const handleSizeChange = (val) => {
  limit = val;
  getTableDataList(curPage, val);
};

const handleCurrentChange = (val) => {
  getTableDataList(val, limit);
};

// 编辑
const handleCreate = (
  row,
  isAll,
  isApi,
  isHandle,
  isModel,
  isVue,
  isLogicWrite,
  isMenu = 0,
) => {
  tableForm = { ...row };
  tableForm.isAll = isAll;
  tableForm.isApi = isApi;
  tableForm.isHandle = isHandle;
  tableForm.isModel = isModel;
  tableForm.isVue = isVue;
  tableForm.isLogicWrite = isLogicWrite;
  tableForm.isMenu = isMenu;
  sysTableApi
    .create(tableForm)
    .then((res) => {
      if (res.code == 200) {
        ElMessage({
          message: "操作成功",
          type: "success",
          plain: true,
        });
      } else {
        ElMessage({
          message: "操作失败",
          type: "erorr",
          plain: true,
        });
      }
    })
    .catch(() => {});
};
</script>

<style scoped>
.query-box {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}
</style>
