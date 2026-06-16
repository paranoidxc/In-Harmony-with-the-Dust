# 方法 1：临时屏蔽 claude（只对当前终端有效）
alias claude=""  # 让脚本检测不到 claude
go-dev "你的需求"
unalias claude   # 恢复

# 方法 2（更优雅）：我给脚本加个手动开关
# 如果你想要这个功能，可以把脚本第 15 行附近改成：
# if [ "$USE_CODEX" = "1" ]; then AI_CMD="codex"; ... 
# 然后执行 USE_CODEX=1 go-dev "需求"
