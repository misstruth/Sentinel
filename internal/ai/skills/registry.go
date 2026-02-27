package skills

import "sync"

var (
	registry = make(map[string]*Skill)
	mu       sync.RWMutex
)

// Register 注册 Skill
func Register(skill *Skill) {
	mu.Lock()
	defer mu.Unlock()
	registry[skill.ID] = skill
}

// Get 获取 Skill
func Get(id string) *Skill {
	mu.RLock()
	defer mu.RUnlock()
	return registry[id]
}

// List 获取所有 Skills
func List() []*Skill {
	mu.RLock()
	defer mu.RUnlock()
	skills := make([]*Skill, 0, len(registry))
	for _, s := range registry {
		if s.Enabled {
			skills = append(skills, s)
		}
	}
	return skills
}
