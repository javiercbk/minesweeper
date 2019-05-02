import _ from "lodash";
import { mapGetters } from "vuex";
import MsTopBar from "@/components/top-bar/ms-top-bar.vue";

export default {
  components: {
    MsTopBar
  },
  computed: {
    ...mapGetters("session", ["user"]),
    userName() {
      return _.get(this, "user.name");
    }
  }
};
