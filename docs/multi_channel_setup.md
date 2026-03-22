# Multi-Channel Communication Setup Guide

This guide provides instructions for setting up and configuring the multi-channel communication system in your `aigen-cms` powered application.

## Prerequisites

1.  **Aigen CMS Instance**: A running instance of `aigen-cms` (version with multi-channel support).
2.  **External ADKs (for WhatsApp/Email)**:
    -   **WhatsApp**: [whatsadk](https://github.com/innomon/whatsadk) gateway for Ed25519 JWT authentication.
    -   **Email**: [mailadk](https://github.com/innomon/mailadk) for IMAP polling and email-based identity verification.
3.  **Public Keys**: The public keys from your `whatsadk` and `mailadk` instances for JWT verification.

## Step-by-Step Configuration

### 1. Update `config.yaml`

Enable and configure the channels you wish to support in your `config.yaml` file:

```yaml
channels:
  whatsapp:
    enabled: true
    gateway_url: "https://your-whatsadk-gateway.com"
    public_key: "YOUR_ED25519_PUBLIC_KEY_BASE64URL"
  email:
    enabled: true
    imap_server: "imap.gmail.com"
    verification_required: true
  guest_access:
    allowed_channels: ["whatsapp", "email"]
    default_role: "guest"
```

### 2. User Profile Setup

Users can link their channel identities through the API. For example, to link a WhatsApp number:

**Endpoint**: `POST /api/channels`
**Body**:
```json
{
  "channelType": "whatsapp",
  "identifier": "+1234567890",
  "metadata": {
    "preferred_name": "John Doe"
  }
}
```

### 3. Channel Verification

Once a channel is registered, it must be verified. For WhatsApp, this typically involves sending an `AUTH` message to the `whatsadk` gateway to receive a JWT.

**Endpoint**: `POST /api/channels/verify`
**Body**:
```json
{
  "channelType": "whatsapp",
  "token": "YOUR_RECEIVED_JWT_TOKEN"
}
```

### 4. Authenticating via Channel

Users can log in directly using a verified channel token:

**Endpoint**: `POST /api/auth/login/channel`
**Body**:
```json
{
  "channelType": "whatsapp",
  "identifier": "+1234567890",
  "token": "YOUR_RECEIVED_JWT_TOKEN"
}
```

### 5. Monitoring E-trail Logs

You can view the authentication history and secure audit trail for non-repudiation:

**Endpoint**: `GET /api/channels/logs?limit=10&offset=0`

## Implementation for Downstream Apps

If you are building an app on top of the `aigen-cms` framework:

1.  **Inject the Service**: The `ChannelService` is automatically initialized if configurations are present in `config.yaml`.
2.  **Send Notifications**: Use the `SendNotification` method in `IChannelService` to reach users across multiple platforms.
    ```go
    err := channelService.SendNotification(ctx, userId, "Your order #123 has been shipped!", []descriptors.ChannelType{descriptors.ChannelWhatsApp, descriptors.ChannelEmail})
    ```
3.  **Handle Inbound**: Implement custom logic for `HandleInbound` if you need to process platform-specific messages (e.g., "TRACK ORDER").

## Supported Channels
- **WhatsApp**: (via whatsadk)
- **Email**: (via mailadk)
- **Signal**
- **Telegram**
- **X (Twitter)**
- **Bluesky**
