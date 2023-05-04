<script lang="ts">
  import { createEventDispatcher } from "svelte";

  import { EventType, ToTitleCase, type SearchResult } from "./DataInterface";
  import { appService } from "./DataService";

  // Event dispatcher for all broadcasts to parent
  const dispatch = createEventDispatcher();

  export let searchTags: (searchInput: string) => Promise<SearchResult[]>;
  // export let addTag: (tagName: string) => Promise<boolean>;
  // export let removeTag: (tagName: string) => Promise<boolean>;

  export let tags: string[] = [];

  // The results that are returned by searching
  let results: SearchResult[] = [];

  // The current search input text
  let searchInput: string = "";

  let displayAddFrame: boolean = false;

  appService.appEvents.subscribe((value) => {
    if (value.type == EventType.Cancel) {
      results = [];
      searchInput = "";
      displayAddFrame = false;
    }
  });

  // Sends the search event to the parent
  function doSearch(event) {
    if (event.key == "Escape") {
      results = [];
      searchInput = "";
      displayAddFrame = false;
    } else if (event.key == "Enter") {
      onClick(searchInput.toLowerCase());
    } else {
      console.log("search: " + searchInput);
      searchTags(searchInput.toLowerCase()).then((searchResults) => {
        results = searchResults;
      });
    }
  }

  // Sets bold highlighting on the text input string
  function getHighlight(input: string, count: number): string {
    let result = "";
    let pieces = input.toLowerCase().split(searchInput.toLowerCase());

    for (let i = 0; i < pieces.length; i++) {
      result = result + pieces[i];

      if (i != pieces.length - 1) {
        // result = result + "<b>" + searchInput + "</b>";
        result = result + "" + searchInput + "";
      }
    }

    if (count) {
      result = result + " (" + count + ")";
    }

    return ToTitleCase(result);
  }

  function onClick(tagName: string) {
    searchInput = "";

    displayAddFrame = false;

    if (!tags.includes(tagName.toLowerCase())) {
      let tempTags = tags;
      tempTags.push(tagName.toLowerCase().replaceAll(",", ""));
      tags = tempTags;
    }

    // dispatch("addTag", {
    //   name: tagName.toLowerCase(),
    // });

    results = [];
  }

  function removeTag(tagName: string) {
    if (tags.includes(tagName.toLowerCase())) {
      let tempTags = tags;
      tempTags = tempTags.filter((e) => e !== tagName.toLowerCase());
      tags = tempTags;
    }
  }

  function addTagClick() {
    displayAddFrame = !displayAddFrame;
  }

  function observeClicks() {}
</script>

<div class="container">
  <div class="tags_header">
    <span class="add_label">Tags:</span>
    <span class="tags_list">
      {#each tags as tag}
        <span class="tag">
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <span
            class="delete_tag_button"
            on:click|stopPropagation={() => removeTag(tag)}>x</span
          >{ToTitleCase(tag)}</span
        >
      {/each}
    </span>
  </div>
  {#if displayAddFrame}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="add_box" on:click|stopPropagation={observeClicks}>
      <div class="search_box">
        <div class="input_box">
          <input
            class="input"
            bind:value={searchInput}
            on:keyup|stopPropagation={doSearch}
            placeholder="Add tags..."
            autofocus
          />
        </div>
      </div>
      {#if results.length > 0}
        <div class="results_box">
          <div class="arrow" />
          <div class="results_list">
            <div class="results_inner_list">
              {#each results as res, i}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <div
                  class="result"
                  on:click={() => {
                    onClick(res.title);
                  }}
                >
                  {@html getHighlight(res.title, res.count)}
                </div>
              {/each}
            </div>
          </div>
        </div>
      {/if}
    </div>
  {:else}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="add_button" on:click|stopPropagation={addTagClick}>+ Add</div>
  {/if}
</div>

<style>
  .container {
    margin-left: 8px;

    font-size: 15px;
    font-weight: 500;
    color: gray;
  }

  .search_box {
    display: flex;
    width: 240px;
    background: rgba(250, 250, 250, 1);
    border-color: rgba(0, 0, 0, 0.15) !important;
    border: 1px solid;
    /* border-radius: 20px; */
    margin: 12px 5px 10px 0px;
    font-family: sohne, "Helvetica Neue", Helvetica, Arial, sans-serif;
  }

  .input_box {
    margin: 0px 12px 0px 12px;
    padding: 4px 0px 0px 0px;
    color: darkgray;
    display: flex;
  }

  .add_button {
    margin-top: 10px;
    user-select: none;
    cursor: pointer;
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
    line-height: 14px;
    font-size: 14px;
  }

  .results_box {
    /* position: absolute; */
  }

  .results_list {
    position: relative;
    left: 0px;
    top: -10px;
    max-height: 200px;
    width: 252px;
    overflow-y: auto;
    overflow-x: auto;
    border-radius: 3px;
    background: rgb(255, 255, 255);
    box-shadow: rgba(0, 0, 0, 0.15) 0px 2px 10px 0px;
    border: 1px solid rgb(242, 242, 242);
    border-radius: 4px;
  }

  .results_inner_list {
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
    left: 32px;
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

  .tags_list {
    margin-top: 14px;
  }

  .tag {
    margin-right: 6px;
    background-color: rgb(231, 231, 231);
    border-radius: 25px;
    padding: 4px 12px 4px 4px;
    font-size: 14px;
    color: gray;
    user-select: none;
    cursor: pointer;
    /* text-transform: capitalize;
    display: inline-block; */
  }

  .delete_tag_button {
    background-color: darkgray;
    color: white;
    border-radius: 50px;
    padding: 0px 4px 0px 4px;
    margin-right: 6px;
  }
</style>
