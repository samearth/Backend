package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleMentor UserRole = "mentor"
	RoleMentee UserRole = "mentee"
	RoleAdmin  UserRole = "admin"
)

type VerificationStatus string

const (
	VerificationUnverified VerificationStatus = "unverified"
	VerificationPending    VerificationStatus = "pending"
	VerificationVerified   VerificationStatus = "verified"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseModel
	Email         string   `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string   `gorm:"not null"`
	Role          UserRole `gorm:"type:varchar(20);not null"`
	IsVerified    bool     `gorm:"default:false"`
	IsActive      bool     `gorm:"default:true"`
	LastLoginAt   *time.Time
	ProfileID     *uuid.UUID     `gorm:"type:char(36);"`
	Profile       *Profile       `gorm:"foreignKey:ProfileID;references:ID"`
	MentorProfile *MentorProfile `gorm:"foreignKey:UserID"`
	MenteeProfile *MenteeProfile `gorm:"foreignKey:UserID"`
	Skills        []Skill        `gorm:"many2many:user_skills;"`
}

type Profile struct {
	BaseModel
	FirstName   string `gorm:"size:100;not null"`
	LastName    string `gorm:"size:100;not null"`
	AvatarURL   string
	ImageURL    string
	Bio         string
	Headline    string
	WebsiteURL  string
	LinkedInURL string
	Twitter     string `gorm:"size:100"`
	Timezone    string `gorm:"not null"`
}

type MentorProfile struct {
	BaseModel
	UserID                uuid.UUID `gorm:"type:char(36);not null"`
	ProfileID             uuid.UUID `gorm:"type:char(36);not null"`
	IntroVideoURL         string
	PodcastURL            string
	VerificationStatus    VerificationStatus `gorm:"type:varchar(20);default:'unverified'"`
	ExpertiseArea         string             `gorm:"size:100;not null"`
	Industry              string             `gorm:"size:100;not null"`
	YearsOfExperience     *int
	CalendlyLink          string
	HourlyRate            float64
	IsAcceptingNewMentees bool    `gorm:"default:true"`
	AvailableDays         JSONMap `gorm:"type:json"`
	AvailableTimeSlots    JSONMap `gorm:"type:json"`
}

type MenteeProfile struct {
	BaseModel
	UserID         uuid.UUID `gorm:"type:char(36);uniqueIndex;not null"`
	ProfileID      uuid.UUID `gorm:"type:char(36);not null"`
	CurrentRole    string    `gorm:"size:100"`
	CurrentCompany string    `gorm:"size:100"`
	LearningGoals  JSONMap   `gorm:"type:json"`
	Interests      JSONMap   `gorm:"type:json"`
	SkillLevel     string    `gorm:"size:50"`
}

type Skill struct {
	BaseModel
	Name     string `gorm:"size:100;uniqueIndex;not null"`
	Category string `gorm:"size:100;not null"`
}

type UserSkill struct {
	UserID           uuid.UUID `gorm:"primaryKey;type:char(36)"`
	SkillID          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	ProficiencyLevel string    `gorm:"size:50"`
}

type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
