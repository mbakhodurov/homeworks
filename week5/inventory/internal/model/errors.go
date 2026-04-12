package model

import "errors"

var ErrInventoryNotFound = errors.New("inventory not found")
var ErrInventoryAlreadyExists = errors.New("inventory already exists")
