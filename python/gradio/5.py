from fastapi import FastAPI
from fastapi.responses import JSONResponse, Response, StreamingResponse
import uvicorn
import asyncio

# 创建 FastAPI 应用实例
app = FastAPI()

# 定义 /ping 接口
@app.get("/ping")
async def pong():
    print("get ping")
    #await asyncio.sleep(10)  # 使用 asyncio.sleep 代替 time.sleep

    async def streaming_resp():
        message = 'hello world'
        for i in range(len(message)):
            await asyncio.sleep(0.3)
            yield "You typed: " + message[: i+1]
    return StreamingResponse(streaming_resp())

    return JSONResponse(content={"text": "pong"})

# 启动服务器
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8888, reload=False)
