package database

import (
	"database/sql"
	"log"
	"math/rand"
	"time"
)

func Alive(db *sql.DB) {
	log.Println("Connecting to database... ")
	for {
		// Ping by itself is un-reliable, the connections are cached. This
		// ensures that the database is still running by executing a harmless
		// dummy query against it.
		//
		// Also, we perform an exponential backoff when querying the database
		// to spread our retries.
		_, err := db.Exec("SELECT true")
		if err == nil {
			log.Println("Database connected")
			return
		}

		base, cap := time.Second, time.Minute
		for backoff := base; err != nil; backoff <<= 1 {
			if backoff > cap {
				backoff = cap
			}

			// A pseudo-random number generator here is fine. No need to be
			// cryptographically secure. Ignore with the following comment:
			/* #nosec */
			jitter := rand.Int63n(int64(backoff * 3))
			sleep := base + time.Duration(jitter)
			time.Sleep(sleep)
			_, err := db.Exec("SELECT true")
			if err == nil {
				log.Println("Database connected")
				return
			}
		}
	}
}
