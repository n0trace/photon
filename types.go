package photon

type (
	ActionType   int64
	ContextStage int64
)

const (
	OnRequest ActionType = iota
	OnResponse
	OnError

	StageStructure ContextStage = iota
	StageDispatchBefore
	StageDispatchAfter
	StageDownloadBefore
	StageDownloadAfter
)
