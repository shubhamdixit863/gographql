package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saisri/gographql/graph/generated"
	"github.com/saisri/gographql/graph/model"
	"github.com/saisri/gographql/internal/auth"
	"github.com/saisri/gographql/internal/models"
	"github.com/saisri/gographql/internal/pkg/jwt"
)

func (r *mutationResolver) CreateBook(ctx context.Context, input model.NewBook) (*model.Book, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Book{}, fmt.Errorf("access denied")
	}

	var book models.Book
	book.Title = input.Title
	book.User.ID = user.ID
	graphqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}

	book.Save()
	return &model.Book{Title: book.Title, User: graphqlUser}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user models.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user models.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		// 1
		return "", &models.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (r *mutationResolver) GetUserDetails(ctx context.Context) (*model.User, error) {
	userauth := auth.ForContext(ctx)
	if userauth == nil {
		return nil, fmt.Errorf("access denied")
	}

	var user models.User
	dbUSer, _ := user.GetUserByUsername(user.Username)
	return &model.User{ID: string(dbUSer)}, nil
}

func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}
	var resultBooks []*model.Book
	var book models.Book
	var dbBooks []models.Book
	dbBooks = book.GetAll()
	for _, book := range dbBooks {
		graphqlUser := &model.User{
			ID:   book.User.ID,
			Name: book.User.Username,
		}
		resultBooks = append(resultBooks, &model.Book{ID: book.ID, Title: book.Title, User: graphqlUser})
	}
	return resultBooks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
