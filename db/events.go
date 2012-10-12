package db

import(
	"database/sql"
	"datagram.io/data"
)

func StoreEvents(eventStream *data.EventStream) *data.EventStream {

	go func () {
		for {
			if err := StoreEvent(<- eventStream.Events); err != nil {
				panic(err)
			}
		}
	}()	

	return newEvents
}

func StoreEvent(event *data.Event) (err error) {
	
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

	if err = tx.Stmt(insertEventStmt).QueryRow(event.Desc).Scan(&event.Id); err != nil {
		return
	}

	if err = tagEvents(tx, event.Id, tagIds); err != nil {
		return
	}

	newEvents.Events <- event

	return
}