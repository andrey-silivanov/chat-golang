<template>
  <div class="container-fluid">
    <div class="clearfix content-wrapper">
      <h3 style="margin-top: 50px;">User {{ _user }}</h3>
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
  name: "MainView",
  components: {
    PeopleList,
    ChatComponent
  },
  computed: {
    _user() {
      return this.$auth.user() || {};
    }
  },
  created() {
    console.log("TOKKKEN :", this.$auth.token())
   // this.getUser()
  },

  methods: {
    getUser() {
      this.$auth.fetch({}).then(res => {
        let user = res.data.body.user
        console.log(user)
        this.$auth.user(user);
      })
    }
  }
}
</script>

<style scoped lang="scss">


</style>