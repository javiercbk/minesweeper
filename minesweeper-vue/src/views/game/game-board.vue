<template>
  <div class="container-fluid p-0">
    <div class="container">
      <ms-game-board :board="game.board" @operation="onOperation" :state="gameStatus"/>
    </div>
  </div>
</template>

<script src="./game-board.js"></script>
