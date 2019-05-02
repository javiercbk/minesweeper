import { mapActions } from "vuex";
import MsUserForm from "@/components/forms/user/ms-user-form.vue";
import PlayerService from "@/services/player-service";

const playerService = new PlayerService();

export default {
  components: {
    MsUserForm
  },
  data() {
    return {
      creatingUser: false
    };
  },
  methods: {
    ...mapActions("notifications", ["pushNotification"]),
    handleUserCreate(newPlayer) {
      this.creatingUser = true;
      return playerService
        .create(newPlayer)
        .then(() => {
          this.pushNotification({
            title: "Success",
            message: "Player was created successfully",
            variant: "success"
          });
          this.$router.push({
            name: "index"
          });
        })
        .catch(() => {
          this.pushNotification({
            title: "Error",
            message: "could not create player",
            variant: "danger"
          });
        })
        .finally(() => {
          this.creatingUser = false;
        });
    }
  }
};
