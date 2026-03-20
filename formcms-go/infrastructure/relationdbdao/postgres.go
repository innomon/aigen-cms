package relationdbdao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/formcms/formcms-go/utils/datamodels"
	_ "github.com/lib/pq"
)

type PostgresDao struct {
	Dao
}

func NewPostgresDao(connectionString string) (*PostgresDao, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &PostgresDao{
		Dao: Dao{
			db:      db,
			builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		},
	}, nil
}

func (d *PostgresDao) GetColumnDefinitions(ctx context.Context, table string) ([]datamodels.Column, error) {
	query := `SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1`
	rows, err := d.db.QueryContext(ctx, query, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []datamodels.Column
	for rows.Next() {
		var name, dataType string
		if err := rows.Scan(&name, &dataType); err != nil {
			return nil, err
		}
		cols = append(cols, datamodels.Column{
			Name: name,
			Type: d.stringToColType(dataType),
		})
	}
	return cols, nil
}

func (d *PostgresDao) CreateTable(ctx context.Context, table string, cols []datamodels.Column) error {
	var parts []string
	var updatedAtField string

	for _, col := range cols {
		if col.Type == datamodels.UpdatedTime {
			updatedAtField = col.Name
		}
		parts = append(parts, fmt.Sprintf(`"%s" %s`, col.Name, d.colTypeToString(col)))
	}

	sqlStr := fmt.Sprintf(`CREATE TABLE "%s" (%s);`, table, strings.Join(parts, ", "))

	if updatedAtField != "" {
		triggerSql := fmt.Sprintf(`
			CREATE OR REPLACE FUNCTION __update_%s_column()
			RETURNS TRIGGER AS $$
			BEGIN
				NEW."%s" = timezone('UTC', now());
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;

			CREATE TRIGGER update_%s_%s
			BEFORE UPDATE ON "%s"
			FOR EACH ROW
			EXECUTE FUNCTION __update_%s_column();
		`, updatedAtField, updatedAtField, table, updatedAtField, table, updatedAtField)
		sqlStr += triggerSql
	}

	_, err := d.db.ExecContext(ctx, sqlStr)
	return err
}

func (d *PostgresDao) AddColumns(ctx context.Context, table string, cols []datamodels.Column) error {
	var parts []string
	for _, col := range cols {
		parts = append(parts, fmt.Sprintf(`ADD COLUMN "%s" %s`, col.Name, d.colTypeToString(col)))
	}
	sqlStr := fmt.Sprintf(`ALTER TABLE "%s" %s;`, table, strings.Join(parts, ", "))
	_, err := d.db.ExecContext(ctx, sqlStr)
	return err
}

func (d *PostgresDao) CreateForeignKey(ctx context.Context, table, col, refTable, refCol string) error {
	sqlStr := fmt.Sprintf(`ALTER TABLE "%s" ADD CONSTRAINT "fk_%s_%s" FOREIGN KEY ("%s") REFERENCES "%s"("%s") ON DELETE CASCADE;`,
		table, table, col, col, refTable, refCol)
	_, err := d.db.ExecContext(ctx, sqlStr)
	return err
}

func (d *PostgresDao) CreateIndex(ctx context.Context, table string, fields []string, isUniq bool) error {
	unique := ""
	if isUniq {
		unique = "UNIQUE"
	}
	quotedFields := make([]string, len(fields))
	for i, f := range fields {
		quotedFields[i] = fmt.Sprintf(`"%s"`, f)
	}
	sqlStr := fmt.Sprintf(`CREATE %s INDEX "idx_%s_%s" ON "%s" (%s);`,
		unique, table, strings.Join(fields, "_"), table, strings.Join(quotedFields, ", "))
	_, err := d.db.ExecContext(ctx, sqlStr)
	return err
}

func (d *PostgresDao) CastDate(field string) string {
	return fmt.Sprintf(`"%s"::date`, field)
}

func (d *PostgresDao) stringToColType(s string) datamodels.ColumnType {
	switch strings.ToLower(s) {
	case "integer", "bigint":
		return datamodels.Int
	case "boolean":
		return datamodels.Boolean
	case "timestamp with time zone", "timestamp without time zone":
		return datamodels.Datetime
	case "text":
		return datamodels.Text
	case "character varying":
		return datamodels.String
	default:
		return datamodels.String
	}
}

func (d *PostgresDao) colTypeToString(col datamodels.Column) string {
	switch col.Type {
	case datamodels.Id:
		return "SERIAL PRIMARY KEY"
	case datamodels.StringPrimaryKey:
		return fmt.Sprintf("VARCHAR(%d) PRIMARY KEY", col.Length)
	case datamodels.Int:
		return "BIGINT"
	case datamodels.Boolean:
		return "BOOLEAN"
	case datamodels.Datetime:
		return "TIMESTAMP WITH TIME ZONE"
	case datamodels.CreatedTime:
		return "TIMESTAMP WITH TIME ZONE DEFAULT (timezone('UTC', now()))"
	case datamodels.UpdatedTime:
		return "TIMESTAMP WITH TIME ZONE DEFAULT (timezone('UTC', now()))"
	case datamodels.Text:
		return "TEXT"
	case datamodels.String:
		length := col.Length
		if length == 0 {
			length = 255
		}
		return fmt.Sprintf("VARCHAR(%d)", length)
	default:
		return "TEXT"
	}
}

// Implement other methods like UpdateOnConflict, BatchUpdateOnConflict, Increase, RenameTable, RenameColumn, DropForeignKey
func (d *PostgresDao) UpdateOnConflict(ctx context.Context, tableName string, data datamodels.Record, keyFields []string) (bool, error) {
	// Implementation for ON CONFLICT DO UPDATE
	return false, nil
}

func (d *PostgresDao) BatchUpdateOnConflict(ctx context.Context, tableName string, records []datamodels.Record, keyFields []string) error {
	return nil
}

func (d *PostgresDao) Increase(ctx context.Context, tableName string, keyConditions datamodels.Record, valueField string, initVal, delta int64) (int64, error) {
	return 0, nil
}

func (d *PostgresDao) RenameTable(ctx context.Context, oldName, newName string) error {
	sqlStr := fmt.Sprintf(`ALTER TABLE "%s" RENAME TO "%s";`, oldName, newName)
	_, err := d.db.ExecContext(ctx, sqlStr)
	return err
}

func (d *PostgresDao) RenameColumn(ctx context.Context, table, oldName, newName string) error {
	sqlStr := fmt.Sprintf(`ALTER TABLE "%s" RENAME COLUMN "%s" TO "%s";`, table, oldName, newName)
	_, err := d.db.ExecContext(ctx, sqlStr)
	return err
}

func (d *PostgresDao) DropForeignKey(ctx context.Context, table, fkName string) error {
	sqlStr := fmt.Sprintf(`ALTER TABLE "%s" DROP CONSTRAINT "%s";`, table, fkName)
	_, err := d.db.ExecContext(ctx, sqlStr)
	return err
}
