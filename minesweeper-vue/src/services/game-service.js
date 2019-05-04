import _ from "lodash";
import { jsonRequest, authenticatedRequest } from "./service-helpers";

export default class GameService {
  find() {
    return authenticatedRequest(authHeaders =>
      fetch(
        "/api/games",
        _.merge(jsonRequest(), authHeaders, {
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
        _.merge({}, jsonRequest(), authHeaders, {
          method: "GET"
        })
      )
    );
  }

  create(game) {
    return authenticatedRequest(authHeaders =>
      fetch(
        "/api/games",
        _.merge(
          {},
          jsonRequest({ withBody: true }),
          {
            method: "POST",
            body: JSON.stringify(game)
          },
          authHeaders
        )
      )
    );
  }

  patch(operation) {
    return authenticatedRequest(authHeaders =>
      fetch(
        `/api/games/${encodeURIComponent(operation.gameId)}`,
        _.merge(
          {},
          jsonRequest({ withBody: true }),
          {
            method: "PATCH",
            body: JSON.stringify(operation)
          },
          authHeaders
        )
      )
    );
  }
}
