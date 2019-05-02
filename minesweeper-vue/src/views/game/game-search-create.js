import GameService from "@/services/game-service";
import MsGameForm from "@/components/forms/game/ms-game-form.vue";
import { mapActions } from "vuex";

const gameService = new GameService();

export default {
  components: {
    MsGameForm
  },
  created() {
    this.findGames();
  },
  data() {
    return {
      creatingGame: false,
      loadingGames: false,
      availableGames: []
    };
  },
  methods: {
    ...mapActions("notifications", ["pushNotification"]),
    findGames() {
      this.loadingGames = true;
      return gameService
        .find()
        .then(res => res.json())
        .then(body => {
          this.availableGames = body.data.games;
        })
        .catch(() => {
          this.pushNotification({
            title: "Error",
            message: "could not retrieve games",
            variant: "danger"
          });
        })
        .finally(() => {
          this.loadingGames = false;
        });
    },
    onGameSubmit(game) {
      this.creatingGame = true;
      return gameService
        .create(game)
        .then(res => res.json())
        .then(body => {
          // when game is created navigate to the new game automatically
          const gameCreated = body.data.games;
          this.$router.push({
            name: "game-board",
            params: {
              entityId: gameCreated.id
            }
          });
        })
        .catch(() => {
          this.pushNotification({
            title: "Error",
            message: "could not create a new game",
            variant: "danger"
          });
        })
        .finally(() => {
          this.creatingGame = false;
        });
    },
    playGame(gameId) {
      this.$router.push({
        name: "game-board",
        params: {
          entityId: gameId
        }
      });
    }
  }
};
