"""
AI 交互组件
"""
import gradio as gr


def create_chat_component():
    """创建聊天组件"""
    return gr.ChatInterface(
        fn=lambda msg, hist: "",
        title="安全分析助手"
    )
