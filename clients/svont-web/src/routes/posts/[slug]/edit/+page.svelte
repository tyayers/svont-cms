<script lang="ts">
  import { HeaderButton } from "../../../../lib/DataInterface";

  import Header from "../../../../lib/Header.svelte";
  import PostNew from "../../../../lib/Post.edit.svelte";
  import type {
    PostOverviewCollection,
    PostOverview,
    Post,
  } from "../../../../lib/DataInterface";

  export let data: Post = undefined;

  let form;

  function submitClick(event) {
    if (form) {
      form.submit();
    }
  }
</script>

<div class="newpost">
  <Header small={true} actionType={HeaderButton.Submit} on:submit={submitClick} />

  <PostNew bind:this={form} post={data} />
</div>

<style>
.newpost {
  width: 100vw;
  height: 100vh;
}
</style>

