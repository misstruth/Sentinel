package rbac

// Role 角色
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleViewer Role = "viewer"
)

// Permission 权限
type Permission string

const (
	PermRead   Permission = "read"
	PermWrite  Permission = "write"
	PermDelete Permission = "delete"
)

// Manager 权限管理器
type Manager struct {
	roles map[Role][]Permission
}

// NewManager 创建管理器
func NewManager() *Manager {
	m := &Manager{
		roles: make(map[Role][]Permission),
	}
	m.initDefaults()
	return m
}

// initDefaults 初始化默认权限
func (m *Manager) initDefaults() {
	m.roles[RoleAdmin] = []Permission{PermRead, PermWrite, PermDelete}
	m.roles[RoleUser] = []Permission{PermRead, PermWrite}
	m.roles[RoleViewer] = []Permission{PermRead}
}

// HasPermission 检查权限
func (m *Manager) HasPermission(role Role, perm Permission) bool {
	perms, ok := m.roles[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
