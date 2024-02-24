# Chat flow with Websocket

## Send a message

1. User send and `send message` event to client, with a `resolveId` which is used to resolve the ack response from server

2. Server constructs a message document, concurrently does:
    - store the message to database
    - send ack response event with the `resolveId`
    - retrieve sessions of all members in the conversation
        - if session exists, send to recipients
        - else dispatch to notification service

3. Sender resolves the response with `resolveId`

4. Recipients add messages to correspond conversation