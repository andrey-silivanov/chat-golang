<template>

  <div>

      <main-page v-if="isUserAuthorized"/>
      <login-page v-else/>

    <!-- page-body-wrapper ends -->
  </div>
</template>

<script>

import LoginPage from "@/page/LoginPage";
import MainPage from "@/page/MainPage";

export default {
  name: 'App',
  data:() => ({
    userId: null
  }),
  computed: {
    isUserAuthorized() {
      return this.userId != null
    }
  },
  components: {
    MainPage,
    LoginPage
  }
}
</script>

<style lang="scss">
  @import "src/assets/scss/style.scss";
</style>
