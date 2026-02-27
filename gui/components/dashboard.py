"""
数据统计面板
"""
import gradio as gr


def create_dashboard():
    """创建统计面板"""
    with gr.Row():
        gr.Number(label="总事件数")
        gr.Number(label="高危事件")
    return None
