<template>
  <div style="margin: 20px">
    <div class="">
      <el-form :model="tableSearchForm" inline>
        <el-form-item label="应用名称">
          <el-input v-model="tableSearchForm.name" placeholder="" clearable />
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
      <el-table-column prop="name" label="应用名称" width="" />
      <el-table-column prop="appId" label="AppID" width="200" />
      <el-table-column prop="appSecret" label="AppSecret" width="200" />
      <el-table-column prop="createdAt" label="创建时间" width="150" />
      <el-table-column prop="updatedAt" label="更新时间" width="150" />
      <el-table-column prop="typo" label="类型" width="100">
        <template #default="scope">
          {{ getOptLabel(typoOptions, scope.row.typo) }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" width="120">
        <template #default="scope">
          <div class="flex flex-wrap items-center">
            <el-dropdown :hide-on-click="false">
              <el-button type="primary">
                操 作
                <el-icon class="el-icon--right">
                  <arrow-down />
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleAuthUrl(scope.row)"
                    >授权地址
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleQrCode(scope.row)"
                    >授权二维码
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleEdit(scope.row)"
                    >编 辑
                  </el-dropdown-item>
                  <el-dropdown-item
                    divided
                    @click="handleRefreshClientToken(scope.row)"
                    >刷新ClientToken
                  </el-dropdown-item>
                  <el-dropdown-item divided>
                    <el-popconfirm
                      title="确定要删除么?"
                      @confirm="handleRowDel(scope.row)"
                    >
                      <template #reference>
                        <el-text class="mx-1" type="danger" style="width: 100%"
                          >删 除
                        </el-text>
                      </template>
                    </el-popconfirm>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
        <!--

      <template #default="scope">
        <el-button size="small" @click="handleAuthUrl(scope.row)">
          授权地址
        </el-button>
        <el-button size="small" @click="handleQrCode(scope.row)">
          授权二维码
        </el-button>
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
      </template> -->
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
        <el-form-item label="应用名称" prop="name" :label-width="100">
          <el-input v-model="tableForm.name" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="AppID" prop="appId" :label-width="100">
          <el-input v-model="tableForm.appId" placeholder="" clearable />
        </el-form-item>
        <el-form-item label="AppSecret" prop="appSecret" :label-width="100">
          <el-input v-model="tableForm.appSecret" placeholder="" clearable />
        </el-form-item>

        <el-form-item label="类型" prop="typo" :label-width="100">
          <el-select
            :disabled="dialogType === 'edit'"
            v-model.number="tableForm.typo"
            placeholder="请选择类型"
          >
            <el-option
              v-for="item in typoOptions"
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

    <el-dialog
      v-model="dialogVisibleUrl"
      :title="dialogAuthUrlTitle"
      width="500"
    >
      <span>{{ authUrl }}</span>
    </el-dialog>

    <el-dialog
      v-model="dialogVisibleUrlButton"
      :title="dialogAuthUrlTitle"
      width="500"
    >
      <span>{{ authUrl }}</span>
      <!-- <button @click="gotoauth(authUrl)">去授权</button> -->
    </el-dialog>
  </div>
</template>
<script setup>
import { getCurrentInstance, proxyRefs } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

const { proxy } = getCurrentInstance();

import sysTableApi from "@/api/feat/third_part_dev_conf.js";

let tableSearchForm = $ref({});
let tableData = $ref([]); // 表格数据
let tableForm = $ref({
  status: 1,
});
let dialogFormVisible = $ref(false);
let dialogVisibleUrl = $ref(false);
let dialogVisibleUrlButton = $ref(false);
let authUrl = $ref("");
let dialogType = $ref("add");
let multipleSelection = $ref([]);
let limit = $ref(10);
let total = $ref(0);
let curPage = $ref(1);
let dialogAuthUrlTitle = $ref("");

const rules = $ref({
  name: [{ required: true, message: "请输入 应用名称", trigger: "blur" }],
  appId: [{ required: true, message: "请输入 AppID", trigger: "blur" }],
  appSecret: [{ required: true, message: "请输入 AppSecret", trigger: "blur" }],
  typo: [{ required: true, message: "请输入 类型", trigger: "blur" }],
});

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

const handleAuthUrl = async (row) => {
  let result = await sysTableApi.authUrl(row.id);
  if (result.code == 200) {
    //console.log("======");
    //console.log(result);
    authUrl = result.data.authUrl;
    dialogAuthUrlTitle = row.name + "授权地址";
    if (row.typo == 2) {
      dialogVisibleUrlButton = true;
    } else {
      dialogVisibleUrl = true;
    }
    /*
    ElMessageBox.alert(
        '<div style="text-align: center; width: 400px; height: 200px;"><textarea style="width: 100%; height: 100%;" class="el-textarea__inner">'+result.data.authUrl+'</textarea></div>',
        //'<div style="text-align: center; width: 400px; height: 200px;" v-html="url">'+result.data.authUrl+'</div>',
        '授权二维码',
        {
          dangerouslyUseHTMLString: false,
          center: true,
        }
    )
     */
  }
};

const handleQrCode = async (row) => {
  let result = await sysTableApi.qrCode(row.id);
  if (result.code == 200) {
    ElMessageBox.confirm(
      '<div style="text-align: center"><img src="' +
        result.data.qrCode +
        '" /></div>',
      row.name + "授权二维码",
      {
        confirmButtonText: "确定",
        dangerouslyUseHTMLString: true,
        center: true,
        showCancelButton: false,
        //cancelButtonText: "Cancel",
        //type: "info",
      },
    )
      .then(() => {})
      .catch(() => {});

    /*
    ElMessageBox.alert(
      '<div style="text-align: center"><img src="' +
        result.data.qrCode +
        '" /></div>',
      row.name + "授权二维码",
      {
        dangerouslyUseHTMLString: true,
        cancelButtonText: "Cancel",
        center: true,
      },
    ).then(() => {});
     */
  }
};

// 涮新token
const handleRefreshClientToken = async (row) => {
  let result = await sysTableApi.refreshClientToken(row.id);
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
// 去授权
const gotoauth = async (authUrl) => {
  // proxy.$router.push({ path: '/auth-page', query: { authUrl: authUrl } });
  window.open(`#/open-auth?authUrl=` + encodeURIComponent(authUrl), "_blank");
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
