<script lang="ts">
  import { goto } from "$app/navigation";

  import Header from "../../lib/Header.svelte";
  import SignIn from "../../lib/SignIn.svelte";
</script>

<Header />

<SignIn
  on:cancel={() => {
    //goto("/");
    history.back();
  }}
/>
