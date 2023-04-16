<script lang="ts">
  import { navigate } from "../App.svelte";
  import type { PostOverviewCollection } from "./DataInterface";

  import PostCard from "./Post.card.svelte";
  import type { AppUser } from "./DataInterface";

  export let posts: PostOverviewCollection;

  console.log(JSON.stringify(posts));

  function loadPost(url: string) {
    navigate(url);
  }
</script>

<div class="cf">
  {#if posts}
    <div class="fl">
      <div class="cc">
        <div>
          {#each Object.values(posts) as post}
            <div>
              <PostCard {post} />
            </div>
          {/each}
        </div>
      </div>
      <div class="sb">
        <div class="sc" />
      </div>
    </div>
  {/if}
</div>

<style>
  .cf {
    max-width: 1336px;
    text-align: left;
    margin: auto;
    height: 100%;
  }

  .fl {
    display: flex;
    justify-content: space-evenly;
    flex-direction: row;
  }

  .cc {
    max-width: 728px;
    flex: 1 1 auto;
    justify-content: center;
  }

  .sb {
    min-height: 100vh;
    border-left: 1px solid rgba(242, 242, 242, 1);
  }

  .sc {
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
