import Vue from "vue";
import Vuelidate from "vuelidate";
import Promise from "bluebird";
import { library } from "@fortawesome/fontawesome-svg-core";
import {
  faCog,
  faSmile,
  faSurprise,
  faDizzy,
  faLaughSquint,
  faFlag,
  faQuestion,
  faSpinner
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import App from "@/minesweeper-app.vue";
import router from "@/router";
import store from "@/store";
import TransparentDirective from "@/directives/ms-transparent.js";

window.Promise = Promise;

library.add(
  faSmile,
  faDizzy,
  faSurprise,
  faFlag,
  faLaughSquint,
  faQuestion,
  faCog,
  faSpinner
);
Vue.component("font-awesome-icon", FontAwesomeIcon);

Vue.config.productionTip = false;

Vue.use(Vuelidate);

Vue.directive("transparent", TransparentDirective);

new Vue({
  router,
  store,
  validations: {},
  render: h => h(App)
}).$mount("#app");
