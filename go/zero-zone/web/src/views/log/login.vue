<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="请求信息">
          <el-input
            v-model="tableSearchForm.name"
            clearable
            placeholder="请求信息"
          ></el-input>
        </el-form-item>
        <el-form-item label="响应信息">
          <el-input
            v-model="tableSearchForm.fullName"
            clearable
            placeholder="请输入响应信息"
          ></el-input>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSearchSubmit">
            <el-icon class="el-icon--left">
              <Search />
            </el-icon>
            查询
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- query -->

    <el-table
      ref="multipleTableRef"
      @selection-change="handleSelectionChange"
      :data="tableData"
      height="460"
      stripe
      style="width: 100%"
      border
    >
      <el-table-column fixed type="selection" width="55" />
      <el-table-column prop="id" label="编号" width="100">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">
            查看 #{{ scope.row.id }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="account" label="账号" width="100" />
      <el-table-column prop="ip" label="IP" width="100" />
      <el-table-column prop="url" label="URL" width="200" />
      <el-table-column prop="createTime" label="时间" width="160" />
      <el-table-column prop="request" label="请求信息" width="" />
      <el-table-column prop="status" label="状态" width="60">
        <template #default="scope">
          {{ getStatusLabel(scope.row.status) }}
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      background
      style="display: flex; justify-content: right; margin-top: 10px"
      v-model:current-page="curPage"
      v-model:page-size="limit"
      :page-sizes="[10, 20, 50, 100, 200, 300, 400, 500]"
      :small="small"
      :disabled="disabled"
      :background="background"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />

    <!-- dialog -->
    <el-dialog
      v-model="dialogFormVisible"
      :title="dialogType === 'add' ? '新增' : '编辑'"
    >
      <el-form :model="tableForm">
        <el-form-item
          v-if="dialogType === 'edit'"
          label="编号"
          :label-width="80"
        >
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>
        <el-form-item label="部门简称" :label-width="80">
          <el-input v-model="tableForm.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="部门全称" :label-width="80">
          <el-input v-model="tableForm.fullName" autocomplete="off" />
        </el-form-item>
        <el-form-item label="唯一值" :label-width="80">
          <el-input v-model="tableForm.uniqueKey" autocomplete="off" />
        </el-form-item>

        <el-form-item label="状态" :label-width="80">
          <el-tooltip
            :content="getStatusLabel(tableForm.status)"
            placement="top"
          >
            <el-switch
              v-model.number="tableForm.status"
              class="mt-2"
              inline-prompt
              :active-value="1"
              :inactive-value="0"
            />
          </el-tooltip>
        </el-form-item>

        <el-form-item label="类型" :label-width="80">
          <el-select v-model.number="tableForm.type" placeholder="Select">
            <el-option
              v-for="item in typeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="备注" :label-width="80">
          <el-input v-model="tableForm.remark" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogConfirm"> 确认 </el-button>
        </span>
      </template>
    </el-dialog>

    <el-drawer v-model="drawer" title="日志信息" size="50%">
      <div>
        <el-text class="mx-1" type="primary"
          >账号：{{ tableForm.name }}
        </el-text>
      </div>
      <div>
        <el-text class="mx-1" type="primary"
          >操作时间{{ tableForm.createTime }}
        </el-text>
      </div>
      <div>
        <el-text class="mx-1" type="primary">{{ tableForm.request }}</el-text>
      </div>
    </el-drawer>
  </div>
</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";

const { proxy } = getCurrentInstance();

import tableDataApi from "@/api/system/sys_log.js";

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
let total = $ref(15);
let curPage = $ref(1);
let drawer = $ref(false);

const small = ref(false);
const background = ref(false);
const disabled = ref(false);
const statusOptions = [
  {
    value: 1,
    label: "启用",
  },
  {
    value: 0,
    label: "禁用",
  },
];
const typeOptions = [];

const getStatusLabel = (idx) => {
  if (idx == 0) {
    return "禁用";
  } else if (idx == 1) {
    return "启用";
  }
};

const getTypeLabel = (idx) => {};

const getTableDataList = async (cur = 1, limit) => {
  let result = await tableDataApi.listPage({
    page: cur,
    limit: limit,
  });
  if (result.code == 200) {
    tableData = result.data.list;
    curPage = result.data.pagination.page;
    total = result.data.pagination.total;
  }
};
getTableDataList(curPage, limit);
/* 请求分页 */

const handleSizeChange = (val) => {
  limit = val;
  getTableDataList(curPage, val);
};

const handleCurrentChange = (val) => {
  getTableDataList(val, limit);
};

// 删除一条
// 选中
const handleSelectionChange = (val) => {
  multipleSelection = [];
  val.forEach((item) => {
    multipleSelection.push(item.id);
  });
};

// 编辑
const handleEdit = (row) => {
  //dialogFormVisible = true;
  //dialogType = "edit";
  drawer = true;
  tableForm = { ...row };
};

// 新增
const handleAdd = () => {
  dialogFormVisible = true;
  tableForm = {
    status: 1,
  };
  dialogType = "add";
};

//查询
const onSearchSubmit = async () => {
  tableSearchForm.page = 1;
  tableSearchForm.limit = limit;

  tableDataApi.search(tableSearchForm).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
      curPage = res.data.page;
      total = res.data.total;
    }
  });
};

// 确认
const dialogConfirm = async () => {
  dialogFormVisible = false;
};
</script>

<style scoped></style>
