<template>
  <div class="container-fluid">
  <div class="clearfix content-wrapper">
    <div class="row">

    <people-list/>
    <chat-component/>
    </div>

  </div> <!-- end container -->

  </div>


</template>

<script>

import PeopleList from "@/components/PeopleList";
import ChatComponent from "@/components/ChatComponent";

export default {
  name: "MainPage",
  components: {
    PeopleList,
    ChatComponent
  }
}
</script>

<style scoped lang="scss">


</style>