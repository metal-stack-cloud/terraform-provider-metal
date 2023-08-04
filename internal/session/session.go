package session

import (
	mclient "github.com/metal-stack-cloud/api/go/client"
)

type Session struct {
	Client       mclient.Client
	Organization string
	Project      string
}
