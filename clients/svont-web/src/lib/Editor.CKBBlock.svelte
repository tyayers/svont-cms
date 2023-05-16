<script lang="ts" context="module">
  let editor;
  let initialData;
  let loaded: boolean = false;

  export function getData(): string {
    return editor.getData();
  }

  export function setData(data: string) {
    initialData = data;
    if (editor) {
      console.log("setting loaded to true in setData");
      setLoadedStatus(true);
      editor.setData(data);
    }
  }

  function setLoadedStatus(newStatus: boolean) {
    loaded = newStatus;
    if (loaded) {
      //document.getElementById("loading_frame").style.display = "none";
      //editor.focus();
    }
  }
</script>

<script lang="ts">
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";

  import { appService } from "./DataService";

  export let imageUploadPath: string = "";

  export let saveDraft: () => void;

  function initCkeditor() {
    appService.GetIdToken().then((token) => {
      BalloonBlockEditor.create(document.querySelector("#content"), {
        placeholder: "Add your content here...",
        simpleUpload: {
          uploadUrl: appService.GetServer("IMAGE") + imageUploadPath,
          withCredentials: true,
          headers: {
            Authorization: "Bearer " + token,
          },
        },
        autosave: {
          save(editor) {
            //return saveData( editor.getData() );
            if (saveDraft) saveDraft();
            return true;
          },
        },
      }).then((newEditor) => {
        editor = newEditor;
        setLoadedStatus(true);
        if (initialData) {
          console.log("Setting loaded to true in editor create");
          setLoadedStatus(true);
          editor.setData(initialData);
        }
      });
    });
  }

  onMount(async function () {
    //console.log("ready")
    initCkeditor();
  });
</script>

<!-- <svelte:head>
	<script src="https://cdn.ckeditor.com/ckeditor5/35.4.0/balloon-block/ckeditor.js" on:load={initCkeditor}></script>
</svelte:head> -->

<!-- {#if loaded} -->
<div id="content" class="post_content" />

<!-- <div id="loading_frame"></div> -->
<style>
  .post_content {
    /* position:fixed;
    top: 1000px; */
    /* border-style: solid;
    border-width: 1px;
    border-color: lightgray;
    height: 132px; */
  }

  #loading_frame {
    position: fixed;
    top: 80px;
    bottom: 0px;
    width: 100vw;
    background-color: white;
  }

  .hidden {
    visibility: hidden;
    opacity: 0;
    transition: visibility 0s 2s, opacity 2s linear;
  }
</style>
