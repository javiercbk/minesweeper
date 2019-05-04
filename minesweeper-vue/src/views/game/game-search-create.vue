<template>
  <div class="container">
    <div class="row">
      <h1>Games</h1>
      <hr>
    </div>
    <div class="row">
      <h2>Playable games</h2>
    </div>
    <div class="row" v-for="g in availableGames" :key="g.id">
      <div class="col-2">
        <span class="sr-only">id</span>
        {{g.id}}
      </div>
      <div class="col-6">
        <span class="sr-only">creator</span>
        {{g.creator.name}}
      </div>
      <div class="col-4">
        <button class="btn btn-primary" @click="playGame(g.id)">Play</button>
      </div>
    </div>
    <div class="row">
      <h2>Create new game</h2>
    </div>
    <div class="row">
      <ms-game-form @game-submit="onGameSubmit" :disabled="creatingGame"/>
    </div>
  </div>
</template>

<script src="./game-search-create.js"></script>
