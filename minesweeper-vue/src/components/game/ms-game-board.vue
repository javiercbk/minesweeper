<template>
  <div class="container-fluid board">
    <div class="row justify-content-center align-items-center board-status">
      <font-awesome-icon :icon="faceIcon" class="display-4"/>
    </div>
    <div class="w-100 h-100 d-flex flex-row" v-for="(r, row) in board" :key="row">
      <ms-board-tile
        v-for="(c, col) in board[row]"
        :key="col"
        :row="row"
        :col="col"
        :mineProximity="board[row][col]"
        @reveal="sendRevealCommand"
        @mark="sendMarkCommand"
        @tile-mouse-down="onMouseDown"
        @tile-mouse-up="onMouseUp"
      />
    </div>
  </div>
</template>

<script src="./ms-game-board.js"></script>
