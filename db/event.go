package db

import(
	"database/sql"
	"bytes"
	"fmt"
)

type Event struct {
	//Uid		 int `PK` //if the table's PrimaryKey is not id ,should add `PK` to ident
	Id	 int
	Desc string
	Tags []string
	//Occurred
	//Created		 time.Time
}

var (
	insertEventStmt     *sql.Stmt
	insertTagStmt	      *sql.Stmt
	findTagIdStmt	      *sql.Stmt
	insertEventTagStmt  *sql.Stmt
)

func prepareEventStatements() {
	
	insertEventStmt = prepareOrPanic(`
		INSERT INTO "event" ("desc") VALUES ($1) RETURNING "id";`)

	insertTagStmt = prepareOrPanic(`
		INSERT INTO "tag" ("key") VALUES ($1) RETURNING "id";`)

	findTagIdStmt = prepareOrPanic(`
		SELECT "id" FROM "tag" WHERE "tag"."key"=$1;`)

	insertEventTagStmt = prepareOrPanic(`
		INSERT INTO "event_tag" ("event_id", "tag_id") VALUES ($1, $2);`)
}

func findTag (tag string) (id int64, notFound bool, err error) {
	
	if err = findTagIdStmt.QueryRow(tag).Scan(&id); err != nil {

		if err.Error() == "sql: no rows in result set" {
			notFound = true
			err = nil
		} else {
			err = fmt.Errorf("Error finding tag '%s': %s", tag, err.Error())
		}
	} 

	return
}

func findOrCreateTags(tags []string) (tagIds []int64, err error) {

	for _, tag := range tags {

		var ( tagId int64
				  notFound bool )

		if tagId, notFound, err = findTag(tag); err != nil {
			return
		}
		

		if notFound {

	    if err = insertTagStmt.QueryRow(tag).Scan(&tagId); err != nil {
				err = fmt.Errorf("Error inserting tag '%s': %s", tag, err.Error())
				return
			}
		}

		tagIds = append(tagIds, tagId)
  }
  
	return 
}

func tagEvents(tx *sql.Tx, eventId int64, tagIds []int64) (err error) {
	
	if len(tagIds) == 0 {
		return
	}

	var buffer bytes.Buffer

	buffer.WriteString(`INSERT INTO "event_tag" ("event_id", "tag_id") VALUES `)

	for i, tagId := range tagIds {
		
		buffer.WriteString(fmt.Sprintf("(%d, %d)", eventId, tagId))
		
		if (i < (len(tagIds) - 1)) {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString(";")

	if _, err = tx.Exec(buffer.String()); err != nil {
		return fmt.Errorf("Error tagging event: %s\nQuery:\n\n%s", err.Error(), buffer.String())
	}

	return
}

func StoreEvent(event *Event) (err error) {
	
	// insert the tags and bail if we get an error
	// TODO, do this within a transaction
	var tagIds []int64
	if tagIds, err = findOrCreateTags(event.Tags); err != nil {
		return
	}

  var tx *sql.Tx

	if tx, err = pg.Begin(); err != nil {
		return
	}

  defer /* panic-recover */ func() {
		
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var eventId int64
	if err = tx.Stmt(insertEventStmt).QueryRow(event.Desc).Scan(&eventId); err != nil {
		return
	}

	if err = tagEvents(tx, eventId, tagIds); err != nil {
		return
	}

	return
}
