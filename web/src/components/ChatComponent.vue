<template>
  <div class="column md-67 lg-75 chat">

    <div class="chat-header clearfix">
      <div class="test">
        <button>Click</button>
      </div>
      <img src="https://s3-us-west-2.amazonaws.com/s.cdpn.io/195612/chat_avatar_01_green.jpg" alt="avatar" />

      <div class="chat-about">
        <div class="chat-with">Chat with Vincent Porter</div>
        <div class="chat-num-messages">already 1 902 messages</div>
      </div>
      <i class="fa fa-star"></i>
    </div> <!-- end chat-header -->

    <div class="chat-history">
      <ul>
        <li class="clearfix">
          <div class="message-data align-right">
            <span class="message-data-time" >10:10 AM, Today</span> &nbsp; &nbsp;
            <span class="message-data-name" >Olia</span> <i class="fa fa-circle me"></i>

          </div>
          <div class="message other-message float-right">
            Hi Vincent, how are you? How is the project coming along?
          </div>
        </li>

        <li>
          <div class="message-data">
            <span class="message-data-name"><i class="fa fa-circle online"></i> Vincent</span>
            <span class="message-data-time">10:12 AM, Today</span>
          </div>
          <div class="message my-message">
            Are we meeting today? Project has been already finished and I have results to show you.
          </div>
        </li>

        <li class="clearfix">
          <div class="message-data align-right">
            <span class="message-data-time" >10:14 AM, Today</span> &nbsp; &nbsp;
            <span class="message-data-name" >Olia</span> <i class="fa fa-circle me"></i>

          </div>
          <div class="message other-message float-right">
            Well I am not sure. The rest of the team is not here yet. Maybe in an hour or so? Have you faced any problems at the last phase of the project?
          </div>
        </li>

        <li>
          <div class="message-data">
            <span class="message-data-name"><i class="fa fa-circle online"></i> Vincent</span>
            <span class="message-data-time">10:20 AM, Today</span>
          </div>
          <div class="message my-message">
            Actually everything was fine. I'm very excited to show this to our team.
          </div>
        </li>

        <li class="clearfix">
          <div class="message-data align-right">
            <span class="message-data-time" >10:14 AM, Today</span> &nbsp; &nbsp;
            <span class="message-data-name" >Olia</span> <i class="fa fa-circle me"></i>

          </div>
          <div class="message other-message float-right">
            Well I am not sure. The rest of the team is not here yet. Maybe in an hour or so? Have you faced any problems at the last phase of the project?
          </div>
        </li>
        <li class="clearfix">
          <div class="message-data align-right">
            <span class="message-data-time" >10:14 AM, Today</span> &nbsp; &nbsp;
            <span class="message-data-name" >Olia</span> <i class="fa fa-circle me"></i>

          </div>
          <div class="message other-message float-right">
            Well I am not sure. The rest of the team is not here yet. Maybe in an hour or so? Have you faced any problems at the last phase of the project?
          </div>
        </li>
        <li class="clearfix">
          <div class="message-data align-right">
            <span class="message-data-time" >10:14 AM, Today</span> &nbsp; &nbsp;
            <span class="message-data-name" >Olia</span> <i class="fa fa-circle me"></i>

          </div>
          <div class="message other-message float-right">
            Well I am not sure. The rest of the team is not here yet. Maybe in an hour or so? Have you faced any problems at the last phase of the project?
          </div>
        </li>
        <li class="clearfix">
          <div class="message-data align-right">
            <span class="message-data-time" >10:14 AM, Today</span> &nbsp; &nbsp;
            <span class="message-data-name" >Olia</span> <i class="fa fa-circle me"></i>

          </div>
          <div class="message other-message float-right">
            Well I am not sure. The rest of the team is not here yet. Maybe in an hour or so? Have you faced any problems at the last phase of the project?
          </div>
        </li>
        <li class="clearfix">
          <div class="message-data align-right">
            <span class="message-data-time" >10:14 AM, Today</span> &nbsp; &nbsp;
            <span class="message-data-name" >Olia</span> <i class="fa fa-circle me"></i>

          </div>
          <div class="message other-message float-right">
            Well I am not sure. The rest of the team is not here yet. Maybe in an hour or so? Have you faced any problems at the last phase of the project?
          </div>
        </li>
        <li class="clearfix">
          <div class="message-data align-right">
            <span class="message-data-time" >10:14 AM, Today</span> &nbsp; &nbsp;
            <span class="message-data-name" >Olia</span> <i class="fa fa-circle me"></i>

          </div>
          <div class="message other-message float-right">
            Well I am not sure. The rest of the team is not here yet. Maybe in an hour or so? Have you faced any problems at the last phase of the project?
          </div>
        </li>
        <li>
          <div class="message-data">
            <span class="message-data-name"><i class="fa fa-circle online"></i> Vincent</span>
            <span class="message-data-time">10:31 AM, Today</span>
          </div>
          <i class="fa fa-circle online"></i>
          <i class="fa fa-circle online" style="color: #AED2A6"></i>
          <i class="fa fa-circle online" style="color:#DAE9DA"></i>
        </li>

      </ul>

    </div> <!-- end chat-history -->

    <div class="chat-message clearfix">
      <textarea name="message-to-send" id="message-to-send" placeholder ="Type your message" rows="3"></textarea>

      <i class="fa fa-file-o"></i> &nbsp;&nbsp;&nbsp;
      <i class="fa fa-file-image-o"></i>

      <button>Send</button>

    </div> <!-- end chat-message -->

  </div> <!-- end chat -->
</template>

<script>
export default {
  name: "ChatComponent",
  mounted() {
    this.ws()
  },
  methods: {
    ws() {

      let socket = new WebSocket("ws://localhost:3000/ws/22");

      socket.onopen = function() {
        console.log("[open] Соединение установлено");
        console.log("Отправляем данные на сервер");
        socket.send(
            JSON.stringify({
                  email: "test@gmail.com",
                  username: "Vue man",
                  message: "Hello from vue"
                }
            ));

      };

      socket.onmessage = function(event) {
        console.log(`[message] Данные получены с сервера: ${event.data}`);
      };

      socket.onclose = function(event) {
        if (event.wasClean) {
          console.log(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
        } else {
          // например, сервер убил процесс или сеть недоступна
          // обычно в этом случае event.code 1006
          console.log('[close] Соединение прервано');
        }
      };

      socket.onerror = function(error) {
        console.log(`[error] ${error.message}`);
      };
    }
  }
}
</script>

<style scoped>

</style>