package events

type Event struct {
	ID           string   `json:"id,omitempty" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Title        string   `json:"title" gorm:"column:title"`
	Image        string   `json:"image" gorm:"column:image_url"`
	LocationID   string   `json:"locationId,omitempty" gorm:"column:location_id"`
	Location     Location `json:"location,omitempty"`
	Date         string   `json:"date,omitempty" gorm:"column:event_date"`
	Time         string   `json:"time,omitempty" gorm:"column:event_time"`
	Price        int      `json:"price,omitempty" gorm:"column:price"`
	Detail       string   `json:"detail,omitempty" gorm:"column:description"`
	IsBanner     bool     `json:"-" gorm:"column:is_banner"`
	IsRecommend  bool     `json:"-" gorm:"column:is_recommend"`
	IsComingSoon bool     `json:"-" gorm:"column:is_coming_soon"`
}

type Location struct {
	ID            string  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Name          string  `json:"name" gorm:"column:name"`
	Latitude      float64 `json:"latitude" gorm:"column:latitude"`
	Longitude     float64 `json:"longitude" gorm:"column:longitude"`
	City          string  `json:"city" gorm:"column:city"`
	StateProvince string  `json:"stateProvince" gorm:"column:state_province"`
	Country       string  `json:"country" gorm:"column:country"`
	PostCode      string  `json:"postCode" gorm:"column:post_code"`
	IsActive      bool    `json:"isActive" gorm:"column:is_active"`
}
