import { render } from 'preact'
import { App } from './app'
import Posts from './components/Posts/Posts'
import './index.css'

render(<Posts />, document.getElementById('app') as HTMLElement)
