import random
import gradio as gr
import os
import time
from urllib.request import getproxies


proxies = getproxies()
os.environ["HTTP_PROXY"]  = os.environ["http_proxy"]  = proxies["http"]
os.environ["HTTPS_PROXY"] = os.environ["https_proxy"] = proxies["https"]
os.environ["NO_PROXY"]    = os.environ["no_proxy"]    = "localhost, 127.0.0.1/8, ::1"

def slow_echo(message, history):
    for i in range(len(message)):
        time.sleep(0.3)
        yield "You typed: " + message[: i+1]





# gr.ChatInterface.queue(self, default_concurrency_limit=100)
app = gr.ChatInterface(slow_echo) #//.launch(max_threads=500)
app.queue(200)
app.launch(max_threads=500)

