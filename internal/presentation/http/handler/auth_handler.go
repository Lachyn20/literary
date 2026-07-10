package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/presentation/http/middleware"
	usecaseauth "github.com/hemra-siirow/literary/internal/usecase/auth"
)

type AuthHandler struct {
	loginUseCase          *usecaseauth.LoginUseCase
	refreshUseCase        *usecaseauth.RefreshTokenUseCase
	logoutUseCase         *usecaseauth.LogoutUseCase
	createUserUseCase     *usecaseauth.CreateUserUseCase
	changePasswordUseCase *usecaseauth.ChangePasswordUseCase
}

func NewAuthHandler(login *usecaseauth.LoginUseCase, refresh *usecaseauth.RefreshTokenUseCase, logout *usecaseauth.LogoutUseCase, createUser *usecaseauth.CreateUserUseCase, changePassword *usecaseauth.ChangePasswordUseCase) *AuthHandler {
	return &AuthHandler{loginUseCase: login, refreshUseCase: refresh, logoutUseCase: logout, createUserUseCase: createUser, changePasswordUseCase: changePassword}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/auth/login", h.Login)
	r.Post("/api/auth/refresh", h.Refresh)
	r.Post("/api/auth/logout", h.Logout)
	r.With(middleware.RequireRoles(entity.RoleAdmin, entity.RoleEditor)).Post("/api/auth/change-password", h.ChangePassword)
	r.With(middleware.RequireRoles(entity.RoleAdmin)).Post("/api/admin/users", h.CreateUser)
	r.With(middleware.RequireRoles(entity.RoleAdmin)).Get("/api/admin/users", h.ListUsers)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	user, accessToken, refreshToken, err := h.loginUseCase.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "data": map[string]interface{}{"user": authUserResponse(user), "access_token": accessToken, "refresh_token": refreshToken}})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	accessToken, refreshToken, err := h.refreshUseCase.Execute(r.Context(), req.RefreshToken)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "data": map[string]interface{}{"access_token": accessToken, "refresh_token": refreshToken}})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.logoutUseCase.Execute(r.Context(), req.RefreshToken); err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.changePasswordUseCase.Execute(r.Context(), claims.UserID, req.CurrentPassword, req.NewPassword); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil || claims.Role != entity.RoleAdmin {
		WriteError(w, http.StatusForbidden, "admin role required")
		return
	}
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	user, err := h.createUserUseCase.Execute(r.Context(), claims.Role, req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]interface{}{"status": "ok", "data": authUserResponse(user)})
}

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil || claims.Role != entity.RoleAdmin {
		WriteError(w, http.StatusForbidden, "admin role required")
		return
	}
	users, err := h.createUserUseCase.ListUsers(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		out = append(out, authUserResponse(u))
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "data": out})
}

func authUserResponse(user *entity.User) map[string]interface{} {
	return map[string]interface{}{"id": user.ID.String(), "name": user.Name, "email": user.Email, "role": user.Role, "active": user.Active, "created_at": user.CreatedAt}
}

func (h *AuthHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil || claims.Role != entity.RoleAdmin {
		WriteError(w, http.StatusForbidden, "admin role required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.createUserUseCase.SetActive(r.Context(), id, false); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
}
