package protos

type BoardSpeaker interface {
	ReciveFromBoard(*BoardMessage)
}

type MBSpeaker interface {
	ReciveFromMainboard(*MainboardMessage)
}

type ServerCommandBus interface {
	GetCommandsBus() MainboardCommandChannel
}

type Speaker interface {
	BoardSpeaker
	MBSpeaker
}

type GlobalBus interface {
	Speaker
	ServerCommandBus
}
