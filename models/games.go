// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Game is an object representing the database table.
type Game struct {
	ID         int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	Private    bool      `boil:"private" json:"private" toml:"private" yaml:"private"`
	Cols       int16     `boil:"cols" json:"cols" toml:"cols" yaml:"cols"`
	Rows       int16     `boil:"rows" json:"rows" toml:"rows" yaml:"rows"`
	Mines      int16     `boil:"mines" json:"mines" toml:"mines" yaml:"mines"`
	StartedAt  null.Time `boil:"started_at" json:"startedAt,omitempty" toml:"startedAt" yaml:"startedAt,omitempty"`
	FinishedAt null.Time `boil:"finished_at" json:"finishedAt,omitempty" toml:"finishedAt" yaml:"finishedAt,omitempty"`
	Won        null.Bool `boil:"won" json:"won,omitempty" toml:"won" yaml:"won,omitempty"`
	CreatorID  int64     `boil:"creator_id" json:"creatorID" toml:"creatorID" yaml:"creatorID"`
	CreatedAt  null.Time `boil:"created_at" json:"createdAt,omitempty" toml:"createdAt" yaml:"createdAt,omitempty"`
	UpdatedAt  null.Time `boil:"updated_at" json:"updatedAt,omitempty" toml:"updatedAt" yaml:"updatedAt,omitempty"`
	R          *gameR    `boil:"-" json:"-" toml:"-" yaml:"-"`
	L          gameL     `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GameColumns = struct {
	ID         string
	Private    string
	Cols       string
	Rows       string
	Mines      string
	StartedAt  string
	FinishedAt string
	Won        string
	CreatorID  string
	CreatedAt  string
	UpdatedAt  string
}{
	ID:         "id",
	Private:    "private",
	Cols:       "cols",
	Rows:       "rows",
	Mines:      "mines",
	StartedAt:  "started_at",
	FinishedAt: "finished_at",
	Won:        "won",
	CreatorID:  "creator_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// Generated where

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelpernull_Bool struct{ field string }

func (w whereHelpernull_Bool) EQ(x null.Bool) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Bool) NEQ(x null.Bool) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Bool) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Bool) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_Bool) LT(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Bool) LTE(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Bool) GT(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Bool) GTE(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var GameWhere = struct {
	ID         whereHelperint64
	Private    whereHelperbool
	Cols       whereHelperint16
	Rows       whereHelperint16
	Mines      whereHelperint16
	StartedAt  whereHelpernull_Time
	FinishedAt whereHelpernull_Time
	Won        whereHelpernull_Bool
	CreatorID  whereHelperint64
	CreatedAt  whereHelpernull_Time
	UpdatedAt  whereHelpernull_Time
}{
	ID:         whereHelperint64{field: `id`},
	Private:    whereHelperbool{field: `private`},
	Cols:       whereHelperint16{field: `cols`},
	Rows:       whereHelperint16{field: `rows`},
	Mines:      whereHelperint16{field: `mines`},
	StartedAt:  whereHelpernull_Time{field: `started_at`},
	FinishedAt: whereHelpernull_Time{field: `finished_at`},
	Won:        whereHelpernull_Bool{field: `won`},
	CreatorID:  whereHelperint64{field: `creator_id`},
	CreatedAt:  whereHelpernull_Time{field: `created_at`},
	UpdatedAt:  whereHelpernull_Time{field: `updated_at`},
}

// GameRels is where relationship names are stored.
var GameRels = struct {
	Creator         string
	GameBoardPoints string
	GameOperations  string
}{
	Creator:         "Creator",
	GameBoardPoints: "GameBoardPoints",
	GameOperations:  "GameOperations",
}

// gameR is where relationships are stored.
type gameR struct {
	Creator         *Player
	GameBoardPoints GameBoardPointSlice
	GameOperations  GameOperationSlice
}

// NewStruct creates a new relationship struct
func (*gameR) NewStruct() *gameR {
	return &gameR{}
}

// gameL is where Load methods for each relationship are stored.
type gameL struct{}

var (
	gameColumns               = []string{"id", "private", "cols", "rows", "mines", "started_at", "finished_at", "won", "creator_id", "created_at", "updated_at"}
	gameColumnsWithoutDefault = []string{"cols", "rows", "mines", "started_at", "finished_at", "creator_id", "created_at", "updated_at"}
	gameColumnsWithDefault    = []string{"id", "private", "won"}
	gamePrimaryKeyColumns     = []string{"id"}
)

type (
	// GameSlice is an alias for a slice of pointers to Game.
	// This should generally be used opposed to []Game.
	GameSlice []*Game

	gameQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	gameType                 = reflect.TypeOf(&Game{})
	gameMapping              = queries.MakeStructMapping(gameType)
	gamePrimaryKeyMapping, _ = queries.BindMapping(gameType, gameMapping, gamePrimaryKeyColumns)
	gameInsertCacheMut       sync.RWMutex
	gameInsertCache          = make(map[string]insertCache)
	gameUpdateCacheMut       sync.RWMutex
	gameUpdateCache          = make(map[string]updateCache)
	gameUpsertCacheMut       sync.RWMutex
	gameUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single game record from the query.
func (q gameQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Game, error) {
	o := &Game{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for games")
	}

	return o, nil
}

// All returns all Game records from the query.
func (q gameQuery) All(ctx context.Context, exec boil.ContextExecutor) (GameSlice, error) {
	var o []*Game

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Game slice")
	}

	return o, nil
}

// Count returns the count of all Game records in the query.
func (q gameQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count games rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q gameQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if games exists")
	}

	return count > 0, nil
}

// Creator pointed to by the foreign key.
func (o *Game) Creator(mods ...qm.QueryMod) playerQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.CreatorID),
	}

	queryMods = append(queryMods, mods...)

	query := Players(queryMods...)
	queries.SetFrom(query.Query, "\"players\"")

	return query
}

// GameBoardPoints retrieves all the game_board_point's GameBoardPoints with an executor.
func (o *Game) GameBoardPoints(mods ...qm.QueryMod) gameBoardPointQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"game_board_points\".\"game_id\"=?", o.ID),
	)

	query := GameBoardPoints(queryMods...)
	queries.SetFrom(query.Query, "\"game_board_points\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"game_board_points\".*"})
	}

	return query
}

// GameOperations retrieves all the game_operation's GameOperations with an executor.
func (o *Game) GameOperations(mods ...qm.QueryMod) gameOperationQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"game_operations\".\"game_id\"=?", o.ID),
	)

	query := GameOperations(queryMods...)
	queries.SetFrom(query.Query, "\"game_operations\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"game_operations\".*"})
	}

	return query
}

// LoadCreator allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (gameL) LoadCreator(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGame interface{}, mods queries.Applicator) error {
	var slice []*Game
	var object *Game

	if singular {
		object = maybeGame.(*Game)
	} else {
		slice = *maybeGame.(*[]*Game)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &gameR{}
		}
		args = append(args, object.CreatorID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &gameR{}
			}

			for _, a := range args {
				if a == obj.CreatorID {
					continue Outer
				}
			}

			args = append(args, obj.CreatorID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`players`), qm.WhereIn(`id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Player")
	}

	var resultSlice []*Player
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Player")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for players")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for players")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Creator = foreign
		if foreign.R == nil {
			foreign.R = &playerR{}
		}
		foreign.R.CreatorGames = append(foreign.R.CreatorGames, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CreatorID == foreign.ID {
				local.R.Creator = foreign
				if foreign.R == nil {
					foreign.R = &playerR{}
				}
				foreign.R.CreatorGames = append(foreign.R.CreatorGames, local)
				break
			}
		}
	}

	return nil
}

