package data

type Event struct {
  //Uid    int `PK` //if the table's PrimaryKey is not id ,should add `PK` to ident
  Id   int64
  Desc string
  Tags []string
  //Occurred
  //Created    time.Time
}
