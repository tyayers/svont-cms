<script lang="ts">
  import CKEditor from "ckeditor5-svelte";
  import DecoupledEditor from "@ckeditor/ckeditor5-build-decoupled-document/build/ckeditor";

  // Setting up editor prop to be sent to wrapper component
  let editor = DecoupledEditor;
  // Reference to initialised editor instance
  let editorInstance = null;
  // Setting up any initial data for the editor
  let editorData = "Hello World";

  // If needed, custom editor config can be passed through to the component
  // Uncomment the custom editor config if you need to customise the editor.
  // Note: If you don't pass toolbar object then Document editor will use default set of toolbar items.
  let editorConfig = {
    toolbar: {
      items: [
        "heading",
        "|",
        "fontFamily",
        "fontSize",
        "bold",
        "italic",
        "underline"
      ]
    }
  };

  function onReady({ detail: editor }) {
    // Insert the toolbar before the editable area.
    editorInstance = editor;
    editor.ui
      .getEditableElement()
      .parentElement.insertBefore(
        editor.ui.view.toolbar.element,
        editor.ui.getEditableElement()
      );
  }
</script>

<div class="new-container">
  <CKEditor
  bind:editor
  on:ready={onReady}
  bind:config={editorConfig}
  bind:value={editorData} />
</div>

<style>
  .new-container {
    width: 500px;
    height: 700px;
    text-align: left;
  }
</style>