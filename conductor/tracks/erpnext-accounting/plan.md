# Implementation Plan: ERPNext Accounting Transformation

## Phase 1: Research & Mapping (Done)
- [x] Analyze ERPNext accounting DocTypes.
- [x] Analyze `formcms-go` Entity/Schema system.
- [x] Define mapping rules.

## Phase 2: Schema Definition
- [ ] Create `Entity` descriptors in Go for:
    - [ ] Company
    - [ ] Currency
    - [ ] Account
    - [ ] GL Entry
    - [ ] Journal Entry
    - [ ] Journal Entry Account (Child Table)
    - [ ] Fiscal Year
    - [ ] Cost Center
- [ ] Create JSON Schema files for these entities.

## Phase 3: Framework Enhancements
- [ ] Add `Float` and `Decimal` support to `formcms-go`.
    - [ ] Update `descriptors.DataType`.
    - [ ] Update `datamodels.ColumnType`.
    - [ ] Update `sqlite.go` and `postgres.go` to support FLOAT/DECIMAL.
    - [ ] Update `displaymodels.DisplayType` if needed.

## Phase 4: Implementation
- [ ] Implement the `Entity` registration in `formcms-go`.
- [ ] Verify schema creation in the database.

## Phase 5: Validation
- [ ] Create test data for each entity.
- [ ] Verify that basic CRUD operations work for the translated accounting entities.
