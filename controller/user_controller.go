package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo/infrastructure"
	"todo/model"
	"todo/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/harlow/authtoken"
)

type UserController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	ChangeRole(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	LoginWithToken(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userService service.UserService
}

// GetAll gets all users currently in table "users"
// @tags user-manager-apis
// @Summary get all users
// @Description get all users
// @Accept json
// @Produce json
// @Success 200 {object} model.Response
// @Router /user/all [get]
func (c *userController) GetAll(w http.ResponseWriter, r *http.Request) {
	var jsonResponse *model.Response
	users, err := c.userService.GetAll()
	if err != nil {
		log.Println(err.Error())
		jsonResponse = &model.Response{
			Data:    nil,
			Message: err.Error(),
			Code:    "failed",
			Success: false,
		}
	} else {
		jsonResponse = &model.Response{
			Data:    users,
			Message: "Get All Successful",
			Code:    "200",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)
}

// GetByID gets user info of user with UID in URL
// @tags user-manager-apis
// @Summary gets user info
// @Description gets user info
// @Accept json
// @Produce json
// @Param uid path integer true "User ID"
// @Success 200 {object} model.Response
// @Router /user/{uid} [get]
func (c *userController) GetByID(w http.ResponseWriter, r *http.Request) {
	var jsonResponse *model.Response

	strID := chi.URLParam(r, "uid")
	uid, err := strconv.Atoi(strID)
	// error converting ID(string) to int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		infrastructure.ErrLog.Println(err)
		jsonResponse = &model.Response{
			Data:    nil,
			Message: "Bad UID in URL",
			Code:    "400",
			Success: false,
		}
		render.JSON(w, r, jsonResponse)
	}
	userInfo, err := c.userService.GetByID(uid)
	if err != nil {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: err.Error(),
			Code:    "404",
			Success: false,
		}
	} else {
		jsonResponse = &model.Response{
			Data:    userInfo,
			Message: "UID " + strID + " :Get info successful",
			Code:    "200",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)

}

// CreateUser creates an user with given data
// @tags user-manager-apis
// @Summary	creates new user
// @Description creates new user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UserInfo body model.UserPayload true "User information"
// @Success 200 {object} model.Response
// @Router /user/create [post]
func (c *userController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInfo *model.UserPayload
	var jsonResponse *model.Response
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInfo); err != nil {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: "error decoding data" + err.Error(),
			Code:    "400",
			Success: false,
		}
		render.JSON(w, r, jsonResponse)
		return
	}

	newUser, err := c.userService.CreateUser(&model.User{
		Username: userInfo.Username,
		Password: userInfo.Password,
		Role:     userInfo.Role,
	})
	if err != nil {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: err.Error(),
			Code:    "500",
			Success: false,
		}

	} else {
		jsonResponse = &model.Response{
			Data:    newUser,
			Message: "User created",
			Code:    "200",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)
}

func (c *userController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// // just change password, what's so difference?
	// strID := chi.URLParam(r, "uid")
	// uid, _ := strconv.Atoi(strID)
	panic("")
}

func (c *userController) ChangeRole(w http.ResponseWriter, r *http.Request) {
	panic("")
}

func (c *userController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	panic("")
}

func (c *userController) Delete(w http.ResponseWriter, r *http.Request) {
	panic("")
}

// Login log user in if they have valid credential
// @tags user-manager-apis
// @Summary log user in
// @Description log user in
// @Accept json
// @Produce json
// @Param LoginPayload body controller.LoginPayload true "username & password"
// @Success 200
// @Router /user/login [post]
func (c *userController) Login(w http.ResponseWriter, r *http.Request) {
	var jsonResponse *model.LoginResponse

	var loginDetail LoginPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginDetail); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		jsonResponse = &model.LoginResponse{
			Token:        "",
			RefreshToken: "",
			Data:         nil,
			Message:      "Bad Request(loginDetail)",
			Code:         "400",
			Success:      false,
		}
		render.JSON(w, r, jsonResponse)
		return
	}
	// loginDetail model valid.
	user, token, refreshToken, err := c.userService.LoginRequest(loginDetail.Username, loginDetail.Password)
	if err != nil {
		jsonResponse = &model.LoginResponse{
			Token:        token,
			RefreshToken: refreshToken,
			Data:         nil,
			Message:      "Invalid Credentials, " + err.Error(),
			Code:         "400",
			Success:      false,
		}
	} else {
		jsonResponse = &model.LoginResponse{
			Token:        token,
			RefreshToken: refreshToken,
			Data: &model.UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Role:     user.Role,
			},
			Message: "Login successful " + user.Role,
			Code:    "200",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)
}

// LoginWithToken provides token each login attempt
// @tags user-manager-apis
// @Summary login user
// @Description login user, return new token string jwt
// @Accept json
// @Produce json
// @Success 200 {object} model.LoginResponse
// @Router /user/login/token [post]
func (c *userController) LoginWithToken(w http.ResponseWriter, r *http.Request) {
	var jsonResponse *model.LoginResponse

	token, err := authtoken.FromRequest(r)
	log.Println(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		// http.Error(w, http.StatusText(401), 401)
		return
	}

	user, token, success, err := c.userService.LoginWithToken(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if !success {
		jsonResponse = &model.LoginResponse{
			Token:   token,
			Data:    nil,
			Message: err.Error(),
			Code:    "400",
			Success: false,
		}
	} else {
		jsonResponse = &model.LoginResponse{
			Token: token,
			Data: &model.UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Role:     user.Role,
			},
			Message: "Login Successful",
			Code:    "200",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)

}

func NewUserController() UserController {
	userService := service.NewUserService()
	return &userController{
		userService: userService,
	}
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
