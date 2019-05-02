/* eslint-disable no-param-reassign */
import moment from "moment";

const isSameNotification = o => n =>
  n.title === o.title && n.message === o.message;

const findIndex = (n, func) => {
  const len = n.length;
  for (let i = 0; i < len; i++) {
    if (func(n[i])) {
      return i;
    }
  }
  return -1;
};

const timestampedId = () => {
  const unixTime = new Date().getTime();
  const randomNumber = Math.floor(Math.random() * 1000);
  return `${unixTime}-${randomNumber}`;
};

const _state = {
  notifications: []
};

const getters = {
  notifications: storeState => storeState.notifications
};

const mutations = {
  pushNotification: (storeState, newNotification) => {
    // avoid the global event handler to remove the notification
    // on arrival
    newNotification.id = timestampedId();
    if (typeof newNotification.sticky !== "boolean") {
      newNotification.sticky = false;
    }
    setTimeout(() => {
      storeState.notifications.push(newNotification);
    }, 0);
  },
  deleteNotification: (storeState, payload) => {
    const toRemoveIndex = findIndex(
      storeState.notifications,
      isSameNotification(payload)
    );
    if (toRemoveIndex >= 0) {
      storeState.notifications.splice(toRemoveIndex, 1);
    }
  },
  clearNotifications: storeState => {
    storeState.notifications = [];
  },
  clearScopedNotifications: storeState => {
    const notificationsClone = storeState.notifications.slice(0);
    for (let i = 0; i < notificationsClone.length; i++) {
      if (!notificationsClone[i].sticky && !notificationsClone[i]._id) {
        notificationsClone.splice(i, 1);
      }
    }
    storeState.notifications = notificationsClone;
  }
};

const actions = {
  pushNotification: ({ commit, state }, payload) => {
    const existingIndex = findIndex(
      state.notifications,
      isSameNotification(payload)
    );
    if (existingIndex === -1) {
      if (payload.dismissible === undefined) {
        payload.dismissible = true;
      }
      payload.createdAt = moment();
      commit("pushNotification", payload);
    }
  },
  deleteNotification: ({ commit }, payload) => {
    commit("deleteNotification", payload);
  },
  clearNotifications: ({ commit }) => {
    commit("clearNotifications");
  },
  clearScopedNotifications: ({ commit }) => {
    commit("clearScopedNotifications");
  }
};

export default {
  state: _state,
  getters,
  mutations,
  actions,
  namespaced: true
};
