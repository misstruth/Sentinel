package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// Config SMTP 配置
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

// Client SMTP 客户端
type Client struct {
	config *Config
}

// NewClient 创建 SMTP 客户端
func NewClient(cfg *Config) *Client {
	return &Client{config: cfg}
}

// Message 邮件消息
type Message struct {
	To      []string
	Cc      []string
	Subject string
	Body    string
	IsHTML  bool
}

// Send 发送邮件
func (c *Client) Send(msg *Message) error {
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	// 构建邮件内容
	header := c.buildHeader(msg)
	body := header + "\r\n" + msg.Body

	recipients := append(msg.To, msg.Cc...)

	if c.config.UseTLS {
		return c.sendWithTLS(addr, auth, recipients, []byte(body))
	}
	return smtp.SendMail(addr, auth, c.config.From, recipients, []byte(body))
}

// buildHeader 构建邮件头
func (c *Client) buildHeader(msg *Message) string {
	header := fmt.Sprintf("From: %s\r\n", c.config.From)
	header += fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ","))
	if len(msg.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.Cc, ","))
	}
	header += fmt.Sprintf("Subject: %s\r\n", msg.Subject)
	if msg.IsHTML {
		header += "Content-Type: text/html; charset=UTF-8\r\n"
	} else {
		header += "Content-Type: text/plain; charset=UTF-8\r\n"
	}
	return header
}

// sendWithTLS 使用 TLS 发送邮件
func (c *Client) sendWithTLS(addr string, auth smtp.Auth, to []string, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName: c.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, c.config.Host)
	if err != nil {
		return fmt.Errorf("创建 SMTP 客户端失败: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %w", err)
	}

	if err = client.Mail(c.config.From); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %w", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取写入器失败: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭写入器失败: %w", err)
	}

	return client.Quit()
}
