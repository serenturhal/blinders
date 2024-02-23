package main

// distributeMessage
// TODO: receive an event with message and conversation id
// query all users' sessions
// dispatch message to all connected users
// else: dispatch via notification

// sendMessage
// TODO: user send message via sendMessage route
// store message to the conversation
// dispatch message by invoking distributeMessage

// updateMessageStatus
// TODO: receive events
// conversationID, messageID, status -> ""
// store to database
// distribute to user
