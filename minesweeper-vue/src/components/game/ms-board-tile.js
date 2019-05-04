import { faCog, faQuestion, faFlag } from "@fortawesome/free-solid-svg-icons";

export default {
  props: {
    row: {
      type: Number,
      required: true
    },
    col: {
      type: Number,
      required: true
    },
    mineProximity: {
      type: [Number, String]
    }
  },
  computed: {
    icon() {
      if (this.mineProximity === 9) {
        return faCog;
      } else if (this.mineProximity === "?") {
        return faQuestion;
      } else if (this.mineProximity === "!") {
        return faFlag;
      }
      return null;
    },
    isMine() {
      return this.mineProximity === 9;
    },
    isRevealed() {
      return typeof this.mineProximity === "number";
    },
    isMarked() {
      return typeof this.mineProximity === "string";
    },
    mineProximityClass() {
      if (this.mineProximity === null) {
        return "mine-unrevealed";
      } else if (this.mineProximity < 9) {
        return `mine-${this.mineProximity}`;
      }
    }
  },
  methods: {
    click() {
      if (!this.isRevealed && !this.isMarked) {
        this._emitMouseEvent("reveal");
      }
    },
    rightClick() {
      if (!this.isRevealed) {
        this._emitMouseEvent("mark");
      }
    },
    mouseDown() {
      if (!this.isRevealed && !this.isMarked) {
        this._emitMouseEvent("tile-mouse-down");
      }
    },
    mouseUp() {
      if (!this.isRevealed && !this.isMarked) {
        this._emitMouseEvent("tile-mouse-up");
      }
    },
    mouseRightUp() {
      if (!this.isRevealed) {
        this._emitMouseEvent("tile-mouse-right-up");
      }
    },
    mouseRightDown() {
      if (!this.isRevealed) {
        this._emitMouseEvent("tile-mouse-down-up");
      }
    },
    _emitMouseEvent(evt) {
      this.$emit(evt, {
        row: this.row,
        col: this.col
      });
    }
  }
};
