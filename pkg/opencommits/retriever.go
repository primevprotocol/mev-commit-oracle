package opencommits

type PreConfCommitment struct {
	TxnHash        string
	CommitmentHash string
}

type Retriver interface {
	GetOpenCommits(builder string) ([]PreConfCommitment, error)
}

type retriver struct {
	endpoint string
}
