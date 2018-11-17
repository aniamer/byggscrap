package models

type Result struct {
	Result []ObjectInfo
	TotalCount int
}

type ObjectInfo struct {
	ObjectTypeDescription string
	ObjectArea string
	ObjectFloor string
	RentPerMonth string
	SeekAreaDescription string
	StreetName string
	EndPeriodMPDateString string

}

