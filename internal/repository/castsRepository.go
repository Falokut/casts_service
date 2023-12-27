package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type postgreRepository struct {
	db *sqlx.DB
}

const (
	castTableName = "casts"
)

func NewCastsRepository(db *sqlx.DB) *postgreRepository {
	return &postgreRepository{db: db}
}

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *postgreRepository) Shutdown() {
	r.db.Close()
}

func (r *postgreRepository) GetCast(ctx context.Context, id int32) (Cast, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgreRepository.GetCast")
	defer span.Finish()

	query := fmt.Sprintf("SELECT movie_id, array_agg(actor_id) AS actors_ids FROM %s WHERE movie_id=$1 GROUP BY movie_id", castTableName)

	var actors string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&id, &actors)
	actors = strings.Trim(actors, "{}")
	if errors.Is(err, sql.ErrNoRows) || strings.EqualFold(actors, "NULL") {
		return Cast{}, ErrNotFound
	}
	if err != nil {
		return Cast{}, err
	}

	actorsIDs := strings.Split(actors, ",")
	ids := make([]int32, 0, len(actorsIDs))
	for _, id := range actorsIDs {
		actorID, err := strconv.Atoi(id)
		if err != nil {
			return Cast{}, err
		}
		ids = append(ids, int32(actorID))
	}
	return Cast{Actors: ids}, nil
}
