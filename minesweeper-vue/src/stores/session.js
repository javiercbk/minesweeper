/* eslint-disable no-param-reassign */
// global window
import AuthService from "@/services/auth-service";

let authService;
authService;

const ensureAuthService = () => {
  if (!authService) {
    authService = new AuthService();
  }
};

const JWT_TOKEN_KEY = "jwt-token";

const _state = {
  user: null,
  viewport: null,
  userRequested: false,
  jwtToken: window.localStorage.getItem(JWT_TOKEN_KEY),
  versionChanged: false,
  firstRoute: null,
  version: null
};

const getters = {
  user: storeState => storeState.user,
  jwtToken: storeState => storeState.jwtToken,
  userRequested: storeState => storeState.userRequested,
  versionChanged: storeState => storeState.versionChanged,
  firstRoute: storeState => storeState.firstRoute,
  version: storeState => storeState.version,
  viewport: storeState => storeState.viewport
};

const mutations = {
  setUser: (storeState, payload) => {
    storeState.user = payload;
  },
  setUserRequested: (storeState, payload) => {
    storeState.userRequested = payload;
  },
  setjwtToken: (storeState, payload) => {
    window.localStorage.setItem(JWT_TOKEN_KEY, payload);
    storeState.jwtToken = payload;
  },
  setFirstRoute: (storeState, payload) => {
    storeState.firstRoute = payload;
  },
  setVersion: (storeState, payload) => {
    // only update the version once
    if (!storeState.version) {
      storeState.version = payload;
    }
  },
  setVersionChanged: storeState => {
    // only update the version once
    storeState.versionChanged = true;
  },
  setViewport: (storeState, payload) => {
    storeState.viewport = payload;
  }
};

const actions = {
  requestUserLogged: ({ commit }) => {
    ensureAuthService();
    return authService
      .getCurrentUser()
      .then(response => response.json())
      .then(body => {
        commit("setUser", body.user);
      })
      .catch(err => {
        if (err.status !== 401) {
          console.log(JSON.stringify(err));
        }
        // if it fails, do nothing
      })
      .finally(() => commit("setUserRequested", true));
  },
  setUser: ({ commit }, payload) => {
    commit("setUser", payload);
  },
  setFirstRoute: ({ commit }, payload) => {
    commit("setFirstRoute", payload);
  },
  setjwtToken: ({ commit }, payload) => {
    commit("setjwtToken", payload);
  },
  setVersion: ({ commit, state }, payload) => {
    if (!state.version || state.version === payload) {
      commit("setVersion", payload);
    } else {
      commit("setVersionChanged");
    }
  },
  logout: ({
    commit
    // dispatch
  }) => {
    commit("setUser", null);
    commit("setjwtToken", null);
  },
  setViewport: ({ commit }, payload) => {
    commit("setViewport", payload);
  }
};

export default {
  state: _state,
  getters,
  mutations,
  actions,
  namespaced: true
};
