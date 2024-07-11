<template>
  1111

</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";

const { proxy } = getCurrentInstance();

import tableDataApi from "@/api/system/sys_log.js";

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
let total = $ref(15);
let curPage = $ref(1);
let drawer = $ref(false);

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
const typeOptions = [];

const getStatusLabel = (idx) => {
  if (idx == 0) {
    return "禁用";
  } else if (idx == 1) {
    return "启用";
  }
};

const getTypeLabel = (idx) => {};

const getTableDataList = async (cur = 1, limit) => {
  let result = await tableDataApi.listPage({
    page: cur,
    limit: limit,
  });
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
// 选中
const handleSelectionChange = (val) => {
  multipleSelection = [];
  val.forEach((item) => {
    multipleSelection.push(item.id);
  });
};

// 编辑
const handleEdit = (row) => {
  //dialogFormVisible = true;
  //dialogType = "edit";
  drawer = true;
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

// 确认
const dialogConfirm = async () => {
  dialogFormVisible = false;
};
</script>

<style scoped></style>
