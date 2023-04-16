<script lang="ts">
  import { UserState, type AppUser, type Post } from "./DataInterface";
  import { appService } from "./DataService";

  export let post: Post = undefined;

  let localUser: AppUser = undefined;
  let localUserState: UserState = UserState.Unknown;

  appService.userState.subscribe((value) => {
    localUserState = value;
  });

  appService.user.subscribe((value) => {
    localUser = value;
  });

  function deletePost() {
    appService.DeletePost(post.header.id).then((result) => {
      appService.Navigate("/");
    });
  }
</script>

<div class="container">
  {#if post}
    <div class="frame">
      <div class="profile">
        <img
          alt="user profile"
          class="profile"
          width="24"
          src= {post.header.authorProfilePic}
          referrerpolicy="no-referrer"
        >
      </div>
      <div class="bylines">
        <div class="byline_left">
          <div class="byline_author">
            {post.header.authorDisplayName}
          </div>
          <div class="byline_details">
            {new Date(post.header.created).toDateString()} · 2 min read · Listen
          </div>
        </div>
        <div class="byline_right">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none"
            ><path
              d="M20 5.34c-.67.41-1.4.7-2.18.87a3.45 3.45 0 0 0-5.02-.1 3.49 3.49 0 0 0-1.02 2.47c0 .28.03.54.07.8a9.91 9.91 0 0 1-7.17-3.66 3.9 3.9 0 0 0-.5 1.74 3.6 3.6 0 0 0 1.56 2.92 3.36 3.36 0 0 1-1.55-.44V10c0 1.67 1.2 3.08 2.8 3.42-.3.06-.6.1-.94.12l-.62-.06a3.5 3.5 0 0 0 3.24 2.43 7.34 7.34 0 0 1-4.36 1.49l-.81-.05a9.96 9.96 0 0 0 5.36 1.56c6.4 0 9.91-5.32 9.9-9.9v-.5c.69-.49 1.28-1.1 1.74-1.81-.63.3-1.3.48-2 .56A3.33 3.33 0 0 0 20 5.33"
              fill="#A8A8A8"
            /></svg
          >
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none"
            ><path
              d="M19.75 12.04c0-4.3-3.47-7.79-7.75-7.79a7.77 7.77 0 0 0-5.9 12.84 7.77 7.77 0 0 0 4.69 2.63v-5.49h-1.9v-2.2h1.9v-1.62c0-1.88 1.14-2.9 2.8-2.9.8 0 1.49.06 1.69.08v1.97h-1.15c-.91 0-1.1.43-1.1 1.07v1.4h2.17l-.28 2.2h-1.88v5.52a7.77 7.77 0 0 0 6.7-7.71"
              fill="#A8A8A8"
            /></svg
          >

          <svg width="24" height="24" viewBox="0 0 24 24" fill="none"
            ><path
              d="M19.75 5.39v13.22a1.14 1.14 0 0 1-1.14 1.14H5.39a1.14 1.14 0 0 1-1.14-1.14V5.39a1.14 1.14 0 0 1 1.14-1.14h13.22a1.14 1.14 0 0 1 1.14 1.14zM8.81 10.18H6.53v7.3H8.8v-7.3zM9 7.67a1.31 1.31 0 0 0-1.3-1.32h-.04a1.32 1.32 0 0 0 0 2.64A1.31 1.31 0 0 0 9 7.71v-.04zm8.46 5.37c0-2.2-1.4-3.05-2.78-3.05a2.6 2.6 0 0 0-2.3 1.18h-.07v-1h-2.14v7.3h2.28V13.6a1.51 1.51 0 0 1 1.36-1.63h.09c.72 0 1.26.45 1.26 1.6v3.91h2.28l.02-4.43z"
              fill="#A8A8A8"
            /></svg
          >
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none"
            ><path
              fill-rule="evenodd"
              clip-rule="evenodd"
              d="M3.57 14.67c0-.57.13-1.11.38-1.6l.02-.02v-.02l.02-.02c0-.02 0-.02.02-.02.12-.26.3-.52.57-.8L7.78 9v-.02l.01-.02c.44-.41.91-.7 1.44-.85a4.87 4.87 0 0 0-1.19 2.36A5.04 5.04 0 0 0 8 11.6L6.04 13.6c-.19.19-.32.4-.38.65a2 2 0 0 0 0 .9c.08.2.2.4.38.57l1.29 1.31c.27.28.62.42 1.03.42.42 0 .78-.14 1.06-.42l1.23-1.25.79-.78 1.15-1.16c.08-.09.19-.22.28-.4.1-.2.15-.42.15-.67 0-.16-.02-.3-.06-.45l-.02-.02v-.02l-.07-.14s0-.03-.04-.06l-.06-.13-.02-.02c0-.02 0-.03-.02-.05a.6.6 0 0 0-.14-.16l-.48-.5c0-.04.02-.1.04-.15l.06-.12 1.17-1.14.09-.09.56.57c.02.04.08.1.16.18l.05.04.03.06.04.05.03.04.04.06.1.14.02.02c0 .02.01.03.03.04l.1.2v.02c.1.16.2.38.3.68a1 1 0 0 1 .04.25 3.2 3.2 0 0 1 .02 1.33 3.49 3.49 0 0 1-.95 1.87l-.66.67-.97.97-1.56 1.57a3.4 3.4 0 0 1-2.47 1.02c-.97 0-1.8-.34-2.49-1.03l-1.3-1.3a3.55 3.55 0 0 1-1-2.51v-.01h-.02v.02zm5.39-3.43c0-.19.02-.4.07-.63.13-.74.44-1.37.95-1.87l.66-.67.97-.98 1.56-1.56c.68-.69 1.5-1.03 2.47-1.03.97 0 1.8.34 2.48 1.02l1.3 1.32a3.48 3.48 0 0 1 1 2.48c0 .58-.11 1.11-.37 1.6l-.02.02v.02l-.02.04c-.14.27-.35.54-.6.8L16.23 15l-.01.02-.01.02c-.44.42-.92.7-1.43.83a4.55 4.55 0 0 0 1.23-3.52L18 10.38c.18-.21.3-.42.35-.65a2.03 2.03 0 0 0-.01-.9 1.96 1.96 0 0 0-.36-.58l-1.3-1.3a1.49 1.49 0 0 0-1.06-.42c-.42 0-.77.14-1.06.4l-1.2 1.27-.8.8-1.16 1.15c-.08.08-.18.21-.29.4a1.66 1.66 0 0 0-.08 1.12l.02.03v.02l.06.14s.01.03.05.06l.06.13.02.02.01.02.01.02c.05.08.1.13.14.16l.47.5c0 .04-.02.09-.04.15l-.06.12-1.15 1.15-.1.08-.56-.56a2.3 2.3 0 0 0-.18-.19c-.02-.01-.02-.03-.02-.04l-.02-.02a.37.37 0 0 1-.1-.12c-.03-.03-.05-.04-.05-.06l-.1-.15-.02-.02-.02-.04-.08-.17v-.02a5.1 5.1 0 0 1-.28-.69 1.03 1.03 0 0 1-.04-.26c-.06-.23-.1-.46-.1-.7v.01z"
              fill="#A8A8A8"
            /></svg
          >
          {#if localUser && (localUser.uid === post.header.authorId)}
            <a href={"/posts/" + post.header.id + "/edit"} class="edit_button" title="Edit post">
              <!-- <span class="text">Edit</span> -->
              <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                aria-label="Write"
                ><path
                  d="M14 4a.5.5 0 0 0 0-1v1zm7 6a.5.5 0 0 0-1 0h1zm-7-7H4v1h10V3zM3 4v16h1V4H3zm1 17h16v-1H4v1zm17-1V10h-1v10h1zm-1 1a1 1 0 0 0 1-1h-1v1zM3 20a1 1 0 0 0 1 1v-1H3zM4 3a1 1 0 0 0-1 1h1V3z"
                  fill="currentColor"
                /><path
                  d="M17.5 4.5l-8.46 8.46a.25.25 0 0 0-.06.1l-.82 2.47c-.07.2.12.38.31.31l2.47-.82a.25.25 0 0 0 .1-.06L19.5 6.5m-2-2l2.32-2.32c.1-.1.26-.1.36 0l1.64 1.64c.1.1.1.26 0 .36L19.5 6.5m-2-2l2 2"
                  stroke="currentColor"
                /></svg
              >
            </a>
            <span class="delete_button" title="Delete post" on:click={deletePost} on:keydown={deletePost}>
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" class="ov ow"><path d="M12 21a9 9 0 1 0 0-18 9 9 0 0 0 0 18zM8.25 12h7.5" stroke="#888" stroke-linecap="round" stroke-linejoin="round"></path></svg>
            </span>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .container {
    cursor: pointer;
    display: flex;
    margin-bottom: 30px;
    margin-top: 15px;
    padding-top: 30px;
  }

  .frame {
    margin: 0px 24px;
    max-width: 680px;
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
    height: 48px;
    width: 48px;
    border-radius: 50%;
  }

  .byline_left {
    display: flex;
    flex-direction: column;
    width: 100%;
    font-size: 19px;
  }

  .byline_details {
    width: 100%;
    margin-top: 4px;
    font-size: 14px;
    color: rgba(117, 117, 117, 1);
  }

  .byline_right {
    width: 100%;
    text-align: right;
  }

  .delete_button {
    color: rgba(117, 117, 117, 1);
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
</style>
