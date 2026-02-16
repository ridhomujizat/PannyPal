# Webhook Payload Guide

## Overview

This application now sends **simplified, clean webhook payloads** to your configured webhook URLs. The new format is easy to parse and consistent across all message types.

---

## What Changed?

### Before (Complex Baileys Proto)
```json
{
  "sessionId": "my-session",
  "from": "628123456789@s.whatsapp.net",
  "messageType": "conversation",
  "message": {
    "conversation": "Hello!",
    "messageContextInfo": { ... },
    "...": "many nested proto fields"
  },
  "timestamp": 1708172400,
  "key": {
    "remoteJid": "628123456789@s.whatsapp.net",
    "id": "3EB0C6D3F2A4E8B9",
    "fromMe": false
  }
}
```

### After (Simplified Clean Format)
```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "text",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "text",
    "text": "Hello!"
  }
}
```

---

## Webhook Payload Structure

All webhook payloads follow this consistent structure:

```typescript
{
  // Session & Basic Info
  sessionId: string;        // Your session identifier
  messageId: string;        // Unique message ID
  timestamp: number;        // Unix timestamp (seconds)

  // Sender Info
  from: string;             // Clean phone number (no @s.whatsapp.net)
  fromName?: string;        // Sender's WhatsApp display name

  // Message Type
  messageType: 'text' | 'image' | 'video' | 'audio' | 'document' |
               'location' | 'contact' | 'sticker' | 'unknown';

  // Chat Context
  isGroup: boolean;         // true if from group chat
  chatId: string;           // Full JID for reference
  participant?: string;     // Sender's number in group (for groups only)

  // Content (varies by message type)
  content: { ... };

  // Reply Context (optional)
  quotedMessage?: {
    messageId: string;
    text?: string;
  };
}
```

---

## Message Types & Examples

### 1. Text Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "text",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "text",
    "text": "Hello, how are you?"
  }
}
```

### 2. Image Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "image",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "image",
    "caption": "Check this out!",
    "mimetype": "image/jpeg",
    "fileSize": 524288,
    "downloadInstructions": {
      "method": "POST",
      "endpoint": "/api/message/download-media",
      "body": {
        "sessionId": "my-session",
        "messageId": "3EB0C6D3F2A4E8B9"
      }
    }
  }
}
```

### 3. Video Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "video",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "video",
    "caption": "Watch this!",
    "mimetype": "video/mp4",
    "fileSize": 2097152,
    "duration": 30,
    "downloadInstructions": {
      "method": "POST",
      "endpoint": "/api/message/download-media",
      "body": {
        "sessionId": "my-session",
        "messageId": "3EB0C6D3F2A4E8B9"
      }
    }
  }
}
```

### 4. Audio/Voice Note Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "audio",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "audio",
    "mimetype": "audio/ogg; codecs=opus",
    "fileSize": 65536,
    "duration": 15,
    "downloadInstructions": {
      "method": "POST",
      "endpoint": "/api/message/download-media",
      "body": {
        "sessionId": "my-session",
        "messageId": "3EB0C6D3F2A4E8B9"
      }
    }
  }
}
```

### 5. Document Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "document",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "document",
    "caption": "Here's the report",
    "mimetype": "application/pdf",
    "fileName": "report.pdf",
    "fileSize": 1048576,
    "downloadInstructions": {
      "method": "POST",
      "endpoint": "/api/message/download-media",
      "body": {
        "sessionId": "my-session",
        "messageId": "3EB0C6D3F2A4E8B9"
      }
    }
  }
}
```

### 6. Location Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "location",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "location",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "name": "Monas",
    "address": "Jakarta, Indonesia"
  }
}
```

### 7. Contact Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "contact",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "contact",
    "displayName": "Jane Smith",
    "vcard": "BEGIN:VCARD\nVERSION:3.0\nFN:Jane Smith\nTEL:+628987654321\nEND:VCARD",
    "phones": ["+628987654321"]
  }
}
```

### 8. Group Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0C6D3F2A4E8B9",
  "timestamp": 1708172400,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "text",
  "isGroup": true,
  "chatId": "120363123456789012@g.us",
  "participant": "628123456789",
  "content": {
    "type": "text",
    "text": "Hello everyone!"
  }
}
```

### 9. Reply/Quoted Message

```json
{
  "sessionId": "my-session",
  "messageId": "3EB0D4E5F6A7B8C9",
  "timestamp": 1708172500,
  "from": "628123456789",
  "fromName": "John Doe",
  "messageType": "text",
  "isGroup": false,
  "chatId": "628123456789@s.whatsapp.net",
  "content": {
    "type": "text",
    "text": "Yes, I agree!"
  },
  "quotedMessage": {
    "messageId": "3EB0AABBCCDDEE",
    "text": "Should we proceed with the plan?"
  }
}
```

---

## Downloading Media

Media messages (image, video, audio, document, sticker) include `downloadInstructions` that tell you how to download the actual file.

### Method 1: Using the Simplified Format (Recommended)

```bash
curl -X POST http://your-server/api/message/download-media \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "my-session",
    "messageId": "3EB0C6D3F2A4E8B9"
  }' \
  --output downloaded-media.jpg
```

### Method 2: Get as Base64 JSON

