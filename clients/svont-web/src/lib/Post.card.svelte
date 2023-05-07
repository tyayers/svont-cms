<script lang="ts">
  import { onMount } from "svelte";

  import type { PostOverview } from "./DataInterface";
  import { ToTitleCase } from "./DataInterface";
  import { appService } from "./DataService";
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
          <div class="content_body">
            <span
              class:content_left={post.image}
              class:content_left_full={!post.image}
            >
              <h2 class="title">{post.title}</h2>
              {@html post.summary}
            </span>

            {#if post.image}
              <span class="content_right">
                <img
                  class="content_right_image"
                  src={appService.GetServer("IMAGE") +
                    "/posts/" +
                    post.id +
                    "/files/" +
                    post.image}
                  alt="Post preview"
                />
              </span>
            {/if}
          </div>
          <br />
          {#if post.tags && post.tags.length > 0}
            <div class="tags_box">
              Tags:
              {#each post.tags as tag}
                {#if tag}
                  <a class="tag" href={"/tags/" + tag} title={ToTitleCase(tag)}
                    >{ToTitleCase(tag)}</a
                  >
                {/if}
              {/each}
            </div>
          {/if}
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

  .content_body {
    display: flex;
  }

  .content_left {
    width: 70%;
  }

  .content_left_full {
    width: 100%;
  }

  .content_right {
    width: 25%;
    margin-left: 10px;
    margin-top: auto;
    height: 100%;
    /* margin-bottom: auto; */
  }

  .tags_box {
    /* margin-left: 23px; */
    /* margin-top: 44px; */
    margin-bottom: 18px;
    font-size: 15px;
    font-weight: 500;
    color: gray;
  }

  .tag {
    margin-right: 6px;
    background-color: rgb(231, 231, 231);
    border-radius: 25px;
    padding: 4px 10px 4px 10px;
    font-size: 14px;
    color: gray;
    user-select: none;
    cursor: pointer;
    /* text-transform: capitalize;
    display: inline-block; */
  }

  .footer {
    color: gray;
    font-size: 12px;
  }
</style>
