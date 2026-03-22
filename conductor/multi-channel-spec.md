# Specification: Multi-Channel Communication & Authentication

## 1. Overview
The Multi-Channel Communication system enables `aigen-cms` to interact with users across various platforms (WhatsApp, Email, Signal, Telegram, X.com, Bluesky). It provides a unified way to link channel identities to user profiles, authenticate via these channels, and maintain a secure audit trail (e-trail) for all interactions.

## 2. Core Components

### 2.1. Channel Types
Supported channels include:
- **WhatsApp**: Integrated via `whatsadk` (Ed25519 JWT auth).
- **Email**: Integrated via `mailadk` (IMAP polling, JWT verification).
- **Signal/Telegram**: Messaging-based interaction and auth.
- **X.com (Twitter) / Bluesky**: Social platform interaction for support/tracking.

### 2.2. User Channel Mapping (`UserChannel`)
Each user can have multiple optional channels.
- `UserId`: Link to the core `User`.
- `ChannelType`: Type of channel (e.g., `whatsapp`, `email`).
- `Identifier`: Channel-specific ID (phone number, email address, handle).
- `IsAuthenticated`: Boolean flag for verification status.
- `Metadata`: JSON field storing tokens, public keys, or session data.

### 2.3. Authentication & Non-Repudiation
- **Multi-Channel Auth**: Users can log in or verify their identity using a channel (e.g., clicking a WhatsApp deep link).
- **Auth Log (E-trail)**: Every authentication attempt is logged with:
    - User ID (if known)
    - Channel Type
    - Action (login, verify, etc.)
    - IP Address
    - User Agent
    - Success Status
    - Payload/Nonce hash for non-repudiation.

### 2.4. Guest Access
- Channels can be configured to allow "guest" users (e.g., for initial support queries).
- Guest users can be automatically promoted to registered users upon successful channel verification.

### 2.5. Configuration
Downstream apps can configure channels via `config.yaml`:
```yaml
channels:
  whatsapp:
    enabled: true
    public_key: "..." # For verifying whatsadk JWTs
  email:
    enabled: true
    verification_required: true
  guest_access:
    allowed_channels: ["whatsapp", "email"]
    default_role: "guest"
```

## 3. Architecture
- **Service**: `ChannelService` handles registration, verification, and routing of inbound/outbound messages.
- **Security**: Uses `Ed25519` for compact, secure tokens (especially for WhatsApp).
- **Storage**: New tables `__user_channels` and `__auth_logs`.
- **Extensibility**: Each channel can have its own provider implementation.
