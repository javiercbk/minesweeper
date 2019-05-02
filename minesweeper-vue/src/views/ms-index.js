import _ from "lodash";
import store from "@/store";
import MsLogin from "@/components/forms/login/ms-login.vue";
import MsHome from "@/views/home/ms-home.vue";

export default {
  functional: true,
  name: "MsIndex",
  render(h) {
    const userLogged = store.getters["app/userLogged"];
    if (_.isNil(userLogged)) {
      return h(MsLogin);
    }
    return h(MsHome);
  }
};
