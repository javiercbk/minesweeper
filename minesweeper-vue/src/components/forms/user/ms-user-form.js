import { required, minLength } from "vuelidate/lib/validators";

export default {
  props: {
    disabled: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      name: "",
      password: ""
    };
  },
  validations: {
    name: {
      required,
      minLength: minLength(1)
    },
    password: {
      required,
      minLength: minLength(1)
    }
  },
  computed: {
    nameValid() {
      return !this.$v.name.$invalid;
    },
    nameInvalidDirty() {
      return this.$v.name.$invalid && this.$v.name.$dirty;
    },
    passwordValid() {
      return !this.$v.password.$invalid;
    },
    passwordInvalidDirty() {
      return this.$v.password.$invalid && this.$v.password.$dirty;
    }
  },
  methods: {
    submit() {
      if (!this.$v.$invalid && !this.disabled) {
        this.$emit("user-submit", {
          name: this.name,
          password: this.password
        });
      }
    }
  }
};
