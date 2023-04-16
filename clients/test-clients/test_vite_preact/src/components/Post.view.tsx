import { useState, useEffect} from 'preact/hooks'
import { route } from 'preact-router';
import preactLogo from '../assets/preact.svg'
import './Posts.css'

export function PostView(props: any) {
  const [count, setCount] = useState(0)
  const [post, setPost] = useState({} as Post)
  
	useEffect(() => {
    fetch("http://localhost:8080/posts/" + props["id"],
		{
			method: "get",
			headers: {
				"Accept": "application/json"
			}
		})
		.then((response) => {
			return response.json()
		})
		.then((data: Post) => {
      setPost(data)
		})
	}, []);

  return (
    <>
      <div class="new-container">
        <h1>{post.title}</h1>
        <div>{post.content}</div>
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

type Post = {
  id: string;
  title: string;
  author: string;
  content: string
  upvotes: number;
  created: string;
  updated: string;
  files: string[];
}
