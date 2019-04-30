import _ from "lodash";
import { required, minValue } from "vuelidate/lib/validators";

export default {
  props: {
    disabled: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      rows: 6,
      cols: 6,
      mines: 10,
      private: false
    };
  },
  validations: {
    rows: {
      required,
      minValue: minValue(2)
    },
    cols: {
      required,
      minValue: minValue(2)
    },
    mines: {
      required,
      mineOverflow: (value, vm) => {
        if (
          typeof value === "number" &&
          typeof vm.rows === "number" &&
          typeof vm.cols === "number"
        ) {
          const boardSize = vm.rows * vm.cols;
          return value < boardSize;
        }
        return false;
      }
    }
  },

  methods: {
    submit() {
      if (!this.$v.$invalid) {
        this.$emit(
          "game-submit",
          _.pick(this, ["rows", "cols", "mines", "private"])
        );
      }
    }
  }
};
