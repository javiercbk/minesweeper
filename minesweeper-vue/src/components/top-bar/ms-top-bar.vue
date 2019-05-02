<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <router-link class="navbar-brand" :to="{name: 'index'}">Minesweeper</router-link>
    <div class="collapse navbar-collapse" id="navbarSupportedContent">
      <ul class="navbar-nav mr-auto">
        <li class="nav-item active" :class="{'active': isGameCreation}">
          <router-link class="navbar-brand" :to="{name: 'game-search-create'}">
            Games
            <span v-if="isGameCreation" class="sr-only">current</span>
          </router-link>
        </li>
        <li class="nav-item dropdown">
          <a
            @click.prevent
            class="nav-link dropdown-toggle"
            href="#"
            id="navbarDropdown"
            role="button"
            data-toggle="dropdown"
            aria-haspopup="true"
            aria-expanded="false"
          >{{userName}}</a>
          <div class="dropdown-menu" aria-labelledby="navbarDropdown">
            <a class="dropdown-item" @click.prevent="performLogout" href="#">logout</a>
          </div>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script src="./ms-top-bar.js"></script>
