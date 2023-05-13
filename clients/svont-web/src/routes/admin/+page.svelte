<script lang="ts">
  import { appService } from "$lib/DataService";
  import { HeaderButton } from "../../lib/DataInterface";

  import Header from "../../lib/Header.svelte";

  export let data;

  let status: string = "";

  function doRefresh() {
    status = "Refreshing data...";
    appService.DoRefresh().then((result) => {
      status = "Refresh complete.";
      data.metadata = result;
    });
  }

  function doPersist() {
    status = "Persisting data...";
    appService.DoRefresh().then((result) => {
      status = "Persisting complete.";
      data.metadata = result;
    });
  }
</script>

<div class="container_box">
  <Header actionType={HeaderButton.NewPost} />

  <div class="full_box">
    <div class="content_box">
      <div class="controls_box">
        <button class="control" on:click={doRefresh}>Refresh Data</button>
        <button class="control" on:click={doPersist}>Persist Data</button>
        <span class="control">{status}</span>
      </div>
      {#if data && data.metadata}
        <div class="stats_box">
          <span class="stat">{data.metadata.postCount} Total Posts</span>
          <span class="stat">{data.metadata.draftCount} Drafts</span>
          <span class="stat">{data.metadata.deletedCount} Deleted</span>
          <span class="stat"
            >{data.metadata.postCount -
              data.metadata.draftCount -
              data.metadata.deletedCount} Live Posts</span
          >
        </div>
      {/if}
      <div />
      <div />
    </div>
  </div>
</div>

<style>
  .container_box {
    overflow-y: scroll;
    height: 100vh;
  }

  .full_box {
    /* width: 100vw; */
    /* height: 100vh; */
    /* overflow-x: auto; */
  }

  .content_box {
    padding-top: 46px;
    padding-bottom: 46px;
    /* max-width: 1336px; */
    max-width: 736px;
    padding-left: 24px;
    padding-right: 24px;
    text-align: left;
    margin: auto;
    display: flex;
    flex-wrap: wrap;
  }

  .controls_box {
    display: flex;
    width: 100%;
    align-content: space-between;
    align-items: stretch;
    justify-content: flex-start;
  }

  .control {
    margin-right: 6px;
  }

  .stats_box {
    display: flex;
    margin-top: 38px;
    width: 100%;
    align-items: stretch;
    align-content: space-around;
    justify-content: space-evenly;
  }

  .stat {
    /* margin-right: 18px; */
  }
</style>
