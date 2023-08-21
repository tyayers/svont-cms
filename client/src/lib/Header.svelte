<script lang="ts">
  import { goto } from "$app/navigation";
  import { createEventDispatcher } from "svelte";
  import {
    type AppUser,
    type SearchResult,
    HeaderButton,
    UserState,
  } from "./DataInterface";
  import { appService, LogoPath } from "./DataService";

  import SearchBox from "./Search.box.svelte";
  import UserMenu from "./User.menu.svelte";

  // Event dispatcher for all broadcasts to parent
  const dispatch = createEventDispatcher();

  export let small: boolean = false;
  export let showSearch: boolean = true;
  export let statusText: string = "";
  export let actionType: HeaderButton = HeaderButton.NewPost;

  let localUser: AppUser = undefined;
  let localUserState: UserState = UserState.Unknown;

  appService.userState.subscribe((value) => {
    localUserState = value;
  });

  appService.user.subscribe((value) => {
    localUser = value;
  });

  function goHome() {
    if (localUserState === UserState.SignedIn) goto("/home");
    else goto("/");
  }

  function signOutClick() {
    appService.SignOut();
  }

  function newPost() {
    goto("/posts/new");
  }

  function searchPosts(input: string): Promise<SearchResult[]> {
    return appService.SearchPosts(input);
  }

  function searchClick(event) {
    goto("/posts/" + event.detail.id);
  }

  function submitClick(event) {
    console.log("header_submitclick");
    dispatch("submit");
  }
</script>

<div class:headersmall={small === true} class="header">
  <div class="left">
    <div class="logo_box" on:click={goHome} on:keydown={goHome}>
      <!-- <img class="logo" src={LogoPath} alt="Site logo" /> -->
      <!-- <span class="title">ggo-1 blog</span> -->
      <span class="title">Nup</span>
    </div>
    {#if showSearch}
      <SearchBox search={searchPosts} on:click={searchClick} />
    {/if}
    {#if statusText}
      <span class="status_text">{statusText}</span>
    {/if}
  </div>
  <div class="right">
    {#if localUserState == UserState.SignedOut}
      <a href="/signin" class="signin">Sign In</a>
    {:else if localUserState == UserState.SignedIn}
      {#if actionType === HeaderButton.NewPost}
        <div class="post">
          <div class="postbutton" on:mousedown={newPost}>
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
            <div class="posttext">Post</div>
          </div>
        </div>
      {:else if actionType === HeaderButton.Submit}
        <button class="publishbutton" on:click={submitClick}>Publish</button>
      {/if}
      <UserMenu user={localUser} on:signOut={signOutClick} />
    {/if}
  </div>
</div>

<style>
  .header {
    height: var(--header-height);
    background-color: var(--main-background-color);
    border-bottom: var(--header-bottom-border);
    padding: var(--header-padding);
    align-items: center;
    display: flex;
    position: sticky;
    top: 0;
    z-index: 1;
  }

  .headersmall {
    max-width: 1050px;
    margin: auto;
    border-bottom: 0px;
  }

  .left {
    display: flex;
    flex: 1 0 auto;
  }

  .right {
    display: flex;
    position: relative;
    right: 40px;
  }

  .logo_box {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 34px;
    cursor: pointer;
  }

  .logo {
    position: relative;
    /* top: 10px; */
    height: 44px;
  }

  .title {
    font-weight: 600;
    font-size: 22px;
    color: darkgray;
  }

  .status_text {
    /* height: var(--header-height); */
    color: gray;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 14px;
  }

  .post {
    font-size: 14px;
    align-items: center;
    display: flex;
    font-weight: 400;
    line-height: 20px;
  }

  .postbutton {
    color: rgba(117, 117, 117, 1);
    display: flex;
    margin-right: 32px;
    cursor: pointer;
    user-select: none;
  }

  .posttext {
    margin-left: 8px;
  }

  .publishbutton {
    cursor: pointer;
    user-select: none;
    border-radius: 99em;
    border-width: 1px;
    border-style: solid;
    margin-right: 20px;
    height: 25px;
    width: 65px;
    margin-top: 7px;
    background: #1a8917;
    border-color: #1a8917;
    color: white;
  }

  .signin {
    cursor: pointer;
    border-radius: 99em;
    border-width: 1px;
    border-style: solid;
    text-decoration: none;
    background: rgba(25, 25, 25, 1);
    fill: rgba(255, 255, 255, 1);
    padding: 7px 16px 9px;
    color: rgba(255, 255, 255, 1);
    line-height: 20px;
    font-size: 14px;
    font-weight: 400;
    width: 60px;
    text-align: center;
    height: 22px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
