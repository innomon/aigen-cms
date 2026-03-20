# Implementation Plan
## Phase 1: Core Framework & Database Abstraction
- Setup Go module and project structure.
- Define data structures (Entity, Attribute, Settings).
- Implement database abstraction layer (`PrimaryDao`, `ReplicaDao`) using `database/sql` + `squirrel`.
- Implement dynamic table creation and migration (`MigrateTable`, `AddColumns`).

## Phase 2: CMS Core APIs
- Implement REST API endpoints for schema management (`/api/schemas`).
- Implement generic CRUD API endpoints for entities (`/api/entities/{name}`).
- Implement junction and collection API endpoints.
- Implement GraphQL parsing and query mapping.

## Phase 3: File Storage & Assets
- Implement local file storage and interface for cloud storage.
- Implement asset tracking in database (`__assets`, `__assetLinks`).
- Support chunked file uploads.
- Implement image resizing (using `bimg` or `disintegration/imaging`).

## Phase 4: Social Engagement & Comments
- Implement `__engagements` tracking (buffered writes).
- Implement comment system and tree-structured replies.
- Implement user portal tracking and authentication (JWT/OAuth).
- Implement Notification and Audit logs.

## Implementation Checklist

### Phase 1: Core Framework & Database Abstraction
- [x] Initialize Go module and project structure
- [x] Define core data structures (`Entity`, `Attribute`, `Settings`)
- [/] Implement database abstraction layer (`PrimaryDao`, `ReplicaDao`) (Partially: Interfaces and PostgresDao base implemented)
- [x] Setup `squirrel` for dynamic SQL building
- [/] Implement dynamic table creation and migration (`MigrateTable`) (Partially: CreateTable, AddColumns in PostgresDao)
- [ ] Implement database connection pooling and health checks

### Phase 2: CMS Core APIs
- [ ] Implement REST API endpoints for schema management (`/api/schemas`)
- [ ] Implement generic CRUD API endpoints for entities (`/api/entities/{name}`)
- [ ] Support basic filtering, sorting, and pagination in REST APIs
- [ ] Implement junction and collection API endpoints (N:N relationships)
- [ ] Implement GraphQL schema generation from dynamic entities
- [ ] Implement GraphQL query execution and N+1 query optimization
- [ ] Support converting GraphQL queries to cached REST endpoints

### Phase 3: File Storage & Assets
- [ ] Define file storage interface (`StorageProvider`)
- [ ] Implement local file storage provider
- [ ] Support extensible cloud storage providers (S3, etc.)
- [ ] Implement asset tracking in database (`__assets`, `__assetLinks`)
- [ ] Implement chunked file upload API
- [ ] Implement image resizing and compression (using `disintegration/imaging`)

### Phase 4: Social Engagement & Comments
- [ ] Implement `__engagements` tracking with buffered writes (Go channels/routines)
- [ ] Implement comment system with tree-structured replies
- [ ] Implement user portal tracking and JWT/OAuth authentication
- [ ] Implement Notification system
- [ ] Implement Audit logs for all schema and data changes

### Phase 5: Page Designer Integration
- [ ] Implement page routing based on dynamic templates
- [ ] Integrate `aymerick/raymond` for Handlebars template compilation
- [ ] Implement GraphQL data fetching within templates
- [ ] Embed React Admin panel and Schema Builder UI into Go binary (`embed`)
- [ ] Implement GrapesJS visual page builder integration

### Phase 6: Final Polish & Deployment
- [ ] Comprehensive unit and integration testing
- [ ] Performance benchmarking (comparison with C# implementation)
- [ ] Documentation and example project setup
- [ ] Dockerize the application
- [ ] CI/CD pipeline setup
