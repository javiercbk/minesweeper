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
      username: "",
      password: ""
    };
  },
  validations: {
    username: {
      required,
      minLength: minLength(1)
    },
    password: {
      required,
      minLength: minLength(1)
    }
  },
  computed: {
    usernameValid() {
      return !this.$v.username.$invalid;
    },
    usernameInvalidDirty() {
      return this.$v.username.$invalid && this.$v.username.$dirty;
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
          username: this.username,
          password: this.password
        });
      }
    }
  }
};
