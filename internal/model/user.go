package model

const (
	UserStatusActive   = 1
	UserStatusDisabled = 0
)

type User struct {
	BaseModel
	Username string `bun:"username,unique" json:"username"`
	Password string `bun:"password_hash" json:"-"`
	Status   int    `bun:"status" json:"status"`
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) TableName() string {
	return "users"
}

type Role struct {
	BaseModel
	Name        string `bun:"name" json:"name"`
	Code        string `bun:"code,unique" json:"code"`
	Description string `bun:"description" json:"description"`
	Status      int    `bun:"status" json:"status"`
}

func (r *Role) TableName() string {
	return "roles"
}

type Permission struct {
	BaseModel
	Name        string `bun:"name" json:"name"`
	Code        string `bun:"code,unique" json:"code"`
	Resource    string `bun:"resource" json:"resource"`
	Action      string `bun:"action" json:"action"`
	Description string `bun:"description" json:"description"`
}

func (p *Permission) TableName() string {
	return "permissions"
}

type UserRole struct {
	BaseModel
	UserID int64 `bun:"user_id" json:"user_id"`
	RoleID int64 `bun:"role_id" json:"role_id"`
}

func (ur *UserRole) TableName() string {
	return "user_roles"
}

type RolePermission struct {
	BaseModel
	RoleID       int64 `bun:"role_id" json:"role_id"`
	PermissionID int64 `bun:"permission_id" json:"permission_id"`
}

func (rp *RolePermission) TableName() string {
	return "role_permissions"
}
