"""
Gradio 主应用 - Fo-Sentinel 安全哨兵
"""
import gradio as gr
import requests
import uuid

API_BASE = "http://localhost:6872/api"


# ==================== API 调用函数 ====================

def get_subscriptions():
    """获取订阅列表"""
    try:
        resp = requests.get(f"{API_BASE}/subscriptions")
        data = resp.json()
        if data.get("message") == "OK":
            return data.get("data", {}).get("items", [])
        return []
    except:
        return []


def create_subscription(name, description, source_type, source_url, cron_expr):
    """创建订阅"""
    try:
        resp = requests.post(f"{API_BASE}/subscriptions", json={
            "name": name,
            "description": description,
            "source_type": source_type,
            "source_url": source_url,
            "cron_expr": cron_expr
        })
        data = resp.json()
        if data.get("message") == "OK":
            return f"✅ 创建成功，ID: {data.get('data', {}).get('id')}"
        return f"❌ 创建失败: {data.get('message')}"
    except Exception as e:
        return f"❌ 请求失败: {str(e)}"


def delete_subscription(sub_id):
    """删除订阅"""
    try:
        resp = requests.delete(f"{API_BASE}/subscriptions/{int(sub_id)}")
        data = resp.json()
        if data.get("message") == "OK":
            return "✅ 删除成功"
        return f"❌ 删除失败: {data.get('message')}"
    except Exception as e:
        return f"❌ 请求失败: {str(e)}"


def get_events(page=1, page_size=20):
    """获取安全事件"""
    try:
        resp = requests.get(f"{API_BASE}/event", params={"page": page, "page_size": page_size})
        data = resp.json()
        if data.get("message") == "OK":
            return data.get("data", {}).get("list", []), data.get("data", {}).get("total", 0)
        return [], 0
    except:
        return [], 0


def get_reports(page=1, page_size=20):
    """获取报告列表"""
    try:
        resp = requests.get(f"{API_BASE}/report", params={"page": page, "page_size": page_size})
        data = resp.json()
        if data.get("message") == "OK":
            return data.get("data", {}).get("list", []), data.get("data", {}).get("total", 0)
        return [], 0
    except:
        return [], 0


def get_report_detail(report_id):
    """获取报告详情"""
    try:
        resp = requests.get(f"{API_BASE}/report/{int(report_id)}")
        data = resp.json()
        if data.get("message") == "OK":
            report = data.get("data", {})
            return report.get("content", "无内容")
        return f"获取失败: {data.get('message')}"
    except Exception as e:
        return f"请求失败: {str(e)}"


def generate_report(title, report_type, start_time, end_time):
    """生成报告"""
    try:
        resp = requests.post(f"{API_BASE}/report/generate", json={
            "title": title,
            "type": report_type,
            "start_time": start_time,
            "end_time": end_time,
            "template_id": 1
        }, timeout=120)
        data = resp.json()
        if data.get("message") == "OK":
            report_data = data.get("data", {})
            return f"✅ 报告生成成功\n\nID: {report_data.get('report_id')}\n\n**摘要:**\n{report_data.get('summary', '')}"
        return f"❌ 生成失败: {data.get('message')}"
    except Exception as e:
        return f"❌ 请求失败: {str(e)}"


def get_templates():
    """获取报告模板"""
    try:
        resp = requests.get(f"{API_BASE}/report/template")
        data = resp.json()
        if data.get("message") == "OK":
            return data.get("data", {}).get("list", [])
        return []
    except:
        return []


# 会话ID管理
session_ids = {}

def chat_with_ai(message, history):
    """AI 对话"""
    try:
        # 使用固定的会话ID保持上下文
        session_id = "gradio-session-001"

        resp = requests.post(
            f"{API_BASE}/chat",
            json={
                "id": session_id,
                "question": message
            },
            timeout=90
        )
        data = resp.json()
        if data.get("message") == "OK":
            return data.get("data", {}).get("answer", "无响应")
        return f"错误: {data.get('message')}"
    except Exception as e:
        return f"服务暂不可用: {str(e)}"


# ==================== UI 组件 ====================

def refresh_subscriptions():
    """刷新订阅列表"""
    subs = get_subscriptions()
    if not subs:
        return "暂无订阅"

    result = "| ID | 名称 | 类型 | 状态 | 创建时间 |\n|---|---|---|---|---|\n"
    for sub in subs:
        result += f"| {sub.get('id')} | {sub.get('name')} | {sub.get('source_type')} | {sub.get('status')} | {sub.get('created_at', '')[:19]} |\n"
    return result


def refresh_events():
    """刷新事件列表"""
    events, total = get_events()
    if not events:
        return f"暂无安全事件 (共 {total} 条)"

    result = f"**共 {total} 条事件**\n\n| ID | 标题 | 严重程度 | 状态 | 时间 |\n|---|---|---|---|---|\n"
    for event in events:
        result += f"| {event.get('id')} | {event.get('title', '')[:30]} | {event.get('severity')} | {event.get('status')} | {event.get('event_time', '')[:19]} |\n"
    return result


def refresh_reports():
    """刷新报告列表"""
    reports, total = get_reports()
    if not reports:
        return f"暂无报告 (共 {total} 条)"

    result = f"**共 {total} 条报告**\n\n| ID | 标题 | 类型 | 状态 | 创建时间 |\n|---|---|---|---|---|\n"
    for report in reports:
        result += f"| {report.get('id')} | {report.get('title', '')} | {report.get('type')} | {report.get('status')} | {report.get('created_at', '')[:19]} |\n"
    return result


