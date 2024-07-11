<template>
  <el-row :gutter="20" style="margin: 10px">
    <el-col :span="6">
      <div class="grid-content ep-bg-purple">
        <el-table
          :data="tableData"
          :row-class-name="tableRowClassName"
          @row-click="getDictTableDataList"
          style="width: 100%"
          height="460"
          border
        >
          <el-table-column prop="name" label="字典项名称">
            <template #header>
              <div style="display: flex; justify-content: space-between">
                <span>字典项名称</span>
                <el-button
                  @click="handleAdd"
                  size="small"
                  type="primary"
                  :icon="Plus"
                  circle
                />
              </div>
            </template>
            <template #default="scope">
              {{ scope.row.name }}
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-col>

    <el-col :span="18">
      <div class="grid-content ep-bg-purple">
        <div v-if="isMainDictSel" class="btn-list" style="margin-bottom: 10px">
          <el-button type="primary" @click="handleSubDictAdd">
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

        <el-table
          ref="multipleTableRef"
          @selection-change="handleSelectionChange"
          :data="tableDictData"
          height="460"
          stripe
          border
          style="width: 100%"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="id" label="编号" width="60" />
          <el-table-column prop="name" label="字典项名称" width="120" />
          <el-table-column prop="uniqueKey" label="字典项标识" width="120" />
          <el-table-column prop="type" label="字典项类型" width="100">
            <template #default="scope">
              {{ getTypeLabel(scope.row.type) }}
            </template>
          </el-table-column>
          <el-table-column prop="value" label="字典项值" width="100" />
          <el-table-column prop="orderNum" label="排序值" width="80" />
          <el-table-column prop="status" label="状态" width="60">
            <template #default="scope">
              {{ getStatusLabel(scope.row.status) }}
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="140">
            <template #default="scope">
              <el-button size="small" @click="handleSubDictEdit(scope.row)">
                编 辑
              </el-button>
              <el-popconfirm
                title="确定要删除么?"
                @confirm="handleRowDel(scope.row)"
              >
                <template #reference>
                  <el-button size="small" type="danger">删 除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-col>
  </el-row>

  <!-- dialog -->
  <el-dialog
    draggable
    v-model="dialogFormVisible"
    :title="dialogType === 'add' ? '新增' : '编辑'"
  >
    <el-form :model="tableForm">
      <el-form-item
        style="display: none"
        v-if="dialogType === 'edit'"
        label="编号"
        :label-width="120"
      >
        <el-input v-model="tableForm.id" autocomplete="off" />
      </el-form-item>
      <el-form-item label="名称" :label-width="120">
        <el-input v-model="tableForm.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="标识" :label-width="120">
        <el-input v-model="tableForm.uniqueKey" autocomplete="off" />
      </el-form-item>
      <el-form-item v-if="selDictId" label="字典项类型" :label-width="120">
        <el-select v-model.number="tableForm.type" placeholder="Select">
          <el-option
            v-for="item in typeOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item v-if="selDictId" label="字典项值" :label-width="120">
        <el-input v-model="tableForm.value" autocomplete="off" />
      </el-form-item>
      <el-form-item label="状态" :label-width="120">
        <el-tooltip :content="getStatusLabel(tableForm.status)" placement="top">
          <el-switch
            v-model.number="tableForm.status"
            class="mt-2"
            inline-prompt
            :active-value="1"
            :inactive-value="0"
          />
        </el-tooltip>
      </el-form-item>

      <el-form-item label="排序值" :label-width="120">
        <el-input-number
          v-model.number="tableForm.orderNum"
          autocomplete="off"
        />
      </el-form-item>

      <el-form-item label="备注" :label-width="120">
        <el-input v-model="tableForm.remark" autocomplete="off" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="dialogConfirm"> 确 认 </el-button>
      </span>
    </template>
  </el-dialog>
</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";

