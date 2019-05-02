<template>
  <span
    class="board-tile"
    :class="{'revealed': isRevealed}"
    @mousedown.right="mouseRightDown"
    @mouseup.right="mouseRightUp"
    @mousedown.left="mouseDown"
    @mouseup.left="mouseUp"
  >
    <font-awesome-icon v-if="icon !== null" :icon="icon"/>
    <span
      v-if="isRevealed && !isMine"
      class="text-bold"
      :class="mineProximityClass"
    >{{mineProximity}}</span>
  </span>
</template>

<style lang="scss" scoped src="ms-board-tile.scss"></style>
