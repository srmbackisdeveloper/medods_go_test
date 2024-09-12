package helpers

const (
	WrongAT      = "wrong access token"
	FailureAT    = "access token failure"
	WrongRT      = "wrong refresh token"
	FailureRT    = "refresh token failure"
	WrongRequest = "wrong request"
	GuidNotFound = "uid not found"
	GuidRequired = "user id is required (guid param)"
	WrongDB      = "database error"
)

func DbError(e string) string {
	return WrongDB + ": " + e
}
