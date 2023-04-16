<script>
  import '../assets/app.css'
  import {EventType} from "$lib/DataInterface"
  import {appService} from "$lib/DataService"

  function sendCancel() {
    appService.SendEvent(EventType.Cancel);
  }
</script>

<div on:click={sendCancel} on:keyup={sendCancel}>
  <slot />
</div>

