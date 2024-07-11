<template>
  <div style="margin: 20px">
    <div class="">
      <el-form
        ref="tableSearchFormRef"
        :model="tableSearchForm"
        inline
        :rules="rules"
      >
        <el-form-item label="下单时间" prop="dateRangeTmp">
          <el-date-picker
            v-model="tableSearchForm.dateRangeTmp"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            type="daterange"
            unlink-panels
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-value="thisMonthRange"
            :shortcuts="shortcuts"
          />
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

    <div>
      <el-card style="">
        <template #header>
          <div class="card-header">
            <h1>数据统计</h1>
          </div>
        </template>
        <el-row>
          <el-col :span="12">
            <el-statistic title="核销数量" :value="tableData.hxOrderCnt" />
          </el-col>
          <el-col :span="12">
            <el-statistic title="撤销数量" :value="tableData.uhxOrderCnt" />
          </el-col>
        </el-row>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { getCurrentInstance, proxyRefs } from "vue";

const { proxy } = getCurrentInstance();

import sysTableApi from "@/api/feat/dashboard.js";
import sysTableAuthApi from "@/api/feat/saas_cooperate_auth.js";

let myDate = new Date();
let weekAgoDate = new Date(myDate.getTime() - 3600 * 1000 * 24 * 7);

let dateStart =
  weekAgoDate.getFullYear() +
  "-" +
  ("0" + (weekAgoDate.getMonth() + 1)).slice(-2) +
  "-" +
  ("0" + weekAgoDate.getDate()).slice(-2);

let dateEnd =
  myDate.getFullYear() +
  "-" +
  ("0" + (myDate.getMonth() + 1)).slice(-2) +
  "-" +
  ("0" + myDate.getDate()).slice(-2);

let tableSearchForm = $ref({
  dateRangeTmp: [dateStart, dateEnd],
});
let tableData = $ref({}); // 表格数据
let thisMonthRange = $ref(); // 表格数据
const rules = $ref({
  dateRangeTmp: [
    { required: true, message: "请输入下单时间", trigger: "blur" },
  ],
});

const onSearchSubmit = async () => {
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

  try {
    proxy.$refs.tableSearchFormRef.validate((valid) => {
      if (valid) {
        sysTableApi
          .index(tableSearchForm)
          .then((result) => {
            if (result.code == 200) {
              tableData = result.data;
            }
          })
          .catch();
      }
    });
  } catch (e) {}
};

const getTableData = async () => {
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
  let result = await sysTableApi.index(tableSearchForm);
  if (result.code == 200) {
    tableData = result.data;
  }
};
getTableData();

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
<style lang="scss" scoped></style>
