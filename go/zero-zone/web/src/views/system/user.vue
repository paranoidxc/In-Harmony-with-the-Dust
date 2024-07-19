<template>
  <div style="margin: 20px">
    <div class="">
      <!--
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="账号">
          <el-input
            v-model="tableSearchForm.account"
            clearable
            placeholder="账号"
          ></el-input>
        </el-form-item>

        <el-form-item label="状态">
          <el-select
            v-model.number="tableSearchForm.status"
            clearable
            placeholder="请选择"
            style="width: 100px"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
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
      -->
    </div>

    <!-- query -->
    <div
      class="query-box"
      style="display: flex; justify-content: space-between; margin-bottom: 10px"
    >
      <div class="btn-list">
        <el-button type="primary" @click="handleAdd">
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
      <el-table-column fixed prop="id" label="编号" width="60" />
      <el-table-column fixed prop="account" label="账号" width="100" />
      <el-table-column prop="username" label="姓名" width="" />
      <el-table-column prop="roles" label="角色" width="">
        <template #default="scope">
          <div class="flex gap-2">
            <!--<el-tag type="primary">Tag 1</el-tag>-->
            {{ showRoleNames(scope.row.roles) }}
          </div>
        </template>
      </el-table-column>
      <!--
      <el-table-column prop="nickname" label="昵称" width="100" />
        <el-table-column class="dsN" prop="gender" label="性别" width="100">
          <template #default="scope">
            {{ getGenderLabel(scope.row.gender) }}
          </template>
        </el-table-column>
        <el-table-column class="dsN" prop="mobile" label="手机号" width="120" />
        <el-table-column class="dsN" prop="orderNum" label="排序值" width="100" />
        -->
      <el-table-column prop="status" label="状态" width="60">
        <template #default="scope">
          {{ getStatusLabel(scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" width="240">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">
            编 辑
          </el-button>

          <el-popconfirm
            :title="newRestPwd"
            @confirm="handleResetPwd(scope.row)"
          >
            <template #reference>
              <el-button size="small">重置密码</el-button>
            </template>
          </el-popconfirm>

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
      v-model="dialogFormVisible"
      @open="handleDialogOpen"
      :title="dialogType === 'add' ? '新增' : '编辑'"
    >
      <el-form ref="tableFormRef" :model="tableForm">
        <el-form-item label="编号" :label-width="140" style="display: none">
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>

        <el-form-item
          v-if="dialogType !== 'edit'"
          label="账号"
          :label-width="140"
          prop="account"
          :rules="[{ required: true, trigger: 'blur', message: '请输入账号' }]"
        >
          <el-input v-model="tableForm.account" autocomplete="off" />
        </el-form-item>

        <el-form-item
          v-if="dialogType === 'edit'"
          label="账号"
          :label-width="140"
          prop="account"
          :rules="[{ required: true, trigger: 'blur', message: '请输入账号' }]"
        >
          <el-input disabled v-model="tableForm.account" autocomplete="off" />
        </el-form-item>

        <el-form-item v-if="dialogType !== 'edit'">
          <label
            id="el-id-9717-256"
            for="el-id-9717-285"
            class="el-form-item__label"
            style="width: 140px"
            >新用户密码</label
          >
          {{ newPwd }}
        </el-form-item>

        <el-form-item
          label="姓名"
          :label-width="140"
          prop="username"
          :rules="[{ required: true, trigger: 'blur', message: '请输入姓名' }]"
        >
          <el-input v-model="tableForm.username" autocomplete="off" />
        </el-form-item>

        <el-form-item
          class="dsN"
          label="昵称"
          :label-width="140"
          prop="nickname"
        >
          <el-input v-model="tableForm.nickname" autocomplete="off" />
        </el-form-item>
        <el-form-item class="dsN" label="电子邮箱" :label-width="140">
          <el-input v-model="tableForm.email" autocomplete="off" />
        </el-form-item>
        <!--
        <el-form-item
          label="手机号"
          :label-width="140"
          prop="mobile"
          :rules="[
            { required: true, trigger: 'blur', message: '请输入手机号' },
          ]"
        >
          <el-input v-model="tableForm.mobile" autocomplete="off" />
        </el-form-item>
        -->

        <el-form-item label="角色" :label-width="140">
          <el-select
            v-model="tableForm.roles"
            multiple
            clearable
            placeholder="请选择角色"
          >
            <el-option
              v-for="item in roleOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="状态" :label-width="140">
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

        <el-form-item class="dsN" label="性别" :label-width="140">
          <el-select v-model.number="tableForm.gender" placeholder="Select">
            <el-option
              v-for="item in genderOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item class="dsN" label="排序值" :label-width="140">
          <el-input-number
            v-model.number="tableForm.orderNum"
            autocomplete="off"
          />
        </el-form-item>

        <el-form-item class="dsN" label="备注" :label-width="140">
          <el-input v-model="tableForm.remark" autocomplete="off" />
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
import { getCurrentInstance, proxyRefs, nextTick } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

const { proxy } = getCurrentInstance();

import tableDataApi from "@/api/system/sys_user.js";
import tableDataRoleApi from "@/api/system/sys_role.js";

let newPwd = $ref();
const today = new Date();
const year = today.getFullYear(); // 获取年份
const month = today.getMonth() + 1; // 获取月份，+1因为getMonth()返回0-11
const day = today.getDate(); // 获取日
newPwd =
  year +
  "" +
  month.toString().padStart(2, "0") +
  "" +
  day.toString().padStart(2, "0");

let newRestPwd = $ref("");
newRestPwd = "重置后密码为:" + newPwd;
let newUserPwd = $ref("");
newUserPwd = "新账户密码为:" + newPwd;

let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableDataRole = $ref([]); // 表格数据
let roleOptions = $ref([]);
let tableForm = $ref({});
let dialogFormVisible = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(10);
let total = $ref(15);
let curPage = $ref(1);

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

const getTableDataList = async (cur = 1, limit = 10) => {
  let result = await tableDataApi.listPage({
    deptId: 0,
    page: cur,
    limit: limit,
  });
  if (result.code == 200) {
    tableData = result.data.list;
    curPage = result.data.pagination.page;
    total = result.data.pagination.total;
  }
};
getTableDataList(1, limit);

const getTableDataRoleList = async () => {
  let result = await tableDataRoleApi.listPage({});
  if (result.code == 200) {
    tableDataRole = result.data.list;
    tableDataRole.forEach((row) => {
      roleOptions.push({
        value: row.id,
        label: row.name,
      });
    });
  }
};
getTableDataRoleList();
/* 请求分页 */

const handleSizeChange = (val) => {
  limit = val;
  getTableDataList(curPage, val);
};

const handleCurrentChange = (val) => {
  getTableDataList(val, limit);
};

// 删除一条
const RowDel = async (id) => {
  await tableDataApi.delete({ id: id });
};

const handleRowDel = async ({ id }) => {
  await tableDataApi.delete({ id: id });
  ElMessage({
    message: "删除成功",
    type: "success",
    plain: true,
  });
  await getTableDataList(curPage, limit);
};

const handleDelList = async () => {
  let promises = [];
  multipleSelection.forEach((id) => {
    promises.push(RowDel(id));
  });
  multipleSelection = [];
  Promise.all(promises).then((result) => {
    ElMessage({
      message: "批量删除成功",
      type: "success",
      plain: true,
    });
    getTableDataList();
  });

  nextTick(() => {});
};

// 选中
const handleSelectionChange = (val) => {
  multipleSelection = [];
  val.forEach((item) => {
    multipleSelection.push(item.id);
  });
};

const handleResetPwd = async (row) => {
  let result = await tableDataApi.resetPassword(row.id);
  if (result.code == 200) {
    ElMessage({
      message: "重置密码成功",
      type: "success",
      plain: true,
    });
    getTableDataList(curPage, limit);
  }
};

// 编辑
const handleEdit = async (row) => {
  dialogType = "edit";
  //tableForm = { ...row };
  /*
  tableForm.roles = [];
  row.roles.forEach((item) => {
    if (item.id) {
      tableForm.roles.push(item.id);
    }
  });
  */
  let result = await tableDataApi.detail(row.id);
  if (result.code == 200) {
    tableForm = { ...result.data };
    tableForm.roles = [];
    result.data.roles.forEach((item) => {
      if (item.id) {
        tableForm.roles.push(item.id);
      }
    });
    dialogFormVisible = true;
  }
};

// 新增
const handleAdd = () => {
  dialogFormVisible = true;
  tableForm = {
    roleIds: [],
    status: 1,
    orderNum: 0,
    gender: 0,
    job: {
      id: 1,
    },
    dept: {
      id: 1,
    },
  };
  dialogType = "add";
};

//查询
const handleDialogOpen = () => {
  nextTick(() => {
    proxy.$refs.tableFormRef.clearValidate();
  });
};

// 确认
const dialogConfirm = async () => {
  if (dialogType === "add") {
    // 添加数据
    tableForm.avatar = "";
    tableForm.email = "";
    //tableForm.professionId = tableForm.profession.id;
    //tableForm.jobId = tableForm.job.id;
    //tableForm.deptId = tableForm.dept.id;
    tableForm.roleIds = tableForm.roles;
  } else {
    // 修改 内容
    //tableForm.professionId = 1;
    //tableForm.jobId = 1;
    //tableForm.deptId = 1;
    tableForm.roleIds = tableForm.roles;
  }

  proxy.$refs.tableFormRef.validate((valid) => {
    if (valid) {
      if (dialogType === "add") {
        tableDataApi.add(tableForm).then((res) => {
          if (res.code == 200) {
            ElMessage({
              message: "创建成功",
              type: "success",
              plain: true,
            });
            dialogFormVisible = false;
            getTableDataList(curPage, limit);
          }
        });
      } else {
        tableDataApi.update(tableForm).then((res) => {
          if (res.code == 200) {
            ElMessage({
              message: "更新成功",
              type: "success",
              plain: true,
            });
            dialogFormVisible = false;
            getTableDataList(curPage, limit);
          }
        });
      }
    }
  });
};

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

const showRoleNames = (roles) => {
  let names = [];
  if (roles.length) {
    roles.forEach((item) => {
      if (item.id) {
        //names.push(`<el-tag type="primary">{{ item.name }} </el-tag>`);
        names.push(item.name);
      }
    });
    return names.join(", ");
  }
  return "";
};

const getStatusLabel = (idx) => {
  const index = statusOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return statusOptions[index].label;
  }
  return "";
};

const genderOptions = [
  {
    value: 0,
    label: "保密",
  },
  {
    value: 1,
    label: "女",
  },
  {
    value: 2,
    label: "男",
  },
];
const getGenderLabel = (idx) => {
  const index = genderOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return genderOptions[index].label;
  }
  return "";
};
</script>

<style scoped></style>
