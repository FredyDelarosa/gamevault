package services

type NotificationService interface {
	SaveDeviceToken(userID, fcmToken string) error
	SendNotificationToUser(userID, title, body, channel string) error
}