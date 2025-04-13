package syntax

type TokenRole int

const (
	TokenRoleNone = TokenRole(iota)
	TokenRoleOperator
	TokenRoleKeyword
	TokenRoleNumber
	TokenRoleString
	TokenRoleComment
)
