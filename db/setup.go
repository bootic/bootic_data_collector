package db

import (
	"fmt"
)

const setup string = `
  
  DROP TABLE IF EXISTS event_tag;
  DROP TABLE IF EXISTS event;
  DROP TABLE IF EXISTS tag;

  CREATE TABLE event (
    "id"   SERIAL PRIMARY KEY, 
    "desc" text
  );
  
  CREATE TABLE tag (
    "id"  SERIAL PRIMARY KEY,
    "key" text
  );

  CREATE UNIQUE INDEX tag_key ON tag (key);

  CREATE TABLE event_tag (
    "id"       SERIAL PRIMARY KEY,
    "event_id" integer,
    "tag_id"   integer
  );
  
  CREATE UNIQUE INDEX event_tag_pairing ON event_tag (event_id, tag_id);`

func SetupDB() (err error) {

	if _, err = pg.Exec(setup); err != nil {
		return
	}

	fmt.Println("done.")

	return
}
