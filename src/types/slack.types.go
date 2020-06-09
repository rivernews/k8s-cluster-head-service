package types

type SlackRequestType struct {
	Token       string `form:"token"`
	TriggerWord string `form:"trigger_word"`
	Text        string `form:"text"`
}
