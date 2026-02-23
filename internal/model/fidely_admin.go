package model

// FidelyAdmin represents the fidely_admin table.
// Platform-level administrators (not store admins).
type FidelyAdmin struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"-"` // Never expose password in JSON
	Role     int    `json:"role"`
}

// FidelyAdminCreate is used when creating a new Fidely admin.
type FidelyAdminCreate struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

// FidelyAdminUpdate is used when updating an existing Fidely admin.
type FidelyAdminUpdate struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}
