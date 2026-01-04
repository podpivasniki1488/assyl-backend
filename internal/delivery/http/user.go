package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/protopb"
)

// TODO: here we need to add user handlers and first handler must be giving users cinema reservations

func (h *httpDelivery) registerUserHandlers(v1 *echo.Group) {
	user := v1.Group("/user")

	user.Use(h.registerJWTMiddleware())
	user.DELETE("", h.deleteUser, h.getJWTData())

}

// deleteUser godoc
//
//	@Summary		Delete user (admin/god only)
//	@Description	Удаляет пользователя по username. Доступно только ролям ADMIN и GOD. Нельзя удалить самого себя.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	deleteUserRequest	true	"Delete user request"
//	@Success		204		"Пользователь удалён"
//	@Failure		400		{object}	DefaultResponse[error]	"Невалидный запрос"
//	@Failure		401		{object}	DefaultResponse[error]	"Не авторизован / некорректный токен"
//	@Failure		403		{object}	DefaultResponse[error]	"Недостаточно прав / попытка удалить себя"
//	@Failure		500		{object}	DefaultResponse[error]	"Внутренняя ошибка сервера"
//	@Router			/user [delete]
func (h *httpDelivery) deleteUser(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "httpDelivery.deleteUser")
	defer span.End()

	var req deleteUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	role, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("role not found in context"))
	}

	fmt.Println(protopb.Role_ADMIN.String())

	if role != protopb.Role_ADMIN.String() && role != protopb.Role_GOD.String() {
		return c.JSON(http.StatusForbidden, ErrorResponse("role not allowed"))
	}

	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("username not found in context"))
	}

	if username == req.Username {
		return c.JSON(http.StatusForbidden, ErrorResponse("cannot delete yourself bro)"))
	}

	if err := h.service.UserManagement.DeleteUserByUsername(ctx, req.Username); err != nil {
		return h.handleErrResponse(c, err)
	}

	return nil
}

type deleteUserRequest struct {
	Username string `json:"username" validate:"required"`
}
