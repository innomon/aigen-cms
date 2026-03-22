# Implementation Plan: Multi-Channel Communication

## Phase 1: Data Modeling & Schema
- [ ] Create `core/descriptors/channel.go` defining `UserChannel`, `AuthLog`, and `ChannelType`.
- [ ] Update `infrastructure/relationdbdao/setup.go` or `framework/init.go` to create `__user_channels` and `__auth_logs` tables.
- [ ] Add `Channels` field to `descriptors.User` (optional/lazy loaded).

## Phase 2: Configuration & Initialization
- [ ] Update `framework/config.go` to include `ChannelsConfig`.
- [ ] Initialize `ChannelService` in `framework/init.go`.
- [ ] Implement `IChannelService` interface in `core/services/interface.go`.

## Phase 3: Channel Service Implementation
- [ ] Create `core/services/channel_service.go`.
- [ ] Implement `RegisterChannel`, `VerifyChannel`, and `GetAuthLogs`.
- [ ] Implement `HandleInbound` (generic entry point for webhooks/pollers).
- [ ] Implement `SendNotification` with multi-channel routing logic.

## Phase 4: Authentication Integration
- [ ] Update `AuthService` to support channel-based authentication.
- [ ] Implement Ed25519 JWT verification for WhatsApp (referencing `whatsadk`).
- [ ] Implement Email JWT verification (referencing `mailadk`).
- [ ] Implement `AuthLog` recording for all channel auth attempts.

## Phase 5: Guest & App-Specific Config
- [ ] Implement guest user creation/handling for specific channels.
- [ ] Ensure `ChannelService` respects per-app configurations.
- [ ] Add `api/channel_api.go` for channel-related endpoints.

## Phase 6: Validation & Testing
- [ ] Unit tests for `ChannelService` logic.
- [ ] Mock tests for WhatsApp/Email auth flows.
- [ ] Verify e-trail logging (IP, UA, etc.).

## Checklist
- [ ] `core/descriptors/channel.go`
- [ ] Database migration for new tables
- [ ] `core/services/channel_service.go`
- [ ] `AuthService` updates
- [ ] `api/channel_api.go`
- [ ] Configuration updates in `framework/config.go`
- [ ] Unit tests
