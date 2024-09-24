package repository

import (
	"context"
	"errors"
	"golang-grpc/golang-grpc/proto"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	proto.UnimplementedUserServiceServer
	DB *pgxpool.Pool
}

func (r *UserRepository)ListUser(ctx context.Context, req *proto.Empty) (*proto.UserListResponse, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*proto.User
	for rows.Next() {
		var user proto.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return &proto.UserListResponse{Users: users}, nil
}

func (r *UserRepository)GetUser(ctx context.Context, req *proto.UserRequest) (*proto.User, error) {
	var user proto.User

	row := r.DB.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE id=$1", req.Id)
	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Println("Error fetching user: ", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository)CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	var user proto.User

	err := r.DB.QueryRow(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email", req.Name, req.Email).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (r *UserRepository)UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.User, error) {
	var user proto.User

	err  := r.DB.QueryRow(context.Background(), "UPDATE users SET name=$1, email=$2 WHERE id=$3 RETURNING id, name, email", req.Name, req.Email, req.Id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
            return nil, errors.New("user not found")
        }
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository)DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.Empty, error) {
	result, err := r.DB.Exec(context.Background(), "DELETE FROM users WHERE id=$1", req.Id)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, errors.New("user not found")
	}

	return &proto.Empty{}, nil
}
