<script lang="ts">
  import { navigate } from "../App.svelte";

  import EditorCKBBlock, { getData } from "./Editor.CKBBlock.svelte";
  import type {
    PostOverviewCollection,
    PostOverview,
    Post,
  } from "./DataInterface";

  function onSubmit(e) {
    const formData = new FormData(e.target);
    let data: string = getData();
    const regex = /<br>(?=(?:\s*<[^>]*>)*$)|(<br>)|<[^>]*>/gi;
    let text: string = data.replace(regex, (x, y) => (y ? " & " : ""));

    formData.set("content", data);
    formData.set("summary", text.substring(0, 200) + "...");

    fetch(import.meta.env.VITE_CMS_SERVICE + "/posts", {
      body: formData,
      method: "post",
      headers: {
        Accept: "application/json",
      },
    })
      .then((response) => {
        return response.json();
      })
      .then((data) => {
        navigate("/");
      });
  }
</script>

<div class="new-container">
  <!-- <h1>New Post</h1> -->
  <form on:submit|preventDefault={onSubmit}>
    <!-- <label for="title">Enter the post title: </label> -->
    <input type="text" name="title" id="title" placeholder="Title" required />

    <div>
      <!-- <label for="content">Enter the post content: </label><br/> -->
      <!-- <textarea name="content" id="content" class="post_content"></textarea> -->
      <br />
      <EditorCKBBlock />
    </div>
    <br />
    <div>
      <label for="files">Select a file:</label>
      <input type="file" id="files" name="files" multiple />
    </div>

    <button type="submit">Submit</button>
  </form>
</div>

<style>
  .new-container {
    max-width: 904px;
    text-align: left;
    margin: auto;
    margin-top: 20px;
    padding: 0px 24px 0px 24px;
  }

  #title {
    font-size: 48px;
    border-style: none;
    outline: none;
  }
</style>
