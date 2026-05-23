package dto

type HealthRecordTypeInput struct {
	Name string
	Code string
	Unit *string
}

type HealthRecordTypeUpdateInput struct {
	Name *string
	Code *string
	Unit *string
}
