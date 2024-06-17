package handlers

import (
    "encoding/json"
    "github.com/boldyrevn/mod-example/internal/adapters/db"
    "github.com/boldyrevn/mod-example/internal/model"
    "github.com/boldyrevn/mod-example/internal/ports/httpserver/dto"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "io"
    "log/slog"
    "net/http"
)

type UserHandlers struct {
    db     db.DB
    logger *slog.Logger
}

func NewUserHandlers(db db.DB, logger *slog.Logger) *UserHandlers {
    return &UserHandlers{
        db:     db,
        logger: logger,
    }
}

func (h *UserHandlers) CreateUser(c *gin.Context) {
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        h.logger.Error("failed to read request body", slog.String("error", err.Error()))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    var req dto.CreateUserRequest
    if err := json.Unmarshal(body, &req); err != nil {
        h.logger.Error("failed to unmarshal request body", slog.String("error", err.Error()))
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    newUser := model.User{
        ID:        uuid.New().String(),
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Age:       req.Age,
        Profile:   req.Profile,
    }
    if err := h.db.CreateUser(c.Request.Context(), newUser); err != nil {
        h.logger.Error("failed to create user", slog.String("error", err.Error()))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    h.logger.Info("user successfully created",
        slog.Any("user", newUser),
        slog.String("requestID", c.GetString("requestID")))
    c.JSON(http.StatusOK, newUser)
}

func (h *UserHandlers) UpdateUser(c *gin.Context) {
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        h.logger.Error("failed to read request body", slog.String("error", err.Error()))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    var req dto.UpdateUserRequest
    if err := json.Unmarshal(body, &req); err != nil {
        h.logger.Error("failed to unmarshal request body", slog.String("error", err.Error()))
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    updateUser := model.User{
        ID:        req.ID,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Age:       req.Age,
        Profile:   req.Profile,
    }
    if err := h.db.UpdateUser(c.Request.Context(), updateUser); err != nil {
        h.logger.Error("failed to update user", slog.String("error", err.Error()))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    h.logger.Info("user successfully updated",
        slog.Any("user", updateUser),
        slog.String("requestID", c.GetString("requestID")))
    c.Status(http.StatusOK)
}

func (h *UserHandlers) DeleteUser(c *gin.Context) {
    userID := c.Request.URL.Query().Get("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, dto.RequestMessage{Message: "userID param is not set"})
        return
    }

    if err := h.db.DeleteUser(c.Request.Context(), userID); err != nil {
        h.logger.Error("failed to delete user", slog.String("error", err.Error()))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    h.logger.Info("user successfully deleted",
        slog.String("userID", userID),
        slog.String("requestID", c.GetString("requestID")))
    c.Status(http.StatusOK)
}

func (h *UserHandlers) GetUser(c *gin.Context) {
    userID := c.Request.URL.Query().Get("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, dto.RequestMessage{Message: "userID param is not set"})
        return
    }

    user, err := h.db.GetUser(c.Request.Context(), userID)
    if err != nil {
        h.logger.Error("failed to get user", slog.String("error", err.Error()))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    h.logger.Info("user successfully selected",
        slog.Any("user", user),
        slog.String("requestID", c.GetString("requestID")))
    c.JSON(http.StatusOK, user)
}
