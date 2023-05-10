<script lang="ts">
  import { appService } from "$lib/DataService";
  import { HeaderButton } from "../../lib/DataInterface";

  import Header from "../../lib/Header.svelte";

  export let data;

  function doRefresh() {
    appService.DoRefresh().then((data) => {
      alert(JSON.stringify(data));
    });
  }

  function doPersist() {
    appService.DoRefresh().then((data) => {
      alert(JSON.stringify(data));
    });
  }
</script>

<div class="newpost">
  <Header actionType={HeaderButton.NewPost} />

  <div class="full_box">
    <div class="content_box">
      {#if data && data.metadata}
        <span>Total Posts: </span><span>{data.metadata.postCount}</span>
      {/if}
      <div>
        <button on:click={doRefresh}>Refresh Data</button>
      </div>
      <div>
        <button on:click={doPersist}>Persist Data</button>
      </div>
    </div>
  </div>
</div>

<style>
  .full_box {
    width: 100vw;
    height: 100vh;
    overflow-x: auto;
  }

  .content_box {
    padding-top: 46px;
    /* max-width: 1336px; */
    max-width: 736px;
    padding-left: 24px;
    padding-right: 24px;
    text-align: left;
    margin: auto;
  }
</style>