// LoadGameBoardPoints allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (gameL) LoadGameBoardPoints(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGame interface{}, mods queries.Applicator) error {
	var slice []*Game
	var object *Game

	if singular {
		object = maybeGame.(*Game)
	} else {
		slice = *maybeGame.(*[]*Game)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &gameR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &gameR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`game_board_points`), qm.WhereIn(`game_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load game_board_points")
	}

	var resultSlice []*GameBoardPoint
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice game_board_points")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on game_board_points")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for game_board_points")
	}

	if singular {
		object.R.GameBoardPoints = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &gameBoardPointR{}
			}
			foreign.R.Game = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.GameID {
				local.R.GameBoardPoints = append(local.R.GameBoardPoints, foreign)
				if foreign.R == nil {
					foreign.R = &gameBoardPointR{}
				}
				foreign.R.Game = local
				break
			}
		}
	}

	return nil
}

// LoadGameOperations allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (gameL) LoadGameOperations(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGame interface{}, mods queries.Applicator) error {
	var slice []*Game
	var object *Game

	if singular {
		object = maybeGame.(*Game)
	} else {
		slice = *maybeGame.(*[]*Game)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &gameR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &gameR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`game_operations`), qm.WhereIn(`game_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load game_operations")
	}

	var resultSlice []*GameOperation
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice game_operations")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on game_operations")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for game_operations")
	}

	if singular {
		object.R.GameOperations = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &gameOperationR{}
			}
			foreign.R.Game = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.GameID {
				local.R.GameOperations = append(local.R.GameOperations, foreign)
				if foreign.R == nil {
					foreign.R = &gameOperationR{}
				}
				foreign.R.Game = local
				break
			}
		}
	}

	return nil
}

