import { required, minLength } from "vuelidate/lib/validators";

export default {
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
  methods: {
    submit() {
      if (!this.$v.$invalid) {
        this.$emit("user-submit", {
          username: this.username,
          password: this.password
        });
      }
    }
  }
};
