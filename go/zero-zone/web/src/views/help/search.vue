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

        <el-form-item label="日期范围">
          <el-date-picker
            v-model="tableSearchForm.dateRange"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            type="daterange"
            unlink-panels
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :disabled-date="disabledDateFun"
            @calendar-Change="calendarChange"
            :shortcuts="shortcuts"
          />
        </el-form-item>

        <el-form-item label="选择框">
          <el-checkbox v-model="tableSearchForm.checkBox" label="启用" />
        </el-form-item>

        <el-form-item label="多选">
          <el-select
            v-model="tableSearchForm.mulSel"
            multiple
            clearable
            collapse-tags
            collapse-tags-tooltip
            :max-collapse-tags="3"
            placeholder="全部"
            style="width: 240px"
          >
            <el-option
              v-for="item in mulSelOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="xx多选">
          <el-select
            v-model="value"
            multiple
            clearable
            collapse-tags
            placeholder="Select"
            popper-class="custom-header"
            :max-collapse-tags="1"
            style="width: 240px"
          >
            <template #header>
              <el-checkbox
                v-model="checkAll"
                :indeterminate="indeterminate"
                @change="handleCheckAll"
              >
                All
              </el-checkbox>
            </template>
            <el-option
              v-for="item in cities"
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

const checkAll = ref(false);
const indeterminate = ref(false);
const value = ref([]);

const mulSelOptions = [
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
  {
    value: 4,
    label: "4部门",
  },
  {
    value: 5,
    label: "5部门",
  },
];

const cities = ref([
  {
    value: "Beijing",
    label: "Beijing",
  },
  {
    value: "Shanghai",
    label: "Shanghai",
  },
  {
    value: "Nanjing",
    label: "Nanjing",
  },
  {
    value: "Chengdu",
    label: "Chengdu",
  },
  {
    value: "Shenzhen",
    label: "Shenzhen",
  },
  {
    value: "Guangzhou",
    label: "Guangzhou",
  },
]);
watch(value, (val) => {
  console.log("Watch");
  if (val.length === 0) {
    checkAll.value = false;
    indeterminate.value = false;
  } else if (val.length === cities.value.length) {
    checkAll.value = true;
    indeterminate.value = false;
  } else {
    indeterminate.value = true;
  }
});
const handleCheckAll = (val) => {
  console.log("val", val);
  indeterminate.value = false;
  if (val) {
    value.value = cities.value.map((_) => _.value);
  } else {
    value.value = [];
  }
};

const shortcuts = [
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

const choiceDate = $ref({});
const disabledDateFun = (time) => {
  // 如何选择了一个日期
  if (choiceDate.value) {
    // 7天的时间戳
    const one = 6 * 24 * 3600 * 1000;
    // 当前日期 - one = 7天之前
    const minTime = choiceDate.value - one;
    // 当前日期 + one = 7天之后
    const maxTime = choiceDate.value + one;
    return time.getTime() < minTime || time.getTime() > maxTime;
    // 限制不能选择今天及以后
    // || time.getTime() + 1 * 24 * 3600 * 1000 > Date.now()
  } else {
    //默认只能选择今天以及今天之前的日期
    //return time.getTime() > Date.now();
    //return time.getTime() >= new Date().getTime();
    // 如果没有选择日期，就要限制不能选择今天及以后
    // return time.getTime() + 1 * 24 * 3600 * 1000 > Date.now();
  }
};

const calendarChange = (obj) => {
  const minDate = obj[0];
  const maxDate = obj[1];
  choiceDate.value = minDate.getTime();
  if (maxDate) {
    choiceDate.value = null;
  }
};

const getStatusLabel = (idx) => {
  if (idx == 0) {
    return "禁用";
  } else if (idx == 1) {
    return "启用";
  }
};

//查询
const onSearchSubmit = async () => {
  console.log("tableSearchForm", tableSearchForm);
  sysDeptApi.search(tableSearchForm).then((res) => {
    tableSearchForm.dtStart = tableSearchForm.dateRange[0];
    tableSearchForm.dtEnd = tableSearchForm.dateRange[1];
    if (res.code === 200) {
    }
  });
};

// 确认
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
