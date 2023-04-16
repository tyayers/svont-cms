<script lang="ts">
  import type { AppUser, SearchResult } from "./DataInterface";
  import { SearchPosts, LogoPath } from "./DataService";

  import { navigate, signUserOut, user } from "../App.svelte";

  import SignIn from "./SignIn.svelte";
  import SearchBox from "./Search.box.svelte";

  export let small = false;

  let localUser: AppUser = undefined;
  let userSignedIn: boolean = primitiveToBoolean(
    localStorage.getItem("UserSignedIn")
  );
  let showSignIn = false;

  user.subscribe((value) => {
    localUser = value;
    if (value != null && value != undefined) {
      console.log("user is there");
      // user is signed in
      userSignedIn = true;
      showSignIn = false;
    } else if (value === null) {
      // user is signed out
      console.log("no user");
      userSignedIn = false;
      showSignIn = false;
    } else {
      // we don't know the user state
      console.log("we just don't know if the user is there or not");
      showSignIn = false;
    }
  });

  function goHome() {
    navigate("/");
  }

  function signOut() {
    signUserOut();
  }

  function newPost() {
    navigate("/new");
  }

  function searchClick(event) {
    navigate("/posts/" + event.detail.id);
  }

  function primitiveToBoolean(
    value: string | number | boolean | null | undefined
  ): boolean {
    if (typeof value === "string") {
      return value.toLowerCase() === "true" || !!+value; // here we parse to number first
    }

    return !!value;
  }
</script>

<div class:headersmall={small === true} class="header">
  <div class="left">
    <div on:click={goHome} on:keydown={goHome}>
      <img class="logo" src={LogoPath} alt="Site logo" />
    </div>
    <SearchBox search={SearchPosts} on:click={searchClick} />
  </div>
  <div class="right">
    {#if !userSignedIn}
      <button class="signin" on:click={() => (showSignIn = true)}
        >Sign In</button
      >
    {:else}
      <div class="post">
        <div class="postbutton" on:click={newPost}>
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
      <button class="signin" on:click={signOut}>Sign Out</button>
    {/if}
  </div>
</div>

{#if showSignIn}
  <SignIn
    on:cancel={() => {
      showSignIn = false;
    }}
  />
{/if}

<style>
  .header {
    height: 57px;
    background-color: rgba(255, 255, 255, 1);
    /* width: 100vw; */
    border-bottom: solid 1px rgba(242, 242, 242, 1);
    padding: 0 24px;
    align-items: center;
    display: flex;
    position: sticky;
    top: 0;
    z-index: 1;
    /* top: 16px;
    position: relative; */
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

  .logo {
    position: relative;
    top: 10px;
    height: 44px;
    cursor: pointer;
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
    min-width: 89px;
  }
</style>
