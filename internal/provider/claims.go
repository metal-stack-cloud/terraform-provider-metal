package provider

import "github.com/golang-jwt/jwt/v4"

type (
	// MethodPermissions is a map from project or organization ->[]methods, e.g.: "/api.v1.ClusterService/List".
	MethodPermissions map[string][]string
	// TokenRoles maps the role to subject
	// subject can be either * (Wildcard) or the concrete Organization or Project
	// role can be one of admin, owner, editor, viewer.
	TokenRoles map[string]string
	Claims     struct {
		jwt.RegisteredClaims

		Roles       TokenRoles        `json:"roles,omitempty"`
		Permissions MethodPermissions `json:"permissions,omitempty"`
		Type        string            `json:"type"`
	}
)
