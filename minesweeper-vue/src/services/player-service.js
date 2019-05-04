import _ from "lodash";
import { jsonRequest } from "./service-helpers";

export default class PlayerService {
  create(player) {
    return fetch(
      "/api/players",
      _.merge(jsonRequest({ withBody: true }), {
        method: "POST",
        body: JSON.stringify(player)
      })
    );
  }
}
