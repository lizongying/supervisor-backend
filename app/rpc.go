package app

import (
	"github.com/kolo/xmlrpc"
)

type SupervisorRpc struct {
	Url    string
	Client *xmlrpc.Client
}

func Rpc(url string) *SupervisorRpc {
	supervisorRpc := &SupervisorRpc{Url: url}
	client, err := supervisorRpc.GetClient()
	if err == nil {
		supervisorRpc.Client = client
	}
	return supervisorRpc
}
