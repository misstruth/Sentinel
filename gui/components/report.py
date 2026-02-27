"""
报告预览组件
"""
import gradio as gr


def create_report_preview():
    """创建报告预览组件"""
    with gr.Column():
        report_id = gr.Number(label="报告ID")
        preview_btn = gr.Button("预览")
        output = gr.Markdown()
    return report_id, preview_btn, output
