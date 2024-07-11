<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="门店名称">
          <el-input v-model="tableSearchForm.name" placeholder="" clearable />
        </el-form-item>

        <el-form-item label="商户号">
          <el-input
            v-model="tableSearchForm.account"
            placeholder=""
            clearable
          />
        </el-form-item>

        <el-form-item label="渠道商户">
          <el-input v-model="tableSearchForm.qdName" placeholder="" clearable />
        </el-form-item>

        <el-form-item label="渠道">
          <el-select
            v-model.number="tableSearchForm.typo"
            clearable
            placeholder="全部"
            style="width: 100px"
          >
            <el-option
              v-for="item in typoOptions"
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
      height="600"
      stripe
      style="width: 100%"
      border
    >
      <el-table-column fixed type="selection" width="55" />
      <el-table-column fixed prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="门店名称" width="" />
      <el-table-column prop="account" label="商户号" width="200" />
      <el-table-column prop="apiKey" label="商户KEY" width="180" />
      <el-table-column prop="saasCooperateAuthId" label="渠道商户" width="140">
        <template #default="scope">
          {{ getOptLabel(authAllOptions, scope.row.saasCooperateAuthId) }}
        </template>
      </el-table-column>
      <el-table-column prop="typo" label="渠道" width="100">
        <template #default="scope">
          {{ getTypoLabel(scope.row.typo) }}
        </template>
      </el-table-column>

      <el-table-column prop="shopId" label="渠道商户门店" width="160">
        <template #default="scope">
          {{ getOptLabel(shopAllOptions, scope.row.shopId) }}
        </template>
      </el-table-column>

      <el-table-column label="状态" width="80">
        <template #default="scope">
          {{ getOptLabel(statusOptions, scope.row.status) }}
        </template>
      </el-table-column>

      <el-table-column prop="updatedAt" label="更新时间" width="150" />
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
        label-width="auto"
      >
        <el-form-item
          style="display: none"
          v-if="dialogType === 'edit'"
          label="编号"
        >
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>

        <el-form-item label="商户号" prop="account" label-width="auto">
          <el-input v-model="tableForm.account" placeholder="" clearable />
        </el-form-item>

        <el-form-item
          v-if="dialogType === 'edit'"
          label="接口key"
          prop="apiKey"
          label-width="auto"
        >
          <el-input
            disabled
            v-model="tableForm.apiKey"
            placeholder=""
            clearable
          />
        </el-form-item>

        <el-form-item label="门店信息" prop="name" label-width="auto">
          <el-input v-model="tableForm.name" placeholder="" clearable />
        </el-form-item>

        <el-form-item label="渠道" prop="typo" label-width="auto">
          <el-select
            v-if="dialogType === 'edit'"
            disabled
            v-model.number="tableForm.typo"
            placeholder="请选择渠道"
            @change="handleChangeTypo"
          >
            <el-option
              v-for="item in typoOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <el-select
            v-if="dialogType !== 'edit'"
            v-model.number="tableForm.typo"
            placeholder="请选择渠道"
            @change="handleChangeTypo"
          >
            <el-option
              v-for="item in typoOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          label="渠道商户"
          prop="saasCooperateAuthId"
          label-width="auto"
        >
          <el-select
            v-model.number="tableForm.saasCooperateAuthId"
            placeholder="请选择渠道商户"
            @change="handleChangeCooperateAuth"
          >
            <el-option
              v-for="item in authOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="渠道商户门店" prop="shopId" label-width="auto">
          <el-select
            v-model="tableForm.shopId"
            placeholder="请选择渠道商户门店"
          >
            <el-option
              v-for="item in shopOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="状态" prop="shopId" label-width="auto">
          <el-select v-model="tableForm.status" placeholder="请选择状态">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
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

import sysTableApi from "@/api/feat/cooperate_shop.js";
import sysAuthApi from "@/api/feat/saas_cooperate_auth.js";

// 授权列表
let authListData = $ref([]); // 表格数据
let authOptions = $ref([]); // 渠道商户
let authAllOptions = $ref([]); // 所有渠道商户
let shopOptions = $ref([]); //渠道商户门店
let shopAllOptions = $ref([]); // 所有渠道商户门店

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

