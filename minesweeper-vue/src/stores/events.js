const _state = {};

const getters = {};

const mutations = {};

const actions = {
  fireGlobalEvent: ({ dispatch }, payload) => {
    dispatch("notifications/clearScopedNotifications", payload, { root: true });
  }
};

export default {
  state: _state,
  getters,
  mutations,
  actions,
  namespaced: true
};
