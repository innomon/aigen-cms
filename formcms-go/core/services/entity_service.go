package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/formcms/formcms-go/infrastructure/relationdbdao"
	"github.com/formcms/formcms-go/utils/datamodels"
)

type EntityService struct {
	schemaService ISchemaService
	dao           relationdbdao.IPrimaryDao
}

func NewEntityService(schemaService ISchemaService, dao relationdbdao.IPrimaryDao) *EntityService {
	return &EntityService{
		schemaService: schemaService,
		dao:           dao,
	}
}

func (s *EntityService) List(ctx context.Context, name string, pagination datamodels.Pagination, filters []datamodels.Filter, sorts []datamodels.Sort) ([]datamodels.Record, int64, error) {
	entity, err := s.schemaService.LoadEntity(ctx, name)
	if err != nil {
		return nil, 0, err
	}

	sb := entity.SelectQuery(s.dao.GetBuilder())
	sb = s.applyFilters(sb, filters)
	sb = s.applySorts(sb, sorts)
	sb = s.applyPagination(sb, pagination, uint64(entity.DefaultPageSize))

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.dao.GetDb().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results, err := s.scanRows(rows)
	if err != nil {
		return nil, 0, err
	}

	// TODO: Get total count
	return results, int64(len(results)), nil
}

func (s *EntityService) Single(ctx context.Context, name string, id interface{}) (datamodels.Record, error) {
	entity, err := s.schemaService.LoadEntity(ctx, name)
	if err != nil {
		return nil, err
	}

	query, args, err := entity.SelectQuery(s.dao.GetBuilder()).Where(squirrel.Eq{entity.PrimaryKey: id}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.dao.GetDb().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := s.scanRows(rows)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("record not found")
	}

	return results[0], nil
}

func (s *EntityService) Insert(ctx context.Context, name string, data datamodels.Record) (datamodels.Record, error) {
	entity, err := s.schemaService.LoadEntity(ctx, name)
	if err != nil {
		return nil, err
	}

	var columns []string
	var values []interface{}
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
	}

	query, args, err := s.dao.GetBuilder().Insert(entity.TableName).Columns(columns...).Values(values...).ToSql()
	if err != nil {
		return nil, err
	}

	if s.dao.GetBuilder().PlaceholderFormat(squirrel.Dollar) == squirrel.Dollar {
		var newId interface{}
		err = s.dao.GetDb().QueryRowContext(ctx, query+" RETURNING "+entity.PrimaryKey, args...).Scan(&newId)
		if err != nil {
			return nil, err
		}
		return s.Single(ctx, name, newId)
	} else {
		res, err := s.dao.GetDb().ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		newId, _ := res.LastInsertId()
		return s.Single(ctx, name, newId)
	}
}

func (s *EntityService) Update(ctx context.Context, name string, data datamodels.Record) (datamodels.Record, error) {
	entity, err := s.schemaService.LoadEntity(ctx, name)
	if err != nil {
		return nil, err
	}

	id := data[entity.PrimaryKey]
	sb := s.dao.GetBuilder().Update(entity.TableName)
	for k, v := range data {
		if k == entity.PrimaryKey {
			continue
		}
		sb = sb.Set(k, v)
	}

	query, args, err := sb.Where(squirrel.Eq{entity.PrimaryKey: id}).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = s.dao.GetDb().ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.Single(ctx, name, id)
}

