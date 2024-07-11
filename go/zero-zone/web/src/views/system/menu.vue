<template>
  <div style="margin: 20px">
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
          增加根目录
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
      height="100%"
      stripe
      border
      style="width: 100%; margin-bottom: 20px"
      ref="multipleTableRef"
      row-key="id"
      :data="tableData"
      @selection-change="handleSelectionChange"
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="菜单名称" width="" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="scope">
          <el-tag>{{ getTypeLabel(scope.row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="icon" label="图标" width="60">
        <template #default="scope">
          <el-icon v-if="scope.row.icon">
            <component :is="scope.row.icon" />
          </el-icon>
        </template>
      </el-table-column>
      <el-table-column prop="keepAlive" label="页面缓存" width="90">
        <template #default="scope">
          {{ getYesNoLabel(scope.row.keepAlive) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="60">
        <template #default="scope">
          {{ getStatusLabel(scope.row.isShow) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="240">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">
            编辑
          </el-button>
          <el-button
            size="small"
            type="primary"
            @click="handleAddSub(scope.row)"
          >
            <el-icon class="el-icon--left">
              <Plus />
            </el-icon>
            增加子菜单
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- dialog -->
    <el-dialog
      draggable
      v-model="dialogFormVisible"
      :title="dialogType === 'add' ? '新增' : '编辑'"
    >
      <el-form :model="tableForm">
        <el-form-item
          class="dsN"
          v-if="dialogType === 'edit'"
          label="编号"
          :label-width="140"
        >
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>

        <el-form-item label="类型" :label-width="140">
          <el-radio-group v-model="tableForm.type" @change="handleSwitchType">
            <el-radio v-for="item in typeOptions" :value="item.value"
              >{{ item.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="菜单名称" :label-width="140">
          <el-input v-model="tableForm.name" autocomplete="off" />
        </el-form-item>

        <el-form-item
          v-if="isAddRoot"
          label="上级菜单"
          :label-width="140"
          style="display: none"
        >
          <el-input v-model.number="tableForm.parentId" autocomplete="off" />
        </el-form-item>

        <el-form-item v-if="isAddRoot" label="上级菜单" :label-width="140">
          <el-input value="根菜单" disabled autocomplete="off" />
        </el-form-item>

        <el-form-item
          v-if="!isAddRoot"
          label="上级菜单"
          :label-width="140"
          style="width: 100%"
        >
          <el-cascader
            style="width: 100%"
            v-model="tableForm.parentId"
            :options="parentOptions"
            :props="{ expandTrigger: 'hover', checkStrictly: true }"
            @change="handleParentChange"
          />
        </el-form-item>

        <el-form-item
          v-if="tableForm.type != 2"
          label="路由"
          :label-width="140"
        >
          <el-input v-model="tableForm.router" autocomplete="off" />
        </el-form-item>

        <el-form-item
          v-if="tableForm.type == 1"
          label="视图路径"
          :label-width="140"
        >
          <el-input v-model="tableForm.viewPath" autocomplete="off" />
        </el-form-item>

        <el-form-item
          v-if="tableForm.type != 2"
          label="图标"
          :label-width="140"
        >
          <el-input v-model="tableForm.icon" autocomplete="off" />
        </el-form-item>

        <el-form-item
          v-if="tableForm.type == 2"
          label="权限"
          :label-width="140"
        >
          <el-select
            v-model="tableForm.perms"
            placeholder="选择分配的权限"
            multiple
          >
            <el-option
              v-for="item in typeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.label"
            />
          </el-select>
        </el-form-item>

        <el-form-item
          v-if="tableForm.type != 2"
          label="KeepAlive"
          :label-width="140"
        >
          <el-tooltip
            :content="getStatusLabel(tableForm.keepAlive)"
            placement="top"
          >
            <el-switch
              v-model.number="tableForm.keepAlive"
              class="mt-2"
              inline-prompt
              :active-value="1"
              :inactive-value="0"
            />
          </el-tooltip>
        </el-form-item>

        <el-form-item
          v-if="tableForm.type != 2"
          label="状态"
          :label-width="140"
        >
          <el-tooltip
            :content="getStatusLabel(tableForm.isShow)"
            placement="top"
          >
            <el-switch
              v-model.number="tableForm.isShow"
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

        <el-form-item class="dsN" label="备注" :label-width="140">
          <el-input v-model="tableForm.remark" autocomplete="off" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogConfirm"> 确认 </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";

const { proxy } = getCurrentInstance();

import tableDataApi from "@/api/system/sys_menu.js";

let tableData = $ref([]); // 表格数据
let tableForm = $ref({
  type: 0,
  status: 1,
  parentId: 0,
});
let dialogFormVisible = $ref(false);
let isAddRoot = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(2);
let total = $ref(15);
let curPage = $ref(1);
let parentOptions = $ref([]);

const getTableDataList = async () => {
  let result = await tableDataApi.listPage();
  if (result.code == 200) {
    tableData = result.data.list;
    parentOptions = result.data.list;
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

// 删除一条
const handleRowDel = async ({ id }) => {
  await tableDataApi.delete({ id: id });
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
  if (tableForm.parentId === 0) {
    isAddRoot = true;
  } else {
    isAddRoot = false;
  }
};

// 新增
const handleAdd = () => {
  dialogFormVisible = true;
  isAddRoot = true;
  tableForm = {
    isShow: 1,
    parentId: 0,
    remark: "",
    type: 0,
    icon: "House",
    orderNum: 100,
    viewPath: "views/dashboard/index",
  };
  dialogType = "add";
};

const handleAddSub = (row) => {
  dialogFormVisible = true;
  isAddRoot = false;
  tableForm = {
    isShow: 1,
    parentId: row.id,
    remark: "",
    type: 1,
    icon: "House",
    orderNum: 100,
    viewPath: "views/dashboard/index",
  };
  dialogType = "add";
};

const handleParentChange = (value) => {
  console.log(value);
};

//查询

// 确认
const dialogConfirm = async () => {
  if (tableForm.type == 2) {
    tableForm.router = "";
  }
  if (tableForm.type != 2) {
    tableForm.perms = [];
  }

  if (dialogType === "add") {
    // 添加数据
    tableForm.activeRouter = "";

    if (!isAddRoot) {
      if (tableForm.parentId.length) {
        let parentIds = tableForm.parentId.toString().split(",");
        if (parentIds.length > 1) {
          tableForm.parentId = parseInt(parentIds[parentIds.length - 1]);
        }
      }
    }

    tableDataApi
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
  } else {
    // 修改 内容
    tableForm.activeRouter = "";
    tableDataApi.update(tableForm).then((res) => {
      if (res.code == 200) {
        dialogFormVisible = false;
        getTableDataList(curPage);
      }
    });
  }
};

const handleSwitchType = (val) => {
  console.log("VAL", val);
};

const parentOptionsTmp = [
  {
    value: 1,
    label: "Guide",
    children: [
      {
        value: 2,
        label: "Disciplines",
        children: [
          {
            value: 3,
            label: "Consistency",
          },
          {
            value: "feedback",
            label: "Feedback",
          },
          {
            value: "efficiency",
            label: "Efficiency",
          },
          {
            value: "controllability",
            label: "Controllability",
          },
        ],
      },
      {
        value: "navigation",
        label: "Navigation",
        children: [
          {
            value: "side nav",
            label: "Side Navigation",
          },
          {
            value: "top nav",
            label: "Top Navigation",
          },
        ],
      },
    ],
  },
  {
    value: "component",
    label: "Component",
    children: [
      {
        value: "basic",
        label: "Basic",
        children: [
          {
            value: "layout",
            label: "Layout",
          },
          {
            value: "color",
            label: "Color",
          },
          {
            value: "typography",
            label: "Typography",
          },
          {
            value: "icon",
            label: "Icon",
          },
          {
            value: "button",
            label: "Button",
          },
        ],
      },
      {
        value: "form",
        label: "Form",
        children: [
          {
            value: "radio",
            label: "Radio",
          },
          {
            value: "checkbox",
            label: "Checkbox",
          },
          {
            value: "input",
            label: "Input",
          },
          {
            value: "input-number",
            label: "InputNumber",
          },
          {
            value: "select",
            label: "Select",
          },
          {
            value: "cascader",
            label: "Cascader",
          },
          {
            value: "switch",
            label: "Switch",
          },
          {
            value: "slider",
            label: "Slider",
          },
          {
            value: "time-picker",
            label: "TimePicker",
          },
          {
            value: "date-picker",
            label: "DatePicker",
          },
          {
            value: "datetime-picker",
            label: "DateTimePicker",
          },
          {
            value: "upload",
            label: "Upload",
          },
          {
            value: "rate",
            label: "Rate",
          },
          {
            value: "form",
            label: "Form",
          },
        ],
      },
      {
        value: "data",
        label: "Data",
        children: [
          {
            value: "table",
            label: "Table",
          },
          {
            value: "tag",
            label: "Tag",
          },
          {
            value: "progress",
            label: "Progress",
          },
          {
            value: "tree",
            label: "Tree",
          },
          {
            value: "pagination",
            label: "Pagination",
          },
          {
            value: "badge",
            label: "Badge",
          },
        ],
      },
      {
        value: "notice",
        label: "Notice",
        children: [
          {
            value: "alert",
            label: "Alert",
          },
          {
            value: "loading",
            label: "Loading",
          },
          {
            value: "message",
            label: "Message",
          },
          {
            value: "message-box",
            label: "MessageBox",
          },
          {
            value: "notification",
            label: "Notification",
          },
        ],
      },
      {
        value: "navigation",
        label: "Navigation",
        children: [
          {
            value: "menu",
            label: "Menu",
          },
          {
            value: "tabs",
            label: "Tabs",
          },
          {
            value: "breadcrumb",
            label: "Breadcrumb",
          },
          {
            value: "dropdown",
            label: "Dropdown",
          },
          {
            value: "steps",
            label: "Steps",
          },
        ],
      },
      {
        value: "others",
        label: "Others",
        children: [
          {
            value: "dialog",
            label: "Dialog",
          },
          {
            value: "tooltip",
            label: "Tooltip",
          },
          {
            value: "popover",
            label: "Popover",
          },
          {
            value: "card",
            label: "Card",
          },
          {
            value: "carousel",
            label: "Carousel",
          },
          {
            value: "collapse",
            label: "Collapse",
          },
        ],
      },
    ],
  },
  {
    value: "resource",
    label: "Resource",
    children: [
      {
        value: "axure",
        label: "Axure Components",
      },
      {
        value: "sketch",
        label: "Sketch Templates",
      },
      {
        value: "docs",
        label: "Design Documentation",
      },
    ],
  },
];

const yesNoOptions = [
  {
    value: 1,
    label: "是",
  },
  {
    value: 0,
    label: "否",
  },
];
const statusOptions = [
  {
    value: 1,
    label: "显示",
  },
  {
    value: 0,
    label: "隐藏",
  },
];
const typeOptions = [
  {
    value: 0,
    label: "目录",
  },
  {
    value: 1,
    label: "菜单",
  },
  {
    value: 2,
    label: "权限",
  },
];

const getYesNoLabel = (idx) => {
  const index = yesNoOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return yesNoOptions[index].label;
  } else {
  }
};

const getStatusLabel = (idx) => {
  const index = statusOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return statusOptions[index].label;
  } else {
  }
};
const getTypeLabel = (idx) => {
  const index = typeOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return typeOptions[index].label;
  } else {
  }
};
</script>

<style scoped></style>
