package domain

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrEmptyEmail = errors.New("email not provided")

var ErrHealthRecordTypeAlreadyExists = errors.New("code already exists")
var ErrHealthRecordTypeNotFound = errors.New("record type not found")
var ErrCodeAlreadyExists = errors.New("code already exists")
var ErrHealthRecordTypeImmutable = errors.New("health record type can not be modified")

var ErrNoInformation = errors.New("no information provided")
var ErrInternalServer = errors.New("database internal server error")
var ErrEmptyId = errors.New("id not provided")

var ErrCodeRequired = errors.New("code attribute is required")
var ErrCodeTooLong = errors.New("code attribute is too long (limit 50)")
var ErrCodeInvalid = errors.New("code is invalid")

var ErrInvalidUserOrHealthRecordType = errors.New("invalid UserID or HealthRecordTypeID")
var ErrHealthRecordAlreadyExists = errors.New("health record already exists")
var ErrHealthRecordNotFound = errors.New("no health record found")
var ErrorUserIDNotProvided = errors.New("userid not provided")
