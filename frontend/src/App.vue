<template>

  <div class="container-fluid">

      <main-block/>

    <!-- page-body-wrapper ends -->
  </div>
</template>

<script>

import MainBlock from "@/components/MainBlock";

export default {
  name: 'App',
  components: {
    MainBlock
  }
}
</script>

<style lang="scss">
  @import "src/assets/scss/style.scss";
</style>
