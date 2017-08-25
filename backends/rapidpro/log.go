package rapidpro

import (
	"fmt"
	"strings"

	"time"

	"github.com/nyaruka/courier"
)

const insertLogSQL = `
INSERT INTO channels_channellog("channel_id", "msg_id", "description", "is_error", "method", "url", "request", "response", "response_status", "created_on", "request_time")
                         VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

// WriteChannelLog writes the passed in channel log to the database, we do not queue on errors but instead just throw away the log
func writeChannelLog(b *backend, log *courier.ChannelLog) error {
	// cast our channel to our own channel type
	dbChan, isChan := log.Channel.(*DBChannel)
	if !isChan {
		return fmt.Errorf("unable to write non-rapidpro channel logs")
	}

	description := "Success"
	if log.Error != "" {
		description = "Error"

		// we append our error to our response as it can be long
		log.Response += "\n\nError: " + log.Error
	}

	// strip null chars from request and response, postgres doesn't like that
	log.Request = strings.Trim(log.Request, "\x00")
	log.Response = strings.Trim(log.Request, "\x00")

	_, err := b.db.Exec(insertLogSQL, dbChan.ID(), log.MsgID, description, log.Error != "", log.Method, log.URL,
		log.Request, log.Response, log.StatusCode, log.CreatedOn, log.Elapsed/time.Millisecond)

	return err
}
