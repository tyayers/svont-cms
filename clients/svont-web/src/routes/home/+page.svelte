<script lang="ts">
  import Header from "../../lib/Header.svelte";
  import PostCard from "../../lib/Post.card.svelte";
  import PostPopularWidget from "../../lib/Post.popular.svelte"

  import {appService} from "$lib/DataService"
  import type { PostOverviewCollection } from "$lib/DataInterface";

  export let data;

</script>

<div>
  <Header />

  <div class="content">
    {#if data && data.posts}
      <div class="container">
        <div class="panel_left">
          <div class="pannel_left_inner">
            {#each Object.values(data.posts).reverse() as post}
              <div>
                <PostCard {post} />
              </div>
            {/each}
            <div class="pannel_left_footer">
              
            </div>
          </div>
        </div>
        <div class="panel_right">
          <div class="widget1">
            <PostPopularWidget posts={data.popular} />
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .content {
    max-width: 1336px;
    text-align: left;
    margin: auto;
    height: 100%;
    height: calc(100vh - 58px);
    overflow-y: hidden;
  }

  .container {
    display: flex;
    justify-content: space-evenly;
    flex-direction: row;
    height: 100%;
  }

  .panel_left {
    width: 100%;
    flex: 1 1 auto;
    justify-content: center;
    display: inline-flex;
    height: 100%;
    overflow-y: scroll;
  }

  .pannel_left_inner {
    max-width: 728px;
    width: 100%;
    padding-bottom: 54px;
  }

  .pannel_left_footer {
    height: 24px;
  }

  .panel_right {
    min-height: 100vh;
    border-left: 1px solid rgba(242, 242, 242, 1);
    padding-left: 32px;
    min-width: 420px;
  }

  .widget1 {
    height: 100%;
    width: 100%;
  }

  @media (max-width: 903.98px) {
    .cc {
      min-width: 0;
    }

    .sb {
      min-width: 0;
    }

    .panel_right {
      display: none;
    }
  }

  @media (min-width: 904px) and (max-width: 1079.98px) {
    .cc {
      max-width: 680px;
    }

    .sb {
      max-width: 352px;
      min-width: 310px;
    }
  }

  @media (min-width: 1080px) {
    .sb {
      max-width: 352px;
      min-width: 352px;
      padding-right: 24px;
    }
  }

  @media (min-width: 1080px) {
    .sb {
      padding-left: clamp(24px, 24px + 100vw - 1080px, 40px);
    }
  }
</style>
