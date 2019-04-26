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
    private BOOLEAN DEFAULT FALSE,
    board_x_size SMALLINT NOT NULL,
    board_y_size SMALLINT NOT NULL,
    map [][] SMALLINT,
    mines SMALLINT,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    won BOOLEAN DEFAULT FALSE,
    creator_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT cnst_games_board CHECK (board_x_size > 0 AND board_y_size > 0 AND board_x_size <= 100 AND board_y_size <= 100),
    CONSTRAINT cnst_games_mines CHECK (mines > 0 AND (board_x_size * board_y_size) - 1 > mines),
    CONSTRAINT fk_games_creator FOREIGN KEY (creator_id) REFERENCES players (id)
);

CREATE TABLE game_operations(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    game_id BIGINT NOT NULL,
    player_id BIGINT NOT NULL,
    operation_id INTEGER NOT NULL,
    operation mine_operation NOT NULL,
    CONSTRAINT fk_game_operation_game FOREIGN KEY (game_id) REFERENCES games (id),
    CONSTRAINT fk_games_creator FOREIGN KEY (player_id) REFERENCES players (id)
);

CREATE UNIQUE INDEX idx_game_operation ON players (game_id, operation_id);

CREATE TABLE game_board_points(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    game_id BIGINT NOT NULL,
    x SMALLINT NOT NULL,
    y SMALLINT NOT NULL,
    mine_proximity SMALLINT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT cnst_games_map_x_y CHECK (x > 0 AND y > 0 AND x <= 100 AND y <= 100),
    CONSTRAINT fk_games_map_game FOREIGN KEY (game_id) REFERENCES games (id)
);