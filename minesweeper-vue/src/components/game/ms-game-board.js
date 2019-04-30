import {
  faSmile,
  faSurprise,
  faDizzy,
  faLaughSquint
} from "@fortawesome/free-solid-svg-icons";

export const OP_REVEAL = 1;
export const OP_MARK = 2;

export default {
  props: {
    board: {
      type: Array,
      default: () => []
    },
    state: {
      type: String
    }
  },
  data() {
    return {
      mouseDown: false
    };
  },
  computed: {
    faceIcon() {
      if (this.state === "won") {
        return faLaughSquint;
      } else if (this.state === "lost") {
        return faDizzy;
      } else if (this.mouseDown) {
        return faSurprise;
      }
      return faSmile;
    }
  },
  methods: {
    onMouseDown() {
      this.mouseDown = true;
    },
    onMouseUp() {
      this.mouseFalse = true;
    },
    sendMarkCommand(data) {
      this._sendCommand(data.row, data.col, OP_MARK);
    },
    sendRevealCommand(data) {
      this._sendCommand(data.row, data.col, OP_REVEAL);
    },
    _sendCommand(row, col, op) {
      this.$emit("command", {
        row: row,
        col: col,
        op: op
      });
    }
  }
};
