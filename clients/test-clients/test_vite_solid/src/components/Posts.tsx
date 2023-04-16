import type { Component } from 'solid-js';
import {createSignal, onMount, For} from 'solid-js'

import styles from './Posts.module.css';

export const Posts: Component = () => {
  const [posts, setPosts] = createSignal({})

  onMount(async () => {
    const res = await fetch("http://localhost:8080/posts");
    setPosts(await res.json())
  })

  return (
    <div class={styles.App}>
      <h2>Posts</h2>
      <For each={Object.values(posts())}>{(post: PostOverview, i) =>
          <div>{post.title}</div>
      }</For>
      
    </div>
  );
};

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
