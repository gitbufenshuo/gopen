package commmsg

type JumpMSGTurn struct {
	Turn int64
	List []JumpMSGOne
}

type JumpMSGOne struct {
	Kind     string // move jump choose login
	UID      string
	Which    int64 // 哪一个
	MoveValX int64
	MoveValZ int64
}
