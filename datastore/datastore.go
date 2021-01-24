package datastore

import (
    "github.com/jmoiron/sqlx"
)

type Post struct {
	ID        int
	Title     string
	Body      string
}

type PostRepository interface {
    // 表示可能な記事をidから1件取得
    GetByID(id int) (*Post, error)
    // 表示可能な記事をidの降順で取得
    List() ([]*Post, error)
}

type postRepository struct {
    db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
    return &postRepository{db}
}

func (r *postRepository) GetByID(id int) (*Post, error) {
    var post Post
    err := r.db.Get(&post, "select * from posts where id=?", id)
    if err != nil {
        return nil, err
    }
    return &post, nil
}

func (r *postRepository) List() ([]*Post, error) {
    posts := []*Post{}
    err := r.db.Select(&posts, "select * from posts order by id desc")
    if err != nil {
        return nil, err
    }
    return posts, nil
}
