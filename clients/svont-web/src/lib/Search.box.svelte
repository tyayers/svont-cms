<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { type SearchResult, EventType } from "./DataInterface";
  import {appService} from "./DataService";

  export let search: (searchInput: string) => Promise<SearchResult[]>;

  // Event dispatcher for all broadcasts to parent
  const dispatch = createEventDispatcher();

  // The results that are returned by searching
  let results: SearchResult[] = [] as SearchResult[];

  // The current search input text
  let searchInput: string = "";

  appService.appEvents.subscribe((value) => {
    if (value.type == EventType.Cancel )
      results = [];
  });

  // Sends the search event to the parent
  function doSearch() {
    console.log("search: " + searchInput);
    search(searchInput).then((searchResults: SearchResult[]) => {
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

    dispatch("click", {
      id: id,
      title: title,
    });

    results = [];
  }

  function onFocus() {
    console.log("onfocus");
  }

  function onBlur() {
    console.log("onblur");
  }
</script>

<div class="bg">
  <div class="search">
    <div class="box">
      <svg class="icon" width="24" height="24" viewBox="0 0 24 24" fill="none"
        ><path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M4.1 11.06a6.95 6.95 0 1 1 13.9 0 6.95 6.95 0 0 1-13.9 0zm6.94-8.05a8.05 8.05 0 1 0 5.13 14.26l3.75 3.75a.56.56 0 1 0 .8-.79l-3.74-3.73A8.05 8.05 0 0 0 11.04 3v.01z"
          fill="currentColor"
        /></svg
      >
      <input class="input" bind:value={searchInput} on:keyup|stopPropagation={doSearch} on:click|stopPropagation={doSearch} on:focus={doSearch} />
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
              on:keydown={() => {
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
    border-radius: 20px;
    margin: 10px 5px 10px 16px;
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
