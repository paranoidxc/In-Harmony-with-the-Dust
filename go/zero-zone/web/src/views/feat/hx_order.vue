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

        <el-form-item label="核销单号">
          <el-input v-model="tableSearchForm.no" placeholder="" clearable />
        </el-form-item>

        <el-form-item label="渠道订单号">
          <el-input
            v-model="tableSearchForm.qdOrderId"
            placeholder=""
            clearable
          />
        </el-form-item>

        <el-form-item label="渠道券号">
          <el-input
            v-model="tableSearchForm.certificateId"
            placeholder=""
            clearable
          />
        </el-form-item>

        <el-form-item label="核销时间">
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

        <el-form-item label="核销状态">
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

        <el-form-item label="是否撤销">
          <el-select
            v-model.number="tableSearchForm.isUnverify"
            clearable
            style="width: 100px"
          >
            <el-option
              v-for="item in yesNoOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
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
      <el-table-column prop="updatedAt" label="核销时间" width="150" />
      <el-table-column prop="no" label="核销单号" width="" />
      <el-table-column prop="cooperateShopId" label="商户" width="150">
        <template #default="scope">
          {{ getAuthNameOptLabel(scope.row.cooperateShopId) }}
        </template>
      </el-table-column>
      <el-table-column prop="cooperateShopId" label="门店" width="150">
        <template #default="scope">
          {{ getOptLabel(shopOptions, scope.row.cooperateShopId) }}
        </template>
      </el-table-column>
      <el-table-column prop="openNo" label="门店订单号 " width="170" />
      <el-table-column prop="typo" label="渠道" width="60">
        <template #default="scope">
          {{ getOptLabel(typoOptions, scope.row.typo) }}
        </template>
      </el-table-column>
      <el-table-column prop="qdOrderId" label="渠道订单号" width="190" />
      <el-table-column prop="qdGoodName" label="渠道商品名称" width="200" />
      <el-table-column prop="qdPrice" label="渠道金额" width="90" />
      <el-table-column prop="status" label="核销状态" width="90">
        <template #default="scope">
          {{ getOptLabel(statusOptions, scope.row.status) }}
        </template>
      </el-table-column>
      <el-table-column prop="isUnverify" label="是否撤销" width="90">
        <template #default="scope">
          {{ getOptLabel(yesNoOptions, scope.row.isUnverify) }}
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
      :page-sizes="[10, 20, 50, 100, 200, 300, 400, 500]"
      :small="small"
      :disabled="disabled"
      :background="background"
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

            <el-form-item label="核销单号" :label-width="100">
              <el-input v-model="tableForm.no" />
            </el-form-item>

            <el-form-item label="渠道" :label-width="100">
              <el-input :value="getOptLabel(typoOptions, tableForm.typo)" />
            </el-form-item>

            <el-form-item label="渠道订单号" :label-width="100">
              <el-input v-model="tableForm.qdOrderId" />
            </el-form-item>

            <el-form-item label="渠道券号" :label-width="100">
              <el-input :value="getQdIdentCode(tableForm)" />
            </el-form-item>

            <el-form-item label="渠道交易金额" :label-width="100">
              <el-input v-model="tableForm.qdPrice" />
            </el-form-item>

            <el-form-item label="渠道商品名称" :label-width="100">
              <el-input v-model="tableForm.qdGoodName" />
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

            <el-form-item label="订单时间" :label-width="100">
              <el-input v-model="tableForm.createdAt" />
            </el-form-item>

            <el-form-item label="确认时间" :label-width="100">
              <el-input v-model="tableForm.qdNotifyDatetime" />
            </el-form-item>

            <el-form-item label="通知时间" :label-width="100">
              <el-input v-model="tableForm.openRespDatetime" />
            </el-form-item>

            <el-form-item label="核销状态" :label-width="100">
              <el-input :value="getOptLabel(statusOptions, tableForm.status)" />
            </el-form-item>

            <el-form-item label="是否撤销" :label-width="100">
              <el-input
                :value="getOptLabel(yesNoOptions, tableForm.isUnverify)"
              />
            </el-form-item>

            <el-form-item label="撤销单号" :label-width="100">
              <el-input v-model="tableForm.uhxOrderNo" />
            </el-form-item>

            <el-form-item label="渠道参数" :label-width="100">
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

        <el-form-item label="类型" :label-width="80">
          <el-select v-model.number="tableForm.type" placeholder="Select">
            <el-option
              v-for="item in typeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
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

    <el-drawer v-model="drawer" title="日志信息" size="50%">
      <div>
        <el-text class="mx-1" type="primary"
          >账号：{{ tableForm.name }}
        </el-text>
      </div>
      <div>
        <el-text class="mx-1" type="primary"
          >操作时间{{ tableForm.createTime }}
        </el-text>
      </div>
      <div>
        <el-text class="mx-1" type="primary">{{ tableForm.request }}</el-text>
      </div>
    </el-drawer>
  </div>
</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";
import Axios from "axios";
import store from "@/store";

const { proxy } = getCurrentInstance();

import tableDataApi from "@/api/system/sys_log.js";
import sysTableApi from "@/api/feat/hx_order.js";
import sysTableShopApi from "@/api/feat/cooperate_shop.js";
import sysTableAuthApi from "@/api/feat/saas_cooperate_auth.js";

let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableForm = $ref({});
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
const typeOptions = [];

let tableDataShop = $ref([]);
let shopOptions = $ref([]); // 渠道商户

let tableDataAuth = $ref([]);
let authOptions = $ref([]); // 授权商户

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
  if (tableSearchForm.isUnverify === undefined) {
    tableSearchForm.isUnverify = -1;
  }
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
  download("/vsapi/admin/feat/hxOrder/export", tableSearchForm);
};

const onSearchSubmit = async () => {
  tableSearchForm.page = 1;
  tableSearchForm.limit = limit;
  if (tableSearchForm.isUnverify === undefined) {
    tableSearchForm.isUnverify = -1;
  }
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
  let result = await sysTableApi.page(tableSearchForm);
  if (result.code == 200) {
    tableData = result.data.list;
    curPage = result.data.pagination.page;
    total = result.data.pagination.total;
  }
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
const getTableDataList = async (cur = 1, limit) => {
  let cp = {};
  cp = tableSearchForm;
  if (cp.isUnverify === undefined) {
    cp.isUnverify = -1;
  }
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

const handleView = async (row) => {
  dialogViewFormVisible = true;
  let result = await sysTableApi.detail(row.id);
  if (result.code == 200) {
    tableForm = { ...result.data };
    tableForm.uhxOrderNo = result.data.uhxOrder.no;
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
const yesNoOptions = [
  {
    value: -1,
    label: "全部",
  },
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
      return authOptions[tmpId].label;
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

const getQdIdentCode = (detail) => {
  if (detail.typo == 1) {
    return detail.certificateId;
  } else if (detail.typo == 2) {
    return detail.identCode;
  }
  return "";
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

<style scoped></style>
