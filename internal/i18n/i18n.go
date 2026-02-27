package i18n

// Language 语言类型
type Language string

const (
	LangZhCN Language = "zh-CN"
	LangEnUS Language = "en-US"
)

// Messages 消息集合
type Messages map[string]string

// Translator 翻译器
type Translator struct {
	lang     Language
	messages map[Language]Messages
}

// NewTranslator 创建翻译器
func NewTranslator(lang Language) *Translator {
	return &Translator{
		lang:     lang,
		messages: make(map[Language]Messages),
	}
}

// Load 加载语言包
func (t *Translator) Load(lang Language, msgs Messages) {
	t.messages[lang] = msgs
}

// T 翻译
func (t *Translator) T(key string) string {
	if msgs, ok := t.messages[t.lang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	return key
}
