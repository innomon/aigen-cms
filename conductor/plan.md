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

## Phase 5: Page Designer Integration
- Implement page routing.
- Implement Handlebars compilation and GraphQL data fetching.
- Embed React Admin panel and Schema Builder UI.
