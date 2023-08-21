<script lang="ts">
  import { appService } from "./DataService";
  import type { PostComment, AppUser } from './DataInterface';

  import Comment from './Comment.svelte'
  import NewComment from './Comment.new.svelte';

  export let data: PostComment[] = [];
  export let createComment: (commentFormData: FormData) => Promise<boolean>;
  export let upvoteComment: (commentId: string) => Promise<PostComment>;

  let localUser: AppUser = undefined;
  
  appService.user.subscribe((value) => {
    localUser = value;
  });
</script>

<div>
  {#if localUser}
    <NewComment {createComment} />
  {:else}
    <div class="sign_in_msg">
      <a class="sign_in_link" href="/signin">Sign in</a> to comment.
    </div>
  {/if}
  {#if data && data.length > 0}
    {#each data as comment, i}
      <Comment {createComment} {upvoteComment} data={comment} />
    {/each}
  {/if}
</div>

<style>
  .sign_in_msg {
    text-align: center;
    color: gray;
  }

  .sign_in_link {
    text-decoration: underline;
  }
</style>