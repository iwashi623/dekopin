package dekopin

var (
	ErrGetCommitHashInLocal = &GetCommitHashErrorInLocal{
		Message: "failed to get commit hash in local",
	}
)

type GetCommitHashErrorInLocal struct {
	Message string
}

func (e *GetCommitHashErrorInLocal) Error() string {
	return e.Message
}
