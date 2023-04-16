<script lang="ts">
  import { navigate } from "../App.svelte";
  import type { AppUser } from "./DataInterface";
  import Header from "./Header.svelte";

  let localUser = undefined;

  export let user: AppUser = undefined;

  if (user) {
    console.log("User exists, going home...");
    navigate("/home");
  }
</script>

{#if user == undefined}
  <div class="pc">
    <h1>Welcome to the Cloud Hub.</h1>
    <h3>Log in above with your Bayer, Quantiphi, or Google account.</h3>
  </div>
{/if}

<style>
  .pc {
    text-align: center;
    font-family: sans-serif;
    padding-top: 150px;
  }

  h3 {
    color: #888;
  }
</style>
