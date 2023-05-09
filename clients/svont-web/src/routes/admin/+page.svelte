<script lang="ts">
  import { HeaderButton } from "../../lib/DataInterface";

  import Header from "../../lib/Header.svelte";

  export let data;
</script>

<div class="newpost">
  <Header actionType={HeaderButton.NewPost} />

  <div class="full_box">
    {#if data && data.metadata}
      <span>Total Posts: </span><span>{data.metadata.postCount}</span>
    {/if}
  </div>
</div>

<style>
  .full_box {
    width: 100vw;
    height: 100vh;
    overflow-x: auto;
  }
</style>
