package mogo

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"reup/internal/models"
	"time"

	"reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	mog "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r implRepository) GetAvailableUsernames(ctx context.Context) ([]string, error) {
	col := r.getBiliSpaceCollection()

	// Sử dụng aggregation để so sánh count_video và count_upload
	pipeline := mog.Pipeline{
		{{Key: "$match", Value: bson.M{
			"$expr": bson.M{
				"$gt": bson.A{
					"$count_video",
					"$count_upload"}}}}},
		{{Key: "$project", Value: bson.M{
			"username": 1,
			"_id":      0}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Username string `bson:"username"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	usernames := make([]string, len(results))
	for i, result := range results {
		usernames[i] = result.Username
	}
	return usernames, nil
}

func (r implRepository) GetTestCodeChannel(ctx context.Context) (*models.Channel, error) {
	log.Println("🐾 Start GetTestCodeChannel...")

	// 1. Lấy toàn bộ video_uploadings để lọc username
	log.Println("📥 Fetching all video_uploadings...")

	cursor, err := r.database.Collection("video_uploading").Find(ctx, bson.M{})
	if err != nil {
		log.Printf("❌ Failed to execute Find query on video_uploading collection: %v\n", err)
		return nil, err
	}
	log.Println("📥 Successfully initiated cursor for video_uploadings")

	// Kiểm tra cursor trước khi sử dụng
	if cursor == nil {
		log.Println("❌ Cursor is nil")
		return nil, fmt.Errorf("cursor is nil")
	}

	// Kiểm tra xem có dữ liệu từ cursor
	if !cursor.Next(ctx) {
		log.Println("❌ No documents found in video_uploading collection")
		return nil, fmt.Errorf("111no documents found")
	}

	log.Println("📥 Cursor has documents")

	var uploadingList []models.VideoUploading
	if err := cursor.All(ctx, &uploadingList); err != nil {
		log.Printf("❌ Failed to decode video_uploadings: %v\n", err)
		return nil, err
	}
	log.Printf("✅ Retrieved %d video_uploadings\n", len(uploadingList))

	// 2. Lọc danh sách username đã upload
	usernames := make([]string, 0, len(uploadingList))
	for _, upload := range uploadingList {
		usernames = append(usernames, upload.Username)
	}
	log.Printf("📃 Extracted %d usernames for exclusion (sample: %v)\n", len(usernames), usernames[:min(3, len(usernames))])

	// 3. Chuẩn bị filter truy vấn chính
	filter := bson.M{
		"$and": bson.A{
			bson.M{"confirm": bson.M{"$ne": 1}},
			bson.M{"subscribe": bson.M{"$ne": -13}},
			bson.M{"username": "debojanipathak623@gmail.com"},
			bson.M{"video_left": 0},
			bson.M{"get_today": bson.M{"$lte": 2}},
			bson.M{"is_upload": 1},
			bson.M{"is_banned": 0},
			bson.M{"username": bson.M{"$nin": usernames}},
		},
	}
	log.Printf("🔍 Query filter: %+v\n", filter)

	// 4. Sắp xếp theo last_get tăng dần
	opts := options.FindOne().SetSort(bson.D{{Key: "last_get", Value: 1}})
	log.Println("📑 Sorting by last_get ASC")

	// 5. Thực thi truy vấn
	var result models.Channel
	err = r.database.Collection("channel").FindOneWithOpt(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("⚠️ No matching channel found (ErrNoDocuments)")
			return nil, nil
		}
		log.Printf("❌ Error while querying channel: %v\n", err)
		return nil, err
	}

	log.Printf("✅ Channel found: %+v\n", result)
	return &result, nil
}

func (r implRepository) AssignProxyIfNeeded(ctx context.Context, channel *models.Channel) error {
	if channel.Proxy == "" {
		// Lấy danh sách proxy live, sắp xếp theo count_used tăng dần, giới hạn 30 proxy
		opts := options.Find()
		opts.SetSort(bson.D{{Key: "count_used", Value: 1}})
		opts.SetLimit(30)

		cursor, err := r.database.Collection("proxy").Find(ctx, bson.M{"live": 1}, opts)
		if err != nil {
			return err
		}

		var proxies []models.Proxy
		if err := cursor.All(ctx, &proxies); err != nil {
			return err
		}

		if len(proxies) == 0 {
			return fmt.Errorf("proxy not available")
		}

		// Random chọn một proxy trong danh sách
		selected := proxies[rand.Intn(len(proxies))]
		channel.Proxy = selected.ProxyIP

		// Tăng count_used của proxy đã chọn
		_, err = r.database.Collection("proxy").UpdateOne(
			ctx,
			bson.M{"proxy_ip": selected.ProxyIP},
			bson.M{"$inc": bson.M{"count_used": 1}},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r implRepository) UpdateChannelData(ctx context.Context, channel *models.Channel) error {
	now := time.Now().Unix()

	channel.LastGet = now
	channel.GetToday += 1

	update := bson.M{
		"$set": bson.M{
			"last_get":  now,
			"get_today": channel.GetToday,
		},
	}

	_, err := r.database.Collection("channel").UpdateOne(
		ctx,
		bson.M{"id": channel.ID},
		update,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r implRepository) GetTestCodeVideo(ctx context.Context, biliSpace *models.BiliSpace) (*models.BiliVideo, error) {
	// 1. Lấy video_id từ video_uploading
	cursorUploading, err := r.database.Collection("video_uploading").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var uploading []struct {
		VideoID string `bson:"video_id"`
	}
	if err := cursorUploading.All(ctx, &uploading); err != nil {
		return nil, err
	}
	uploadingIDs := make([]string, 0, len(uploading))
	for _, v := range uploading {
		uploadingIDs = append(uploadingIDs, v.VideoID)
	}

	// 2. Lấy video_id từ video_uploaded
	cursorUploaded, err := r.database.Collection("video_uploaded").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var uploaded []struct {
		VideoID string `bson:"video_id"`
	}
	if err := cursorUploaded.All(ctx, &uploaded); err != nil {
		return nil, err
	}
	uploadedIDs := make([]string, 0, len(uploaded))
	for _, v := range uploaded {
		uploadedIDs = append(uploadedIDs, v.VideoID)
	}

	// 3. Ghép truy vấn chính
	filter := bson.M{
		"mid":      biliSpace.Mid,
		"duration": bson.M{"$lt": 1800, "$ne": -1},
		"video_id": bson.M{
			"$nin": append(uploadingIDs, uploadedIDs...),
		},
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "date_make", Value: -1}})

	var result models.BiliVideo
	err = r.database.Collection("douyin_video").FindOneWithOpt(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r implRepository) SaveVideoUploading(ctx context.Context, biliVideo *models.BiliVideo, biliSpace *models.BiliSpace, channel *models.Channel) error {

	// Chuyển đổi từ string sang int64
	videoID := biliVideo.VideoID
	videoUploading := models.VideoUploading{
		VideoID:   videoID,
		Username:  channel.Username,
		StartTime: time.Now().Unix(),
		Mid:       biliSpace.Mid,
		Computer:  channel.Computer,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.database.Collection("video_uploading").InsertOne(ctx, videoUploading)
	if err != nil {
		return err
	}
	return nil
}

func (r implRepository) GetNewVideo(ctx context.Context, biliSpace *models.BiliSpace, queCountG bson.M) (*models.BiliVideo, error) {
	// Base filter
	baseFilter := bson.M{
		"mid":           biliSpace.Mid,
		"duration":      bson.M{"$lt": 1800, "$ne": -1},
		"download_fail": bson.M{"$lt": 3},
		"next":          0,
	}
	// Gộp queCountG vào baseFilter
	for k, v := range queCountG {
		baseFilter[k] = v
	}

	// 1. Tìm video có download_url
	filterWithURL := bson.M{}
	for k, v := range baseFilter {
		filterWithURL[k] = v
	}
	filterWithURL["download_url"] = bson.M{"$ne": nil}

	opts := options.FindOne().SetSort(bson.D{{Key: "date_make", Value: -1}})

	var video models.BiliVideo
	data := r.database.Collection("douyin_video").FindOneWithOpt(ctx, filterWithURL, opts).Decode(&video)
	if data == nil {
		return nil, nil
	}

	// 2. Nếu không có video có URL, tìm video download_url = null hoặc rỗng
	filterNoURL := bson.M{}
	for k, v := range baseFilter {
		filterNoURL[k] = v
	}
	filterNoURL["$or"] = bson.A{
		bson.M{"download_url": bson.M{"$eq": nil}},
		bson.M{"download_url": ""},
	}

	data = r.database.Collection("douyin_video").FindOneWithOpt(ctx, filterNoURL, opts).Decode(&video)
	if data == nil {

		return nil, nil
	}
	return &video, nil
}

func (r implRepository) GetAnyVideo(ctx context.Context, biliSpace *models.BiliSpace, queCountG bson.M) (*models.BiliVideo, error) {
	// Tạo filter
	filter := bson.M{
		"mid":           biliSpace.Mid,
		"duration":      bson.M{"$lt": 1800, "$ne": -1},
		"download_fail": bson.M{"$lt": 3},
		"next":          0,
	}

	//Thêm các điều kiện động từ queCountG
	for k, v := range queCountG {
		filter[k] = v
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "date_make", Value: -1}})

	var video models.BiliVideo
	err := r.database.Collection("douyin_video").FindOneWithOpt(ctx, filter, opts).Decode(&video)
	if err != nil {

		return nil, err
	}

	return &video, nil
}

func (r implRepository) GetRegularVideo(ctx context.Context, biliSpace *models.BiliSpace, queCountG bson.M, channel *models.Channel) (*models.BiliVideo, error) {
	if channel.NewOnly == 1 {
		return r.GetNewVideo(ctx, biliSpace, queCountG)
	}
	return r.GetAnyVideo(ctx, biliSpace, queCountG)
}

func (r implRepository) GetRegularChannel(ctx context.Context, computer string, dataAva []string, appsetting *models.AppSetting) (*models.Channel, error) {
	currentHour := time.Now().Hour()
	log.Printf("🐾 Start GetRegularChannel - Current Hour: %d\n", currentHour)

	// Lấy danh sách username đang upload
	log.Println("📥 Fetching all video_uploadings...")
	cursor, err := r.database.Collection("video_uploading").Find(ctx, bson.M{})
	if err != nil {
		log.Printf("❌ Failed to execute Find query on video_uploading collection: %v\n", err)
		return nil, err
	}
	log.Println("✅ Successfully initiated cursor for video_uploadings", cursor)

	var uploading []models.VideoUploading
	if err := cursor.All(ctx, &uploading); err != nil {
		log.Printf("❌ Failed to decode video_uploadings: %v\n", err)
		return nil, err
	}
	log.Printf("✅ Retrieved %d video_uploadings\n", len(uploading))

	uploadingUsernames := make([]string, 0)
	for _, v := range uploading {
		uploadingUsernames = append(uploadingUsernames, v.Username)
	}
	log.Printf("📃 Extracted %d usernames currently uploading (sample: %v)\n", len(uploadingUsernames), uploadingUsernames[:min(3, len(uploadingUsernames))])

	// Filter chính theo điều kiện giờ (check_from/to)
	filter := bson.M{
		"computer":   computer,
		"confirm":    bson.M{"$ne": 1},
		"subscribe":  bson.M{"$ne": -13},
		"video_left": 0,
		"is_upload":  1,
		// "check_from":    bson.M{"$lte": currentHour},
		// "check_to":      bson.M{"$gte": currentHour},
		"upload_today": 0,
		// "get_today":    bson.M{"$lt": 2},
		"is_banned": 0,
		"username": bson.M{
			"$in":  dataAva,
			"$nin": uploadingUsernames,
		},
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "last_upload", Value: 1}})
	log.Println("📑 Sorting by last_upload ASC")

	// Thực thi truy vấn chính
	var channel models.Channel
	err = r.database.Collection("channel").FindOneWithOpt(ctx, filter, opts).Decode(&channel)
	if err == nil {
		log.Printf("✅ Channel found: %+v\n", channel)
		return &channel, nil
	}

	// Nếu không có channel phù hợp trong khung giờ, thì thử với điều kiện đơn giản hơn
	log.Println("⚠️ No channel found with specific time filters. Trying fallback filter...")

	filterFallback := bson.M{
		"computer":     computer,
		"confirm":      bson.M{"$ne": 1},
		"subscribe":    bson.M{"$ne": -13},
		"video_left":   0,
		"upload_today": 0,
		// "get_today":    bson.M{"$lt": 2},
		"is_upload": 1,
		"is_banned": 0,
		"username": bson.M{
			"$nin": uploadingUsernames,
			"$in":  dataAva,
		},
	}

	log.Printf("🔍 Fallback query filter: %+v\n", filterFallback)

	// Thực thi truy vấn fallback
	err = r.database.Collection("channel").FindOneWithOpt(ctx, filterFallback, opts).Decode(&channel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("⚠️ No matching channel found in fallback filter")
			return nil, nil
		}
		log.Printf("❌ Error while querying fallback channel: %v\n", err)
		return nil, err
	}
	log.Printf("✅ Channel found in fallback filter: %+v\n", channel)

	return &channel, nil
}

func (r implRepository) GetChannel(ctx context.Context, computer string, dataAva []string, appsetting *models.AppSetting) (*models.Channel, error) {
	log.Printf("📡 GetChannel called with computer=%s, %d available usernames\n", computer, len(dataAva))

	var channel *models.Channel
	var err error

	if computer == "test_code" {
		log.Println("🧪 Using test_code logic to fetch test channel...")
		channel, err = r.GetTestCodeChannel(ctx)
		if err != nil {
			log.Printf("❌ Error while getting test code channel: %v\n", err)
			return nil, err
		}
		if channel == nil {
			log.Println("⚠️ No test code channel returned.")
		} else {
			log.Printf("✅ Test code channel retrieved: %+v\n", channel)
		}
	} else {
		log.Println("🔍 Using regular logic to fetch channel...")
		log.Printf("📦 AppSetting used: VideoPerDay=%d, MaxNextDay=%d\n", appsetting.VideoPerDay, appsetting.MaxNextDay)
		channel, err = r.GetRegularChannel(ctx, computer, dataAva, appsetting)
		if err != nil {
			log.Printf("❌ Error while getting regular channel: %v\n", err)
			return nil, err
		}
		if channel == nil {
			log.Println("⚠️ No regular channel matched.")
		} else {
			log.Printf("✅ Regular channel retrieved: %+v\n", channel)
		}
	}

	return channel, nil
}

func (r implRepository) GetBiliVideo(ctx context.Context, channel *models.Channel, biliSpace *models.BiliSpace, appsetting *models.AppSetting, computer string) (*models.BiliVideo, error) {
	var queCountG bson.M

	if channel.NextToday > appsetting.MaxNextDay {
		queCountG = bson.M{}
	} else {
		// check count-get =1
		queCountG = bson.M{
			"count_get": bson.M{"$lt": 2},
		}
	}

	if computer == "test_code" {
		return r.GetTestCodeVideo(ctx, biliSpace)
	}
	return r.GetRegularVideo(ctx, biliSpace, queCountG, channel)
}
func (r implRepository) GetAppSetting(ctx context.Context) (*models.AppSetting, error) {
	var appsetting models.AppSetting

	err := r.database.Collection("app-settings").FindOne(ctx, bson.M{}).Decode(&appsetting)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không có bản ghi nào
		}
		return nil, fmt.Errorf("failed to get appsetting: %v", err)
	}

	return &appsetting, nil
}
func (r implRepository) GetBiliSpaceByUsername(ctx context.Context, username string) (*models.BiliSpace, error) {
	var biliSpace models.BiliSpace

	data := r.database.Collection("bili_space").FindOne(ctx, bson.M{
		"username": username,
	}).Decode(&biliSpace)

	if data == nil {
		return nil, data
	}

	return &biliSpace, nil
}
func (r implRepository) UpdateBiliSpaceModel(ctx context.Context, biliSpace *models.BiliSpace) error {
	filter := bson.M{"username": biliSpace.Username}

	update := bson.M{
		"$set": bson.M{
			"count_upload": biliSpace.CountUpload,
			//"updated_at":   time.Now(), // optional nếu có dùng UpdatedAt
		},
	}

	_, err := r.database.Collection("bili_space").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update biliSpace: %v", err)
	}

	return nil
}
func (r implRepository) UpdateChannel(ctx context.Context, channel *models.Channel) error {
	filter := bson.M{"username": channel.Username}

	update := bson.M{
		"$set": bson.M{
			"video_left": channel.VideoLeft,
			//"updated_at": time.Now(), // optional nếu model có UpdatedAt
		},
	}

	_, err := r.database.Collection("channel").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update channel: %v", err)
	}

	return nil
}

// Helper cho limit log
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r implRepository) ProcessVideo(ctx context.Context, biliVideo *models.BiliVideo, channel *models.Channel, dataADV map[string]interface{}) error {
	// 1. Kiểm tra thời lượng video
	if biliVideo.Duration > 1800 {
		return fmt.Errorf("video max 30 mins")
	}

	// 2. (Optional) Nếu muốn dùng nội bộ cho logging/response thì gán tạm
	// Không lưu xuống MongoDB vì không có các trường này trong model
	// biliVideo.Username = channel.Username
	// biliVideo.ProxyIP = channel.Proxy
	// biliVideo.Ads = channel.Ads
	// biliVideo.CheckCopyright = channel.CheckCopyright

	// 3. Lấy duration mới từ dataADV
	durationVal, ok := dataADV["duration"]
	if !ok {
		return fmt.Errorf("duration missing from dataADV")
	}

	var duration int
	switch v := durationVal.(type) {
	case int:
		duration = v
	case float64:
		duration = int(v)
	default:
		return fmt.Errorf("invalid duration type in dataADV")
	}

	// 4. Cập nhật vào DB (chỉ update field có trong model)
	update := bson.M{
		"$set": bson.M{
			"duration":  duration,
			"count_get": biliVideo.CountGet + 1,
		},
	}

	_, err := r.database.Collection("douyin_video").UpdateOne(
		ctx,
		bson.M{"video_id": biliVideo.VideoID},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update biliVideo: %v", err)
	}

	return nil
}

func (r implRepository) GetQuest(ctx context.Context, computer string) (*models.BiliVideo, map[string]interface{}, error) {
	r.l.Info(ctx, "👉 Start GetQuest")

	// 1. Lấy appsetting
	log.Println("📦 Getting app settings...")
	appsetting, err := r.GetAppSetting(ctx)
	if err != nil {
		log.Printf("❌ Failed to get appsetting: %v\n", err)
		return nil, nil, fmt.Errorf("failed to get appsetting: %v", err)
	}
	log.Printf("✅ Appsetting: %+v\n", appsetting)

	// 2. Lấy danh sách username khả dụng
	log.Println("🧍 Getting available usernames...")
	dataAva, err := r.GetAvailableUsernames(ctx)
	if err != nil {
		log.Printf("❌ Failed to get available usernames: %v\n", err)
		return nil, nil, fmt.Errorf("failed to get available usernames: %v", err)
	}
	log.Printf("✅ %d available usernames received\n", len(dataAva))
	log.Printf("🔍 Sample usernames: %v\n", dataAva[:min(3, len(dataAva))]) // show up to 3

	// 3. Lấy channel phù hợp
	log.Printf("📡 Getting channel for computer: %s...\n", computer)
	channel, err := r.GetChannel(ctx, computer, dataAva, appsetting)
	if err != nil {
		log.Printf("❌ Failed to get channel: %v\n", err)
		return nil, nil, err
	}
	if channel == nil {
		r.l.Info(ctx, "⚠️ Channel not available")
		return nil, nil, fmt.Errorf("account not available")
	}
	log.Printf("✅ Got channel: %+v\n", channel)

	// 4. Gán proxy nếu cần và cập nhật dữ liệu channel
	log.Println("🔗 Assigning proxy if needed...")
	err = r.AssignProxyIfNeeded(ctx, channel)
	if err != nil {
		log.Printf("❌ Failed to assign proxy: %v\n", err)
		return nil, nil, err
	}
	log.Println("📦 Updating channel data...")
	err = r.UpdateChannelData(ctx, channel)
	if err != nil {
		log.Printf("❌ Failed to update channel data: %v\n", err)
		return nil, nil, err
	}
	log.Println("✅ Channel updated")

	// 5. Lấy thông tin BiliSpace theo username
	log.Printf("🌌 Getting BiliSpace by username: %s...\n", channel.Username)
	biliSpace, err := r.GetBiliSpaceByUsername(ctx, channel.Username)
	if err != nil {
		log.Printf("❌ Failed to get BiliSpace: %v\n", err)
		return nil, nil, err
	}
	if biliSpace == nil {
		log.Println("⚠️ BiliSpace not available")
		return nil, nil, fmt.Errorf("biliSpace not available")
	}
	log.Printf("✅ BiliSpace: MID=%d, CountVideo=%d, CountUpload=%d\n", biliSpace.Mid, biliSpace.CountVideo, biliSpace.CountUpload)

	// 6. Lấy video phù hợp
	log.Printf("🎬 Getting BiliVideo for MID: %d...\n", biliSpace.Mid)
	biliVideo, err := r.GetBiliVideo(ctx, channel, biliSpace, appsetting, computer)
	if err != nil {
		log.Printf("❌ Failed to get BiliVideo: %v\n", err)
		return nil, nil, err
	}
	if biliVideo == nil {
		log.Println("⚠️ BiliVideo not available, updating counts...")

		biliSpace.CountUpload = biliSpace.CountVideo
		if err := r.UpdateBiliSpaceModel(ctx, biliSpace); err != nil {
			log.Printf("❌ Failed to update BiliSpace: %v\n", err)
			return nil, nil, err
		}

		channel.VideoLeft = 1
		if err := r.UpdateChannel(ctx, channel); err != nil {
			log.Printf("❌ Failed to update channel video left: %v\n", err)
			return nil, nil, err
		}

		return nil, nil, fmt.Errorf("%s video bili not available %d", channel.Username, biliSpace.Mid)
	}
	log.Printf("✅ Got BiliVideo: %+v\n", biliVideo)

	// 7. Ghi log upload và xử lý video
	log.Printf("📝 Saving video uploading: video ID %d\n", biliVideo.VideoID)
	err = r.SaveVideoUploading(ctx, biliVideo, biliSpace, channel)
	if err != nil {
		log.Printf("❌ Failed to save video uploading: %v\n", err)
		return nil, nil, err
	}

	dataADV := map[string]interface{}{
		"duration": biliVideo.Duration,
	}
	log.Printf("🛠️ Processing video BVID: %d with ADV: %+v\n", biliVideo.VideoID, dataADV)
	err = r.ProcessVideo(ctx, biliVideo, channel, dataADV)
	if err != nil {
		log.Printf("❌ Failed to process video: %v\n", err)
		return nil, nil, err
	}

	// 8. Trả về kết quả thành công
	log.Printf("🎉 Successfully got BiliVideo: %d\n", biliVideo.VideoID)
	userProxyInfo := map[string]interface{}{
		"proxy":    channel.Proxy,
		"username": biliSpace.Username,
	}
	r.l.Infof(ctx, "🎉 Successfully got BiliVideo: %d", userProxyInfo)

	return biliVideo, userProxyInfo, nil
}
