<script lang="ts">
  import type { PostOverview } from "./DataInterface";


  export let posts: PostOverview[] = [];
</script>

<div class="frame">
  <span class="frame_title">Popular✨</span>
  <div class="posts_frame">
    {#each posts as post}
      <div class="post_frame">
        <a href={"/posts/" + post.id}>
          <div class="user">
            <div class="profile">
              <img
                alt="user profile"
                class="profile"
                width="24"
                src={post.authorProfilePic}
              />
            </div>
            <div class="byline">
              {post.authorDisplayName}
            </div>
          </div>
          <div class="post_title">
            {post.title}
          </div>
        </a>
      </div>
    {/each}
  </div>
</div>

<style>
  .frame {
    margin-top: 42px;
    max-width: 280px;
  }

  .frame_title {
    color: rgba(41, 41, 41, 1);
    font-weight: 600;
    font-size: 17px;
    margin-bottom: 6px;
  }

  .posts_frame {

  }

  .post_frame {
    margin-top: 14px;
    margin-bottom: 28px;
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
    font-size: 14px;
    top: 3px;
    left: 6px;
    font-weight: 600;
  }

  .post_title {
    margin-top: 8px;
    margin-left: 2px;
    font-weight: 800;
  }
</style>