const small = ref(false);
const background = ref(false);
const disabled = ref(false);
const rules = $ref({
  account: [{ required: true, message: "请输入 商户号信息", trigger: "blur" }],
  name: [{ required: true, message: "请输入 门店信息", trigger: "blur" }],
  typo: [{ required: true, message: "请输入 渠道", trigger: "blur" }],
  saasCooperateAuthId: [
    { required: true, message: "请输入 渠道商户", trigger: "blur" },
  ],
  shopId: [{ required: true, message: "请输入 渠道商户门店", trigger: "blur" }],
});

// 取得所有的授权列表
const getAuthDataList = async () => {
  let res = await sysAuthApi.all({ includeDeleted: 1 });
  if (res.code == 200) {
    authListData = res.data.list;
    authAllOptions = [];
    shopAllOptions = [];
    authListData.forEach(function (ele) {
      authAllOptions.push({
        value: ele.id,
        label: ele.name,
      });

      if (ele.syncContent.length) {
        let shopJson = JSON.parse(ele.syncContent);
        if (ele.typo == 1) {
          shopJson.data.pois.forEach(function (poi) {
            shopAllOptions.push({
              value: poi.poi.poi_id,
              label: poi.poi.poi_name,
            });
          });
        } else if (ele.typo == 2) {
          shopJson.data.forEach(function (poi) {
            shopAllOptions.push({
              value: poi.open_shop_uuid,
              label: poi.shopname,
            });
          });
        }
      }
    });
  }
};
getAuthDataList();

// 切换渠道
const handleChangeTypo = function (value) {
  authOptions = [];
  shopOptions = [];
  tableForm.saasCooperateAuthId = "";
  tableForm.shopId = "";
  if (value == 1) {
    authListData.forEach(function (ele) {
      if (ele.deletedAt.length == 0) {
        if (ele.typo == value) {
          authOptions.push({
            value: ele.id,
            label: ele.name,
          });
        }
      }
    });
  } else if (value == 2) {
    authListData.forEach(function (ele) {
      if (ele.deletedAt.length == 0) {
        if (ele.typo == value) {
          authOptions.push({
            value: ele.id,
            label: ele.name,
          });
        }
      }
    });
  }
};

const handleChangeCooperateAuth = function (value) {
  shopOptions = [];
  authListData.forEach(function (ele) {
    if (ele.deletedAt.length == 0) {
      if (value == ele.id) {
        let shopJson = JSON.parse(ele.syncContent);
        if (ele.typo == 1) {
          shopJson.data.pois.forEach(function (poi) {
            shopOptions.push({
              value: poi.poi.poi_id,
              label: poi.poi.poi_name,
            });
          });
        } else if (ele.typo == 2) {
          shopJson.data.forEach(function (item) {
            shopOptions.push({
              value: item.open_shop_uuid,
              label: item.shopname,
            });
          });
        }
      }
    }
  });
};

//查询
const onSearchSubmit = async () => {
  tableSearchForm.page = 1;
  tableSearchForm.limit = limit;
  let result = await sysTableApi.page(tableSearchForm);
  if (result.code == 200) {
    tableData = result.data.list;
    curPage = result.data.pagination.page;
    total = result.data.pagination.total;
  }
};
const getTableDataList = async (cur = 1, limit) => {
  let cp = {};
  cp = tableSearchForm;
  cp.page = cur;
  cp.limit = limit;

  let result = await sysTableApi.page(cp);
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

// 编辑
const handleEdit = async (row) => {
  dialogFormVisible = true;
  dialogType = "edit";

  let result = await sysTableApi.detail(row.id);
  if (result.code == 200) {
    tableForm = { ...result.data };
    handleChangeTypo(result.data.typo);
    handleChangeCooperateAuth(result.data.saasCooperateAuthId);
    tableForm.saasCooperateAuthId = result.data.saasCooperateAuthId;
    tableForm.shopId = result.data.shopId;
  }
};

// 新增
const handleCreate = () => {
  dialogFormVisible = true;
  tableForm = {
    status: 1,
  };
  authOptions = [];
  shopOptions = [];
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

const statusOptions = [
  {
    value: 1,
    label: "正常",
  },
  {
    value: 2,
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

const typoOptions = [
  {
    value: 1,
    label: "抖音",
  },
  {
    value: 2,
    label: "美团",
  },
];

const getTypoLabel = (idx) => {
  const index = typoOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return typoOptions[index].label;
  } else {
  }
};

const getOptLabel = (options, idx) => {
  const index = options.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return options[index].label;
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
