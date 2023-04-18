<script lang="ts">
  import { goto } from "$app/navigation";
  import Editor, { getData, setData } from "./Editor.CKBBlock.svelte";
  import { appService } from "./DataService";

  import type {
    Post,
    AppUser,
  } from "./DataInterface";

  let localUser: AppUser = undefined;

  appService.user.subscribe((value) => {
    localUser = value;
  });

  export let post: Post = undefined;

  if (post) {
    setData(post.content);
  }
  else
    setData("")

  export function submit() {
    var myForm: HTMLFormElement = document.getElementById(
      "new_post_form"
    ) as HTMLFormElement;

    const formData = new FormData(myForm);

    // Set content
    let content: string = getData();
    const regex = /<br>(?=(?:\s*<[^>]*>)*$)|(<br>)|<[^>]*>/gi;
    let text: string = content.replace(regex, (x, y) => (y ? " & " : ""));
    let textPieces: string[] = text.split(" ");
    let wordLimit = textPieces.length > 50 ? 50 : textPieces.length;
    let summaryText: string = "";
    for (var i = 0; i < wordLimit; i++) {
      summaryText += " " + textPieces[i];
    }
    if (textPieces.length > wordLimit) summaryText += "...";
    formData.set("content", content);

    // Set summary
    //formData.set("summary", text.substring(0, 200) + "...");
    formData.set("summary", summaryText);

    if (post) {
      appService.UpdatePost(post.header.id, formData).then((post: Post) => {
        goto("/home");
      })
    }
    else {
      // Set user for new post
      formData.set("authorId", localUser.uid);
      formData.set("authorDisplayName", localUser.displayName);
      formData.set("authorProfilePic", localUser.photoURL);
      
      appService.CreatePost(formData).then((post: Post) => {
        goto("/home");
      });
    }
  }
</script>

<div class="new-container">
  <!-- <h1>New Post</h1> -->
  <form id="new_post_form">
    <!-- <label for="title">Enter the post title: </label> -->
    {#if post}
      <input type="text" name="title" id="title" placeholder="Title" required value={post.header.title} />

    {:else}
      <input type="text" name="title" id="title" placeholder="Title" required />

    {/if}

    <div>
      <!-- <label for="content">Enter the post content: </label><br/> -->
      <!-- <textarea name="content" id="content" class="post_content"></textarea> -->
      <br />
      <Editor />
    </div>
    <br />
    <div>
      <br />
      <label class="attachlabel" for="files">Attachments:</label>
      <br /><br />
      <input
        class="attachbutton"
        type="file"
        id="files"
        name="files"
        multiple
      />
    </div>

    <!-- <button type="submit">Submit</button> -->
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
    margin-left: 4px;
  }

  .attachlabel {
    margin-left: 10px;
  }

  input[type="file"]::file-selector-button {
    cursor: pointer;
    user-select: none;
    border-radius: 99em;
    border-width: 1px;
    border-style: solid;
    margin-right: 20px;
    height: 25px;
    width: 95px;
    margin-top: 7px;
    margin-left: 7px;
    background: #efefef;
    border-color: #b2b2b2;
    color: black;
  }
</style>
