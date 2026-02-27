package template

import (
	"bytes"
	"html/template"
)

// Template 邮件模板
type Template struct {
	Name    string
	Subject string
	Body    string
}

// Engine 模板引擎
type Engine struct {
	templates map[string]*Template
}

// NewEngine 创建模板引擎
func NewEngine() *Engine {
	return &Engine{
		templates: make(map[string]*Template),
	}
}

// Register 注册模板
func (e *Engine) Register(tpl *Template) {
	e.templates[tpl.Name] = tpl
}

// Get 获取模板
func (e *Engine) Get(name string) *Template {
	return e.templates[name]
}

// Render 渲染模板
func (e *Engine) Render(name string, data interface{}) (string, string, error) {
	tpl := e.templates[name]
	if tpl == nil {
		return "", "", nil
	}

	subject, err := renderString(tpl.Subject, data)
	if err != nil {
		return "", "", err
	}

	body, err := renderString(tpl.Body, data)
	if err != nil {
		return "", "", err
	}

	return subject, body, nil
}

// renderString 渲染字符串模板
func renderString(text string, data interface{}) (string, error) {
	tpl, err := template.New("").Parse(text)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
