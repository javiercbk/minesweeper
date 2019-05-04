<template>
  <span
    class="f-grow board-tile minesweeper-clickable"
    :class="{'revealed': isRevealed}"
    @click.left="click"
    @click.right.prevent="rightClick"
    @mousedown.right="mouseRightDown"
    @mouseup.right="mouseRightUp"
    @mousedown.left="mouseDown"
    @mouseup.left="mouseUp"
  >
    <font-awesome-icon v-if="icon !== null" :icon="icon"/>
    <span
      v-if="isRevealed && !isMine"
      class="w-100 h-100 text-bold"
      :class="mineProximityClass"
    >{{mineProximity}}</span>
    <span v-else class="w-100 h-100" :class="mineProximityClass">&nbsp;</span>
  </span>
</template>

<style lang="scss" scoped src="./ms-board-tile.scss"></style>
<script src="./ms-board-tile.js"></script>