package repository

type ChannelGroup struct {
	Computer  string `bson:"_id"`        // Nhóm theo trường Computer
	TotalUser int    `bson:"total_user"` // Đếm theo số lượng Username
}
