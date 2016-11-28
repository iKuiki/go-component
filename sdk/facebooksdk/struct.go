package facebooksdk

import (
	"time"
)

type FacebookAccessToken struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type FacebookUserInfo struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	Gender      string `json:"gender"`
	Id          string `json:"id"`
	LastName    string `json:"last_name"`
	Link        string `json:"link"`
	Locale      string `json:"locale"`
	Name        string `json:"name"`
	Timezone    int32  `json:"timezone"`
	UpdatedTime string `json:"updated_time"`
}
