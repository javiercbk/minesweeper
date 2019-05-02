import { jsonRequest } from "./service-helpers";

export default class PlayerService {
  create(player) {
    return fetch(
      "/api/players",
      Object.assign(jsonRequest({ withBody: true }), {
        method: "POST",
        body: JSON.stringify(player)
      })
    );
  }
}
