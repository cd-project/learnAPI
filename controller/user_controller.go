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

var usingRole string = "client"

type UserController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	ChangeRole(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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

// ChangePassword changes password of user with UID
// @tags user-manager-apis
// @Summary change password
// @Description change password
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UserPasswordInfo body controller.UserPasswordPayload true "User and password info"
// @Success 200 {object} model.Response
// @Router /user/modify/pwd [put]
func (c *userController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var jsonResponse *model.Response
	var data UserPasswordPayload

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)

		jsonResponse = &model.Response{
			Data:    nil,
			Message: "Bad Request",
			Code:    "400",
			Success: false,
		}
		render.JSON(w, r, jsonResponse)
		return
	}

	userResponse, err := c.userService.ChangePassword(data.UserID, data.NewPassword, data.OldPassword)
	if data.NewPassword == data.OldPassword {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: "New password cannot be the same as old password. Please choose another password!",
			Code:    "400",
			Success: false,
		}
		render.JSON(w, r, jsonResponse)
		return
	}
	if err != nil {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: "failed!" + err.Error(),
			Code:    "200",
			Success: false,
		}
	} else {
		jsonResponse = &model.Response{
			Data:    userResponse,
			Message: "Password change successfully!",
			Code:    "400",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)

}

// ChangeRole changes role of user with UID
// @tags user-manager-apis
// @Summary change role
// @Description change role
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UserRoleInfo body controller.UserRolePayload true "UserID and role"
// @Success 200
// @Router /user/modify/role [put]
func (c *userController) ChangeRole(w http.ResponseWriter, r *http.Request) {
	var jsonResponse *model.Response
	var data UserRolePayload

	decoder := json.NewDecoder(r.Body)
	// bad request
	if err := decoder.Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		jsonResponse = &model.Response{
			Data:    nil,
			Message: "Bad Request",
			Code:    "400",
			Success: false,
		}
		render.JSON(w, r, jsonResponse)
		return
	}

	// request valid, check for same role
	if userdata, err := c.userService.GetByID(data.UserID); userdata.Role == data.NewRole && err == nil {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: "New role cannot be the same as old role. Please choose another role!",
			Code:    "200",
			Success: false,
		}
		render.JSON(w, r, jsonResponse)
		return
	}
	// valid NewRole
	userResponse, err := c.userService.ChangeRole(data.UserID, data.NewRole)
	if err != nil {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: err.Error(),
			Code:    "200",
			Success: false,
		}
	} else {
		jsonResponse = &model.Response{
			Data:    userResponse,
			Message: "Role changed. Your new role is: " + data.NewRole,
			Code:    "200",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)

}

// ResetPassword resets password of user with UID
// @tags user-manager-apis
// @Summary reset password
// @Description reset password
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uid path integer true "User ID"
// @Success 200 {object} model.Response
// @Router /user/reset/{uid} [put]
func (c *userController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var jsonResponse *model.Response

	strID := chi.URLParam(r, "uid")
	uid, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, http.StatusText(404), 404)
		jsonResponse = &model.Response{
			Data:    nil,
			Message: "Bad URL",
			Code:    "404",
			Success: false,
		}
		render.JSON(w, r, jsonResponse)
		return
	}

	// get userID from URL done, validate user
	userData, err := c.userService.ResetPassword(uid)
	if err != nil {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: err.Error(),
			Code:    "200",
			Success: false,
		}
	} else {
		jsonResponse = &model.Response{
			Data:    userData,
			Message: "Reset password successfully. Your new password is 0000.",
			Code:    "200",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)
}

// DeleteUser deletes user with UserID
// @tags user-manager-apis
// @Summary delete user
// @Description delete user
// @Accept json
// @Produce json
// @Param uid path integer true "User ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.Response
// @Router /user/delete/{uid} [delete]
func (c *userController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var jsonResponse *model.Response

	strID := chi.URLParam(r, "uid")
	uid, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, http.StatusText(404), 404)
		jsonResponse = &model.Response{
			Data:    nil,
			Message: "Bad URL",
			Code:    "404",
			Success: false,
		}
		render.JSON(w, r, jsonResponse)
		return
	}

	// get UserID from URL done, deleteUser
	if err := c.userService.DeleteUser(uid); err != nil {
		jsonResponse = &model.Response{
			Data:    nil,
			Message: err.Error(),
			Code:    "200",
			Success: false,
		}
	} else {
		jsonResponse = &model.Response{
			Message: "User successfully deleted!",
			Code:    "200",
			Success: true,
		}
	}
	render.JSON(w, r, jsonResponse)
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
	usingRole = user.Role
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

type UserPasswordPayload struct {
	UserID      int    `json:"userID"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserRolePayload struct {
	UserID  int    `json:"userID"`
	NewRole string `json:"newRole"`
}
