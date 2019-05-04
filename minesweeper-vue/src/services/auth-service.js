import _ from "lodash";
import { jsonRequest, authenticatedRequest } from "./service-helpers";

export default class AuthService {
  getCurrentUser() {
    return authenticatedRequest(authHeaders =>
      fetch(
        "/api/players/current",
        _.merge(jsonRequest(), authHeaders, {
          method: "GET"
        })
      )
    );
  }

  authenticate(credentials) {
    return fetch(
      "/api/auth",
      _.merge(jsonRequest({ withBody: true }), {
        method: "POST",
        body: JSON.stringify(credentials)
      })
    );
  }
}
