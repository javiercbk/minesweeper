<template>
  <form @submit.prevent="submit" novalidate>
    <div class="form-group">
      <label for="login-name">Username</label>
      <input
        type="text"
        class="form-control"
        id="login-username"
        v-model="$v.username.$model"
        aria-describedby="login-username-help"
        placeholder="Enter username"
      >
      <small id="login-username-help" class="form-text text-muted">The username that you've chosen</small>
    </div>
    <div class="form-group">
      <label for="login-password">Password</label>
      <input
        type="password"
        class="form-control"
        id="login-password"
        v-model="$v.password.$model"
        placeholder="Password"
      >
    </div>
    <button type="submit" class="btn btn-primary" :disabled="$v.invalid">Submit</button>
  </form>
</template>

<script src="./ms-user-form.js"></script>