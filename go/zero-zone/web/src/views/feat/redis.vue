<template>
  <div style="margin: 20px">
    <div class="">
      <el-input
        class="el-inp"
        v-model="inputQuery"
        placeholder="搜索cacheKey"
      />
    </div>

    <!-- query -->
    <div class="query-box">
      <div class="btn-list">
        <el-button
          type="danger"
          @click="handleDelList"
          v-if="multipleSelection.length > 0"
        >
          <el-icon class="el-icon--left">
            <Delete />
          </el-icon>
          删除多选
        </el-button>
      </div>
      <div class="btn-list"></div>
    </div>

    <el-table
      ref="multipleTableRef"
      @selection-change="handleSelectionChange"
      :data="tableData"
      height="600"
      stripe
      style="width: 100%"
      border
    >
      <el-table-column fixed type="selection" width="55" />
      <el-table-column prop="key" label="Cache Key" width="" />
      <el-table-column fixed="right" label="操作" width="160">
        <template #default="scope">
          <el-popconfirm
            title="确定要删除么?"
            @confirm="handleRowDel(scope.row)"
          >
            <template #reference>
              <el-button size="small" type="danger"> 删 除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- dialog -->
    <el-dialog
      draggable
      @open="handleDialogOpen"
      v-model="dialogFormVisible"
      :title="dialogType === 'create' ? '新增' : '编辑'"
    >
      <el-form :model="tableForm" ref="tableFormRef" :rules="rules">
        <el-form-item
          style="display: none"
          v-if="dialogType === 'edit'"
          label="编号"
          :label-width="80"
        >
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>
        <el-form-item label="应用名称" prop="name">
          <el-input v-model="tableForm.name" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="AppID" prop="appId">
          <el-input v-model="tableForm.appId" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="AppSecret" prop="appSecret">
          <el-input v-model="tableForm.appSecret" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="client_token" prop="clientToken">
          <el-input v-model="tableForm.clientToken" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="临时调用凭证" prop="xcode">
          <el-input v-model="tableForm.xcode" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="接口调用凭证" prop="accessToken">
          <el-input v-model="tableForm.accessToken" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="刷新令牌" prop="refreshToken">
          <el-input v-model="tableForm.refreshToken" placeholder="" clearable />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogConfirm"> 确 认 </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";

const { proxy } = getCurrentInstance();

import sysTableApi from "@/api/feat/redis.js";

let inputQuery = ref(""); // 搜索框数据
let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let copyTableData = $ref([]);
let tableForm = $ref({
  status: 1,
});
let dialogFormVisible = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(10);
let total = $ref(0);
let curPage = $ref(1);

const rules = $ref({});

// 监听搜索框
watch(inputQuery, (val) => {
  if (val.length > 0) {
    tableData = copyTableData;
    console.log("val", val);
    // 过滤自己的name然后使用match正则匹配输入的name
    // toLowerCase 是为了统一名称都是小写, 方便检索
    tableData = tableData.filter((item) =>
      item.key.toLowerCase().match(val.toLowerCase()),
    );
  } else {
    tableData = copyTableData;
  }
});

//查询
const onSearchSubmit = async () => {
  tableSearchForm.page = 1;
  tableSearchForm.limit = limit;

  sysTableApi.list(tableSearchForm).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
    }
  });
};

/* 请求分页 */
const getTableDataList = async (cur, limit) => {
  let res = await sysTableApi.list();
  if (res.code == 200) {
    tableData = res.data.list;
    copyTableData = tableData;
  }
};
getTableDataList(1, limit);

const handleSizeChange = (val) => {
  limit = val;
  getTableDataList(curPage, val);
};

const handleCurrentChange = (val) => {
  getTableDataList(val, limit);
};

// 删除一条
const reqRowDel = async (id) => {
  await sysTableApi.delete(id);
};

const handleRowDel = async (row) => {
  await sysTableApi.delete(row.key);
  await getTableDataList();
};

const handleDelList = async () => {
  await sysTableApi.deletes(multipleSelection);
  multipleSelection = [];
  await getTableDataList();
};

// 选中
const handleSelectionChange = (val) => {
  multipleSelection = [];
  val.forEach((item) => {
    multipleSelection.push(item.key);
  });
};

// 编辑
const handleEdit = async (row) => {
  dialogFormVisible = true;
  dialogType = "edit";

  let result = await sysTableApi.detail(row.id);
  if (result.code == 200) {
    tableForm = { ...result.data };
  }
};

// 新增
const handleCreate = () => {
  dialogFormVisible = true;
  tableForm = {};
  dialogType = "create";
};

const handleDialogOpen = () => {
  nextTick(() => {
    proxy.$refs.tableFormRef.clearValidate();
  });
};

// 确认
const dialogConfirm = async () => {
  if (dialogType === "create") {
    // 添加数据
    proxy.$refs.tableFormRef.validate((valid) => {
      if (valid) {
        sysTableApi.create(tableForm).then((res) => {
          if (res.code == 200) {
            dialogFormVisible = false;
            getTableDataList(curPage, limit);
          }
        });
      }
    });
  } else {
    // 修改 内容
    proxy.$refs.tableFormRef.validate((valid) => {
      if (valid) {
        sysTableApi.update(tableForm).then((res) => {
          if (res.code == 200) {
            dialogFormVisible = false;
            getTableDataList(curPage, limit);
          }
        });
      }
    });
  }
};
</script>

<style scoped>
.query-box {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}
</style>
