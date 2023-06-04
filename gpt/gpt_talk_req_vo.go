package gpt

type TalkVO struct {
	Prompt            string `json:"prompt"`
	Model             string `json:"model"`
	Message_id        string `json:"message_id"`
	Parent_message_id string `json:"parent_message_id"`
	Conversation_id   string `json:"conversation_id"`
	Stream            bool   `json:"stream"`
}

type ConversationsVO struct {
	Offset string `json:"offset"`
	Limit  string `json:"limit"`
}

type Json struct {
	Model      string `json:"model"`
	Message_id string `json:"message_id"`
}
