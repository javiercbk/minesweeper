<template>
  <div id="app" @click="onGlobalEvent($event, 'home')" @keydown="onGlobalEvent($event, 'home')">
    <div class="notification-manager">
      <notification-manager/>
    </div>
    <router-view/>
  </div>
</template>

<script src="./minesweeper-app.js"></script>
<style lang="scss" scoped="true" src="styles/minesweeper-app.scss"></style>
<style lang="scss" src="styles/minesweeper-app_global.scss"></style>