import {
  Check,
  Plus,
  Delete,
  Edit,
  Message,
  Search,
  Star,
} from "@element-plus/icons-vue";

const { proxy } = getCurrentInstance();

import tableDataApi from "@/api/system/sys_dict.js";
import tableTypeDataApi from "@/api/system/sys_dict_type.js";

let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableDictData = $ref([]); // 表格数据
let tableForm = $ref({});
let dialogFormVisible = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(2);
let total = $ref(15);
let curPage = $ref(1);

let selDictId = $ref(0);
let isMainDictSel = $ref(false);

const dropdown1 = $ref();

function showClick() {
  if (!dropdown1.value) return;
  dropdown1.value.handleOpen();
}

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

const tableRowClassName = ({ row, rowIndex }) => {
  //console.log("row", row);
  //console.log("row", row.id);
  //console.log("rowIndex", rowIndex);
  if (row.id === selDictId) {
    return "success-row";
  }
  return "";
};

const getTableDataList = async () => {
  let result = await tableDataApi.list({});
  if (result.code == 200) {
    tableData = result.data.list;
  }
};
getTableDataList();

const getDictTableDataList = async (row) => {
  selDictId = row.id;
  isMainDictSel = true;
  let result = await tableDataApi.listPage({
    parentId: row.id,
    page: 1,
    limit: 1000000,
  });
  if (result.code == 200) {
    tableDictData = result.data.list;
  }
};
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
  await tableDataApi.delete({ id: id }).then((res) => {
    if (res.code == 200) {
      if (selDictId != 0) {
        getDictTableDataList(selDictId);
      }
    }
  });
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

const handleSubDictEdit = (row) => {
  dialogFormVisible = true;
  dialogType = "edit";
  tableForm = { ...row };
};

// 新增
const handleAdd = () => {
  selDictId = 0;
  dialogFormVisible = true;
  tableForm = {
    orderNum: 0,
    status: 1,
  };
  dialogType = "add";
};

const handleSubDictAdd = () => {
  dialogFormVisible = true;
  tableForm = {
    orderNum: 0,
    status: 1,
  };
  dialogType = "add";
};

//查询

// 确认
const dialogConfirm = async () => {
  if (dialogType === "add") {
    if (selDictId == 0) {
      tableForm.parentId = 0;
      tableForm.value = "";
    } else {
      tableForm.parentId = selDictId;
    }
    tableForm.type = 1;
    // 添加数据
    tableDataApi.add(tableForm).then((res) => {
      if (res.code == 200) {
        dialogFormVisible = false;
        if (selDictId == 0) {
          getTableDataList();
        } else {
          getDictTableDataList(selDictId);
        }
      }
    });
  } else {
    // 修改 内容
    if (selDictId == 0) {
      tableForm.parentId = 0;
      tableForm.value = "";
    } else {
      tableForm.parentId = selDictId;
    }
    tableForm.type = 1;
    tableDataApi.update(tableForm).then((res) => {
      if (res.code == 200) {
        dialogFormVisible = false;
        if (selDictId == 0) {
          getTableDataList();
        } else {
          getDictTableDataList(selDictId);
        }
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
  } else {
  }
};
const typeOptions = [
  {
    value: 1,
    label: "文本",
  },
  {
    value: 2,
    label: "数字",
  },
  {
    value: 3,
    label: "数组",
  },
  {
    value: 4,
    label: "单选",
  },
  {
    value: 5,
    label: "多选",
  },
  {
    value: 6,
    label: "下拉",
  },
  {
    value: 7,
    label: "日期",
  },
  {
    value: 8,
    label: "时间",
  },
  {
    value: 9,
    label: "9单图",
  },
  {
    value: 10,
    label: "多图",
  },
  {
    value: 11,
    label: "单文件",
  },

  {
    value: 12,
    label: "多文件",
  },
];
const getTypeLabel = (idx) => {
  const index = typeOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return typeOptions[index].label;
  } else {
  }
};
</script>

<style scoped></style>
