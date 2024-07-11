<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="部门简称">
          <el-input
            v-model="tableSearchForm.name"
            placeholder="请输入部门简称"
          ></el-input>
        </el-form-item>
        <el-form-item label="部门全称">
          <el-input
            v-model="tableSearchForm.fullName"
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

        <el-form-item>
          <el-button type="primary" @click="onSearchSubmit">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- query -->
    <div class="query-box">
      <el-input
        class="query-input"
        v-model="queryInput"
        placeholder="请输入姓名搜索🔍"
        @change="handleQueryName"
      />

      <div class="btn-list">
        <el-button
          type="danger"
          @click="handleDelList"
          v-if="multipleSelection.length > 0"
          >删除多选
        </el-button>
        <el-button type="primary" @click="handleAdd">增加</el-button>
      </div>
    </div>

    <el-table
      ref="multipleTableRef"
      @selection-change="handleSelectionChange"
      :data="tableData"
      height="250"
      stripe
      style="width: 100%"
      border
    >
      <el-table-column fixed type="selection" width="55" />
      <el-table-column fixed prop="id" label="  编号" width="60" />
      <el-table-column fixed prop="name" label="部门简称" width="100" />
      <el-table-column prop="fullName" label="部门全称" width="100" />
      <el-table-column prop="uniqueKey" label="唯一值" width="100" />
      <el-table-column prop="orderNum" label="排序值" width="100" />
      <el-table-column prop="parentId" label="上级ID" width="80" />
      <el-table-column prop="status" label="状态" width="60">
        <template #default="scope">
          {{ getStatusLabel(scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column prop="type" label="类型" width="80">
        <template #default="scope">
          <el-tag>{{ getTypeLabel(scope.row.type) }}</el-tag>
          <!--<el-button size="small">
                      {{ getTypeLabel(scope.row.type) }}
                    </el-button>-->
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" width="400" />
      <el-table-column fixed="right" label="操作" width="160">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">
            编辑
          </el-button>
          <el-button
            size="small"
            type="danger"
            @click="handleRowDel(scope.row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!--<el-pagination
                                      background
                                      layout="prev, pager, next"
                                      style="display: flex; justify-content: center; margin-top: 10px"
                                      :total="total"
                                      :page-size="limit"
                                      v-model:current-page="curPage"
                                      @current-change="handleChangePage"
                                    />-->

    <el-pagination
      background
      style="display: flex; justify-content: center; margin-top: 10px"
      v-model:current-page="curPage"
      v-model:page-size="limit"
      :page-sizes="[limit, 10, 20, 50, 100, 200, 300, 400, 500]"
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
          <el-select
            v-model.number="tableForm.status"
            placeholder="Select"
            size="large"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" :label-width="80">
          <el-input v-model.number="tableForm.type" autocomplete="off" />
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
  </div>
</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";
import sysDeptApi from "@/api/system/sys_dept.js";

const { proxy } = getCurrentInstance();

let queryInput = $ref("");
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
let limit = $ref(2);
let total = $ref(15);
let curPage = $ref(1);

//const currentPage4 = ref(1);
//const pageSize4 = ref(2);

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

const getDeptList = async (cur = 1, limit = 2) => {
  let result = await sysDeptApi.listPage({ page: cur, limit: limit });
  console.log(result);
  tableData = result.data.list;
  curPage = result.data.page;
  total = result.data.total;
};
getDeptList();
/* 请求分页 */

const handleChangePage = (val) => {
  getDeptList(curPage);
};

const handleSizeChange = (val) => {
  limit = val;
  getDeptList(curPage, val);
};

const handleCurrentChange = (val) => {
  getDeptList(val, limit);
};

// 删除一条
const handleRowDel = async ({ id }) => {
  await sysDeptApi.delete({ id: id });
  await getDeptList(curPage);
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

// 搜索
const handleQueryName = async (val) => {
  if (val.length > 0) {
    let query = {
      name: val,
    };
    let result = await sysDeptApi.search(query);
    tableData = result.data.list;
  } else {
    getDeptList(curPage);
  }
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

//查询
const onSearchSubmit = async () => {
  sysDeptApi.search(tableSearchForm).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
    }
  });
};

// 确认
const dialogConfirm = async () => {
  if (dialogType === "add") {
    // 添加数据
    tableForm.parentId = 0;
    tableForm.orderNum = 100;
    sysDeptApi
      .add(tableForm)
      .then((res) => {
        if (res.code == 200) {
          dialogFormVisible = false;
          getDeptList(curPage);
        }
      })
      .catch(() => {
        // console.log("CCCC");
      });
  } else {
    // 修改 内容
    sysDeptApi.update(tableForm).then((res) => {
      if (res.code == 200) {
        dialogFormVisible = false;
        getDeptList(curPage);
      }
    });
  }
};
</script>

<style scoped>
.table-box {
  margin: 200px auto;
  width: 800px;
}

.title {
  text-align: center;
}

.query-box {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.query-input {
  width: 200px;
}

.demo-pagination-block + .demo-pagination-block {
  margin-top: 10px;
}

.demo-pagination-block .demonstration {
  margin-bottom: 16px;
}
</style>
