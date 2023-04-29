<script lang="ts">
  import Header from "../../../lib/Header.svelte";
  import PostHeader from "../../../lib/Post.header.svelte";
  import PostFooter from "../../../lib/Post.footer.svelte";
  import Comments from "../../../lib/Comments.svelte";
  import PostPopularWidget from "../../../lib/Post.popular.svelte";

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
      appService
        .CreateComment(data.post.header.id, commentFormData)
        .then((result) => {
          data.post.header.commentCount++;
          data.comments = result;
          resolve(true);
        });
    });
  }

  function upvoteComment(commentId: string): Promise<PostComment> {
    return new Promise<PostComment>((resolve, reject) => {
      appService
        .UpvoteComment(data.post.header.id, commentId)
        .then((result: PostComment) => {
          resolve(result);
        });
    });
  }
</script>

<div>
  <Header />

  <div class="content">
    <div class="container">
      <div class="panel_left">
        <div class="pannel_left_inner">
          {#if data}
            <div class="post-article">
              <PostHeader post={data.post} />
              <div class="post-content">
                <h1>{data.post.header.title}</h1>
                <div>{@html data.post.content}</div>
              </div>
              <div class="tags_box">
                {#if data.post.header.tags && data.post.header.tags.length > 0}
                  Tags:
                  {#each data.post.header.tags as tag}
                    {tag}
                  {/each}
                {/if}
              </div>
              <div class="post_attachments">
                {#if data.post.files && data.post.files.length > 0}
                  <div class="post_attachments_header">Attachments:</div>
                  {#each data.post.files as file}
                    <span class="post_attachment">
                      <a
                        href={appService.defaultServer +
                          "/posts/" +
                          data.post.header.id +
                          "/files/" +
                          file}
                      >
                        <span style="white-space: nowrap;">
                          <svg
                            width="16px"
                            height="16px"
                            style="position: relative; top: 3px;"
                            viewBox="0 0 16 16"
                            version="1.1"
                            xmlns="http://www.w3.org/2000/svg"
                            xmlns:xlink="http://www.w3.org/1999/xlink"
                            class="si-glyph si-glyph-paper-clip"
                            fill="#222"
                            ><g id="SVGRepo_bgCarrier" stroke-width="0" /><g
                              id="SVGRepo_tracerCarrier"
                              stroke-linecap="round"
                              stroke-linejoin="round"
                            /><g id="SVGRepo_iconCarrier">
                              <title>Paper-clip</title> <defs />
                              <g
                                stroke="none"
                                stroke-width="1"
                                fill="none"
                                fill-rule="evenodd"
                              >
                                <path
                                  d="M6.346,16 C5.009,16 4,14.907 4,13.725 L4,3.99799991 C4,1.63391113 5.25378418,0 7.69795109,0 L8.3671875,0 C11.046,0 12,1.56054688 12,3.99799991 L12,11.0050049 L11.046,11.0050049 L11.046,3.99799991 C11.046,2.4140625 10.4089355,1 8.3671875,1 L7.68199992,1 C5.87280273,1 5,2.31750488 5,3.99799991 L5,13.725 C5,14.463 5.448,14.999 6.345,14.999 L7.683,14.999 C8.535,14.999 9.062,14.511 9.062,13.725 L9.062,5.756 C9.062,5.225 8.98100008,5.03984473 7.94300008,4.99084473 C6.88400008,5.04284473 7,5.262 7,5.756 L6.99999995,10.0100098 L5.99899995,10.0100098 L5.999,5.756 C5.999,4.635 6.635,4.06 7.943,3.998 C9.249,4.058 10,4.616 10,5.756 L10,13.725 C10,14.947 8.966,16 7.682,16 L6.346,16 Z"
                                  fill="#aaa"
                                  class="si-glyph-fill"
                                />
                              </g>
                            </g></svg
                          >
                          {file}
                        </span>
                      </a>
                    </span>
                  {/each}
                {/if}
              </div>
              <PostFooter post={data.post} />
            </div>
            <div class="post-comments">
              <h3 class="post-comments-header">Comments</h3>
              <Comments {createComment} {upvoteComment} data={data.comments} />
            </div>
            <div class="pannel_left_footer" />
          {/if}
        </div>
      </div>
      <div class="panel_right">
        <div class="widget1">
          <PostPopularWidget posts={data.popular} />
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .content {
    max-width: 1336px;
    text-align: left;
    margin: auto;
    /* height: 100%;
    height: calc(100vh - 58px);
    overflow-y: hidden; */
  }

  .container {
    display: flex;
    justify-content: space-evenly;
    flex-direction: row;
    height: 100%;
    height: calc(100vh - 58px);
    overflow-y: auto;
  }

  .panel_left {
    width: 100%;
    flex: 1 1 auto;
    justify-content: center;
    display: inline-flex;
    /* height: 100%;
    overflow-y: auto; */
  }

  .pannel_left_inner {
    max-width: 728px;
    width: 100%;
    padding-bottom: 54px;
  }

  .pannel_left_footer {
    height: 84px;
  }

  .panel_right {
    min-height: 100vh;
    border-left: 1px solid rgba(242, 242, 242, 1);
    padding-left: 32px;
    min-width: 420px;
    position: sticky;
    top: 0px;
  }

  .widget1 {
    height: 100%;
    width: 100%;
  }

  .post-content {
    max-width: 728px;
    flex: 1 1 auto;
    justify-content: center;
    margin: 0px 24px;
  }

  .post-article {
    max-width: 680px;
    margin: 0 24px;
    min-width: 0px;
    /* width: 100%; */
  }

  .tags_box {
    margin-left: 23px;

    font-size: 15px;
    font-weight: 500;
    color: gray;
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

  .post_attachments {
    margin-top: 32px;
    margin-left: 24px;
    color: rgba(117, 117, 117, 1);
    overflow-wrap: break-word;
  }

  .post_attachments_header {
    font-size: 15px;
    font-weight: 500;
    margin-bottom: 8px;
  }

  .post_attachment {
    font-size: 14px;
    margin-right: 8px;
  }

  @media (max-width: 903.98px) {
    .post-content {
      min-width: 0;
    }
    /* .post-article {
      margin: 0 0;
    } */
    .panel_right {
      display: none;
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
