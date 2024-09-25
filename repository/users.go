package repository

import (
	"context"

	marketmodels "github.com/KingSupermarket/model/marketModels"
	usermodels "github.com/KingSupermarket/model/userModels"
	requsermodel "github.com/KingSupermarket/model/userModels/reqUserModel"
	"github.com/KingSupermarket/pkg/logger"
)

type UserRepo interface {
	SignIn(ctx context.Context, data *requsermodel.SigninModel) (*usermodels.Users, error)
	SignUp(ctx context.Context, data *usermodels.Users) (*usermodels.Users, error)
	ProfileUser(ctx context.Context, id string) (*usermodels.Users, error)
	UpdateUser(ctx context.Context, id string, data *usermodels.Users) error
	DeleteUser(ctx context.Context, id string) error
	HistoryPurchases(ctx context.Context, id string) (*usermodels.Users, []marketmodels.OrderItems, error)
}

type userRepoIml struct {
	u      UserRepo
	logger logger.Logger
}

func NewUserRepoIml(u UserRepo) *userRepoIml {
	return &userRepoIml{
		u:      u,
		logger: logger.GetLogger(),
	}
}
func (l *userRepoIml) NewSignUp(ctx context.Context, data *usermodels.Users) (*usermodels.Users, error) {
	dataSignUp, err := l.u.SignUp(ctx, data)
	if err != nil {
		l.logger.Errorf("SignUp failed for email %s: %v", data.Email, err)
		return nil, err
	}
	l.logger.Infof("User signed up successfully: %s", data.Email)
	return dataSignUp, nil
}
func (l *userRepoIml) NewHistoryPurchases(ctx context.Context, id string) (*usermodels.Users, []marketmodels.OrderItems, error) {
	dataProfile, history, err := l.u.HistoryPurchases(ctx, id)
	if err != nil {
		l.logger.Errorf("Error retrieving purchase history for user ID %s: %v", id, err)
		return nil, nil, err
	}
	l.logger.Infof("Retrieved purchase history for user ID %s", id)
	return dataProfile, history, nil
}

func (l *userRepoIml) NewProfileUser(ctx context.Context, id string) (*usermodels.Users, error) {
	dataProfile, err := l.u.ProfileUser(ctx, id)
	if err != nil {
		l.logger.Errorf("Error retrieving profile for user ID %s: %v", id, err)
		return nil, err
	}
	l.logger.Infof("Retrieved profile for user ID %s", id)
	return dataProfile, nil
}

func (l *userRepoIml) NewSignIn(ctx context.Context, data *requsermodel.SigninModel) (*usermodels.Users, error) {
	dataUser, err := l.u.SignIn(ctx, data)
	if err != nil {
		l.logger.Errorf("SignIn failed for email %s: %v", data.Email, err)
		return nil, err
	}
	l.logger.Infof("User signed in successfully: %s", data.Email)
	return dataUser, nil
}


func (l *userRepoIml) NewUpdateUser(ctx context.Context, id string, data *usermodels.Users) error {
	if err := l.u.UpdateUser(ctx, id, data); err != nil {
		l.logger.Errorf("Update user failed for user ID %s: %v", id, err)
		return err
	}
	l.logger.Infof("User updated successfully: %s", id)
	return nil
}

func (l *userRepoIml) NewDeleteUser(ctx context.Context, id string) error {
	if err := l.u.DeleteUser(ctx, id); err != nil {
		l.logger.Errorf("Failed to delete user with ID %s: %v", id, err)
		return err
	}
	l.logger.Infof("Deleted user with ID %s successfully", id)
	return nil
}
