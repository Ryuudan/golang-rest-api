package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ryuudan/golang-rest-api/ent/generated"
	"github.com/ryuudan/golang-rest-api/src/constants"
	"github.com/ryuudan/golang-rest-api/src/database"
	"github.com/ryuudan/golang-rest-api/src/internal/services"
	"github.com/ryuudan/golang-rest-api/src/utils"
	"github.com/ryuudan/golang-rest-api/src/utils/render"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	user  services.UserService
	cache *database.RedisCache
}

func NewUserHandler(userService services.UserService, cache *database.RedisCache) *UserHandler {
	return &UserHandler{
		user:  userService,
		cache: cache,
	}
}

func (handler *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Create a validator instance for input validation
	validate := render.Validator()

	var user generated.User
	var validationErrors []render.ValidationErrorDetails

	// Decode the JSON request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		render.Error(w, r, http.StatusUnprocessableEntity, "Invalid JSON: "+err.Error())
		return
	}

	// Struct level validation of the user object
	if err := validate.Struct(user); err != nil {
		render.ValidationError(w, r, err)
		return
	}

	// Check if a user with the same email already exists
	existingUser, _ := handler.user.GetUserByEmail(r.Context(), user.Email)

	if existingUser != nil {
		validationErrors = append(validationErrors, render.ValidationErrorDetails{
			Field:   "email",
			Message: "email already exists, please try another one",
		})
	}

	// If there are validation errors, return a custom validation error response
	if len(validationErrors) > 0 {
		render.CustomValidationError(w, r, validationErrors)
		return
	}

	// Generate a salted and hashed password
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Error(w, r, http.StatusInternalServerError, "failed to generate hashed password")
		return
	}

	// Replace the plaintext password with the hashed password
	user.Password = string(password)

	// Register the user in the system
	newUser, err := handler.user.CreateUser(r.Context(), &user)

	if err != nil {
		render.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = handler.cache.SetCache(
		fmt.Sprintf("users:%d", newUser.ID),
		newUser,
		constants.DEFAULT_CACHE_EXPIRATION,
	)

	if err != nil {
		render.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, newUser)
}

func (handler *UserHandler) GetOneByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.StringToInt(chi.URLParam(r, "id"))

	if err != nil {
		render.Error(w, r, http.StatusBadRequest, constants.INVALID_FORMAT_ID)
		return
	}

	// check cache
	cachedUser, err := handler.cache.GetCache(fmt.Sprintf("users:%d", id))
	if err == nil {
		var user []generated.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			render.JSON(w, http.StatusOK, user)
			return
		}
	}

	user, err := handler.user.GetUserByID(r.Context(), id)

	if err != nil {
		if generated.IsNotFound(err) {
			render.Error(w, r, http.StatusNotFound, "user not found")
			return
		} else {
			render.Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// set cache
	err = handler.cache.SetCache(
		fmt.Sprintf("users:%d", user.ID),
		user,
		constants.DEFAULT_CACHE_EXPIRATION,
	)

	if err != nil {
		render.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, user)
}
