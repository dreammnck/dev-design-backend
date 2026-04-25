package events

import (
	"mime/multipart"
	"time"
)

type PublishStatus string

const (
	PublishStatusPending  PublishStatus = "pending"
	PublishStatusApproved PublishStatus = "approved"
	PublishStatusRejected PublishStatus = "rejected"
)

const (
	LocationTypeA = "a"
	LocationTypeB = "b"
	LocationTypeC = "c"
	LocationTypeD = "d"
)

type PayoutStatus string

const (
	PayoutStatusRequested  PayoutStatus = "requested"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusRejected   PayoutStatus = "rejected"
)

type Event struct {
	ID             string        `json:"id,omitempty"             gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Title          string        `json:"title"                    gorm:"column:title"`
	Image          string        `json:"image"                    gorm:"column:image_url"`
	LocationID     string        `json:"locationId,omitempty"     gorm:"column:location_id"`
	Location       Location      `json:"location,omitempty"`
	Date           string        `json:"date,omitempty"           gorm:"column:event_date"`
	Time           *string       `json:"time,omitempty"           gorm:"column:event_time"`
	Price          int           `json:"price,omitempty"          gorm:"column:price"`
	Detail         string        `json:"detail,omitempty"         gorm:"column:description"`
	IsBanner       bool          `json:"-"                        gorm:"column:is_banner"`
	IsRecommend    bool          `json:"-"                        gorm:"column:is_recommend"`
	IsComingSoon   bool          `json:"-"                        gorm:"column:is_coming_soon"`
	OrganizationID *string       `json:"organizationId,omitempty" gorm:"type:uuid;column:organization_id"`
	UserID         *string       `json:"userId,omitempty"         gorm:"type:uuid;column:user_id"`
	PublishStatus  PublishStatus `json:"publishStatus,omitempty"  gorm:"type:event_publish_status;column:publish_status;default:pending"`
	RejectReason   *string       `json:"rejectReason,omitempty"   gorm:"column:reject_reason"`
	PublishedAt    *time.Time    `json:"publishedAt,omitempty"    gorm:"column:published_at"`
	IsFav          bool          `json:"isFav"                    gorm:"-"`
	CreatedAt      time.Time     `json:"createdAt,omitempty"      gorm:"column:created_at"`
	UpdatedAt      time.Time     `json:"updatedAt,omitempty"      gorm:"column:updated_at"`
}

func (Event) TableName() string { return "events" }

type Location struct {
	ID            string  `json:"id"            gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Name          string  `json:"name"          gorm:"column:name"`
	Latitude      float64 `json:"latitude"      gorm:"column:latitude"`
	Longitude     float64 `json:"longitude"     gorm:"column:longitude"`
	City          string  `json:"city"          gorm:"column:city"`
	StateProvince string  `json:"stateProvince" gorm:"column:state_province"`
	Country       string  `json:"country"       gorm:"column:country"`
	PostCode      string  `json:"postCode"      gorm:"column:post_code"`
	Type          string  `json:"type"          gorm:"type:location_type;column:type;not null;default:a"`
	IsActive      bool    `json:"isActive"      gorm:"column:is_active"`
}

func (Location) TableName() string { return "locations" }

// ── Request / Response DTOs ──────────────────────────────────────────────────

// EventCreateReq is used by organization to create a new event
type EventCreateReq struct {
	Title         string                `form:"title"        binding:"required"`
	Description   string                `form:"description"`
	ImageFile     *multipart.FileHeader `form:"imageUrl"     swaggerignore:"true"`
	LocationName  string                `form:"locationName" binding:"required"`
	Latitude      float64               `form:"latitude"`
	Longitude     float64               `form:"longitude"`
	City          string                `form:"city"`
	StateProvince string                `form:"stateProvince"`
	Country       string                `form:"country"`
	PostCode      string                `form:"postCode"`
	LocationType  string                `form:"locationType"`
	EventDate     string                `form:"eventDate"    binding:"required"` // YYYY-MM-DD
	EventTime     string                `form:"eventTime"`                       // HH:MM
	Price         int                   `form:"price"`
	IsBanner      bool                  `form:"isBanner"`
	IsRecommend   bool                  `form:"isRecommend"`
	IsComingSoon  bool                  `form:"isComingSoon"`
}

