import Vue from "vue";
import Vuex from "vuex";
import events from "./stores/events";
import notifications from "./stores/notifications";
import session from "./stores/session";

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    events,
    notifications,
    session
  }
});
