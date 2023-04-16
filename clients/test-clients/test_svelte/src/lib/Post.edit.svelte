<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import EditorCKBBlock, {getData, setData} from "./Editor.CKBBlock.svelte";
  import Header from "./Header.svelte";
  import type { PostOverviewCollection, PostOverview, Post } from "./Types";
  
  export let id
  
  let post: Post

  onMount(async function () {
    const response = await fetch(import.meta.env.VITE_CMS_SERVICE + "/posts/" + id);
    const data: Post = await response.json();
    post = data;

    var titleInput: HTMLInputElement = document.getElementById("title") as HTMLInputElement
    if (titleInput) titleInput.value = post.title

    setData(post.content)
  });

  function onSubmit(e) {

    const formData = new FormData(e.target);
    formData.set("content", getData())
    
    fetch(import.meta.env.VITE_CMS_SERVICE + "/posts/" + id,
      {
        body: formData,
        method: "put",
        headers: {
          'Accept': 'application/json'
        }
      })
      .then((response) => {
        return response.json()
       })
      .then((data) => {
        navigate("/", { replace: true });
      })
  }
</script>

<Header></Header>
<div class="new-container">
  <!-- <h1>New Post</h1> -->

    <form on:submit|preventDefault={onSubmit}>
      <!-- <label for="title">Enter the post title: </label> -->
      <input type="text" name="title" id="title" placeholder="" required>
      
      {#if post}
        <div>
          <!-- <label for="content">Enter the post content: </label><br/> -->
          <!-- <textarea name="content" id="content" class="post_content"></textarea> -->
          <br/>
          <EditorCKBBlock />
        </div>
        <br/>
        <div>
          <label for="files">Select a file:</label>
          <input type="file" id="files" name="files" multiple>  
        </div>
      {/if}
      <button id="save-button" type="submit">Submit</button>
    </form>
</div>

<style>
  .new-container {
    max-width: 904px;
    text-align: left;
    margin: auto;
    margin-top: 20px;
  }

  #title {
    font-size: 48px;
    border-style: none;
    outline: none;
  }

  #save-button {
    position: fixed;
    top: 28px;
    right: 175px;
  }

</style>