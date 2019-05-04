<template>
  <form @submit.prevent="submit" novalidate>
    <div class="form-group">
      <label for="game-rows">Rows</label>
      <input
        type="number"
        min="1"
        max="100"
        class="form-control"
        id="game-rows"
        v-model="$v.rows.$model"
        :disabled="disabled"
        aria-describedby="game-rows-help"
      >
      <small id="game-rows-help" class="form-text text-muted">Maximum 100</small>
    </div>
    <div class="form-group">
      <label for="game-cols">Columns</label>
      <input
        type="number"
        min="1"
        max="100"
        class="form-control"
        id="game-cols"
        v-model="$v.cols.$model"
        :disabled="disabled"
        aria-describedby="game-cols-help"
      >
      <small id="game-cols-help" class="form-text text-muted">Maximum 100</small>
    </div>
    <div class="form-group">
      <label for="game-mines">Mines</label>
      <input
        type="number"
        min="1"
        max="100"
        class="form-control"
        id="game-mines"
        v-model="$v.mines.$model"
        :disabled="disabled"
        aria-describedby="game-mines-help"
      >
      <small
        id="game-cols-help"
        class="form-text text-muted"
      >The amount of mines must be less than the amount of rows and cols</small>
    </div>
    <div class="form-check">
      <input type="checkbox" class="form-check-input" id="game-private" v-model="$v.private.$model">
      <label class="form-check-label" for="game-private">Private game</label>
    </div>
    <button type="submit" class="btn btn-primary" :disabled="$v.invalid || disabled">Submit</button>
  </form>
</template>

<script src="./ms-game-form.js"></script>