<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="商户">
          <el-input v-model="tableSearchForm.shName" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="门店">
          <el-input
            v-model="tableSearchForm.shopName"
            placeholder=""
            clearable
          />
        </el-form-item>

        <el-form-item label="订单号">
          <el-input
            v-model="tableSearchForm.qdOrderId"
            placeholder=""
            clearable
          />
        </el-form-item>

        <el-form-item label="撤销单号">
          <el-input v-model="tableSearchForm.no" placeholder="" clearable />
        </el-form-item>

        <el-form-item label="撤销时间">
          <el-date-picker
            v-model="tableSearchForm.dateRangeTmp"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            type="daterange"
            unlink-panels
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :shortcuts="shortcuts"
          />
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

        <el-form-item label="结果">
          <el-select
            v-model.number="tableSearchForm.status"
            clearable
            placeholder="全部"
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
          <el-button type="primary" @click="onExportSubmit">
            <el-icon class="el-icon--left">
              <Download />
            </el-icon>
            导出
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- query -->
    <div class="query-box">
      <div class="btn-list"></div>
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
      <!--<el-table-column fixed prop="id" label="ID" width="60" />-->
      <el-table-column prop="updatedAt" label="撤销时间" width="150" />
      <el-table-column prop="no" label="撤销单号" width="" />
      <el-table-column label="商户" width="150">
        <template #default="scope">
          {{ getAuthNameOptLabel(scope.row.cooperateShopId) }}
        </template>
      </el-table-column>
      <el-table-column label="门店" width="150">
        <template #default="scope">
          {{ getOptLabel(shopOptions, scope.row.cooperateShopId) }}
        </template>
      </el-table-column>
      <el-table-column prop="typo" label="渠道" width="60">
        <template #default="scope">
          {{ getOptLabel(typoOptions, scope.row.typo) }}
        </template>
      </el-table-column>
      <el-table-column prop="qdGoodName" label="渠道商品名称" width="200" />
      <el-table-column prop="identCode" label="渠道券号" width="160" />
      <el-table-column prop="qdNotifyDatetime" label="通知时间" width="150" />
      <el-table-column prop="status" label="结果" width="80">
        <template #default="scope">
          {{ getOptLabel(statusOptions, scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" width="100">
        <template #default="scope">
          <el-button size="small" @click="handleView(scope.row)">
            查 看
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
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />

    <el-dialog
      draggable
      @open="handleDialogOpen"
      v-model="dialogViewFormVisible"
      destroy-on-close
      title="查看"
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

            <el-form-item label="门店代码" :label-width="100">
              <el-input
                :value="getShopKeyOptLabel(tableForm.cooperateShopId)"
              />
            </el-form-item>

            <el-form-item label="门店订单号" :label-width="100">
              <el-input v-model="tableForm.openNo" autocomplete="off" />
            </el-form-item>

            <el-form-item label="撤销单号" prop="no" :label-width="100">
              <el-input v-model="tableForm.no" />
            </el-form-item>

            <el-form-item label="渠道" :label-width="100">
              <el-input :value="getOptLabel(typoOptions, tableForm.typo)" />
            </el-form-item>

            <el-form-item label="渠道券号" prop="identCode" :label-width="100">
              <el-input v-model="tableForm.identCode" />
            </el-form-item>

            <el-form-item label="渠道单号" prop="qdOrderId" :label-width="100">
              <el-input v-model="tableForm.qdOrderId" />
            </el-form-item>

            <el-form-item label="核销单号" :label-width="100">
              <el-input v-model="tableForm.hxNo" />
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item
              label="商户名称"
              prop="cooperateShopId"
              :label-width="100"
            >
              <el-input
                :value="getAuthNameOptLabel(tableForm.cooperateShopId)"
              />
            </el-form-item>

            <el-form-item
              label="门店名称"
              prop="cooperateShopId"
              :label-width="100"
            >
              <el-input
                :value="getOptLabel(shopOptions, tableForm.cooperateShopId)"
              />
            </el-form-item>

            <el-form-item label="订单时间" prop="createdAt" :label-width="100">
              <el-input v-model="tableForm.createdAt" />
            </el-form-item>

            <el-form-item label="确认时间" prop="updatedAt" :label-width="100">
              <el-input v-model="tableForm.qdNotifyDatetime" />
            </el-form-item>

            <el-form-item label="通知时间" prop="updatedAt" :label-width="100">
              <el-input v-model="tableForm.openRespDatetime" />
            </el-form-item>

            <el-form-item label="撤销状态" :label-width="100">
              <el-input :value="getOptLabel(statusOptions, tableForm.status)" />
            </el-form-item>

            <el-form-item label="渠道参数" prop="content" :label-width="100">
              <el-input
                :value="formatJson(tableForm.content)"
                type="textarea"
                rows="10"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogViewFormVisible = false">
            关 闭
          </el-button>
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
          label-width="auto"
        >
          <el-input v-model="tableForm.id" autocomplete="off" />
        </el-form-item>
        <el-form-item label="撤销单号" prop="no" :label-width="100">
          <el-input v-model="tableForm.no" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="渠道" prop="typo" :label-width="100">
          <el-input v-model="tableForm.typo" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="门店ID" prop="cooperateShopId" :label-width="100">
          <el-input
            v-model="tableForm.cooperateShopId"
            placeholder=""
            clearable
          />
        </el-form-item>
        <el-form-item label="第三方返回信息" prop="content" :label-width="100">
          <el-input v-model="tableForm.content" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="平台订单" prop="openNo" :label-width="100">
          <el-input v-model="tableForm.openNo" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="状态" prop="status" :label-width="100">
          <el-input v-model="tableForm.status" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="券号" prop="identCode" :label-width="100">
          <el-input v-model="tableForm.identCode" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="渠道订单号" prop="qdOrderId" :label-width="100">
          <el-input v-model="tableForm.qdOrderId" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="渠道商品名称" prop="qdGoodName" :label-width="100">
          <el-input v-model="tableForm.qdGoodName" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="渠道金额" prop="qdPrice" :label-width="100">
          <el-input v-model="tableForm.qdPrice" placeholder="" clearable />
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
import Axios from "axios";
import store from "@/store";

const { proxy } = getCurrentInstance();

import sysTableApi from "@/api/feat/uhx_order.js";
import sysTableShopApi from "@/api/feat/cooperate_shop.js";
import sysTableAuthApi from "@/api/feat/saas_cooperate_auth.js";

let tableDataAuth = $ref([]);
let authOptions = $ref([]); // 授权商户

let tableDataShop = $ref([]);
let shopOptions = $ref([]); // 渠道商户

let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableForm = $ref({
  empty: "",
});
let dialogFormVisible = $ref(false);
let dialogViewFormVisible = $ref(false);
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(10);
let total = $ref(0);
let curPage = $ref(1);
let drawer = $ref(false);

const small = ref(false);
const background = ref(false);
const disabled = ref(false);
const rules = $ref({});

const download = (url, query = {}) => {
  const { isLogin, tokenObj } = toRefs(store.user.useUserStore());
  const queryArgs = {
    url: url,
    method: "get",
    params: query,
    responseType: "blob",
    headers: {
      Authorization: tokenObj.value.tokenValue,
      Accept: "application/json",
      "Content-Type": "application/json; charset=utf-8",
      withCredentials: true,
    },
  };
  return Axios.request(queryArgs)
    .then((res) => {
      const fileName =
        res.headers["content-disposition"].match(/filename=(.*)/)[1];
      const content = res.data;
      const blob = new Blob([content]);
      if ("download" in document.createElement("a")) {
        // 非IE下载
        const elink = document.createElement("a");
        elink.download = decodeURIComponent(fileName);
        elink.style.display = "none";
        elink.href = URL.createObjectURL(blob);
        document.body.appendChild(elink);
        elink.click();
        URL.revokeObjectURL(elink.href); // 释放URL 对象
        document.body.removeChild(elink);
      } else {
        // IE10+下载
        navigator.msSaveBlob(blob, fileName);
      }
    })
    .catch((err) => console.log(err));
};

const onExportSubmit = async () => {
  tableSearchForm.dateRangeStart = "";
  tableSearchForm.dateRangeEnd = "";
  if (tableSearchForm.dateRangeTmp !== undefined) {
    if (
      Array.isArray(tableSearchForm.dateRangeTmp) &&
      tableSearchForm.dateRangeTmp.length > 0
    ) {
      tableSearchForm.dateRangeStart = tableSearchForm.dateRangeTmp[0];
      tableSearchForm.dateRangeEnd = tableSearchForm.dateRangeTmp[1];
    }
  }
  download("/vsapi/admin/feat/uhxOrder/export", tableSearchForm);
};
//查询
const onSearchSubmit = async () => {
  tableSearchForm.page = 1;
  tableSearchForm.limit = limit;

  tableSearchForm.dateRangeStart = "";
  tableSearchForm.dateRangeEnd = "";
  if (tableSearchForm.dateRangeTmp !== undefined) {
    if (
      Array.isArray(tableSearchForm.dateRangeTmp) &&
      tableSearchForm.dateRangeTmp.length > 0
    ) {
      tableSearchForm.dateRangeStart = tableSearchForm.dateRangeTmp[0];
      tableSearchForm.dateRangeEnd = tableSearchForm.dateRangeTmp[1];
    }
  }
  sysTableApi.page(tableSearchForm).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
      curPage = res.data.pagination.page;
      total = res.data.pagination.total;
    }
  });
};

