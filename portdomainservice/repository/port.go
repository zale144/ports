package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zale144/ports/portdomainservice/internal/config"
	"github.com/zale144/ports/portdomainservice/internal/model"
	db "github.com/zale144/ports/portdomainservice/pkg/database"
)

type Port struct {
	cfg *config.Config
	db  *db.DB
}

func NewPort(cfg *config.Config, db *db.DB) Port {
	return Port{
		cfg: cfg,
		db:  db,
	}
}

// SavePorts batch inserts a number of ports to the database using a single query
func (pr Port) SavePorts(ctx context.Context, ports []model.Port) error {

	query := "INSERT INTO port(id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code) VALUES "

	// concatenate for each batch's insert query string
	var values []string
	for _, p := range ports {
		values = append(values, fmt.Sprintf(" ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') ",
			p.ID, p.Name, p.City, p.Country, toJson(p.Alias), toJson(p.Regions), toJson(p.Coordinates), p.Province, p.Timezone, toJson(p.Unlocs), p.Code))
	}
	query += strings.Join(values, ",") + ` ON CONFLICT (id) DO UPDATE 
			SET name = EXCLUDED.name, city = EXCLUDED.city, country = EXCLUDED.country, alias = EXCLUDED.alias, regions = EXCLUDED.regions, coordinates = EXCLUDED.coordinates, 
			province = EXCLUDED.province, timezone = EXCLUDED.timezone, unlocs = EXCLUDED.unlocs, code = EXCLUDED.code;`

	stmt, err := pr.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}
	return nil
}
// GetPorts returns a list of ports
func (pr Port) GetPorts(ctx context.Context) ([]model.Port, error) {
	var ports []model.Port
	if err := pr.db.SelectContext(ctx, &ports,"SELECT id, name, city, country, province, timezone, code FROM port"); err != nil {
		return nil, err
	}
	return ports, nil
}

func toJson(i interface{}) string {
	s, _ := json.Marshal(&i)
	return string(s)
}
