package nomad

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigAccess(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/access",
		Fields: map[string]*framework.FieldSchema{
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Nomad server address",
			},

			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Token for API calls",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigAccessRead,
			logical.UpdateOperation: b.pathConfigAccessWrite,
		},
	}
}

func (b *backend) readConfigAccess(storage logical.Storage) (*accessConfig, error) {
	entry, err := storage.Get("config/access")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf(
			"Access credentials for the backend itself haven't been configured. Please configure them at the '/config/access' endpoint")
	}

	conf := &accessConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, fmt.Errorf("error reading nomad access configuration: %s", err)
	}

	return conf, nil
}

func (b *backend) pathConfigAccessRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.readConfigAccess(req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, fmt.Errorf("no user or internal error reported but nomad access configuration not found")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": conf.Address,
		},
	}, nil
}

func (b *backend) pathConfigAccessWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/access", accessConfig{
		Address: data.Get("address").(string),
		Token:   data.Get("token").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type accessConfig struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}
