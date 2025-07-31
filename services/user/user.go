package user

import (
	"errors"
	"nexus/models"
	"nexus/services/auth"
	"time"
)

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UserDetails struct {
	Account       string    `json:"account"`
	Username      string    `json:"username"`
	Region        string    `json:"region"`
	Location      *Point    `json:"location"`
	Icon          string    `json:"icon"`
	Role          string    `json:"role"`
	LastLoginTime time.Time `json:"last_login_time"`
	SignUpAt      time.Time `json:"sign_up_at"`
}

type UserService struct {
	userModels models.IUserModels
}

type ModifyUserRequest struct {
	ID       string  `json:"id"`
	Username *string  `json:"username"`
	Location *Point  `json:"location"`
	Icon     *string `json:"icon"`
}

type SignUpRequest struct {
	Account  string  `json:"account"`
	Password string  `json:"password"`
	Username string  `json:"username"`
	Region   string  `json:"region"`
	Location *Point  `json:"location"`
}

type SignUpResponse struct {
	Token string      `json:"token"`
	User  *UserDetails `json:"user"`
}

type SignInRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string      `json:"token"`
	User  *UserDetails `json:"user"`
}

type UserResponse struct {
	User *UserDetails `json:"user"`
}

type UserSearchRequest struct {
	LikeAccount string `json:"like_account"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	OrderBy    string `json:"order_by"`
}

type UserSearchResponse struct {
	Total int64           `json:"total"`
	Users []UserDetails `json:"users"`
}

type IUserService interface {
	SignUp(request SignUpRequest) (SignUpResponse, error)
	SignIn(request SignInRequest) (SignInResponse, error)
	Get(id string) (UserResponse, error)
	GetByAccount(account string) (UserResponse, error)
	SearchByAccount(request UserSearchRequest) (UserSearchResponse, error)
	Update(modify ModifyUserRequest) error
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

func (service *UserService) SignUp(request SignUpRequest) (SignUpResponse, error) {
	user := models.User{
		ID:            models.GenerateID(),
		Account:       request.Account,
		Password:      request.Password,
		Username:      request.Username,
		Region:        request.Region,
		Location:      toModelLocation(request.Location),
		Icon:          "",
		Role:          "user",
		Status:        true,
		LastLoginTime: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err := service.userModels.Create(&user)
	if err != nil {
		return SignUpResponse{}, err
	}

	token, err := auth.GenerateToken("HS256", user.ID, 3600)
	if err != nil {
		return SignUpResponse{}, err
	}

	return SignUpResponse{
		Token: token,
		User: &UserDetails{
			Account:       user.Account,
			Username:      user.Username,
			Region:        user.Region,
			Location:      toUserLocation(user.Location),
			Icon:          user.Icon,
			Role:          user.Role,
			LastLoginTime: user.LastLoginTime,
			SignUpAt:      user.CreatedAt,
		},
	}, nil
}

func (service *UserService) SignIn(request SignInRequest) (SignInResponse, error) {
	user, err := service.userModels.GetByAccount(request.Account)
	if err != nil {
		return SignInResponse{}, errors.New("account and password is not correct")
	}
	if user.Password != request.Password {
		return SignInResponse{}, errors.New("account and password is not correct")
	}
	token, err := auth.GenerateToken("HS256", user.ID, 3600)
	if err != nil {
		return SignInResponse{}, err
	}
	return SignInResponse{Token: token, User: &UserDetails{
		Account:       user.Account,
		Username:      user.Username,
		Region:        user.Region,
		Location:      &Point{Latitude: user.Location.Latitude, Longitude: user.Location.Longitude},
		Icon:          user.Icon,
		Role:          user.Role,
		LastLoginTime: user.LastLoginTime,
		SignUpAt:      user.CreatedAt,
	}}, nil
}

func (service *UserService) Get(id string) (UserResponse, error) {
	user, err := service.userModels.Get(id)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		User: &UserDetails{
		Account:       user.Account,
		Username:      user.Username,
		Region:        user.Region,
		Location:      &Point{Latitude: user.Location.Latitude, Longitude: user.Location.Longitude},
		Icon:          user.Icon,
		Role:          user.Role,
		LastLoginTime: user.LastLoginTime,
		SignUpAt:      user.CreatedAt,
	},
	}, nil
}

func (service *UserService) GetByAccount(account string) (UserResponse, error) {
	user, err := service.userModels.GetByAccount(account)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		User: &UserDetails{
		Account:       user.Account,
		Username:      user.Username,
		Region:        user.Region,
		Location:      &Point{Latitude: user.Location.Latitude, Longitude: user.Location.Longitude},
		Icon:          user.Icon,
		Role:          user.Role,
		LastLoginTime: user.LastLoginTime,
		SignUpAt:      user.CreatedAt,
	},
	}, nil
}

func (service *UserService) SearchByAccount(request UserSearchRequest) (UserSearchResponse, error) {
	total, err := service.userModels.Count(request.LikeAccount)
	if err != nil {
		return UserSearchResponse{}, err
	}

	users, err := service.userModels.List(request.LikeAccount, request.Limit, request.Offset, request.OrderBy)
	if err != nil {
		return UserSearchResponse{}, err
	}
	var userDetails []UserDetails
	for _, user := range users {
		userDetails = append(userDetails, UserDetails{
			Account:       user.Account,
			Username:      user.Username,
			Region:        user.Region,
			Location:      &Point{Latitude: user.Location.Latitude, Longitude: user.Location.Longitude},
			Icon:          user.Icon,
			Role:          user.Role,
			LastLoginTime: user.LastLoginTime,
			SignUpAt:      user.CreatedAt,
		})
	}
	return UserSearchResponse{
		Users: userDetails,
		Total: total,
	}, nil
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
		user.Location = &models.Point{
			Latitude:  modify.Location.Latitude,
			Longitude: modify.Location.Longitude,
		}
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

func toUserLocation(location *models.Point) *Point {
	if location == nil {
		return nil
	}
	return &Point{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
}

func toModelLocation(location *Point) *models.Point {
	if location == nil {
		return nil
	}
	return &models.Point{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
}
