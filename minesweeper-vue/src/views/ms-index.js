import _ from "lodash";
import store from "@/store";
import MsLogin from "@/components/forms/login/ms-login.vue";
import GameSearchCreate from "@/views/game/game-search-create.vue";

export default {
  functional: true,
  name: "MsIndex",
  render(h) {
    const user = store.getters["session/user"];
    if (_.isNil(user)) {
      return h(MsLogin);
    }
    return h(GameSearchCreate);
  }
};
