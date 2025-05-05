<template>
  <div id="app">
    <nav class="navbar navbar-expand-lg navbar-light bg-light mb-4">
      <div class="container">
        <a class="navbar-brand" href="#">Social Media App</a>
        <div>
          <!-- <router-link to="/" class="btn btn-outline-primary me-2">Sign Up</router-link>
          <router-link to="/login" class="btn btn-outline-success">Log In</router-link> -->
        </div>
      </div>
    </nav>

    <div class="container">
      <router-view />
    </div>
  </div>
</template>

<script lang="ts" setup>
// empty for now
</script>
