<template>
  <form @submit.prevent="submit" novalidate>
    <div class="form-group">
      <label for="login-name">Username</label>
      <input
        type="text"
        class="form-control"
        :class="{'is-valid': usernameValid, 'is-invalid': usernameInvalidDirty}"
        id="login-username"
        v-model="$v.username.$model"
        :disabled="disabled"
        placeholder="Enter username"
      >
      <div
        v-transparent="!usernameInvalidDirty"
        class="invalid-feedback d-block"
      >Username is mandatory</div>
    </div>
    <div class="form-group">
      <label for="login-password">Password</label>
      <input
        type="password"
        class="form-control"
        :class="{'is-valid': passwordValid, 'is-invalid': passwordInvalidDirty}"
        id="login-password"
        v-model="$v.password.$model"
        :disabled="disabled"
        placeholder="Password"
      >
      <div
        v-transparent="!passwordInvalidDirty"
        class="invalid-feedback d-block"
      >Password is mandatory</div>
    </div>
    <button type="submit" class="btn btn-primary w-100" :disabled="$v.invalid || disabled">Submit</button>
  </form>
</template>

<script src="./ms-user-form.js"></script>