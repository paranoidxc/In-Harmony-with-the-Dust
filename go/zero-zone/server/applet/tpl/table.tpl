<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        {{- range $i, $v := .MinVueFields }}
        <el-form-item label="{{ $v.Label }}">
            <el-input v-model="tableSearchForm.{{ $v.Key }}" placeholder="" clearable />
        </el-form-item>
        {{- end }}
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
      <el-table-column fixed prop="{{ .PrimaryKeyJson}}" label="ID" width="60" />
      {{- range $i, $v := .VueFields }}
      <el-table-column prop="{{ $v.Key}}" label="{{ $v.Label}}" width="100" />
      {{- end }}
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
      :page-sizes="[limit, 20, 50, 100, 200, 300, 400, 500]"
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
      <el-form :model="tableForm"
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

        {{- range $i, $v := .MinVueFields }}
        <el-form-item label="{{ $v.Label }}" prop="{{ $v.Key }}" :label-width="100" >
            <el-input v-model="tableForm.{{ $v.Key }}" placeholder="" clearable />
        </el-form-item>
        {{- end }}

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

import sysTableApi from "@/api/{{ .DirectoryLowerCase }}/{{ .UnderlineName }}.js";

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

const rules = $ref({
    /*
    {{- range $i, $v := .MinVueFields }}
    {{ $v.Key }}: [
        { required: true, message: '请输入 {{ $v.Label }}', trigger: 'blur' },
    ],
    {{- end }}
    */
})

//查询
const onSearchSubmit = async () => {
  tableSearchForm.page = 1;
  tableSearchForm.limit = limit;

  sysTableApi.page(tableSearchForm).then((res) => {
    if (res.code === 200) {
      tableData = res.data.list;
      curPage = res.data.pagination.page;
      total = res.data.pagination.total;
    }
  });
};

/* 请求分页 */
const getTableDataList = async (cur, limit) => {
  let cpSearchForm =  tableSearchForm;
  cpSearchForm.page = cur;
  cpSearchForm.limit = limit;
  let res = await sysTableApi.page(cpSearchForm);
  if (res.code == 200) {
    tableData = res.data.list;
    curPage = res.data.pagination.page;
    total = res.data.pagination.total;
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
  await sysTableApi.delete(row.{{ .PrimaryKeyJson }});
  ElMessage({
    message: "删除成功",
    type: "success",
    plain: true,
  });
  await getTableDataList(curPage, limit);
};

const handleDelList = async() => {
  /*
  multipleSelection.forEach((id) => {
	reqRowDel(id)
  });
  */
  await sysTableApi.deletes(multipleSelection);
  multipleSelection = [];
  ElMessage({
    message: "批量删除成功",
    type: "success",
    plain: true,
  });
  await getTableDataList(curPage, limit);
};

// 选中
const handleSelectionChange = (val) => {
  multipleSelection = [];
  val.forEach((item) => {
    multipleSelection.push(item.{{ .PrimaryKeyJson}});
  });
};

// 编辑
const handleEdit = async (row) => {
  dialogFormVisible = true;
  dialogType = "edit";

  let result = await sysTableApi.detail(row.{{ .PrimaryKeyJson}});
  if (result.code == 200) {
    tableForm = { ...result.data };
  }
};

// 新增
const handleCreate = () => {
  dialogFormVisible = true;
  tableForm = {
  };
  dialogType = "create";
};

const handleDialogOpen = () => {
  nextTick( () => {
    proxy.$refs.tableFormRef.clearValidate()
  })
}

// 确认
const dialogConfirm = async () => {
  if (dialogType === "create") {
    // 添加数据
    proxy.$refs.tableFormRef.validate((valid) => {
        if (valid) {
            sysTableApi
              .create(tableForm)
              .then((res) => {
                if (res.code == 200) {
                    ElMessage({
                      message: "创建成功",
                      type: "success",
                      plain: true,
                    });
                    dialogFormVisible = false;
                    getTableDataList(curPage, limit);
                }
              })
       }
   });
  } else {
    // 修改 内容
    proxy.$refs.tableFormRef.validate((valid) => {
        if (valid) {
            sysTableApi.update(tableForm).then((res) => {
                if (res.code == 200) {
                    ElMessage({
                      message: "更新成功",
                      type: "success",
                      plain: true,
                    });
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