func (s *EntityService) Delete(ctx context.Context, name string, id interface{}) error {
	entity, err := s.schemaService.LoadEntity(ctx, name)
	if err != nil {
		return err
	}

	query, args, err := s.dao.GetBuilder().Delete(entity.TableName).Where(squirrel.Eq{entity.PrimaryKey: id}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.dao.GetDb().ExecContext(ctx, query, args...)
	return err
}

func (s *EntityService) CollectionList(ctx context.Context, name, id, attrName string, pagination datamodels.Pagination, filters []datamodels.Filter, sorts []datamodels.Sort) ([]datamodels.Record, int64, error) {
	le, err := s.schemaService.LoadLoadedEntity(ctx, name)
	if err != nil {
		return nil, 0, err
	}

	var collectionAttr *relationdbdao.LoadedAttribute
	for i := range le.Attributes {
		if le.Attributes[i].Field == attrName {
			collectionAttr = &le.Attributes[i]
			break
		}
	}
	if collectionAttr == nil || collectionAttr.Collection == nil {
		return nil, 0, fmt.Errorf("collection attribute %s not found in entity %s", attrName, name)
	}

	collection := collectionAttr.Collection
	targetEntity := collection.TargetEntity

	sb := s.dao.GetBuilder().Select("*").From(targetEntity.TableName).
		Where(squirrel.Eq{collection.LinkAttribute.Field: id})

	sb = s.applyFilters(sb, filters)
	sb = s.applySorts(sb, sorts)
	sb = s.applyPagination(sb, pagination, uint64(targetEntity.DefaultPageSize))

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.dao.GetDb().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results, err := s.scanRows(rows)
	return results, int64(len(results)), err
}

func (s *EntityService) CollectionInsert(ctx context.Context, name, id, attrName string, data datamodels.Record) (datamodels.Record, error) {
	le, err := s.schemaService.LoadLoadedEntity(ctx, name)
	if err != nil {
		return nil, err
	}

	var collectionAttr *relationdbdao.LoadedAttribute
	for i := range le.Attributes {
		if le.Attributes[i].Field == attrName {
			collectionAttr = &le.Attributes[i]
			break
		}
	}
	if collectionAttr == nil || collectionAttr.Collection == nil {
		return nil, fmt.Errorf("collection attribute %s not found", attrName)
	}

	collection := collectionAttr.Collection
	data[collection.LinkAttribute.Field] = id
	return s.Insert(ctx, collection.TargetEntity.Name, data)
}

func (s *EntityService) JunctionList(ctx context.Context, name, id, attrName string, exclude bool, pagination datamodels.Pagination, filters []datamodels.Filter, sorts []datamodels.Sort) ([]datamodels.Record, int64, error) {
	le, err := s.schemaService.LoadLoadedEntity(ctx, name)
	if err != nil {
		return nil, 0, err
	}

	var junctionAttr *relationdbdao.LoadedAttribute
	for i := range le.Attributes {
		if le.Attributes[i].Field == attrName {
			junctionAttr = &le.Attributes[i]
			break
		}
	}
	if junctionAttr == nil || junctionAttr.Junction == nil {
		return nil, 0, fmt.Errorf("junction attribute %s not found", attrName)
	}

	junction := junctionAttr.Junction
	targetEntity := junction.TargetEntity

	// subquery to find target ids
	subQuery, subArgs, _ := s.dao.GetBuilder().Select(junction.TargetAttribute.Field).
		From(junction.JunctionEntity.TableName).
		Where(squirrel.Eq{junction.SourceAttribute.Field: id}).ToSql()

	sb := s.dao.GetBuilder().Select("*").From(targetEntity.TableName)
	if exclude {
		sb = sb.Where(fmt.Sprintf("%s NOT IN (%s)", targetEntity.PrimaryKey, subQuery), subArgs...)
	} else {
		sb = sb.Where(fmt.Sprintf("%s IN (%s)", targetEntity.PrimaryKey, subQuery), subArgs...)
	}

	sb = s.applyFilters(sb, filters)
	sb = s.applySorts(sb, sorts)
	sb = s.applyPagination(sb, pagination, uint64(targetEntity.DefaultPageSize))

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.dao.GetDb().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results, err := s.scanRows(rows)
	return results, int64(len(results)), err
}

func (s *EntityService) JunctionSave(ctx context.Context, name, id, attrName string, targetIds []interface{}) error {
	le, err := s.schemaService.LoadLoadedEntity(ctx, name)
	if err != nil {
		return err
	}

	var junctionAttr *relationdbdao.LoadedAttribute
	for i := range le.Attributes {
		if le.Attributes[i].Field == attrName {
			junctionAttr = &le.Attributes[i]
			break
		}
	}
	if junctionAttr == nil || junctionAttr.Junction == nil {
		return fmt.Errorf("junction attribute %s not found", attrName)
	}

	junction := junctionAttr.Junction
	tx, err := s.dao.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Delete existing
	delQuery, delArgs, _ := s.dao.GetBuilder().Delete(junction.JunctionEntity.TableName).
		Where(squirrel.Eq{junction.SourceAttribute.Field: id}).ToSql()
	if _, err := tx.ExecContext(ctx, delQuery, delArgs...); err != nil {
		return err
	}

	// 2. Insert new
	for _, tid := range targetIds {
		insQuery, insArgs, _ := s.dao.GetBuilder().Insert(junction.JunctionEntity.TableName).
			Columns(junction.SourceAttribute.Field, junction.TargetAttribute.Field).
			Values(id, tid).ToSql()
		if _, err := tx.ExecContext(ctx, insQuery, insArgs...); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *EntityService) JunctionDelete(ctx context.Context, name, id, attrName string, targetIds []interface{}) error {
	le, err := s.schemaService.LoadLoadedEntity(ctx, name)
	if err != nil {
		return err
	}

	var junctionAttr *relationdbdao.LoadedAttribute
	for i := range le.Attributes {
		if le.Attributes[i].Field == attrName {
			junctionAttr = &le.Attributes[i]
			break
		}
	}
	if junctionAttr == nil || junctionAttr.Junction == nil {
		return fmt.Errorf("junction attribute %s not found", attrName)
	}

	junction := junctionAttr.Junction
	query, args, err := s.dao.GetBuilder().Delete(junction.JunctionEntity.TableName).
		Where(squirrel.Eq{junction.SourceAttribute.Field: id, junction.TargetAttribute.Field: targetIds}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.dao.GetDb().ExecContext(ctx, query, args...)
	return err
}

func (s *EntityService) applyFilters(sb squirrel.SelectBuilder, filters []datamodels.Filter) squirrel.SelectBuilder {
	for _, f := range filters {
		for _, c := range f.Constraints {
			if c.Match == "equals" && len(c.Values) > 0 {
				sb = sb.Where(squirrel.Eq{f.FieldName: *c.Values[0]})
			}
			// Add more match types as needed
		}
	}
	return sb
}

func (s *EntityService) applySorts(sb squirrel.SelectBuilder, sorts []datamodels.Sort) squirrel.SelectBuilder {
	for _, sort := range sorts {
		order := "ASC"
		if sort.Order == datamodels.SortOrderDesc {
			order = "DESC"
		}
		sb = sb.OrderBy(fmt.Sprintf("%s %s", sort.Field, order))
	}
	return sb
}

func (s *EntityService) applyPagination(sb squirrel.SelectBuilder, pagination datamodels.Pagination, defaultLimit uint64) squirrel.SelectBuilder {
	limit := defaultLimit
	if pagination.Limit != nil {
		if l, err := strconv.ParseUint(*pagination.Limit, 10, 64); err == nil {
			limit = l
		}
	}
	sb = sb.Limit(limit)

	if pagination.Offset != nil {
		if o, err := strconv.ParseUint(*pagination.Offset, 10, 64); err == nil {
			sb = sb.Offset(o)
		}
	}
	return sb
}

func (s *EntityService) scanRows(rows *sql.Rows) ([]datamodels.Record, error) {
	var results []datamodels.Record
	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		record := make(datamodels.Record)
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				record[col] = string(b)
			} else {
				record[col] = val
			}
		}
		results = append(results, record)
	}
	return results, nil
}
