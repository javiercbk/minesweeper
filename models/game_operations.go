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
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// GameOperation is an object representing the database table.
type GameOperation struct {
	ID            int64           `boil:"id" json:"id" toml:"id" yaml:"id"`
	GameID        int64           `boil:"game_id" json:"gameID" toml:"gameID" yaml:"gameID"`
	PlayerID      int64           `boil:"player_id" json:"playerID" toml:"playerID" yaml:"playerID"`
	OperationID   int             `boil:"operation_id" json:"operationID" toml:"operationID" yaml:"operationID"`
	Operation     string          `boil:"operation" json:"operation" toml:"operation" yaml:"operation"`
	Row           int16           `boil:"row" json:"row" toml:"row" yaml:"row"`
	Col           int16           `boil:"col" json:"col" toml:"col" yaml:"col"`
	MineProximity int16           `boil:"mine_proximity" json:"mineProximity" toml:"mineProximity" yaml:"mineProximity"`
	R             *gameOperationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L             gameOperationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GameOperationColumns = struct {
	ID            string
	GameID        string
	PlayerID      string
	OperationID   string
	Operation     string
	Row           string
	Col           string
	MineProximity string
}{
	ID:            "id",
	GameID:        "game_id",
	PlayerID:      "player_id",
	OperationID:   "operation_id",
	Operation:     "operation",
	Row:           "row",
	Col:           "col",
	MineProximity: "mine_proximity",
}

// Generated where

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

var GameOperationWhere = struct {
	ID            whereHelperint64
	GameID        whereHelperint64
	PlayerID      whereHelperint64
	OperationID   whereHelperint
	Operation     whereHelperstring
	Row           whereHelperint16
	Col           whereHelperint16
	MineProximity whereHelperint16
}{
	ID:            whereHelperint64{field: `id`},
	GameID:        whereHelperint64{field: `game_id`},
	PlayerID:      whereHelperint64{field: `player_id`},
	OperationID:   whereHelperint{field: `operation_id`},
	Operation:     whereHelperstring{field: `operation`},
	Row:           whereHelperint16{field: `row`},
	Col:           whereHelperint16{field: `col`},
	MineProximity: whereHelperint16{field: `mine_proximity`},
}

// GameOperationRels is where relationship names are stored.
var GameOperationRels = struct {
	Game   string
	Player string
}{
	Game:   "Game",
	Player: "Player",
}

// gameOperationR is where relationships are stored.
type gameOperationR struct {
	Game   *Game
	Player *Player
}

// NewStruct creates a new relationship struct
func (*gameOperationR) NewStruct() *gameOperationR {
	return &gameOperationR{}
}

// gameOperationL is where Load methods for each relationship are stored.
type gameOperationL struct{}

var (
	gameOperationColumns               = []string{"id", "game_id", "player_id", "operation_id", "operation", "row", "col", "mine_proximity"}
	gameOperationColumnsWithoutDefault = []string{"game_id", "player_id", "operation_id", "operation", "row", "col", "mine_proximity"}
	gameOperationColumnsWithDefault    = []string{"id"}
	gameOperationPrimaryKeyColumns     = []string{"id"}
)

type (
	// GameOperationSlice is an alias for a slice of pointers to GameOperation.
	// This should generally be used opposed to []GameOperation.
	GameOperationSlice []*GameOperation

	gameOperationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	gameOperationType                 = reflect.TypeOf(&GameOperation{})
	gameOperationMapping              = queries.MakeStructMapping(gameOperationType)
	gameOperationPrimaryKeyMapping, _ = queries.BindMapping(gameOperationType, gameOperationMapping, gameOperationPrimaryKeyColumns)
	gameOperationInsertCacheMut       sync.RWMutex
	gameOperationInsertCache          = make(map[string]insertCache)
	gameOperationUpdateCacheMut       sync.RWMutex
	gameOperationUpdateCache          = make(map[string]updateCache)
	gameOperationUpsertCacheMut       sync.RWMutex
	gameOperationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single gameOperation record from the query.
func (q gameOperationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*GameOperation, error) {
	o := &GameOperation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for game_operations")
	}

	return o, nil
}

// All returns all GameOperation records from the query.
func (q gameOperationQuery) All(ctx context.Context, exec boil.ContextExecutor) (GameOperationSlice, error) {
	var o []*GameOperation

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to GameOperation slice")
	}

	return o, nil
}

