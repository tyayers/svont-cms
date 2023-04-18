<script lang="ts">
  import { onMount } from "svelte";

  import type { PostOverview } from "./DataInterface";

  export let post: PostOverview = undefined;

  let dateString = "";

  onMount(async function () {
    console.log(post.created);
    let nd = new Date(post.created);
    dateString = nd.toDateString();
  });
</script>

<div class="container">
  {#if post}
    <a href={"/posts/" + post.id}>
      <div class="frame">
        <div class="user">
          <div class="profile">
            <img
              class="profile"
              width="24"
              alt="Profile of user"
              src={post.authorProfilePic}
            />
          </div>
          <div class="byline">
            {post.authorDisplayName} · {new Date(post.created).toDateString()}
          </div>
        </div>
        <div class="content">
          <!-- <Link to="/posts/{id}" -->
          <h2 class="title">{post.title}</h2>
          <div>
            {post.summary}
          </div>
          <!--/Link> -->
          <br />
          <span class="footer"
            >{post.upvotes} Likes - {post.commentCount} Comments - {post.fileCount}
            Attachments</span
          >
        </div>
      </div>
    </a>
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
