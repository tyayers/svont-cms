import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'

/**
 * An example element.
 *
 * @slot - This element has a slot
 * @csspart button - The button
 */
@customElement('posts-overview')
export class Posts extends LitElement {
  /**
   * Copy for the read the docs hint.
   */
  @property()
  docsHint = 'Click on the Vite and Lit logos to learn more'

  
  /**
   * The number of times the button has been clicked.
   */
  @property({ type: Number })
  count = 0

  // @property({ type: Array<PostOverview> })
  // posts: Array<PostOverview> = []
  @property({ type: Object })
  posts: PostOverviewCollection = {}

  connectedCallback() {

    super.connectedCallback();

    fetch("http://localhost:8080/posts",
		{
			method: "get",
			headers: {
				'Accept': 'application/json'
			}
		})
		.then((response) => {
			return response.json()
		})
		.then((data: PostOverviewCollection) => {
      this.posts = data
      // for (const [key, value] of Object.entries(data)) {
      //   this.posts.push(value)
      // }
      
      this.requestUpdate();
		})
  }

  render() {
    return html`
      <div>
        <h1>hello world</h1>
        ${Object.values(this.posts).map((post: PostOverview) => {
          return html`<div>${post.title}</div>`
        })}
      </div>
    `
  }

  // private _onClick() {
  //   this.count++
  // }

  static styles = css`
    :host {
      max-width: 1280px;
      margin: 0 auto;
      padding: 2rem;
      text-align: center;
    }

    .logo {
      height: 6em;
      padding: 1.5em;
      will-change: filter;
    }
    .logo:hover {
      filter: drop-shadow(0 0 2em #646cffaa);
    }
    .logo.lit:hover {
      filter: drop-shadow(0 0 2em #325cffaa);
    }

    .card {
      padding: 2em;
    }

    .read-the-docs {
      color: #888;
    }

    h1 {
      font-size: 3.2em;
      line-height: 1.1;
    }

    a {
      font-weight: 500;
      color: #646cff;
      text-decoration: inherit;
    }
    a:hover {
      color: #535bf2;
    }

    button {
      border-radius: 8px;
      border: 1px solid transparent;
      padding: 0.6em 1.2em;
      font-size: 1em;
      font-weight: 500;
      font-family: inherit;
      background-color: #1a1a1a;
      cursor: pointer;
      transition: border-color 0.25s;
    }
    button:hover {
      border-color: #646cff;
    }
    button:focus,
    button:focus-visible {
      outline: 4px auto -webkit-focus-ring-color;
    }

    @media (prefers-color-scheme: light) {
      a:hover {
        color: #747bff;
      }
      button {
        background-color: #f9f9f9;
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'posts-overview': Posts
  }
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
};