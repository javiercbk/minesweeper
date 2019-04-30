import { mapGetters, mapActions } from "vuex";
import LoadingApp from "./views/loading-app.vue";
import NotificationManager from "./components/notification-manager/notification-manager.vue";

export default {
  components: {
    LoadingApp,
    NotificationManager
  },
  computed: {
    ...mapGetters("session", ["user", "userRequested"])
  },
  methods: {
    ...mapActions("events", ["fireGlobalEvent"]),
    onGlobalEvent(event, where) {
      this.fireGlobalEvent({
        event,
        where
      });
    }
  }
};
