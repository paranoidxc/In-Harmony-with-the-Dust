# AI 协作指南

## 当用户提出新功能需求时
AI 必须按顺序输出：
1. `spec.md`（模板见 .specify/templates/）
2. 等待用户确认（用户会说“通过”或“开始写代码”）
3. 生成符合 Go/Gin 和 Vue3 CDN 规范的代码

## 绝对红线
- 不要动 `go.mod` 增加不必要依赖（除非用户要求）。
- 不要生成 `package.json` 里的新依赖（前端是 CDN）。
- 所有 API 返回格式必须统一为 `{ code, message, data }`。# AI 协作指南

## 当用户提出新功能需求时
AI 必须按顺序输出：
1. `spec.md`（模板见 .specify/templates/）
2. 等待用户确认（用户会说“通过”或“开始写代码”）
3. 生成符合 Go/Gin 和 Vue3 CDN 规范的代码

## 绝对红线
- 不要动 `go.mod` 增加不必要依赖（除非用户要求）。
- 不要生成 `package.json` 里的新依赖（前端是 CDN）。
- 所有 API 返回格式必须统一为 `{ code, message, data }`。
