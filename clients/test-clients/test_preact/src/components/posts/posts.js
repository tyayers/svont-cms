import { h } from 'preact';
import {useEffect, useState} from "preact/hooks";
import style from './style.css';

// Note: `user` comes from the URL, courtesy of our router
const Posts = () => {
	const [posts, setPosts] = useState({});

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
		.then((data) => {
      setPosts(data)
		})
	}, []);

	return (
		<div>
			<h1>Posts</h1>
			{Object.values(posts).map(post => (
				<div key={post.id}>{post.title}</div>
      ))}
		</div>
	);
}

export default Posts;
