package builtin

import "errors"


var ErrMissingPath error = errors.New("cd: Missing path")
var ErrTooManyValues error = errors.New("cd: too many arguments")