//取的授权列表
const getTableDataAuthList = async () => {
  let res = await sysTableAuthApi.all({ includeDeleted: 1 });
  if (res.code == 200) {
    tableDataAuth = res.data.list;
    tableDataAuth.forEach(function (ele) {
      authOptions.push({
        value: ele.id,
        label: ele.name,
      });
    });
  }
};
getTableDataAuthList();

//取的门店数据
const getTableDataShopList = async () => {
  let res = await sysTableShopApi.all({ includeDeleted: 1 });
  if (res.code == 200) {
    tableDataShop = res.data.list;
    tableDataShop.forEach(function (ele) {
      shopOptions.push({
        value: ele.id,
        label: ele.name,
      });
    });
  }
};
getTableDataShopList();

/* 请求分页 */
const getTableDataList = async (cur, limit) => {
  let cp = {};
  cp = tableSearchForm;
  cp.page = cur;
  cp.limit = limit;
  cp.dateRangeStart = "";
  cp.dateRangeEnd = "";
  if (cp.dateRangeTmp !== undefined) {
    if (Array.isArray(cp.dateRangeTmp) && cp.dateRangeTmp.length > 0) {
      cp.dateRangeStart = cp.dateRangeTmp[0];
      cp.dateRangeEnd = cp.dateRangeTmp[1];
    }
  }

  let result = await sysTableApi.page(cp);
  if (result.code == 200) {
    tableData = result.data.list;
    curPage = result.data.pagination.page;
    total = result.data.pagination.total;
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
    label: "成功",
  },
  {
    value: 2,
    label: "失败",
  },
  {
    value: 3,
    label: "异常",
  },
];

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

