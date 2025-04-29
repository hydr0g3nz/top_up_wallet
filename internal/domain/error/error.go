package err

import "errors"

var ErrNegativeAmount = errors.New("amount cannot be negative")
var ErrInsufficientBalance = errors.New("insufficient balance")
