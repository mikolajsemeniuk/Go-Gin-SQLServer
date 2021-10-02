package services

import (
	"context"
	"database/sql"
	"errors"
	"go-gin-sqlserver/program/database"
	"go-gin-sqlserver/program/payloads"
	"time"
)

const (
	GET_POSTLIKES_QUERY     = "select username from users join post_likes on users.user_id = post_likes.user_id where post_id = @PostId;"
	EXISTS_POST_QUERY       = "select count(*) from posts where user_id = @UserId and post_id = @PostId;"
	EXISTS_POSTLIKE_QUERY   = "select count(*) from post_likes where user_id = @UserId and post_id = @PostId;"
	ADD_POSTLIKE_COMMAND    = "insert into post_likes (user_id, post_id) values (@UserId, @PostId);"
	REMOVE_POSTLIKE_COMMAND = "delete from post_likes where user_id = @UserId and post_id = @PostId;"
)

func GetPostLikes(postId int64) ([]payloads.PostLike, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var likes []payloads.PostLike
	rows, error := database.Client.QueryContext(context, GET_POSTLIKES_QUERY, sql.Named("PostId", postId))
	if error != nil {
		return likes, error
	}
	defer rows.Close()

	for rows.Next() {
		var like payloads.PostLike
		if error := rows.Scan(&like.Username); error != nil {
			return likes, error
		}
		likes = append(likes, like)
	}

	if error = rows.Err(); error != nil {
		return likes, error
	}

	if len(likes) == 0 {
		likes = []payloads.PostLike{}
	}

	return likes, nil
}

func PostExists(userId int64, postId int64) (bool, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exists := 0

	if _, error := GetUser(userId); error != nil {
		return false, error
	}

	if _, error := GetPost(postId); error != nil {
		return false, error
	}

	if error := database.Client.
		QueryRowContext(context, EXISTS_POST_QUERY, sql.Named("UserId", userId), sql.Named("PostId", postId)).
		Scan(&exists); error != nil {
		return false, error
	}

	return exists != 0, nil
}

func PostLikeExists(userId int64, postId int64) (bool, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exists := 0

	if _, error := GetUser(userId); error != nil {
		return false, error
	}

	if _, error := GetPost(postId); error != nil {
		return false, error
	}

	if error := database.Client.
		QueryRowContext(context, EXISTS_POSTLIKE_QUERY, sql.Named("UserId", userId), sql.Named("PostId", postId)).
		Scan(&exists); error != nil {
		return false, error
	}

	return exists != 0, nil
}

func SetPostLike(userId int64, postId int64) error {
	exists, error := PostExists(userId, postId)
	if error != nil {
		return error
	}

	if exists {
		return errors.New("you cannot like your own posts")
	}

	exists, error = PostLikeExists(userId, postId)
	if error != nil {
		return error
	}

	query := ADD_POSTLIKE_COMMAND
	if exists {
		query = REMOVE_POSTLIKE_COMMAND
	}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement, error := database.Client.PrepareContext(context, query)
	if error != nil {
		return error
	}
	defer statement.Close()

	result, error := statement.ExecContext(context, sql.Named("UserId", userId), sql.Named("PostId", postId))
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
