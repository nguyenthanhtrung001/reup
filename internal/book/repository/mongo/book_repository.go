package mogo

import (
	"context"
	"reup/internal/book/repository"
	"reup/internal/models"
	"reup/pkg/mongo"
	"reup/pkg/paginator"

	"time"

	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	bookCollection = "books"
)

func (repo implRepository) getBookCollection() mongo.Collection {
	return repo.database.Collection(bookCollection)
}

// CreateBook implements repository.Repository.
func (repo implRepository) CreateBook(ctx context.Context, input repository.CreateOption) error {
	col := repo.getBookCollection()
	newBook := models.Book{
		ID:          primitive.NewObjectID(),
		Title:       input.Title,
		Author:      input.Author,
		PublishedAt: input.PublishedAt,
		CreateAt:    time.Now(),
	}

	if _, err := col.InsertOne(ctx, newBook); err != nil {
		repo.l.Warnf(ctx, "book.repo.mogo.create")
		return err
	}
	return nil
}

// DeleteBook implements repository.Repository.
func (repo implRepository) DeleteBook(ctx context.Context, idBook string) error {
	col := repo.getBookCollection()
	id, err := primitive.ObjectIDFromHex(idBook)
	if err != nil {
		return err
	}
	_, err = col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

// DetailBook implements repository.Repository.
func (repo implRepository) DetailBook(ctx context.Context, idBook string) (repository.DetailOutputOption, error) {
	col := repo.getBookCollection()
	id, err := primitive.ObjectIDFromHex(idBook)

	if err != nil {
		return repository.DetailOutputOption{}, err
	}
	var book repository.DetailOutputOption
	if err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&book); err != nil {
		repo.l.Errorf(ctx, "repo.DetailBook:%v", err)
		return repository.DetailOutputOption{}, err
	}
	return book, nil

}

// UpdateBook implements repository.Repository.
func (repo implRepository) UpdateBook(ctx context.Context, input repository.UpdateOpition) error {
	col := repo.getBookCollection()
	id, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"title":        input.Title,
			"author":       input.Author,
			"published_at": input.PublishedAt,
			"update_at":    time.Now(),
		},
	}
	_, err = col.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil

}

// ListBook implements repository.Repository.
func (repo implRepository) ListBook(ctx context.Context, sc models.Scope, opt repository.ListOptions) ([]models.Book, paginator.Paginator, error) {
	col := repo.getBookCollection()

	filter, err := repo.buildListQuery(ctx, sc, true, opt)
	if err != nil {
		repo.l.Warnf(ctx, "book.repository.List.buildListQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	findOptions := options.Find()
	findOptions.SetSkip(opt.PaginatorQuery.Offset()) // Phân trang: bỏ qua số lượng bản ghi theo Offset
	findOptions.SetLimit(opt.PaginatorQuery.Limit)   // Phân trang: giới hạn số lượng bản ghi theo Limit
	findOptions.SetSort(bson.D{                      // Sắp xếp: theo "created_at" và "_id"
		{Key: "created_at", Value: -1}, // Sắp xếp theo trường "created_at" giảm dần
		{Key: "_id", Value: -1},        // Sắp xếp theo trường "_id" giảm dần
	})

	cursor, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		repo.l.Warnf(ctx, "book.repository.List.col.Find: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var books []models.Book
	if err := cursor.All(ctx, &books); err != nil {
		repo.l.Warnf(ctx, "book.repository.List.cursor.All: %v", err)
		return nil, paginator.Paginator{}, err
	}

	total, err := col.CountDocuments(ctx, filter)
	if err != nil {
		repo.l.Warnf(ctx, "book.repository.List.col.CountDocuments: %v", err)
		return nil, paginator.Paginator{}, err
	}

	return books, paginator.Paginator{
		Total:       total,
		Count:       int64(len(books)),
		PerPage:     opt.PaginatorQuery.Limit,
		CurrentPage: opt.PaginatorQuery.Page,
	}, nil
}

// GetRandomBook implements repository.Repository.
func (repo implRepository) GetRandomBook() (models.Book, error) {
	col := repo.getBookCollection()
	// tạo ngẫu nhiên
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result models.Book
	// ngẫu nhiên 1-9
	skip := randomSource.Intn(10)
	opts := options.FindOne().SetSkip(int64(skip))
	err := col.FindOneWithOpt(context.Background(), bson.M{}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Book{
				Title: "NOT BOOK",
			}, nil // Không có sách nào
		}
		return models.Book{}, err // Lỗi khi truy vấn
	}

	return result, nil
}
