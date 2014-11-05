package yorm

//yes-orm

/**
type Inner struct{
	Id int64
}
type Table struct{
	In Inner `yorm:"-"`
	//this struct's expected fields will represent as column too
	I1 Inner
	//this situation means Inner must implent the interface to decode value
	//and must implent the value to string
	I2 Inner `yorm:"column(name)"`
	//use i2 as column name
	I2 Inner `yorm:"column()"`
	//yorm assiagn default value to basic type ,as int string float64...
	//also the type implents Byte2Struct
	I3 int `yorm:"column();default()"`
}
**/
