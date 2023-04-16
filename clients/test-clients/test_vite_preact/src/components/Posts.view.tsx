import { useState, useEffect} from 'preact/hooks'
import preactLogo from '../assets/preact.svg'
import './Posts.css'

export function Posts() {
  const [count, setCount] = useState(0)
  const [posts, setPosts] = useState({} as PostOverviewCollection)

	useEffect(() => {
    fetch("http://localhost:8080/posts",
		{
			method: "get",
			headers: {
				"Accept": "application/json"
			}
		})
		.then((response) => {
			return response.json()
		})
		.then((data: PostOverviewCollection) => {
      setPosts(data)
		})
	}, []);

  return (
    <>
      <div>
        <h1>Posts</h1>
        {Object.values(posts).map(post => (
          <div key={post.id}>
            <a href={"/posts/" + post.id}>{post.title}</a>
          </div>
        ))}
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
