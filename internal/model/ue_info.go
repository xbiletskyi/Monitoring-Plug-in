package model

type UeInfo struct {
	CellID int64
	Tac    [3]byte
	PlmnID [6]byte
	Msin   [10]byte
	Imei   [14]byte
}
