package repository

import (
	"auth-service/database"
	"auth-service/logger"
	"auth-service/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

var jwtKey = []byte("my_secret")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthRepository interface {
	Register(models.AuthUser, context.Context) (error, string)
	Login(models.AuthUser, context.Context) (string, string, time.Time, error)
	ChangePassword(models.AuthUser, context.Context) (string, error)
	ChangeUsername(string, string, context.Context) (error, string)
	UpdateUser(models.AuthUser, string, context.Context) (models.AuthUser, string, error)
	DeleteUser(string, context.Context) (error, string)
}

type dbRepo struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &dbRepo{db: db}
}

func (d *dbRepo) ChangePassword(updatedUser models.AuthUser, ctx context.Context) (string, error) {
	loggerInst := ctx.Value(logger.LoggerKey).(*logger.LogInstance)
	var errMsg string
	var checkUser models.AuthUser
	result := database.DB.DB.Where("username = ?", updatedUser.Username).First(&checkUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errMsg = "The user does not exist, you cant change password.."
			loggerInst.Error(ctx, "The user does not exist, you cant change password..", result.Error)
		} else {
			// another error is occured
			errMsg = "Error checking user registration"
			loggerInst.Error(ctx, "Error checking user registration", result.Error)
		}
		return errMsg, result.Error
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		loggerInst.Error(ctx, "Error in generating hashed password.")
		errMsg = "Error in generating hashed password."
		return errMsg, err
	}
	checkUser.Password = string(hashedPassword)
	if err := database.DB.DB.Save(&checkUser).Error; err != nil {
		loggerInst.Error(ctx, "Error while updating password.", err)
		errMsg = "Error while updating password."
		return errMsg, err
	}

	return "", nil
}
func (d *dbRepo) ChangeUsername(username, newUsername string, ctx context.Context) (error, string) {
	loggerInst := ctx.Value(logger.LoggerKey).(*logger.LogInstance)
	var errMsg string
	var existsUser models.AuthUser
	result := database.DB.DB.Where("username = ?", username).First(&existsUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			loggerInst.Error(ctx, "The user does not exist, you cant change username..", result.Error)
			errMsg = "The user does not exist, you cant change username.."
		} else {
			// another error is occured
			loggerInst.Error(ctx, "Error checking user registration.", result.Error)
			errMsg = "Error checking user registration."
		}
		return result.Error, errMsg
	}
	var duplicateUser models.AuthUser
	result = database.DB.DB.Where("username = ?", newUsername).First(&duplicateUser)
	if result.Error == nil {
		loggerInst.Error(ctx, "The new username is already taken.")
		errMsg = "The new username is already taken."
		return result.Error, errMsg

	} else if result.Error != gorm.ErrRecordNotFound {
		loggerInst.Error(ctx, "Error checking new username.", result.Error)
		errMsg = "Error checking new username."
		return result.Error, errMsg
	}

	// Update the username
	existsUser.Username = newUsername
	if err := database.DB.DB.Save(&existsUser).Error; err != nil {
		loggerInst.Error(ctx, "Error while updating username.", err)
		errMsg = "Error while updating username."
		return err, errMsg
	}
	return nil, ""

}

func (d *dbRepo) UpdateUser(updatedUser models.AuthUser, username string, ctx context.Context) (models.AuthUser, string, error) {
	loggerInst := ctx.Value(logger.LoggerKey).(logger.LogInstance)
	var errMsg string
	//check the user is in the db or not
	var foundedUser models.AuthUser
	if err := database.DB.DB.Where("username = ?", username).First(&foundedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errMsg = "User not found."
			loggerInst.Error(ctx, "User not found.", err)

		} else {
			errMsg = "Error fetching user from the database."
			loggerInst.Error(ctx, "Error fetching user from the database.", err)
		}
		return models.AuthUser{}, errMsg, err
	}

	var okNewUsername models.AuthUser
	result := database.DB.DB.Where("username = ?", updatedUser.Username).First(&okNewUsername)
	if result.Error == nil {
		errMsg = "The new username is already taken."
		loggerInst.Error(ctx, "The new username is already taken.")
		return models.AuthUser{}, errMsg, result.Error
	} else if result.Error != gorm.ErrRecordNotFound {
		errMsg = "Error checking new username."
		loggerInst.Error(ctx, "Error checking new username.", result.Error)
		return models.AuthUser{}, errMsg, result.Error
	}

	foundedUser.Username = updatedUser.Username
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		errMsg = "Error hashing password"
		loggerInst.Error(ctx, "Error hashing password", err)
		return models.AuthUser{}, errMsg, err
	}
	foundedUser.Password = string(hashedPassword)

	if err := database.DB.DB.Save(&foundedUser).Error; err != nil {
		loggerInst.Error(ctx, "Error updating user in the database", err)
		errMsg = "Error updating user in the database"
		return models.AuthUser{}, errMsg, err
	}
	return foundedUser, errMsg, err
}
func (d *dbRepo) DeleteUser(username string, ctx context.Context) (error, string) {
	loggerInst := ctx.Value(logger.LoggerKey).(*logger.LogInstance)
	var errMsg string
	var foundedUser models.AuthUser
	result := database.DB.DB.Where("username = ?", username).First(&foundedUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errMsg = "The user is not found, so you can't delete it"
			loggerInst.Error(ctx, "The user is not found, so you can't delete it", result.Error)
		} else {
			errMsg = "Error fetching user from the database"
			loggerInst.Error(ctx, "Error fetching user from the database", result.Error)
		}
		return result.Error, errMsg
	}
	if err := database.DB.DB.Delete(&foundedUser).Error; err != nil {
		errMsg = "Error in deleting user"
		loggerInst.Error(ctx, "Error in deleting user", err)
		return err, errMsg
	}
	return nil, ""

}
func (d *dbRepo) Register(user models.AuthUser, ctx context.Context) (error, string) {
	loggerInst := ctx.Value(logger.LoggerKey).(*logger.LogInstance)
	var errMsg string

	//check whether the user is already registered
	var existingUser models.AuthUser
	result := database.DB.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		errMsg = "Error checking user registration"
		loggerInst.Error(ctx, "Error checking user registration", result.Error)
		return result.Error, errMsg
	}

	if result.RowsAffected > 0 {
		errMsg = "User already exists"
		return errors.New("User already exists"), errMsg
	}

	//if the user is not registered, register it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errMsg = "Error in generating hashed password."
		loggerInst.Error(ctx, "Error in generating hashed password.", err)
		return err, errMsg
	}

	user.Password = string(hashedPassword)
	if err := database.DB.DB.Create(&user).Error; err != nil {
		errMsg = "Could not create user"
		loggerInst.Error(ctx, "Could not create user", err)
		return err, errMsg
	}
	return nil, ""
}

func (d *dbRepo) Login(input models.AuthUser, ctx context.Context) (string, string, time.Time, error) {
	loggerInst := ctx.Value(logger.LoggerKey).(*logger.LogInstance)
	var errMsg string
	var user models.AuthUser
	if err := database.DB.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		loggerInst.Error(ctx, "There is no such user.", err)
		errMsg = "There is no such user."
		return "", errMsg, time.Time{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		loggerInst.Error(ctx, "Invalid username or password.", err)
		errMsg = "Invalid username or password."
		return "", errMsg, time.Time{}, err
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		errMsg = "Internal Error."
		loggerInst.Error(ctx, "Internal Error.", err)
		return "", errMsg, time.Time{}, err
	}

	return tokenString, "", expirationTime, nil
}
