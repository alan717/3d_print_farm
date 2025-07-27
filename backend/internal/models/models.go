package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Material struct {
	ID           uint                   `json:"id" gorm:"primaryKey"`
	Name         string                 `json:"name" gorm:"uniqueIndex;not null"`
	Description  string                 `json:"description"`
	PricePerGram float64                `json:"price_per_gram" gorm:"not null"`
	Properties   map[string]interface{} `json:"properties" gorm:"type:jsonb"`
	Available    bool                   `json:"available" gorm:"default:true"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type File struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Filename   string    `json:"filename" gorm:"not null"`
	OriginalName string  `json:"original_name" gorm:"not null"`
	Size       int64     `json:"size" gorm:"not null"`
	MimeType   string    `json:"mime_type" gorm:"not null"`
	Path       string    `json:"path" gorm:"not null"`
	Volume     float64   `json:"volume"` // 模型体积（立方厘米）
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"    // 待处理
	OrderStatusPrinting   OrderStatus = "printing"   // 打印中
	OrderStatusCompleted  OrderStatus = "completed"  // 已完成
	OrderStatusShipped    OrderStatus = "shipped"    // 已发货
	OrderStatusCancelled  OrderStatus = "cancelled"  // 已取消
)

type Order struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	User          User        `json:"user" gorm:"foreignKey:UserID"`
	Status        OrderStatus `json:"status" gorm:"default:'pending'"`
	TotalAmount   float64     `json:"total_amount" gorm:"not null"`
	ShippingCost  float64     `json:"shipping_cost" gorm:"default:0"`
	TrackingNumber string     `json:"tracking_number"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type Quality string

const (
	QualityFast     Quality = "fast"     // 快速
	QualityStandard Quality = "standard" // 标准
	QualityHigh     Quality = "high"     // 精细
)

type OrderItem struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	OrderID      uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	Order        Order     `json:"order" gorm:"foreignKey:OrderID"`
	FileID       uuid.UUID `json:"file_id" gorm:"type:uuid;not null"`
	File         File      `json:"file" gorm:"foreignKey:FileID"`
	MaterialID   uint      `json:"material_id" gorm:"not null"`
	Material     Material  `json:"material" gorm:"foreignKey:MaterialID"`
	Quality      Quality   `json:"quality" gorm:"not null"`
	Quantity     int       `json:"quantity" gorm:"default:1"`
	SupportEnabled bool    `json:"support_enabled" gorm:"default:false"`
	Price        float64   `json:"price" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GORM回调：在创建用户前生成UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}