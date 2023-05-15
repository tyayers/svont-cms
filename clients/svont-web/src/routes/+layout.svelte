<script>
  import "../assets/app.css";
  import { EventType } from "$lib/DataInterface";
  import { appService } from "$lib/DataService";

  function sendCancel() {
    appService.SendEvent(EventType.Cancel);
  }

  // This works at runtime
  // document.documentElement.style.setProperty(
  //   "--main-background-color",
  //   "brown"
  // );
</script>

<div on:click={sendCancel} on:keyup={sendCancel}>
  <slot />
</div>

<style>
  :root {
    --main-background-color: rgba(255, 255, 255, 1);

    --header-height: 57px;
    --header-padding: 0 24px;
    --header-bottom-border: solid 1px rgba(242, 242, 242, 1);
  }
</style>
