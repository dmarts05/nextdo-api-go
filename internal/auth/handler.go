package auth

import (
	"errors"
	"net/http"

	"github.com/dmarts05/nextdo-api-go/internal/customer"
	"github.com/dmarts05/nextdo-api-go/internal/shared/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	CustomerRepo customer.Repository
}

func (h Handler) Login(c echo.Context) error {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if the customer exists
	customer, err := h.CustomerRepo.GetByEmail(c.Request().Context(), body.Email)
	switch {
	case errors.Is(err, repository.ErrNotFound{}):
		return echo.NewHTTPError(http.StatusUnauthorized, "customer does not exist")
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, "an error occurred while processing the request, please try again later")
	}

	// Compare password
	ok := IsPasswordValid(customer.Password, body.Password)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
	}

	// Generate token pair
	tokenPair, err := generateCustomerTokenPair(customer.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "an error occurred while processing the request, please try again later")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

func (h Handler) Register(c echo.Context) error {
	var body struct {
		FirstName string `json:"first_name" validate:"required,min=2,max=255"`
		LastName  string `json:"last_name" validate:"required,min=2,max=255"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if the customer exists
	_, err := h.CustomerRepo.GetByEmail(c.Request().Context(), body.Email)
	switch {
	case errors.Is(err, repository.ErrNotFound{}):
		// Hash the password
		hash, err := HashPassword(body.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "an error occurred while processing the request, please try again later")
		}

		// Create the customer
		customer := customer.Customer{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Email:     body.Email,
			Password:  hash,
		}
		if err := h.CustomerRepo.Create(c.Request().Context(), customer); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "an error occurred while processing the request, please try again later")
		}

		return c.NoContent(http.StatusNoContent)
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, "an error occurred while processing the request, please try again later")
	// No error means the customer already exists
	default:
		return echo.NewHTTPError(http.StatusConflict, "customer already exists")
	}
}

func (h Handler) Refresh(c echo.Context) error {
	var body struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get customer ID from refresh token
	customerID, err := getCustomerIDFromToken(body.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}

	// Check if the customer exists
	_, err = h.CustomerRepo.GetByID(c.Request().Context(), customerID)
	switch {
	case errors.Is(err, repository.ErrNotFound{}):
		return echo.NewHTTPError(http.StatusUnauthorized, "customer does not exist")
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, "an error occurred while processing the request, please try again later")
	}

	// Generate token pair
	tokenPair, err := generateCustomerTokenPair(customerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "an error occurred while processing the request, please try again later")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}
