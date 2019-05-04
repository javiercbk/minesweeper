import _ from "lodash";
import { mapActions } from "vuex";
import GameService from "@/services/game-service";
import MsGameBoard from "@/components/game/ms-game-board.vue";
import moment from "moment";

const DEFAULT_STATUS = {
  won: false,
  lost: false
};

const STATE_NOT_REVEALED = 0;
const STATE_SUSPECT_MINE = 1;
const STATE_MARKED_MINE = 2;
const STATE_REVEALED = 3;

const gameService = new GameService();

export default {
  components: {
    MsGameBoard
  },
  props: {
    entityId: {
      type: [String, Number],
      required: true
    }
  },
  created() {
    this.retrieveGame();
  },
  data() {
    return {
      loadingGame: false,
      synchronizing: false,
      game: {
        rows: 0,
        cols: 0,
        mines: 0,
        board: [],
        latestOperationId: 0,
        startedAt: null,
        finishedAt: null,
        won: false
      }
    };
  },
  computed: {
    gameId() {
      let id = parseInt(this.entityId, 10);
      if (_.isNaN(id)) {
        id = null;
      }
      return id;
    },
    gameStatus() {
      if (!_.isNil(this.game.finishedAt)) {
        if (this.game.won) {
          return "won";
        } else {
          return "lost";
        }
      }
      return null;
    }
  },
  methods: {
    ...mapActions("notifications", ["pushNotification"]),
    retrieveGame() {
      const id = parseInt(this.entityId, 10);
      if (_.isNaN(id)) {
        this.pushNotification({
          title: "Error",
          message: "The game does not exist",
          variant: "danger"
        });
        return;
      }
      this.loadingGame = true;
      gameService
        .retrieve(this.entityId)
        .then(res => res.json())
        .then(body => {
          this.game = body.data.game;
          if (_.isNil(this.game.latestOperationId)) {
            this.game.latestOperationId = 0;
          }
        })
        .catch(() => {
          this.pushNotification({
            title: "Error",
            message: "could not retrieve game",
            variant: "danger"
          });
        })
        .finally(() => {
          this.loadingGame = false;
        });
    },
    onOperation(operation) {
      this.synchronizing = true;
      operation.id = this.game.latestOperationId + 1;
      operation.gameId = this.gameId;
      gameService
        .patch(operation)
        .then(res => res.json())
        .then(body => {
          this._processConfirmation(body.data.confirmation);
        })
        .catch(() => {
          this.pushNotification({
            title: "Error",
            message: "could not apply operation, please try again.",
            variant: "warning"
          });
        })
        .finally(() => {
          this.synchronizing = false;
        });
    },
    _processConfirmation(confirmation) {
      let latestOperationId = this.game.latestOperationId;
      let boardChanged = false;
      const board = _.cloneDeep(this.game.board);
      const applyOperation = function(id) {
        return function(r) {
          const { row, col, mineProximity, pointState } = r;
          if (!_.isNil(id) && id > latestOperationId) {
            latestOperationId = id;
          }
          if (pointState === STATE_SUSPECT_MINE && board[row][col] !== "?") {
            boardChanged = true;
            board[row][col] = "?";
          } else if (
            pointState === STATE_MARKED_MINE &&
            board[row][col] !== "!"
          ) {
            boardChanged = true;
            board[row][col] = "!";
          } else if (
            pointState === STATE_REVEALED &&
            typeof board[row][col] !== "number"
          ) {
            boardChanged = true;
            board[row][col] = mineProximity;
          } else if (
            pointState === STATE_NOT_REVEALED &&
            board[row][col] !== null
          ) {
            boardChanged = true;
            board[row][col] = null;
          }
        };
      };
      const newStatus = _.get(confirmation, "status", DEFAULT_STATUS);
      const operationId = _.get(confirmation, "operation.id", 0);
      if (operationId > 0) {
        latestOperationId = operationId;
      }
      const opResults = _.get(confirmation, "operation.result", []);
      const deltaOperations = _.get(confirmation, "deltaOperations", []);
      if (newStatus.lost || newStatus.won) {
        this.game.won = !newStatus.lost && newStatus.won;
        this.game.finishedAt = _.get(
          newStatus,
          "finishedAt".moment.utc().format()
        );
        this.game.board = newStatus.board;
      }
      _.forEach(opResults, applyOperation(operationId));
      _.forEach(deltaOperations, d => {
        const deltaResult = _.get(d, "result", []);
        _.forEach(deltaResult, r => {
          applyOperation(d.id)(r);
        });
      });
      if (boardChanged) {
        this.game.board = board;
      }
      this.game.latestOperationId = latestOperationId;
    }
  }
};
