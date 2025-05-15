from fastapi import FastAPI
from fastapi.responses import JSONResponse, Response, StreamingResponse
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import asyncio
import requests, json
import httpx


# 创建 FastAPI 应用实例
app = FastAPI()

# 定义 /ping 接口
@app.get("/ping")
async def pong():
    print("get ping")

    async def streaming_resp():
        async with httpx.AsyncClient() as client:
            r = await client.post(
                'http://localhost:11434/api/generate',
                json={
                    'model': 'llama3.2:3b',
                    'prompt': 'Why is the sky blue?',
                    'context': [],
                    'options': {
                        'top_k': 40,
                        'temperature': 0.5,
                        'top_p': 0.9
                    }
                },
                timeout=None  # 设置为 None 以避免超时
            )
            r.raise_for_status()

            response = ""

            async for line in r.aiter_lines():
                if line:
                    body = json.loads(line)
                    response_part = body.get('response', '')
                    print(response_part)
                    if 'error' in body:
                        raise Exception(body['error'])

                    response += response_part

                    if body.get('done', False):
                        context = body.get('context', [])
                        yield response
                        return

    return StreamingResponse(streaming_resp())


    # async def streaming_resp():
    #     r = requests.post('http://localhost:11434/api/generate',
    #                     json={
    #                         'model': 'llama3.2:3b',
    #                         'prompt': 'Why is the sky blue?',
    #                         'context': [],
    #                         'options':{
    #                             'top_k': 40,
    #                             'temperature':0.5,
    #                             'top_p': 0.9
    #                         }
    #                     },
    #                     stream=False)
    #     r.raise_for_status()
    #
    #     response = ""
    #
    #     for line in r.iter_lines():
    #         body = json.loads(line)
    #         response_part = body.get('response', '')
    #         print(response_part)
    #         if 'error' in body:
    #             raise Exception(body['error'])
    #
    #         response += response_part
    #
    #         if body.get('done', False):
    #             context = body.get('context', [])
    #             return response, context
    #
    # return StreamingResponse(streaming_resp())

    #await asyncio.sleep(10)  # 使用 asyncio.sleep 代替 time.sleep

    # async def streaming_resp():
    #     message = 'hello world'
    #     for i in range(len(message)):
    #         await asyncio.sleep(0.3)
    #         yield "You typed: " + message[: i+1]
    # return StreamingResponse(streaming_resp())

    return JSONResponse(content={"text": "pong"})
    #return "pong"

# 启动服务器
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8888, reload=False)
