import Vue from "vue";
import Router from "vue-router";
import Index from "@/views/ms-index.js";
import store from "./store";

Vue.use(Router);

const GameBoard = () =>
  import(/* webpackChunkName: 'game' */ "./views/game/game-board.vue");
const PlayerCreate = () =>
  import(/* webpackChunkName: 'player' */ "./views/player/player-create.vue");

const router = new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [
    {
      path: "/",
      name: "index",
      component: Index,
      meta: {
        public: true
      }
    },
    {
      path: "/players/new",
      name: "player-create",
      component: PlayerCreate,
      meta: {
        public: true
      }
    },
    {
      path: "/games/:entityId",
      name: "game-board",
      component: GameBoard,
      props: true
    }
  ]
});

const routerSecurityCheck = (user, to, from, next, isFirstTime) => {
  let nextResolve;
  if (user) {
    if (to.meta && to.meta.public && to.name !== "index") {
      nextResolve = {
        name: "index"
      };
    } else if (to.meta.role) {
      const missingRole = user.roles.indexOf(to.meta.role) === -1;
      if (missingRole) {
        // user has not enough privileged to see the view
        // TODO define what to do here
        nextResolve = {
          name: "index"
        };
      }
    }
  } else if (to.meta && !to.meta.public) {
    if (isFirstTime) {
      store._actions["session/setFirstRoute"][0](to);
    }
    nextResolve = {
      name: "index"
    };
  }
  next(nextResolve);
};

router.beforeEach((to, from, next) => {
  const userRequested = store.getters["session/userRequested"];
  store._actions["notifications/clearScopedNotifications"][0]().then(() => {});
  if (!userRequested) {
    store._actions["session/requestUserLogged"][0]().then(() => {
      // the first route will be setted here
      const user = store.getters["session/user"];
      routerSecurityCheck(user, to, from, next, true);
    });
    return;
  }
  const user = store.getters["session/user"];
  routerSecurityCheck(user, to, from, next, false);
});

export default router;
