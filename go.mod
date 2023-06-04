module sf-bot

go 1.20

require github.com/eatmoreapple/openwechat v1.4.3

require (
	github.com/go-ini/ini v1.67.0
	github.com/sf-bot/gpt v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
)

replace github.com/sf-bot/handler => ./handler

replace github.com/sf-bot/gpt => ./gpt