```bash
curl -X POST http://your-server/api/message/download-media \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "my-session",
    "messageId": "3EB0C6D3F2A4E8B9",
    "returnBase64": true
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "mimetype": "image/jpeg",
    "base64": "/9j/4AAQSkZJRgABAQEAYABgAAD...",
    "size": 524288
  }
}
```

### Media Cache Duration

- Messages are cached for **24 hours** after receipt
- After 24 hours, you'll get a "Message not found in cache" error
- Download media promptly or store it in your own system

---

## Integration Examples

### Python (Flask)

```python
from flask import Flask, request, jsonify
import requests

app = Flask(__name__)

@app.route('/webhook', methods=['POST'])
def webhook():
    payload = request.json

    # Extract basic info
    session_id = payload['sessionId']
    message_type = payload['messageType']
    from_number = payload['from']
    from_name = payload.get('fromName', 'Unknown')

    # Handle different message types
    if message_type == 'text':
        text = payload['content']['text']
        print(f"Text from {from_name}: {text}")

    elif message_type == 'image':
        caption = payload['content'].get('caption', '')
        message_id = payload['messageId']

        # Download the image
        download_response = requests.post(
            'http://your-server/api/message/download-media',
            json={
                'sessionId': session_id,
                'messageId': message_id,
                'returnBase64': True
            }
        )

        if download_response.ok:
            data = download_response.json()['data']
            print(f"Image from {from_name}: {caption}, Size: {data['size']} bytes")

    elif message_type == 'location':
        lat = payload['content']['latitude']
        lon = payload['content']['longitude']
        print(f"Location from {from_name}: {lat}, {lon}")

    return jsonify({'status': 'received'}), 200

if __name__ == '__main__':
    app.run(port=3000)
```

### Node.js (Express)

```javascript
const express = require('express');
const axios = require('axios');

const app = express();
app.use(express.json());

app.post('/webhook', async (req, res) => {
  const payload = req.body;

  const { sessionId, messageType, from, fromName, content, messageId } = payload;

  console.log(`Message from ${fromName || from}: ${messageType}`);

  // Handle text messages
  if (messageType === 'text') {
    console.log(`Text: ${content.text}`);
  }

  // Handle media messages
  if (['image', 'video', 'audio', 'document'].includes(messageType)) {
    try {
      const response = await axios.post(
        'http://your-server/api/message/download-media',
        { sessionId, messageId },
        { responseType: 'arraybuffer' }
      );

      console.log(`Downloaded ${messageType}, Size: ${response.data.length} bytes`);
      // Save to file, upload to S3, etc.
    } catch (error) {
      console.error('Failed to download media:', error.message);
    }
  }

  // Handle location
  if (messageType === 'location') {
    const { latitude, longitude, name } = content;
    console.log(`Location: ${name} (${latitude}, ${longitude})`);
  }

  res.json({ status: 'received' });
});

app.listen(3000, () => {
  console.log('Webhook server listening on port 3000');
});
```

---

## Benefits

✅ **Clean phone numbers** - No more `@s.whatsapp.net` suffixes
✅ **Predictable structure** - Same top-level fields for all message types
✅ **Easy parsing** - No need to understand Baileys proto objects
✅ **Type safety** - Clear TypeScript types available
✅ **Simple media download** - Just sessionId + messageId
✅ **Group support** - Clear `isGroup` flag and `participant` field
✅ **Reply detection** - Automatic extraction of quoted messages

---

## TypeScript Types

If you're building a TypeScript webhook consumer, you can use these types:

```typescript
interface SimplifiedWebhookPayload {
  sessionId: string;
  messageId: string;
  timestamp: number;
  from: string;
  fromName?: string;
  messageType: 'text' | 'image' | 'video' | 'audio' | 'document' | 'location' | 'contact' | 'sticker' | 'unknown';
  isGroup: boolean;
  chatId: string;
  participant?: string;
  content: WebhookMessageContent;
  quotedMessage?: {
    messageId: string;
    text?: string;
  };
}

type WebhookMessageContent =
  | TextContent
  | MediaContent
  | LocationContent
  | ContactContent
  | UnknownContent;

interface TextContent {
  type: 'text';
  text: string;
}

interface MediaContent {
  type: 'image' | 'video' | 'audio' | 'document' | 'sticker';
  caption?: string;
  mimetype?: string;
  fileName?: string;
  fileSize?: number;
  duration?: number;
  downloadInstructions: {
    method: 'POST';
    endpoint: string;
    body: {
      sessionId: string;
      messageId: string;
    };
  };
}

interface LocationContent {
  type: 'location';
  latitude: number;
  longitude: number;
  name?: string;
  address?: string;
}

interface ContactContent {
  type: 'contact';
  displayName: string;
  vcard: string;
  phones?: string[];
}

interface UnknownContent {
  type: 'unknown';
  description: string;
}
```

---

## Migration Notes

### Backward Compatibility

The old `IncomingMessage` format is still available in the TypeScript types for reference, but **all new webhook payloads use the simplified format**.

### If You Need the Old Format

If you have existing integrations that require the old Baileys proto format, you'll need to update your webhook consumers to use the new simplified format. The new format is much easier to work with and includes all the essential information.

---

## Support

For issues or questions about webhook payloads, please check:
1. This guide
2. TypeScript types in `src/types/index.ts`
3. The GitHub repository issues section
