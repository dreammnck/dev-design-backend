package auth

import "time"

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleOrganization UserRole = "organization"
	RoleCustomer     UserRole = "customer"
)

// User is the GORM model mapped to the users table
type User struct {
	ID           string     `json:"id"           gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Username     string     `json:"username"     gorm:"column:username;uniqueIndex;not null"`
	Email        string     `json:"email"        gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string     `json:"-"            gorm:"column:password_hash;not null"`
	Role         UserRole   `json:"role"         gorm:"type:user_role;column:role;not null;default:customer"`
	IsActive     bool       `json:"isActive"     gorm:"column:is_active;default:true"`
	CreatedAt    time.Time  `json:"createdAt"    gorm:"column:created_at"`
	UpdatedAt    time.Time  `json:"updatedAt"    gorm:"column:updated_at"`
	DeletedAt    *time.Time `json:"-"            gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}

// LoginReq is the request body for POST /auth/login
type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginRes is the response body for POST /auth/login
type LoginRes struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo is the public user information returned in the token response
type UserInfo struct {
	ID       string    `json:"id"`
	Username string    `json:"username"` // Match DB column name
	Email    string    `json:"email"`
	Role     UserRole  `json:"role"`
	Status   string    `json:"status"`   // active | suspended
	JoinedAt time.Time `json:"joinedAt"` // CreatedAt
}

// JWTClaims are the custom claims embedded in the JWT
type JWTClaims struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

// RegisterReq is the request body for POST /auth/register
type RegisterReq struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Role     UserRole `json:"role"` // Optional: default customer
}

// UpdateUserReq is the request body for PATCH /users/:id
type UpdateUserReq struct {
	Username string   `json:"username"` // Match DB column name
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
	Status   string   `json:"status"` // map to is_active (active | suspended)
}

// UpdateRoleReq is the request body for PUT /users/:id/role
type UpdateRoleReq struct {
	Role UserRole `json:"role" binding:"required"`
}
