<template>
  <div class="container-fluid board">
    <div class="row board-status">
      <font-awesome-icon :icon="faceIcon"/>
    </div>
    <div class="row p-0 board-row" v-for="(r, row) in board" :key="row">
      <div v-for="(c, col) in board[row]" :key="col" class="board-col">
        <ms-board-tile
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
  </div>
</template>

<script src="./ms-game-board.js"></script>
