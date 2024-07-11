// 抽取公用的实例 - 操作成功与失败消息提醒内容等
export default {
  data() {
    return {
      qdTypoOptions: [
        {
          value: 1,
          label: "抖音",
        },
        {
          value: 2,
          label: "美团",
        },
      ],
      hxStatusOptions: [
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
      ],
      statusOptions: [
        {
          value: 1,
          label: "启用",
        },
        {
          value: 0,
          label: "禁用",
        },
      ],
      sexList: [
        { name: "未知", value: 0 },
        { name: "男", value: 1 },
        { name: "女", value: 2 },
      ],
      // 弹出框标题
      dialogTitleObj: {
        add: "添加",
        update: "编辑",
        detail: "详情",
      },
    };
  },
  methods: {
    // 操作成功消息提醒内容
    submitOk(msg, cb) {
      console.log("okay");
      this.$notify({
        title: "成功",
        message: msg || "操作成功！",
        type: "success",
        duration: 2000,
        onClose: function () {
          cb && cb();
        },
      });
    },
    // 操作失败消息提醒内容
    submitFail(msg) {
      this.$message({
        message: msg || "网络异常，请稍后重试！",
        type: "error",
      });
    },
  },
};
