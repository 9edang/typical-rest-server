package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/postgresdb"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"github.com/typical-go/typical-rest-server/pkg/reflectkit"
	"go.uber.org/dig"
)

var (
	// BookTableName is table name for book entity
	BookTableName = "books"
	// BookTable is columns for book entity
	BookTable = struct {
		ID        string
		Title     string
		Author    string
		UpdatedAt string
		CreatedAt string
	}{
		ID:        "id",
		Title:     "title",
		Author:    "author",
		UpdatedAt: "updated_at",
		CreatedAt: "created_at",
	}
)

type (
	// BookRepo to get book data from database
	BookRepo interface {
		Find(context.Context, ...dbkit.SelectOption) ([]*postgresdb.Book, error)
		Create(context.Context, *postgresdb.Book) (int64, error)
		Delete(context.Context, dbkit.DeleteOption) (int64, error)
		Update(context.Context, *postgresdb.Book, dbkit.UpdateOption) (int64, error)
		Patch(context.Context, *postgresdb.Book, dbkit.UpdateOption) (int64, error)
	}
	// BookRepoImpl is implementation book repository
	BookRepoImpl struct {
		dig.In
		*sql.DB `name:"pg"`
	}
)

// NewBookRepo return new instance of BookRepo
func NewBookRepo(impl BookRepoImpl) BookRepo {
	return &impl
}

// Find book
func (r *BookRepoImpl) Find(ctx context.Context, opts ...dbkit.SelectOption) (list []*postgresdb.Book, err error) {
	builder := sq.
		Select(
			BookTable.ID,
			BookTable.Title,
			BookTable.Author,
			BookTable.UpdatedAt,
			BookTable.CreatedAt,
		).
		From(BookTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(r)

	for _, opt := range opts {
		if builder, err = opt.CompileSelect(builder); err != nil {
			return nil, err
		}
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	list = make([]*postgresdb.Book, 0)
	for rows.Next() {
		book := new(postgresdb.Book)
		if err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.UpdatedAt,
			&book.CreatedAt,
		); err != nil {
			return
		}
		list = append(list, book)
	}
	return
}

// Create book
func (r *BookRepoImpl) Create(ctx context.Context, ent *postgresdb.Book) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	scanner := sq.
		Insert(BookTableName).
		Columns(
			BookTable.Title,
			BookTable.Author,
			BookTable.CreatedAt,
			BookTable.UpdatedAt,
		).
		Values(
			ent.Title,
			ent.Author,
			time.Now(),
			time.Now(),
		).
		Suffix(
			fmt.Sprintf("RETURNING \"%s\"", BookTable.ID),
		).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn.DB).
		QueryRowContext(ctx)

	var id int64
	if err := scanner.Scan(&id); err != nil {
		txn.SetError(err)
		return -1, err
	}
	return id, nil
}

// Update book
func (r *BookRepoImpl) Update(ctx context.Context, ent *postgresdb.Book, opt dbkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update(BookTableName).
		Set(BookTable.Title, ent.Title).
		Set(BookTable.Author, ent.Author).
		Set(BookTable.UpdatedAt, time.Now()).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn.DB)

	if builder, err = opt.CompileUpdate(builder); err != nil {
		txn.SetError(err)
		return -1, err
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}
	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}

// Patch book to update field of book if available
func (r *BookRepoImpl) Patch(ctx context.Context, ent *postgresdb.Book, opt dbkit.UpdateOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Update(BookTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn.DB)

	if !reflectkit.IsZero(ent.Title) {
		builder = builder.Set(BookTable.Title, ent.Title)
	}
	if !reflectkit.IsZero(ent.Author) {
		builder = builder.Set(BookTable.Author, ent.Author)
	}
	builder = builder.Set(BookTable.UpdatedAt, time.Now())

	if builder, err = opt.CompileUpdate(builder); err != nil {
		txn.SetError(err)
		return -1, err
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	affectedRow, err := res.RowsAffected()
	txn.SetError(err)
	return affectedRow, err
}

// Delete book
func (r *BookRepoImpl) Delete(ctx context.Context, opt dbkit.DeleteOption) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	builder := sq.
		Delete(BookTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn.DB)

	if builder, err = opt.CompileDelete(builder); err != nil {
		txn.SetError(err)
		return -1, err
	}

	res, err := builder.ExecContext(ctx)
	if err != nil {
		txn.SetError(err)
		return -1, err
	}

	return res.RowsAffected()
}
