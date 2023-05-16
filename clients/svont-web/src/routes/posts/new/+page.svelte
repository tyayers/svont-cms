<script lang="ts">
  import { HeaderButton } from "../../../lib/DataInterface";

  import Header from "../../../lib/Header.svelte";
  import PostNew from "../../../lib/Post.edit.svelte";

  let form;

  let statusText: string = "";

  function submitClick(event) {
    if (form) {
      form.submit();
    }
  }

  function statusUpdate(newText: string) {
    statusText = newText;
  }
</script>

<div class="newpost">
  <Header
    small={true}
    showSearch={false}
    {statusText}
    actionType={HeaderButton.Submit}
    on:submit={submitClick}
  />

  <PostNew bind:this={form} {statusUpdate} />
</div>

<style>
  .newpost {
    width: 100vw;
    height: 100vh;
    overflow-x: auto;
  }
</style>
