import { jsonRequest, authenticatedRequest } from "./service-helpers";

export default class GameService {
  find() {
    return authenticatedRequest(authHeaders =>
      fetch(
        "/api/games",
        Object.assign(jsonRequest(), authHeaders, {
          method: "GET"
        })
      )
    );
  }

  retrieve(gameId) {
    const escapedGameId = encodeURIComponent(gameId);
    return authenticatedRequest(authHeaders =>
      fetch(
        `/api/games/${escapedGameId}`,
        Object.assign(jsonRequest(), authHeaders, {
          method: "GET"
        })
      )
    );
  }

  create(game) {
    return authenticatedRequest(authHeaders =>
      fetch(
        "/api/games",
        Object.assign(jsonRequest({ withBody: true }), authHeaders, {
          method: "POST",
          body: JSON.stringify(game)
        })
      )
    );
  }
}
