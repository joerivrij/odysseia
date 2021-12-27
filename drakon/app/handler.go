package app

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
)

func (d *DrakonConfig) CreateRoles() (bool, error) {
	glg.Debug("creating elastic roles based on labels")

	var created bool
	for _, index := range d.Indexes {
		for _, role := range d.Roles {
			glg.Debugf("creating a role for index %s with role %s", index, role)

			var privileges []string
			if role == "seeder" {
				privileges = append(privileges, "delete_index")
				privileges = append(privileges, "create_index")
				privileges = append(privileges, "create")
			} else {
				privileges = append(privileges, "read")
			}

			names := []string{index}

			indices := []models.Index{
				{
					Names:      names,
					Privileges: privileges,
					Query:      "",
				},
			}

			application := []models.Application{}

			putRole := models.CreateRoleRequest{
				Cluster:      []string{"all"},
				Indices:      indices,
				Applications: application,
				RunAs:        nil,
				Metadata:     models.Metadata{Version: 1},
			}

			roleName := fmt.Sprintf("%s_%s", index, role)
			roleCreated, err := elastic.CreateRole(&d.ElasticClient, roleName, putRole)
			if err != nil {
				glg.Error(err)
				return false, err
			}

			created = roleCreated
		}
	}

	return created, nil
}
