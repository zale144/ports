package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
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

	valueStrings := make([]string, 0, len(ports))
	valueArgs := make([]interface{}, 0, len(ports)*11)
	for i, p := range ports {
		var sb strings.Builder
		for j := i*11 + 1; j <= i*11+11; j++ {
			sb.WriteString(fmt.Sprintf("$%d,", j))
		}
		valStr := "(" + strings.TrimSuffix(sb.String(), ",") + ")"
		valueStrings = append(valueStrings, valStr)
		valueArgs = append(valueArgs, p.ID)
		valueArgs = append(valueArgs, p.Name)
		valueArgs = append(valueArgs, p.City)
		valueArgs = append(valueArgs, p.Country)
		valueArgs = append(valueArgs, toJson(p.Alias))
		valueArgs = append(valueArgs, toJson(p.Regions))
		valueArgs = append(valueArgs, toJson(p.Coordinates))
		valueArgs = append(valueArgs, p.Province)
		valueArgs = append(valueArgs, p.Timezone)
		valueArgs = append(valueArgs, toJson(p.Unlocs))
		valueArgs = append(valueArgs, p.Code)
	}
	query := fmt.Sprintf(`INSERT INTO port(id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code) VALUES %s ON CONFLICT (id) DO UPDATE 
			SET name = EXCLUDED.name, city = EXCLUDED.city, country = EXCLUDED.country, alias = EXCLUDED.alias, regions = EXCLUDED.regions, coordinates = EXCLUDED.coordinates, 
			province = EXCLUDED.province, timezone = EXCLUDED.timezone, unlocs = EXCLUDED.unlocs, code = EXCLUDED.code;`, strings.Join(valueStrings, ","))

	stmt, err := pr.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	if _, err := stmt.ExecContext(ctx, valueArgs...); err != nil {
		return err
	}
	return nil
}

// GetPorts returns a list of ports
func (pr Port) GetPorts(ctx context.Context) ([]model.Port, error) {
	var ports []model.Port
	if err := pr.db.SelectContext(ctx, &ports, "SELECT id, name, city, country, province, timezone, code FROM port"); err != nil {
		return nil, err
	}
	return ports, nil
}

func toJson(i interface{}) string {
	if reflect.ValueOf(i).IsNil() {
		return "[]"
	}
	s, _ := json.Marshal(&i)
	return string(s)
}