// EventUpdateReq is used by organization to edit an owned event (all fields optional)
type EventUpdateReq struct {
	Title        *string `json:"title"`
	Detail       *string `json:"description"`
	Image        *string `json:"image"`
	LocationID   *string `json:"locationId"`
	Date         *string `json:"eventDate"`
	Time         *string `json:"eventTime"`
	Price        *int    `json:"price"`
	IsBanner     *bool   `json:"isBanner"`
	IsRecommend  *bool   `json:"isRecommend"`
	IsComingSoon *bool   `json:"isComingSoon"` // Supports the typo fallback at the frontend or correct it
}

// ReviewReq is used by admin to approve or reject a pending event
type ReviewReq struct {
	Status       PublishStatus `json:"status"       binding:"required"` // approved | rejected
	RejectReason string        `json:"rejectReason"`                    // required when rejected
}

// AdminEditEventReq is used by admin to edit events (all fields optional, only sent fields are updated)
type AdminEditEventReq struct {
	Title         *string        `json:"title"`
	Description   *string        `json:"description"`
	LocationID    *string        `json:"locationId"`
	EventDate     *string        `json:"eventDate"`
	EventTime     *string        `json:"eventTime"`
	IsBanner      *bool          `json:"isBanner"`
	IsRecommend   *bool          `json:"isRecommend"`
	IsComingSoon  *bool          `json:"isComingSoon"`
	PublishStatus *PublishStatus `json:"publishStatus"`
}

// Payout is the GORM model for payout requests
type Payout struct {
	ID             string       `json:"id"             gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	OrganizationID string       `json:"organizationId" gorm:"type:uuid;column:organization_id;not null"`
	EventID        *string      `json:"eventId"        gorm:"type:uuid;column:event_id"`
	Amount         int          `json:"amount"         gorm:"column:amount;not null"`
	Status         PayoutStatus `json:"status"         gorm:"type:payout_status;column:status;default:requested"`
	AccountName    string       `json:"accountName"    gorm:"column:account_name"`
	BankAccount    string       `json:"bankAccount"    gorm:"column:bank_account"`
	BankName       string       `json:"bankName"       gorm:"column:bank_name"`
	RejectReason   *string      `json:"rejectReason"   gorm:"column:reject_reason"`
	RequestedAt    time.Time    `json:"requestedAt"    gorm:"column:requested_at"`
	ProcessedAt    *time.Time   `json:"processedAt"    gorm:"column:processed_at"`
	CreatedAt      time.Time    `json:"createdAt"      gorm:"column:created_at"`
	UpdatedAt      time.Time    `json:"updatedAt"      gorm:"column:updated_at"`
}

func (Payout) TableName() string { return "payouts" }

// PayoutReq is the request body for requesting a payout
type PayoutReq struct {
	EventID     string `json:"eventId"` // Optional
	Amount      int    `json:"amount"      binding:"required"`
	AccountName string `json:"accountName" binding:"required"`
	BankAccount string `json:"bankAccount" binding:"required"`
	BankName    string `json:"bankName"    binding:"required"`
}

// SeatBatchCreateReq is used by organization to bulk-add seats to an event
type SeatBatchCreateReq struct {
	Seats []SeatInput `json:"seats" binding:"required,min=1"`
}

type SeatInput struct {
	SeatNumber string `json:"seatNumber"` // Optional if Capacity is provided
	Price      int    `json:"price"`
	SeatType   string `json:"seatType"`
	Capacity   int    `json:"capacity"` // If > 0, system will auto-generate N seats
}

type UserFavorite struct {
	UserID    string    `gorm:"type:uuid;primaryKey;column:user_id"`
	EventID   string    `gorm:"type:uuid;primaryKey;column:event_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserFavorite) TableName() string { return "user_favorites" }