const getStatusLabel = (idx) => {
  const index = statusOptions.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return statusOptions[index].label;
  } else {
  }
};

const getAuthNameOptLabel = (idx) => {
  if (idx > 0) {
    let authId = 0;
    tableDataShop.forEach((ele) => {
      if (ele.id == idx) {
        authId = ele.saasCooperateAuthId;
        return;
      }
    });
    if (authId > 0) {
      let tmpId = authOptions.findIndex((option) => option.value === authId);
      return authOptions[tmpId] ? authOptions[tmpId].label : "";
    }
    return "";
  }
};

const getShopKeyOptLabel = (idx) => {
  let apiKey = "";
  tableDataShop.forEach((ele) => {
    if (ele.id == idx) {
      apiKey = ele.apiKey;
      return;
    }
  });
  return apiKey;
};

const getOptLabel = (options, idx) => {
  const index = options.findIndex((option) => option.value === idx);
  if (index !== -1) {
    return options[index].label;
  }
};

const formatJson = (jsonStr) => {
  try {
    return JSON.stringify(JSON.parse(jsonStr), null, 4);
  } catch (e) {}
};
const shortcuts = [
  {
    text: "今天",
    value: [new Date(), new Date()],
  },
  {
    text: "昨天",
    value: () => {
      const date = new Date();
      date.setTime(date.getTime() - 3600 * 1000 * 24);
      return [date, date];
    },
  },
  {
    text: "最近一周",
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
      return [start, end];
    },
  },
  {
    text: "最近一个月",
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
      return [start, end];
    },
  },
  {
    text: "最近三个月",
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
      return [start, end];
    },
  },
];
</script>

<style scoped>
.query-box {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.el-row {
  margin-bottom: 20px;
}

.el-row:last-child {
  margin-bottom: 0;
}

.el-col {
  border-radius: 4px;
}

.grid-content {
  border-radius: 4px;
  min-height: 36px;
}
</style>
