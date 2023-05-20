<script lang="ts">
  import { goto } from "$app/navigation";
  import { onDestroy } from "svelte";
  import Editor, { getData, setData } from "./Editor.CKBBlock.svelte";
  import TagBox from "./Tag.box.svelte";
  import { appService } from "./DataService";

  import type { Post, AppUser, SearchResult } from "./DataInterface";
  import { append } from "svelte/internal";

  export let post: Post = undefined;
  export let statusUpdate: (status: string) => void;

  let localUser: AppUser = undefined;

  appService.user.subscribe((value) => {
    localUser = value;
  });

  let tags: string[] = [];
  let isNewPost: boolean = true;

  if (post) {
    isNewPost = false;
    setData(post.content);
    if (post.header.tags && post.header.tags.length > 0)
      tags = post.header.tags;
  } else {
    setData("");
    const formData = new FormData();
    formData.set("authorId", localUser.uid);
    formData.set("authorDisplayName", localUser.displayName);
    formData.set("authorProfilePic", localUser.photoURL);
    // Mark as draft
    formData.set("draft", "true");

    appService.CreatePost(formData).then((newPost: Post) => {
      post = newPost;
    });
  }

  onDestroy(() => {
    if (isNewPost && post.header.title == "") {
      // This post was never really started, so delete before leaving...
      appService.DeletePost(post.header.id);
    }
  });

  export function submit(draft: boolean = false) {
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
    let tempTags: string = tags.toString();
    if (tempTags) formData.set("tags", tempTags);

    if (post) {
      if (draft) formData.set("draft", "true");
      else formData.set("draft", "false");

      //let newDraft: boolean = draft;
      appService.UpdatePost(post.header.id, formData).then((post: Post) => {
        if (!draft) {
          //goto("/posts/" + post.header.id);
          goto("/home");
        } else if (statusUpdate) statusUpdate("Draft saved");
      });
    } else {
      // Set user for new post
      formData.set("authorId", localUser.uid);
      formData.set("authorDisplayName", localUser.displayName);
      formData.set("authorProfilePic", localUser.photoURL);

      appService.CreatePost(formData).then((post: Post) => {
        //goto("/posts/" + post.header.id);
        goto("/home");
      });
    }
  }

  function init(el) {
    el.focus();
  }

  function searchTags(searchInput: string): Promise<SearchResult[]> {
    return appService.SearchTags(searchInput);
  }

  function addTag(event) {
    tags.push(event.detail.name);
  }

  function saveDraft() {
    if (statusUpdate) statusUpdate("Saving draft...");
    submit(true);
  }
</script>

<div class="new-container">
  <!-- <h1>New Post</h1> -->
  <form id="new_post_form">
    <!-- <label for="title">Enter the post title: </label> -->
    {#if post}
      <input
        type="text"
        name="title"
        id="title"
        placeholder="Title"
        required
        bind:value={post.header.title}
        autofocus
      />
    {:else}
      <input
        type="text"
        name="title"
        id="title"
        placeholder="Title"
        required
        autofocus
      />
    {/if}

    <div>
      <br />
      {#if post}
        <Editor
          imageUploadPath={"/posts/" + post.header.id + "/files"}
          {saveDraft}
        />
      {/if}
    </div>
    <br />
    <div class="tag_frame">
      <TagBox {searchTags} bind:tags on:addTag={addTag} />
    </div>
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
    margin-bottom: 148px;
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
    color: gray;
    font-size: 15px;
    font-weight: 500;
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

  .tag_frame {
    margin-top: 30px;
    margin-bottom: 20px;
  }
</style>
