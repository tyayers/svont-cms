<script lang="ts">
  import type { PostComment } from "./DataInterface";
  import NewComment from "./Comment.new.svelte";
  import { appService } from "./DataService";

  export let data: PostComment = undefined;
  export let createComment: (commentFormData: FormData) => Promise<boolean>;
  export let upvoteComment: (commentId: string) => Promise<PostComment>;

  let editMode: boolean = false;
  let replyMode: boolean = false;
  let backupContent: string = "";

  function startEdit() {
    // backupContent = data.content;
    // editMode = true;
  }

  function cancelEdit() {
    data.content = backupContent;
    editMode = false;
  }

  function submitEdit() {
    editMode = false;
  }

  function startReply() {
    replyMode = true;
  }

  function cancelReply() {
    replyMode = false;
  }

  function init(el){
    el.focus()
  }

  function upvote() {
    data.upvotes++;

    upvoteComment(data.id).then((postComment: PostComment) => {
      //data = postComment;
    });
  }

  function addChildComment(commentFormData: FormData): Promise<boolean> {
    return new Promise((resolve, reject) => {
      createComment(commentFormData);
      replyMode = false;

      resolve(true);
    })
  }

</script>

<div class="outer_container">
  <div>
    <div class="container">
      {#if data}
        <div class="frame">
          <div class="profile">
            <img
              alt="user profile"
              class="profile"
              width="16"
              src={data.authorProfilePic}
              referrerpolicy="no-referrer"
            />
          </div>
          <div class="bylines">
            <div class="byline_left">
              <div class="byline_author">
                {data.authorDisplayName}
              </div>
              <div class="byline_details">
                {new Date(data.created).toDateString()}
              </div>
            </div>
            <div class="byline_right">
              {#if editMode}
                <button class="cancelbutton" on:click={cancelEdit}>Cancel</button>
                <button class="publishbutton" on:click={submitEdit}>Submit</button>
              {:else}
                <span class="edit_button" on:click={startEdit} on:keydown={startEdit}>
                  <span class="text">···</span>
                </span>
              {/if}
            </div>
          </div>
        </div>
      {/if}
    </div>
    <div class="content">
      {#if editMode}
        <textarea class="content_text" rows="4" bind:value={data.content} use:init></textarea>
      {:else}
        <div>{data.content}</div>
      {/if}
      <div class="comment_footer">
        <span class="like_box" on:click={upvote} on:keydown={upvote}>
          <svg version="1.1" id="Layer_1" width="18px" fill="gray" viewBox="0 0 122.88 104.19" style="enable-background:new 0 0 122.88 104.19" xml:space="preserve"><g><path d="M62.63,6.25c0.56-2.85,2.03-4.68,4.04-5.61c1.63-0.76,3.54-0.83,5.52-0.31c1.72,0.45,3.53,1.37,5.26,2.69 c4.69,3.57,9.08,10.3,9.64,18.71c0.17,2.59,0.12,5.35-0.12,8.29c-0.16,1.94-0.41,3.98-0.75,6.1h19.95l0.09,0.01 c3.24,0.13,6.38,0.92,9.03,2.3c2.29,1.2,4.22,2.84,5.56,4.88c1.38,2.1,2.13,4.6,2.02,7.46c-0.08,2.12-0.65,4.42-1.81,6.87 c0.66,2.76,0.97,5.72,0.54,8.32c-0.36,2.21-1.22,4.17-2.76,5.63c0.08,3.65-0.41,6.71-1.39,9.36c-1.01,2.72-2.52,4.98-4.46,6.98 c-0.17,1.75-0.45,3.42-0.89,4.98c-0.55,1.96-1.36,3.76-2.49,5.35l0,0c-3.4,4.8-6.12,4.69-10.43,4.51c-0.6-0.02-1.24-0.05-2.24-0.05 l-39.03,0c-3.51,0-6.27-0.51-8.79-1.77c-2.49-1.25-4.62-3.17-6.89-6.01l-0.58-1.66V47.78l1.98-0.53 c5.03-1.36,8.99-5.66,12.07-10.81c3.16-5.3,5.38-11.5,6.9-16.51V6.76L62.63,6.25L62.63,6.25L62.63,6.25z M4,43.02h29.08 c2.2,0,4,1.8,4,4v53.17c0,2.2-1.8,4-4,4l-29.08,0c-2.2,0-4-1.8-4-4V47.02C0,44.82,1.8,43.02,4,43.02L4,43.02L4,43.02z M68.9,5.48 c-0.43,0.2-0.78,0.7-0.99,1.56V20.3l-0.12,0.76c-1.61,5.37-4.01,12.17-7.55,18.1c-3.33,5.57-7.65,10.36-13.27,12.57v40.61 c1.54,1.83,2.96,3.07,4.52,3.85c1.72,0.86,3.74,1.2,6.42,1.2l39.03,0c0.7,0,1.6,0.04,2.45,0.07c2.56,0.1,4.17,0.17,5.9-2.27v-0.01 c0.75-1.06,1.3-2.31,1.69-3.71c0.42-1.49,0.67-3.15,0.79-4.92l0.82-1.75c1.72-1.63,3.03-3.46,3.87-5.71 c0.86-2.32,1.23-5.11,1.02-8.61l-0.09-1.58l1.34-0.83c0.92-0.57,1.42-1.65,1.63-2.97c0.34-2.1-0.02-4.67-0.67-7.06l0.21-1.93 c1.08-2.07,1.6-3.92,1.67-5.54c0.06-1.68-0.37-3.14-1.17-4.35c-0.84-1.27-2.07-2.31-3.56-3.09c-1.92-1.01-4.24-1.59-6.66-1.69v0.01 l-26.32,0l0.59-3.15c0.57-3.05,0.98-5.96,1.22-8.72c0.23-2.7,0.27-5.21,0.12-7.49c-0.45-6.72-3.89-12.04-7.56-14.83 c-1.17-0.89-2.33-1.5-3.38-1.77C70.04,5.27,69.38,5.26,68.9,5.48L68.9,5.48L68.9,5.48z"></path></g></svg>
          <span class="upvote_text">{data.upvotes} likes</span>
        </span>
        <span class="reply_box">
          {#if !replyMode}
            <span class="reply_link" on:click={startReply} on:keydown={startReply}>
              <span class="reply_text">Reply</span>
            </span>
          {:else}
            <span class="reply_link" on:click={cancelReply} on:keydown={cancelReply}>
              <span class="reply_text">Cancel</span>
            </span>
          {/if}
        </span>
      </div>
    </div>
    {#if replyMode}
      <div class="reply_frame">
        <NewComment parentId={data.id} createComment={addChildComment} />
      </div>
    {/if}
  </div>
  {#each data.children as child_comment}
    <div class="child_frame">
      <svelte:self data={child_comment} createComment={addChildComment} {upvoteComment} />
    </div>
  {/each}
</div>

<style>
  .outer_container {
    max-width: 680px;
  }

  .container {
    cursor: pointer;
    display: flex;
    margin-bottom: 16px;
    margin-top: 16px;
    padding-top: 16px;
  }

  .frame {
    margin: 0px 24px;
    /* max-width: 680px; */
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: row;
  }

  .bylines {
    margin-left: 10px;
    width: 100%;
    display: flex;
    flex-direction: row;
  }

  .profile {
    height: 34px;
    width: 34px;
    border-radius: 50%;
  }

  .byline_left {
    display: flex;
    flex-direction: column;
    width: 100%;
    font-size: 15px;
  }

  .byline_details {
    width: 100%;
    margin-top: 4px;
    font-size: 12px;
    color: rgba(117, 117, 117, 1);
  }

  .byline_right {
    width: 100%;
    text-align: right;
  }

  .edit_button {
    margin-left: 20px;
    color: rgba(117, 117, 117, 1);
  }

  .edit_button .text {
    position: relative;
    font-size: 15px;
    top: -6px;
  }

  .publishbutton {
    cursor: pointer;
    user-select: none;
    border-radius: 5px;
    border-width: 1px;
    border-style: solid;
    margin-right: 20px;
    height: 25px;
    width: 65px;
    /* margin-top: 7px; */
    background: #1a8917;
    border-color: #1a8917;
    color: white;
  }

  .cancelbutton {
    cursor: pointer;
    user-select: none;
    border-radius: 5px;
    border-width: 1px;
    border-style: solid;
    /* margin-right: 20px; */
    height: 25px;
    width: 65px;
    /* margin-top: 7px; */
    background: #ffffff;
    border-color: #d0d0d0;
    color: #5a5a5a;
    margin-right: 2px;
  }

  .content {
    margin: 0px 24px;
    max-width: 680px;
    width: 100%;
    font-size: 15px;
    /* color: rgba(117, 117, 117, 1); */
    color: black;
  }

  .content_text {
    width: 90%;
    border-radius: 5px;
    border: 1px solid lightgray;
    padding: 5px;
    box-sizing: border-box;
    font-family: sohne, "Helvetica Neue", Helvetica, Arial, sans-serif;
  }

  .content_text:focus {
    outline: none !important;
    border: 1px solid lightgray;
  }

  .comment_footer {
    margin-top: 12px;
    cursor: pointer;
    user-select: none;
    display: flex;
    flex-direction: row;
  }

  .like_box {
    width: 150px;
  }

  .upvote_text {
    width: 100px;
    position: relative;
    top: 0px;
    left: 6px;
    font-size: 14px;
    color: rgba(117, 117, 117, 1);
  }

  .reply_box {
    width: 90%;
    text-align: right;
  }

  .reply_text {
    margin-right: 44px;
    color: black;
    font-size: 14px;
  }

  .reply_frame {
    margin-top: 12px;
    margin-left: 20px;
  }

  .child_frame {
    margin-left: 20px;
  }
</style>
