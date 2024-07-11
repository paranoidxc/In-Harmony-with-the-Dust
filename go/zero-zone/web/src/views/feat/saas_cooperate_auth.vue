<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="渠道商户">
          <el-input v-model="tableSearchForm.name" placeholder="" clearable />
        </el-form-item>

        <el-form-item label="渠道">
          <el-select
            v-model.number="tableSearchForm.typo"
            clearable
            placeholder="全部"
            style="width: 100px"
          >
            <el-option
              v-for="item in mixin.data().qdTypoOptions"
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
      <el-table-column fixed prop="id" label="ID" width="60" />
      <el-table-column prop="devConfId" label="服务商平台" width="">
        <template #default="scope">
          {{ filters.optLabelOrName(devConfOptions, scope.row.devConfId) }}
        </template>
      </el-table-column>
      <el-table-column prop="accountId" label="accountID" width="240" />
      <el-table-column prop="name" label="渠道商户" width="200" />
      <el-table-column prop="typo" label="渠道" width="100">
        <template #default="scope">
          {{ filters.typoName(scope.row.typo) }}
        </template>
      </el-table-column>
      <!--
      <el-table-column prop="content" label="授权信息" width="" />
      <el-table-column prop="syncContent" label="同步信息" width="" />
      -->
      <el-table-column prop="updatedAt" label="更新时间" width="150" />
      <el-table-column prop="status" label="状态" width="150">
        <template #default="scope">
          {{ filters.optLabelOrName(statusOptions, scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" width="200">
        <template #default="scope">
          <el-button size="small" @click="handleView(scope.row)">
            编 辑
          </el-button>

          <el-popconfirm
            title="同步后，可能会更新门店信息?"
            @confirm="handleSync(scope.row)"
          >
            <template #reference>
              <el-button size="small"> 同步</el-button>
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

    <el-dialog
      draggable
      @open="handleDialogOpen"
      v-model="dialogViewFormVisible"
      :title="dialogType === 'create' ? '新增' : '编辑'"
      destroy-on-close
    >
      <el-form :model="tableForm" ref="tableFormRef" label-position="right">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item
              style="display: none"
              v-if="dialogType === 'edit'"
              label="编号"
              label-width="auto"
            >
              <el-input v-model="tableForm.id" autocomplete="off" />
            </el-form-item>
            <el-form-item label="服务商平台" :label-width="100">
              <el-input
                :value="
                  filters.optLabelOrName(devConfOptions, tableForm.devConfId)
                "
              />
            </el-form-item>
            <el-form-item label="accountId" :label-width="100">
              <el-input :value="tableForm.accountId" />
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="渠道商户" :label-width="100">
              <el-input v-model="tableForm.name" />
              编辑只可修改 渠道商户
            </el-form-item>

            <el-form-item label="渠道" :label-width="100">
              <el-input :value="filters.typoName(tableForm.typo)" />
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item label="授权信息" :label-width="100">
              <el-input
                :value="formatJson(tableForm.content)"
                type="textarea"
                rows="10"
              />
            </el-form-item>

            <el-form-item label="同步信息" :label-width="100">
              <el-input
                :value="formatJson(tableForm.syncContent)"
                type="textarea"
                rows="10"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <!--
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogViewFormVisible = false">
            关 闭
          </el-button>
        </span>
        -->
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogConfirm"> 确 认 </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- dialog -->
    <el-dialog
      draggable
      @open="handleDialogOpen"
      v-model="dialogFormVisible"
      :title="dialogType === 'create' ? '新增' : '编辑'"
    >
      <el-form
        :model="tableForm"
        ref="tableFormRef"
        :rules="rules"
        label-position="right"
      >
        <el-form-item
          style="display: none"
          v-if="dialogType === 'edit'"
          label="编号"
          :label-width="100"
        >
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>
        <el-form-item label="DEV应用ID" prop="devConfId" :label-width="100">
          <el-input v-model="tableForm.devConfId" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="授权信息" prop="content" :label-width="100">
          <el-input v-model="tableForm.content" placeholder="" clearable />
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
import { ElMessage, ElMessageBox } from "element-plus";

const { proxy } = getCurrentInstance();

import mixin from "@/utils/mixin.js";
import { filters } from "@/utils/filters.js";

import sysTableApi from "@/api/feat/saas_cooperate_auth.js";
import sysTableDevConfApi from "@/api/feat/third_part_dev_conf.js";

let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableDataDevConf = $ref([]);
let devConfOptions = $ref([]); // 渠道商户
let tableForm = $ref({
  status: 1,
});
let dialogFormVisible = $ref(false);
let dialogViewFormVisible = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(10);
let total = $ref(0);
let curPage = $ref(1);

const rules = $ref({});

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

const getTableDataDevConfList = async () => {
  let res = await sysTableDevConfApi.all({ includeDeleted: 1 });
  devConfOptions = [];
  if (res.code == 200) {
    tableDataDevConf = res.data.list;
    tableDataDevConf.forEach(function (ele) {
      devConfOptions.push({
        value: ele.id,
        label: ele.name,
      });
    });
  }
};
getTableDataDevConfList();

/* 请求分页 */
const getTableDataList = async (cur, limit) => {
  let cp = {};
  cp = tableSearchForm;
  cp.page = cur;
  cp.limit = limit;
  let res = await sysTableApi.page(cp);
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

const handleDelList = async () => {
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

const handleView = async (row) => {
  dialogViewFormVisible = true;
  let result = await sysTableApi.detail(row.id);
  if (result.code == 200) {
    tableForm = { ...result.data };
  }
};

// 同步信息
const handleSync = async (row) => {
  let result = await sysTableApi.syncAuth(row.id);
  if (result.code == 200) {
    getTableDataList(curPage, limit);
    ElMessage({
      message: "操作成功",
      type: "success",
      plain: true,
    });
  }
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
    /*
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
     */
  } else {
    // 修改 内容
    proxy.$refs.tableFormRef.validate((valid) => {
      if (valid) {
        sysTableApi.update(tableForm).then((res) => {
          if (res.code == 200) {
            //dialogFormVisible = false;
            dialogViewFormVisible = false;
            ElMessage({
              message: "操作成功",
              type: "success",
              plain: true,
            });
            getTableDataList(curPage, limit);
          }
        });
      }
    });
  }
};

const statusOptions = [
  {
    value: 0,
    label: "已授权",
  },
  {
    value: 1,
    label: "异常",
  },
];

const formatJson = (jsonStr) => {
  try {
    return JSON.stringify(JSON.parse(jsonStr), null, 4);
  } catch (e) {}
};
</script>

<style scoped>
.query-box {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}
</style>
