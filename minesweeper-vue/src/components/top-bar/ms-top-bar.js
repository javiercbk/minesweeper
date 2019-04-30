import _ from "lodash";
import { mapActions, mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters("session", ["user"]),
    userName() {
      return _.get(this, "user.name");
    },
    isGameCreation() {
      return this.$router.currentRoute.name === "game-creation";
    }
  },
  methods: {
    ...mapActions("session", ["logout"]),
    performLogout() {
      this.logout().then(() => {
        this.$router.replace({
          name: "index"
        });
      });
    }
  }
};
