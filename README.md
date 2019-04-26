# Minesweeper

A web based minesweeper game. It can create boards of n x m sizes (n and m are integers larger then 0 and smaller or equal than 100). A player can create private or public games and solve the mine sweeping with a friend. Games are saved upon interaction and can be resumed.

# Game Rules

- This is a 2D minesweeper.
- To create a game you need to specify a board size, and an amount of mines.
- The mines are placed randomly at game creation, meaning that an unlucky player CAN hit a mine in the first click.
- The first action that the player makes, a timer will start ticking. The timer will continue ticking regardless if the user is looking at the board or not.
- The player can mark unrevealed places where a mine might be.
- The game is saved at every step and the player cannot rollback any action.
- Upon winning or losing a game, all mines are shown but no further actions can be performed on the game.

# Conventions

- Dates are stored in UTC-0 timezone, no exception.

# Data Entities

###### Player

Using a unique name and a password, the user can start creating games or join public games.

###### Game

A game defines whether a game is public or not, when it was started, which player is the creator and whether it has finished or not and what was the end result. Also it contains the amount of mines and the board size.

* See "Database model" section for more details

###### Game board points

A game board 2D points (game map) contains all the points in space for a game and a `mine_proximity` value associated.

* See "Database model" inside the Architecture section for more details

## Minesweep algebra

The `mine_proximity` value is an integer that describes the state of a point in space. To interpret that number there are certain rules:

- A positive number means that the state of that point in space is revealed.
- A negative number means that the state of that point in space is unrevealed.
- Values from -8 to 6 mean that the point has no mine and the absolute value signals how many mines are next (in any direction) to the point in space.
- Values -8 and 7 means that there is an unrevealed or revealed mine in that point of space.
- Values from -11 to -18 means that there is a marked *possible* mine point in space. Marking a point is the same as subtracting 10, thus unmarking it is the same as adding 10.
- Values from -21 to -28 means that there is a marked mine point in space. Marking a point is the same as subtracting -20, thus unmarking it is the same as adding 20.

With this in mind you can define the *minesweep algebra*:

Let:

- **x, y** as integers which are valid points inside the game board.
- **GSID** as an integer which is the game state id.
- **O** as the game's state **open**, meaning that the game is not won nor lost.
- **W** as the game's state **won**.
- **L** as the game's state **lost**.
- **M** as the mines amount.
- **U** as the amount of unrevealed points.
- **z** as the `mine_proximity`.
- There are only three operations, **reveal**, **mark**, and **compose**.
- (x, y, GSID) reveal (O, z) / z (E in [-7..6]) => |z| IF U == M then game status changes to W (game won)
- (x, y, GSID) reveal (O, z) / z is -8 => 8 AND game status changes to L (game lost)
- (x, y, GSID) reveal (W, z) / z => Operation not allowed, game has concluded (game won)
- (x, y, GSID) reveal (L, z) / z => Operation not allowed, game has concluded (game lost)
- (x, y, GSID) mark (O, z) / z (E in [0..7]) => z
- (x, y, GSID) mark (O, z) / z (E in [-8..-1]) => z - 10
- (x, y, GSID) mark (O, z) / z (E in [-18..-11]) => z - 10
- (x, y, GSID) mark (O, z) / z (E in [-28..-21]) => z + 20
- (x, y, GSID) mark (O, z) / z (E not in [-27..7]) => Operation not allowed, cannot mark a field which has already been revealed or marked
- (x, y, GSID) mark (W, z) / z => Operation not allowed, game has concluded (game won)
- (x, y, GSID) mark (L, z) / z => Operation not allowed, game has concluded (game lost)
- compose(reveal(x1,y1), reveal(x1,y1)) = reveal(x1,y1)
- compose(reveal(x1,y1), reveal(x2,y2)) = [reveal(x1,y1), reveal(x2,y2)]
- compose(mark(x1,y1), mark(x1,y1)) = mark(x1,y1)
- compose(mark(x1,y1), mark(x2,y2)) = [mark(x1,y1), mark(x2,y2)]
- compose(reveal(x1,y1), mark(x1,y1)) = reveal(x1,y1)
- compose(mark(x1,y1), reveal(x1,y1)) = mark(x1,y1)
- compose(reveal(x1,y1), mark(x2,y2)) = [reveal(x1,y1), mark(x2,y2)]
- compose(mark(x1,y1), reveal(x2,y2)) = [mark(x1,y1), reveal(x2,y2)]

