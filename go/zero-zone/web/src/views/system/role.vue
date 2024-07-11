<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline style="display: none">
        <el-form-item label="部门简称">
          <el-input
            v-model="tableSearchForm.name"
            clearable
            placeholder="请输入部门简称"
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
      <el-table-column fixed prop="name" label="角色名称" width="100" />
      <el-table-column prop="uniqueKey" label="角色标识" width="100" />
      <el-table-column prop="orderNum" label="排序值" width="100" />
      <el-table-column prop="status" label="状态" width="60">
        <template #default="scope">
          {{ getStatusLabel(scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" width="" />
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

    <!-- dialog -->
    <el-dialog
      v-model="dialogFormVisible"
      @open="handleDialogOpen"
      @close="handleDialogClose"
      draggable
      :title="dialogTitle()"
    >
      <el-form :model="tableForm" ref="tableFormRef">
        <el-form-item
          style="display: none"
          v-if="dialogType === 'edit'"
          label="编号"
          :label-width="140"
        >
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>
        <el-form-item
          label="角色名称"
          :label-width="140"
          prop="name"
          :rules="[
            { required: true, trigger: 'blur', message: '请输入角色名称' },
          ]"
        >
          <el-input v-model="tableForm.name" autocomplete="off" />
        </el-form-item>
        <el-form-item
          label="角色标识"
          :label-width="140"
          prop="uniqueKey"
          :rules="[
            { required: true, trigger: 'blur', message: '请输入角色标识' },
          ]"
        >
          <el-input v-model="tableForm.uniqueKey" autocomplete="off" />
        </el-form-item>
        <el-form-item label="分配权限" :label-width="140">
          <el-scrollbar
            height="200px"
            style="border: 1px solid #dddddd; border-radius: 5px; width: 100%"
          >
            <el-tree
              style="max-width: 600px"
              ref="treeRef"
              :data="tableMenuData"
              node-key="id"
              :props="defaultProps"
              :default-checked-keys="tableForm.permMenuIds"
              show-checkbox
              highlight-current
            />
          </el-scrollbar>
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

        <el-form-item label="排序值" :label-width="140">
          <el-input-number
            v-model.number="tableForm.orderNum"
            autocomplete="off"
          />
        </el-form-item>

        <el-form-item label="备注" :label-width="140" prop="remark">
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

const { proxy } = getCurrentInstance();

import tableDataApi from "@/api/system/sys_role.js";
import tableMenuDataApi from "@/api/system/sys_menu.js";
import { ElMessage } from "element-plus";
import sysTableApi from "@/api/feat/hx_order.js";

const defaultProps = {
  children: "children",
  label: "label",
};

let treeRef = $ref();
let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableMenuData = $ref([]);
let tableForm = $ref({
  type: 0,
  status: 1,
  parentId: 0,
});
let dialogFormVisible = $ref(false);
let handleAddIsClick = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);

const onSearchSubmit = async () => {
  tableDataApi.search(tableSearchForm).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
    }
  });
};

const getTableDataList = async () => {
  tableDataApi.listPage({}).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
      getMenuDataList();
    }
  });
};
getTableDataList();

const getMenuDataList = async () => {
  tableMenuDataApi.listPage().then((res) => {
    if (res.code === 200) {
      tableMenuData = res.data.list;
    }
  });
};

// 删除一条
const RowDel = async (id) => {
  await tableDataApi.delete(id);
};

const handleRowDel = async ({ id }) => {
  await tableDataApi.delete(id);
  ElMessage({
    message: "删除成功",
    type: "success",
    plain: true,
  });
  await getTableDataList();
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

// 编辑
const handleEdit = (row) => {
  dialogFormVisible = true;
  tableForm = {
    permMenuIds: [],
  };
  if (handleAddIsClick) {
    treeRef.setCheckedKeys([], false);
  } else {
    handleAddIsClick = true;
  }
  tableForm = { ...row };
  dialogType = "edit";
};

// 新增
const handleAdd = () => {
  dialogFormVisible = true;
  handleAddIsClick = true;
  tableForm = {
    status: 1,
    orderNum: 0,
    permMenuIds: [],
  };
  dialogType = "add";
  nextTick(() => {
    treeRef.setCheckedKeys([], false);
  });
};

const handleDialogClose = () => {
  tableForm.permMenuIds = [];
  tableForm.permMenuIdsAll = [];
};

const handleDialogOpen = () => {
  nextTick(() => {
    proxy.$refs.tableFormRef.clearValidate();
  });
};

const dialogTitle = () => {
  let title = "新增";
  if (dialogType !== "add") {
    title = "编辑";
    if (tableForm.name) {
      title += " #" + tableForm.name;
    } else if (tableForm.id) {
      title += " #" + tableForm.id;
    }
  }
  return title;
};
const dialogConfirm = async () => {
  tableForm.permMenuIds = treeRef.getCheckedKeys(false);
  tableForm.permMenuIdsAll = treeRef.getCheckedKeys(false);
  let tmps = treeRef.getHalfCheckedKeys(true);
  tableForm.permMenuIdsAll.push(...tmps);

  if (dialogType === "add") {
    // 添加数据
    tableForm.parentId = 0;
    proxy.$refs.tableFormRef.validate((valid) => {
      if (valid) {
        tableDataApi.add(tableForm).then((res) => {
          if (res.code === 200) {
            ElMessage({
              message: "创建成功",
              type: "success",
              plain: true,
            });
            dialogFormVisible = false;
            getTableDataList();
          }
        });
      }
    });
  } else {
    // 修改 内容
    proxy.$refs.tableFormRef.validate((valid) => {
      if (valid) {
        tableDataApi.update(tableForm).then((res) => {
          if (res.code === 200) {
            ElMessage({
              message: "更新成功",
              type: "success",
              plain: true,
            });
            dialogFormVisible = false;
            getTableDataList();
          }
        });
      }
    });
  }
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

const getStatusLabel = (idx) => {
  const index = statusOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return statusOptions[index].label;
  }
  return "";
};
</script>

<style scoped></style>
