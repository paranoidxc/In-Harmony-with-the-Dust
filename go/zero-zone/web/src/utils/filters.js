export const filters = {
  sexName: (sex) => {
    const { proxy } = getCurrentInstance();
    let result = proxy.sexList.find((obj) => obj.value == sex);
    return result ? result.name : "数据丢失";
  },
  typoName: (typo) => {
    const { proxy } = getCurrentInstance();
    let result = proxy.qdTypoOptions.find((obj) => obj.value == typo);
    return result ? (result.name ? result.name : result.label) : "数据丢失";
  },
  optLabelOrName: (options, idx) => {
    const index = options.findIndex((option) => option.value === idx);
    if (index !== -1) {
      return options[index].label
        ? options[index].label
        : options[index].name
          ? options[index].name
          : "数据丢失";
    }
    return "数据丢失";
  },
};
