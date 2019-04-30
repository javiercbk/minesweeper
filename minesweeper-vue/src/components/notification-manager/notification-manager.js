import _ from "lodash";
import { mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters("notifications", ["notifications"]),
    notificationClasses() {
      return _.map(this.notifications, n => `alert-${n.variant}`);
    }
  }
};
