package storage

import (
	"forum/internal/models"
	"strings"
)

func (p *PostStorage) Category() ([]models.Category, error) {
	query := `
		SELECT hashtag
		FROM hashtags;
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (p *PostStorage) GetAllPostsByCategories(category string) ([]models.Post, error) {
	query := `
		SELECT p.id, p.title, p.description,p.imageURL, u.username, p.likes, p.dislikes, p.category, p.created_at
		FROM post p
		LEFT JOIN user u
		ON u.id = p.author_id
		WHERE category LIKE '%' || $1 || '%';
	`

	rows, err := p.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	var cats string
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Image, &post.Author, &post.Likes, &post.Dislikes, &cats, &post.CreateAt)
		if err != nil {
			return nil, err
		}
		post.Category = strings.Split(cats, ",")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
