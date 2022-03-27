<template>
  <div>
  <div class="clearfix content-wrapper">
    <div class="row">

    <people-list :show="true"/>
    <chat-component/>
    </div>

  </div> <!-- end container -->

  </div>


</template>

<script>

import PeopleList from "@/components/PeopleList";
import ChatComponent from "@/components/ChatComponent";

export default {
  name: "MainBlock",
  components: {
    PeopleList,
    ChatComponent
  }
}
</script>

<style scoped lang="scss">


</style>