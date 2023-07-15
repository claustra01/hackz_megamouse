package handler

import (
	"net/http"

	"github.com/claustra01/hackz-megamouse/server/db"
	"github.com/claustra01/hackz-megamouse/server/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SignUP(c echo.Context) error {

	type Body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// parse json
	obj := new(Body)
	if err := c.Bind(obj); err != nil {
		// return 400
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Json Format Error: " + err.Error(),
		})
	}

	// check field
	if util.HasEmptyField(obj, "Username", "Email", "Password") {
		// return 400
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Missing Required Field",
		})
	}

	var user db.User
	if err := db.DB.Where("email = ?", obj.Email).First(&user).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			// create new user
			new := db.User{
				Username: obj.Username,
				Email:    obj.Email,
				Password: obj.Password,
			}
			db.DB.Create(&new)
			return c.JSON(http.StatusCreated, new)

		} else {
			// return 500
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database Error: " + err.Error(),
			})
		}

	} else {
		// return 409
		return c.JSON(http.StatusConflict, echo.Map{
			"message": "Email Conflict",
		})
	}
}
