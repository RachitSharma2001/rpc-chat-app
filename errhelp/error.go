package errhelp

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrUnableToDeleteUser = Error("Unable to delete user - user does not exist!")
	ErrUnableToUpdateUser = Error("Unable to update user - user does not exist!")
)

func ErrorExists(err error) bool {
	return err != nil
}
