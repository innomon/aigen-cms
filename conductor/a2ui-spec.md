# A2UI Integration Specification

## Overview
This specification outlines the integration of the **Agent-to-User Interface (A2UI)** protocol into `formcms-go`. A2UI allows the backend ("The Agent") to stream rich, interactive UI structures as JSON data to the frontend ("The Client"), which renders them using a trusted component catalog.

## Core Architectural Components

### 1. The A2UI Protocol (Adjacency List Model)
Instead of nested JSON, we will implement the **Adjacency List** format (v0.8/v0.9). 
- **Message Structure**: A flat array of component objects.
- **Component ID**: Each component must have a unique `id`.
- **Children**: Container components (e.g., `Row`, `Column`, `Card`) will have a `children` attribute containing an array of child `id`s.

### 2. Backend: The Agent (Go)
- **A2UIService**: Manages the construction of UI payloads. It tracks the "Source of Truth" for the current UI state.
- **A2UIApi**: 
    - **Streaming Endpoint**: `/api/a2ui/stream` (using Server-Sent Events - SSE).
    - **Action Endpoint**: `/api/a2ui/action` (receives `userAction` signals).
- **Go Chi Integration**: Leverages existing middleware (Auth, RBAC) for secure streaming.

### 3. Frontend: The Renderer (JavaScript)
- **A2UI-Renderer.js**: A new core library that:
    1. Subscribes to the SSE stream.
    2. Maintains a local "Component Map" of the adjacency list.
    3. Recursively renders components starting from a root `Surface`.
- **Component Catalog**: A mapping of A2UI types to Bootstrap 5 HTML templates.
    - `Button` -> `<button class="btn btn-primary">`
    - `TextField` -> `<input class="form-control">`
    - `Card` -> `<div class="card"><div class="card-body">...</div></div>`

### 4. Interaction Loop
1. **User Action**: User clicks a button in an A2UI-rendered component.
2. **Signal**: Client sends a JSON payload to `/api/a2ui/action` containing `{ componentId, actionType, data }`.
3. **Reasoning**: The backend processes the action and streams back a "Delta" update (new components or modified attributes) via SSE.

## Data Schema Examples

### A2UI Component Payload
```json
[
  { "id": "root", "type": "Column", "children": ["title", "form"] },
  { "id": "title", "type": "Text", "attributes": { "content": "New Lead Entry", "style": "h3" } },
  { "id": "form", "type": "Card", "children": ["input-name", "submit-btn"] },
  { "id": "input-name", "type": "TextField", "attributes": { "label": "Full Name", "placeholder": "John Doe" } },
  { "id": "submit-btn", "type": "Button", "attributes": { "label": "Save Lead", "variant": "primary" } }
]
```

## Security Considerations
- **Sanitization**: The renderer MUST NOT use `eval()` or `innerHTML` directly with untrusted strings. Attributes like `content` will be set via `textContent`.
- **Trusted Catalog**: Only components registered in the client-side `ComponentCatalog` will be rendered. Unknown types will be ignored.
- **Authentication**: SSE connections will require valid JWT tokens.
