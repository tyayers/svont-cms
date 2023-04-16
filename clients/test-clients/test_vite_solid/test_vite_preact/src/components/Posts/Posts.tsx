import { h } from 'preact';
import {useEffect, useState} from "preact/hooks";
import './Posts.css';

// Note: `user` comes from the URL, courtesy of our router
const Posts = () => {
	const [posts, setPosts] = useState({} as PostOverviewCollection);

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
		<div class="posts-container">
			<h1>Posts</h1>
			{Object.values(posts).map(post => (
				<div key={post.id}>{post.title}</div>
      ))}
		</div>
	);
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

export default Posts;
