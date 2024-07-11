<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="应用名称">
            <el-input v-model="tableSearchForm.name" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="AppID">
            <el-input v-model="tableSearchForm.appId" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="AppSecret">
            <el-input v-model="tableSearchForm.appSecret" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="client_token">
            <el-input v-model="tableSearchForm.clientToken" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="临时调用凭证">
            <el-input v-model="tableSearchForm.xcode" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="接口调用凭证">
            <el-input v-model="tableSearchForm.accessToken" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="刷新令牌">
            <el-input v-model="tableSearchForm.refreshToken" placeholder="" clearable />
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
    <div class="query-box">
      <div class="btn-list">
        <el-button type="primary" @click="handleCreate">
          <el-icon class="el-icon--left">
            <Plus />
          </el-icon>
          增加
        </el-button>
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
      height="460"
      stripe
      style="width: 100%"
      border
    >
      <el-table-column fixed type="selection" width="55" />
      <el-table-column fixed prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="应用名称" width="100" />
      <el-table-column prop="appId" label="AppID" width="100" />
      <el-table-column prop="appSecret" label="AppSecret" width="100" />
      <el-table-column prop="clientToken" label="client_token" width="100" />
      <el-table-column prop="xcode" label="临时调用凭证" width="100" />
      <el-table-column prop="accessToken" label="接口调用凭证" width="100" />
      <el-table-column prop="refreshToken" label="刷新令牌" width="100" />
      <el-table-column prop="createdAt" label="创建时间" width="100" />
      <el-table-column prop="updatedAt" label="更新时间" width="100" />
      <el-table-column prop="deletedAt" label="删除时间" width="100" />
      <el-table-column fixed="right" label="操作" width="160">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">
            编 辑
          </el-button>
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

    <el-pagination
      background
      style="display: flex; justify-content: right; margin-top: 10px"
      v-model:current-page="curPage"
      v-model:page-size="limit"
      :page-sizes="[limit, 20, 50, 100, 200, 300, 400, 500]"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />

    <!-- dialog -->
    <el-dialog
      draggable
      @open="handleDialogOpen"
      v-model="dialogFormVisible"
      :title="dialogType === 'create' ? '新增' : '编辑'"
    >
      <el-form :model="tableForm"
        ref="tableFormRef"
        :rules="rules"
      >
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

import sysTableApi from "@/api/feat/douyin.js";

let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableForm = $ref({
  status: 1,
});
let dialogFormVisible = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(10);
let total = $ref(0);
let curPage = $ref(1);

const rules = $ref({
    /*
    name: [
        { required: true, message: '请输入 应用名称', trigger: 'blur' },
    ],
    appId: [
        { required: true, message: '请输入 AppID', trigger: 'blur' },
    ],
    appSecret: [
        { required: true, message: '请输入 AppSecret', trigger: 'blur' },
    ],
    clientToken: [
        { required: true, message: '请输入 client_token', trigger: 'blur' },
    ],
    xcode: [
        { required: true, message: '请输入 临时调用凭证', trigger: 'blur' },
    ],
    accessToken: [
        { required: true, message: '请输入 接口调用凭证', trigger: 'blur' },
    ],
    refreshToken: [
        { required: true, message: '请输入 刷新令牌', trigger: 'blur' },
    ],
    */
})

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

const getStatusLabel = (idx) => {
  const index = statusOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return statusOptions[index].label;
  } else {
  }
};

//查询
const onSearchSubmit = async () => {
  tableSearchForm.page = 1;
  tableSearchForm.limit = limit;

  sysTableApi.page(tableSearchForm).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
      curPage = res.data.pagination.page;
      total = res.data.pagination.page;
    }
  });
};

/* 请求分页 */
const getTableDataList = async (cur, limit) => {
  let res = await sysTableApi.page({ page: cur, limit: limit });
  if (res.code == 200) {
    tableData = res.data.list;
    curPage = res.data.pagination.page;
    total = res.data.pagination.page;
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
  await sysTableApi.delete(row.id);
  await getTableDataList(curPage, limit);
};

const handleDelList = async() => {
  /*
  multipleSelection.forEach((id) => {
	reqRowDel(id)
  });
  */
  await sysTableApi.deletes(multipleSelection);
  multipleSelection = [];
  await getTableDataList(curPage, limit);
};

// 选中
const handleSelectionChange = (val) => {
  multipleSelection = [];
  val.forEach((item) => {
    multipleSelection.push(item.id);
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
  tableForm = {
  };
  dialogType = "create";
};

const handleDialogOpen = () => {
  nextTick( () => {
    proxy.$refs.tableFormRef.clearValidate()
  })
}

// 确认
const dialogConfirm = async () => {
  if (dialogType === "create") {
    // 添加数据
    proxy.$refs.tableFormRef.validate((valid) => {
        if (valid) {
            sysTableApi
              .create(tableForm)
              .then((res) => {
                if (res.code == 200) {
                  dialogFormVisible = false;
                  getTableDataList(curPage, limit);
                }
              })
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