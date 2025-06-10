package models

type RolePermission struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}
