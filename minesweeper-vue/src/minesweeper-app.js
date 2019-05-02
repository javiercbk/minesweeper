import { mapGetters, mapActions } from "vuex";
import MsTopBar from "@/components/top-bar/ms-top-bar.vue";
import LoadingApp from "@/components/loading/loading-animation.vue";
import NotificationManager from "@/components/notification-manager/notification-manager.vue";

export default {
  components: {
    MsTopBar,
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
