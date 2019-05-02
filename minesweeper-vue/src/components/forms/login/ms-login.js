import { mapActions } from "vuex";
import MsUserForm from "@/components/forms/user/ms-user-form.vue";
import AuthService from "@/services/auth-service";

const authService = new AuthService();

export default {
  components: {
    MsUserForm
  },
  methods: {
    ...mapActions("session", ["setUser", "setJWTToken"]),
    ...mapActions("notifications", ["pushNotification"]),
    handleUserLogin(credentials) {
      authService
        .authenticate(credentials)
        .then(res => res.json())
        .then(body => {
          // store the user info
          this.setJWTToken(body.data.token);
          this.setUser(body.data.user);
        })
        .catch(res => {
          if (res.status === 401) {
            this.pushNotification({
              title: "Bad credentials",
              message: "Username or password are incorrect",
              variant: "warning"
            });
          } else {
            this.pushNotification({
              title: "Error",
              message: "could not login",
              variant: "danger"
            });
          }
        });
    }
  }
};
