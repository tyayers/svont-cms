<script lang="ts">
  // import { browser } from "$app/environment";
  import { goto } from "$app/navigation";
  import { UserState } from "$lib/DataInterface";

  import { appService } from "../lib/DataService";
  import Header from "../lib/Header.svelte";

  goto("/home");

  // appService.userState.subscribe((value) => {
  //   console.log(`home page received userstate ${value}`);
  //   //if (value === UserState.SignedIn) {
  //     console.log("sending user home.");
  //     goto("/home");
  //   //}
  // });
</script>

<Header />
