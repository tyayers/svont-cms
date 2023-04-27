<script lang="ts">
  import type { SearchResult } from "./DataInterface";

  export let search: (searchInput: string) => Promise<SearchResult[]>;;

  // The results that are returned by searching
  let results: SearchResult[] = [];

  // The current search input text
  let searchInput: string = "";

  // Sends the search event to the parent
  function doSearch() {
    console.log("search: " + searchInput);
    search(searchInput).then((searchResults) => {
      results = searchResults;
    });
  }

  // Sets bold highlighting on the text input string
  function getHighlight(input: string): string {
    let result = "";
    let pieces = input.split(searchInput);

    for (let i = 0; i < pieces.length; i++) {
      result = result + pieces[i];
      if (i != pieces.length - 1)
        result = result + "<b>" + searchInput + "</b>";
    }

    return result;
  }

  function onClick(id: string, title: string) {
    searchInput = "";

    // dispatch("click", {
    //   id: id,
    //   title: title,
    // });

    results = [];
  }

</script>

<div class="container">
  <div class="search">
    <div class="box">
      <input class="input" bind:value={searchInput} on:keyup={doSearch} placeholder="Add tags..." />
    </div>
  </div>
  {#if results.length > 0}
    <div class="resultsPanel">
      <div class="arrow" />
      <div class="results">
        <div class="panel">
          {#each results as res, i}
            <div
              class="result"
              on:click={() => {
                onClick(res.id, res.title);
              }}
            >
              {@html getHighlight(res.title)}
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .search {
    display: flex;
    width: 240px;
    background: rgba(250, 250, 250, 1);
    border-color: rgba(0,0,0,.15)!important;
    border: 1px solid;
    /* border-radius: 20px; */
    margin: 0px 5px 10px 8px;
    font-family: sohne, "Helvetica Neue", Helvetica, Arial, sans-serif;
  }

  .box {
    margin: 0px 12px 0px 12px;
    padding: 4px 0px 0px 0px;
    color: darkgray;
    display: flex;
  }

  .icon {
    position: relative;
    top: 3px;
  }

  .input {
    position: relative;
    top: -4px;
    left: 3px;
    padding: 10px 20px 10px 0;
    background-color: transparent;
    outline: none;
    border: none;
    color: rgba(41, 41, 41, 1);
    line-height: 20px;
    font-size: 14px;
  }

  .resultsPanel {
    position: absolute;
  }

  .results {
    position: relative;
    left: 20px;
    top: -10px;
    max-height: 200px;
    width: 316px;
    overflow-y: scroll;
    border-radius: 3px;
    background: rgb(255, 255, 255);
    box-shadow: rgba(0, 0, 0, 0.15) 0px 2px 10px 0px;
    border: 1px solid rgb(242, 242, 242);
    border-radius: 4px;
  }

  .panel {
    position: relative;
    background: rgb(255, 255, 255);
    width: 100%;
    height: 100%;
    z-index: 2;
    padding-top: 20px;
    padding-bottom: 20px;
    font-family: sohne, "Helvetica Neue", Helvetica, Arial, sans-serif;
  }

  .arrow {
    position: relative;
    top: -22px;
    left: 52px;
    z-index: 1;
    border: 1px solid rgb(242, 242, 242);
    box-shadow: rgba(0, 0, 0, 0.15) -1px -1px 1px -1px;
    transform: rotate(45deg) translate(16px, 16px);
    background: rgb(255, 255, 255);
    height: 14px;
    width: 14px;
    display: block;
    content: "";
    pointer-events: none;
  }

  .result {
    padding-top: 5px;
    padding-bottom: 5px;
    padding-left: 10px;
    border-bottom: 1px dashed rgb(242, 242, 242);
    cursor: pointer;
  }
</style>
