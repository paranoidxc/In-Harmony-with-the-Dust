import random
import gradio as gr
import os
from urllib.request import getproxies

proxies = getproxies()
os.environ["HTTP_PROXY"]  = os.environ["http_proxy"]  = proxies["http"]
os.environ["HTTPS_PROXY"] = os.environ["https_proxy"] = proxies["https"]
os.environ["NO_PROXY"]    = os.environ["no_proxy"]    = "localhost, 127.0.0.1/8, ::1"

def alternatingly_agree(message, history):
    print(history)
    if len(history) % 2 == 0:
        return f"Yes, I do think that '{message}'"
    else:
        return "I don't think so"

gr.ChatInterface(alternatingly_agree).launch()
