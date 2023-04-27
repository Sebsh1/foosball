package controllers

import (
	"foosball/internal/authentication"
	"foosball/internal/invite"
	"foosball/internal/match"
	"foosball/internal/organization"
	"foosball/internal/rating"
	"foosball/internal/user"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	logger              *logrus.Entry
	authService         authentication.Service
	userService         user.Service
	organizationService organization.Service
	inviteService       invite.Service
	matchService        match.Service
	ratingService       rating.Service
}

func Register(
	e *echo.Group,
	logger *logrus.Entry,
	authService authentication.Service,
	userService user.Service,
	organizationService organization.Service,
	inviteService invite.Service,
	matchService match.Service,
	ratingService rating.Service,
) {
	h := &Handlers{
		logger:              logger,
		authService:         authService,
		userService:         userService,
		organizationService: organizationService,
		inviteService:       inviteService,
		matchService:        matchService,
		ratingService:       ratingService,
	}

	// Authentication
	e.POST("/login", h.Login)
	e.POST("/signup", h.Signup)

	// Users
	e.DELETE("/user/:userId", h.DeleteUser)
	e.GET("/user/:userId/invites", h.GetUserInvites)
	e.POST("/user/:userId/invite/:inviteId/accept", h.AcceptInvite)
	e.POST("/user/:userId/invite/:inviteId/decline", h.DeclineInvite)

	// Organizations
	e.GET("/organization/:orgId/users", h.GetUsersInOrganization)
	e.DELETE("/organization/:orgId", h.DeleteOrganization)
	e.POST("/organization", h.CreateOrganization)
	e.POST("/organization/:orgId", h.UpdateOrganization)
	e.POST("/organization/:orgId/invite", h.InviteUserToOrganization)
	e.POST("/organization/:orgId/user/:userId/remove", h.RemoveUserFromOrganization)
	e.POST("/organization/:orgId/user/:userId/admin", h.SetUserAsAdmin)
	e.POST("/organization/:orgId/user/:userId/admin/remove", h.RemoveAdminFromUser)

	// Matches
	e.POST("/match", h.PostMatch)
}
