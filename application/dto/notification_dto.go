package dto

type SaveFcmTokenRequest struct {
	FcmToken string `json:"fcm_token" binding:"required"`
}

type SendNotificationRequest struct {
	Title   string `json:"title" binding:"required"`
	Body    string `json:"body" binding:"required"`
	Channel string `json:"channel"`
}