// Count returns the count of all GameOperation records in the query.
func (q gameOperationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count game_operations rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q gameOperationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if game_operations exists")
	}

	return count > 0, nil
}

// Game pointed to by the foreign key.
func (o *GameOperation) Game(mods ...qm.QueryMod) gameQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.GameID),
	}

	queryMods = append(queryMods, mods...)

	query := Games(queryMods...)
	queries.SetFrom(query.Query, "\"games\"")

	return query
}

// Player pointed to by the foreign key.
func (o *GameOperation) Player(mods ...qm.QueryMod) playerQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.PlayerID),
	}

	queryMods = append(queryMods, mods...)

	query := Players(queryMods...)
	queries.SetFrom(query.Query, "\"players\"")

	return query
}

// LoadGame allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (gameOperationL) LoadGame(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGameOperation interface{}, mods queries.Applicator) error {
	var slice []*GameOperation
	var object *GameOperation

	if singular {
		object = maybeGameOperation.(*GameOperation)
	} else {
		slice = *maybeGameOperation.(*[]*GameOperation)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &gameOperationR{}
		}
		args = append(args, object.GameID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &gameOperationR{}
			}

			for _, a := range args {
				if a == obj.GameID {
					continue Outer
				}
			}

			args = append(args, obj.GameID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`games`), qm.WhereIn(`id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Game")
	}

	var resultSlice []*Game
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Game")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for games")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for games")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Game = foreign
		if foreign.R == nil {
			foreign.R = &gameR{}
		}
		foreign.R.GameOperations = append(foreign.R.GameOperations, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.GameID == foreign.ID {
				local.R.Game = foreign
				if foreign.R == nil {
					foreign.R = &gameR{}
				}
				foreign.R.GameOperations = append(foreign.R.GameOperations, local)
				break
			}
		}
	}

	return nil
}

// LoadPlayer allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (gameOperationL) LoadPlayer(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGameOperation interface{}, mods queries.Applicator) error {
	var slice []*GameOperation
	var object *GameOperation

	if singular {
		object = maybeGameOperation.(*GameOperation)
	} else {
		slice = *maybeGameOperation.(*[]*GameOperation)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &gameOperationR{}
		}
		args = append(args, object.PlayerID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &gameOperationR{}
			}

			for _, a := range args {
				if a == obj.PlayerID {
					continue Outer
				}
			}

			args = append(args, obj.PlayerID)

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
		object.R.Player = foreign
		if foreign.R == nil {
			foreign.R = &playerR{}
		}
		foreign.R.GameOperations = append(foreign.R.GameOperations, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.PlayerID == foreign.ID {
				local.R.Player = foreign
				if foreign.R == nil {
					foreign.R = &playerR{}
				}
				foreign.R.GameOperations = append(foreign.R.GameOperations, local)
				break
			}
		}
	}

	return nil
}

// SetGame of the gameOperation to the related item.
// Sets o.R.Game to related.
// Adds o to related.R.GameOperations.
func (o *GameOperation) SetGame(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Game) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"game_operations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"game_id"}),
		strmangle.WhereClause("\"", "\"", 2, gameOperationPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.GameID = related.ID
	if o.R == nil {
		o.R = &gameOperationR{
			Game: related,
		}
	} else {
		o.R.Game = related
	}

	if related.R == nil {
		related.R = &gameR{
			GameOperations: GameOperationSlice{o},
		}
	} else {
		related.R.GameOperations = append(related.R.GameOperations, o)
	}

	return nil
}

// SetPlayer of the gameOperation to the related item.
// Sets o.R.Player to related.
// Adds o to related.R.GameOperations.
func (o *GameOperation) SetPlayer(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Player) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"game_operations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"player_id"}),
		strmangle.WhereClause("\"", "\"", 2, gameOperationPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PlayerID = related.ID
	if o.R == nil {
		o.R = &gameOperationR{
			Player: related,
		}
	} else {
		o.R.Player = related
	}

	if related.R == nil {
		related.R = &playerR{
			GameOperations: GameOperationSlice{o},
		}
	} else {
		related.R.GameOperations = append(related.R.GameOperations, o)
	}

	return nil
}

// GameOperations retrieves all the records using an executor.
func GameOperations(mods ...qm.QueryMod) gameOperationQuery {
	mods = append(mods, qm.From("\"game_operations\""))
	return gameOperationQuery{NewQuery(mods...)}
}

// FindGameOperation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGameOperation(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*GameOperation, error) {
	gameOperationObj := &GameOperation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"game_operations\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, gameOperationObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from game_operations")
	}

	return gameOperationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GameOperation) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no game_operations provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(gameOperationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	gameOperationInsertCacheMut.RLock()
	cache, cached := gameOperationInsertCache[key]
	gameOperationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			gameOperationColumns,
			gameOperationColumnsWithDefault,
			gameOperationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(gameOperationType, gameOperationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(gameOperationType, gameOperationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"game_operations\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"game_operations\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into game_operations")
	}

	if !cached {
		gameOperationInsertCacheMut.Lock()
		gameOperationInsertCache[key] = cache
		gameOperationInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the GameOperation.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GameOperation) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	gameOperationUpdateCacheMut.RLock()
	cache, cached := gameOperationUpdateCache[key]
	gameOperationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			gameOperationColumns,
			gameOperationPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update game_operations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"game_operations\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, gameOperationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(gameOperationType, gameOperationMapping, append(wl, gameOperationPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update game_operations row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for game_operations")
	}

	if !cached {
		gameOperationUpdateCacheMut.Lock()
		gameOperationUpdateCache[key] = cache
		gameOperationUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q gameOperationQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for game_operations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for game_operations")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GameOperationSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gameOperationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"game_operations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, gameOperationPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in gameOperation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all gameOperation")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GameOperation) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no game_operations provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(gameOperationColumnsWithDefault, o)

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

	gameOperationUpsertCacheMut.RLock()
	cache, cached := gameOperationUpsertCache[key]
	gameOperationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			gameOperationColumns,
			gameOperationColumnsWithDefault,
			gameOperationColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			gameOperationColumns,
			gameOperationPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert game_operations, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(gameOperationPrimaryKeyColumns))
			copy(conflict, gameOperationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"game_operations\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(gameOperationType, gameOperationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(gameOperationType, gameOperationMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert game_operations")
	}

	if !cached {
		gameOperationUpsertCacheMut.Lock()
		gameOperationUpsertCache[key] = cache
		gameOperationUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single GameOperation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GameOperation) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no GameOperation provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), gameOperationPrimaryKeyMapping)
	sql := "DELETE FROM \"game_operations\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from game_operations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for game_operations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q gameOperationQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no gameOperationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from game_operations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for game_operations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GameOperationSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no GameOperation slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gameOperationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"game_operations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, gameOperationPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from gameOperation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for game_operations")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *GameOperation) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGameOperation(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GameOperationSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GameOperationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gameOperationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"game_operations\".* FROM \"game_operations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, gameOperationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in GameOperationSlice")
	}

	*o = slice

	return nil
}

// GameOperationExists checks if the GameOperation row exists.
func GameOperationExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"game_operations\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if game_operations exists")
	}

	return exists, nil
}
