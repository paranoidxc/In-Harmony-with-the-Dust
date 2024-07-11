<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="部门简称">
          <el-input
            v-model="tableSearchForm.name"
            clearable
            placeholder="请输入部门简称"
          ></el-input>
        </el-form-item>
        <!--
        <el-form-item label="部门全称">
          <el-input
            v-model="tableSearchForm.fullName"
            clearable
            placeholder="请输入部门全称"
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
        -->

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
      <el-table-column fixed prop="name" label="部门简称" width="100" />
      <el-table-column prop="fullName" label="部门全称" width="100" />
      <el-table-column prop="uniqueKey" label="部门标识" width="100" />
      <el-table-column prop="orderNum" label="排序值" width="100" />
      <el-table-column prop="parentId" label="上级ID" width="80" />
      <el-table-column prop="status" label="状态" width="60">
        <template #default="scope">
          {{ getStatusLabel(scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column prop="type" label="类型" width="100">
        <template #default="scope">
          <el-tag>{{ getTypeLabel(scope.row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" width="" />
      <el-table-column fixed="right" label="操作" width="160">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">
            编 辑
          </el-button>
          <el-button
            size="small"
            type="danger"
            @click="handleRowDel(scope.row)"
          >
            删 除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      background
      style="display: flex; justify-content: right; margin-top: 10px"
      v-model:current-page="curPage"
      v-model:page-size="limit"
      :page-sizes="[limit, 20, 50, 100, 200, 300, 400, 500]"
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
      draggable
      v-model="dialogFormVisible"
      @open="handleDialogOpen"
      :title="dialogType === 'add' ? '新增' : '编辑'"
    >
      <el-form ref="tableFormRef" :model="tableForm"
      :validate-on-rule-change="false"
      >
        <el-form-item
          style="display: none"
          v-if="dialogType === 'edit'"
          label="编号"
          :label-width="80"
        >
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>
        <el-form-item
          label="部门简称"
          :label-width="80"
          prop="name"
          :rules="[
            { required: true, trigger: 'blur', message: '请输入部门简称' },
          ]"
        >
          <el-input v-model="tableForm.name" autocomplete="off" />
        </el-form-item>
        <el-form-item
          label="部门全称"
          :label-width="80"
          prop="fullName"
          :rules="[
            { required: true, trigger: 'blur', message: '请输入部门全称' },
          ]"
        >
          <el-input v-model="tableForm.fullName" autocomplete="off" />
        </el-form-item>
        <el-form-item
          label="部门标识"
          :label-width="80"
          prop="uniqueKey"
          :rules="[
            { required: true, trigger: 'blur', message: '请输入部门标识' },
          ]"
        >
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

        <el-form-item label="类型" :label-width="80"
            prop="type"
            :rules="[
                { required: true, trigger: 'change', message: '请输入部门类型' },
            ]"
        >
          <el-select
            v-model.number="tableForm.type"
            placeholder="请选择"
          >
            <el-option
              v-for="item in typeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item
          label="备注"
          :label-width="80"
          prop="remark"
          :rules="[{ required: true, trigger: 'blur', message: '请输入备注' }]"
        >
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

import sysDeptApi from "@/api/system/sys_dept.js";

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
const typeOptions = [
  {
    value: 1,
    label: "公司",
  },
  {
    value: 2,
    label: "子公司",
  },
  {
    value: 3,
    label: "部门",
  },
];

const getStatusLabel = (idx) => {
  if (idx == 0) {
    return "禁用";
  } else if (idx == 1) {
    return "启用";
  }
};
const getTypeLabel = (idx) => {
  const index = typeOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return typeOptions[index].label;
  } else {
  }
};

const getTableDataList = async (cur = 1, limit = 10) => {
  let result = await sysDeptApi.listPage({ page: cur, limit: limit });
  if (result.code == 200) {
    tableData = result.data.list;
    curPage = result.data.page;
    total = result.data.total;
  }
};
getTableDataList(1, limit);
/* 请求分页 */

const handleSizeChange = (val) => {
  limit = val;
  getTableDataList(curPage, val);
};

const handleCurrentChange = (val) => {
  getTableDataList(val, limit);
};

// 删除一条
const handleRowDel = async ({ id }) => {
  await sysDeptApi.delete({ id: id });
  await getTableDataList(curPage);
};

const handleDelList = () => {
  multipleSelection.forEach((id) => {
    handleRowDel({ id });
  });
  multipleSelection = [];
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
  dialogType = "edit";
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

const handleDialogOpen = () => {
    nextTick( () => {
        proxy.$refs.tableFormRef.clearValidate()
    })
}

//查询
const onSearchSubmit = async () => {
  tableSearchForm.page = 1;
  tableSearchForm.limit = limit;

  sysDeptApi.search(tableSearchForm).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
      curPage = res.data.page;
      total = res.data.total;
    }
  });
};

// 确认
const dialogConfirm = async () => {
  if (dialogType === "add") {
    // 添加数据
    tableForm.parentId = 0;
    tableForm.orderNum = 100;
    proxy.$refs.tableFormRef.validate((valid) => {
      if (valid) {
        sysDeptApi
          .add(tableForm)
          .then((res) => {
            if (res.code == 200) {
              dialogFormVisible = false;
              getTableDataList(curPage);
            }
          })
          .catch(() => {
            // console.log("CCCC");
          });
      }
    });
  } else {
    // 修改 内容
    proxy.$refs.tableFormRef.validate((valid) => {
      if (valid) {
        sysDeptApi.update(tableForm).then((res) => {
          if (res.code == 200) {
            dialogFormVisible = false;
            getTableDataList(curPage);
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