def get_dashboard_stats():
    """获取仪表盘统计"""
    subs = get_subscriptions()
    events, event_total = get_events()
    reports, report_total = get_reports()

    stats = f"""
## 📊 系统概览

| 指标 | 数量 |
|------|------|
| 📡 订阅源 | {len(subs)} |
| 🚨 安全事件 | {event_total} |
| 📋 分析报告 | {report_total} |

---

### 最近订阅
"""
    if subs:
        for sub in subs[:5]:
            status_icon = "🟢" if sub.get('status') == 'active' else "🟡"
            stats += f"- {status_icon} **{sub.get('name')}** ({sub.get('source_type')})\n"
    else:
        stats += "暂无订阅\n"

    return stats


def create_app():
    """创建 Gradio 应用"""
    with gr.Blocks(title="Fo-Sentinel 安全哨兵", theme=gr.themes.Soft()) as app:
        gr.Markdown("# 🛡️ Fo-Sentinel 安全哨兵")
        gr.Markdown("基于多 Agent 协作的智能安全分析系统")

        with gr.Tabs():
            # ==================== 仪表盘 ====================
            with gr.Tab("📊 仪表盘"):
                dashboard_output = gr.Markdown(get_dashboard_stats())
                refresh_dashboard_btn = gr.Button("🔄 刷新", variant="secondary")
                refresh_dashboard_btn.click(get_dashboard_stats, outputs=dashboard_output)

            # ==================== AI 对话 ====================
            with gr.Tab("🤖 AI 助手"):
                gr.Markdown("### 智能安全分析助手")
                gr.Markdown("可以询问当前时间、查询数据库、分析安全事件等")
                chatbot = gr.ChatInterface(
                    chat_with_ai,
                    examples=[
                        "现在几点了？",
                        "帮我分析一下最近的安全态势",
                        "有哪些常见的安全漏洞类型？"
                    ],
                    retry_btn=None,
                    undo_btn=None,
                )

            # ==================== 订阅管理 ====================
            with gr.Tab("📡 订阅管理"):
                with gr.Row():
                    with gr.Column(scale=2):
                        gr.Markdown("### 订阅列表")
                        sub_list_output = gr.Markdown(refresh_subscriptions())
                        refresh_sub_btn = gr.Button("🔄 刷新列表", variant="secondary")
                        refresh_sub_btn.click(refresh_subscriptions, outputs=sub_list_output)

                    with gr.Column(scale=1):
                        gr.Markdown("### 添加订阅")
                        sub_name = gr.Textbox(label="名称", placeholder="如: GitHub Security Advisory")
                        sub_desc = gr.Textbox(label="描述", placeholder="订阅描述")
                        sub_type = gr.Dropdown(
                            label="类型",
                            choices=["github_repo", "rss", "nvd", "cve", "vulnerability", "threat_intel"],
                            value="github_repo"
                        )
                        sub_url = gr.Textbox(label="源地址", placeholder="https://...")
                        sub_cron = gr.Textbox(label="Cron表达式", value="0 */6 * * *", placeholder="0 */6 * * *")
                        create_sub_btn = gr.Button("➕ 创建订阅", variant="primary")
                        create_result = gr.Markdown()
                        create_sub_btn.click(
                            create_subscription,
                            inputs=[sub_name, sub_desc, sub_type, sub_url, sub_cron],
                            outputs=create_result
                        )

                with gr.Row():
                    gr.Markdown("### 删除订阅")
                    del_sub_id = gr.Number(label="订阅ID", precision=0)
                    del_sub_btn = gr.Button("🗑️ 删除", variant="stop")
                    del_result = gr.Markdown()
                    del_sub_btn.click(delete_subscription, inputs=del_sub_id, outputs=del_result)

            # ==================== 安全事件 ====================
            with gr.Tab("🚨 安全事件"):
                gr.Markdown("### 安全事件列表")
                event_output = gr.Markdown(refresh_events())
                refresh_event_btn = gr.Button("🔄 刷新", variant="secondary")
                refresh_event_btn.click(refresh_events, outputs=event_output)

            # ==================== 分析报告 ====================
            with gr.Tab("📋 分析报告"):
                with gr.Row():
                    with gr.Column(scale=1):
                        gr.Markdown("### 报告列表")
                        report_list_output = gr.Markdown(refresh_reports())
                        refresh_report_btn = gr.Button("🔄 刷新列表", variant="secondary")
                        refresh_report_btn.click(refresh_reports, outputs=report_list_output)

                        gr.Markdown("### 查看报告")
                        report_id_input = gr.Number(label="报告ID", precision=0)
                        view_report_btn = gr.Button("📖 查看详情")
                        report_detail_output = gr.Markdown()
                        view_report_btn.click(get_report_detail, inputs=report_id_input, outputs=report_detail_output)

                    with gr.Column(scale=1):
                        gr.Markdown("### 生成新报告")
                        report_title = gr.Textbox(label="报告标题", placeholder="如: 2026年2月安全周报")
                        report_type = gr.Dropdown(
                            label="报告类型",
                            choices=["daily", "weekly", "monthly", "custom"],
                            value="weekly"
                        )
                        report_start = gr.Textbox(label="开始时间", placeholder="2026-02-01 00:00:00")
                        report_end = gr.Textbox(label="结束时间", placeholder="2026-02-12 23:59:59")
                        generate_btn = gr.Button("🚀 生成报告", variant="primary")
                        generate_result = gr.Markdown()
                        generate_btn.click(
                            generate_report,
                            inputs=[report_title, report_type, report_start, report_end],
                            outputs=generate_result
                        )

    return app


if __name__ == "__main__":
    app = create_app()
    app.launch(server_name="0.0.0.0", server_port=7860)
