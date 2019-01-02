package photon

type (
	ActionType   int64
	ContextStage int64
)

const (
	StageStructure ContextStage = iota
	StageDispatchBefore
	StageDispatchAfter
	StageDownloadBefore
	StageDownloadAfter
	StateCallbackBefore
	StateCallbackAfter
)
