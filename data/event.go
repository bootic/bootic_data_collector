package data

import "time"

type Event struct {
	//Uid    int `PK` //if the table's PrimaryKey is not id ,should add `PK` to ident
	Id   int64
	Name string
	App string
	User string
	Tags []string
	//Occurred
	CreatedOn time.Time
}
