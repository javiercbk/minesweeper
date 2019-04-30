import { jsonRequest, authenticatedRequest } from "./service-helpers";

export default class AuthService {
  getCurrentUser() {
    return authenticatedRequest(authHeaders =>
      fetch(
        "/api/players/current",
        Object.assign(jsonRequest(), authHeaders, {
          method: "GET"
        })
      )
    );
  }

  authenticate(credentials) {
    return fetch(
      "/api/auth",
      Object.assign(jsonRequest({ withBody: true }), {
        method: "POST",
        body: JSON.stringify(credentials)
      })
    );
  }
}
