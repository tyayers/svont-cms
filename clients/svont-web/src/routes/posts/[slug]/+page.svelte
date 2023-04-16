<script lang="ts">
  import Header from "../../../lib/Header.svelte";
  import PostHeader from "../../../lib/Post.header.svelte";
  import PostFooter from "../../../lib/Post.footer.svelte";
  import Comments from "../../../lib/Comments.svelte";
  import type {
    PostOverviewCollection,
    Post,
    PostComment,
  } from "../../../lib/DataInterface";
  import { appService } from "$lib/DataService";

  export let data;

  function deletePost() {
    // fetch(import.meta.env.VITE_CMS_SERVICE + "/posts/" + id, {
    //   method: "delete",
    //   headers: {
    //     Accept: "application/json",
    //   },
    // }).then((response) => {
    //   navigate("/", { replace: true });
    // });
  }

  function createComment(commentFormData: FormData): Promise<boolean> {
    return new Promise((resolve, reject) => {

      appService.CreateComment(data.post.header.id, commentFormData).then((result) => {
        data.post.header.commentCount++;
        data.comments = result;
        resolve(true);
      });
    });
  }

  function upvoteComment(commentId: string): Promise<PostComment> {
    return new Promise<PostComment>((resolve, reject) => {
      appService.UpvoteComment(data.post.header.id, commentId).then((result: PostComment) => {
        resolve(result);
      });
    });
  }
</script>

<div>
  <Header />

  <div class="post-container">
    <div class="post-flexbox">
      <div class="post-content">
        {#if data}
          <div class="post-article">
            <PostHeader post={data.post} />
            <div class="content">
              <h1>{data.post.header.title}</h1>
              <div>{@html data.post.content}</div>
            </div>
            <PostFooter post={data.post} />
          </div>
          <div class="post-comments">
            <h3 class="post-comments-header">Comments</h3>
            <Comments {createComment} {upvoteComment} data={data.comments} />
          </div>
        {/if}
      </div>
      <div class="post-sidebar">
        <div class="post-sidebar-content" />
      </div>
    </div>
  </div>
</div>

<style>
  .post-container {
    max-width: 1336px;
    text-align: left;
    margin: auto;
    height: 100%;
    padding-bottom: 120px;
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

  .post-comments {
    max-width: 680px;
    margin-left: 24px;
    margin-right: 24px;
    margin-top: 40px;
    min-width: 0px;
  }

  .post-comments-header {
    margin-left: 18px;
  }

  .post-sidebar {
    min-height: 100vh;
    border-left: 1px solid rgba(242, 242, 242, 1);
  }

  .post-sidebar-content {
    height: 100%;
    width: 100%;
  }

  .content {
    margin: 0px 24px;
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
