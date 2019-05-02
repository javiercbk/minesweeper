<template>
  <div class="container-fluid h-100">
    <div class="row justify-content-center align-items-center h-100">
      <div class="col-12 col-md-4">
        <ms-user-form @user-submit="handleUserCreate" class="mb-4"/>
        <div class="w-100 text-center">
          <router-link :to="{name: 'player-create'}">Already have an account? Log in</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script src="./player-create.js"></script>