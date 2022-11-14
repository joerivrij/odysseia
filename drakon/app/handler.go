package app

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia-greek/plato/elastic"
)

type DrakonHandler struct {
	Config *configs.DrakonConfig
}

const (
	seederRole string = "seeder"
	hybridRole string = "hybrid"
	apiRole    string = "api"
)

func (d *DrakonHandler) CreateRoles() (bool, error) {
	glg.Debug("creating elastic roles based on labels")

	var created bool
	for _, index := range d.Config.Indexes {
		for _, role := range d.Config.Roles {
			glg.Debugf("creating a role for index %s with role %s", index, role)

			var privileges []string
			switch role {
			case seederRole:
				privileges = append(privileges, "delete_index")
				privileges = append(privileges, "create_index")
				privileges = append(privileges, "create")
			case hybridRole:
				privileges = append(privileges, "create")
				privileges = append(privileges, "read")
			case apiRole:
				privileges = append(privileges, "read")
			}

			names := []string{index}

			indices := []elastic.Indices{
				{
					Names:      names,
					Privileges: privileges,
					Query:      "",
				},
			}

			putRole := elastic.CreateRoleRequest{
				Cluster:      []string{"all"},
				Indices:      indices,
				Applications: []elastic.Application{},
				RunAs:        nil,
				Metadata:     elastic.Metadata{Version: 1},
			}

			roleName := fmt.Sprintf("%s_%s", index, role)

			glg.Info(roleName)
			roleCreated, err := d.Config.Elastic.Access().CreateRole(roleName, putRole)
			if err != nil {
				glg.Error(err)
				return false, err
			}

			created = roleCreated
		}
	}

	return created, nil
}
