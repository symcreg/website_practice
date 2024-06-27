package handler

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"strings"
	"test/internal/middleware"
	"test/internal/model"
	"test/internal/utility"
)

type UserHandler struct {
	userNameExp *regexp.Regexp
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		namePattern     = `^[a-zA-Z0-9._%+-]{3,}$`
		emailPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		passwordPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`
	)
	nameExp := regexp.MustCompile(namePattern, regexp.None)
	emailExp := regexp.MustCompile(emailPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordPattern, regexp.None)
	return &UserHandler{
		userNameExp: nameExp,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) Register(router *gin.Engine) {
	ug := router.Group("/user")
	protected := ug.Group("/protected")
	protected.Use(middleware.JwtToken())
	ug.POST("/register", u.RegisterUser)
	protected.POST("/login", u.Login)
	protected.GET("/profile", u.Profile)
	protected.GET("/logout", u.Logout)
	protected.PUT("/profile_update", u.UpdateProfile)
	protected.PUT("/change_password", u.ChangePassword)
	//protected.GET("/cancel_account", u.CancelAccount)
}

func (u *UserHandler) Profile(c *gin.Context) {
	email, _ := c.Get("email")
	if email == nil || email == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	type ProfileResponse struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	user, err := model.GetUserByEmail(email.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, ProfileResponse{Name: user.Name, Email: user.Email})
}
func (u *UserHandler) RegisterUser(c *gin.Context) {
	type RegisterRequest struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid email"})
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(400, gin.H{"error": "Password and confirm password do not match"})
		return
	}
	// Check if email is already registered
	var user model.User
	if model.IsUserExist(req.Email) {
		c.JSON(400, gin.H{"error": "Email is already registered"})
		return
	}
	// convert password to hash
	user.Name = req.Name
	user.Email = req.Email
	user.Password = utility.HashPassword(req.Password)
	user.ResisterTime = utility.GetCurrentTime()
	// Save user to database
	if err := model.InsertUser(&user); err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "User registered successfully"})
}
func (u *UserHandler) Login(c *gin.Context) {
	email, _ := c.Get("email")
	if email != nil && email != "" {
		c.JSON(400, gin.H{"error": "Already logged in"})
		return
	}
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid email or password"})
		return

	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return

	}
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid password or password"})
		return
	}
	user, err := model.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid email or password"})
		return
	}
	if !utility.VerifyPassword(user.Password, req.Password) {
		c.JSON(400, gin.H{"error": "Invalid email or password"})
		return
	}
	c.JSON(200, gin.H{"message": "Login successful"})
	// Generate token
	type LoginResponse struct {
		Token string `json:"token"`
	}
	token, err := utility.GenerateToken(req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, LoginResponse{Token: token})
}
func (u *UserHandler) Logout(c *gin.Context) {
	email, _ := c.Get("email")
	if email == nil || email == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	auth := c.GetHeader("Authorization")
	token := auth[len("Bearer"):]
	token = strings.TrimSpace(token)
	err := utility.DeleteToken(token)
	if err != nil {
		println(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "Logout successful"})
}
func (u *UserHandler) UpdateProfile(c *gin.Context) {
	email, _ := c.Get("email")
	if email == nil || email == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := model.GetUserByEmail(email.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	type UpdateProfileRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return

	}
	if req.Email != "" {
		ok, err := u.emailExp.MatchString(req.Email)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		if !ok {
			c.JSON(400, gin.H{"error": "Invalid email"})
			return
		}
		user.Email = req.Email
	}
	if req.Name != "" {
		ok, err := u.userNameExp.MatchString(req.Name)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		if !ok {
			c.JSON(400, gin.H{"error": "Invalid name"})
			return
		}
		user.Name = req.Name
	}
	if err := model.UpdateUser(user); err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "Profile updated successfully"})

}
func (u *UserHandler) ChangePassword(c *gin.Context) {
	email, _ := c.Get("email")
	if email == nil || email == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := model.GetUserByEmail(email.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	type ChangePasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return

	}
	ok, err := u.passwordExp.MatchString(req.NewPassword)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return

	}
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}
	if !utility.VerifyPassword(user.Password, req.OldPassword) {
		c.JSON(400, gin.H{"error": "Wrong old password"})
		return
	}
	user.Password = utility.HashPassword(req.NewPassword)
	if err := model.UpdateUser(user); err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return

	}
	c.JSON(200, gin.H{"message": "Password changed successfully"})
}

//func (u *UserHandler) CancelAccount(c *gin.Context) {
//	email, _ := c.Get("email")
//	if email == nil || email == "" {
//		c.JSON(401, gin.H{"error": "Unauthorized"})
//		return
//	}
//	user, err := model.GetUserByEmail(email.(string))
//	if err != nil {
//		c.JSON(500, gin.H{"error": "Internal server error"})
//		return
//	}
//	if err := model.DeleteUser(user); err != nil {
//		c.JSON(500, gin.H{"error": "Internal server error"})
//		return
//	}
//	c.JSON(200, gin.H{"message": "Account deleted successfully"})
//}
