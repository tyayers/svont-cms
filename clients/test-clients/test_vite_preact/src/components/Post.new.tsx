import { useState, useEffect} from 'preact/hooks'
import { route } from 'preact-router';
import preactLogo from '../assets/preact.svg'
import './Posts.css'

export function PostsNew() {
  const [count, setCount] = useState(0)
  const [posts, setPosts] = useState({} as PostOverviewCollection)

  function onSubmit(e: any) {
    e.preventDefault();
    const formData = new FormData(e.target);
    fetch("http://localhost:8080/posts",
      {
        body: formData,
        method: "post",
        headers: {
          'Accept': 'application/json'
        }
      })
      .then((response) => {
        return response.json()
       })
      .then((data) => {
        route('/', true);
      })
  }	

  return (
    <>
      <div class="new-container">
        <h1>New Post</h1>
        <form onSubmit={onSubmit}>
          <label for="title">Enter the post title: </label>
          <input type="text" name="title" id="title" required />

          <div>
            <label for="content">Enter the post content: </label><br/>
            <textarea name="content" id="content" class="post_content"></textarea>
          </div>

          <div>
            <label for="files">Select a file:</label>
            <input type="file" id="files" name="files" multiple />
          </div>

          <button type="submit">Submit</button>
        </form>
      </div>
    </>
  )
}

type PostOverviewCollection = { 
  [id: string]: PostOverview; 
}

type PostOverview = {
  id: string;
  title: string;
  author: string;
  upvotes: number;
  created: string;
  updated: string;
  fileCount: number;
}
