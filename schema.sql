-- CREATE DATABASE minesweeper WITH OWNER 'minesweeper' ENCODING 'UTF8';

CREATE TYPE mine_operation AS ENUM ('reveal', 'mark');


CREATE TABLE players(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_players_name ON players (name);

CREATE TABLE games(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    private BOOLEAN NOT NULL DEFAULT FALSE,
    rows SMALLINT NOT NULL,
    cols SMALLINT NOT NULL,
    mines SMALLINT NOT NULL,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    won BOOLEAN DEFAULT FALSE,
    creator_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT cnst_games_board CHECK (cols > 0 AND rows > 0 AND cols <= 100 AND rows <= 100),
    CONSTRAINT cnst_games_mines CHECK (mines > 0 AND (rows * cols) - 1 > mines),
    CONSTRAINT fk_games_creator FOREIGN KEY (creator_id) REFERENCES players (id)
);

CREATE TABLE game_operations(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    game_id BIGINT NOT NULL,
    player_id BIGINT NOT NULL,
    row SMALLINT NOT NULL,
    col SMALLINT NOT NULL,
    operation_id INTEGER NOT NULL,
    mine_proximity SMALLINT NOT NULL,
    operation mine_operation NOT NULL,
    CONSTRAINT fk_game_operation_game FOREIGN KEY (game_id) REFERENCES games (id),
    CONSTRAINT fk_games_creator FOREIGN KEY (player_id) REFERENCES players (id)
);

CREATE INDEX idx_game_operation_game ON game_operations (game_id);
CREATE UNIQUE INDEX idx_game_operation ON game_operations (game_id, operation_id);

CREATE TABLE game_board_points(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    game_id BIGINT NOT NULL,
    row SMALLINT NOT NULL,
    col SMALLINT NOT NULL,
    mine_proximity SMALLINT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT cnst_games_map_x_y CHECK (row >= 0 AND col >= 0 AND row < 100 AND col < 100),
    CONSTRAINT fk_games_map_game FOREIGN KEY (game_id) REFERENCES games (id)
);

CREATE UNIQUE INDEX idx_game_board ON game_board_points (game_id, row, col);
CREATE INDEX idx_game_board_mine_proximity ON game_board_points (game_id, mine_proximity);