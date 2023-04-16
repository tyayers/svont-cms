<script lang="ts">
  import { onMount } from "svelte";

  import type { AppUser, PostOverview } from "./DataInterface";
  import { navigate, user } from "../App.svelte";

  let localUser: AppUser = undefined;
  user.subscribe((value) => {
    localUser = value;
  });

  export let post: PostOverview = undefined;

  let dateString = "";

  function loadPost(postId: string) {
    navigate("/posts/" + postId);
  }

  onMount(async function () {
    console.log(post.created);
    let nd = new Date(post.created);
    dateString = nd.toDateString();
  });
</script>

<div class="container">
  {#if post}
    <div class="frame">
      <div class="user">
        <div class="profile">
          <img
            class="profile"
            width="24"
            src="https://images.unsplash.com/photo-1494790108377-be9c29b29330?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8dXNlciUyMHByb2ZpbGV8ZW58MHx8MHx8&w=1000&q=80"
          />
        </div>
        <div class="byline">
          {post.author} · {new Date(post.created).toDateString()}
        </div>
      </div>
      <div class="content" on:click={loadPost(post.id)}>
        <!-- <Link to="/posts/{id}" -->
        <h2 class="title">{post.title}</h2>
        <div>
          {post.summary}
        </div>
        <!--/Link> -->
        <br />
        <span class="footer"
          >{post.commentCount} Comments - {post.upvotes} Likes - {post.fileCount}
          Attachments</span
        >
      </div>
    </div>
  {/if}
</div>

<style>
  .container {
    cursor: pointer;
    display: flex;
    margin-bottom: 40px;
    margin-top: 15px;
    padding-top: 30px;
  }

  .frame {
    margin: 0px 24px;
    max-width: 680px;
    width: 100%;
    height: 100%;
  }

  .user {
    display: flex;
  }

  .profile {
    height: 24px;
    width: 24px;
    border-radius: 50%;
  }

  .byline {
    position: relative;
    top: 3px;
    left: 6px;
  }

  .content {
    text-align: left;
    text-decoration: none;
  }

  .title {
    text-decoration: none;
  }

  .footer {
    color: gray;
    font-size: 12px;
  }
</style>
