<script lang="ts">
  export let onEnter: () => void;
  export let onLeave: () => void;

  function onMouseEnter() {
    if (onEnter) onEnter();
  }

  function onMouseLeave() {
    if (onLeave) onLeave();
  }
</script>

<div
  class="user_popup"
  on:mouseenter={onMouseEnter}
  on:mouseleave={onMouseLeave}
>
  <div class="user_popup_inner">hello world!</div>
</div>

<style>
  .user_popup {
    position: absolute;
    z-index: 999;
    width: 250px;
    height: 200px;
  }

  .user_popup_inner {
    background-color: red;
    position: absolute;
    top: 24px;
    width: 100%;
    height: 179px;
  }
</style>