By defining the *minesweep algebra*, you can define a **minesweep game** as a game state that has been applied a finite amount of operations that can be composed until you calculate the current game state.

The *minesweep algebra* has several advantages:

- It allows the application to rely on operation composition rather than locks to manage the game state.
- It prevents data races by making explicit conflicts resolutions.
- It allows a game to be decomposed into several steps that can be undone or replayed.

It only has one disadvantage

- Both client and server need to implement the full *minesweep algebra* to work properly and any change made to the algebra will force changes in both server and client.


## Architecture

The application is basically a REST API that validates and executes operations of the *minesweep algebra* on a game state, such validations are:

- Validate that the player owns the game or is playing on a public game.
- Validate if the *minesweep operation* is valid within the board (do not allow out of bounds operations)
- Validate if the *minesweep operation* is valid within the game state (do not allow operations on games won or lost).

If the operation is valid it will always result in a board change **unless** the reveal operation is performed on an already revealed field. This is important when taking into account multiple users.

Upon connecting to a game, the player receives the game state and starts receiving operations to update the local game state.

The synchronization strategy between players is optimistic, applying [OT](https://en.wikipedia.org/wiki/Operational_transformation) to keep the local and server game state synchronized. The operations to be transformed are the defined in the *minesweep algebra*.

For each operation sent and received there is an operation id attached. Such id increases with each operation and it allows the client and the server to know whether they are synchronized, if the server notices that there is a client sending an operation with an already taken operation id, it can send back all the operations that the client has missed plus it can apply the operation without forcing the client to re-send the operation.

The server must not give away any more information than the revealed 2D points, the board size, the amount of mines and the first operation date (this last value is used to calculate the time spent playing). The whole board cannot be stored in the client because it would allow cheaters to read the board and know where the mines are, that is why clients only receive the points already revealed. Only when the game has finished the server can send all the board to the clients.

Clients are authenticated using a JWT token that expires on 15 minutes (without any token blacklist). Client's Passwords are hashed using bcrypt.

Clients will send operations to the server via websockets or a REST API.

### Database model

One of the aims of this project is to benchmark which data modeling approach works best. There are two ways of modeling the game board, normalized and denormalized.

The **normalized** approach stores every 2D point of the game board as a row in the `game_board_points` table. This is standard relational modeling and it is guarantee to work.

The **denormalized** approach stores every 2D point of the game board in a matrix inside the `games` table. No subquery is needed to retrieve the game board and the data is easily mapped into the application.

I have certain hypothesis:

- Denormalized approach would be faster to read and first write (creation would be faster) but slower to update.
- Denormalized approach would be slower to check the `mine_proximity` of a single point.
- Denormalized approach would make the algorithm to auto reveal *zero mines*.
- Normalized approach would be faster to read a single 2D point and faster to apply operations on a row level.

Are my hypothesis correct? I don't know, but I can quickly benchmark this.

## TODO

- [x] Analyze and write down the solution specification.
- [x] Write the database schema
- [x] Implement the *minesweep algebra* in the backend (GO).
- [ ] Implement the game backend API and benchmark the normalized and denormalized approach.
- [ ] Implement authentication.
- [ ] Create a login page in the frontend .
- [ ] Implement the *minesweep algebra* in javascript.
- [ ] Implement a client for the backend API.
- [ ] Create the game board in the frontend using the client backend API.
