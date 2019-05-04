import _ from "lodash";
import Promise from "bluebird";
import store from "@/store";

export const jsonRequest = function(options) {
  const fetchOptions = {
    headers: {
      Accept: "application/json"
    }
  };
  if (_.get(options, "withBody")) {
    fetchOptions.headers["Content-Type"] = "application/json";
  }
  return fetchOptions;
};

export const getAuthHeaders = function() {
  const token = store.getters["session/jwtToken"];
  if (_.isNil(token) || token === "") {
    return null;
  }
  return {
    headers: {
      authorization: `Bearer ${token}`
    }
  };
};

export const authenticatedRequest = function(fetchFactory) {
  const authHeaders = getAuthHeaders();
  if (authHeaders) {
    return fetchFactory(authHeaders);
  } else {
    return buildUnauthorizedResponse();
  }
};

export const buildUnauthorizedResponse = function() {
  return Promise.reject({
    ok: false,
    status: 401
  });
};

export const getCurrentUser = function() {
  return store.getters["session/user"];
};

export const getCurrentUserId = function() {
  const user = getCurrentUser();
  if (!_.isNil(user)) {
    return encodeURIComponent(user._id);
  }
  return null;
};
