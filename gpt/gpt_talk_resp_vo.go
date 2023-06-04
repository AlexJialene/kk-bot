package gpt

type TalkRespVO struct {
	Conversation_id string
	Error           string
	Message         TalkMessageRespVO
}

type TalkMessageRespVO struct {
	Id        string
	End_turn  bool
	Recipient string
	Status    string
	Content   TalkMessageContentRespVO
}

type TalkMessageContentRespVO struct {
	Content_type string
	Parts        []string
}

type Title struct {
	Title string
}
