package repository

import (
	"context"

	usermodels "github.com/KingSupermarket/model/userModels"
	requsermodel "github.com/KingSupermarket/model/userModels/reqUserModel"
)

type UserRepo interface {
	SignIn(ctx context.Context, data *requsermodel.SigninModel) (*usermodels.Users, error)
	SignUp(ctx context.Context, data *usermodels.Users) (*usermodels.Users, error)
	ProfileUser(ctx context.Context, id string) (*usermodels.Users, error)
	UpdateUser(ctx context.Context, id string, data *usermodels.Users) error
	DeleteUser(ctx context.Context, id string) error
}

type userRepoIml struct {
	u UserRepo
}

func NewUserRepoIml(u UserRepo) *userRepoIml {
	return &userRepoIml{u: u}
}
func (l *userRepoIml) NewProfileUser(ctx context.Context, id string) (*usermodels.Users, error) {
	dataProfile, err := l.u.ProfileUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return dataProfile, nil
}
func (l *userRepoIml) NewSignIn(ctx context.Context, data *requsermodel.SigninModel) (*usermodels.Users, error) {
	dataUser, err := l.u.SignIn(ctx, data)
	if err != nil {
		return nil, err
	}
	return dataUser, nil
}
func (l *userRepoIml) NewSignUp(ctx context.Context, data *usermodels.Users) (*usermodels.Users, error) {
	dataSignUp, err := l.u.SignUp(ctx, data)
	if err != nil {
		return nil, err
	}
	return dataSignUp, nil
}
func (l *userRepoIml) NewUpdateUser(ctx context.Context, id string, data *usermodels.Users) error {
	if err := l.u.UpdateUser(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (l *userRepoIml) NewDeleteUser(ctx context.Context, id string) error {
	if err := l.u.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}
