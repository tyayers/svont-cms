<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { PostOverviewCollection, Post } from "./Types";
  import Quill from 'quill';

  import Toolbar from 'quill/modules/toolbar';
  import Snow from 'quill/themes/snow';

  import Bold from 'quill/formats/bold';
  import Italic from 'quill/formats/italic';
  import Header from 'quill/formats/header';

  let post: Post
  let editor;

	let toolbarOptions = [
		[{ header: 1 }, { header: 2 }, "blockquote", "link", "image", "video"],
		["bold", "italic", "underline", "strike"],
		[{ list: "ordered" }, { list: "ordered" }],
		[{ align: [] }],
		["clean"]
	];

  onMount(async function () {
		// const { default: Quill } = await import("quill");
	
    let quill = new Quill(editor, {
      theme: "snow",
      placeholder: "Write your story..."
    });
  });
</script>

<div class="new-container">
  <div bind:this={editor} >

  </div>
</div>

<style>
  @import 'https://cdn.quilljs.com/1.3.6/quill.snow.css';

  .new-container {
    width: 500px;
    height: 700px;
    text-align: left;
  }
</style>