<script lang="ts">

  import type { AppUser, PostComment } from "./DataInterface";
  import { appService } from "./DataService";

  export let parentId: string = "";
  export let createComment: (commentFormData: FormData) => Promise<boolean>;

  let localUser: AppUser = undefined;
  let commentText: string = ""; 

  appService.user.subscribe((value) => {
    localUser = value;
  });

  function doSubmit() {
    const formData = new FormData();

    formData.set("authorId", localUser.uid);
    formData.set("authorDisplayName", localUser.displayName);
    formData.set("authorProfilePic", localUser.photoURL);
    formData.set("content", commentText);

    if (parentId)
      formData.set("parentCommentId", parentId);

    createComment(formData).then((result: boolean) => {
      commentText = "";
    })
  }

  function makeid(length: number): string {
    let result = "";
    const characters =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    const charactersLength = characters.length;
    let counter = 0;
    while (counter < length) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
      counter += 1;
    }
    return result;
  }  
</script>

<div class="component">
  {#if localUser}
    <div class="frame">
      <div class="profile">
        <img
          alt="user profile"
          class="profile"
          width="16"
          src={localUser.photoURL}
          referrerpolicy="no-referrer"
        />
      </div>
      <div class="bylines">
        <div class="byline_left">
          <div class="byline_author">
            {localUser.displayName}
          </div>
        </div>
      </div>
    </div>
    <textarea class="content" placeholder="What are your thoughts?" rows="4" bind:value={commentText}></textarea>
    <div class="controls">
      <button disabled={commentText ? false : true} class="publishbutton" on:click={doSubmit}>Respond</button>
    </div>
  {/if}
</div>

<style>

  .component {
    /* width: 100%; */
    max-width: 680px;
    min-height: 170px;
    margin-left: 18px;
    box-shadow: rgba(0, 0, 0, 0.12) 0px 2px 8px;
  }

  .frame {
    padding-top: 14px;
    padding-left: 18px;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: row;
  }

  .bylines {
    margin-left: 10px;
    margin-top: 8px;
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

  .content {
    margin-left: 18px;
    margin-top: 12px;
    box-sizing: border-box;
    border: 0px solid lightgray;
    font-family: sohne, "Helvetica Neue", Helvetica, Arial, sans-serif;
    width: 90%;
    font-size: 14px;
  }

  .content:focus {
    outline: none !important;
  }

  .controls {
    padding-right: 18px;
    padding-top: 8px;
    padding-bottom: 12px;
    text-align: right;
  }

  .publishbutton {
    cursor: pointer;
    user-select: none;
    border-radius: 5px;
    border-width: 1px;
    border-style: solid;
    /* margin-right: 20px; */
    height: 26px;
    width: 65px;
    /* margin-top: 7px; */
    background: #1a8917;
    border-color: #1a8917;
    color: white;
  }

  .publishbutton:disabled {
    opacity: .3;
    cursor: default;
  }
</style>