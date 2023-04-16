<script lang="ts">
  import { navigate } from "../App.svelte";
  import type { AppUser } from "./DataInterface";
  import Header from "./Header.svelte";
  import type { PostOverviewCollection, Post } from "./DataInterface";

  // The Post Id to view
  export let id: string;

  export let post: Post;

  console.log(JSON.stringify(post));

  // $: newId = (() => {
  //   // Load post data
  //   console.log("Loading content for post " + id);
  //   fetch(import.meta.env.VITE_CMS_SERVICE + "/posts/" + id)
  //     .then((response) => {
  //       return response.json();
  //     })
  //     .then((data: Post) => {
  //       post = data;
  //     });
  // })();

  // The Post content loaded from the server
  // let post: Post;

  function deletePost() {
    fetch(import.meta.env.VITE_CMS_SERVICE + "/posts/" + id, {
      method: "delete",
      headers: {
        Accept: "application/json",
      },
    }).then((response) => {
      navigate("/");
    });
  }
</script>

<Header />
<div class="post-container">
  <div class="post-flexbox">
    <div class="post-content">
      <!-- <div class="post-article">
        <div class="post-header">
          <Link to="/posts/{post.header.id}/edit">Edit</Link>
          <button
            class=""
            on:click={() => {
              deletePost();
            }}>Delete</button
          >
        </div>
        <h1>{post.header.title}</h1>
        <div>{@html post.content}</div>
      </div> -->

      {#if post}
        <div class="post-article">
          <div class="post-header">
            <!-- <Link to="/posts/{post.header.id}/edit">Edit</Link> -->
            <button
              class=""
              on:click={() => {
                deletePost();
              }}>Delete</button
            >
          </div>
          <h1>{post.header.title}</h1>
          <div>{@html post.content}</div>
        </div>
      {/if}
    </div>
    <div class="post-sidebar">
      <div class="post-sidebar-content" />
    </div>
  </div>
</div>

<style>
  .post-container {
    max-width: 1336px;
    text-align: left;
    margin: auto;
    height: 100%;
  }

  .post-flexbox {
    display: flex;
    justify-content: space-evenly;
    flex-direction: row;
  }

  .post-content {
    max-width: 728px;
    flex: 1 1 auto;
    justify-content: center;
  }

  .post-article {
    max-width: 680px;
    margin: 0 24px;
    min-width: 0px;
    /* width: 100%; */
  }

  .post-header {
    margin-top: 56px;
    margin-bottom: 32px;
  }

  .post-sidebar {
    min-height: 100vh;
    border-left: 1px solid rgba(242, 242, 242, 1);
  }

  .post-sidebar-content {
    height: 100%;
    width: 100%;
  }

  @media (max-width: 903.98px) {
    .post-content {
      min-width: 0;
    }
    /* .post-article {
      margin: 0 0;
    } */
    .post-sidebar {
      min-width: 0;
    }
  }

  @media (min-width: 904px) and (max-width: 1079.98px) {
    .post-content {
      max-width: 680px;
    }

    .post-sidebar {
      max-width: 352px;
      min-width: 310px;
    }
  }

  @media (min-width: 1080px) {
    .post-sidebar {
      max-width: 352px;
      min-width: 352px;
      padding-right: 24px;
    }
  }

  @media (min-width: 1080px) {
    .post-sidebar {
      padding-left: clamp(24px, 24px + 100vw - 1080px, 40px);
    }
  }
</style>
