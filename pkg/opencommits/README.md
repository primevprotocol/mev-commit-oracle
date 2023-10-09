## OpenCommits

Open Commits pulls relevant data from our chain about active commits related to a builder at any given time.

```go
type Retriver interface {
	GetOpenCommits(builder string) ([]PreConfCommitment, error)
}
```