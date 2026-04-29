package model

type Position string

const (
	Goalkeeper Position = "goalkeeper"
	Defender   Position = "defender"
	Midfielder Position = "midfielder"
	Forward    Position = "forward"
	Any        Position = "any"
)
