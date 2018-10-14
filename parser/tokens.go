package parser

type userName struct {
	name         string
	discriminant int
}

type userReference struct {
	uid string
}

type channelReference struct {
	channelID string
}
