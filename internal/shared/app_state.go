package shared

import (
	"database/sql"
	"log"

	amqp_utils "github.com/ishola-faazele/taskflow/internal/utils/amqp"
	utils_db "github.com/ishola-faazele/taskflow/internal/utils/db"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AppState struct {
	DB       *sql.DB
	AmqpConn *amqp.Connection
}

func NewAppState() *AppState {
	db, err := sqlx.Connect("pgx", "user=taskflow_user password=taskflow_password dbname=taskflow_db sslmode=disable port=5432 host=localhost")
	if err != nil {
		log.Fatalln("FAILED_TO_CONNECT_TO_DB:", err)
	}
	// connect to rabbitmq
	conn := amqp_utils.InitAMQP()
	// initialize tables
	migrationMgr := utils_db.NewMigrationManager(db.DB)
	if err := migrationMgr.EnsureTablesExist(); err != nil {
		log.Fatalln("FAILED_TO_INITIALIZE_TABLES:", err)
	}
	return &AppState{
		DB:       db.DB,
		AmqpConn: conn,
	}
}
func (as *AppState) Clean() {
	as.DB.Close()
	as.AmqpConn.Close()
}
