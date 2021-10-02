package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-gin-sqlserver/program/database"
	"go-gin-sqlserver/program/inputs"
	"go-gin-sqlserver/program/payloads"
	"time"
)

const (
	GET_POSTS_QUERY     = "SELECT [post_id], [title] FROM [db].[dbo].[posts] WHERE [user_id] = @UserId;"
	GET_POST_QUERY      = "SELECT [post_id], [title] FROM [db].[dbo].[posts] WHERE [post_id] = @PostId;"
	ADD_POST_COMMAND    = "INSERT INTO [db].[dbo].[posts] ([user_id], [title]) VALUES (@UserId, @Title);"
	UPDATE_POST_COMMAND = "UPDATE [db].[dbo].[posts] SET [title] = @Title WHERE [post_id] = @PostId;"
	REMOVE_POST_COMMAND = "DELETE FROM [db].[dbo].[posts] WHERE [post_id] = @PostId;"
)

func GetPosts(userId int64) ([]payloads.Post, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var posts []payloads.Post
	rows, error := database.Client.QueryContext(context, GET_POSTS_QUERY, sql.Named("UserId", userId))
	if error != nil {
		return posts, error
	}
	defer rows.Close()

	for rows.Next() {
		var post payloads.Post
		if error := rows.Scan(&post.PostId, &post.Title); error != nil {
			return posts, error
		}

		likes, error := GetPostLikes(post.PostId)
		if error != nil {
			return posts, error
		}

		if len(likes) == 0 {
			likes = []payloads.PostLike{}
		}

		post.Likes = likes
		posts = append(posts, post)
	}

	if error = rows.Err(); error != nil {
		return posts, error
	}

	if len(posts) == 0 {
		posts = []payloads.Post{}
	}

	return posts, nil
}

func GetPost(postId int64) (payloads.Post, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var post payloads.Post
	if error := database.Client.
		QueryRowContext(context, GET_POST_QUERY, sql.Named("PostId", postId)).
		Scan(&post.PostId, &post.Title); error != nil {
		if error == sql.ErrNoRows {
			return payloads.Post{}, fmt.Errorf("no post with id of: %d", postId)
		} else {
			return payloads.Post{}, error
		}
	}

	return post, nil
}

func AddPost(userId int64, input inputs.Post) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement, error := database.Client.PrepareContext(context, ADD_POST_COMMAND)
	if error != nil {
		return error
	}
	defer statement.Close()

	result, error := statement.ExecContext(context, sql.Named("UserId", userId), sql.Named("Title", input.Title))
	if error != nil {
		return error
	}

	rows, error := result.RowsAffected()
	if error != nil {
		return error
	}

	if rows == 0 {
		return errors.New("something went wrong")
	}

	return nil
}

func UpdatePost(postId int64, input inputs.Post) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, error := GetPost(postId); error != nil {
		return error
	}

	statement, error := database.Client.PrepareContext(context, UPDATE_POST_COMMAND)
	if error != nil {
		return error
	}

	result, error := statement.ExecContext(context, sql.Named("PostId", postId), sql.Named("Title", input.Title))
	if error != nil {
		return error
	}

	rows, error := result.RowsAffected()
	if error != nil {
		return error
	}

	if rows == 0 {
		return errors.New("something went wrong")
	}

	return nil
}

func RemovePost(postId int64) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, error := GetPost(postId); error != nil {
		return error
	}

	statement, error := database.Client.PrepareContext(context, REMOVE_POST_COMMAND)
	if error != nil {
		return error
	}

	result, error := statement.ExecContext(context, sql.Named("PostId", postId))
	if error != nil {
		return error
	}

	rows, error := result.RowsAffected()
	if error != nil {
		return error
	}

	if rows == 0 {
		return errors.New("something went wrong")
	}

	return nil
}
