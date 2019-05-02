<template>
  <div id="app" @click="onGlobalEvent($event, 'home')" @keydown="onGlobalEvent($event, 'home')">
    <div class="notification-manager">
      <notification-manager/>
    </div>
    <ms-top-bar v-if="user !== null"/>
    <div class="container-fluid p-0 h-100">
      <router-view/>
    </div>
  </div>
</template>

<script src="./minesweeper-app.js"></script>
<style lang="scss" scoped="true" src="./minesweeper-app.scss"></style>
<style lang="scss" src="./minesweeper-app_global.scss"></style>

