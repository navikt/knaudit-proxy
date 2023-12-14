package main

import (
	"database/sql/driver"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	goora "github.com/sijms/go-ora/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/report", func(c *gin.Context) {
		auditDataBytes, err := c.GetRawData()
		if err != nil {
			log.WithError(err).Error("unable to read audit data")
			c.JSON(400, gin.H{"status": "bad request"})
			return
		}

		if err := sendAuditDataToDVH(string(auditDataBytes)); err != nil {
			log.WithError(err).Error("error storing audit data")
			c.JSON(500, gin.H{"status": "error storing audit data"})
			return
		}

		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

func sendAuditDataToDVH(blob string) error {
	connection, err := goora.NewConnection(os.Getenv("ORACLE_URL"))
	if err != nil {
		return fmt.Errorf("failed creating new connection to Oracle: %v", err)
	}

	if err = connection.Open(); err != nil {
		return fmt.Errorf("failed opening connection to Oracle: %v", err)
	}

	defer connection.Close()

	stmt := goora.NewStmt("begin dvh_dmo.knaudit_api.log(p_event_document => :1); end;", connection)
	defer stmt.Close()

	rows, err := stmt.Query([]driver.Value{blob})
	if err != nil {
		return fmt.Errorf("failed executing statement: %v", err)
	}
	defer rows.Close()

	return nil
}
