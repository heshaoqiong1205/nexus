package user

import (
	"errors"
	"nexus/models"
	"nexus/services/auth"
)

type UserDetails struct {
	Account       string `json:"account"`
	Username      string `json:"username"`
	Region        string `json:"region"`
	Location      string `json:"location"`
	Icon          string `json:"icon"`
	Role          string `json:"role"`
	LastLoginTime string `json:"last_login_time"`
	SginUpAt      string `json:"sgin_up_at"`
}

type UserService struct {
	userModels models.IUserModels
}

type ModifyUserRequest struct {
	ID       string  `json:"id"`
	Username *string `json:"username"`
	Location *string `json:"location"`
	Icon     *string `json:"icon"`
}

type IUserService interface {
	SginUp(account string, password string, username string, region string) (models.User, error)
	SginIn(account string, password string) ([]models.User, error)
}

func NewUserService() *UserService {
	return &UserService{
		userModels: &models.UserModels{},
	}
}

func NewUserServiceWithModel(userModels models.IUserModels) *UserService {
	return &UserService{
		userModels: userModels,
	}
}

func (service *UserService) SginUp(account string, password string, username string, region string) (*UserDetails, error) {
	user := models.User{
		Account:  account,
		Password: password,
		Username: username,
		Region:   region,
	}
	err := service.userModels.Create(&user)
	if err != nil {
		return nil, err
	}
	return &UserDetails{
		Account:       user.Account,
		Username:      user.Username,
		Region:        user.Region,
		LastLoginTime: user.LastLoginTime.String(),
		SginUpAt:      user.CreatedAt.String(),
	}, nil
}

func (service *UserService) SginIn(account string, password string) (string, error) {
	user, err := service.userModels.GetByAccount(account)
	if err != nil {
		return "", errors.New("account and password is not correct")
	}
	if user.Password != password {
		return "", errors.New("account and password is not correct")
	}
	return auth.GenerateToken("HS256", user.ID, 3600)
}

func (service *UserService) Get(id string) (*UserDetails, error) {
	user, err := service.userModels.Get(id)
	if err != nil {
		return nil, err
	}
	return &UserDetails{
		Account:       user.Account,
		Username:      user.Username,
		Region:        user.Region,
		Location:      user.Location,
		Icon:          user.Icon,
		Role:          user.Role,
		LastLoginTime: user.LastLoginTime.String(),
		SginUpAt:      user.CreatedAt.String(),
	}, nil
}

func (service *UserService) GetByAccount(account string) (*UserDetails, error) {
	user, err := service.userModels.GetByAccount(account)
	if err != nil {
		return nil, err
	}
	return &UserDetails{
		Account:       user.Account,
		Username:      user.Username,
		Region:        user.Region,
		Location:      user.Location,
		Icon:          user.Icon,
		Role:          user.Role,
		LastLoginTime: user.LastLoginTime.String(),
		SginUpAt:      user.CreatedAt.String(),
	}, nil
}

func (service *UserService) SearchByAccount(LikeAccount string, limt int, offset int, orderBY string) ([]UserDetails, error) {
	users, err := service.userModels.List(LikeAccount, limt, offset, orderBY)
	if err != nil {
		return nil, err
	}
	var userDetails []UserDetails
	for _, user := range users {
		userDetails = append(userDetails, UserDetails{
			Account:       user.Account,
			Username:      user.Username,
			Region:        user.Region,
			Location:      user.Location,
			Icon:          user.Icon,
			Role:          user.Role,
			LastLoginTime: user.LastLoginTime.String(),
			SginUpAt:      user.CreatedAt.String(),
		})
	}
	return userDetails, nil
}

func (service *UserService) Update(modify ModifyUserRequest) error {
	user, err := service.userModels.Get(modify.ID)
	if err != nil {
		return err
	}
	if modify.Username != nil {
		user.Username = *modify.Username
	}
	if modify.Location != nil {
		user.Location = *modify.Location
	}
	if modify.Icon != nil {
		user.Icon = *modify.Icon
	}
	err = service.userModels.Update(&user)
	if err != nil {
		return err
	}
	return nil
}
