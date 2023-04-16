import { useState } from 'preact/hooks'
import { Router, Route } from 'preact-router';
import { Posts } from './components/Posts.view'
import { PostsNew } from './components/Post.new';
import preactLogo from './assets/preact.svg'
import './app.css'
import { PostView } from './components/Post.view';

export function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <h1>Vite + Preact</h1>
      <div class="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <a href="/">Posts</a>
        <a href="/new">New Post</a>
        <Router>
          <Route path="/" component={Posts}></Route>
          <Route path="/new" component={PostsNew}></Route>
          <Route path="/posts/:id" component={PostView}></Route>
        </Router>
      </div>
    </>
  )
}