// SetCreator of the game to the related item.
// Sets o.R.Creator to related.
// Adds o to related.R.CreatorGames.
func (o *Game) SetCreator(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Player) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"games\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"creator_id"}),
		strmangle.WhereClause("\"", "\"", 2, gamePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CreatorID = related.ID
	if o.R == nil {
		o.R = &gameR{
			Creator: related,
		}
	} else {
		o.R.Creator = related
	}

	if related.R == nil {
		related.R = &playerR{
			CreatorGames: GameSlice{o},
		}
	} else {
		related.R.CreatorGames = append(related.R.CreatorGames, o)
	}

	return nil
}

// AddGameBoardPoints adds the given related objects to the existing relationships
// of the game, optionally inserting them as new records.
// Appends related to o.R.GameBoardPoints.
// Sets related.R.Game appropriately.
func (o *Game) AddGameBoardPoints(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*GameBoardPoint) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.GameID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"game_board_points\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"game_id"}),
				strmangle.WhereClause("\"", "\"", 2, gameBoardPointPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.GameID = o.ID
		}
	}

	if o.R == nil {
		o.R = &gameR{
			GameBoardPoints: related,
		}
	} else {
		o.R.GameBoardPoints = append(o.R.GameBoardPoints, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &gameBoardPointR{
				Game: o,
			}
		} else {
			rel.R.Game = o
		}
	}
	return nil
}

// AddGameOperations adds the given related objects to the existing relationships
// of the game, optionally inserting them as new records.
// Appends related to o.R.GameOperations.
// Sets related.R.Game appropriately.
func (o *Game) AddGameOperations(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*GameOperation) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.GameID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"game_operations\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"game_id"}),
				strmangle.WhereClause("\"", "\"", 2, gameOperationPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.GameID = o.ID
		}
	}

	if o.R == nil {
		o.R = &gameR{
			GameOperations: related,
		}
	} else {
		o.R.GameOperations = append(o.R.GameOperations, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &gameOperationR{
				Game: o,
			}
		} else {
			rel.R.Game = o
		}
	}
	return nil
}

// Games retrieves all the records using an executor.
func Games(mods ...qm.QueryMod) gameQuery {
	mods = append(mods, qm.From("\"games\""))
	return gameQuery{NewQuery(mods...)}
}

// FindGame retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGame(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Game, error) {
	gameObj := &Game{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"games\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, gameObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from games")
	}

	return gameObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Game) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no games provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		if queries.MustTime(o.UpdatedAt).IsZero() {
			queries.SetScanner(&o.UpdatedAt, currTime)
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(gameColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	gameInsertCacheMut.RLock()
	cache, cached := gameInsertCache[key]
	gameInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			gameColumns,
			gameColumnsWithDefault,
			gameColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(gameType, gameMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(gameType, gameMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"games\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"games\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into games")
	}

	if !cached {
		gameInsertCacheMut.Lock()
		gameInsertCache[key] = cache
		gameInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Game.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Game) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	key := makeCacheKey(columns, nil)
	gameUpdateCacheMut.RLock()
	cache, cached := gameUpdateCache[key]
	gameUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			gameColumns,
			gamePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update games, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"games\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, gamePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(gameType, gameMapping, append(wl, gamePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update games row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for games")
	}

	if !cached {
		gameUpdateCacheMut.Lock()
		gameUpdateCache[key] = cache
		gameUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q gameQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for games")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for games")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GameSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gamePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"games\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, gamePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in game slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all game")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Game) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no games provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	nzDefaults := queries.NonZeroDefaultSet(gameColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	gameUpsertCacheMut.RLock()
	cache, cached := gameUpsertCache[key]
	gameUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			gameColumns,
			gameColumnsWithDefault,
			gameColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			gameColumns,
			gamePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert games, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(gamePrimaryKeyColumns))
			copy(conflict, gamePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"games\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(gameType, gameMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(gameType, gameMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert games")
	}

	if !cached {
		gameUpsertCacheMut.Lock()
		gameUpsertCache[key] = cache
		gameUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Game record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Game) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Game provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), gamePrimaryKeyMapping)
	sql := "DELETE FROM \"games\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from games")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for games")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q gameQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no gameQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from games")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for games")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GameSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Game slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gamePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"games\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, gamePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from game slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for games")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Game) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGame(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GameSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GameSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gamePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"games\".* FROM \"games\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, gamePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in GameSlice")
	}

	*o = slice

	return nil
}

// GameExists checks if the Game row exists.
func GameExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"games\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if games exists")
	}

	return exists, nil
}
